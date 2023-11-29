import dataTable from "../../jspkg/datatable";

const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
const tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
    return new bootstrap.Tooltip(tooltipTriggerEl);
});

const formID = "#filterAbTesting";
const tableID = "#tableAbTestings";
const currentURL = "/ab-testing";

let domainOldSelected = [];
let formatOldSelected = [];
let sizeOldSelected = [];
let tagOldSelected = [];
let deviceOldSelected = [];
let geoOldSelected = [];

$(document).ready(function () {

    $(formID).find(".submit").on("click", function (e) {
        dataTable.Render(formID, true);
    });
    dataTable.Render(formID);

    $(formID).keypress(function (e) {
        const key = e.which;
        if (key === 13) { // the enter key code
            dataTable.Render(formID, true);
            return false;
        }
    });

    $(tableID).on("click", "tbody td div.btn-group a.remove", function (e) {
        let id = $(this).data("id");
        let isCheck = confirm('Are you sure delete!"');
        if (isCheck) {
            Delete(id);
        }
    });

    InitSearch();

});

function InitSearch() {
    SearchDomain();
    SearchAdFormat();
    SearchAdTag();
    SearchAdSize();
    SearchCountry();
    SearchDevice();
}

function Delete(id) {
    let url = "/ab-testing/del";
    $.ajax({
        url: url,
        type: "POST",
        dataType: "JSON",
        contentType: "application/json",
        data: JSON.stringify({
            id: id
        }),
        beforeSend: function (xhr) {
            xhr.overrideMimeType("text/plain; charset=x-user-defined");
        },
        error: function (jqXHR, exception) {
            var msg = AjaxErrorMessage(jqXHR, exception);
            new AlertError("AJAX ERROR: " + msg);
        }
    }).done(function (result) {
        switch (result.status) {
            case "success":
                NoticeSuccess("A/B Testing has been removed successfully");
                dataTable.Render(formID, false);
                break;
            case "err":
                new AlertError(result.message);
        }
    });
}

function SearchDomain() {
    try {
        $('select.search-domain').selectpicker({
            liveSearch: true
        }).ajaxSelectPicker({
            ajax: {
                url: "/line-item/searchDomain",
                type: "GET",
                data: function () {
                    return {
                        q: '{{{q}}}'
                    };
                }
            },
            locale: {
                emptyTitle: 'All',
                statusInitialized: ''
            },
            cache: false,
            clearOnEmpty: true,
            clearOnError: true,
            emptyRequest: true,
            preserveSelected: false,
            preprocessData: function (data) {
                domainOldSelected.map((item, index) => {
                    let currentIdx = InArray(data, item.id);
                    //nếu lựa chọn cũ không tồn tại trong data search mới -> append data
                    if (currentIdx === -1) {
                        data.push(domainOldSelected[index]);
                    } else {
                        data[currentIdx].selected = true;
                    }
                });
                let contacts = [];
                let len = data.length;
                for (let i = 0; i < len; i++) {
                    let curr = data[i];
                    contacts.push({
                        value: curr.id,
                        text: curr.name,
                        disabled: false,
                        selected: curr.selected
                    });
                }
                return contacts;
            },
        });
    } catch (e) {
        console.log(e);
    }
}

function SearchAdFormat() {
    try {
        $('select.search-adformat').selectpicker({
            liveSearch: true
        }).ajaxSelectPicker({
            ajax: {
                url: "/line-item/searchAdFormat",
                type: "GET",
                data: function () {
                    return {
                        q: '{{{q}}}'
                    };
                }
            },
            locale: {
                emptyTitle: 'All',
                statusInitialized: ''
            },
            cache: false,
            clearOnEmpty: true,
            clearOnError: true,
            emptyRequest: true,
            preserveSelected: false,
            preprocessData: function (data) {
                formatOldSelected.map((item, index) => {
                    let currentIdx = InArray(data, item.id);
                    //nếu lựa chọn cũ không tồn tại trong data search mới -> append data
                    if (currentIdx === -1) {
                        data.push(formatOldSelected[index]);
                    } else {
                        data[currentIdx].selected = true;
                    }
                });
                let contacts = [];
                let len = data.length;
                for (let i = 0; i < len; i++) {
                    let curr = data[i];
                    contacts.push({
                        value: curr.id,
                        text: curr.name,
                        disabled: false,
                        selected: curr.selected
                    });
                }
                return contacts;
            },
        });
    } catch (e) {
        console.log(e);
    }
}

