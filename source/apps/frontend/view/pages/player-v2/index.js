const formID = "#formFilterTemplate"
const tableID = "#tableTemplates"
const filterURL = "/player-v2/template";
const currentURL = "/player-v2/template";
let firstLoad = true;

$(document).ready(function () {
    //Click delete
    $(tableID).on("click", "tbody td div.btn-group a.remove", function (e) {
        var id = $(this).data("id")
        var isCheck = confirm('Are you sure delete!"');
        if (isCheck) {
            Delete(id)
        }
    });

    //Click duplicate
    $(tableID).on("click", "tbody td div.btn-group a.duplicate", function (e) {
        var id = $(this).data("id")
        var isCheck = confirm('Are you sure duplicate!');
        if (isCheck) {
            Duplicate(id)
        }
    });

    $(formID).find(".submit").on("click", function (e) {
        e.preventDefault();
        GetTable(true);
    });
    GetTable(false);
});


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
        // dom: '<"top float-start px-4 pb-3"fl><"clearfix"><"row"<"col-sm-12"tr>><"bottom d-flex align-items-center p-3 p-md-4 border-top border-gray-200"<"col-sm-12 col-md-5"<"ml-2"i>><"col-sm-12 col-md-7"<"mr-2"p>>><"clear">',
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
            { data: "name", name: "Name" },
            { data: "type", name: "Player Type" },
            { data: "action", name: "Action" },
        ],
        drawCallback: function () {
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
    let url = "/player/template/del"
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
            const msg = AjaxErrorMessage(jqXHR, exception);
            new AlertError("AJAX ERROR: " + msg);
        }
    }).done(function (result) {
        switch (result.status) {
            case "success":
                NoticeSuccess("Template has been removed successfully")
                GetTable(false)
                break
            case "err":
                new AlertError(result.message);
        }
    });
}

function Duplicate(id) {
    let url = "/player/template/duplicate"
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
                NoticeSuccess("Duplicate template success!")
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