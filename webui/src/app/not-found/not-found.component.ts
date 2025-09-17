import { Component } from "@angular/core";
import { AppModule } from "../app.module";
import { DocumentTitleComponent } from "../layout/document-title.component";

@Component({
  selector: "app-not-found",
  template: `
    <ng-container *transloco="let t">
      <app-document-title [parts]="[t('general.page_not_found')]" />
      <mat-card class="card-not-found">
        <mat-card-header>
          <mat-card-title>
            <h2>{{ t("general.page_not_found") }}</h2>
          </mat-card-title>
        </mat-card-header>
      </mat-card>
    </ng-container>
  `,
  styles: [
    `
      .card-not-found {
        max-width: 960px;
        margin: 20px auto;
        h2 {
          margin-top: 10px;
        }
      }
    `,
  ],
  standalone: true,
  imports: [AppModule, DocumentTitleComponent],
})
export class NotFoundComponent {}
