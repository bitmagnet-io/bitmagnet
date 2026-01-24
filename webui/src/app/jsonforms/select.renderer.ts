import {
  ChangeDetectionStrategy,
  Component,
  Input,
  OnInit,
} from "@angular/core";
import { JsonFormsAngularService, JsonFormsControl } from "@jsonforms/angular";
import {
  isEnumControl,
  OwnPropsOfControl,
  RankedTester,
  rankWith,
} from "@jsonforms/core";

import { CommonModule } from "@angular/common";
import { ReactiveFormsModule } from "@angular/forms";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatSelectModule } from "@angular/material/select";

@Component({
  selector: "SelectControlRenderer",
  template: `
    <mat-form-field [ngStyle]="{ display: hidden ? 'none' : '' }">
      <mat-label>{{ label }}</mat-label>
      <mat-select
        [id]="id"
        [formControl]="form"
        (selectionChange)="onChange($event)"
      >
        @for (option of options || scopedSchema.enum || []; track option) {
          <mat-option [value]="option">
            {{ option }}
          </mat-option>
        }
      </mat-select>
      <mat-hint *ngIf="shouldShowUnfocusedDescription()">{{
        description
      }}</mat-hint>
      <mat-error>{{ error }}</mat-error>
    </mat-form-field>
  `,
  styles: [
    `
      :host {
        display: flex;
        flex-direction: row;
      }
      mat-form-field {
        flex: 1 1 auto;
      }
    `,
  ],
  // changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatSelectModule,
  ],
})
export class SelectControlRenderer extends JsonFormsControl implements OnInit {
  @Input() options: string[];

  constructor(jsonformsService: JsonFormsAngularService) {
    super(jsonformsService);
  }

  override ngOnInit() {
    super.ngOnInit();
    console.log(this.schema, this.scopedSchema, this.uischema);
  }

  override getEventValue = (event: any) => event.value;

  protected override getOwnProps(): OwnPropsOfAutoComplete {
    return {
      ...super.getOwnProps(),
      options: this.options,
    };
  }
}

export const enumControlTester: RankedTester = rankWith(2, isEnumControl);

interface OwnPropsOfAutoComplete extends OwnPropsOfControl {
  options: string[];
}
