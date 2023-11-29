const formID = "#formFilterHistory";
const filterURL = "/history/";
const currentURL = "/history/";
const loadObjectByPage = "/history/object-by-page/";

let firstLoad = true;

const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
const tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
    return new bootstrap.Tooltip(tooltipTriggerEl)
});

// if ($(".pubpower-main").attr("data-windows") == "false") {
//     filterURL = "https://be.valueimpression.com/history/";
//     currentURL = "https://be.valueimpression.com/history/";
//     loadObjectByPage = "https://be.valueimpression.com/history/object-by-page/";
// }

$(document).ready(function () {
    $(".list-history").css("max-height", (screen.height - 650) + "px");

    // var time_zone = 'America/New_York';
    // moment.tz.setDefault(time_zone);
    var defaultStartDate = moment().subtract(100, 'days');
    var defaultEndDate = moment().subtract(0, 'days');

    $('#date_filter').daterangepicker({
        opens: 'left',
        ranges: {
            'Default': [moment().subtract(13, 'days'), moment()],
            'Today': [moment(), moment()],
            'Yesterday': [moment().subtract(1, 'days'), moment().subtract(1, 'days')],
            'Last 7 days': [moment().subtract(7, 'days'), moment().subtract(1, 'days')],
            'Last 30 days': [moment().subtract(30, 'days'), moment().subtract(1, 'days')],
            'This Month': [moment().startOf('month'), moment().endOf('month')],
            'Last Month': [moment().subtract(1, 'month').startOf('month'), moment().subtract(1, 'month').endOf('month')],
            'All Time': [moment.unix(1611196933), moment()]
        },
        autoUpdateInput: true,
        alwaysShowCalendars: true,
        autoApply: false,
        startDate: defaultStartDate,
        endDate: defaultEndDate
    }, function (start, end, label) {
        sm(start, end, label);
    });


    $(formID).find(".submit").on("click", function (e) {
        $(this).text("Loading...")
        e.preventDefault();
        $("#page").val(0);
        FiltersHistory();
    });
    $(".load-more").on("click", function (e) {
        $(this).text("Loading...")
        e.preventDefault();
        $("#page").val(parseInt($("#page").val()) + 1);
        FiltersHistory(false);
        // GetTable(true);
    });
    FiltersHistory();

    // Load History
    $('.list-history').on("click", '.load-history', function () {
        var id = $(this).closest(".list-group-item").attr("data-id");
        loadHistory(id)
    });

    // select object
    $("#object_page").change(function () {
        var page = $(this).val();
        if (!page) {
            $(".filter_objectId").addClass("d-none");
        }
        $.ajax({
            type: 'POST',
            url: loadObjectByPage,
            data: {object_page: page}
        })
            .done(function (result) {
                if (result.status == 'success') {
                    if (!result.data) {
                        $(".filter_objectId").addClass("d-none");
                    } else {
                        if (result.data.length > 0) {
                            $("#object_id").find("option").remove()
                            $("#object_id").append('<option value="">All</option>')
                            result.data.forEach(function (value, index) {
                                console.log(value);
                                var option = '<option value="' + value.ID + '" >' + value.Name + '</option>';
                                $("#object_id").append(option)
                                $('#object_id').selectpicker('refresh');
                                $(".filter_objectId").removeClass("d-none")
                            })
                        } else {
                            $(".filter_objectId").addClass("d-none");
                        }
                    }

                }
                if (result.status == 'error') {
                    new AlertError(result.message);
                }

                $(".filter_objectId .filter-option-inner-inner").text("All " + page)
            })
    })

    $(".result-history-item").on("change", ".loadType", function () {
        checkTypeHistory($(".result-history-item"), 'on');
    });
});

function loadHistory(id) {
    // $(".result-history-item .modal-body").css("opacity", 0.1)
    $(".result-history-item").find(".table-history tbody").css("opacity", "0");
    $.ajax({
        url: filterURL + id,
        type: "POST",
        dataType: "JSON",
        contentType: "application/json; charset=utf-8",
    })
        .done(function (result) {
            if (result.status == true) {
                // $(".modal-content").html(renderModalHistory(result.data.Compare, result.data.Row, result.data.CreateTime));
                $(".result-history-item").html(renderHtmlHistory(result.data.Compare, result.data.Row, result.data.CreateTime));
                checkTypeHistory($(".result-history-item"), 'on', true);
                $('[data-bs-toggle="tooltip"]').tooltip();
                $(".pubpower-main").find(".table-responsive").css("max-height", (screen.height - 600) + "px");
            }
        });
}

