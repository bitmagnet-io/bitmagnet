import { Component, inject, TemplateRef, ViewChild } from "@angular/core";
import { TranslocoService } from "@jsverse/transloco";
import { MatRadioButton, MatRadioGroup } from "@angular/material/radio";
import { MatSliderModule } from "@angular/material/slider";
import { FormsModule, ReactiveFormsModule, FormControl } from "@angular/forms";
import {
  MatBottomSheet,
  MatBottomSheetRef,
} from "@angular/material/bottom-sheet";
import { AppModule } from "../../app.module";
import {
  AggregationBudgetService,
  AggregationBudgetSelection,
} from "./aggregation-budget.service";

enum BudgetTypeEnum {
  Default = 1,
  High,
  Custom,
}

@Component({
  selector: "aggregation-budget",
  standalone: true,
  imports: [
    AppModule,
    MatRadioButton,
    MatRadioGroup,
    MatSliderModule,
    FormsModule,
    ReactiveFormsModule,
  ],
  templateUrl: "./aggregation-budget.component.html",
  styleUrl: "./aggregation-budget.component.scss",
})
export class AggregationBudgetComponent {
  @ViewChild("aggregationBudgetSheet", { static: false })
  _sheetTemplate: TemplateRef<void>;
  private _sheet = inject(MatBottomSheet);
  private _sheetRef: MatBottomSheetRef<void>;
  transloco = inject(TranslocoService);
  aggregationBudget = inject(AggregationBudgetService);
  budgetTypeEnum = BudgetTypeEnum;
  budgetType: BudgetTypeEnum = BudgetTypeEnum.Default;
  budgetTypes: BudgetTypeEnum[] = [
    BudgetTypeEnum.Default,
    BudgetTypeEnum.High,
    BudgetTypeEnum.Custom,
  ];
  sliderBudget = new FormControl<number>(0);
  sliderMax: number = 100000;
  sliderStep: number = this.sliderMax / 10;

  ngOnInit(): void {
    this.aggregationBudget.getCost().subscribe({
      next: (highBudget: AggregationBudgetSelection) => {
        this.sliderBudget.setValue(highBudget);
        if (highBudget) {
          this.sliderMax = highBudget * 10;
          this.sliderStep = highBudget / 10;
        }
      },
    });
    this.sliderBudget.setValue(this.aggregationBudget.highBudget);
  }

  open(): void {
    this._sheetRef = this._sheet.open(this._sheetTemplate);
  }

  selectBudgetType(): void {
    switch (this.budgetType) {
      case BudgetTypeEnum.High: {
        this.aggregationBudget.setBudget(this.aggregationBudget.highBudget);
        break;
      }
      case BudgetTypeEnum.Custom: {
        this.aggregationBudget.setBudget(this.sliderBudget.value);
        break;
      }
      default: {
        this.aggregationBudget.setBudget(null);
      }
    }
  }

  customBudgetLabel(): string {
    if (this.sliderBudget.value == null) {
      return "";
    } else {
      return new Intl.NumberFormat(this.transloco.getActiveLang(), {
        notation: "compact",
      }).format(this.sliderBudget.value);
    }
  }

  highValueText(): string {
    return this.aggregationBudget.highBudget != null
      ? " ( " +
          new Intl.NumberFormat(this.transloco.getActiveLang(), {
            notation: "compact",
          }).format(this.aggregationBudget.highBudget) +
          " )"
      : "";
  }

  itemIcon(budgetType: BudgetTypeEnum): string {
    switch (budgetType) {
      case BudgetTypeEnum.Default:
        return "rocket_launch";
      case BudgetTypeEnum.High:
        return "hourglass";
      case BudgetTypeEnum.Custom:
        return "settings_slow_motion";
    }
  }
}
