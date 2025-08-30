import {
  __spreadValues
} from "./chunk-DMMUMX3A.js";

// src/app/torrents/content-types.ts
var contentTypeMap = {
  movie: {
    singular: "Movie",
    plural: "Movies",
    icon: "movie"
  },
  tv_show: {
    singular: "TV Show",
    plural: "TV Shows",
    icon: "live_tv"
  },
  music: {
    singular: "Music",
    plural: "Music",
    icon: "music_note"
  },
  ebook: {
    singular: "E-Book",
    plural: "E-Books",
    icon: "auto_stories"
  },
  comic: {
    singular: "Comic",
    plural: "Comics",
    icon: "comic_bubble"
  },
  audiobook: {
    singular: "Audiobook",
    plural: "Audiobooks",
    icon: "mic"
  },
  software: {
    singular: "Software",
    plural: "Software",
    icon: "desktop_windows"
  },
  game: {
    singular: "Game",
    plural: "Games",
    icon: "sports_esports"
  },
  xxx: {
    singular: "XXX",
    plural: "XXX",
    icon: "18_up_rating"
  },
  null: {
    singular: "Unknown",
    plural: "Unknown",
    icon: "question_mark"
  }
};
var contentTypeList = Object.entries(contentTypeMap).map(([key, info]) => __spreadValues({
  key
}, info));
var contentTypeInfo = (key) => key ? contentTypeMap[key] : void 0;

export {
  contentTypeMap,
  contentTypeList,
  contentTypeInfo
};
//# sourceMappingURL=chunk-UGVUNZOV.js.map
