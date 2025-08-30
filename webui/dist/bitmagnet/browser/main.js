import {
  ThemeInfoService,
  provideCharts,
  themeColors,
  withDefaultRegisterables
} from "./chunk-DCY4KWPQ.js";
import {
  HealthModule,
  HealthService,
  HealthWidgetComponent
} from "./chunk-WNZLSZUT.js";
import {
  BreakpointsService
} from "./chunk-NQ6E5D5R.js";
import {
  Apollo,
  ApolloLink,
  AppModule,
  GraphQLModule,
  InMemoryCache,
  MatAnchor,
  MatIcon,
  MatIconAnchor,
  MatIconButton,
  MatIconRegistry,
  MatMenu,
  MatMenuItem,
  MatMenuTrigger,
  MatToolbar,
  MatTooltip,
  Observable as Observable2,
  TranslocoDirective,
  TranslocoService,
  VersionDocument,
  print,
  provideApollo,
  provideTransloco
} from "./chunk-WWRDQTKJ.js";
import "./chunk-VSVMRYN2.js";
import {
  DomRendererFactory2,
  DomSanitizer,
  HttpClient,
  HttpHeaders,
  Router,
  RouterLink,
  RouterLinkActive,
  RouterOutlet,
  Title,
  bootstrapApplication,
  provideHttpClient,
  provideRouter,
  withComponentInputBinding,
  withInterceptorsFromDi
} from "./chunk-Y2ZC5Z2X.js";
import {
  ANIMATION_MODULE_TYPE,
  BehaviorSubject,
  ChangeDetectionScheduler,
  DOCUMENT,
  Injectable,
  InjectionToken,
  NgZone,
  Observable,
  RendererFactory2,
  RuntimeError,
  __assign,
  __async,
  __extends,
  __spreadValues,
  inject,
  makeEnvironmentProviders,
  map,
  performanceMarkFeature,
  provideZoneChangeDetection,
  setClassMetadata,
  ɵsetClassDebugInfo,
  ɵɵStandaloneFeature,
  ɵɵadvance,
  ɵɵclassMap,
  ɵɵconditional,
  ɵɵdefineComponent,
  ɵɵdefineInjectable,
  ɵɵdirectiveInject,
  ɵɵelement,
  ɵɵelementContainerEnd,
  ɵɵelementContainerStart,
  ɵɵelementEnd,
  ɵɵelementStart,
  ɵɵgetCurrentView,
  ɵɵinject,
  ɵɵinvalidFactory,
  ɵɵlistener,
  ɵɵloadQuery,
  ɵɵnextContext,
  ɵɵprojection,
  ɵɵprojectionDef,
  ɵɵproperty,
  ɵɵpureFunction0,
  ɵɵpureFunction1,
  ɵɵqueryRefresh,
  ɵɵreference,
  ɵɵrepeater,
  ɵɵrepeaterCreate,
  ɵɵrepeaterTrackByIdentity,
  ɵɵresetView,
  ɵɵrestoreView,
  ɵɵtemplate,
  ɵɵtext,
  ɵɵtextInterpolate,
  ɵɵtextInterpolate1,
  ɵɵviewQuery
} from "./chunk-DMMUMX3A.js";

// node_modules/@angular/platform-browser/fesm2022/animations/async.mjs
var ANIMATION_PREFIX = "@";
var AsyncAnimationRendererFactory = class _AsyncAnimationRendererFactory {
  /**
   *
   * @param moduleImpl allows to provide a mock implmentation (or will load the animation module)
   */
  constructor(doc, delegate, zone, animationType, moduleImpl) {
    this.doc = doc;
    this.delegate = delegate;
    this.zone = zone;
    this.animationType = animationType;
    this.moduleImpl = moduleImpl;
    this._rendererFactoryPromise = null;
    this.scheduler = inject(ChangeDetectionScheduler, {
      optional: true
    });
    this.loadingSchedulerFn = inject(\u0275ASYNC_ANIMATION_LOADING_SCHEDULER_FN, {
      optional: true
    });
  }
  /** @nodoc */
  ngOnDestroy() {
    this._engine?.flush();
  }
  /**
   * @internal
   */
  loadImpl() {
    const loadFn = () => this.moduleImpl ?? import("./chunk-PJWTV352.js").then((m) => m);
    let moduleImplPromise;
    if (this.loadingSchedulerFn) {
      moduleImplPromise = this.loadingSchedulerFn(loadFn);
    } else {
      moduleImplPromise = loadFn();
    }
    return moduleImplPromise.catch((e) => {
      throw new RuntimeError(5300, (typeof ngDevMode === "undefined" || ngDevMode) && "Async loading for animations package was enabled, but loading failed. Angular falls back to using regular rendering. No animations will be displayed and their styles won't be applied.");
    }).then(({
      \u0275createEngine,
      \u0275AnimationRendererFactory
    }) => {
      this._engine = \u0275createEngine(this.animationType, this.doc);
      const rendererFactory = new \u0275AnimationRendererFactory(this.delegate, this._engine, this.zone);
      this.delegate = rendererFactory;
      return rendererFactory;
    });
  }
  /**
   * This method is delegating the renderer creation to the factories.
   * It uses default factory while the animation factory isn't loaded
   * and will rely on the animation factory once it is loaded.
   *
   * Calling this method will trigger as side effect the loading of the animation module
   * if the renderered component uses animations.
   */
  createRenderer(hostElement, rendererType) {
    const renderer = this.delegate.createRenderer(hostElement, rendererType);
    if (renderer.\u0275type === 0) {
      return renderer;
    }
    if (typeof renderer.throwOnSyntheticProps === "boolean") {
      renderer.throwOnSyntheticProps = false;
    }
    const dynamicRenderer = new DynamicDelegationRenderer(renderer);
    if (rendererType?.data?.["animation"] && !this._rendererFactoryPromise) {
      this._rendererFactoryPromise = this.loadImpl();
    }
    this._rendererFactoryPromise?.then((animationRendererFactory) => {
      const animationRenderer = animationRendererFactory.createRenderer(hostElement, rendererType);
      dynamicRenderer.use(animationRenderer);
      this.scheduler?.notify(
        10
        /* NotificationSource.AsyncAnimationsLoaded */
      );
    }).catch((e) => {
      dynamicRenderer.use(renderer);
    });
    return dynamicRenderer;
  }
  begin() {
    this.delegate.begin?.();
  }
  end() {
    this.delegate.end?.();
  }
  whenRenderingDone() {
    return this.delegate.whenRenderingDone?.() ?? Promise.resolve();
  }
  static {
    this.\u0275fac = function AsyncAnimationRendererFactory_Factory(__ngFactoryType__) {
      \u0275\u0275invalidFactory();
    };
  }
  static {
    this.\u0275prov = /* @__PURE__ */ \u0275\u0275defineInjectable({
      token: _AsyncAnimationRendererFactory,
      factory: _AsyncAnimationRendererFactory.\u0275fac
    });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && setClassMetadata(AsyncAnimationRendererFactory, [{
    type: Injectable
  }], () => [{
    type: Document
  }, {
    type: RendererFactory2
  }, {
    type: NgZone
  }, {
    type: void 0
  }, {
    type: Promise
  }], null);
})();
var DynamicDelegationRenderer = class {
  constructor(delegate) {
    this.delegate = delegate;
    this.replay = [];
    this.\u0275type = 1;
  }
  use(impl) {
    this.delegate = impl;
    if (this.replay !== null) {
      for (const fn of this.replay) {
        fn(impl);
      }
      this.replay = null;
    }
  }
  get data() {
    return this.delegate.data;
  }
  destroy() {
    this.replay = null;
    this.delegate.destroy();
  }
  createElement(name, namespace) {
    return this.delegate.createElement(name, namespace);
  }
  createComment(value) {
    return this.delegate.createComment(value);
  }
  createText(value) {
    return this.delegate.createText(value);
  }
  get destroyNode() {
    return this.delegate.destroyNode;
  }
  appendChild(parent, newChild) {
    this.delegate.appendChild(parent, newChild);
  }
  insertBefore(parent, newChild, refChild, isMove) {
    this.delegate.insertBefore(parent, newChild, refChild, isMove);
  }
  removeChild(parent, oldChild, isHostElement) {
    this.delegate.removeChild(parent, oldChild, isHostElement);
  }
  selectRootElement(selectorOrNode, preserveContent) {
    return this.delegate.selectRootElement(selectorOrNode, preserveContent);
  }
  parentNode(node) {
    return this.delegate.parentNode(node);
  }
  nextSibling(node) {
    return this.delegate.nextSibling(node);
  }
  setAttribute(el, name, value, namespace) {
    this.delegate.setAttribute(el, name, value, namespace);
  }
  removeAttribute(el, name, namespace) {
    this.delegate.removeAttribute(el, name, namespace);
  }
  addClass(el, name) {
    this.delegate.addClass(el, name);
  }
  removeClass(el, name) {
    this.delegate.removeClass(el, name);
  }
  setStyle(el, style, value, flags) {
    this.delegate.setStyle(el, style, value, flags);
  }
  removeStyle(el, style, flags) {
    this.delegate.removeStyle(el, style, flags);
  }
  setProperty(el, name, value) {
    if (this.shouldReplay(name)) {
      this.replay.push((renderer) => renderer.setProperty(el, name, value));
    }
    this.delegate.setProperty(el, name, value);
  }
  setValue(node, value) {
    this.delegate.setValue(node, value);
  }
  listen(target, eventName, callback) {
    if (this.shouldReplay(eventName)) {
      this.replay.push((renderer) => renderer.listen(target, eventName, callback));
    }
    return this.delegate.listen(target, eventName, callback);
  }
  shouldReplay(propOrEventName) {
    return this.replay !== null && propOrEventName.startsWith(ANIMATION_PREFIX);
  }
};
var \u0275ASYNC_ANIMATION_LOADING_SCHEDULER_FN = new InjectionToken(ngDevMode ? "async_animation_loading_scheduler_fn" : "");
function provideAnimationsAsync(type = "animations") {
  performanceMarkFeature("NgAsyncAnimations");
  return makeEnvironmentProviders([{
    provide: RendererFactory2,
    useFactory: (doc, renderer, zone) => {
      return new AsyncAnimationRendererFactory(doc, renderer, zone, type);
    },
    deps: [DOCUMENT, DomRendererFactory2, NgZone]
  }, {
    provide: ANIMATION_MODULE_TYPE,
    useValue: type === "noop" ? "NoopAnimations" : "BrowserAnimations"
  }]);
}

// node_modules/@apollo/client/link/batch/batching.js
var OperationBatcher = (
  /** @class */
  function() {
    function OperationBatcher2(_a) {
      var batchDebounce = _a.batchDebounce, batchInterval = _a.batchInterval, batchMax = _a.batchMax, batchHandler = _a.batchHandler, batchKey = _a.batchKey;
      this.batchesByKey = /* @__PURE__ */ new Map();
      this.scheduledBatchTimerByKey = /* @__PURE__ */ new Map();
      this.batchDebounce = batchDebounce;
      this.batchInterval = batchInterval;
      this.batchMax = batchMax || 0;
      this.batchHandler = batchHandler;
      this.batchKey = batchKey || function() {
        return "";
      };
    }
    OperationBatcher2.prototype.enqueueRequest = function(request) {
      var _this = this;
      var requestCopy = __assign(__assign({}, request), {
        next: [],
        error: [],
        complete: [],
        subscribers: /* @__PURE__ */ new Set()
      });
      var key = this.batchKey(request.operation);
      if (!requestCopy.observable) {
        requestCopy.observable = new Observable2(function(observer) {
          var batch = _this.batchesByKey.get(key);
          if (!batch) _this.batchesByKey.set(key, batch = /* @__PURE__ */ new Set());
          var isFirstEnqueuedRequest = batch.size === 0;
          var isFirstSubscriber = requestCopy.subscribers.size === 0;
          requestCopy.subscribers.add(observer);
          if (isFirstSubscriber) {
            batch.add(requestCopy);
          }
          if (observer.next) {
            requestCopy.next.push(observer.next.bind(observer));
          }
          if (observer.error) {
            requestCopy.error.push(observer.error.bind(observer));
          }
          if (observer.complete) {
            requestCopy.complete.push(observer.complete.bind(observer));
          }
          if (isFirstEnqueuedRequest || _this.batchDebounce) {
            _this.scheduleQueueConsumption(key);
          }
          if (batch.size === _this.batchMax) {
            _this.consumeQueue(key);
          }
          return function() {
            var _a;
            if (requestCopy.subscribers.delete(observer) && requestCopy.subscribers.size < 1) {
              if (batch.delete(requestCopy) && batch.size < 1) {
                _this.consumeQueue(key);
                (_a = batch.subscription) === null || _a === void 0 ? void 0 : _a.unsubscribe();
              }
            }
          };
        });
      }
      return requestCopy.observable;
    };
    OperationBatcher2.prototype.consumeQueue = function(key) {
      if (key === void 0) {
        key = "";
      }
      var batch = this.batchesByKey.get(key);
      this.batchesByKey.delete(key);
      if (!batch || !batch.size) {
        return;
      }
      var operations = [];
      var forwards = [];
      var observables = [];
      var nexts = [];
      var errors = [];
      var completes = [];
      batch.forEach(function(request) {
        operations.push(request.operation);
        forwards.push(request.forward);
        observables.push(request.observable);
        nexts.push(request.next);
        errors.push(request.error);
        completes.push(request.complete);
      });
      var batchedObservable = this.batchHandler(operations, forwards) || Observable2.of();
      var onError = function(error) {
        errors.forEach(function(rejecters) {
          if (rejecters) {
            rejecters.forEach(function(e) {
              return e(error);
            });
          }
        });
      };
      batch.subscription = batchedObservable.subscribe({
        next: function(results) {
          if (!Array.isArray(results)) {
            results = [results];
          }
          if (nexts.length !== results.length) {
            var error = new Error("server returned results with length ".concat(results.length, ", expected length of ").concat(nexts.length));
            error.result = results;
            return onError(error);
          }
          results.forEach(function(result, index) {
            if (nexts[index]) {
              nexts[index].forEach(function(next) {
                return next(result);
              });
            }
          });
        },
        error: onError,
        complete: function() {
          completes.forEach(function(complete) {
            if (complete) {
              complete.forEach(function(c) {
                return c();
              });
            }
          });
        }
      });
      return observables;
    };
    OperationBatcher2.prototype.scheduleQueueConsumption = function(key) {
      var _this = this;
      clearTimeout(this.scheduledBatchTimerByKey.get(key));
      this.scheduledBatchTimerByKey.set(key, setTimeout(function() {
        _this.consumeQueue(key);
        _this.scheduledBatchTimerByKey.delete(key);
      }, this.batchInterval));
    };
    return OperationBatcher2;
  }()
);

// node_modules/@apollo/client/link/batch/batchLink.js
var BatchLink = (
  /** @class */
  function(_super) {
    __extends(BatchLink2, _super);
    function BatchLink2(fetchParams) {
      var _this = _super.call(this) || this;
      var _a = fetchParams || {}, batchDebounce = _a.batchDebounce, _b = _a.batchInterval, batchInterval = _b === void 0 ? 10 : _b, _c = _a.batchMax, batchMax = _c === void 0 ? 0 : _c, _d = _a.batchHandler, batchHandler = _d === void 0 ? function() {
        return null;
      } : _d, _e = _a.batchKey, batchKey = _e === void 0 ? function() {
        return "";
      } : _e;
      _this.batcher = new OperationBatcher({
        batchDebounce,
        batchInterval,
        batchMax,
        batchHandler,
        batchKey
      });
      if (fetchParams.batchHandler.length <= 1) {
        _this.request = function(operation) {
          return _this.batcher.enqueueRequest({
            operation
          });
        };
      }
      return _this;
    }
    BatchLink2.prototype.request = function(operation, forward) {
      return this.batcher.enqueueRequest({
        operation,
        forward
      });
    };
    return BatchLink2;
  }(ApolloLink)
);

// node_modules/apollo-angular/fesm2022/ngApolloLinkHttp.mjs
var fetch = (req, httpClient, extractFiles) => {
  const shouldUseBody = ["POST", "PUT", "PATCH"].indexOf(req.method.toUpperCase()) !== -1;
  const shouldStringify = (param) => ["variables", "extensions"].indexOf(param.toLowerCase()) !== -1;
  const isBatching = req.body.length;
  let shouldUseMultipart = req.options && req.options.useMultipart;
  let multipartInfo;
  if (shouldUseMultipart) {
    if (isBatching) {
      return new Observable((observer) => observer.error(new Error("File upload is not available when combined with Batching")));
    }
    if (!shouldUseBody) {
      return new Observable((observer) => observer.error(new Error("File upload is not available when GET is used")));
    }
    if (!extractFiles) {
      return new Observable((observer) => observer.error(new Error(`To use File upload you need to pass "extractFiles" function from "extract-files" library to HttpLink's options`)));
    }
    multipartInfo = extractFiles(req.body);
    shouldUseMultipart = !!multipartInfo.files.size;
  }
  let bodyOrParams = {};
  if (isBatching) {
    if (!shouldUseBody) {
      return new Observable((observer) => observer.error(new Error("Batching is not available for GET requests")));
    }
    bodyOrParams = {
      body: req.body
    };
  } else {
    const body = shouldUseMultipart ? multipartInfo.clone : req.body;
    if (shouldUseBody) {
      bodyOrParams = {
        body
      };
    } else {
      const params = Object.keys(req.body).reduce((obj, param) => {
        const value = req.body[param];
        obj[param] = shouldStringify(param) ? JSON.stringify(value) : value;
        return obj;
      }, {});
      bodyOrParams = {
        params
      };
    }
  }
  if (shouldUseMultipart && shouldUseBody) {
    const form = new FormData();
    form.append("operations", JSON.stringify(bodyOrParams.body));
    const map2 = {};
    const files = multipartInfo.files;
    let i = 0;
    files.forEach((paths) => {
      map2[++i] = paths;
    });
    form.append("map", JSON.stringify(map2));
    i = 0;
    files.forEach((_, file) => {
      form.append(++i + "", file, file.name);
    });
    bodyOrParams.body = form;
  }
  return httpClient.request(req.method, req.url, __spreadValues(__spreadValues({
    observe: "response",
    responseType: "json",
    reportProgress: false
  }, bodyOrParams), req.options));
};
var mergeHeaders = (source, destination) => {
  if (source && destination) {
    const merged = destination.keys().reduce((headers, name) => headers.set(name, destination.getAll(name)), source);
    return merged;
  }
  return destination || source;
};
function prioritize(...values) {
  return values.find((val) => typeof val !== "undefined");
}
function createHeadersWithClientAwareness(context) {
  let headers = context.headers && context.headers instanceof HttpHeaders ? context.headers : new HttpHeaders(context.headers);
  if (context.clientAwareness) {
    const {
      name,
      version
    } = context.clientAwareness;
    if (name && !headers.has("apollographql-client-name")) {
      headers = headers.set("apollographql-client-name", name);
    }
    if (version && !headers.has("apollographql-client-version")) {
      headers = headers.set("apollographql-client-version", version);
    }
  }
  return headers;
}
var defaults = {
  batchInterval: 10,
  batchMax: 10,
  uri: "graphql",
  method: "POST",
  withCredentials: false,
  includeQuery: true,
  includeExtensions: false,
  useMultipart: false
};
function pick(context, options, key) {
  return prioritize(context[key], options[key], defaults[key]);
}
var HttpBatchLinkHandler = class extends ApolloLink {
  httpClient;
  options;
  batcher;
  batchInterval;
  batchMax;
  print = print;
  constructor(httpClient, options) {
    super();
    this.httpClient = httpClient;
    this.options = options;
    this.batchInterval = options.batchInterval || defaults.batchInterval;
    this.batchMax = options.batchMax || defaults.batchMax;
    if (this.options.operationPrinter) {
      this.print = this.options.operationPrinter;
    }
    const batchHandler = (operations) => {
      return new Observable2((observer) => {
        const body = this.createBody(operations);
        const headers = this.createHeaders(operations);
        const {
          method,
          uri,
          withCredentials
        } = this.createOptions(operations);
        if (typeof uri === "function") {
          throw new Error(`Option 'uri' is a function, should be a string`);
        }
        const req = {
          method,
          url: uri,
          body,
          options: {
            withCredentials,
            headers
          }
        };
        const sub = fetch(req, this.httpClient, () => {
          throw new Error("File upload is not available when combined with Batching");
        }).subscribe({
          next: (result) => observer.next(result.body),
          error: (err) => observer.error(err),
          complete: () => observer.complete()
        });
        return () => {
          if (!sub.closed) {
            sub.unsubscribe();
          }
        };
      });
    };
    const batchKey = options.batchKey || ((operation) => {
      return this.createBatchKey(operation);
    });
    this.batcher = new BatchLink({
      batchInterval: this.batchInterval,
      batchMax: this.batchMax,
      batchKey,
      batchHandler
    });
  }
  createOptions(operations) {
    const context = operations[0].getContext();
    return {
      method: pick(context, this.options, "method"),
      uri: pick(context, this.options, "uri"),
      withCredentials: pick(context, this.options, "withCredentials")
    };
  }
  createBody(operations) {
    return operations.map((operation) => {
      const includeExtensions = prioritize(operation.getContext().includeExtensions, this.options.includeExtensions, false);
      const includeQuery = prioritize(operation.getContext().includeQuery, this.options.includeQuery, true);
      const body = {
        operationName: operation.operationName,
        variables: operation.variables
      };
      if (includeExtensions) {
        body.extensions = operation.extensions;
      }
      if (includeQuery) {
        body.query = this.print(operation.query);
      }
      return body;
    });
  }
  createHeaders(operations) {
    return operations.reduce((headers, operation) => {
      return mergeHeaders(headers, operation.getContext().headers);
    }, createHeadersWithClientAwareness({
      headers: this.options.headers,
      clientAwareness: operations[0]?.getContext()?.clientAwareness
    }));
  }
  createBatchKey(operation) {
    const context = operation.getContext();
    if (context.skipBatching) {
      return Math.random().toString(36).substring(2, 11);
    }
    const headers = context.headers && context.headers.keys().map((k) => context.headers.get(k));
    const opts = JSON.stringify({
      includeQuery: context.includeQuery,
      includeExtensions: context.includeExtensions,
      headers
    });
    return prioritize(context.uri, this.options.uri, "") + opts;
  }
  request(op) {
    return this.batcher.request(op);
  }
};
var HttpBatchLink = class _HttpBatchLink {
  httpClient;
  constructor(httpClient) {
    this.httpClient = httpClient;
  }
  create(options) {
    return new HttpBatchLinkHandler(this.httpClient, options);
  }
  static \u0275fac = function HttpBatchLink_Factory(__ngFactoryType__) {
    return new (__ngFactoryType__ || _HttpBatchLink)(\u0275\u0275inject(HttpClient));
  };
  static \u0275prov = /* @__PURE__ */ \u0275\u0275defineInjectable({
    token: _HttpBatchLink,
    factory: _HttpBatchLink.\u0275fac,
    providedIn: "root"
  });
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && setClassMetadata(HttpBatchLink, [{
    type: Injectable,
    args: [{
      providedIn: "root"
    }]
  }], () => [{
    type: HttpClient
  }], null);
})();
var HttpLinkHandler = class extends ApolloLink {
  httpClient;
  options;
  requester;
  print = print;
  constructor(httpClient, options) {
    super();
    this.httpClient = httpClient;
    this.options = options;
    if (this.options.operationPrinter) {
      this.print = this.options.operationPrinter;
    }
    this.requester = (operation) => new Observable2((observer) => {
      const context = operation.getContext();
      let method = pick(context, this.options, "method");
      const includeQuery = pick(context, this.options, "includeQuery");
      const includeExtensions = pick(context, this.options, "includeExtensions");
      const url = pick(context, this.options, "uri");
      const withCredentials = pick(context, this.options, "withCredentials");
      const useMultipart = pick(context, this.options, "useMultipart");
      const useGETForQueries = this.options.useGETForQueries === true;
      const isQuery = operation.query.definitions.some((def) => def.kind === "OperationDefinition" && def.operation === "query");
      if (useGETForQueries && isQuery) {
        method = "GET";
      }
      const req = {
        method,
        url: typeof url === "function" ? url(operation) : url,
        body: {
          operationName: operation.operationName,
          variables: operation.variables
        },
        options: {
          withCredentials,
          useMultipart,
          headers: this.options.headers
        }
      };
      if (includeExtensions) {
        req.body.extensions = operation.extensions;
      }
      if (includeQuery) {
        req.body.query = this.print(operation.query);
      }
      const headers = createHeadersWithClientAwareness(context);
      req.options.headers = mergeHeaders(req.options.headers, headers);
      const sub = fetch(req, this.httpClient, this.options.extractFiles).subscribe({
        next: (response) => {
          operation.setContext({
            response
          });
          observer.next(response.body);
        },
        error: (err) => observer.error(err),
        complete: () => observer.complete()
      });
      return () => {
        if (!sub.closed) {
          sub.unsubscribe();
        }
      };
    });
  }
  request(op) {
    return this.requester(op);
  }
};
var HttpLink = class _HttpLink {
  httpClient;
  constructor(httpClient) {
    this.httpClient = httpClient;
  }
  create(options) {
    return new HttpLinkHandler(this.httpClient, options);
  }
  static \u0275fac = function HttpLink_Factory(__ngFactoryType__) {
    return new (__ngFactoryType__ || _HttpLink)(\u0275\u0275inject(HttpClient));
  };
  static \u0275prov = /* @__PURE__ */ \u0275\u0275defineInjectable({
    token: _HttpLink,
    factory: _HttpLink.\u0275fac,
    providedIn: "root"
  });
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && setClassMetadata(HttpLink, [{
    type: Injectable,
    args: [{
      providedIn: "root"
    }]
  }], () => [{
    type: HttpClient
  }], null);
})();

// src/environments/environment.ts
var graphqlEndpoint = "http://localhost:3333/graphql";

