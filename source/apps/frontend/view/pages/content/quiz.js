module.exports = {
    firstLoad, getDataQuestion
}

let _validFileImageExtensions = [".jpg", ".jpeg", ".bmp", ".gif", ".png"];
let _validFileVideoExtensions = [".m4v", ".avi", ".mpg", ".mp4", ".3gp"];
let api_upload_video = "https://ul.pubpowerplatform.io/api/v1/video"
let api_upload_image = "https://ul.pubpowerplatform.io/api/v1/image"
let domainUpload = "https://ul.pubpowerplatform.io"
let IsUploadThumbFinish = true

let Number = 1
let maxSize = 15 * 1024 * 1024

$(document).ready(function () {
    // $("#tag").select2({
    //     tags: true,
    //     "language": {
    //         "noResults": function () {
    //             return "Please add tag and enter";
    //         }
    //     },
    //     escapeMarkup: function (markup) {
    //         return markup;
    //     }
    // })
})

function getDataQuestion(el) {
    const data = {};
    data.id = parseInt(el.attr("data-question"))
    data.title = el.find('input[name="question_title"]').val()
    data.background_type = parseInt(el.find('select[name="background_type"]').val())
    switch (data.background_type) {
        case 2:
            data.background = el.find('input[name="background_upload"]').val()
            break
        case 3:
            data.background = el.find('input[name="background_link"]').val()
            break
        default:
            data.background = el.find('input[name="background_color"]').val()
            break
    }
    data.type = parseInt(el.find('select[name="type"]').val())
    data.answer = []
    el.find('input[name="answer[]"]').each(function () {
        if ($(this).val()) {
            data.answer.push($(this).val())
        }
    })
    // data.answer = el.find('input[name="answer[]"]').val()
    data.picture_type = parseInt(el.find('select[name="picture_type"]').val())
    switch (data.picture_type) {
        case 2:
            data.picture = el.find('input[name="picture_upload"]').val()
            break
        case 3:
            data.picture = el.find('input[name="picture_link"]').val()
            break
        default:
            data.picture = el.find('input[name="picture_color"]').val()
            break
    }
    return data
}

function firstLoad() {
    $("#QuestionsCard").on("mouseover", ".align-items-center", function () {
        const number = numberAnswer($(this));
        maximumAnswer($(this), number)
        $("#QuestionsCard").find('.btn-add-answer').closest('.align-items-center').addClass('d-none')
        $("#QuestionsCard").find('.btn-remove-answer').closest('.align-items-center').addClass('d-none')
        $(this).find('.align-items-center').removeClass('d-none')
    })
    $("#QuestionsCard").on("mouseout", ".align-items-center", function () {
        const number = numberAnswer($(this));
        maximumAnswer($(this), number)
        $("#QuestionsCard").find('.btn-add-answer').closest('.align-items-center').addClass('d-none')
        $("#QuestionsCard").find('.btn-remove-answer').closest('.align-items-center').addClass('d-none')
    })

    refreshListQuestion()
    loadQuestion();
    eventClickUpload();
}

function loadQuestion() {
    $('.list-questions').on("change", ".change-color", function () {
        $(this).closest("div").find('.input-color').val($(this).val())
        console.log($(this).val());
    })
    $(".list-questions").on("change", ".background_type", function () {
        BackgroundChange($(this))
    })
    $(".list-questions").on("change", ".picture_type", function () {
        PictureChange($(this))
    })
    $(".list-questions").on("change", ".question-type", function () {
        const select = $(this).val();
        switch (select) {
            case '1':
                $(this).closest('.bidder-box').find('.multiple-choice').removeClass('d-none')
                $(this).closest('.bidder-box').find('.picture-choice').addClass('d-none')
                break;
            default:
                // console.log($('#picture').closest('.align-items-center').find('input').width());
                $('.picture_type').select2({width: $('#title').outerWidth() + "px"});
                $(this).closest('.bidder-box').find('.multiple-choice').addClass('d-none')
                $(this).closest('.bidder-box').find('.picture-choice').removeClass('d-none')
                break;
        }

        // xoá toàn bộ text đc điền vào input answer
        $(this).closest(".question-item").find("input[name='answer[]']").val("")
        $(this).closest(".question-item").find("input[name='answer[]']").val("")
    })
    // Add question
    $(".btn-add-question").on("click", function () {
        $(".list-questions").append(htmlQuestion(Number))
        refreshListQuestion()
        Number++
    })
    // Remove Question
    $(".list-questions").on("click", ".box-c .remove-question", function () {
        $(this).closest(".box-c").remove()
        refreshListQuestion()
    })
    // Add Answer
    $("#QuestionsCard").on("click", ".btn-remove-answer", function () {
        removeAnswer($(this))
    })
    // Remove Answer
    $("#QuestionsCard").on("click", ".btn-add-answer", function () {
        addAnswer($(this))
    })
}

