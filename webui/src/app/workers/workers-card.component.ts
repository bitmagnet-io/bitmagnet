import { Component, inject } from "@angular/core";
import { WorkersService } from "./workers.service";

@Component({
  selector: "app-workers-card",
  standalone: false,
  templateUrl: "./workers-card.component.html",
  styleUrl: "./workers-card.component.scss",
})
export class WorkersCardComponent {
  workers = inject(WorkersService);
}
