import { MatIconRegistry } from "@angular/material/icon";
import { DomSanitizer } from "@angular/platform-browser";

export const initializeIcons = (
  iconRegistry: MatIconRegistry,
  domSanitizer: DomSanitizer,
) =>
  iconRegistry
    .setDefaultFontSetClass(
      "material-icons-outlined",
      "material-symbols-outlined",
    )
    .addSvgIcon(
      "magnet",
      domSanitizer.bypassSecurityTrustResourceUrl("magnet.svg"),
    )
    .addSvgIcon(
      "external-link",
      domSanitizer.bypassSecurityTrustResourceUrl("external-link.svg"),
    )
    .addSvgIcon(
      "binary",
      domSanitizer.bypassSecurityTrustResourceUrl("binary.svg"),
    )
    .addSvgIcon(
      "queue",
      domSanitizer.bypassSecurityTrustResourceUrl("queue.svg"),
    );
