import { Injectable } from "@angular/core";
import { BehaviorSubject } from "rxjs";
import { ThemeColors, ThemeInfo, ThemeType } from "./theme-types";
import { emptyThemeInfo } from "./theme-constants";

@Injectable({ providedIn: "root" })
export class ThemeInfoService {
  private infoSubject = new BehaviorSubject<ThemeInfo>(emptyThemeInfo);

  public info$ = this.infoSubject.asObservable();

  public get info(): ThemeInfo {
    return this.infoSubject.getValue();
  }

  public get colors(): ThemeColors {
    return this.info.colors;
  }

  public get type(): ThemeType {
    return this.info.type;
  }

  public get isDark(): boolean {
    return this.type === "dark";
  }

  public setInfo(info: ThemeInfo) {
    this.infoSubject.next(info);
  }
}
