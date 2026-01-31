import { Component } from "@angular/core";
import { AppModule } from "../../app.module";
import { UsersTableComponent } from "./users-table.component";

@Component({
  selector: "app-users-admin",
  template: `
    <ng-container *transloco="let t">
      <app-document-title [parts]="[t('routes.users')]" />
      <mat-card>
        <mat-card-header>
          <mat-toolbar>
            <h2><mat-icon>groups</mat-icon>{{ t("routes.users") }}</h2>
          </mat-toolbar>
        </mat-card-header>
        <mat-card-content>
          <app-users-table />
        </mat-card-content>
      </mat-card>
    </ng-container>
  `,
  styles: [
    `
      mat-toolbar h2 mat-icon {
        position: relative;
        top: 3px;
        margin-right: 14px;
        margin-left: 20px;
      }
    `,
  ],

  standalone: true,
  imports: [AppModule, UsersTableComponent],
})
export class UsersAdminComponent {}
