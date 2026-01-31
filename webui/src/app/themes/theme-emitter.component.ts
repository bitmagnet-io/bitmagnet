import {
  AfterViewInit,
  Component,
  ElementRef,
  inject,
  ViewChild,
  ViewChildren,
} from "@angular/core";
import { themeColors } from "./theme-constants";
import { ThemeEmitterColorComponent } from "./theme-emitter-color.component";
import { ThemeColors } from "./theme-types";
import { ThemeInfoService } from "./theme-info.service";
import { ThemeManager } from "./theme-manager.service";

@Component({
  selector: "app-theme-emitter",
  template: `
    <ng-container>
      @for (c of themeColors; track c) {
        <app-theme-emitter-color [color]="c" />
      }
      <div class="theme-emitter-lightdark" #lightdark></div>
    </ng-container>
  `,
  styles: [
    `
      :host {
        display: none;
      }

      .theme-emitter-color.background {
        color: var(--mat-app-background-color);
      }

      .theme-emitter-color.foreground {
        color: var(--mat-app-text-color);
      }
    `,
  ],
  standalone: true,
  imports: [ThemeEmitterColorComponent],
})
export class ThemeEmitterComponent implements AfterViewInit {
  private service = inject(ThemeInfoService);
  private themeManager = inject(ThemeManager);

  themeColors = themeColors;

  @ViewChildren(ThemeEmitterColorComponent)
  elements: ThemeEmitterColorComponent[];
  @ViewChild("lightdark") lightdark?: ElementRef;

  constructor() {
    this.themeManager.selectedTheme$.subscribe(() => {
      this.updateThemeColors();
    });
  }

  ngAfterViewInit() {
    this.updateThemeColors();
  }

  updateThemeColors() {
    const colors: Partial<ThemeColors> = {};
    for (const color of this.elements ?? []) {
      colors[color.color] = getComputedStyle(
        color.element.nativeElement as Element,
      ).color;
    }
    const type =
      this.lightdark &&
      getComputedStyle(this.lightdark.nativeElement as Element).color ===
        "rgb(0, 0, 0)"
        ? "dark"
        : "light";
    this.service.setInfo({
      colors: colors as ThemeColors,
      type,
    });
  }
}
