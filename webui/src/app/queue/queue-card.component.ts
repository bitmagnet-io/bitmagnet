import { Component, inject, OnDestroy, OnInit } from '@angular/core';
import { Apollo } from 'apollo-angular';
import {
  MatCard,
  MatCardActions,
  MatCardContent,
  MatCardFooter,
  MatCardHeader,
  MatCardTitle,
} from '@angular/material/card';
import { MatIcon } from '@angular/material/icon';
import { TranslocoDirective } from '@jsverse/transloco';
import { MatFormField, MatLabel } from '@angular/material/form-field';
import { MatInput } from '@angular/material/input';
import { MatOption, MatSelect } from '@angular/material/select';
import { MatRadioButton, MatRadioGroup } from '@angular/material/radio';
import { AsyncPipe } from '@angular/common';
import {
  MatAnchor,
  MatButton, MatIconAnchor,
  MatIconButton,
  MatMiniFabButton,
} from '@angular/material/button';
import { MatTooltip } from '@angular/material/tooltip';
import { MatSlider, MatSliderThumb } from '@angular/material/slider';
import { MatGridList, MatGridTile } from '@angular/material/grid-list';
import { MatProgressBar } from '@angular/material/progress-bar';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatToolbar } from '@angular/material/toolbar';
import { MatMenu, MatMenuItem, MatMenuTrigger } from '@angular/material/menu';
import { GraphQLModule } from '../graphql/graphql.module';
import { ChartComponent } from '../charting/chart.component';
import { BreakpointsService } from '../layout/breakpoints.service';
import { ErrorsService } from '../errors/errors.service';
import { QueuePurgeJobsDialog } from './queue-purge-jobs-dialog.component';
import {
  autoRefreshIntervalNames,
  availableQueueNames,
  eventNames,
  resolutionNames,
  timeframeNames,
} from './queue.constants';
import { QueueEnqueueReprocessTorrentsBatchDialog } from './queue-enqueue-reprocess-torrents-batch-dialog.component';
import { QueueModule } from './queue.module';
import { QueueChartAdapterTotals } from './queue-chart-adapter.totals';
import { QueueMetricsController } from './queue-metrics.controller';
import { QueueChartAdapterTimeline } from './queue-chart-adapter.timeline';

@Component({
  selector: 'app-queue-card',
  standalone: true,
  templateUrl: './queue-card.component.html',
  styleUrl: './queue-card.component.scss',
  imports: [
    QueueModule,
    MatCardContent,
    MatDialogModule,
    ChartComponent,
    MatIcon,
    MatCardTitle,
    MatCardHeader,
    MatCard,
    GraphQLModule,
    TranslocoDirective,
    MatFormField,
    MatLabel,
    MatInput,
    MatSelect,
    MatOption,
    MatRadioGroup,
    MatRadioButton,
    AsyncPipe,
    MatMiniFabButton,
    MatTooltip,
    MatSlider,
    MatSliderThumb,
    MatIconButton,
    MatGridList,
    MatGridTile,
    MatCardFooter,
    MatCardActions,
    MatProgressBar,
    MatButton,
    MatToolbar,
    MatMenu,
    MatMenuItem,
    MatMenuTrigger,
    MatAnchor,
    MatIconAnchor,
  ],
})
export class QueueCardComponent implements OnInit, OnDestroy {
  breakpoints = inject(BreakpointsService);
  private apollo = inject(Apollo);
  readonly dialog = inject(MatDialog);
  queueMetricsController = new QueueMetricsController(
    this.apollo,
    {
      buckets: {
        duration: 'AUTO',
        multiplier: 'AUTO',
        timeframe: 'all',
      },
      autoRefresh: 'seconds_30',
    },
    inject(ErrorsService),
  );
  protected readonly timeline = inject(QueueChartAdapterTimeline);
  protected readonly totals = inject(QueueChartAdapterTotals);

  protected readonly resolutionNames = resolutionNames;
  protected readonly timeframeNames = timeframeNames;
  protected readonly availableQueueNames = availableQueueNames;
  protected readonly autoRefreshIntervalNames = autoRefreshIntervalNames;

  ngOnInit() {
    this.dialog.afterAllClosed.subscribe(() => {
      this.queueMetricsController.refresh();
    });
    this.queueMetricsController.result$.subscribe((result) => {
      // change the default settings to more sensible ones if there is <12 hours of data to show
      if (
        this.queueMetricsController.params.buckets.timeframe === 'all' &&
        this.queueMetricsController.params.buckets.duration === 'AUTO' &&
        result.params.buckets.duration === 'hour'
      ) {
        const span = result.bucketSpan;
        if (span && span.latestBucket - span.earliestBucket < 12) {
          this.queueMetricsController.setBucketDuration('minute');
        }
      }
    });
  }

  ngOnDestroy() {
    this.queueMetricsController.setAutoRefreshInterval('off');
  }

  openDialogPurgeJobs() {
    this.dialog.open(QueuePurgeJobsDialog, {
      data: {
        onPurged: () => {
          this.queueMetricsController.refresh();
        },
      },
    });
  }

  openDialogEnqueueReprocessTorrentsBatch() {
    this.dialog.open(QueueEnqueueReprocessTorrentsBatchDialog, {
      data: {
        onEnqueued: () => {
          this.queueMetricsController.refresh();
        },
      },
    });
  }

  protected readonly eventNames = eventNames;

  handleMultiplierEvent(event: Event) {
    const value = (event.currentTarget as HTMLInputElement).value;
    this.queueMetricsController.setBucketMultiplier(
      /^\d+$/.test(value) ? parseInt(value) : 'AUTO',
    );
  }
}
