const ALERT = "alert", NOTICE = "notice"
module.exports = {
    ALERT: ALERT, NOTICE: NOTICE,
    DoneWithNotify: function (message) {
        done.Notice(message)
    },
    DoneWithAlert: function (message, functionCallback) {
        done.Alert(message, functionCallback)
    },
    WarningWithNotify: function (message) {
        warning.Notice(message)
    },
    WarningWithAlert: function (message, functionCallback) {
        warning.Alert(message, functionCallback)
    },
    ErrorWithNotify: function (message) {
        error.Notice(message)
    },
    ErrorWithAlert: function (message, functionCallback) {
        error.Alert(message, functionCallback)
    },
}

let done = {
    Notice: function (message) {
        const content = {
            title: "Good Job",
            message: message,
            icon: "fa fa-check"
        };
        const state = 'success';
        const placementFrom = 'top';
        const placementAlign = 'right';
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
    },
    Alert: function (message, functionCallback) {
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
}

let warning = {
    Notice: function (message) {
        const content = {
            title: "Warning",
            message: message,
            icon: "fa fa-close"
        };
        const state = 'warning';
        const placementFrom = 'top';
        const placementAlign = 'right';
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
    },
    Alert: function (message, functionCallback) {
        let alertModal = swal("Ohhh !!!", message, {
            icon: "warning",
            // html: false,
            buttons: {
                confirm: {
                    className: 'btn btn-lg btn-warning'
                }
            }
        });
        if (functionCallback !== undefined && IsFunction(functionCallback)) {
            alertModal.then(() => {
                functionCallback();
            });
        }
    }
}

let error = {
    Notice: function (message) {
        const content = {
            title: "Error",
            message: message,
            icon: "fa fa-close"
        };
        const state = 'danger';
        const placementFrom = 'top';
        const placementAlign = 'right';
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
    },
    Alert: function (message, functionCallback) {
        let alertModal = swal("Ohhh, Sorry !!!", message, {
            icon: "error",
            // html: false,
            buttons: {
                confirm: {
                    className: 'btn btn-lg btn-danger'
                }
            }
        });
        if (functionCallback !== undefined && IsFunction(functionCallback)) {
            alertModal.then(() => {
                functionCallback();
            });
        }
    }
}