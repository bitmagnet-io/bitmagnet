import { Component, inject } from "@angular/core";
import { JsonFormsAngularMaterialModule } from "@jsonforms/angular-material";
import { map } from "rxjs";
import { DocumentTitleComponent } from "../../layout/document-title.component";
import { AppModule } from "../../app.module";
import { ConfigService } from "../../config/config.service";
import { ConfigFormComponent } from "./config-form.component";

@Component({
  selector: "app-config-admin",
  template: `
    <ng-container *transloco="let t">
      <app-document-title [parts]="[t('routes.config')]" />
      <mat-card class="admin-card">
        <mat-card-header>
          <mat-toolbar>
            <h2><mat-icon>settings</mat-icon>{{ t("routes.config") }}</h2>
          </mat-toolbar>
        </mat-card-header>
        <mat-card-content>
          <mat-accordion>
            @for (plugin of pluginRefs$ | async; track plugin.ref) {
              <mat-expansion-panel hideToggle>
                <mat-expansion-panel-header>
                  <mat-panel-title>
                    {{ plugin.ref }}
                  </mat-panel-title>
                  <mat-panel-description>
                    {{ plugin.description }}
                  </mat-panel-description>
                </mat-expansion-panel-header>
                <mat-divider />
                @for (ref of plugin.params; track ref) {
                  <app-config-form [ref]="ref"> </app-config-form>
                }
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
  imports: [
    AppModule,
    DocumentTitleComponent,
    JsonFormsAngularMaterialModule,
    ConfigFormComponent,
  ],
})
export class AdminConfigComponent {
  private configService = inject(ConfigService);

  pluginRefs$ = this.configService.config$.pipe(
    map((state) =>
      Object.values(state.plugins).flatMap((plugin) => {
        const params = Object.values(state.params).filter(
          (p) => p.plugin === plugin.ref,
        );
        if (!params.length) {
          return [];
        }
        return [
          {
            ref: plugin.ref,
            description: plugin.description,
            params: params
              .sort((p) => (p.ref.endsWith(".activation") ? -1 : 0))
              .map((p) => p.ref),
          },
        ];
      }),
    ),
  );
}
