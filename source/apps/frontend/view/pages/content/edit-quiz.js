let _validFileImageExtensions = [".jpg", ".jpeg", ".bmp", ".gif", ".png"];
let _validFileVideoExtensions = [".m4v", ".avi", ".mpg", ".mp4", ".3gp"];
let api_upload_video = "https://ul.pubpowerplatform.io/api/v1/video";
let api_upload_image = "https://ul.pubpowerplatform.io/api/v1/image";
let domainUpload = "https://ul.pubpowerplatform.io";
// let api_upload_video = "http://127.0.0.1:8543/api/v1/video";
// let api_upload_image = "http://127.0.0.1:8543/api/v1/image";
// let domainUpload = "http://127.0.0.1:8543";
let IsUploadThumbFinish = true;

let number = 1;
let maxSize = 15 * 1024 * 1024;

$(document).ready(function () {
    firstLoad();
    SubmitFormCreate("formQuiz", Added, "/content/quiz/edit");

    $(".quiz-tab").on("click", ".nav-link", function () {
        $(".quiz-tab").find(".nav-link").removeClass("pp-4");
        $(this).addClass("pp-4");
        var tab = $(this).attr("data-tab");
        if (tab != "1") {
            $(this).addClass("at-1");
        } else {
            $(".quiz-tab").find('.at-1').removeClass('at-1');
        }
    });

    $("#content_type").on("change", function (e) {
        checkContentType();
    });
});

function firstLoad() {
    // load image
    $(".show-image-answer").each(function (index, element) {
        if ( $(this).attr("data-link") ){
            $(this).attr("src", domainUpload + $(this).attr("data-link"));
        }
    });
    refreshListQuestion();
    loadQuestion();
    checkContentType();
    clickAddImage();
    clickAddIllustration();
}

function clickAddIllustration() {
    $("#upload_illustration").on("change", function (e) {
        validateAndUploadLogo($(this), _validFileImageExtensions, e);
    });
}

function checkContentType() {
    let contentType = $("#content_type").val();
    if (contentType === "2") {
        $(".box-answer-quiz2").addClass("d-none");
        $(".box-answer-quiz1").removeClass("d-none");
    } else {
        $(".box-answer-quiz1").addClass("d-none");
        $(".box-answer-quiz2").removeClass("d-none");
    }
}

function checkAnswer(element) {
    if ($(element).is(":checked")) {
        $(element).closest(".box-answer").find(".answer-correct").prop('checked', false);
    }
    $(element).prop('checked', true);
}

function loadQuestion() {
    // Add question
    $(".btn-add-question").on("click", function () {
        $(".list-questions").append(htmlQuestion(number));
        $(".selectpicker").selectpicker("refresh");
        refreshListQuestion();
        checkContentType();
        number++;
    });
    // Remove Question
    $(".list-questions").on("click", ".box-c .remove-question", function () {
        $(this).closest(".box-c").remove();
        refreshListQuestion();
    });
    // Add Answer
    $(".list-questions").on("click", ".box-c .btn-add-answer-quiz1", function () {
        let NumberAnswer = 0;
        $(this).parent().prev(".box-answer-quiz1").find('.answer-item').each(function () {
            NumberAnswer++;
        });
        if (NumberAnswer < 4) {
            $(this).parent().prev(".box-answer-quiz1").append(htmlAnswer(number));
            refreshListQuestion();
        }
    });
    // Remove Answer
    $(".list-questions").on("click", ".box-c .box-answer-quiz1 .btn-remove-answer", function () {
        let NumberAnswer = 0;
        $(this).closest(".box-answer-quiz1").find('.answer-item').each(function () {
            NumberAnswer++;
        });
        if (NumberAnswer > 2) {
            $(this).closest(".answer-item").remove();
            refreshListQuestion();
        }
    });

    // Check correct answer
    $(".list-questions").on("change", ".box-c .answer-correct", function () {
        checkAnswer(this);
    });
}

