import { Injectable, inject } from "@angular/core";
import { BrowserStorageService } from "../browser-storage/browser-storage.service";
import { BehaviorSubject, distinctUntilChanged, Observable } from "rxjs";

const localStorageKey = "bitmagnet-jwt";
const pollInterval = 10000;

@Injectable({ providedIn: "root" })
export class AuthTokenService {
  private browserStorage = inject(BrowserStorageService);
  private tokenSubject: BehaviorSubject<string | null>;

  token$: Observable<string | null>;

  constructor() {
    this.tokenSubject = new BehaviorSubject<string | null>(this.readToken());
    this.token$ = this.tokenSubject.asObservable().pipe(distinctUntilChanged());
    setInterval(() => this.tokenSubject.next(this.readToken()), pollInterval);
  }

  private readToken(): string | null {
    return this.browserStorage.get(localStorageKey) || null;
  }

  getToken(): string | null {
    return this.tokenSubject.getValue();
  }

  setToken(token: string) {
    this.browserStorage.set(localStorageKey, token);
    this.tokenSubject.next(token);
  }

  clearToken() {
    this.browserStorage.remove(localStorageKey);
    this.tokenSubject.next(null);
  }
}
