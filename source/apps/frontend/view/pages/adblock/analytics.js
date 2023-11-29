const formID = "#formFilterAdblock"
const tableID = "#tableBidders"
const filterURL = "/adblock"
const currentURL = "/adblock/analytics";
let firstLoad = true;

const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
const tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
    return new bootstrap.Tooltip(tooltipTriggerEl)
});

$(document).ready(function () {
    SubmitFilter();
    $(formID).find(".submit").on("click", function (e) {
        e.preventDefault();
        // GetTable(true);
        SubmitFilter();
    });


    let today = new Date();
    var currentYear = today.getFullYear();
    var currentMonth = today.getMonth();
    var currentDay = today.getDate();
    var daterangepicker = new ej.calendars.DateRangePicker({
        placeholder: 'Select a range',
        format: "dd-MM-yyyy",
        // value: "10/01/2023 - 10/30/2023",
        startDate: new Date(new Date().setDate(new Date().getDate() - 10)),
        endDate: new Date(),
        max: new Date(),
        presets: [
            {label: 'Today', start: new Date(), end: new Date()},
            {label: 'This Month', start: new Date(new Date().setDate(1)), end: new Date()},
            {
                label: 'Last Month',
                start: new Date(new Date(new Date().setMonth(new Date().getMonth() - 1)).setDate(1)),
                end: new Date(new Date().setDate(0))
            },
            {label: 'Last Year', start: new Date(new Date().getFullYear() - 1, 0, 1), end: new Date()},

        ],
        change: function (args) {
            if (args.isInteracted) { // Đảm bảo rằng sự kiện là do tương tác người dùng
                console.log(args);
                // Cập nhật giá trị của input "startDate" và "endDate"
                var startDateInput = document.querySelector('[name="startDate"]');
                var endDateInput = document.querySelector('[name="endDate"]');
                if (args.startDate) {
                    const startYear = args.startDate.getFullYear();
                    const startMonth = (args.startDate.getMonth() + 1).toString().padStart(2, '0'); // Thêm 0 đằng trước nếu cần
                    const startDay = args.startDate.getDate().toString().padStart(2, '0'); // Thêm 0 đằng trước nếu cần
                    startDateInput.value = `${startYear}-${startMonth}-${startDay}`; // Thay đổi định dạng ngày tháng tùy theo nhu cầu
                } else {
                    startDateInput.value = ""
                }
                if (args.endDate) {
                    const endyYear = args.endDate.getFullYear();
                    const endMonth = (args.endDate.getMonth() + 1).toString().padStart(2, '0'); // Thêm 0 đằng trước nếu cần
                    const endDay = args.endDate.getDate().toString().padStart(2, '0'); // Thêm 0 đằng trước nếu cần
                    endDateInput.value = `${endyYear}-${endMonth}-${endDay}`; // Thay đổi định dạng ngày tháng tùy theo nhu cầu
                } else {
                    endDateInput.value = ""
                }
            }
        }
    });
    daterangepicker.appendTo('#date_syncfusion');
});

function SubmitFilter() {
    const formElement = $(formID);
    let buttonElement = formElement.find(".submit");
    buttonElement.text("Loading...");
    let postData = formElement.serializeObject();
    delete postData.date_syncfusion;
    if (!postData.inventory_id) {
        postData.inventory_id = 0;
    }
    postData.inventory_id = parseInt(postData.inventory_id);
    console.log(postData);

    makeParamsUrl(postData);

    $.ajax({
        type: 'POST',
        url: '/adblock/analytics',
        data: postData
    }).done(function (result) {
        buttonElement.text("Run");
        console.log(result);
        highcharts(result.AdblockAnalytics);
    });
}

function highcharts(data) {
    Highcharts.chart('highchart-adblock', {
        chart: {
            type: 'column'
        },
        title: {
            text: ''
        },
        subtitle: {
            text: 'The total page views with adblock plugins detected'
        },
        xAxis: {
            categories: data.date,
            crosshair: true
        },
        yAxis: {
            min: 0,
            title: {
                text: ''
            }
        },
        tooltip: {
            headerFormat: '<span style="font-size:10px">{point.key}</span><table>',
            pointFormat: '<tr><td style="color:{series.color};padding:0">{series.name}: </td>' +
                '<td style="padding:0"><b>{point.y}</b></td></tr>',
            footerFormat: '</table>',
            shared: true,
            useHTML: true
        },
        plotOptions: {
            column: {
                pointPadding: 0.2,
                borderWidth: 0
            }
        },
        colors: [
            'rgb(213, 219, 219)',
        ],
        series: [{
            name: 'Found',
            data: data.found

        }]
    });
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
            // {data: "display_name", name: "Display Name"},
            {data: "name", name: "Bidder"},
            {data: "type", name: "Type"},
            {data: "status", name: "Status"},
            {data: "media_type", name: "Media Type"},
            {data: "bid_adjustment", name: "Bid Adjustment"},
            {data: "action", name: "Action"},
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
    let url = "/bidder/del"
    $.ajax({
        url: url,
        type: "POST",
        dataType: "JSON",
        contentType: "application/json",
        data: JSON.stringify({
            id: id
        }),
        beforeSend: function (xhr) {
            // xhr.overrideMimeType("text/plain; charset=x-user-defined");
        },
        error: function (jqXHR, exception) {
            var msg = AjaxErrorMessage(jqXHR, exception);
            new AlertError("AJAX ERROR: " + msg);
        }
    }).done(function (result) {
        switch (result.status) {
            case "success":
                NoticeSuccess("Bidder has been removed successfully")
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
