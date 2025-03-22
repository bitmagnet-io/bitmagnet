import {
  ChangeDetectionStrategy,
  Component,
  EventEmitter,
  Input,
  OnChanges,
  OnDestroy,
  OnInit,
  Output,
  SimpleChanges,
} from "@angular/core";
import { CommonModule, JsonPipe } from "@angular/common";
import { FormControl, FormGroup, ReactiveFormsModule } from "@angular/forms";
import { MatButtonModule } from "@angular/material/button";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatSelectModule } from "@angular/material/select";
import { MatChipsModule } from "@angular/material/chips";
import { MatIconModule } from "@angular/material/icon";
import { MatTooltipModule } from "@angular/material/tooltip";
import { MatDatepickerModule } from "@angular/material/datepicker";
import {
  MatNativeDateModule,
  provideNativeDateAdapter,
  MAT_DATE_FORMATS,
} from "@angular/material/core";
import { MatTabsModule } from "@angular/material/tabs";
import { MatDividerModule } from "@angular/material/divider";
import { TranslocoModule } from "@jsverse/transloco";
import {
  Subject,
  Subscription,
  debounceTime,
  distinctUntilChanged,
} from "rxjs";
import {
  TimeFrame,
  formatTimeFrameDescription,
  parseTimeFrame,
  timeFramePresets,
} from "./parse-timeframe";

// Custom date format to match our app's style
export const MY_DATE_FORMATS = {
  parse: {
    dateInput: "MM/DD/YYYY",
  },
  display: {
    dateInput: "MMM D, YYYY",
    monthYearLabel: "MMM YYYY",
    dateA11yLabel: "LL",
    monthYearA11yLabel: "MMMM YYYY",
  },
};

