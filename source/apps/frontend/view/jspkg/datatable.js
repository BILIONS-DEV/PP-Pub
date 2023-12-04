const {ErrorWithAlert} = require("./ajax-response");
module.exports = {
    Render
};

let firstLoad = true;

function Render(formID = "", isClickForm = false, refresh = false) {
    if (formID === "") {
        ErrorWithAlert("filter could not be found");
        return;
    }
    const formElement = $(formID);
    // Button
    let buttonElement = formElement.find(".submit");
    let buttonText = buttonElement.text();
    let buttonTextLoading = buttonElement.data("text-loading") || "Loading...";
    // TableID
    let tableID = formElement.data("table-id");
    let url = formElement.data("url");
    let columns = [];
    $(`table#${tableID} > thead > tr`).each(function () {
        $(this).find('th').each(function () {
            const name = $(this).text().trim(); // also tried val() here
            const field = $(this).data("obj-field");
            columns.push({data: field, name: name});
        });
    });
    // PostData
    let postData = formElement.serializeObject();
    let setting = {
        width: "100%",
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
        order: [parseInt(postData.order_column), postData.order_dir],
        columns: columns,
        columnDefs: [
            { "orderSequence": [ "desc", "asc"],"targets": '_all' },
        ],
        dom: '<"row"<"col-sm-12"tr>><"bottom d-flex p-3 border-top border-gray-200"<"row ms-auto align-items-center px-3"<"col-auto"fl><"col-auto"i><"col-auto"p>>>',
        ajax: {
            url: url,
            type: "POST",
            contentType: "application/json; charset=utf-8",
            data: function (d) {
                if (isClickForm) {
                    postData = formElement.serializeObject();
                    d.length = parseInt(postData.length);
                    d.start = parseInt(postData.start);
                    d.order[0].column = parseInt(postData.order_column);
                    d.order[0].dir = postData.order_dir;
                }
                if (refresh) {
                    d.start = 0;
                    d.order[0].column = 0;
                    d.order[0].dir = "desc";
                }
                formElement.find("[name='length']").val(d.length);
                formElement.find("[name='start']").val(d.start);
                formElement.find("[name='order_column']").val(d.order[0].column);
                formElement.find("[name='order_dir']").val(d.order[0].dir);
                postData.length = d.length;
                postData.start = d.start;
                d.postData = postData;
                return JSON.stringify(d);
            },
            beforeSend: function (xhr) {
                if (isClickForm) {
                    buttonElement.attr("disabled", true).text(buttonTextLoading);
                }
            },
            dataSrc: function (json) {
                if (isClickForm) {
                    buttonElement.attr("disabled", false).text(buttonText);
                }
                isClickForm = false;
                postData = formElement.serializeObject();
                if (!firstLoad) {
                    makeParamsUrl(url, postData);
                } else {
                    firstLoad = false;
                }
                if (json.data === null) {
                    return [];
                }
                return json.data;
            },
            error: function (jqXHR, exception) {
                new AlertError(jqXHR.responseText);
                if (isClickForm) {
                    buttonElement.attr("disabled", false).text(buttonText);
                }
            },
        },
        drawCallback: function () {
            $(".dataTables_paginate > ul.pagination > li > a.page-link").addClass("text-secondary");
            $(".dataTables_paginate > ul.pagination > li.active > a").addClass('bg-warning border-warning').css("color", "#111111");
            $(".dataTables_length > label > select").removeClass().addClass("form-select form-select-sm");
            if (!isClickForm) {
                $(".table-responsive").css("height", "");
            }
            const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
            const tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
                return new bootstrap.Tooltip(tooltipTriggerEl);
            });
            $("[data-bs-toggle=popover]").popover({
                trigger: "manual",
                html: true,
                animation: false,
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
                }, 300);
            });
        }
    };
    let pageLength;
    pageLength = parseFloat(postData.length);
    if (jQuery.inArray(pageLength, [10, 25, 50, 100]) >= 0) {
        setting.pageLength = pageLength;
    }
    setting.displayStart = parseFloat(postData.start);
    $(`#${tableID}`).DataTable(setting);
}

function makeParamsUrl(url, obj) {
    let params = jQuery.param(obj).replaceAll("%5B%5D", "");
    // let params = jQuery.param(obj)
    let newUrl = url + "?" + params;
    window.history.pushState("object or string", "Title", newUrl);
    window.history.replaceState("object or string", "Title", newUrl);
}