import { Component } from "@angular/core";
import { HealthModule } from "../health/health.module";
import { AppModule } from "../app.module";

@Component({
  selector: "app-dashboard",
  templateUrl: "./dashboard-home.component.html",
  styleUrl: "./dashboard-home.component.scss",
  standalone: true,
  imports: [AppModule, HealthModule],
})
export class DashboardHomeComponent {}
