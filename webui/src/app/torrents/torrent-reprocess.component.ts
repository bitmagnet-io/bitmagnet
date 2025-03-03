import { Component, EventEmitter, inject, Input, Output } from "@angular/core";
import { Apollo } from "apollo-angular";
import { map } from "rxjs/operators";
import * as generated from "../graphql/generated";
import { AppModule } from "../app.module";
import { ErrorsService } from "../errors/errors.service";

@Component({
  selector: "app-torrent-reprocess",
  standalone: true,
  imports: [AppModule],
  templateUrl: "./torrent-reprocess.component.html",
  styleUrl: "./torrent-reprocess.component.scss",
})
export class TorrentReprocessComponent {
  @Input() infoHashes: string[];
  apollo = inject(Apollo);
  errors = inject(ErrorsService);

  classifierRematch = false;
  apisDisabled = true;
  localSearchDisabled = true;

  @Output() updated = new EventEmitter<null>();

  reprocess() {
    this.apollo
      .mutate<
        generated.TorrentReprocessMutation,
        generated.TorrentReprocessMutationVariables
      >({
        mutation: generated.TorrentReprocessDocument,
        variables: {
          input: {
            infoHashes: this.infoHashes,
            classifierRematch: this.classifierRematch,
            apisDisabled: this.apisDisabled,
            localSearchDisabled: this.localSearchDisabled,
          },
        },
      })
      .pipe(
        map(() => {
          this.updated.emit(null);
        }),
      )
      .subscribe();
  }
}