// function renderModalHistory(Compare, history, createTime) {
//     var html = '<div class="modal-header border-0 pb-0 align-items-start px-4">  ' +
//         '  <h5 class="modal-title" id="exampleModalLabel">History: ' + history.Title + '</h5>  ' +
//         '  <button type="button" class="btn btn-icon p-0" data-bs-dismiss="modal" aria-label="Close">  ' +
//         '      <svg data-name="icons/tabler/close" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16">  ' +
//         '          <rect data-name="Icons/Tabler/Close background" width="16" height="16" fill="none"></rect>  ' +
//         '          <path d="M.82.1l.058.05L6,5.272,11.122.151A.514.514,0,0,1,11.9.82l-.05.058L6.728,6l5.122,5.122a.514.514,0,0,1-.67.777l-.058-.05L6,6.728.878,11.849A.514.514,0,0,1,.1,11.18l.05-.058L5.272,6,.151.878A.514.514,0,0,1,.75.057Z" transform="translate(2 2)" fill="#1e1e1e"></path>  ' +
//         '      </svg>  ' +
//         '  </button></div><div class="modal-body p-4 pb-0">  ' +
//         '  <div class="header-history">  ' +
//         '      <div class="header-h">  ' +
//         '          <span class="pr-2">Show</span> <span class="">  ' +
//         '              <label class="" for="is_default"> </label>  ' +
//         '              <div class="form-check form-check-inline form-check-rounded">  ' +
//         '                  <label for="Everything">  ' +
//         '                      Everything  ' +
//         '                  </label>  ' +
//         '                  <input class="form-check-input loadType" type="radio" value="Everything" name="loadType" id="Everything">  ' +
//         '              </div>  ' +
//         '          </span> <span class="">  ' +
//         '              <label class="" for="is_default"> </label>  ' +
//         '              <div class="form-check form-check-inline form-check-rounded">  ' +
//         '                  <label for="ChangesOnly">  ' +
//         '                      Changes only  ' +
//         '                  </label>  ' +
//         '                  <input class="form-check-input loadType" type="radio" value="ChangesOnly" name="loadType" checked id="ChangesOnly">  ' +
//         '              </div>  ' +
//         '          </span>  ' +
//         '      </div>  ' +
//         '  </div>  ' +
//         '  <div class="table-responsive border" style="max-height: 500px">  ' +
//         '      <table class="table table-hover mb-0 table-history text-white">  ' +
//         '          <thead>  ' +
//         '          <tr>  ' +
//         '              <th scope="col" style=" min-width: 300px;">Field</th>  ' +
//         '              <th scope="col"></th>  ' +
//         '              <th scope="col" style="min-width: 120px;">Previous state</th>  ' +
//         '              <th scope="col">Current state</th>  ' +
//         '          </tr>  ' +
//         '          </thead>  ' +
//         '          <tbody>  ' +
//         '              <tr class="item-' + history.ID + ' title-item">  ' +
//         '                  <td colspan="4"><b>' + history.Title + ' on ' + createTime + '</b></td>  ' +
//         '              </tr>  ' +
//         compaseHistory(Compare, history) +
//         '          </tbody>  ' +
//         '      </table>  ' +
//         '  </div>  ' +
//         '  </div><div class="modal-footer border-top-0">  ' +
//         '  <button type="button" class="btn px-2" data-bs-dismiss="modal">  ' +
//         '      <span class="px-1 text-primary">CLOSE</span>  ' +
//         '  </button></div>';
//     return html;
// }

