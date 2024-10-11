import { Component, ElementRef, Input, ViewChild } from "@angular/core";
import { ThemeColor } from "./theme-types";

@Component({
  selector: "app-theme-emitter-color",
  standalone: true,
  template: `<div [class]="'theme-emitter-color ' + color" #element></div>`,
})
export class ThemeEmitterColorComponent {
  @Input() color: ThemeColor;
  @ViewChild("element") element: ElementRef;
}
