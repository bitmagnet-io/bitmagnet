// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gen

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics/queuemetrics"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics/torrentmetrics"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

type ContentTypeAgg struct {
	Value      *model.ContentType `json:"value,omitempty"`
	Label      string             `json:"label"`
	Count      int                `json:"count"`
	IsEstimate bool               `json:"isEstimate"`
}

type ContentTypeFacetInput struct {
	Aggregate graphql.Omittable[*bool]                `json:"aggregate,omitempty"`
	Filter    graphql.Omittable[[]*model.ContentType] `json:"filter,omitempty"`
}

type GenreAgg struct {
	Value      string `json:"value"`
	Label      string `json:"label"`
	Count      int    `json:"count"`
	IsEstimate bool   `json:"isEstimate"`
}

type GenreFacetInput struct {
	Aggregate graphql.Omittable[*bool]             `json:"aggregate,omitempty"`
	Logic     graphql.Omittable[*model.FacetLogic] `json:"logic,omitempty"`
	Filter    graphql.Omittable[[]string]          `json:"filter,omitempty"`
}

type HealthCheck struct {
	Key       string       `json:"key"`
	Status    HealthStatus `json:"status"`
	Timestamp time.Time    `json:"timestamp"`
	Error     *string      `json:"error,omitempty"`
}

type HealthQuery struct {
	Status HealthStatus  `json:"status"`
	Checks []HealthCheck `json:"checks"`
}

type LanguageAgg struct {
	Value      model.Language `json:"value"`
	Label      string         `json:"label"`
	Count      int            `json:"count"`
	IsEstimate bool           `json:"isEstimate"`
}

type LanguageFacetInput struct {
	Aggregate graphql.Omittable[*bool]            `json:"aggregate,omitempty"`
	Filter    graphql.Omittable[[]model.Language] `json:"filter,omitempty"`
}

type Mutation struct {
}

type Query struct {
}

type QueueJobQueueAgg struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

type QueueJobQueueFacetInput struct {
	Aggregate graphql.Omittable[*bool]    `json:"aggregate,omitempty"`
	Filter    graphql.Omittable[[]string] `json:"filter,omitempty"`
}

type QueueJobStatusAgg struct {
	Value model.QueueJobStatus `json:"value"`
	Label string               `json:"label"`
	Count int                  `json:"count"`
}

type QueueJobStatusFacetInput struct {
	Aggregate graphql.Omittable[*bool]                  `json:"aggregate,omitempty"`
	Filter    graphql.Omittable[[]model.QueueJobStatus] `json:"filter,omitempty"`
}

type QueueJobsAggregations struct {
	Queue  []QueueJobQueueAgg  `json:"queue,omitempty"`
	Status []QueueJobStatusAgg `json:"status,omitempty"`
}

type QueueJobsFacetsInput struct {
	Status graphql.Omittable[*QueueJobStatusFacetInput] `json:"status,omitempty"`
	Queue  graphql.Omittable[*QueueJobQueueFacetInput]  `json:"queue,omitempty"`
}

type QueueJobsOrderByInput struct {
	Field      QueueJobsOrderByField    `json:"field"`
	Descending graphql.Omittable[*bool] `json:"descending,omitempty"`
}

type QueueMetricsQueryInput struct {
	BucketDuration MetricsBucketDuration                     `json:"bucketDuration"`
	Statuses       graphql.Omittable[[]model.QueueJobStatus] `json:"statuses,omitempty"`
	Queues         graphql.Omittable[[]string]               `json:"queues,omitempty"`
	StartTime      graphql.Omittable[*time.Time]             `json:"startTime,omitempty"`
	EndTime        graphql.Omittable[*time.Time]             `json:"endTime,omitempty"`
}

type QueueMetricsQueryResult struct {
	Buckets []queuemetrics.Bucket `json:"buckets"`
}

type ReleaseYearAgg struct {
	Value      *model.Year `json:"value,omitempty"`
	Label      string      `json:"label"`
	Count      int         `json:"count"`
	IsEstimate bool        `json:"isEstimate"`
}

type ReleaseYearFacetInput struct {
	Aggregate graphql.Omittable[*bool]         `json:"aggregate,omitempty"`
	Filter    graphql.Omittable[[]*model.Year] `json:"filter,omitempty"`
}