function refreshListQuestion() {
    let NumberQuestion = 1;
    $('.list-questions > div.box-c > div.bidder-box').each(function () {
        let lbQuestion = "QUESTION " + NumberQuestion
        $(this).find(".lb-question").html(lbQuestion)
        $(this).find('input[name="question_title"]').attr("id", "question" + NumberQuestion + "_title")
        $(this).find('input[name="background_type"]').attr("id", "question" + NumberQuestion + "_background_type")
        $(this).find('input[name="background_color"]').closest('div').attr("id", "question" + NumberQuestion + "_background_color")
        $(this).find('input[name="background_upload"]').closest('.box-upload-thumb').attr("id", "question" + NumberQuestion + "_background_upload")
        $(this).find('input[name="background_link"]').attr("id", "question" + NumberQuestion + "_background_link")

        $(this).find('.answer-question').each(function () {
            $(this).find("input").addClass("question" + NumberQuestion + "_answer")
        })
        $(this).find('input[name="picture_type"]').closest('div').attr("id", "question" + NumberQuestion + "_picture_type")
        $(this).find('input[name="picture_color"]').closest('div').attr("id", "question" + NumberQuestion + "_picture_color")
        $(this).find('input[name="picture_upload"]').closest('.box-upload-thumb').attr("id", "question" + NumberQuestion + "_picture_upload")
        $(this).find('input[name="picture_link"]').attr("id", "question" + NumberQuestion + "_picture_link")
        // $(this).find('.box-upload-thumb > label').attr("for", "upload_picture" + NumberQuestion)
        NumberQuestion++
    })


    $('.answer-question').each(function () {
        let NumberAnswer = 1;
        $(this).find("div.input-answer").each(function () {
            let label = "ANSWER " + NumberAnswer
            $(this).find("label").html(label)
            NumberAnswer++
        })
    })
}

function BackgroundChange(el) {
    const select = el.val();
    switch (select) {
        case '2':
            el.closest('.bidder-box').find('.background-link').addClass('d-none')
            el.closest('.bidder-box').find('.background-upload').removeClass('d-none')
            el.closest('.bidder-box').find('.background-color').addClass('d-none')
            break;
        case '3':
            el.closest('.bidder-box').find('.background-link').removeClass('d-none')
            el.closest('.bidder-box').find('.background-upload').addClass('d-none')
            el.closest('.bidder-box').find('.background-color').addClass('d-none')
            break;
        default:
            el.closest('.bidder-box').find('.background-link').addClass('d-none')
            el.closest('.bidder-box').find('.background-upload').addClass('d-none')
            el.closest('.bidder-box').find('.background-color').removeClass('d-none')
            break;
    }
}

