$(document).ready(function () {

    $(".blocking-tab .nav-tabs-blocking").on("click", ".nav-link", function () {
        $(".blocking-tab .nav-tabs-blocking").find(".nav-link").removeClass("pp-4");
        $(this).addClass("pp-4");
        var tab = $(this).attr("data-tab");
        if (tab != "1") {
            $(this).addClass("at-1");
        } else {
            $(".blocking-tab .nav-tabs-blocking").find('.at-1').removeClass('at-1');
        }
    });

    $("#submitBlocking").find("textarea").on("input", function (e) {
        let inputElement = $(this);
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });

    $(".box-list-domain").on("click", ".btn-remove", function (e) {
        // console.log($(this).closest(".target-item"))
        $(this).closest(".target-item").remove();
    });

    $(".box-list-creative").on("click", ".btn-remove", function (e) {
        // console.log($(this).closest(".target-item"))

        $(this).closest(".target-item").remove();
    });

    $(".textarea-add-custom").on("keypress", function (e) {
        if(e.which === 13){
            ValidateDomainRestrictions();
        }
    });

    $(".textarea-add-creative").on("keypress", function (e) {
        if(e.which === 13){
            ValidateCreativeIds();
        }
    });

    SubmitFormCreate("submitBlocking", AddBlocking);
    // LoadHistory();
});

function validateDomain(listDomain) {
    let url = "/blocking/validateDomain";
    let result;
    $.ajax({
        url: url,
        type: "POST",
        dataType: "JSON",
        async: false,
        contentType: "application/json",
        data: JSON.stringify({
            listDomain: listDomain
        }),
        beforeSend: function (xhr) {
            // xhr.overrideMimeType("text/plain; charset=x-user-defined");
        },
        error: function (jqXHR, exception) {
            var msg = AjaxErrorMessage(jqXHR, exception);
            new AlertError("AJAX ERROR: " + msg);
        }
    }).done(function (json) {
        result = json;
    });
    return result;
}

function DisplaySelected(name, element) {
    element.append(`<div class="d-flex flex-row align-items-center p-2 border-bottom target-item item_selected">
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
    </div>`);
}

function ValidateDomainRestrictions() {
    let array;
    let domainName = $(".textarea-add-custom").val();
    if (!domainName.trim()) {
        return true;
    }

    let sp = "";
    if (domainName.trim().indexOf('\n') !== -1) {
        array = domainName.trim().split("\n");
        sp = "\n";
    } else {
        array = domainName.trim().split(",");
        sp = ",";
    }
    let domainErr = "";
    let checkErr = true;
    let countErr = 0;
    let txtDomain = "";
    let result = validateDomain(array);
    if (result.listValid) {
        result.listValid.forEach(function (name) {
            name = name.replace(/\s+/g, '');
            if (name !== "") {
                let flag = true;
                if ($('.box-list-domain').find('.item_selected').length) {
                    $('.box-list-domain').find('.item_selected').each(function () {
                        // console.log($(this).find('.domain-remove').attr('data-name').trim());
                        if (name.trim() === $(this).find('.btn-remove').attr('data-name').trim()) {
                            flag = false;
                        }
                    });
                }
                if (flag) {
                    DisplaySelected(name, $('.box-list-domain'));
                }
            }
        });
    }
    if (result.listError) {
        checkErr = false;
        result.listError.forEach(function (name) {
            name = name.replace(/\s+/g, '');
            if (name !== "") {
                if (countErr === 0) {
                    domainErr = name;
                    txtDomain = name;
                } else if (countErr < 5) {
                    domainErr += ", " + name;
                    txtDomain += sp + name;
                } else {
                    txtDomain += sp + name;
                }
                countErr++;
                checkErr = false;
            }
        });
    }
    if (!checkErr) {
        if (countErr > 5) {
            domainErr = domainErr + ",... ";
        }
        $(".textarea-add-custom").addClass("is-invalid").nextAll("span.invalid-feedback").text(domainErr + " aren't valid, please check again!");
    }
    // console.log(array)
    $(".textarea-add-custom").val(txtDomain);
    return checkErr;
}

