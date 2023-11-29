function CopyTextToClipboard(text) {
    const textArea = document.createElement("textarea");

    //
    // *** This styling is an extra step which is likely not required. ***
    //
    // Why is it here? To ensure:
    // 1. the element is able to have focus and selection.
    // 2. if the element was to flash render it has minimal visual impact.
    // 3. less flakyness with selection and copying which **might** occur if
    //    the textarea element is not visible.
    //
    // The likelihood is the element won't even render, not even a
    // flash, so some of these are just precautions. However in
    // Internet Explorer the element is visible whilst the popup
    // box asking the user for permission for the web page to
    // copy to the clipboard.
    //

    // Place in the top-left corner of screen regardless of scroll position.
    textArea.style.position = 'fixed';
    textArea.style.top = 0;
    textArea.style.left = 0;

    // Ensure it has a small width and height. Setting to 1px / 1em
    // doesn't work as this gives a negative w/h on some browsers.
    textArea.style.width = '2em';
    textArea.style.height = '2em';

    // We don't need padding, reducing the size if it does flash render.
    textArea.style.padding = 0;

    // Clean up any borders.
    textArea.style.border = 'none';
    textArea.style.outline = 'none';
    textArea.style.boxShadow = 'none';

    // Avoid flash of the white box if rendered for any reason.
    textArea.style.background = 'transparent';


    textArea.value = text;

    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();
    let successful;
    try {
        successful = document.execCommand('copy');
        let msg = successful ? 'successful' : 'unsuccessful';
        // console.log('Copying text command was ' + msg);
    } catch (err) {
        console.log('Oops, unable to copy');
    }

    document.body.removeChild(textArea);
    return successful;
}

/**
 * Copy text from element
 * @param elem
 * @returns {boolean}
 */
function CopyToClipboardWithElement(elem) {
    // create hidden text element, if it doesn't already exist
    let targetId = "_hiddenCopyText_";
    let isInput = elem.tagName === "INPUT" || elem.tagName === "TEXTAREA";
    let origSelectionStart, origSelectionEnd;
    if (isInput) {
        // can just use the original source element for the selection and copy
        target = elem;
        origSelectionStart = elem.selectionStart;
        origSelectionEnd = elem.selectionEnd;
    } else {
        // must use a temporary form element for the selection and copy
        target = document.getElementById(targetId);
        if (!target) {
            var target = document.createElement("textarea");
            target.style.position = "absolute";
            target.style.left = "-9999px";
            target.style.top = "0";
            target.id = targetId;
            document.body.appendChild(target);
        }
        target.textContent = elem.textContent;
    }
    // select the content
    let currentFocus = document.activeElement;
    target.focus();
    target.setSelectionRange(0, target.value.length);

    // copy the selection
    let succeed;
    try {
        succeed = document.execCommand("copy");
    } catch (e) {
        succeed = false;
    }
    // restore original focus
    if (currentFocus && typeof currentFocus.focus === "function") {
        currentFocus.focus();
    }

    if (isInput) {
        // restore prior selection
        elem.setSelectionRange(origSelectionStart, origSelectionEnd);
    } else {
        // clear temporary content
        target.textContent = "";
    }
    return succeed;
}

