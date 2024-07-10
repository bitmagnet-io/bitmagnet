import { Component, inject, OnInit } from '@angular/core';
import { MatTooltip } from '@angular/material/tooltip';
import { TranslocoDirective } from '@jsverse/transloco';
import { Apollo } from 'apollo-angular';
import { map } from 'rxjs';
import * as generated from '../graphql/generated';
import { GraphQLModule } from '../graphql/graphql.module';

const defaultVersionName = 'v-unknown';

@Component({
  selector: 'app-version',
  standalone: true,
  templateUrl: './version.component.html',
  imports: [GraphQLModule, MatTooltip, TranslocoDirective],
})
export class VersionComponent implements OnInit {
  private apollo = inject(Apollo);

  version: string = defaultVersionName;
  versionUnknown = true;

  ngOnInit(): void {
    this.apollo
      .query<generated.VersionQuery, generated.VersionQueryVariables>({
        query: generated.VersionDocument,
      })
      .pipe(map((r) => r.data.version))
      .subscribe({
        next: (version: string) => {
          if (version) {
            this.version = version;
            this.versionUnknown = false;
          } else {
            this.version = defaultVersionName;
            this.versionUnknown = true;
          }
        },
        error: () => {
          this.version = defaultVersionName;
        },
      });
  }
}
