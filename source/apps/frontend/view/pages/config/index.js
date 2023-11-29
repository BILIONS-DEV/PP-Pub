$(document).ready(function () {
    SubmitFormSaveConfig("formConfig", Saved, "/config/save");
    let currency = $("#isSetCurrency").val()
    // if (currency === "") {
    //     let message = $("#isSetCurrency").data("message")
    //     ShowAlert(message)
    // }

    var elmnt = $("#ConfigBox").find(".select2-container")[0]
    $("#ConfigBox").find(".media").css("max-width", 250 + elmnt.offsetWidth + "px")
})

function ShowAlert(message) {
    let alertModal = swal("Oops...", message, {
        icon: "warning",
        dangerMode: true,
        html: true,
        buttons: {
            confirm: {
                className: 'btn btn-lg btn-warning'
            }
        }
    });
}

function SubmitFormSaveConfig(formID, functionCallback, ajxURL = "") {
    let formElement = $("#" + formID);
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
    formElement.find(".select2").on("select2:open", function () {
        let selectElement = $(this)
        if (selectElement.next().find(".select2-selection").hasClass("select2-is-invalid")) {
            selectElement.next().find(".select2-selection").removeClass("select2-is-invalid")
            selectElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    })
    formElement.find(".submit").on("click", function (e) {
        e.preventDefault();
        const buttonElement = $(this);
        const submitButtonText = buttonElement.text();
        const submitButtonTextLoading = "Loading...";
        var postData = formElement.serializeObject();
        postData.ad_refresh_time = parseInt(postData.ad_refresh_time)
        postData.prebid_time_out = parseInt(postData.prebid_time_out)
        postData.currency = parseInt(postData.currency)
        $.ajax({
            url: ajxURL,
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

function Saved(response, formElement) {
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
                    inputElement.next().find(".select2-selection").addClass("select2-is-invalid")
                });
                $("#" + response.errors[0].id).focus();
                $("#" + response.errors[0].id).prev('label').focus();
                // new AlertError(response.errors[0].message);
            } else {
                new AlertError("Error!");
            }
            break
        case "success":
            NoticeSuccess("Config has been saved successfully");
            setTimeout(function(){
                window.location.href = window.location.origin + "/supply"
            }, 1000)
            // location.reload();
            break
        default:
            new AlertError("Undefined");
            break
    }
}