// src/app/i18n/translations/ar.json
var ar_default = {
  content_types: {
    plural: {
      all: "\u0627\u0644\u0643\u0644",
      audiobook: "\u0627\u0644\u0643\u062A\u0628 \u0627\u0644\u0635\u0648\u062A\u064A\u0629",
      comic: "\u0627\u0644\u0642\u0635\u0635 \u0627\u0644\u0645\u0635\u0648\u0631\u0629",
      ebook: "\u0627\u0644\u0643\u062A\u0628 \u0627\u0644\u0625\u0644\u0643\u062A\u0631\u0648\u0646\u064A\u0629",
      game: "\u0627\u0644\u0623\u0644\u0639\u0627\u0628",
      movie: "\u0627\u0644\u0623\u0641\u0644\u0627\u0645",
      music: "\u0627\u0644\u0645\u0648\u0633\u064A\u0642\u0649",
      null: "\u063A\u064A\u0631 \u0645\u0639\u0631\u0648\u0641",
      software: "\u0627\u0644\u0628\u0631\u0645\u062C\u064A\u0627\u062A",
      tv_show: "\u0628\u0631\u0627\u0645\u062C \u0627\u0644\u062A\u0644\u0641\u0627\u0632",
      xxx: "\u0627\u0644\u0625\u0628\u0627\u062D\u064A\u0629"
    },
    singular: {
      audiobook: "\u0643\u062A\u0627\u0628 \u0635\u0648\u062A\u064A",
      comic: "\u0642\u0635\u0629 \u0645\u0635\u0648\u0631\u0629",
      ebook: "\u0643\u062A\u0627\u0628 \u0625\u0644\u0643\u062A\u0631\u0648\u0646\u064A",
      game: "\u0644\u0639\u0628\u0629",
      movie: "\u0641\u064A\u0644\u0645",
      music: "\u0645\u0648\u0633\u064A\u0642\u0649",
      null: "\u063A\u064A\u0631 \u0645\u0639\u0631\u0648\u0641",
      software: "\u0628\u0631\u0645\u062C\u064A\u0627\u062A",
      tv_show: "\u0628\u0631\u0646\u0627\u0645\u062C \u062A\u0644\u0641\u0632\u064A\u0648\u0646\u064A",
      xxx: "\u0627\u0644\u0625\u0628\u0627\u062D\u064A\u0629"
    }
  },
  dashboard: {
    event: {
      created: "\u062A\u0645 \u0627\u0644\u0625\u0646\u0634\u0627\u0621",
      failed: "\u0641\u0634\u0644",
      processed: "\u062A\u0645\u062A \u0627\u0644\u0645\u0639\u0627\u0644\u062C\u0629",
      updated: "\u062A\u0645 \u0627\u0644\u062A\u062D\u062F\u064A\u062B"
    },
    interval: {
      all: "\u0627\u0644\u0643\u0644",
      days: "\u064A\u0648\u0645",
      days_1: "\u064A\u0648\u0645 \u0648\u0627\u062D\u062F",
      hours: "\u0633\u0627\u0639\u0629",
      hours_1: "\u0633\u0627\u0639\u0629 \u0648\u0627\u062D\u062F\u0629",
      hours_12: "12 \u0633\u0627\u0639\u0629",
      hours_6: "6 \u0633\u0627\u0639\u0629",
      minutes: "\u062F\u0642\u064A\u0642\u0629",
      minutes_1: "1 \u062F\u0642\u064A\u0642\u0629",
      minutes_15: "15 \u062F\u0642\u0627\u0626\u0642",
      minutes_30: "30 \u062F\u0642\u0627\u0626\u0642",
      minutes_5: "5 \u062F\u0642\u0627\u0626\u0642",
      off: "\u0625\u064A\u0642\u0627\u0641",
      seconds_10: "10 \u062B\u0648\u0627\u0646\u064A",
      seconds_30: "30 \u062B\u0627\u0646\u064A\u0629",
      weeks_1: "1 \u0623\u0633\u0628\u0648\u0639"
    },
    metrics: {
      event: "\u062D\u062F\u062B",
      resolution: "\u0627\u0644\u062F\u0642\u0629",
      throughput: "\u0645\u0639\u062F\u0644 \u0627\u0644\u0646\u0642\u0644",
      timeframe: "\u0627\u0644\u0625\u0637\u0627\u0631 \u0627\u0644\u0632\u0645\u0646\u064A",
      toggle_legend: "\u062A\u0628\u062F\u064A\u0644 \u0627\u0644\u0623\u0633\u0637\u0648\u0631\u0629"
    },
    queues: {
      created: "\u062A\u0645 \u0627\u0644\u0625\u0646\u0634\u0627\u0621",
      created_at: "\u062A\u0645 \u0627\u0644\u0625\u0646\u0634\u0627\u0621 \u0641\u064A",
      enqueue_jobs: "\u0625\u062F\u0631\u0627\u062C \u0627\u0644\u0648\u0638\u0627\u0626\u0641 \u0641\u064A \u0627\u0644\u0637\u0627\u0628\u0648\u0631",
      enqueue_torrent_processing_batch: "\u0625\u062F\u0631\u0627\u062C \u062F\u0641\u0639\u0629 \u0645\u0639\u0627\u0644\u062C\u0629 \u0627\u0644\u062A\u0648\u0631\u0646\u062A \u0641\u064A \u0627\u0644\u0637\u0627\u0628\u0648\u0631",
      failed: "\u0641\u0634\u0644",
      force_rematch: "\u0641\u0631\u0636 \u0625\u0639\u0627\u062F\u0629 \u0627\u0644\u0645\u0637\u0627\u0628\u0642\u0629 \u0644\u0644\u0645\u062D\u062A\u0648\u0649 \u0627\u0644\u0645\u0637\u0627\u0628\u0642 \u0628\u0627\u0644\u0641\u0639\u0644",
      jobs_enqueued: "\u0627\u0644\u0648\u0638\u0627\u0626\u0641 \u0627\u0644\u0645\u062F\u0631\u062C\u0629 \u0641\u064A \u0627\u0644\u0637\u0627\u0628\u0648\u0631",
      latency: "\u0627\u0644\u062A\u0623\u062E\u064A\u0631",
      match_content_by_external_api_search: "\u0645\u0637\u0627\u0628\u0642\u0629 \u0627\u0644\u0645\u062D\u062A\u0648\u0649 \u0645\u0646 \u062E\u0644\u0627\u0644 \u0627\u0644\u0628\u062D\u062B \u0641\u064A API \u0627\u0644\u062E\u0627\u0631\u062C\u064A\u0629",
      match_content_by_local_search: "\u0645\u0637\u0627\u0628\u0642\u0629 \u0627\u0644\u0645\u062D\u062A\u0648\u0649 \u0645\u0646 \u062E\u0644\u0627\u0644 \u0627\u0644\u0628\u062D\u062B \u0627\u0644\u0645\u062D\u0644\u064A",
      payload: "\u0627\u0644\u062D\u0645\u0648\u0644\u0629",
      pending: "\u0642\u064A\u062F \u0627\u0644\u0627\u0646\u062A\u0638\u0627\u0631",
      priority: "\u0627\u0644\u0623\u0648\u0644\u0648\u064A\u0629",
      process_orphaned_torrents_only: "\u0645\u0639\u0627\u0644\u062C\u0629 \u0627\u0644\u062A\u0648\u0631\u0646\u062A \u0627\u0644\u064A\u062A\u064A\u0645 \u0641\u0642\u0637",
      processed: "\u062A\u0645\u062A \u0627\u0644\u0645\u0639\u0627\u0644\u062C\u0629",
      purge_jobs: "\u062A\u0646\u0638\u064A\u0641 \u0627\u0644\u0648\u0638\u0627\u0626\u0641",
      purge_queue_jobs: "\u062A\u0646\u0638\u064A\u0641 \u0648\u0638\u0627\u0626\u0641 \u0627\u0644\u0637\u0627\u0628\u0648\u0631",
      queue: "\u0627\u0644\u0637\u0627\u0628\u0648\u0631",
      queue_purged: "\u062A\u0645 \u062A\u0646\u0638\u064A\u0641 \u0627\u0644\u0637\u0627\u0628\u0648\u0631",
      queues: "\u0627\u0644\u0637\u0648\u0627\u0628\u064A\u0631",
      ran_at: "\u062A\u0645 \u0627\u0644\u062A\u0634\u063A\u064A\u0644 \u0641\u064A",
      retry: "\u0625\u0639\u0627\u062F\u0629 \u0627\u0644\u0645\u062D\u0627\u0648\u0644\u0629",
      total_counts_by_status: "\u0625\u062C\u0645\u0627\u0644\u064A \u0627\u0644\u0639\u062F \u062D\u0633\u0628 \u0627\u0644\u062D\u0627\u0644\u0629"
    }
  },
  facets: {
    content_type: "\u0646\u0648\u0639 \u0627\u0644\u0645\u062D\u062A\u0648\u0649",
    file_type: "\u0646\u0648\u0639 \u0627\u0644\u0645\u0644\u0641",
    genre: "\u0627\u0644\u0646\u0648\u0639",
    language: "\u0627\u0644\u0644\u063A\u0629",
    queue: "\u0627\u0644\u0637\u0627\u0628\u0648\u0631",
    status: "\u0627\u0644\u062D\u0627\u0644\u0629",
    torrent_source: "\u0645\u0635\u062F\u0631 \u0627\u0644\u062A\u0648\u0631\u0646\u062A",
    torrent_tag: "\u0639\u0644\u0627\u0645\u0629 \u0627\u0644\u062A\u0648\u0631\u0646\u062A",
    video_resolution: "\u062F\u0642\u0629 \u0627\u0644\u0641\u064A\u062F\u064A\u0648",
    video_source: "\u0645\u0635\u062F\u0631 \u0627\u0644\u0641\u064A\u062F\u064A\u0648"
  },
  file_types: {
    archive: "\u0623\u0631\u0634\u064A\u0641",
    audio: "\u0635\u0648\u062A",
    data: "\u0628\u064A\u0627\u0646\u0627\u062A",
    document: "\u0648\u062B\u064A\u0642\u0629",
    image: "\u0635\u0648\u0631\u0629",
    software: "\u0628\u0631\u0645\u062C\u064A\u0627\u062A",
    subtitles: "\u062A\u0631\u062C\u0645\u0627\u062A",
    unknown: "\u063A\u064A\u0631 \u0645\u0639\u0631\u0648\u0641",
    video: "\u0641\u064A\u062F\u064A\u0648"
  },
  general: {
    all: "\u0627\u0644\u0643\u0644",
    dismiss: "\u0631\u0641\u0636",
    error: "\u062E\u0637\u0623",
    none: "\u0644\u0627 \u0634\u064A\u0621",
    page_not_found: "\u0627\u0644\u0635\u0641\u062D\u0629 \u063A\u064A\u0631 \u0645\u0648\u062C\u0648\u062F\u0629",
    refresh: "\u062A\u062D\u062F\u064A\u062B",
    status: "\u0627\u0644\u062D\u0627\u0644\u0629"
  },
  health: {
    bitmagnet_is_status: "bitmagnet \u0647\u0648 {{status}}",
    check_failed_with_error: "\u0641\u0634\u0644 \u0627\u0644\u062A\u062D\u0642\u0642 \u0645\u0639 \u062E\u0637\u0623",
    component: "\u0645\u0643\u0648\u0646",
    components: {
      dht: "DHT",
      postgres: "Postgres",
      tmdb: "TMDB"
    },
    dismiss: "\u0631\u0641\u0636",
    error: "\u062E\u0637\u0623",
    status: "\u0627\u0644\u062D\u0627\u0644\u0629",
    statuses: {
      degraded: "\u0645\u062A\u062F\u0647\u0648\u0631",
      down: "\u0645\u0639\u0637\u0644",
      error: "\u062E\u0637\u0623",
      inactive: "\u063A\u064A\u0631 \u0646\u0634\u0637",
      started: "\u0628\u062F\u0623",
      unknown: "\u0645\u0639\u0644\u0642",
      up: "\u0646\u0634\u0637"
    },
    summary: "\u0645\u0644\u062E\u0635 \u0627\u0644\u0635\u062D\u0629",
    worker: "\u0639\u0627\u0645\u0644",
    workers: {
      dht_crawler: "\u0632\u0627\u062D\u0641 DHT",
      http_server: "\u062E\u0627\u062F\u0645 HTTP",
      queue_server: "\u062E\u0627\u062F\u0645 \u0627\u0644\u0637\u0627\u0628\u0648\u0631"
    }
  },
  languages: {
    af: "\u0627\u0644\u0623\u0641\u0631\u064A\u0643\u0627\u0646\u064A\u0629",
    ar: "\u0627\u0644\u0639\u0631\u0628\u064A\u0629",
    az: "\u0627\u0644\u0623\u0630\u0631\u0628\u064A\u062C\u0627\u0646\u064A\u0629",
    be: "\u0627\u0644\u0628\u064A\u0644\u0627\u0631\u0648\u0633\u064A\u0629",
    bg: "\u0627\u0644\u0628\u0644\u063A\u0627\u0631\u064A\u0629",
    bs: "\u0627\u0644\u0628\u0648\u0633\u0646\u064A\u0629",
    ca: "\u0627\u0644\u0643\u0627\u062A\u0627\u0644\u0627\u0646\u064A\u0629",
    ce: "\u0627\u0644\u0634\u064A\u0634\u0627\u0646\u064A\u0629",
    co: "\u0627\u0644\u0643\u0648\u0631\u0633\u064A\u0643\u064A\u0629",
    cs: "\u0627\u0644\u062A\u0634\u064A\u0643\u064A\u0629",
    cy: "\u0627\u0644\u0648\u064A\u0644\u0632\u064A\u0629",
    da: "\u0627\u0644\u062F\u0627\u0646\u0645\u0627\u0631\u0643\u064A\u0629",
    de: "\u0627\u0644\u0623\u0644\u0645\u0627\u0646\u064A\u0629",
    el: "\u0627\u0644\u064A\u0648\u0646\u0627\u0646\u064A\u0629",
    en: "\u0627\u0644\u0625\u0646\u062C\u0644\u064A\u0632\u064A\u0629",
    es: "\u0627\u0644\u0625\u0633\u0628\u0627\u0646\u064A\u0629",
    et: "\u0627\u0644\u0625\u0633\u062A\u0648\u0646\u064A\u0629",
    eu: "\u0627\u0644\u0628\u0627\u0633\u0643\u064A\u0629",
    fa: "\u0627\u0644\u0641\u0627\u0631\u0633\u064A\u0629",
    fi: "\u0627\u0644\u0641\u0646\u0644\u0646\u062F\u064A\u0629",
    fr: "\u0627\u0644\u0641\u0631\u0646\u0633\u064A\u0629",
    he: "\u0627\u0644\u0639\u0628\u0631\u064A\u0629",
    hi: "\u0627\u0644\u0647\u0646\u062F\u064A\u0629",
    hr: "\u0627\u0644\u0643\u0631\u0648\u0627\u062A\u064A\u0629",
    hu: "\u0627\u0644\u0647\u0646\u063A\u0627\u0631\u064A\u0629",
    hy: "\u0627\u0644\u0623\u0631\u0645\u0646\u064A\u0629",
    id: "\u0627\u0644\u0625\u0646\u062F\u0648\u0646\u064A\u0633\u064A\u0629",
    is: "\u0627\u0644\u0623\u064A\u0633\u0644\u0646\u062F\u064A\u0629",
    it: "\u0627\u0644\u0625\u064A\u0637\u0627\u0644\u064A\u0629",
    ja: "\u0627\u0644\u064A\u0627\u0628\u0627\u0646\u064A\u0629",
    ka: "\u0627\u0644\u062C\u0648\u0631\u062C\u064A\u0629",
    ko: "\u0627\u0644\u0643\u0648\u0631\u064A\u0629",
    ku: "\u0627\u0644\u0643\u0631\u062F\u064A\u0629",
    lt: "\u0627\u0644\u0644\u064A\u062A\u0648\u0627\u0646\u064A\u0629",
    lv: "\u0627\u0644\u0644\u0627\u062A\u0641\u064A\u0629",
    mi: "\u0627\u0644\u0645\u0627\u0648\u0631\u064A\u0629",
    mk: "\u0627\u0644\u0645\u0642\u062F\u0648\u0646\u064A\u0629",
    ml: "\u0627\u0644\u0645\u0627\u0644\u0627\u064A\u0627\u0644\u0627\u0645\u064A\u0629",
    mn: "\u0627\u0644\u0645\u0646\u063A\u0648\u0644\u064A\u0629",
    ms: "\u0627\u0644\u0645\u0644\u0627\u064A\u0648\u064A\u0629",
    mt: "\u0627\u0644\u0645\u0627\u0644\u0637\u064A\u0629",
    nl: "\u0627\u0644\u0647\u0648\u0644\u0646\u062F\u064A\u0629",
    no: "\u0627\u0644\u0646\u0631\u0648\u064A\u062C\u064A\u0629",
    pl: "\u0627\u0644\u0628\u0648\u0644\u0646\u062F\u064A\u0629",
    pt: "\u0627\u0644\u0628\u0631\u062A\u063A\u0627\u0644\u064A\u0629",
    ro: "\u0627\u0644\u0631\u0648\u0645\u0627\u0646\u064A\u0629",
    ru: "\u0627\u0644\u0631\u0648\u0633\u064A\u0629",
    sa: "\u0627\u0644\u0633\u0646\u0633\u0643\u0631\u064A\u062A\u064A\u0629",
    sk: "\u0627\u0644\u0633\u0644\u0648\u0641\u0627\u0643\u064A\u0629",
    sl: "\u0627\u0644\u0633\u0644\u0648\u0641\u064A\u0646\u064A\u0629",
    sm: "\u0627\u0644\u0633\u0627\u0645\u0648\u064A\u0629",
    so: "\u0627\u0644\u0635\u0648\u0645\u0627\u0644\u064A\u0629",
    sr: "\u0627\u0644\u0635\u0631\u0628\u064A\u0629",
    sv: "\u0627\u0644\u0633\u0648\u064A\u062F\u064A\u0629",
    ta: "\u0627\u0644\u062A\u0627\u0645\u064A\u0644\u064A\u0629",
    th: "\u0627\u0644\u062A\u0627\u064A\u0644\u0627\u0646\u062F\u064A\u0629",
    tr: "\u0627\u0644\u062A\u0631\u0643\u064A\u0629",
    uk: "\u0627\u0644\u0623\u0648\u0643\u0631\u0627\u0646\u064A\u0629",
    vi: "\u0627\u0644\u0641\u064A\u062A\u0646\u0627\u0645\u064A\u0629",
    yi: "\u0627\u0644\u064A\u062F\u064A\u0634\u064A\u0629",
    zh: "\u0627\u0644\u0635\u064A\u0646\u064A\u0629",
    zu: "\u0627\u0644\u0632\u0648\u0644\u0648"
  },
  layout: {
    bitmagnet_on_service: "bitmagnet \u0639\u0644\u0649 {{service}}",
    change_theme: "\u062A\u063A\u064A\u064A\u0631 \u0627\u0644\u0633\u0645\u0629",
    external_links: "\u0631\u0648\u0627\u0628\u0637 \u062E\u0627\u0631\u062C\u064A\u0629",
    sponsor: "\u0627\u0644\u0631\u0627\u0639\u064A",
    support_bitmagnet: "\u062F\u0639\u0645 bitmagnet",
    translate: "\u062A\u0631\u062C\u0645\u0629"
  },
  paginator: {
    first_page: "\u0627\u0644\u0635\u0641\u062D\u0629 \u0627\u0644\u0623\u0648\u0644\u0649",
    items_per_page: "\u0627\u0644\u0639\u0646\u0627\u0635\u0631 \u0644\u0643\u0644 \u0635\u0641\u062D\u0629",
    last_page: "\u0627\u0644\u0635\u0641\u062D\u0629 \u0627\u0644\u0623\u062E\u064A\u0631\u0629",
    next_page: "\u0627\u0644\u0635\u0641\u062D\u0629 \u0627\u0644\u062A\u0627\u0644\u064A\u0629",
    page_x: "\u0627\u0644\u0635\u0641\u062D\u0629 {{x}}",
    previous_page: "\u0627\u0644\u0635\u0641\u062D\u0629 \u0627\u0644\u0633\u0627\u0628\u0642\u0629",
    x_to_y: "{{x}} \u0625\u0644\u0649 {{y}}",
    x_to_y_of_z: "{{x}} \u0625\u0644\u0649 {{y}} \u0645\u0646 {{z}}"
  },
  routes: {
    admin: "\u0627\u0644\u0645\u0633\u0624\u0648\u0644",
    dashboard: "\u0644\u0648\u062D\u0629 \u0627\u0644\u062A\u062D\u0643\u0645",
    home: "\u0627\u0644\u0635\u0641\u062D\u0629 \u0627\u0644\u0631\u0626\u064A\u0633\u064A\u0629",
    jobs: "\u0627\u0644\u0648\u0638\u0627\u0626\u0641",
    queues: "\u0627\u0644\u0637\u0648\u0627\u0628\u064A\u0631",
    torrents: "\u0627\u0644\u062A\u0648\u0631\u0646\u062A",
    visualize: "\u062A\u0635\u0648\u0631"
  },
  torrents: {
    classification: "\u0627\u0644\u062A\u0635\u0646\u064A\u0641",
    clear_search: "\u0645\u0633\u062D \u0627\u0644\u0628\u062D\u062B",
    copy: "\u0646\u0633\u062E",
    copy_to_clipboard: "\u0646\u0633\u062E \u0625\u0644\u0649 \u0627\u0644\u062D\u0627\u0641\u0638\u0629",
    delete: "\u062D\u0630\u0641",
    delete_action_cannot_be_undone: "\u0644\u0627 \u064A\u0645\u0643\u0646 \u0627\u0644\u062A\u0631\u0627\u062C\u0639 \u0639\u0646 \u0647\u0630\u0627 \u0627\u0644\u0625\u062C\u0631\u0627\u0621",
    delete_are_you_sure: "\u0647\u0644 \u0623\u0646\u062A \u0645\u062A\u0623\u0643\u062F \u0623\u0646\u0643 \u062A\u0631\u064A\u062F \u062D\u0630\u0641 \u0647\u0630\u0627 \u0627\u0644\u062A\u0648\u0631\u0646\u062A\u061F",
    deselect_all: "\u0625\u0644\u063A\u0627\u0621 \u062A\u062D\u062F\u064A\u062F \u0627\u0644\u0643\u0644",
    edit_tags: "\u062A\u062D\u0631\u064A\u0631 \u0627\u0644\u0639\u0644\u0627\u0645\u0627\u062A",
    episodes: "\u0627\u0644\u062D\u0644\u0642\u0627\u062A",
    external_links: "\u0631\u0648\u0627\u0628\u0637 \u062E\u0627\u0631\u062C\u064A\u0629",
    file_index: "\u0641\u0647\u0631\u0633 \u0627\u0644\u0645\u0644\u0641\u0627\u062A",
    file_path: "\u0645\u0633\u0627\u0631 \u0627\u0644\u0645\u0644\u0641",
    file_size: "\u062D\u062C\u0645 \u0627\u0644\u0645\u0644\u0641",
    file_type: "\u0646\u0648\u0639 \u0627\u0644\u0645\u0644\u0641",
    files: "\u0627\u0644\u0645\u0644\u0641\u0627\u062A",
    files_count_n: "{{count}} \u0645\u0644\u0641\u0627\u062A",
    files_no_info: "\u0644\u0627 \u062A\u0648\u062C\u062F \u0645\u0639\u0644\u0648\u0645\u0627\u062A \u0639\u0646 \u0627\u0644\u0645\u0644\u0641\u0627\u062A",
    files_single: "\u0645\u0644\u0641 \u0648\u0627\u062D\u062F",
    genres: "\u0627\u0644\u0623\u0646\u0648\u0627\u0639",
    info_hash: "\u062A\u062C\u0632\u0626\u0629 \u0627\u0644\u0645\u0639\u0644\u0648\u0645\u0627\u062A",
    info_hashes: "\u062A\u062C\u0632\u0626\u0627\u062A \u0627\u0644\u0645\u0639\u0644\u0648\u0645\u0627\u062A",
    languages: "\u0627\u0644\u0644\u063A\u0627\u062A",
    leechers: "\u0627\u0644\u0645\u0633\u062A\u0641\u064A\u062F\u0648\u0646",
    magnet: "\u0645\u063A\u0646\u0627\u0637\u064A\u0633",
    magnet_links: "\u0631\u0648\u0627\u0628\u0637 \u0645\u063A\u0646\u0627\u0637\u064A\u0633\u064A\u0629",
    new_tag: "\u0639\u0644\u0627\u0645\u0629 \u062C\u062F\u064A\u062F\u0629",
    order_by: "\u062A\u0631\u062A\u064A\u0628 \u062D\u0633\u0628",
    order_direction_toggle: "\u062A\u0628\u062F\u064A\u0644 \u0627\u0644\u0627\u062A\u062C\u0627\u0647",
    ordering: {
      files_count: "\u0639\u062F\u062F \u0627\u0644\u0645\u0644\u0641\u0627\u062A",
      info_hash: "\u062A\u062C\u0632\u0626\u0629 \u0627\u0644\u0645\u0639\u0644\u0648\u0645\u0627\u062A",
      leechers: "\u0627\u0644\u0645\u0633\u062A\u0641\u064A\u062F\u0648\u0646",
      name: "\u0627\u0644\u0627\u0633\u0645",
      published_at: "\u0646\u0634\u0631 \u0641\u064A",
      relevance: "\u0627\u0644\u0635\u0644\u0629",
      seeders: "\u0627\u0644\u0645\u0632\u0627\u0631\u0639\u0648\u0646",
      size: "\u0627\u0644\u062D\u062C\u0645",
      updated_at: "\u062A\u0645 \u0627\u0644\u062A\u062D\u062F\u064A\u062B \u0641\u064A"
    },
    original_release_date: "\u062A\u0627\u0631\u064A\u062E \u0627\u0644\u0625\u0635\u062F\u0627\u0631 \u0627\u0644\u0623\u0635\u0644\u064A",
    permalink: "\u0631\u0627\u0628\u0637 \u062F\u0627\u0626\u0645",
    poster: "\u0645\u0644\u0635\u0642",
    published: "\u0645\u0646\u0634\u0648\u0631",
    rating: "\u0627\u0644\u062A\u0642\u064A\u064A\u0645",
    refresh: "\u062A\u062D\u062F\u064A\u062B \u0627\u0644\u0646\u062A\u0627\u0626\u062C",
    reprocess: {
      force_rematch: "\u0641\u0631\u0636 \u0625\u0639\u0627\u062F\u0629 \u0627\u0644\u0645\u0637\u0627\u0628\u0642\u0629 \u0644\u0644\u0645\u062D\u062A\u0648\u0649 \u0627\u0644\u0645\u0637\u0627\u0628\u0642 \u0628\u0627\u0644\u0641\u0639\u0644",
      match_content_by_external_api_search: "\u0645\u0637\u0627\u0628\u0642\u0629 \u0627\u0644\u0645\u062D\u062A\u0648\u0649 \u0645\u0646 \u062E\u0644\u0627\u0644 \u0627\u0644\u0628\u062D\u062B \u0641\u064A API \u0627\u0644\u062E\u0627\u0631\u062C\u064A\u0629",
      match_content_by_local_search: "\u0645\u0637\u0627\u0628\u0642\u0629 \u0627\u0644\u0645\u062D\u062A\u0648\u0649 \u0645\u0646 \u062E\u0644\u0627\u0644 \u0627\u0644\u0628\u062D\u062B \u0627\u0644\u0645\u062D\u0644\u064A",
      reprocess: "\u0625\u0639\u0627\u062F\u0629 \u0627\u0644\u0645\u0639\u0627\u0644\u062C\u0629"
    },
    s_l: "S / L",
    search: "\u0628\u062D\u062B",
    seeders: "\u0627\u0644\u0628\u0627\u0630\u0631\u0648\u0646",
    select_all: "\u062A\u062D\u062F\u064A\u062F \u0627\u0644\u0643\u0644",
    showing_x_of_y_files: "\u0639\u0631\u0636 {{x}} \u0645\u0646 {{y}} \u0645\u0644\u0641\u0627\u062A",
    size: "\u0627\u0644\u062D\u062C\u0645",
    source: "\u0645\u0635\u062F\u0631 \u0627\u0644\u062A\u0648\u0631\u0646\u062A",
    summary: "\u0627\u0644\u0645\u0644\u062E\u0635",
    tags: {
      delete: "\u062D\u0630\u0641 \u0627\u0644\u0639\u0644\u0627\u0645\u0627\u062A",
      delete_tip: "\u0625\u0632\u0627\u0644\u0629 \u0627\u0644\u0639\u0644\u0627\u0645\u0627\u062A \u0645\u0646 \u0645\u0644\u0641\u0627\u062A \u0627\u0644\u062A\u0648\u0631\u0646\u062A \u0627\u0644\u0645\u062D\u062F\u062F\u0629",
      placeholder: "\u0627\u0644\u0639\u0644\u0627\u0645\u0629...",
      put: "\u0648\u0636\u0639 \u0627\u0644\u0639\u0644\u0627\u0645\u0627\u062A",
      put_tip: "\u0625\u0636\u0627\u0641\u0629 \u0627\u0644\u0639\u0644\u0627\u0645\u0627\u062A \u0625\u0644\u0649 \u0645\u0644\u0641\u0627\u062A \u0627\u0644\u062A\u0648\u0631\u0646\u062A \u0627\u0644\u0645\u062D\u062F\u062F\u0629",
      set: "\u062A\u0639\u064A\u064A\u0646 \u0627\u0644\u0639\u0644\u0627\u0645\u0627\u062A",
      set_tip: "\u0627\u0633\u062A\u0628\u062F\u0627\u0644 \u0627\u0644\u0639\u0644\u0627\u0645\u0627\u062A \u0641\u064A \u0645\u0644\u0641\u0627\u062A \u0627\u0644\u062A\u0648\u0631\u0646\u062A \u0627\u0644\u0645\u062D\u062F\u062F\u0629"
    },
    title: "\u0627\u0644\u0639\u0646\u0648\u0627\u0646",
    toggle_drawer: "\u062A\u0628\u062F\u064A\u0644 \u0627\u0644\u062F\u0631\u062C",
    votes_count_n: "{{count}} \u0623\u0635\u0648\u0627\u062A"
  },
  version: {
    bitmagnet_version: "\u0625\u0635\u062F\u0627\u0631 bitmagnet {{version}}",
    unknown: "\u063A\u064A\u0631 \u0645\u0639\u0631\u0648\u0641"
  }
};

// src/app/i18n/translations/ca.json
var ca_default = {
  content_types: {
    plural: {
      all: "Tot",
      audiobook: "Audiollibres",
      comic: "C\xF2mics",
      ebook: "Llibres electr\xF2nics",
      movie: "Pel\xB7l\xEDcules",
      music: "M\xFAsica",
      null: "Desconegut",
      software: "Programari",
      tv_show: "Programes de TV",
      xxx: "XXX"
    },
    singular: {
      audiobook: "Audiollibre",
      comic: "C\xF2mic",
      ebook: "Llibre electr\xF2nic",
      movie: "Pel\xB7l\xEDcula",
      music: "M\xFAsica",
      software: "Desconegut",
      tv_show: "Programa de TV",
      xxx: "XXX"
    }
  },
  dashboard: {
    interval: {
      all: "Tot",
      days: "Dies",
      days_1: "1 dia",
      hours: "Hores",
      hours_1: "1 hora",
      hours_12: "12 hores",
      hours_6: "6 hores",
      minutes: "Minuts",
      minutes_1: "1 minut",
      minutes_15: "15 minuts",
      minutes_30: "30 minuts",
      minutes_5: "5 minuts",
      off: "Apagat",
      seconds_10: "10 segons",
      seconds_30: "30 segons",
      weeks_1: "1 setmana"
    },
    metrics: {
      event: "Esdeveniment",
      resolution: "Resoluci\xF3",
      throughput: "Rendiment",
      timeframe: "Per\xEDode de temps",
      toggle_legend: "Mostra o oculta la llegenda"
    },
    queues: {
      created: "Creat",
      created_at: "Creaci\xF3",
      enqueue_jobs: "Encua les tasques",
      enqueue_torrent_processing_batch: "Encua el lot de processament de torrents",
      failed: "Fallat",
      jobs_enqueued: "Tasques encuades",
      latency: "Lat\xE8ncia",
      payload: "Contingut",
      priority: "Prioritat",
      process_orphaned_torrents_only: "Processa nom\xE9s els torrents orfes",
      processed: "Processat",
      purge_jobs: "Purga les tasques",
      purge_queue_jobs: "Purga les cues de tasques",
      queue: "Cua",
      queue_purged: "Cua purgada",
      queues: "Cues",
      ran_at: "Executat a",
      total_counts_by_status: "Recompte total per estat"
    }
  },
  facets: {
    content_type: "Tipus de Contingut",
    file_type: "Tipus de Fitxer",
    genre: "G\xE8nere",
    language: "Idioma",
    torrent_source: "Origen del Torrent",
    torrent_tag: "Etiqueta del Torrent",
    video_resolution: "Resoluci\xF3 del V\xEDdeo",
    video_source: "Origen del V\xEDdeo"
  },
  file_types: {
    archive: "Arxiu",
    audio: "\xC0udio",
    data: "Dades",
    document: "Document",
    image: "Imatge",
    software: "Programari",
    subtitles: "Subt\xEDtols",
    unknown: "Desconegut",
    video: "V\xEDdeo"
  },
  general: {
    all: "Tot",
    dismiss: "Descarta",
    error: "Error",
    none: "Cap",
    page_not_found: "P\xE0gina no trobada",
    refresh: "Actualitza",
    status: "Estat"
  },
  health: {
    bitmagnet_is_status: "bitmagnet est\xE0 {{status}}",
    check_failed_with_error: "Ha fallat la comprovaci\xF3 amb un error",
    component: "Component",
    components: {
      dht: "DHT",
      postgres: "Postgres",
      tmdb: "TMDB"
    },
    statuses: {
      degraded: "Degradat",
      down: "Caigut",
      error: "Error",
      inactive: "Inactiu",
      started: "Iniciat",
      unknown: "Pendent",
      up: "Actiu"
    },
    summary: "Resum de salut",
    worker: "Treballador",
    workers: {
      dht_crawler: "Rastrejador DHT",
      http_server: "Servidor HTTP",
      queue_server: "Servidor de cues"
    }
  },
  languages: {
    af: "Afrikaans",
    ar: "\xC0rab",
    az: "\xC0zeri",
    be: "Belar\xFAs",
    bg: "B\xFAlgar",
    bs: "Bosni\xE0",
    ca: "Catal\xE0",
    ce: "Txetx\xE8",
    co: "Cors",
    cs: "Txec",
    cy: "Gal\xB7l\xE8s",
    da: "Dan\xE8s",
    de: "Alemany",
    el: "Grec",
    en: "Angl\xE8s",
    es: "Castell\xE0",
    et: "Estoni\xE0",
    eu: "Basc",
    fa: "Persa",
    fi: "Fin\xE8s",
    fr: "Franc\xE8s",
    he: "Hebreu",
    hi: "Hindi",
    hr: "Croat",
    hu: "Hongar\xE8s",
    hy: "Armeni",
    id: "Indonesi",
    is: "Island\xE8s",
    it: "Itali\xE0",
    ja: "Japon\xE8s",
    ka: "Georgi\xE0",
    ko: "Core\xE0",
    ku: "Kurd",
    lt: "Litu\xE0",
    lv: "Let\xF3",
    mi: "Maori",
    mk: "Maced\xF2nic",
    ml: "Malai\xE0lam",
    mn: "Mongol",
    ms: "Malai",
    mt: "Malt\xE8s",
    nl: "Neerland\xE8s",
    no: "Noruec",
    pl: "Polon\xE8s",
    pt: "Portugu\xE8s",
    ro: "Roman\xE8s",
    ru: "Rus",
    sa: "S\xE0nscrit",
    sk: "Eslovac",
    sl: "Eslov\xE8",
    sm: "Samo\xE0",
    so: "Somali",
    sr: "Serbi",
    sv: "Suec",
    ta: "T\xE0mil",
    th: "Tailand\xE8s",
    tr: "Turc",
    uk: "Ucra\xEFn\xE8s",
    vi: "Vietnamita",
    yi: "\xCDdix",
    zh: "Xin\xE8s",
    zu: "Zul\xFA"
  },
  layout: {
    bitmagnet_on_service: "bitmagnet a {{service}}",
    change_theme: "Canviar el tema",
    external_links: "Enlla\xE7os externs",
    sponsor: "Patrocina",
    support_bitmagnet: "Ajuda a bitmagnet",
    translate: "Tradueix"
  },
  paginator: {
    first_page: "Primera p\xE0gina",
    last_page: "Elements per p\xE0gina",
    next_page: "P\xE0gina seg\xFCent",
    page_x: "P\xE0gina {{x}}",
    previous_page: "P\xE0gina anterior",
    x_to_y: "{{x}} a {{y}}",
    x_to_y_of_z: "{{x}} a {{y}} de {{z}}"
  },
  routes: {
    admin: "Administraci\xF3",
    dashboard: "Tauler de control",
    home: "Inici",
    jobs: "Tasques",
    queues: "Cues",
    torrents: "Torrents",
    visualize: "Visualitza"
  },
  torrents: {
    classification: "Classificaci\xF3",
    clear_search: "Esborrar la cerca",
    copy: "Copiar",
    copy_to_clipboard: "Copia al porta-retalls",
    delete: "Esborrar",
    delete_action_cannot_be_undone: "Aquesta acci\xF3 no es pot desfer",
    delete_are_you_sure: "Est\xE0s segur que vols esborrar aquest torrent?",
    deselect_all: "Desselecciona-ho tot",
    edit_tags: "Edita l'etiqueta",
    episodes: "Episodis",
    external_links: "Enlla\xE7os externs",
    file_index: "\xCDndex del fitxer",
    file_path: "Cam\xED del fitxer",
    file_size: "Mida del fitxer",
    file_type: "Tipus de fitxer",
    files: "Fitxers",
    files_no_info: "{{count}} fitxers",
    genres: "G\xE8neres",
    info_hash: "Informaci\xF3 del resum",
    info_hashes: "Informaci\xF3 dels resums",
    languages: "Idiomes",
    leechers: "Sangoneres",
    magnet: "Magnet",
    magnet_links: "Enlla\xE7os magnet",
    new_tag: "Etiqueta nova",
    order_by: "Ordenar per",
    order_direction_toggle: "Commuta la direcci\xF3",
    ordering: {
      files_count: "Recompte de fitxers",
      info_hash: "Informaci\xF3 del resum",
      leechers: "Sangoneres",
      name: "Nom",
      published_at: "Publicat el",
      relevance: "Rellev\xE0ncia",
      seeders: "Sembradors",
      size: "Mida",
      updated_at: "Actualitzat el"
    },
    original_release_date: "Data de llan\xE7ament original",
    permalink: "Enlla\xE7 permanent",
    poster: "P\xF2ster",
    published: "Publicat",
    rating: "Valoraci\xF3",
    refresh: "Actualitza els resultats",
    reprocess: {
      force_rematch: "For\xE7ar nova coincid\xE8ncia del contingut ja coincident",
      match_content_by_external_api_search: "Fes coincidir el contingut a trav\xE9s d'una cerca d'API externa",
      match_content_by_local_search: "Fes conicidir a trav\xE9s d'una cerca local",
      reprocess: "Tornar a processar"
    },
    s_l: "S / S",
    search: "Cerca",
    seeders: "Sembradors",
    select_all: "Seleccionar-ho tot",
    showing_x_of_y_files: "Mostrant {{x}} de {{y}} fitxers",
    size: "Mida",
    source: "Origen del torrent",
    summary: "Resum",
    tags: {
      delete: "Esborrar etiquetes",
      delete_tip: "Esborrar etiquetes dels torrents seleccionats",
      placeholder: "Etiqueta...",
      put: "Afegeix etiquetes",
      put_tip: "Afegeix etiquetes als torrents seleccionats",
      set: "Reempla\xE7a etiquetes",
      set_tip: "Reempla\xE7a les etiquetes dels torrents seleccionats"
    },
    title: "T\xEDtol",
    toggle_drawer: "Mostra o oculta el calaix",
    votes_count_n: "{{count}} vots"
  },
  version: {
    bitmagnet_version: "versi\xF3 de bitmagnet {{version}}",
    unknown: "desconegut"
  }
};

// src/app/i18n/translations/de.json
var de_default = {
  content_types: {
    plural: {
      all: "Alle",
      audiobook: "H\xF6rb\xFCcher",
      comic: "Comics",
      ebook: "E-B\xFCcher",
      game: "Spiele",
      movie: "Filme",
      music: "Musik",
      null: "Unbekannt",
      software: "Software",
      tv_show: "TV-Shows",
      xxx: "XXX"
    },
    singular: {
      audiobook: "H\xF6rbuch",
      comic: "Comic",
      ebook: "E-Buch",
      game: "Spiel",
      movie: "Film",
      music: "Musik",
      null: "Unbekannt",
      software: "Software",
      tv_show: "TV-Show",
      xxx: "XXX"
    }
  },
  dashboard: {
    event: {
      created: "Erstellt",
      failed: "Fehlgeschlagen",
      processed: "Verarbeitet",
      updated: "Aktualisiert"
    },
    interval: {
      all: "Alle",
      days: "Tage",
      days_1: "1 Tag",
      hours: "Stunden",
      hours_1: "1 Stunde",
      hours_12: "12 Stunden",
      hours_6: "6 Stunden",
      minutes: "Minuten",
      minutes_1: "1 Minute",
      minutes_15: "15 Minuten",
      minutes_30: "30 Minuten",
      minutes_5: "5 Minuten",
      off: "Aus",
      seconds_10: "10 Sekunden",
      seconds_30: "30 Sekunden",
      weeks_1: "1 Woche"
    },
    metrics: {
      event: "Ereignis",
      resolution: "Aufl\xF6sung",
      throughput: "Durchsatz",
      timeframe: "Zeitrahmen",
      toggle_legend: "Legende umschalten"
    },
    queues: {
      created: "Erstellt",
      created_at: "Erstellt am",
      enqueue_jobs: "Jobs einreihen",
      enqueue_torrent_processing_batch: "Torrent-Verarbeitungsbatch einreihen",
      failed: "Fehlgeschlagen",
      force_rematch: "Erneutes Zuordnen von bereits zugeordneten Inhalten erzwingen",
      jobs_enqueued: "Jobs eingereiht",
      latency: "Latenz",
      match_content_by_external_api_search: "Inhalte durch externe API-Suche zuordnen",
      match_content_by_local_search: "Inhalte durch lokale Suche zuordnen",
      payload: "Nutzlast",
      pending: "Ausstehend",
      priority: "Priorit\xE4t",
      process_orphaned_torrents_only: "Nur verwaiste Torrents verarbeiten",
      processed: "Verarbeitet",
      purge_jobs: "Jobs bereinigen",
      purge_queue_jobs: "Warteschlangen-Jobs bereinigen",
      queue: "Warteschlange",
      queue_purged: "Warteschlange bereinigt",
      queues: "Warteschlangen",
      ran_at: "Ausgef\xFChrt am",
      retry: "Erneut versuchen",
      total_counts_by_status: "Gesamtanzahl nach Status"
    }
  },
  facets: {
    content_type: "Inhaltstyp",
    file_type: "Dateityp",
    genre: "Genre",
    language: "Sprache",
    queue: "Warteschlange",
    status: "Status",
    torrent_source: "Torrent-Quelle",
    torrent_tag: "Torrent-Tag",
    video_resolution: "Videoaufl\xF6sung",
    video_source: "Videoquelle"
  },
  file_types: {
    archive: "Archiv",
    audio: "Audio",
    data: "Daten",
    document: "Dokument",
    image: "Bild",
    software: "Software",
    subtitles: "Untertitel",
    unknown: "Unbekannt",
    video: "Video"
  },
  general: {
    all: "Alle",
    dismiss: "Verwerfen",
    error: "Fehler",
    none: "Keine",
    page_not_found: "Seite nicht gefunden",
    refresh: "Aktualisieren",
    status: "Status"
  },
  health: {
    bitmagnet_is_status: "bitmagnet ist {{status}}",
    check_failed_with_error: "\xDCberpr\xFCfung mit Fehler fehlgeschlagen",
    component: "Komponente",
    components: {
      dht: "DHT",
      postgres: "Postgres",
      tmdb: "TMDB"
    },
    dismiss: "Verwerfen",
    error: "Fehler",
    status: "Status",
    statuses: {
      degraded: "Verschlechtert",
      down: "Aus",
      error: "Fehler",
      inactive: "Inaktiv",
      started: "Gestartet",
      unknown: "Ausstehend",
      up: "An"
    },
    summary: "Gesundheits\xFCbersicht",
    worker: "Arbeiter",
    workers: {
      dht_crawler: "DHT-Crawler",
      http_server: "HTTP-Server",
      queue_server: "Warteschlangen-Server"
    }
  },
  languages: {
    af: "Afrikaans",
    ar: "Arabisch",
    az: "Aserbaidschanisch",
    be: "Wei\xDFrussisch",
    bg: "Bulgarisch",
    bs: "Bosnisch",
    ca: "Katalanisch",
    ce: "Tschetschenisch",
    co: "Korsisch",
    cs: "Tschechisch",
    cy: "Walisisch",
    da: "D\xE4nisch",
    de: "Deutsch",
    el: "Griechisch",
    en: "Englisch",
    es: "Spanisch",
    et: "Estnisch",
    eu: "Baskisch",
    fa: "Persisch",
    fi: "Finnisch",
    fr: "Franz\xF6sisch",
    he: "Hebr\xE4isch",
    hi: "Hindi",
    hr: "Kroatisch",
    hu: "Ungarisch",
    hy: "Armenisch",
    id: "Indonesisch",
    is: "Isl\xE4ndisch",
    it: "Italienisch",
    ja: "Japanisch",
    ka: "Georgisch",
    ko: "Koreanisch",
    ku: "Kurdisch",
    lt: "Litauisch",
    lv: "Lettisch",
    mi: "Maori",
    mk: "Mazedonisch",
    ml: "Malayalam",
    mn: "Mongolisch",
    ms: "Malaiisch",
    mt: "Maltesisch",
    nl: "Niederl\xE4ndisch",
    no: "Norwegisch",
    pl: "Polnisch",
    pt: "Portugiesisch",
    ro: "Rum\xE4nisch",
    ru: "Russisch",
    sa: "Sanskrit",
    sk: "Slowakisch",
    sl: "Slowenisch",
    sm: "Samoanisch",
    so: "Somalisch",
    sr: "Serbisch",
    sv: "Schwedisch",
    ta: "Tamil",
    th: "Thai",
    tr: "T\xFCrkisch",
    uk: "Ukrainisch",
    vi: "Vietnamesisch",
    yi: "Jiddisch",
    zh: "Chinesisch",
    zu: "Zulu"
  },
  layout: {
    bitmagnet_on_service: "bitmagnet auf {{service}}",
    change_theme: "Thema \xE4ndern",
    external_links: "Externe Links",
    sponsor: "Sponsor",
    support_bitmagnet: "bitmagnet unterst\xFCtzen",
    translate: "\xDCbersetzen"
  },
  paginator: {
    first_page: "Erste Seite",
    items_per_page: "Artikel pro Seite",
    last_page: "Letzte Seite",
    next_page: "N\xE4chste Seite",
    page_x: "Seite {{x}}",
    previous_page: "Vorherige Seite",
    x_to_y: "{{x}} bis {{y}}",
    x_to_y_of_z: "{{x}} bis {{y}} von {{z}}"
  },
  routes: {
    admin: "Admin",
    dashboard: "Dashboard",
    home: "Startseite",
    jobs: "Aufgaben",
    queues: "Warteschlangen",
    torrents: "Torrents",
    visualize: "Visualisieren"
  },
  torrents: {
    classification: "Klassifizierung",
    clear_search: "Suche l\xF6schen",
    copy: "Kopieren",
    copy_to_clipboard: "In die Zwischenablage kopieren",
    delete: "L\xF6schen",
    delete_action_cannot_be_undone: "Diese Aktion kann nicht r\xFCckg\xE4ngig gemacht werden",
    delete_are_you_sure: "Sind Sie sicher, dass Sie diesen Torrent l\xF6schen m\xF6chten?",
    deselect_all: "Alle abw\xE4hlen",
    edit_tags: "Tags bearbeiten",
    episodes: "Episoden",
    external_links: "Externe Links",
    file_index: "Dateiindex",
    file_path: "Dateipfad",
    file_size: "Dateigr\xF6\xDFe",
    file_type: "Dateityp",
    files: "Dateien",
    files_count_n: "{{count}} Dateien",
    files_no_info: "Keine Dateiinformationen verf\xFCgbar",
    files_single: "Einzelne Datei",
    genres: "Genres",
    info_hash: "Info-Hash",
    info_hashes: "Info-Hashes",
    languages: "Sprachen",
    leechers: "Leechers",
    magnet: "Magnet",
    magnet_links: "Magnet-Links",
    new_tag: "Neuer Tag",
    order_by: "Sortieren nach",
    order_direction_toggle: "Richtung umschalten",
    ordering: {
      files_count: "Dateienanzahl",
      info_hash: "Info-Hash",
      leechers: "Leechers",
      name: "Name",
      published_at: "Ver\xF6ffentlicht am",
      relevance: "Relevanz",
      seeders: "Seeders",
      size: "Gr\xF6\xDFe",
      updated_at: "Aktualisiert am"
    },
    original_release_date: "Originales Ver\xF6ffentlichungsdatum",
    permalink: "Permalink",
    poster: "Poster",
    published: "Ver\xF6ffentlicht",
    rating: "Bewertung",
    refresh: "Ergebnisse aktualisieren",
    reprocess: {
      force_rematch: "Erneutes Zuordnen von bereits zugeordneten Inhalten erzwingen",
      match_content_by_external_api_search: "Inhalte durch externe API-Suche zuordnen",
      match_content_by_local_search: "Inhalte durch lokale Suche zuordnen",
      reprocess: "Erneut verarbeiten"
    },
    s_l: "S / L",
    search: "Suche",
    seeders: "Seeders",
    select_all: "Alle ausw\xE4hlen",
    showing_x_of_y_files: "{{x}} von {{y}} Dateien anzeigen",
    size: "Gr\xF6\xDFe",
    source: "Torrent-Quelle",
    summary: "Zusammenfassung",
    tags: {
      delete: "Tags l\xF6schen",
      delete_tip: "Tags aus den ausgew\xE4hlten Torrents entfernen",
      placeholder: "Tag...",
      put: "Tags setzen",
      put_tip: "Tags zu den ausgew\xE4hlten Torrents hinzuf\xFCgen",
      set: "Tags setzen",
      set_tip: "Tags der ausgew\xE4hlten Torrents ersetzen"
    },
    title: "Titel",
    toggle_drawer: "Schublade umschalten",
    votes_count_n: "{{count}} Stimmen"
  },
  version: {
    bitmagnet_version: "bitmagnet Version {{version}}",
    unknown: "unbekannt"
  }
};

