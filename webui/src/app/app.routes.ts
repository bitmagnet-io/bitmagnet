import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'dashboard',
  },
  {
    path: 'dashboard',
    loadComponent: () =>
      import('./dashboard/dashboard.component').then(
        (c) => c.DashboardComponent,
      ),
    children: [
      {
        path: 'home',
        loadComponent: () =>
          import('./dashboard/dashboard-home.component').then(
            (c) => c.DashboardHomeComponent,
          ),
      },
      {
        path: 'queues',
        loadComponent: () =>
          import('./queue/queue-dashboard.component').then(
            (c) => c.QueueDashboardComponent,
          ),
        children: [
          {
            path: 'visualize',
            loadComponent: () =>
              import('./queue/queue-visualize.component').then(
                (c) => c.QueueVisualizeComponent,
              ),
          },
          {
            path: 'jobs',
            loadComponent: () =>
              import('./queue/queue-jobs.component').then(
                (c) => c.QueueJobsComponent,
              ),
          },
          {
            path: 'admin',
            loadComponent: () =>
              import('./queue/queue-admin.component').then(
                (c) => c.QueueAdminComponent,
              ),
          }
        ]
      },
    ],
  },
  {
    path: 'torrents',
    loadComponent: () =>
      import('./torrents-search/torrents-search.component').then(
        (c) => c.TorrentsSearchComponent,
      ),
  },
  {
    path: 'torrents/:infoHash',
    loadComponent: () =>
      import('./torrent-permalink/torrent-permalink.component').then(
        (c) => c.TorrentPermalinkComponent,
      ),
  },
  // {
  //   path: 'queue',
  //   loadComponent: () =>
  //     import('./queue/queue-card.component').then(
  //       (c) => c.QueueCardComponent,
  //     ),
  //   title: 'queue',
  // },
];
