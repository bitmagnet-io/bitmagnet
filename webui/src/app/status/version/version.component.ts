import { Component, OnInit } from "@angular/core";
import * as generated from "../../graphql/generated";
import { GraphQLService } from "src/app/graphql/graphql.service";

const defaultVersionName = "";

@Component({
  selector: "app-version",
  standalone: false,
  templateUrl: "./version.component.html",
  styleUrl: "./version.component.scss",
})
export class VersionComponent implements OnInit {
  version: string = defaultVersionName;

  constructor(private graphQLService: GraphQLService) {}

  ngOnInit(): void {
    this.graphQLService.systemQuery().subscribe({
      next: this.parseVersionFromRsp.bind(this),
      error: () => {
        this.version = defaultVersionName;
      },
    });
  }

  parseVersionFromRsp(data: generated.SystemQuery) {
    this.version = data?.version || defaultVersionName;
  }
}
