let _validFileImageExtensions = [".jpg", ".jpeg", ".bmp", ".gif", ".png"];
let maxSize = 15 * 1024 * 1024;
let api_upload_image = "https://ul.pubpowerplatform.io/api/v1/image";
let domainUpload = "https://ul.pubpowerplatform.io";
// let api_upload_image = "http://127.0.0.1:8543/api/v1/image";
// let domainUpload = "http://127.0.0.1:8543";
$(document).ready(function () {
    $('[data-bs-toggle="popover"]').popover({
        html: true,
        sanitize: false,
    });
    $("#slider-container").css("cursor","not-allowed");
    $(".slide").css("pointer-events","none");
    $(".slide img").css("cursor","not-allowed");
    firstLoad();
    selectTab();
    selectPlayerType();
    selectPlayerLayout();
    selectPlayerSize();
    selectPlayMode();
    selectFloating();
    selectPosition();
    selectShowMainTitle();
    selectColor();
    clickAddLogo();
    changeAutoSkip();
    eventMakePreview();
    var myFieldset = document.getElementById("myFieldset");
    myFieldset.disabled = true;
});

function changeAutoSkip() {
    $("#auto_skip").change(function () {
        checkShowAdvertisement();
        makePreview();
    });
}

function checkShowAdvertisement() {
    let adType = $("#type").val();
    let layout = $("#player_layout").val();
    let isAutoSkip = $("#auto_skip").is(":checked");
    let boxAutoSkip = $(".box-auto-skip");
    let boxTimeToSkip = $(".box-time-to-skip");
    let boxShowAutoSkipButton = $(".box-show-auto-skips-button");
    let boxNumberOfPreRoll = $(".box-number-of-pre-roll");
    let boxDelay = $(".box-delay");
    if (adType === "2" && layout !== "8"){
        boxAutoSkip.removeClass("d-none");
        if (isAutoSkip){
            boxTimeToSkip.removeClass("d-none");
            boxShowAutoSkipButton.removeClass("d-none");
            boxNumberOfPreRoll.removeClass("d-none");
            boxDelay.addClass("d-none");
        } else {
            boxTimeToSkip.addClass("d-none");
            boxShowAutoSkipButton.addClass("d-none");
            boxNumberOfPreRoll.addClass("d-none");
            boxDelay.removeClass("d-none");
        }
    } else {
        boxAutoSkip.addClass("d-none");
        boxTimeToSkip.addClass("d-none");
        boxShowAutoSkipButton.addClass("d-none");
        boxNumberOfPreRoll.addClass("d-none");
        boxDelay.removeClass("d-none");
    }
}

function checkShowColumnsPosition() {
    let adType = $("#type").val();
    let layout = $("#player_layout").val();
    let box = $(".box-columns-position");
    let positionRight = $("#columns_position_right");
    if (adType === "2" && (layout === "6" || layout === "7")) {
        box.removeClass("d-none");
        positionRight.prop("checked", true);
    } else {
        box.addClass("d-none");
    }
}

function checkShowColumnsNumber() {
    let adType = $("#type").val();
    let layout = $("#player_layout").val();
    let box = $(".box-columns-number");
    let columnsTwo = $("#columns_number_two");
    if (adType === "2" && layout === "6") {
        box.removeClass("d-none");
        columnsTwo.prop("checked", true);
    } else {
        box.addClass("d-none");
    }
}

function checkDisableLogoConfig() {
    let adType = $("#type").val();
    // Logo Config
    let logoConfig = $(".logo-config");
    if (adType === "2") {
        // logoConfig.removeClass("disabled");
        // logoConfig.prop("disabled", false);
        // logoConfig.next().removeClass("disabled");
    } else {
        // logoConfig.addClass("disabled");
        // logoConfig.prop("disabled", true);
        // logoConfig.next().addClass("disabled");
    }
}

function checkDisableControls() {
    let adType = $("#type").val();
    // Controls Color
    let controlsColor = $(".controls");
    if (adType === "2") {
        // controlsColor.removeClass("disabled");
        // controlsColor.prop("disabled", false);
        // controlsColor.next().removeClass("disabled");
    } else {
        // controlsColor.addClass("disabled");
        // controlsColor.prop("disabled", true);
        // controlsColor.next().addClass("disabled");
    }
}

