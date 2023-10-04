package query

import (
	"context"
	"errors"
)

type HasOneHydrator[Root any, Sub any, SubID comparable] interface {
	RootToSubID(root Root) (SubID, bool)
	GetSubs(ctx context.Context, dbCtx DbContext, ids []SubID) ([]Sub, error)
	SubID(sub Sub) SubID
	Hydrate(root *Root, sub Sub)
	MustSucceed() bool
}

func HydrateHasOne[Root any, Sub any, SubID comparable](h HasOneHydrator[Root, Sub, SubID]) Option {
	return func(b OptionBuilder) (OptionBuilder, error) {
		return b.Callback(func(ctx context.Context, cbCtx CallbackContext, result any) error {
			items, ok := result.([]Root)
			if !ok {
				return errors.New("invalid result type")
			}
			return hydrateHasOne(ctx, cbCtx, h, items)
		}), nil
	}
}

func hydrateHasOne[Root any, Sub any, SubID comparable](
	ctx context.Context,
	cbCtx CallbackContext,
	h HasOneHydrator[Root, Sub, SubID],
	roots []Root,
) error {
	indexMap := make(map[int]SubID, len(roots))
	idMap := make(map[SubID]struct{}, len(roots))
	ids := make([]SubID, 0, len(roots))
	cbCtx.Lock()
	for i, root := range roots {
		if subID, ok1 := h.RootToSubID(root); ok1 {
			indexMap[i] = subID
			if _, ok2 := idMap[subID]; !ok2 {
				ids = append(ids, subID)
				idMap[subID] = struct{}{}
			}
		}
	}
	cbCtx.Unlock()
	if len(ids) == 0 {
		return nil
	}
	subs, err := h.GetSubs(ctx, cbCtx, ids)
	if err != nil {
		return err
	}
	subsMap := make(map[SubID]Sub, len(subs))
	for _, sub := range subs {
		subsMap[h.SubID(sub)] = sub
	}
	cbCtx.Lock()
	defer cbCtx.Unlock()
	for i, subID := range indexMap {
		if sub, ok := subsMap[subID]; ok {
			h.Hydrate(&roots[i], sub)
		} else if h.MustSucceed() {
			return errors.New("failed to hydrateHasOne")
		}
	}
	return nil
}

type HasManyHydrator[Root any, RootID comparable, JoinSub any, Sub any] interface {
	RootID(root Root) (RootID, bool)
	GetJoinSubs(ctx context.Context, dbCtx DbContext, ids []RootID) ([]JoinSub, error)
	JoinSubToRootIDAndSub(j JoinSub) (RootID, Sub)
	Hydrate(root *Root, subs []Sub)
}

func HydrateHasMany[Root any, RootID comparable, JoinSub any, Sub any](h HasManyHydrator[Root, RootID, JoinSub, Sub]) Option {
	return func(b OptionBuilder) (OptionBuilder, error) {
		return b.Callback(func(ctx context.Context, cbCtx CallbackContext, result any) error {
			items, ok := result.([]Root)
			if !ok {
				return errors.New("invalid result type")
			}
			return hydrateHasMany(ctx, cbCtx, h, items)
		}), nil
	}
}

func hydrateHasMany[Root any, RootID comparable, JoinSub any, Sub any](
	ctx context.Context,
	cbCtx CallbackContext,
	h HasManyHydrator[Root, RootID, JoinSub, Sub],
	roots []Root,
) error {
	indexMap := make(map[int]RootID, len(roots))
	idMap := make(map[RootID]struct{}, len(roots))
	ids := make([]RootID, 0, len(roots))
	cbCtx.Lock()
	for i, root := range roots {
		if rootID, ok1 := h.RootID(root); ok1 {
			indexMap[i] = rootID
			if _, ok2 := idMap[rootID]; !ok2 {
				ids = append(ids, rootID)
				idMap[rootID] = struct{}{}
			}
		}
	}
	cbCtx.Unlock()
	if len(ids) == 0 {
		return nil
	}
	joinSubs, err := h.GetJoinSubs(ctx, cbCtx, ids)
	if err != nil {
		return err
	}
	subsMap := make(map[RootID][]Sub, len(roots))
	for _, j := range joinSubs {
		rootID, sub := h.JoinSubToRootIDAndSub(j)
		subsMap[rootID] = append(subsMap[rootID], sub)
	}
	cbCtx.Lock()
	defer cbCtx.Unlock()
	for i, rootID := range indexMap {
		if thisSubs, ok := subsMap[rootID]; ok {
			h.Hydrate(&roots[i], thisSubs)
		}
	}
	return nil
}