function refreshListQuestion(data) {
    let NumberQuestion = 1;
    $('.list-questions > div.box-c > div.question-box').each(function (index, element) {
        let lbQuestion = "<span class='ms-3'> QUESTION " + NumberQuestion + "</span>";
        $(this).find(".lb-question").html(lbQuestion);
        $(this).closest(".question-item").attr("data-question", NumberQuestion);
        if (typeof data !== 'undefined') {
            let questionId = data.questions[index].id;
            $(this).closest(".question-item").find(".question-id").val(questionId);
        }
        let checkCorrectQuiz1 = false;
        let NumberAnswerQuiz1 = 1;
        $(this).find('.box-answer-quiz1 .answer-item').each(function (indexAnswer, elementAnswer) {
            let label = "Answer " + NumberAnswerQuiz1;
            $(this).find(".answer-name").html(label);
            $(this).attr("data-answer", NumberAnswerQuiz1);
            let idAnswer = "quiz1_question_" + NumberQuestion + "_answer_" + NumberAnswerQuiz1;
            $(this).find(".answer-correct").attr("id", idAnswer);
            $(this).find(".answer-name").attr("for", idAnswer);
            if (typeof data !== 'undefined') {
                let answerId = data.questions[index].answers[indexAnswer].id;
                $(this).find(".answer-id").val(answerId);
            }
            if (!checkCorrectQuiz1) {
                checkCorrectQuiz1 = $(this).find(".answer-correct").is(":checked");
            }
            NumberAnswerQuiz1++;
        });
        if (!checkCorrectQuiz1) {
            $(this).find(".box-answer-quiz1 .answer-item[data-answer=1] .answer-correct").prop("checked", true);
        }
        let checkCorrectQuiz2 = false;
        let NumberAnswerQuiz2 = 1;
        $(this).find('.box-answer-quiz2 .answer-item').each(function (indexAnswer, elementAnswer) {
            let label = "Answer " + NumberAnswerQuiz2;
            $(this).find(".answer-name").html(label);
            $(this).attr("data-answer", NumberAnswerQuiz2);
            let idAnswer = "quiz2_question_" + NumberQuestion + "_answer_" + NumberAnswerQuiz2;
            $(this).find(".answer-correct").attr("id", idAnswer);
            $(this).find(".answer-name").attr("for", idAnswer);
            if (typeof data !== 'undefined') {
                let answerId = data.questions[index].answers[indexAnswer].id;
                $(this).find(".answer-id").val(answerId);
            }
            if (!checkCorrectQuiz2) {
                checkCorrectQuiz2 = $(this).find(".answer-correct").is(":checked");
            }
            NumberAnswerQuiz2++;
        });
        if (!checkCorrectQuiz2) {
            $(this).find(".box-answer-quiz2 .answer-item[data-answer=1] .answer-correct").prop("checked", true);
        }
        NumberQuestion++;
    });

}

function getQuestions() {
    let questions = [];
    let contentType = $("#content_type").val();
    $(".list-questions .question-item").each(function () {
        let question = {};
        question.id = parseInt($(this).find(".question-id").val());
        question.question = $(this).find(".text_question").val();
        question.background = $(this).find("#background").val();
        let answers = [];
        if (contentType === "2") {
            $(this).find(".box-answer-quiz1 .answer-item").each(function () {
                let answer = {};
                answer.id = parseInt($(this).find(".answer-id").val());
                answer.correct = $(this).find(".answer-correct").is(":checked");
                answer.text = $(this).find(".text_answer").val();
                answers.push(answer);
            });
        } else {
            $(this).find(".box-answer-quiz2 .answer-item").each(function () {
                let answer = {};
                answer.id = parseInt($(this).find(".answer-id").val());
                answer.correct = $(this).find(".answer-correct").is(":checked");
                answer.text = $(this).find(".text_answer").val();
                answer.img = $(this).find(".show-image-answer").attr("data-link");
                answers.push(answer);
            });
        }
        question.answer = answers;
        questions.push(question);
    });
    return questions;
}

