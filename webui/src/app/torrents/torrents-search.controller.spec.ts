import {
  TorrentsSearchController,
  defaultOrderBy,
  TorrentSearchControls,
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

  describe("setPublishedAt", () => {
    it("should update controls with publishedAt time frame", () => {
      controller.setPublishedAt("7d");

      let currentControls: TorrentSearchControls = {} as TorrentSearchControls;
      const subscription = controller.controls$.subscribe((controls) => {
        currentControls = controls;
      });

      expect(currentControls.publishedAt).toBe("7d");
      expect(currentControls.page).toBe(1); // Should reset to page 1

      subscription.unsubscribe();
    });

    it("should update controls with special time frame", () => {
      controller.setPublishedAt("this month");

      let currentControls: TorrentSearchControls = {} as TorrentSearchControls;
      const subscription = controller.controls$.subscribe((controls) => {
        currentControls = controls;
      });

      expect(currentControls.publishedAt).toBe("this month");

      subscription.unsubscribe();
    });

    it("should update controls with date range", () => {
      controller.setPublishedAt("2023-01-01 to 2023-01-31");

      let currentControls: TorrentSearchControls = {} as TorrentSearchControls;
      const subscription = controller.controls$.subscribe((controls) => {
        currentControls = controls;
      });

      expect(currentControls.publishedAt).toBe("2023-01-01 to 2023-01-31");

      subscription.unsubscribe();
    });

    it("should remove publishedAt when value is undefined or empty", () => {
      // First set a time frame
      controller.setPublishedAt("7d");

      let currentControls: TorrentSearchControls = {} as TorrentSearchControls;
      const subscription = controller.controls$.subscribe((controls) => {
        currentControls = controls;
      });

      expect(currentControls.publishedAt).toBe("7d");

      // Then clear it
      controller.setPublishedAt(undefined);
      expect(currentControls.publishedAt).toBeUndefined();

      // Set it again and clear with empty string
      controller.setPublishedAt("30d");
      expect(currentControls.publishedAt).toBe("30d");

      controller.setPublishedAt("");
      expect(currentControls.publishedAt).toBeUndefined();

      subscription.unsubscribe();
    });
  });

  describe("controlsToQueryVariables", () => {
    it("should add publishedAt to facets when set", (done) => {
      let foundPublishedAt = false;
      let checkCount = 0;

      // Set up subscription
      const subscription = controller.params$.subscribe((params) => {
        checkCount++;

        // Check structure exists
        expect(params).toBeDefined();
        expect(params.input).toBeDefined();
        expect(params.input.facets).toBeDefined();

        // Type assertion for accessing facets properties
        // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/no-unsafe-assignment
        const facets = params.input.facets as any;

        // On the second emission, we should have publishedAt
        if (checkCount === 2) {
          // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
          expect(facets.publishedAt).toBe("7d");
          foundPublishedAt = true;
          done(); // Signal test completion
        }
      });

      // Initial check should not have publishedAt

      // Now set the published date - this should trigger the params$ observable
      controller.setPublishedAt("7d");

      // Cleanup subscription in case the test times out
      setTimeout(() => {
        subscription.unsubscribe();
        if (!foundPublishedAt) {
          done.fail("Timeout: did not find publishedAt in facets");
        }
      }, 2000);
    });

    it("should set the correct value of publishedAt in facets", (done) => {
      const testValue = "this month";
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      let latestParams: any = null;

      // Set up subscription
      const subscription = controller.params$.subscribe((params) => {
        latestParams = params;
      });

      // Set the published date
      controller.setPublishedAt(testValue);

      // Check the result
      setTimeout(() => {
        subscription.unsubscribe();
        try {
          expect(latestParams).toBeDefined();
          // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
          expect(latestParams.input.facets.publishedAt).toBe(testValue);
          done();
        } catch (e) {
          done.fail(e instanceof Error ? e.message : String(e));
        }
      }, 500);
    });
  });
});
