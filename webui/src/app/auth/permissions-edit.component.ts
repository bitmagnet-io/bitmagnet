import { Component, EventEmitter, Input, Output } from "@angular/core";
import { MatCheckboxChange } from "@angular/material/checkbox";
import { AppModule } from "../app.module";
import { BehaviorSubject } from "rxjs";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import * as generated from "../graphql/generated";

export type Permission = {
  key: string;
  objectAction: generated.AuthObjectAction;
  active: boolean;
  core: boolean;
};

export type Changes = Record<
  string,
  { objectAction: generated.AuthObjectAction; active: boolean }
>;

export type ChangeEvent = {
  changes: Changes;
  permissions: Permission[];
};

@Component({
  selector: "app-permissions-edit",
  template: `
    <ul>
      @for (perm of permissions; track perm.key) {
        <li>
          <mat-checkbox
            [checked]="perm.active"
            (change)="onChecked($event, perm)"
            [disabled]="perm.core"
          >
            {{ perm.key }}
          </mat-checkbox>
        </li>
      }
    </ul>
  `,
  styles: [
    `
      ul {
        list-style-type: none;
        padding-bottom: 20px;
        padding-left: 0;
      }
    `,
  ],
  standalone: true,
  imports: [AppModule],
})
export class PermissionsEditComponent {
  private changes = new BehaviorSubject<Changes>({});

  @Input() permissions: Permission[];
  @Output() change = new EventEmitter<ChangeEvent>();

  constructor() {
    this.changes
      .asObservable()
      .pipe(takeUntilDestroyed())
      .subscribe((changes) =>
        this.change.emit({
          changes,
          permissions: this.permissions.map((p) =>
            p.key in changes
              ? {
                  key: p.key,
                  active: changes[p.key].active,
                  core: p.core ?? false,
                  objectAction: p.objectAction,
                }
              : p,
          ),
        }),
      );
  }

  onChecked(event: MatCheckboxChange, perm: Permission) {
    const changes = { ...this.changes.getValue() };

    if (event.checked) {
      if (perm.active) {
        delete changes[perm.key];
      } else {
        changes[perm.key] = { ...perm, active: true };
      }
    } else {
      if (perm.active) {
        changes[perm.key] = { ...perm, active: false };
      } else {
        delete changes[perm.key];
      }
    }

    this.changes.next(changes);
  }
}
