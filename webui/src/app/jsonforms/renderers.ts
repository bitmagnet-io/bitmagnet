import { RankedTester } from "@jsonforms/core";
import {
  ArrayLayoutRenderer,
  ArrayLayoutRendererTester,
  BooleanControlRenderer,
  booleanControlTester,
  CategorizationTabLayoutRenderer,
  categorizationTester,
  DateControlRenderer,
  DateControlRendererTester,
  enumControlTester,
  GroupLayoutRenderer,
  groupLayoutTester,
  HorizontalLayoutRenderer,
  horizontalLayoutTester,
  LabelRenderer,
  LabelRendererTester,
  masterDetailTester,
  MasterListComponent,
  NumberControlRenderer,
  NumberControlRendererTester,
  ObjectControlRenderer,
  ObjectControlRendererTester,
  RangeControlRenderer,
  RangeControlRendererTester,
  TableRenderer,
  TableRendererTester,
  TextAreaRenderer,
  TextAreaRendererTester,
  TextControlRenderer,
  TextControlRendererTester,
  ToggleControlRenderer,
  ToggleControlRendererTester,
  VerticalLayoutRenderer,
  verticalLayoutTester,
} from "@jsonforms/angular-material";
import { SelectControlRenderer } from "./select.renderer";

export const angularMaterialRenderers: {
  tester: RankedTester;
  renderer: any;
}[] = [
  // controls
  { tester: booleanControlTester, renderer: BooleanControlRenderer },
  { tester: TextControlRendererTester, renderer: TextControlRenderer },
  { tester: TextAreaRendererTester, renderer: TextAreaRenderer },
  { tester: NumberControlRendererTester, renderer: NumberControlRenderer },
  { tester: RangeControlRendererTester, renderer: RangeControlRenderer },
  { tester: DateControlRendererTester, renderer: DateControlRenderer },
  { tester: ToggleControlRendererTester, renderer: ToggleControlRenderer },
  { tester: enumControlTester, renderer: SelectControlRenderer },
  { tester: ObjectControlRendererTester, renderer: ObjectControlRenderer },
  // layouts
  { tester: verticalLayoutTester, renderer: VerticalLayoutRenderer },
  { tester: groupLayoutTester, renderer: GroupLayoutRenderer },
  { tester: horizontalLayoutTester, renderer: HorizontalLayoutRenderer },
  { tester: categorizationTester, renderer: CategorizationTabLayoutRenderer },
  { tester: LabelRendererTester, renderer: LabelRenderer },
  { tester: ArrayLayoutRendererTester, renderer: ArrayLayoutRenderer },
  // other
  { tester: masterDetailTester, renderer: MasterListComponent },
  { tester: TableRendererTester, renderer: TableRenderer },
];
