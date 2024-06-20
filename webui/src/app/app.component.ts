import { Component } from "@angular/core";

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.scss"],
})
export class AppComponent {
  title = "bitmagnet";
  theme: Theme = "light";

  constructor() {
    if (window.matchMedia("(prefers-color-scheme)").media !== "not all") {
      const darkModeMediaQuery = window.matchMedia(
        "(prefers-color-scheme: dark)",
      );
      this.theme = darkModeMediaQuery.matches ? "dark" : "light";
      darkModeMediaQuery.addEventListener("change", (e) => {
        this.theme = e.matches ? "dark" : "light";
      });
    }
  }
}

type Theme = "light" | "dark";
