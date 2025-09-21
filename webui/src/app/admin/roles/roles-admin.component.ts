import { Component, inject } from "@angular/core";
import { AppModule } from "../../app.module";
import { RolesService } from "../../auth/roles.service";
import { RoleEditComponent } from "./role-edit.component";
import { MatDialog } from "@angular/material/dialog";
import { RoleAddComponent } from "./role-add.component";
import { Observable } from "rxjs";
import * as generated from "../../graphql/generated";

@Component({
  selector: "app-roles-admin",
  template: `
    <ng-container *transloco="let t">
      <app-document-title [parts]="[t('routes.roles')]" />
      <mat-card>
        <mat-card-header>
          <mat-toolbar>
            <h2>
              <mat-icon>approval_delegation</mat-icon>{{ t("routes.roles") }}
            </h2>
            <button mat-stroked-button (click)="newRole()">
              <mat-icon>add</mat-icon>{{ t("auth.new_role") }}
            </button>
          </mat-toolbar>
        </mat-card-header>
        <mat-card-content>
          <mat-accordion>
            @for (role of roles$ | async; track role.name) {
              <mat-expansion-panel
                hideToggle
                [expanded]="activeRole === role.name"
                (opened)="activeRole = role.name"
                (closed)="activeRole = undefined"
              >
                <mat-expansion-panel-header>
                  <mat-panel-title>{{ role.name }}</mat-panel-title>
                </mat-expansion-panel-header>
                <mat-divider />
                <app-role-edit [role]="role" (updated)="refresh()" />
              </mat-expansion-panel>
            }
          </mat-accordion>
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

      mat-expansion-panel {
        mat-divider {
          margin-bottom: 20px;
        }
      }
    `,
  ],

  standalone: true,
  imports: [AppModule, RoleEditComponent],
})
export class RolesAdminComponent {
  private rolesService = inject(RolesService);
  readonly dialog = inject(MatDialog);

  roles$: Observable<generated.Role[]>;

  activeRole?: string;

  constructor() {
    this.refresh();
  }

  refresh() {
    this.roles$ = this.rolesService.listRoles();
  }

  newRole() {
    const ref = this.dialog.open(RoleAddComponent);
    ref.afterClosed().subscribe(() => {
      this.roles$ = this.rolesService.listRoles();
    });
  }
}
