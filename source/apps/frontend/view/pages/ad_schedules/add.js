$(document).ready(function () {
    // CaptilizeBidder()
    SubmitForm("formAdvertisingSchedules", Added);

    checkSelectAdClient();
    checkSelectConfigAdBreaks();
    clickAddRoll();
    clickCancelRoll();
});

let countAdd = 0;

function getTypeRoll(elementBoxRoll) {
    let type;
    if (elementBoxRoll.hasClass("pre-roll-break")){
        type = "preroll";
    } else if (elementBoxRoll.hasClass("mid-roll-break")){
        type = "midroll";
    } else if (elementBoxRoll.hasClass("mid-roll-break")){
        type = "postroll";
    }
    return type;
}

function checkSelectAdClient() {
    function check() {
        let adClient = $('input[name="ad_client"]:checked').val();
        if (adClient === "1") {
            $(".box_vpaid_mode").addClass("d-none");
        } else {
            $(".box_vpaid_mode").removeClass("d-none");
        }
    }

    check();
    $('input[type=radio][name=ad_client]').change(function () {
        check();
    });
}

function checkSelectConfigAdBreaks() {
    function check() {
        let adBreaks = $('input[name="ad_breaks"]:checked').val();
        if (adBreaks === "1") {
            $(".vmap").addClass("d-none");
            $(".add-break-manually").removeClass("d-none");
        } else {
            $(".add-break-manually").addClass("d-none");
            $(".vmap").removeClass("d-none");
        }
    }

    check();
    $('input[type=radio][name=ad_breaks]').change(function () {
        check();
    });
}

function clickAddRoll() {
    $('.btn-add-roll').on("click", function () {
        let elementBoxRoll = $(this).closest(".box-roll");
        $(htmlAddRoll(getTypeRoll(elementBoxRoll))).insertAfter(elementBoxRoll);
        elementBoxRoll.find(".box-step-roll").addClass("_editing_z87zl_152");
        elementBoxRoll.find(".box-btn-add").addClass("d-none");
        elementBoxRoll.find(".box-btn-action").addClass("d-none");
        elementBoxRoll.find(".box-btn-edit").removeClass("d-none");

        elementBoxRoll.attr("data-action","add");
        $(".selectpicker").selectpicker("refresh");
    });
}

function clickCancelRoll() {
    $('.btn-cancel-roll').on("click", function () {
        let elementBoxRoll = $(this).closest(".box-roll");
        elementBoxRoll.find(".box-btn-edit").addClass("d-none");

        let action = elementBoxRoll.attr("data-action");
        if (action === "add"){
            elementBoxRoll.find(".box-step-roll").removeClass("_editing_z87zl_152");
            elementBoxRoll.find(".box-btn-add").removeClass("d-none");
        } else if (action === "done") {
            elementBoxRoll.find("._inactive_iyx7e_22").removeClass("_inactive_iyx7e_22").addClass("_active_z87zl_147");
            elementBoxRoll.find(".box-btn-action").removeClass("d-none");
        }

        let elementBoxAdTagUrl = elementBoxRoll.next();
        if (elementBoxAdTagUrl.hasClass("box-adtag-url")) {
            elementBoxAdTagUrl.remove();
        }

        $(".selectpicker").selectpicker("refresh");
    });
}

function clickEditRoll() {
    $('.btn-add-roll').on("click", function () {
        let elementBoxRoll = $(this).closest(".box-roll");
        $(htmlAddRoll()).insertAfter(elementBoxRoll);
        elementBoxRoll.find(".box-step-roll").addClass("_editing_z87zl_152");
        elementBoxRoll.find(".box-btn-add").addClass("d-none");
        elementBoxRoll.find(".box-btn-edit").removeClass("d-none");
        let elementBoxAdTagUrl = elementBoxRoll.next();
        if (elementBoxAdTagUrl.hasClass("box-adtag-url")) {
            $(elementBoxAdTagUrl.clone()).insertAfter(elementBoxAdTagUrl);
            let elementBoxAdTagUrlBackup = elementBoxAdTagUrl.next();
            if (elementBoxAdTagUrlBackup.hasClass("box-adtag-url")) {
                elementBoxAdTagUrlBackup.addClass("box-adtag-url-backup").removeClass("box-adtag-url").addClass("d-none");
            }
        }

        $(".selectpicker").selectpicker("refresh");
    });
}

