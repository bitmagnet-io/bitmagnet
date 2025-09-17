import { Component, inject } from "@angular/core";
import { AppModule } from "../../app.module";
import { RolesService } from "../../auth/roles.service";
import { RoleEditComponent } from "./role-edit.component";

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
          </mat-toolbar>
        </mat-card-header>
        <mat-card-content>
          <mat-accordion>
            @for (role of roles$ | async; track role.name) {
              <mat-expansion-panel hideToggle>
                <mat-expansion-panel-header>
                  <mat-panel-title>{{ role.name }}</mat-panel-title>
                </mat-expansion-panel-header>
                <mat-divider />
                <app-role-edit [role]="role" />
              </mat-expansion-panel>
            }
          </mat-accordion>
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

  roles$ = this.rolesService.listRoles();
  // objectActions$ = this.rolesService.listObjectActions();
  // rolesEnforcers$ = this.rolesService.listRoles().pipe(
  //   map((roles) =>
  //     roles.map((role) => ({
  //       role,
  //       enforcer: newEnforcer(role.permissions),
  //     })),
  //   ),
  // );
}
