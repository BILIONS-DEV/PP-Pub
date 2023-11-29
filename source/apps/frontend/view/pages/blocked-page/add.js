let _validCSV = [".csv"];
let api_upload_csv = "/blocked-page/importCSV";
let maxSize = 15 * 1024 * 1024;
$(document).ready(function () {
    importCSV();
    handleBlockedPage();

    submitFormCreate("formBlockedPage", added);
});

function handleBlockedPage() {
    $(".box-page-selected").on("click", ".target-item .btn-remove", function (e) {
        $(this).closest(".target-item").remove();
        checkButtonClear();
    });
    $(".remove_all").on("click", function () {
        $(".box-page-selected").html("");
        checkButtonClear();
    });
    $("#input-page").on('keypress', function (e) {
        var keyCode = e.keyCode || e.which;
        if (keyCode === 13) {
            e.preventDefault();
            validatePages();
            checkButtonClear();
        }
    });
}

function checkButtonClear() {
    if ($('.item_selected').length > 0){
        $(".btn-clear-target").removeClass("d-none");
    } else {
        $(".btn-clear-target").addClass("d-none");
    }
}

function importCSV() {
    // Upload file csv
    $("#upload_csv").on("change", function (e) {
        $("#loading_csv").removeClass("d-none");
        validate($(this), _validCSV, e);
    });
}

function validate(element, validFile, event) {
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
                $("#loading_csv").addClass("d-none");
                new AlertError("Sorry, " + sFileName + " is invalid, allowed extensions are: " + validFile.join(", "));
                return false;
            }

        } else {
            $("#loading_csv").addClass("d-none");
        }
    }
    checkButtonClear();
    return true;
}

function uploadFile(event, element) {
    var fd = new FormData();
    var file = event.target.files[0];
    if (file.size > maxSize) {
        new AlertError("You uploaded file over 15mb, please choose another file!");
        $("#loading_csv").addClass("d-none");
        return;
    }
    fd.append('file', file);

    $.ajax({
        url: api_upload_csv,
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
            $("#loading_csv").addClass("d-none");
            // buttonElement.attr('disabled', false).text(submitButtonText);
        },
        success: function (responseJSON) {
            $("#loading_csv").addClass("d-none");
            // buttonElement.attr('disabled', false).text(submitButtonText);
        },
        complete: function (res) {
            $("#loading_csv").addClass("d-none");
            afterUpload(res.responseJSON, element);
        }
    });
    element.nextAll(".form-control").val(file.name);
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
            let pages = [];
            $.each(response.data_object, function (i, item) {
                pages.push(item.issue_location);
            });
            validatePages(pages);
            NoticeSuccess("Import .csv success!");
    }
}

function submitFormCreate(formID, functionCallback, ajaxUrl = "") {
    let formElement = $("#" + formID);
    formElement.find("input").on("input", function (e) {
        let inputElement = $(this);
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
        if (inputElement.closest("#box-upload-csv").hasClass("is-invalid")) {
            inputElement.closest("#box-upload-csv").removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    formElement.find("textarea").on("input", function (e) {
        let inputElement = $(this);
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
        if (inputElement.hasClass("textarea-add-custom")) {
            $(".notify-domain-err").addClass("invalid-feedback");
            $(".domain-err").empty();
        } else {
            $(".notify-creative-err").addClass("invalid-feedback");
            $(".creative-err").empty();
        }
    });

    formElement.on('keypress', "input", function (e) {
        let noSubmit = false;
        noSubmit = $(this).data("no-submit");
        var keyCode = e.keyCode || e.which;
        if (keyCode === 13 && !noSubmit) {
            e.preventDefault();
            submit();
        }
    });
    formElement.find(".submit").on("click", function (e) {
        e.preventDefault();
        submit();
    });

    function submit() {
        let errV = validatePages();
        if (errV === false) {
            AlertError("Please check the page list again!");
            $("#input-page").focus();
            return;
        }
        const buttonElement = $(".submit");
        const submitButtonText = buttonElement.text();
        const submitButtonTextLoading = "Loading...";
        var postData = formElement.serializeObject();
        let pages = [];
        if ($('.box-page-selected').find('.item_selected').length) {
            $('.box-page-selected').find('.item_selected').each(function () {
                pages.push($(this).find('.btn-remove').attr('data-name').trim());
            });
        }
        postData.pages = pages;
        postData.id = parseInt(postData.id);

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
                functionCallback(res.responseJSON, formElement);
            }
        });
    }
}

