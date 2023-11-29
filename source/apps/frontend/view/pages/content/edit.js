let _validFileImageExtensions = [".jpg", ".jpeg", ".bmp", ".gif", ".png"];
let _validFileVideoExtensions = [".m4v", ".avi", ".mpg", ".mp4", ".3gp"];
let api_upload_video = "https://ul.pubpowerplatform.io/api/v1/video";
let api_upload_image = "https://ul.pubpowerplatform.io/api/v1/image";
let domainUpload = "https://ul.pubpowerplatform.io";
// let api_upload_video = "http://127.0.0.1:8543/api/v1/video"
// let api_upload_image = "http://127.0.0.1:8543/api/v1/image"
// let domainUpload = "http://127.0.0.1:8543"
let IsUploadVideoFinish = true;
let IsUploadThumbFinish = true;
let CheckImgSelected = false;
let IsCheckRedirect = false;
let indexAdBreak = 1;

let maxSize = 200 * 1024 * 1024;


$(function () {
    $(".content-tab").on("click", ".nav-link", function () {
        $(".content-tab").find(".nav-link").removeClass("pp-4");
        $(this).addClass("pp-4");
        var tab = $(this).attr("data-tab");
        if (tab != "1") {
            $(this).addClass("at-1");
        } else {
            $(".content-tab").find('.at-1').removeClass('at-1');
        }
    });

    firstLoad();

    SubmitFormEdit("editContent", Edit, "/content/edit");

    eventClickAndChangeUpload();

    $(".select-thumb").parent().on("click", function () {
        // $(".select-thumb").parent().removeClass("selected")
        if ($(this).hasClass("selected")) {
            $(this).removeClass("selected");
            $("#thumb").val("");
            $(".preview-thumb").attr("src", "");
            CheckImgSelected = false;
        } else {
            // if (CheckImgSelected) {
            //
            // }
            $(".select-thumb").parent().removeClass("selected");
            $(this).addClass("selected");
            let src = $(this).find(".select-thumb").attr("src");
            $(".preview-thumb").attr("src", src).removeClass("d-none");
            $("#thumb").val(src.replace(domainUpload, ""));
            $("#box-upload-thumb").removeClass("is-invalid").next(".invalid-feedback").empty();
            CheckImgSelected = true;
        }
    });

    adBreak();
    handleTags();
    checkRedirect();
    // LoadHistory();

    // Ad Schedule
    initAdSchedule();
});

function initAdSchedule() {
    checkSelectConfigAdBreaks();
    clickAddRoll();
    clickCancelRoll();
    clickDoneRoll();
    clickEditRoll();
    clickDeleteRoll();
    checkInputTimeBreak();
    clickCopyContentId();
}

function clickCopyContentId() {
    $(".copy-content-id").on("click", function () {
        var $temp = $("<input>");
        $("body").append($temp);
        $temp.val($(".content-uuid").text()).select();
        document.execCommand("copy");
        $temp.remove();
        NoticeSuccess("copied");
    });
}

function getTypeRoll(elementBoxRoll) {
    let type;
    if (elementBoxRoll.hasClass("pre-roll-break")) {
        type = "preroll";
    } else if (elementBoxRoll.hasClass("mid-roll-break")) {
        type = "midroll";
    } else if (elementBoxRoll.hasClass("post-roll-break")) {
        type = "postroll";
    }
    return type;
}

function checkInputTimeBreak() {
    $('.box-config-adbreak').on("input", ".time_break", function () {
        let elementBoxAdTagUrl = $(this).closest(".box-adtag-url");
        if ($(this).val().length > 0) {
            elementBoxAdTagUrl.prev().find(".btn-done-roll").removeClass("-disabled").find("button").prop("disabled", false);
        } else {
            elementBoxAdTagUrl.prev().find(".btn-done-roll").addClass("-disabled").find("button").prop("disabled", true);
        }
    });
}