type SuggestTagsQueryInput struct {
	Prefix     graphql.Omittable[*string]  `json:"prefix,omitempty"`
	Exclusions graphql.Omittable[[]string] `json:"exclusions,omitempty"`
}

type TorrentContentAggregations struct {
	ContentType     []ContentTypeAgg     `json:"contentType,omitempty"`
	TorrentSource   []TorrentSourceAgg   `json:"torrentSource,omitempty"`
	TorrentTag      []TorrentTagAgg      `json:"torrentTag,omitempty"`
	TorrentFileType []TorrentFileTypeAgg `json:"torrentFileType,omitempty"`
	Language        []LanguageAgg        `json:"language,omitempty"`
	Genre           []GenreAgg           `json:"genre,omitempty"`
	ReleaseYear     []ReleaseYearAgg     `json:"releaseYear,omitempty"`
	VideoResolution []VideoResolutionAgg `json:"videoResolution,omitempty"`
	VideoSource     []VideoSourceAgg     `json:"videoSource,omitempty"`
}

type TorrentContentFacetsInput struct {
	ContentType     graphql.Omittable[*ContentTypeFacetInput]     `json:"contentType,omitempty"`
	TorrentSource   graphql.Omittable[*TorrentSourceFacetInput]   `json:"torrentSource,omitempty"`
	TorrentTag      graphql.Omittable[*TorrentTagFacetInput]      `json:"torrentTag,omitempty"`
	TorrentFileType graphql.Omittable[*TorrentFileTypeFacetInput] `json:"torrentFileType,omitempty"`
	Language        graphql.Omittable[*LanguageFacetInput]        `json:"language,omitempty"`
	Genre           graphql.Omittable[*GenreFacetInput]           `json:"genre,omitempty"`
	ReleaseYear     graphql.Omittable[*ReleaseYearFacetInput]     `json:"releaseYear,omitempty"`
	VideoResolution graphql.Omittable[*VideoResolutionFacetInput] `json:"videoResolution,omitempty"`
	VideoSource     graphql.Omittable[*VideoSourceFacetInput]     `json:"videoSource,omitempty"`
	PublishedAt     graphql.Omittable[*string]                    `json:"publishedAt,omitempty"`
}

type TorrentContentOrderByInput struct {
	Field      TorrentContentOrderByField `json:"field"`
	Descending graphql.Omittable[*bool]   `json:"descending,omitempty"`
}

type TorrentFileTypeAgg struct {
	Value      model.FileType `json:"value"`
	Label      string         `json:"label"`
	Count      int            `json:"count"`
	IsEstimate bool           `json:"isEstimate"`
}

type TorrentFileTypeFacetInput struct {
	Aggregate graphql.Omittable[*bool]             `json:"aggregate,omitempty"`
	Logic     graphql.Omittable[*model.FacetLogic] `json:"logic,omitempty"`
	Filter    graphql.Omittable[[]model.FileType]  `json:"filter,omitempty"`
}

type TorrentFilesOrderByInput struct {
	Field      TorrentFilesOrderByField `json:"field"`
	Descending graphql.Omittable[*bool] `json:"descending,omitempty"`
}

type TorrentListSourcesResult struct {
	Sources []model.TorrentSource `json:"sources"`
}

type TorrentMetricsQueryInput struct {
	BucketDuration MetricsBucketDuration         `json:"bucketDuration"`
	Sources        graphql.Omittable[[]string]   `json:"sources,omitempty"`
	StartTime      graphql.Omittable[*time.Time] `json:"startTime,omitempty"`
	EndTime        graphql.Omittable[*time.Time] `json:"endTime,omitempty"`
}

type TorrentMetricsQueryResult struct {
	Buckets []torrentmetrics.Bucket `json:"buckets"`
}

type TorrentReprocessInput struct {
	InfoHashes          []protocol.ID              `json:"infoHashes"`
	ClassifierRematch   graphql.Omittable[*bool]   `json:"classifierRematch,omitempty"`
	ClassifierWorkflow  graphql.Omittable[*string] `json:"classifierWorkflow,omitempty"`
	ApisDisabled        graphql.Omittable[*bool]   `json:"apisDisabled,omitempty"`
	LocalSearchDisabled graphql.Omittable[*bool]   `json:"localSearchDisabled,omitempty"`
}