function checkDisableColor() {
    let adType = $("#type").val();
    let layout = $("#player_layout").val();
    // Controls Color
    let controlsColor = $(".box-controls-color");
    if (adType === "2") {
        if (layout !== "8") {
            controlsColor.removeClass("d-none");
            controlsColor = $("#control_color");
            // controlsColor.removeClass("disabled");
            // controlsColor.prop("disabled", false);
            // controlsColor.prev().removeClass("disabled");
        } else {
            controlsColor.addClass("d-none");
        }
    } else if (adType === "3") {
        controlsColor.removeClass("d-none");
        controlsColor = $("#control_color");
        // controlsColor.addClass("disabled");
        // controlsColor.prop("disabled", true);
        // controlsColor.prev().addClass("disabled");
    } else {
        controlsColor.addClass("d-none");
    }

    // Background Color
    let backgroundColor = $(".box-background-color");
    if (adType === "2") {
        backgroundColor.removeClass("d-none");
        backgroundColor = $("#background_color");
        // backgroundColor.removeClass("disabled");
        // backgroundColor.prop("disabled", false);
        // backgroundColor.prev().removeClass("disabled");
    } else if (adType === "3") {
        backgroundColor.removeClass("d-none");
        backgroundColor = $("#background_color");
        // backgroundColor.addClass("disabled");
        // backgroundColor.prop("disabled", true);
        // backgroundColor.prev().addClass("disabled");
    }

    // Title Color
    let titleColor = $(".box-title-color");
    if (adType === "2") {
        titleColor.removeClass("d-none");
        titleColor = $("#title_color");
        // titleColor.removeClass("disabled");
        // titleColor.prop("disabled", false);
        // titleColor.prev().removeClass("disabled");
    } else if (adType === "3") {
        titleColor.removeClass("d-none");
        titleColor = $("#title_color");
        // titleColor.addClass("disabled");
        // titleColor.prop("disabled", true);
        // titleColor.prev().addClass("disabled");
    }

    // Description Color
    let descriptionColor = $(".box-description-color");
    if (adType === "2") {
        if (layout === "3" ||layout === "4" ||layout === "5" ||layout === "7") {
            descriptionColor.removeClass("d-none");
            descriptionColor = $("#description_color");
            // descriptionColor.removeClass("disabled");
            // descriptionColor.prop("disabled", false);
            // descriptionColor.prev().removeClass("disabled");
        } else {
            descriptionColor.addClass("d-none");
        }
    } else if (adType === "3") {
        descriptionColor.removeClass("d-none");
        descriptionColor = $("#description_color");
        // descriptionColor.addClass("disabled");
        // descriptionColor.prop("disabled", true);
        // descriptionColor.prev().addClass("disabled");
    } else {
        descriptionColor.addClass("d-none");
    }

    // Theme Color
    let themeColor = $(".box-theme-color");
    if (adType === "2" && layout === "8") {
        themeColor.removeClass("d-none");
    } else {
        themeColor.addClass("d-none");
    }

    // Title Background Color
    let titileBackgroundColor = $(".box-title-background-color");
    if (adType === "2" && layout === "8") {
        titileBackgroundColor.removeClass("d-none");
    } else {
        titileBackgroundColor.addClass("d-none");
    }
}

function selectShowMainTitle() {
    $("#show_main_title").change(function () {
        checkDisableDisplayOption();
    });
}

