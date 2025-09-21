import {
  Component,
  ComponentRef,
  ViewContainerRef,
  ViewChild,
  inject,
  Input,
  OnInit,
  Output,
  EventEmitter,
} from "@angular/core";
import { AppModule } from "../../app.module";
import * as generated from "../../graphql/generated";
import { Enforcer, newEnforcer, ObjectAction } from "../../auth/enforcer";
import { RolesService } from "../../auth/roles.service";
import { map, Observable, take } from "rxjs";
import { MatCheckboxChange } from "@angular/material/checkbox";
import {
  ChangeEvent,
  Changes,
  Permission,
  PermissionsEditComponent,
} from "../../auth/permissions-edit.component";
import { objectActionKey } from "../../auth/util";

type ObjectActionPermission = ObjectAction & {
  key: string;
  permission: Permission | { subject: false; core: false };
};

@Component({
  selector: "app-role-edit",
  template: `
    <ng-container *transloco="let t">
      <mat-card>
        <mat-card-content>
          <app-permissions-edit
            [permissions]="(permissions$ | async) ?? []"
            (change)="this.changes = $event"
          />
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
  standalone: true,
  imports: [AppModule, PermissionsEditComponent],
})
export class RoleEditComponent implements OnInit {
  private rolesService = inject(RolesService);

  permissions$: Observable<Permission[]>;

  @Input() role: generated.Role;
  @Output() updated = new EventEmitter<void>();

  enforcer: Enforcer;

  changes?: ChangeEvent;

  ngOnInit(): void {
    this.enforcer = newEnforcer(this.role.permissions);
    this.permissions$ = this.rolesService.objectActions$.pipe(
      map((objActs) =>
        objActs.map((objectAction) => {
          const enf = this.enforcer(objectAction);

          return {
            key: objectActionKey(objectAction),
            objectAction,
            active: !!enf,
            core: enf?.core ?? false,
          };
        }),
      ),
    );
  }

  get hasChanges(): boolean {
    return Object.keys(this.changes?.changes ?? {}).length > 0;
  }

  save() {
    const changes = this.changes;
    if (!changes) {
      return;
    }

    this.rolesService
      .putRole({
        role: this.role.name,
        objectActions: changes.permissions.flatMap((permission) =>
          permission.active && !permission.core
            ? [permission.objectAction]
            : [],
        ),
      })
      .pipe(
        take(1),
        map(() => {
          this.updated.emit();
        }),
      )
      .subscribe();
  }
}
