import {Component, inject} from "@angular/core";
import { AppModule } from "../../app.module";
import { DocumentTitleComponent } from "../../layout/document-title.component";
import {BreakpointsService} from "../../layout/breakpoints.service";
import {WorkersModule} from "../../workers/workers.module";

@Component({
  selector: "app-workers-dashboard",
  templateUrl: "./dashboard-workers.component.html",
  styleUrl: "./dashboard-workers.component.scss",
  standalone: true,
  imports: [AppModule, DocumentTitleComponent, WorkersModule],
})
export class DashboardWorkersComponent {
  breakpoints = inject(BreakpointsService);
}
