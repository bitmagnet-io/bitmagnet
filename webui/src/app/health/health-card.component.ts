import { Component, inject } from "@angular/core";
import { HealthService } from "./health.service";

@Component({
  selector: "app-health-card",
  standalone: false,
  templateUrl: "./health-card.component.html",
  styleUrl: "./health-card.component.scss",
})
export class HealthCardComponent {
  health = inject(HealthService);
}
