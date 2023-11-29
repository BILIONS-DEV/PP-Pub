$(document).ready(function () {
    // CaptilizeBidder()
    SubmitFormCreate("formAdvertisingSchedules", Added);

    checkSelectAdClient();
});

function checkSelectAdClient() {
   function check() {
       let adClient = $('input[name="ad_client"]:checked').val();
       if (adClient === "1"){
           $(".box_vpaid_mode").removeClass("d-none");
       } else {
           $(".box_vpaid_mode").addClass("d-none");
       }
   }
   
   check();
   
   $('input[type=radio][name=ad_client]').change(function() {
        console.log($(this).val());
    });
}

function Added(response, formElement) {
    switch (response.status) {
        case false:
            if (response.errors.length === 1 && response.errors[0].id === "") {
                new AlertError(response.errors[0].message);
            } else if (response.errors.length > 0) {
                $.each(response.errors, function (key, value) {
                    let inputElement = $("body").find("#" + value.id);
                    inputElement.addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
                    inputElement.parent().addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
                });
                $("#" + response.errors[0].id).focus();
                $("#" + response.errors[0].id).prev('label').focus();
            } else {
                new AlertError("Error!");
            }
            break;
        case true:
            if ($("#select_line_item_type").val() !== "") {
                $("#select_line_item_type").prop("disabled", true);
            }
            NoticeSuccess("Advertising Schedules has been created successfully");
            setTimeout(function () {
                window.location.href = "/line-item";
            },1000);
            break;
        default:
            new AlertError("Undefined");
            break;
    }
}

function SubmitFormBidder(formID, functionCallback, ajxURL = "") {
    let formElement = $("#" + formID);
    let button = formElement.find(".submit");
    let submitButtonText = button.text();
    let submitButtonTextLoading = "Loading...";
    formElement.find("input").on("input", function (e) {
        let inputElement = $(this);
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    formElement.find("textarea").on("input", function (e) {
        let inputElement = $(this);
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    $(formElement).on('keypress', "input", function (e) {
        var keyCode = e.keyCode || e.which;
        if (keyCode === 13) {
            e.preventDefault();
            submit();
        }
    });
    formElement.find(".submit").on("click", function (e) {
        e.preventDefault();
        $.each(formElement.find("input.is-invalid"), function () {
            $(this).removeClass("is-invalid").next(".invalid-feedback").empty();
        });
        submit();
    });

    function submit() {
        let postData = formElement.serializeObject();
        
        
        const data = JSON.stringify(postData);
        const buttonElement = $(".submit");
        $.ajax({
            url: ajxURL,
            type: "POST",
            dataType: "JSON",
            contentType: "application/json",
            data: data,
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
    }
}