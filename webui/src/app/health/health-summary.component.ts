import { Component, inject } from "@angular/core";
import { HealthService } from "./health.service";

@Component({
  selector: "app-health-summary",
  standalone: false,
  templateUrl: "./health-summary.component.html",
  styleUrl: "./health-summary.component.scss",
})
export class HealthSummaryComponent {
  health = inject(HealthService);
}
