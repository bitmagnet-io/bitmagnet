import { BehaviorSubject, Observable } from "rxjs";
import * as generated from "../../graphql/generated";

export class ExpandedItem {
  private itemsSubject = new BehaviorSubject<generated.TorrentContent[]>([]);
  private itemSubject = new BehaviorSubject<
    generated.TorrentContent | undefined
  >(undefined);
  private editedTagsSubject = new BehaviorSubject<EditedTags | undefined>(
    undefined,
  );

  constructor(items: Observable<generated.TorrentContent[]>) {
    items.subscribe((items) => {
      this.itemsSubject.next(items);
      const item = this.itemSubject.getValue();
      if (item && !items.some((i) => i.id === item.id)) {
        this.itemSubject.next(undefined);
      }
    });
    this.itemSubject.subscribe((item) => {
      if (!item) {
        this.editedTagsSubject.next(undefined);
        return;
      }
      const editedTags = this.editedTagsSubject.getValue();
      if (editedTags?.infoHash !== item.infoHash) {
        this.editedTagsSubject.next({
          infoHash: item.infoHash,
        });
      }
    });
  }

  select(id: string) {
    this.itemSubject.next(
      this.itemsSubject.getValue().find((i) => i.id === id),
    );
  }

  toggle(id: string) {
    const item = this.itemSubject.getValue();
    if (item?.id === id) {
      this.itemSubject.next(undefined);
    } else {
      this.select(id);
    }
  }

  get id(): string | undefined {
    return this.itemSubject.getValue()?.id;
  }

  get editedTags(): EditedTags | undefined {
    return this.editedTagsSubject.getValue();
  }

  addTag(tagName: string) {
    this.editTags((tags) => [...tags, tagName]);
  }

  removeTag(tagName: string) {
    this.editTags((tags) => tags.filter((t) => t !== tagName));
  }

  renameTag(oldTagName: string, newTagName: string) {
    this.editTags((tags) =>
      tags.map((t) => (t === oldTagName ? newTagName : t)),
    );
  }

  private editTags(fn: (tagNames: string[]) => string[]) {
    const item = this.itemSubject.getValue();
    if (!item) {
      return;
    }
    const editedTags = {
      ...(this.editedTagsSubject.getValue() ?? { infoHash: item.infoHash }),
    };
    this.editedTagsSubject.next({
      ...editedTags,
      tagNames: fn(editedTags.tagNames ?? []),
    });
  }
}

type EditedTags = {
  infoHash: string;
  tagNames?: string[];
};
