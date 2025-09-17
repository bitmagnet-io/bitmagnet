import { Component, inject } from "@angular/core";
import { AppModule } from "../../app.module";
import { InvitationsTableComponent } from "./invitations-table.component";
import { MatDialog } from "@angular/material/dialog";
import { InvitationsAddComponent } from "./invitations-add-dialog.component";
import { BehaviorSubject } from "rxjs";

@Component({
  selector: "app-invitations-admin",
  template: `
    <ng-container *transloco="let t">
      <app-document-title [parts]="[t('routes.invitations')]" />
      <mat-card>
        <mat-card-header>
          <mat-toolbar>
            <h2><mat-icon>rsvp</mat-icon>{{ t("routes.invitations") }}</h2>
            <button mat-stroked-button (click)="newInvitation()">
              <mat-icon>add</mat-icon>{{ t("auth.new_invitation") }}
            </button>
          </mat-toolbar>
        </mat-card-header>
        <mat-card-content>
          <app-invitations-table [update]="update$" />
        </mat-card-content>
      </mat-card>
    </ng-container>
  `,
  styles: [
    `
      mat-toolbar {
        h2 mat-icon {
          position: relative;
          top: 3px;
          margin-right: 14px;
          margin-left: 20px;
        }

        justify-content: space-between;
      }
    `,
  ],

  standalone: true,
  imports: [AppModule, InvitationsTableComponent],
})
export class InvitationsAdminComponent {
  readonly dialog = inject(MatDialog);
  private update = new BehaviorSubject<void>(void 0);
  update$ = this.update.asObservable();

  newInvitation() {
    const ref = this.dialog.open(InvitationsAddComponent);
    ref.afterClosed().subscribe(() => this.update.next(void 0));
  }
}
