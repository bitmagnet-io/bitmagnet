import { Component, inject } from "@angular/core";
import { BreakpointsService } from "../layout/breakpoints.service";
import { AppModule } from "../app.module";

@Component({
  selector: "app-admin",
  template: `
    <ng-container *transloco="let t">
      <mat-drawer-container class="drawer-container">
        <mat-drawer
          #drawer
          class="drawer"
          [attr.role]="
            breakpoints.sizeAtLeast('Medium') ? 'navigation' : 'dialog'
          "
          [mode]="breakpoints.sizeAtLeast('Medium') ? 'side' : 'over'"
          [opened]="breakpoints.sizeAtLeast('Medium')"
        >
          <nav>
            <ul>
              <li>
                <a
                  mat-button
                  routerLink="/admin"
                  routerLinkActive
                  [routerLinkActiveOptions]="{ exact: true }"
                  #linkHome="routerLinkActive"
                  [class]="linkHome.isActive ? 'active' : ''"
                  ><mat-icon>admin</mat-icon>{{ t("routes.home") }}</a
                >
              </li>
              <li>
                <a
                  mat-button
                  routerLink="workers"
                  routerLinkActive
                  #linkWorkers="routerLinkActive"
                  [class]="linkWorkers.isActive ? 'active' : ''"
                  ><mat-icon>manufacturing</mat-icon
                  >{{ t("routes.workers") }}</a
                >
              </li>
              <li>
                <a
                  mat-button
                  routerLink="queues"
                  routerLinkActive
                  #linkQueues="routerLinkActive"
                  [class]="linkQueues.isActive ? 'active' : ''"
                  ><mat-icon svgIcon="queue" />{{ t("routes.queues") }}</a
                >
              </li>
              <!-- <li>
                <a
                  mat-button
                  routerLink="torrents"
                  routerLinkActive
                  #linkTorrents="routerLinkActive"
                  [class]="linkTorrents.isActive ? 'active' : ''"
                  ><mat-icon svgIcon="magnet" />{{ t("routes.torrents") }}</a
                >
              </li> -->
              <li>
                <a
                  mat-button
                  routerLink="users"
                  routerLinkActive
                  #linkUsers="routerLinkActive"
                  [class]="linkUsers.isActive ? 'active' : ''"
                  ><mat-icon>groups</mat-icon>{{ t("routes.users") }}</a
                >
              </li>
              <li>
                <a
                  mat-button
                  routerLink="roles"
                  routerLinkActive
                  #linkRoles="routerLinkActive"
                  [class]="linkRoles.isActive ? 'active' : ''"
                  ><mat-icon>approval_delegation</mat-icon
                  >{{ t("routes.roles") }}</a
                >
              </li>
              <li>
                <a
                  mat-button
                  routerLink="invitations"
                  routerLinkActive
                  #linkInvitations="routerLinkActive"
                  [class]="linkInvitations.isActive ? 'active' : ''"
                  ><mat-icon>rsvp</mat-icon>{{ t("routes.invitations") }}</a
                >
              </li>
              <li>
                <a
                  mat-button
                  routerLink="config"
                  routerLinkActive
                  #linkConfig="routerLinkActive"
                  [class]="linkConfig.isActive ? 'active' : ''"
                  ><mat-icon>settings</mat-icon>{{ t("routes.config") }}</a
                >
              </li>
            </ul>
          </nav>
        </mat-drawer>
        <mat-drawer-content>
          <div
            class="form-field-container button-container button-container-toggle-drawer"
          >
            <button
              type="button"
              class="button-toggle-drawer"
              mat-icon-button
              (click)="drawer.toggle()"
              [matTooltip]="t('torrents.toggle_drawer')"
            >
              <mat-icon
                aria-label="Side nav toggle icon"
                fontSet="material-icons"
                >{{
                  drawer.opened ? "arrow_circle_left" : "arrow_circle_right"
                }}</mat-icon
              >
            </button>
          </div>
          <router-outlet #outlet />
        </mat-drawer-content>
      </mat-drawer-container>
    </ng-container>
  `,
  styles: [
    `
      mat-drawer {
        nav {
          padding-top: 12px;
          --mat-text-button-icon-spacing: 14px;
          ul {
            list-style-type: none;
            padding-left: 0;
            a {
              width: 100%;
              font-size: var(--mat-expansion-container-text-size);
              justify-content: flex-start;
              padding-left: 20px;
            }
            li {
              margin-bottom: 6px;
            }
          }
        }
      }

      mat-drawer-content {
        margin-left: 10px;
        .button-container-toggle-drawer {
          position: absolute;
          left: 0;
          top: 28px;
          z-index: 100;
        }
      }

      .drawer {
        width: 220px;
      }
    `,
  ],
  standalone: true,
  imports: [AppModule],
})
export class AdminComponent {
  breakpoints = inject(BreakpointsService);
}
