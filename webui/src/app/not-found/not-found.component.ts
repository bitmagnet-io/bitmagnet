import { Component } from "@angular/core";
import { AppModule } from "../app.module";

@Component({
  selector: "app-not-found",
  standalone: true,
  imports: [AppModule],
  templateUrl: "./not-found.component.html",
  styleUrl: "./not-found.component.scss",
})
export class NotFoundComponent {}
