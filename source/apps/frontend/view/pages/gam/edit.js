$(document).ready(function () {
    firstLoad()
    $(".buttonPushLineItem").on("click", ".btn-push-line-item", function () {
        let GamNetworkId = [];
        // $('.box-push-line-item table tr').each(function () {
        //     let id = $(this).attr("data-id")
        //     if (id !== undefined) {
        //         GamNetworkId.push(parseInt(id))
        //     }
        // });
        $('table > tbody.networkPushLineItem  > tr').each(function (index, tr) {
            let status = $(this).data("status-push")
            if (status !== 1 && status !== 3) {
                GamNetworkId.push(parseInt($(this).data("id")))
            }
        });
        PushLineItem(GamNetworkId)
    })

    $(".buttonCheckApi").on("click", function () {
        const buttonElement = $(this);
        CheckApiAcess(buttonElement)
    })

    SelectNetwork();
    // LoadHistory();
});

function firstLoad() {
    let gamId = parseInt($("#flexSwitchCheckChecked").data("gam-id"));
    $(".load-history").attr("data-id", gamId)
    Reload(gamId)
}

function CheckApiAcess(buttonElement) {
    const submitButtonText = buttonElement.text();
    const submitButtonTextLoading = "Loading...";
    let gamId = 0;
    let listNetwork = [];
    $('input[id=flexSwitchCheckChecked]').each(function () {
        if ($(this).prop("checked")) {
            gamId = parseInt($(this).data("gam-id"));
            listNetwork.push(parseInt($(this).val()));
        }
    })
    $.ajax({
        url: "/gam/checkApiAccess",
        type: "POST",
        dataType: "json",
        contentType: "application/json",
        data: JSON.stringify({
            gam_id: gamId,
            list_network: listNetwork,
        }),
        beforeSend: function (xhr) {
            buttonElement.attr('disabled', true).text(submitButtonTextLoading);
        },
        success: function (json) {
            switch (json.status) {
                case "error":
                    new AlertError(json.message);
                    break
                case "success":
                    json.networks.forEach(function (network) {
                        if (network.status === 2 && network.api_access === 1) {
                            new NoticeSuccess("Check API Access: " + network.network_name + " API Access Enable");
                        } else if (network.status === 2 && network.api_access !== 1) {
                            new NoticeError("Check API Access: " + network.network_name + " API Access Disabled");
                        }
                    });
                    // ReloadApiAccess(json.networks);
                    ReloadApiAccess(json.networks);
                    ReloadPushLineItemBox(json.networks)
                    ReloadButtonPushLineItem(json.networks)
                    $(".box-push-line-item").removeClass("d-none");
                    $(".table-api-access").removeClass("d-none");
                    // $(buttonElement).addClass("d-none");
                    CheckStep(json.networks)
                    break
                default:
                    new AlertError("Undefined");
                    break
            }
            buttonElement.attr('disabled', false).text(submitButtonText);
        },
        error: function (xhr) {
            console.log(xhr)
            buttonElement.attr('disabled', false).text(submitButtonText);
        }
    })
}

function SelectNetwork() {
    $(".networkSelect").on("change", ".select-network", function () {
        let gamId = parseFloat($(this).data("gam-id"));
        let networkId = parseFloat($(this).val());
        let select = $(this).prop("checked")
        $.ajax({
            url: "/gam/select-network",
            type: "POST",
            dataType: "json",
            contentType: "application/json",
            data: JSON.stringify({
                gam_id: gamId,
                network_id: networkId,
                select: select
            }),
            success: function (json) {
                switch (json.status) {
                    case "error":
                        new AlertError(json.message);
                        break
                    case "success":
                        new NoticeSuccess(json.message);
                        Reload(gamId)
                        break
                    default:
                        new AlertError("Undefined");
                        break
                }
            },
            error: function (xhr) {
                console.log(xhr)
            }
        })
    });

    function ReloadSelectBox(networks) {
        let html = ""
        networks.forEach(function (network) {
            html += `
             <tr>
                <td class="fw-normal">${network.network_id}</td>
                <td class="fw-normal">${network.network_name}</td>
                <td class="clearfix">
                    <div class="form-check form-switch mb-0 text-end float-end">
                        <input class="form-check-input select-network" type="checkbox" id="flexSwitchCheckChecked"
                               data-gam-id="${network.gam_id}"
                               value="${network.id}"
                               ${network.status === 2 ? 'checked="" disabled' : ''} />
                    </div>
                </td>
            </tr>`
        })
        $(".networkSelect").html(html.trim())
    }
}

