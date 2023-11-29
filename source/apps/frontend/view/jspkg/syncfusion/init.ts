// Tool build theme: https://ej2.syncfusion.com/themestudio/?theme=material

// @ts-ignore
import { CheckBoxSelection, MultiSelect } from "@syncfusion/ej2-dropdowns";

MultiSelect.Inject(CheckBoxSelection);
let defaultSetting = {
  mode: "CheckBox",
  placeholder: " Select All",
  showSelectAll: true,
  selectAllText: "Select All",
  // popupWidth: 300,
  showDropDownIcon: true,
};

/**
 * Init select with checkbox
 * @constructor
 */
export function InitSelectWithCheckbox(
  className: string = "sfSelectWithCheckbox"
) {
  appendCss([
    // "http://127.0.0.1:8545/assets/test/index.css?v=2"
    // "https://cdn.syncfusion.com/ej2/ej2-base/styles/material.css",
    // "https://cdn.syncfusion.com/ej2/ej2-inputs/styles/material.css",
    // "https://cdn.syncfusion.com/ej2/ej2-dropdowns/styles/material.css",
    // "https://cdn.syncfusion.com/ej2/ej2-buttons/styles/material.css",
  ]);
  // @ts-ignore
  $(`.${className}`).map(function () {
    let divID = guidGenerator();
    let setting = defaultSetting;

    // @ts-ignore
    const element = $(this);
    let multiselectObj: MultiSelect = new MultiSelect(setting);

    element.attr("id", divID);
    multiselectObj.appendTo(`#${divID}`);
  });
}

/**
 * Append CSS to <head>
 * @param CSS
 */
function appendCss(CSS: Array<string>) {
  CSS.forEach(function (css) {
    // @ts-ignore
    $("head").append(`<link rel="stylesheet" href="${css}" type="text/css" />`);
  });
}

/**
 * Generator random id
 */
function guidGenerator() {
  return Math.random().toString(36).substr(2, 9);
}