// src/app/i18n/translations/en.json
var en_default = {
  content_types: {
    plural: {
      all: "All",
      audiobook: "Audiobooks",
      comic: "Comics",
      ebook: "E-Books",
      game: "Games",
      movie: "Movies",
      music: "Music",
      null: "Unknown",
      software: "Software",
      tv_show: "TV Shows",
      xxx: "XXX"
    },
    singular: {
      audiobook: "Audiobook",
      comic: "Comic",
      ebook: "E-Book",
      game: "Game",
      movie: "Movie",
      music: "Music",
      null: "Unknown",
      software: "Software",
      tv_show: "TV Show",
      xxx: "XXX"
    }
  },
  dashboard: {
    event: {
      created: "Created",
      failed: "Failed",
      processed: "Processed",
      updated: "Updated"
    },
    interval: {
      all: "All",
      days: "Days",
      days_1: "1 day",
      hours: "Hours",
      hours_1: "1 hour",
      hours_12: "12 hours",
      hours_6: "6 hours",
      minutes: "Minutes",
      minutes_1: "1 minute",
      minutes_15: "15 minutes",
      minutes_30: "30 minutes",
      minutes_5: "5 minutes",
      off: "Off",
      seconds_10: "10 seconds",
      seconds_30: "30 seconds",
      weeks_1: "1 week"
    },
    metrics: {
      event: "Event",
      resolution: "Resolution",
      throughput: "Throughput",
      timeframe: "Timeframe",
      toggle_legend: "Toggle legend"
    },
    queues: {
      created: "Created",
      created_at: "Created at",
      enqueue_jobs: "Enqueue jobs",
      enqueue_torrent_processing_batch: "Enqueue Torrent Processing Batch",
      failed: "Failed",
      jobs_enqueued: "Jobs enqueued",
      latency: "Latency",
      payload: "Payload",
      pending: "Pending",
      priority: "Priority",
      process_orphaned_torrents_only: "Process orphaned torrents only",
      processed: "Processed",
      purge_jobs: "Purge jobs",
      purge_queue_jobs: "Purge queue jobs",
      queue: "Queue",
      queue_purged: "Queue purged",
      queues: "Queues",
      ran_at: "Ran at",
      retry: "Retry",
      total_counts_by_status: "Total counts by status"
    }
  },
  facets: {
    content_type: "Content Type",
    file_type: "File Type",
    genre: "Genre",
    language: "Language",
    queue: "Queue",
    status: "Status",
    torrent_source: "Torrent Source",
    torrent_tag: "Torrent Tag",
    video_resolution: "Video Resolution",
    video_source: "Video Source"
  },
  file_types: {
    archive: "Archive",
    audio: "Audio",
    data: "Data",
    document: "Document",
    image: "Image",
    software: "Software",
    subtitles: "Subtitles",
    unknown: "Unknown",
    video: "Video"
  },
  general: {
    all: "All",
    confirm: "Confirm",
    cancel: "Cancel",
    dismiss: "Dismiss",
    error: "Error",
    none: "None",
    page_not_found: "Page not found",
    refresh: "Refresh",
    status: "Status"
  },
  health: {
    bitmagnet_is_status: "bitmagnet is {{status}}",
    check_failed_with_error: "Check failed with error",
    component: "Component",
    components: {
      dht: "DHT",
      postgres: "Postgres",
      tmdb: "TMDB",
      workers: "Workers"
    },
    dismiss: "Dismiss",
    error: "Error",
    status: "Status",
    statuses: {
      degraded: "Degraded",
      down: "Down",
      error: "Error",
      inactive: "Inactive",
      unknown: "Pending",
      up: "Up"
    },
    summary: "Health summary"
  },
  languages: {
    af: "Afrikaans",
    ar: "Arabic",
    az: "Azerbaijani",
    be: "Belarusian",
    bg: "Bulgarian",
    bs: "Bosnian",
    ca: "Catalan",
    ce: "Chechen",
    co: "Corsican",
    cs: "Czech",
    cy: "Welsh",
    da: "Danish",
    de: "German",
    el: "Greek",
    en: "English",
    es: "Spanish",
    et: "Estonian",
    eu: "Basque",
    fa: "Persian",
    fi: "Finnish",
    fr: "French",
    he: "Hebrew",
    hi: "Hindi",
    hr: "Croatian",
    hu: "Hungarian",
    hy: "Armenian",
    id: "Indonesian",
    is: "Icelandic",
    it: "Italian",
    ja: "Japanese",
    ka: "Georgian",
    ko: "Korean",
    ku: "Kurdish",
    lt: "Lithuanian",
    lv: "Latvian",
    mi: "Maori",
    mk: "Macedonian",
    ml: "Malayalam",
    mn: "Mongolian",
    ms: "Malay",
    mt: "Maltese",
    nl: "Dutch",
    no: "Norwegian",
    pl: "Polish",
    pt: "Portuguese",
    ro: "Romanian",
    ru: "Russian",
    sa: "Sanskrit",
    sk: "Slovak",
    sl: "Slovenian",
    sm: "Samoan",
    so: "Somali",
    sr: "Serbian",
    sv: "Swedish",
    ta: "Tamil",
    th: "Thai",
    tr: "Turkish",
    uk: "Ukrainian",
    vi: "Vietnamese",
    yi: "Yiddish",
    zh: "Chinese",
    zu: "Zulu"
  },
  layout: {
    bitmagnet_on_service: "bitmagnet on {{service}}",
    change_theme: "Change theme",
    external_links: "External links",
    sponsor: "Sponsor",
    support_bitmagnet: "Support bitmagnet",
    translate: "Translate"
  },
  paginator: {
    first_page: "First page",
    items_per_page: "Items per page",
    last_page: "Last page",
    next_page: "Next page",
    page_x: "Page {{x}}",
    previous_page: "Previous page",
    x_to_y: "{{x}} to {{y}}",
    x_to_y_of_z: "{{x}} to {{y}} of {{z}}"
  },
  routes: {
    admin: "Admin",
    dashboard: "Dashboard",
    home: "Home",
    jobs: "Jobs",
    queues: "Queues",
    torrents: "Torrents",
    visualize: "Visualize",
    workers: "Workers"
  },
  torrents: {
    classification: "Classification",
    clear_search: "Clear Search",
    copy: "Copy",
    copy_to_clipboard: "Copy to clipboard",
    delete: "Delete",
    delete_action_cannot_be_undone: "This action cannot be undone",
    delete_are_you_sure: "Are you sure you want to delete this torrent?",
    deselect_all: "Deselect All",
    edit_tags: "Edit tags",
    episodes: "Episodes",
    external_links: "External links",
    file_index: "File index",
    file_path: "File path",
    file_size: "File size",
    file_type: "File type",
    files: "Files",
    files_count_n: "{{count}} files",
    files_no_info: "No files information available",
    files_single: "Single file",
    genres: "Genres",
    info_hash: "Info hash",
    info_hashes: "Info hashes",
    languages: "Languages",
    leechers: "Leechers",
    magnet: "Magnet",
    magnet_links: "Magnet links",
    new_tag: "New tag",
    order_by: "Order by",
    order_direction_toggle: "Toggle direction",
    ordering: {
      files_count: "Files count",
      info_hash: "Info hash",
      leechers: "Leechers",
      name: "Name",
      published_at: "Published at",
      relevance: "Relevance",
      seeders: "Seeders",
      size: "Size",
      updated_at: "Updated at"
    },
    original_release_date: "Original release date",
    permalink: "Permalink",
    poster: "Poster",
    published: "Published",
    rating: "Rating",
    refresh: "Refresh results",
    reprocess: {
      force_rematch: "Force rematch of already matched content",
      match_content_by_external_api_search: "Match content by external API search",
      match_content_by_local_search: "Match content by local search",
      reprocess: "Reprocess"
    },
    s_l: "S / L",
    search: "Search",
    seeders: "Seeders",
    select_all: "Select All",
    showing_x_of_y_files: "Showing {{x}} of {{y}} files",
    size: "Size",
    source: "Torrent Source",
    summary: "Summary",
    tags: {
      delete: "Delete tags",
      delete_tip: "Remove tags from the selected torrents",
      placeholder: "Tag...",
      put: "Put tags",
      put_tip: "Add tags to the selected torrents",
      set: "Set tags",
      set_tip: "Replace tags of the selected torrents"
    },
    title: "Title",
    toggle_drawer: "Toggle Drawer",
    votes_count_n: "{{count}} votes"
  },
  version: {
    bitmagnet_version: "bitmagnet version {{version}}",
    unknown: "unknown"
  },
  workers: {
    worker: "Worker",
    names: {
      dht_crawler: "DHT crawler",
      http_server: "HTTP server",
      queue_server: "Queue server",
      blocker: "Blocker",
      database: "Database",
      dht_server: "DHT server",
      dht_socket: "DHT socket",
      health_checker: "Health checker",
      logger: "Logger"
    },
    states: {
      idle: "Idle",
      startup: "Startup",
      running: "Running",
      shutdown: "Shutdown",
      error: "Error"
    },
    actions: {
      start: "Start",
      shutdown: "Shutdown",
      restart: "Restart"
    },
    warnings: {
      depends_on: "This action will also start the following services:",
      required_by: "This action will affect the following dependent services:"
    }
  }
};

// src/app/i18n/translations/es.json
var es_default = {
  content_types: {
    plural: {
      all: "Todos",
      audiobook: "Audiolibros",
      comic: "C\xF3mics",
      ebook: "E-Libros",
      game: "Juegos",
      movie: "Pel\xEDculas",
      music: "M\xFAsica",
      null: "Desconocido",
      software: "Software",
      tv_show: "Programas de TV",
      xxx: "XXX"
    },
    singular: {
      audiobook: "Audiolibro",
      comic: "C\xF3mic",
      ebook: "E-Libro",
      game: "Juego",
      movie: "Pel\xEDcula",
      music: "M\xFAsica",
      null: "Desconocido",
      software: "Software",
      tv_show: "Programa de TV",
      xxx: "XXX"
    }
  },
  dashboard: {
    event: {
      created: "Creado",
      failed: "Fallido",
      processed: "Procesado",
      updated: "Actualizado"
    },
    interval: {
      all: "Todos",
      days: "D\xEDas",
      days_1: "1 d\xEDa",
      hours: "Horas",
      hours_1: "1 hora",
      hours_12: "12 horas",
      hours_6: "6 horas",
      minutes: "Minutos",
      minutes_1: "1 minuto",
      minutes_15: "15 minutos",
      minutes_30: "30 minutos",
      minutes_5: "5 minutos",
      off: "Apagado",
      seconds_10: "10 segundos",
      seconds_30: "30 segundos",
      weeks_1: "1 semana"
    },
    metrics: {
      event: "Evento",
      resolution: "Resoluci\xF3n",
      throughput: "Rendimiento",
      timeframe: "Periodo de tiempo",
      toggle_legend: "Alternar leyenda"
    },
    queues: {
      created: "Creado",
      created_at: "Creado en",
      enqueue_jobs: "Encolar trabajos",
      enqueue_torrent_processing_batch: "Encolar lote de procesamiento de torrents",
      failed: "Fallido",
      force_rematch: "Forzar nueva coincidencia de contenido ya coincidente",
      jobs_enqueued: "Trabajos encolados",
      latency: "Latencia",
      match_content_by_external_api_search: "Coincidir contenido por b\xFAsqueda de API externa",
      match_content_by_local_search: "Coincidir contenido por b\xFAsqueda local",
      payload: "Carga \xFAtil",
      pending: "Pendiente",
      priority: "Prioridad",
      process_orphaned_torrents_only: "Procesar solo torrents hu\xE9rfanos",
      processed: "Procesado",
      purge_jobs: "Purgar trabajos",
      purge_queue_jobs: "Purgar trabajos de la cola",
      queue: "Cola",
      queue_purged: "Cola purgada",
      queues: "Colas",
      ran_at: "Ejecutado en",
      retry: "Reintentar",
      total_counts_by_status: "Recuento total por estado"
    }
  },
  facets: {
    content_type: "Tipo de contenido",
    file_type: "Tipo de archivo",
    genre: "G\xE9nero",
    language: "Idioma",
    queue: "Cola",
    status: "Estado",
    torrent_source: "Fuente del torrent",
    torrent_tag: "Etiqueta del torrent",
    video_resolution: "Resoluci\xF3n de video",
    video_source: "Fuente de video"
  },
  file_types: {
    archive: "Archivo",
    audio: "Audio",
    data: "Datos",
    document: "Documento",
    image: "Imagen",
    software: "Software",
    subtitles: "Subt\xEDtulos",
    unknown: "Desconocido",
    video: "Video"
  },
  general: {
    all: "Todos",
    dismiss: "Descartar",
    error: "Error",
    none: "Ninguno",
    page_not_found: "P\xE1gina no encontrada",
    refresh: "Actualizar",
    status: "Estado"
  },
  health: {
    bitmagnet_is_status: "bitmagnet est\xE1 {{status}}",
    check_failed_with_error: "La verificaci\xF3n fall\xF3 con error",
    component: "Componente",
    components: {
      dht: "DHT",
      postgres: "Postgres",
      tmdb: "TMDB"
    },
    dismiss: "Descartar",
    error: "Error",
    status: "Estado",
    statuses: {
      degraded: "Degradado",
      down: "Ca\xEDdo",
      error: "Error",
      inactive: "Inactivo",
      started: "Iniciado",
      unknown: "Pendiente",
      up: "Activo"
    },
    summary: "Resumen de salud",
    worker: "Trabajador",
    workers: {
      dht_crawler: "Rastreador DHT",
      http_server: "Servidor HTTP",
      queue_server: "Servidor de colas"
    }
  },
  languages: {
    af: "Afrik\xE1ans",
    ar: "\xC1rabe",
    az: "Azerbaiyano",
    be: "Bielorruso",
    bg: "B\xFAlgaro",
    bs: "Bosnio",
    ca: "Catal\xE1n",
    ce: "Checheno",
    co: "Corso",
    cs: "Checo",
    cy: "Gal\xE9s",
    da: "Dan\xE9s",
    de: "Alem\xE1n",
    el: "Griego",
    en: "Ingl\xE9s",
    es: "Espa\xF1ol",
    et: "Estonio",
    eu: "Vasco",
    fa: "Persa",
    fi: "Finland\xE9s",
    fr: "Franc\xE9s",
    he: "Hebreo",
    hi: "Hindi",
    hr: "Croata",
    hu: "H\xFAngaro",
    hy: "Armenio",
    id: "Indonesio",
    is: "Island\xE9s",
    it: "Italiano",
    ja: "Japon\xE9s",
    ka: "Georgiano",
    ko: "Coreano",
    ku: "Kurdo",
    lt: "Lituano",
    lv: "Let\xF3n",
    mi: "Maor\xED",
    mk: "Macedonio",
    ml: "Malayalam",
    mn: "Mongol",
    ms: "Malayo",
    mt: "Malt\xE9s",
    nl: "Neerland\xE9s",
    no: "Noruego",
    pl: "Polaco",
    pt: "Portugu\xE9s",
    ro: "Rumano",
    ru: "Ruso",
    sa: "S\xE1nscrito",
    sk: "Eslovaco",
    sl: "Esloveno",
    sm: "Samoano",
    so: "Somal\xED",
    sr: "Serbio",
    sv: "Sueco",
    ta: "Tamil",
    th: "Tailand\xE9s",
    tr: "Turco",
    uk: "Ucraniano",
    vi: "Vietnamita",
    yi: "Yidis",
    zh: "Chino",
    zu: "Zul\xFA"
  },
  layout: {
    bitmagnet_on_service: "bitmagnet en {{service}}",
    change_theme: "Cambiar tema",
    external_links: "Enlaces externos",
    sponsor: "Patrocinador",
    support_bitmagnet: "Apoyar bitmagnet",
    translate: "Traducir"
  },
  paginator: {
    first_page: "Primera p\xE1gina",
    items_per_page: "Elementos por p\xE1gina",
    last_page: "\xDAltima p\xE1gina",
    next_page: "P\xE1gina siguiente",
    page_x: "P\xE1gina {{x}}",
    previous_page: "P\xE1gina anterior",
    x_to_y: "{{x}} a {{y}}",
    x_to_y_of_z: "{{x}} a {{y}} de {{z}}"
  },
  routes: {
    admin: "Admin",
    dashboard: "Tablero",
    home: "Inicio",
    jobs: "Trabajos",
    queues: "Colas",
    torrents: "Torrents",
    visualize: "Visualizar"
  },
  torrents: {
    classification: "Clasificaci\xF3n",
    clear_search: "Borrar b\xFAsqueda",
    copy: "Copiar",
    copy_to_clipboard: "Copiar al portapapeles",
    delete: "Eliminar",
    delete_action_cannot_be_undone: "Esta acci\xF3n no se puede deshacer",
    delete_are_you_sure: "\xBFEst\xE1 seguro de que desea eliminar este torrent?",
    deselect_all: "Deseleccionar todo",
    edit_tags: "Editar etiquetas",
    episodes: "Episodios",
    external_links: "Enlaces externos",
    file_index: "\xCDndice de archivos",
    file_path: "Ruta del archivo",
    file_size: "Tama\xF1o del archivo",
    file_type: "Tipo de archivo",
    files: "Archivos",
    files_count_n: "{{count}} archivos",
    files_no_info: "No hay informaci\xF3n de archivos disponible",
    files_single: "Archivo \xFAnico",
    genres: "G\xE9neros",
    info_hash: "Hash de informaci\xF3n",
    info_hashes: "Hashes de informaci\xF3n",
    languages: "Idiomas",
    leechers: "Leechers",
    magnet: "Magnet",
    magnet_links: "Enlaces Magnet",
    new_tag: "Nueva etiqueta",
    order_by: "Ordenar por",
    order_direction_toggle: "Alternar direcci\xF3n",
    ordering: {
      files_count: "Recuento de archivos",
      info_hash: "Hash de informaci\xF3n",
      leechers: "Leechers",
      name: "Nombre",
      published_at: "Publicado en",
      relevance: "Relevancia",
      seeders: "Seeders",
      size: "Tama\xF1o",
      updated_at: "Actualizado en"
    },
    original_release_date: "Fecha de lanzamiento original",
    permalink: "Enlace permanente",
    poster: "P\xF3ster",
    published: "Publicado",
    rating: "Calificaci\xF3n",
    refresh: "Actualizar resultados",
    reprocess: {
      force_rematch: "Forzar nueva coincidencia de contenido ya coincidente",
      match_content_by_external_api_search: "Coincidir contenido por b\xFAsqueda de API externa",
      match_content_by_local_search: "Coincidir contenido por b\xFAsqueda local",
      reprocess: "Volver a procesar"
    },
    s_l: "S / L",
    search: "Buscar",
    seeders: "Seeders",
    select_all: "Seleccionar todo",
    showing_x_of_y_files: "Mostrando {{x}} de {{y}} archivos",
    size: "Tama\xF1o",
    source: "Fuente del torrent",
    summary: "Resumen",
    tags: {
      delete: "Eliminar etiquetas",
      delete_tip: "Eliminar etiquetas de los torrents seleccionados",
      placeholder: "Etiqueta...",
      put: "Poner etiquetas",
      put_tip: "A\xF1adir etiquetas a los torrents seleccionados",
      set: "Establecer etiquetas",
      set_tip: "Reemplazar etiquetas de los torrents seleccionados"
    },
    title: "T\xEDtulo",
    toggle_drawer: "Alternar caj\xF3n",
    votes_count_n: "{{count}} votos"
  },
  version: {
    bitmagnet_version: "versi\xF3n de bitmagnet {{version}}",
    unknown: "desconocido"
  }
};

// src/app/i18n/translations/fr.json
var fr_default = {
  content_types: {
    plural: {
      all: "Tous",
      audiobook: "Livres audio",
      comic: "Bandes dessin\xE9es",
      ebook: "E-livres",
      game: "Jeux",
      movie: "Films",
      music: "Musique",
      null: "Inconnu",
      software: "Logiciels",
      tv_show: "S\xE9ries TV",
      xxx: "XXX"
    },
    singular: {
      audiobook: "Livre audio",
      comic: "Bande dessin\xE9e",
      ebook: "E-livre",
      game: "Jeu",
      movie: "Film",
      music: "Musique",
      null: "Inconnu",
      software: "Logiciel",
      tv_show: "S\xE9rie TV",
      xxx: "XXX"
    }
  },
  dashboard: {
    event: {
      created: "Cr\xE9\xE9",
      failed: "\xC9chou\xE9",
      processed: "Trait\xE9",
      updated: "Mis \xE0 jour"
    },
    interval: {
      all: "Tous",
      days: "Jours",
      days_1: "1 jour",
      hours: "Heures",
      hours_1: "1 heure",
      hours_12: "12 heures",
      hours_6: "6 heures",
      minutes: "Minutes",
      minutes_1: "1 minute",
      minutes_15: "15 minutes",
      minutes_30: "30 minutes",
      minutes_5: "5 minutes",
      off: "D\xE9sactiv\xE9",
      seconds_10: "10 secondes",
      seconds_30: "30 secondes",
      weeks_1: "1 semaine"
    },
    metrics: {
      event: "\xC9v\xE9nement",
      resolution: "R\xE9solution",
      throughput: "D\xE9bit",
      timeframe: "P\xE9riode",
      toggle_legend: "Basculer la l\xE9gende"
    },
    queues: {
      created: "Cr\xE9\xE9",
      created_at: "Cr\xE9\xE9 \xE0",
      enqueue_jobs: "Mettre les taches en file d'attente",
      enqueue_torrent_processing_batch: "Mettre le traitement des torrents en file d'attente",
      failed: "\xC9chou\xE9",
      force_rematch: "Forcer le rematch du contenu d\xE9j\xE0 appari\xE9",
      jobs_enqueued: "T\xE2ches en attente",
      latency: "Latence",
      match_content_by_external_api_search: "Apparier le contenu par recherche API externe",
      match_content_by_local_search: "Apparier le contenu par recherche locale",
      payload: "Charge utile",
      pending: "En attente",
      priority: "Priorit\xE9",
      process_orphaned_torrents_only: "Traiter uniquement les torrents orphelins",
      processed: "Trait\xE9",
      purge_jobs: "Purger les t\xE2ches",
      purge_queue_jobs: "Purger les t\xE2ches de la file d'attente",
      queue: "File d'attente",
      queue_purged: "File d'attente purg\xE9e",
      queues: "Files d'attente",
      ran_at: "Ex\xE9cut\xE9 \xE0",
      retry: "R\xE9essayer",
      total_counts_by_status: "Totaux par statut"
    }
  },
  facets: {
    content_type: "Type de contenu",
    file_type: "Type de fichier",
    genre: "Genre",
    language: "Langue",
    queue: "File d'attente",
    status: "Statut",
    torrent_source: "Source du torrent",
    torrent_tag: "Tag du torrent",
    video_resolution: "R\xE9solution vid\xE9o",
    video_source: "Source vid\xE9o"
  },
  file_types: {
    archive: "Archive",
    audio: "Audio",
    data: "Donn\xE9es",
    document: "Document",
    image: "Image",
    software: "Logiciel",
    subtitles: "Sous-titres",
    unknown: "Inconnu",
    video: "Vid\xE9o"
  },
  general: {
    all: "Tous",
    dismiss: "Fermer",
    error: "Erreur",
    none: "Aucun",
    page_not_found: "Page non trouv\xE9e",
    refresh: "Rafra\xEEchir",
    status: "Statut"
  },
  health: {
    bitmagnet_is_status: "bitmagnet est {{status}}",
    check_failed_with_error: "V\xE9rification \xE9chou\xE9e avec erreur",
    component: "Composant",
    components: {
      dht: "DHT",
      postgres: "Postgres",
      tmdb: "TMDB"
    },
    dismiss: "Fermer",
    error: "Erreur",
    status: "Statut",
    statuses: {
      degraded: "D\xE9grad\xE9",
      down: "Hors ligne",
      error: "Erreur",
      inactive: "Inactif",
      started: "D\xE9marr\xE9",
      unknown: "En attente",
      up: "En ligne"
    },
    summary: "R\xE9sum\xE9 de la sant\xE9",
    worker: "Travailleur",
    workers: {
      dht_crawler: "Explorateur DHT",
      http_server: "Serveur HTTP",
      queue_server: "Serveur de file d'attente"
    }
  },
  languages: {
    af: "Afrikaans",
    ar: "Arabe",
    az: "Azerba\xEFdjanais",
    be: "Bi\xE9lorusse",
    bg: "Bulgare",
    bs: "Bosniaque",
    ca: "Catalan",
    ce: "Tch\xE9tch\xE8ne",
    co: "Corse",
    cs: "Tch\xE8que",
    cy: "Gallois",
    da: "Danois",
    de: "Allemand",
    el: "Grec",
    en: "Anglais",
    es: "Espagnol",
    et: "Estonien",
    eu: "Basque",
    fa: "Persan",
    fi: "Finnois",
    fr: "Fran\xE7ais",
    he: "H\xE9breu",
    hi: "Hindi",
    hr: "Croate",
    hu: "Hongrois",
    hy: "Arm\xE9nien",
    id: "Indon\xE9sien",
    is: "Islandais",
    it: "Italien",
    ja: "Japonais",
    ka: "G\xE9orgien",
    ko: "Cor\xE9en",
    ku: "Kurde",
    lt: "Lituanien",
    lv: "Letton",
    mi: "Maori",
    mk: "Mac\xE9donien",
    ml: "Malayalam",
    mn: "Mongol",
    ms: "Malais",
    mt: "Maltais",
    nl: "N\xE9erlandais",
    no: "Norv\xE9gien",
    pl: "Polonais",
    pt: "Portugais",
    ro: "Roumain",
    ru: "Russe",
    sa: "Sanskrit",
    sk: "Slovaque",
    sl: "Slov\xE8ne",
    sm: "Samoan",
    so: "Somali",
    sr: "Serbe",
    sv: "Su\xE9dois",
    ta: "Tamoul",
    th: "Tha\xEF",
    tr: "Turc",
    uk: "Ukrainien",
    vi: "Vietnamien",
    yi: "Yiddish",
    zh: "Chinois",
    zu: "Zoulou"
  },
  layout: {
    bitmagnet_on_service: "bitmagnet sur {{service}}",
    change_theme: "Changer de th\xE8me",
    external_links: "Liens externes",
    sponsor: "Sponsor",
    support_bitmagnet: "Soutenir bitmagnet",
    translate: "Traduire"
  },
  paginator: {
    first_page: "Premi\xE8re page",
    items_per_page: "Articles par page",
    last_page: "Derni\xE8re page",
    next_page: "Page suivante",
    page_x: "Page {{x}}",
    previous_page: "Page pr\xE9c\xE9dente",
    x_to_y: "{{x}} \xE0 {{y}}",
    x_to_y_of_z: "{{x}} \xE0 {{y}} sur {{z}}"
  },
  routes: {
    admin: "Admin",
    dashboard: "Tableau de bord",
    home: "Accueil",
    jobs: "T\xE2ches",
    queues: "Files d'attente",
    torrents: "Torrents",
    visualize: "Visualiser"
  },
  torrents: {
    classification: "Classification",
    clear_search: "Effacer la recherche",
    copy: "Copier",
    copy_to_clipboard: "Copier dans le presse-papiers",
    delete: "Supprimer",
    delete_action_cannot_be_undone: "Cette action ne peut pas \xEAtre annul\xE9e",
    delete_are_you_sure: "\xCAtes-vous s\xFBr de vouloir supprimer ce torrent?",
    deselect_all: "Tout d\xE9s\xE9lectionner",
    edit_tags: "Modifier les tags",
    episodes: "\xC9pisodes",
    external_links: "Liens externes",
    file_index: "Index de fichier",
    file_path: "Chemin de fichier",
    file_size: "Taille de fichier",
    file_type: "Type de fichier",
    files: "Fichiers",
    files_count_n: "{{count}} fichiers",
    files_no_info: "Aucune information sur les fichiers disponible",
    files_single: "Fichier unique",
    genres: "Genres",
    info_hash: "Hash d'info",
    info_hashes: "Hashes d'info",
    languages: "Langues",
    leechers: "Leechers",
    magnet: "Magnet",
    magnet_links: "Liens Magnet",
    new_tag: "Nouveau tag",
    order_by: "Trier par",
    order_direction_toggle: "Inverser le sens",
    ordering: {
      files_count: "Nombre de fichiers",
      info_hash: "Hash d'info",
      leechers: "Leechers",
      name: "Nom",
      published_at: "Publi\xE9 \xE0",
      relevance: "Pertinence",
      seeders: "Seeders",
      size: "Taille",
      updated_at: "Mis \xE0 jour \xE0"
    },
    original_release_date: "Date de sortie originale",
    permalink: "Permalien",
    poster: "Affiche",
    published: "Publi\xE9",
    rating: "\xC9valuation",
    refresh: "Rafra\xEEchir les r\xE9sultats",
    reprocess: {
      force_rematch: "Forcer le rematch du contenu d\xE9j\xE0 appari\xE9",
      match_content_by_external_api_search: "Apparier le contenu par recherche API externe",
      match_content_by_local_search: "Apparier le contenu par recherche locale",
      reprocess: "Retraitement"
    },
    s_l: "S / L",
    search: "Rechercher",
    seeders: "Seeders",
    select_all: "Tout s\xE9lectionner",
    showing_x_of_y_files: "Affichage de {{x}} sur {{y}} fichiers",
    size: "Taille",
    source: "Source du torrent",
    summary: "R\xE9sum\xE9",
    tags: {
      delete: "Supprimer les tags",
      delete_tip: "Supprimer les tags des torrents s\xE9lectionn\xE9s",
      placeholder: "Tag...",
      put: "Mettre des tags",
      put_tip: "Ajouter des tags aux torrents s\xE9lectionn\xE9s",
      set: "D\xE9finir des tags",
      set_tip: "Remplacer les tags des torrents s\xE9lectionn\xE9s"
    },
    title: "Titre",
    toggle_drawer: "Basculer le tiroir",
    votes_count_n: "{{count}} votes"
  },
  version: {
    bitmagnet_version: "version bitmagnet {{version}}",
    unknown: "inconnu"
  }
};