@Component({
  selector: "app-time-frame-selector",
  standalone: true,
  imports: [
    CommonModule,
    JsonPipe,
    ReactiveFormsModule,
    MatButtonModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatChipsModule,
    MatIconModule,
    MatTooltipModule,
    MatDatepickerModule,
    MatNativeDateModule,
    MatTabsModule,
    MatDividerModule,
    TranslocoModule,
  ],
  providers: [
    provideNativeDateAdapter(),
    { provide: MAT_DATE_FORMATS, useValue: MY_DATE_FORMATS },
  ],
  template: `
    <ng-container *transloco="let t">
      <div class="time-frame-container">
        <!-- Tabbed interface for filtering options -->
        <mat-tab-group
          animationDuration="0ms"
          (selectedTabChange)="onTabChange($event)"
          [selectedIndex]="activeTabIndex"
          class="filter-tabs"
        >
          <!-- Quick Presets Tab -->
          <mat-tab [label]="t('dates.quick_presets')">
            <div class="preset-buttons">
              <button
                *ngFor="let preset of commonPresets"
                mat-stroked-button
                [class.selected]="timeFrameControl.value === preset.value"
                (click)="onPresetSelected({ value: preset.value })"
              >
                {{ preset.label }}
              </button>
            </div>

            <mat-divider></mat-divider>

            <div class="more-presets-section">
              <mat-form-field appearance="outline" class="full-width">
                <mat-label>{{ t("dates.more_presets") }}</mat-label>
                <mat-select (selectionChange)="onPresetSelected($event)">
                  <mat-option
                    *ngFor="let preset of morePresets"
                    [value]="preset.value"
                  >
                    {{ preset.label }}
                  </mat-option>
                </mat-select>
              </mat-form-field>
            </div>
          </mat-tab>

          <!-- Calendar Tab -->
          <mat-tab [label]="t('dates.date_range')">
            <div class="date-range-container">
              <mat-form-field appearance="outline" class="date-range-field">
                <mat-label>{{ t("dates.select_date_range") }}</mat-label>
                <mat-date-range-input
                  [formGroup]="dateRange"
                  [rangePicker]="rangePicker"
                >
                  <input
                    matStartDate
                    formControlName="start"
                    [placeholder]="t('dates.start_date')"
                  />
                  <input
                    matEndDate
                    formControlName="end"
                    [placeholder]="t('dates.end_date')"
                  />
                </mat-date-range-input>
                <mat-hint>{{ t("dates.date_range_hint") }}</mat-hint>
                <mat-datepicker-toggle
                  matIconSuffix
                  [for]="rangePicker"
                ></mat-datepicker-toggle>
                <mat-date-range-picker #rangePicker></mat-date-range-picker>

                <mat-error
                  *ngIf="
                    dateRange.controls.start.hasError('matStartDateInvalid')
                  "
                >
                  {{ t("dates.invalid_start_date") }}
                </mat-error>
                <mat-error
                  *ngIf="dateRange.controls.end.hasError('matEndDateInvalid')"
                >
                  {{ t("dates.invalid_end_date") }}
                </mat-error>
              </mat-form-field>

              <button
                mat-raised-button
                color="primary"
                class="apply-button"
                [disabled]="
                  !dateRange.valid ||
                  !dateRange.value.start ||
                  !dateRange.value.end
                "
                (click)="applyDateRangePicker()"
              >
                {{ t("general.apply") }}
              </button>
            </div>
          </mat-tab>

          <!-- Custom Expression Tab -->
          <mat-tab [label]="t('dates.custom')">
            <div class="custom-expression">
              <mat-form-field appearance="outline" class="full-width">
                <mat-label>{{ t("dates.time_frame") }}</mat-label>
                <input
                  matInput
                  type="text"
                  [formControl]="timeFrameControl"
                  [placeholder]="t('dates.time_frame_placeholder')"
                  (keyup.enter)="updateTimeFrameAndEmit(timeFrameControl.value)"
                />
                <mat-hint>{{ t("dates.custom_time_hint") }}</mat-hint>
              </mat-form-field>

              <button
                mat-raised-button
                color="primary"
                class="apply-button"
                (click)="updateTimeFrameAndEmit(timeFrameControl.value)"
              >
                {{ t("general.apply") }}
              </button>
            </div>
          </mat-tab>
        </mat-tab-group>

        <!-- Active time frame display -->
        <div
          class="selected-time-frame"
          *ngIf="
            currentTimeFrame &&
            currentTimeFrame.isValid &&
            timeFrameControl.value
          "
        >
          <mat-chip highlighted color="primary" class="range-chip">
            <mat-icon>event</mat-icon>
            <span class="time-frame-value">{{
              formatTimeFrameDescription(currentTimeFrame)
            }}</span>
            <button matChipRemove (click)="clearTimeFrame()">
              <mat-icon>cancel</mat-icon>
            </button>
          </mat-chip>
        </div>

        <div
          class="error-message"
          *ngIf="currentTimeFrame && !currentTimeFrame.isValid"
        >
          <mat-chip highlighted color="warn">
            {{ t("dates.time_frame_error") }}: {{ currentTimeFrame.error }}
          </mat-chip>
        </div>
      </div>
    </ng-container>
  `,
  styles: [
    `
      .time-frame-container {
        display: flex;
        flex-direction: column;
        gap: 16px;
      }

      .filter-tabs {
        border-radius: 8px;
        overflow: hidden;
      }

      .preset-buttons {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
        margin: 16px 0;
      }

      .preset-buttons button {
        margin: 4px;
        min-width: 110px;
      }

      .preset-buttons button.selected {
        background-color: rgba(0, 0, 0, 0.08);
        font-weight: 500;
      }

      .more-presets-section {
        margin: 16px 0;
      }

      .full-width {
        width: 100%;
      }

      .date-range-container {
        padding: 16px 0;
        display: flex;
        flex-direction: column;
        gap: 16px;
      }

      .date-range-field {
        width: 100%;
      }

      .custom-expression {
        padding: 16px 0;
        display: flex;
        flex-direction: column;
        gap: 16px;
      }

      .apply-button {
        align-self: flex-start;
      }

      .selected-time-frame {
        margin-top: 8px;
      }

      .range-chip {
        padding: 8px 12px;
        gap: 8px;
      }

      .time-frame-value {
        font-weight: normal;
      }

      .error-message {
        margin-top: 8px;
      }

      .error-message mat-chip {
        padding: 8px 12px;
      }

      @media (max-width: 599px) {
        .preset-buttons {
          flex-direction: column;
        }

        .preset-buttons button {
          width: 100%;
        }
      }
    `,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TimeFrameSelectorComponent
  implements OnInit, OnDestroy, OnChanges
{
  @Input() initialTimeFrame: string = "";
  @Output() timeFrameChanged = new EventEmitter<TimeFrame>();

  timeFrameControl = new FormControl<string>("");

  // Date range form group
  dateRange = new FormGroup({
    start: new FormControl<Date | null>(null),
    end: new FormControl<Date | null>(null),
  });

  activeTabIndex = 0;
  currentTimeFrame: TimeFrame | null = null;

  // Split presets into common and more categories
  commonPresets = timeFramePresets.slice(0, 6); // First 6 presets for quick access
  morePresets = timeFramePresets.slice(6); // Remaining presets in dropdown

  private _manualChange = new Subject<string>();
  private _subscriptions: Subscription[] = [];

  ngOnInit(): void {
    if (this.initialTimeFrame) {
      this.timeFrameControl.setValue(this.initialTimeFrame, {
        emitEvent: false,
      });
      this.updateTimeFrame(this.initialTimeFrame, true); // Pass true to emit initial value

      // If there's an initial time frame, set the date range control accordingly
      if (this.currentTimeFrame?.isValid) {
        this.dateRange.setValue({
          start: this.currentTimeFrame.startDate,
          end: this.currentTimeFrame.endDate,
        });
      }
    }

    // Handle form control changes - only update UI, don't emit events
    this._subscriptions.push(
      this.timeFrameControl.valueChanges
        .pipe(debounceTime(500), distinctUntilChanged()) // Increased debounce time
        .subscribe((value) => {
          if (value !== null) {
            this.updateTimeFrame(value, false); // Don't emit events during typing

            // Update date range picker when text input changes
            if (this.currentTimeFrame?.isValid) {
              this.dateRange.setValue({
                start: this.currentTimeFrame.startDate,
                end: this.currentTimeFrame.endDate,
              });
            }
          }
        }),
    );

    // Handle manual changes (for presets)
    this._subscriptions.push(
      this._manualChange.pipe(distinctUntilChanged()).subscribe((value) => {
        this.timeFrameControl.setValue(value);
        this.updateTimeFrame(value, true); // Emit events for preset selection
      }),
    );
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (
      changes["initialTimeFrame"] &&
      !changes["initialTimeFrame"].firstChange
    ) {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
      const newValue = changes["initialTimeFrame"].currentValue;

      // Update the form control without triggering valueChanges
      // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
      this.timeFrameControl.setValue(newValue || "", {
        emitEvent: false,
      });

      if (newValue) {
        // Update the time frame - emit=true for query param changes
        // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
        this.updateTimeFrame(newValue, true);

        // If valid, update date range
        if (this.currentTimeFrame?.isValid) {
          this.dateRange.setValue({
            start: this.currentTimeFrame.startDate,
            end: this.currentTimeFrame.endDate,
          });
        }
      } else {
        // Clear the time frame if empty
        this.currentTimeFrame = null;
        this.dateRange.reset();
      }
    }
  }

  ngOnDestroy(): void {
    this._subscriptions.forEach((sub) => sub.unsubscribe());
    this._subscriptions = [];
  }

  onTabChange(event: { index: number }): void {
    this.activeTabIndex = event.index;
  }

  onPresetSelected(event: { value: string }): void {
    this._manualChange.next(event.value);
  }

  applyDateRangePicker(): void {
    const { start, end } = this.dateRange.value;

    if (!start || !end) {
      return;
    }

    // Format dates for the expression
    const formatDate = (date: Date) => {
      return date.toLocaleDateString("en-US", {
        year: "numeric",
        month: "short",
        day: "numeric",
      });
    };

    // Create a date range expression like "Jan 1, 2023 to Jan 31, 2023"
    const expression = `${formatDate(start)} to ${formatDate(end)}`;

    // Update the time frame control without emitting change events
    this.timeFrameControl.setValue(expression, { emitEvent: false });

    // Parse the expression and emit directly
    const timeFrame = parseTimeFrame(expression);
    this.currentTimeFrame = timeFrame;

    if (timeFrame.isValid) {
      this.timeFrameChanged.emit(timeFrame);
    }
  }

  updateTimeFrameAndEmit(value: string | null): void {
    if (!value) {
      this.clearTimeFrame();
      return;
    }

    // First update the time frame without emitting
    this.updateTimeFrame(value, false);

    // Then only emit if it's valid
    if (this.currentTimeFrame && this.currentTimeFrame.isValid) {
      this.timeFrameChanged.emit(this.currentTimeFrame);
    }
  }

  clearTimeFrame(): void {
    // Reset all controls
    this.timeFrameControl.setValue("");
    this.dateRange.reset();
    this.currentTimeFrame = null;

    this.timeFrameChanged.emit({
      startDate: new Date(),
      endDate: new Date(),
      expression: "",
      isValid: true,
    });
  }

  updateTimeFrame(value: string | null, emit: boolean): void {
    if (!value) {
      return;
    }

    const timeFrame = parseTimeFrame(value);
    this.currentTimeFrame = timeFrame;

    // Only emit if requested and the time frame is valid
    if (emit && timeFrame.isValid) {
      this.timeFrameChanged.emit(timeFrame);
    }
  }

  // Helper to format time frame description
  formatTimeFrameDescription(timeFrame: TimeFrame): string {
    return formatTimeFrameDescription(timeFrame);
  }
}
