import { Component } from "@angular/core";
import { AppModule } from "../../app.module";

@Component({
  selector: "app-queue-dashboard",
  standalone: true,
  imports: [AppModule],
  templateUrl: "./queue-dashboard.component.html",
  styleUrl: "./queue-dashboard.component.scss",
})
export class QueueDashboardComponent {}
