const ajaxResponse = require("./ajax-response");
module.exports = {
    SubmitForm, Post, PrintResponse
}

function SubmitForm(formID, functionCallback, ajxURL = "") {
    let formElement = $("#" + formID);
    let button = formElement.find(".submit");
    let submitButtonText = button.text();
    let submitButtonTextLoading = "Loading...";
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
    formElement.on('keypress', function (e) {
        if (e.which === 13 && !button.is(":disabled") && formElement.find("input").is(':focus')) {
            button.click();
        }
    });
    button.on("click", function (e) {
        e.preventDefault();
        $.each(formElement.find("input.is-invalid"), function () {
            $(this).removeClass("is-invalid").next(".invalid-feedback").empty();
        });
        const postData = formElement.serializeArray();
        const buttonElement = $(this);
        $.ajax({
            url: ajxURL,
            type: "POST",
            dataType: "JSON",
            data: postData,
            beforeSend: function (xhr) {
                buttonElement.attr('disabled', true).text(submitButtonTextLoading);
            },
            error: function (jqXHR, exception) {
                const msg = AjaxErrorMessage(jqXHR, exception);
                ajaxResponse.ErrorWithAlert("AJAX ERROR: " + msg)
                buttonElement.attr('disabled', false).text(submitButtonText)
            },
            success: function (responseJSON) {
                buttonElement.attr('disabled', false).text(submitButtonText);
            },
            complete: function (resp) {
                if (functionCallback !== undefined && IsFunction(functionCallback)) {
                    functionCallback(resp.responseJSON, formElement);
                } else {
                    PrintResponse(resp.responseJSON)
                }
            }
        });
    });
}

function Post(formElement, postData, functionCallback) {
    let url = formElement.attr("action")
    let buttonElement = formElement.find(".submit")
    let textButtonLoading = buttonElement.data("text-loading") || "Loading..."
    let textButtonSubmit = buttonElement.text()
    $.ajax({
        url: url,
        data: postData,
        type: "POST",
        dataType: "JSON",
        contentType: "application/json",
        beforeSend: function (xhr) {
            buttonElement.attr('disabled', true).text(textButtonLoading);
        },
        error: function (jqXHR, exception) {
            const msg = AjaxErrorMessage(jqXHR, exception);
            new AlertError("AJAX ERROR: " + msg);
            buttonElement.attr('disabled', false).text(textButtonSubmit);
        },
        success: function (responseJSON) {
            buttonElement.attr('disabled', false).text(textButtonSubmit);
        },
        complete: function (resp) {
            if (functionCallback !== undefined && IsFunction(functionCallback)) {
                functionCallback(resp.responseJSON, formElement);
            } else {
                PrintResponse(resp.responseJSON)
            }
        }
    });
}


function PrintResponse(resp) {
    switch (resp.status) {
        case "error":
            if (resp.errors.length === 1 && resp.errors[0].id === "") {
                ajaxResponse.ErrorWithAlert(resp.errors[0].message)
            } else if (resp.errors.length > 0) {
                $.each(resp.errors, function (key, value) {
                    let inputElement = $("#" + value.id);
                    inputElement.addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
                    inputElement.next().find(".select2-selection").addClass("select2-is-invalid")
                    inputElement.closest('.box-select2').addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message)
                    if (value.id === "text_for_domain") {
                        let msg = value.message
                        inputElement.html(`<small class="text-danger">${msg}</small>`)
                        $('.domain-card').addClass("domain-card-invalid")
                    }
                });
                $("#" + resp.errors[0].id).focus();
                $("#" + resp.errors[0].id).prev('label').focus();
            } else {
                ajaxResponse.ErrorWithAlert("Error!!!")
            }
            break

        case "success":
            ajaxResponse.DoneWithNotify(resp.message)
            break

        default:
            ajaxResponse.ErrorWithAlert("Undefined!!!")
            break
    }
}