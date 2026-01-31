import { Component, inject, Input, OnInit } from "@angular/core";
import { AppModule } from "../app.module";
import { TargetsService } from "./targets.service";
import { JsonFormsAngularMaterialModule } from "@jsonforms/angular-material";
import * as generated from "../graphql/generated";
import { angularMaterialRenderers } from "../jsonforms/renderers";
import { UISchemaElement } from "@jsonforms/core";

@Component({
  selector: "app-torrents-send-target",
  standalone: true,
  imports: [AppModule, JsonFormsAngularMaterialModule],
  template: `
    <ng-container *transloco="let t">
      <mat-card>
        <mat-card-header>
          <mat-card-title
            ><h4>{{ target.name }}</h4></mat-card-title
          >
        </mat-card-header>
        @if (target.dataSchema) {
          <mat-card-content>
            <jsonforms
              [schema]="target.dataSchema"
              [uischema]="uischema"
              [renderers]="renderers"
              [data]="data"
              (dataChange)="setData($event)"
              (errors)="setErrors($event)"
            ></jsonforms>
          </mat-card-content>
        }
        <mat-card-actions>
          <button
            mat-raised-button
            color="primary"
            [disabled]="infoHashes.length === 0 || hasError"
            (click)="send()"
          >
            <mat-icon>send</mat-icon>
            {{ t("torrents.send.sendToTarget", { target: target.name }) }}
          </button>
        </mat-card-actions>
      </mat-card>
    </ng-container>
  `,
  styles: [
    `
      h4 {
        font-size: 16px;
        padding: 0;
        margin-top: 8px;
        margin-bottom: 12px;
      }
      mat-card-actions {
        padding-top: 12px;
        padding-bottom: 16px;
      }
    `,
  ],
})
export class TorrentsSendTargetComponent implements OnInit {
  private targetsService = inject(TargetsService);

  @Input() target: generated.TargetFragment;

  uischema: UISchemaElement;

  renderers = angularMaterialRenderers;

  @Input() index: string | null = null;
  @Input() infoHashes: string[];

  hasError: boolean = false;
  data: unknown = {};

  ngOnInit() {
    this.uischema =
      (this.target.uiSchema as UISchemaElement) || defaultUISchema;
  }

  setData(value: unknown) {
    this.data = value;
  }

  setErrors(errors: unknown[]) {
    this.hasError = errors.length > 0;
  }

  send() {
    this.targetsService.send({
      infoHashes: this.infoHashes,
      target: this.target.ref,
      data: this.data,
      index: this.index || undefined,
    });
  }
}

const defaultUISchema = {
  type: "VerticalLayout",
  elements: [
    {
      type: "Control",
      scope: "#/",
    },
  ],
};
