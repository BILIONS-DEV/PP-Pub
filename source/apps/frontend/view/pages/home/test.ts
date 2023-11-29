// @ts-ignore
import {CheckBoxSelection, MultiSelect} from '@syncfusion/ej2-dropdowns';
// @ts-ignore
import {CheckBox} from '@syncfusion/ej2-buttons';
// @ts-ignore
import {enableRipple} from '@syncfusion/ej2-base';

enableRipple(true);

//Label position - Left.
let checkbox: CheckBox = new CheckBox({label: 'Left Side Label', labelPosition: 'Before'});
checkbox.appendTo('#checkbox1');

//Label position - Right.
checkbox = new CheckBox({label: 'Right Side Label', checked: true});
checkbox.appendTo('#checkbox2');

MultiSelect.Inject(CheckBoxSelection);

//define the array of complex data
let sportsData: { [key: string]: Object }[] = [
    {id: 'game1', sports: 'Badminton'},
    {id: 'game2', sports: 'Football'},
    {id: 'game3', sports: 'Tennis'},
    {id: 'game4', sports: 'Golf'},
    {id: 'game5', sports: 'Cricket'},
    {id: 'game6', sports: 'Handball'},
    {id: 'game7', sports: 'Karate'},
    {id: 'game8', sports: 'Fencing'},
    {id: 'game9', sports: 'Boxing'}
];

//initiate the MultiSelect
let msObject: MultiSelect = new MultiSelect({
    // bind the sports Data to datasource property
    dataSource: sportsData,
    // maps the appropriate column to fields property
    fields: {text: 'sports', value: 'id'},
    //set the placeholder to MultiSelect input
    placeholder: "Select games",
    // set the type of mode for checkbox to visualized the checkbox added in li element.
    mode: 'CheckBox',
    // set true for enable the selectAll support.
    showSelectAll: true,
    // set the select all text to MultiSelect checkbox label.
    selectAllText: "Select All"

});
//render the component
msObject.appendTo('#select');