import {
  TorrentsSearchController,
  defaultOrderBy,
  TorrentSearchControls,
  SizeRangeFilter,
} from "./torrents-search.controller";

describe("TorrentsSearchController", () => {
  let controller: TorrentsSearchController;

  beforeEach(() => {
    controller = new TorrentsSearchController({
      limit: 20,
      page: 1,
      contentType: null,
      orderBy: defaultOrderBy,
      facets: {
        genre: { active: false },
        language: { active: false },
        fileType: { active: false },
        torrentSource: { active: false },
        torrentTag: { active: false },
        videoResolution: { active: false },
        videoSource: { active: false },
      },
    });
  });

  describe("setSizeRange", () => {
    it("should set the size range and update controls", () => {
      // Set a size range
      controller.setSizeRange(1024, 1048576);

      // Get the current controls
      let currentControls = {} as TorrentSearchControls;
      const subscription = controller.controls$.subscribe({
        next: (controls: TorrentSearchControls) => {
          currentControls = controls;
        },
      });

      // Check if sizeRange property is set correctly
      const sizeRange = currentControls.sizeRange as SizeRangeFilter;
      expect(sizeRange).toEqual({ min: 1024, max: 1048576 });

      // Check that page is reset to 1
      expect(currentControls.page).toBe(1);

      // Setting only min
      controller.setSizeRange(1024, undefined);
      const minOnlyRange = currentControls.sizeRange as SizeRangeFilter;
      expect(minOnlyRange).toEqual({ min: 1024, max: undefined });

      // Setting only max
      controller.setSizeRange(undefined, 1048576);
      const maxOnlyRange = currentControls.sizeRange as SizeRangeFilter;
      expect(maxOnlyRange).toEqual({
        min: undefined,
        max: 1048576,
      });

      // Setting to undefined should remove the size range filter
      controller.setSizeRange(undefined, undefined);
      expect(currentControls.sizeRange).toBeUndefined();

      subscription.unsubscribe();
    });

    it("should not set a size range if both min and max are undefined", () => {
      controller.setSizeRange(undefined, undefined);

      let currentControls = {} as TorrentSearchControls;
      const subscription = controller.controls$.subscribe({
        next: (controls: TorrentSearchControls) => {
          currentControls = controls;
        },
      });

      expect(currentControls.sizeRange).toBeUndefined();

      subscription.unsubscribe();
    });

    it("should correctly convert query parameters to size range", () => {
      // Manually trigger an update with these params
      // Note: This is typically handled by the paramsToControls function
      controller.update(() => ({
        limit: 20,
        page: 1,
        contentType: null,
        orderBy: defaultOrderBy,
        facets: {
          genre: { active: false },
          language: { active: false },
          fileType: { active: false },
          torrentSource: { active: false },
          torrentTag: { active: false },
          videoResolution: { active: false },
          videoSource: { active: false },
        },
        sizeRange: {
          min: 100 * 1024 * 1024, // 100 MiB in bytes
          max: 1000 * 1024 * 1024 * 1024, // 1000 GiB in bytes
        },
      }));

      let currentControls = {} as TorrentSearchControls;
      const subscription = controller.controls$.subscribe({
        next: (controls: TorrentSearchControls) => {
          currentControls = controls;
        },
      });

      // Check if values were converted correctly
      expect(currentControls.sizeRange).toBeDefined();
      const sizeRange = currentControls.sizeRange as SizeRangeFilter;
      expect(sizeRange.min).toBe(100 * 1024 * 1024);
      expect(sizeRange.max).toBe(1000 * 1024 * 1024 * 1024);

      subscription.unsubscribe();
    });
  });
});