type TorrentSourceAgg struct {
	Value      string `json:"value"`
	Label      string `json:"label"`
	Count      int    `json:"count"`
	IsEstimate bool   `json:"isEstimate"`
}

type TorrentSourceFacetInput struct {
	Aggregate graphql.Omittable[*bool]             `json:"aggregate,omitempty"`
	Logic     graphql.Omittable[*model.FacetLogic] `json:"logic,omitempty"`
	Filter    graphql.Omittable[[]string]          `json:"filter,omitempty"`
}

type TorrentTagAgg struct {
	Value      string `json:"value"`
	Label      string `json:"label"`
	Count      int    `json:"count"`
	IsEstimate bool   `json:"isEstimate"`
}

type TorrentTagFacetInput struct {
	Aggregate graphql.Omittable[*bool]             `json:"aggregate,omitempty"`
	Logic     graphql.Omittable[*model.FacetLogic] `json:"logic,omitempty"`
	Filter    graphql.Omittable[[]string]          `json:"filter,omitempty"`
}

type VideoResolutionAgg struct {
	Value      *model.VideoResolution `json:"value,omitempty"`
	Label      string                 `json:"label"`
	Count      int                    `json:"count"`
	IsEstimate bool                   `json:"isEstimate"`
}

type VideoResolutionFacetInput struct {
	Aggregate graphql.Omittable[*bool]                    `json:"aggregate,omitempty"`
	Filter    graphql.Omittable[[]*model.VideoResolution] `json:"filter,omitempty"`
}

type VideoSourceAgg struct {
	Value      *model.VideoSource `json:"value,omitempty"`
	Label      string             `json:"label"`
	Count      int                `json:"count"`
	IsEstimate bool               `json:"isEstimate"`
}

type VideoSourceFacetInput struct {
	Aggregate graphql.Omittable[*bool]                `json:"aggregate,omitempty"`
	Filter    graphql.Omittable[[]*model.VideoSource] `json:"filter,omitempty"`
}

type Worker struct {
	Key     string `json:"key"`
	Started bool   `json:"started"`
}

type WorkersListAllQueryResult struct {
	Workers []Worker `json:"workers"`
}

type WorkersQuery struct {
	ListAll WorkersListAllQueryResult `json:"listAll"`
}

type HealthStatus string

const (
	HealthStatusUnknown  HealthStatus = "unknown"
	HealthStatusInactive HealthStatus = "inactive"
	HealthStatusUp       HealthStatus = "up"
	HealthStatusDown     HealthStatus = "down"
)

var AllHealthStatus = []HealthStatus{
	HealthStatusUnknown,
	HealthStatusInactive,
	HealthStatusUp,
	HealthStatusDown,
}

func (e HealthStatus) IsValid() bool {
	switch e {
	case HealthStatusUnknown, HealthStatusInactive, HealthStatusUp, HealthStatusDown:
		return true
	}
	return false
}

func (e HealthStatus) String() string {
	return string(e)
}

func (e *HealthStatus) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = HealthStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid HealthStatus", str)
	}
	return nil
}

func (e HealthStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type MetricsBucketDuration string

const (
	MetricsBucketDurationMinute MetricsBucketDuration = "minute"
	MetricsBucketDurationHour   MetricsBucketDuration = "hour"
	MetricsBucketDurationDay    MetricsBucketDuration = "day"
)

var AllMetricsBucketDuration = []MetricsBucketDuration{
	MetricsBucketDurationMinute,
	MetricsBucketDurationHour,
	MetricsBucketDurationDay,
}

func (e MetricsBucketDuration) IsValid() bool {
	switch e {
	case MetricsBucketDurationMinute, MetricsBucketDurationHour, MetricsBucketDurationDay:
		return true
	}
	return false
}

func (e MetricsBucketDuration) String() string {
	return string(e)
}

func (e *MetricsBucketDuration) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MetricsBucketDuration(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MetricsBucketDuration", str)
	}
	return nil
}

func (e MetricsBucketDuration) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type QueueJobsOrderByField string

const (
	QueueJobsOrderByFieldCreatedAt QueueJobsOrderByField = "created_at"
	QueueJobsOrderByFieldRanAt     QueueJobsOrderByField = "ran_at"
	QueueJobsOrderByFieldPriority  QueueJobsOrderByField = "priority"
)