function eventClickUpload() {
    $(".list-questions").on("click", ".upload_thumb", function () {
        $(this).closest(".box-upload-thumb").removeClass("is-invalid").next(".invalid-feedback").empty();
    })
    $(".list-questions").on("change", ".upload_thumb", function (e) {
        $(this).closest('.box-upload-thumb').find('.loading_thumb').removeClass("d-none")
        Validate($(this), _validFileImageExtensions, e, "img")
    })
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
                    UploadFile(event, element, typeUpload)
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

function UploadFile(event, element, typeUpload) {
    const fd = new FormData();
    const file = event.target.files[0];
    if (file.size > maxSize) {
        new AlertError("You uploaded file over 10mb, please choose another file!");
        element.closest(".box-upload-thumb").find(".loading_thumb").addClass("d-none")
        element.closest(".box-upload-video").find(".loading_video").addClass("d-none")
        return
    }
    fd.append('file', file);
    switch (typeUpload) {
        case"video":
            let IsUploadVideoFinish = false
            let blobURL = URL.createObjectURL(file);
            $(".preview-video").html("")
            $(".preview-video").append(`<video src="${blobURL}" style="height:150px;width:300px" controls="controls" type="video/mp4"></video>`)
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
                    $("#loading_video").addClass("d-none")
                    AfterUpload(res.responseJSON, element, "video");
                }
            });
            break
        case "img":
        case "thumb":
            IsUploadThumbFinish = false
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
                    element.closest(".box-upload-thumb").find(".loading_thumb").addClass("d-none")
                    AfterUpload(res.responseJSON, element, "img");
                }
            });
            break
    }
    element.nextAll(".form-control").val(file.name)
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
            break
        case "success":
            switch (typeUpload) {
                case "img":
                case "thumb":
                    // new AlertSuccess("Upload thumb successfully!");
                    console.log(element.closest(".box-upload-thumb").find(".thumb"));
                    element.closest(".box-upload-thumb").find(".thumb").val(response.data_object.thumb[0])
                    element.closest(".box-upload-thumb").find(".preview-thumb").removeClass("d-none").attr("src", domainUpload + response.data_object.thumb[0])
                    element.closest(".box-upload-thumb").find(".thumb").removeClass("is-invalid").next(".invalid-feedback").empty();
                    CheckImgSelected = false;
                    IsUploadThumbFinish = true
                    break
            }
            break
        default:
            new AlertError("Undefined");
            break
    }
}


function PictureChange(el) {
    var select = el.val()
    switch (select) {
        case '2':
            el.closest('.bidder-box').find('.picture-link').addClass('d-none')
            el.closest('.bidder-box').find('.picture-upload').removeClass('d-none')
            el.closest('.bidder-box').find('.picture-color').addClass('d-none')
            break;
        case '3':
            el.closest('.bidder-box').find('.picture-link').removeClass('d-none')
            el.closest('.bidder-box').find('.picture-upload').addClass('d-none')
            el.closest('.bidder-box').find('.picture-color').addClass('d-none')
            break;
        default:
            el.closest('.bidder-box').find('.picture-link').addClass('d-none')
            el.closest('.bidder-box').find('.picture-upload').addClass('d-none')
            el.closest('.bidder-box').find('.picture-color').removeClass('d-none')
            break;
    }
}

function numberAnswer(e) {
    var number = 0
    e.closest(".answer-question").find(".input-answer").each(function () {
        number += 1
    })
    return number
}

function removeAnswer(e) {
    var AnswerQuestion = e.closest(".answer-question")
    var number = numberAnswer(e)
    e.closest(".input-answer").remove()
    maximumAnswer(e, number - 1)

    var number_answer = 0
    AnswerQuestion.find(".input-answer").each(function () {
        number_answer += 1
        $(this).find("label").text("Answer " + number_answer)
    })
}