function added(response, formElement) {
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
            NoticeSuccess("Blocked has been created successfully");
            setTimeout(function () {
                window.location.replace(`/rule`);
            }, 1000);
            break;
        default:
            new AlertError("Undefined");
            break;
    }
}

function displaySelected(name, element) {
    name = name.replace(/(^\w+:|^)\/\//, '');
    element.append(`
    <div class="d-flex flex-row align-items-center p-2 border-bottom target-item item_selected">
        <div class="col p-0" style="width: 97%">
            <span class="m-0" style="overflow-wrap: break-word;font-size: 14px">${name}</span>
        </div>
        <div class="col-auto px-0">
            <button type="button" data-name="${name}" style="width: 16px; height: 16px"
                    class="btn d-flex align-items-center btn-outline-danger btn-icon rounded-circle p-0 btn-remove">
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor" class="bi bi-dash-lg"
                     viewBox="0 0 16 16">
                    <path d="M0 8a1 1 0 0 1 1-1h14a1 1 0 1 1 0 2H1a1 1 0 0 1-1-1z"></path>
                </svg>
            </button>
        </div>
    </div>                                           
    `);
}

function validatePages(pages) {
    let array = [];
    let sp = ",";
    if (pages) {
        array = pages;
    } else {
        let domainName = $("#input-page").val();
        if (!domainName.trim()) {
            return true;
        }

        if (domainName.trim().indexOf('\n') !== -1) {
            array = domainName.trim().split("\n");
            sp = "\n";
        } else if (domainName.trim().indexOf(',') !== -1) {
            array = domainName.trim().split(",");
            sp = ",";
        } else {
            array.push(domainName.trim());
        }
    }
    let domainErr = "";
    let checkErr = true;
    let countErr = 0;
    let txtDomain = "";
    let listError = [];
    if (array.length > 0) {
        array.forEach(function (name) {
            name = name.replace(/(^\w+:|^)\/\//, '');
            if (name !== "") {
                let flag = true;
                if ($('.box-page-selected').find('.item_selected').length) {
                    $('.box-page-selected').find('.item_selected').each(function () {
                        // console.log($(this).find('.domain-remove').attr('data-name').trim());
                        if (name.trim() === $(this).find('.btn-remove').attr('data-name').trim()) {
                            flag = false;
                        }
                    });
                }
                if (!validURL(name)) {
                    flag = false;
                    listError.push(name);
                }
                if (flag) {
                    displaySelected(name, $('.box-page-selected'));
                }
            }
        });
    }
    if (listError.length > 0) {
        checkErr = false;
        listError.forEach(function (name) {
            name = name.replace(/\s+/g, '');
            if (name !== "") {
                if (countErr === 0) {
                    domainErr = name;
                    txtDomain = name;
                } else {
                    txtDomain += sp + name;
                }
                countErr++;
            }
            NoticeError(name + " fail");
        });
    }
    $("#input-page").val(txtDomain);
    return checkErr;
}

function validURL(str) {
    var pattern = new RegExp(
        '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|' + // domain name
        '((\\d{1,3}\\.){3}\\d{1,3}))' + // OR ip (v4) address
        '(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*' + // port and path
        '(\\?[;&a-z\\d%_.~+=-]*)?' + // query string
        '(\\#[-a-z\\d_]*)?$', 'i'); // fragment locator
    return !!pattern.test(str);
}