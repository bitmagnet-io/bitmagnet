import { Routes } from "@angular/router";

export const routes: Routes = [
  {
    path: "",
    pathMatch: "full",
    redirectTo: "torrents",
  },
  {
    path: "torrents",
    loadComponent: () =>
      import("./torrents/torrents.component").then((c) => c.TorrentsComponent),
    children: [
      {
        path: "",
        loadComponent: () =>
          import("./torrents-search/torrents-search.component").then(
            (c) => c.TorrentsSearchComponent,
          ),
      },
      {
        path: "permalink/:infoHash",
        loadComponent: () =>
          import("./torrent-permalink/torrent-permalink.component").then(
            (c) => c.TorrentPermalinkComponent,
          ),
      },
    ],
  },
  {
    path: "dashboard",
    loadComponent: () =>
      import("./dashboard/dashboard.component").then(
        (c) => c.DashboardComponent,
      ),
    children: [
      {
        path: "",
        loadComponent: () =>
          import("./dashboard/dashboard-home.component").then(
            (c) => c.DashboardHomeComponent,
          ),
      },
      {
        path: "torrents",
        loadComponent: () =>
          import("./torrents-dashboard/torrents-dashboard.component").then(
            (c) => c.TorrentsDashboardComponent,
          ),
      },
      {
        path: "queues",
        loadComponent: () =>
          import("./queue/queue-dashboard.component").then(
            (c) => c.QueueDashboardComponent,
          ),
        children: [
          {
            path: "visualize",
            loadComponent: () =>
              import("./queue/queue-visualize.component").then(
                (c) => c.QueueVisualizeComponent,
              ),
          },
          {
            path: "jobs",
            loadComponent: () =>
              import("./queue/queue-jobs.component").then(
                (c) => c.QueueJobsComponent,
              ),
          },
          {
            path: "admin",
            loadComponent: () =>
              import("./queue/queue-admin.component").then(
                (c) => c.QueueAdminComponent,
              ),
          },
        ],
      },
    ],
  },
];