function CheckStep(networks) {
    let checkNetworkSelected = 0
    let checkApiAccess = 0
    let checkPushLineItem = 0
    networks.forEach(function (network) {
        if (network.status === 2) {
            checkNetworkSelected++
            if (network.api_access === 1) {
                checkApiAccess++
            }
            if (network.push_line_item === 1) {
                checkPushLineItem++
            }
        }
    })
    if (checkNetworkSelected > 0) {
        $(".step-select").addClass("active").find("span.circle").removeClass("bg-gray-200").addClass("bg-primary")
        if (checkApiAccess === checkNetworkSelected) {
            $(".step-api-access").addClass("active").find("span.circle").removeClass("bg-gray-200").addClass("bg-primary")
            if (checkPushLineItem === checkNetworkSelected) {
                $(".step-push-line").addClass("active").find("span.circle").removeClass("bg-gray-200").addClass("bg-primary")
            } else {
                if ($(".step-push-line").hasClass("active")) {
                    $(".step-push-line").removeClass("active").find("span.circle").removeClass("bg-primary").addClass("bg-gray-200")
                }
            }
        } else {
            if ($(".step-api-access").hasClass("active")) {
                $(".step-api-access").removeClass("active").find("span.circle").removeClass("bg-primary").addClass("bg-gray-200")
            }
            if ($(".step-push-line").hasClass("active")) {
                $(".step-push-line").removeClass("active").find("span.circle").removeClass("bg-primary").addClass("bg-gray-200")
            }
        }
    } else {
        if ($(".step-select").hasClass("active")) {
            $(".step-select").removeClass("active").find("span.circle").removeClass("bg-primary").addClass("bg-gray-200")
        }
        if ($(".step-api-access").hasClass("active")) {
            $(".step-api-access").removeClass("active").find("span.circle").removeClass("bg-primary").addClass("bg-gray-200")
        }
        if ($(".step-push-line").hasClass("active")) {
            $(".step-push-line").removeClass("active").find("span.circle").removeClass("bg-primary").addClass("bg-gray-200")
        }
    }
}

function Reload(gamId) {
    $.ajax({
        url: "/gam/get-networks",
        type: "POST",
        dataType: "json",
        contentType: "application/json",
        data: JSON.stringify({
            gam_id: parseFloat(gamId)
        }),
        success: function (networks) {
            if (networks.length > 0) {
                // ReloadSelectBox(networks)
                ReloadApiAccess(networks);
                ReloadPushLineItemBox(networks)
                ReloadButtonPushLineItem(networks)
                CheckSelectedNetwork(networks)
                CheckStep(networks)
            }
        },
        error: function (xhr) {
            console.log(xhr)
        }
    })
}

function CheckSelectedNetwork(networks) {
    let checkSelected = false
    networks.forEach(function (network) {
        if (network.status === 2) {
            checkSelected = true
        }
    })
    if (!checkSelected) {
        $(".buttonCheckApi").addClass("d-none");
        $(".table-api-access").addClass("d-none");
        $(".btn-push-line-item").addClass("d-none");
        $(".box-push-line-item").addClass("d-none");
    }
}

function ReloadButtonPushLineItem(networks) {
    let button = ""
    let textButton = ""
    let checkNetworkSelected = 0
    let checkPushLineItem = 0
    let checkDisableButton = 0
    let checkHiddenButton = 0
    networks.forEach(function (network) {
        if (network.status === 2 && network.api_access === 1) {
            checkNetworkSelected++
            if (network.push_line_item !== 1 && network.push_line_item !== 3) {
                checkPushLineItem++
            } else if (network.push_line_item === 3) {
                checkDisableButton++
            } else if (network.push_line_item === 1) {
                checkHiddenButton++
            }
        }
    })
    if (checkNetworkSelected > 0) {
        if (checkPushLineItem > 0) {
            textButton = "Push Line Item"
            $(".btn-push-line-item").removeAttr("disabled");
            $(".btn-push-line-item").removeClass("d-none");
        } else {
            if (checkDisableButton > 0) {
                textButton = "In Processing"
                $(".btn-push-line-item").attr("disabled", true);
                $(".btn-push-line-item").removeClass("d-none");
            } else {
                textButton = "Push Line Item"
                $(".btn-push-line-item").removeAttr("disabled");
                if (checkHiddenButton === checkNetworkSelected) {
                    $(".btn-push-line-item").addClass("d-none");
                } else {
                    $(".btn-push-line-item").removeClass("d-none");
                }
            }
        }
        $(".btn-push-line-item").text(textButton)
    } else {
        $(".btn-push-line-item").addClass("d-none");
    }
}