function checkSelectConfigAdBreaks() {
    function check() {
        let adBreaks = $('input[name="config_ad_break"]:checked').val();
        if (adBreaks === "1") {
            $(".add-break-manually").removeClass("d-none");
        } else {
            $(".add-break-manually").addClass("d-none");
        }
    }

    check();
    $('input[type=radio][name=config_ad_break]').change(function () {
        check();
    });
}

function clickAddRoll() {
    $('.box-config-adbreak').on("click", ".box-btn-add", function () {
        let elementBoxRoll = $(this).closest(".box-roll");
        let type = getTypeRoll(elementBoxRoll);
        if (type === "midroll"){
            $(htmlAddRoll(type)).insertAfter(elementBoxRoll);
            elementBoxRoll.find(".box-step-roll").addClass("_editing_z87zl_152");
            elementBoxRoll.find(".box-btn-add").addClass("d-none");
            elementBoxRoll.find(".box-btn-action").addClass("d-none");
            elementBoxRoll.find(".box-btn-edit").removeClass("d-none");

            elementBoxRoll.attr("data-action", "add");
            elementBoxRoll.next().find(".selectpicker").selectpicker("refresh");
        } else {
            elementBoxRoll.find(".box-btn-edit").addClass("d-none");
            elementBoxRoll.find(".box-btn-add").addClass("d-none");
            elementBoxRoll.find(".box-btn-action").removeClass("d-none");
            elementBoxRoll.find(".box-step-roll").removeClass().addClass("box-step-roll").addClass("_active_z87zl_147");
            elementBoxRoll.find("._inactive_iyx7e_22").removeClass("_inactive_iyx7e_22").addClass("_active_iyx7e_28");
            elementBoxRoll.attr("data-action", "done");
        }
    });
}

function clickCancelRoll() {
    $('.box-config-adbreak').on("click", ".btn-cancel-roll", function () {
        let elementBoxRoll = $(this).closest(".box-roll");
        elementBoxRoll.find(".box-btn-edit").addClass("d-none");

        let elementBoxAdTagUrl = elementBoxRoll.next();
        let action = elementBoxRoll.attr("data-action");
        if (action === "add") {
            elementBoxRoll.find(".box-step-roll").removeClass("_editing_z87zl_152");
            elementBoxRoll.find(".box-btn-add").removeClass("d-none");
            if (elementBoxAdTagUrl.hasClass("box-adtag-url")) {
                elementBoxAdTagUrl.remove();
            }
        } else if (action === "done") {
            elementBoxRoll.find(".box-step-roll").removeClass().addClass("box-step-roll").addClass("_active_z87zl_147");
            elementBoxRoll.find("._inactive_iyx7e_22").removeClass("_inactive_iyx7e_22").addClass("_active_iyx7e_28");
            elementBoxRoll.find(".box-btn-action").removeClass("d-none");

            if (elementBoxAdTagUrl.hasClass("box-adtag-url")) {

                let elementBoxAdTagUrlBackup = elementBoxAdTagUrl.next();
                if (elementBoxAdTagUrlBackup.hasClass("box-adtag-url-backup")) {
                    elementBoxAdTagUrlBackup.addClass("box-adtag-url").removeClass("box-adtag-url-backup");
                }
                elementBoxAdTagUrl.remove();
            }
        }
    });
}

