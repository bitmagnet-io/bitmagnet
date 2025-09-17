import {
  Router,
  Routes,
  CanActivateFn,
  RouterStateSnapshot,
} from "@angular/router";
import { AuthService, ObjectAction } from "./auth/auth.service";
import { inject } from "@angular/core";
import { map } from "rxjs";

const navigateToLogin = (router: Router, state: RouterStateSnapshot) => {
  router.navigate(["/login"], {
    queryParams: { returnUrl: state.url },
  });
};

const authGuard = (...objectActions: ObjectAction[]): CanActivateFn => {
  return ({}, state) => {
    const authService = inject(AuthService);
    const router = inject(Router);
    return authService.enforce(...objectActions).pipe(
      map((allowed) => {
        if (allowed) {
          return true;
        }

        navigateToLogin(router, state);

        return false;
      }),
    );
  };
};

const requireUserGuard =
  (
    require: boolean,
    onFailed: (router: Router, state: RouterStateSnapshot) => void,
  ): CanActivateFn =>
  ({}, state) => {
    const authService = inject(AuthService);
    const router = inject(Router);
    return authService.self$.pipe(
      map(({ user }) => {
        if (require !== !user) {
          return true;
        }

        onFailed(router, state);

        return false;
      }),
    );
  };

export const routes: Routes = [
  {
    path: "",
    pathMatch: "full",
    redirectTo: "torrents",
  },
  {
    path: "torrents",
    canActivate: Array<CanActivateFn>(
      authGuard({
        object: "torrent",
        action: "query",
      }),
    ),
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
    path: "login",
    canActivate: [
      requireUserGuard(false, (router) => router.navigate(["/account"])),
    ],
    loadComponent: () =>
      import("./auth/login.component").then((c) => c.LoginComponent),
  },
  {
    path: "register",
    canActivate: [
      requireUserGuard(false, (router) => router.navigate(["/account"])),
    ],
    loadComponent: () =>
      import("./auth/register.component").then((c) => c.RegisterComponent),
  },
  {
    path: "account",
    canActivate: [requireUserGuard(true, navigateToLogin)],
    loadComponent: () =>
      import("./account/account.component").then((c) => c.AccountComponent),
    children: [
      {
        path: "",
        loadComponent: () =>
          import("./account/account-home.component").then(
            (c) => c.AccountHomeComponent,
          ),
      },
    ],
  },
  {
    path: "admin",
    loadComponent: () =>
      import("./admin/admin.component").then((c) => c.AdminComponent),
    children: [
      {
        path: "",
        loadComponent: () =>
          import("./admin/admin-home.component").then(
            (c) => c.AdminHomeComponent,
          ),
      },
      {
        path: "config",
        canActivate: Array<CanActivateFn>(
          authGuard({
            object: "config",
            action: "query",
          }),
        ),
        loadComponent: () =>
          import("./admin/config/config-admin.component").then(
            (c) => c.AdminConfigComponent,
          ),
      },
      {
        path: "workers",
        canActivate: Array<CanActivateFn>(
          authGuard({
            object: "worker",
            action: "query",
          }),
        ),
        loadComponent: () =>
          import("./admin/workers/admin-workers.component").then(
            (c) => c.AdminWorkersComponent,
          ),
      },
      {
        path: "queues",
        pathMatch: "full",
        redirectTo: "queues/visualize",
      },
      {
        path: "queues",
        canActivate: Array<CanActivateFn>(
          authGuard({
            object: "queue",
            action: "query",
          }),
        ),
        loadComponent: () =>
          import("./admin/queue/queue-admin.component").then(
            (c) => c.QueueAdminComponent,
          ),
        children: [
          {
            path: "visualize",
            loadComponent: () =>
              import("./admin/queue/queue-visualize.component").then(
                (c) => c.QueueVisualizeComponent,
              ),
          },
          {
            path: "jobs",
            loadComponent: () =>
              import("./admin/queue/queue-jobs.component").then(
                (c) => c.QueueJobsComponent,
              ),
          },
          {
            path: "admin",
            loadComponent: () =>
              import("./admin/queue/queue-manager.component").then(
                (c) => c.QueueManageComponent,
              ),
          },
        ],
      },
      {
        path: "torrents",
        loadComponent: () =>
          import("./admin/torrents/torrents-admin.component").then(
            (c) => c.TorrentsAdminComponent,
          ),
      },
      {
        path: "users",
        loadComponent: () =>
          import("./admin/users/users-admin.component").then(
            (c) => c.UsersAdminComponent,
          ),
      },
      {
        path: "roles",
        loadComponent: () =>
          import("./admin/roles/roles-admin.component").then(
            (c) => c.RolesAdminComponent,
          ),
      },
      {
        path: "invitations",
        loadComponent: () =>
          import("./admin/invitations/invitations-admin.component").then(
            (c) => c.InvitationsAdminComponent,
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
