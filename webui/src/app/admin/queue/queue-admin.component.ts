import { Component } from "@angular/core";
import { AppModule } from "../../app.module";

@Component({
  selector: "app-queue-admin",
  template: `
    <ng-container *transloco="let t">
      <mat-card class="admin-card">
        <mat-card-header>
          <mat-toolbar>
            <nav
              mat-tab-nav-bar
              [tabPanel]="tabPanel"
              mat-stretch-tabs="false"
              mat-align-tabs="start"
            >
              <h2><mat-icon svgIcon="queue" />{{ t("routes.queues") }}</h2>
              <a
                mat-tab-link
                routerLink="visualize"
                #linkVisualize="routerLinkActive"
                routerLinkActive
                [active]="linkVisualize.isActive"
                ><mat-icon>monitoring</mat-icon>{{ t("routes.visualize") }}</a
              >
              <a
                mat-tab-link
                routerLink="jobs"
                #linkJobs="routerLinkActive"
                routerLinkActive
                [active]="linkJobs.isActive"
                ><mat-icon>toc</mat-icon>{{ t("routes.jobs") }}</a
              >
              <a
                mat-tab-link
                routerLink="admin"
                #linkAdmin="routerLinkActive"
                routerLinkActive
                [active]="linkAdmin.isActive"
                ><mat-icon>construction</mat-icon>{{ t("routes.admin") }}</a
              >
            </nav>
            <mat-tab-nav-panel #tabPanel></mat-tab-nav-panel>
          </mat-toolbar>
        </mat-card-header>
        <mat-card-content>
          <router-outlet #outlet />
        </mat-card-content>
      </mat-card>
    </ng-container>
  `,
  styles: [
    `
      mat-card-header {
        flex-wrap: wrap;
        h2 {
          font-size: 18px;
          margin: 0 60px 0 48px;
          height: 48px;
          line-height: 48px;
          mat-icon {
            position: relative;
            top: 6px;
            margin-right: 14px;
            line-height: 1.25rem;
          }
        }
        nav {
          flex: 0 0 100%;
          a {
            margin-top: 2px;
            mat-icon {
              margin-right: 12px;
            }
          }
        }
      }
    `,
  ],
  standalone: true,
  imports: [AppModule],
})
export class QueueAdminComponent {}
