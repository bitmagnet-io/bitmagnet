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
        (c) => c.DashboardComponent
      ),
    title: 'Dashboard'
  },
  // {
  //   path: 'address',
  //   loadComponent: () =>
  //     import('./address-form/address-form.component').then(
  //       (c) => c.AddressFormComponent
  //     ),
  //   title: 'Address'
  // },
  {
    path: 'torrents',
    loadComponent: () =>
      import('./torrents-search/torrents-search.component').then(
        (c) => c.TorrentsSearchComponent
      ),
    title: 'Torrents'
  },
  // {
  //   path: 'tree',
  //   loadComponent: () =>
  //     import('./tree/tree.component').then(
  //       (c) => c.TreeComponent
  //     ),
  //   title: 'Tree'
  // },
  // {
  //   path: 'drag-drop',
  //   loadComponent: () =>
  //     import('./drag-drop/drag-drop.component').then(
  //       (c) => c.DragDropComponent
  //     ),
  //   title: 'Drag-Drop'
  // },
];
