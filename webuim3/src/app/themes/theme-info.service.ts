import {Injectable} from "@angular/core";
import {BehaviorSubject} from "rxjs";
import {ThemeColors, ThemeInfo, ThemeType} from "./theme-types";

@Injectable({providedIn: "root"})
export class ThemeInfoService {
  private infoSubject = new BehaviorSubject<ThemeInfo>({
    type: "light",
    colors: {}  as ThemeColors
  });

  public info$ = this.infoSubject.asObservable()

  public get colors(): ThemeColors {
    return this.infoSubject.getValue().colors
  }

  public get type(): ThemeType {
    return this.infoSubject.getValue().type
  }

  public get isDark(): boolean {
    return this.type === 'dark'
  }

  public setInfo(info: ThemeInfo) {
    this.infoSubject.next(info);
  }
}
