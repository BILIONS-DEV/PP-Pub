$(document).ready(function () {
    SubmitFormForget("forgetPassword", "/user/forgot-password", Callback)
});


function SubmitFormForget(formID, ajaxUrl, callback) {
    let formElement = $("#" + formID);
    formElement.find("input").on("click change blur", function (e) {
        let inputElement = $(this)
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next("div.invalid-feedback").empty();
        }
    });
    formElement.find(".submit").on("click", function (e) {
        e.preventDefault();
        const buttonElement = $(this);
        const submitButtonText = buttonElement.text();
        const submitButtonTextLoading = "Loading...";
        let postData = formElement.serializeObject();
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
                // new AlertError("AJAX ERROR: " + msg);
                buttonElement.attr('disabled', false).text(submitButtonText);
            },
            success: function (responseJSON) {
                buttonElement.attr('disabled', false).text(submitButtonText);
            },
            complete: function (res) {
                callback(res.responseJSON, formElement);
            }
        });
    });
}

function Callback(response, formElement) {
    switch (response.status) {
        case "error":
            if (response.errors.length > 0) {
                $.each(response.errors, function (key, value) {
                    let inputElement = $("input#" + value.id);
                    if (key === 0) {
                        inputElement.select().focus();
                    }
                    inputElement.addClass("is-invalid").nextAll(".invalid-feedback").text(value.message);
                });
            }
            new AlertError(response.errors[0].message);
            break
        case "success":
            let buttonElement = formElement.find(".submit");
            buttonElement.remove()
            let div = formElement.find("#send-link")
            div.attr("hidden", false)
            let email = formElement.find("#email")
            email.attr("disabled", true)
            break
        default:
            new AlertError("undefined");
            break
    }
}