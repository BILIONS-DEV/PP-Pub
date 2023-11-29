// Show clock on left sidebar
$(document).ready(function () {
    var customtimestamp = parseInt($("#jqclock").data("time"));
    $("#jqclock").clock({"langSet": "en", "timestamp": customtimestamp, "timeFormat": "H:i:s", "dateFormat": "D, M j, Y"});
});

// Set select2 default
$.fn.select2.defaults.set("theme", "bootstrap");
if ($(".select2").length > 0) {
    $(document).ready(function () {
        $(".select2").select2({
            width: "auto",
//            closeOnSelect: false,
        });
        $(".select2-multiple").select2({
            width: "auto",
            closeOnSelect: false,
        });
    });
//    container.on('select', function () {
//        if (!container.isOpen()) {
//            return;
//        }
//
//        self.setClasses();
//
//        if (!(self.options.get('multiple') && !self.options.get('closeOnSelect'))) {
//            self.highlightFirstItem();
//        }
//    } );
}

function get_On_Off_CheckBox(element) {
    if (element.closest('.toggle').hasClass('off')) {
        return false;
    } else {
        return true;
    }
}

//function notifiSuccess(success = '') {
//    swal("Good job!", success, {
//        icon: "success",
//        html: true,
//        buttons: {
//            confirm: {
//                className: 'btn btn-success'
//            }
//        }
//    });
//}

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

function notifiResult(result) {
    if (result.success == true || result == true) {
        show_notifi(result);
    } else if (result.success) {
        show_notifi(result);
    } else {
        if (result.warning) {
            notifiWarning(result.warning);
        } else if (result.error) {
            notifiError(result.error);
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
    }, 1500);
}

function ShowNotifiSuccess(result) {
    var content = {};


    var state = 'success';
    content.title = 'Good Job';
    if (result == true) {
        content.message = 'TRUE';
    } else {
        content.message = result;
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
    }, 1500);
}

function loading() {
    return '<div class="__blur"><img class="avatar-small" src="/themes/atlantis/assets/img/Loading.svg" alt="loading"/></div>';
}
function _loading() {
    return '<div class="_blur"><img class="avatar-sm" style="transform: scale(2); position: absolute;top: 44%;" src="/themes/atlantis/assets/img/Loading.svg" alt="loading"/></div>';
}
function LoadingCenter() {
    return '<div class="text-center"><img class="avatar-sm" style="transform: scale(2);" src="/themes/atlantis/assets/img/Loading.svg" alt="loading"/></div>';
}
function __loading() {
    return '<img class="avatar-sm" src="/themes/atlantis/assets/img/loadding.svg" alt="loadding"/>';
}

//function extractEmails(text) {
//    return text.match(/([a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\.[a-zA-Z0-9_-]+)/gi);
//}


// copy text in element
function copyToClipboard(element) {
    var $temp = $("<input>");
    $("body").append($temp);
    $temp.val($(element).html()).select();
    document.execCommand("copy");
    $temp.remove();
}