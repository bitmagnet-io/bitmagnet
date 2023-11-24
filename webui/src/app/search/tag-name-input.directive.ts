import { Directive, ElementRef, HostListener } from "@angular/core";

/**
 * TagNameInputDirective transforms tag name input to valid kebab-case.
 */
@Directive({
  selector: "input[tagNameInput]",
})
export class TagNameInputDirective {
  constructor(private el: ElementRef<HTMLInputElement>) {}

  @HostListener("input", ["$event"]) onInputChange($event: {
    target: { value: string };
  }) {
    this.el.nativeElement.value = normalizeTagInput($event.target.value);
  }
}

const normalizeTagInput = (value: string): string =>
  value
    .toLowerCase()
    .replaceAll(/[^a-z0-9\-]/g, "-")
    .replace(/^-+/, "")
    .replaceAll(/-+/g, "-");