function renderHtmlHistory(Compare, history, createTime) {
    var html = '<div class="modal-header border-0 pb-0 align-items-start px-4">  ' +
        '  <h5 class="modal-title text-white" id="exampleModalLabel">History: ' + history.ObjectName + ' ' + history.Title + '</h5>  ' +
        '  </button></div><div class="modal-body p-4 pb-0">  ' +
        '  <div class="header-history">  ' +
        '      <div class="header-h">  ' +
        '          <span class="pr-2">Show</span> <span class="">  ' +
        '              <label class="" for="is_default"> </label>  ' +
        '              <div class="form-check form-check-inline form-check-rounded">  ' +
        '                  <label for="Everything">  ' +
        '                      Everything  ' +
        '                  </label>  ' +
        '                  <input class="form-check-input loadType" type="radio" value="Everything" name="loadType" id="Everything">  ' +
        '              </div>  ' +
        '          </span> <span class="">  ' +
        '              <label class="" for="is_default"> </label>  ' +
        '              <div class="form-check form-check-inline form-check-rounded">  ' +
        '                  <label for="ChangesOnly">  ' +
        '                      Changes only  ' +
        '                  </label>  ' +
        '                  <input class="form-check-input loadType" type="radio" value="ChangesOnly" name="loadType" checked id="ChangesOnly">  ' +
        '              </div>  ' +
        '          </span>  ' +
        '      </div>  ' +
        '  </div>  ' +
        '  <div class="table-responsive" style="overflow: auto;">  ' +
        '      <table class="table table-hover table-striped border-0 mb-0 table-history">  ' +
        '          <thead>  ' +
        '          <tr>  ' +
        '              <th scope="col" style=" min-width: 300px;">Field</th>  ' +
        '              <th scope="col"></th>  ' +
        '              <th scope="col" style="min-width: 120px;">Previous state</th>  ' +
        '              <th scope="col">Current state</th>  ' +
        '          </tr>  ' +
        '          </thead>  ' +
        '          <tbody>  ' +
        '              <tr class="item-' + history.ID + ' title-item">  ' +
        '                  <td colspan="4"><b>' + history.Title + ' on ' + createTime + '</b></td>  ' +
        '              </tr>  ' +
        compaseHistory(Compare, history) +
        '          </tbody>  ' +
        '      </table>  ' +
        '  </div>  ' +
        '  </div>';
    return html;
}

function compaseHistory(Compare, history) {
    if (Compare.length == 0) {
        return "";
    }

    var html = "";

    Compare.forEach(function (value, index) {
        var Action = "";
        if (value.Action == "none") {
            Action = 'everything';
        }
        var icon = "";
        switch (value.Action) {
            case "add":
                icon = '<span data-bs-toggle="tooltip" data-bs-placement="left" title="Add">  ' +
                    '                 <i class="fa fa-plus-circle text-success " aria-hidden="true"></i>  ' +
                    '            </span>';
                break;
            case "update":
                icon = '<span data-bs-toggle="tooltip" data-bs-placement="left" title="Update">  ' +
                    '                 <i class="fa fa-refresh" aria-hidden="true"></i>  ' +
                    '            </span>';
                break;
            case "delete":
                icon = '<span data-bs-toggle="tooltip" data-bs-placement="left" title="Delete">  ' +
                    '                 <i class="fa fa-minus-circle text-danger" aria-hidden="true"></i>  ' +
                    '            </span>';
                break;
            default:
                icon = '<i class="fa fa-circle-o" aria-hidden="true" style="opacity: 0.3;"></i>';
                break;
        }
        var oldData = "";
        if (value.OldData == "") {
            oldData = '<span class="opacity-50">N/A</span>';
        } else {
            oldData = htmlEntities(value.OldData);
        }
        var newData = "";
        if (value.NewData == "") {
            newData = '<span class="opacity-50">N/A</span>';
        } else {
            newData = htmlEntities(value.NewData);
        }
        html = html +
            '<tr class="value-item ' + Action + ' value-' + history.ID + '" data-id="' + history.ID + '">  ' +
            '    <td class="text-nowrap">' + value.Text + '</td>  ' +
            '    <td>' + icon + '</td>  ' +
            '    <td>' + oldData + '</td>' +
            '    <td>' + newData + '</td>' +
            '</tr>  ';
    });
    return html;
}

