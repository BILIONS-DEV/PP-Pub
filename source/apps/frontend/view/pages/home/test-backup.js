window._APDOBJ = "--{{.DataObj}}--";
(function (doc, w) {
    let s = doc.createElement('script');
    s.src = '//"--{{.domain}}--"/p/"--{{.cacheParam}}--"/"--{{.uuid}}--".js?geo=' + _APDOBJ.geo + '&device=' + getDeviceType();
    let h = doc.getElementsByTagName('head')[0];
    h.appendChild(s);
})(document, window);

function getDeviceType() {
    let ua = navigator.userAgent;
    if (/(tablet|ipad|playbook|silk)|(android(?!.*mobi))/i.test(ua)) {
        return "tablet";
    }
    if (/Mobile|iP(hone|od)|Android|BlackBerry|IEMobile|Kindle|Silk-Accelerated|(hpw|web)OS|Opera M(obi|ini)/.test(ua)) {
        return "mobile";
    }
    return "desktop";
}

// import "../../quickstart/src/styles/styles.css"
// // import "../../quickstart/src/system.config"
// import {MultiSelect} from '@syncfusion/ej2-dropdowns';
//
//
// $(document).ready(function () {
//
//     // define the array of data
//     let sportsData = ['Badminton', 'Basketball', 'Cricket', 'Football', 'Golf', 'Gymnastics', 'Hockey', 'Rugby', 'Snooker', 'Tennis'];
//
//     // initialize MultiSelect component
//     let msObject = new MultiSelect({
//         //set the data to dataSource property
//         dataSource: sportsData
//     });
//
//     // render initialized MultiSelect
//     msObject.appendTo('#select');
//
//     console.log(msObject)
//     // DoneWithNotify("This is message with notify")
//     // DoneWithAlert("message alert", myCallback)
//     // PrintResponse()
//     // SubmitForms()
// })
//
// function myCallback(mess) {
//     alert("My callback: " + mess)
// }
//
// // // const {Hello} = require("./test_include_function")
// // // const func = require("./test_import_function")
// // const {The, SCopyToClipboardWithElement} = require("../.js_libs");
// // // const {Say} = require("./test_import_function");
// //
// // $(document).ready(function () {
// //     // The()
// //     SCopyToClipboardWithElement("hihi")
// //
// //     // Hello()
// //     // World()
// //
// //     // Say()
// //     // func.Say()
// //     // func.Hi()
// //
// // })