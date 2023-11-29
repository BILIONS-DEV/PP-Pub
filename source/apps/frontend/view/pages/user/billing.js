$(document).ready(function () {
    $("#method").on("change", function () {
        ChangeSelectBilling($(this))
    })
    GetMethodDefaultBilling()
    SumitFormBilling("submitBilling", "/user/billing", Submit)
})

function GetMethodDefaultBilling() {
    let method = $('select[name="method"] option').filter(':selected').val()
    ShowHideByMethod(method)
}

function ChangeSelectBilling(element) {
    let selectValue = element.val();
    ShowHideByMethod(selectValue)
}

function ShowHideByMethod(method) {
    switch (method) {
        case "bank":
            $("#method_bank").removeClass("d-none")
            $("#method_payoneer").addClass("d-none")
            $("#method_paypal").addClass("d-none")
            break
        case "paypal":
            $("#method_bank").addClass("d-none")
            $("#method_payoneer").addClass("d-none")
            $("#method_paypal").removeClass("d-none")
            break
        case "payoneer":
            $("#method_bank").addClass("d-none")
            $("#method_payoneer").removeClass("d-none")
            $("#method_paypal").addClass("d-none")
            break
        default:
            $("#method_bank").removeClass("d-none")
            $("#method_payoneer").addClass("d-none")
            $("#method_paypal").addClass("d-none")
            break
    }
}

function SumitFormBilling(formID, ajaxUrl, callback) {
    let formElement = $("#" + formID);
    formElement.find("input").on("click change blur", function (e) {
        let inputElement = $(this)
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    formElement.find(".submit").on("click", function (e) {
        e.preventDefault();
        const buttonElement = $(this);
        const submitButtonText = buttonElement.text();
        const submitButtonTextLoading = "Loading...";
        let postData = formElement.serializeObject();
        postData.id = parseInt(postData.id)
        postData.user_id = parseInt(postData.user_id)
        $.ajax({
            url: ajaxUrl,
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
                callback(res.responseJSON);
            }
        });
    });
}

function Submit(response) {
    switch (response.status) {
        case "error":
            if (response.errors.length === 1 && response.errors[0].id === "") {
                new AlertError(response.errors[0].message);
            } else if (response.errors.length > 0) {
                $.each(response.errors, function (key, value) {
                    let inputElement = $("#" + value.id);
                    if (key === 0) {
                    }
                    inputElement.addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
                });
                new AlertError(response.errors[0].message);
            } else {
                new AlertError("Error!");
            }
            break
        case "success":
            NoticeSuccess("Billing has been updated successfully")
            break
        default:
            new AlertError("Undefined");
            break
    }
}
