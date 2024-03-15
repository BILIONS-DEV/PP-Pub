$(document).ready(function () {
    SumitFormAccount("profile", "/user/account", Submit)
    SumitFormAccount("changePassword", "/user/changePassword", Submit)
    SumitFormBilling("submitBilling", "/user/billing", Submit);

        $("#method").on("change", function () {
        ChangeSelectBilling($(this))
    })
    ChangeTab()
    GetMethodDefaultBilling()
})

function SumitFormAccount(formID, ajaxUrl, callback) {
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
                new AlertError(response.errors[0].message);
                $.each(response.errors, function (key, value) {
                    let inputElement = $("#" + value.id);
                    if (key === 0) {
                    }
                    inputElement.addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
                });
            } else {
                new AlertError("Error!");
            }
            break
        case "success":
            NoticeSuccess("Account has been updated successfully")
            break
        default:
            new AlertError("Undefined");
            break
    }
}

// ***********************  Billing  ****************************
function GetMethodDefaultBilling() {
    let method = $('select[name="method"] option').filter(':selected').val()
    ShowHideByMethod(method)
}

function ChangeSelectBilling(element) {
    let selectValue = element.val();
    ShowHideByMethod(selectValue)
    ChangeMinPayCcurrency(element)
}

function ShowHideByMethod(method) {
    switch (method) {
        case "bank":
            $("#method_bank").removeClass("d-none");
            $("#method_payoneer").addClass("d-none");
            $("#method_paypal").addClass("d-none");
            $("#method_currency").addClass("d-none");
            break;
        case "paypal":
            $("#method_bank").addClass("d-none");
            $("#method_payoneer").addClass("d-none");
            $("#method_paypal").removeClass("d-none");
            $("#method_currency").addClass("d-none");
            break;
        case "payoneer":
            $("#method_bank").addClass("d-none");
            $("#method_payoneer").removeClass("d-none");
            $("#method_paypal").addClass("d-none");
            $("#method_currency").addClass("d-none");
            break;
        case "currency":
            $("#method_currency").removeClass("d-none");
            $("#method_bank").addClass("d-none");
            $("#method_payoneer").addClass("d-none");
            $("#method_paypal").addClass("d-none");
            break;
        case  "BTC":
            $("#method_currency").removeClass("d-none");
            $("#method_bank").addClass("d-none");
            $("#method_payoneer").addClass("d-none");
            $("#method_paypal").addClass("d-none");
            break;
        case  "BCH":
            $("#method_currency").removeClass("d-none");
            $("#method_bank").addClass("d-none");
            $("#method_payoneer").addClass("d-none");
            $("#method_paypal").addClass("d-none");
            break;
        case   "USDT":
            $("#method_currency").removeClass("d-none");
            $("#method_bank").addClass("d-none");
            $("#method_payoneer").addClass("d-none");
            $("#method_paypal").addClass("d-none");
            break;
        default:
            $("#method_bank").removeClass("d-none");
            $("#method_payoneer").addClass("d-none");
            $("#method_paypal").addClass("d-none");
            $("#method_currency").addClass("d-none");
            break;
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

function ChangeTab() {
    $(".custom-tab").on("click", ".nav-link", function () {
        $(".custom-tab").find(".nav-link").removeClass("pp-4").removeClass('at-1')
        $(this).addClass("pp-4")
        var tab = $(this).attr("data-tab")
        if (tab != "1") {
            $(this).addClass("at-1")
        }
    })

    $("#tab-profile").on("click", function () {
        $("#tab-profile").addClass("active")
        $("#tab-change-password").removeClass("active")
        $("#tab-billing").removeClass("active")
    })
    $("#tab-change-password").on("click", function () {
        $("#tab-profile").removeClass("active")
        $("#tab-change-password").addClass("active")
        $("#tab-billing").removeClass("active")
    })
    $("#tab-billing").on("click", function () {
        $("#tab-profile").removeClass("active")
        $("#tab-change-password").removeClass("active")
        $("#tab-billing").addClass("active")
    })
}

function ChangeMinPayCcurrency(element) {
    let selectValue = element.val();
    $(".currency-minpay").addClass("d-none");
    $(".currency-" + selectValue).removeClass("d-none");
}