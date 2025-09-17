import { Component } from "@angular/core";
import { AppModule } from "../../app.module";
import { TorrentMetricsComponent } from "./torrent-metrics.component";

@Component({
  selector: "app-torrents",
  template: `
    <ng-container *transloco="let t">
      <app-document-title [parts]="[t('routes.torrents'), t('routes.admin')]" />
      <mat-card class="admin-card">
        <mat-card-header>
          <mat-toolbar>
            <h2><mat-icon svgIcon="magnet" />{{ t("routes.torrents") }}</h2>
          </mat-toolbar>
        </mat-card-header>
        <mat-card-content>
          <app-torrent-metrics />
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
  imports: [AppModule, TorrentMetricsComponent],
})
export class TorrentsAdminComponent {}
