import { Component, inject, Input, OnInit } from "@angular/core";
import { WorkersService } from "./workers.service";
import { MatDialogRef } from "@angular/material/dialog";
import { Action } from "./types";

@Component({
  selector: "app-workers-confirm-action-dialog",
  standalone: false,
  templateUrl: "./workers-confirm-action-dialog.component.html",
})
export class WorkersConfirmActionDialogComponent implements OnInit {
  private workers = inject(WorkersService);
  private dialogRef = inject(MatDialogRef<WorkersConfirmActionDialogComponent>);

  @Input() action: Action;
  @Input() worker: string;

  dependsOn = Array<string>();
  requiredBy = Array<string>();

  ngOnInit() {
    this.workers.result$.subscribe((result) => {
      const worker = result.workers.find((w) => w.ref === this.worker);
      if (worker) {
        this.dependsOn = worker.dependsOn;
        this.requiredBy = worker.requiredBy;
      }
    });
  }

  confirm() {
    switch (this.action) {
      case "start":
        this.workers.startWorkers(this.worker);
        break;
      case "shutdown":
        this.workers.shutdownWorkers(this.worker);
        break;
      case "restart":
        this.workers.restartWorkers(this.worker);
        break;
    }
    this.dialogRef.close();
  }

  cancel() {
    this.dialogRef.close();
  }
}
