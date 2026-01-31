import { Component, inject } from "@angular/core";
import { AppModule } from "../../app.module";
import { DocumentTitleComponent } from "../../layout/document-title.component";
import { BreakpointsService } from "../../layout/breakpoints.service";
import { WorkersModule } from "../../workers/workers.module";

@Component({
  selector: "app-workers-admin",
  template: `
    <ng-container *transloco="let t">
      <app-document-title [parts]="[t('routes.admin')]" />
      <mat-card class="admin-card">
        <mat-card-header>
          <mat-toolbar>
            <h2><mat-icon>manufacturing</mat-icon>{{ t("routes.workers") }}</h2>
          </mat-toolbar>
        </mat-card-header>
        <mat-card-content>
          <mat-divider></mat-divider>
          <div class="grid-container">
            <mat-grid-list
              [cols]="breakpoints.sizeAtLeast('Medium') ? 2 : 1"
              rowHeight="600px"
            >
              <mat-grid-tile [colspan]="1" [rowspan]="1">
                <app-workers-card></app-workers-card>
              </mat-grid-tile>
            </mat-grid-list>
          </div>
        </mat-card-content>
      </mat-card>
    </ng-container>
  `,
  styles: [
    `
      .grid-container {
        margin: 20px;
      }

      .more-button {
        position: absolute;
        top: 5px;
        right: 10px;
      }

      app-health-card {
        width: 100%;
        height: 100%;
        mat-card {
          height: 100%;
        }
      }

      mat-grid-tile mat-card {
        width: 100%;
      }

      mat-toolbar h2 mat-icon {
        position: relative;
        top: 3px;
        margin-right: 14px;
        margin-left: 32px;
      }
    `,
  ],
  standalone: true,
  imports: [AppModule, DocumentTitleComponent, WorkersModule],
})
export class AdminWorkersComponent {
  breakpoints = inject(BreakpointsService);
}
