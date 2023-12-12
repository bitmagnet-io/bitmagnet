import { Observable } from "rxjs";
import * as generated from "../../graphql/generated";

export type Agg<T, _allowNull extends boolean> = {
  count: number;
  label: string;
  value: FacetValue<T, _allowNull>;
};

export type GenreAgg = Agg<string, false>;
export type VideoResolutionAgg = Agg<generated.VideoResolution, true>;
export type VideoSourceAgg = Agg<generated.VideoSource, true>;

export type FacetValue<T = unknown, _allowNull extends boolean = true> =
  | T
  | (_allowNull extends true ? null : T);

export class Facet<T = unknown, _allowNull extends boolean = true> {
  private relevantContentTypes: Set<generated.ContentType> | null;

  aggregations: Agg<T, _allowNull>[] = [];

  private active: boolean = false;
  private selected = new Set<FacetValue<T, _allowNull>>();

  constructor(
    public name: string,
    public icon: string,
    contentTypes: generated.ContentType[] | null,
    obs: Observable<Agg<T, _allowNull>[]>,
  ) {
    this.relevantContentTypes =
      contentTypes === null ? null : new Set(contentTypes);
    obs.subscribe((aggs) => {
      this.aggregations = aggs;
    });
  }

  isRelevant(contentType?: generated.ContentType | "null" | null): boolean {
    return (
      this.relevantContentTypes === null ||
      (!!contentType &&
        contentType !== "null" &&
        this.relevantContentTypes.has(contentType))
    );
  }

  isActive(): boolean {
    return this.active;
  }

  isEmpty(): boolean {
    return this.selected.size === 0;
  }

  facetParams() {
    return this.active
      ? {
          aggregate: {
            limit: 10,
          },
          filter: this.selected.size ? [...this.selected] : undefined,
        }
      : undefined;
  }

  isSelected(value: FacetValue<T, _allowNull> | undefined): boolean {
    return value !== undefined && this.selected.has(value);
  }

  activate() {
    this.active = true;
  }

  deactivate() {
    this.active = false;
  }

  select(value: FacetValue<T, _allowNull> | undefined) {
    if (value !== undefined) {
      this.selected.add(value);
    }
  }

  deselect(value: FacetValue<T, _allowNull> | undefined) {
    if (value !== undefined) {
      this.selected.delete(value);
    }
  }

  toggle(value: FacetValue<T, _allowNull> | undefined) {
    if (value === undefined) {
      return;
    }
    if (this.selected.has(value)) {
      this.deselect(value);
    } else {
      this.select(value);
    }
  }

  reset() {
    this.selected.clear();
  }

  deactivateAndReset() {
    this.deactivate();
    this.reset();
  }
}
