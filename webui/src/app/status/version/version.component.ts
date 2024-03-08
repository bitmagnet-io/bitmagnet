import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { statusFetchEndpoint } from "../../../environments/environment";
import { Status } from '../status.type';

const defaultVersionName = 'N/A';

@Component({
  selector: 'app-version',
  standalone: false,
  templateUrl: './version.component.html',
  styleUrl: './version.component.scss'
})
export class VersionComponent implements OnInit {

  version: string = defaultVersionName;

  constructor(private http: HttpClient) { }

  ngOnInit(): void {
    this.http.get<Status>(statusFetchEndpoint).subscribe({
      next: this.parseVersionFromRsp.bind(this),
      error: () => {
        this.version = defaultVersionName;
      },
    });
  }

  parseVersionFromRsp(data: Status) {
    this.version = data?.component?.version || 'N/A';
  }
}