function htmlAddRoll(typeRoll) {
    countAdd++;
    return `
<div class="_breakConfig_z87zl_166 pp_ad_schedules_7 -orientation-vertical -gutter-s -wrap hydrated w-100 box-adtag-url" data-type-roll="${typeRoll}">
    <div style="padding: 12px 20px;">
        <div class="_tagContainer_z87zl_192 pp_ad_schedules_7 -orientation-horizontal -align-center -gutter-s -wrap hydrated adtag-url">
            <div>
                <div class="wui-input -size-m -empty hydrated">
                    <div class="wui-input__box">
                        <input type="text" placeholder="Enter an ad tag URL"
                               autocomplete="off" required="">
                    </div>
                </div>
            </div>
        </div>
        
        <div class="wui-button -type-tertiary hydrated box-btn-add-tag-url">
            <button type="button" class="btn-add-waterfall">
                <div class="wui-icon wui-icon--dashboard_add -size-m hydrated">
                    <svg xmlns="http://www.w3.org/2000/svg"
                         viewBox="0 0 24 24" id="ds-icon-dashboard-add">
                        <path d="M20 10h-6V4a2 2 0 0 0-4 0v6H4a2 2 0 0 0 0 4h6v6a2 2 0 0 0 4 0v-6h6a2 2 0 0 0 0-4Z"></path>
                    </svg>
                </div>
                <span>Add Waterfall Tag</span>
            </button>
        </div>
        <hr>
        <div class="pp_ad_schedules_7 -orientation-horizontal -justify-space-between -wrap hydrated">
            <div style="padding: 12px 0px;">
                <div class="_breakOptions_z87zl_180 pp_ad_schedules_7 -orientation-vertical -gutter-s -wrap hydrated">
                    <div>
                        <div class="pp_ad_schedules_7 -orientation-horizontal -align-center -wrap hydrated">
                            <div>
                                <div class="d-flex align-items-center">
                                    <input class="mt-0 me-2"
                                           id="skipable_after_${countAdd}"
                                           name="skipable_after"
                                           type="checkbox" checked=""
                                           style="width: 16px; height: 16px; cursor: pointer">
                                    <label class="light-font-label pp_ad_schedules_font_text"
                                           for="skipable_after_${countAdd}">Skippable after</label>
                                    <div data-test="set-skip-offset"
                                         class="_numberInput_z87zl_124 wui-input -size-m -empty hydrated"
                                         type="number">
                                        <div class="wui-input__box">
                                            <input type="number"
                                                   autocomplete="off">
                                        </div>
                                    </div>
                                    <p class="pp_ad_schedules_font_text ms-1">
                                        seconds</p>
                                </div>
                            </div>
                        </div>
                        <div class="pp_ad_schedules_7 -orientation-horizontal -align-center -wrap hydrated">
                            <div>
                                <div class="d-flex align-items-center">
                                    <input class="mt-0 me-2 non-overlay-ad" id="overlay_ad_${countAdd}"
                                           type="checkbox" checked=""
                                           style="width: 16px; height: 16px; cursor: pointer">
                                    <label class="light-font-label pp_ad_schedules_font_text"
                                           for="overlay_ad_${countAdd}">Non-linear overlay ad</label>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="pp_ad_schedules_7 -orientation-vertical -gutter-s -wrap hydrated">
                    <div>
                        <h4>Break Timing</h4>
                        <div class="_breakTiming_z87zl_175 pp_ad_schedules_7 -orientation-horizontal -gutter-s -wrap hydrated">
                            <div>
                                <div class="_numberInput_z87zl_124 wui-input -size-m -empty hydrated"
                                     type="number">
                                    <div class="wui-input__box">
                                        <input type="number" autocomplete="off" min="0" required="">
                                    </div>
                                </div>
                                <div class="border border-gray-500 rounded-1 select-ad-schedules">
                                    <select name="break_mode" class="selectpicker p-0 ">
                                        <option value="1">Seconds into Video</option>
                                        <option value="2">Timecode</option>
                                        <option value="3">% of Video</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
    `;
}

function htmlAddAdtagUrl() {
    return `
 <div class="_tagContainer_z87zl_192 pp_ad_schedules_7 -orientation-horizontal -align-center -gutter-s -wrap hydrated">
                <div>
                    <div class="wui-input -size-m -empty hydrated">
                        <div class="wui-input__box">
                            <input type="text" placeholder="Enter an ad tag URL" autocomplete="off">
                        </div>
                    </div>
                    <div data-test="delete-waterfall-tag" class="_deleteButton_z87zl_192 wui-button -type-icon hydrated" type="icon">
                        <button type="button">
                            <div class="wui-icon wui-icon--dashboard_trash -size-m hydrated">
                                <svg xmlns="http://www.w3.org/2000/svg"
                                     viewBox="0 0 24 24"
                                     id="ds-icon-dashboard-trash">
                                    <path d="M21 4h-4a1 1 0 0 0 0-2H7a1 1 0 0 0 0 2H3a1 1 0 0 0 0 2h18a1 1 0 0 0 0-2Z"></path>
                                    <rect x="6" y="8.91" width="12"
                                          height="13.09" rx="1" ry="1"></rect>
                                </svg>
                            </div>
                        </button>
                    </div>
                </div>
            </div>
 `;
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
            }, 1000);
            break;
        default:
            new AlertError("Undefined");
            break;
    }
}

function SubmitForm(formID, functionCallback, ajxURL = "") {
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