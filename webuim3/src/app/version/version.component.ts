import { Component, OnInit } from '@angular/core';
import { MatTooltip } from '@angular/material/tooltip';
import { TranslocoDirective } from '@jsverse/transloco';
import * as generated from '../graphql/generated';
import { GraphQLService } from '../graphql/graphql.service';
import { GraphQLModule } from '../graphql/graphql.module';

const defaultVersionName = 'v-unknown';

@Component({
  selector: 'app-version',
  standalone: true,
  templateUrl: './version.component.html',
  imports: [GraphQLModule, MatTooltip, TranslocoDirective],
})
export class VersionComponent implements OnInit {
  version: string = defaultVersionName;
  versionUnknown = true;

  constructor(private graphQLService: GraphQLService) {}

  ngOnInit(): void {
    this.graphQLService.systemQuery().subscribe({
      next: (data: generated.SystemQuery) => {
        if (data.version) {
          this.version = data.version;
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
