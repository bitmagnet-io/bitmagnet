package processor

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/blocker"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"gorm.io/gen/field"
)

type Processor interface {
	NewJob(params MessageParams) runner.Runner
}

type processor struct {
	searchClient     search.Search
	blocker          blocker.Blocker
	runner           classifier.Runner
	queueJobProvider queue.JobProvider[MessageParams]
	persister        persister.Adder
	defaultWorkflow  classifier.Workflow
}

const concurrency = 10

func (p processor) NewJob(params MessageParams) runner.Runner {
	return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
		shutdown := make(chan struct{})

		setupCtx, setupCancel := context.WithCancel(ctx)

		go func() {
			// abort setup if shutdown initiated:
			select {
			case <-shutdown:
				cancel(fmt.Errorf("%w: %w: %w", Err, ErrSetup, ErrInterrupted))
			case <-setupCtx.Done():
			}

			setupCancel()
		}()

		searchResult, err := p.getTorrentsToClassify(setupCtx, params)

		setupCancel()

		if err != nil {
			err = fmt.Errorf("%w: %w: %w", Err, ErrSetup, err)
			cancel(err)

			return runner.NopShutdowner, err
		}

		if len(searchResult.Torrents)+len(searchResult.MissingInfoHashes) == 0 {
			cancel(runner.ErrCompleted)

			return runner.NopShutdowner, nil
		}

		go func() {
			defer cancel(runner.ErrCompleted)

			err := p.handleTorrents(
				ctx,
				shutdown,
				params.ClassifierParams,
				searchResult,
			)
			if err != nil {
				// todo: Check error handling
				err = fmt.Errorf("%w: %w: %w", Err, ErrClassify, err)
				cancel(err)

				return
			}
		}()

		return func(shutdownCtx context.Context) error {
			close(shutdown)

			select {
			case <-ctx.Done():
			case <-shutdownCtx.Done():
			}

			cause := context.Cause(ctx)
			if cause != nil && !errors.Is(cause, runner.ErrCompleted) {
				return fmt.Errorf("%w: %w: %w", Err, ErrShutdown, cause)
			}

			return nil
		}, nil
	}
}

func (p processor) getTorrentsToClassify(
	ctx context.Context,
	params MessageParams,
) (search.TorrentsWithMissingInfoHashesResult, error) {
	searchResult, searchErr := p.searchClient.TorrentsWithMissingInfoHashes(
		ctx,
		params.InfoHashes,
		query.Preload(func(q *dao.Query) []field.RelationField {
			return []field.RelationField{
				q.Torrent.Files.RelationField,
				q.Torrent.Hint.RelationField,
				q.Torrent.Sources.RelationField,
				q.Torrent.Tags.RelationField,
			}
		}),
	)
	if searchErr != nil {
		return searchResult, searchErr
	}

	tcResult, tcErr := p.searchClient.TorrentContent(
		ctx,
		query.Where(search.TorrentContentInfoHashCriteria(params.InfoHashes...)),
		search.HydrateTorrentContentContent(),
	)
	if tcErr != nil {
		return searchResult, tcErr
	}

	for _, tc := range tcResult.Items {
		for ti, t := range searchResult.Torrents {
			if t.InfoHash == tc.InfoHash {
				searchResult.Torrents[ti].Contents = append(
					searchResult.Torrents[ti].Contents,
					tc.TorrentContent,
				)

				break
			}
		}
	}

	return searchResult, nil
}

