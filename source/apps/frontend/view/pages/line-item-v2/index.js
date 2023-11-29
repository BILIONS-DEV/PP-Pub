const formID = "#formFilterBidder"
const tableID = "#tableBidders"
const filterURL = "/line-item-v2";
const currentURL = "/line-item-v2";
let firstLoad = true;

let domainOldSelected = [];
let formatOldSelected = [];
let sizeOldSelected = [];
let tagOldSelected = [];
let deviceOldSelected = [];
let geoOldSelected = [];

$(document).ready(function () {
    $("select.selectpicker").selectpicker('refresh');
    $(formID).find(".submit").on("click", function (e) {
        e.preventDefault();
        GetTable(true);
    });

    $(tableID).on("click", "tbody td div.btn-group a.remove", function (e) {
        var id = $(this).data("id")
        var isCheck = confirm('Are you sure delete!"');
        if (isCheck) {
            Delete(id)
        }
    });

    GetTable(false);
    InitSearch()
    InitValueSelected()
    HandleSelectOption()
});

function InitSearch() {
    SearchDomain()
    SearchAdFormat()
    SearchAdTag()
    SearchAdSize()
    SearchCountry()
    SearchDevice()
}

function InitValueSelected() {
    $('select.search-domain :selected').each(function () {
        domainOldSelected.push({
            id: parseInt($(this).val()),
            name: $(this).text(),
            selected: true
        })
    });

    $('select.search-adformat :selected').each(function () {
        formatOldSelected.push({
            id: parseInt($(this).val()),
            name: $(this).text(),
            selected: true
        })
    });

    $('select.search-adsize :selected').each(function () {
        sizeOldSelected.push({
            id: parseInt($(this).val()),
            name: $(this).text(),
            selected: true
        })
    });

    $('select.search-adtag :selected').each(function () {
        tagOldSelected.push({
            id: parseInt($(this).val()),
            name: $(this).text(),
            selected: true
        })
    });

    $('select.search-device :selected').each(function () {
        deviceOldSelected.push({
            id: parseInt($(this).val()),
            name: $(this).text(),
            selected: true
        })
    });

    $('select.search-country :selected').each(function () {
        geoOldSelected.push({
            id: parseInt($(this).val()),
            name: $(this).text(),
            selected: true
        })
    });

    if (tagOldSelected.length > 0){
        $("#inventory").val('default').attr('disabled',true).selectpicker("refresh");
        $("#ad_format").val('default').attr('disabled',true).selectpicker("refresh");
        $("#ad_size").val('default').attr('disabled',true).selectpicker("refresh");
    } else {
        $("#inventory").attr('disabled',false).selectpicker("refresh");
        $("#ad_format").attr('disabled',false).selectpicker("refresh");
        $("#ad_size").attr('disabled',false).selectpicker("refresh");
    }
}