// src/app/i18n/translations/hi.json
var hi_default = {
  content_types: {
    plural: {
      all: "\u0938\u092D\u0940",
      audiobook: "\u0911\u0921\u093F\u092F\u094B\u092C\u0941\u0915\u094D\u0938",
      comic: "\u0915\u0949\u092E\u093F\u0915\u094D\u0938",
      ebook: "\u0908-\u092C\u0941\u0915\u094D\u0938",
      game: "\u0917\u0947\u092E\u094D\u0938",
      movie: "\u092B\u093C\u093F\u0932\u094D\u092E\u0947\u0902",
      music: "\u0938\u0902\u0917\u0940\u0924",
      null: "\u0905\u091C\u094D\u091E\u093E\u0924",
      software: "\u0938\u0949\u092B\u093C\u094D\u091F\u0935\u0947\u092F\u0930",
      tv_show: "\u091F\u0940\u0935\u0940 \u0936\u094B",
      xxx: "XXX"
    },
    singular: {
      audiobook: "\u0911\u0921\u093F\u092F\u094B\u092C\u0941\u0915",
      comic: "\u0915\u0949\u092E\u093F\u0915",
      ebook: "\u0908-\u092C\u0941\u0915",
      game: "\u0917\u0947\u092E",
      movie: "\u092B\u093C\u093F\u0932\u094D\u092E",
      music: "\u0938\u0902\u0917\u0940\u0924",
      null: "\u0905\u091C\u094D\u091E\u093E\u0924",
      software: "\u0938\u0949\u092B\u093C\u094D\u091F\u0935\u0947\u092F\u0930",
      tv_show: "\u091F\u0940\u0935\u0940 \u0936\u094B",
      xxx: "XXX"
    }
  },
  dashboard: {
    event: {
      created: "\u092C\u0928\u093E\u092F\u093E \u0917\u092F\u093E",
      failed: "\u0905\u0938\u092B\u0932",
      processed: "\u092A\u094D\u0930\u0938\u0902\u0938\u094D\u0915\u0943\u0924",
      updated: "\u0905\u092A\u0921\u0947\u091F \u0915\u093F\u092F\u093E \u0917\u092F\u093E"
    },
    interval: {
      all: "\u0938\u092D\u0940",
      days: "\u0926\u093F\u0928",
      days_1: "1 \u0926\u093F\u0928",
      hours: "\u0918\u0902\u091F\u0947",
      hours_1: "1 \u0918\u0902\u091F\u093E",
      hours_12: "12 \u0918\u0902\u091F\u0947",
      hours_6: "6 \u0918\u0902\u091F\u0947",
      minutes: "\u092E\u093F\u0928\u091F",
      minutes_1: "1 \u092E\u093F\u0928\u091F",
      minutes_15: "15 \u092E\u093F\u0928\u091F",
      minutes_30: "30 \u092E\u093F\u0928\u091F",
      minutes_5: "5 \u092E\u093F\u0928\u091F",
      off: "\u092C\u0902\u0926",
      seconds_10: "10 \u0938\u0947\u0915\u0902\u0921",
      seconds_30: "30 \u0938\u0947\u0915\u0902\u0921",
      weeks_1: "1 \u0938\u092A\u094D\u0924\u093E\u0939"
    },
    metrics: {
      event: "\u0918\u091F\u0928\u093E",
      resolution: "\u0930\u093F\u091C\u093C\u0949\u0932\u094D\u092F\u0942\u0936\u0928",
      throughput: "\u0925\u094D\u0930\u0942\u092A\u0941\u091F",
      timeframe: "\u0938\u092E\u092F \u0938\u0940\u092E\u093E",
      toggle_legend: "\u0932\u0940\u091C\u0947\u0902\u0921 \u091F\u0949\u0917\u0932 \u0915\u0930\u0947\u0902"
    },
    queues: {
      created: "\u092C\u0928\u093E\u092F\u093E \u0917\u092F\u093E",
      created_at: "\u092C\u0928\u093E\u092F\u093E \u0917\u092F\u093E \u0938\u092E\u092F",
      enqueue_jobs: "\u091C\u0949\u092C\u094D\u0938 \u0915\u094B \u0915\u0924\u093E\u0930 \u092E\u0947\u0902 \u0932\u0917\u093E\u090F\u0902",
      enqueue_torrent_processing_batch: "\u091F\u094B\u0930\u0947\u0902\u091F \u092A\u094D\u0930\u094B\u0938\u0947\u0938\u093F\u0902\u0917 \u092C\u0948\u091A \u0915\u0924\u093E\u0930 \u092E\u0947\u0902 \u0932\u0917\u093E\u090F\u0902",
      failed: "\u0905\u0938\u092B\u0932",
      force_rematch: "\u092A\u0939\u0932\u0947 \u0938\u0947 \u092E\u0947\u0932 \u0916\u093E\u0908 \u0938\u093E\u092E\u0917\u094D\u0930\u0940 \u0915\u094B \u092B\u093F\u0930 \u0938\u0947 \u092E\u093F\u0932\u093E\u090F\u0902",
      jobs_enqueued: "\u091C\u0949\u092C\u094D\u0938 \u0915\u0924\u093E\u0930 \u092E\u0947\u0902 \u0932\u0917\u093E\u0908 \u0917\u0908\u0902",
      latency: "\u0932\u0947\u091F\u0947\u0902\u0938\u0940",
      match_content_by_external_api_search: "\u092C\u093E\u0939\u0930\u0940 API \u0916\u094B\u091C \u0938\u0947 \u0938\u093E\u092E\u0917\u094D\u0930\u0940 \u0915\u093E \u092E\u093F\u0932\u093E\u0928 \u0915\u0930\u0947\u0902",
      match_content_by_local_search: "\u0938\u094D\u0925\u093E\u0928\u0940\u092F \u0916\u094B\u091C \u0938\u0947 \u0938\u093E\u092E\u0917\u094D\u0930\u0940 \u0915\u093E \u092E\u093F\u0932\u093E\u0928 \u0915\u0930\u0947\u0902",
      payload: "\u092A\u0947\u0932\u094B\u0921",
      pending: "\u092C\u0915\u093E\u092F\u093E",
      priority: "\u092A\u094D\u0930\u093E\u0925\u092E\u093F\u0915\u0924\u093E",
      process_orphaned_torrents_only: "\u0915\u0947\u0935\u0932 \u0905\u0928\u093E\u0925 \u091F\u094B\u0930\u0947\u0902\u091F\u094D\u0938 \u0915\u094B \u092A\u094D\u0930\u094B\u0938\u0947\u0938 \u0915\u0930\u0947\u0902",
      processed: "\u092A\u094D\u0930\u0938\u0902\u0938\u094D\u0915\u0943\u0924",
      purge_jobs: "\u091C\u0949\u092C\u094D\u0938 \u0915\u094B \u0938\u093E\u092B\u093C \u0915\u0930\u0947\u0902",
      purge_queue_jobs: "\u0915\u0924\u093E\u0930 \u0915\u0947 \u091C\u0949\u092C\u094D\u0938 \u0915\u094B \u0938\u093E\u092B\u093C \u0915\u0930\u0947\u0902",
      queue: "\u0915\u0924\u093E\u0930",
      queue_purged: "\u0915\u0924\u093E\u0930 \u0915\u094B \u0938\u093E\u092B\u093C \u0915\u093F\u092F\u093E \u0917\u092F\u093E",
      queues: "\u0915\u0924\u093E\u0930\u0947\u0902",
      ran_at: "\u091A\u0932\u093E\u092F\u093E \u0917\u092F\u093E \u0938\u092E\u092F",
      retry: "\u092A\u0941\u0928\u0903 \u092A\u094D\u0930\u092F\u093E\u0938 \u0915\u0930\u0947\u0902",
      total_counts_by_status: "\u0938\u094D\u0925\u093F\u0924\u093F \u0926\u094D\u0935\u093E\u0930\u093E \u0915\u0941\u0932 \u0917\u0923\u0928\u093E"
    }
  },
  facets: {
    content_type: "\u0938\u093E\u092E\u0917\u094D\u0930\u0940 \u092A\u094D\u0930\u0915\u093E\u0930",
    file_type: "\u092B\u093C\u093E\u0907\u0932 \u092A\u094D\u0930\u0915\u093E\u0930",
    genre: "\u0936\u0948\u0932\u0940",
    language: "\u092D\u093E\u0937\u093E",
    queue: "\u0915\u0924\u093E\u0930",
    status: "\u0938\u094D\u0925\u093F\u0924\u093F",
    torrent_source: "\u091F\u094B\u0930\u0947\u0902\u091F \u0938\u094D\u0930\u094B\u0924",
    torrent_tag: "\u091F\u094B\u0930\u0947\u0902\u091F \u091F\u0948\u0917",
    video_resolution: "\u0935\u0940\u0921\u093F\u092F\u094B \u0930\u093F\u091C\u093C\u0949\u0932\u094D\u092F\u0942\u0936\u0928",
    video_source: "\u0935\u0940\u0921\u093F\u092F\u094B \u0938\u094D\u0930\u094B\u0924"
  },
  file_types: {
    archive: "\u0906\u0930\u094D\u0915\u093E\u0907\u0935",
    audio: "\u0911\u0921\u093F\u092F\u094B",
    data: "\u0921\u0947\u091F\u093E",
    document: "\u0926\u0938\u094D\u0924\u093E\u0935\u0947\u091C\u093C",
    image: "\u091B\u0935\u093F",
    software: "\u0938\u0949\u092B\u093C\u094D\u091F\u0935\u0947\u092F\u0930",
    subtitles: "\u0909\u092A\u0936\u0940\u0930\u094D\u0937\u0915",
    unknown: "\u0905\u091C\u094D\u091E\u093E\u0924",
    video: "\u0935\u0940\u0921\u093F\u092F\u094B"
  },
  general: {
    all: "\u0938\u092D\u0940",
    dismiss: "\u0916\u093E\u0930\u093F\u091C \u0915\u0930\u0947\u0902",
    error: "\u0924\u094D\u0930\u0941\u091F\u093F",
    none: "\u0915\u094B\u0908 \u0928\u0939\u0940\u0902",
    page_not_found: "\u092A\u0943\u0937\u094D\u0920 \u0928\u0939\u0940\u0902 \u092E\u093F\u0932\u093E",
    refresh: "\u0924\u093E\u091C\u093C\u093E \u0915\u0930\u0947\u0902",
    status: "\u0938\u094D\u0925\u093F\u0924\u093F"
  },
  health: {
    bitmagnet_is_status: "bitmagnet {{status}} \u0939\u0948",
    check_failed_with_error: "\u0924\u094D\u0930\u0941\u091F\u093F \u0915\u0947 \u0938\u093E\u0925 \u091C\u093E\u0902\u091A \u0935\u093F\u092B\u0932",
    component: "\u0918\u091F\u0915",
    components: {
      dht: "DHT",
      postgres: "Postgres",
      tmdb: "TMDB"
    },
    dismiss: "\u0916\u093E\u0930\u093F\u091C \u0915\u0930\u0947\u0902",
    error: "\u0924\u094D\u0930\u0941\u091F\u093F",
    status: "\u0938\u094D\u0925\u093F\u0924\u093F",
    statuses: {
      degraded: "\u0915\u094D\u0937\u0940\u0923",
      down: "\u092C\u0902\u0926",
      error: "\u0924\u094D\u0930\u0941\u091F\u093F",
      inactive: "\u0928\u093F\u0937\u094D\u0915\u094D\u0930\u093F\u092F",
      started: "\u0936\u0941\u0930\u0942 \u0915\u093F\u092F\u093E \u0917\u092F\u093E",
      unknown: "\u0932\u0902\u092C\u093F\u0924",
      up: "\u091A\u093E\u0932\u0942"
    },
    summary: "\u0938\u094D\u0935\u093E\u0938\u094D\u0925\u094D\u092F \u0938\u093E\u0930\u093E\u0902\u0936",
    worker: "\u0915\u093E\u0930\u094D\u092F\u0915\u0930\u094D\u0924\u093E",
    workers: {
      dht_crawler: "DHT \u0915\u094D\u0930\u0949\u0932\u0930",
      http_server: "HTTP \u0938\u0930\u094D\u0935\u0930",
      queue_server: "\u0915\u0924\u093E\u0930 \u0938\u0930\u094D\u0935\u0930"
    }
  },
  languages: {
    af: "\u0905\u092B\u094D\u0930\u0940\u0915\u0940",
    ar: "\u0905\u0930\u092C\u0940",
    az: "\u0905\u091C\u093C\u0947\u0930\u0940",
    be: "\u092C\u0947\u0932\u093E\u0930\u0942\u0938\u0940",
    bg: "\u092C\u0941\u0932\u094D\u0917\u093E\u0930\u093F\u092F\u093E\u0908",
    bs: "\u092C\u094B\u0938\u094D\u0928\u093F\u092F\u093E\u0908",
    ca: "\u0915\u0948\u091F\u0932\u0928",
    ce: "\u091A\u0947\u091A\u0947\u0928",
    co: "\u0915\u094B\u0930\u094D\u0938\u093F\u0915\u0928",
    cs: "\u091A\u0947\u0915",
    cy: "\u0935\u0947\u0932\u094D\u0936",
    da: "\u0921\u0947\u0928\u093F\u0936",
    de: "\u091C\u0930\u094D\u092E\u0928",
    el: "\u0917\u094D\u0930\u0940\u0915",
    en: "\u0905\u0902\u0917\u094D\u0930\u0947\u091C\u0940",
    es: "\u0938\u094D\u092A\u0948\u0928\u093F\u0936",
    et: "\u090F\u0938\u094D\u0924\u094B\u0928\u093F\u092F\u093E\u0908",
    eu: "\u092C\u093E\u0938\u094D\u0915",
    fa: "\u092B\u093E\u0930\u0938\u0940",
    fi: "\u092B\u093C\u093F\u0928\u093F\u0936",
    fr: "\u092B\u094D\u0930\u0947\u0902\u091A",
    he: "\u0939\u093F\u092C\u094D\u0930\u0942",
    hi: "\u0939\u093F\u0902\u0926\u0940",
    hr: "\u0915\u094D\u0930\u094B\u090F\u0936\u093F\u092F\u093E\u0908",
    hu: "\u0939\u0902\u0917\u0947\u0930\u093F\u092F\u0928",
    hy: "\u0906\u0930\u094D\u092E\u0947\u0928\u093F\u092F\u093E\u0908",
    id: "\u0907\u0902\u0921\u094B\u0928\u0947\u0936\u093F\u092F\u093E\u0908",
    is: "\u0906\u0907\u0938\u0932\u0948\u0902\u0921\u093F\u0915",
    it: "\u0907\u0924\u093E\u0932\u0935\u0940",
    ja: "\u091C\u093E\u092A\u093E\u0928\u0940",
    ka: "\u091C\u0949\u0930\u094D\u091C\u093F\u092F\u093E\u0908",
    ko: "\u0915\u094B\u0930\u093F\u092F\u093E\u0908",
    ku: "\u0915\u0941\u0930\u094D\u0926\u0940",
    lt: "\u0932\u093F\u0925\u0941\u0906\u0928\u093F\u092F\u093E\u0908",
    lv: "\u0932\u093E\u0924\u0935\u093F\u092F\u093E\u0908",
    mi: "\u092E\u093E\u0913\u0930\u0940",
    mk: "\u092E\u0948\u0938\u093F\u0921\u094B\u0928\u093F\u092F\u093E\u0908",
    ml: "\u092E\u0932\u092F\u093E\u0932\u092E",
    mn: "\u092E\u0902\u0917\u094B\u0932\u093F\u092F\u093E\u0908",
    ms: "\u092E\u0932\u092F",
    mt: "\u092E\u093E\u0932\u094D\u091F\u0940\u091C\u093C",
    nl: "\u0921\u091A",
    no: "\u0928\u0949\u0930\u094D\u0935\u0947\u091C\u093F\u092F\u0928",
    pl: "\u092A\u094B\u0932\u093F\u0936",
    pt: "\u092A\u0941\u0930\u094D\u0924\u0917\u093E\u0932\u0940",
    ro: "\u0930\u094B\u092E\u093E\u0928\u093F\u092F\u093E\u0908",
    ru: "\u0930\u0942\u0938\u0940",
    sa: "\u0938\u0902\u0938\u094D\u0915\u0943\u0924",
    sk: "\u0938\u094D\u0932\u094B\u0935\u093E\u0915",
    sl: "\u0938\u094D\u0932\u094B\u0935\u0947\u0928\u093F\u092F\u093E\u0908",
    sm: "\u0938\u093E\u092E\u094B\u0928",
    so: "\u0938\u094B\u092E\u093E\u0932\u0940",
    sr: "\u0938\u0930\u094D\u092C\u093F\u092F\u093E\u0908",
    sv: "\u0938\u094D\u0935\u0940\u0921\u093F\u0936",
    ta: "\u0924\u092E\u093F\u0932",
    th: "\u0925\u093E\u0908",
    tr: "\u0924\u0941\u0930\u094D\u0915\u0940",
    uk: "\u092F\u0942\u0915\u094D\u0930\u0947\u0928\u0940",
    vi: "\u0935\u093F\u092F\u0924\u0928\u093E\u092E\u0940",
    yi: "\u092F\u093F\u0926\u094D\u0926\u093F\u0936",
    zh: "\u091A\u0940\u0928\u0940",
    zu: "\u091C\u093C\u0941\u0932\u0941"
  },
  layout: {
    bitmagnet_on_service: "{{service}} \u092A\u0930 bitmagnet",
    change_theme: "\u0925\u0940\u092E \u092C\u0926\u0932\u0947\u0902",
    external_links: "\u092C\u093E\u0939\u0930\u0940 \u0932\u093F\u0902\u0915",
    sponsor: "\u092A\u094D\u0930\u093E\u092F\u094B\u091C\u0915",
    support_bitmagnet: "bitmagnet \u0915\u094B \u0938\u092E\u0930\u094D\u0925\u0928 \u0926\u0947\u0902",
    translate: "\u0905\u0928\u0941\u0935\u093E\u0926 \u0915\u0930\u0947\u0902"
  },
  paginator: {
    first_page: "\u092A\u0939\u0932\u093E \u092A\u0943\u0937\u094D\u0920",
    items_per_page: "\u092A\u094D\u0930\u0924\u093F \u092A\u0943\u0937\u094D\u0920 \u0906\u0907\u091F\u092E",
    last_page: "\u0905\u0902\u0924\u093F\u092E \u092A\u0943\u0937\u094D\u0920",
    next_page: "\u0905\u0917\u0932\u093E \u092A\u0943\u0937\u094D\u0920",
    page_x: "\u092A\u0943\u0937\u094D\u0920 {{x}}",
    previous_page: "\u092A\u093F\u091B\u0932\u093E \u092A\u0943\u0937\u094D\u0920",
    x_to_y: "{{x}} \u0938\u0947 {{y}} \u0924\u0915",
    x_to_y_of_z: "{{x}} \u0938\u0947 {{y}} \u0924\u0915 {{z}} \u092E\u0947\u0902"
  },
  routes: {
    admin: "\u092A\u094D\u0930\u0936\u093E\u0938\u0915",
    dashboard: "\u0921\u0948\u0936\u092C\u094B\u0930\u094D\u0921",
    home: "\u0939\u094B\u092E",
    jobs: "\u0928\u094C\u0915\u0930\u093F\u092F\u093E\u0902",
    queues: "\u0915\u0924\u093E\u0930\u0947\u0902",
    torrents: "\u091F\u094B\u0930\u0947\u0902\u091F\u094D\u0938",
    visualize: "\u0926\u0943\u0936\u094D\u092F \u0915\u0930\u0947\u0902"
  },
  torrents: {
    classification: "\u0935\u0930\u094D\u0917\u0940\u0915\u0930\u0923",
    clear_search: "\u0916\u094B\u091C \u0938\u093E\u092B\u093C \u0915\u0930\u0947\u0902",
    copy: "\u0915\u0949\u092A\u0940 \u0915\u0930\u0947\u0902",
    copy_to_clipboard: "\u0915\u094D\u0932\u093F\u092A\u092C\u094B\u0930\u094D\u0921 \u092E\u0947\u0902 \u0915\u0949\u092A\u0940 \u0915\u0930\u0947\u0902",
    delete: "\u0939\u091F\u093E\u090F\u0902",
    delete_action_cannot_be_undone: "\u092F\u0939 \u0915\u094D\u0930\u093F\u092F\u093E \u092A\u0942\u0930\u094D\u0935\u0935\u0924 \u0928\u0939\u0940\u0902 \u0915\u0940 \u091C\u093E \u0938\u0915\u0924\u0940",
    delete_are_you_sure: "\u0915\u094D\u092F\u093E \u0906\u092A \u0935\u093E\u0915\u0908 \u0907\u0938 \u091F\u094B\u0930\u0947\u0902\u091F \u0915\u094B \u0939\u091F\u093E\u0928\u093E \u091A\u093E\u0939\u0924\u0947 \u0939\u0948\u0902?",
    deselect_all: "\u0938\u092D\u0940 \u0915\u094B \u0905\u091A\u092F\u0928\u093F\u0924 \u0915\u0930\u0947\u0902",
    edit_tags: "\u091F\u0948\u0917 \u0938\u0902\u092A\u093E\u0926\u093F\u0924 \u0915\u0930\u0947\u0902",
    episodes: "\u090F\u092A\u093F\u0938\u094B\u0921\u094D\u0938",
    external_links: "\u092C\u093E\u0939\u0930\u0940 \u0932\u093F\u0902\u0915",
    file_index: "\u092B\u093C\u093E\u0907\u0932 \u0938\u0942\u091A\u0915\u093E\u0902\u0915",
    file_path: "\u092B\u093C\u093E\u0907\u0932 \u092A\u0925",
    file_size: "\u092B\u093C\u093E\u0907\u0932 \u0906\u0915\u093E\u0930",
    file_type: "\u092B\u093C\u093E\u0907\u0932 \u092A\u094D\u0930\u0915\u093E\u0930",
    files: "\u092B\u093C\u093E\u0907\u0932\u0947\u0902",
    files_count_n: "{{count}} \u092B\u093C\u093E\u0907\u0932\u0947\u0902",
    files_no_info: "\u0915\u094B\u0908 \u092B\u093C\u093E\u0907\u0932 \u091C\u093E\u0928\u0915\u093E\u0930\u0940 \u0909\u092A\u0932\u092C\u094D\u0927 \u0928\u0939\u0940\u0902 \u0939\u0948",
    files_single: "\u090F\u0915\u0932 \u092B\u093C\u093E\u0907\u0932",
    genres: "\u0936\u0948\u0932\u093F\u092F\u093E\u0901",
    info_hash: "\u0938\u0942\u091A\u0928\u093E \u0939\u0948\u0936",
    info_hashes: "\u0938\u0942\u091A\u0928\u093E \u0939\u0948\u0936\u0947\u091C",
    languages: "\u092D\u093E\u0937\u093E\u090F\u0901",
    leechers: "\u0932\u0940\u091A\u0930\u094D\u0938",
    magnet: "\u092E\u0948\u0917\u094D\u0928\u0947\u091F",
    magnet_links: "\u092E\u0948\u0917\u094D\u0928\u0947\u091F \u0932\u093F\u0902\u0915",
    new_tag: "\u0928\u092F\u093E \u091F\u0948\u0917",
    order_by: "\u0915\u094D\u0930\u092E\u092C\u0926\u094D\u0927 \u0915\u0930\u0947\u0902",
    order_direction_toggle: "\u0926\u093F\u0936\u093E \u092C\u0926\u0932\u0947\u0902",
    ordering: {
      files_count: "\u092B\u093C\u093E\u0907\u0932\u094B\u0902 \u0915\u0940 \u0938\u0902\u0916\u094D\u092F\u093E",
      info_hash: "\u0938\u0942\u091A\u0928\u093E \u0939\u0948\u0936",
      leechers: "\u0932\u0940\u091A\u0930\u094D\u0938",
      name: "\u0928\u093E\u092E",
      published_at: "\u092A\u094D\u0930\u0915\u093E\u0936\u093F\u0924 \u0938\u092E\u092F",
      relevance: "\u092A\u094D\u0930\u093E\u0938\u0902\u0917\u093F\u0915\u0924\u093E",
      seeders: "\u0938\u0940\u0921\u0930\u094D\u0938",
      size: "\u0906\u0915\u093E\u0930",
      updated_at: "\u0905\u092A\u0921\u0947\u091F \u0938\u092E\u092F"
    },
    original_release_date: "\u092E\u0942\u0932 \u0930\u093F\u0932\u0940\u091C\u093C \u0924\u093F\u0925\u093F",
    permalink: "\u0938\u094D\u0925\u093E\u092F\u0940 \u0932\u093F\u0902\u0915",
    poster: "\u092A\u094B\u0938\u094D\u091F\u0930",
    published: "\u092A\u094D\u0930\u0915\u093E\u0936\u093F\u0924",
    rating: "\u0930\u0947\u091F\u093F\u0902\u0917",
    refresh: "\u092A\u0930\u093F\u0923\u093E\u092E \u0924\u093E\u091C\u093C\u093E \u0915\u0930\u0947\u0902",
    reprocess: {
      force_rematch: "\u092A\u0939\u0932\u0947 \u0938\u0947 \u092E\u0947\u0932 \u0916\u093E\u0908 \u0938\u093E\u092E\u0917\u094D\u0930\u0940 \u0915\u094B \u092B\u093F\u0930 \u0938\u0947 \u092E\u093F\u0932\u093E\u090F\u0902",
      match_content_by_external_api_search: "\u092C\u093E\u0939\u0930\u0940 API \u0916\u094B\u091C \u0938\u0947 \u0938\u093E\u092E\u0917\u094D\u0930\u0940 \u0915\u093E \u092E\u093F\u0932\u093E\u0928 \u0915\u0930\u0947\u0902",
      match_content_by_local_search: "\u0938\u094D\u0925\u093E\u0928\u0940\u092F \u0916\u094B\u091C \u0938\u0947 \u0938\u093E\u092E\u0917\u094D\u0930\u0940 \u0915\u093E \u092E\u093F\u0932\u093E\u0928 \u0915\u0930\u0947\u0902",
      reprocess: "\u092A\u0941\u0928\u0903 \u092A\u094D\u0930\u0915\u094D\u0930\u093F\u092F\u093E \u0915\u0930\u0947\u0902"
    },
    s_l: "S / L",
    search: "\u0916\u094B\u091C",
    seeders: "\u0938\u0940\u0921\u0930\u094D\u0938",
    select_all: "\u0938\u092D\u0940 \u0915\u093E \u091A\u092F\u0928 \u0915\u0930\u0947\u0902",
    showing_x_of_y_files: "{{x}} \u092E\u0947\u0902 \u0938\u0947 {{y}} \u092B\u093C\u093E\u0907\u0932\u0947\u0902 \u0926\u093F\u0916\u093E \u0930\u0939\u093E \u0939\u0948",
    size: "\u0906\u0915\u093E\u0930",
    source: "\u091F\u094B\u0930\u0947\u0902\u091F \u0938\u094D\u0930\u094B\u0924",
    summary: "\u0938\u093E\u0930\u093E\u0902\u0936",
    tags: {
      delete: "\u091F\u0948\u0917 \u0939\u091F\u093E\u090F\u0902",
      delete_tip: "\u091A\u092F\u0928\u093F\u0924 \u091F\u094B\u0930\u0947\u0902\u091F \u0938\u0947 \u091F\u0948\u0917 \u0939\u091F\u093E\u090F\u0902",
      placeholder: "\u091F\u0948\u0917...",
      put: "\u091F\u0948\u0917 \u0921\u093E\u0932\u0947\u0902",
      put_tip: "\u091A\u092F\u0928\u093F\u0924 \u091F\u094B\u0930\u0947\u0902\u091F \u092E\u0947\u0902 \u091F\u0948\u0917 \u091C\u094B\u0921\u093C\u0947\u0902",
      set: "\u091F\u0948\u0917 \u0938\u0947\u091F \u0915\u0930\u0947\u0902",
      set_tip: "\u091A\u092F\u0928\u093F\u0924 \u091F\u094B\u0930\u0947\u0902\u091F \u0915\u0947 \u091F\u0948\u0917 \u092C\u0926\u0932\u0947\u0902"
    },
    title: "\u0936\u0940\u0930\u094D\u0937\u0915",
    toggle_drawer: "\u0921\u094D\u0930\u0949\u0905\u0930 \u092C\u0926\u0932\u0947\u0902",
    votes_count_n: "{{count}} \u0935\u094B\u091F"
  },
  version: {
    bitmagnet_version: "bitmagnet \u0938\u0902\u0938\u094D\u0915\u0930\u0923 {{version}}",
    unknown: "\u0905\u091C\u094D\u091E\u093E\u0924"
  }
};

// src/app/i18n/translations/ja.json
var ja_default = {
  content_types: {
    plural: {
      all: "\u3059\u3079\u3066",
      audiobook: "\u30AA\u30FC\u30C7\u30A3\u30AA\u30D6\u30C3\u30AF",
      comic: "\u30B3\u30DF\u30C3\u30AF",
      ebook: "\u96FB\u5B50\u66F8\u7C4D",
      game: "\u30B2\u30FC\u30E0",
      movie: "\u6620\u753B",
      music: "\u97F3\u697D",
      null: "\u4E0D\u660E",
      software: "\u30BD\u30D5\u30C8\u30A6\u30A7\u30A2",
      tv_show: "\u30C6\u30EC\u30D3\u756A\u7D44",
      xxx: "\u30DD\u30EB\u30CE"
    },
    singular: {
      audiobook: "\u30AA\u30FC\u30C7\u30A3\u30AA\u30D6\u30C3\u30AF",
      comic: "\u30B3\u30DF\u30C3\u30AF",
      ebook: "\u96FB\u5B50\u66F8\u7C4D",
      game: "\u30B2\u30FC\u30E0",
      movie: "\u6620\u753B",
      music: "\u97F3\u697D",
      null: "\u4E0D\u660E",
      software: "\u30BD\u30D5\u30C8\u30A6\u30A7\u30A2",
      tv_show: "\u30C6\u30EC\u30D3\u756A\u7D44",
      xxx: "\u30DD\u30EB\u30CE"
    }
  },
  dashboard: {
    event: {
      created: "\u4F5C\u6210\u6E08\u307F",
      failed: "\u5931\u6557",
      processed: "\u51E6\u7406\u6E08\u307F",
      updated: "\u66F4\u65B0\u6E08\u307F"
    },
    interval: {
      all: "\u3059\u3079\u3066",
      days: "\u65E5",
      days_1: "1\u65E5",
      hours: "\u6642\u9593",
      hours_1: "1\u6642\u9593",
      hours_12: "12\u6642\u9593",
      hours_6: "6\u6642\u9593",
      minutes: "\u5206",
      minutes_1: "1\u5206",
      minutes_15: "15\u5206",
      minutes_30: "30\u5206",
      minutes_5: "5\u5206",
      off: "\u30AA\u30D5",
      seconds_10: "10\u79D2",
      seconds_30: "30\u79D2",
      weeks_1: "1\u9031\u9593"
    },
    metrics: {
      event: "\u30A4\u30D9\u30F3\u30C8",
      resolution: "\u89E3\u50CF\u5EA6",
      throughput: "\u30B9\u30EB\u30FC\u30D7\u30C3\u30C8",
      timeframe: "\u671F\u9593",
      toggle_legend: "\u51E1\u4F8B\u3092\u5207\u308A\u66FF\u3048"
    },
    queues: {
      created: "\u4F5C\u6210\u6E08\u307F",
      created_at: "\u4F5C\u6210\u65E5\u6642",
      enqueue_jobs: "\u30B8\u30E7\u30D6\u3092\u30AD\u30E5\u30FC\u306B\u5165\u308C\u308B",
      enqueue_torrent_processing_batch: "\u30C8\u30EC\u30F3\u30C8\u51E6\u7406\u30D0\u30C3\u30C1\u3092\u30AD\u30E5\u30FC\u306B\u5165\u308C\u308B",
      failed: "\u5931\u6557",
      force_rematch: "\u65E2\u306B\u4E00\u81F4\u3057\u305F\u30B3\u30F3\u30C6\u30F3\u30C4\u3092\u5F37\u5236\u7684\u306B\u518D\u4E00\u81F4\u3055\u305B\u308B",
      jobs_enqueued: "\u30AD\u30E5\u30FC\u306B\u8FFD\u52A0\u3055\u308C\u305F\u30B8\u30E7\u30D6",
      latency: "\u30EC\u30A4\u30C6\u30F3\u30B7",
      match_content_by_external_api_search: "\u5916\u90E8API\u691C\u7D22\u3067\u30B3\u30F3\u30C6\u30F3\u30C4\u3092\u4E00\u81F4\u3055\u305B\u308B",
      match_content_by_local_search: "\u30ED\u30FC\u30AB\u30EB\u691C\u7D22\u3067\u30B3\u30F3\u30C6\u30F3\u30C4\u3092\u4E00\u81F4\u3055\u305B\u308B",
      payload: "\u30DA\u30A4\u30ED\u30FC\u30C9",
      pending: "\u4FDD\u7559\u4E2D",
      priority: "\u512A\u5148\u9806\u4F4D",
      process_orphaned_torrents_only: "\u5B64\u7ACB\u3057\u305F\u30C8\u30EC\u30F3\u30C8\u306E\u307F\u51E6\u7406\u3059\u308B",
      processed: "\u51E6\u7406\u6E08\u307F",
      purge_jobs: "\u30B8\u30E7\u30D6\u3092\u524A\u9664",
      purge_queue_jobs: "\u30AD\u30E5\u30FC\u5185\u306E\u30B8\u30E7\u30D6\u3092\u524A\u9664",
      queue: "\u30AD\u30E5\u30FC",
      queue_purged: "\u30AD\u30E5\u30FC\u304C\u30AF\u30EA\u30A2\u3055\u308C\u307E\u3057\u305F",
      queues: "\u30AD\u30E5\u30FC",
      ran_at: "\u5B9F\u884C\u6642\u523B",
      retry: "\u518D\u8A66\u884C",
      total_counts_by_status: "\u30B9\u30C6\u30FC\u30BF\u30B9\u3054\u3068\u306E\u7DCF\u6570"
    }
  },
  facets: {
    content_type: "\u30B3\u30F3\u30C6\u30F3\u30C4\u30BF\u30A4\u30D7",
    file_type: "\u30D5\u30A1\u30A4\u30EB\u30BF\u30A4\u30D7",
    genre: "\u30B8\u30E3\u30F3\u30EB",
    language: "\u8A00\u8A9E",
    queue: "\u30AD\u30E5\u30FC",
    status: "\u30B9\u30C6\u30FC\u30BF\u30B9",
    torrent_source: "\u30C8\u30EC\u30F3\u30C8\u30BD\u30FC\u30B9",
    torrent_tag: "\u30C8\u30EC\u30F3\u30C8\u30BF\u30B0",
    video_resolution: "\u30D3\u30C7\u30AA\u89E3\u50CF\u5EA6",
    video_source: "\u30D3\u30C7\u30AA\u30BD\u30FC\u30B9"
  },
  file_types: {
    archive: "\u30A2\u30FC\u30AB\u30A4\u30D6",
    audio: "\u30AA\u30FC\u30C7\u30A3\u30AA",
    data: "\u30C7\u30FC\u30BF",
    document: "\u30C9\u30AD\u30E5\u30E1\u30F3\u30C8",
    image: "\u753B\u50CF",
    software: "\u30BD\u30D5\u30C8\u30A6\u30A7\u30A2",
    subtitles: "\u5B57\u5E55",
    unknown: "\u4E0D\u660E",
    video: "\u30D3\u30C7\u30AA"
  },
  general: {
    all: "\u3059\u3079\u3066",
    dismiss: "\u9589\u3058\u308B",
    error: "\u30A8\u30E9\u30FC",
    none: "\u306A\u3057",
    page_not_found: "\u30DA\u30FC\u30B8\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093",
    refresh: "\u66F4\u65B0",
    status: "\u30B9\u30C6\u30FC\u30BF\u30B9"
  },
  health: {
    bitmagnet_is_status: "bitmagnet\u306F{{status}}\u3067\u3059",
    check_failed_with_error: "\u30A8\u30E9\u30FC\u3067\u30C1\u30A7\u30C3\u30AF\u306B\u5931\u6557\u3057\u307E\u3057\u305F",
    component: "\u30B3\u30F3\u30DD\u30FC\u30CD\u30F3\u30C8",
    components: {
      dht: "DHT",
      postgres: "Postgres",
      tmdb: "TMDB"
    },
    dismiss: "\u9589\u3058\u308B",
    error: "\u30A8\u30E9\u30FC",
    status: "\u30B9\u30C6\u30FC\u30BF\u30B9",
    statuses: {
      degraded: "\u4F4E\u4E0B",
      down: "\u505C\u6B62",
      error: "\u30A8\u30E9\u30FC",
      inactive: "\u975E\u30A2\u30AF\u30C6\u30A3\u30D6",
      started: "\u958B\u59CB",
      unknown: "\u4FDD\u7559\u4E2D",
      up: "\u7A3C\u50CD\u4E2D"
    },
    summary: "\u5065\u5EB7\u72B6\u614B\u306E\u6982\u8981",
    worker: "\u30EF\u30FC\u30AB\u30FC",
    workers: {
      dht_crawler: "DHT\u30AF\u30ED\u30FC\u30E9\u30FC",
      http_server: "HTTP\u30B5\u30FC\u30D0\u30FC",
      queue_server: "\u30AD\u30E5\u30FC\u30B5\u30FC\u30D0\u30FC"
    }
  },
  languages: {
    af: "\u30A2\u30D5\u30EA\u30AB\u30FC\u30F3\u30B9\u8A9E",
    ar: "\u30A2\u30E9\u30D3\u30A2\u8A9E",
    az: "\u30A2\u30BC\u30EB\u30D0\u30A4\u30B8\u30E3\u30F3\u8A9E",
    be: "\u30D9\u30E9\u30EB\u30FC\u30B7\u8A9E",
    bg: "\u30D6\u30EB\u30AC\u30EA\u30A2\u8A9E",
    bs: "\u30DC\u30B9\u30CB\u30A2\u8A9E",
    ca: "\u30AB\u30BF\u30EB\u30FC\u30CB\u30E3\u8A9E",
    ce: "\u30C1\u30A7\u30C1\u30A7\u30F3\u8A9E",
    co: "\u30B3\u30EB\u30B7\u30AB\u8A9E",
    cs: "\u30C1\u30A7\u30B3\u8A9E",
    cy: "\u30A6\u30A7\u30FC\u30EB\u30BA\u8A9E",
    da: "\u30C7\u30F3\u30DE\u30FC\u30AF\u8A9E",
    de: "\u30C9\u30A4\u30C4\u8A9E",
    el: "\u30AE\u30EA\u30B7\u30E3\u8A9E",
    en: "\u82F1\u8A9E",
    es: "\u30B9\u30DA\u30A4\u30F3\u8A9E",
    et: "\u30A8\u30B9\u30C8\u30CB\u30A2\u8A9E",
    eu: "\u30D0\u30B9\u30AF\u8A9E",
    fa: "\u30DA\u30EB\u30B7\u30E3\u8A9E",
    fi: "\u30D5\u30A3\u30F3\u30E9\u30F3\u30C9\u8A9E",
    fr: "\u30D5\u30E9\u30F3\u30B9\u8A9E",
    he: "\u30D8\u30D6\u30E9\u30A4\u8A9E",
    hi: "\u30D2\u30F3\u30C7\u30A3\u30FC\u8A9E",
    hr: "\u30AF\u30ED\u30A2\u30C1\u30A2\u8A9E",
    hu: "\u30CF\u30F3\u30AC\u30EA\u30FC\u8A9E",
    hy: "\u30A2\u30EB\u30E1\u30CB\u30A2\u8A9E",
    id: "\u30A4\u30F3\u30C9\u30CD\u30B7\u30A2\u8A9E",
    is: "\u30A2\u30A4\u30B9\u30E9\u30F3\u30C9\u8A9E",
    it: "\u30A4\u30BF\u30EA\u30A2\u8A9E",
    ja: "\u65E5\u672C\u8A9E",
    ka: "\u30B0\u30EB\u30B8\u30A2\u8A9E",
    ko: "\u97D3\u56FD\u8A9E",
    ku: "\u30AF\u30EB\u30C9\u8A9E",
    lt: "\u30EA\u30C8\u30A2\u30CB\u30A2\u8A9E",
    lv: "\u30E9\u30C8\u30D3\u30A2\u8A9E",
    mi: "\u30DE\u30AA\u30EA\u8A9E",
    mk: "\u30DE\u30B1\u30C9\u30CB\u30A2\u8A9E",
    ml: "\u30DE\u30E9\u30E4\u30FC\u30E9\u30E0\u8A9E",
    mn: "\u30E2\u30F3\u30B4\u30EB\u8A9E",
    ms: "\u30DE\u30EC\u30FC\u8A9E",
    mt: "\u30DE\u30EB\u30BF\u8A9E",
    nl: "\u30AA\u30E9\u30F3\u30C0\u8A9E",
    no: "\u30CE\u30EB\u30A6\u30A7\u30FC\u8A9E",
    pl: "\u30DD\u30FC\u30E9\u30F3\u30C9\u8A9E",
    pt: "\u30DD\u30EB\u30C8\u30AC\u30EB\u8A9E",
    ro: "\u30EB\u30FC\u30DE\u30CB\u30A2\u8A9E",
    ru: "\u30ED\u30B7\u30A2\u8A9E",
    sa: "\u30B5\u30F3\u30B9\u30AF\u30EA\u30C3\u30C8",
    sk: "\u30B9\u30ED\u30D0\u30AD\u30A2\u8A9E",
    sl: "\u30B9\u30ED\u30D9\u30CB\u30A2\u8A9E",
    sm: "\u30B5\u30E2\u30A2\u8A9E",
    so: "\u30BD\u30DE\u30EA\u8A9E",
    sr: "\u30BB\u30EB\u30D3\u30A2\u8A9E",
    sv: "\u30B9\u30A6\u30A7\u30FC\u30C7\u30F3\u8A9E",
    ta: "\u30BF\u30DF\u30EB\u8A9E",
    th: "\u30BF\u30A4\u8A9E",
    tr: "\u30C8\u30EB\u30B3\u8A9E",
    uk: "\u30A6\u30AF\u30E9\u30A4\u30CA\u8A9E",
    vi: "\u30D9\u30C8\u30CA\u30E0\u8A9E",
    yi: "\u30A4\u30C7\u30A3\u30C3\u30B7\u30E5\u8A9E",
    zh: "\u4E2D\u56FD\u8A9E",
    zu: "\u30BA\u30FC\u30EB\u30FC\u8A9E"
  },
  layout: {
    bitmagnet_on_service: "{{service}}\u306Ebitmagnet",
    change_theme: "\u30C6\u30FC\u30DE\u3092\u5909\u66F4",
    external_links: "\u5916\u90E8\u30EA\u30F3\u30AF",
    sponsor: "\u30B9\u30DD\u30F3\u30B5\u30FC",
    support_bitmagnet: "bitmagnet\u3092\u30B5\u30DD\u30FC\u30C8",
    translate: "\u7FFB\u8A33\u3059\u308B"
  },
  paginator: {
    first_page: "\u6700\u521D\u306E\u30DA\u30FC\u30B8",
    items_per_page: "\u30DA\u30FC\u30B8\u3042\u305F\u308A\u306E\u9805\u76EE\u6570",
    last_page: "\u6700\u5F8C\u306E\u30DA\u30FC\u30B8",
    next_page: "\u6B21\u306E\u30DA\u30FC\u30B8",
    page_x: "\u30DA\u30FC\u30B8 {{x}}",
    previous_page: "\u524D\u306E\u30DA\u30FC\u30B8",
    x_to_y: "{{x}} \u304B\u3089 {{y}} \u307E\u3067",
    x_to_y_of_z: "{{x}} \u304B\u3089 {{y}} \u307E\u3067\u306E {{z}}"
  },
  routes: {
    admin: "\u7BA1\u7406\u8005",
    dashboard: "\u30C0\u30C3\u30B7\u30E5\u30DC\u30FC\u30C9",
    home: "\u30DB\u30FC\u30E0",
    jobs: "\u30B8\u30E7\u30D6",
    queues: "\u30AD\u30E5\u30FC",
    torrents: "\u30C8\u30EC\u30F3\u30C8",
    visualize: "\u53EF\u8996\u5316"
  },
  torrents: {
    classification: "\u5206\u985E",
    clear_search: "\u691C\u7D22\u3092\u30AF\u30EA\u30A2",
    copy: "\u30B3\u30D4\u30FC",
    copy_to_clipboard: "\u30AF\u30EA\u30C3\u30D7\u30DC\u30FC\u30C9\u306B\u30B3\u30D4\u30FC",
    delete: "\u524A\u9664",
    delete_action_cannot_be_undone: "\u3053\u306E\u64CD\u4F5C\u306F\u5143\u306B\u623B\u305B\u307E\u305B\u3093",
    delete_are_you_sure: "\u3053\u306E\u30C8\u30EC\u30F3\u30C8\u3092\u524A\u9664\u3057\u3066\u3082\u3088\u308D\u3057\u3044\u3067\u3059\u304B\uFF1F",
    deselect_all: "\u3059\u3079\u3066\u306E\u9078\u629E\u3092\u89E3\u9664",
    edit_tags: "\u30BF\u30B0\u3092\u7DE8\u96C6",
    episodes: "\u30A8\u30D4\u30BD\u30FC\u30C9",
    external_links: "\u5916\u90E8\u30EA\u30F3\u30AF",
    file_index: "\u30D5\u30A1\u30A4\u30EB\u30A4\u30F3\u30C7\u30C3\u30AF\u30B9",
    file_path: "\u30D5\u30A1\u30A4\u30EB\u30D1\u30B9",
    file_size: "\u30D5\u30A1\u30A4\u30EB\u30B5\u30A4\u30BA",
    file_type: "\u30D5\u30A1\u30A4\u30EB\u30BF\u30A4\u30D7",
    files: "\u30D5\u30A1\u30A4\u30EB",
    files_count_n: "{{count}} \u4EF6\u306E\u30D5\u30A1\u30A4\u30EB",
    files_no_info: "\u30D5\u30A1\u30A4\u30EB\u60C5\u5831\u306F\u3042\u308A\u307E\u305B\u3093",
    files_single: "\u5358\u4E00\u30D5\u30A1\u30A4\u30EB",
    genres: "\u30B8\u30E3\u30F3\u30EB",
    info_hash: "\u60C5\u5831\u30CF\u30C3\u30B7\u30E5",
    info_hashes: "\u60C5\u5831\u30CF\u30C3\u30B7\u30E5",
    languages: "\u8A00\u8A9E",
    leechers: "\u30EA\u30FC\u30C1\u30E3\u30FC",
    magnet: "\u30DE\u30B0\u30CD\u30C3\u30C8",
    magnet_links: "\u30DE\u30B0\u30CD\u30C3\u30C8\u30EA\u30F3\u30AF",
    new_tag: "\u65B0\u3057\u3044\u30BF\u30B0",
    order_by: "\u4E26\u3073\u66FF\u3048",
    order_direction_toggle: "\u4E26\u3073\u9806\u3092\u5207\u308A\u66FF\u3048",
    ordering: {
      files_count: "\u30D5\u30A1\u30A4\u30EB\u6570",
      info_hash: "\u60C5\u5831\u30CF\u30C3\u30B7\u30E5",
      leechers: "\u30EA\u30FC\u30C1\u30E3\u30FC",
      name: "\u540D\u524D",
      published_at: "\u516C\u958B\u65E5\u6642",
      relevance: "\u95A2\u9023\u6027",
      seeders: "\u30B7\u30FC\u30C0\u30FC",
      size: "\u30B5\u30A4\u30BA",
      updated_at: "\u66F4\u65B0\u65E5\u6642"
    },
    original_release_date: "\u30AA\u30EA\u30B8\u30CA\u30EB\u306E\u767A\u58F2\u65E5",
    permalink: "\u30D1\u30FC\u30DE\u30EA\u30F3\u30AF",
    poster: "\u30DD\u30B9\u30BF\u30FC",
    published: "\u516C\u958B\u6E08\u307F",
    rating: "\u8A55\u4FA1",
    refresh: "\u7D50\u679C\u3092\u66F4\u65B0",
    reprocess: {
      force_rematch: "\u65E2\u306B\u4E00\u81F4\u3057\u305F\u30B3\u30F3\u30C6\u30F3\u30C4\u3092\u5F37\u5236\u7684\u306B\u518D\u4E00\u81F4\u3055\u305B\u308B",
      match_content_by_external_api_search: "\u5916\u90E8API\u691C\u7D22\u3067\u30B3\u30F3\u30C6\u30F3\u30C4\u3092\u4E00\u81F4\u3055\u305B\u308B",
      match_content_by_local_search: "\u30ED\u30FC\u30AB\u30EB\u691C\u7D22\u3067\u30B3\u30F3\u30C6\u30F3\u30C4\u3092\u4E00\u81F4\u3055\u305B\u308B",
      reprocess: "\u518D\u51E6\u7406"
    },
    s_l: "S / L",
    search: "\u691C\u7D22",
    seeders: "\u30B7\u30FC\u30C0\u30FC",
    select_all: "\u3059\u3079\u3066\u9078\u629E",
    showing_x_of_y_files: "{{x}} \u4EF6\u4E2D {{y}} \u4EF6\u8868\u793A",
    size: "\u30B5\u30A4\u30BA",
    source: "\u30C8\u30EC\u30F3\u30C8\u30BD\u30FC\u30B9",
    summary: "\u6982\u8981",
    tags: {
      delete: "\u30BF\u30B0\u3092\u524A\u9664",
      delete_tip: "\u9078\u629E\u3057\u305F\u30C8\u30EC\u30F3\u30C8\u304B\u3089\u30BF\u30B0\u3092\u524A\u9664",
      placeholder: "\u30BF\u30B0...",
      put: "\u30BF\u30B0\u3092\u4ED8\u3051\u308B",
      put_tip: "\u9078\u629E\u3057\u305F\u30C8\u30EC\u30F3\u30C8\u306B\u30BF\u30B0\u3092\u8FFD\u52A0\u3059\u308B",
      set: "\u30BF\u30B0\u3092\u8A2D\u5B9A\u3059\u308B",
      set_tip: "\u9078\u629E\u3057\u305F\u30C8\u30EC\u30F3\u30C8\u306E\u30BF\u30B0\u3092\u7F6E\u304D\u63DB\u3048\u308B"
    },
    title: "\u30BF\u30A4\u30C8\u30EB",
    toggle_drawer: "\u5F15\u304D\u51FA\u3057\u3092\u5207\u308A\u66FF\u3048",
    votes_count_n: "{{count}} \u7968"
  },
  version: {
    bitmagnet_version: "bitmagnet\u30D0\u30FC\u30B8\u30E7\u30F3 {{version}}",
    unknown: "\u4E0D\u660E"
  }
};