func (p processor) handleTorrents(
	ctx context.Context,
	shutdown <-chan struct{},
	params ClassifierParams,
	searchResult search.TorrentsWithMissingInfoHashesResult,
) error {
	var (
		mtx         sync.Mutex
		interrupted bool
		hasInputs   bool
		errs        []error
	)

	addErr := func(err error) {
		if err == nil {
			return
		}

		mtx.Lock()
		defer mtx.Unlock()

		errs = append(errs, err)
	}

	failedHashes := make([]protocol.ID, 0, len(searchResult.MissingInfoHashes))
	failedHashes = append(failedHashes, searchResult.MissingInfoHashes...)

	sem := make(chan struct{}, concurrency)

	classifyCtx, classifyCancel := context.WithCancel(ctx)
	defer classifyCancel()

	for _, torrent := range searchResult.Torrents {
		select {
		case <-classifyCtx.Done():
			return classifyCtx.Err()
		case <-shutdown:
			mtx.Lock()
			failedHashes = append(failedHashes, torrent.InfoHash)
			interrupted = true
			mtx.Unlock()

			continue
		case sem <- struct{}{}:
		}

		go func(torrent model.Torrent) {
			defer func() {
				<-sem
			}()

			persisterInputs, err := p.handleTorrent(classifyCtx, params, torrent)
			if err != nil {
				failedHashes = append(failedHashes, torrent.InfoHash)

				addErr(err)
			}

			if len(persisterInputs) > 0 {
				err := p.persister.Add(ctx, persisterInputs.Input())
				if err != nil {
					addErr(err)
				} else {
					hasInputs = true
				}
			}
		}(torrent)
	}

	// wait for all classifications:
	for range concurrency {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case sem <- struct{}{}:
		}
	}

	if !hasInputs {
		var err error

		if interrupted {
			err = ErrInterrupted
		} else {
			err = errors.Join(errs...)
			if err != nil {
				err = fmt.Errorf("%w: %w: %w", Err, ErrAllTorrentsFailed, err)
			}
		}

		return err
	}

	if len(failedHashes) > 0 {
		var inputs persister.Inputs

		for _, infoHash := range failedHashes {
			job, err := p.queueJobProvider(MessageParams{
				InfoHashes:       []protocol.ID{infoHash},
				ClassifierParams: params,
			})
			if err != nil {
				return fmt.Errorf("%w: %w", Err, err)
			}

			inputs = append(inputs, persister.InputQueueJobs(job))
		}

		err := p.persister.Add(ctx, inputs.Input())
		if err != nil {
			return fmt.Errorf("%w: %w", Err, err)
		}
	}

	return nil
}

type handleTorrentResult struct {
	torrentContent model.TorrentContent
	tcRefsToDelete map[model.TorrentContentRef]struct{}
	tags           map[string]bool
}

func (p processor) handleTorrent(
	ctx context.Context,
	params ClassifierParams,
	torrent model.Torrent,
) (persister.Inputs, error) {
	thisDeleteIDs := make(map[model.TorrentContentRef]struct{}, len(torrent.Contents))
	foundMatch := false

	for _, tc := range torrent.Contents {
		thisDeleteIDs[tc.Ref()] = struct{}{}

		if !foundMatch &&
			!torrent.Hint.ContentSource.Valid &&
			params.ClassifyMode != ClassifyModeRematch &&
			tc.ContentType.Valid &&
			tc.ContentSource.Valid &&
			(torrent.Hint.IsNil() || torrent.Hint.ContentType == tc.ContentType.ContentType) {
			torrent.Hint.ContentType = tc.ContentType.ContentType
			torrent.Hint.ContentSource = tc.ContentSource
			torrent.Hint.ContentID = tc.ContentID
			foundMatch = true
		}
	}

	workflowName := params.ClassifierWorkflow
	if workflowName == "" {
		workflowName = p.defaultWorkflow
	}

	cl, err := p.runner.Run(ctx, workflowName, params.ClassifierFlags, torrent)

	if err != nil {
		return errorToInputs(torrent.InfoHash, err)
	}

	torrentContent := cl.ToTorrentContent()

	delete(thisDeleteIDs, torrentContent.Ref())

	return resultToInputs(handleTorrentResult{
		torrentContent: torrentContent,
		tcRefsToDelete: thisDeleteIDs,
		tags:           cl.Tags,
	}), nil
}

func errorToInputs(infoHash protocol.ID, err error) (persister.Inputs, error) {
	if errors.Is(err, classification.ErrDeleteTorrent) {
		return persister.Inputs{persister.InputDeleteInfoHashes(infoHash)}, nil
	}

	return nil, err
}

func resultToInputs(result handleTorrentResult) persister.Inputs {
	var inputs persister.Inputs

	for ref := range result.tcRefsToDelete {
		inputs = append(inputs, persister.InputDeleteTorrentContent(ref))
	}

	inputs = append(inputs, persister.InputTorrentContent(result.torrentContent))

	var tagsToAdd, tagsToRemove []string
	for tagName, addRemove := range result.tags {
		if addRemove {
			tagsToAdd = append(tagsToAdd, tagName)
		} else {
			tagsToRemove = append(tagsToRemove, tagName)
		}
	}

	if len(tagsToAdd) > 0 {
		inputs = append(
			inputs,
			persister.InputTorrentTags(result.torrentContent.InfoHash, tagsToAdd...),
		)
	}

	if len(tagsToRemove) > 0 {
		inputs = append(
			inputs,
			persister.InputDeleteTorrentTags(result.torrentContent.InfoHash, tagsToRemove...),
		)
	}

	return inputs
}
