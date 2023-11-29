let idx = 0;

$(document).ready(function () {
    $('[data-bs-toggle="popover"]').popover({
        html: true,
        sanitize: false,
    });
    // CaptilizeBidder()
    SubmitFormCreate("submitBidder", Added);

    //clone params
    $('.btn-add-param').click(function () {
        $('.list_params').append($(".box-param-default").html());
        let demand = $("#select-bidder").find(':selected').text();
        addParam(demand);
        intiSelectPicker();
    });

    //delete params
    $('.list_params').on('click', '.btn-remove-param', function () {
        $(this).closest('.box-param').remove();
    });

    $('.list_params').on('input', '.param-add .param .param_name', function () {
        let inputElement = $(this);
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid");
            inputElement.closest(".param").removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    selectBidder();
});

function intiSelectPicker() {
    $(".list_params").find('.selectpicker').selectpicker();
}

function selectBidder() {
    $('#select-bidder').on('change', function () {
        idx = 0;
        var newOption = new Option("Native", '3', false, false);
        //If val = 1 is google
        if ($(this).val() === "1") {
            //Empty media type select
            $(".btn-add-param").addClass('disabled');
            $(".bidder-amz").addClass('d-none');
            $(".bidder-google").removeClass("d-none");
            $('#media_type').prop('disabled', false);
            $("#media_type option[value='3']").remove();
            $("#media_type").selectpicker('refresh');
            $("#pub_id_des_amz").addClass("d-none");
            $("#pub_id_des_google").removeClass("d-none");
            let listParams = $('.list_params');
            listParams.html("");
            listParams.prepend(paramHtmlGoogle());
            $(".box-params").addClass("d-none");
        } else if ($(this).val() === "2") {
            $(".box-params").removeClass("d-none");
            $(".btn-add-param").addClass('disabled');
            $(".bidder-google").addClass("d-none");
            $('#media_type').prop('disabled', false);
            if (!$('#media_type').find("option[value='3']").length) {
                // Append it to the select
                $('#media_type').append(newOption).trigger('change');
            }
            $("#media_type").selectpicker('refresh');
            $(".bidder-amz").removeClass("d-none");
            $("#pub_id_des_google").addClass("d-none");
            $("#pub_id_des_amz").removeClass("d-none");
            addTemplateBidder($(this).val());
        } else {
            $(".box-params").removeClass("d-none");
            $(".btn-add-param").removeClass('disabled');
            $(".bidder-amz").addClass('d-none');
            $(".bidder-google").addClass("d-none");
            $('#media_type').prop('disabled', true);
            if (!$('#media_type').find("option[value='3']").length) {
                // Append it to the select
                $('#media_type').append(newOption).trigger('change');
            }
            $("#media_type").selectpicker('refresh');
            addTemplateBidder($(this).val());
        }
    });
}

function CaptilizeBidder() {
    var capitalizeMe = "";
    $('#select-bidder option').each(function () {
        capitalizeMe = $(this).text();
        $(this).text(capitalizeMe.charAt(0).toUpperCase() + capitalizeMe.substring(1));
    });
}

function paramHtmlGoogle() {
    return `
    <div class="text-white m-2">No Param</div>
`;
}

function addParamBidder() {
    let params = [];
    $(".list_params").find(".param_value").each(function () {
        if ($(this).val()){
            let param = {};
            param['name'] = $(this).val();
            param['type'] = $(this).find(':selected').data('type');
            params.push(param);
        }
    });
    return params;
}

function addTemplateBidder(id) {
    $.ajax({
        url: "/bidder/addTemplate",
        type: "GET",
        data: {
            id: id,
        },
        success: function (json) {
            let listParams = $('.list_params');
            listParams.html("");
            listParams.prepend(json.params);

            //Change select
            $('#media_type').val(json.list_media_type).trigger('change');
            intiSelectPicker();
        },
        error: function (xhr) {
            console.log(xhr);
        }
    });
}

function addParam(demand) {
    $.ajax({
        url: "/bidder/addParam",
        type: "GET",
        data: {
            demand: demand,
        },
        success: function (json) {
            let listParams = $('.list_params');
            listParams.append(json);
            intiSelectPicker();
        },
        error: function (xhr) {
            console.log(xhr);
        }
    });
}

function SubmitFormCreate(formID, functionCallback, ajxURL = "") {
    let formElement = $("#" + formID);
    formElement.find("input").on("input", function (e) {
        let inputElement = $(this);
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
            inputElement.closest(".param").removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    formElement.find("textarea").on("input", function (e) {
        let inputElement = $(this);
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    formElement.find(".select2").on("select2:open", function () {
        let selectElement = $(this);
        if (selectElement.next().find(".select2-selection").hasClass("select2-is-invalid")) {
            selectElement.next().find(".select2-selection").removeClass("select2-is-invalid");
            selectElement.removeClass("is-invalid").next(".invalid-feedback").empty();
            selectElement.closest('.box-select2').next(".invalid-feedback").empty();
        }
    })
    formElement.find(".selectpicker").on("change", function() {
        let selectElement = $(this);
        if (selectElement.next(".dropdown-toggle").hasClass("select2-is-invalid")) {
            selectElement.next(".dropdown-toggle").addClass("btn-light").removeClass("select2-is-invalid");
            selectElement.removeClass("is-invalid").next(".invalid-feedback").empty();
            selectElement.closest('.box-selectpicker').next(".invalid-feedback").empty();
        }
    });
    formElement.find(".submit").on("click", function (e) {
        e.preventDefault();
        const buttonElement = $(this);
        const submitButtonText = buttonElement.text();
        const submitButtonTextLoading = "Loading...";
        const postData = formElement.serializeObject();
        postData.bidder_id = parseInt($("#select-bidder").val());
        // if (postData.bidder_id === 1) {
        postData.media_type = $("#media_type").val();
        // }
        postData.bid_adjustment = parseFloat(postData.bid_adjustment);
        postData.linked_gam = parseInt(postData.linked_gam);
        postData.account_type = parseInt(postData.account_type);
        postData.rpm = parseFloat(postData.rpm);
        postData.params = addParamBidder();
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

function Added(response, formElement) {
    switch (response.status) {
        case "error":
            if (response.errors.length === 1 && response.errors[0].id === "") {
                new AlertError(response.errors[0].message);
            } else if (response.errors.length > 0) {
                $.each(response.errors, function (key, value) {
                    let inputElement = $("#" + value.id);
                    if (value.id === "display_name") {
                        let inputWidth = inputElement.width()
                        let labelWidth = inputElement.prev("label.col-form-label").width()
                        let fullWidth = inputElement.closest("div.d-flex").width()
                        let widthForErr = fullWidth - labelWidth - inputWidth - 50
                        inputElement.next("span.invalid-feedback").css("max-width", widthForErr)
                    }
                    inputElement.addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
                    inputElement.next().find(".select2-selection").addClass("select2-is-invalid")
                    inputElement.closest('.box-select2').addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message)
                    inputElement.next(".dropdown-toggle").removeClass("btn-light").addClass("select2-is-invalid")
                    inputElement.closest('.box-selectpicker').addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message)
                    inputElement.find(".param_name").addClass("is-invalid")
                });
                $("#" + response.errors[0].id).focus();
                $("#" + response.errors[0].id).prev('label').focus();
                // new AlertError(response.errors[0].message, function () {
                //     $("#" + response.errors[0].id).focus();
                //     $("#" + response.errors[0].id).prev('label').focus();
                // })
            } else {
                new AlertError("Error!");
            }
            break
        case "success":
            NoticeSuccess("Bidder has been created successfully")
            setTimeout(function () {
                window.location.replace("/bidder");
            }, 1000);
            break
        default:
            new AlertError("Undefined");
            break
    }
}