function ReloadApiAccess(networks) {
    let html = ``;
    let EnableCheckApiAccess = false;
    networks.forEach(function (network) {
        if (network.status === 2 && network.api_access !== 1) {
            EnableCheckApiAccess = true;
        }
        if (network.status === 2) {
            html += `
            <tr>
                    <td class="fw-normal">${network.network_id}</td>
                    <td class="fw-normal">${network.network_name}</td>
                    <td class="clear">
            ${(() => {
                if (network.api_access === 1) {
                    return `
                                <button class="btn btn-sm btn-icon btn-success circle circle-sm float-end" disabled>
                                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16"
                                             fill="currentColor" class="bi bi-check-lg" viewBox="0 0 16 16">
                                            <path d="M13.485 1.431a1.473 1.473 0 0 1 2.104 2.062l-7.84 9.801a1.473 1.473 0 0 1-2.12.04L.431 8.138a1.473 1.473 0 0 1 2.084-2.083l4.111 4.112 6.82-8.69a.486.486 0 0 1 .04-.045z"/>
                                        </svg>
                                </button>
                                `
                } else {
                    return `
                                <button class="btn btn-sm btn-icon btn-danger circle circle-sm float-end"
                                    disabled>
                                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16"
                                         fill="currentColor" class="bi bi-x-lg" viewBox="0 0 16 16">
                                        <path d="M1.293 1.293a1 1 0 0 1 1.414 0L8 6.586l5.293-5.293a1 1 0 1 1 1.414 1.414L9.414 8l5.293 5.293a1 1 0 0 1-1.414 1.414L8 9.414l-5.293 5.293a1 1 0 0 1-1.414-1.414L6.586 8 1.293 2.707a1 1 0 0 1 0-1.414z"/>
                                    </svg>
                                </button>
                                `
                }
            })()}
                    </td>
            </tr>`
        }
        if (EnableCheckApiAccess) {
            $(".buttonCheckApi").removeClass("d-none");
        } else {
            $(".buttonCheckApi").addClass("d-none");
        }
    })
    $(".table-api-access").removeClass("d-none");
    $(".btn-push-line-item").removeClass("d-none");
    $(".box-push-line-item").removeClass("d-none");
    $(".api_access").html(html.trim());
}

function ReloadPushLineItemBox(networks) {
    let html = ""
    networks.forEach(function (network) {
        if (network.status === 2 && network.api_access === 1) {
            let text = ""
            switch (network.push_line_item) {
                case 1:
                    text = `<span class="text-success">Pushed</span>`
                    break
                case 2:
                    text = `<span class="text-danger">Unachievable</span>`
                    break
                case 3:
                    text = `<span>In Process</span>`
                    break
                default:
                    text = `<span>Not been pushed ever</span>`
                    break
            }
            html += `
                 <tr data-id="${network.id}" data-status-push="${network.push_line_item}">
                    <td class="fw-normal">${network.network_id}</td>
                    <td class="fw-normal">${network.network_name}</td>
                    <td class="fw-normal">
                        <span class="fs-12 float-end notify-push-line-item">${text}</span>
                    </td>
                </tr>`
        } else if (network.status === 2 && network.api_access !== 1) {
            let text = "Please Check Api Enable"
            html += `
                 <tr data-id="${network.id}" data-status-push="5">
                    <td class="fw-normal">${network.network_id}</td>
                    <td class="fw-normal">${network.network_name}</td>
                    <td class="fw-normal">
                        <span class="fs-12 float-end notify-push-line-item">${text}</span>
                    </td>
                </tr>`
        }
    })
    $(".networkPushLineItem").html(html.trim())
}

function PushLineItem(GamNetworkId) {
    $.ajax({
        url: "/gam/pushLine",
        type: "POST",
        data: JSON.stringify({
            gam_network_ids: GamNetworkId,
        }),
        success: function (json) {
            switch (json.status) {
                case "error":
                    new AlertError("Error!");
                    break
                case "success":
                    NoticeSuccess("Line item is in process")
                    // $(".btn-push-line-item").addClass("d-none")
                    $(".btn-push-line-item").attr("disabled", true).text("In Processing")
                    $('.box-push-line-item table tr').each(function () {
                        let notify = $(this).find('.notify-push-line-item')
                        let status = $(this).data("status-push")
                        if (status !== 1 && status !== 3) {
                            notify.text("In Progress")
                        }
                    });
                    break
                default:
                    new AlertError("Undefined");
                    break
            }
        },
        error: function (xhr) {
            console.log(xhr)
        }
    })
}