function clickDoneRoll() {
    $('.box-config-adbreak').on("click", ".btn-done-roll", function () {
        if (!$(this).hasClass("-disabled")) {
            let elementBoxRoll = $(this).closest(".box-roll");
            elementBoxRoll.find(".box-btn-edit").addClass("d-none");
            elementBoxRoll.find(".box-btn-action").removeClass("d-none");
            elementBoxRoll.find(".box-step-roll").removeClass().addClass("box-step-roll").addClass("_active_z87zl_147");
            elementBoxRoll.find("._inactive_iyx7e_22").removeClass("_inactive_iyx7e_22").addClass("_active_iyx7e_28");

            let elementBoxAdTagUrl = elementBoxRoll.next();
            let action = elementBoxRoll.attr("data-action");
            if (action === "add") {
                if (elementBoxAdTagUrl.hasClass("box-adtag-url")) {
                    elementBoxAdTagUrl.addClass("d-none");
                    elementBoxRoll.attr("data-action", "done");
                }
            } else if (action === "done") {
                if (elementBoxAdTagUrl.hasClass("box-adtag-url")) {
                    elementBoxAdTagUrl.addClass("d-none");

                    let elementBoxAdTagUrlBackup = elementBoxAdTagUrl.next();
                    if (elementBoxAdTagUrlBackup.hasClass("box-adtag-url-backup")) {
                        elementBoxAdTagUrlBackup.remove();
                    }
                }
            }
            let selectedMode = elementBoxAdTagUrl.find(".break_mode").find("option:selected").val();
            let timeBreak = elementBoxAdTagUrl.find(".time_break").val();

            let text = "";
            if (selectedMode === "1") {
                text = "Seconds into video: " + timeBreak + "s";
            } else if (selectedMode === "2") {
                text = "Timecode: " + timeBreak;
            } else if (selectedMode === "3") {
                text = "% Of Video: " + timeBreak + "%";
            }
            elementBoxRoll.find(".text-adtag-url").text(text);
            let type = getTypeRoll(elementBoxRoll);
            checkAddBoxMidRoll(type);
        }
    });
}

function checkAddBoxMidRoll(type) {
    function addBoxMidRoll() {
        let elementLastMidRoll = $(".mid-roll-break").last();
        let elementLastMidRollBoxAdTagUrl = elementLastMidRoll.next();
        if (elementLastMidRollBoxAdTagUrl.hasClass("box-adtag-url")) {
            let elementLastMidRollBoxAdTagUrlBackup = elementLastMidRollBoxAdTagUrl.next();
            if (elementLastMidRollBoxAdTagUrlBackup.hasClass("box-adtag-url-backup")) {
                $(htmlBoxRoll("midroll")).insertAfter(elementLastMidRollBoxAdTagUrlBackup);
                elementLastMidRollBoxAdTagUrlBackup.next().find(".selectpicker").selectpicker("refresh");
            } else {
                $(htmlBoxRoll("midroll")).insertAfter(elementLastMidRollBoxAdTagUrl);
                elementLastMidRollBoxAdTagUrl.next().find(".selectpicker").selectpicker("refresh");
            }
        } else {
            $(htmlBoxRoll("midroll")).insertAfter(elementLastMidRoll);
            elementLastMidRoll.next().find(".selectpicker").selectpicker("refresh");
        }
    }

    if (type === "midroll") {
        let isAdd = true;
        $(".mid-roll-break").each(function (index) {
            if ($(this).attr("data-action") !== "done" || index > 8) {
                isAdd = false;
            }
        });
        if (isAdd) {
            addBoxMidRoll();
        }
    }
}

function clickEditRoll() {
    $('.box-config-adbreak').on("click", ".btn-edit-roll", function () {
        let elementBoxRoll = $(this).closest(".box-roll");
        elementBoxRoll.find(".box-step-roll").removeClass().addClass("box-step-roll").addClass("_editing_z87zl_152");
        elementBoxRoll.find("._active_iyx7e_28").addClass("_inactive_iyx7e_22").removeClass("_active_iyx7e_28");
        elementBoxRoll.find(".box-btn-add").addClass("d-none");
        elementBoxRoll.find(".box-btn-action").addClass("d-none");
        elementBoxRoll.find(".box-btn-edit").removeClass("d-none");
        let elementBoxAdTagUrl = elementBoxRoll.next();
        if (elementBoxAdTagUrl.hasClass("box-adtag-url")) {
            elementBoxAdTagUrl.removeClass("d-none");
            let selectedMode = elementBoxAdTagUrl.find(".break_mode").find("option:selected").val();
            $(elementBoxAdTagUrl.clone()).insertAfter(elementBoxAdTagUrl);
            let elementBoxAdTagUrlBackup = elementBoxAdTagUrl.next();
            if (elementBoxAdTagUrlBackup.hasClass("box-adtag-url")) {
                elementBoxAdTagUrlBackup.addClass("box-adtag-url-backup").removeClass("box-adtag-url").addClass("d-none");
                elementBoxAdTagUrlBackup.find('.bootstrap-select').replaceWith(function () {
                    return $('select', this);
                });
                elementBoxAdTagUrlBackup.find(".break_mode").val(selectedMode);
                elementBoxAdTagUrlBackup.find('.selectpicker').selectpicker('render');
            }
        }

        $(".selectpicker").selectpicker("refresh");
    });
}