function HandleSelectOption() {
    $("select.search-domain").on("changed.bs.select", function (e, clickedIndex, isSelected, oldValue) {
        if (clickedIndex !== null && isSelected !== null) {
            var name = $(this).find('option').eq(clickedIndex).text();
            var id = parseInt($(this).find('option').eq(clickedIndex).val());
            let currentIdx = domainOldSelected.map(item => {
                return item.id
            }).indexOf(id)
            if (currentIdx === -1) {
                domainOldSelected.push({
                    id: id,
                    name: name,
                    selected: true
                })
            } else {
                domainOldSelected.splice(currentIdx, 1)
            }
        }
    });

    $("select.search-adformat").on("changed.bs.select", function (e, clickedIndex, isSelected, oldValue) {
        if (clickedIndex !== null && isSelected !== null) {
            var name = $(this).find('option').eq(clickedIndex).text();
            var id = parseInt($(this).find('option').eq(clickedIndex).val());
            let currentIdx = formatOldSelected.map(item => {
                return item.id
            }).indexOf(id)
            if (currentIdx === -1) {
                formatOldSelected.push({
                    id: id,
                    name: name,
                    selected: true
                })
            } else {
                formatOldSelected.splice(currentIdx, 1)
            }
        }
    });

    $("select.search-adsize").on("changed.bs.select", function (e, clickedIndex, isSelected, oldValue) {
        if (clickedIndex !== null && isSelected !== null) {
            var name = $(this).find('option').eq(clickedIndex).text();
            var id = parseInt($(this).find('option').eq(clickedIndex).val());
            let currentIdx = sizeOldSelected.map(item => {
                return item.id
            }).indexOf(id)
            if (currentIdx === -1) {
                sizeOldSelected.push({
                    id: id,
                    name: name,
                    selected: true
                })
            } else {
                sizeOldSelected.splice(currentIdx, 1)
            }
        }
    });
    
    $("select.search-adtag").on("changed.bs.select", function (e, clickedIndex, isSelected, oldValue) {
        if (clickedIndex !== null && isSelected !== null) {
            var name = $(this).find('option').eq(clickedIndex).text();
            var id = parseInt($(this).find('option').eq(clickedIndex).val());
            let currentIdx = tagOldSelected.map(item => {
                return item.id
            }).indexOf(id);
            if (currentIdx === -1) {
                tagOldSelected.push({
                    id: id,
                    name: name,
                    selected: true
                });
            } else {
                tagOldSelected.splice(currentIdx, 1);
            }
            if (tagOldSelected.length > 0){
                $("#inventory").val('default').attr('disabled',true).selectpicker("refresh");
                $("#ad_format").val('default').attr('disabled',true).selectpicker("refresh");
                $("#ad_size").val('default').attr('disabled',true).selectpicker("refresh");
            } else {
                $("#inventory").attr('disabled',false).selectpicker("refresh");
                $("#ad_format").attr('disabled',false).selectpicker("refresh");
                $("#ad_size").attr('disabled',false).selectpicker("refresh");
            }
        }
    });

    $("select.search-country").on("changed.bs.select", function (e, clickedIndex, isSelected, oldValue) {
        if (clickedIndex !== null && isSelected !== null) {
            var name = $(this).find('option').eq(clickedIndex).text();
            var id = parseInt($(this).find('option').eq(clickedIndex).val());
            let currentIdx = geoOldSelected.map(item => {
                return item.id
            }).indexOf(id)
            if (currentIdx === -1) {
                geoOldSelected.push({
                    id: id,
                    name: name,
                    selected: true
                })
            } else {
                geoOldSelected.splice(currentIdx, 1)
            }
        }
    });

    $("select.search-device").on("changed.bs.select", function (e, clickedIndex, isSelected, oldValue) {
        if (clickedIndex !== null && isSelected !== null) {
            var name = $(this).find('option').eq(clickedIndex).text();
            var id = parseInt($(this).find('option').eq(clickedIndex).val());
            let currentIdx = deviceOldSelected.map(item => {
                return item.id
            }).indexOf(id)
            if (currentIdx === -1) {
                deviceOldSelected.push({
                    id: id,
                    name: name,
                    selected: true
                })
            } else {
                deviceOldSelected.splice(currentIdx, 1)
            }
        }
    });
}

