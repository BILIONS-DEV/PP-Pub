let liSuccessTemp = `<li>
    <span class="avatar avatar-xs bg-teal-50 rounded-circle d-flex align-items-center justify-content-center">
        <img src="/static/svg/icons/right-check2.svg" width="14" height="14">
    </span>
    <span>
        {{domain}}
    </span>
</li>`;
let liErrorTemp = `<li>
    <span class="avatar avatar-xs bg-red-50 rounded-circle d-flex align-items-center justify-content-center">
        <img src="/static/svg/icons/close-danger.svg" width="14" height="14">
    </span>
    <span>{{domain}} <small class="text-muted">{{message}}</small></span>
</li>`;

const formID = "#formFilterInventory"
const tableID = "#tableInventories"
const filterURL = "/supply";
const currentURL = "/supply";
let firstLoad = true;

$(document).ready(function () {
    SubmitFormInventory("formSubmitInventory", SubmitInventoryResponse, "/supply/submit");
    $(formID).find(".submit").on("click", function (e) {
        e.preventDefault();
        GetTable(true);
    });
    GetTable(false);

    $(tableID).on("click", "tbody td div.btn-group a.remove", function (e) {
        var id = $(this).data("id")
        var isCheck = confirm('Are you sure delete!"');
        if (isCheck) {
            Delete(id)
        }
    });

    $(tableID).on("click", "a.copy-tag", function (e) {
        const id = $(this).data("id");
        CopyTag(id)
    });

    // window.onclick = function (event) {
    //     if (event.target == modal[0]) {
    //         $(`#copy-adtag-modal`).css("display", "none")
    //     }
    // }

    // $(".close").on("click", function () {
    //     $(`#copy-adtag-modal`).css("display", "none")
    // })

    // $('#exampleModalFullscreen').on('shown.bs.modal', function () {
    //     $('input#inventories').focus();
    // })
    $("input#inventories").on("input", function (e) {
        let text = e.target.value
        if (text !== "") {
            $(".save-domain").attr("disabled", false)
        } else {
            $(".save-domain").attr("disabled", true)
        }
    })

    $("#formSubmitInventory").keypress(function (e) {
        const key = e.which;
        if (key === 13) { // the enter key code
            e.preventDefault();
            const buttonElement = $(".save-domain");
            const submitButtonText = buttonElement.text();
            const submitButtonTextLoading = "Loading...";
            var postData = $(this).serializeObject();
            // console.log(postData)
            $.ajax({
                url: "/supply/submit",
                type: "POST",
                dataType: "JSON",
                contentType: "application/json",
                data: JSON.stringify(postData),
                beforeSend: function (xhr) {
                    buttonElement.attr('disabled', true).text(submitButtonTextLoading);
                },
                error: function (jqXHR, exception) {
                    const msg = AjaxErrorMessage(jqXHR, exception);
                    new AlertError("AJAX ERROR: " + msg);
                    buttonElement.attr('disabled', false).text(submitButtonText);
                },
                success: function (responseJSON) {
                    buttonElement.attr('disabled', false).text(submitButtonText);
                },
                complete: function (res) {
                    SubmitInventoryResponse(res.responseJSON, this);
                }
            });
            // SubmitFormInventory("formSubmitInventory", SubmitInventoryResponse, "/supply/submit");
            return false;
        }
    });

});

function SubmitFormInventory(formID, functionCallback, url = "") {
    const formElement = $("#" + formID);
    formElement.find(".submit").on("click", function (e) {
        e.preventDefault();
        let validate = true
        const buttonElement = $(this);
        const submitButtonText = buttonElement.text();
        const submitButtonTextLoading = "Loading...";
        formElement.find("input").on("click change blur", function (e) {
            let inputElement = $(this)
            if (inputElement.hasClass("is-invalid")) {
                inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
            }
        });
        formElement.find("textarea").on("click change blur", function (e) {
            let inputElement = $(this)
            if (inputElement.hasClass("is-invalid")) {
                inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
            }
        });
        const postData = formElement.serializeObject();
        // console.log(postData)
        $.ajax({
            url: url,
            type: "POST",
            dataType: "JSON",
            contentType: "application/json",
            data: JSON.stringify(postData),
            beforeSend: function (xhr) {
                buttonElement.attr('disabled', true).text(submitButtonTextLoading);
            },
            error: function (jqXHR, exception) {
                const msg = AjaxErrorMessage(jqXHR, exception);
                new AlertError("AJAX ERROR: " + msg);
                buttonElement.attr('disabled', false).text(submitButtonText);
            },
            success: function (responseJSON) {
                buttonElement.attr('disabled', false).text(submitButtonText);
            },
            complete: function (res) {
                functionCallback(res.responseJSON, formElement);
            }
        });
    });
}

/**
 *
 * @param response
 * @param formElement
 * @constructor
 */
function SubmitInventoryResponse(response, formElement) {
    if (response === undefined || response === null) {
        return
    }
    let ulElement = $("ul#respSubmit").empty();
    $.each(response, function (key, obj) {
        let liElement;
        if (obj.message === "") {
            // liElement = liSuccessTemp.replace("{{domain}}", obj.id);
            // ulElement.prepend(liElement)
            // $("#submitDomainModalFullscreen").modal('hide');
            $("#inventories").val("")
            NoticeSuccess("Your domain name has been submitted successfully")
            setTimeout(function () {
                $("#submitDomainModalFullscreen").modal('hide')
            }, 1000);
        } else {
            // console.log(obj)
            // liElement = liErrorTemp.replace("{{domain}}", obj.id).replace("{{message}}", obj.message);
            // ulElement.prepend(liElement)
            AlertError(obj.message)
        }
    });
    GetTable(false);
}

function GetTable(isClickForm = false) {
    const formElement = $(formID);
    let buttonElement = formElement.find(".submit");
    let submitButtonText = buttonElement.text();
    let submitButtonTextLoading = "Loading...";
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
        order: [],
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
                    buttonElement.attr("disabled", true).text(submitButtonTextLoading);
                }
            },
            dataSrc: function (json) {
                if (isClickForm) {
                    buttonElement.attr("disabled", false).text(submitButtonText);
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
                    buttonElement.attr("disabled", false).text(submitButtonText);
                }
            },
        },
        columns: [
            {data: "name", name: "Domain"},
            {data: "status", name: "Status"},
            {data: "sync_ads_txt", name: "Ads.txt in Sync"},
            {data: "live", name: "Website Live"},
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

function CopyTag(id) {
    let url = "/supply/copyAdTag"
    // $(`#copy-adtag-modal`).css("display", "block")
    $.ajax({
        url: url,
        type: "GET",
        dataType: "JSON",
        contentType: "application/json",
        data: {id: id},
        beforeSend: function (xhr) {
        },
        error: function (jqXHR, exception) {
            const msg = AjaxErrorMessage(jqXHR, exception);
            new AlertError("AJAX ERROR: " + msg);
        }
    }).done(function (result) {
        $('#copyAdTagModal .modal-body').html(result);
    });
}

function Delete(id) {
    let url = "/supply/del"
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
                NoticeSuccess("Inventory has been removed successfully")
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