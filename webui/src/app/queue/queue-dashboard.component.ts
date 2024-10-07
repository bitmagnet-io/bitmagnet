import { Component, inject, OnDestroy, OnInit } from '@angular/core';
import { MatCard, MatCardContent, MatCardHeader } from '@angular/material/card';
import { TranslocoDirective } from '@jsverse/transloco';
import { MatIcon } from '@angular/material/icon';
import { MatToolbar } from '@angular/material/toolbar';
import {
  ActivatedRoute,
  EventType,
  Router,
  RouterLink,
  RouterLinkActive,
  RouterOutlet,
} from '@angular/router';
import { EMPTY, Subscription } from 'rxjs';
import { MatAnchor } from '@angular/material/button';
import { MatTabLink, MatTabNav, MatTabNavPanel } from '@angular/material/tabs';

@Component({
  selector: 'app-queue-dashboard',
  standalone: true,
  imports: [
    MatCard,
    MatCardHeader,
    TranslocoDirective,
    MatCardContent,
    MatIcon,
    MatToolbar,
    RouterOutlet,
    RouterLink,
    MatAnchor,
    RouterLinkActive,
    MatTabNav,
    MatTabLink,
    MatTabNavPanel,
  ],
  templateUrl: './queue-dashboard.component.html',
  styleUrl: './queue-dashboard.component.scss',
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
          event.urlAfterRedirects === '/dashboard/queue'
        ) {
          this.redirectVisualize();
        }
        return EMPTY;
      }),
    );
  }

  private redirectVisualize(): void {
    void this.router.navigate(['visualize'], {
      relativeTo: this.route,
    });
  }

  ngOnDestroy() {
    this.subscriptions.forEach((subscription) => subscription.unsubscribe());
    this.subscriptions = new Array<Subscription>();
  }
}