function checkDisableDisplayOption() {
    let adType = $("#type").val();
    let layout = $("#player_layout").val();
    // Show main title
    let showMainTitleElement = $("#show_main_title");
    if (adType === "2") {
        // showMainTitleElement.removeClass("disabled");
        // showMainTitleElement.prop("disabled", false);
        // showMainTitleElement.next().removeClass("disabled");
    } else {
        // showMainTitleElement.addClass("disabled");
        // showMainTitleElement.prop("disabled", true);
        // showMainTitleElement.next().addClass("disabled");
    }
    // Show content title
    let showContentTitleElement = $("#show_content_title");
    if (adType === "2") {
        // showContentTitleElement.removeClass("disabled");
        // showContentTitleElement.prop("disabled", false);
        // showContentTitleElement.next().removeClass("disabled");
    } else {
        // showContentTitleElement.addClass("disabled");
        // showContentTitleElement.prop("disabled", true);
        // showContentTitleElement.next().addClass("disabled");
    }
    // Show content description
    let showContentDescElement = $("#show_content_description");
    let boxShowContentDesc = $(".box-show-content-desc");
    if (adType === "2") {
        if (layout === "3" ||layout === "4" ||layout === "5" ||layout === "7" || layout === "8") {
            boxShowContentDesc.removeClass("d-none");
            // showContentDescElement.removeClass("disabled");
            // showContentDescElement.prop("disabled", false);
            // showContentDescElement.next().removeClass("disabled");
        } else {
            boxShowContentDesc.addClass("d-none");
            // showContentDescElement.removeClass("disabled");
            // showContentDescElement.prop("disabled", false);
            // showContentDescElement.next().removeClass("disabled");
        }
    } else {
        boxShowContentDesc.removeClass("d-none");
        // showContentDescElement.addClass("disabled");
        // showContentDescElement.prop("disabled", true);
        // showContentDescElement.next().addClass("disabled");
    }
    // Show controls
    let showControlsElement = $("#show_controls");
    if (adType === "2") {
        // showControlsElement.removeClass("disabled");
        // showControlsElement.prop("disabled", false);
        // showControlsElement.next().removeClass("disabled");
    } else {
        // showControlsElement.addClass("disabled");
        // showControlsElement.prop("disabled", true);
        // showControlsElement.next().addClass("disabled");
    }
    // Main title text
    let mainTitleTextElement = $("#main_title_text");
    if ((adType === "2" && showMainTitleElement.is(":checked"))) {
        // mainTitleTextElement.removeClass("disabled");
        // mainTitleTextElement.prop("disabled", false);
        // mainTitleTextElement.next().removeClass("disabled");
    } else {
        // mainTitleTextElement.addClass("disabled");
        // mainTitleTextElement.prop("disabled", true);
        // mainTitleTextElement.next().addClass("disabled");
    }
}

function clickAddLogo() {
    $(".btn-add-logo").on("click", function () {
        $("#logo").click();
    });
    $("#logo").on("change", function (e) {
        validateAndUploadLogo($(this), _validFileImageExtensions, e);
    });
}

function selectPosition() {
    $("#position_desktop").on("change", function () {
        checkShowPositionDesktop();
        makePreview();
    });
    $("#position_mobile").on("change", function () {
        checkShowPositionMobile();
        makePreview();
    });
}

function checkShowPositionDesktop() {
    let position = $("#position_desktop").val();
    let elementPosition = $(".position_desktop");
    elementPosition.addClass("d-none");
    if (position === "1") {
        elementPosition = $(".position_desktop_bottom_right");
    } else if (position === "2") {
        elementPosition = $(".position_desktop_bottom_left");
    } else if (position === "3") {
        elementPosition = $(".position_desktop_top_right");
    } else if (position === "4") {
        elementPosition = $(".position_desktop_top_left");
    }
    elementPosition.removeClass("d-none");
}

function checkShowPositionMobile() {
    let position = $("#position_mobile").val();
    let elementPosition = $(".position_mobile");
    elementPosition.addClass("d-none");
    if (position === "1") {
        elementPosition = $(".position_mobile_bottom_right");
    } else if (position === "2") {
        elementPosition = $(".position_mobile_bottom_left");
    }
    elementPosition.removeClass("d-none");
}

function selectFloating() {
    $("#floating_on_desktop").change(function () {
        checkDisableConfigFloatingDesktop();
        makePreview();
    });
    $("#floating_on_mobile").change(function () {
        checkDisableConfigFloatingMobile();
        makePreview();
    });
}

function checkDisableConfigFloatingDesktop() {
    let isFloatingDesktop = $('#floating_on_desktop').is(":checked");
    // Close Floating Button
    let elementFloatingOnDesktop = $(".floating_on_desktop");
    if (isFloatingDesktop) {
        // elementFloatingOnDesktop.removeClass("disabled");
        // elementFloatingOnDesktop.prop("disabled", false);
        // elementFloatingOnDesktop.next().removeClass("disabled");
    } else {
        // elementFloatingOnDesktop.addClass("disabled");
        // elementFloatingOnDesktop.prop("disabled", true);
        // elementFloatingOnDesktop.next().addClass("disabled");
    }
}

