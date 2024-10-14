import { Component, inject, OnDestroy, OnInit } from "@angular/core";
import { ActivatedRoute, EventType, Router } from "@angular/router";
import { EMPTY, Subscription } from "rxjs";
import { AppModule } from "../../app.module";

@Component({
  selector: "app-queue-dashboard",
  standalone: true,
  imports: [AppModule],
  templateUrl: "./queue-dashboard.component.html",
  styleUrl: "./queue-dashboard.component.scss",
})
export class QueueDashboardComponent implements OnInit, OnDestroy {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private subscriptions = new Array<Subscription>();

  ngOnInit() {
    this.subscriptions.push(
      this.route.url.subscribe(() => {
        if (!this.route.firstChild) {
          this.redirectVisualize();
        }
        return EMPTY;
      }),
      this.router.events.subscribe((event) => {
        if (
          event.type === EventType.NavigationEnd &&
          event.urlAfterRedirects === "/dashboard/queue"
        ) {
          this.redirectVisualize();
        }
        return EMPTY;
      }),
    );
  }

  private redirectVisualize(): void {
    void this.router.navigate(["visualize"], {
      relativeTo: this.route,
    });
  }

  ngOnDestroy() {
    this.subscriptions.forEach((subscription) => subscription.unsubscribe());
    this.subscriptions = new Array<Subscription>();
  }
}