function FiltersHistory(flag = true) {
    const formElement = $(formID);
    let buttonElement = formElement.find(".submit");
    let submitButtonText = buttonElement.text();
    let submitButtonTextLoading = "Loading...";
    let postData = formElement.serializeObject();
    postData.page = parseInt(postData.page) + 1;
    postData.start_date = $("#input_start_date").val();
    postData.end_date = $("#input_end_date").val();

    $.ajax({
        url: filterURL,
        type: "POST",
        dataType: "JSON",
        contentType: "application/json; charset=utf-8",
        data: JSON.stringify(postData),
        beforeSend: function (xhr) {
            // xhr.overrideMimeType("text/plain; charset=x-user-defined");
        },
        error: function (jqXHR, exception) {
            var msg = AjaxErrorMessage(jqXHR, exception);
            new AlertError("AJAX ERROR: " + msg);
        }
    }).done(function (result) {
        switch (result.status) {
            case true:
                if (result.data.records.length > 0) {
                    if (postData.page < 2) {
                        $('.list-history').html("");
                    }
                    var date = "";
                    result.data.records.forEach(function (value, index) {
                        if (date != new Date(value.CreatedAt).toLocaleDateString('en-us', {
                            weekday: "long",
                            year: "numeric",
                            month: "short",
                            day: "numeric"
                        })) {
                            date = new Date(value.CreatedAt).toLocaleDateString('en-us', {
                                weekday: "long",
                                year: "numeric",
                                month: "short",
                                day: "numeric"
                            });
                            $('.list-history').append('<div class="list-group-item px-md-4 border-right" data-id="2117"><div class="row"> <div class="col-auto"><span class="font-weight-bolt" data-bs-toggle="tooltip" data-bs-placement="left" title="" data-bs-original-title="Add" aria-label="Add" style="font-weight: 700;">'
                                + date + '\n' +
                                '</span></div>   </div></div>');
                        }
                        $('.list-history').append(historyItem(value));
                        console.log(flag);
                        if (flag) {
                            loadHistory(value.ID)
                            flag = false;
                        }
                    });
                } else {
                    $('.list-history').html("<div class=\"list-group-item px-md-4\" data-id=\"\">" +
                        "                        <div class=\"row px-3\">" +
                        "                           No Activity" +
                        "                       </div>" +
                        "                   </div>");
                }

                if (result.data.totalRecord <= postData.page * 30) {
                    $(".load-more").closest("div").addClass("d-none");
                } else {
                    $(".load-more").closest("div").removeClass("d-none");
                }
                $('[data-bs-toggle="tooltip"]').tooltip();
                // NoticeSuccess("Bidder has been removed successfully");
                break;
            case false:
                if (result.errors.length > 0) {
                    result.errors.forEach(function (value, index) {
                        $('.list-history').html("").append("<div class=\"list-group-item px-md-4\" data-id=\"\">" +
                            "                        <div class=\"row px-3\">" +
                            value.message +
                            "                       </div>" +
                            "                   </div>");
                    })
                }
                break;
            case "error":
                console.log(result);
            // new AlertError(result.message);
        }
        $(".load-more").text("Load more")
        $(formID).find(".submit").text("Run")
    });
}