function clickDeleteRoll() {
    $('.box-config-adbreak').on("click", ".btn-del-roll", function () {
        let elementBoxRoll = $(this).closest(".box-roll");
        let elementBoxAdTagUrl = elementBoxRoll.next();
        if (elementBoxAdTagUrl.hasClass("box-adtag-url")) {
            elementBoxAdTagUrl.remove();
        }
        let type = getTypeRoll(elementBoxRoll);
        if (type === "midroll") {
            checkAddBoxMidRoll(type);
        } else {
            $(htmlBoxRoll(type)).insertAfter(elementBoxRoll);
        }
        elementBoxRoll.remove();

    });
}

function htmlAddRoll(typeRoll) {
    return `
<div class="_breakConfig_z87zl_166 pp_ad_schedules_7 -orientation-vertical -gutter-s -wrap hydrated w-100 box-adtag-url" data-type-roll="${typeRoll}">
    <div style="padding: 12px 20px;">
        <div class="pp_ad_schedules_7 -orientation-horizontal -justify-space-between -wrap hydrated">
            <div style="padding: 12px 0px;">
                <div class="pp_ad_schedules_7 -orientation-vertical -gutter-s -wrap hydrated">
                    <div>
                        <h4>Break Timing</h4>
                        <div class="_breakTiming_z87zl_175 pp_ad_schedules_7 -orientation-horizontal -gutter-s -wrap hydrated">
                            <div>
                                <div class="_numberInput_z87zl_124 wui-input -size-m -empty hydrated"
                                     type="number">
                                    <div class="wui-input__box">
                                        <input class="time_break" type="number" autocomplete="off" min="0" required="">
                                    </div>
                                </div>
                                <div class="border border-gray-500 rounded-1 select-ad-schedules">
                                    <select name="break_mode" class="break_mode selectpicker p-0 ">
                                        <option value="1" selected>Seconds into Video</option>
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

function htmlBoxRoll(type) {
    let classRoll = "";
    let textRoll = "";
    if (type === "preroll") {
        classRoll = "pre-roll-break";
        textRoll = "Preroll";
    } else if (type === "midroll") {
        classRoll = "mid-roll-break";
        textRoll = "Midroll";
    } else if (type === "postroll") {
        classRoll = "post-roll-break";
        textRoll = "Postroll";
    }
    return `
<div class="_breakItem_z87zl_75 pp_ad_schedules_7 -orientation-horizontal -align-center -wrap hydrated w-100 box-roll ${classRoll}">
                    <div class="" style="padding: 20px;flex-wrap: nowrap;">
                        <div class="box-step-roll">
                            <div data-test="show-break-label"
                                 class="_breakLabel_z87zl_137 pp_ad_schedules_7 -orientation-horizontal -align-center -justify-space-between -wrap hydrated">
                                <div>
                                    <div data-test="show-break-label-text">${textRoll}</div>
                                    <svg class="_stepper_iyx7e_16" width="40" height="6"
                                         fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <circle class="${(() => {if (type==="preroll") {return `_inactive_iyx7e_22`;} else {return ``;}})()}" cx="3" cy="3" r="3"
                                                data-test="show-preroll-step"></circle>
                                        <circle class="${(() => {if (type==="midroll") {return `_inactive_iyx7e_22`;} else {return ``;}})()}" cx="20" cy="3"
                                                r="3"
                                                data-test="show-midroll-step"></circle>
                                        <circle class="${(() => {if (type==="postroll") {return `_inactive_iyx7e_22`;} else {return ``;}})()}" cx="37" cy="3" r="3"
                                                data-test="show-postroll-step"></circle>
                                        <path stroke-width="2" d="M8 3h7M25 3h7"></path>
                                    </svg>
                                </div>
                            </div>
                        </div>
                        <p class="_font_83r5j_16 _s_83r5j_23 _baselineShift_83r5j_42 text-adtag-url"
                           style="font-family: 'Open Sans',Arial,sans-serif; color: #475470;"></p>
                        <div style="flex: 1 1 0%;"></div>
                        <div class="wui-button -type-icon hydrated box-btn-add">
                            <button type="button" class="btn-add-roll">
                                <div class="wui-icon wui-icon--dashboard_add -size-m hydrated">
                                    <svg xmlns="http://www.w3.org/2000/svg"
                                         viewBox="0 0 24 24"
                                         id="ds-icon-dashboard-add">
                                        <path d="M20 10h-6V4a2 2 0 0 0-4 0v6H4a2 2 0 0 0 0 4h6v6a2 2 0 0 0 4 0v-6h6a2 2 0 0 0 0-4Z"></path>
                                    </svg>
                                </div>
                            </button>
                        </div>
                        <div class="pp_ad_schedules_7 -orientation-horizontal -gutter-s -wrap hydrated box-btn-edit d-none">
                            <div>
                                <div data-test="cancel-editing-state" type="tertiary"
                                     class="wui-button -type-tertiary hydrated btn-cancel-roll">
                                    <button type="button"><span>Cancel</span></button>
                                </div>
                                <div data-test="save-break" type="secondary"
                                     class="wui-button -type-secondary -disabled hydrated btn-done-roll">
                                    <button type="button" disabled=""><span>Done</span>
                                    </button>
                                </div>
                            </div>
                        </div>
                        <div class="dropdown wui-more-menu -type-icon hydrated box-btn-action d-none">
                            <a class="dropdown-toggle wui-button -type-icon hydrated" type="button"
                               id="dropdownUserButton" data-toggle="dropdown"
                               aria-haspopup="true" aria-expanded="false" title="Email of publisher">
                                <div class="wui-icon wui-icon--dashboard_more -size-m hydrated">
                                    <svg xmlns="http://www.w3.org/2000/svg"
                                         viewBox="0 0 24 24"
                                         id="ds-icon-dashboard-more">
                                        <path d="M6 12a2 2 0 1 1-2-2 2 2 0 0 1 2 2Zm6-2a2 2 0 1 0 2 2 2 2 0 0 0-2-2Zm8 0a2 2 0 1 0 2 2 2 2 0 0 0-2-2Z"></path>
                                    </svg>
                                </div>
                            </a>
                            <div class="dropdown-menu" aria-labelledby="dropdownUserButton">
                                ${(() => {if (type==="midroll") {return `<a class="dropdown-item btn-edit-roll" type="button">Edit</a>`;} else {return ``;}})()}
                                <a class="dropdown-item text-danger btn-del-roll" type="button">Delete</a>
                            </div>
                        </div>
                    </div>
                </div>
`;
}

function checkRedirect() {
    $("form").on("change", "input", function () {
        if ($(this).attr("name")) {
            IsCheckRedirect = true;
        }
    });
    $("form").on("change", "select", function () {
        if ($(this).attr("name")) {
            IsCheckRedirect = true;
        }
    });
    $("form").on("change", "textarea", function () {
        if ($(this).attr("name")) {
            IsCheckRedirect = true;
        }
    });

    $("a").click(function () {
        var href = $(this).attr("href");
        if ($(this).attr("target") === "_blank" || href === "javascript:void(0)" || href === "#" || href.charAt(0) === "#") {
            return true;
        }
        if (href && IsCheckRedirect) {
            swal("apps.valueimpression.com says\n" +
                "Changes that you made may not be saved.", {
                className: "red-bg",
                buttons: {
                    cancel: "Cancel",
                    catch: {
                        text: "OK",
                        value: "ok",
                    },
                },
            }).then((value) => {
                switch (value) {
                    case "ok":
                        window.location = href;
                        break;

                    default:
                        return false;
                }
            });
        } else {
            return true;
        }
        return false;
    });
}

function firstLoad() {
    CheckImgSelected = $(".box-upload-thumb").data("check-img-selected") === true;
    LoadVideo();
}

function adBreak() {
    $(".btn-add-break").on("click", function () {
        $(".list-ad-breaks").append(htmlAdBreak(indexAdBreak));
        indexAdBreak++;
    });

    $(".list-ad-breaks").on("click", ".box-c .remove-ad-break", function () {
        $(this).closest(".box-c").remove();
        refreshListAdBreak();
    });
}

function handleTags() {
    _donetyping($(".input-tags"), 10, function () {
        if ($(".input-tags").val() !== "") {
            $(".add-tag").removeClass("cursor-not-allowed");
        } else (
            $(".add-tag").addClass("cursor-not-allowed")
        );
    });

    //add tag
    $(".add-tag").click(function () {
        var tag = $(".input-tags").val();
        if (!tag) {
            return;
        }
        var html = `<div class="av2 tag-item">
                        <div style="padding: 0 16px;display: inline-flex;">
                            <span class="p-0">` + tag + `</span>
                             <a href="javascript:void(0)" class="remove-tag">
                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" id="ds-icon-dashboard-cross">
                                    <path d="M14.83 12l6.58-6.59a2 2 0 0 0-2.82-2.82L12 9.17 5.41 2.59a2 2 0 0 0-2.82 2.82L9.17 12l-6.58 6.59a2 2 0 1 0 2.82 2.82L12 14.83l6.59 6.58a2 2 0 0 0 2.82-2.82z">
                                    </path>
                                </svg>
                            </a>
                        </div>
                    </div>`;
        $("#list-tags").prepend(html);
        $(".input-tags").val("");
        $(".add-tag").addClass("cursor-not-allowed");
    });
    document.getElementById('input-tags').addEventListener("keyup", function (e) {
        e.preventDefault();
        if (e.keyCode === 13) {
            $(".add-tag").click();
        }
    });

    // remove tag
    $("#list-tags").on('click', '.remove-tag', function () {
        $(this).closest(".av2").remove();
    });
}

function LoadVideo() {
    let video_url = $("#video_url").data("url");
    let id = $("#id").val();
    $(".preview-video").html("");
    $(".preview-video").append(`
                     <video id="video_${id}" src="${video_url}" style="max-height:535px;max-width: 870px; width: 100%;height: 100%;" controls="controls" type="video/mp4"></video>`);
    if (Hls.isSupported()) {
        var video = document.getElementById(`video_${id}`);
        var hls = new Hls();
        // bind them together
        hls.attachMedia(video);
        hls.on(Hls.Events.MEDIA_ATTACHED, function () {
            // console.log('video and hls.js are now bound together !');
            hls.loadSource(`${video_url}`);
            hls.on(Hls.Events.MANIFEST_PARSED, function (event, data) {
                // console.log(
                //     'manifest loaded, found ' + data.levels.length + ' quality level'
                // );
            });
        });
        // video.play();
    }
}

function eventClickAndChangeUpload() {
    $("#upload_video").on("click", function () {
        $("#box-upload-video").removeClass("is-invalid").next(".invalid-feedback").empty();
    });

    $("#upload_thumb").on("click", function () {
        $("#box-upload-thumb").removeClass("is-invalid").next(".invalid-feedback").empty();
    });

    $("#upload_video").on("change", function (e) {
        $("#loading_video").removeClass("d-none");
        Validate($(this), _validFileVideoExtensions, e, "video");
    });
    $("#upload_thumb").on("change", function (e) {
        $("#loading_thumb").removeClass("d-none");
        Validate($(this), _validFileImageExtensions, e, "img");
    });
}

function UploadFile(event, element, typeUpload) {
    var fd = new FormData();
    var file = event.target.files[0];
    if (file.size > maxSize) {
        new AlertError("You uploaded file over 10mb, please choose another file!");
        $("#loading_video").addClass("d-none");
        $("#loading_thumb").addClass("d-none");
        return;
    }
    fd.append('file', file);
    switch (typeUpload) {
        case"video":
            IsUploadVideoFinish = false;
            let blobURL = URL.createObjectURL(file);
            $(".preview-video").html("");
            $(".preview-video").append(`
                     <video src="${blobURL}" style="max-height:535px;max-width: 870px" controls="controls" type="video/mp4"></video>`);
            $.ajax({
                url: api_upload_video,
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
                    $("#loading_video").addClass("d-none");
                    AfterUpload(res.responseJSON, element, "video");
                }
            });
            break;
        case "img":
            IsUploadThumbFinish = false;
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
                    $("#loading_thumb").addClass("d-none");
                    AfterUpload(res.responseJSON, element, "img");
                }
            });
            break;
    }
    element.nextAll(".form-control").val(file.name);
}

function AfterUpload(response, element, typeUpload) {
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
                    inputElement.find("span.invalid-feedback").text(value.message);
                });
                new AlertError(response.errors[0].message);
            } else {
                new AlertError("Error!");
            }
            break;
        case "success":
            switch (typeUpload) {
                case "video":
                    // new AlertSuccess("Upload video successfully,automatically create thumb!");
                    $("#video_url").val(response.data_object.video_url);
                    $(".select-thumb").each(function (i) {
                        if (i === 0) {
                            if (CheckImgSelected) {
                                $(".select-thumb").parent().removeClass("selected");
                            }
                            $(this).parent().addClass("selected");
                            $("#thumb").val(response.data_object.thumb[i]);
                            $(".preview-thumb").attr("src", domainUpload + response.data_object.thumb[i]);
                            CheckImgSelected = true;
                        }
                        $(this).attr("src", domainUpload + response.data_object.thumb[i]);
                    });
                    $("#thumb").removeClass("is-invalid").next(".invalid-feedback").empty();
                    $("#box-upload-thumb").removeClass("is-invalid").next(".invalid-feedback").empty();
                    let nameFile = response.data_object.video_url.replace("/assets/video/", "").replace(".m3u8", "");
                    $("#name_file").attr("value", nameFile);
                    $("#duration").val(response.data_object.duration);
                    IsUploadVideoFinish = true;
                    break;
                case "img":
                    // new AlertSuccess("Upload thumb successfully!");
                    $("#thumb").val(response.data_object.thumb[0]);
                    $(".preview-thumb").removeClass("d-none").attr("src", domainUpload + response.data_object.thumb[0]);
                    $("#thumb").removeClass("is-invalid").next(".invalid-feedback").empty();
                    $(".select-thumb").parent().removeClass("selected");
                    CheckImgSelected = false;
                    IsUploadThumbFinish = true;
                    break;
            }
            break;
        default:
            new AlertError("Undefined");
            break;
    }
}


function Validate(element, validFile, event, typeUpload) {
    let oInput = element[0];
    if (oInput.type === "file") {
        let sFileName = oInput.value;
        if (sFileName.length > 0) {
            let blnValid = false;
            for (let j = 0; j < validFile.length; j++) {
                let sCurExtension = validFile[j];
                if (sFileName.substr(sFileName.length - sCurExtension.length, sCurExtension.length).toLowerCase() === sCurExtension.toLowerCase()) {
                    blnValid = true;
                    UploadFile(event, element, typeUpload);
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

function GetTags() {
    var tags = [];
    $("#list-tags").find(".tag-item").each(function () {
        tags.push($(this).find('span').text());
    });
    return tags;
}

function getAdBreak() {
    let adBreaks = [];
    let actionPreRoll = $(".pre-roll-break").attr("data-action");
    if (actionPreRoll === "done"){
        let adBreak = {};
        adBreak.type = "preroll";
        adBreak.break_mode = 1;
        adBreak.time_break = "";
        adBreaks.push(adBreak);
    }
    $(".box-adtag-url").each(function () {
        let adBreak = {};
        adBreak.type = $(this).attr("data-type-roll");
        adBreak.break_mode = parseInt($(this).find(".break_mode").find("option:selected").val());
        adBreak.time_break = $(this).find(".time_break").val();
        adBreaks.push(adBreak);
    });
    let actionPostRoll = $(".post-roll-break").attr("data-action");
    if (actionPostRoll === "done"){
        let adBreak = {};
        adBreak.type = "postroll";
        adBreak.break_mode = 1;
        adBreak.time_break = "";
        adBreaks.push(adBreak);
    }
    return adBreaks;
}

function SubmitFormEdit(formID, functionCallback, ajxURL = "") {
    let formElement = $("#" + formID);
    formElement.find("input").on("click change blur", function (e) {
        let inputElement = $(this);
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    formElement.find(".list-ad-breaks").on("input", "input", function (e) {
        let inputElement = $(this);
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    formElement.find("textarea").on("click change blur", function (e) {
        let inputElement = $(this);
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    formElement.find(".submit").on("click", function (e) {
        e.preventDefault();
        const buttonElement = $(this);
        const submitButtonText = buttonElement.text();
        const submitButtonTextLoading = "Loading...";
        var postData = formElement.serializeObject();
        postData.id = parseInt(postData.id);
        // postData.category = parseInt(postData.category)
        postData.duration = parseInt(postData.duration);
        postData.channels = parseInt(postData.channels);
        postData.video_type = parseInt(postData.video_type);
        postData.config_ad_break = parseInt(postData.config_ad_break);
        // if (postData.link_video !== "") {
        //     postData.video_url = postData.link_video;
        // }

        postData.video_url = $("#video_url").val();
        if ($(".list-ad-breaks .box-c .bidder-box").length === 1) {
            postData.ad_start_time = [postData.ad_start_time];
        }

        if (typeof postData.tag === 'string') {
            postData.tag = [postData.tag];
        }
        postData.keyword = GetTags();
        postData.ad_breaks = getAdBreak();
        // if (typeof postData.keyword === 'string') {
        //     postData.keyword = [postData.keyword];
        // }
        if (!IsUploadVideoFinish) {
            new AlertError("Video is processing, please wait a moment!");
            return;
        }
        if (!IsUploadThumbFinish) {
            new AlertError("Thumb is processing, please wait a moment!");
            return;
        }
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

function Edit(response, formElement) {
    switch (response.status) {
        case "error":
            if (response.errors.length === 1 && response.errors[0].id === "") {
                new AlertError(response.errors[0].message);
            } else if (response.errors.length > 0) {
                $.each(response.errors, function (key, value) {
                    if ($("#" + value.id).length) {
                        let inputElement = $("#" + value.id);
                        inputElement.addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
                        inputElement.find("span.invalid-feedback").text(value.message);
                        inputElement.addClass("is-invalid").closest(".bootstrap-select").addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
                    } else {
                        let inputElement = $("." + value.id);
                        inputElement.each(function () {
                            if (!$(this).val()) {
                                $(this).addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
                            }
                        });
                        // if ( !inputElement.val() ){
                        //     inputElement.addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
                        // }
                    }
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
            break;
        case "success":
            NoticeSuccess("Content has been updated successfully")
            IsCheckRedirect = false;
            // setTimeout(function () {
            //     window.location.replace("/content");
            // }, 1500);
            break;
        default:
            new AlertError("Undefined");
            break;
    }
}
