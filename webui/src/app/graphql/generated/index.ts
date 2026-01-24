// THIS FILE IS GENERATED, DO NOT EDIT!

import { gql } from 'apollo-angular';
import { Injectable } from '@angular/core';
import * as Apollo from 'apollo-angular';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Date: { input: string; output: string; }
  DateTime: { input: string; output: string; }
  Duration: { input: string; output: string; }
  Hash20: { input: string; output: string; }
  JSON: { input: unknown; output: unknown; }
  Ref: { input: string; output: string; }
  Void: { input: void; output: void; }
  Year: { input: number; output: number; }
};

export type ApiKey = {
  __typename?: 'APIKey';
  createdAt: Scalars['DateTime']['output'];
  expiresAt?: Maybe<Scalars['DateTime']['output']>;
  id: Scalars['Int']['output'];
  name: Scalars['String']['output'];
  user: User;
  userId: Scalars['Int']['output'];
};

export type AuthMutation = {
  __typename?: 'AuthMutation';
  deleteInvitation?: Maybe<Scalars['Void']['output']>;
  deleteRole?: Maybe<Scalars['Void']['output']>;
  deleteUser?: Maybe<Scalars['Void']['output']>;
  invite: Invitation;
  putRole: Role;
  setUserEnabled: User;
  setUserRole: User;
};


export type AuthMutationDeleteInvitationArgs = {
  code: Scalars['String']['input'];
};


export type AuthMutationDeleteRoleArgs = {
  role: Scalars['String']['input'];
};


export type AuthMutationDeleteUserArgs = {
  userId: Scalars['Int']['input'];
};


export type AuthMutationInviteArgs = {
  input: InviteInput;
};


export type AuthMutationPutRoleArgs = {
  objectActions: Array<AuthObjectActionInput>;
  role: Scalars['String']['input'];
};


export type AuthMutationSetUserEnabledArgs = {
  enabled: Scalars['Boolean']['input'];
  userId: Scalars['Int']['input'];
};


export type AuthMutationSetUserRoleArgs = {
  roleName: Scalars['String']['input'];
  userId: Scalars['Int']['input'];
};

export type AuthObjectAction = {
  __typename?: 'AuthObjectAction';
  action: Scalars['String']['output'];
  namespace: Scalars['String']['output'];
  object: Scalars['String']['output'];
};

export type AuthObjectActionInput = {
  action: Scalars['String']['input'];
  namespace: Scalars['String']['input'];
  object: Scalars['String']['input'];
};

export type AuthQuery = {
  __typename?: 'AuthQuery';
  listInvitations: ListInvitationsResult;
  listObjectActions: Array<AuthObjectAction>;
  listRoles: Array<Role>;
  listUsers: ListUsersResult;
};


export type AuthQueryListInvitationsArgs = {
  input?: InputMaybe<ListInvitationsInput>;
};


export type AuthQueryListUsersArgs = {
  input?: InputMaybe<ListUsersInput>;
};

export type AuthSubject = {
  __typename?: 'AuthSubject';
  name: Scalars['String']['output'];
  type: AuthSubjectType;
};

export type AuthSubjectType =
  | 'role';

export type ConfigMutation = {
  __typename?: 'ConfigMutation';
  delete: ConfigParam;
  save: ConfigParam;
};


export type ConfigMutationDeleteArgs = {
  ref: Scalars['Ref']['input'];
};


export type ConfigMutationSaveArgs = {
  ref: Scalars['Ref']['input'];
  value: Scalars['JSON']['input'];
};

export type ConfigParam = {
  __typename?: 'ConfigParam';
  default: Scalars['JSON']['output'];
  description?: Maybe<Scalars['String']['output']>;
  dynamic: Scalars['Boolean']['output'];
  jsonSchema: Scalars['JSON']['output'];
  pending: Scalars['Boolean']['output'];
  plugin: Scalars['Ref']['output'];
  ref: Scalars['Ref']['output'];
  source: Scalars['String']['output'];
  value: Scalars['JSON']['output'];
};

export type ConfigQuery = {
  __typename?: 'ConfigQuery';
  params: Array<ConfigParam>;
  pending: Scalars['Boolean']['output'];
};

export type Content = {
  __typename?: 'Content';
  adult?: Maybe<Scalars['Boolean']['output']>;
  attributes: Array<ContentAttribute>;
  collections: Array<ContentCollection>;
  createdAt: Scalars['DateTime']['output'];
  externalLinks: Array<ExternalLink>;
  id: Scalars['String']['output'];
  metadataSource: MetadataSource;
  originalLanguage?: Maybe<LanguageInfo>;
  originalTitle?: Maybe<Scalars['String']['output']>;
  overview?: Maybe<Scalars['String']['output']>;
  popularity?: Maybe<Scalars['Float']['output']>;
  releaseDate?: Maybe<Scalars['Date']['output']>;
  releaseYear?: Maybe<Scalars['Year']['output']>;
  runtime?: Maybe<Scalars['Int']['output']>;
  source: Scalars['String']['output'];
  title: Scalars['String']['output'];
  type: ContentType;
  updatedAt: Scalars['DateTime']['output'];
  voteAverage?: Maybe<Scalars['Float']['output']>;
  voteCount?: Maybe<Scalars['Int']['output']>;
};

export type ContentAttribute = {
  __typename?: 'ContentAttribute';
  createdAt: Scalars['DateTime']['output'];
  key: Scalars['String']['output'];
  metadataSource: MetadataSource;
  source: Scalars['String']['output'];
  updatedAt: Scalars['DateTime']['output'];
  value: Scalars['String']['output'];
};

export type ContentCollection = {
  __typename?: 'ContentCollection';
  createdAt: Scalars['DateTime']['output'];
  id: Scalars['String']['output'];
  metadataSource: MetadataSource;
  name: Scalars['String']['output'];
  source: Scalars['String']['output'];
  type: Scalars['String']['output'];
  updatedAt: Scalars['DateTime']['output'];
};

export type ContentType =
  | 'audiobook'
  | 'comic'
  | 'ebook'
  | 'game'
  | 'movie'
  | 'music'
  | 'software'
  | 'tv_show'
  | 'xxx';

export type CreateApiKeyInput = {
  expiry?: InputMaybe<Scalars['Duration']['input']>;
  name: Scalars['String']['input'];
  permissions: Array<AuthObjectActionInput>;
};

export type CreateApiKeyResult = {
  __typename?: 'CreateAPIKeyResult';
  apiKey: Scalars['String']['output'];
  expiresAt?: Maybe<Scalars['DateTime']['output']>;
  id: Scalars['Int']['output'];
  name: Scalars['String']['output'];
};

export type Episodes = {
  __typename?: 'Episodes';
  label: Scalars['String']['output'];
  seasons: Array<Season>;
};

export type ExternalLink = {
  __typename?: 'ExternalLink';
  metadataSource: MetadataSource;
  url: Scalars['String']['output'];
};

export type FacetInput = {
  aggregate?: InputMaybe<Scalars['Boolean']['input']>;
  filter?: InputMaybe<Array<Scalars['String']['input']>>;
  key: FacetKey;
  logic?: InputMaybe<FacetLogic>;
};

export type FacetKey =
  | 'content_genre'
  | 'content_type'
  | 'file_type'
  | 'language'
  | 'release_year'
  | 'tag'
  | 'torrent_source'
  | 'video_3d'
  | 'video_codec'
  | 'video_modifier'
  | 'video_resolution'
  | 'video_source';

export type FacetLogic =
  | 'and'
  | 'or';

export type FacetResult = {
  __typename?: 'FacetResult';
  items: Array<FacetResultItem>;
  key: FacetKey;
  label: Scalars['String']['output'];
  logic: FacetLogic;
};

