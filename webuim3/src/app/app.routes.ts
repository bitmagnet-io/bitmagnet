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
    title: 'dashboard',
    children: [
      {
        path: 'home',
        loadComponent: () =>
          import('./dashboard/dashboard-home.component').then(
            (c) => c.DashboardHomeComponent,
          ),
        title: "home",
      },
      {
        path: 'queues',
        loadComponent: () =>
          import('./queue/queue-card.component').then(
            (c) => c.QueueCardComponent,
          ),
        title: 'queues',
      },
    ]
  },
  {
    path: 'torrents',
    loadComponent: () =>
      import('./torrents-search/torrents-search.component').then(
        (c) => c.TorrentsSearchComponent,
      ),
    title: 'torrents',
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
