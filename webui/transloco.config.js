module.exports = {
  rootTranslationsPath: "src/app/i18n/translations/",
  langs: ["ar", "es", "de", "en", "fr", "ja", "nl", "ru", "uk", "zh"],
  keysManager: {
    output: "app/i18n/translations",
    translationsPath: "app/i18n/translations",
    addMissingKeys: true,
    emitErrorOnExtraKeys: true,
    unflat: true,
    sort: true,
    defaultValue: "__missing__",
  },
};