function checkDisableConfigFloatingMobile() {
    let isFloatingMobile = $('#floating_on_mobile').is(":checked");
    // Close Floating Button
    let elementFloatingOnMobile = $(".floating_on_mobile");
    if (isFloatingMobile) {
        // elementFloatingOnMobile.removeClass("disabled");
        // elementFloatingOnMobile.prop("disabled", false);
        // elementFloatingOnMobile.next().removeClass("disabled");
    } else {
        // elementFloatingOnMobile.addClass("disabled");
        // elementFloatingOnMobile.prop("disabled", true);
        // elementFloatingOnMobile.next().addClass("disabled");
    }
}

function selectPlayerSize() {
    $('input[type=radio][name=size]').change(function () {
        checkShowPlayerSize();
    });
}

function checkShowPlayerSize() {
    let type = $('input[name="size"]:checked').val();
    let boxWidth = $(".box-width");
    let boxRatio = $(".box-ratio");
    if (type === "1") {
        boxWidth.addClass("d-none");
        boxRatio.removeClass("d-none");
    } else if (type === "2") {
        boxRatio.addClass("d-none");
        boxWidth.removeClass("d-none");
    }
}

function selectPlayerType() {
    $("#type").on("change", function () {
        firstLoad();
    });
}

function checkShowLayout() {
    if ($("#type").val() === "2") {
        $(".box-player-layout").removeClass("d-none");
    } else {
        $(".box-player-layout").addClass("d-none");
    }
}

function checkShowFloatOnView() {
    let adType = $("#type").val();
    if (adType === "2") {
        $(".box-float-on-view").removeClass("d-none");
        $(".box-float-on-view-mobile").removeClass("d-none");
    } else {
        $(".box-float-on-view").addClass("d-none");
        $(".box-float-on-view-mobile").addClass("d-none");
    }
}

function selectPlayerLayout() {
    $(".slide").on("click", function () {
        $(".slide").removeClass("selected");
        $(this).addClass("selected");
        let value = $(this).data("value");
        $(this).prevAll("#player_layout").attr("value", value);
        firstLoad();
    });
}

function selectPlayMode() {
    $(".play_mode").on("click", function () {
        $(".play_mode").removeClass("pt52");
        $(this).addClass("pt52");
        let type = $(this).data("type");
        $(".box-play-mode").addClass("d-none");
        $(".box-" + type).removeClass("d-none");
    });
}

function selectTab() {
    $(".template-tab").on("click", ".nav-link", function () {
        $(".template-tab").find(".nav-link").removeClass("pp-4");
        $(this).addClass("pp-4");
        $(".nav-link").removeClass("at-1");
        var tab = $(this).attr("data-tab");
        if (tab !== "1") {
            $(this).addClass("at-1");
        }
    });
}

function selectColor() {
    let elementInputColor = $(".input-color");
    elementInputColor.on("input", function () {
        $(this).next("input").val($(this).val());
    });
    elementInputColor.on("change", function () {
        makePreview();
    });
    elementInputColor.next("input").on("input", function (e) {
        $(this).prev("input").val(e.target.value);
    });
    elementInputColor.next("change", function () {
        makePreview();
    });
}

function eventMakePreview() {
    $("input").on("change", function () {
        makePreview();
    });
}