function SearchAdSize() {
    try {
        $('select.search-adsize').selectpicker({
            liveSearch: true
        }).ajaxSelectPicker({
            ajax: {
                url: "/line-item/searchAdSize",
                type: "GET",
                data: function () {
                    return {
                        q: '{{{q}}}'
                    };
                }
            },
            locale: {
                emptyTitle: 'All',
                statusInitialized: ''
            },
            cache: false,
            clearOnEmpty: true,
            clearOnError: true,
            emptyRequest: true,
            preserveSelected: false,
            preprocessData: function (data) {
                sizeOldSelected.map((item, index) => {
                    let currentIdx = InArray(data, item.id);
                    //nếu lựa chọn cũ không tồn tại trong data search mới -> append data
                    if (currentIdx === -1) {
                        data.push(sizeOldSelected[index]);
                    } else {
                        data[currentIdx].selected = true;
                    }
                });
                let contacts = [];
                let len = data.length;
                for (let i = 0; i < len; i++) {
                    let curr = data[i];
                    contacts.push({
                        value: curr.id,
                        text: curr.name,
                        disabled: false,
                        selected: curr.selected
                    });
                }
                return contacts;
            },
        });
    } catch (e) {
        console.log(e);
    }
}

function SearchAdTag() {
    try {
        $('select.search-adtag').selectpicker({
            liveSearch: true
        }).ajaxSelectPicker({
            ajax: {
                url: "/line-item/searchAdTag",
                type: "GET",
                data: function () {
                    return {
                        q: '{{{q}}}'
                    };
                }
            },
            locale: {
                emptyTitle: 'All',
                statusInitialized: ''
            },
            cache: false,
            clearOnEmpty: true,
            clearOnError: true,
            emptyRequest: true,
            preserveSelected: false,
            preprocessData: function (data) {
                tagOldSelected.map((item, index) => {
                    let currentIdx = InArray(data, item.id);
                    //nếu lựa chọn cũ không tồn tại trong data search mới -> append data
                    if (currentIdx === -1) {
                        data.push(tagOldSelected[index]);
                    } else {
                        data[currentIdx].selected = true;
                    }
                });
                let contacts = [];
                let len = data.length;
                for (let i = 0; i < len; i++) {
                    let curr = data[i];
                    contacts.push({
                        value: curr.id,
                        text: curr.name,
                        disabled: false,
                        selected: curr.selected
                    });
                }
                return contacts;
            },
        });
    } catch (e) {
        console.log(e);
    }
}

function SearchDevice() {
    try {
        $('select.search-device').selectpicker({
            liveSearch: true
        }).ajaxSelectPicker({
            ajax: {
                url: "/line-item/searchDevice",
                type: "GET",
                data: function () {
                    return {
                        q: '{{{q}}}'
                    };
                }
            },
            locale: {
                emptyTitle: 'All',
                statusInitialized: ''
            },
            cache: false,
            clearOnEmpty: true,
            clearOnError: true,
            emptyRequest: true,
            preserveSelected: false,
            preprocessData: function (data) {
                deviceOldSelected.map((item, index) => {
                    let currentIdx = InArray(data, item.id);
                    //nếu lựa chọn cũ không tồn tại trong data search mới -> append data
                    if (currentIdx === -1) {
                        data.push(deviceOldSelected[index]);
                    } else {
                        data[currentIdx].selected = true;
                    }
                });
                let contacts = [];
                let len = data.length;
                for (let i = 0; i < len; i++) {
                    let curr = data[i];
                    contacts.push({
                        value: curr.id,
                        text: curr.name,
                        disabled: false,
                        selected: curr.selected
                    });
                }
                return contacts;
            },
        });
    } catch (e) {
        console.log(e);
    }
}

function SearchCountry() {
    try {
        $('select.search-country').selectpicker({
            liveSearch: true
        }).ajaxSelectPicker({
            ajax: {
                url: "/line-item/searchCountry",
                type: "GET",
                data: function () {
                    return {
                        q: '{{{q}}}'
                    };
                }
            },
            locale: {
                emptyTitle: 'All',
                statusInitialized: ''
            },
            cache: false,
            clearOnEmpty: true,
            clearOnError: true,
            emptyRequest: true,
            preserveSelected: false,
            preprocessData: function (data) {
                geoOldSelected.map((item, index) => {
                    let currentIdx = InArray(data, item.id);
                    //nếu lựa chọn cũ không tồn tại trong data search mới -> append data
                    if (currentIdx === -1) {
                        data.push(geoOldSelected[index]);
                    } else {
                        data[currentIdx].selected = true;
                    }
                })
                let contacts = [];
                let len = data.length;
                for (let i = 0; i < len; i++) {
                    let curr = data[i];
                    contacts.push({
                        value: curr.id,
                        text: curr.name,
                        disabled: false,
                        selected: curr.selected
                    });
                }
                return contacts;
            },
        });
    } catch (e) {
        console.log(e);
    }
}

function makeParamsUrl(obj) {
    let params = jQuery.param(obj).replaceAll("%5B%5D", "");
    let newUrl = currentURL + "?" + params;
    window.history.pushState("object or string", "Title", newUrl);
    window.history.replaceState("object or string", "Title", newUrl);
}