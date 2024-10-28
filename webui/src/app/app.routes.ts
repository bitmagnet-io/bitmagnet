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
          import("./torrents/torrents-search.component").then(
            (c) => c.TorrentsSearchComponent,
          ),
      },
      {
        path: "permalink/:infoHash",
        loadComponent: () =>
          import("./torrents/torrent-permalink.component").then(
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
        path: "queues",
        loadComponent: () =>
          import("./dashboard/queue/queue-dashboard.component").then(
            (c) => c.QueueDashboardComponent,
          ),
        children: [
          {
            path: "visualize",
            loadComponent: () =>
              import("./dashboard/queue/queue-visualize.component").then(
                (c) => c.QueueVisualizeComponent,
              ),
          },
          {
            path: "jobs",
            loadComponent: () =>
              import("./dashboard/queue/queue-jobs.component").then(
                (c) => c.QueueJobsComponent,
              ),
          },
          {
            path: "admin",
            loadComponent: () =>
              import("./dashboard/queue/queue-admin.component").then(
                (c) => c.QueueAdminComponent,
              ),
          },
        ],
      },
      {
        path: "torrents",
        loadComponent: () =>
          import("./dashboard/torrents/torrents-dashboard.component").then(
            (c) => c.TorrentsDashboardComponent,
          ),
      },
    ],
  },
  {
    path: "**",
    loadComponent: () =>
      import("./not-found/not-found.component").then(
        (c) => c.NotFoundComponent,
      ),
  },
];