function makeData(postData) {
    postData.id = parseInt(postData.id);
    postData.type = parseInt(postData.type);
    postData.vast_retry = parseInt(postData.vast_retry);
    postData.delay = parseInt(postData.delay);
    postData.time_to_skip = parseInt(postData.time_to_skip);
    postData.max_width = parseInt(postData.max_width);
    postData.width = parseInt(postData.width);
    postData.floating_width = parseInt(postData.floating_width);
    postData.floating_width_mobile = parseInt(postData.floating_width_mobile);
    postData.margin_bottom_desktop = parseInt(postData.margin_bottom_desktop);
    postData.margin_bottom_mobile = parseInt(postData.margin_bottom_mobile);
    postData.margin_left_desktop = parseInt(postData.margin_left_desktop);
    postData.margin_left_mobile = parseInt(postData.margin_left_mobile);
    postData.margin_right_desktop = parseInt(postData.margin_right_desktop);
    postData.margin_right_mobile = parseInt(postData.margin_right_mobile);
    postData.margin_top_desktop = parseInt(postData.margin_top_desktop);
    postData.columns_number = parseInt(postData.columns_number);
    postData.columns_position = parseInt(postData.columns_position);
    postData.columns_number = parseInt(postData.columns_number);
    postData.player_layout = parseInt(postData.player_layout);
    postData.show_auto_skip_button = parseInt(postData.show_auto_skip_button);
    postData.number_of_pre_roll_ads = parseInt(postData.number_of_pre_roll_ads);
    postData.floating_position_desktop = parseInt(postData.floating_position_desktop);
    postData.floating_position_mobile = parseInt(postData.floating_position_mobile);
    postData.play_mode = parseInt(postData.play_mode);
    postData.advertisement_scenario = parseInt(postData.advertisement_scenario);
    postData.size = parseInt(postData.size);
    postData.auto_start = parseInt(postData.auto_start);
    postData.main_title_top_article = postData.main_title_text;
    postData.link = $("#preview-logo").attr("src");
    if (postData.close_floating_button_mobile === "on") {
        postData.close_floating_button_mobile = 1;
    } else {
        postData.close_floating_button_mobile = 2;
    }

    if (postData.close_floating_button_desktop === "on") {
        postData.close_floating_button_desktop = 1;
    } else {
        postData.close_floating_button_desktop = 2;
    }

    if (postData.default_sound_mode === "on") {
        postData.default_sound_mode = 1;
    } else {
        postData.default_sound_mode = 2;
    }

    if (postData.description_enable === "on") {
        postData.description_enable = 1;
    } else {
        postData.description_enable = 2;
    }

    if (postData.show_controls === "on") {
        postData.show_controls = 1;
    } else {
        postData.show_controls = 2;
    }

    if (postData.floating_on_desktop === "on") {
        postData.floating_on_desktop = 1;
    } else {
        postData.floating_on_desktop = 2;
    }

    if (postData.floating_on_mobile === "on") {
        postData.floating_on_mobile = 1;
    } else {
        postData.floating_on_mobile = 2;
    }

    if (postData.main_title === "on") {
        postData.main_title = 1;
    } else {
        postData.main_title = 2;
    }

    if (postData.title_enable === "on") {
        postData.title_enable = 1;
    } else {
        postData.title_enable = 2;
    }

    if (postData.powered_by === "on") {
        postData.powered_by = 1;
    } else {
        postData.powered_by = 2;
    }

    if (postData.powered_by_top_article === "on") {
        postData.powered_by_top_article = 1;
    } else {
        postData.powered_by_top_article = 2;
    }

    if (postData.share_button === "on") {
        postData.share_button = 1;
    } else {
        postData.share_button = 2;
    }

    if (postData.video_config === "on") {
        postData.video_config = 1;
    } else {
        postData.video_config = 2;
    }

    if (postData.show_stats === "on") {
        postData.show_stats = 1;
    } else {
        postData.show_stats = 2;
    }

    if (postData.fullscreen_button === "on") {
        postData.fullscreen_button = 1;
    } else {
        postData.fullscreen_button = 2;
    }

    if (postData.next_prev_arrows_button === "on") {
        postData.next_prev_arrows_button = 1;
    } else {
        postData.next_prev_arrows_button = 2;
    }

    if (postData.next_prev_time === "on") {
        postData.next_prev_time = 1;
    } else {
        postData.next_prev_time = 2;
    }

    postData.custom_logo_top_article = 1;
    postData.custom_logo = 1;
    if (postData.enable_logo === "on") {
        postData.enable_logo = 1;
    } else {
        postData.enable_logo = 2;
    }
    postData.enable_logo_top_article = postData.enable_logo;

    if (postData.auto_skip === "on") {
        postData.auto_skip = 1;
    } else {
        postData.auto_skip = 2;
    }

    if (postData.float_on_bottom === "on") {
        postData.float_on_bottom = 1;
    } else {
        postData.float_on_bottom = 2;
    }
    if (postData.floating_on_view === "on") {
        postData.floating_on_view = 1;
    } else {
        postData.floating_on_view = 2;
    }
    if (postData.float_on_bottom_mobile === "on") {
        postData.float_on_bottom_mobile = 1;
    } else {
        postData.float_on_bottom_mobile = 2;
    }
    if (postData.floating_on_view_mobile === "on") {
        postData.floating_on_view_mobile = 1;
    } else {
        postData.floating_on_view_mobile = 2;
    }
    if (postData.floating_on_impression === "on") {
        postData.floating_on_impression = 1;
    } else {
        postData.floating_on_impression = 2;
    }
    if (postData.wait_for_ad === "on") {
        postData.wait_for_ad = 1;
    } else {
        postData.wait_for_ad = 2;
    }
    if (postData.pre_roll === "on") {
        postData.pre_roll = 1;
    } else {
        postData.pre_roll = 2;
    }
    if (postData.mid_roll === "on") {
        postData.mid_roll = 1;
    } else {
        postData.mid_roll = 2;
    }
    if (postData.post_roll === "on") {
        postData.post_roll = 1;
    } else {
        postData.post_roll = 2;
    }
    return postData;
}

