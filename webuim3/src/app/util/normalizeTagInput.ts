const normalizeTagInput = (value: string): string =>
  value
    .toLowerCase()
    .replaceAll(/[^a-z0-9\-]/g, '-')
    .replace(/^-+/, '')
    .replaceAll(/-+/g, '-');

export default normalizeTagInput;
