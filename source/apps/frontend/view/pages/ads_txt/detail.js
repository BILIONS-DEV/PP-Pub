$(document).ready(function () {
    if ($("#resultLoadMissingLine").get(0).clientHeight > 600) {
        $("#resultLoadMissingLine").css({"max-height": "600px","height": "600px","overflow-y": "auto"})
    } else {
        $("#resultLoadMissingLine").css({"max-height": "","height": "","overflow-y": ""})
    }

    $("#resultLoadMissingLine").on("click", ".copy-line", function () {
        let copyType = $(this).data("type");
        Copy(copyType);
    });

    new SubmitForm("scanAds", function (resp) {
        if (resp.status === "success") {
            NoticeSuccess(resp.message)
        } else {
            new NoticeError(resp.message);
        }
        Load();
    }, "/ads_txt/scan");

    new SubmitForm("SaveAdsTxt", SaveAdsTxt);
});

function SaveAdsTxt(response) {
    switch (response.status) {
        case "error":
            new AlertError(response.message)
            break
        case "success":
            NoticeSuccess(response.message)
            break
        default:
            new AlertError("Undefined");
            break
    }
    Load();
}

function Copy(copyType) {
    let lines = ""
    let className = ".ads-line"
    if (copyType !== "entire") {
        className = ".ads-line-missing"
    }
    $(className).each(function (i) {
        let line = $(this).text(); // This is your rel value
        if (line) {
            lines += line + "\n"
        }
    });
    CopyTextToClipboard(lines);
    let buttonCopy = $("#CopyAdsTxt");
    let buttonCopyText = buttonCopy.html();
    buttonCopy.attr('disabled', true).text("Copied!");
    setTimeout(function () {
        buttonCopy.attr('disabled', false).html(buttonCopyText);
    }, 1000);
}

function Load() {
    let postData = {
        did: did
    }
    $.ajax({
        url: "/ads_txt/load",
        type: "POST",
        dataType: "JSON",
        data: postData,
        beforeSend: function (xhr) {
            // alert("Loading")
        },
        error: function (jqXHR, exception) {
            const msg = AjaxErrorMessage(jqXHR, exception);
            new AlertError("AJAX ERROR: " + msg);
        },
        complete: function (res) {
            $("#resultLoadMissingLine").html(res.responseJSON.data_object.html);
            if ($("#resultLoadMissingLine").get(0).clientHeight > 600) {
                $("#resultLoadMissingLine").css({"max-height": "600px","height": "600px","overflow-y": "auto"})
            } else {
                $("#resultLoadMissingLine").css({"max-height": "","height": "","overflow-y": ""})
            }
            $("#lastScanAds").html(res.responseJSON.data_object.lastScanAdsTxt);
        }
    });
}