function makePreview() {
    let postData = $("#formTemplate").serializeObject();
    postData = makeData(postData);
    $.ajax({
        url: "/player/template/preview",
        type: "POST",
        dataType: "JSON",
        contentType: "application/json",
        data: JSON.stringify(postData),
        success: function (json) {
            document.querySelector("#videocontainer").innerHTML = "";
            (vitag.Init = window.vitag.Init || []).push(function () {
                // var config = ""
                viAPItag.initPowerVideoContainer(json);
            });
        },
    });
}

function firstLoad() {
    checkShowLayout();
    checkShowPlayerSize();
    checkDisableConfigFloatingDesktop();
    checkDisableConfigFloatingMobile();
    checkShowPositionDesktop();
    checkShowPositionMobile();
    checkDisableDisplayOption();
    checkDisableColor();
    checkDisableControls();
    checkDisableLogoConfig();
    checkShowColumnsNumber();
    checkShowColumnsPosition();
    checkShowAdvertisement();
    // makePreview();
}


function uploadFile(event, element) {
    var fd = new FormData();
    var file = event.target.files[0];
    if (file.size > maxSize) {
        new AlertError("You uploaded file over 10mb, please choose another file!");
        return;
    }
    fd.append('file', file);
    $.ajax({
        url: api_upload_image,
        type: "POST",
        // dataType: 'json',
        contentType: false,
        processData: false,
        data: fd,
        beforeSend: function (xhr) {
            // buttonElement.attr('disabled', true).text(submitButtonTextLoading);
        },
        error: function (jqXHR, exception) {
            const msg = AjaxErrorMessage(jqXHR, exception);
            new AlertError("AJAX ERROR: " + msg);
            // buttonElement.attr('disabled', false).text(submitButtonText);
        },
        success: function (responseJSON) {
            // buttonElement.attr('disabled', false).text(submitButtonText);
        },
        complete: function (res) {
            afterUpload(res.responseJSON, element);
        }
    });
}

function afterUpload(response, element) {
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
            break;
        case "success":
            $("#preview-logo").attr("src", domainUpload + response.data_object.thumb[0]);
            $("#logo").attr("value", response.data_object.thumb[0]);
            break;
        default:
            new AlertError("Undefined");
            break;
    }
}

function validateAndUploadLogo(element, validFile, event) {
    let oInput = element[0];
    if (oInput.type === "file") {
        let sFileName = oInput.value;
        if (sFileName.length > 0) {
            let blnValid = false;
            for (let j = 0; j < validFile.length; j++) {
                let sCurExtension = validFile[j];
                if (sFileName.substr(sFileName.length - sCurExtension.length, sCurExtension.length).toLowerCase() === sCurExtension.toLowerCase()) {
                    blnValid = true;
                    uploadFile(event, element);
                    break;
                }
            }
            if (!blnValid) {
                new AlertError("Sorry, " + sFileName + " is invalid, allowed extensions are: " + validFile.join(", "));
                return false;
            }

        }
    }
    return true;
}