function addAnswer(e) {
    var number_answer = numberAnswer(e) + 1
    var question = e.closest(".question-item").attr('data-number')
    var AnswerHTML = '<div class="row my-4 input-answer">\n' +
        '<div class="d-flex align-items-center">\n ' +
        '<label class="col-form-label form-label form-label-lg w-input-group">ANSWER ' + number_answer + '</label>\n' +
        '<input class="w-50-custom form-control question' + question + '_answer" name="answer[]" value="">\n ' +
        '<div class="ms-2 ps-0 d-flex align-items-center d-none">\n ' +
        '<button type="button" class="btn btn-outline-warning btn-sm px-2 rounded-2 btn-add-answer d-none">\n ' +
        '<svg xmlns="http://www.w3.org/2000/svg" width="14" fill="currentColor" height="14" viewBox="0 0 14 14">\n ' +
        '<rect data-name="Icons/Tabler/Add background" width="14" height="14" fill="none"></rect>\n' +
        '<path d="M6.329,13.414l-.006-.091V7.677H.677A.677.677,0,0,1,.585,6.329l.092-.006H6.323V.677A.677.677,0,0,1,7.671.585l.006.092V6.323h5.646a.677.677,0,0,1,.091,1.348l-.091.006H7.677v5.646a.677.677,0,0,1-1.348.091Z" fill="currentColor"></path>\n' +
        '</svg>\n' +
        '</button>\n ' +
        '</div>\n' +
        '<div class="ms-2 ps-0 d-flex align-items-center d-none">\n' +
        '<button type="button" class="btn btn-outline-danger btn-sm px-2 rounded-2 btn-remove-answer d-none">\n' +
        '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-trash" viewBox="0 0 16 16">\n' +
        '<path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z"></path>\n ' +
        '<path fill-rule="evenodd" d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z"></path>\n' +
        '</svg>\n' +
        '</button>\n' +
        '</div>\n' +
        '<span class="form-text ms-2 w-auto invalid-feedback"></span>\n' +
        '</div>\n' +
        '</div>';
    e.closest(".answer-question").append(AnswerHTML)
    maximumAnswer(e, number_answer)
}

function maximumAnswer(e, number) {
    if (number >= 4) {
        e.closest('.answer-question').find('.btn-add-answer').addClass('d-none')
    } else {
        e.closest('.answer-question').find('.btn-add-answer').removeClass('d-none')
    }
    if (number <= 2) {
        e.closest('.answer-question').find('.btn-remove-answer').addClass('d-none')
    } else {
        e.closest('.answer-question').find('.btn-remove-answer').removeClass('d-none')
    }
}

function makePreview() {
    let postData = $("#editTemplate").serializeObject();
    postData = makeData(postData)
    $.ajax({
        url: "/player/template/preview",
        type: "POST",
        dataType: "JSON",
        contentType: "application/json",
        data: JSON.stringify(postData),
        success: function (json) {
            document.querySelector("#videocontainer").innerHTML = "";
            (vitag.Init = window.vitag.Init || []).push(function () {
                // console.log(JSON.stringify(config))
                viAPItag.initVliVideoContainer(json);
            });
        },
    });
}

