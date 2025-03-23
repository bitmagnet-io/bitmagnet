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
import { MatExpansionModule } from "@angular/material/expansion";
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
    MatExpansionModule,
    TranslocoModule,
  ],
  providers: [
    provideNativeDateAdapter(),
    { provide: MAT_DATE_FORMATS, useValue: MY_DATE_FORMATS },
  ],
  template: `
    <ng-container *transloco="let t">
      <div class="time-frame-container">
        <!-- Quick Preset Buttons -->
        <div class="quick-presets-grid">
          <button
            *ngFor="let preset of quickPresets"
            mat-stroked-button
            [class.selected]="timeFrameControl.value === preset.value"
            (click)="onPresetSelected({ value: preset.value })"
            class="preset-button"
          >
            <span class="preset-label">{{ preset.label }}</span>
          </button>
        </div>

        <!-- Advanced Filter Accordion -->
        <mat-accordion class="filter-accordion">
          <!-- More Presets Panel -->
          <mat-expansion-panel>
            <mat-expansion-panel-header>
              <mat-panel-title>
                <mat-icon>date_range</mat-icon> {{ t("dates.more_presets") }}
              </mat-panel-title>
            </mat-expansion-panel-header>
            <div class="more-presets-grid">
              <button
                *ngFor="let preset of morePresets"
                mat-stroked-button
                [class.selected]="timeFrameControl.value === preset.value"
                (click)="onPresetSelected({ value: preset.value })"
                class="preset-button"
              >
                {{ preset.label }}
              </button>
            </div>
          </mat-expansion-panel>

          <!-- Calendar Range Panel -->
          <mat-expansion-panel>
            <mat-expansion-panel-header>
              <mat-panel-title>
                <mat-icon>calendar_month</mat-icon> {{ t("dates.date_range") }}
              </mat-panel-title>
            </mat-expansion-panel-header>
            <div class="date-range-container">
              <mat-form-field appearance="outline" class="date-range-field">
                <mat-label>{{ t("dates.select_date_range") }}</mat-label>
                <mat-date-range-input
                  [formGroup]="dateRange"
                  [rangePicker]="rangePicker"
                  [comparisonStart]="comparisonStart"
                  [comparisonEnd]="comparisonEnd"
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
                <mat-date-range-picker 
                  #rangePicker 
                  [dateClass]="dateClass"
                >
                  <mat-date-range-picker-actions>
                    <button mat-button matDateRangePickerCancel>{{ t("general.dismiss") }}</button>
                    <button mat-raised-button color="primary" matDateRangePickerApply>{{ t("general.apply") }}</button>
                  </mat-date-range-picker-actions>
                </mat-date-range-picker>

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

              <div class="date-shortcuts">
                <button 
                  mat-button 
                  (click)="selectThisMonth()"
                  [class.active]="isThisMonthSelected"
                >
                  {{ t("dates.this_month") }}
                </button>
                <button 
                  mat-button 
                  (click)="selectLastMonth()"
                  [class.active]="isLastMonthSelected"
                >
                  {{ t("dates.last_month") }}
                </button>
                <button 
                  mat-button 
                  (click)="selectThisYear()"
                  [class.active]="isThisYearSelected"
                >
                  {{ t("dates.this_year") }}
                </button>
                <button 
                  mat-button 
                  (click)="selectLastYear()"
                  [class.active]="isLastYearSelected"
                >
                  {{ t("dates.last_year") }}
                </button>
              </div>

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
          </mat-expansion-panel>

          <!-- Custom Expression Panel -->
          <mat-expansion-panel>
            <mat-expansion-panel-header>
              <mat-panel-title>
                <mat-icon>edit</mat-icon> {{ t("dates.custom") }}
              </mat-panel-title>
            </mat-expansion-panel-header>
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
                <mat-icon 
                  matSuffix 
                  [matTooltip]="t('dates.time_frame_tooltip')"
                >help_outline</mat-icon>
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
          </mat-expansion-panel>
        </mat-accordion>

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

      .quick-presets-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
        gap: 8px;
      }

      .more-presets-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
        gap: 8px;
        margin-top: 8px;
      }

      .preset-button {
        height: auto;
        line-height: 1.2;
        padding: 8px 12px;
        white-space: normal;
        text-align: center;
      }

      .preset-button.selected {
        background-color: rgba(103, 58, 183, 0.1);
        border-color: rgba(103, 58, 183, 0.5);
        font-weight: 500;
      }

      .preset-label {
        white-space: normal;
        display: block;
      }

      .filter-accordion {
        width: 100%;
      }

      .full-width {
        width: 100%;
      }

      .date-range-container {
        padding: 8px 0;
        display: flex;
        flex-direction: column;
        gap: 16px;
      }

      .date-shortcuts {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
      }
      
      .date-shortcuts button {
        min-width: auto;
        padding: 0 8px;
      }
      
      .date-shortcuts button.active {
        background-color: rgba(103, 58, 183, 0.1);
        font-weight: 500;
      }

      .date-range-field {
        width: 100%;
      }

      .custom-expression {
        padding: 8px 0;
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
        .quick-presets-grid {
          grid-template-columns: repeat(2, 1fr);
        }
        
        .more-presets-grid {
          grid-template-columns: repeat(2, 1fr);
        }
        
        .date-shortcuts {
          flex-direction: column;
          align-items: flex-start;
        }
        
        .date-shortcuts button {
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

  // Comparison date range for highlighting in the calendar
  comparisonStart: Date | null = null;
  comparisonEnd: Date | null = null;

  // State tracking for date shortcut active states
  isThisMonthSelected = false;
  isLastMonthSelected = false;
  isThisYearSelected = false;
  isLastYearSelected = false;

  currentTimeFrame: TimeFrame | null = null;

  // Split presets into quick and more categories
  quickPresets = timeFramePresets.slice(0, 4); // First 4 presets for immediate access
  morePresets = timeFramePresets.slice(4); // Remaining presets in expandable panel

  private _manualChange = new Subject<string>();
  private _subscriptions: Subscription[] = [];

  // Function to apply custom date class for highlighting in the calendar
  dateClass = (date: Date): string => {
    // Highlight today's date
    if (this.isToday(date)) {
      return 'today-date';
    }
    
    // Highlight dates in the current selection range
    if (this.comparisonStart && this.comparisonEnd) {
      if (date >= this.comparisonStart && date <= this.comparisonEnd) {
        return 'comparison-date';
      }
    }
    
    return '';
  };

  // Helper to check if a date is today
  private isToday(date: Date): boolean {
    const today = new Date();
    return date.getDate() === today.getDate() &&
      date.getMonth() === today.getMonth() &&
      date.getFullYear() === today.getFullYear();
  }

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
        
        // Set comparison dates for highlighting
        this.comparisonStart = this.currentTimeFrame.startDate;
        this.comparisonEnd = this.currentTimeFrame.endDate;
        
        // Update shortcut button states
        this.updateShortcutButtonStates();
      }
    }

    // Handle form control changes - only update UI, don't emit events
    this._subscriptions.push(
      this.timeFrameControl.valueChanges
        .pipe(debounceTime(300), distinctUntilChanged())
        .subscribe((value) => {
          if (value !== null) {
            this.updateTimeFrame(value, false); // Don't emit events during typing

            // Update date range picker when text input changes
            if (this.currentTimeFrame?.isValid) {
              this.dateRange.setValue({
                start: this.currentTimeFrame.startDate,
                end: this.currentTimeFrame.endDate,
              });
              
              // Update comparison dates for highlighting
              this.comparisonStart = this.currentTimeFrame.startDate;
              this.comparisonEnd = this.currentTimeFrame.endDate;
              
              // Update shortcut button states
              this.updateShortcutButtonStates();
            }
          }
        }),
    );

    // Handle date range form changes
    this._subscriptions.push(
      this.dateRange.valueChanges.subscribe(range => {
        if (range.start && range.end) {
          // Update comparison dates for highlighting
          this.comparisonStart = range.start;
          this.comparisonEnd = range.end;
          
          // Update shortcut button states when user selects dates
          this.updateShortcutButtonStates();
        }
      })
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
          
          // Update comparison dates for highlighting
          this.comparisonStart = this.currentTimeFrame.startDate;
          this.comparisonEnd = this.currentTimeFrame.endDate;
          
          // Update shortcut button states
          this.updateShortcutButtonStates();
        }
      } else {
        // Clear the time frame if empty
        this.currentTimeFrame = null;
        this.dateRange.reset();
        this.comparisonStart = null;
        this.comparisonEnd = null;
        this.resetShortcutButtonStates();
      }
    }
  }

  ngOnDestroy(): void {
    this._subscriptions.forEach((sub) => sub.unsubscribe());
    this._subscriptions = [];
  }

  onPresetSelected(event: { value: string }): void {
    this._manualChange.next(event.value);
  }

  // Date range preset methods
  selectThisMonth(): void {
    const now = new Date();
    const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1);
    const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0);
    endOfMonth.setHours(23, 59, 59, 999);
    
    this.dateRange.setValue({
      start: startOfMonth,
      end: endOfMonth
    });
    
    this.applyDateRangePicker();
  }

  selectLastMonth(): void {
    const now = new Date();
    const startOfLastMonth = new Date(now.getFullYear(), now.getMonth() - 1, 1);
    const endOfLastMonth = new Date(now.getFullYear(), now.getMonth(), 0);
    endOfLastMonth.setHours(23, 59, 59, 999);
    
    this.dateRange.setValue({
      start: startOfLastMonth,
      end: endOfLastMonth
    });
    
    this.applyDateRangePicker();
  }

  selectThisYear(): void {
    const now = new Date();
    const startOfYear = new Date(now.getFullYear(), 0, 1);
    const endOfYear = new Date(now.getFullYear(), 11, 31);
    endOfYear.setHours(23, 59, 59, 999);
    
    this.dateRange.setValue({
      start: startOfYear,
      end: endOfYear
    });
    
    this.applyDateRangePicker();
  }

  selectLastYear(): void {
    const now = new Date();
    const startOfLastYear = new Date(now.getFullYear() - 1, 0, 1);
    const endOfLastYear = new Date(now.getFullYear() - 1, 11, 31);
    endOfLastYear.setHours(23, 59, 59, 999);
    
    this.dateRange.setValue({
      start: startOfLastYear,
      end: endOfLastYear
    });
    
    this.applyDateRangePicker();
  }
  
  // Update the shortcut button active states based on current date range
  updateShortcutButtonStates(): void {
    this.resetShortcutButtonStates();
    
    const { start, end } = this.dateRange.value;
    if (!start || !end) return;
    
    const now = new Date();
    
    // Check for this month
    const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1);
    const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0);
    endOfMonth.setHours(23, 59, 59, 999);
    
    if (this.isSameDay(start, startOfMonth) && this.isSameDay(end, endOfMonth)) {
      this.isThisMonthSelected = true;
    }
    
    // Check for last month
    const startOfLastMonth = new Date(now.getFullYear(), now.getMonth() - 1, 1);
    const endOfLastMonth = new Date(now.getFullYear(), now.getMonth(), 0);
    endOfLastMonth.setHours(23, 59, 59, 999);
    
    if (this.isSameDay(start, startOfLastMonth) && this.isSameDay(end, endOfLastMonth)) {
      this.isLastMonthSelected = true;
    }
    
    // Check for this year
    const startOfYear = new Date(now.getFullYear(), 0, 1);
    const endOfYear = new Date(now.getFullYear(), 11, 31);
    endOfYear.setHours(23, 59, 59, 999);
    
    if (this.isSameDay(start, startOfYear) && this.isSameDay(end, endOfYear)) {
      this.isThisYearSelected = true;
    }
    
    // Check for last year
    const startOfLastYear = new Date(now.getFullYear() - 1, 0, 1);
    const endOfLastYear = new Date(now.getFullYear() - 1, 11, 31);
    endOfLastYear.setHours(23, 59, 59, 999);
    
    if (this.isSameDay(start, startOfLastYear) && this.isSameDay(end, endOfLastYear)) {
      this.isLastYearSelected = true;
    }
  }
  
  // Helper to check if two dates represent the same day
  private isSameDay(date1: Date, date2: Date): boolean {
    return date1.getDate() === date2.getDate() &&
           date1.getMonth() === date2.getMonth() &&
           date1.getFullYear() === date2.getFullYear();
  }
  
  // Reset all shortcut button states
  resetShortcutButtonStates(): void {
    this.isThisMonthSelected = false;
    this.isLastMonthSelected = false;
    this.isThisYearSelected = false;
    this.isLastYearSelected = false;
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
      
      // Update comparison dates for highlighting
      this.comparisonStart = timeFrame.startDate;
      this.comparisonEnd = timeFrame.endDate;
      
      // Update shortcut button states
      this.updateShortcutButtonStates();
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
    this.comparisonStart = null;
    this.comparisonEnd = null;
    this.resetShortcutButtonStates();

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
