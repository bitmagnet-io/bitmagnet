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
import { MatToolbar } from '@angular/material/toolbar';
import { MatMenu, MatMenuItem, MatMenuTrigger } from '@angular/material/menu';
import { GraphQLModule } from '../graphql/graphql.module';
import { ChartComponent } from '../charting/chart.component';
import { BreakpointsService } from '../layout/breakpoints.service';
import { ErrorsService } from '../errors/errors.service';
import {
  autoRefreshIntervalNames,
  eventNames,
  resolutionNames,
  timeframeNames,
} from './torrent-metrics.constants';
import {TorrentMetricsController} from "./torrent-metrics.controller";
import {TorrentChartAdapterTimeline} from "./torrent-chart-adapter.timeline";
// import { QueueChartAdapterTotals } from './queue-chart-adapter.totals';
// import { QueueMetricsController } from './queue-metrics.controller';
// import { QueueChartAdapterTimeline } from './queue-chart-adapter.timeline';

@Component({
  selector: 'app-torrent-metrics',
  standalone: true,
  templateUrl: './torrent-metrics.component.html',
  styleUrl: './torrent-metrics.component.scss',
  imports: [
    MatCardContent,
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
export class TorrentMetricsComponent implements OnDestroy {
  breakpoints = inject(BreakpointsService);
  private apollo = inject(Apollo);
  torrentMetricsController = new TorrentMetricsController(
    this.apollo,
    {
      buckets: {
        duration: 'AUTO',
        multiplier: 'AUTO',
        timeframe: 'days_1',
      },
      autoRefresh: 'seconds_30',
    },
    inject(ErrorsService),
  );
  protected readonly timeline = inject(TorrentChartAdapterTimeline);

  protected readonly resolutionNames = resolutionNames;
  protected readonly timeframeNames = timeframeNames;
  protected readonly autoRefreshIntervalNames = autoRefreshIntervalNames;

  ngOnDestroy() {
    this.torrentMetricsController.setAutoRefreshInterval('off');
  }

  protected readonly eventNames = eventNames;

  handleMultiplierEvent(event: Event) {
    const value = (event.currentTarget as HTMLInputElement).value;
    this.torrentMetricsController.setBucketMultiplier(
      /^\d+$/.test(value) ? parseInt(value) : 'AUTO',
    );
  }
}
