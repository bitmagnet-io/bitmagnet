import {Component, inject, Input} from "@angular/core";
import {WorkersService} from "./workers.service";
import {MatDialog} from "@angular/material/dialog";
import {Action} from "./types";
import {WorkersConfirmActionDialogComponent} from "./workers-confirm-action-dialog.component";

@Component({
  selector: "app-workers-table",
  standalone: false,
  templateUrl: "./workers-table.component.html",
  styleUrl: "./workers-table.component.scss",
})
export class WorkersTableComponent {
  readonly workers = inject(WorkersService);
  readonly dialog = inject(MatDialog);

  @Input() actions: boolean = false;

  confirm(action: Action, worker: string) {
    const ref = this.dialog.open(WorkersConfirmActionDialogComponent)
    ref.componentInstance.action = action;
    ref.componentInstance.worker = worker;
  }
}
