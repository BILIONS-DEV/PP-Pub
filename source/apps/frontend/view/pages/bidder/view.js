let idx = 0;

function firstLoad() {
    let valueBidder = $("#select-bidder").val();
    var newOption = new Option("Native", '3', false, false);
    if (valueBidder === "1") {
        //Empty media type select
        $(".btn-add-param").addClass('disabled');
        $(".bidder-amz").addClass('d-none');
        $(".bidder-google").removeClass("d-none");
        $('#media_type').prop('disabled', false);
        $("#media_type option[value='3']").remove();
        $("#media_type").selectpicker('refresh');
        // $('#media_type').val(null).trigger('change');
        let listParams = $('.list_params');
        listParams.html("");
        listParams.prepend(paramHtmlGoogle());
        $(".box-params").addClass("d-none");
    } else if (valueBidder === "2") {
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
        let desAmz = $("#pub_id").data("content-amz");
        $("#pub_id_des").attr("data-bs-content", desAmz);
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
    }
    // idx = $('.box-params').data("last-index");
}

$(document).ready(function () {
    firstLoad();

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
    var myFieldset = document.getElementById("myFieldset");
    myFieldset.disabled = true;
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
            // $('#media_type').val(null).trigger('change');
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
            let desAmz = $("#pub_id").data("content-amz");
            $("#pub_id_des").attr("data-bs-content", desAmz);
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