function submitForm(formID, url = "") {
    var formElement = $("#" + formID);
    formElement.find(".submit").on("click", function (e) {
        e.preventDefault();
        const buttonElement = $(this);
        const submitButtonText = buttonElement.text();
        const submitButtonTextLoading = "Loading...";
        var postData = formElement.serializeArray();
        $.ajax({
            url: url,
            type: "POST",
            dataType: "JSON",
            data: postData,
            error: function (jqXHR, exception) {
                let msg = "";
                switch (jqXHR.status) {
                    case 0:
                        msg = 'Not connect.\n Verify Network.';
                        break;
                    case 404:
                        msg = 'Requested page not found. [404]';
                        break;
                    case 500:
                        msg = 'Internal Server Error [500].';
                        break;
                    default:
                        switch (exception) {
                            case "parsererror":
                                msg = 'Requested JSON parse failed.';
                                break;
                            case "timeout":
                                msg = 'Time out error.';
                                break;
                            case "abort":
                                msg = 'Ajax request aborted.';
                                break;
                            default:
                                msg = 'Uncaught Error.\n' + jqXHR.responseText;
                        }

                }
                alertError(msg);
                buttonElement.attr('disabled', false).text(submitButtonText);
            },
            beforeSend: function (xhr) {
                xhr.overrideMimeType("text/plain; charset=x-user-defined");
                buttonElement.attr('disabled', true).text(submitButtonTextLoading);
                removeWarning(formElement);
            }
        }).done(function (result) {
            buttonElement.attr('disabled', false).text(submitButtonText);
            if (!result) {
                alertError("ERROR");
                return;
            }
            if (!result.status) {
                alertError("ERROR: not have status");
                return;
            }
            switch (result.status) {
                case "success":
                    if (result.message) {
                        alertSuccess(result.message);
                    }
                    break;

                case "error":
                default:
                    if (result.message) {
                        alertError(result.message);
                    }
                    if ($.isArray(result.errors)) {
                        $.each(result.errors, function (key, err) {
                            $("#" + err.field).closest('.form-group').addClass('has-error').find('.error').removeClass('d-none').text(err.message);
                        });
                    }
                    break;
            }
        });
    });
}

function removeWarning(formElement) {
    formElement.find(".form-group").removeClass("has-error").find(".error").addClass("d-none").text("error");
}

/**
 *
 * @param object formId
 * @param string id
 * @param string value
 * @returns {undefined}
 */
function addInputToFormSubmit(formId, id, value) {
    var input = $(formId).find('#' + id);
    if (input.length) {
        input.val(value);
    } else {
        $('<input>').attr({type: 'hidden', id: id, name: id, value: value}).appendTo(formId);
    }
}

/**
 * Show modal warning error
 * @param string message
 * @returns HTML
 */
function alertError(message) {
    swal("Ohhh, Sorry !!!", message, {
        icon: "error",
        // html: false,
        buttons: {
            confirm: {
                className: 'btn btn-danger'
            }
        }
    });
}

/**
 * Show modal success
 * @param string message
 * @returns HTML
 */
function alertSuccess(message) {
    swal("Congratulations !!!", message, {
        icon: "success",
        buttons: {
            confirm: {
                className: 'btn btn-success'
            }
        }
    });
}

/**
 * Show notification success
 * @param string message
 * @returns HTML
 */
function notifiSuccess(message) {
    var content = {};
    var state = 'success';
    content.title = 'Good Job';
    content.message = message;
    var placementFrom = 'top';
    var placementAlign = 'right';
    content.icon = 'fa fa-check';
    var notify = $.notify(content, {
        type: state,
        placement: {
            from: placementFrom,
            align: placementAlign
        },
        time: 1000,
        delay: 0,
        z_index: 2000
    });
    setTimeout(function () {
        notify.close();
    }, 3000);
}

function ajaxErrorMessage(jqXHR, exception) {
    let msg = '';
    if (jqXHR.status === 0) {
        msg = 'Not connect.\n Verify Network.';
    } else if (jqXHR.status == 404) {
        msg = 'Requested page not found. [404]';
    } else if (jqXHR.status == 500) {
        msg = 'Internal Server Error [500].';
    } else if (exception === 'parsererror') {
        msg = 'Requested JSON parse failed.';
    } else if (exception === 'timeout') {
        msg = 'Time out error.';
    } else if (exception === 'abort') {
        msg = 'Ajax request aborted.';
    } else {
        msg = 'Uncaught Error.\n' + jqXHR.responseText;
    }
    return msg;
}

// Converts object with "name" and "value" keys
// into object with "name" key having "value" as value
// See http://stackoverflow.com/a/12399106/3549014 for more details
$.fn.serializeObject = function () {
    let obj = {};
    $.each(this.serializeArray(), function (i, o) {
        let n = o.name, v = !isNumeric(o.value) ? o.value : Number(o.value);
        obj[n] = obj[n] === undefined ? v
            : $.isArray(obj[n]) ? obj[n].concat(v)
                : [obj[n], v];
    });
    return obj;
};

function isNumeric(str) {
    if (typeof str != "string") return false // we only process strings!
    return !isNaN(str) && // use type coercion to parse the _entirety_ of the string (`parseFloat` alone does not do this)...
        !isNaN(parseFloat(str)) // ...and ensure strings of whitespace fail
}
