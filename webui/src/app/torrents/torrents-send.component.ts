import { Component, inject, Input } from "@angular/core";
import { AppModule } from "../app.module";
import { TargetsService } from "./targets.service";
import { TorrentsSendTargetComponent } from "./torrents-send-target.component";

@Component({
  selector: "app-torrents-send",
  standalone: true,
  imports: [AppModule, TorrentsSendTargetComponent],
  template: `
    <ng-container *transloco="let t">
      <mat-card>
        @for (target of targetsService.targets$ | async; track target.ref) {
          <app-torrents-send-target
            [index]="index"
            [target]="target"
            [infoHashes]="infoHashes"
          ></app-torrents-send-target>
        }
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
export class TorrentsSendComponent {
  targetsService = inject(TargetsService);

  @Input() index: string | null = null;
  @Input() infoHashes: string[];
}
