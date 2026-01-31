import { Component, inject, OnInit } from "@angular/core";
import { Router } from "@angular/router";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { map } from "rxjs";
import { AppModule } from "../app.module";
import { AuthService } from "../auth/auth.service";
import { BreakpointsService } from "../layout/breakpoints.service";

@Component({
  selector: "app-account",
  template: ` <ng-container *transloco="let t">
    <app-document-title [parts]="[t('account.account')]" />
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
                routerLink="/account"
                routerLinkActive
                [routerLinkActiveOptions]="{ exact: true }"
                #linkHome="routerLinkActive"
                [class]="linkHome.isActive ? 'active' : ''"
                ><mat-icon>person</mat-icon>{{ t("routes.account") }}</a
              >
            </li>
            <li>
              <a
                mat-button
                routerLink="/account/api-keys"
                routerLinkActive
                [routerLinkActiveOptions]="{ exact: true }"
                #linkAPIKeys="routerLinkActive"
                [class]="linkAPIKeys.isActive ? 'active' : ''"
                ><mat-icon>key</mat-icon>{{ t("routes.api_keys") }}</a
              >
            </li>
            <li>
              <a mat-button (click)="logout()"
                ><mat-icon>logout</mat-icon>{{ t("auth.logout") }}</a
              >
            </li>
          </ul>
        </nav> </mat-drawer
      >>
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
  </ng-container>`,
  styles: [
    `
      mat-drawer nav {
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

      mat-drawer-content {
        .button-container-toggle-drawer {
          position: absolute;
          left: 20px;
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
export class AccountComponent implements OnInit {
  private authService = inject(AuthService);
  private router = inject(Router);
  breakpoints = inject(BreakpointsService);

  user$ = this.authService.self$.pipe(
    takeUntilDestroyed(),
    map((self) => self.user || null),
  );

  ngOnInit(): void {
    this.user$.subscribe((user) => {
      if (!user) {
        void this.router.navigate(["/"]);
      }
    });
  }

  logout() {
    this.authService.logout().subscribe();
  }
}
