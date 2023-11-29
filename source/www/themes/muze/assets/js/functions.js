/**
 * Redirect to url
 * @param url
 * @constructor
 */
function Redirect(url) {
    // similar behavior as an HTTP redirect
    window.location.replace(url);
    // similar behavior as clicking on a link
    window.location.href = url;
}

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

/**
 * Submit form
 * @param formID
 * @param ajxURL
 * @param callback
 * @constructor
 */
function SubmitForm(formID, functionCallback, ajxURL = "") {
    let formElement = $("#" + formID);
    let button = formElement.find(".submit");
    let submitButtonText = button.text();
    let submitButtonTextLoading = "Loading...";
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
    formElement.on('keypress', function (e) {
        if (e.which === 13 && !button.is(":disabled") && formElement.find("input").is(':focus')) {
            button.click();
        }
    });
    button.on("click", function (e) {
        e.preventDefault();
        $.each(formElement.find("input.is-invalid"), function () {
            $(this).removeClass("is-invalid").next(".invalid-feedback").empty();
        });
        const postData = formElement.serializeArray();
        const buttonElement = $(this);
        $.ajax({
            url: ajxURL,
            type: "POST",
            dataType: "JSON",
            data: postData,
            beforeSend: function (xhr) {
                buttonElement.attr('disabled', true).text(submitButtonTextLoading);
            },
            error: function (jqXHR, exception) {
                const msg = AjaxErrorMessage(jqXHR, exception);
                NoticeError("AJAX ERROR: " + msg);
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

/**
 * Check if obj is function
 * @param functionToCheck
 * @returns {boolean}
 * @constructor
 */
function IsFunction(functionToCheck) {
    return functionToCheck && {}.toString.call(functionToCheck) === '[object Function]';
}

/**
 * Show modal warning error
 * @param message
 * @param functionCallback
 * @constructor
 */
function AlertError(message, functionCallback) {
    // let alertModal = swal("Ohhh, Sorry !!!", message, {
    //     icon: "error",
    //     // html: false,
    //     buttons: {
    //         confirm: {
    //             className: 'btn btn-lg btn-danger'
    //         }
    //     }
    // });
    // if (functionCallback !== undefined && IsFunction(functionCallback)) {
    //     alertModal.then(() => {
    //         functionCallback();
    //     });
    // }

    const content = {};
    const state = 'danger';
    // content.title = 'Oh No!!! ';
    content.message = message;
    const placementFrom = 'top';
    const placementAlign = 'right';
    content.icon = 'fa fa-check';
    const notify = $.notify(content, {
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

/**
 * Show modal success
 * @param message
 * @param functionCallback
 * @constructor
 */
function AlertSuccess(message, functionCallback) {
    let alertModal = swal("Congratulations !!!", message, {
        icon: "success",
        buttons: {
            confirm: {
                className: 'btn btn-lg btn-success'
            }
        }
    });
    if (functionCallback !== undefined && IsFunction(functionCallback)) {
        alertModal.then(() => {
            functionCallback();
        });
    }
}

/**
 * Show notification success
 * @param string message
 * @returns HTML
 */
function NoticeSuccess(message) {
    const content = {};
    const state = 'success';
    content.title = 'Good Job!!! ';
    content.message = message;
    content.icon = 'fa fa-check';
    const notify = $.notify(content, {
        type: state,
        placement: {
            from: "bottom",
            align: "center"
        },
        time: 1000,
        delay: 0,
        z_index: 2000,
        animate: {
            enter: 'animated fadeInDown',
            exit: 'animated fadeOutRight'
        },
        onClose: true,
        allow_dismiss: true,
        icon_type: 'class',
    });
    setTimeout(function () {
        notify.close();
    }, 3000);
}

/**
 * Show notification Error
 * @param message
 * @returns HTML
 */
function NoticeError(message) {
    const content = {};
    const state = 'danger';
    content.title = 'Error!!! ';
    content.message = message;
    content.icon = 'fa fa-check';
    const notify = $.notify(content, {
        type: state,
        placement: {
            from: "top",
            align: "right"
        },
        time: 1000,
        delay: 0,
        z_index: 2000,
        onClose: null,
        animate: {
            enter: 'animated fadeInDown',
            exit: 'animated fadeOutRight'
        },
    });
    setTimeout(function () {
        notify.close();
    }, 4000);
}

/**
 *
 * @param jqXHR
 * @param exception
 * @returns {string}
 */
function AjaxErrorMessage(jqXHR, exception) {
    let msg = '';
    if (jqXHR.status === 0) {
        msg = 'Not connect.\n Verify Network.';
    } else if (jqXHR.status === 404) {
        msg = 'Requested page not found. [404]';
    } else if (jqXHR.status === 500) {
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

function notifiResult(result) {
    if (result.success == true || result.status == true || result == true || result.success || result.status) {
        var state = 'success';
        if (typeof result.message !== "undefined" && result.message) {
            state = result.message
        }
        NoticeSuccess(state);
    } else {
        if (result.warning) {
            notifiWarning(result.warning);
        } else if (result.error) {
            notifiError(result.error);
        } else if (result.status == false) {
            notifiError(result.message);
        } else {
            notifiError(result);
        }
    }
}

function show_notifi(result) {
    var content = {};

    if (typeof result.success !== "undefined" || result === true) {
        var state = 'success';
        content.title = 'Good Job';
        if (result.success) {
            content.message = result.success;
        } else {
            content.message = 'TRUE';
        }
    } else if (typeof result.status !== "undefined") {
        var state = 'success';
        if (result.message) {
            state = result.message
        }
        content.title = 'Good Job';
        if (result.success) {
            content.message = result.success;
        } else {
            content.message = 'TRUE';
        }
    } else {
        var state = 'danger';
        content.title = 'Ohh No';
        content.message = result.error;
    }
    var placementFrom = 'top';
    var placementAlign = 'right';
    content.icon = 'fa fa-bell';
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

function notifiError(error) {
    swal("Fail!", error, {
        icon: "error",
        html: true,
        buttons: {
            confirm: {
                className: 'btn btn-danger'
            }
        }
    });
}

function notifiWarning(warning) {
    swal("Warning!", warning, {
        icon: "warning",
        html: true,
        buttons: {
            confirm: {
                className: 'btn btn-warning'
            }
        }
    });
}

// Converts object with "name" and "value" keys
// into object with "name" key having "value" as value
// See http://stackoverflow.com/a/12399106/3549014 for more details
$.fn.serializeObject = function () {
    let obj = {};
    $.each(this.serializeArray(), function (i, o) {
        let n = o.name, v = o.value;
        obj[n] = obj[n] === undefined ? v
            : $.isArray(obj[n]) ? obj[n].concat(v)
                : [obj[n], v];
    });
    return obj;
};

// Set select2 default
$(document).ready(function () {
    if ($.fn.select2 !== undefined) {
        $.fn.select2.defaults.set("theme", "bootstrap");
        if ($("select.select2").length > 0) {
            $(".select2").select2({
                // width: "auto",
            });
            $(".select2").on('select2:open', (e) => {
                if (!e.target.multiple) {
                    $('.select2-container--open .select2-search--dropdown .select2-search__field').last()[0].focus()
                }
            });
            $(".select2-multiple").select2({
                // width: "auto",
                closeOnSelect: false,
            });
        }
    }
});

function _donetyping($input, doneTypingInterval, callback) {
    var typingTimer;

    $input.on('keyup', function () {
        clearTimeout(typingTimer);
        typingTimer = setTimeout(callback, doneTypingInterval);
    });

    $input.on('keydown', function () {
        clearTimeout(typingTimer);
    });

}

// load history
function LoadHistory() {
    // $(".history-a").removeClass("d-none")
    $("body").on("click", ".load-history", function () {
        var id = ($(this).attr("data-id"))
        var object = $(this).attr("data-object")
        // var object = "line_item"
        $.ajax({
            type: 'POST',
            url: '/history/load-histories',
            data: { id: id, object: object }
        })
            .done(function (result) {
                if (result.error) {
                    alert(result.error);
                } else {
                    $("#modal-history").find(".modal-content").html(result)
                    checkTypeHistory($("#modal-history"), 'off')
                    $('[data-bs-toggle="tooltip"]').tooltip()
                }
            })
    })
    // checkTypeHistory('off')
}

LoadHistory()

$("#modal-history").on("change", ".loadType", function () {
    checkTypeHistory($("#modal-history"), 'on');
});

function checkTypeHistory(motherBox, delay, loadHistory = false) {
    // var motherBox = $("#modal-history");
    var loadType = motherBox.find("input[name='loadType']:checked").val();
    if (loadHistory == true) {
        loadType = "ChangesOnly"
    }
    if (loadType == 'ChangesOnly') {
        if (delay == 'on') {
            motherBox.find(".table-history tbody").css("opacity", 0);
            setTimeout(function () {
                motherBox.find(".everything").addClass("d-none");
                motherBox.find(".value-item").each(function () {
                    var id = $(this).attr('data-id');
                    if (!id) {
                        return;
                    }

                    var showItem = false;
                    motherBox.find(".value-" + id).each(function () {
                        if (!$(this).hasClass('d-none')) {
                            showItem = true;
                        }
                    });
                    if (showItem) {
                        motherBox.find(".item-" + id).removeClass('d-none');
                    } else {
                        motherBox.find(".item-" + id).addClass('d-none');
                    }
                });
            }, 10);
            motherBox.find(".table-history tbody").animate({ opacity: 1 }, 50);

        } else {
            motherBox.find(".everything").addClass("d-none");
            motherBox.find(".table-history tbody").animate({ opacity: 1 }, 130);
            motherBox.find(".value-item").each(function () {
                var id = $(this).attr('data-id');
                if (!id) {
                    return;
                }

                var showItem = false;
                motherBox.find(".value-" + id).each(function () {
                    if (!$(this).hasClass('everything')) {
                        showItem = true;
                    }
                });
                if (showItem) {
                    motherBox.find(".item-" + id).removeClass('d-none');
                } else {
                    motherBox.find(".item-" + id).addClass('d-none');
                }
            });
        }
    } else {
        if (delay == 'on') {
            motherBox.find(".table-history tbody").animate({ opacity: 0 }, 130);
            setTimeout(function () {
                motherBox.find(".everything").removeClass("d-none");
                motherBox.find(".title-item").removeClass("d-none");
            }, 130);
            motherBox.find(".table-history tbody").animate({ opacity: 1 }, 130);
        } else {
            motherBox.find(".title-item").removeClass("d-none");
            motherBox.find(".everything").removeClass("d-none");
            motherBox.find(".table-history tbody").animate({ opacity: 1 }, 130);
        }
    }
}

function htmlEntities(str) {
    return String(str).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}

var progressbar = $('.progress-bar'),
    i = 0,
    throttle = 5; // 0-1
var progressing = function () {

    if (i <= 100) {
        if (i > 80) {
            throttle = 0.1; // 0-1
            // bar.addClass("done")
        } else {
            throttle = 5; // 0-1
        }
        var r = Math.random();
        if (i < 50) {
            r = r + 2;
        }
        progressbar.css("width", i + "%");
        if (r < throttle) { // Simulate d/l speed and uneven bitrate
            i = i + r;
        }
        if (progressbar.hasClass("done")) {
            progressbar.css("width", "100%");
            return;
        }
        requestAnimationFrame(progressing);
    } else {
        // progressbar.addClass("done");
        progressbar.css("width", "100%");
        requestAnimationFrame(progressing);
    }
};

function ProgressStart() {
    i = 0;
    progressbar.removeClass("done");
    progressbar.removeClass("d-none")
    progressbar.css("width", "0");
    progressing();
}

function ProgressDone() {
    progressbar.addClass("done");
    setTimeout(function () {
        progressbar.addClass("d-none")
    }, 320);
}