function SubmitFormCreate(formID, functionCallback, ajxURL = "") {
    let formElement = $("#" + formID);
    formElement.find("input").on("click change blur", function (e) {
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
    formElement.find(".select2").on("select2:open", function (e) {
        let selectElement = $(this);
        if (selectElement.next().find(".select2-selection").hasClass("select2-is-invalid")) {
            selectElement.next().find(".select2-selection").removeClass("select2-is-invalid");
            selectElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    formElement.find(".submit").on("click", function (e) {
        e.preventDefault();
        const buttonElement = $(this);
        const submitButtonText = buttonElement.text();
        const submitButtonTextLoading = "Loading...";
        var postData = formElement.serializeObject();
        postData.id = parseInt(postData.id);
        postData.category = parseInt(postData.category);
        postData.content_type = parseInt(postData.content_type);
        postData.questions = getQuestions();

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
        case false:
            if (response.errors.length === 1 && response.errors[0].id === "") {
                new AlertError(response.errors[0].message);
            } else if (response.errors.length > 0) {
                $.each(response.errors, function (key, value) {
                    let inputElement = $("#" + value.id);
                    if (key === 0) {
                    }
                    inputElement.addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
                    inputElement.next().find(".select2-selection").addClass("select2-is-invalid");
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
        case true:
            NoticeSuccess("Quiz has been updated successfully");
            refreshListQuestion(response.data_object);
            // setTimeout(function () {
            //     window.location.reload();
            // }, 1000);
            break;
        default:
            new AlertError("Undefined");
            break;
    }
}

function htmlQuestion(index) {
    return `
        <div class="box-c question-item" data-question="0">
            <input type="hidden" class="question-id" value="">
            <hr class="mt-3 hr_custom bg-gray-400" style="margin: 0 -20px 0 -50px !important;">
            <div class="question-box">
                <div class="d-flex align-items-center question">
                    <label class="col-form-label form-label form-label-lg w-input-group text-uppercase lb-question"><span
                                class="ms-3">QUESTION 1</span></label>
                    <a class="btn p-1 ms-2 remove-question">
                        <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" fill="currentColor"
                             class="bi bi-x-lg" viewBox="0 0 16 16">
                            <path d="M1.293 1.293a1 1 0 0 1 1.414 0L8 6.586l5.293-5.293a1 1 0 1 1 1.414 1.414L9.414 8l5.293 5.293a1 1 0 0 1-1.414 1.414L8 9.414l-5.293 5.293a1 1 0 0 1-1.414-1.414L6.586 8 1.293 2.707a1 1 0 0 1 0-1.414z"></path>
                        </svg>
                    </a>
                </div>
                
                <div class="pp15">
                    <div class="pt2">Question</div>
                </div>
                <div class="pp7">
                    <div class="pp-9">
                        <input class="pp-10 form-control text_question" type="text">
                        <span class="form-text ms-2 w-auto invalid-feedback"></span>
                    </div>
                </div>

                <div class="pp15">
                    <div class="pt2">Background</div>
                </div>
                <div class="pp7">
                    <div class="box-upload-thumb w-100" >
                        <input type="file" id="upload_background` + index + `" name="upload_background" class="d-none upload_background" title="Choose a background please" aria-describedby="thumb" aria-label="Upload" accept="image/*">
                        <div class="d-flex" id="box-upload-thumb">
                            <label class="pp-11 mb-0" for="upload_background` + index + `" style="max-width: 100px;min-width: 100px;border-right: 0;border-left: 1px solid #aab4c8;">Choose file</label>
                            <input class="pp10 w-100 background" id="background" name="background" type="text" value="" readonly style="border-top-left-radius: 0;border-bottom-left-radius: 0;cursor: not-allowed;">
                        </div>
                        <span class="form-text ms-2 w-auto invalid-feedback"></span>
                    </div>
                </div>
                <div class="pp7">
                    <div class="preview-image-answer d-flex justify-content-center">
                        <img class="show-background" src="" data-link="" style="height: 92px">
                    </div>
                    <span class="ms-2 text_answer" style="height: 92px;max-height: 92px; width: 100%;padding: 10px"></span>
                </div>

                <div class="box-answer box-answer-quiz1">
                    <div class="answer-item" data-answer="0">
                        <input type="hidden" class="answer-id" value="">
                        <div class="pp15">
                            <div class="m-0 d-flex align-items-center justify-content-center" style="height: 15px; width: 15px">
                                <input class="m-0 answer-correct" type="checkbox" checked="checked" value="" id="" style="height: 15px; width: 15px;cursor: pointer;border-color: #ced4da;">
                            </div>
                            <label class="pt2 answer-name ms-2" for="">Answer 0</label>
                        </div>
                        <div class="pp7">
                            <div class="w-100">
                                <input style="border-right: 0;" type="text" class="pp-10 text_answer">
                                <span class="form-text ms-2 w-auto invalid-feedback"></span>
                            </div>
                            <a href="javascript:void(0)" class="dg1 btn-remove-answer" title="Remove Answer">
                                <svg xmlns="http://www.w3.org/2000/svg" style="vertical-align: initial;" width="15"
                                     height="15" fill="currentColor" class="bi bi-x-lg" viewBox="0 0 16 16">
                                    <path d="M1.293 1.293a1 1 0 0 1 1.414 0L8 6.586l5.293-5.293a1 1 0 1 1 1.414 1.414L9.414 8l5.293 5.293a1 1 0 0 1-1.414 1.414L8 9.414l-5.293 5.293a1 1 0 0 1-1.414-1.414L6.586 8 1.293 2.707a1 1 0 0 1 0-1.414z"></path>
                                </svg>
                            </a>
                        </div>
                    </div>
                    <div class="answer-item" data-answer="0">
                        <input type="hidden" class="answer-id" value="">
                        <div class="pp15">
                            <div class="m-0 d-flex align-items-center justify-content-center" style="height: 15px; width: 15px">
                                <input class="m-0 answer-correct" type="checkbox" value="" id="" style="height: 15px; width: 15px;cursor: pointer;border-color: #ced4da;">
                            </div>
                            <label class="pt2 answer-name ms-2" for="">Answer 0</label>
                        </div>
                        <div class="pp7">
                            <div class="w-100">
                                <input style="border-right: 0;" type="text" class="pp-10 text_answer">
                                <span class="form-text ms-2 w-auto invalid-feedback"></span>
                            </div>
                            <a href="javascript:void(0)" class="dg1 btn-remove-answer" title="Remove Answer">
                                <svg xmlns="http://www.w3.org/2000/svg" style="vertical-align: initial;" width="15"
                                     height="15" fill="currentColor" class="bi bi-x-lg" viewBox="0 0 16 16">
                                    <path d="M1.293 1.293a1 1 0 0 1 1.414 0L8 6.586l5.293-5.293a1 1 0 1 1 1.414 1.414L9.414 8l5.293 5.293a1 1 0 0 1-1.414 1.414L8 9.414l-5.293 5.293a1 1 0 0 1-1.414-1.414L6.586 8 1.293 2.707a1 1 0 0 1 0-1.414z"></path>
                                </svg>
                            </a>
                        </div>
                    </div>
                </div>
                <div class="pp7 box-answer-quiz1">
                    <a class="pt98 btn-add-answer-quiz1">Add Answer</a>
                </div>

                <div class="box-answer box-answer-quiz2">
                    <div class="answer-item" data-answer="1">
                        <input type="hidden" class="answer-id" value="">
                        <div class="pp15">
                            <div class="m-0 d-flex align-items-center justify-content-center" style="height: 15px; width: 15px">
                                <input class="answer-correct m-0" type="checkbox" checked="checked" value="" id="" style="height: 15px; width: 15px">
                            </div>
                            <label class="pt2 answer-name ms-2" for="">Answer 1</label>
                        </div>
                        <div class="pp7">
                            <div class="preview-image-answer d-flex justify-content-center">
                                <img class="show-image-answer" src="" data-link="" style="height: 92px">
                                <input type="file" class="d-none link logo" value="">
                            </div>
                            <textarea class="ms-2 text_answer" type="text"
                                      style="height: 92px;max-height: 92px; width: 100%;padding: 10px"></textarea>
                        </div>
                        <div class="pp7">
                            <a class="pt98 btn-add-logo logo-config">ADD</a>
                        </div>
                    </div>
                    <div class="answer-item" data-answer="2">
                        <input type="hidden" class="answer-id" value="">
                        <div class="pp15">
                            <div class="m-0 d-flex align-items-center justify-content-center" style="height: 15px; width: 15px">
                                <input class="answer-correct m-0" type="checkbox" value="" id="" style="height: 15px; width: 15px">
                            </div>
                            <label class="pt2 answer-name ms-2" for="">Answer 2</label>
                        </div>
                        <div class="pp7" data-answer="2">
                            <div class="preview-image-answer d-flex justify-content-center">
                                <img class="show-image-answer" src="" data-link="" style="height: 92px">
                                <input type="file" class="d-none link logo" value="">
                            </div>
                            <textarea class="ms-2 text_answer" type="text"
                                      style="height: 92px;max-height: 92px; width: 100%;padding: 10px"></textarea>
                        </div>
                        <div class="pp7">
                            <a class="pt98 btn-add-logo logo-config">ADD</a>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `;
}

function htmlAnswer(numberQuestion, numberAnswer) {
    return `
                    <div class="answer-item" data-answer="0">
                        <input type="hidden" class="answer-id" value="">
                        <div class="pp15">
                            <div class="m-0 d-flex align-items-center justify-content-center" style="height: 15px; width: 15px">
                                <input class="m-0 answer-correct" type="checkbox" value="" id="" style="height: 15px; width: 15px;cursor: pointer;border-color: #ced4da;">
                            </div>
                            <label class="pt2 answer-name ms-2" for="">Answer 0</label>
                        </div>
                        <div class="pp7">
                            <div class="w-100">
                                <input style="border-right: 0;" type="text" class="pp-10 text_answer">
                                <span class="form-text ms-2 w-auto invalid-feedback"></span>
                            </div>
                            <a href="javascript:void(0)" class="dg1 btn-remove-answer" title="Remove Answer">
                                <svg xmlns="http://www.w3.org/2000/svg" style="vertical-align: initial;" width="15"
                                     height="15" fill="currentColor" class="bi bi-x-lg" viewBox="0 0 16 16">
                                    <path d="M1.293 1.293a1 1 0 0 1 1.414 0L8 6.586l5.293-5.293a1 1 0 1 1 1.414 1.414L9.414 8l5.293 5.293a1 1 0 0 1-1.414 1.414L8 9.414l-5.293 5.293a1 1 0 0 1-1.414-1.414L6.586 8 1.293 2.707a1 1 0 0 1 0-1.414z"></path>
                                </svg>
                            </a>
                        </div>
                    </div>
    `;
}

function clickAddImage() {
    $(".list-questions").on("click", ".btn-add-logo", function () {
        $(this).closest(".answer-item").find(".logo").click();
    });
    $(".list-questions").on("change", ".logo", function (e) {
        validateAndUploadLogo($(this), _validFileImageExtensions, e);
    });
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
            $(element).prev(".show-image-answer").attr("src", domainUpload + response.data_object.thumb[0]).attr("data-link", response.data_object.thumb[0]);
            $(element).attr("value", response.data_object.thumb[0]);
            if (element.hasClass("upload_background")) {
                element.closest(".question-item").find(".background").val(response.data_object.thumb[0]);
                element.closest(".question-item").find(".show-background").attr("src", domainUpload + response.data_object.thumb[0]);
            }
            if (element.hasClass("upload_illustration")) {
                element.closest(".box-illustration").find(".illustration").val(response.data_object.thumb[0]);
                element.closest(".box-illustration").find(".show-illustration").attr("src", domainUpload + response.data_object.thumb[0]);
            }
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

$("#nav-question").on("change", ".upload_background", function (e) {
    // $("#upload_background").removeClass("d-none");
    validateAndUploadLogo($(this), _validFileImageExtensions, e);
});