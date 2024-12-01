import { Component, inject, Input, OnChanges, OnInit } from "@angular/core";
import { Title } from "@angular/platform-browser";

@Component({
  selector: "app-document-title",
  standalone: true,
  template: "<ng-container></ng-container>",
})
export class DocumentTitleComponent implements OnInit, OnChanges {
  private title = inject(Title);
  @Input() parts: Array<string | null | undefined> = [];

  ngOnInit() {
    this.updateTitle();
  }

  ngOnChanges() {
    this.updateTitle();
  }

  private updateTitle() {
    this.title.setTitle(
      [...this.parts.filter(Boolean), "bitmagnet"].join(" - "),
    );
  }
}
