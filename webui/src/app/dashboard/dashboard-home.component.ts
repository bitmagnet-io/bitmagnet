import {Component, inject} from "@angular/core";
import { HealthModule } from "../health/health.module";
import { AppModule } from "../app.module";
import { DocumentTitleComponent } from "../layout/document-title.component";
import {BreakpointsService} from "../layout/breakpoints.service";

@Component({
  selector: "app-dashboard",
  templateUrl: "./dashboard-home.component.html",
  styleUrl: "./dashboard-home.component.scss",
  standalone: true,
  imports: [AppModule, HealthModule, DocumentTitleComponent],
})
export class DashboardHomeComponent {
  breakpoints = inject(BreakpointsService);
}