function ValidateCreativeIds() {
    let array;
    let creativeId = $(".textarea-add-creative").val();
    if (!creativeId.trim()) {
        return true;
    }

    let sp = "";
    if (creativeId.trim().indexOf('\n') !== -1) {
        array = creativeId.trim().split("\n");
        sp = "\n";
    } else {
        array = creativeId.trim().split(",");
        sp = ",";
    }
    let creativeErr = "";
    let checkErr = true;
    let countErr = 0;
    let txtcreative = "";
    let listErr = [];
    array.forEach(function (value) {
        value = value.replace(/\s+/g, '');
        if (value !== "") {
            let flag = true;
            if ($('.box-list-creative').find('.item_selected').length) {
                $('.box-list-creative').find('.item_selected').each(function () {
                    // console.log($(this).find('.domain-remove').attr('data-name').trim());
                    if (value.trim() === $(this).find('.btn-remove').attr('data-name').trim()) {
                        flag = false;
                    }
                });
            }
            if (flag) {
                DisplaySelected(value.trim(), $('.box-list-creative'));
            }
        }
    });
    if (listErr.length > 0) {
        checkErr = false;
        listErr.forEach(function (name) {
            name = name.replace(/\s+/g, '');
            if (name !== "") {
                if (countErr === 0) {
                    creativeErr = name;
                    txtcreative = name;
                } else if (countErr < 5) {
                    creativeErr += ", " + name;
                    txtcreative += sp + name;
                } else {
                    txtcreative += sp + name;
                }
                countErr++;
                checkErr = false;
            }
        });
    }
    if (!checkErr) {
        if (countErr > 5) {
            creativeErr = creativeErr + ",... ";
        }
        $(".textarea-add-creative").addClass("is-invalid").nextAll("span.invalid-feedback").text(creativeErr);
    }
    // console.log(array)
    $(".textarea-add-creative").val(txtcreative);
    return checkErr;
}

function SubmitFormCreate(formID, functionCallback, ajaxUrl = "") {
    let formElement = $("#" + formID);
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
        if (inputElement.hasClass("textarea-add-custom")){
            $(".notify-domain-err").addClass("invalid-feedback");
            $(".domain-err").empty();
        } else {
            $(".notify-creative-err").addClass("invalid-feedback");
            $(".creative-err").empty();
        }
    });
    formElement.find(".submit").on("click", function (e) {
        e.preventDefault();
        let errV = ValidateDomainRestrictions();
        if (errV === false) {
            AlertError("Please check the restriction list again!");
            $(".textarea-add-custom").focus();
            return;
        }
        errV = ValidateCreativeIds();
        if (errV === false) {
            AlertError("Please check the creative list again!");
            $(".textarea-add-creative").focus();
            return;
        }
        const buttonElement = $(this);
        const submitButtonText = buttonElement.text();
        const submitButtonTextLoading = "Loading...";
        let postData = formElement.serializeObject();
        let inventories = [];
        if (postData.inventories && postData.inventories.constructor === Array) {
            $.each(postData.inventories, function (index, value) {
                inventories.push(parseInt(value));
            });
        } else {
            inventories.push(parseInt(postData.inventories));
        }
        postData.inventories = inventories;
        let advertiseDomains = [];
        if ($('.box-list-domain').find('.item_selected').length) {
            $('.box-list-domain').find('.item_selected').each(function () {
                // console.log($(this).find('.domain-remove').attr('data-name').trim());
                advertiseDomains.push($(this).find('.btn-remove').attr('data-name').trim());
            });
        }
        postData.advertise_domains = advertiseDomains;
        let creativeIds = [];
        if ($('.box-list-creative').find('.item_selected').length) {
            $('.box-list-creative').find('.item_selected').each(function () {
                // console.log($(this).find('.domain-remove').attr('data-name').trim());
                creativeIds.push($(this).find('.btn-remove').attr('data-name').trim());
            });
        }
        postData.creative_ids = creativeIds;
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
    });
}

function AddBlocking(response, formElement) {
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
            NoticeSuccess("Blocking has been updated successfully");
            break;
        default:
            new AlertError("Undefined");
            break;
    }
}