// src/app/i18n/translations/nl.json
var nl_default = {
  content_types: {
    plural: {
      all: "Alle",
      audiobook: "Audioboeken",
      comic: "Strips",
      ebook: "E-Boeken",
      game: "Spellen",
      movie: "Films",
      music: "Muziek",
      null: "Onbekend",
      software: "Software",
      tv_show: "Tv-programma's",
      xxx: "XXX"
    },
    singular: {
      audiobook: "Audioboek",
      comic: "Strip",
      ebook: "E-Book",
      game: "Spel",
      movie: "Film",
      music: "Muziek",
      null: "Onbekend",
      software: "Software",
      tv_show: "Tv-programma",
      xxx: "XXX"
    }
  },
  dashboard: {
    event: {
      created: "Gemaakt",
      failed: "Mislukt",
      processed: "Verwerkt",
      updated: "Bijgewerkt"
    },
    interval: {
      all: "Alle",
      days: "Dagen",
      days_1: "1 dag",
      hours: "Uren",
      hours_1: "1 uur",
      hours_12: "12 uur",
      hours_6: "6 uur",
      minutes: "Minuten",
      minutes_1: "1 minuut",
      minutes_15: "15 minuten",
      minutes_30: "30 minuten",
      minutes_5: "5 minuten",
      off: "Uit",
      seconds_10: "10 seconden",
      seconds_30: "30 seconden",
      weeks_1: "1 week"
    },
    metrics: {
      event: "Gebeurtenis",
      resolution: "Resolutie",
      throughput: "Doorvoer",
      timeframe: "Tijdsperiode",
      toggle_legend: "Legenda in-/uitschakelen"
    },
    queues: {
      created: "Gemaakt",
      created_at: "Gemaakt op",
      enqueue_jobs: "Voeg taken toe aan wachtrij",
      enqueue_torrent_processing_batch: "Torrentverwerkingsbatch toevoegen",
      failed: "Mislukt",
      force_rematch: "Forceer nieuwe match van al gematchte inhoud",
      jobs_enqueued: "Taken in wachtrij gezet",
      latency: "Latentie",
      match_content_by_external_api_search: "Match inhoud via externe API-zoekopdracht",
      match_content_by_local_search: "Match inhoud via lokale zoekopdracht",
      payload: "Inhoud",
      pending: "In afwachting",
      priority: "Prioriteit",
      process_orphaned_torrents_only: "Verwerk alleen verweesde torrents",
      processed: "Verwerkt",
      purge_jobs: "Wis taken",
      purge_queue_jobs: "Wis wachtrijtaken",
      queue: "Wachtrij",
      queue_purged: "Wachtrij gewist",
      queues: "Wachtrijen",
      ran_at: "Uitgevoerd op",
      retry: "Opnieuw proberen",
      total_counts_by_status: "Totaal aantal per status"
    }
  },
  facets: {
    content_type: "Inhoudstype",
    file_type: "Bestandstype",
    genre: "Genre",
    language: "Taal",
    queue: "Wachtrij",
    status: "Status",
    torrent_source: "Torrentbron",
    torrent_tag: "Torrenttag",
    video_resolution: "Videoresolutie",
    video_source: "Videobron"
  },
  file_types: {
    archive: "Archief",
    audio: "Audio",
    data: "Data",
    document: "Document",
    image: "Afbeelding",
    software: "Software",
    subtitles: "Ondertitels",
    unknown: "Onbekend",
    video: "Video"
  },
  general: {
    all: "Alle",
    dismiss: "Sluiten",
    error: "Fout",
    none: "Geen",
    page_not_found: "Pagina niet gevonden",
    refresh: "Verversen",
    status: "Status"
  },
  health: {
    bitmagnet_is_status: "bitmagnet is {{status}}",
    check_failed_with_error: "Controle mislukt met foutmelding",
    component: "Component",
    components: {
      dht: "DHT",
      postgres: "Postgres",
      tmdb: "TMDB"
    },
    dismiss: "Sluiten",
    error: "Fout",
    status: "Status",
    statuses: {
      degraded: "Verlaagd",
      down: "Niet beschikbaar",
      error: "Fout",
      inactive: "Inactief",
      started: "Gestart",
      unknown: "Onbekend",
      up: "Beschikbaar"
    },
    summary: "Gezondheidsrapport",
    worker: "Werker",
    workers: {
      dht_crawler: "DHT crawler",
      http_server: "HTTP server",
      queue_server: "Wachtrijserver"
    }
  },
  languages: {
    af: "Afrikaans",
    ar: "Arabisch",
    az: "Azerbeidzjaans",
    be: "Wit-Russisch",
    bg: "Bulgaars",
    bs: "Bosnisch",
    ca: "Catalaans",
    ce: "Tsjetsjeens",
    co: "Corsicaans",
    cs: "Tsjechisch",
    cy: "Welsh",
    da: "Deens",
    de: "Duits",
    el: "Grieks",
    en: "Engels",
    es: "Spaans",
    et: "Ests",
    eu: "Baskisch",
    fa: "Perzisch",
    fi: "Fins",
    fr: "Frans",
    he: "Hebreeuws",
    hi: "Hindi",
    hr: "Kroatisch",
    hu: "Hongaars",
    hy: "Armeens",
    id: "Indonesisch",
    is: "IJslands",
    it: "Italiaans",
    ja: "Japans",
    ka: "Georgisch",
    ko: "Koreaans",
    ku: "Koerdisch",
    lt: "Litouws",
    lv: "Lets",
    mi: "Maori",
    mk: "Macedonisch",
    ml: "Malayalam",
    mn: "Mongools",
    ms: "Maleis",
    mt: "Maltees",
    nl: "Nederlands",
    no: "Noors",
    pl: "Pools",
    pt: "Portugees",
    ro: "Roemeens",
    ru: "Russisch",
    sa: "Sanskriet",
    sk: "Slowaaks",
    sl: "Sloveens",
    sm: "Samoaans",
    so: "Somalisch",
    sr: "Servisch",
    sv: "Zweeds",
    ta: "Tamil",
    th: "Thais",
    tr: "Turks",
    uk: "Oekra\xEFens",
    vi: "Vietnamees",
    yi: "Jiddisch",
    zh: "Chinees",
    zu: "Zoeloe"
  },
  layout: {
    bitmagnet_on_service: "bitmagnet op {{service}}",
    change_theme: "Wijzig thema",
    external_links: "Externe links",
    sponsor: "Sponsor",
    support_bitmagnet: "Ondersteun bitmagnet",
    translate: "Vertalen"
  },
  paginator: {
    first_page: "Eerste pagina",
    items_per_page: "Items per pagina",
    last_page: "Laatste pagina",
    next_page: "Volgende pagina",
    page_x: "Pagina {{x}}",
    previous_page: "Vorige pagina",
    x_to_y: "{{x}} tot {{y}}",
    x_to_y_of_z: "{{x}} tot {{y}} van {{z}}"
  },
  routes: {
    admin: "Beheerder",
    dashboard: "Dashboard",
    home: "Home",
    jobs: "Taken",
    queues: "Wachtrijen",
    torrents: "Torrents",
    visualize: "Visualiseren"
  },
  torrents: {
    classification: "Classificatie",
    clear_search: "Zoekopdracht wissen",
    copy: "Kopi\xEBren",
    copy_to_clipboard: "Kopi\xEBren naar klembord",
    delete: "Verwijderen",
    delete_action_cannot_be_undone: "Deze actie kan niet ongedaan worden gemaakt",
    delete_are_you_sure: "Weet u zeker dat u deze torrent wilt verwijderen?",
    deselect_all: "Deselecteer alles",
    edit_tags: "Tags bewerken",
    episodes: "Afleveringen",
    external_links: "Externe links",
    file_index: "Bestandsindex",
    file_path: "Bestandspad",
    file_size: "Bestandsgrootte",
    file_type: "Bestandstype",
    files: "Bestanden",
    files_count_n: "{{count}} bestanden",
    files_no_info: "Geen informatie over bestanden beschikbaar",
    files_single: "Enkel bestand",
    genres: "Genres",
    info_hash: "Info hash",
    info_hashes: "Info hashes",
    languages: "Talen",
    leechers: "Leechers",
    magnet: "Magnet",
    magnet_links: "Magnet-links",
    new_tag: "Nieuwe tag",
    order_by: "Sorteren op",
    order_direction_toggle: "Richting omkeren",
    ordering: {
      files_count: "Aantal bestanden",
      info_hash: "Info hash",
      leechers: "Leechers",
      name: "Naam",
      published_at: "Gepubliceerd op",
      relevance: "Relevantie",
      seeders: "Seeders",
      size: "Grootte",
      updated_at: "Bijgewerkt op"
    },
    original_release_date: "Oorspronkelijke releasedatum",
    permalink: "Permalink",
    poster: "Poster",
    published: "Gepubliceerd",
    rating: "Beoordeling",
    refresh: "Vernieuw resultaten",
    reprocess: {
      force_rematch: "Forceer nieuwe match van al gematchte inhoud",
      match_content_by_external_api_search: "Match inhoud via externe API-zoekopdracht",
      match_content_by_local_search: "Match inhoud via lokale zoekopdracht",
      reprocess: "Opnieuw verwerken"
    },
    s_l: "S / L",
    search: "Zoeken",
    seeders: "Seeders",
    select_all: "Alles selecteren",
    showing_x_of_y_files: "{{x}} van {{y}} bestanden weergegeven",
    size: "Grootte",
    source: "Torrentbron",
    summary: "Samenvatting",
    tags: {
      delete: "Tags verwijderen",
      delete_tip: "Tags verwijderen van de geselecteerde torrents",
      placeholder: "Tag...",
      put: "Tags plaatsen",
      put_tip: "Tags toevoegen aan de geselecteerde torrents",
      set: "Tags instellen",
      set_tip: "Tags van de geselecteerde torrents vervangen"
    },
    title: "Titel",
    toggle_drawer: "Zijpaneel in-/uitschakelen",
    votes_count_n: "{{count}} stemmen"
  },
  version: {
    bitmagnet_version: "bitmagnet versie {{version}}",
    unknown: "onbekend"
  }
};

// src/app/i18n/translations/pt.json
var pt_default = {
  content_types: {
    plural: {
      all: "Todos",
      audiobook: "Audiolivros",
      comic: "Quadrinhos",
      ebook: "E-books",
      game: "Jogos",
      movie: "Filmes",
      music: "M\xFAsica",
      null: "Desconhecido",
      software: "Software",
      tv_show: "Programas de TV",
      xxx: "XXX"
    },
    singular: {
      audiobook: "Audiolivro",
      comic: "Quadrinho",
      ebook: "E-book",
      game: "Jogo",
      movie: "Filme",
      music: "M\xFAsica",
      null: "Desconhecido",
      software: "Software",
      tv_show: "Programa de TV",
      xxx: "XXX"
    }
  },
  dashboard: {
    event: {
      created: "Criado",
      failed: "Falhou",
      processed: "Processado",
      updated: "Atualizado"
    },
    interval: {
      all: "Todos",
      days: "Dias",
      days_1: "1 dia",
      hours: "Horas",
      hours_1: "1 hora",
      hours_12: "12 horas",
      hours_6: "6 horas",
      minutes: "Minutos",
      minutes_1: "1 minuto",
      minutes_15: "15 minutos",
      minutes_30: "30 minutos",
      minutes_5: "5 minutos",
      off: "Desligado",
      seconds_10: "10 segundos",
      seconds_30: "30 segundos",
      weeks_1: "1 semana"
    },
    metrics: {
      event: "Evento",
      resolution: "Resolu\xE7\xE3o",
      throughput: "Taxa de transfer\xEAncia",
      timeframe: "Intervalo de tempo",
      toggle_legend: "Alternar legenda"
    },
    queues: {
      created: "Criado",
      created_at: "Criado em",
      enqueue_jobs: "Enfileirar trabalhos",
      enqueue_torrent_processing_batch: "Enfileirar lote de processamento de torrents",
      failed: "Falhou",
      force_rematch: "For\xE7ar nova correspond\xEAncia de conte\xFAdo j\xE1 correspondido",
      jobs_enqueued: "Trabalhos enfileirados",
      latency: "Lat\xEAncia",
      match_content_by_external_api_search: "Correspond\xEAncia de conte\xFAdo por pesquisa de API externa",
      match_content_by_local_search: "Correspond\xEAncia de conte\xFAdo por pesquisa local",
      payload: "Carga \xFAtil",
      pending: "Pendente",
      priority: "Prioridade",
      process_orphaned_torrents_only: "Processar apenas torrents \xF3rf\xE3os",
      processed: "Processado",
      purge_jobs: "Limpar trabalhos",
      purge_queue_jobs: "Limpar trabalhos da fila",
      queue: "Fila",
      queue_purged: "Fila limpa",
      queues: "Filas",
      ran_at: "Executado em",
      retry: "Tentar novamente",
      total_counts_by_status: "Contagens totais por status"
    }
  },
  facets: {
    content_type: "Tipo de Conte\xFAdo",
    file_type: "Tipo de Arquivo",
    genre: "G\xEAnero",
    language: "Idioma",
    queue: "Fila",
    status: "Status",
    torrent_source: "Fonte do Torrent",
    torrent_tag: "Tag do Torrent",
    video_resolution: "Resolu\xE7\xE3o de V\xEDdeo",
    video_source: "Fonte de V\xEDdeo"
  },
  file_types: {
    archive: "Arquivo",
    audio: "\xC1udio",
    data: "Dados",
    document: "Documento",
    image: "Imagem",
    software: "Software",
    subtitles: "Legendas",
    unknown: "Desconhecido",
    video: "V\xEDdeo"
  },
  general: {
    all: "Todos",
    dismiss: "Fechar",
    error: "Erro",
    none: "Nenhum",
    page_not_found: "P\xE1gina n\xE3o encontrada",
    refresh: "Atualizar",
    status: "Status"
  },
  health: {
    bitmagnet_is_status: "bitmagnet est\xE1 {{status}}",
    check_failed_with_error: "Verifica\xE7\xE3o falhou com erro",
    component: "Componente",
    components: {
      dht: "DHT",
      postgres: "Postgres",
      tmdb: "TMDB"
    },
    dismiss: "Fechar",
    error: "Erro",
    status: "Status",
    statuses: {
      degraded: "Degradado",
      down: "Fora do ar",
      error: "Erro",
      inactive: "Inativo",
      started: "Iniciado",
      unknown: "Pendente",
      up: "Funcionando"
    },
    summary: "Resumo de Sa\xFAde",
    worker: "Trabalhador",
    workers: {
      dht_crawler: "Rastreamento DHT",
      http_server: "Servidor HTTP",
      queue_server: "Servidor de Fila"
    }
  },
  languages: {
    af: "Afric\xE2ner",
    ar: "\xC1rabe",
    az: "Azeri",
    be: "Bielorrusso",
    bg: "B\xFAlgaro",
    bs: "B\xF3snio",
    ca: "Catal\xE3o",
    ce: "Checheno",
    co: "Corso",
    cs: "Tcheco",
    cy: "Gal\xEAs",
    da: "Dinamarqu\xEAs",
    de: "Alem\xE3o",
    el: "Grego",
    en: "Ingl\xEAs",
    es: "Espanhol",
    et: "Estoniano",
    eu: "Basco",
    fa: "Persa",
    fi: "Finland\xEAs",
    fr: "Franc\xEAs",
    he: "Hebraico",
    hi: "Hindi",
    hr: "Croata",
    hu: "H\xFAngaro",
    hy: "Arm\xEAnio",
    id: "Indon\xE9sio",
    is: "Island\xEAs",
    it: "Italiano",
    ja: "Japon\xEAs",
    ka: "Georgiano",
    ko: "Coreano",
    ku: "Curdo",
    lt: "Lituano",
    lv: "Let\xE3o",
    mi: "Maori",
    mk: "Maced\xF4nio",
    ml: "Malaiala",
    mn: "Mongol",
    ms: "Malaio",
    mt: "Malt\xEAs",
    nl: "Holand\xEAs",
    no: "Noruegu\xEAs",
    pl: "Polon\xEAs",
    pt: "Portugu\xEAs",
    ro: "Romeno",
    ru: "Russo",
    sa: "S\xE2nscrito",
    sk: "Eslovaco",
    sl: "Esloveno",
    sm: "Samoano",
    so: "Somali",
    sr: "S\xE9rvio",
    sv: "Sueco",
    ta: "T\xE2mil",
    th: "Tailand\xEAs",
    tr: "Turco",
    uk: "Ucraniano",
    vi: "Vietnamita",
    yi: "I\xEDdiche",
    zh: "Chin\xEAs",
    zu: "Zulu"
  },
  layout: {
    bitmagnet_on_service: "bitmagnet em {{service}}",
    change_theme: "Alterar tema",
    external_links: "Links externos",
    sponsor: "Patrocinador",
    support_bitmagnet: "Apoiar bitmagnet",
    translate: "Traduzir"
  },
  paginator: {
    first_page: "Primeira p\xE1gina",
    items_per_page: "Itens por p\xE1gina",
    last_page: "\xDAltima p\xE1gina",
    next_page: "Pr\xF3xima p\xE1gina",
    page_x: "P\xE1gina {{x}}",
    previous_page: "P\xE1gina anterior",
    x_to_y: "{{x}} a {{y}}",
    x_to_y_of_z: "{{x}} a {{y}} de {{z}}"
  },
  routes: {
    admin: "Administra\xE7\xE3o",
    dashboard: "Painel",
    home: "In\xEDcio",
    jobs: "Tarefas",
    queues: "Filas",
    torrents: "Torrents",
    visualize: "Visualizar"
  },
  torrents: {
    classification: "Classifica\xE7\xE3o",
    clear_search: "Limpar Pesquisa",
    copy: "Copiar",
    copy_to_clipboard: "Copiar para a \xE1rea de transfer\xEAncia",
    delete: "Excluir",
    delete_action_cannot_be_undone: "Esta a\xE7\xE3o n\xE3o pode ser desfeita",
    delete_are_you_sure: "Tem certeza de que deseja excluir este torrent?",
    deselect_all: "Desmarcar todos",
    edit_tags: "Editar tags",
    episodes: "Epis\xF3dios",
    external_links: "Links externos",
    file_index: "\xCDndice de arquivo",
    file_path: "Caminho do arquivo",
    file_size: "Tamanho do arquivo",
    file_type: "Tipo de arquivo",
    files: "Arquivos",
    files_count_n: "{{count}} arquivos",
    files_no_info: "Sem informa\xE7\xF5es de arquivos dispon\xEDveis",
    files_single: "Arquivo \xFAnico",
    genres: "G\xEAneros",
    info_hash: "Hash de informa\xE7\xE3o",
    info_hashes: "Hashes de informa\xE7\xE3o",
    languages: "Idiomas",
    leechers: "Leechers",
    magnet: "Magnet",
    magnet_links: "Links magnet",
    new_tag: "Nova tag",
    order_by: "Ordenar por",
    order_direction_toggle: "Inverter dire\xE7\xE3o",
    ordering: {
      files_count: "Contagem de arquivos",
      info_hash: "Hash de informa\xE7\xE3o",
      leechers: "Leechers",
      name: "Nome",
      published_at: "Publicado em",
      relevance: "Relev\xE2ncia",
      seeders: "Seeders",
      size: "Tamanho",
      updated_at: "Atualizado em"
    },
    original_release_date: "Data de lan\xE7amento original",
    permalink: "Link permanente",
    poster: "P\xF4ster",
    published: "Publicado",
    rating: "Classifica\xE7\xE3o",
    refresh: "Atualizar resultados",
    reprocess: {
      force_rematch: "For\xE7ar nova correspond\xEAncia de conte\xFAdo j\xE1 correspondido",
      match_content_by_external_api_search: "Corresponder conte\xFAdo por pesquisa de API externa",
      match_content_by_local_search: "Corresponder conte\xFAdo por pesquisa local",
      reprocess: "Reprocessar"
    },
    s_l: "S / L",
    search: "Buscar",
    seeders: "Seeders",
    select_all: "Selecionar tudo",
    showing_x_of_y_files: "Mostrando {{x}} de {{y}} arquivos",
    size: "Tamanho",
    source: "Fonte do torrent",
    summary: "Resumo",
    tags: {
      delete: "Excluir tags",
      delete_tip: "Remover tags dos torrents selecionados",
      placeholder: "Tag...",
      put: "Colocar tags",
      put_tip: "Adicionar tags aos torrents selecionados",
      set: "Definir tags",
      set_tip: "Substituir tags dos torrents selecionados"
    },
    title: "T\xEDtulo",
    toggle_drawer: "Alternar gaveta",
    votes_count_n: "{{count}} votos"
  },
  version: {
    bitmagnet_version: "Vers\xE3o do bitmagnet {{version}}",
    unknown: "desconhecido"
  }
};

