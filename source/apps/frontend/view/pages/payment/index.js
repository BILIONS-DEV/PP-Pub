const formID = "#formFilter";
const tableID = "#tableResponse";
const filterURL = "/payment";
const currentURL = "/payment";
let firstLoad = true;

$(document).ready(function () {
    $(formID).find(".submit").on("click", function (e) {
        e.preventDefault();
        GetTable(true);
    });
    GetTable(false);

    $("#tableResponse").on("click", ".load-pdf", function () {
        PreviewInvoice($(this));
    });
});

function PreviewInvoice(el) {
    // var id = el.attr("data-id")
    // if (!id) {
    //     return
    // }
    // $.ajax({
    //     type: 'GET',
    //     url: '/payment/preview',
    //     data: {invoiceId: parseInt(id)}
    // })
    //     .done(function (result) {
    //         if (result.error) {
    //             return;
    //         }
    //         $("#preview-invoice").find(".modal-body").html(result)
    //     })
    var PDF = el.attr("data-pdf");
    var height = window.innerHeight;
    $("#InvoicePDFModal").find(".modal-body").html('<iframe src="' + PDF + '" style="min-height: '+(height - 90)+'px"></iframe>');
}

function GetTable(isClickForm = false) {
    const formElement = $(formID);
    let buttonElement = formElement.find(".submit");
    let submitButtonText = buttonElement.text();
    let submitButtonTextLoading = "Loading...";
    let postData = formElement.serializeObject();
    let setting = {
        "ordering": false,
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
        order: [0, 'desc'],
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
            {data: "period", name: "Period"},
            {data: "amount", name: "Amount"},
            {data: "info", name: "Info"},
            {data: "status", name: "Status"},
            // {data: "action", name: "Action"},
        ],
        drawCallback: function () {
            $(".dataTables_paginate > ul.pagination > li > a.page-link").addClass("text-secondary");
            $(".dataTables_paginate > ul.pagination > li.active > a").addClass('bg-warning border-warning').css("color", "#111111");
            $(".dataTables_length > label > select").removeClass().addClass("form-select form-select-sm");
            if (!isClickForm) {
                $(".table-responsive").css("height", "");
            }
            $("[data-bs-toggle=popover]").popover({
                html: true,
            })
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

function makeParamsUrl(obj) {
    let params = jQuery.param(obj).replaceAll("%5B%5D", "")
    // let params = jQuery.param(obj)
    let newUrl = currentURL + "?" + params
    window.history.pushState("object or string", "Title", newUrl);
    window.history.replaceState("object or string", "Title", newUrl);
}