function htmlQuestion(index) {
    return `
    <div class="box-c question-item" data-number="` + index + `" data-question="0">
        <hr class="mt-3 hr_custom bg-gray-400" style="margin: 0 -20px 0 -50px !important;">
        <div class="bidder-box">
            <div class="row my-4">
                <div class="d-flex align-items-center bidder-name">
                    <label class="col-form-label form-label form-label-lg w-input-group text-uppercase lb-question"> QUESTION ` + index + ` </label>
                    <a class="btn p-1 ms-2 remove-question">
                        <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" fill="currentColor" class="bi bi-x-lg" viewBox="0 0 16 16">
                            <path d="M1.293 1.293a1 1 0 0 1 1.414 0L8 6.586l5.293-5.293a1 1 0 1 1 1.414 1.414L9.414 8l5.293 5.293a1 1 0 0 1-1.414 1.414L8 9.414l-5.293 5.293a1 1 0 0 1-1.414-1.414L6.586 8 1.293 2.707a1 1 0 0 1 0-1.414z"></path>
                        </svg>
                    </a>
                </div>
            </div>
    
            <div class="row my-4">
                <div class="d-flex align-items-center">
                    <label class="col-form-label form-label form-label-lg w-input-group">Title</label>
                    <input type="text" class="form-control w-50-custom" name="question_title" id="question` + index + `_title">
                    <span class="form-text ms-2 w-auto invalid-feedback"></span>
                </div>
            </div>
            <div class="row my-4">
                <div class="d-flex align-items-center">
                    <label class="col-form-label form-label form-label-lg w-input-group">Background</label>
                    <select name="background_type" id="question` + index + `_background_type" class="form-select w-50-custom select2 background_type" data-minimum-results-for-search="-1">
                        <option value="1">Color</option>
                        <option value="2">Upload File</option>
                        <option value="3">Link Gifphy</option>
                    </select>
                </div>
            </div>
    
            <div class="row my-4 top-articles background-color">
                <div class="d-flex align-items-center">
                    <label class="col-form-label form-label form-label-lg w-input-group"></label>
                    <div class="d-flex w-50-custom">
                        <input type="color" class="form-control form-control-color border-gray-300 border-end-0 rounded-0 rounded-start change-color" value="#ffc107" title="Choose your color">
                        <input type="text" class="form-control border-gray-300 border-start-0 rounded-0 rounded-end input-color" name="background_color" id="question` + index + `_background_color" value="#ffc107">
                        <span class="form-text ms-2 w-auto invalid-feedback"></span>
                    </div>
                </div>
            </div>
            <div class="row my-4 background-upload d-none">
                <div class="d-flex align-items-center">
                    <label class="col-form-label form-label form-label-lg w-input-group"></label>
                    <div class="w-50-custom is-invalid box-upload-thumb">
                        <input type="file" title="Choose a image please" name="background-upload" class="w-50-custom form-control upload_thumb d-none" id="upload_background` + index + `" aria-describedby="background_url" aria-label="Upload" accept="iamge/*">
                        <div class="w-100 d-flex flex-row form-control p-0">
                            <label class="m-0 px-2 d-flex align-items-center border-end border-gray-300 bg-gray-200" for="upload_background` + index + `">Choose file</label>
                            <input class="border-0 form-control bg-white is-invalid thumb" style="width: 80%" name="background_upload" id="question` + index + `_background_upload" value="" readonly>
                            <span class="form-text ms-2 w-auto invalid-feedback"></span>
                        </div>
                        <div class="loading_thumb d-none">
                            <span class="spinner-border text-secondary" role="status"></span>
                            <span class="text-muted">Background is processing, please wait a moment!</span>
                        </div>
                          <div class="d-block ">
                        <img class="preview-thumb d-none mt-2" src="" alt=""/></div>
                    </div>
                    <span class="form-text ms-2 w-auto invalid-feedback"></span>
                    <input hidden id="name_file" name="name_file" value="">
                </div>
            </div>
            <div class="row my-4 background-link d-none">
                <div class="d-flex align-items-center">
                    <label class="col-form-label form-label form-label-lg w-input-group"></label>
                    <input class="w-50-custom form-control" name="background_link" value="" id="question` + index + `_background_link" placeholder="Link gifphy">
                    <span class="form-text ms-2 w-auto invalid-feedback"></span>
                </div>
            </div>
            <div class="row my-4">
                <div class="d-flex align-items-center">
                    <label class="col-form-label form-label form-label-lg w-input-group" for="tag">Type</label>
                    <select name="type" class="form-select w-50-custom select2 question-type" id="question` + index + `_type" data-minimum-results-for-search="-1">
                        <option value="1">Multiple Choice</option>
                        <option value="2">Picture Choice</option>
                    </select>
                </div>
            </div>
            <div class="multiple-choice answer-question">
                <div class="row my-4 input-answer">
                    <div class="d-flex align-items-center">
                        <label class="col-form-label form-label form-label-lg w-input-group">Answer 1</label>
                        <input class="w-50-custom form-control" name="answer[]" value="">
                        <div class="ms-2 ps-0 d-flex align-items-center ">
                            <button type="button" class="btn btn-outline-warning btn-sm px-2 rounded-2 btn-add-answer d-none">
                                <svg xmlns="http://www.w3.org/2000/svg" width="14" fill="currentColor" height="14" viewBox="0 0 14 14">
                                    <rect data-name="Icons/Tabler/Add background" width="14" height="14" fill="none"></rect>
                                    <path d="M6.329,13.414l-.006-.091V7.677H.677A.677.677,0,0,1,.585,6.329l.092-.006H6.323V.677A.677.677,0,0,1,7.671.585l.006.092V6.323h5.646a.677.677,0,0,1,.091,1.348l-.091.006H7.677v5.646a.677.677,0,0,1-1.348.091Z" fill="currentColor"></path>
                                </svg>
                            </button>
                        </div>
                        <div class="ms-2 ps-0 d-flex align-items-center d-none">
                            <button type="button" class="btn btn-outline-danger btn-sm px-2 rounded-2 btn-remove-answer">
                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-trash" viewBox="0 0 16 16">
                                    <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z"></path>
                                    <path fill-rule="evenodd" d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z"></path>
                                </svg>
                            </button>
                        </div>
                        <span class="form-text ms-2 w-auto invalid-feedback"></span>
                    </div>
                </div>
    
                <div class="row my-4 input-answer">
                    <div class="d-flex align-items-center">
                        <label class="col-form-label form-label form-label-lg w-input-group">Answer 2</label>
                        <input class="w-50-custom form-control" name="answer[]" value="">
                        <div class="ms-2 ps-0 d-flex align-items-center ">
                            <button type="button" class="btn btn-outline-warning btn-sm px-2 rounded-2 btn-add-answer d-none">
                                <svg xmlns="http://www.w3.org/2000/svg" width="14" fill="currentColor" height="14" viewBox="0 0 14 14">
                                    <rect data-name="Icons/Tabler/Add background" width="14" height="14" fill="none"></rect>
                                    <path d="M6.329,13.414l-.006-.091V7.677H.677A.677.677,0,0,1,.585,6.329l.092-.006H6.323V.677A.677.677,0,0,1,7.671.585l.006.092V6.323h5.646a.677.677,0,0,1,.091,1.348l-.091.006H7.677v5.646a.677.677,0,0,1-1.348.091Z" fill="currentColor"></path>
                                </svg>
                            </button>
                        </div>
                        <div class="ms-2 ps-0 d-flex align-items-center d-none">
                            <button type="button" class="btn btn-outline-danger btn-sm px-2 rounded-2 btn-remove-answer">
                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-trash" viewBox="0 0 16 16">
                                    <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z"></path>
                                    <path fill-rule="evenodd" d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z"></path>
                                </svg>
                            </button>
                        </div>
                       <span class="form-text ms-2 w-auto invalid-feedback"></span>                         
                    </div>
                </div>
            </div>
            <div class="picture-choice d-none">
                <div class="answer-question">
                    <div class="row my-4 input-answer">
                        <div class="d-flex align-items-center">
                            <label class="col-form-label form-label form-label-lg w-input-group">Answer 1</label>
                            <input class="w-50-custom form-control" name="answer[]" value="">
                            <div class="ms-2 ps-0 d-flex align-items-center ">
                                <button type="button" class="btn btn-outline-warning btn-sm px-2 rounded-2 btn-add-answer d-none">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="14" fill="currentColor" height="14" viewBox="0 0 14 14">
                                        <rect data-name="Icons/Tabler/Add background" width="14" height="14" fill="none"></rect>
                                        <path d="M6.329,13.414l-.006-.091V7.677H.677A.677.677,0,0,1,.585,6.329l.092-.006H6.323V.677A.677.677,0,0,1,7.671.585l.006.092V6.323h5.646a.677.677,0,0,1,.091,1.348l-.091.006H7.677v5.646a.677.677,0,0,1-1.348.091Z" fill="currentColor"></path>
                                    </svg>
                                </button>
                            </div>
                            <div class="ms-2 ps-0 d-flex align-items-center d-none">
                                <button type="button" class="btn btn-outline-danger btn-sm px-2 rounded-2 btn-remove-answer">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-trash" viewBox="0 0 16 16">
                                        <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z"></path>
                                        <path fill-rule="evenodd" d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z"></path>
                                    </svg>
                                </button>
                            </div>
                            <span class="form-text ms-2 w-auto invalid-feedback"></span>                            
                        </div>
                    </div>
    
                    <div class="row my-4 input-answer">
                        <div class="d-flex align-items-center">
                            <label class="col-form-label form-label form-label-lg w-input-group">Answer 2</label>
                            <input class="w-50-custom form-control" name="answer[]" value="">
                            
                            <div class="ms-2 ps-0 d-flex align-items-center ">
                                <button type="button" class="btn btn-outline-warning btn-sm px-2 rounded-2 btn-add-answer d-none">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="14" fill="currentColor" height="14" viewBox="0 0 14 14">
                                        <rect data-name="Icons/Tabler/Add background" width="14" height="14" fill="none"></rect>
                                        <path d="M6.329,13.414l-.006-.091V7.677H.677A.677.677,0,0,1,.585,6.329l.092-.006H6.323V.677A.677.677,0,0,1,7.671.585l.006.092V6.323h5.646a.677.677,0,0,1,.091,1.348l-.091.006H7.677v5.646a.677.677,0,0,1-1.348.091Z" fill="currentColor"></path>
                                    </svg>
                                </button>
                            </div>
                            <div class="ms-2 ps-0 d-flex align-items-center d-none">
                                <button type="button" class="btn btn-outline-danger btn-sm px-2 rounded-2 btn-remove-answer">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-trash" viewBox="0 0 16 16">
                                        <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z"></path>
                                        <path fill-rule="evenodd" d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z"></path>
                                    </svg>
                                </button>
                            </div>
                            <span class="form-text ms-2 w-auto invalid-feedback"></span>
                        </div>
                    </div>
                </div>
    
                <div class="row my-4">
                    <div class="d-flex align-items-center">
                        <label class="col-form-label form-label form-label-lg w-input-group">Picture</label>
                        <select name="picture_type" id="picture" class="form-select w-50-custom picture_type" id="question` + index + `_picture_type" data-minimum-results-for-search="-1">
                            <option value="1">Color</option>
                            <option value="2">Upload File</option>
                            <option value="3">Link Gifphy</option>
                        </select>
                    </div>
                </div>
    
                <div class="row my-4 top-articles picture-color">
                    <div class="d-flex align-items-center">
                        <label class="col-form-label form-label form-label-lg w-input-group"></label>
                        <div class="d-flex w-50-custom">
                            <input type="color" class="form-control form-control-color border-gray-300 border-end-0 rounded-0 rounded-start change-color" value="#ffc107" title="Choose your color">
                            <input type="text" class="form-control border-gray-300 border-start-0 rounded-0 rounded-end input-color" name="picture_color" id="question` + index + `_picture_color" value="#ffc107">
                        </div>
                        <span class="form-text ms-2 w-auto invalid-feedback"></span>
                    </div>
                </div>
                <div class="row my-4 picture-upload d-none">
                    <div class="d-flex align-items-center">
                        <label class="col-form-label form-label form-label-lg w-input-group"></label>
                        <div class="w-50-custom box-upload-thumb">
                            <input type="file" title="Choose a thumb please" id="upload_picture` + index + `" class="w-50-custom form-control upload_thumb d-none" aria-describedby="picture_url" aria-label="Upload" accept="image/*">
                            <div class="w-100 d-flex flex-row form-control p-0">
                                <label class="m-0 px-2 d-flex align-items-center border-end border-gray-300 bg-gray-200" for="upload_picture` + index + `">Choose file</label>
                                <input class="border-0 form-control bg-white thumb" style="width: 80%" name="picture_upload" id="question` + index + `_picture_upload" value="" readonly>
                                <span class="form-text ms-2 w-auto invalid-feedback"></span>
                            </div>
                            <div class="loading_thumb d-none">
                                <span class="spinner-border text-secondary" role="status"></span>
                                <span class="text-muted">Picture is processing, please wait a moment!</span>
                            </div>
                            <div class="d-block "><img class="preview-thumb d-none mt-2" src="" alt=""/></div>
                        </div>
                        <span class="form-text ms-2 w-auto invalid-feedback"></span>
                        <input hidden id="name_file" name="name_file" value="">
                    </div>
                </div>
                <div class="row my-4 picture-link d-none">
                    <div class="d-flex align-items-center">
                        <label class="col-form-label form-label form-label-lg w-input-group"></label>
                        <input class="w-50-custom form-control" name="picture_link" id="question` + index + `_picture_link" value="" placeholder="Link gifphy">
                        <span class="form-text ms-2 w-auto invalid-feedback"></span>
                    </div>
                </div>
            </div>
        </div>
    </div>
    `
}
