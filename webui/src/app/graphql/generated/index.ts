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
  Hash20: { input: string; output: string; }
  Void: { input: void; output: void; }
  Year: { input: number; output: number; }
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
  | 'book'
  | 'game'
  | 'movie'
  | 'music'
  | 'software'
  | 'tv_show'
  | 'xxx';

export type ContentTypeAgg = {
  __typename?: 'ContentTypeAgg';
  count: Scalars['Int']['output'];
  label: Scalars['String']['output'];
  value?: Maybe<ContentType>;
};

export type ContentTypeFacetInput = {
  aggregate?: InputMaybe<Scalars['Boolean']['input']>;
  filter?: InputMaybe<Array<InputMaybe<ContentType>>>;
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

export type FacetLogic =
  | 'and'
  | 'or';

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

export type GenreAgg = {
  __typename?: 'GenreAgg';
  count: Scalars['Int']['output'];
  label: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export type GenreFacetInput = {
  aggregate?: InputMaybe<Scalars['Boolean']['input']>;
  filter?: InputMaybe<Array<Scalars['String']['input']>>;
  logic?: InputMaybe<FacetLogic>;
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

export type LanguageAgg = {
  __typename?: 'LanguageAgg';
  count: Scalars['Int']['output'];
  label: Scalars['String']['output'];
  value: Language;
};

export type LanguageFacetInput = {
  aggregate?: InputMaybe<Scalars['Boolean']['input']>;
  filter?: InputMaybe<Array<Language>>;
};

export type LanguageInfo = {
  __typename?: 'LanguageInfo';
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type MetadataSource = {
  __typename?: 'MetadataSource';
  key: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type Mutation = {
  __typename?: 'Mutation';
  torrent: TorrentMutation;
};

export type Query = {
  __typename?: 'Query';
  torrent: TorrentQuery;
  torrentContent: TorrentContentQuery;
};

export type ReleaseYearAgg = {
  __typename?: 'ReleaseYearAgg';
  count: Scalars['Int']['output'];
  label: Scalars['String']['output'];
  value?: Maybe<Scalars['Year']['output']>;
};

export type ReleaseYearFacetInput = {
  aggregate?: InputMaybe<Scalars['Boolean']['input']>;
  filter?: InputMaybe<Array<InputMaybe<Scalars['Year']['input']>>>;
};

export type SearchQueryInput = {
  cached?: InputMaybe<Scalars['Boolean']['input']>;
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
  queryString?: InputMaybe<Scalars['String']['input']>;
  totalCount?: InputMaybe<Scalars['Boolean']['input']>;
};

export type Season = {
  __typename?: 'Season';
  episodes?: Maybe<Array<Scalars['Int']['output']>>;
  season: Scalars['Int']['output'];
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

export type Torrent = {
  __typename?: 'Torrent';
  createdAt: Scalars['DateTime']['output'];
  extension?: Maybe<Scalars['String']['output']>;
  fileType?: Maybe<FileType>;
  fileTypes?: Maybe<Array<FileType>>;
  files?: Maybe<Array<TorrentFile>>;
  filesStatus: FilesStatus;
  hasFilesInfo: Scalars['Boolean']['output'];
  infoHash: Scalars['Hash20']['output'];
  leechers?: Maybe<Scalars['Int']['output']>;
  magnetUri: Scalars['String']['output'];
  name: Scalars['String']['output'];
  private: Scalars['Boolean']['output'];
  seeders?: Maybe<Scalars['Int']['output']>;
  singleFile?: Maybe<Scalars['Boolean']['output']>;
  size: Scalars['Int']['output'];
  sources: Array<TorrentSource>;
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
  releaseDate?: Maybe<Scalars['Date']['output']>;
  releaseGroup?: Maybe<Scalars['String']['output']>;
  releaseYear?: Maybe<Scalars['Year']['output']>;
  title: Scalars['String']['output'];
  torrent: Torrent;
  updatedAt: Scalars['DateTime']['output'];
  video3d?: Maybe<Video3d>;
  videoCodec?: Maybe<VideoCodec>;
  videoModifier?: Maybe<VideoModifier>;
  videoResolution?: Maybe<VideoResolution>;
  videoSource?: Maybe<VideoSource>;
};

export type TorrentContentAggregations = {
  __typename?: 'TorrentContentAggregations';
  contentType?: Maybe<Array<ContentTypeAgg>>;
  genre?: Maybe<Array<GenreAgg>>;
  language?: Maybe<Array<LanguageAgg>>;
  releaseYear?: Maybe<Array<ReleaseYearAgg>>;
  torrentFileType?: Maybe<Array<TorrentFileTypeAgg>>;
  torrentSource?: Maybe<Array<TorrentSourceAgg>>;
  torrentTag?: Maybe<Array<TorrentTagAgg>>;
  videoResolution?: Maybe<Array<VideoResolutionAgg>>;
  videoSource?: Maybe<Array<VideoSourceAgg>>;
};

export type TorrentContentFacetsInput = {
  contentType?: InputMaybe<ContentTypeFacetInput>;
  genre?: InputMaybe<GenreFacetInput>;
  language?: InputMaybe<LanguageFacetInput>;
  releaseYear?: InputMaybe<ReleaseYearFacetInput>;
  torrentFileType?: InputMaybe<TorrentFileTypeFacetInput>;
  torrentSource?: InputMaybe<TorrentSourceFacetInput>;
  torrentTag?: InputMaybe<TorrentTagFacetInput>;
  videoResolution?: InputMaybe<VideoResolutionFacetInput>;
  videoSource?: InputMaybe<VideoSourceFacetInput>;
};

export type TorrentContentQuery = {
  __typename?: 'TorrentContentQuery';
  search: TorrentContentSearchResult;
};


export type TorrentContentQuerySearchArgs = {
  facets?: InputMaybe<TorrentContentFacetsInput>;
  query?: InputMaybe<SearchQueryInput>;
};

export type TorrentContentSearchResult = {
  __typename?: 'TorrentContentSearchResult';
  aggregations: TorrentContentAggregations;
  items: Array<TorrentContent>;
  totalCount: Scalars['Int']['output'];
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

export type TorrentFileTypeAgg = {
  __typename?: 'TorrentFileTypeAgg';
  count: Scalars['Int']['output'];
  label: Scalars['String']['output'];
  value: FileType;
};

export type TorrentFileTypeFacetInput = {
  aggregate?: InputMaybe<Scalars['Boolean']['input']>;
  filter?: InputMaybe<Array<FileType>>;
  logic?: InputMaybe<FacetLogic>;
};

export type TorrentMutation = {
  __typename?: 'TorrentMutation';
  delete?: Maybe<Scalars['Void']['output']>;
  deleteTags?: Maybe<Scalars['Void']['output']>;
  putTags?: Maybe<Scalars['Void']['output']>;
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


export type TorrentMutationSetTagsArgs = {
  infoHashes: Array<Scalars['Hash20']['input']>;
  tagNames: Array<Scalars['String']['input']>;
};

export type TorrentQuery = {
  __typename?: 'TorrentQuery';
  suggestTags: TorrentSuggestTagsResult;
};


export type TorrentQuerySuggestTagsArgs = {
  query?: InputMaybe<SuggestTagsQueryInput>;
};

export type TorrentSource = {
  __typename?: 'TorrentSource';
  importId?: Maybe<Scalars['String']['output']>;
  key: Scalars['String']['output'];
  leechers?: Maybe<Scalars['Int']['output']>;
  name: Scalars['String']['output'];
  seeders?: Maybe<Scalars['Int']['output']>;
};

export type TorrentSourceAgg = {
  __typename?: 'TorrentSourceAgg';
  count: Scalars['Int']['output'];
  label: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export type TorrentSourceFacetInput = {
  aggregate?: InputMaybe<Scalars['Boolean']['input']>;
  filter?: InputMaybe<Array<Scalars['String']['input']>>;
  logic?: InputMaybe<FacetLogic>;
};

export type TorrentSuggestTagsResult = {
  __typename?: 'TorrentSuggestTagsResult';
  suggestions: Array<SuggestedTag>;
};

export type TorrentTagAgg = {
  __typename?: 'TorrentTagAgg';
  count: Scalars['Int']['output'];
  label: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export type TorrentTagFacetInput = {
  aggregate?: InputMaybe<Scalars['Boolean']['input']>;
  filter?: InputMaybe<Array<Scalars['String']['input']>>;
  logic?: InputMaybe<FacetLogic>;
};

export type Video3d =
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

export type VideoResolutionAgg = {
  __typename?: 'VideoResolutionAgg';
  count: Scalars['Int']['output'];
  label: Scalars['String']['output'];
  value?: Maybe<VideoResolution>;
};

export type VideoResolutionFacetInput = {
  aggregate?: InputMaybe<Scalars['Boolean']['input']>;
  filter?: InputMaybe<Array<InputMaybe<VideoResolution>>>;
};

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

export type VideoSourceAgg = {
  __typename?: 'VideoSourceAgg';
  count: Scalars['Int']['output'];
  label: Scalars['String']['output'];
  value?: Maybe<VideoSource>;
};

export type VideoSourceFacetInput = {
  aggregate?: InputMaybe<Scalars['Boolean']['input']>;
  filter?: InputMaybe<Array<InputMaybe<VideoSource>>>;
};

export type ContentFragment = { __typename?: 'Content', type: ContentType, source: string, id: string, title: string, releaseDate?: string | null, releaseYear?: number | null, overview?: string | null, runtime?: number | null, voteAverage?: number | null, voteCount?: number | null, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string }, originalLanguage?: { __typename?: 'LanguageInfo', id: string, name: string } | null, attributes: Array<{ __typename?: 'ContentAttribute', source: string, key: string, value: string, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }>, collections: Array<{ __typename?: 'ContentCollection', type: string, source: string, id: string, name: string, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }>, externalLinks: Array<{ __typename?: 'ExternalLink', url: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }> };

export type TorrentFragment = { __typename?: 'Torrent', infoHash: string, name: string, size: number, private: boolean, filesStatus: FilesStatus, hasFilesInfo: boolean, singleFile?: boolean | null, fileType?: FileType | null, seeders?: number | null, leechers?: number | null, tagNames: Array<string>, magnetUri: string, createdAt: string, updatedAt: string, files?: Array<{ __typename?: 'TorrentFile', infoHash: string, index: number, path: string, size: number, fileType?: FileType | null, createdAt: string, updatedAt: string }> | null, sources: Array<{ __typename?: 'TorrentSource', key: string, name: string }> };

export type TorrentContentFragment = { __typename?: 'TorrentContent', id: string, infoHash: string, contentType?: ContentType | null, title: string, video3d?: Video3d | null, videoCodec?: VideoCodec | null, videoModifier?: VideoModifier | null, videoResolution?: VideoResolution | null, videoSource?: VideoSource | null, releaseDate?: string | null, releaseYear?: number | null, createdAt: string, updatedAt: string, torrent: { __typename?: 'Torrent', infoHash: string, name: string, size: number, private: boolean, filesStatus: FilesStatus, hasFilesInfo: boolean, singleFile?: boolean | null, fileType?: FileType | null, seeders?: number | null, leechers?: number | null, tagNames: Array<string>, magnetUri: string, createdAt: string, updatedAt: string, files?: Array<{ __typename?: 'TorrentFile', infoHash: string, index: number, path: string, size: number, fileType?: FileType | null, createdAt: string, updatedAt: string }> | null, sources: Array<{ __typename?: 'TorrentSource', key: string, name: string }> }, content?: { __typename?: 'Content', type: ContentType, source: string, id: string, title: string, releaseDate?: string | null, releaseYear?: number | null, overview?: string | null, runtime?: number | null, voteAverage?: number | null, voteCount?: number | null, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string }, originalLanguage?: { __typename?: 'LanguageInfo', id: string, name: string } | null, attributes: Array<{ __typename?: 'ContentAttribute', source: string, key: string, value: string, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }>, collections: Array<{ __typename?: 'ContentCollection', type: string, source: string, id: string, name: string, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }>, externalLinks: Array<{ __typename?: 'ExternalLink', url: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }> } | null, languages?: Array<{ __typename?: 'LanguageInfo', id: string, name: string }> | null, episodes?: { __typename?: 'Episodes', label: string, seasons: Array<{ __typename?: 'Season', season: number, episodes?: Array<number> | null }> } | null };

export type TorrentContentSearchResultFragment = { __typename?: 'TorrentContentSearchResult', totalCount: number, items: Array<{ __typename?: 'TorrentContent', id: string, infoHash: string, contentType?: ContentType | null, title: string, video3d?: Video3d | null, videoCodec?: VideoCodec | null, videoModifier?: VideoModifier | null, videoResolution?: VideoResolution | null, videoSource?: VideoSource | null, releaseDate?: string | null, releaseYear?: number | null, createdAt: string, updatedAt: string, torrent: { __typename?: 'Torrent', infoHash: string, name: string, size: number, private: boolean, filesStatus: FilesStatus, hasFilesInfo: boolean, singleFile?: boolean | null, fileType?: FileType | null, seeders?: number | null, leechers?: number | null, tagNames: Array<string>, magnetUri: string, createdAt: string, updatedAt: string, files?: Array<{ __typename?: 'TorrentFile', infoHash: string, index: number, path: string, size: number, fileType?: FileType | null, createdAt: string, updatedAt: string }> | null, sources: Array<{ __typename?: 'TorrentSource', key: string, name: string }> }, content?: { __typename?: 'Content', type: ContentType, source: string, id: string, title: string, releaseDate?: string | null, releaseYear?: number | null, overview?: string | null, runtime?: number | null, voteAverage?: number | null, voteCount?: number | null, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string }, originalLanguage?: { __typename?: 'LanguageInfo', id: string, name: string } | null, attributes: Array<{ __typename?: 'ContentAttribute', source: string, key: string, value: string, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }>, collections: Array<{ __typename?: 'ContentCollection', type: string, source: string, id: string, name: string, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }>, externalLinks: Array<{ __typename?: 'ExternalLink', url: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }> } | null, languages?: Array<{ __typename?: 'LanguageInfo', id: string, name: string }> | null, episodes?: { __typename?: 'Episodes', label: string, seasons: Array<{ __typename?: 'Season', season: number, episodes?: Array<number> | null }> } | null }>, aggregations: { __typename?: 'TorrentContentAggregations', contentType?: Array<{ __typename?: 'ContentTypeAgg', value?: ContentType | null, label: string, count: number }> | null, torrentSource?: Array<{ __typename?: 'TorrentSourceAgg', value: string, label: string, count: number }> | null, torrentTag?: Array<{ __typename?: 'TorrentTagAgg', value: string, label: string, count: number }> | null, torrentFileType?: Array<{ __typename?: 'TorrentFileTypeAgg', value: FileType, label: string, count: number }> | null, language?: Array<{ __typename?: 'LanguageAgg', value: Language, label: string, count: number }> | null, genre?: Array<{ __typename?: 'GenreAgg', value: string, label: string, count: number }> | null, releaseYear?: Array<{ __typename?: 'ReleaseYearAgg', value?: number | null, label: string, count: number }> | null, videoResolution?: Array<{ __typename?: 'VideoResolutionAgg', value?: VideoResolution | null, label: string, count: number }> | null, videoSource?: Array<{ __typename?: 'VideoSourceAgg', value?: VideoSource | null, label: string, count: number }> | null } };

export type TorrentFileFragment = { __typename?: 'TorrentFile', infoHash: string, index: number, path: string, size: number, fileType?: FileType | null, createdAt: string, updatedAt: string };

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

export type TorrentSetTagsMutationVariables = Exact<{
  infoHashes: Array<Scalars['Hash20']['input']> | Scalars['Hash20']['input'];
  tagNames: Array<Scalars['String']['input']> | Scalars['String']['input'];
}>;


export type TorrentSetTagsMutation = { __typename?: 'Mutation', torrent: { __typename?: 'TorrentMutation', setTags?: void | null } };

export type TorrentContentSearchQueryVariables = Exact<{
  query?: InputMaybe<SearchQueryInput>;
  facets?: InputMaybe<TorrentContentFacetsInput>;
}>;


export type TorrentContentSearchQuery = { __typename?: 'Query', torrentContent: { __typename?: 'TorrentContentQuery', search: { __typename?: 'TorrentContentSearchResult', totalCount: number, items: Array<{ __typename?: 'TorrentContent', id: string, infoHash: string, contentType?: ContentType | null, title: string, video3d?: Video3d | null, videoCodec?: VideoCodec | null, videoModifier?: VideoModifier | null, videoResolution?: VideoResolution | null, videoSource?: VideoSource | null, releaseDate?: string | null, releaseYear?: number | null, createdAt: string, updatedAt: string, torrent: { __typename?: 'Torrent', infoHash: string, name: string, size: number, private: boolean, filesStatus: FilesStatus, hasFilesInfo: boolean, singleFile?: boolean | null, fileType?: FileType | null, seeders?: number | null, leechers?: number | null, tagNames: Array<string>, magnetUri: string, createdAt: string, updatedAt: string, files?: Array<{ __typename?: 'TorrentFile', infoHash: string, index: number, path: string, size: number, fileType?: FileType | null, createdAt: string, updatedAt: string }> | null, sources: Array<{ __typename?: 'TorrentSource', key: string, name: string }> }, content?: { __typename?: 'Content', type: ContentType, source: string, id: string, title: string, releaseDate?: string | null, releaseYear?: number | null, overview?: string | null, runtime?: number | null, voteAverage?: number | null, voteCount?: number | null, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string }, originalLanguage?: { __typename?: 'LanguageInfo', id: string, name: string } | null, attributes: Array<{ __typename?: 'ContentAttribute', source: string, key: string, value: string, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }>, collections: Array<{ __typename?: 'ContentCollection', type: string, source: string, id: string, name: string, createdAt: string, updatedAt: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }>, externalLinks: Array<{ __typename?: 'ExternalLink', url: string, metadataSource: { __typename?: 'MetadataSource', key: string, name: string } }> } | null, languages?: Array<{ __typename?: 'LanguageInfo', id: string, name: string }> | null, episodes?: { __typename?: 'Episodes', label: string, seasons: Array<{ __typename?: 'Season', season: number, episodes?: Array<number> | null }> } | null }>, aggregations: { __typename?: 'TorrentContentAggregations', contentType?: Array<{ __typename?: 'ContentTypeAgg', value?: ContentType | null, label: string, count: number }> | null, torrentSource?: Array<{ __typename?: 'TorrentSourceAgg', value: string, label: string, count: number }> | null, torrentTag?: Array<{ __typename?: 'TorrentTagAgg', value: string, label: string, count: number }> | null, torrentFileType?: Array<{ __typename?: 'TorrentFileTypeAgg', value: FileType, label: string, count: number }> | null, language?: Array<{ __typename?: 'LanguageAgg', value: Language, label: string, count: number }> | null, genre?: Array<{ __typename?: 'GenreAgg', value: string, label: string, count: number }> | null, releaseYear?: Array<{ __typename?: 'ReleaseYearAgg', value?: number | null, label: string, count: number }> | null, videoResolution?: Array<{ __typename?: 'VideoResolutionAgg', value?: VideoResolution | null, label: string, count: number }> | null, videoSource?: Array<{ __typename?: 'VideoSourceAgg', value?: VideoSource | null, label: string, count: number }> | null } } } };

export type TorrentSuggestTagsQueryVariables = Exact<{
  query: SuggestTagsQueryInput;
}>;


export type TorrentSuggestTagsQuery = { __typename?: 'Query', torrent: { __typename?: 'TorrentQuery', suggestTags: { __typename?: 'TorrentSuggestTagsResult', suggestions: Array<{ __typename?: 'SuggestedTag', name: string, count: number }> } } };

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
export const TorrentFragmentDoc = gql`
    fragment Torrent on Torrent {
  infoHash
  name
  size
  private
  filesStatus
  hasFilesInfo
  singleFile
  fileType
  files {
    ...TorrentFile
  }
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
    ${TorrentFileFragmentDoc}`;
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
  releaseDate
  releaseYear
  createdAt
  updatedAt
}
    ${TorrentFragmentDoc}
${ContentFragmentDoc}`;
export const TorrentContentSearchResultFragmentDoc = gql`
    fragment TorrentContentSearchResult on TorrentContentSearchResult {
  totalCount
  items {
    ...TorrentContent
  }
  aggregations {
    contentType {
      value
      label
      count
    }
    torrentSource {
      value
      label
      count
    }
    torrentTag {
      value
      label
      count
    }
    torrentFileType {
      value
      label
      count
    }
    language {
      value
      label
      count
    }
    genre {
      value
      label
      count
    }
    releaseYear {
      value
      label
      count
    }
    videoResolution {
      value
      label
      count
    }
    videoSource {
      value
      label
      count
    }
  }
}
    ${TorrentContentFragmentDoc}`;
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
export const TorrentContentSearchDocument = gql`
    query TorrentContentSearch($query: SearchQueryInput, $facets: TorrentContentFacetsInput) {
  torrentContent {
    search(query: $query, facets: $facets) {
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
export const TorrentSuggestTagsDocument = gql`
    query TorrentSuggestTags($query: SuggestTagsQueryInput!) {
  torrent {
    suggestTags(query: $query) {
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