$(document).ready(function () {
    SubmitFormCreate("editChannels", Added, "/channels/edit");
    // $.fn.select2.defaults.set("theme", "bootstrap");
    // $("#keyword").select2({
    //     tags: true,
    //     "language": {
    //         "noResults": function () {
    //             return "Please add keyword and enter";
    //         }
    //     },
    //     escapeMarkup: function (markup) {
    //         return markup;
    //     }
    // });
    handleTags()
    // LoadHistory();
});


function SubmitFormCreate(formID, functionCallback, ajxURL = "") {
    let formElement = $("#" + formID);
    formElement.find("input").on("click change blur", function (e) {
        let inputElement = $(this)
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    formElement.find("textarea").on("click change blur", function (e) {
        let inputElement = $(this)
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    formElement.find(".select2").on("select2:open", function (e) {
        let selectElement = $(this)
        if (selectElement.next().find(".select2-selection").hasClass("select2-is-invalid")) {
            selectElement.next().find(".select2-selection").removeClass("select2-is-invalid")
            selectElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    formElement.find(".submit").on("click", function (e) {
        e.preventDefault();
        const buttonElement = $(this);
        const submitButtonText = buttonElement.text();
        const submitButtonTextLoading = "Loading...";
        var postData = formElement.serializeObject();
        postData.id = parseInt(postData.id)
        postData.category = parseInt(postData.category)
        postData.language = parseInt(postData.language)
        postData.keyword = GetTags();
        // if (typeof postData.keyword === "string") {
        //     postData.keyword = [postData.keyword]
        // }
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

function Added(response) {
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
                    inputElement.find(".bs-placeholder").attr("style", "border-color:#e35d6a!important");
                    inputElement.addClass("is-invalid").closest(".bootstrap-select").addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
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
            NoticeSuccess("Channels has been updated successfully")
            break
        default:
            new AlertError("Undefined");
            break
    }
}

function GetTags() {
    var tags = [];
    $("#list-tags").find(".tag-item").each(function () {
        tags.push($(this).find('span').text());
    });
    return tags;
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
