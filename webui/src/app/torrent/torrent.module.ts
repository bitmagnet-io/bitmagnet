import { NgModule } from "@angular/core";
import { CommonModule, DecimalPipe, NgOptimizedImage } from "@angular/common";
import {
  MatAutocomplete,
  MatAutocompleteTrigger,
  MatOption,
} from "@angular/material/autocomplete";
import { MatButton } from "@angular/material/button";
import {
  MatCard,
  MatCardActions,
  MatCardContent,
} from "@angular/material/card";
import {
  MatChipGrid,
  MatChipInput,
  MatChipRemove,
  MatChipRow,
} from "@angular/material/chips";
import { MatDivider } from "@angular/material/divider";
import { MatFormField } from "@angular/material/form-field";
import { MatIcon } from "@angular/material/icon";
import {
  MatTab,
  MatTabContent,
  MatTabGroup,
  MatTabLabel,
} from "@angular/material/tabs";
import { MatTooltip } from "@angular/material/tooltip";
import { NgxFilesizeModule } from "ngx-filesize";
import { ReactiveFormsModule } from "@angular/forms";
import { CdkCopyToClipboard } from "@angular/cdk/clipboard";
import { TorrentContentComponent } from "./torrent-content.component";

@NgModule({
  declarations: [TorrentContentComponent],
  imports: [
    CommonModule,
    DecimalPipe,
    MatAutocomplete,
    MatAutocompleteTrigger,
    MatButton,
    MatCard,
    MatCardActions,
    MatCardContent,
    MatChipGrid,
    MatChipInput,
    MatChipRemove,
    MatChipRow,
    MatDivider,
    MatFormField,
    MatIcon,
    MatOption,
    MatTab,
    MatTabContent,
    MatTabGroup,
    MatTabLabel,
    MatTooltip,
    NgOptimizedImage,
    NgxFilesizeModule,
    ReactiveFormsModule,
    CdkCopyToClipboard,
  ],
  exports: [TorrentContentComponent],
})
export class TorrentModule {}
