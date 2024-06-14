module.exports = {
  rootTranslationsPath: 'src/app/i18n/',
  langs: ['en', 'fr', 'es', 'de', 'zh'],
  keysManager: {
    output: "app/i18n",
    translationsPath: "app/i18n",
    addMissingKeys: true,
    emitErrorOnExtraKeys: true,
    unflat: true,
    sort: true,
    defaultValue: "__missing__"
  }
};
