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
];
