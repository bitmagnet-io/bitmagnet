import { Component, Input, inject, OnInit, OnDestroy } from "@angular/core";
import {
  angularMaterialRenderers,
  JsonFormsAngularMaterialModule,
} from "@jsonforms/angular-material";
import { Subscription } from "rxjs";
import { JsonSchema, UISchemaElement } from "@jsonforms/core";
import { ConfigParam } from "../../graphql/generated";
import { AppModule } from "../../app.module";
import { ConfigService } from "../../config/config.service";
import * as generated from "../../graphql/generated";

@Component({
  selector: "app-config-form",
  template: `
    <div class="form-container">
      <div class="form-content">
        <mat-label
          >{{ param.description }}
          <small>default: {{ defaultLabel }}</small></mat-label
        >
        <jsonforms
          [(data)]="data"
          [schema]="jsonSchema"
          [uischema]="uischema"
          [renderers]="renderers"
          [config]="{ useGrouping: false }"
          (dataChange)="setData($event)"
          (errors)="setErrors($event)"
        ></jsonforms>
      </div>
      <div class="form-actions">
        <button
          mat-mini-fab
          color="primary"
          [matTooltip]="'Save'"
          [disabled]="!changed || hasError"
          (click)="save()"
        >
          <mat-icon>save</mat-icon>
        </button>
        <button
          mat-mini-fab
          [matTooltip]="'Reset'"
          [disabled]="!changed"
          (click)="reset()"
        >
          <mat-icon>undo</mat-icon>
        </button>
        <button
          mat-mini-fab
          color="warning"
          [matTooltip]="'Delete persisted value'"
          [disabled]="!persisted"
          (click)="delete()"
        >
          <mat-icon>delete</mat-icon>
        </button>
      </div>
    </div>
  `,
  styles: [
    `
      .form-container {
        display: flex;
        align-items: flex-start;
        gap: 16px;
        margin-bottom: 24px;
      }

      .form-content {
        flex: 1;
      }

      .form-actions {
        display: flex;
        flex-direction: row;
        gap: 8px;
        margin-top: 36px; /* Align with form fields */
      }

      mat-label {
        display: block;
        margin-bottom: 10px;
        small {
          margin-left: 10px;
        }
      }

      ::ng-deep mat-form-field {
        --mat-form-field-filled-with-label-container-padding-top: 16px;
        mat-label {
          display: none;
        }
      }
    `,
  ],
  standalone: true,
  imports: [AppModule, JsonFormsAngularMaterialModule],
})
export class ConfigFormComponent implements OnInit, OnDestroy {
  @Input() ref: string;

  readonly renderers = angularMaterialRenderers;

  configService = inject(ConfigService);
  private subscription: Subscription;

  param: ConfigParam;
  jsonSchema: JsonSchema;
  uischema: UISchemaElement;

  serverData: unknown;
  data: unknown;
  changed: boolean;
  persisted: boolean;
  hasError: boolean = false;

  ngOnInit(): void {
    this.uischema = {
      type: "VerticalLayout",
      elements: [
        {
          type: "Control",
          scope: "#/",
          label: this.ref,
        },
      ],
    };
    this.subscription = this.configService.config$.subscribe((state) => {
      const param = state.params[this.ref];
      if (JSON.stringify(param) === JSON.stringify(this.param)) {
        return;
      }
      this.param = param;
      this.jsonSchema = transformJSONSchema(this.param.jsonSchema);
      if (this.serverData !== param.value) {
        this.serverData = param.value;
        this.data = this.serverData;
        this.changed = false;
      }
      this.persisted = ["dynamic", "persisted", "pending"].includes(
        param.source,
      );
    });
  }

  ngOnDestroy(): void {
    this.subscription.unsubscribe();
  }

  setData(value: unknown) {
    this.data = value;
    this.changed = this.data !== this.param.value;
  }

  save() {
    this.configService.save(this.ref, this.data);
  }

  reset() {
    this.data = this.param.value;
    this.changed = false;
  }

  delete() {
    this.configService.delete(this.ref);
  }

  setErrors(errors: unknown[]) {
    this.hasError = errors.length > 0;
  }

  get defaultLabel(): string {
    let label = "";
    switch (this.param.jsonSchema.type) {
      case "string":
      case "number":
      case "integer":
      case "boolean":
        label = `${this.param.default as string}`;
    }

    if (!label) {
      label = "[empty]";
    }
    return label;
  }
}

const transformJSONSchema = (s: generated.JsonSchemaFragment): JsonSchema =>
  ({
    type: s.type,
    default: s.default ?? undefined,
    enum: s.enum ?? undefined,
    pattern: s.pattern ?? undefined,
    multipleOf: s.multipleOf ?? undefined,
    maximum: s.maximum ?? undefined,
    exclusiveMaximum: s.exclusiveMaximum ?? undefined,
    minimum: s.minimum ?? undefined,
    exclusiveMinimum: s.exclusiveMinimum ?? undefined,
    maxLength: s.maxLength ?? undefined,
    minItems: s.minItems ?? undefined,
    maxItems: s.maxItems ?? undefined,
    uniqueItems: s.uniqueItems ?? undefined,
    items: s.items ? transformJSONSchema(s.items) : undefined,
  }) as JsonSchema;
