import { Component, inject } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { HealthService } from "./health.service";

@Component({
  selector: "app-health-widget",
  standalone: false,
  templateUrl: "./health-widget.component.html",
  styleUrl: "./health-widget.component.scss",
})
export class HealthWidgetComponent {
  health = inject(HealthService);
  dialog = inject(MatDialog);
}