// src/app/i18n/translations/ru.json
var ru_default = {
  content_types: {
    plural: {
      all: "\u0412\u0441\u0435",
      audiobook: "\u0410\u0443\u0434\u0438\u043E\u043A\u043D\u0438\u0433\u0438",
      comic: "\u041A\u043E\u043C\u0438\u043A\u0441\u044B",
      ebook: "\u042D\u043B\u0435\u043A\u0442\u0440\u043E\u043D\u043D\u044B\u0435 \u043A\u043D\u0438\u0433\u0438",
      game: "\u0418\u0433\u0440\u044B",
      movie: "\u0424\u0438\u043B\u044C\u043C\u044B",
      music: "\u041C\u0443\u0437\u044B\u043A\u0430",
      null: "\u041D\u0435\u0438\u0437\u0432\u0435\u0441\u0442\u043D\u043E",
      software: "\u041F\u0440\u043E\u0433\u0440\u0430\u043C\u043C\u044B",
      tv_show: "\u0422\u0435\u043B\u0435\u043F\u0435\u0440\u0435\u0434\u0430\u0447\u0438",
      xxx: "\u041F\u043E\u0440\u043D\u043E"
    },
    singular: {
      audiobook: "\u0410\u0443\u0434\u0438\u043E\u043A\u043D\u0438\u0433\u0430",
      comic: "\u041A\u043E\u043C\u0438\u043A\u0441",
      ebook: "\u042D\u043B\u0435\u043A\u0442\u0440\u043E\u043D\u043D\u0430\u044F \u043A\u043D\u0438\u0433\u0430",
      game: "\u0418\u0433\u0440\u0430",
      movie: "\u0424\u0438\u043B\u044C\u043C",
      music: "\u041C\u0443\u0437\u044B\u043A\u0430",
      null: "\u041D\u0435\u0438\u0437\u0432\u0435\u0441\u0442\u043D\u043E",
      software: "\u041F\u0440\u043E\u0433\u0440\u0430\u043C\u043C\u0430",
      tv_show: "\u0422\u0435\u043B\u0435\u043F\u0435\u0440\u0435\u0434\u0430\u0447\u0430",
      xxx: "\u041F\u043E\u0440\u043D\u043E"
    }
  },
  dashboard: {
    event: {
      created: "\u0421\u043E\u0437\u0434\u0430\u043D\u043E",
      failed: "\u041E\u0448\u0438\u0431\u043A\u0430",
      processed: "\u041E\u0431\u0440\u0430\u0431\u043E\u0442\u0430\u043D\u043E",
      updated: "\u041E\u0431\u043D\u043E\u0432\u043B\u0435\u043D\u043E"
    },
    interval: {
      all: "\u0412\u0441\u0435",
      days: "\u0414\u043D\u0438",
      days_1: "1 \u0434\u0435\u043D\u044C",
      hours: "\u0427\u0430\u0441\u044B",
      hours_1: "1 \u0447\u0430\u0441",
      hours_12: "12 \u0447\u0430\u0441\u043E\u0432",
      hours_6: "6 \u0447\u0430\u0441\u043E\u0432",
      minutes: "\u041C\u0438\u043D\u0443\u0442\u044B",
      minutes_1: "1 \u043C\u0438\u043D\u0443\u0442\u0430",
      minutes_15: "15 \u043C\u0438\u043D\u0443\u0442",
      minutes_30: "30 \u043C\u0438\u043D\u0443\u0442",
      minutes_5: "5 \u043C\u0438\u043D\u0443\u0442",
      off: "\u041E\u0442\u043A\u043B\u044E\u0447\u0435\u043D\u043E",
      seconds_10: "10 \u0441\u0435\u043A\u0443\u043D\u0434",
      seconds_30: "30 \u0441\u0435\u043A\u0443\u043D\u0434",
      weeks_1: "1 \u043D\u0435\u0434\u0435\u043B\u044F"
    },
    metrics: {
      event: "\u0421\u043E\u0431\u044B\u0442\u0438\u0435",
      resolution: "\u0420\u0430\u0437\u0440\u0435\u0448\u0435\u043D\u0438\u0435",
      throughput: "\u041F\u0440\u043E\u043F\u0443\u0441\u043A\u043D\u0430\u044F \u0441\u043F\u043E\u0441\u043E\u0431\u043D\u043E\u0441\u0442\u044C",
      timeframe: "\u041F\u0435\u0440\u0438\u043E\u0434 \u0432\u0440\u0435\u043C\u0435\u043D\u0438",
      toggle_legend: "\u041F\u0435\u0440\u0435\u043A\u043B\u044E\u0447\u0438\u0442\u044C \u043B\u0435\u0433\u0435\u043D\u0434\u0443"
    },
    queues: {
      created: "\u0421\u043E\u0437\u0434\u0430\u043D\u043E",
      created_at: "\u0421\u043E\u0437\u0434\u0430\u043D\u043E \u0432",
      enqueue_jobs: "\u0414\u043E\u0431\u0430\u0432\u0438\u0442\u044C \u0437\u0430\u0434\u0430\u043D\u0438\u044F \u0432 \u043E\u0447\u0435\u0440\u0435\u0434\u044C",
      enqueue_torrent_processing_batch: "\u0414\u043E\u0431\u0430\u0432\u0438\u0442\u044C \u043F\u0430\u0440\u0442\u0438\u044E \u0434\u043B\u044F \u043E\u0431\u0440\u0430\u0431\u043E\u0442\u043A\u0438 \u0442\u043E\u0440\u0440\u0435\u043D\u0442\u043E\u0432",
      failed: "\u041E\u0448\u0438\u0431\u043A\u0430",
      force_rematch: "\u041F\u0440\u0438\u043D\u0443\u0434\u0438\u0442\u0435\u043B\u044C\u043D\u043E\u0435 \u043F\u043E\u0432\u0442\u043E\u0440\u043D\u043E\u0435 \u0441\u043E\u043F\u043E\u0441\u0442\u0430\u0432\u043B\u0435\u043D\u0438\u0435 \u0443\u0436\u0435 \u0441\u043E\u043F\u043E\u0441\u0442\u0430\u0432\u043B\u0435\u043D\u043D\u043E\u0433\u043E \u043A\u043E\u043D\u0442\u0435\u043D\u0442\u0430",
      jobs_enqueued: "\u0417\u0430\u0434\u0430\u043D\u0438\u044F \u0434\u043E\u0431\u0430\u0432\u043B\u0435\u043D\u044B \u0432 \u043E\u0447\u0435\u0440\u0435\u0434\u044C",
      latency: "\u0417\u0430\u0434\u0435\u0440\u0436\u043A\u0430",
      match_content_by_external_api_search: "\u0421\u043E\u043F\u043E\u0441\u0442\u0430\u0432\u0438\u0442\u044C \u043A\u043E\u043D\u0442\u0435\u043D\u0442 \u0447\u0435\u0440\u0435\u0437 \u0432\u043D\u0435\u0448\u043D\u0438\u0439 API",
      match_content_by_local_search: "\u0421\u043E\u043F\u043E\u0441\u0442\u0430\u0432\u0438\u0442\u044C \u043A\u043E\u043D\u0442\u0435\u043D\u0442 \u0447\u0435\u0440\u0435\u0437 \u043B\u043E\u043A\u0430\u043B\u044C\u043D\u044B\u0439 \u043F\u043E\u0438\u0441\u043A",
      payload: "\u0414\u0430\u043D\u043D\u044B\u0435",
      pending: "\u0412 \u043E\u0436\u0438\u0434\u0430\u043D\u0438\u0438",
      priority: "\u041F\u0440\u0438\u043E\u0440\u0438\u0442\u0435\u0442",
      process_orphaned_torrents_only: "\u041E\u0431\u0440\u0430\u0431\u0430\u0442\u044B\u0432\u0430\u0442\u044C \u0442\u043E\u043B\u044C\u043A\u043E \xAB\u043E\u0441\u0438\u0440\u043E\u0442\u0435\u0432\u0448\u0438\u0435\xBB \u0442\u043E\u0440\u0440\u0435\u043D\u0442\u044B",
      processed: "\u041E\u0431\u0440\u0430\u0431\u043E\u0442\u0430\u043D\u043E",
      purge_jobs: "\u041E\u0447\u0438\u0441\u0442\u0438\u0442\u044C \u0437\u0430\u0434\u0430\u043D\u0438\u044F",
      purge_queue_jobs: "\u041E\u0447\u0438\u0441\u0442\u0438\u0442\u044C \u0437\u0430\u0434\u0430\u043D\u0438\u044F \u043E\u0447\u0435\u0440\u0435\u0434\u0438",
      queue: "\u041E\u0447\u0435\u0440\u0435\u0434\u044C",
      queue_purged: "\u041E\u0447\u0435\u0440\u0435\u0434\u044C \u043E\u0447\u0438\u0449\u0435\u043D\u0430",
      queues: "\u041E\u0447\u0435\u0440\u0435\u0434\u0438",
      ran_at: "\u0412\u044B\u043F\u043E\u043B\u043D\u0435\u043D\u043E \u0432",
      retry: "\u041F\u043E\u0432\u0442\u043E\u0440\u0438\u0442\u044C",
      total_counts_by_status: "\u041E\u0431\u0449\u0435\u0435 \u043A\u043E\u043B\u0438\u0447\u0435\u0441\u0442\u0432\u043E \u043F\u043E \u0441\u0442\u0430\u0442\u0443\u0441\u0430\u043C"
    }
  },
  facets: {
    content_type: "\u0422\u0438\u043F \u043A\u043E\u043D\u0442\u0435\u043D\u0442\u0430",
    file_type: "\u0422\u0438\u043F \u0444\u0430\u0439\u043B\u0430",
    genre: "\u0416\u0430\u043D\u0440",
    language: "\u042F\u0437\u044B\u043A",
    queue: "\u041E\u0447\u0435\u0440\u0435\u0434\u044C",
    status: "\u0421\u0442\u0430\u0442\u0443\u0441",
    torrent_source: "\u0418\u0441\u0442\u043E\u0447\u043D\u0438\u043A \u0442\u043E\u0440\u0440\u0435\u043D\u0442\u0430",
    torrent_tag: "\u0422\u0435\u0433 \u0442\u043E\u0440\u0440\u0435\u043D\u0442\u0430",
    video_resolution: "\u0420\u0430\u0437\u0440\u0435\u0448\u0435\u043D\u0438\u0435 \u0432\u0438\u0434\u0435\u043E",
    video_source: "\u0418\u0441\u0442\u043E\u0447\u043D\u0438\u043A \u0432\u0438\u0434\u0435\u043E"
  },
  file_types: {
    archive: "\u0410\u0440\u0445\u0438\u0432",
    audio: "\u0410\u0443\u0434\u0438\u043E",
    data: "\u0414\u0430\u043D\u043D\u044B\u0435",
    document: "\u0414\u043E\u043A\u0443\u043C\u0435\u043D\u0442",
    image: "\u0418\u0437\u043E\u0431\u0440\u0430\u0436\u0435\u043D\u0438\u0435",
    software: "\u041F\u0440\u043E\u0433\u0440\u0430\u043C\u043C\u044B",
    subtitles: "\u0421\u0443\u0431\u0442\u0438\u0442\u0440\u044B",
    unknown: "\u041D\u0435\u0438\u0437\u0432\u0435\u0441\u0442\u043D\u043E",
    video: "\u0412\u0438\u0434\u0435\u043E"
  },
  general: {
    all: "\u0412\u0441\u0435",
    dismiss: "\u0417\u0430\u043A\u0440\u044B\u0442\u044C",
    error: "\u041E\u0448\u0438\u0431\u043A\u0430",
    none: "\u041D\u0435\u0442",
    page_not_found: "\u0421\u0442\u0440\u0430\u043D\u0438\u0446\u0430 \u043D\u0435 \u043D\u0430\u0439\u0434\u0435\u043D\u0430",
    refresh: "\u041E\u0431\u043D\u043E\u0432\u0438\u0442\u044C",
    status: "\u0421\u0442\u0430\u0442\u0443\u0441"
  },
  health: {
    bitmagnet_is_status: "bitmagnet {{status}}",
    check_failed_with_error: "\u041F\u0440\u043E\u0432\u0435\u0440\u043A\u0430 \u0437\u0430\u0432\u0435\u0440\u0448\u0438\u043B\u0430\u0441\u044C \u0441 \u043E\u0448\u0438\u0431\u043A\u043E\u0439",
    component: "\u041A\u043E\u043C\u043F\u043E\u043D\u0435\u043D\u0442",
    components: {
      dht: "DHT",
      postgres: "Postgres",
      tmdb: "TMDB"
    },
    dismiss: "\u0417\u0430\u043A\u0440\u044B\u0442\u044C",
    error: "\u041E\u0448\u0438\u0431\u043A\u0430",
    status: "\u0421\u0442\u0430\u0442\u0443\u0441",
    statuses: {
      degraded: "\u0421\u043D\u0438\u0436\u0435\u043D\u0430 \u043F\u0440\u043E\u0438\u0437\u0432\u043E\u0434\u0438\u0442\u0435\u043B\u044C\u043D\u043E\u0441\u0442\u044C",
      down: "\u041D\u0435 \u0440\u0430\u0431\u043E\u0442\u0430\u0435\u0442",
      error: "\u041E\u0448\u0438\u0431\u043A\u0430",
      inactive: "\u041D\u0435\u0430\u043A\u0442\u0438\u0432\u0435\u043D",
      started: "\u0417\u0430\u043F\u0443\u0449\u0435\u043D\u043E",
      unknown: "\u041D\u0435\u0438\u0437\u0432\u0435\u0441\u0442\u043D\u043E",
      up: "\u0420\u0430\u0431\u043E\u0442\u0430\u0435\u0442"
    },
    summary: "\u0421\u0432\u043E\u0434\u043A\u0430 \u0441\u043E\u0441\u0442\u043E\u044F\u043D\u0438\u044F",
    worker: "\u0420\u0430\u0431\u043E\u0447\u0438\u0439",
    workers: {
      dht_crawler: "DHT \u043E\u0431\u0445\u043E\u0434\u0447\u0438\u043A",
      http_server: "HTTP \u0441\u0435\u0440\u0432\u0435\u0440",
      queue_server: "\u0421\u0435\u0440\u0432\u0435\u0440 \u043E\u0447\u0435\u0440\u0435\u0434\u0435\u0439"
    }
  },
  languages: {
    af: "\u0410\u0444\u0440\u0438\u043A\u0430\u0430\u043D\u0441",
    ar: "\u0410\u0440\u0430\u0431\u0441\u043A\u0438\u0439",
    az: "\u0410\u0437\u0435\u0440\u0431\u0430\u0439\u0434\u0436\u0430\u043D\u0441\u043A\u0438\u0439",
    be: "\u0411\u0435\u043B\u043E\u0440\u0443\u0441\u0441\u043A\u0438\u0439",
    bg: "\u0411\u043E\u043B\u0433\u0430\u0440\u0441\u043A\u0438\u0439",
    bs: "\u0411\u043E\u0441\u043D\u0438\u0439\u0441\u043A\u0438\u0439",
    ca: "\u041A\u0430\u0442\u0430\u043B\u0430\u043D\u0441\u043A\u0438\u0439",
    ce: "\u0427\u0435\u0447\u0435\u043D\u0441\u043A\u0438\u0439",
    co: "\u041A\u043E\u0440\u0441\u0438\u043A\u0430\u043D\u0441\u043A\u0438\u0439",
    cs: "\u0427\u0435\u0448\u0441\u043A\u0438\u0439",
    cy: "\u0412\u0430\u043B\u043B\u0438\u0439\u0441\u043A\u0438\u0439",
    da: "\u0414\u0430\u0442\u0441\u043A\u0438\u0439",
    de: "\u041D\u0435\u043C\u0435\u0446\u043A\u0438\u0439",
    el: "\u0413\u0440\u0435\u0447\u0435\u0441\u043A\u0438\u0439",
    en: "\u0410\u043D\u0433\u043B\u0438\u0439\u0441\u043A\u0438\u0439",
    es: "\u0418\u0441\u043F\u0430\u043D\u0441\u043A\u0438\u0439",
    et: "\u042D\u0441\u0442\u043E\u043D\u0441\u043A\u0438\u0439",
    eu: "\u0411\u0430\u0441\u043A\u0441\u043A\u0438\u0439",
    fa: "\u041F\u0435\u0440\u0441\u0438\u0434\u0441\u043A\u0438\u0439",
    fi: "\u0424\u0438\u043D\u0441\u043A\u0438\u0439",
    fr: "\u0424\u0440\u0430\u043D\u0446\u0443\u0437\u0441\u043A\u0438\u0439",
    he: "\u0418\u0432\u0440\u0438\u0442",
    hi: "\u0425\u0438\u043D\u0434\u0438",
    hr: "\u0425\u043E\u0440\u0432\u0430\u0442\u0441\u043A\u0438\u0439",
    hu: "\u0412\u0435\u043D\u0433\u0435\u0440\u0441\u043A\u0438\u0439",
    hy: "\u0410\u0440\u043C\u044F\u043D\u0441\u043A\u0438\u0439",
    id: "\u0418\u043D\u0434\u043E\u043D\u0435\u0437\u0438\u0439\u0441\u043A\u0438\u0439",
    is: "\u0418\u0441\u043B\u0430\u043D\u0434\u0441\u043A\u0438\u0439",
    it: "\u0418\u0442\u0430\u043B\u044C\u044F\u043D\u0441\u043A\u0438\u0439",
    ja: "\u042F\u043F\u043E\u043D\u0441\u043A\u0438\u0439",
    ka: "\u0413\u0440\u0443\u0437\u0438\u043D\u0441\u043A\u0438\u0439",
    ko: "\u041A\u043E\u0440\u0435\u0439\u0441\u043A\u0438\u0439",
    ku: "\u041A\u0443\u0440\u0434\u0441\u043A\u0438\u0439",
    lt: "\u041B\u0438\u0442\u043E\u0432\u0441\u043A\u0438\u0439",
    lv: "\u041B\u0430\u0442\u044B\u0448\u0441\u043A\u0438\u0439",
    mi: "\u041C\u0430\u043E\u0440\u0438",
    mk: "\u041C\u0430\u043A\u0435\u0434\u043E\u043D\u0441\u043A\u0438\u0439",
    ml: "\u041C\u0430\u043B\u0430\u044F\u043B\u0430\u043C",
    mn: "\u041C\u043E\u043D\u0433\u043E\u043B\u044C\u0441\u043A\u0438\u0439",
    ms: "\u041C\u0430\u043B\u0430\u0439\u0441\u043A\u0438\u0439",
    mt: "\u041C\u0430\u043B\u044C\u0442\u0438\u0439\u0441\u043A\u0438\u0439",
    nl: "\u041D\u0438\u0434\u0435\u0440\u043B\u0430\u043D\u0434\u0441\u043A\u0438\u0439",
    no: "\u041D\u043E\u0440\u0432\u0435\u0436\u0441\u043A\u0438\u0439",
    pl: "\u041F\u043E\u043B\u044C\u0441\u043A\u0438\u0439",
    pt: "\u041F\u043E\u0440\u0442\u0443\u0433\u0430\u043B\u044C\u0441\u043A\u0438\u0439",
    ro: "\u0420\u0443\u043C\u044B\u043D\u0441\u043A\u0438\u0439",
    ru: "\u0420\u0443\u0441\u0441\u043A\u0438\u0439",
    sa: "\u0421\u0430\u043D\u0441\u043A\u0440\u0438\u0442",
    sk: "\u0421\u043B\u043E\u0432\u0430\u0446\u043A\u0438\u0439",
    sl: "\u0421\u043B\u043E\u0432\u0435\u043D\u0441\u043A\u0438\u0439",
    sm: "\u0421\u0430\u043C\u043E\u0430\u043D\u0441\u043A\u0438\u0439",
    so: "\u0421\u043E\u043C\u0430\u043B\u0438\u0439\u0441\u043A\u0438\u0439",
    sr: "\u0421\u0435\u0440\u0431\u0441\u043A\u0438\u0439",
    sv: "\u0428\u0432\u0435\u0434\u0441\u043A\u0438\u0439",
    ta: "\u0422\u0430\u043C\u0438\u043B\u044C\u0441\u043A\u0438\u0439",
    th: "\u0422\u0430\u0439\u0441\u043A\u0438\u0439",
    tr: "\u0422\u0443\u0440\u0435\u0446\u043A\u0438\u0439",
    uk: "\u0423\u043A\u0440\u0430\u0438\u043D\u0441\u043A\u0438\u0439",
    vi: "\u0412\u044C\u0435\u0442\u043D\u0430\u043C\u0441\u043A\u0438\u0439",
    yi: "\u0418\u0434\u0438\u0448",
    zh: "\u041A\u0438\u0442\u0430\u0439\u0441\u043A\u0438\u0439",
    zu: "\u0417\u0443\u043B\u0443\u0441\u0441\u043A\u0438\u0439"
  },
  layout: {
    bitmagnet_on_service: "bitmagnet \u043D\u0430 {{service}}",
    change_theme: "\u0418\u0437\u043C\u0435\u043D\u0438\u0442\u044C \u0442\u0435\u043C\u0443",
    external_links: "\u0412\u043D\u0435\u0448\u043D\u0438\u0435 \u0441\u0441\u044B\u043B\u043A\u0438",
    sponsor: "\u0421\u043F\u043E\u043D\u0441\u043E\u0440",
    support_bitmagnet: "\u041F\u043E\u0434\u0434\u0435\u0440\u0436\u0430\u0442\u044C bitmagnet",
    translate: "\u041F\u0435\u0440\u0435\u0432\u0435\u0441\u0442\u0438"
  },
  paginator: {
    first_page: "\u041F\u0435\u0440\u0432\u0430\u044F \u0441\u0442\u0440\u0430\u043D\u0438\u0446\u0430",
    items_per_page: "\u042D\u043B\u0435\u043C\u0435\u043D\u0442\u043E\u0432 \u043D\u0430 \u0441\u0442\u0440\u0430\u043D\u0438\u0446\u0435",
    last_page: "\u041F\u043E\u0441\u043B\u0435\u0434\u043D\u044F\u044F \u0441\u0442\u0440\u0430\u043D\u0438\u0446\u0430",
    next_page: "\u0421\u043B\u0435\u0434\u0443\u044E\u0449\u0430\u044F \u0441\u0442\u0440\u0430\u043D\u0438\u0446\u0430",
    page_x: "\u0421\u0442\u0440\u0430\u043D\u0438\u0446\u0430 {{x}}",
    previous_page: "\u041F\u0440\u0435\u0434\u044B\u0434\u0443\u0449\u0430\u044F \u0441\u0442\u0440\u0430\u043D\u0438\u0446\u0430",
    x_to_y: "{{x}} \u0434\u043E {{y}}",
    x_to_y_of_z: "{{x}} \u0434\u043E {{y}} \u0438\u0437 {{z}}"
  },
  routes: {
    admin: "\u0410\u0434\u043C\u0438\u043D\u0438\u0441\u0442\u0440\u0430\u0442\u043E\u0440",
    dashboard: "\u041F\u0430\u043D\u0435\u043B\u044C \u0443\u043F\u0440\u0430\u0432\u043B\u0435\u043D\u0438\u044F",
    home: "\u0413\u043B\u0430\u0432\u043D\u0430\u044F",
    jobs: "\u0417\u0430\u0434\u0430\u043D\u0438\u044F",
    queues: "\u041E\u0447\u0435\u0440\u0435\u0434\u0438",
    torrents: "\u0422\u043E\u0440\u0440\u0435\u043D\u0442\u044B",
    visualize: "\u0412\u0438\u0437\u0443\u0430\u043B\u0438\u0437\u0430\u0446\u0438\u044F"
  },
  torrents: {
    classification: "\u041A\u043B\u0430\u0441\u0441\u0438\u0444\u0438\u043A\u0430\u0446\u0438\u044F",
    clear_search: "\u041E\u0447\u0438\u0441\u0442\u0438\u0442\u044C \u043F\u043E\u0438\u0441\u043A",
    copy: "\u041A\u043E\u043F\u0438\u0440\u043E\u0432\u0430\u0442\u044C",
    copy_to_clipboard: "\u041A\u043E\u043F\u0438\u0440\u043E\u0432\u0430\u0442\u044C \u0432 \u0431\u0443\u0444\u0435\u0440 \u043E\u0431\u043C\u0435\u043D\u0430",
    delete: "\u0423\u0434\u0430\u043B\u0438\u0442\u044C",
    delete_action_cannot_be_undone: "\u042D\u0442\u043E \u0434\u0435\u0439\u0441\u0442\u0432\u0438\u0435 \u043D\u0435\u043B\u044C\u0437\u044F \u043E\u0442\u043C\u0435\u043D\u0438\u0442\u044C",
    delete_are_you_sure: "\u0412\u044B \u0443\u0432\u0435\u0440\u0435\u043D\u044B, \u0447\u0442\u043E \u0445\u043E\u0442\u0438\u0442\u0435 \u0443\u0434\u0430\u043B\u0438\u0442\u044C \u044D\u0442\u043E\u0442 \u0442\u043E\u0440\u0440\u0435\u043D\u0442?",
    deselect_all: "\u0421\u043D\u044F\u0442\u044C \u0432\u044B\u0434\u0435\u043B\u0435\u043D\u0438\u0435",
    edit_tags: "\u0420\u0435\u0434\u0430\u043A\u0442\u0438\u0440\u043E\u0432\u0430\u0442\u044C \u0442\u0435\u0433\u0438",
    episodes: "\u042D\u043F\u0438\u0437\u043E\u0434\u044B",
    external_links: "\u0412\u043D\u0435\u0448\u043D\u0438\u0435 \u0441\u0441\u044B\u043B\u043A\u0438",
    file_index: "\u0418\u043D\u0434\u0435\u043A\u0441 \u0444\u0430\u0439\u043B\u0430",
    file_path: "\u041F\u0443\u0442\u044C \u043A \u0444\u0430\u0439\u043B\u0443",
    file_size: "\u0420\u0430\u0437\u043C\u0435\u0440 \u0444\u0430\u0439\u043B\u0430",
    file_type: "\u0422\u0438\u043F \u0444\u0430\u0439\u043B\u0430",
    files: "\u0424\u0430\u0439\u043B\u044B",
    files_count_n: "{{count}} \u0444\u0430\u0439\u043B\u043E\u0432",
    files_no_info: "\u0418\u043D\u0444\u043E\u0440\u043C\u0430\u0446\u0438\u044F \u043E \u0444\u0430\u0439\u043B\u0430\u0445 \u043D\u0435\u0434\u043E\u0441\u0442\u0443\u043F\u043D\u0430",
    files_single: "\u041E\u0434\u0438\u043D \u0444\u0430\u0439\u043B",
    genres: "\u0416\u0430\u043D\u0440\u044B",
    info_hash: "\u0425\u044D\u0448 \u0438\u043D\u0444\u043E\u0440\u043C\u0430\u0446\u0438\u0438",
    info_hashes: "\u0425\u044D\u0448\u0438 \u0438\u043D\u0444\u043E\u0440\u043C\u0430\u0446\u0438\u0438",
    languages: "\u042F\u0437\u044B\u043A\u0438",
    leechers: "\u041B\u0438\u0447\u0435\u0440\u044B",
    magnet: "\u041C\u0430\u0433\u043D\u0435\u0442",
    magnet_links: "\u041C\u0430\u0433\u043D\u0435\u0442 \u0441\u0441\u044B\u043B\u043A\u0438",
    new_tag: "\u041D\u043E\u0432\u044B\u0439 \u0442\u0435\u0433",
    order_by: "\u0421\u043E\u0440\u0442\u0438\u0440\u043E\u0432\u0430\u0442\u044C \u043F\u043E",
    order_direction_toggle: "\u041F\u0435\u0440\u0435\u043A\u043B\u044E\u0447\u0438\u0442\u044C \u043D\u0430\u043F\u0440\u0430\u0432\u043B\u0435\u043D\u0438\u0435",
    ordering: {
      files_count: "\u041A\u043E\u043B\u0438\u0447\u0435\u0441\u0442\u0432\u043E \u0444\u0430\u0439\u043B\u043E\u0432",
      info_hash: "\u0425\u044D\u0448 \u0438\u043D\u0444\u043E\u0440\u043C\u0430\u0446\u0438\u0438",
      leechers: "\u041B\u0438\u0447\u0435\u0440\u044B",
      name: "\u0418\u043C\u044F",
      published_at: "\u0414\u0430\u0442\u0430 \u043F\u0443\u0431\u043B\u0438\u043A\u0430\u0446\u0438\u0438",
      relevance: "\u0410\u043A\u0442\u0443\u0430\u043B\u044C\u043D\u043E\u0441\u0442\u044C",
      seeders: "\u0421\u0438\u0434\u0435\u0440\u044B",
      size: "\u0420\u0430\u0437\u043C\u0435\u0440",
      updated_at: "\u0414\u0430\u0442\u0430 \u043E\u0431\u043D\u043E\u0432\u043B\u0435\u043D\u0438\u044F"
    },
    original_release_date: "\u041E\u0440\u0438\u0433\u0438\u043D\u0430\u043B\u044C\u043D\u0430\u044F \u0434\u0430\u0442\u0430 \u0432\u044B\u043F\u0443\u0441\u043A\u0430",
    permalink: "\u041F\u043E\u0441\u0442\u043E\u044F\u043D\u043D\u0430\u044F \u0441\u0441\u044B\u043B\u043A\u0430",
    poster: "\u041F\u043E\u0441\u0442\u0435\u0440",
    published: "\u041E\u043F\u0443\u0431\u043B\u0438\u043A\u043E\u0432\u0430\u043D\u043E",
    rating: "\u0420\u0435\u0439\u0442\u0438\u043D\u0433",
    refresh: "\u041E\u0431\u043D\u043E\u0432\u0438\u0442\u044C \u0440\u0435\u0437\u0443\u043B\u044C\u0442\u0430\u0442\u044B",
    reprocess: {
      force_rematch: "\u041F\u0440\u0438\u043D\u0443\u0434\u0438\u0442\u0435\u043B\u044C\u043D\u043E\u0435 \u043F\u043E\u0432\u0442\u043E\u0440\u043D\u043E\u0435 \u0441\u043E\u043F\u043E\u0441\u0442\u0430\u0432\u043B\u0435\u043D\u0438\u0435 \u0443\u0436\u0435 \u0441\u043E\u043F\u043E\u0441\u0442\u0430\u0432\u043B\u0435\u043D\u043D\u043E\u0433\u043E \u043A\u043E\u043D\u0442\u0435\u043D\u0442\u0430",
      match_content_by_external_api_search: "\u0421\u043E\u043F\u043E\u0441\u0442\u0430\u0432\u0438\u0442\u044C \u043A\u043E\u043D\u0442\u0435\u043D\u0442 \u0447\u0435\u0440\u0435\u0437 \u0432\u043D\u0435\u0448\u043D\u0438\u0439 API",
      match_content_by_local_search: "\u0421\u043E\u043F\u043E\u0441\u0442\u0430\u0432\u0438\u0442\u044C \u043A\u043E\u043D\u0442\u0435\u043D\u0442 \u0447\u0435\u0440\u0435\u0437 \u043B\u043E\u043A\u0430\u043B\u044C\u043D\u044B\u0439 \u043F\u043E\u0438\u0441\u043A",
      reprocess: "\u041F\u0435\u0440\u0435\u0440\u0430\u0431\u043E\u0442\u0430\u0442\u044C"
    },
    s_l: "S / L",
    search: "\u041F\u043E\u0438\u0441\u043A",
    seeders: "\u0421\u0438\u0434\u0435\u0440\u044B",
    select_all: "\u0412\u044B\u0431\u0440\u0430\u0442\u044C \u0432\u0441\u0435",
    showing_x_of_y_files: "\u041F\u043E\u043A\u0430\u0437\u0430\u043D\u043E {{x}} \u0438\u0437 {{y}} \u0444\u0430\u0439\u043B\u043E\u0432",
    size: "\u0420\u0430\u0437\u043C\u0435\u0440",
    source: "\u0418\u0441\u0442\u043E\u0447\u043D\u0438\u043A \u0442\u043E\u0440\u0440\u0435\u043D\u0442\u0430",
    summary: "\u0421\u0432\u043E\u0434\u043A\u0430",
    tags: {
      delete: "\u0423\u0434\u0430\u043B\u0438\u0442\u044C \u0442\u0435\u0433\u0438",
      delete_tip: "\u0423\u0434\u0430\u043B\u0438\u0442\u044C \u0442\u0435\u0433\u0438 \u0438\u0437 \u0432\u044B\u0431\u0440\u0430\u043D\u043D\u044B\u0445 \u0442\u043E\u0440\u0440\u0435\u043D\u0442\u043E\u0432",
      placeholder: "\u0422\u0435\u0433...",
      put: "\u041F\u043E\u043C\u0435\u0441\u0442\u0438\u0442\u044C \u0442\u0435\u0433\u0438",
      put_tip: "\u0414\u043E\u0431\u0430\u0432\u0438\u0442\u044C \u0442\u0435\u0433\u0438 \u043A \u0432\u044B\u0431\u0440\u0430\u043D\u043D\u044B\u043C \u0442\u043E\u0440\u0440\u0435\u043D\u0442\u0430\u043C",
      set: "\u0423\u0441\u0442\u0430\u043D\u043E\u0432\u0438\u0442\u044C \u0442\u0435\u0433\u0438",
      set_tip: "\u0417\u0430\u043C\u0435\u043D\u0438\u0442\u044C \u0442\u0435\u0433\u0438 \u0432\u044B\u0431\u0440\u0430\u043D\u043D\u044B\u0445 \u0442\u043E\u0440\u0440\u0435\u043D\u0442\u043E\u0432"
    },
    title: "\u041D\u0430\u0437\u0432\u0430\u043D\u0438\u0435",
    toggle_drawer: "\u041F\u0435\u0440\u0435\u043A\u043B\u044E\u0447\u0438\u0442\u044C \u043F\u0430\u043D\u0435\u043B\u044C",
    votes_count_n: "{{count}} \u0433\u043E\u043B\u043E\u0441\u043E\u0432"
  },
  version: {
    bitmagnet_version: "\u0412\u0435\u0440\u0441\u0438\u044F bitmagnet {{version}}",
    unknown: "\u043D\u0435\u0438\u0437\u0432\u0435\u0441\u0442\u043D\u043E"
  }
};

// src/app/i18n/translations/tr.json
var tr_default = {
  content_types: {
    plural: {
      all: "T\xFCm\xFC",
      audiobook: "Sesli Kitaplar",
      comic: "\xC7izgi Romanlar",
      ebook: "E-Kitaplar",
      game: "Oyunlar",
      movie: "Filmler",
      music: "M\xFCzik",
      null: "Bilinmiyor",
      software: "Yaz\u0131l\u0131m",
      tv_show: "TV Programlar\u0131",
      xxx: "XXX"
    },
    singular: {
      audiobook: "Sesli Kitap",
      comic: "\xC7izgi Roman",
      ebook: "E-Kitap",
      game: "Oyun",
      movie: "Film",
      music: "M\xFCzik",
      null: "Bilinmiyor",
      software: "Yaz\u0131l\u0131m",
      tv_show: "TV Program\u0131",
      xxx: "XXX"
    }
  },
  dashboard: {
    event: {
      created: "Olu\u015Fturuldu",
      failed: "Ba\u015Far\u0131s\u0131z",
      processed: "\u0130\u015Flendi",
      updated: "G\xFCncellendi"
    },
    interval: {
      all: "T\xFCm\xFC",
      days: "G\xFCnler",
      days_1: "1 g\xFCn",
      hours: "Saatler",
      hours_1: "1 saat",
      hours_12: "12 saat",
      hours_6: "6 saat",
      minutes: "Dakikalar",
      minutes_1: "1 dakika",
      minutes_15: "15 dakika",
      minutes_30: "30 dakika",
      minutes_5: "5 dakika",
      off: "Kapal\u0131",
      seconds_10: "10 saniye",
      seconds_30: "30 saniye",
      weeks_1: "1 hafta"
    },
    metrics: {
      event: "Olay",
      resolution: "\xC7\xF6z\xFCn\xFCrl\xFCk",
      throughput: "Verim",
      timeframe: "Zaman Dilimi",
      toggle_legend: "Efsaneyi De\u011Fi\u015Ftir"
    },
    queues: {
      created: "Olu\u015Fturuldu",
      created_at: "Olu\u015Fturulma tarihi",
      enqueue_jobs: "\u0130\u015Fleri Kuyru\u011Fa Al",
      enqueue_torrent_processing_batch: "Torrent \u0130\u015Fleme Paketini Kuyru\u011Fa Al",
      failed: "Ba\u015Far\u0131s\u0131z",
      force_rematch: "E\u015Fle\u015Fmi\u015F i\xE7eri\u011Fi yeniden e\u015Fle\u015Ftir",
      jobs_enqueued: "Kuyru\u011Fa Al\u0131nan \u0130\u015Fler",
      latency: "Gecikme",
      match_content_by_external_api_search: "D\u0131\u015F API aramas\u0131yla i\xE7eri\u011Fi e\u015Fle\u015Ftir",
      match_content_by_local_search: "Yerel aramayla i\xE7eri\u011Fi e\u015Fle\u015Ftir",
      payload: "Veri Y\xFCk\xFC",
      pending: "Beklemede",
      priority: "\xD6ncelik",
      process_orphaned_torrents_only: "Sadece sahipsiz torrentleri i\u015Fle",
      processed: "\u0130\u015Flendi",
      purge_jobs: "\u0130\u015Fleri Temizle",
      purge_queue_jobs: "Kuyruk \u0130\u015Flerini Temizle",
      queue: "Kuyruk",
      queue_purged: "Kuyruk Temizlendi",
      queues: "Kuyruklar",
      ran_at: "\xC7al\u0131\u015Ft\u0131r\u0131ld\u0131",
      retry: "Tekrar Dene",
      total_counts_by_status: "Duruma G\xF6re Toplam Say\u0131lar"
    }
  },
  facets: {
    content_type: "\u0130\xE7erik T\xFCr\xFC",
    file_type: "Dosya T\xFCr\xFC",
    genre: "T\xFCr",
    language: "Dil",
    queue: "Kuyruk",
    status: "Durum",
    torrent_source: "Torrent Kayna\u011F\u0131",
    torrent_tag: "Torrent Etiketi",
    video_resolution: "Video \xC7\xF6z\xFCn\xFCrl\xFC\u011F\xFC",
    video_source: "Video Kayna\u011F\u0131"
  },
  file_types: {
    archive: "Ar\u015Fiv",
    audio: "Ses",
    data: "Veri",
    document: "Belge",
    image: "G\xF6r\xFCnt\xFC",
    software: "Yaz\u0131l\u0131m",
    subtitles: "Altyaz\u0131lar",
    unknown: "Bilinmiyor",
    video: "Video"
  },
  general: {
    all: "T\xFCm\xFC",
    dismiss: "Kapat",
    error: "Hata",
    none: "Hi\xE7biri",
    page_not_found: "Sayfa Bulunamad\u0131",
    refresh: "Yenile",
    status: "Durum"
  },
  health: {
    bitmagnet_is_status: "bitmagnet durumu {{status}}",
    check_failed_with_error: "Hata ile kontrol ba\u015Far\u0131s\u0131z oldu",
    component: "Bile\u015Fen",
    components: {
      dht: "DHT",
      postgres: "Postgres",
      tmdb: "TMDB"
    },
    dismiss: "Kapat",
    error: "Hata",
    status: "Durum",
    statuses: {
      degraded: "Azalm\u0131\u015F",
      down: "Kapal\u0131",
      error: "Hata",
      inactive: "Etkin de\u011Fil",
      started: "Ba\u015Flad\u0131",
      unknown: "Bilinmiyor",
      up: "\xC7al\u0131\u015F\u0131yor"
    },
    summary: "Sa\u011Fl\u0131k \xD6zeti",
    worker: "\u0130\u015F\xE7i",
    workers: {
      dht_crawler: "DHT taray\u0131c\u0131",
      http_server: "HTTP sunucusu",
      queue_server: "Kuyruk sunucusu"
    }
  },
  languages: {
    af: "Afrikanca",
    ar: "Arap\xE7a",
    az: "Azerice",
    be: "Beyaz Rus\xE7a",
    bg: "Bulgarca",
    bs: "Bo\u015Fnak\xE7a",
    ca: "Katalanca",
    ce: "\xC7e\xE7ence",
    co: "Korsikaca",
    cs: "\xC7ek\xE7e",
    cy: "Galce",
    da: "Danca",
    de: "Almanca",
    el: "Yunanca",
    en: "\u0130ngilizce",
    es: "\u0130spanyolca",
    et: "Estonca",
    eu: "Bask\xE7a",
    fa: "Fars\xE7a",
    fi: "Fince",
    fr: "Frans\u0131zca",
    he: "\u0130branice",
    hi: "Hint\xE7e",
    hr: "H\u0131rvat\xE7a",
    hu: "Macarca",
    hy: "Ermenice",
    id: "Endonezce",
    is: "\u0130zlandaca",
    it: "\u0130talyanca",
    ja: "Japonca",
    ka: "G\xFCrc\xFCce",
    ko: "Korece",
    ku: "K\xFCrt\xE7e",
    lt: "Litvanca",
    lv: "Letonca",
    mi: "Maorice",
    mk: "Makedonca",
    ml: "Malayalamca",
    mn: "Mo\u011Folca",
    ms: "Malayca",
    mt: "Maltaca",
    nl: "Flemenk\xE7e",
    no: "Norve\xE7\xE7e",
    pl: "Leh\xE7e",
    pt: "Portekizce",
    ro: "Romence",
    ru: "Rus\xE7a",
    sa: "Sanskrit\xE7e",
    sk: "Slovak\xE7a",
    sl: "Slovence",
    sm: "Samoaca",
    so: "Somalice",
    sr: "S\u0131rp\xE7a",
    sv: "\u0130sve\xE7\xE7e",
    ta: "Tamilce",
    th: "Tayca",
    tr: "T\xFCrk\xE7e",
    uk: "Ukraynaca",
    vi: "Vietnamca",
    yi: "Yidi\u015F",
    zh: "\xC7ince",
    zu: "Zuluca"
  },
  layout: {
    bitmagnet_on_service: "{{service}} \xFCzerinde bitmagnet",
    change_theme: "Temay\u0131 De\u011Fi\u015Ftir",
    external_links: "D\u0131\u015F Ba\u011Flant\u0131lar",
    sponsor: "Sponsor",
    support_bitmagnet: "bitmagnet'i Destekle",
    translate: "\xC7evir"
  },
  paginator: {
    first_page: "\u0130lk Sayfa",
    items_per_page: "Sayfa ba\u015F\u0131na \xF6\u011Fe",
    last_page: "Son Sayfa",
    next_page: "Sonraki Sayfa",
    page_x: "{{x}}. Sayfa",
    previous_page: "\xD6nceki Sayfa",
    x_to_y: "{{x}} - {{y}} aras\u0131",
    x_to_y_of_z: "{{x}} ile {{y}} aras\u0131, toplam {{z}}"
  },
  routes: {
    admin: "Y\xF6netici",
    dashboard: "Kontrol Paneli",
    home: "Ana Sayfa",
    jobs: "G\xF6revler",
    queues: "Kuyruklar",
    torrents: "Torrentler",
    visualize: "G\xF6rselle\u015Ftir"
  },
  torrents: {
    classification: "S\u0131n\u0131fland\u0131rma",
    clear_search: "Aramay\u0131 Temizle",
    copy: "Kopyala",
    copy_to_clipboard: "Panoya Kopyala",
    delete: "Sil",
    delete_action_cannot_be_undone: "Bu i\u015Flem geri al\u0131namaz",
    delete_are_you_sure: "Bu torrent'i silmek istedi\u011Finizden emin misiniz?",
    deselect_all: "T\xFCm Se\xE7imleri Kald\u0131r",
    edit_tags: "Etiketleri D\xFCzenle",
    episodes: "B\xF6l\xFCmler",
    external_links: "D\u0131\u015F Ba\u011Flant\u0131lar",
    file_index: "Dosya Dizini",
    file_path: "Dosya Yolu",
    file_size: "Dosya Boyutu",
    file_type: "Dosya T\xFCr\xFC",
    files: "Dosyalar",
    files_count_n: "{{count}} dosya",
    files_no_info: "Dosya bilgisi mevcut de\u011Fil",
    files_single: "Tek dosya",
    genres: "T\xFCrler",
    info_hash: "Bilgi hash'i",
    info_hashes: "Bilgi hash'leri",
    languages: "Diller",
    leechers: "Leechers",
    magnet: "Magnet",
    magnet_links: "Magnet Ba\u011Flant\u0131lar",
    new_tag: "Yeni etiket",
    order_by: "S\u0131ralama \xF6l\xE7\xFCt\xFC",
    order_direction_toggle: "Y\xF6n\xFC de\u011Fi\u015Ftir",
    ordering: {
      files_count: "Dosya say\u0131s\u0131",
      info_hash: "Bilgi hash'i",
      leechers: "Leechers",
      name: "Ad",
      published_at: "Yay\u0131nlanma tarihi",
      relevance: "Alaka d\xFCzeyi",
      seeders: "Seeders",
      size: "Boyut",
      updated_at: "G\xFCncellenme tarihi"
    },
    original_release_date: "Orijinal \xE7\u0131k\u0131\u015F tarihi",
    permalink: "Kal\u0131c\u0131 Ba\u011Flant\u0131",
    poster: "Poster",
    published: "Yay\u0131nland\u0131",
    rating: "Puan",
    refresh: "Sonu\xE7lar\u0131 Yenile",
    reprocess: {
      force_rematch: "Zaten e\u015Fle\u015Fen i\xE7eri\u011Fi yeniden e\u015Fle\u015Ftir",
      match_content_by_external_api_search: "D\u0131\u015F API aramas\u0131yla i\xE7eri\u011Fi e\u015Fle\u015Ftir",
      match_content_by_local_search: "Yerel aramayla i\xE7eri\u011Fi e\u015Fle\u015Ftir",
      reprocess: "Yeniden i\u015Fle"
    },
    s_l: "S / L",
    search: "Ara",
    seeders: "Seeders",
    select_all: "T\xFCm\xFCn\xFC Se\xE7",
    showing_x_of_y_files: "{{y}} dosyan\u0131n {{x}} tanesi g\xF6steriliyor",
    size: "Boyut",
    source: "Torrent Kayna\u011F\u0131",
    summary: "\xD6zet",
    tags: {
      delete: "Etiketleri sil",
      delete_tip: "Se\xE7ili torrentlerden etiketleri kald\u0131r",
      placeholder: "Etiket...",
      put: "Etiket koy",
      put_tip: "Se\xE7ili torrentlere etiket ekle",
      set: "Etiketleri ayarla",
      set_tip: "Se\xE7ili torrentlerin etiketlerini de\u011Fi\u015Ftir"
    },
    title: "Ba\u015Fl\u0131k",
    toggle_drawer: "\xC7ekmeceyi A\xE7/Kapat",
    votes_count_n: "{{count}} oy"
  },
  version: {
    bitmagnet_version: "bitmagnet versiyonu {{version}}",
    unknown: "bilinmiyor"
  }
};

