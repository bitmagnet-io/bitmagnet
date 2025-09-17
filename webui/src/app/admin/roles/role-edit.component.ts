import { Component, inject, Input, OnInit } from "@angular/core";
import { AppModule } from "../../app.module";
import * as generated from "../../graphql/generated";
import { Enforcer, newEnforcer, ObjectAction } from "../../auth/enforcer";
import { RolesService } from "../../auth/roles.service";
import { map, Observable, take } from "rxjs";
import { MatCheckboxChange } from "@angular/material/checkbox";

type ObjectActionPermission = ObjectAction & {
  key: string;
  permission: generated.Permission | { subject: false; core: false };
};

@Component({
  selector: "app-role-edit",
  template: `
    <ng-container *transloco="let t">
      <mat-card>
        <mat-card-content>
          <ul>
            @for (objAct of objectActions$ | async; track objAct) {
              <li>
                <mat-checkbox
                  [checked]="objAct.permission.subject"
                  (change)="onChecked($event, objAct)"
                  [disabled]="objAct.permission.core"
                >
                  {{ objAct.key }}
                </mat-checkbox>
              </li>
            }
          </ul>
        </mat-card-content>
        <mat-card-actions>
          <button mat-stroked-button [disabled]="!hasChanges" (click)="save()">
            <mat-icon>save</mat-icon>
            {{ t("general.save") }}
          </button>
        </mat-card-actions>
      </mat-card>
    </ng-container>
  `,
  styles: [
    `
      ul {
        list-style-type: none;
        padding-bottom: 20px;
      }
    `,
  ],
  standalone: true,
  imports: [AppModule],
})
export class RoleEditComponent implements OnInit {
  private rolesService = inject(RolesService);
  private changes: Record<string, ObjectAction & { checked: boolean }> = {};

  objectActions$: Observable<ObjectActionPermission[]>;

  @Input() role: generated.Role;

  enforcer: Enforcer;

  ngOnInit(): void {
    this.enforcer = newEnforcer(this.role.permissions);
    this.objectActions$ = this.rolesService.objectActions$.pipe(
      map((objActs) =>
        objActs.map((objAct) => ({
          key: `${objAct.namespace}::${objAct.object}::${objAct.action}`,
          ...objAct,
          permission: this.enforcer(objAct) || {
            subject: false,
            core: false,
          },
        })),
      ),
    );
  }

  onChecked(event: MatCheckboxChange, perm: ObjectActionPermission) {
    if (event.checked) {
      if (perm.permission.subject) {
        delete this.changes[perm.key];
      } else {
        this.changes[perm.key] = { ...perm, checked: true };
      }
    } else {
      if (perm.permission.subject) {
        this.changes[perm.key] = { ...perm, checked: false };
      } else {
        delete this.changes[perm.key];
      }
    }
  }

  get hasChanges(): boolean {
    return Object.keys(this.changes).length > 0;
  }

  save() {
    this.rolesService
      .putRole({
        role: this.role.name,
        objectActions: Object.values(this.changes).flatMap((change) =>
          change.checked
            ? [
                {
                  namespace: change.namespace as string,
                  object: change.object,
                  action: change.action,
                },
              ]
            : [],
        ),
      })
      .pipe(
        take(1),
        map(() => void 0),
      )
      .subscribe();
  }
}
