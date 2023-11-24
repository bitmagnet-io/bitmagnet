import { Injectable } from "@angular/core";
import { MatSnackBar } from "@angular/material/snack-bar";

@Injectable()
export class AppErrorsService {
  public readonly expiry = 1000 * 10;

  constructor(private snackBar: MatSnackBar) {}

  addError(message: string, expiry = this.expiry) {
    this.snackBar.open(message, "Dismiss", {
      duration: expiry,
      panelClass: ["snack-bar-error"],
    });
  }
}