// src/app/i18n/translations/uk.json
var uk_default = {
  content_types: {
    plural: {
      all: "\u0423\u0441\u0456",
      audiobook: "\u0410\u0443\u0434\u0456\u043E\u043A\u043D\u0438\u0433\u0438",
      comic: "\u041A\u043E\u043C\u0456\u043A\u0441\u0438",
      ebook: "\u0415\u043B\u0435\u043A\u0442\u0440\u043E\u043D\u043D\u0456 \u043A\u043D\u0438\u0433\u0438",
      game: "\u0406\u0433\u0440\u0438",
      movie: "\u0424\u0456\u043B\u044C\u043C\u0438",
      music: "\u041C\u0443\u0437\u0438\u043A\u0430",
      null: "\u041D\u0435\u0432\u0456\u0434\u043E\u043C\u043E",
      software: "\u041F\u0440\u043E\u0433\u0440\u0430\u043C\u0438",
      tv_show: "\u0422\u0435\u043B\u0435\u0448\u043E\u0443",
      xxx: "\u041F\u043E\u0440\u043D\u043E"
    },
    singular: {
      audiobook: "\u0410\u0443\u0434\u0456\u043E\u043A\u043D\u0438\u0433\u0430",
      comic: "\u041A\u043E\u043C\u0456\u043A\u0441",
      ebook: "\u0415\u043B\u0435\u043A\u0442\u0440\u043E\u043D\u043D\u0430 \u043A\u043D\u0438\u0433\u0430",
      game: "\u0413\u0440\u0430",
      movie: "\u0424\u0456\u043B\u044C\u043C",
      music: "\u041C\u0443\u0437\u0438\u043A\u0430",
      null: "\u041D\u0435\u0432\u0456\u0434\u043E\u043C\u043E",
      software: "\u041F\u0440\u043E\u0433\u0440\u0430\u043C\u0430",
      tv_show: "\u0422\u0435\u043B\u0435\u0448\u043E\u0443",
      xxx: "\u041F\u043E\u0440\u043D\u043E"
    }
  },
  dashboard: {
    event: {
      created: "\u0421\u0442\u0432\u043E\u0440\u0435\u043D\u043E",
      failed: "\u041F\u043E\u043C\u0438\u043B\u043A\u0430",
      processed: "\u041E\u0431\u0440\u043E\u0431\u043B\u0435\u043D\u043E",
      updated: "\u041E\u043D\u043E\u0432\u043B\u0435\u043D\u043E"
    },
    interval: {
      all: "\u0423\u0441\u0456",
      days: "\u0414\u043D\u0456",
      days_1: "1 \u0434\u0435\u043D\u044C",
      hours: "\u0413\u043E\u0434\u0438\u043D\u0438",
      hours_1: "1 \u0433\u043E\u0434\u0438\u043D\u0430",
      hours_12: "12 \u0433\u043E\u0434\u0438\u043D",
      hours_6: "6 \u0433\u043E\u0434\u0438\u043D",
      minutes: "\u0425\u0432\u0438\u043B\u0438\u043D\u0438",
      minutes_1: "1 \u0445\u0432\u0438\u043B\u0438\u043D\u0430",
      minutes_15: "15 \u0445\u0432\u0438\u043B\u0438\u043D",
      minutes_30: "30 \u0445\u0432\u0438\u043B\u0438\u043D",
      minutes_5: "5 \u0445\u0432\u0438\u043B\u0438\u043D",
      off: "\u0412\u0438\u043C\u043A\u043D\u0435\u043D\u043E",
      seconds_10: "10 \u0441\u0435\u043A\u0443\u043D\u0434",
      seconds_30: "30 \u0441\u0435\u043A\u0443\u043D\u0434",
      weeks_1: "1 \u0442\u0438\u0436\u0434\u0435\u043D\u044C"
    },
    metrics: {
      event: "\u041F\u043E\u0434\u0456\u044F",
      resolution: "\u0420\u043E\u0437\u0434\u0456\u043B\u044C\u043D\u0430 \u0437\u0434\u0430\u0442\u043D\u0456\u0441\u0442\u044C",
      throughput: "\u041F\u0440\u043E\u043F\u0443\u0441\u043A\u043D\u0430 \u0437\u0434\u0430\u0442\u043D\u0456\u0441\u0442\u044C",
      timeframe: "\u041F\u0440\u043E\u043C\u0456\u0436\u043E\u043A \u0447\u0430\u0441\u0443",
      toggle_legend: "\u041F\u0435\u0440\u0435\u043C\u043A\u043D\u0443\u0442\u0438 \u043B\u0435\u0433\u0435\u043D\u0434\u0443"
    },
    queues: {
      created: "\u0421\u0442\u0432\u043E\u0440\u0435\u043D\u043E",
      created_at: "\u0421\u0442\u0432\u043E\u0440\u0435\u043D\u043E \u043E",
      enqueue_jobs: "\u0414\u043E\u0434\u0430\u0442\u0438 \u0437\u0430\u0432\u0434\u0430\u043D\u043D\u044F \u0434\u043E \u0447\u0435\u0440\u0433\u0438",
      enqueue_torrent_processing_batch: "\u0414\u043E\u0434\u0430\u0442\u0438 \u043F\u0430\u043A\u0435\u0442 \u0434\u043B\u044F \u043E\u0431\u0440\u043E\u0431\u043A\u0438 \u0442\u043E\u0440\u0440\u0435\u043D\u0442\u0456\u0432",
      failed: "\u041F\u043E\u043C\u0438\u043B\u043A\u0430",
      force_rematch: "\u041F\u0440\u0438\u043C\u0443\u0441\u043E\u0432\u043E \u043F\u043E\u0432\u0442\u043E\u0440\u043D\u043E \u0437\u0456\u0441\u0442\u0430\u0432\u0438\u0442\u0438 \u0432\u0436\u0435 \u0437\u0456\u0441\u0442\u0430\u0432\u043B\u0435\u043D\u0438\u0439 \u043A\u043E\u043D\u0442\u0435\u043D\u0442",
      jobs_enqueued: "\u0417\u0430\u0432\u0434\u0430\u043D\u043D\u044F \u0434\u043E\u0434\u0430\u043D\u043E \u0432 \u0447\u0435\u0440\u0433\u0443",
      latency: "\u0417\u0430\u0442\u0440\u0438\u043C\u043A\u0430",
      match_content_by_external_api_search: "\u0417\u0456\u0441\u0442\u0430\u0432\u0438\u0442\u0438 \u043A\u043E\u043D\u0442\u0435\u043D\u0442 \u0447\u0435\u0440\u0435\u0437 \u0437\u043E\u0432\u043D\u0456\u0448\u043D\u0456\u0439 API",
      match_content_by_local_search: "\u0417\u0456\u0441\u0442\u0430\u0432\u0438\u0442\u0438 \u043A\u043E\u043D\u0442\u0435\u043D\u0442 \u0447\u0435\u0440\u0435\u0437 \u043B\u043E\u043A\u0430\u043B\u044C\u043D\u0438\u0439 \u043F\u043E\u0448\u0443\u043A",
      payload: "\u0414\u0430\u043D\u0456",
      pending: "\u0412 \u043E\u0447\u0456\u043A\u0443\u0432\u0430\u043D\u043D\u0456",
      priority: "\u041F\u0440\u0456\u043E\u0440\u0438\u0442\u0435\u0442",
      process_orphaned_torrents_only: "\u041E\u0431\u0440\u043E\u0431\u043B\u044F\u0442\u0438 \u043B\u0438\u0448\u0435 \xAB\u0441\u0438\u0440\u043E\u0442\u043B\u0438\u0432\u0456\xBB \u0442\u043E\u0440\u0440\u0435\u043D\u0442\u0438",
      processed: "\u041E\u0431\u0440\u043E\u0431\u043B\u0435\u043D\u043E",
      purge_jobs: "\u041E\u0447\u0438\u0441\u0442\u0438\u0442\u0438 \u0437\u0430\u0432\u0434\u0430\u043D\u043D\u044F",
      purge_queue_jobs: "\u041E\u0447\u0438\u0441\u0442\u0438\u0442\u0438 \u0437\u0430\u0432\u0434\u0430\u043D\u043D\u044F \u0447\u0435\u0440\u0433\u0438",
      queue: "\u0427\u0435\u0440\u0433\u0430",
      queue_purged: "\u0427\u0435\u0440\u0433\u0430 \u043E\u0447\u0438\u0449\u0435\u043D\u0430",
      queues: "\u0427\u0435\u0440\u0433\u0438",
      ran_at: "\u0417\u0430\u043F\u0443\u0449\u0435\u043D\u043E \u043E",
      retry: "\u041F\u043E\u0432\u0442\u043E\u0440\u0438\u0442\u0438",
      total_counts_by_status: "\u0417\u0430\u0433\u0430\u043B\u044C\u043D\u0430 \u043A\u0456\u043B\u044C\u043A\u0456\u0441\u0442\u044C \u0437\u0430 \u0441\u0442\u0430\u0442\u0443\u0441\u0430\u043C\u0438"
    }
  },
  facets: {
    content_type: "\u0422\u0438\u043F \u043A\u043E\u043D\u0442\u0435\u043D\u0442\u0443",
    file_type: "\u0422\u0438\u043F \u0444\u0430\u0439\u043B\u0443",
    genre: "\u0416\u0430\u043D\u0440",
    language: "\u041C\u043E\u0432\u0430",
    queue: "\u0427\u0435\u0440\u0433\u0430",
    status: "\u0421\u0442\u0430\u0442\u0443\u0441",
    torrent_source: "\u0414\u0436\u0435\u0440\u0435\u043B\u043E \u0442\u043E\u0440\u0440\u0435\u043D\u0442\u0430",
    torrent_tag: "\u0422\u0435\u0433 \u0442\u043E\u0440\u0440\u0435\u043D\u0442\u0430",
    video_resolution: "\u0420\u043E\u0437\u0434\u0456\u043B\u044C\u043D\u0430 \u0437\u0434\u0430\u0442\u043D\u0456\u0441\u0442\u044C \u0432\u0456\u0434\u0435\u043E",
    video_source: "\u0414\u0436\u0435\u0440\u0435\u043B\u043E \u0432\u0456\u0434\u0435\u043E"
  },
  file_types: {
    archive: "\u0410\u0440\u0445\u0456\u0432",
    audio: "\u0410\u0443\u0434\u0456\u043E",
    data: "\u0414\u0430\u043D\u0456",
    document: "\u0414\u043E\u043A\u0443\u043C\u0435\u043D\u0442",
    image: "\u0417\u043E\u0431\u0440\u0430\u0436\u0435\u043D\u043D\u044F",
    software: "\u041F\u0440\u043E\u0433\u0440\u0430\u043C\u0438",
    subtitles: "\u0421\u0443\u0431\u0442\u0438\u0442\u0440\u0438",
    unknown: "\u041D\u0435\u0432\u0456\u0434\u043E\u043C\u043E",
    video: "\u0412\u0456\u0434\u0435\u043E"
  },
  general: {
    all: "\u0423\u0441\u0456",
    dismiss: "\u0417\u0430\u043A\u0440\u0438\u0442\u0438",
    error: "\u041F\u043E\u043C\u0438\u043B\u043A\u0430",
    none: "\u041D\u0435\u043C\u0430\u0454",
    page_not_found: "\u0421\u0442\u043E\u0440\u0456\u043D\u043A\u0430 \u043D\u0435 \u0437\u043D\u0430\u0439\u0434\u0435\u043D\u0430",
    refresh: "\u041E\u043D\u043E\u0432\u0438\u0442\u0438",
    status: "\u0421\u0442\u0430\u0442\u0443\u0441"
  },
  health: {
    bitmagnet_is_status: "bitmagnet {{status}}",
    check_failed_with_error: "\u041F\u0435\u0440\u0435\u0432\u0456\u0440\u043A\u0430 \u0437\u0430\u0432\u0435\u0440\u0448\u0438\u043B\u0430\u0441\u044F \u043F\u043E\u043C\u0438\u043B\u043A\u043E\u044E",
    component: "\u041A\u043E\u043C\u043F\u043E\u043D\u0435\u043D\u0442",
    components: {
      dht: "DHT",
      postgres: "Postgres",
      tmdb: "TMDB"
    },
    dismiss: "\u0417\u0430\u043A\u0440\u0438\u0442\u0438",
    error: "\u041F\u043E\u043C\u0438\u043B\u043A\u0430",
    status: "\u0421\u0442\u0430\u0442\u0443\u0441",
    statuses: {
      degraded: "\u0417\u043D\u0438\u0436\u0435\u043D\u0430 \u043F\u0440\u043E\u0434\u0443\u043A\u0442\u0438\u0432\u043D\u0456\u0441\u0442\u044C",
      down: "\u041D\u0435 \u043F\u0440\u0430\u0446\u044E\u0454",
      error: "\u041F\u043E\u043C\u0438\u043B\u043A\u0430",
      inactive: "\u041D\u0435\u0430\u043A\u0442\u0438\u0432\u043D\u0438\u0439",
      started: "\u0417\u0430\u043F\u0443\u0449\u0435\u043D\u043E",
      unknown: "\u041E\u0447\u0456\u043A\u0443\u0454\u0442\u044C\u0441\u044F",
      up: "\u041F\u0440\u0430\u0446\u044E\u0454"
    },
    summary: "\u0417\u0432\u0435\u0434\u0435\u043D\u043D\u044F \u0441\u0442\u0430\u043D\u0443",
    worker: "\u041F\u0440\u0430\u0446\u0456\u0432\u043D\u0438\u043A",
    workers: {
      dht_crawler: "DHT \u0441\u043A\u0430\u043D\u0435\u0440",
      http_server: "HTTP \u0441\u0435\u0440\u0432\u0435\u0440",
      queue_server: "\u0421\u0435\u0440\u0432\u0435\u0440 \u0447\u0435\u0440\u0433\u0438"
    }
  },
  languages: {
    af: "\u0410\u0444\u0440\u0438\u043A\u0430\u0430\u043D\u0441",
    ar: "\u0410\u0440\u0430\u0431\u0441\u044C\u043A\u0430",
    az: "\u0410\u0437\u0435\u0440\u0431\u0430\u0439\u0434\u0436\u0430\u043D\u0441\u044C\u043A\u0430",
    be: "\u0411\u0456\u043B\u043E\u0440\u0443\u0441\u044C\u043A\u0430",
    bg: "\u0411\u043E\u043B\u0433\u0430\u0440\u0441\u044C\u043A\u0430",
    bs: "\u0411\u043E\u0441\u043D\u0456\u0439\u0441\u044C\u043A\u0430",
    ca: "\u041A\u0430\u0442\u0430\u043B\u043E\u043D\u0441\u044C\u043A\u0430",
    ce: "\u0427\u0435\u0447\u0435\u043D\u0441\u044C\u043A\u0430",
    co: "\u041A\u043E\u0440\u0441\u0438\u043A\u0430\u043D\u0441\u044C\u043A\u0430",
    cs: "\u0427\u0435\u0441\u044C\u043A\u0430",
    cy: "\u0412\u0430\u043B\u043B\u0456\u0439\u0441\u044C\u043A\u0430",
    da: "\u0414\u0430\u043D\u0441\u044C\u043A\u0430",
    de: "\u041D\u0456\u043C\u0435\u0446\u044C\u043A\u0430",
    el: "\u0413\u0440\u0435\u0446\u044C\u043A\u0430",
    en: "\u0410\u043D\u0433\u043B\u0456\u0439\u0441\u044C\u043A\u0430",
    es: "\u0406\u0441\u043F\u0430\u043D\u0441\u044C\u043A\u0430",
    et: "\u0415\u0441\u0442\u043E\u043D\u0441\u044C\u043A\u0430",
    eu: "\u0411\u0430\u0441\u043A\u0441\u044C\u043A\u0430",
    fa: "\u041F\u0435\u0440\u0441\u044C\u043A\u0430",
    fi: "\u0424\u0456\u043D\u0441\u044C\u043A\u0430",
    fr: "\u0424\u0440\u0430\u043D\u0446\u0443\u0437\u044C\u043A\u0430",
    he: "\u0406\u0432\u0440\u0438\u0442",
    hi: "\u0425\u0456\u043D\u0434\u0456",
    hr: "\u0425\u043E\u0440\u0432\u0430\u0442\u0441\u044C\u043A\u0430",
    hu: "\u0423\u0433\u043E\u0440\u0441\u044C\u043A\u0430",
    hy: "\u0412\u0456\u0440\u043C\u0435\u043D\u0441\u044C\u043A\u0430",
    id: "\u0406\u043D\u0434\u043E\u043D\u0435\u0437\u0456\u0439\u0441\u044C\u043A\u0430",
    is: "\u0406\u0441\u043B\u0430\u043D\u0434\u0441\u044C\u043A\u0430",
    it: "\u0406\u0442\u0430\u043B\u0456\u0439\u0441\u044C\u043A\u0430",
    ja: "\u042F\u043F\u043E\u043D\u0441\u044C\u043A\u0430",
    ka: "\u0413\u0440\u0443\u0437\u0438\u043D\u0441\u044C\u043A\u0430",
    ko: "\u041A\u043E\u0440\u0435\u0439\u0441\u044C\u043A\u0430",
    ku: "\u041A\u0443\u0440\u0434\u0441\u044C\u043A\u0430",
    lt: "\u041B\u0438\u0442\u043E\u0432\u0441\u044C\u043A\u0430",
    lv: "\u041B\u0430\u0442\u0438\u0441\u044C\u043A\u0430",
    mi: "\u041C\u0430\u043E\u0440\u0456",
    mk: "\u041C\u0430\u043A\u0435\u0434\u043E\u043D\u0441\u044C\u043A\u0430",
    ml: "\u041C\u0430\u043B\u0430\u044F\u043B\u0430\u043C",
    mn: "\u041C\u043E\u043D\u0433\u043E\u043B\u044C\u0441\u044C\u043A\u0430",
    ms: "\u041C\u0430\u043B\u0430\u0439\u0441\u044C\u043A\u0430",
    mt: "\u041C\u0430\u043B\u044C\u0442\u0456\u0439\u0441\u044C\u043A\u0430",
    nl: "\u041D\u0456\u0434\u0435\u0440\u043B\u0430\u043D\u0434\u0441\u044C\u043A\u0430",
    no: "\u041D\u043E\u0440\u0432\u0435\u0437\u044C\u043A\u0430",
    pl: "\u041F\u043E\u043B\u044C\u0441\u044C\u043A\u0430",
    pt: "\u041F\u043E\u0440\u0442\u0443\u0433\u0430\u043B\u044C\u0441\u044C\u043A\u0430",
    ro: "\u0420\u0443\u043C\u0443\u043D\u0441\u044C\u043A\u0430",
    ru: "\u0420\u043E\u0441\u0456\u0439\u0441\u044C\u043A\u0430",
    sa: "\u0421\u0430\u043D\u0441\u043A\u0440\u0438\u0442",
    sk: "\u0421\u043B\u043E\u0432\u0430\u0446\u044C\u043A\u0430",
    sl: "\u0421\u043B\u043E\u0432\u0435\u043D\u0441\u044C\u043A\u0430",
    sm: "\u0421\u0430\u043C\u043E\u0430\u043D\u0441\u044C\u043A\u0430",
    so: "\u0421\u043E\u043C\u0430\u043B\u0456\u0439\u0441\u044C\u043A\u0430",
    sr: "\u0421\u0435\u0440\u0431\u0441\u044C\u043A\u0430",
    sv: "\u0428\u0432\u0435\u0434\u0441\u044C\u043A\u0430",
    ta: "\u0422\u0430\u043C\u0456\u043B\u044C\u0441\u044C\u043A\u0430",
    th: "\u0422\u0430\u0439\u0441\u044C\u043A\u0430",
    tr: "\u0422\u0443\u0440\u0435\u0446\u044C\u043A\u0430",
    uk: "\u0423\u043A\u0440\u0430\u0457\u043D\u0441\u044C\u043A\u0430",
    vi: "\u0412'\u0454\u0442\u043D\u0430\u043C\u0441\u044C\u043A\u0430",
    yi: "\u0407\u0434\u0438\u0448",
    zh: "\u041A\u0438\u0442\u0430\u0439\u0441\u044C\u043A\u0430",
    zu: "\u0417\u0443\u043B\u0443\u0441\u044C\u043A\u0430"
  },
  layout: {
    bitmagnet_on_service: "bitmagnet \u043D\u0430 {{service}}",
    change_theme: "\u0417\u043C\u0456\u043D\u0438\u0442\u0438 \u0442\u0435\u043C\u0443",
    external_links: "\u0417\u043E\u0432\u043D\u0456\u0448\u043D\u0456 \u043F\u043E\u0441\u0438\u043B\u0430\u043D\u043D\u044F",
    sponsor: "\u0421\u043F\u043E\u043D\u0441\u043E\u0440",
    support_bitmagnet: "\u041F\u0456\u0434\u0442\u0440\u0438\u043C\u0430\u0442\u0438 bitmagnet",
    translate: "\u041F\u0435\u0440\u0435\u043A\u043B\u0430\u0441\u0442\u0438"
  },
  paginator: {
    first_page: "\u041F\u0435\u0440\u0448\u0430 \u0441\u0442\u043E\u0440\u0456\u043D\u043A\u0430",
    items_per_page: "\u0415\u043B\u0435\u043C\u0435\u043D\u0442\u0456\u0432 \u043D\u0430 \u0441\u0442\u043E\u0440\u0456\u043D\u0446\u0456",
    last_page: "\u041E\u0441\u0442\u0430\u043D\u043D\u044F \u0441\u0442\u043E\u0440\u0456\u043D\u043A\u0430",
    next_page: "\u041D\u0430\u0441\u0442\u0443\u043F\u043D\u0430 \u0441\u0442\u043E\u0440\u0456\u043D\u043A\u0430",
    page_x: "\u0421\u0442\u043E\u0440\u0456\u043D\u043A\u0430 {{x}}",
    previous_page: "\u041F\u043E\u043F\u0435\u0440\u0435\u0434\u043D\u044F \u0441\u0442\u043E\u0440\u0456\u043D\u043A\u0430",
    x_to_y: "{{x}} \u0434\u043E {{y}}",
    x_to_y_of_z: "{{x}} \u0434\u043E {{y}} \u0437 {{z}}"
  },
  routes: {
    admin: "\u0410\u0434\u043C\u0456\u043D\u0456\u0441\u0442\u0440\u0430\u0442\u043E\u0440",
    dashboard: "\u041F\u0430\u043D\u0435\u043B\u044C \u0443\u043F\u0440\u0430\u0432\u043B\u0456\u043D\u043D\u044F",
    home: "\u0413\u043E\u043B\u043E\u0432\u043D\u0430",
    jobs: "\u0417\u0430\u0432\u0434\u0430\u043D\u043D\u044F",
    queues: "\u0427\u0435\u0440\u0433\u0438",
    torrents: "\u0422\u043E\u0440\u0440\u0435\u043D\u0442\u0438",
    visualize: "\u0412\u0456\u0437\u0443\u0430\u043B\u0456\u0437\u0443\u0432\u0430\u0442\u0438"
  },
  torrents: {
    classification: "\u041A\u043B\u0430\u0441\u0438\u0444\u0456\u043A\u0430\u0446\u0456\u044F",
    clear_search: "\u041E\u0447\u0438\u0441\u0442\u0438\u0442\u0438 \u043F\u043E\u0448\u0443\u043A",
    copy: "\u041A\u043E\u043F\u0456\u044E\u0432\u0430\u0442\u0438",
    copy_to_clipboard: "\u041A\u043E\u043F\u0456\u044E\u0432\u0430\u0442\u0438 \u0432 \u0431\u0443\u0444\u0435\u0440 \u043E\u0431\u043C\u0456\u043D\u0443",
    delete: "\u0412\u0438\u0434\u0430\u043B\u0438\u0442\u0438",
    delete_action_cannot_be_undone: "\u0426\u044E \u0434\u0456\u044E \u043D\u0435 \u043C\u043E\u0436\u043D\u0430 \u0441\u043A\u0430\u0441\u0443\u0432\u0430\u0442\u0438",
    delete_are_you_sure: "\u0412\u0438 \u0432\u043F\u0435\u0432\u043D\u0435\u043D\u0456, \u0449\u043E \u0445\u043E\u0447\u0435\u0442\u0435 \u0432\u0438\u0434\u0430\u043B\u0438\u0442\u0438 \u0446\u0435\u0439 \u0442\u043E\u0440\u0440\u0435\u043D\u0442?",
    deselect_all: "\u0417\u043D\u044F\u0442\u0438 \u0432\u0438\u0434\u0456\u043B\u0435\u043D\u043D\u044F",
    edit_tags: "\u0420\u0435\u0434\u0430\u0433\u0443\u0432\u0430\u0442\u0438 \u0442\u0435\u0433\u0438",
    episodes: "\u0415\u043F\u0456\u0437\u043E\u0434\u0438",
    external_links: "\u0417\u043E\u0432\u043D\u0456\u0448\u043D\u0456 \u043F\u043E\u0441\u0438\u043B\u0430\u043D\u043D\u044F",
    file_index: "\u0406\u043D\u0434\u0435\u043A\u0441 \u0444\u0430\u0439\u043B\u0443",
    file_path: "\u0428\u043B\u044F\u0445 \u0434\u043E \u0444\u0430\u0439\u043B\u0443",
    file_size: "\u0420\u043E\u0437\u043C\u0456\u0440 \u0444\u0430\u0439\u043B\u0443",
    file_type: "\u0422\u0438\u043F \u0444\u0430\u0439\u043B\u0443",
    files: "\u0424\u0430\u0439\u043B\u0438",
    files_count_n: "{{count}} \u0444\u0430\u0439\u043B\u0456\u0432",
    files_no_info: "\u0406\u043D\u0444\u043E\u0440\u043C\u0430\u0446\u0456\u044F \u043F\u0440\u043E \u0444\u0430\u0439\u043B\u0438 \u043D\u0435\u0434\u043E\u0441\u0442\u0443\u043F\u043D\u0430",
    files_single: "\u041E\u0434\u0438\u043D \u0444\u0430\u0439\u043B",
    genres: "\u0416\u0430\u043D\u0440\u0438",
    info_hash: "\u0425\u0435\u0448 \u0456\u043D\u0444\u043E\u0440\u043C\u0430\u0446\u0456\u0457",
    info_hashes: "\u0425\u0435\u0448\u0456 \u0456\u043D\u0444\u043E\u0440\u043C\u0430\u0446\u0456\u0457",
    languages: "\u041C\u043E\u0432\u0438",
    leechers: "\u041B\u0456\u0447\u0435\u0440\u0438",
    magnet: "\u041C\u0430\u0433\u043D\u0435\u0442",
    magnet_links: "\u041C\u0430\u0433\u043D\u0435\u0442-\u043F\u043E\u0441\u0438\u043B\u0430\u043D\u043D\u044F",
    new_tag: "\u041D\u043E\u0432\u0438\u0439 \u0442\u0435\u0433",
    order_by: "\u0421\u043E\u0440\u0442\u0443\u0432\u0430\u0442\u0438 \u0437\u0430",
    order_direction_toggle: "\u0417\u043C\u0456\u043D\u0438\u0442\u0438 \u043D\u0430\u043F\u0440\u044F\u043C\u043E\u043A",
    ordering: {
      files_count: "\u041A\u0456\u043B\u044C\u043A\u0456\u0441\u0442\u044C \u0444\u0430\u0439\u043B\u0456\u0432",
      info_hash: "\u0425\u0435\u0448 \u0456\u043D\u0444\u043E\u0440\u043C\u0430\u0446\u0456\u0457",
      leechers: "\u041B\u0456\u0447\u0435\u0440\u0438",
      name: "\u041D\u0430\u0437\u0432\u0430",
      published_at: "\u0414\u0430\u0442\u0430 \u043F\u0443\u0431\u043B\u0456\u043A\u0430\u0446\u0456\u0457",
      relevance: "\u0410\u043A\u0442\u0443\u0430\u043B\u044C\u043D\u0456\u0441\u0442\u044C",
      seeders: "\u0421\u0456\u0434\u0435\u0440\u0438",
      size: "\u0420\u043E\u0437\u043C\u0456\u0440",
      updated_at: "\u0414\u0430\u0442\u0430 \u043E\u043D\u043E\u0432\u043B\u0435\u043D\u043D\u044F"
    },
    original_release_date: "\u0414\u0430\u0442\u0430 \u043E\u0440\u0438\u0433\u0456\u043D\u0430\u043B\u044C\u043D\u043E\u0433\u043E \u0432\u0438\u043F\u0443\u0441\u043A\u0443",
    permalink: "\u041F\u043E\u0441\u0442\u0456\u0439\u043D\u0435 \u043F\u043E\u0441\u0438\u043B\u0430\u043D\u043D\u044F",
    poster: "\u041F\u043E\u0441\u0442\u0435\u0440",
    published: "\u041E\u043F\u0443\u0431\u043B\u0456\u043A\u043E\u0432\u0430\u043D\u043E",
    rating: "\u0420\u0435\u0439\u0442\u0438\u043D\u0433",
    refresh: "\u041E\u043D\u043E\u0432\u0438\u0442\u0438 \u0440\u0435\u0437\u0443\u043B\u044C\u0442\u0430\u0442\u0438",
    reprocess: {
      force_rematch: "\u041F\u0440\u0438\u043C\u0443\u0441\u043E\u0432\u043E \u043F\u043E\u0432\u0442\u043E\u0440\u043D\u043E \u0437\u0456\u0441\u0442\u0430\u0432\u0438\u0442\u0438 \u0432\u0436\u0435 \u0437\u0456\u0441\u0442\u0430\u0432\u043B\u0435\u043D\u0438\u0439 \u043A\u043E\u043D\u0442\u0435\u043D\u0442",
      match_content_by_external_api_search: "\u0417\u0456\u0441\u0442\u0430\u0432\u0438\u0442\u0438 \u043A\u043E\u043D\u0442\u0435\u043D\u0442 \u0447\u0435\u0440\u0435\u0437 \u0437\u043E\u0432\u043D\u0456\u0448\u043D\u0456\u0439 API",
      match_content_by_local_search: "\u0417\u0456\u0441\u0442\u0430\u0432\u0438\u0442\u0438 \u043A\u043E\u043D\u0442\u0435\u043D\u0442 \u0447\u0435\u0440\u0435\u0437 \u043B\u043E\u043A\u0430\u043B\u044C\u043D\u0438\u0439 \u043F\u043E\u0448\u0443\u043A",
      reprocess: "\u041F\u0435\u0440\u0435\u0440\u043E\u0431\u0438\u0442\u0438"
    },
    s_l: "S / L",
    search: "\u041F\u043E\u0448\u0443\u043A",
    seeders: "\u0421\u0456\u0434\u0435\u0440\u0438",
    select_all: "\u0412\u0438\u0431\u0440\u0430\u0442\u0438 \u0432\u0441\u0456",
    showing_x_of_y_files: "\u041F\u043E\u043A\u0430\u0437\u0430\u043D\u043E {{x}} \u0437 {{y}} \u0444\u0430\u0439\u043B\u0456\u0432",
    size: "\u0420\u043E\u0437\u043C\u0456\u0440",
    source: "\u0414\u0436\u0435\u0440\u0435\u043B\u043E \u0442\u043E\u0440\u0440\u0435\u043D\u0442\u0430",
    summary: "\u0417\u0432\u0435\u0434\u0435\u043D\u043D\u044F",
    tags: {
      delete: "\u0412\u0438\u0434\u0430\u043B\u0438\u0442\u0438 \u0442\u0435\u0433\u0438",
      delete_tip: "\u0412\u0438\u0434\u0430\u043B\u0456\u0442\u044C \u0442\u0435\u0433\u0438 \u0437 \u0432\u0438\u0431\u0440\u0430\u043D\u0438\u0445 \u0442\u043E\u0440\u0440\u0435\u043D\u0442\u0456\u0432",
      placeholder: "\u0422\u0435\u0433\u0438...",
      put: "\u0421\u0442\u0430\u0432\u0442\u0435 \u0442\u0435\u0433\u0438",
      put_tip: "\u0414\u043E\u0434\u0430\u0439\u0442\u0435 \u0442\u0435\u0433\u0438 \u0434\u043E \u0432\u0438\u0431\u0440\u0430\u043D\u0438\u0445 \u0442\u043E\u0440\u0440\u0435\u043D\u0442\u0456\u0432",
      set: "\u0412\u0441\u0442\u0430\u043D\u043E\u0432\u0438\u0442\u0438 \u0442\u0435\u0433\u0438",
      set_tip: "\u0417\u0430\u043C\u0456\u043D\u0438\u0442\u0438 \u0442\u0435\u0433\u0438 \u0432\u0438\u0431\u0440\u0430\u043D\u0438\u0445 \u0442\u043E\u0440\u0440\u0435\u043D\u0442\u0456\u0432"
    },
    title: "\u041D\u0430\u0437\u0432\u0430",
    toggle_drawer: "\u041F\u0435\u0440\u0435\u043C\u043A\u043D\u0443\u0442\u0438 \u043F\u0430\u043D\u0435\u043B\u044C",
    votes_count_n: "{{count}} \u0433\u043E\u043B\u043E\u0441\u0456\u0432"
  },
  version: {
    bitmagnet_version: "\u0412\u0435\u0440\u0441\u0456\u044F bitmagnet {{version}}",
    unknown: "\u043D\u0435\u0432\u0456\u0434\u043E\u043C\u043E"
  }
};

// src/app/i18n/translations/zh.json
var zh_default = {
  content_types: {
    plural: {
      all: "\u5168\u90E8",
      audiobook: "\u6709\u58F0\u4E66",
      comic: "\u6F2B\u753B",
      ebook: "\u7535\u5B50\u4E66",
      game: "\u6E38\u620F",
      movie: "\u7535\u5F71",
      music: "\u97F3\u4E50",
      null: "\u672A\u77E5",
      software: "\u8F6F\u4EF6",
      tv_show: "\u7535\u89C6\u8282\u76EE",
      xxx: "\u8272\u60C5"
    },
    singular: {
      audiobook: "\u6709\u58F0\u4E66",
      comic: "\u6F2B\u753B",
      ebook: "\u7535\u5B50\u4E66",
      game: "\u6E38\u620F",
      movie: "\u7535\u5F71",
      music: "\u97F3\u4E50",
      null: "\u672A\u77E5",
      software: "\u8F6F\u4EF6",
      tv_show: "\u7535\u89C6\u8282\u76EE",
      xxx: "\u8272\u60C5"
    }
  },
  dashboard: {
    event: {
      created: "\u5DF2\u521B\u5EFA",
      failed: "\u5931\u8D25",
      processed: "\u5DF2\u5904\u7406",
      updated: "\u5DF2\u66F4\u65B0"
    },
    interval: {
      all: "\u5168\u90E8",
      days: "\u5929",
      days_1: "1\u5929",
      hours: "\u5C0F\u65F6",
      hours_1: "1\u5C0F\u65F6",
      hours_12: "12\u5C0F\u65F6",
      hours_6: "6\u5C0F\u65F6",
      minutes: "\u5206\u949F",
      minutes_1: "1\u5206\u949F",
      minutes_15: "15\u5206\u949F",
      minutes_30: "30\u5206\u949F",
      minutes_5: "5\u5206\u949F",
      off: "\u5173\u95ED",
      seconds_10: "10\u79D2",
      seconds_30: "30\u79D2",
      weeks_1: "1\u5468"
    },
    metrics: {
      event: "\u4E8B\u4EF6",
      resolution: "\u5206\u8FA8\u7387",
      throughput: "\u541E\u5410\u91CF",
      timeframe: "\u65F6\u95F4\u8303\u56F4",
      toggle_legend: "\u5207\u6362\u56FE\u4F8B"
    },
    queues: {
      created: "\u5DF2\u521B\u5EFA",
      created_at: "\u521B\u5EFA\u4E8E",
      enqueue_jobs: "\u52A0\u5165\u961F\u5217\u7684\u4EFB\u52A1",
      enqueue_torrent_processing_batch: "\u52A0\u5165\u961F\u5217\u7684\u79CD\u5B50\u5904\u7406\u6279\u6B21",
      failed: "\u5931\u8D25",
      force_rematch: "\u5F3A\u5236\u91CD\u65B0\u5339\u914D\u5DF2\u5339\u914D\u7684\u5185\u5BB9",
      jobs_enqueued: "\u52A0\u5165\u961F\u5217\u7684\u4EFB\u52A1",
      latency: "\u5EF6\u8FDF",
      match_content_by_external_api_search: "\u901A\u8FC7\u5916\u90E8API\u641C\u7D22\u5339\u914D\u5185\u5BB9",
      match_content_by_local_search: "\u901A\u8FC7\u672C\u5730\u641C\u7D22\u5339\u914D\u5185\u5BB9",
      payload: "\u6709\u6548\u8F7D\u8377",
      pending: "\u5F85\u5904\u7406",
      priority: "\u4F18\u5148\u7EA7",
      process_orphaned_torrents_only: "\u4EC5\u5904\u7406\u5B64\u7ACB\u7684\u79CD\u5B50",
      processed: "\u5DF2\u5904\u7406",
      purge_jobs: "\u6E05\u9664\u4EFB\u52A1",
      purge_queue_jobs: "\u6E05\u9664\u961F\u5217\u4EFB\u52A1",
      queue: "\u961F\u5217",
      queue_purged: "\u961F\u5217\u5DF2\u6E05\u9664",
      queues: "\u961F\u5217",
      ran_at: "\u8FD0\u884C\u4E8E",
      retry: "\u91CD\u8BD5",
      total_counts_by_status: "\u6309\u72B6\u6001\u7EDF\u8BA1\u603B\u6570"
    }
  },
  facets: {
    content_type: "\u5185\u5BB9\u7C7B\u578B",
    file_type: "\u6587\u4EF6\u7C7B\u578B",
    genre: "\u7C7B\u578B",
    language: "\u8BED\u8A00",
    queue: "\u961F\u5217",
    status: "\u72B6\u6001",
    torrent_source: "\u79CD\u5B50\u6765\u6E90",
    torrent_tag: "\u79CD\u5B50\u6807\u7B7E",
    video_resolution: "\u89C6\u9891\u5206\u8FA8\u7387",
    video_source: "\u89C6\u9891\u6765\u6E90"
  },
  file_types: {
    archive: "\u6863\u6848",
    audio: "\u97F3\u9891",
    data: "\u6570\u636E",
    document: "\u6587\u6863",
    image: "\u56FE\u50CF",
    software: "\u8F6F\u4EF6",
    subtitles: "\u5B57\u5E55",
    unknown: "\u672A\u77E5",
    video: "\u89C6\u9891"
  },
  general: {
    all: "\u5168\u90E8",
    dismiss: "\u5FFD\u7565",
    error: "\u9519\u8BEF",
    none: "\u65E0",
    page_not_found: "\u9875\u9762\u672A\u627E\u5230",
    refresh: "\u5237\u65B0",
    status: "\u72B6\u6001"
  },
  health: {
    bitmagnet_is_status: "bitmagnet\u662F{{status}}",
    check_failed_with_error: "\u68C0\u67E5\u5931\u8D25\uFF0C\u9519\u8BEF",
    component: "\u7EC4\u4EF6",
    components: {
      dht: "DHT",
      postgres: "Postgres",
      tmdb: "TMDB"
    },
    dismiss: "\u5FFD\u7565",
    error: "\u9519\u8BEF",
    status: "\u72B6\u6001",
    statuses: {
      degraded: "\u964D\u7EA7",
      down: "\u5173\u95ED",
      error: "\u9519\u8BEF",
      inactive: "\u4E0D\u6D3B\u8DC3",
      started: "\u5DF2\u542F\u52A8",
      unknown: "\u672A\u77E5",
      up: "\u8FD0\u884C\u4E2D"
    },
    summary: "\u5065\u5EB7\u6982\u8FF0",
    worker: "\u5DE5\u4F5C\u8005",
    workers: {
      dht_crawler: "DHT\u722C\u866B",
      http_server: "HTTP\u670D\u52A1\u5668",
      queue_server: "\u961F\u5217\u670D\u52A1\u5668"
    }
  },
  languages: {
    af: "\u5357\u975E\u8377\u5170\u8BED",
    ar: "\u963F\u62C9\u4F2F\u8BED",
    az: "\u963F\u585E\u62DC\u7586\u8BED",
    be: "\u767D\u4FC4\u7F57\u65AF\u8BED",
    bg: "\u4FDD\u52A0\u5229\u4E9A\u8BED",
    bs: "\u6CE2\u65AF\u5C3C\u4E9A\u8BED",
    ca: "\u52A0\u6CF0\u7F57\u5C3C\u4E9A\u8BED",
    ce: "\u8F66\u81E3\u8BED",
    co: "\u79D1\u897F\u5609\u8BED",
    cs: "\u6377\u514B\u8BED",
    cy: "\u5A01\u5C14\u58EB\u8BED",
    da: "\u4E39\u9EA6\u8BED",
    de: "\u5FB7\u8BED",
    el: "\u5E0C\u814A\u8BED",
    en: "\u82F1\u8BED",
    es: "\u897F\u73ED\u7259\u8BED",
    et: "\u7231\u6C99\u5C3C\u4E9A\u8BED",
    eu: "\u5DF4\u65AF\u514B\u8BED",
    fa: "\u6CE2\u65AF\u8BED",
    fi: "\u82AC\u5170\u8BED",
    fr: "\u6CD5\u8BED",
    he: "\u5E0C\u4F2F\u6765\u8BED",
    hi: "\u5370\u5730\u8BED",
    hr: "\u514B\u7F57\u5730\u4E9A\u8BED",
    hu: "\u5308\u7259\u5229\u8BED",
    hy: "\u4E9A\u7F8E\u5C3C\u4E9A\u8BED",
    id: "\u5370\u5EA6\u5C3C\u897F\u4E9A\u8BED",
    is: "\u51B0\u5C9B\u8BED",
    it: "\u610F\u5927\u5229\u8BED",
    ja: "\u65E5\u8BED",
    ka: "\u683C\u9C81\u5409\u4E9A\u8BED",
    ko: "\u97E9\u8BED",
    ku: "\u5E93\u5C14\u5FB7\u8BED",
    lt: "\u7ACB\u9676\u5B9B\u8BED",
    lv: "\u62C9\u8131\u7EF4\u4E9A\u8BED",
    mi: "\u6BDB\u5229\u8BED",
    mk: "\u9A6C\u5176\u987F\u8BED",
    ml: "\u9A6C\u62C9\u96C5\u62C9\u59C6\u8BED",
    mn: "\u8499\u53E4\u8BED",
    ms: "\u9A6C\u6765\u8BED",
    mt: "\u9A6C\u8033\u4ED6\u8BED",
    nl: "\u8377\u5170\u8BED",
    no: "\u632A\u5A01\u8BED",
    pl: "\u6CE2\u5170\u8BED",
    pt: "\u8461\u8404\u7259\u8BED",
    ro: "\u7F57\u9A6C\u5C3C\u4E9A\u8BED",
    ru: "\u4FC4\u8BED",
    sa: "\u68B5\u8BED",
    sk: "\u65AF\u6D1B\u4F10\u514B\u8BED",
    sl: "\u65AF\u6D1B\u6587\u5C3C\u4E9A\u8BED",
    sm: "\u8428\u6469\u4E9A\u8BED",
    so: "\u7D22\u9A6C\u91CC\u8BED",
    sr: "\u585E\u5C14\u7EF4\u4E9A\u8BED",
    sv: "\u745E\u5178\u8BED",
    ta: "\u6CF0\u7C73\u5C14\u8BED",
    th: "\u6CF0\u8BED",
    tr: "\u571F\u8033\u5176\u8BED",
    uk: "\u4E4C\u514B\u5170\u8BED",
    vi: "\u8D8A\u5357\u8BED",
    yi: "\u610F\u7B2C\u7EEA\u8BED",
    zh: "\u4E2D\u6587",
    zu: "\u7956\u9C81\u8BED"
  },
  layout: {
    bitmagnet_on_service: "bitmagnet\u5728{{service}}",
    change_theme: "\u66F4\u6539\u4E3B\u9898",
    external_links: "\u5916\u90E8\u94FE\u63A5",
    sponsor: "\u8D5E\u52A9\u5546",
    support_bitmagnet: "\u652F\u6301bitmagnet",
    translate: "\u7FFB\u8BD1"
  },
  paginator: {
    first_page: "\u7B2C\u4E00\u9875",
    items_per_page: "\u6BCF\u9875\u9879\u76EE\u6570",
    last_page: "\u6700\u540E\u4E00\u9875",
    next_page: "\u4E0B\u4E00\u9875",
    page_x: "\u7B2C{{x}}\u9875",
    previous_page: "\u4E0A\u4E00\u9875",
    x_to_y: "{{x}}\u5230{{y}}",
    x_to_y_of_z: "{{x}}\u5230{{y}}\uFF0C\u5171{{z}}"
  },
  routes: {
    admin: "\u7BA1\u7406\u5458",
    dashboard: "\u4EEA\u8868\u76D8",
    home: "\u9996\u9875",
    jobs: "\u4EFB\u52A1",
    queues: "\u961F\u5217",
    torrents: "\u79CD\u5B50",
    visualize: "\u53EF\u89C6\u5316"
  },
  torrents: {
    classification: "\u5206\u7C7B",
    clear_search: "\u6E05\u9664\u641C\u7D22",
    copy: "\u590D\u5236",
    copy_to_clipboard: "\u590D\u5236\u5230\u526A\u8D34\u677F",
    delete: "\u5220\u9664",
    delete_action_cannot_be_undone: "\u6B64\u64CD\u4F5C\u65E0\u6CD5\u64A4\u9500",
    delete_are_you_sure: "\u60A8\u786E\u5B9A\u8981\u5220\u9664\u6B64\u79CD\u5B50\u5417\uFF1F",
    deselect_all: "\u53D6\u6D88\u5168\u9009",
    edit_tags: "\u7F16\u8F91\u6807\u7B7E",
    episodes: "\u5267\u96C6",
    external_links: "\u5916\u90E8\u94FE\u63A5",
    file_index: "\u6587\u4EF6\u7D22\u5F15",
    file_path: "\u6587\u4EF6\u8DEF\u5F84",
    file_size: "\u6587\u4EF6\u5927\u5C0F",
    file_type: "\u6587\u4EF6\u7C7B\u578B",
    files: "\u6587\u4EF6",
    files_count_n: "{{count}}\u4E2A\u6587\u4EF6",
    files_no_info: "\u65E0\u6587\u4EF6\u4FE1\u606F",
    files_single: "\u5355\u4E2A\u6587\u4EF6",
    genres: "\u7C7B\u578B",
    info_hash: "\u4FE1\u606F\u54C8\u5E0C",
    info_hashes: "\u4FE1\u606F\u54C8\u5E0C",
    languages: "\u8BED\u8A00",
    leechers: "\u4E0B\u8F7D\u8005",
    magnet: "\u78C1\u529B\u94FE\u63A5",
    magnet_links: "\u78C1\u529B\u94FE\u63A5",
    new_tag: "\u65B0\u6807\u7B7E",
    order_by: "\u6392\u5E8F\u4F9D\u636E",
    order_direction_toggle: "\u5207\u6362\u65B9\u5411",
    ordering: {
      files_count: "\u6587\u4EF6\u6570",
      info_hash: "\u4FE1\u606F\u54C8\u5E0C",
      leechers: "\u4E0B\u8F7D\u8005",
      name: "\u540D\u79F0",
      published_at: "\u53D1\u5E03\u4E8E",
      relevance: "\u76F8\u5173\u6027",
      seeders: "\u4E0A\u4F20\u8005",
      size: "\u5927\u5C0F",
      updated_at: "\u66F4\u65B0\u4E8E"
    },
    original_release_date: "\u539F\u59CB\u53D1\u5E03\u65E5\u671F",
    permalink: "\u6C38\u4E45\u94FE\u63A5",
    poster: "\u6D77\u62A5",
    published: "\u5DF2\u53D1\u5E03",
    rating: "\u8BC4\u5206",
    refresh: "\u5237\u65B0\u7ED3\u679C",
    reprocess: {
      force_rematch: "\u5F3A\u5236\u91CD\u65B0\u5339\u914D\u5DF2\u5339\u914D\u7684\u5185\u5BB9",
      match_content_by_external_api_search: "\u901A\u8FC7\u5916\u90E8API\u641C\u7D22\u5339\u914D\u5185\u5BB9",
      match_content_by_local_search: "\u901A\u8FC7\u672C\u5730\u641C\u7D22\u5339\u914D\u5185\u5BB9",
      reprocess: "\u91CD\u65B0\u5904\u7406"
    },
    s_l: "S / L",
    search: "\u641C\u7D22",
    seeders: "\u4E0A\u4F20\u8005",
    select_all: "\u5168\u9009",
    showing_x_of_y_files: "\u663E\u793A{{x}}\u4E2A\uFF0C\u5171{{y}}\u4E2A\u6587\u4EF6",
    size: "\u5927\u5C0F",
    source: "\u79CD\u5B50\u6765\u6E90",
    summary: "\u6458\u8981",
    tags: {
      delete: "\u5220\u9664\u6807\u7B7E",
      delete_tip: "\u4ECE\u9009\u5B9A\u7684\u79CD\u5B50\u4E2D\u5220\u9664\u6807\u7B7E",
      placeholder: "\u6807\u7B7E...",
      put: "\u6DFB\u52A0\u6807\u7B7E",
      put_tip: "\u4E3A\u9009\u5B9A\u7684\u79CD\u5B50\u6DFB\u52A0\u6807\u7B7E",
      set: "\u8BBE\u7F6E\u6807\u7B7E",
      set_tip: "\u66FF\u6362\u9009\u5B9A\u79CD\u5B50\u7684\u6807\u7B7E"
    },
    title: "\u6807\u9898",
    toggle_drawer: "\u5207\u6362\u62BD\u5C49",
    votes_count_n: "{{count}}\u7968"
  },
  version: {
    bitmagnet_version: "bitmagnet\u7248\u672C{{version}}",
    unknown: "\u672A\u77E5"
  }
};

