import dataTable from "../../jspkg/datatable";

const formID = "#formFilterUser";
const tableID = "#tableRules";

$(document).ready(function () {
    /**
     * Load Datatable
     * @type {string}
     */
    LoadDatatable();
    ChooseRule();
    CreateRule();

    $("#createModal").on("click", ".modal-close", function () {
        $("body").removeClass('customize-box');
    });
    
    $(tableID).on("click", "tbody td a.remove", function (e) {
        var id = $(this).data("id");
        var url = $(this).data("url");
        var isCheck = confirm('Are you sure delete!"');
        if (isCheck) {
            Delete(id, url);
        }
    });

});


function LoadDatatable() {
    $(formID).find(".submit").on("click", function (e) {
        dataTable.Render(formID, true);
    });
    dataTable.Render(formID);
}

function ChooseRule() {
    $(".box-rule").on("click", function (e) {
        $(".btn-create").attr("disabled", false);
        $(".box-rule").removeClass("selected");
        $(this).addClass("selected");
    });
}

function CreateRule() {
    $('.btn-create').modal({backdrop: 'static', keyboard: false});

    $(".btn-create").on("click", function (e) {
        let type = $(".modal-body .selected").data("type");
        if (type === "") {
            return;
        }
        let win = window.open('/' + type + '/add', '_blank');
        if (win) {
            //Browser has allowed it to be opened
            win.focus();
        } else {
            //Browser has blocked it
            alertError('Please allow popups for this website');
        }
    });
}

function Delete(id, url) {
    $.ajax({
        url: url,
        type: "POST",
        dataType: "JSON",
        contentType: "application/json",
        data: JSON.stringify({
            id: id
        }),
        beforeSend: function (xhr) {
            xhr.overrideMimeType("text/plain; charset=x-user-defined");
        },
        error: function (jqXHR, exception) {
            const msg = AjaxErrorMessage(jqXHR, exception);
            new AlertError("AJAX ERROR: " + msg);
        }
    }).done(function (result) {
        switch (result.status) {
            case "success":
                NoticeSuccess("Rule has been removed successfully");
                dataTable.Render(formID);
                break;
            case "err":
                new AlertError(result.message);
        }
    });
}