function SearchDomain() {
    try {
        $('select.search-domain').selectpicker({
            liveSearch: true
        }).ajaxSelectPicker({
            ajax: {
                url: "/line-item-v2/searchDomain",
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
                    let currentIdx = InArray(data, item.id)
                    //nếu lựa chọn cũ không tồn tại trong data search mới -> append data
                    if (currentIdx === -1) {
                        data.push(domainOldSelected[index])
                    } else {
                        data[currentIdx].selected = true
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

function SearchAdFormat() {
    try {
        $('select.search-adformat').selectpicker({
            liveSearch: true
        }).ajaxSelectPicker({
            ajax: {
                url: "/line-item-v2/searchAdFormat",
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
                    let currentIdx = InArray(data, item.id)
                    //nếu lựa chọn cũ không tồn tại trong data search mới -> append data
                    if (currentIdx === -1) {
                        data.push(formatOldSelected[index])
                    } else {
                        data[currentIdx].selected = true
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

function SearchAdSize() {
    try {
        $('select.search-adsize').selectpicker({
            liveSearch: true
        }).ajaxSelectPicker({
            ajax: {
                url: "/line-item-v2/searchAdSize",
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
                    let currentIdx = InArray(data, item.id)
                    //nếu lựa chọn cũ không tồn tại trong data search mới -> append data
                    if (currentIdx === -1) {
                        data.push(sizeOldSelected[index])
                    } else {
                        data[currentIdx].selected = true
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

function SearchAdTag() {
    try {
        $('select.search-adtag').selectpicker({
            liveSearch: true
        }).ajaxSelectPicker({
            ajax: {
                url: "/line-item-v2/searchAdTag",
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
                    let currentIdx = InArray(data, item.id)
                    //nếu lựa chọn cũ không tồn tại trong data search mới -> append data
                    if (currentIdx === -1) {
                        data.push(tagOldSelected[index])
                    } else {
                        data[currentIdx].selected = true
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

function SearchDevice() {
    try {
        $('select.search-device').selectpicker({
            liveSearch: true
        }).ajaxSelectPicker({
            ajax: {
                url: "/line-item-v2/searchDevice",
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
                    let currentIdx = InArray(data, item.id)
                    //nếu lựa chọn cũ không tồn tại trong data search mới -> append data
                    if (currentIdx === -1) {
                        data.push(deviceOldSelected[index])
                    } else {
                        data[currentIdx].selected = true
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

function SearchCountry() {
    try {
        $('select.search-country').selectpicker({
            liveSearch: true
        }).ajaxSelectPicker({
            ajax: {
                url: "/line-item-v2/searchCountry",
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
                    let currentIdx = InArray(data, item.id)
                    //nếu lựa chọn cũ không tồn tại trong data search mới -> append data
                    if (currentIdx === -1) {
                        data.push(geoOldSelected[index])
                    } else {
                        data[currentIdx].selected = true
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

function GetTable(isClickForm = false) {
    const formElement = $(formID);
    let buttonElement = formElement.find(".submit");
    let submitButtonText = buttonElement.text();
    let submitButtonTextLoading = "Loading...";
    let postData = formElement.serializeObject();
    let setting = {
        processing: true,
        serverSide: true,
        searching: false,
        destroy: true,
        pagingType: "simple",
        language: {
            info: "_START_ - _END_ of _TOTAL_",
            infoEmpty: "0 - 0 of 0",
            lengthMenu: "Rows per page:  _MENU_",
            paginate: {
                next: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-chevron-right" viewBox="0 0 16 16"> <path fill-rule="evenodd" d="M4.646 1.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1 0 .708l-6 6a.5.5 0 0 1-.708-.708L10.293 8 4.646 2.354a.5.5 0 0 1 0-.708z"/> </svg>', // or '→'
                previous: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-chevron-left" viewBox="0 0 16 16"> <path fill-rule="evenodd" d="M11.354 1.646a.5.5 0 0 1 0 .708L5.707 8l5.647 5.646a.5.5 0 0 1-.708.708l-6-6a.5.5 0 0 1 0-.708l6-6a.5.5 0 0 1 .708 0z"/> </svg>' // or '←'
            }
        },
        order: [],
        // order: [0, 'desc'],
        // dom: '<"clearfix"><"row"<"col-sm-12"tr>><"bottom d-flex align-items-center p-3 p-md-4 border-top border-gray-200"<"col-md-12"<"mr-2"<"float-right"p><"float-right mr-3"i><"float-right mr-3 m10"fl>>>><"clear">',
        dom: '<"row"<"col-sm-12"tr>><"bottom d-flex p-3 border-top border-gray-200"<"row ms-auto align-items-center px-3"<"col-auto"fl><"col-auto"i><"col-auto"p>>>',
        ajax: {
            url: filterURL,
            type: "POST",
            contentType: "application/json; charset=utf-8",
            data: function (d) {
                postData.length = d.length;
                postData.start = d.start;

                d.postData = postData;
                return JSON.stringify(d);
            },
            beforeSend: function (xhr) {
                if (isClickForm) {
                    buttonElement.attr('disabled', true).text(submitButtonTextLoading);
                }
            },
            dataSrc: function (json) {
                if (isClickForm) {
                    buttonElement.attr('disabled', false).text(submitButtonText);
                }
                isClickForm = false;
                if (!firstLoad) {
                    makeParamsUrl(postData);
                } else {
                    firstLoad = false;
                }
                if (json.data === null) {
                    return []
                }
                return json.data;
            },
            error: function (jqXHR, exception) {
                // let msg = AjaxErrorMessage(jqXHR, exception)
                // new AlertError(msg);
                new AlertError(jqXHR.responseText);
                if (isClickForm) {
                    buttonElement.attr('disabled', false).text(submitButtonText);
                }
            },
        },
        columns: [
            // {data: "id", name: "ID"},
            {data: "name", name: "Name"},
            {data: "target", name: "Target"},
            {data: "status", name: "Status"},
            {data: "server_type", name: "Type"},
            {data: "priority", name: "Priority"},
            {data: "action", name: "Action"},
        ],
        drawCallback: function (settings) {
            $(".dataTables_paginate > ul.pagination > li > a.page-link").addClass("text-secondary");
            $(".dataTables_paginate > ul.pagination > li.active > a").addClass('bg-warning border-warning').css("color", "#111111");
            $(".dataTables_length > label > select").removeClass().addClass("form-select form-select-sm");
            if (!isClickForm) {
                $(".table-responsive").css("height", "");
            }
            const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
            const tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
                return new bootstrap.Tooltip(tooltipTriggerEl)
            });
            // $("[data-bs-toggle=popover]").popover({
            //     html: true,
            // });
            $("[data-bs-toggle=popover]").popover({
                trigger: "manual",
                html: true,
                animation: true,
                sanitize: false,
                sanitizeFn: function (content) {
                    return content;
                }
            }).on("mouseenter", function () {
                const _this = this;
                $(this).popover("show");
                $(".popover").on("mouseleave", function () {
                    $(_this).popover('hide');
                });
            }).on("mouseleave", function () {
                const _this = this;
                setTimeout(function () {
                    if (!$(".popover:hover").length) {
                        $(_this).popover("hide");
                    }
                }, 100);
            });
        }
    }
    let pageLength;
    pageLength = parseFloat(postData.length);
    if (jQuery.inArray(pageLength, [10, 25, 50, 100]) >= 0) {
        setting.pageLength = pageLength;
    }
    if (!isClickForm) {
        setting.displayStart = parseFloat(postData.start);
        // setting.order = JSON.parse(postData.order);
        // console.log("postData: ", JSON.parse(postData.order));
    }
    $(tableID).DataTable(setting);
}

function Delete(id) {
    let url = "/line-item-v2/del"
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
                NoticeSuccess("Line item has been removed successfully")
                GetTable(false)
                break
            case "err":
                new AlertError(result.message);
        }
    });
}

function makeParamsUrl(obj) {
    let params = jQuery.param(obj).replaceAll("%5B%5D", "")
    // let params = jQuery.param(obj)
    let newUrl = currentURL + "?" + params
    window.history.pushState("object or string", "Title", newUrl);
    window.history.replaceState("object or string", "Title", newUrl);
}

function InArray(array = [], id) {
    return array.map(item => {
        return item.id
    }).indexOf(id)
}