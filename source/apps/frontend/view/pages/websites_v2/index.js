import dataTable from "../../jspkg/datatable";

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

const formID = "#formFilterInventory";
const tableID = "#tableInventories";

$(document).ready(function () {
    /**
     * Load Datatable
     * @type {string}
     */
    $(formID).find(".submit").on("click", function (e) {
        dataTable.Render(formID, true, true);
    });
    dataTable.Render(formID);

    /**
     * Delete Inventory
     */
    $(tableID).on("click", "tbody td div.btn-group a.remove", function (e) {
        const id = $(this).data("id");
        const isCheck = confirm('Are you sure delete!"');
        if (isCheck) {
            Delete(id);
        }
    });

    /**
     * Copy AdTag
     */
    $(tableID).on("click", "a.copy-tag", function (e) {
        const id = $(this).data("id");
        CopyTag(id);
    });

    /**
     * Submit Inventory
     */
    $("input#inventories").on("input", function (e) {
        let text = e.target.value;
        if (text !== "") {
            $(".save-domain").attr("disabled", false);
        } else {
            $(".save-domain").attr("disabled", true);
        }
    });
    SubmitFormInventory("formSubmitInventory", SubmitInventoryResponse, "/websites/submit");
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
                url: "/websites/submit",
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
            // SubmitFormInventory("formSubmitInventory", SubmitInventoryResponse, "/websites/submit");
            return false;
        }
    });

});

function SubmitFormInventory(formID, functionCallback, url = "") {
    const formElement = $("#" + formID);
    formElement.find(".submit").on("click", function (e) {
        e.preventDefault();
        let validate = true;
        const buttonElement = $(this);
        const submitButtonText = buttonElement.text();
        const submitButtonTextLoading = "Loading...";
        formElement.find("input").on("click change blur", function (e) {
            let inputElement = $(this);
            if (inputElement.hasClass("is-invalid")) {
                inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
            }
        });
        formElement.find("textarea").on("click change blur", function (e) {
            let inputElement = $(this);
            if (inputElement.hasClass("is-invalid")) {
                inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
            }
        });
        const postData = formElement.serializeObject();
        if (postData.inventories.trim() == "") {
            $("#formSubmitInventory").find(".invalid-feedback").text("(*) Required").parent().addClass("is-invalid").find("input").addClass("is-invalid");
            return;
        }
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
        return;
    }
    let ulElement = $("ul#respSubmit").empty();
    $.each(response, function (key, obj) {
        let liElement;
        if (obj.message === "") {
            // liElement = liSuccessTemp.replace("{{domain}}", obj.id);
            // ulElement.prepend(liElement)
            // $("#submitDomainModalFullscreen").modal('hide');
            $("#inventories").val("");
            NoticeSuccess("Your domain name has been submitted successfully");
            setTimeout(function () {
                $("#submitDomainModalFullscreen").modal('hide');
            }, 1000);
        } else {
            // console.log(obj)
            // liElement = liErrorTemp.replace("{{domain}}", obj.id).replace("{{message}}", obj.message);
            // ulElement.prepend(liElement)
            AlertError(obj.message);
        }
    });
    setTimeout(function () {
        $("body").removeClass("customize-box");
    }, 200);
    // GetTable(false);
    dataTable.Render(formID, false, true);
}

function CopyTag(id) {
    let url = "/websites/copyAdTag";
    // $(`#copy-adtag-modal`).css("display", "block")
    $.ajax({
        url: url,
        type: "GET",
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
    let url = "/websites/del";
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
                NoticeSuccess("Inventory has been removed successfully");
                dataTable.Render(formID, true, false);
                // GetTable(false)
                break;
            case "err":
                new AlertError(result.message);
        }
    });
}