function historyItem(history) {
    const time = new Date(history.CreatedAt);

    var icon = "";
    if (history.ObjectType == 1) {
        icon = '<span data-bs-toggle="tooltip" data-bs-placement="left" title="Add">' +
            '<i class="fa fa-plus-circle text-success " aria-hidden="true"></i>' +
            '</span>';
    } else if (history.ObjectType == 2) {
        icon = '<span data-bs-toggle="tooltip" data-bs-placement="left" title="Update">' +
            '<i class="fa fa-refresh" aria-hidden="true"></i>' +
            '</span>';
    } else if (history.ObjectType == 3) {
        icon = '<span data-bs-toggle="tooltip" data-bs-placement="left" title="Delete">' +
            '<i class="fa fa-minus-circle text-danger" aria-hidden="true"></i>' +
            '</span>';
    } else {
        icon = '<i class="fa fa-circle-o" aria-hidden="true" style="opacity: 0.3;"></i>';
    }

    var html = '<div class="list-group-item px-md-4 border-right" data-id="' + history.ID + '">' +
        '<div class="row px-3">' +
        ' <div class="col-auto">' + icon + '</div>' +
        '   <div class="col ps-0">' +
        '   <span class="mb-2 d-block text-gray-800 load-history" role="button">' + history.ObjectName + ' ' + history.Title + '</span>' +
        // '   <span class="mb-2 d-block text-gray-800 load-history" role="button" data-bs-toggle="modal" data-bs-target=".modal-history">' + history.ObjectName + ' ' + history.Title + '</span>' +
        '<p class="card-text text-gray-600 lh-sm">' + time + '</p>' +
        '</div>' +
        // '<div class="col-auto">' +
        // '   <a href="#" class="btn btn-dark-100 btn-icon btn-sm rounded-circle load-history" role="button" data-bs-toggle="modal" data-bs-target=".modal-history">' +
        // '   <svg width="48" height="48" viewBox="0 0 16 16" class="bi bi-code" fill="currentColor" xmlns="http://www.w3.org/2000/svg">' +
        // '   <path fill-rule="evenodd" d="M5.854 4.146a.5.5 0 0 1 0 .708L2.707 8l3.147 3.146a.5.5 0 0 1-.708.708l-3.5-3.5a.5.5 0 0 1 0-.708l3.5-3.5a.5.5 0 0 1 .708 0zm4.292 0a.5.5 0 0 0 0 .708L13.293 8l-3.147 3.146a.5.5 0 0 0 .708.708l3.5-3.5a.5.5 0 0 0 0-.708l-3.5-3.5a.5.5 0 0 0-.708 0z"></path>' +
        // '   </svg>' +

        // '<svg width="48" height="48" viewBox="0 0 16 16" class="bi bi-folder2-open" fill="currentColor" xmlns="http://www.w3.org/2000/svg">\n' +
        // '                        <path fill-rule="evenodd" d="M1 3.5A1.5 1.5 0 0 1 2.5 2h2.764c.958 0 1.76.56 2.311 1.184C7.985 3.648 8.48 4 9 4h4.5A1.5 1.5 0 0 1 15 5.5v.64c.57.265.94.876.856 1.546l-.64 5.124A2.5 2.5 0 0 1 12.733 15H3.266a2.5 2.5 0 0 1-2.481-2.19l-.64-5.124A1.5 1.5 0 0 1 1 6.14V3.5zM2 6h12v-.5a.5.5 0 0 0-.5-.5H9c-.964 0-1.71-.629-2.174-1.154C6.374 3.334 5.82 3 5.264 3H2.5a.5.5 0 0 0-.5.5V6zm-.367 1a.5.5 0 0 0-.496.562l.64 5.124A1.5 1.5 0 0 0 3.266 14h9.468a1.5 1.5 0 0 0 1.489-1.314l.64-5.124A.5.5 0 0 0 14.367 7H1.633z"></path>\n' +
        // '                      </svg>' +
        // '   </a>' +
        // '</div>' +
        '</div>' +
        '</div>';
    return html;
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

function sm(start, end, label) {
    $('#input_start_date').val(start.format('YYYY-MM-DD'));
    $('#input_end_date').val(end.format('YYYY-MM-DD'));
//        $('.result-search-time').text(moment().subtract(start, 'DD/MM/YYYY'));
//        $('.result-search-time').text(start.format('YYYYMMDD') + ' - ' + end.format('YYYYMMDD'));
    if (label == 'Today' || label == 'Yesterday') {
        $('.result-search-time').text(start.format('MMM D, YYYY'));
    } else if (label == 'All Time') {
        $('.result-search-time').text(label);
    } else {
        $('.result-search-time').text(start.format('MMM D, YYYY') + ' - ' + end.format('MMM D, YYYY'));
    }
    FiltersHistory()
}

function formatDate(date) {
    var d = new Date(date),
        month = '' + (d.getMonth() + 1),
        day = '' + d.getDate(),
        year = d.getFullYear();

    if (month.length < 2)
        month = '0' + month;
    if (day.length < 2)
        day = '0' + day;

    return [year, month, day].join('-');
}