// src/app/i18n/translations.ts
var translations_default = {
  ar: ar_default,
  ca: ca_default,
  de: de_default,
  en: en_default,
  es: es_default,
  fr: fr_default,
  hi: hi_default,
  ja: ja_default,
  nl: nl_default,
  pt: pt_default,
  ru: ru_default,
  tr: tr_default,
  uk: uk_default,
  zh: zh_default
};

// src/app/i18n/transloco.loader.ts
var TranslocoImportLoader = class _TranslocoImportLoader {
  getTranslation(lang) {
    return __async(this, null, function* () {
      if (lang in translations_default) {
        const tr = translations_default[lang];
        return stripMissing(tr);
      } else {
        return Promise.reject(new Error(`Translation not found: ${lang}`));
      }
    });
  }
  static {
    this.\u0275fac = function TranslocoImportLoader_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _TranslocoImportLoader)();
    };
  }
  static {
    this.\u0275prov = /* @__PURE__ */ \u0275\u0275defineInjectable({ token: _TranslocoImportLoader, factory: _TranslocoImportLoader.\u0275fac, providedIn: "root" });
  }
};
var missingValues = ["__missing__", "__fallback__"];
var stripMissing = (tr) => Object.fromEntries(Object.entries(tr).flatMap(([k, v]) => {
  if (typeof v === "object") {
    v = stripMissing(v);
  } else if (typeof v === "string" && missingValues.includes(v)) {
    return [];
  }
  return [[k, v]];
}));

// src/app/app.routes.ts
var routes = [
  {
    path: "",
    pathMatch: "full",
    redirectTo: "torrents"
  },
  {
    path: "torrents",
    loadComponent: () => import("./chunk-7KDEQJJU.js").then((c) => c.TorrentsComponent),
    children: [
      {
        path: "",
        loadComponent: () => import("./chunk-EPENVDKZ.js").then((c) => c.TorrentsSearchComponent)
      },
      {
        path: "permalink/:infoHash",
        loadComponent: () => import("./chunk-HFGVMCKP.js").then((c) => c.TorrentPermalinkComponent)
      }
    ]
  },
  {
    path: "dashboard",
    loadComponent: () => import("./chunk-XLDUZ7AY.js").then((c) => c.DashboardComponent),
    children: [
      {
        path: "",
        loadComponent: () => import("./chunk-BI2TVXFB.js").then((c) => c.DashboardHomeComponent)
      },
      {
        path: "config",
        loadComponent: () => import("./chunk-CBWCYPCV.js").then((c) => c.DashboardConfigComponent)
      },
      {
        path: "workers",
        loadComponent: () => import("./chunk-6PZWTJGH.js").then((c) => c.DashboardWorkersComponent)
      },
      {
        path: "queues",
        pathMatch: "full",
        redirectTo: "queues/visualize"
      },
      {
        path: "queues",
        loadComponent: () => import("./chunk-2IVFA22B.js").then((c) => c.QueueDashboardComponent),
        children: [
          {
            path: "visualize",
            loadComponent: () => import("./chunk-JWHWW7YB.js").then((c) => c.QueueVisualizeComponent)
          },
          {
            path: "jobs",
            loadComponent: () => import("./chunk-ORC7SCA4.js").then((c) => c.QueueJobsComponent)
          },
          {
            path: "admin",
            loadComponent: () => import("./chunk-EYUWA2YW.js").then((c) => c.QueueAdminComponent)
          }
        ]
      },
      {
        path: "torrents",
        loadComponent: () => import("./chunk-YZYDLAWD.js").then((c) => c.TorrentsDashboardComponent)
      }
    ]
  },
  {
    path: "**",
    loadComponent: () => import("./chunk-5Q4VGCHF.js").then((c) => c.NotFoundComponent)
  }
];

// src/app/app.config.ts
var appConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(routes, withComponentInputBinding()),
    provideAnimationsAsync("animations"),
    provideHttpClient(withInterceptorsFromDi()),
    provideHttpClient(),
    provideApollo(() => {
      const httpLink = inject(HttpLink);
      return {
        link: httpLink.create({ uri: graphqlEndpoint }),
        cache: new InMemoryCache({
          typePolicies: {
            Query: {
              fields: {
                search: {
                  merge(existing, incoming) {
                    return __spreadValues(__spreadValues({}, existing), incoming);
                  }
                }
              }
            }
          }
        })
      };
    }),
    provideTransloco({
      config: {
        availableLangs: [
          {
            id: "ar",
            label: "\u0627\u0644\u0639\u0631\u0628\u064A\u0629"
          },
          {
            id: "ca",
            label: "Catal\xE0"
          },
          {
            id: "de",
            label: "Deutsch"
          },
          {
            id: "en",
            label: "English"
          },
          {
            id: "es",
            label: "Espa\xF1ol"
          },
          {
            id: "fr",
            label: "Fran\xE7ais"
          },
          {
            id: "hi",
            label: "\u0939\u093F\u0928\u094D\u0926\u0940"
          },
          {
            id: "ja",
            label: "\u65E5\u672C\u8A9E"
          },
          {
            id: "nl",
            label: "Nederlands"
          },
          {
            id: "pt",
            label: "Portugu\xEAs"
          },
          {
            id: "ru",
            label: "\u0420\u0443\u0441\u0441\u043A\u0438\u0439"
          },
          {
            id: "tr",
            label: "T\xFCrk\xE7e"
          },
          {
            id: "uk",
            label: "\u0423\u043A\u0440\u0430\u0457\u043D\u0441\u044C\u043A\u0430"
          },
          {
            id: "zh",
            label: "\u4E2D\u6587"
          }
        ],
        defaultLang: "en",
        fallbackLang: "en",
        missingHandler: {
          // It will use the first language set in the `fallbackLang` property
          useFallbackTranslation: true
        },
        // Remove this option if your application doesn't support changing language in runtime.
        reRenderOnLangChange: true,
        prodMode: false
      },
      loader: TranslocoImportLoader
    }),
    provideCharts(withDefaultRegisterables())
  ]
};

// src/app/browser-storage/browser-storage.service.ts
var BROWSER_STORAGE = new InjectionToken("Browser Storage", {
  providedIn: "root",
  factory: () => localStorage
});
var BrowserStorageService = class _BrowserStorageService {
  constructor(storage) {
    this.storage = storage;
  }
  get(key) {
    return this.storage.getItem(key);
  }
  set(key, value) {
    this.storage.setItem(key, value);
  }
  remove(key) {
    this.storage.removeItem(key);
  }
  clear() {
    this.storage.clear();
  }
  static {
    this.\u0275fac = function BrowserStorageService_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _BrowserStorageService)(\u0275\u0275inject(BROWSER_STORAGE));
    };
  }
  static {
    this.\u0275prov = /* @__PURE__ */ \u0275\u0275defineInjectable({ token: _BrowserStorageService, factory: _BrowserStorageService.\u0275fac, providedIn: "root" });
  }
};

// src/app/themes/theme-registry.ts
var _themes = {
  classic: {
    key: "classic",
    label: "Classic",
    dark: false
  },
  clean: {
    key: "clean",
    label: "Clean",
    dark: false
  },
  neon: {
    key: "neon",
    label: "Neon",
    dark: true
  },
  tundra: {
    key: "tundra",
    label: "Tundra",
    dark: true
  }
};
var themes = _themes;
var defaultLightTheme = "classic";
var defaultDarkTheme = "tundra";

// src/app/themes/theme-manager.service.ts
var LOCAL_STORAGE_KEY = "bitmagnet-theme";
var ThemeManager = class _ThemeManager {
  constructor() {
    this.document = inject(DOCUMENT);
    this.browserStorage = inject(BrowserStorageService);
    this._window = this.document.defaultView;
    this.selectedThemeSubject = new BehaviorSubject(void 0);
    this.selectedTheme$ = this.selectedThemeSubject.asObservable();
    this.themes = Object.values(themes);
    this.getPreferredTheme = () => {
      return this.getStoredTheme() ?? this.getAutoTheme();
    };
    this.getStoredTheme = () => {
      const value = this.browserStorage.get(LOCAL_STORAGE_KEY);
      return value && value in themes ? value : void 0;
    };
    this.getAutoTheme = () => {
      return this.windowMatchMediaPrefersDark()?.matches ? defaultDarkTheme : defaultLightTheme;
    };
    this.setTheme = (theme) => {
      this.setActiveTheme(theme);
      this.setStoredTheme(this.selectedTheme ?? "auto");
    };
    this.setActiveTheme = (theme) => {
      if (theme === "auto" || !(theme in themes)) {
        theme = this.getAutoTheme();
        this.selectedTheme = void 0;
      } else {
        this.selectedTheme = theme;
      }
      this.document.documentElement.setAttribute("data-bitmagnet-theme", theme);
      this.selectedThemeSubject.next(this.selectedTheme);
    };
    this.setStoredTheme = (theme) => {
      if (theme === "auto") {
        this.browserStorage.remove(LOCAL_STORAGE_KEY);
      } else {
        this.browserStorage.set(LOCAL_STORAGE_KEY, theme);
      }
    };
    this.setActiveTheme(this.getPreferredTheme());
    this.windowMatchMediaPrefersDark()?.addEventListener("change", () => {
      const storedTheme = this.getStoredTheme();
      if (!storedTheme) {
        this.setActiveTheme(this.getAutoTheme());
      }
    });
  }
  windowMatchMediaPrefersDark() {
    return this._window && this._window.matchMedia ? this._window.matchMedia("(prefers-color-scheme: dark)") : void 0;
  }
  static {
    this.\u0275fac = function ThemeManager_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _ThemeManager)();
    };
  }
  static {
    this.\u0275prov = /* @__PURE__ */ \u0275\u0275defineInjectable({ token: _ThemeManager, factory: _ThemeManager.\u0275fac, providedIn: "root" });
  }
};

// src/app/version/version.component.ts
var _c0 = (a0) => ({ version: a0 });
function VersionComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275elementStart(1, "span", 1);
    \u0275\u0275text(2);
    \u0275\u0275elementEnd();
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const t_r1 = ctx.$implicit;
    const ctx_r1 = \u0275\u0275nextContext();
    \u0275\u0275advance();
    \u0275\u0275property("matTooltip", t_r1("version.bitmagnet_version", \u0275\u0275pureFunction1(2, _c0, ctx_r1.versionUnknown ? t_r1("version.unknown") : ctx_r1.version)));
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(ctx_r1.version);
  }
}
var defaultVersionName = "v-unknown";
var VersionComponent = class _VersionComponent {
  constructor() {
    this.apollo = inject(Apollo);
    this.version = defaultVersionName;
    this.versionUnknown = true;
  }
  ngOnInit() {
    this.apollo.query({
      query: VersionDocument
    }).pipe(map((r) => r.data.version)).subscribe({
      next: (version) => {
        if (version) {
          this.version = version;
          this.versionUnknown = false;
        } else {
          this.version = defaultVersionName;
          this.versionUnknown = true;
        }
      },
      error: () => {
        this.version = defaultVersionName;
      }
    });
  }
  static {
    this.\u0275fac = function VersionComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _VersionComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _VersionComponent, selectors: [["app-version"]], standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 1, vars: 0, consts: [[4, "transloco"], [1, "version", 3, "matTooltip"]], template: function VersionComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275template(0, VersionComponent_ng_container_0_Template, 3, 4, "ng-container", 0);
      }
    }, dependencies: [AppModule, MatTooltip, TranslocoDirective, GraphQLModule], encapsulation: 2 });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(VersionComponent, { className: "VersionComponent", filePath: "src/app/version/version.component.ts", lineNumber: 16 });
})();

// src/app/i18n/translate-manager.service.ts
var LOCAL_STORAGE_KEY2 = "bitmagnet-language";
var TranslateManager = class _TranslateManager {
  constructor() {
    this.transloco = inject(TranslocoService);
    this.browserStorage = inject(BrowserStorageService);
    this.availableLanguages = this.transloco.getAvailableLangs();
    this.transloco.setActiveLang(this.getPreferredLanguage());
  }
  getPreferredLanguage() {
    return this.getStoredLanguage() ?? this.getAutoLanguage();
  }
  getStoredLanguage() {
    const value = this.browserStorage.get(LOCAL_STORAGE_KEY2);
    return value && this.transloco.isLang(value) ? value : void 0;
  }
  getAutoLanguage() {
    const navLang = navigator?.language?.split("-")?.[0];
    return this.transloco.isLang(navLang) ? navLang : "en";
  }
  setLanguage(lang) {
    this.transloco.setActiveLang(lang);
    this.browserStorage.set(LOCAL_STORAGE_KEY2, lang);
  }
  static {
    this.\u0275fac = function TranslateManager_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _TranslateManager)();
    };
  }
  static {
    this.\u0275prov = /* @__PURE__ */ \u0275\u0275defineInjectable({ token: _TranslateManager, factory: _TranslateManager.\u0275fac, providedIn: "root" });
  }
};

// src/app/themes/theme-emitter-color.component.ts
var _c02 = ["element"];
var ThemeEmitterColorComponent = class _ThemeEmitterColorComponent {
  static {
    this.\u0275fac = function ThemeEmitterColorComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _ThemeEmitterColorComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _ThemeEmitterColorComponent, selectors: [["app-theme-emitter-color"]], viewQuery: function ThemeEmitterColorComponent_Query(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275viewQuery(_c02, 5);
      }
      if (rf & 2) {
        let _t;
        \u0275\u0275queryRefresh(_t = \u0275\u0275loadQuery()) && (ctx.element = _t.first);
      }
    }, inputs: { color: "color" }, standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 2, vars: 2, consts: [["element", ""]], template: function ThemeEmitterColorComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275element(0, "div", null, 0);
      }
      if (rf & 2) {
        \u0275\u0275classMap("theme-emitter-color " + ctx.color);
      }
    }, encapsulation: 2 });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(ThemeEmitterColorComponent, { className: "ThemeEmitterColorComponent", filePath: "src/app/themes/theme-emitter-color.component.ts", lineNumber: 9 });
})();

// src/app/themes/theme-emitter.component.ts
var _c03 = ["lightdark"];
function ThemeEmitterComponent_For_2_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275element(0, "app-theme-emitter-color", 1);
  }
  if (rf & 2) {
    const c_r1 = ctx.$implicit;
    \u0275\u0275property("color", c_r1);
  }
}
var ThemeEmitterComponent = class _ThemeEmitterComponent {
  constructor() {
    this.service = inject(ThemeInfoService);
    this.themeManager = inject(ThemeManager);
    this.themeColors = themeColors;
    this.themeManager.selectedTheme$.subscribe(() => {
      this.updateThemeColors();
    });
  }
  ngAfterViewInit() {
    this.updateThemeColors();
  }
  updateThemeColors() {
    const colors = {};
    for (const color of this.elements ?? []) {
      colors[color.color] = getComputedStyle(color.element.nativeElement).color;
    }
    const type = this.lightdark && getComputedStyle(this.lightdark.nativeElement).color === "rgb(0, 0, 0)" ? "dark" : "light";
    this.service.setInfo({
      colors,
      type
    });
  }
  static {
    this.\u0275fac = function ThemeEmitterComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _ThemeEmitterComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _ThemeEmitterComponent, selectors: [["app-theme-emitter"]], viewQuery: function ThemeEmitterComponent_Query(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275viewQuery(_c03, 5);
        \u0275\u0275viewQuery(ThemeEmitterColorComponent, 5);
      }
      if (rf & 2) {
        let _t;
        \u0275\u0275queryRefresh(_t = \u0275\u0275loadQuery()) && (ctx.lightdark = _t.first);
        \u0275\u0275queryRefresh(_t = \u0275\u0275loadQuery()) && (ctx.elements = _t);
      }
    }, standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 5, vars: 0, consts: [["lightdark", ""], [3, "color"], [1, "theme-emitter-lightdark"]], template: function ThemeEmitterComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275elementContainerStart(0);
        \u0275\u0275repeaterCreate(1, ThemeEmitterComponent_For_2_Template, 1, 1, "app-theme-emitter-color", 1, \u0275\u0275repeaterTrackByIdentity);
        \u0275\u0275element(3, "div", 2, 0);
        \u0275\u0275elementContainerEnd();
      }
      if (rf & 2) {
        \u0275\u0275advance();
        \u0275\u0275repeater(ctx.themeColors);
      }
    }, dependencies: [ThemeEmitterColorComponent], styles: ["\n\n[_nghost-%COMP%] {\n  display: none;\n}\n.theme-emitter-color.background[_ngcontent-%COMP%] {\n  color: var(--mat-app-background-color);\n}\n.theme-emitter-color.foreground[_ngcontent-%COMP%] {\n  color: var(--mat-app-text-color);\n}\n/*# sourceMappingURL=theme-emitter.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(ThemeEmitterComponent, { className: "ThemeEmitterComponent", filePath: "src/app/themes/theme-emitter.component.ts", lineNumber: 22 });
})();

// src/app/layout/layout.component.ts
var _c04 = ["*"];
var _forTrack0 = ($index, $item) => $item.key;
var _forTrack1 = ($index, $item) => $item.id;
var _c1 = () => ({ service: "Discord" });
var _c2 = () => ({ service: "GitHub" });
var _c3 = () => ({ service: "OpenCollective" });
function LayoutComponent_ng_container_0_Conditional_2_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "h1")(1, "a", 22);
    \u0275\u0275element(2, "mat-icon", 23);
    \u0275\u0275elementStart(3, "span", 24);
    \u0275\u0275text(4, "bitmagnet");
    \u0275\u0275elementEnd();
    \u0275\u0275element(5, "app-version");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(6, "nav")(7, "a", 25, 3);
    \u0275\u0275element(9, "mat-icon", 23);
    \u0275\u0275text(10);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(11, "a", 26, 4)(13, "mat-icon");
    \u0275\u0275text(14, "dashboard");
    \u0275\u0275elementEnd();
    \u0275\u0275text(15);
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const linkTorrents_r1 = \u0275\u0275reference(8);
    const linkDashboard_r2 = \u0275\u0275reference(12);
    const t_r3 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance(7);
    \u0275\u0275classMap(linkTorrents_r1.isActive ? "active" : "");
    \u0275\u0275advance(3);
    \u0275\u0275textInterpolate1(" ", t_r3("routes.torrents"), " ");
    \u0275\u0275advance();
    \u0275\u0275classMap(linkDashboard_r2.isActive ? "active" : "");
    \u0275\u0275advance(4);
    \u0275\u0275textInterpolate1(" ", t_r3("routes.dashboard"), " ");
  }
}
function LayoutComponent_ng_container_0_Conditional_3_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "h1")(1, "a", 27);
    \u0275\u0275element(2, "mat-icon", 23);
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(3, "button", 28, 4)(5, "mat-icon");
    \u0275\u0275text(6, "dashboard");
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const linkDashboard_r4 = \u0275\u0275reference(4);
    const t_r3 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275advance(3);
    \u0275\u0275classMap(linkDashboard_r4.isActive ? "active" : "");
    \u0275\u0275property("matTooltip", t_r3("routes.dashboard"));
  }
}
function LayoutComponent_ng_container_0_For_12_Template(rf, ctx) {
  if (rf & 1) {
    const _r5 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "a", 29);
    \u0275\u0275listener("click", function LayoutComponent_ng_container_0_For_12_Template_a_click_0_listener() {
      const th_r6 = \u0275\u0275restoreView(_r5).$implicit;
      const ctx_r6 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r6.themeManager.setTheme(th_r6.key));
    });
    \u0275\u0275elementStart(1, "mat-icon");
    \u0275\u0275text(2);
    \u0275\u0275elementEnd();
    \u0275\u0275text(3);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const th_r6 = ctx.$implicit;
    const ctx_r6 = \u0275\u0275nextContext(2);
    \u0275\u0275classMap(th_r6.key === ctx_r6.themeManager.selectedTheme ? "active" : "");
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(th_r6.dark ? "dark_mode" : "light_mode");
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(th_r6.label);
  }
}
function LayoutComponent_ng_container_0_For_19_Template(rf, ctx) {
  if (rf & 1) {
    const _r8 = \u0275\u0275getCurrentView();
    \u0275\u0275elementStart(0, "a", 30);
    \u0275\u0275listener("click", function LayoutComponent_ng_container_0_For_19_Template_a_click_0_listener() {
      const l_r9 = \u0275\u0275restoreView(_r8).$implicit;
      const ctx_r6 = \u0275\u0275nextContext(2);
      return \u0275\u0275resetView(ctx_r6.translateManager.setLanguage(l_r9.id));
    });
    \u0275\u0275text(1);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const l_r9 = ctx.$implicit;
    const t_r3 = \u0275\u0275nextContext().$implicit;
    const ctx_r6 = \u0275\u0275nextContext();
    \u0275\u0275classMap(l_r9.id === ctx_r6.translateManager.getPreferredLanguage() ? "active" : "");
    \u0275\u0275property("matTooltip", l_r9.id === ctx_r6.translateManager.getPreferredLanguage() ? void 0 : t_r3("languages." + l_r9.id));
    \u0275\u0275advance();
    \u0275\u0275textInterpolate(l_r9.label);
  }
}
function LayoutComponent_ng_container_0_Conditional_32_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "a", 19)(1, "mat-icon");
    \u0275\u0275text(2, "favorite");
    \u0275\u0275elementEnd()();
  }
  if (rf & 2) {
    const t_r3 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275property("matTooltip", t_r3("layout.sponsor"));
  }
}
function LayoutComponent_ng_container_0_Conditional_33_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementStart(0, "a", 20)(1, "mat-icon");
    \u0275\u0275text(2, "favorite");
    \u0275\u0275elementEnd();
    \u0275\u0275text(3);
    \u0275\u0275elementEnd();
  }
  if (rf & 2) {
    const t_r3 = \u0275\u0275nextContext().$implicit;
    \u0275\u0275property("matTooltip", t_r3("layout.support_bitmagnet"));
    \u0275\u0275advance(3);
    \u0275\u0275textInterpolate(t_r3("layout.sponsor"));
  }
}
function LayoutComponent_ng_container_0_Template(rf, ctx) {
  if (rf & 1) {
    \u0275\u0275elementContainerStart(0);
    \u0275\u0275elementStart(1, "mat-toolbar", 6);
    \u0275\u0275template(2, LayoutComponent_ng_container_0_Conditional_2_Template, 16, 6)(3, LayoutComponent_ng_container_0_Conditional_3_Template, 7, 3);
    \u0275\u0275element(4, "span", 7)(5, "app-health-widget");
    \u0275\u0275elementStart(6, "button", 8)(7, "mat-icon");
    \u0275\u0275text(8, "apparel");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(9, "mat-menu", 9, 0);
    \u0275\u0275repeaterCreate(11, LayoutComponent_ng_container_0_For_12_Template, 4, 4, "a", 10, _forTrack0);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(13, "button", 11)(14, "mat-icon");
    \u0275\u0275text(15, "translate");
    \u0275\u0275elementEnd()();
    \u0275\u0275elementStart(16, "mat-menu", 9, 1);
    \u0275\u0275repeaterCreate(18, LayoutComponent_ng_container_0_For_19_Template, 2, 4, "a", 12, _forTrack1);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(20, "button", 13);
    \u0275\u0275element(21, "mat-icon", 14);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(22, "mat-menu", 9, 2)(24, "a", 15);
    \u0275\u0275text(25, "bitmagnet.io");
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(26, "a", 16);
    \u0275\u0275text(27);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(28, "a", 17);
    \u0275\u0275text(29);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(30, "a", 18);
    \u0275\u0275text(31);
    \u0275\u0275elementEnd()();
    \u0275\u0275template(32, LayoutComponent_ng_container_0_Conditional_32_Template, 3, 1, "a", 19)(33, LayoutComponent_ng_container_0_Conditional_33_Template, 4, 2, "a", 20);
    \u0275\u0275elementEnd();
    \u0275\u0275elementStart(34, "div", 21);
    \u0275\u0275projection(35);
    \u0275\u0275elementEnd();
    \u0275\u0275element(36, "app-theme-emitter");
    \u0275\u0275elementContainerEnd();
  }
  if (rf & 2) {
    const t_r3 = ctx.$implicit;
    const themesMenu_r10 = \u0275\u0275reference(10);
    const languagesMenu_r11 = \u0275\u0275reference(17);
    const externalLinksMenu_r12 = \u0275\u0275reference(23);
    const ctx_r6 = \u0275\u0275nextContext();
    \u0275\u0275advance(2);
    \u0275\u0275conditional(ctx_r6.breakpoints.sizeAtLeast("Medium") ? 2 : 3);
    \u0275\u0275advance(4);
    \u0275\u0275property("matMenuTriggerFor", themesMenu_r10)("matTooltip", t_r3("layout.change_theme"));
    \u0275\u0275advance(5);
    \u0275\u0275repeater(ctx_r6.themeManager.themes);
    \u0275\u0275advance(2);
    \u0275\u0275property("matMenuTriggerFor", languagesMenu_r11)("matTooltip", t_r3("layout.translate"));
    \u0275\u0275advance(5);
    \u0275\u0275repeater(ctx_r6.translateManager.availableLanguages);
    \u0275\u0275advance(2);
    \u0275\u0275property("matTooltip", t_r3("layout.external_links"))("matMenuTriggerFor", externalLinksMenu_r12);
    \u0275\u0275advance(7);
    \u0275\u0275textInterpolate(t_r3("layout.bitmagnet_on_service", \u0275\u0275pureFunction0(11, _c1)));
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(t_r3("layout.bitmagnet_on_service", \u0275\u0275pureFunction0(12, _c2)));
    \u0275\u0275advance(2);
    \u0275\u0275textInterpolate(t_r3("layout.bitmagnet_on_service", \u0275\u0275pureFunction0(13, _c3)));
    \u0275\u0275advance();
    \u0275\u0275conditional(!ctx_r6.breakpoints.sizeAtLeast("Medium") ? 32 : 33);
  }
}
var LayoutComponent = class _LayoutComponent {
  constructor() {
    this.themeManager = inject(ThemeManager);
    this.translateManager = inject(TranslateManager);
    this.breakpoints = inject(BreakpointsService);
    this.title = inject(Title);
    this.router = inject(Router);
    this.health = inject(HealthService);
  }
  static {
    this.\u0275fac = function LayoutComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _LayoutComponent)();
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _LayoutComponent, selectors: [["app-layout"]], standalone: true, features: [\u0275\u0275StandaloneFeature], ngContentSelectors: _c04, decls: 1, vars: 0, consts: [["themesMenu", "matMenu"], ["languagesMenu", "matMenu"], ["externalLinksMenu", "matMenu"], ["linkTorrents", "routerLinkActive"], ["linkDashboard", "routerLinkActive"], [4, "transloco"], [1, "toolbar-header"], [1, "spacer"], ["mat-icon-button", "", "aria-label", "Theme", 3, "matMenuTriggerFor", "matTooltip"], [1, "layout-header-menu"], ["mat-menu-item", "", 3, "class"], ["mat-icon-button", "", 3, "matMenuTriggerFor", "matTooltip"], ["mat-menu-item", "", "matTooltipPosition", "right", "matTooltipShowDelay", "500", 3, "matTooltip", "class"], ["mat-icon-button", "", 3, "matTooltip", "matMenuTriggerFor"], ["svgIcon", "external-link"], ["mat-menu-item", "", "href", "https://bitmagnet.io", "target", "_blank"], ["mat-menu-item", "", "href", "https://discord.gg/6mFNszX8qM", "target", "_blank"], ["mat-menu-item", "", "href", "https://github.com/bitmagnet-io/bitmagnet", "target", "_blank"], ["mat-menu-item", "", "href", "https://opencollective.com/bitmagnet", "target", "_blank"], ["mat-icon-button", "", "href", "https://opencollective.com/bitmagnet", "target", "_blank", 1, "button-sponsor", 3, "matTooltip"], ["mat-button", "", "href", "https://opencollective.com/bitmagnet", "target", "_blank", 1, "button-sponsor", 3, "matTooltip"], [1, "app-content"], ["routerLink", "torrents"], ["svgIcon", "magnet"], [1, "name"], ["mat-button", "", "routerLink", "torrents", "routerLinkActive", ""], ["mat-button", "", "routerLink", "dashboard", "routerLinkActive", ""], ["routerLink", "/torrents"], ["mat-icon-button", "", "routerLink", "dashboard", "routerLinkActive", "", 3, "matTooltip"], ["mat-menu-item", "", 3, "click"], ["mat-menu-item", "", "matTooltipPosition", "right", "matTooltipShowDelay", "500", 3, "click", "matTooltip"]], template: function LayoutComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275projectionDef();
        \u0275\u0275template(0, LayoutComponent_ng_container_0_Template, 37, 14, "ng-container", 5);
      }
    }, dependencies: [AppModule, MatAnchor, MatIconAnchor, MatIconButton, MatIcon, MatMenu, MatMenuItem, MatMenuTrigger, MatToolbar, MatTooltip, RouterLink, RouterLinkActive, TranslocoDirective, HealthModule, HealthWidgetComponent, ThemeEmitterComponent, VersionComponent], styles: ["\n\n.mat-toolbar.toolbar-header[_ngcontent-%COMP%] {\n  position: sticky;\n  top: 0;\n  z-index: 10;\n  --mat-toolbar-title-text-size: 22px;\n  --mat-toolbar-standard-height: 56px;\n  --mat-icon-color: #fff;\n  padding: 0 20px;\n  --mdc-text-button-label-text-size: 15px;\n}\n.mat-toolbar.toolbar-header[_ngcontent-%COMP%]   h1[_ngcontent-%COMP%] {\n  margin-top: -2px;\n  margin-right: 20px;\n}\n.mat-toolbar.toolbar-header[_ngcontent-%COMP%]   h1[_ngcontent-%COMP%]   a[_ngcontent-%COMP%] {\n  text-decoration: none;\n}\n.mat-toolbar.toolbar-header[_ngcontent-%COMP%]   h1[_ngcontent-%COMP%]   mat-icon[_ngcontent-%COMP%] {\n  position: relative;\n  top: 4px;\n}\n.mat-toolbar.toolbar-header[_ngcontent-%COMP%]   h1[_ngcontent-%COMP%]   .name[_ngcontent-%COMP%] {\n  margin-left: 10px;\n  margin-right: 16px;\n}\n.mat-toolbar.toolbar-header[_ngcontent-%COMP%]   h1[_ngcontent-%COMP%]   app-version[_ngcontent-%COMP%] {\n  font-size: 13px;\n}\n.mat-toolbar.toolbar-header[_ngcontent-%COMP%]   .mdc-button[_ngcontent-%COMP%] {\n  margin-left: 6px;\n  --mdc-text-button-label-text-weight: bold;\n  --mat-text-button-horizontal-padding: 12px;\n  --mdc-text-button-container-height: 30px;\n}\n.mat-toolbar.toolbar-header[_ngcontent-%COMP%]   nav[_ngcontent-%COMP%] {\n  position: absolute;\n  left: 340px;\n  top: 0;\n  height: 56px;\n  padding-top: 14px;\n}\n.mat-toolbar.toolbar-header[_ngcontent-%COMP%]   nav[_ngcontent-%COMP%]   a[_ngcontent-%COMP%] {\n  margin-right: 10px;\n}\n.mat-toolbar.toolbar-header[_ngcontent-%COMP%]   nav[_ngcontent-%COMP%]   mat-icon[data-mat-icon-name=magnet][_ngcontent-%COMP%] {\n  margin-top: -5px;\n  overflow: visible;\n}\n.app-content[_ngcontent-%COMP%] {\n  z-index: 1;\n  padding-bottom: 10px;\n}\n.layout-header-menu[_ngcontent-%COMP%]   a.active[_ngcontent-%COMP%] {\n  font-weight: bold;\n}\n/*# sourceMappingURL=layout.component.css.map */"] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(LayoutComponent, { className: "LayoutComponent", filePath: "src/app/layout/layout.component.ts", lineNumber: 20 });
})();

// src/app/app.icons.ts
var initializeIcons = (iconRegistry, domSanitizer) => iconRegistry.setDefaultFontSetClass("material-icons-outlined", "material-symbols-outlined").addSvgIcon("magnet", domSanitizer.bypassSecurityTrustResourceUrl("magnet.svg")).addSvgIcon("external-link", domSanitizer.bypassSecurityTrustResourceUrl("external-link.svg")).addSvgIcon("binary", domSanitizer.bypassSecurityTrustResourceUrl("binary.svg")).addSvgIcon("queue", domSanitizer.bypassSecurityTrustResourceUrl("queue.svg"));

// src/app/app.component.ts
var AppComponent = class _AppComponent {
  constructor(iconRegistry, domSanitizer) {
    this.title = "bitmagnet";
    initializeIcons(iconRegistry, domSanitizer);
  }
  static {
    this.\u0275fac = function AppComponent_Factory(__ngFactoryType__) {
      return new (__ngFactoryType__ || _AppComponent)(\u0275\u0275directiveInject(MatIconRegistry), \u0275\u0275directiveInject(DomSanitizer));
    };
  }
  static {
    this.\u0275cmp = /* @__PURE__ */ \u0275\u0275defineComponent({ type: _AppComponent, selectors: [["app-root"]], standalone: true, features: [\u0275\u0275StandaloneFeature], decls: 2, vars: 0, template: function AppComponent_Template(rf, ctx) {
      if (rf & 1) {
        \u0275\u0275elementStart(0, "app-layout");
        \u0275\u0275element(1, "router-outlet");
        \u0275\u0275elementEnd();
      }
    }, dependencies: [RouterOutlet, LayoutComponent] });
  }
};
(() => {
  (typeof ngDevMode === "undefined" || ngDevMode) && \u0275setClassDebugInfo(AppComponent, { className: "AppComponent", filePath: "src/app/app.component.ts", lineNumber: 15 });
})();

// src/main.ts
bootstrapApplication(AppComponent, appConfig).catch((err) => (
  // eslint-disable-next-line no-console
  console.error(err)
));
/*! Bundled license information:

@angular/platform-browser/fesm2022/animations/async.mjs:
  (**
   * @license Angular v18.2.8
   * (c) 2010-2024 Google LLC. https://angular.io/
   * License: MIT
   *)
*/
//# sourceMappingURL=main.js.map