var AllQueueJobsOrderByField = []QueueJobsOrderByField{
	QueueJobsOrderByFieldCreatedAt,
	QueueJobsOrderByFieldRanAt,
	QueueJobsOrderByFieldPriority,
}

func (e QueueJobsOrderByField) IsValid() bool {
	switch e {
	case QueueJobsOrderByFieldCreatedAt, QueueJobsOrderByFieldRanAt, QueueJobsOrderByFieldPriority:
		return true
	}
	return false
}

func (e QueueJobsOrderByField) String() string {
	return string(e)
}

func (e *QueueJobsOrderByField) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = QueueJobsOrderByField(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid QueueJobsOrderByField", str)
	}
	return nil
}

func (e QueueJobsOrderByField) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TorrentContentOrderByField string

const (
	TorrentContentOrderByFieldRelevance   TorrentContentOrderByField = "relevance"
	TorrentContentOrderByFieldPublishedAt TorrentContentOrderByField = "published_at"
	TorrentContentOrderByFieldUpdatedAt   TorrentContentOrderByField = "updated_at"
	TorrentContentOrderByFieldSize        TorrentContentOrderByField = "size"
	TorrentContentOrderByFieldFilesCount  TorrentContentOrderByField = "files_count"
	TorrentContentOrderByFieldSeeders     TorrentContentOrderByField = "seeders"
	TorrentContentOrderByFieldLeechers    TorrentContentOrderByField = "leechers"
	TorrentContentOrderByFieldName        TorrentContentOrderByField = "name"
	TorrentContentOrderByFieldInfoHash    TorrentContentOrderByField = "info_hash"
)

var AllTorrentContentOrderByField = []TorrentContentOrderByField{
	TorrentContentOrderByFieldRelevance,
	TorrentContentOrderByFieldPublishedAt,
	TorrentContentOrderByFieldUpdatedAt,
	TorrentContentOrderByFieldSize,
	TorrentContentOrderByFieldFilesCount,
	TorrentContentOrderByFieldSeeders,
	TorrentContentOrderByFieldLeechers,
	TorrentContentOrderByFieldName,
	TorrentContentOrderByFieldInfoHash,
}

func (e TorrentContentOrderByField) IsValid() bool {
	switch e {
	case TorrentContentOrderByFieldRelevance, TorrentContentOrderByFieldPublishedAt, TorrentContentOrderByFieldUpdatedAt, TorrentContentOrderByFieldSize, TorrentContentOrderByFieldFilesCount, TorrentContentOrderByFieldSeeders, TorrentContentOrderByFieldLeechers, TorrentContentOrderByFieldName, TorrentContentOrderByFieldInfoHash:
		return true
	}
	return false
}

func (e TorrentContentOrderByField) String() string {
	return string(e)
}

func (e *TorrentContentOrderByField) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TorrentContentOrderByField(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TorrentContentOrderByField", str)
	}
	return nil
}

func (e TorrentContentOrderByField) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TorrentFilesOrderByField string

const (
	TorrentFilesOrderByFieldIndex     TorrentFilesOrderByField = "index"
	TorrentFilesOrderByFieldPath      TorrentFilesOrderByField = "path"
	TorrentFilesOrderByFieldExtension TorrentFilesOrderByField = "extension"
	TorrentFilesOrderByFieldSize      TorrentFilesOrderByField = "size"
)

var AllTorrentFilesOrderByField = []TorrentFilesOrderByField{
	TorrentFilesOrderByFieldIndex,
	TorrentFilesOrderByFieldPath,
	TorrentFilesOrderByFieldExtension,
	TorrentFilesOrderByFieldSize,
}

func (e TorrentFilesOrderByField) IsValid() bool {
	switch e {
	case TorrentFilesOrderByFieldIndex, TorrentFilesOrderByFieldPath, TorrentFilesOrderByFieldExtension, TorrentFilesOrderByFieldSize:
		return true
	}
	return false
}

func (e TorrentFilesOrderByField) String() string {
	return string(e)
}

func (e *TorrentFilesOrderByField) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TorrentFilesOrderByField(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TorrentFilesOrderByField", str)
	}
	return nil
}

func (e TorrentFilesOrderByField) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