export type FacetResultItem = {
  __typename?: 'FacetResultItem';
  count: Scalars['Int']['output'];
  isEstimate: Scalars['Boolean']['output'];
  label: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export type FileType =
  | 'archive'
  | 'audio'
  | 'data'
  | 'document'
  | 'image'
  | 'software'
  | 'subtitles'
  | 'video';

export type FilesStatus =
  | 'multi'
  | 'no_info'
  | 'over_threshold'
  | 'single';

export type HealthCheck = {
  __typename?: 'HealthCheck';
  error?: Maybe<Scalars['String']['output']>;
  key: Scalars['String']['output'];
  status: HealthStatus;
  timestamp?: Maybe<Scalars['DateTime']['output']>;
};

export type HealthQuery = {
  __typename?: 'HealthQuery';
  checks: Array<HealthCheck>;
  status: HealthStatus;
};

export type HealthStatus =
  | 'down'
  | 'inactive'
  | 'unknown'
  | 'up';

export type IndexQuery = {
  __typename?: 'IndexQuery';
  default: Scalars['Ref']['output'];
  infos: Array<SearchIndexInfo>;
};

export type Invitation = {
  __typename?: 'Invitation';
  claimedBy?: Maybe<User>;
  code: Scalars['String']['output'];
  createdAt: Scalars['DateTime']['output'];
  createdBy?: Maybe<User>;
  email?: Maybe<Scalars['String']['output']>;
  expiresAt?: Maybe<Scalars['DateTime']['output']>;
  role: Scalars['String']['output'];
};

export type InviteInput = {
  email?: InputMaybe<Scalars['String']['input']>;
  expiry?: InputMaybe<Scalars['Duration']['input']>;
  role?: InputMaybe<Scalars['String']['input']>;
};

export type Language =
  | 'af'
  | 'ar'
  | 'az'
  | 'be'
  | 'bg'
  | 'bs'
  | 'ca'
  | 'ce'
  | 'co'
  | 'cs'
  | 'cy'
  | 'da'
  | 'de'
  | 'el'
  | 'en'
  | 'es'
  | 'et'
  | 'eu'
  | 'fa'
  | 'fi'
  | 'fr'
  | 'he'
  | 'hi'
  | 'hr'
  | 'hu'
  | 'hy'
  | 'id'
  | 'is'
  | 'it'
  | 'ja'
  | 'ka'
  | 'ko'
  | 'ku'
  | 'lt'
  | 'lv'
  | 'mi'
  | 'mk'
  | 'ml'
  | 'mn'
  | 'ms'
  | 'mt'
  | 'nl'
  | 'no'
  | 'pl'
  | 'pt'
  | 'ro'
  | 'ru'
  | 'sa'
  | 'sk'
  | 'sl'
  | 'sm'
  | 'so'
  | 'sr'
  | 'sv'
  | 'ta'
  | 'th'
  | 'tr'
  | 'uk'
  | 'vi'
  | 'yi'
  | 'zh'
  | 'zu';

export type LanguageInfo = {
  __typename?: 'LanguageInfo';
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type ListInvitationsInput = {
  pagination?: InputMaybe<PaginationInput>;
};

export type ListInvitationsResult = {
  __typename?: 'ListInvitationsResult';
  invitations: Array<Invitation>;
  totalCount: Scalars['Int']['output'];
};

export type ListUsersInput = {
  pagination?: InputMaybe<PaginationInput>;
  usernameLike?: InputMaybe<Scalars['String']['input']>;
};

export type ListUsersResult = {
  __typename?: 'ListUsersResult';
  totalCount: Scalars['Int']['output'];
  users: Array<User>;
};

export type LoginResult = {
  __typename?: 'LoginResult';
  permissions: Array<Permission>;
  token: Scalars['String']['output'];
  user: User;
};

export type MetadataSource = {
  __typename?: 'MetadataSource';
  key: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type MetricsBucketDuration =
  | 'day'
  | 'hour'
  | 'minute';

export type Mutation = {
  __typename?: 'Mutation';
  auth: AuthMutation;
  config: ConfigMutation;
  queue: QueueMutation;
  self: SelfMutation;
  target: TargetMutation;
  torrent: TorrentMutation;
  worker: WorkerMutation;
};

export type PaginationInput = {
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
  page?: InputMaybe<Scalars['Int']['input']>;
};

export type PasswordEntropyResult = {
  __typename?: 'PasswordEntropyResult';
  entropy: Scalars['Float']['output'];
  minEntropy: Scalars['Float']['output'];
  valid: Scalars['Boolean']['output'];
};

export type Permission = {
  __typename?: 'Permission';
  core: Scalars['Boolean']['output'];
  objectAction: AuthObjectAction;
  subject: AuthSubject;
};

export type PluginInfo = {
  __typename?: 'PluginInfo';
  dependsOn: Array<Scalars['Ref']['output']>;
  description?: Maybe<Scalars['String']['output']>;
  enabled: Scalars['Boolean']['output'];
  ref: Scalars['Ref']['output'];
  requiredBy: Array<Scalars['Ref']['output']>;
};

export type PluginQuery = {
  __typename?: 'PluginQuery';
  list: Array<PluginInfo>;
};

export type Query = {
  __typename?: 'Query';
  auth: AuthQuery;
  config: ConfigQuery;
  health: HealthQuery;
  index: IndexQuery;
  plugin: PluginQuery;
  queue: QueueQuery;
  self: SelfQuery;
  target: TargetQuery;
  torrent: TorrentQuery;
  version: Scalars['String']['output'];
  worker: WorkerQuery;
};

export type QueueEnqueueReprocessTorrentsBatchInput = {
  apisDisabled?: InputMaybe<Scalars['Boolean']['input']>;
  batchSize?: InputMaybe<Scalars['Int']['input']>;
  chunkSize?: InputMaybe<Scalars['Int']['input']>;
  classifierRematch?: InputMaybe<Scalars['Boolean']['input']>;
  classifierWorkflow?: InputMaybe<Scalars['String']['input']>;
  contentTypes?: InputMaybe<Array<InputMaybe<ContentType>>>;
  localSearchDisabled?: InputMaybe<Scalars['Boolean']['input']>;
  orphans?: InputMaybe<Scalars['Boolean']['input']>;
  purge?: InputMaybe<Scalars['Boolean']['input']>;
};

export type QueueJob = {
  __typename?: 'QueueJob';
  createdAt: Scalars['DateTime']['output'];
  error?: Maybe<Scalars['String']['output']>;
  id: Scalars['ID']['output'];
  maxRetries: Scalars['Int']['output'];
  payload: Scalars['String']['output'];
  priority: Scalars['Int']['output'];
  queue: Scalars['String']['output'];
  ranAt?: Maybe<Scalars['DateTime']['output']>;
  retries: Scalars['Int']['output'];
  runAfter: Scalars['DateTime']['output'];
  status: QueueJobStatus;
};

export type QueueJobQueueAgg = {
  __typename?: 'QueueJobQueueAgg';
  count: Scalars['Int']['output'];
  label: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export type QueueJobQueueFacetInput = {
  aggregate?: InputMaybe<Scalars['Boolean']['input']>;
  filter?: InputMaybe<Array<Scalars['String']['input']>>;
};

export type QueueJobStatus =
  | 'failed'
  | 'pending'
  | 'processed'
  | 'retry';

export type QueueJobStatusAgg = {
  __typename?: 'QueueJobStatusAgg';
  count: Scalars['Int']['output'];
  label: Scalars['String']['output'];
  value: QueueJobStatus;
};

export type QueueJobStatusFacetInput = {
  aggregate?: InputMaybe<Scalars['Boolean']['input']>;
  filter?: InputMaybe<Array<QueueJobStatus>>;
};

export type QueueJobsAggregations = {
  __typename?: 'QueueJobsAggregations';
  queue?: Maybe<Array<QueueJobQueueAgg>>;
  status?: Maybe<Array<QueueJobStatusAgg>>;
};

export type QueueJobsFacetsInput = {
  queue?: InputMaybe<QueueJobQueueFacetInput>;
  status?: InputMaybe<QueueJobStatusFacetInput>;
};

export type QueueJobsOrderByField =
  | 'created_at'
  | 'priority'
  | 'ran_at';

export type QueueJobsOrderByInput = {
  descending?: InputMaybe<Scalars['Boolean']['input']>;
  field: QueueJobsOrderByField;
};

export type QueueJobsQueryInput = {
  facets?: InputMaybe<QueueJobsFacetsInput>;
  hasNextPage?: InputMaybe<Scalars['Boolean']['input']>;
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
  orderBy?: InputMaybe<Array<QueueJobsOrderByInput>>;
  page?: InputMaybe<Scalars['Int']['input']>;
  queues?: InputMaybe<Array<Scalars['String']['input']>>;
  statuses?: InputMaybe<Array<QueueJobStatus>>;
  totalCount?: InputMaybe<Scalars['Boolean']['input']>;
};

export type QueueJobsQueryResult = {
  __typename?: 'QueueJobsQueryResult';
  aggregations: QueueJobsAggregations;
  hasNextPage?: Maybe<Scalars['Boolean']['output']>;
  items: Array<QueueJob>;
  totalCount: Scalars['Int']['output'];
};

export type QueueMetricsBucket = {
  __typename?: 'QueueMetricsBucket';
  count: Scalars['Int']['output'];
  createdAtBucket: Scalars['DateTime']['output'];
  latency?: Maybe<Scalars['Duration']['output']>;
  queue: Scalars['String']['output'];
  ranAtBucket?: Maybe<Scalars['DateTime']['output']>;
  status: QueueJobStatus;
};

export type QueueMetricsQueryInput = {
  bucketDuration: MetricsBucketDuration;
  endTime?: InputMaybe<Scalars['DateTime']['input']>;
  queues?: InputMaybe<Array<Scalars['String']['input']>>;
  startTime?: InputMaybe<Scalars['DateTime']['input']>;
  statuses?: InputMaybe<Array<QueueJobStatus>>;
};

export type QueueMetricsQueryResult = {
  __typename?: 'QueueMetricsQueryResult';
  buckets: Array<QueueMetricsBucket>;
};

export type QueueMutation = {
  __typename?: 'QueueMutation';
  enqueueReprocessTorrentsBatch?: Maybe<Scalars['Void']['output']>;
  purgeJobs?: Maybe<Scalars['Void']['output']>;
};


export type QueueMutationEnqueueReprocessTorrentsBatchArgs = {
  input?: InputMaybe<QueueEnqueueReprocessTorrentsBatchInput>;
};


export type QueueMutationPurgeJobsArgs = {
  input: QueuePurgeJobsInput;
};

export type QueuePurgeJobsInput = {
  queues?: InputMaybe<Array<Scalars['String']['input']>>;
  statuses?: InputMaybe<Array<QueueJobStatus>>;
};

export type QueueQuery = {
  __typename?: 'QueueQuery';
  jobs: QueueJobsQueryResult;
  metrics: QueueMetricsQueryResult;
};


export type QueueQueryJobsArgs = {
  input: QueueJobsQueryInput;
};


export type QueueQueryMetricsArgs = {
  input: QueueMetricsQueryInput;
};

export type RegisterInput = {
  email?: InputMaybe<Scalars['String']['input']>;
  invitationCode?: InputMaybe<Scalars['String']['input']>;
  password: Scalars['String']['input'];
  username: Scalars['String']['input'];
};

export type RegisterResult = {
  __typename?: 'RegisterResult';
  user: User;
};

export type Role = {
  __typename?: 'Role';
  core: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  permissions: Array<Permission>;
};

export type SearchIndexInfo = {
  __typename?: 'SearchIndexInfo';
  name: Scalars['String']['output'];
  ref: Scalars['Ref']['output'];
  resultTypes: Array<SearchResultType>;
};

export type SearchInput = {
  aggregationBudget?: InputMaybe<Scalars['Float']['input']>;
  cached?: InputMaybe<Scalars['Boolean']['input']>;
  criteria?: InputMaybe<Scalars['JSON']['input']>;
  facets?: InputMaybe<Array<FacetInput>>;
  /** hasNextPage if true, the search result will include the hasNextPage field, indicating if there are more results to fetch */
  hasNextPage?: InputMaybe<Scalars['Boolean']['input']>;
  index?: InputMaybe<Scalars['Ref']['input']>;
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
  orderBy?: InputMaybe<Array<TorrentContentOrderByInput>>;
  page?: InputMaybe<Scalars['Int']['input']>;
  queryString?: InputMaybe<Scalars['String']['input']>;
  totalCount?: InputMaybe<Scalars['Boolean']['input']>;
};

export type SearchResult = {
  /** hasNextPage is true if there are more results to fetch */
  hasNextPage?: Maybe<Scalars['Boolean']['output']>;
  totalCount?: Maybe<Scalars['Int']['output']>;
  totalCountIsEstimate: Scalars['Boolean']['output'];
};

export type SearchResultType =
  | 'content'
  | 'torrent'
  | 'torrent_content'
  | 'torrent_file';

export type Season = {
  __typename?: 'Season';
  episodes?: Maybe<Array<Scalars['Int']['output']>>;
  season: Scalars['Int']['output'];
};

export type Self = {
  __typename?: 'Self';
  apiKey?: Maybe<ApiKey>;
  permissions: Array<AuthObjectAction>;
  user?: Maybe<User>;
};

export type SelfMutation = {
  __typename?: 'SelfMutation';
  createAPIKey: CreateApiKeyResult;
  deleteAPIKey?: Maybe<Scalars['Void']['output']>;
  login: LoginResult;
  register: RegisterResult;
};


export type SelfMutationCreateApiKeyArgs = {
  input: CreateApiKeyInput;
};


export type SelfMutationDeleteApiKeyArgs = {
  id: Scalars['Int']['input'];
};


export type SelfMutationLoginArgs = {
  password: Scalars['String']['input'];
  username: Scalars['String']['input'];
};


export type SelfMutationRegisterArgs = {
  input: RegisterInput;
};

export type SelfQuery = {
  __typename?: 'SelfQuery';
  apiKeys: Array<ApiKey>;
  identity: Self;
  passwordEntropy: PasswordEntropyResult;
};


export type SelfQueryPasswordEntropyArgs = {
  password: Scalars['String']['input'];
};

export type SendInput = {
  data?: InputMaybe<Scalars['JSON']['input']>;
  index?: InputMaybe<Scalars['Ref']['input']>;
  infoHashes: Array<Scalars['Hash20']['input']>;
  target: Scalars['Ref']['input'];
};

export type SendResult = {
  __typename?: 'SendResult';
  data?: Maybe<Scalars['JSON']['output']>;
  missingInfoHashes: Array<Scalars['Hash20']['output']>;
};

export type SuggestTagsQueryInput = {
  exclusions?: InputMaybe<Array<Scalars['String']['input']>>;
  prefix?: InputMaybe<Scalars['String']['input']>;
};

export type SuggestedTag = {
  __typename?: 'SuggestedTag';
  count: Scalars['Int']['output'];
  name: Scalars['String']['output'];
};

export type Target = {
  __typename?: 'Target';
  dataSchema?: Maybe<Scalars['JSON']['output']>;
  name: Scalars['String']['output'];
  ref: Scalars['Ref']['output'];
  uiSchema?: Maybe<Scalars['JSON']['output']>;
};

export type TargetMutation = {
  __typename?: 'TargetMutation';
  send: SendResult;
};


export type TargetMutationSendArgs = {
  input: SendInput;
};

export type TargetQuery = {
  __typename?: 'TargetQuery';
  default?: Maybe<Scalars['Ref']['output']>;
  targets: Array<Target>;
};

export type Torrent = {
  __typename?: 'Torrent';
  createdAt: Scalars['DateTime']['output'];
  extension?: Maybe<Scalars['String']['output']>;
  fileType?: Maybe<FileType>;
  fileTypes?: Maybe<Array<FileType>>;
  files?: Maybe<Array<TorrentFile>>;
  filesCount?: Maybe<Scalars['Int']['output']>;
  filesStatus: FilesStatus;
  hasFilesInfo: Scalars['Boolean']['output'];
  infoHash: Scalars['Hash20']['output'];
  leechers?: Maybe<Scalars['Int']['output']>;
  magnetUri: Scalars['String']['output'];
  name: Scalars['String']['output'];
  seeders?: Maybe<Scalars['Int']['output']>;
  singleFile?: Maybe<Scalars['Boolean']['output']>;
  size: Scalars['Int']['output'];
  sources: Array<TorrentSourceInfo>;
  tagNames: Array<Scalars['String']['output']>;
  updatedAt: Scalars['DateTime']['output'];
};

export type TorrentContent = {
  __typename?: 'TorrentContent';
  content?: Maybe<Content>;
  contentId?: Maybe<Scalars['String']['output']>;
  contentSource?: Maybe<Scalars['String']['output']>;
  contentType?: Maybe<ContentType>;
  createdAt: Scalars['DateTime']['output'];
  episodes?: Maybe<Episodes>;
  id: Scalars['ID']['output'];
  infoHash: Scalars['Hash20']['output'];
  languages?: Maybe<Array<LanguageInfo>>;
  leechers?: Maybe<Scalars['Int']['output']>;
  publishedAt: Scalars['DateTime']['output'];
  releaseGroup?: Maybe<Scalars['String']['output']>;
  seeders?: Maybe<Scalars['Int']['output']>;
  title: Scalars['String']['output'];
  torrent: Torrent;
  updatedAt: Scalars['DateTime']['output'];
  video3d?: Maybe<Video3D>;
  videoCodec?: Maybe<VideoCodec>;
  videoModifier?: Maybe<VideoModifier>;
  videoResolution?: Maybe<VideoResolution>;
  videoSource?: Maybe<VideoSource>;
};

export type TorrentContentOrderByField =
  | 'files_count'
  | 'info_hash'
  | 'leechers'
  | 'name'
  | 'published_at'
  | 'relevance'
  | 'seeders'
  | 'size'
  | 'updated_at';

export type TorrentContentOrderByInput = {
  descending?: InputMaybe<Scalars['Boolean']['input']>;
  field: TorrentContentOrderByField;
};

export type TorrentContentSearchResult = SearchResult & {
  __typename?: 'TorrentContentSearchResult';
  facets: Array<FacetResult>;
  /** hasNextPage is true if there are more results to fetch */
  hasNextPage?: Maybe<Scalars['Boolean']['output']>;
  items: Array<TorrentContent>;
  totalCount?: Maybe<Scalars['Int']['output']>;
  totalCountIsEstimate: Scalars['Boolean']['output'];
};

export type TorrentFile = {
  __typename?: 'TorrentFile';
  createdAt: Scalars['DateTime']['output'];
  extension?: Maybe<Scalars['String']['output']>;
  fileType?: Maybe<FileType>;
  index: Scalars['Int']['output'];
  infoHash: Scalars['Hash20']['output'];
  path: Scalars['String']['output'];
  size: Scalars['Int']['output'];
  updatedAt: Scalars['DateTime']['output'];
};

export type TorrentFilesOrderByField =
  | 'extension'
  | 'index'
  | 'path'
  | 'size';

export type TorrentFilesOrderByInput = {
  descending?: InputMaybe<Scalars['Boolean']['input']>;
  field: TorrentFilesOrderByField;
};

export type TorrentFilesSearchResult = SearchResult & {
  __typename?: 'TorrentFilesSearchResult';
  facets: Array<FacetResult>;
  hasNextPage?: Maybe<Scalars['Boolean']['output']>;
  items: Array<TorrentFile>;
  totalCount?: Maybe<Scalars['Int']['output']>;
  totalCountIsEstimate: Scalars['Boolean']['output'];
};

export type TorrentListSourcesResult = {
  __typename?: 'TorrentListSourcesResult';
  sources: Array<TorrentSource>;
};

export type TorrentMetricsBucket = {
  __typename?: 'TorrentMetricsBucket';
  bucket: Scalars['DateTime']['output'];
  count: Scalars['Int']['output'];
  source: Scalars['String']['output'];
  updated: Scalars['Boolean']['output'];
};

export type TorrentMetricsQueryInput = {
  bucketDuration: MetricsBucketDuration;
  endTime?: InputMaybe<Scalars['DateTime']['input']>;
  sources?: InputMaybe<Array<Scalars['String']['input']>>;
  startTime?: InputMaybe<Scalars['DateTime']['input']>;
};

export type TorrentMetricsQueryResult = {
  __typename?: 'TorrentMetricsQueryResult';
  buckets: Array<TorrentMetricsBucket>;
};

export type TorrentMutation = {
  __typename?: 'TorrentMutation';
  delete?: Maybe<Scalars['Void']['output']>;
  deleteTags?: Maybe<Scalars['Void']['output']>;
  putTags?: Maybe<Scalars['Void']['output']>;
  reprocess?: Maybe<Scalars['Void']['output']>;
  setTags?: Maybe<Scalars['Void']['output']>;
};


export type TorrentMutationDeleteArgs = {
  infoHashes: Array<Scalars['Hash20']['input']>;
};


export type TorrentMutationDeleteTagsArgs = {
  infoHashes?: InputMaybe<Array<Scalars['Hash20']['input']>>;
  tagNames?: InputMaybe<Array<Scalars['String']['input']>>;
};


export type TorrentMutationPutTagsArgs = {
  infoHashes: Array<Scalars['Hash20']['input']>;
  tagNames: Array<Scalars['String']['input']>;
};


export type TorrentMutationReprocessArgs = {
  input: TorrentReprocessInput;
};


export type TorrentMutationSetTagsArgs = {
  infoHashes: Array<Scalars['Hash20']['input']>;
  tagNames: Array<Scalars['String']['input']>;
};

export type TorrentQuery = {
  __typename?: 'TorrentQuery';
  listSources: TorrentListSourcesResult;
  metrics: TorrentMetricsQueryResult;
  searchTorrentContent: TorrentContentSearchResult;
  searchTorrentFiles: TorrentFilesSearchResult;
  suggestTags: TorrentSuggestTagsResult;
};


export type TorrentQueryMetricsArgs = {
  input: TorrentMetricsQueryInput;
};


export type TorrentQuerySearchTorrentContentArgs = {
  input: SearchInput;
};


export type TorrentQuerySearchTorrentFilesArgs = {
  input: SearchInput;
};


export type TorrentQuerySuggestTagsArgs = {
  input?: InputMaybe<SuggestTagsQueryInput>;
};

export type TorrentReprocessInput = {
  apisDisabled?: InputMaybe<Scalars['Boolean']['input']>;
  classifierRematch?: InputMaybe<Scalars['Boolean']['input']>;
  classifierWorkflow?: InputMaybe<Scalars['String']['input']>;
  infoHashes: Array<Scalars['Hash20']['input']>;
  localSearchDisabled?: InputMaybe<Scalars['Boolean']['input']>;
};

export type TorrentSource = {
  __typename?: 'TorrentSource';
  key: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type TorrentSourceInfo = {
  __typename?: 'TorrentSourceInfo';
  importId?: Maybe<Scalars['String']['output']>;
  key: Scalars['String']['output'];
  leechers?: Maybe<Scalars['Int']['output']>;
  name: Scalars['String']['output'];
  seeders?: Maybe<Scalars['Int']['output']>;
};

export type TorrentSuggestTagsResult = {
  __typename?: 'TorrentSuggestTagsResult';
  suggestions: Array<SuggestedTag>;
};

export type User = {
  __typename?: 'User';
  createdAt: Scalars['DateTime']['output'];
  email?: Maybe<Scalars['String']['output']>;
  id: Scalars['Int']['output'];
  lastLoginAt?: Maybe<Scalars['DateTime']['output']>;
  role: Scalars['String']['output'];
  updatedAt: Scalars['DateTime']['output'];
  username: Scalars['String']['output'];
};

export type Video3D =
  | 'V3D'
  | 'V3DOU'
  | 'V3DSBS';

export type VideoCodec =
  | 'DivX'
  | 'H264'
  | 'MPEG2'
  | 'MPEG4'
  | 'XviD'
  | 'x264'
  | 'x265';

export type VideoModifier =
  | 'BRDISK'
  | 'RAWHD'
  | 'REGIONAL'
  | 'REMUX'
  | 'SCREENER';

export type VideoResolution =
  | 'V360p'
  | 'V480p'
  | 'V540p'
  | 'V576p'
  | 'V720p'
  | 'V1080p'
  | 'V1440p'
  | 'V2160p'
  | 'V4320p';

export type VideoSource =
  | 'BluRay'
  | 'CAM'
  | 'DVD'
  | 'TELECINE'
  | 'TELESYNC'
  | 'TV'
  | 'WEBDL'
  | 'WEBRip'
  | 'WORKPRINT';

export type Worker = {
  __typename?: 'Worker';
  dependsOn: Array<Scalars['Ref']['output']>;
  error?: Maybe<Scalars['String']['output']>;
  label: Scalars['String']['output'];
  ref: Scalars['Ref']['output'];
  requiredBy: Array<Scalars['Ref']['output']>;
  state: WorkerState;
};

export type WorkerListAllQueryResult = {
  __typename?: 'WorkerListAllQueryResult';
  workers: Array<Worker>;
};

export type WorkerMutation = {
  __typename?: 'WorkerMutation';
  restart: WorkerListAllQueryResult;
  shutdown: WorkerListAllQueryResult;
  start: WorkerListAllQueryResult;
};


export type WorkerMutationRestartArgs = {
  refs: Array<Scalars['Ref']['input']>;
};


export type WorkerMutationShutdownArgs = {
  refs: Array<Scalars['Ref']['input']>;
};


export type WorkerMutationStartArgs = {
  refs: Array<Scalars['Ref']['input']>;
};

export type WorkerQuery = {
  __typename?: 'WorkerQuery';
  listAll: WorkerListAllQueryResult;
};

export type WorkerState =
  | 'error'
  | 'idle'
  | 'running'
  | 'shutdown'
  | 'startup';

export type ApiKeyFragment = { __typename?: 'APIKey', id: number, userId: number, name: string, expiresAt?: string | null, createdAt: string };

export type AuthObjectActionFragment = { __typename?: 'AuthObjectAction', namespace: string, object: string, action: string };

export type ConfigParamFragment = { __typename?: 'ConfigParam', ref: string, plugin: string, description?: string | null, value: unknown, source: string, default: unknown, dynamic: boolean, pending: boolean, jsonSchema: unknown };

export type ContentFragment = { __typename?: 'Content', type: ContentType, source: string, id: string, title: string, releaseDate?: string | null, releaseYear?: number | null, overview?: string | null, runtime?: number | null, voteAverage?: number | null, voteCount?: number | null, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string }, originalLanguage?: { __typename?: 'LanguageInfo', id: string, name: string } | null, attributes: Array<{ __typename?: 'ContentAttribute', source: string, key: string, value: string, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }>, collections: Array<{ __typename?: 'ContentCollection', type: string, source: string, id: string, name: string, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }>, externalLinks: Array<{ __typename?: 'ExternalLink', url: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }> };

export type FacetResultItemFragment = { __typename?: 'FacetResultItem', value: string, label: string, count: number, isEstimate: boolean };

export type FacetResultFragment = { __typename?: 'FacetResult', key: FacetKey, label: string, logic: FacetLogic, items: Array<{ __typename?: 'FacetResultItem', value: string, label: string, count: number, isEstimate: boolean }> };

export type InvitationFragment = { __typename?: 'Invitation', code: string, role: string, email?: string | null, expiresAt?: string | null, createdAt: string, createdBy?: { __typename?: 'User', id: number, username: string, role: string, lastLoginAt?: string | null, createdAt: string, updatedAt: string } | null, claimedBy?: { __typename?: 'User', id: number, username: string, role: string, lastLoginAt?: string | null, createdAt: string, updatedAt: string } | null };

export type PasswordEntropyResultFragment = { __typename?: 'PasswordEntropyResult', entropy: number, minEntropy: number, valid: boolean };

export type PermissionFragment = { __typename?: 'Permission', core: boolean, subject: { __typename?: 'AuthSubject', type: AuthSubjectType, name: string }, objectAction: { __typename?: 'AuthObjectAction', namespace: string, object: string, action: string } };

export type PluginInfoFragment = { __typename?: 'PluginInfo', ref: string, description?: string | null, enabled: boolean, dependsOn: Array<string>, requiredBy: Array<string> };

export type QueueJobFragment = { __typename?: 'QueueJob', id: string, queue: string, status: QueueJobStatus, payload: string, priority: number, retries: number, maxRetries: number, runAfter: string, ranAt?: string | null, error?: string | null, createdAt: string };

export type QueueJobsQueryResultFragment = { __typename?: 'QueueJobsQueryResult', totalCount: number, hasNextPage?: boolean | null, items: Array<{ __typename?: 'QueueJob', id: string, queue: string, status: QueueJobStatus, payload: string, priority: number, retries: number, maxRetries: number, runAfter: string, ranAt?: string | null, error?: string | null, createdAt: string }>, aggregations: { __typename?: 'QueueJobsAggregations', queue?: Array<{ __typename?: 'QueueJobQueueAgg', value: string, label: string, count: number }> | null, status?: Array<{ __typename?: 'QueueJobStatusAgg', value: QueueJobStatus, label: string, count: number }> | null } };

export type RoleFragment = { __typename?: 'Role', name: string, core: boolean, permissions: Array<{ __typename?: 'Permission', core: boolean, subject: { __typename?: 'AuthSubject', type: AuthSubjectType, name: string }, objectAction: { __typename?: 'AuthObjectAction', namespace: string, object: string, action: string } }> };

export type SearchIndexInfoFragment = { __typename?: 'SearchIndexInfo', ref: string, name: string, resultTypes: Array<SearchResultType> };

export type SelfFragment = { __typename?: 'Self', user?: { __typename?: 'User', id: number, username: string, role: string, lastLoginAt?: string | null, createdAt: string, updatedAt: string } | null, apiKey?: { __typename?: 'APIKey', id: number, userId: number, name: string, expiresAt?: string | null, createdAt: string } | null, permissions: Array<{ __typename?: 'AuthObjectAction', namespace: string, object: string, action: string }> };

export type TargetFragment = { __typename?: 'Target', ref: string, name: string, dataSchema?: unknown | null, uiSchema?: unknown | null };

export type TorrentFragment = { __typename?: 'Torrent', infoHash: string, name: string, size: number, filesStatus: FilesStatus, filesCount?: number | null, hasFilesInfo: boolean, singleFile?: boolean | null, fileType?: FileType | null, seeders?: number | null, leechers?: number | null, tagNames: Array<string>, magnetUri: string, createdAt: string, updatedAt: string, sources: Array<{ __typename?: 'TorrentSourceInfo', key: string, name: string }> };

export type TorrentContentFragment = { __typename?: 'TorrentContent', id: string, infoHash: string, contentType?: ContentType | null, title: string, video3d?: Video3D | null, videoCodec?: VideoCodec | null, videoModifier?: VideoModifier | null, videoResolution?: VideoResolution | null, videoSource?: VideoSource | null, seeders?: number | null, leechers?: number | null, publishedAt: string, createdAt: string, updatedAt: string, torrent: { __typename?: 'Torrent', infoHash: string, name: string, size: number, filesStatus: FilesStatus, filesCount?: number | null, hasFilesInfo: boolean, singleFile?: boolean | null, fileType?: FileType | null, seeders?: number | null, leechers?: number | null, tagNames: Array<string>, magnetUri: string, createdAt: string, updatedAt: string, sources: Array<{ __typename?: 'TorrentSourceInfo', key: string, name: string }> }, content?: { __typename?: 'Content', type: ContentType, source: string, id: string, title: string, releaseDate?: string | null, releaseYear?: number | null, overview?: string | null, runtime?: number | null, voteAverage?: number | null, voteCount?: number | null, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string }, originalLanguage?: { __typename?: 'LanguageInfo', id: string, name: string } | null, attributes: Array<{ __typename?: 'ContentAttribute', source: string, key: string, value: string, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }>, collections: Array<{ __typename?: 'ContentCollection', type: string, source: string, id: string, name: string, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }>, externalLinks: Array<{ __typename?: 'ExternalLink', url: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }> } | null, languages?: Array<{ __typename?: 'LanguageInfo', id: string, name: string }> | null, episodes?: { __typename?: 'Episodes', label: string, seasons: Array<{ __typename?: 'Season', season: number, episodes?: Array<number> | null }> } | null };

export type TorrentContentSearchResultFragment = { __typename?: 'TorrentContentSearchResult', totalCount?: number | null, totalCountIsEstimate: boolean, hasNextPage?: boolean | null, items: Array<{ __typename?: 'TorrentContent', id: string, infoHash: string, contentType?: ContentType | null, title: string, video3d?: Video3D | null, videoCodec?: VideoCodec | null, videoModifier?: VideoModifier | null, videoResolution?: VideoResolution | null, videoSource?: VideoSource | null, seeders?: number | null, leechers?: number | null, publishedAt: string, createdAt: string, updatedAt: string, torrent: { __typename?: 'Torrent', infoHash: string, name: string, size: number, filesStatus: FilesStatus, filesCount?: number | null, hasFilesInfo: boolean, singleFile?: boolean | null, fileType?: FileType | null, seeders?: number | null, leechers?: number | null, tagNames: Array<string>, magnetUri: string, createdAt: string, updatedAt: string, sources: Array<{ __typename?: 'TorrentSourceInfo', key: string, name: string }> }, content?: { __typename?: 'Content', type: ContentType, source: string, id: string, title: string, releaseDate?: string | null, releaseYear?: number | null, overview?: string | null, runtime?: number | null, voteAverage?: number | null, voteCount?: number | null, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string }, originalLanguage?: { __typename?: 'LanguageInfo', id: string, name: string } | null, attributes: Array<{ __typename?: 'ContentAttribute', source: string, key: string, value: string, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }>, collections: Array<{ __typename?: 'ContentCollection', type: string, source: string, id: string, name: string, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }>, externalLinks: Array<{ __typename?: 'ExternalLink', url: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }> } | null, languages?: Array<{ __typename?: 'LanguageInfo', id: string, name: string }> | null, episodes?: { __typename?: 'Episodes', label: string, seasons: Array<{ __typename?: 'Season', season: number, episodes?: Array<number> | null }> } | null }>, facets: Array<{ __typename?: 'FacetResult', key: FacetKey, label: string, logic: FacetLogic, items: Array<{ __typename?: 'FacetResultItem', value: string, label: string, count: number, isEstimate: boolean }> }> };

export type TorrentFileFragment = { __typename?: 'TorrentFile', infoHash: string, index: number, path: string, size: number, fileType?: FileType | null, createdAt: string, updatedAt: string };

export type TorrentFilesSearchResultFragment = { __typename?: 'TorrentFilesSearchResult', totalCount?: number | null, hasNextPage?: boolean | null, items: Array<{ __typename?: 'TorrentFile', infoHash: string, index: number, path: string, size: number, fileType?: FileType | null, createdAt: string, updatedAt: string }> };

export type UserFragment = { __typename?: 'User', id: number, username: string, role: string, lastLoginAt?: string | null, createdAt: string, updatedAt: string };

export type WorkerFragment = { __typename?: 'Worker', ref: string, label: string, state: WorkerState, error?: string | null, requiredBy: Array<string>, dependsOn: Array<string> };

export type ConfigDeleteMutationVariables = Exact<{
  ref: Scalars['Ref']['input'];
}>;


export type ConfigDeleteMutation = { __typename?: 'Mutation', config: { __typename?: 'ConfigMutation', delete: { __typename?: 'ConfigParam', ref: string, plugin: string, description?: string | null, value: unknown, source: string, default: unknown, dynamic: boolean, pending: boolean, jsonSchema: unknown } } };

export type ConfigSaveMutationVariables = Exact<{
  ref: Scalars['Ref']['input'];
  value: Scalars['JSON']['input'];
}>;


export type ConfigSaveMutation = { __typename?: 'Mutation', config: { __typename?: 'ConfigMutation', save: { __typename?: 'ConfigParam', ref: string, plugin: string, description?: string | null, value: unknown, source: string, default: unknown, dynamic: boolean, pending: boolean, jsonSchema: unknown } } };

export type CreateApiKeyMutationVariables = Exact<{
  input: CreateApiKeyInput;
}>;


export type CreateApiKeyMutation = { __typename?: 'Mutation', self: { __typename?: 'SelfMutation', createAPIKey: { __typename?: 'CreateAPIKeyResult', id: number, apiKey: string, name: string, expiresAt?: string | null } } };

export type InvitationDeleteMutationVariables = Exact<{
  code: Scalars['String']['input'];
}>;


export type InvitationDeleteMutation = { __typename?: 'Mutation', auth: { __typename?: 'AuthMutation', deleteInvitation?: void | null } };

export type InviteMutationVariables = Exact<{
  input: InviteInput;
}>;


export type InviteMutation = { __typename?: 'Mutation', auth: { __typename?: 'AuthMutation', invite: { __typename?: 'Invitation', code: string, role: string, email?: string | null, expiresAt?: string | null, createdAt: string, createdBy?: { __typename?: 'User', id: number, username: string, role: string, lastLoginAt?: string | null, createdAt: string, updatedAt: string } | null, claimedBy?: { __typename?: 'User', id: number, username: string, role: string, lastLoginAt?: string | null, createdAt: string, updatedAt: string } | null } } };

export type LoginMutationVariables = Exact<{
  username: Scalars['String']['input'];
  password: Scalars['String']['input'];
}>;


export type LoginMutation = { __typename?: 'Mutation', self: { __typename?: 'SelfMutation', login: { __typename?: 'LoginResult', token: string, user: { __typename?: 'User', id: number, username: string, role: string, lastLoginAt?: string | null, createdAt: string, updatedAt: string }, permissions: Array<{ __typename?: 'Permission', core: boolean, subject: { __typename?: 'AuthSubject', type: AuthSubjectType, name: string }, objectAction: { __typename?: 'AuthObjectAction', namespace: string, object: string, action: string } }> } } };

export type PutRoleMutationVariables = Exact<{
  role: Scalars['String']['input'];
  objectActions: Array<AuthObjectActionInput> | AuthObjectActionInput;
}>;


export type PutRoleMutation = { __typename?: 'Mutation', auth: { __typename?: 'AuthMutation', putRole: { __typename?: 'Role', name: string, core: boolean, permissions: Array<{ __typename?: 'Permission', core: boolean, subject: { __typename?: 'AuthSubject', type: AuthSubjectType, name: string }, objectAction: { __typename?: 'AuthObjectAction', namespace: string, object: string, action: string } }> } } };

export type QueueEnqueueReprocessTorrentsBatchMutationVariables = Exact<{
  input: QueueEnqueueReprocessTorrentsBatchInput;
}>;


export type QueueEnqueueReprocessTorrentsBatchMutation = { __typename?: 'Mutation', queue: { __typename?: 'QueueMutation', enqueueReprocessTorrentsBatch?: void | null } };

export type QueuePurgeJobsMutationVariables = Exact<{
  input: QueuePurgeJobsInput;
}>;


export type QueuePurgeJobsMutation = { __typename?: 'Mutation', queue: { __typename?: 'QueueMutation', purgeJobs?: void | null } };

export type RegisterMutationVariables = Exact<{
  input: RegisterInput;
}>;


export type RegisterMutation = { __typename?: 'Mutation', self: { __typename?: 'SelfMutation', register: { __typename?: 'RegisterResult', user: { __typename?: 'User', id: number, username: string, role: string, lastLoginAt?: string | null, createdAt: string, updatedAt: string } } } };

export type SendTorrentMutationVariables = Exact<{
  target: Scalars['Ref']['input'];
  index?: InputMaybe<Scalars['Ref']['input']>;
  infoHashes: Array<Scalars['Hash20']['input']> | Scalars['Hash20']['input'];
  data: Scalars['JSON']['input'];
}>;


export type SendTorrentMutation = { __typename?: 'Mutation', target: { __typename?: 'TargetMutation', send: { __typename?: 'SendResult', data?: unknown | null, missingInfoHashes: Array<string> } } };

export type TorrentDeleteMutationVariables = Exact<{
  infoHashes: Array<Scalars['Hash20']['input']> | Scalars['Hash20']['input'];
}>;


export type TorrentDeleteMutation = { __typename?: 'Mutation', torrent: { __typename?: 'TorrentMutation', delete?: void | null } };

export type TorrentDeleteTagsMutationVariables = Exact<{
  infoHashes?: InputMaybe<Array<Scalars['Hash20']['input']> | Scalars['Hash20']['input']>;
  tagNames?: InputMaybe<Array<Scalars['String']['input']> | Scalars['String']['input']>;
}>;


export type TorrentDeleteTagsMutation = { __typename?: 'Mutation', torrent: { __typename?: 'TorrentMutation', deleteTags?: void | null } };

export type TorrentPutTagsMutationVariables = Exact<{
  infoHashes: Array<Scalars['Hash20']['input']> | Scalars['Hash20']['input'];
  tagNames: Array<Scalars['String']['input']> | Scalars['String']['input'];
}>;


export type TorrentPutTagsMutation = { __typename?: 'Mutation', torrent: { __typename?: 'TorrentMutation', putTags?: void | null } };

export type TorrentReprocessMutationVariables = Exact<{
  input: TorrentReprocessInput;
}>;


export type TorrentReprocessMutation = { __typename?: 'Mutation', torrent: { __typename?: 'TorrentMutation', reprocess?: void | null } };

export type TorrentSetTagsMutationVariables = Exact<{
  infoHashes: Array<Scalars['Hash20']['input']> | Scalars['Hash20']['input'];
  tagNames: Array<Scalars['String']['input']> | Scalars['String']['input'];
}>;


export type TorrentSetTagsMutation = { __typename?: 'Mutation', torrent: { __typename?: 'TorrentMutation', setTags?: void | null } };

export type WorkersRestartMutationVariables = Exact<{
  refs: Array<Scalars['Ref']['input']> | Scalars['Ref']['input'];
}>;


export type WorkersRestartMutation = { __typename?: 'Mutation', worker: { __typename?: 'WorkerMutation', restart: { __typename?: 'WorkerListAllQueryResult', workers: Array<{ __typename?: 'Worker', ref: string, label: string, state: WorkerState, error?: string | null, requiredBy: Array<string>, dependsOn: Array<string> }> } } };

export type WorkersShutdownMutationVariables = Exact<{
  refs: Array<Scalars['Ref']['input']> | Scalars['Ref']['input'];
}>;


export type WorkersShutdownMutation = { __typename?: 'Mutation', worker: { __typename?: 'WorkerMutation', shutdown: { __typename?: 'WorkerListAllQueryResult', workers: Array<{ __typename?: 'Worker', ref: string, label: string, state: WorkerState, error?: string | null, requiredBy: Array<string>, dependsOn: Array<string> }> } } };

export type WorkersStartMutationVariables = Exact<{
  refs: Array<Scalars['Ref']['input']> | Scalars['Ref']['input'];
}>;


export type WorkersStartMutation = { __typename?: 'Mutation', worker: { __typename?: 'WorkerMutation', start: { __typename?: 'WorkerListAllQueryResult', workers: Array<{ __typename?: 'Worker', ref: string, label: string, state: WorkerState, error?: string | null, requiredBy: Array<string>, dependsOn: Array<string> }> } } };

export type ApiKeysQueryVariables = Exact<{ [key: string]: never; }>;


export type ApiKeysQuery = { __typename?: 'Query', self: { __typename?: 'SelfQuery', apiKeys: Array<{ __typename?: 'APIKey', id: number, userId: number, name: string, expiresAt?: string | null, createdAt: string }> } };

export type AuthObjectActionsQueryVariables = Exact<{ [key: string]: never; }>;


export type AuthObjectActionsQuery = { __typename?: 'Query', auth: { __typename?: 'AuthQuery', listObjectActions: Array<{ __typename?: 'AuthObjectAction', namespace: string, object: string, action: string }> } };

export type ConfigParamsQueryVariables = Exact<{ [key: string]: never; }>;


export type ConfigParamsQuery = { __typename?: 'Query', config: { __typename?: 'ConfigQuery', pending: boolean, params: Array<{ __typename?: 'ConfigParam', ref: string, plugin: string, description?: string | null, value: unknown, source: string, default: unknown, dynamic: boolean, pending: boolean, jsonSchema: unknown }> } };

export type HealthCheckQueryVariables = Exact<{ [key: string]: never; }>;


export type HealthCheckQuery = { __typename?: 'Query', health: { __typename?: 'HealthQuery', status: HealthStatus, checks: Array<{ __typename?: 'HealthCheck', key: string, status: HealthStatus, timestamp?: string | null, error?: string | null }> } };

export type IdentityQueryVariables = Exact<{ [key: string]: never; }>;


export type IdentityQuery = { __typename?: 'Query', self: { __typename?: 'SelfQuery', identity: { __typename?: 'Self', user?: { __typename?: 'User', id: number, username: string, role: string, lastLoginAt?: string | null, createdAt: string, updatedAt: string } | null, apiKey?: { __typename?: 'APIKey', id: number, userId: number, name: string, expiresAt?: string | null, createdAt: string } | null, permissions: Array<{ __typename?: 'AuthObjectAction', namespace: string, object: string, action: string }> } } };

export type IndexesQueryVariables = Exact<{ [key: string]: never; }>;


export type IndexesQuery = { __typename?: 'Query', index: { __typename?: 'IndexQuery', default: string, infos: Array<{ __typename?: 'SearchIndexInfo', ref: string, name: string, resultTypes: Array<SearchResultType> }> } };

export type InvitationsQueryVariables = Exact<{
  input?: InputMaybe<ListInvitationsInput>;
}>;


export type InvitationsQuery = { __typename?: 'Query', auth: { __typename?: 'AuthQuery', listInvitations: { __typename?: 'ListInvitationsResult', totalCount: number, invitations: Array<{ __typename?: 'Invitation', code: string, role: string, email?: string | null, expiresAt?: string | null, createdAt: string, createdBy?: { __typename?: 'User', id: number, username: string, role: string, lastLoginAt?: string | null, createdAt: string, updatedAt: string } | null, claimedBy?: { __typename?: 'User', id: number, username: string, role: string, lastLoginAt?: string | null, createdAt: string, updatedAt: string } | null }> } } };

export type PasswordEntropyQueryVariables = Exact<{
  password: Scalars['String']['input'];
}>;


export type PasswordEntropyQuery = { __typename?: 'Query', self: { __typename?: 'SelfQuery', passwordEntropy: { __typename?: 'PasswordEntropyResult', entropy: number, minEntropy: number, valid: boolean } } };

export type PluginsQueryVariables = Exact<{ [key: string]: never; }>;


export type PluginsQuery = { __typename?: 'Query', plugin: { __typename?: 'PluginQuery', list: Array<{ __typename?: 'PluginInfo', ref: string, description?: string | null, enabled: boolean, dependsOn: Array<string>, requiredBy: Array<string> }> } };

export type QueueJobsQueryVariables = Exact<{
  input: QueueJobsQueryInput;
}>;


export type QueueJobsQuery = { __typename?: 'Query', queue: { __typename?: 'QueueQuery', jobs: { __typename?: 'QueueJobsQueryResult', totalCount: number, hasNextPage?: boolean | null, items: Array<{ __typename?: 'QueueJob', id: string, queue: string, status: QueueJobStatus, payload: string, priority: number, retries: number, maxRetries: number, runAfter: string, ranAt?: string | null, error?: string | null, createdAt: string }>, aggregations: { __typename?: 'QueueJobsAggregations', queue?: Array<{ __typename?: 'QueueJobQueueAgg', value: string, label: string, count: number }> | null, status?: Array<{ __typename?: 'QueueJobStatusAgg', value: QueueJobStatus, label: string, count: number }> | null } } } };

export type QueueMetricsQueryVariables = Exact<{
  input: QueueMetricsQueryInput;
}>;


export type QueueMetricsQuery = { __typename?: 'Query', queue: { __typename?: 'QueueQuery', metrics: { __typename?: 'QueueMetricsQueryResult', buckets: Array<{ __typename?: 'QueueMetricsBucket', queue: string, status: QueueJobStatus, createdAtBucket: string, ranAtBucket?: string | null, count: number, latency?: string | null }> } } };

export type RolesQueryVariables = Exact<{ [key: string]: never; }>;


export type RolesQuery = { __typename?: 'Query', auth: { __typename?: 'AuthQuery', listRoles: Array<{ __typename?: 'Role', name: string, core: boolean, permissions: Array<{ __typename?: 'Permission', core: boolean, subject: { __typename?: 'AuthSubject', type: AuthSubjectType, name: string }, objectAction: { __typename?: 'AuthObjectAction', namespace: string, object: string, action: string } }> }> } };

export type TargetsQueryVariables = Exact<{ [key: string]: never; }>;


export type TargetsQuery = { __typename?: 'Query', target: { __typename?: 'TargetQuery', default?: string | null, targets: Array<{ __typename?: 'Target', ref: string, name: string, dataSchema?: unknown | null, uiSchema?: unknown | null }> } };

export type TorrentContentSearchQueryVariables = Exact<{
  input: SearchInput;
}>;


export type TorrentContentSearchQuery = { __typename?: 'Query', torrent: { __typename?: 'TorrentQuery', searchTorrentContent: { __typename?: 'TorrentContentSearchResult', totalCount?: number | null, totalCountIsEstimate: boolean, hasNextPage?: boolean | null, items: Array<{ __typename?: 'TorrentContent', id: string, infoHash: string, contentType?: ContentType | null, title: string, video3d?: Video3D | null, videoCodec?: VideoCodec | null, videoModifier?: VideoModifier | null, videoResolution?: VideoResolution | null, videoSource?: VideoSource | null, seeders?: number | null, leechers?: number | null, publishedAt: string, createdAt: string, updatedAt: string, torrent: { __typename?: 'Torrent', infoHash: string, name: string, size: number, filesStatus: FilesStatus, filesCount?: number | null, hasFilesInfo: boolean, singleFile?: boolean | null, fileType?: FileType | null, seeders?: number | null, leechers?: number | null, tagNames: Array<string>, magnetUri: string, createdAt: string, updatedAt: string, sources: Array<{ __typename?: 'TorrentSourceInfo', key: string, name: string }> }, content?: { __typename?: 'Content', type: ContentType, source: string, id: string, title: string, releaseDate?: string | null, releaseYear?: number | null, overview?: string | null, runtime?: number | null, voteAverage?: number | null, voteCount?: number | null, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string }, originalLanguage?: { __typename?: 'LanguageInfo', id: string, name: string } | null, attributes: Array<{ __typename?: 'ContentAttribute', source: string, key: string, value: string, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }>, collections: Array<{ __typename?: 'ContentCollection', type: string, source: string, id: string, name: string, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }>, externalLinks: Array<{ __typename?: 'ExternalLink', url: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }> } | null, languages?: Array<{ __typename?: 'LanguageInfo', id: string, name: string }> | null, episodes?: { __typename?: 'Episodes', label: string, seasons: Array<{ __typename?: 'Season', season: number, episodes?: Array<number> | null }> } | null }>, facets: Array<{ __typename?: 'FacetResult', key: FacetKey, label: string, logic: FacetLogic, items: Array<{ __typename?: 'FacetResultItem', value: string, label: string, count: number, isEstimate: boolean }> }> } } };

export type TorrentFilesSearchQueryVariables = Exact<{
  input: SearchInput;
}>;


export type TorrentFilesSearchQuery = { __typename?: 'Query', torrent: { __typename?: 'TorrentQuery', searchTorrentFiles: { __typename?: 'TorrentFilesSearchResult', totalCount?: number | null, hasNextPage?: boolean | null, items: Array<{ __typename?: 'TorrentFile', infoHash: string, index: number, path: string, size: number, fileType?: FileType | null, createdAt: string, updatedAt: string }> } } };

export type TorrentMetricsQueryVariables = Exact<{
  input: TorrentMetricsQueryInput;
}>;


export type TorrentMetricsQuery = { __typename?: 'Query', torrent: { __typename?: 'TorrentQuery', metrics: { __typename?: 'TorrentMetricsQueryResult', buckets: Array<{ __typename?: 'TorrentMetricsBucket', source: string, updated: boolean, bucket: string, count: number }> }, listSources: { __typename?: 'TorrentListSourcesResult', sources: Array<{ __typename?: 'TorrentSource', key: string, name: string }> } } };

export type TorrentSuggestTagsQueryVariables = Exact<{
  input: SuggestTagsQueryInput;
}>;


export type TorrentSuggestTagsQuery = { __typename?: 'Query', torrent: { __typename?: 'TorrentQuery', suggestTags: { __typename?: 'TorrentSuggestTagsResult', suggestions: Array<{ __typename?: 'SuggestedTag', name: string, count: number }> } } };

export type UsersQueryVariables = Exact<{
  input: ListUsersInput;
}>;


export type UsersQuery = { __typename?: 'Query', auth: { __typename?: 'AuthQuery', listUsers: { __typename?: 'ListUsersResult', totalCount: number, users: Array<{ __typename?: 'User', id: number, username: string, role: string, lastLoginAt?: string | null, createdAt: string, updatedAt: string }> } } };

export type VersionQueryVariables = Exact<{ [key: string]: never; }>;


export type VersionQuery = { __typename?: 'Query', version: string };

export type WorkersQueryVariables = Exact<{ [key: string]: never; }>;


export type WorkersQuery = { __typename?: 'Query', worker: { __typename?: 'WorkerQuery', listAll: { __typename?: 'WorkerListAllQueryResult', workers: Array<{ __typename?: 'Worker', ref: string, label: string, state: WorkerState, error?: string | null, requiredBy: Array<string>, dependsOn: Array<string> }> } } };

export const ConfigParamFragmentDoc = gql`
    fragment ConfigParam on ConfigParam {
  ref
  plugin
  description
  value
  source
  default
  dynamic
  pending
  jsonSchema
}
    `;
export const UserFragmentDoc = gql`
    fragment User on User {
  id
  username
  role
  lastLoginAt
  createdAt
  updatedAt
}
    `;
export const InvitationFragmentDoc = gql`
    fragment Invitation on Invitation {
  code
  role
  email
  createdBy {
    ...User
  }
  claimedBy {
    ...User
  }
  expiresAt
  createdAt
}
    ${UserFragmentDoc}`;
export const PasswordEntropyResultFragmentDoc = gql`
    fragment PasswordEntropyResult on PasswordEntropyResult {
  entropy
  minEntropy
  valid
}
    `;
export const PluginInfoFragmentDoc = gql`
    fragment PluginInfo on PluginInfo {
  ref
  description
  enabled
  dependsOn
  requiredBy
}
    `;
export const QueueJobFragmentDoc = gql`
    fragment QueueJob on QueueJob {
  id
  queue
  status
  payload
  priority
  retries
  maxRetries
  runAfter
  ranAt
  error
  createdAt
}
    `;
export const QueueJobsQueryResultFragmentDoc = gql`
    fragment QueueJobsQueryResult on QueueJobsQueryResult {
  items {
    ...QueueJob
  }
  totalCount
  hasNextPage
  aggregations {
    queue {
      value
      label
      count
    }
    status {
      value
      label
      count
    }
  }
}
    ${QueueJobFragmentDoc}`;
export const AuthObjectActionFragmentDoc = gql`
    fragment AuthObjectAction on AuthObjectAction {
  namespace
  object
  action
}
    `;
export const PermissionFragmentDoc = gql`
    fragment Permission on Permission {
  subject {
    type
    name
  }
  objectAction {
    ...AuthObjectAction
  }
  core
}
    ${AuthObjectActionFragmentDoc}`;
export const RoleFragmentDoc = gql`
    fragment Role on Role {
  name
  core
  permissions {
    ...Permission
  }
}
    ${PermissionFragmentDoc}`;
export const SearchIndexInfoFragmentDoc = gql`
    fragment SearchIndexInfo on SearchIndexInfo {
  ref
  name
  resultTypes
}
    `;
export const ApiKeyFragmentDoc = gql`
    fragment APIKey on APIKey {
  id
  userId
  name
  expiresAt
  createdAt
}
    `;
export const SelfFragmentDoc = gql`
    fragment Self on Self {
  user {
    ...User
  }
  apiKey {
    ...APIKey
  }
  permissions {
    ...AuthObjectAction
  }
}
    ${UserFragmentDoc}
${ApiKeyFragmentDoc}
${AuthObjectActionFragmentDoc}`;
export const TargetFragmentDoc = gql`
    fragment Target on Target {
  ref
  name
  dataSchema
  uiSchema
}
    `;
export const TorrentFragmentDoc = gql`
    fragment Torrent on Torrent {
  infoHash
  name
  size
  filesStatus
  filesCount
  hasFilesInfo
  singleFile
  fileType
  sources {
    key
    name
  }
  seeders
  leechers
  tagNames
  magnetUri
  createdAt
  updatedAt
}
    `;
export const ContentFragmentDoc = gql`
    fragment Content on Content {
  type
  source
  id
  metadataSource {
    key
    name
  }
  title
  releaseDate
  releaseYear
  overview
  runtime
  voteAverage
  voteCount
  originalLanguage {
    id
    name
  }
  attributes {
    metadataSource {
      key
      name
    }
    source
    key
    value
    createdAt
    updatedAt
  }
  collections {
    metadataSource {
      key
      name
    }
    type
    source
    id
    name
    createdAt
    updatedAt
  }
  externalLinks {
    metadataSource {
      key
      name
    }
    url
  }
  createdAt
  updatedAt
}
    `;
export const TorrentContentFragmentDoc = gql`
    fragment TorrentContent on TorrentContent {
  id
  infoHash
  contentType
  title
  torrent {
    ...Torrent
  }
  content {
    ...Content
  }
  languages {
    id
    name
  }
  episodes {
    label
    seasons {
      season
      episodes
    }
  }
  video3d
  videoCodec
  videoModifier
  videoResolution
  videoSource
  seeders
  leechers
  publishedAt
  createdAt
  updatedAt
}
    ${TorrentFragmentDoc}
${ContentFragmentDoc}`;
export const FacetResultItemFragmentDoc = gql`
    fragment FacetResultItem on FacetResultItem {
  value
  label
  count
  isEstimate
}
    `;
export const FacetResultFragmentDoc = gql`
    fragment FacetResult on FacetResult {
  key
  label
  logic
  items {
    ...FacetResultItem
  }
}
    ${FacetResultItemFragmentDoc}`;
export const TorrentContentSearchResultFragmentDoc = gql`
    fragment TorrentContentSearchResult on TorrentContentSearchResult {
  items {
    ...TorrentContent
  }
  totalCount
  totalCountIsEstimate
  hasNextPage
  facets {
    ...FacetResult
  }
}
    ${TorrentContentFragmentDoc}
${FacetResultFragmentDoc}`;
export const TorrentFileFragmentDoc = gql`
    fragment TorrentFile on TorrentFile {
  infoHash
  index
  path
  size
  fileType
  createdAt
  updatedAt
}
    `;
export const TorrentFilesSearchResultFragmentDoc = gql`
    fragment TorrentFilesSearchResult on TorrentFilesSearchResult {
  items {
    ...TorrentFile
  }
  totalCount
  hasNextPage
}
    ${TorrentFileFragmentDoc}`;
export const WorkerFragmentDoc = gql`
    fragment Worker on Worker {
  ref
  label
  state
  error
  requiredBy
  dependsOn
}
    `;
export const ConfigDeleteDocument = gql`
    mutation ConfigDelete($ref: Ref!) {
  config {
    delete(ref: $ref) {
      ...ConfigParam
    }
  }
}
    ${ConfigParamFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class ConfigDeleteGQL extends Apollo.Mutation<ConfigDeleteMutation, ConfigDeleteMutationVariables> {
    override document = ConfigDeleteDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const ConfigSaveDocument = gql`
    mutation ConfigSave($ref: Ref!, $value: JSON!) {
  config {
    save(ref: $ref, value: $value) {
      ...ConfigParam
    }
  }
}
    ${ConfigParamFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class ConfigSaveGQL extends Apollo.Mutation<ConfigSaveMutation, ConfigSaveMutationVariables> {
    override document = ConfigSaveDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const CreateApiKeyDocument = gql`
    mutation CreateAPIKey($input: CreateAPIKeyInput!) {
  self {
    createAPIKey(input: $input) {
      id
      apiKey
      name
      expiresAt
    }
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class CreateApiKeyGQL extends Apollo.Mutation<CreateApiKeyMutation, CreateApiKeyMutationVariables> {
    override document = CreateApiKeyDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const InvitationDeleteDocument = gql`
    mutation InvitationDelete($code: String!) {
  auth {
    deleteInvitation(code: $code)
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class InvitationDeleteGQL extends Apollo.Mutation<InvitationDeleteMutation, InvitationDeleteMutationVariables> {
    override document = InvitationDeleteDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const InviteDocument = gql`
    mutation Invite($input: InviteInput!) {
  auth {
    invite(input: $input) {
      ...Invitation
    }
  }
}
    ${InvitationFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class InviteGQL extends Apollo.Mutation<InviteMutation, InviteMutationVariables> {
    override document = InviteDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const LoginDocument = gql`
    mutation Login($username: String!, $password: String!) {
  self {
    login(username: $username, password: $password) {
      token
      user {
        ...User
      }
      permissions {
        ...Permission
      }
    }
  }
}
    ${UserFragmentDoc}
${PermissionFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class LoginGQL extends Apollo.Mutation<LoginMutation, LoginMutationVariables> {
    override document = LoginDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const PutRoleDocument = gql`
    mutation PutRole($role: String!, $objectActions: [AuthObjectActionInput!]!) {
  auth {
    putRole(role: $role, objectActions: $objectActions) {
      ...Role
    }
  }
}
    ${RoleFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class PutRoleGQL extends Apollo.Mutation<PutRoleMutation, PutRoleMutationVariables> {
    override document = PutRoleDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const QueueEnqueueReprocessTorrentsBatchDocument = gql`
    mutation QueueEnqueueReprocessTorrentsBatch($input: QueueEnqueueReprocessTorrentsBatchInput!) {
  queue {
    enqueueReprocessTorrentsBatch(input: $input)
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class QueueEnqueueReprocessTorrentsBatchGQL extends Apollo.Mutation<QueueEnqueueReprocessTorrentsBatchMutation, QueueEnqueueReprocessTorrentsBatchMutationVariables> {
    override document = QueueEnqueueReprocessTorrentsBatchDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const QueuePurgeJobsDocument = gql`
    mutation QueuePurgeJobs($input: QueuePurgeJobsInput!) {
  queue {
    purgeJobs(input: $input)
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class QueuePurgeJobsGQL extends Apollo.Mutation<QueuePurgeJobsMutation, QueuePurgeJobsMutationVariables> {
    override document = QueuePurgeJobsDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const RegisterDocument = gql`
    mutation Register($input: RegisterInput!) {
  self {
    register(input: $input) {
      user {
        ...User
      }
    }
  }
}
    ${UserFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class RegisterGQL extends Apollo.Mutation<RegisterMutation, RegisterMutationVariables> {
    override document = RegisterDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const SendTorrentDocument = gql`
    mutation SendTorrent($target: Ref!, $index: Ref, $infoHashes: [Hash20!]!, $data: JSON!) {
  target {
    send(
      input: {target: $target, index: $index, infoHashes: $infoHashes, data: $data}
    ) {
      data
      missingInfoHashes
    }
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class SendTorrentGQL extends Apollo.Mutation<SendTorrentMutation, SendTorrentMutationVariables> {
    override document = SendTorrentDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const TorrentDeleteDocument = gql`
    mutation TorrentDelete($infoHashes: [Hash20!]!) {
  torrent {
    delete(infoHashes: $infoHashes)
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class TorrentDeleteGQL extends Apollo.Mutation<TorrentDeleteMutation, TorrentDeleteMutationVariables> {
    override document = TorrentDeleteDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const TorrentDeleteTagsDocument = gql`
    mutation TorrentDeleteTags($infoHashes: [Hash20!], $tagNames: [String!]) {
  torrent {
    deleteTags(infoHashes: $infoHashes, tagNames: $tagNames)
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class TorrentDeleteTagsGQL extends Apollo.Mutation<TorrentDeleteTagsMutation, TorrentDeleteTagsMutationVariables> {
    override document = TorrentDeleteTagsDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const TorrentPutTagsDocument = gql`
    mutation TorrentPutTags($infoHashes: [Hash20!]!, $tagNames: [String!]!) {
  torrent {
    putTags(infoHashes: $infoHashes, tagNames: $tagNames)
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class TorrentPutTagsGQL extends Apollo.Mutation<TorrentPutTagsMutation, TorrentPutTagsMutationVariables> {
    override document = TorrentPutTagsDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const TorrentReprocessDocument = gql`
    mutation TorrentReprocess($input: TorrentReprocessInput!) {
  torrent {
    reprocess(input: $input)
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class TorrentReprocessGQL extends Apollo.Mutation<TorrentReprocessMutation, TorrentReprocessMutationVariables> {
    override document = TorrentReprocessDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const TorrentSetTagsDocument = gql`
    mutation TorrentSetTags($infoHashes: [Hash20!]!, $tagNames: [String!]!) {
  torrent {
    setTags(infoHashes: $infoHashes, tagNames: $tagNames)
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class TorrentSetTagsGQL extends Apollo.Mutation<TorrentSetTagsMutation, TorrentSetTagsMutationVariables> {
    override document = TorrentSetTagsDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const WorkersRestartDocument = gql`
    mutation WorkersRestart($refs: [Ref!]!) {
  worker {
    restart(refs: $refs) {
      workers {
        ...Worker
      }
    }
  }
}
    ${WorkerFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class WorkersRestartGQL extends Apollo.Mutation<WorkersRestartMutation, WorkersRestartMutationVariables> {
    override document = WorkersRestartDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const WorkersShutdownDocument = gql`
    mutation WorkersShutdown($refs: [Ref!]!) {
  worker {
    shutdown(refs: $refs) {
      workers {
        ...Worker
      }
    }
  }
}
    ${WorkerFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class WorkersShutdownGQL extends Apollo.Mutation<WorkersShutdownMutation, WorkersShutdownMutationVariables> {
    override document = WorkersShutdownDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const WorkersStartDocument = gql`
    mutation WorkersStart($refs: [Ref!]!) {
  worker {
    start(refs: $refs) {
      workers {
        ...Worker
      }
    }
  }
}
    ${WorkerFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class WorkersStartGQL extends Apollo.Mutation<WorkersStartMutation, WorkersStartMutationVariables> {
    override document = WorkersStartDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const ApiKeysDocument = gql`
    query APIKeys {
  self {
    apiKeys {
      ...APIKey
    }
  }
}
    ${ApiKeyFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class ApiKeysGQL extends Apollo.Query<ApiKeysQuery, ApiKeysQueryVariables> {
    override document = ApiKeysDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const AuthObjectActionsDocument = gql`
    query AuthObjectActions {
  auth {
    listObjectActions {
      ...AuthObjectAction
    }
  }
}
    ${AuthObjectActionFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class AuthObjectActionsGQL extends Apollo.Query<AuthObjectActionsQuery, AuthObjectActionsQueryVariables> {
    override document = AuthObjectActionsDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const ConfigParamsDocument = gql`
    query ConfigParams {
  config {
    pending
    params {
      ...ConfigParam
    }
  }
}
    ${ConfigParamFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class ConfigParamsGQL extends Apollo.Query<ConfigParamsQuery, ConfigParamsQueryVariables> {
    override document = ConfigParamsDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const HealthCheckDocument = gql`
    query HealthCheck {
  health {
    status
    checks {
      key
      status
      timestamp
      error
    }
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class HealthCheckGQL extends Apollo.Query<HealthCheckQuery, HealthCheckQueryVariables> {
    override document = HealthCheckDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const IdentityDocument = gql`
    query Identity {
  self {
    identity {
      ...Self
    }
  }
}
    ${SelfFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class IdentityGQL extends Apollo.Query<IdentityQuery, IdentityQueryVariables> {
    override document = IdentityDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const IndexesDocument = gql`
    query Indexes {
  index {
    default
    infos {
      ...SearchIndexInfo
    }
  }
}
    ${SearchIndexInfoFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class IndexesGQL extends Apollo.Query<IndexesQuery, IndexesQueryVariables> {
    override document = IndexesDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const InvitationsDocument = gql`
    query Invitations($input: ListInvitationsInput) {
  auth {
    listInvitations(input: $input) {
      invitations {
        ...Invitation
      }
      totalCount
    }
  }
}
    ${InvitationFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class InvitationsGQL extends Apollo.Query<InvitationsQuery, InvitationsQueryVariables> {
    override document = InvitationsDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const PasswordEntropyDocument = gql`
    query PasswordEntropy($password: String!) {
  self {
    passwordEntropy(password: $password) {
      ...PasswordEntropyResult
    }
  }
}
    ${PasswordEntropyResultFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class PasswordEntropyGQL extends Apollo.Query<PasswordEntropyQuery, PasswordEntropyQueryVariables> {
    override document = PasswordEntropyDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const PluginsDocument = gql`
    query Plugins {
  plugin {
    list {
      ...PluginInfo
    }
  }
}
    ${PluginInfoFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class PluginsGQL extends Apollo.Query<PluginsQuery, PluginsQueryVariables> {
    override document = PluginsDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const QueueJobsDocument = gql`
    query QueueJobs($input: QueueJobsQueryInput!) {
  queue {
    jobs(input: $input) {
      ...QueueJobsQueryResult
    }
  }
}
    ${QueueJobsQueryResultFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class QueueJobsGQL extends Apollo.Query<QueueJobsQuery, QueueJobsQueryVariables> {
    override document = QueueJobsDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const QueueMetricsDocument = gql`
    query QueueMetrics($input: QueueMetricsQueryInput!) {
  queue {
    metrics(input: $input) {
      buckets {
        queue
        status
        createdAtBucket
        ranAtBucket
        count
        latency
      }
    }
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class QueueMetricsGQL extends Apollo.Query<QueueMetricsQuery, QueueMetricsQueryVariables> {
    override document = QueueMetricsDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const RolesDocument = gql`
    query Roles {
  auth {
    listRoles {
      ...Role
    }
  }
}
    ${RoleFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class RolesGQL extends Apollo.Query<RolesQuery, RolesQueryVariables> {
    override document = RolesDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const TargetsDocument = gql`
    query Targets {
  target {
    default
    targets {
      ...Target
    }
  }
}
    ${TargetFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class TargetsGQL extends Apollo.Query<TargetsQuery, TargetsQueryVariables> {
    override document = TargetsDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const TorrentContentSearchDocument = gql`
    query TorrentContentSearch($input: SearchInput!) {
  torrent {
    searchTorrentContent(input: $input) {
      ...TorrentContentSearchResult
    }
  }
}
    ${TorrentContentSearchResultFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class TorrentContentSearchGQL extends Apollo.Query<TorrentContentSearchQuery, TorrentContentSearchQueryVariables> {
    override document = TorrentContentSearchDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const TorrentFilesSearchDocument = gql`
    query TorrentFilesSearch($input: SearchInput!) {
  torrent {
    searchTorrentFiles(input: $input) {
      ...TorrentFilesSearchResult
    }
  }
}
    ${TorrentFilesSearchResultFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class TorrentFilesSearchGQL extends Apollo.Query<TorrentFilesSearchQuery, TorrentFilesSearchQueryVariables> {
    override document = TorrentFilesSearchDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const TorrentMetricsDocument = gql`
    query TorrentMetrics($input: TorrentMetricsQueryInput!) {
  torrent {
    metrics(input: $input) {
      buckets {
        source
        updated
        bucket
        count
      }
    }
    listSources {
      sources {
        key
        name
      }
    }
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class TorrentMetricsGQL extends Apollo.Query<TorrentMetricsQuery, TorrentMetricsQueryVariables> {
    override document = TorrentMetricsDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const TorrentSuggestTagsDocument = gql`
    query TorrentSuggestTags($input: SuggestTagsQueryInput!) {
  torrent {
    suggestTags(input: $input) {
      suggestions {
        name
        count
      }
    }
  }
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class TorrentSuggestTagsGQL extends Apollo.Query<TorrentSuggestTagsQuery, TorrentSuggestTagsQueryVariables> {
    override document = TorrentSuggestTagsDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const UsersDocument = gql`
    query Users($input: ListUsersInput!) {
  auth {
    listUsers(input: $input) {
      users {
        ...User
      }
      totalCount
    }
  }
}
    ${UserFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class UsersGQL extends Apollo.Query<UsersQuery, UsersQueryVariables> {
    override document = UsersDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const VersionDocument = gql`
    query Version {
  version
}
    `;

  @Injectable({
    providedIn: 'root'
  })
  export class VersionGQL extends Apollo.Query<VersionQuery, VersionQueryVariables> {
    override document = VersionDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }
export const WorkersDocument = gql`
    query Workers {
  worker {
    listAll {
      workers {
        ...Worker
      }
    }
  }
}
    ${WorkerFragmentDoc}`;

  @Injectable({
    providedIn: 'root'
  })
  export class WorkersGQL extends Apollo.Query<WorkersQuery, WorkersQueryVariables> {
    override document = WorkersDocument;
    
    constructor(apollo: Apollo.Apollo) {
      super(apollo);
    }
  }