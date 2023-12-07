const formID = "#formFilterAdTag";
const tableID = "#tableAdTags";
const filterURL = "/supply/adtag";
let firstLoad = true;
const formConnectionID = "#formFilterConnection";
const tableConnectionID = "#tableConnection";
const filterConnectionURL = "/supply/connection";
let firstLoadConnection = true;
let currentURL = "/supply-v2/setup";
let currentTab = 1;
let idInventory;
let list_select_module = [];
let moduleParams = [];
let moduleParamsDefault = [];

function selectAdRefresh() {
    function showAdRefreshTime(opts = {}) {
        if ($("#BoxAdtag").find("#ad_refresh").val() === "individual_tag_configuration") {
            $("#BoxAdtag").find('.individual_tag_configuration').removeClass("d-none");

            if (opts.autoChange) {
                $("#BoxAdtag").find("#ad_refresh_time").val(30);
            }

        } else {
            $("#BoxAdtag").find('.individual_tag_configuration').addClass("d-none");
        }
    }

    showAdRefreshTime();
    $("#BoxAdtag").find("#ad_refresh").on("change", function () {
        showAdRefreshTime({ autoChange: true });
    });
}

$(document).ready(function () {
    idInventory = $("#inventory_id").attr("data-inventory-id");
    currentURL += "?id=" + idInventory;
    ChangeTab();
    HandleSubmit();
    HandleEnable();
    GetModuleDefault();
    HandleCollapse();
    HandleTabTool();
    HandleTabPlayZone();
    // LoadModules()
    $('#select-module').on('select2:select', function (e) {
        const selected = $('#select-module').val(); // Select the option with a value of '1'
        const modules = $('#module').val();
        modules.push(selected);
        if (!$("#collapseSettingBox").find('#ModuleParamBox-' + selected).length) {
            $('#module').val(modules);
            $('#module').trigger('change');
            SelectModuleUserId(e);
        }
        $('#select-module').val('');
        $('#select-module').trigger('change');
    });

    $('#module').on('select2:unselect', function (e) {
        let data = e.params.data;
        RemoveModuleParams(data.id);
    });

    $(formID).find(".submit").on("click", function (e) {
        e.preventDefault();
        GetTable(true, true);
    });

    GetTable(false);

    $(tableID).on("click", "tbody td div.btn-group a.remove", function (e) {
        const id = $(this).data("id");
        const isCheck = confirm('Are you sure archive!"');
        if (isCheck) {
            Delete(id);
        }
    });

    $("#module-params").on("click", ".export-dropdown a.remove_module_param", function () {
        let idEl = $(this).data("id");
        RemoveModuleParams(idEl);
        UpdateSelectedModule(idEl.toString());
    });

    $('[data-bs-toggle="popover"]').popover({
        html: true,
        sanitize: false,
    });
    // $("[data-toggle='toggle']").bootstrapToggle('destroy')
    // $("[data-toggle='toggle']").bootstrapToggle();

    $("#module-params").on("click", "div.bidder-box .d a.rm_c", function () {
        const module_id = $(this).attr("data-id");
        RemoveModule(module_id);
    });

    $("input").change(function () {
        $(this).removeClass("border-danger");
        $(this).next("span.invalid-feedback").text("");
    });

    changeLoadAdType();
    loadTabConnection();
    // anno_guided_tour()
    // LoadHistory()
    selectAdSizeCopyTag();
    initPreventClickOverSidebar();
});

function initPreventClickOverSidebar() {
    document.querySelectorAll('body').forEach(muzeHamburger => {
        muzeHamburger.addEventListener('click', (e) => {
            if (document.querySelector('.customize-sidebar') && !document.querySelector('.customize-sidebar').contains(e.target)) {
                if ($("#BoxAdtag").hasClass("close-sidebar")) {
                    document.querySelector('body').classList.remove('customize-box');
                } else {
                    document.querySelector('body').classList.add('customize-box');
                }
            }
        });
    });
    $("#btn-close").on("click", function () {
        $("#BoxAdtag").addClass("close-sidebar");
    });
    $("#BoxAdtag").on("click", ".btn-close-sidebar", function () {
        $("#btn-close").click();
    });
}

function anno_guided_tour() {
    var step1 = new Anno({
        // target: '#myModal .modal-content',
        target: '#step1',
        content: "Because Anno is very flexible, you can make it work with all sorts of other libraries, e.g. Bootstrap",
        position: { right: '165px', top: '-15px' },
        // position: 'left',
        arrowPosition: 'right',
        // showOverlay: function () {}, // the modal already has one, so disable the anno.js one
        buttons: {
            text: 'Next',
            click: function () {
                $("#step1").click();
                step1.hide();
                step2.show();
            }
        }
    });
    step1.show();
    var step2 = new Anno({
        target: '#step2',
        content: "Because Anno is very flexible, you can make it work with all sorts of other libraries, e.g. Bootstrap",
        position: "left",
        buttons: {
            text: 'Next',
            click: function () {
                $("#step2").click();
                step2.hide();
                step3.show();
            }
        }
    });
    $('#copyAdTagModal').on('hide.bs.modal', function () {
        step2.hide();
    });
    var step3 = new Anno({
        target: '#tagTabs',
        content: "Because Anno is very flexible, you can make it work with all sorts of other libraries, e.g. Bootstrap",
        position: "left",
        // buttons: [AnnoButton.BackButton]
    });
}

function loadTabConnection() {
    $(formConnectionID).find(".submit").on("click", function (e) {
        e.preventDefault();
        // GetTableConnection()
    });

    // GetTableConnection(false);

    // changeStatusConnection
    $(".connection").change(function (e) {
        // console.log(this.checked);
        changeStatusConnection(this);
    });
}

function changeStatusConnection(clickElement) {
    ProgressStart();
    let bidderID = $(clickElement).attr("data-bidder-id");
    let inventoryID = $(clickElement).attr("data-inventory-id");
    let status = 0;
    if (clickElement.checked) {
        status = 1;
    } else {
        status = 2;
    }
    $.ajax({
        url: "/supply/change-status-connection",
        type: "POST",
        dataType: "JSON",
        data: {
            inventory_id: parseInt(inventoryID),
            bidder_id: parseInt(bidderID),
            status: status
        },
        beforeSend: function (xhr) {
            $(clickElement).addClass("disabled");
        },
        complete: function (res) {
        },
        error: function (jqXHR, exception) {
            new AlertError("AJAX ERROR: " + ajaxErrorMessage(jqXHR, exception));
        },
        success: function (responseJSON) {
            ProgressDone();
            if (responseJSON.status === "error") {
                AlertError(responseJSON.message);
            } else {
                console.log(responseJSON);
                countBidderWaiting();
                NoticeSuccess("Changed!");
                if (status == 1) {
                    $(clickElement).closest(".item").find(".status").addClass("pp23").removeClass("pp-23");
                    $(clickElement).closest(".item").find(".status-text").text("Live");
                } else {
                    $(clickElement).closest(".item").find(".status").addClass("pp-23").removeClass("pp23");
                    $(clickElement).closest(".item").find(".status-text").text("Waiting");
                }
            }
        }
    });
}

function countBidderWaiting() {
    var waiting = 0;
    $("#nav-connection").find(".item").each(function () {
        if (!$(this).find("input").is(":checked")) {
            waiting = waiting + 1;
        }
    });
    if (waiting) {
        $(".numberCircle ").text(waiting).removeClass("d-none")
    } else {
        $(".numberCircle ").text(waiting).addClass("d-none")
    }
}

function changeLoadAdType() {
    checkLoadAdType();
    $("#load_ad_type").on("change", function () {
        checkLoadAdType();
    });
}

function checkLoadAdType() {
    let value = $("#load_ad_type").val();
    if (value === "lazyload") {
        $(".load_ad_type_lazyload").removeClass("d-none");
        $("#formGeneral").find("#ad_refresh_type option").prop("disabled", true);
        $("#formGeneral").find("#ad_refresh_type").find('option[value="signal reload"]').prop("disabled", false).prop("selected", true);
        $("#formGeneral").find("#ad_refresh_type").selectpicker("refresh");
    } else {
        $(".load_ad_type_lazyload").addClass("d-none");
        $("#formGeneral").find("#ad_refresh_type option").prop("disabled", false);
        $("#formGeneral").find("#ad_refresh_type").selectpicker("refresh");
    }
}

function HandleCollapse() {
    let listId = ["#collapseUserIdBox"];
    let url = "/supply/collapse";
    $('#formUserId').on('show.bs.collapse', '.collapse', function (e) {
        let box = e.target.id;
        SendRequestShow(url, box);
    });
    $('#formUserId').on('hide.bs.collapse', '.collapse', function (e) {
        let box = e.target.id;
        SendRequestHide(url, box);
    });
}

function SendRequestHide(url, box) {
    const urlSearchParams = new URLSearchParams(window.location.search);
    const params = Object.fromEntries(urlSearchParams.entries());
    $.ajax({
        url: url,
        type: "POST",
        dataType: "JSON",
        contentType: "application/json",
        data: JSON.stringify({
            is_collapse: 1,
            box_collapse: box,
            page_type: "edit",
            page_id: parseInt(params.id)
        }),
        beforeSend: function (xhr) {
            xhr.overrideMimeType("text/plain; charset=x-user-defined");
        },
        error: function (jqXHR, exception) {
            var msg = AjaxErrorMessage(jqXHR, exception);
            new AlertError("AJAX ERROR: " + msg);
        }
    }).done(function (result) {
    });
}

function SendRequestShow(url, box) {
    const urlSearchParams = new URLSearchParams(window.location.search);
    const params = Object.fromEntries(urlSearchParams.entries());
    $.ajax({
        url: url,
        type: "POST",
        dataType: "JSON",
        contentType: "application/json",
        data: JSON.stringify({
            is_collapse: 2,
            box_collapse: box,
            page_type: "edit",
            page_id: parseInt(params.id)
        }),
        beforeSend: function (xhr) {
            xhr.overrideMimeType("text/plain; charset=x-user-defined");
        },
        error: function (jqXHR, exception) {
            var msg = AjaxErrorMessage(jqXHR, exception);
            new AlertError("AJAX ERROR: " + msg);
        }
    }).done(function (result) {
    });
}

function UpdateSelectedModule(idRemove) {
    let valuesSelected = $('#module').val();
    let currentId = valuesSelected.indexOf(idRemove)
    if (currentId !== -1) {
        valuesSelected.splice(currentId, 1)
    }
    $('#module').val(valuesSelected).trigger('change')
}

function GetModuleDefault() {
    $('#module-params > div.module_param_box').each(function () {
        $('input', this).each(function () {
            let params = {};
            let option;
            let type = $(this).attr('data-type');
            let name = $(this).attr('data-name');
            option = $(this).attr('data-option');
            switch (type) {
                case "string":
                    params.type = type;
                    params.name = name;
                    params.template = this.value;
                    break;
                case "int":
                    params.type = type;
                    params.name = name;
                    if (this.value !== "") {
                        params.template = parseInt(this.value);
                        if (!params.template) {
                            params.template = 0;
                        }
                    } else {
                        params.template = 0;
                    }
                    break;
                case "float":
                    params.type = type;
                    params.name = name;
                    if (this.value !== "") {
                        params.template = parseFloat(this.value);
                        if (!params.template) {
                            params.template = 0;
                        }
                    } else {
                        params.template = 0;
                    }
                    break;
                case "json":
                    params.type = type;
                    params.name = name;
                    try {
                        params.template = JSON.parse(this.value);
                    } catch (err) {
                        params.template = this.value;
                    }
                    break;
                case "boolean":
                    params.type = type;
                    params.name = name;
                    params.template = this.value;
                    break;
            }

            let bidderParam = {}
            bidderParam.storage = []
            bidderParam.params = []
            let id = parseInt($(this).attr('data-id'))
            let currentId = moduleParamsDefault.map(item => {
                return item.id
            }).indexOf(id)

            if (currentId === -1) {
                bidderParam.id = id
                // bidderParam.name = $(this).attr('data-module-name').toLowerCase()
                bidderParam.name = $(this).attr('data-module-name')
                switch (option) {
                    case "storage":
                        bidderParam.storage = [params, ...bidderParam.storage]

                        break
                    case "params":
                        bidderParam.params = [params, ...bidderParam.params]
                        break
                }
                moduleParamsDefault.push(bidderParam)
            } else {
                switch (option) {
                    case "storage":
                        moduleParamsDefault[currentId].storage.push(params)
                        break
                    case "params":
                        moduleParamsDefault[currentId].params.push(params)
                        break
                }

            }
        })
    });
    moduleParamsDefault = Object.values(moduleParamsDefault);
}

function HandleSubmit() {
    submitFormSetup("formGeneral", "/supply/setup", SetUp)
    submitFormConsent("formConsent", "/supply/setupConsent", SetUp)
    submitFormUser("formUserId", "/supply/setupUserId", SetUp)
}

function ChangeTab() {
    $(".custom-tab").on("click", ".nav-link", function () {
        $(".custom-tab").find(".nav-link").removeClass("pp-4").removeClass('at-1')
        $(this).addClass("pp-4")
        var tab = $(this).attr("data-tab")
        if (tab != "1") {
            $(this).addClass("at-1")
        }
    })

    $("#tab-adtag").on("click", function () {
        let postData = $(formID).serializeObject()
        postData.tab = 1
        currentTab = 1
        $("#btn-add-adtag").removeClass("d-none")
        $("#btn-copy-adtag").removeClass("d-none")
        $(".submit-config").closest(".dm6 ").addClass("d-none")
        $(".submit-consent").closest(".dm6 ").addClass("d-none")
        $(".history-a").addClass("d-none")
        makeParamsUrl(postData)

    })

    $("#tab-general").on("click", function () {
        let postData = $(formID).serializeObject()
        postData.tab = 2
        currentTab = 2
        $("#btn-add-adtag").addClass("d-none")
        $("#btn-copy-adtag").addClass("d-none")
        $(".submit-config").closest(".dm6 ").removeClass("d-none")
        $(".submit-consent").closest(".dm6 ").addClass("d-none")
        $(".load-history").attr("data-object", "inventory_config_fe").removeClass("d-none")
        makeParamsUrl(postData)
    })

    $("#tab-consent").on("click", function () {
        let postData = $(formID).serializeObject()
        postData.tab = 3
        currentTab = 3
        $("#btn-add-adtag").addClass("d-none")
        $("#btn-copy-adtag").addClass("d-none")
        $(".submit-consent").closest(".dm6 ").removeClass("d-none")
        $(".submit-config").closest(".dm6 ").addClass("d-none")
        $(".load-history").attr("data-object", "inventory_consent_fe").removeClass("d-none")
        makeParamsUrl(postData)
    })

    $("#tab-adstxt").on("click", function () {
        let postData = $(formID).serializeObject()
        postData.tab = 4
        currentTab = 4
        $("#btn-add-adtag").addClass("d-none")
        $("#btn-copy-adtag").addClass("d-none")
        $(".submit-config").closest(".dm6 ").addClass("d-none")
        $(".submit-consent").closest(".dm6 ").addClass("d-none")
        $(".load-history").attr("data-object", "inventory_adstxt_fe").removeClass("d-none")
        makeParamsUrl(postData)
    })

    $("#tab-connection").on("click", function () {
        let postData = $(formID).serializeObject()
        postData.tab = 5
        currentTab = 5
        $("#btn-add-adtag").addClass("d-none")
        $("#btn-copy-adtag").addClass("d-none")
        $(".submit-config").closest(".dm6 ").addClass("d-none")
        $(".submit-consent").closest(".dm6 ").addClass("d-none")
        $(".load-history").attr("data-object", "inventory_connection_fe").removeClass("d-none")
        makeParamsUrl(postData)
    })

    $("#tab-integration").on("click", function () {
        let postData = $(formID).serializeObject()
        postData.tab = 6
        currentTab = 6
        $("#btn-add-adtag").addClass("d-none")
        $("#btn-copy-adtag").addClass("d-none")
        $(".submit-config").closest(".dm6 ").addClass("d-none")
        $(".submit-consent").closest(".dm6 ").addClass("d-none")
        $(".load-history").addClass("d-none")
        makeParamsUrl(postData)
    })

}

function SelectModuleUserId(e) {
    const data = e.params.data;
    let currentIdx = list_select_module.map(item => {
        return item.id
    }).lastIndexOf(data.id)
    if (currentIdx > -1) {
        let index = list_select_module[currentIdx].index + 1
        list_select_module.push({ id: data.id, index: index })
        LoadModuleParam(data.id, data.text, index)
    } else {
        list_select_module.push({ id: data.id, index: 1 })
        LoadModuleParam(data.id, data.text, 1)
    }
    //Clear select
    // $('#module').val(null).trigger('change')
}

function LoadModuleParam(id, name, index) {
    $.ajax({
        url: "/supply/loadParam",
        type: "GET",
        data: {
            id: id,
            name: name,
            index: index,
        },
        success: function (json) {
            $("#module-params").prepend(json.data)
            InitAfterLoadModule()
        },
        error: function (xhr) {
            console.log(xhr)
        }
    })
}

function LoadModules() {
    var IDs = $("#module").val();
    IDs.forEach(function (id) {
        LoadModuleParam(id, "", 1)
    });
}

function RemoveModule(moduleID) {
    $("#ModuleParamBox-" + moduleID).remove()
    const modules = [];
    $("#module-params").find('.module_param_box').each(function () {
        modules.push($(this).attr('data-id'));
    })
    $('#module').val(modules);
    $('#module').trigger('change');
}

function RemoveModuleParams(id) {
    list_select_module = list_select_module.filter((item) => {
        return !(item.id === id.toString())
    })
    let list = document.querySelectorAll(`#ModuleParamBox-${id}`)
    list.forEach(item => {
        item.remove()
    })
}

function addParamModule() {
    $('#module-params > div.module_param_box').each(function () {
        $('input', this).each(function () {
            let params = {}
            let option = $(this).attr('data-option')

            params.type = $(this).attr('data-type')
            params.name = $(this).attr('data-name')
            params.template = this.value.trim()

            let bidderParam = {}
            bidderParam.storage = []
            bidderParam.params = []
            bidderParam.ab_testing = 0
            bidderParam.volume = 0

            let id = parseInt($(this).attr('data-id'))
            let currentId = moduleParams.map(item => {
                return item.id
            }).indexOf(id)

            if (currentId === -1) {
                bidderParam.id = id
                // bidderParam.name = $(this).attr('data-module-name').toLowerCase()
                bidderParam.name = $(this).attr('data-module-name')
                switch (option) {
                    case "storage":
                        bidderParam.storage = [params, ...bidderParam.storage]

                        break
                    case "params":
                        bidderParam.params = [params, ...bidderParam.params]
                        break
                    case "ab_testing":
                        if (this.checked) {
                            bidderParam.ab_testing = 1
                        } else {
                            bidderParam.ab_testing = 2
                        }
                        break
                    case "volume":
                        bidderParam.volume = parseInt(this.value)
                        break
                }
                moduleParams.push(bidderParam)
            } else {
                switch (option) {
                    case "storage":
                        moduleParams[currentId].storage.push(params)
                        break
                    case "params":
                        moduleParams[currentId].params.push(params)
                        break
                    case "ab_testing":
                        if (this.checked) {
                            moduleParams[currentId].ab_testing = 1
                        } else {
                            moduleParams[currentId].ab_testing = 2
                        }
                        break
                    case "volume":
                        moduleParams[currentId].volume = parseInt(this.value)
                        break
                }

            }
        })
    });
    return Object.values(moduleParams);
}

function HandleEnable() {
    checkFirstLoad()
    GdprEnable();
    CcpaEnable();
    CustomBrandEnable()
    AdRefreshEnable()
    ModuleEnable()
}

function GdprEnable() {
    $("#gdpr").on("change", function () {
        if (this.checked) {
            $("#enable_custom_gdpr").attr("hidden", false);
        } else {
            $("#enable_custom_gdpr").attr("hidden", true);
        }
    });
}

function CcpaEnable() {
    $("#ccpa").on("change", function () {
        if (this.checked) {
            $("#enable_custom_ccpa").attr("hidden", false);
        } else {
            $("#enable_custom_ccpa").attr("hidden", true);
        }
    });
}

function CustomBrandEnable() {
    $("#custom_brand").on("change", function () {
        if (this.checked) {
            $("#enable_custom_brand").attr("hidden", false);
        } else {
            $("#enable_custom_brand").attr("hidden", true);
        }
    });
}

function AdRefreshEnable() {
    $("#ad_refresh").on("change", function () {
        if (this.checked) {
            $("#enable_ad_refresh").attr("hidden", false);
        } else {
            $("#enable_ad_refresh").attr("hidden", true);
        }
    });
}

function ModuleEnable() {
    $(".enable_testing").on("change", function () {
        let name = $(this).data("module-name")
        if (this.checked) {
            $(`div#ab_testing_${name}`).attr("hidden", false);
        } else {
            $(`div#ab_testing_${name}`).attr("hidden", true);
        }
    });
}

function checkFirstLoad() {
    if ($("#gdpr").is(':checked')) {
        $(".gdpr_sub").attr("hidden", false);
    } else {
        $(".gdpr_sub").attr("hidden", true);
    }

    if ($("#ccpa").is(':checked')) {
        $(".ccpa_sub").attr("hidden", false);
    } else {
        $(".ccpa_sub").attr("hidden", true);
    }
}

function InitAfterLoadModule() {
    $("[data-toggle='toggle']").bootstrapToggle('destroy')
    $("[data-toggle='toggle']").bootstrapToggle();
    ModuleEnable()
    $("[data-bs-toggle=popover]").popover();
    // $("[data-bs-toggle=tooltip]").tooltip();
}

function submitFormSetup(formID, url = "", functionCallback) {
    const formElement = $("#" + formID);
    // formElement.find("input").on("click change blur", function (e) {
    //     let inputElement = $(this)
    //     if (inputElement.hasClass("is-invalid")) {
    // inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
    // inputElement.removeClass("is-invalid").parent().removeClass("is-invalid").next(".invalid-feedback").empty();
    //     }
    // });
    // formElement.find("textarea").on("click change blur", function (e) {
    //     let inputElement = $(this)
    //     if (inputElement.hasClass("is-invalid")) {
    //         inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
    // inputElement.removeClass("is-invalid").parent().removeClass("is-invalid").next(".invalid-feedback").empty();
    //     }
    // });
    // formElement.find('.btn-group-radio').on('change', function (e) {
    //     let box = $(this).closest(".box-radio-btn")
    //     if (box.hasClass("is-invalid")) {
    //         box.removeClass("is-invalid").next(".invalid-feedback").empty();
    //     }
    // });
    $(".submit-config").on("click", function (e) {
        ProgressStart()
        e.preventDefault();
        let validate = true
        const buttonElement = $(this);
        const submitButtonText = buttonElement.text();
        const submitButtonTextLoading = "Loading...";

        var postData = formElement.serializeObject();
        postData.ad_refresh_time = parseInt(postData.ad_refresh_time)
        postData.id = parseInt(postData.id)
        postData.inventory_id = parseInt(postData.inventory_id)
        postData.prebid_timeout = parseInt(postData.prebid_timeout)
        postData.gam_account = parseInt(postData.gam_account)
        postData.pb_render_mode = parseInt(postData.pb_render_mode)
        postData.fetch_margin_percent = parseInt(postData.fetch_margin_percent)
        postData.render_margin_percent = parseInt(postData.render_margin_percent)
        postData.mobile_scaling = parseInt(postData.mobile_scaling)
        // postData.user_id = parseInt(postData.user_id)

        if (postData.ad_refresh === "on") {
            postData.ad_refresh = 1
        } else {
            postData.ad_refresh = 2
        }

        if (postData.gam_auto_create === "on") {
            postData.gam_auto_create = 1
        } else {
            postData.gam_auto_create = 2
        }

        if (postData.safe_frame === "on") {
            postData.safe_frame = 1
        } else {
            postData.safe_frame = 2
        }

        if (postData.direct_sales === "on") {
            postData.direct_sales = 1
        } else {
            postData.direct_sales = 2
        }
        $.ajax({
            url: url,
            type: "POST",
            dataType: "JSON",
            contentType: "application/json",
            data: JSON.stringify(postData),
            beforeSend: function (xhr) {
                buttonElement.addClass("disabled").attr('disabled', true).text(submitButtonTextLoading);
            },
            error: function (jqXHR, exception) {
                const msg = AjaxErrorMessage(jqXHR, exception);
                new AlertError("AJAX ERROR: " + msg);
                buttonElement.removeClass("disabled").attr('disabled', false).text(submitButtonText);
            },
            success: function (responseJSON) {
                buttonElement.removeClass("disabled").attr('disabled', false).text(submitButtonText);
            },
            complete: function (res) {
                ProgressDone()
                functionCallback(res.responseJSON, formElement);
            }
        });
    });
}

function submitFormConsent(formID, url = "", functionCallback) {
    var formElement = $("#" + formID);
    $(".submit-consent").on("click", function (e) {
        ProgressStart()
        e.preventDefault();
        let validate = true
        const buttonElement = $(this);
        const submitButtonText = buttonElement.text();
        const submitButtonTextLoading = "Loading...";
        formElement.find("input").on("click change blur", function (e) {
            let inputElement = $(this)
            if (inputElement.hasClass("is-invalid")) {
                inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
                inputElement.removeClass("is-invalid").parent().removeClass("is-invalid").next(".invalid-feedback").empty();
            }
        });
        formElement.find("textarea").on("click change blur", function (e) {
            let inputElement = $(this)
            if (inputElement.hasClass("is-invalid")) {
                inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
                inputElement.removeClass("is-invalid").parent().removeClass("is-invalid").next(".invalid-feedback").empty();
            }
        });
        var postData = formElement.serializeObject();
        postData.id = parseInt(postData.id)
        postData.inventory_id = parseInt(postData.inventory_id)
        // postData.user_id = parseInt(postData.user_id)
        if (postData.gdpr === "on") {
            postData.gdpr = 1
        } else {
            postData.gdpr = 2
        }
        postData.gdpr_timeout = parseInt(postData.gdpr_timeout)
        if (postData.ccpa === "on") {
            postData.ccpa = 1
        } else {
            postData.ccpa = 2
        }
        if (postData.custom_brand === "on") {
            postData.custom_brand = 1
        } else {
            postData.custom_brand = 2
        }
        if (postData.ad_refresh === "on") {
            postData.ad_refresh = 1
        } else {
            postData.ad_refresh = 2
        }
        if (postData.direct_sales === "on") {
            postData.direct_sales = 1
        } else {
            postData.direct_sales = 2
        }
        postData.ccpa_timeout = parseInt(postData.ccpa_timeout)

        $.ajax({
            url: url,
            type: "POST",
            dataType: "JSON",
            contentType: "application/json",
            data: JSON.stringify(postData),
            beforeSend: function (xhr) {
                buttonElement.addClass("disabled").attr('disabled', true).text(submitButtonTextLoading);
            },
            error: function (jqXHR, exception) {
                const msg = AjaxErrorMessage(jqXHR, exception);
                new AlertError("AJAX ERROR: " + msg);
                buttonElement.removeClass("disabled").attr('disabled', false).text(submitButtonText);
            },
            success: function (responseJSON) {
                buttonElement.removeClass("disabled").attr('disabled', false).text(submitButtonText);
            },
            complete: function (res) {
                ProgressDone()
                functionCallback(res.responseJSON, formElement);
            }
        });
    });
}

function submitFormUser(formID, url = "", functionCallback) {
    var formElement = $("#" + formID);
    formElement.find(".submit").on("click", function (e) {
        e.preventDefault();
        let validate = true
        const buttonElement = $(this);
        const submitButtonText = buttonElement.text();
        const submitButtonTextLoading = "Loading...";
        formElement.find("input").on("click change blur", function (e) {
            let inputElement = $(this)
            if (inputElement.hasClass("is-invalid")) {
                inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
                inputElement.removeClass("is-invalid").parent().removeClass("is-invalid").next(".invalid-feedback").empty();
            }
        });
        formElement.find("textarea").on("click change blur", function (e) {
            let inputElement = $(this)
            if (inputElement.hasClass("is-invalid")) {
                inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
                inputElement.removeClass("is-invalid").parent().removeClass("is-invalid").next(".invalid-feedback").empty();
            }
        });
        var postData = formElement.serializeObject();
        postData.inventory_id = parseInt(postData.inventory_id)
        postData.auction_delay = parseInt(postData.auction_delay)
        postData.sync_delay = parseInt(postData.sync_delay)
        postData.id = parseInt(postData.id)
        postData.module_params = addParamModule()

        let listModuleId = []
        let listModuleIdDefault = []

        moduleParams.map(item => {
            listModuleId = [...listModuleId, item.id]
        })
        moduleParamsDefault.map(item => {
            listModuleIdDefault = [...listModuleIdDefault, item.id]
        })
        postData.list_del = listModuleIdDefault.filter(x => !listModuleId.includes(x))
        $.ajax({
            url: url,
            type: "POST",
            dataType: "JSON",
            contentType: "application/json",
            data: JSON.stringify(postData),
            beforeSend: function (xhr) {
                buttonElement.addClass("disabled").attr('disabled', true).text(submitButtonTextLoading);
            },
            error: function (jqXHR, exception) {
                const msg = AjaxErrorMessage(jqXHR, exception);
                new AlertError("AJAX ERROR: " + msg);
                buttonElement.removeClass("disabled").attr('disabled', false).text(submitButtonText);
            },
            success: function (responseJSON) {
                buttonElement.removeClass("disabled").attr('disabled', false).text(submitButtonText);
                // if (responseJSON.status == "error") {
                //     responseJSON.errors.forEach(function (value, index) {
                //         $("#" + value.id).addClass("border-danger")
                //         $("#" + value.id).nextAll("span.invalid-feedback").text(value.message).addClass('text-danger')
                //     })
                // }
            },
            complete: function (res) {
                moduleParams = []
                moduleParamsDefault = []
                GetModuleDefault()
                functionCallback(res.responseJSON, formElement);
            }
        });
    });
}

function SetUp(response, formElement) {
    switch (response.status) {
        case "error":
            if (response.errors.length === 1 && response.errors[0].id === "") {
                new AlertError(response.errors[0].message);
            } else if (response.errors.length > 0) {
                $.each(response.errors, function (key, value) {
                    let inputElement = $("#" + value.id);
                    if (value.id === "prebid_timeout") {

                    }
                    inputElement.addClass("is-invalid").parent().addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
                    inputElement.addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
                });
                $("#" + response.errors[0].id).focus();
                // new AlertError(response.errors[0].message);
            } else {
                new AlertError("Error!");
            }
            break
        case "success":
            new NoticeSuccess("Setup inventory has been successfully")
            break
        default:
            new AlertError("Undefined");
            break
    }
}

function SetUpSticky(response, formElement) {
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
            break
        case "success":
            let id = formElement.find("input#id")
            $(id).val(response.data_object.id)
            let inventoryId = formElement.find("input#inventory_id")
            $(inventoryId).val(response.data_object.inventory_id)
            NoticeSuccess("Setup inventory has been successfully")
            break
        default:
            new AlertError("Undefined");
            break
    }
}

function GetTable(isClickForm = false, refresh = false) {
    const formElement = $(formID);
    let buttonElement = formElement.find(".submit");
    let submitButtonText = buttonElement.text();
    let submitButtonTextLoading = "Loading...";
    let postData = formElement.serializeObject();
    let inventoryId = $("#inventory_id").data("inventory-id")
    let setting = {
        processing: true,
        serverSide: true,
        searching: false,
        destroy: true,
        pagingType: "simple",
        order: [parseInt(postData.order_column), postData.order_dir],
        language: {
            info: "_START_ - _END_ of _TOTAL_",
            infoEmpty: "0 - 0 of 0",
            lengthMenu: "Rows per page:  _MENU_",
            paginate: {
                next: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-chevron-right" viewBox="0 0 16 16"> <path fill-rule="evenodd" d="M4.646 1.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1 0 .708l-6 6a.5.5 0 0 1-.708-.708L10.293 8 4.646 2.354a.5.5 0 0 1 0-.708z"/> </svg>', // or '→'
                previous: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-chevron-left" viewBox="0 0 16 16"> <path fill-rule="evenodd" d="M11.354 1.646a.5.5 0 0 1 0 .708L5.707 8l5.647 5.646a.5.5 0 0 1-.708.708l-6-6a.5.5 0 0 1 0-.708l6-6a.5.5 0 0 1 .708 0z"/> </svg>' // or '←'
            }
        },
        dom: '<"row"<"col-sm-12"tr>><"bottom d-flex p-3 border-top border-gray-200"<"row ms-auto align-items-center px-3"<"col-auto"fl><"col-auto"i><"col-auto"p>>>',
        ajax: {
            url: filterURL,
            type: "POST",
            contentType: "application/json; charset=utf-8",
            data: function (d) {
                if (isClickForm) {
                    postData = formElement.serializeObject();
                    d.length = parseInt(postData.length);
                    d.start = parseInt(postData.start);
                    d.order[0].column = parseInt(postData.order_column);
                    d.order[0].dir = postData.order_dir;
                }
                if (refresh) {
                    d.start = 0;
                    d.order[0].column = 0;
                    d.order[0].dir = "desc";
                }
                formElement.find("[name='length']").val(d.length);
                formElement.find("[name='start']").val(d.start);
                formElement.find("[name='order_column']").val(d.order[0].column);
                formElement.find("[name='order_dir']").val(d.order[0].dir);
                postData.length = d.length;
                postData.start = d.start;
                postData.inventory_id = parseInt(inventoryId);
                d.postData = postData;
                return JSON.stringify(d);
            },
            beforeSend: function (xhr) {
                if (isClickForm) {
                    buttonElement.attr('disabled', true).text(submitButtonTextLoading);
                }
            },
            dataSrc: function (json) {
                if (isClickForm) {
                    buttonElement.attr('disabled', false).text(submitButtonText);
                }
                isClickForm = false;
                postData = formElement.serializeObject();
                if (!firstLoad) {
                    postData.tab = currentTab
                    makeParamsUrl(postData);
                } else {
                    firstLoad = false;
                }
                if (json.data === null) {
                    return []
                }
                return json.data;
            },
            error: function (jqXHR, exception) {
                let msg = AjaxErrorMessage(jqXHR, exception)
                new AlertError(msg);
                if (isClickForm) {
                    buttonElement.attr('disabled', false).text(submitButtonText);
                }
            },
        },
        columns: [
            { data: "id", name: "ID" },
            { data: "name", name: "Name" },
            { data: "status", name: "Status" },
            { data: "type", name: "Type" },
            { data: "size", name: "Size" },
            { data: "action", name: "Action" },
        ],
        drawCallback: function () {
            $(".dataTables_paginate > ul.pagination > li > a.page-link").addClass("text-secondary");
            $(".dataTables_paginate > ul.pagination > li.active > a").addClass('bg-warning border-warning').css("color", "#111111");
            $(".dataTables_length > label > select").removeClass().addClass("form-select form-select-sm");
            if (!isClickForm) {
                $(".table-responsive").css("height", "");
            }
            const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
            const tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
                return new bootstrap.Tooltip(tooltipTriggerEl)
            });
        }
    }
    let pageLength;
    pageLength = parseFloat(postData.length);
    if (jQuery.inArray(pageLength, [10, 25, 50, 100]) >= 0) {
        setting.pageLength = pageLength;
    }
    setting.displayStart = parseFloat(postData.start);
    $(tableID).DataTable(setting);
}

function Delete(id) {
    let url = "/adtag/del"
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
            var msg = AjaxErrorMessage(jqXHR, exception);
            new AlertError("AJAX ERROR: " + msg);
        }
    }).done(function (result) {
        switch (result.status) {
            case "success":
                NoticeSuccess("Ad tag has been removed successfully")
                GetTable(true, false)
                break
            case "err":
                new AlertError(result.message);
        }
    });
}

function makeParamsUrl(obj) {
    let params = jQuery.param(obj).replaceAll("%5B%5D", "")
    // let params = jQuery.param(obj)
    let newUrl = currentURL + "&" + params
    window.history.pushState("object or string", "Title", newUrl);
    window.history.replaceState("object or string", "Title", newUrl);
}


// *******************************************  Ad Tag  *************************************************
$(document).ready(function () {
    // var BoxAdtag = $("#BoxAdtag")
    var BoxAdtag = $("#BoxAdtag")
    //Submit

    $("#btn-add-adtag").click(function () {
        $("#BoxAdtag").removeClass("close-sidebar");
        LoadCreateAdtag();
    })

    $(".tab-content").on("click", ".edit-adtag", function () {
        $("#BoxAdtag").removeClass("close-sidebar");
        LoadEditAdtag($(this));
    })

    // HandleCollapseAddAdtag()
    // $.fn.select2.defaults.set("theme", "bootstrap");
    // BoxAdtag.on("change", "#player_size", function () {
    //     if ($(this).prop("checked") === true) {
    //         $("#set_player_size").addClass("d-none");
    //     } else {
    //         $("#set_player_size").removeClass("d-none");
    //         $("#player_size_width").focus();
    //     }
    // });

    // Select Type
    BoxAdtag.on("change", "#ad_tag_type", function () {
        let adTagType = $(this).val();
        // if (adTagType === "3" || adTagType === "4" || adTagType === '5') {
        //     $("#gam_ad_display").addClass("d-none")
        //     $("#gam_ad_instream").addClass("d-none")
        // } else {
        //     $("#gam_ad_display").removeClass("d-none")
        //     $("#gam_ad_instream").removeClass("d-none")
        // }

        BoxAdtag.find(".box-config").attr("hidden", true)
        BoxAdtag.find('*[data-adtype="' + adTagType + '"]').attr("hidden", false)
        BoxAdtag.find('span.select2-container--bootstrap').css("width", "100%")

        switch (adTagType) {
            // case "5":
            //     BoxAdtag.find(".box_shift_content").removeClass("d-none");
            //     BoxAdtag.find(".ad_refresh").removeClass("d-none");
            //     break;
            case "5":
                BoxAdtag.find(".box_shift_content").removeClass("d-none");
                BoxAdtag.find(".ad_refresh").removeClass("d-none");
                break;

            default:
                BoxAdtag.find(".box_shift_content").addClass("d-none");
                BoxAdtag.find(".ad_refresh").addClass("d-none");
                break;
        }
        // inputCheckBox()
    });

    // $("#sticky_show_on_mobile").change(function () {
    //     if ($(this).prop('checked')) {
    //         $(".sticky-mobile-config").removeClass("d-none");
    //     } else {
    //         $(".sticky-mobile-config").addClass("d-none");
    //     }
    // });

    $("#BoxAdtag").on("click", ".toggle-collapse", function () {
        $(this).find('a[data-bs-toggle="collapse"]')[0].click()
    })
    $("#BoxAdtag").on("change", "#passback_type_outstream", function () {
        checkPassbackType()
    })
    $("#BoxAdtag").on("change", "#content_source", function () {
        SelectContentSource()
    });
    $("#BoxAdtag").on("change", "#content_source_articles", function () {
        SelectContentSource()
    });
    // $("#BoxAdtag").on("change", "#renderer_instream", function () {
    //     changeRenderer()
    // })
    // $("#BoxAdtag").on("change", "#renderer_outstream", function () {
    //     changeRenderer()
    // }) 

    $("#BoxAdtag").on("change", "#renderer_video", function () {
        changeRenderer()
    })
    $("#BoxAdtag").on("click", "#btn-close", function () {
        $("body").removeClass('customize-box')
    })
    $("body").on("click", ".customize-btn", function () {
        $("body").addClass('customize-box')
    })
    // SelectContentSource();
    // SelectPrimaryAdSize();
    // SelectSizeStickBanner();
    CheckFormError()
})

function checkPassbackType() {
    var PassbackType = parseInt($("#BoxAdtag").find("#passback_type_outstream").val())
    if (PassbackType == 1) {
        $("#BoxAdtag #inline_tag_outstream").closest(".collapsePassbackType").removeClass("d-none")
        $("#BoxAdtag #pass_back_outstream").closest(".collapsePassbackType").addClass("d-none")
    } else if (PassbackType == 2) {
        $("#BoxAdtag #inline_tag_outstream").closest(".collapsePassbackType").addClass("d-none")
        $("#BoxAdtag #inline_tag_outstream").selectpicker("refresh")
        $("#BoxAdtag #pass_back_outstream").closest(".collapsePassbackType").removeClass("d-none")
    } else {
        $("#BoxAdtag #inline_tag_outstream").closest(".collapsePassbackType").addClass("d-none")
        $("#BoxAdtag #inline_tag_outstream").selectpicker("refresh")
        $("#BoxAdtag #pass_back_outstream").closest(".collapsePassbackType").addClass("d-none")
    }
}

// **********************  Add Adtag  **************************

function LoadCreateAdtag() {
    var inventoryID = $("#btn-add-adtag").attr("data-id")
    if (!inventoryID) {
        return
    }
    // $("#BoxAdtag").find(".modal-body").append(Loading())
    $("#BoxAdtag").find("h3 span.sidebar-title").text("Create New Ad Tag")
    $("#BoxAdtag").find(".history-a").addClass("d-none")
    $("#BoxAdtag").append(Loading())

    $.ajax({
        type: 'GET',
        url: '/adtag-v2/add',
        data: { inventoryId: inventoryID }
    })
        .done(function (result) {
            if (result.error) {
                return;
            }
            $("#BoxAdtag").find(".result-adtag").html(result)
            $("#BoxAdtag").find("._blur").remove()
            $("#BoxAdtag").find(".selectpicker").selectpicker("refresh")

            // $.fn.select2.defaults.set("theme", "bootstrap");
            // $('.select2').select2({
            //     dropdownParent: $('#BoxAdtag')
            // });

            // HandleCollapseAddAdtag()
            SelectContentSource();
            SelectPrimaryAdSize();
            SelectAdSizeAdditional()
            // SelectSizeStickBanner();
            EventSelectPositionStickyBanner();
            selectAdSize();
            checkPassbackType();
            selectColorPlayZone();
            checkEnableStickyDesktopAndMobile();
            selectContentType();
            selectTemplatePlayZone();

            selectAdRefresh();
            SubmitFormCreate("AddAdTag", AddAdTag, "/adtag/add");

            eventSelectAdType();

            // inputCheckBox()
        })
}

function eventSelectAdType() {
    $("#ConfigVideoBox").on("change", "#template", function () {
        let type = $(this).find(":selected").data("type");

        switch (type) {
            case "Instream":
                $(".instream-show-wrapper").removeClass("d-none");
                $(".outstream-show-wrapper").addClass("d-none");
                break;
            case "Outstream":
                $(".instream-show-wrapper").addClass("d-none");
                $(".outstream-show-wrapper").removeClass("d-none");
                break;
            default:
                break;
        }
    });

}

function selectTemplatePlayZone() {
    $(".slide").on("click", function () {
        $(".slide").removeClass("selected");
        $(this).addClass("selected");
        let value = $(this).data("value");
        $(this).prevAll("#template_play_zone").attr("value", value);
    });
}

function selectContentType() {
    function changOptionTemplate() {
        if ($("#content_type").val() === "1") {
            $('#ConfigPlayZoneBox').find(".content_type_quiz").addClass("d-none");
            $('#ConfigPlayZoneBox').find(".content_type_related").removeClass("d-none");
            // $('#template_play_zone').val("1");
            $('#slider-container').find(".content_type_related").first().click();
        } else if ($("#content_type").val() === "2") {
            $('#ConfigPlayZoneBox').find(".content_type_related").addClass("d-none");
            $('#ConfigPlayZoneBox').find(".content_type_quiz").removeClass("d-none");
            // $('#template_play_zone').val("1");
            $('#slider-container').find(".content_type_quiz").first().click();
        }
    }

    changOptionTemplate();
    $("#content_type").on("change", function () {
        changOptionTemplate();
    });
}

function checkEnableStickyDesktopAndMobile(isEdit = false) {
    function checkEnableDesktop() {
        if ($("#BoxAdtag").find("#enable_sticky_desktop").is(':checked')) {
            $("#BoxAdtag").find(".sticky_desktop").prop('disabled', false);

            if ($("#BoxAdtag").find("#size_sticky").val() === "") {
                $("#BoxAdtag").find("#additional_ad_size_desktop_stick").prop('disabled', true);
            }
            // if (isEdit && $("#BoxAdtag").find("#size_sticky").val() !== "") {
            //     $("#BoxAdtag").find("#size_sticky").prop('disabled', true);
            // }

            $("#BoxAdtag").find("#enable_sticky_mobile").prop('disabled', false);
        } else {
            $("#BoxAdtag").find(".sticky_desktop").prop('disabled', true);

            $("#BoxAdtag").find("#enable_sticky_mobile").prop('disabled', true);
        }
        $(".selectpicker").selectpicker('refresh');
    }

    function checkEnableMobile() {
        if ($("#BoxAdtag").find("#enable_sticky_mobile").is(':checked')) {
            $("#BoxAdtag").find(".sticky_mobile").prop('disabled', false);

            if ($("#BoxAdtag").find("#size_sticky_mobile").val() === "") {
                $("#BoxAdtag").find("#additional_ad_size_mobile_stick").prop('disabled', true)
            }
            // if (isEdit && $("#BoxAdtag").find("#size_sticky_mobile").val() !== "") {
            //     $("#BoxAdtag").find("#size_sticky_mobile").prop('disabled', true);
            // }

            $("#BoxAdtag").find("#enable_sticky_desktop").prop('disabled', false);
        } else {
            $("#BoxAdtag").find(".sticky_mobile").prop('disabled', true);

            $("#BoxAdtag").find("#enable_sticky_desktop").prop('disabled', true);
        }
        $(".selectpicker").selectpicker('refresh');
    }

    checkEnableDesktop();
    checkEnableMobile();
    $("#BoxAdtag").find("#enable_sticky_desktop").on("change", function () {
        checkEnableDesktop();
    });

    $("#BoxAdtag").find("#enable_sticky_mobile").on("change", function () {
        checkEnableMobile();
    });
}

// position_sticky: 1 - bottom center , 2 - bottom left, 3 - bottom right, 4 - top, 5 - bottom, 6 - top center
function SelectSizeStickBanner() {
    // $("#BoxAdtag").on("change", "#size_sticky", function () {
    //     let value = $(this).val();
    //     if (value === "2" || value === "10") {
    //         $("#BoxAdtag").find('#position_sticky').empty()
    //             .append(appendOption(1)).append(appendOption(6)).selectpicker('refresh')
    //     } else if (value === "1" || value === "3" || value === "4") {
    //         $("#BoxAdtag").find('#position_sticky').empty()
    //             .append(appendOption(2)).append(appendOption(3)).selectpicker('refresh')
    //     } else if (value === "7" || value === "9" || value === "4") {
    //         $("#BoxAdtag").find('#position_sticky').empty()
    //             .append(appendOption(4)).append(appendOption(5)).selectpicker('refresh')
    //     }
    // });
    // $("#BoxAdtag").on("change", "#size_sticky_mobile", function () {
    //     let value = $(this).val();
    //     if (value === "2" || value === "10") {
    //         $("#BoxAdtag").find('#position_sticky_mobile').empty()
    //             .append(appendOption(1)).append(appendOption(6)).selectpicker('refresh')
    //     } else if (value === "1" || value === "3" || value === "4") {
    //         $("#BoxAdtag").find('#position_sticky_mobile').empty()
    //             .append(appendOption(2)).append(appendOption(3)).selectpicker('refresh')
    //     } else if (value === "7" || value === "9" || value === "4") {
    //         $("#BoxAdtag").find('#position_sticky_mobile').empty()
    //             .append(appendOption(4)).append(appendOption(5)).selectpicker('refresh')
    //     }
    // });
}

function EventSelectPositionStickyBanner() {
    $("#BoxAdtag").on("change", "#position_sticky", function () {
        let value = $(this).val();
        switch (value) {
            //Bottom Center
            case "1":
            //Top Center
            case "6":
                $("#BoxAdtag").find('#size_sticky').empty()
                    .append(appendOptionAdSize(2)).append(appendOptionAdSize(10)).selectpicker('refresh')
                break;

            //Bottom Left
            case "2":
            //Bottom Right
            case "3":
                $("#BoxAdtag").find('#size_sticky').empty()
                    .append(appendOptionAdSize(1)).append(appendOptionAdSize(3)).append(appendOptionAdSize(4)).selectpicker('refresh')
                break;

            default:
                break;
        }
    });

    $("#BoxAdtag").on("change", "#position_sticky_mobile", function () {
        let value = $(this).val();
        switch (value) {
            //Top
            case "4":
            //Bottom
            case "5":

                $("#BoxAdtag").find('#size_sticky_mobile').empty()
                    .append(appendOptionAdSize(7)).append(appendOptionAdSize(9)).
                    append(appendOptionAdSize(19)).append(appendOptionAdSize(20)).selectpicker('refresh')
                break;


            default:
                break;
        }
    });
}

function selectAdSize() {
    if ($("#BoxAdtag #ad_size").val() === "1") {
        $("#BoxAdtag .ad_size_responsive").addClass("d-none");
        $("#BoxAdtag .ad_size_fixed").removeClass("d-none");
    } else if ($("#BoxAdtag #ad_size").val() === "2") {
        $("#BoxAdtag .ad_size_fixed").addClass("d-none");
        $("#BoxAdtag .ad_size_responsive").removeClass("d-none");
    }
    $("#BoxAdtag").on("change", "#ad_size", function () {
        if ($(this).val() === "1") {
            $(".ad_size_responsive").addClass("d-none");
            $(".ad_size_fixed").removeClass("d-none");
        } else if ($(this).val() === "2") {
            $(".ad_size_fixed").addClass("d-none");
            $(".ad_size_responsive").removeClass("d-none");
        }
    });
}

function appendOption(position) {
    let value, name;
    switch (position) {
        case 1:
            value = 1;
            name = "Bottom Center";
            break;
        case 2:
            value = 2;
            name = "Bottom Left";
            break;
        case 3:
            value = 3;
            name = "Bottom Right";
            break;
        case 4:
            value = 4;
            name = "Top";
            break;
        case 5:
            value = 5;
            name = "Bottom";
            break;
        case 6:
            value = 6;
            name = "Top Center";
            break;
    }
    return `<option value="${value}">${name}</option>`
}

function appendOptionAdSize(position) {
    let value, name;
    switch (position) {
        case 1:
            value = 1;
            name = "300x250";
            break;
        case 2:
            value = 2;
            name = "728x90";
            break;
        case 3:
            value = 3;
            name = "160x600";
            break;
        case 4:
            value = 4;
            name = "300x600";
            break;
        case 7:
            value = 7;
            name = "320x50";
            break;
        case 9:
            value = 9;
            name = "320x100";
            break;
        case 10:
            value = 10;
            name = "970x90";
            break;
        case 19:
            value = 19;
            name = "300x100";
            break;
        case 20:
            value = 20;
            name = "300x50";
            break;
    }
    return `<option selected value="${value}">${name}</option>`
}

function inputCheckBox() {
    // if ($("#BoxAdtag").find("#close_button_sticky").is(':checked')) {
    //     $("#BoxAdtag").find("#close_button_sticky").bootstrapToggle('on');
    // } else {
    //     $("#BoxAdtag").find("#close_button_sticky").bootstrapToggle('off');
    // }
    // if ($("#BoxAdtag").find("#close_button_sticky_mobile").is(':checked')) {
    //     $("#BoxAdtag").find("#close_button_sticky_mobile").bootstrapToggle('on');
    // } else {
    //     $("#BoxAdtag").find("#close_button_sticky_mobile").bootstrapToggle('off');
    // }
    $("#BoxAdtag").find('input[type="checkbox"]').each(function () {
        if ($(this).is(':checked')) {
            $(this).bootstrapToggle('on');
        } else {
            $(this).bootstrapToggle('off');
        }
    })
}

function SelectPrimaryAdSize() {
    if ($("#BoxAdtag").find("#bid_out_stream").is(':checked')) {
        $("#BoxAdtag").find("#bid_out_stream").bootstrapToggle('on');
    } else {
        $("#BoxAdtag").find("#bid_out_stream").bootstrapToggle('off');
    }
    $("#BoxAdtag").find("#primary_ad_size").on("change", function (e) {
        // data = e.params.data;
        // let value = $("#BoxAdtag").find("#primary_ad_size").val();
        // $("#BoxAdtag").find('#additional_ad_size').selectpicker('refresh')
        // $("#BoxAdtag").find('#additional_ad_size_mobile').selectpicker('refresh')
        // getSizeAdditional(parseInt(value))

        let size = $(this).find("option:selected").text();
        let listSizeEnableBidOutStream = ["300x250", "336x280", "300x600", "970x250"]
        let checkSizeEnableBidOutStream = $.inArray(size.trim(), listSizeEnableBidOutStream)
        if (checkSizeEnableBidOutStream > -1) {
            $("#BoxAdtag").find(".box_bid_out_stream").removeClass('d-none')
            $("#BoxAdtag").find("#bid_out_stream").bootstrapToggle('on');
        } else {
            $("#BoxAdtag").find(".box_bid_out_stream").addClass('d-none')
            $("#BoxAdtag").find("#bid_out_stream").bootstrapToggle('off');
        }

        // khi Primary Ad Size có size > 320px thêm select box Size On Mobile: 300x250, 320x50x, 320x100, 300x100, 300x50
        if (size) {
            var arraySize = size.split("x");
            if (parseInt(arraySize[0]) > 320) {
                $("#BoxAdtag").find(".mobile-config").attr("hidden", false);
            } else {
                $("#BoxAdtag").find(".mobile-config").attr("hidden", true);
                $("#BoxAdtag").find('#size_on_mobile').selectpicker('val', '')
                $("#BoxAdtag").find('#additional_ad_size_mobile').selectpicker('val', '')
                $("#BoxAdtag").find('#pass_back_mobile').val("")
                // $("#BoxAdtag").find("#size_on_mobile").val("").trigger('change');
            }
        }
    });
}

function SelectAdSizeAdditional() {
    $("#BoxAdtag").find("#primary_ad_size").on("change", function (e) {
        let value = $("#BoxAdtag").find("#primary_ad_size").val();
        $("#BoxAdtag").find('#additional_ad_size').selectpicker('refresh')
        getAdditionalAdSizes(parseInt(value), "#additional_ad_size")
    })
    $("#BoxAdtag").find("#size_on_mobile").on("change", function (e) {
        let value = $("#BoxAdtag").find("#size_on_mobile").val();
        $("#BoxAdtag").find('#additional_ad_size_mobile').selectpicker('refresh')
        getAdditionalAdSizes(parseInt(value), "#additional_ad_size_mobile")
    })
    // $("#BoxAdtag").find("#size_sticky").on("change", function (e) {
    //     let value = $("#BoxAdtag").find("#size_sticky").val();
    //     $("#BoxAdtag").find('#additional_ad_size_desktop_stick').selectpicker('refresh')
    //     getAdditionalAdSizes(parseInt(value), "#additional_ad_size_desktop_stick")
    // });
    // $("#BoxAdtag").find("#size_sticky_mobile").on("change", function (e) {
    //     let value = $("#BoxAdtag").find("#size_sticky_mobile").val();
    //     $("#BoxAdtag").find('#additional_ad_size_mobile_stick').selectpicker('refresh')
    //     getAdditionalAdSizes(parseInt(value), "#additional_ad_size_mobile_stick", "sticky_mobile")
    // });
}

function SelectContentSource() {
    $("#BoxAdtag").find(".content_source_sub").attr("hidden", true);
    switch ($("#BoxAdtag").find("#content_source").val()) {
        case "1":
            $("#BoxAdtag").find("div.content_source_sub[data-content_source='1']").attr("hidden", false);
            break;
        case "2":
            $("#BoxAdtag").find("div.content_source_sub[data-content_source='2']").attr("hidden", false);
            break;
        default:
            break;
    }
    $("#BoxAdtag").find(".content_source_articles_sub").attr("hidden", true);
    let contentSourceArticles = $("#BoxAdtag").find("#content_source_articles").val();
    if (contentSourceArticles) {
        $("#BoxAdtag").find("div.content_source_articles_sub[data-content_source_articles=" + contentSourceArticles + "]").attr("hidden", false);
    }
}

function getSizeAdditional(id) {
    if (!id) {
        return
    }
    $.ajax({
        url: "/adtag/getSizeAdditional",
        type: "GET",
        data: {
            id: id,
        },
        success: function (data) {
            $("#BoxAdtag").find('#additional_ad_size').html("").prop("disabled", false)
            $("#BoxAdtag").find('#additional_ad_size_mobile').html("").prop("disabled", false)
            $.each(data, function (index, item) {
                const newOption = new Option(item.name, item.id, false, false);
                $("#BoxAdtag").find('#additional_ad_size').append(newOption).selectpicker('refresh');
            });
            $.each(data, function (index, item) {
                const newOption = new Option(item.name, item.id, false, false);
                $("#BoxAdtag").find('#additional_ad_size_mobile').append(newOption).selectpicker('refresh');
            });
        },
        error: function (xhr) {
            console.log(xhr)
        }
    })
}


function getAdditionalAdSizes(id, element, type) {
    if (!id) {
        return
    }
    $.ajax({
        url: "/adtag/getSizeAdditional",
        type: "GET",
        data: {
            id: id,
            type: type,
        },
        success: function (data) {
            $("#BoxAdtag").find(element).html("").prop("disabled", false)
            $.each(data, function (index, item) {
                const newOption = new Option(item.name, item.id, false, false);
                $("#BoxAdtag").find(element).append(newOption).selectpicker('refresh');
            });
        },
        error: function (xhr) {
            console.log(xhr)
        }
    })
}


function SubmitFormCreate(formID, functionCallback, ajaxUrl = "") {
    let formElement = $("#BoxAdtag").find("#" + formID);
    formElement.find(".selectpicker").on("changed.bs.select", function (e, clickedIndex, newValue, oldValue) {
        $(this).closest(".box-selectpicker").find(".bs-placeholder").attr("style", "border-color:#a0acc2 !important");
        $(this).closest(".box-selectpicker").removeClass("is-invalid");
    });
    formElement.find(".submit-adtag").on("click", function (e) {
        e.preventDefault();
        const buttonElement = $(this);
        const submitButtonText = buttonElement.text();
        const submitButtonTextLoading = "Loading...";
        var postData = formElement.serializeObject();
        let additional_ad_size = []
        if (postData.additional_ad_size && postData.additional_ad_size.constructor === Array) {
            $.each(postData.additional_ad_size, function (index, value) {
                additional_ad_size.push(parseInt(value))
            })
        } else {
            additional_ad_size.push(parseInt(postData.additional_ad_size))
        }
        postData.additional_ad_size = additional_ad_size

        let additional_ad_size_mobile = []
        if (postData.additional_ad_size_mobile && postData.additional_ad_size_mobile.constructor === Array) {
            $.each(postData.additional_ad_size_mobile, function (index, value) {
                additional_ad_size_mobile.push(parseInt(value))
            })
        } else {
            additional_ad_size_mobile.push(parseInt(postData.additional_ad_size_mobile))
        }
        postData.additional_ad_size_mobile = additional_ad_size_mobile

        let additional_ad_size_mobile_stick = []
        if (postData.size_sticky_mobile && postData.size_sticky_mobile.constructor === Array) {
            $.each(postData.size_sticky_mobile, function (index, value) {
                additional_ad_size_mobile_stick.push(parseInt(value))
            })
        } else {
            additional_ad_size_mobile_stick.push(parseInt(postData.size_sticky_mobile))
        }
        postData.additional_ad_size_mobile_stick = additional_ad_size_mobile_stick

        let additional_ad_size_desktop_stick = []
        if (postData.size_sticky && postData.size_sticky.constructor === Array) {
            $.each(postData.size_sticky, function (index, value) {
                additional_ad_size_desktop_stick.push(parseInt(value))
            })
        } else {
            additional_ad_size_desktop_stick.push(parseInt(postData.size_sticky))
        }
        postData.additional_ad_size_desktop_stick = additional_ad_size_desktop_stick

        postData.ad_tag_type = parseInt(postData.ad_tag_type);

        postData.content_source_articles = parseInt(postData.content_source_articles);
        postData.content_type = parseInt(postData.content_type);
        postData.content_source = parseInt(postData.content_source);
        postData.playlist = parseInt(postData.playlist);
        if (!postData.primary_ad_size) {
            postData.primary_ad_size = 0;
        } else {
            postData.primary_ad_size = parseInt(postData.primary_ad_size);
        }
        if (!postData.size_on_mobile) {
            postData.size_on_mobile = 0;
        } else {
            postData.size_on_mobile = parseInt(postData.size_on_mobile);
        }
        postData.passback_type_outstream = parseInt(postData.passback_type_outstream);
        postData.inline_tag_outstream = parseInt(postData.inline_tag_outstream);
        postData.template = parseInt(postData.template);
        postData.template_articles = parseInt(postData.template_articles);
        postData.template_outstream = parseInt(postData.template_outstream);
        postData.inventory_id = parseInt(postData.inventory_id);
        postData.position_sticky = parseInt(postData.position_sticky);
        postData.position_sticky_mobile = parseInt(postData.position_sticky_mobile);
        postData.size_sticky = parseInt(postData.size_sticky);
        postData.size_sticky_mobile = parseInt(postData.size_sticky_mobile);
        postData.renderer_instream = parseInt(postData.renderer_instream);
        postData.renderer_outstream = parseInt(postData.renderer_outstream);
        postData.renderer_video = parseInt(postData.renderer_video);
        postData.total_ads = parseInt(postData.total_ads);
        postData.template_play_zone = parseInt(postData.template_play_zone);
        postData.content_source_play_zone = parseInt(postData.content_source_play_zone);
        postData.ad_size = parseInt(postData.ad_size);
        postData.responsive_type = parseInt(postData.responsive_type);
        postData.template_native = parseInt(postData.template_native);
        postData.ad_refresh_time = parseInt(postData.ad_refresh_time);
        if (postData.close_button_sticky === "on") {
            postData.close_button_sticky = 1
        } else {
            postData.close_button_sticky = 2
        }
        if (postData.close_button_sticky_mobile) {
            if (postData.close_button_sticky_mobile === "on") {
                postData.close_button_sticky_mobile = 1
            } else {
                postData.close_button_sticky_mobile = 2
            }
        }
        if (postData.bid_out_stream) {
            if (postData.bid_out_stream === "on") {
                postData.bid_out_stream = 1
            } else {
                postData.bid_out_stream = 2
            }
        }
        if (postData.status === "on" || postData.status === "ON" || postData.status === "running") {
            postData.status = 1
        } else {
            postData.status = 2
        }
        if (postData.shift_content === "on") {
            postData.shift_content = 1
        } else {
            postData.shift_content = 2
        }
        if ($("#enable_sticky_desktop").is(":checked")) {
            postData.enable_sticky_desktop = 1
        } else {
            postData.enable_sticky_desktop = 2
        }
        if ($("#enable_sticky_mobile").is(":checked")) {
            postData.enable_sticky_mobile = 1
        } else {
            postData.enable_sticky_mobile = 2
        }
        if (postData.banner_ad === "on") {
            postData.banner_ad = 1
        } else {
            postData.banner_ad = 2
        }
        if (postData.video_ad === "on") {
            postData.video_ad = 1
        } else {
            postData.video_ad = 2
        }

        if (postData.ad_tag_type === 1 || postData.ad_tag_type === 5) {
            if (postData.ad_refresh === "domain_configuration") {
                delete postData.ad_refresh_time;
            }
        } else {
            delete postData.ad_refresh;
            delete postData.ad_refresh_time;
        }
        $.ajax({
            url: ajaxUrl,
            type: "POST",
            dataType: "JSON",
            contentType: "application/json",
            data: JSON.stringify(postData),
            beforeSend: function (xhr) {
                buttonElement.prop('disabled', true).text(submitButtonTextLoading);
            },
            error: function (jqXHR, exception) {
                const msg = AjaxErrorMessage(jqXHR, exception);
                new AlertError("AJAX ERROR: " + msg);
                buttonElement.prop('disabled', false).text(submitButtonText);
            },
            success: function (responseJSON) {
                buttonElement.prop('disabled', false).text(submitButtonText);
            },
            complete: function (res) {
                functionCallback(res.responseJSON, formElement);
            }
        });
    });
}

function AddAdTag(response, formElement) {
    switch (response.status) {
        case "error":
            if (response.errors.length === 1 && response.errors[0].id === "") {
                new AlertError(response.errors[0].message);
            } else if (response.errors.length > 0) {
                $.each(response.errors, function (key, value) {
                    let inputElement = $("#BoxAdtag").find("#" + value.id);
                    inputElement.addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
                    inputElement.closest(".box-selectpicker").find(".bs-placeholder").attr("style", "border-color:#e35d6a!important")
                    inputElement.closest(".box-selectpicker").addClass("is-invalid").nextAll("span.invalid-feedback").html(value.message);
                    // if (value.id === "box_primary_ad_size") {
                    // inputElement.find(".select2-container--bootstrap .select2-selection--single").css("border-color", "#e35d6a")
                    // }

                });
                $("#BoxAdtag").find("#" + response.errors[0].id).focus();
                // new AlertError(response.errors[0].message, function () {
                //     $("#" + response.errors[0].id).focus();
                //     $("#" + response.errors[0].id).prev('label').focus();
                // })
            } else {
                new AlertError("Error!");
            }
            break
        case "success":
            let id = response.data_object.inventory_id
            NoticeSuccess("Ad tag has been created successfully");
            GetTable(false, true);
            $("#BoxAdtag").find('[data-bs-dismiss=modal]').click();
            $("#BoxAdtag").addClass('close-sidebar');
            LoadAdTag($("#inventoryTabs").attr('data-id'))
            setTimeout(function () {
                $("body").removeClass("customize-box");
            }, 400)
            break
        default:
            new AlertError("Undefined");
            break
    }
}

// function HandleCollapseAddAdtag() {
//     let url = "/adtag/collapse"
//     const urlSearchParams = new URLSearchParams(window.location.search);
//     const params = Object.fromEntries(urlSearchParams.entries());
//     $("#BoxAdtag").find('.collapse').on('show.bs.collapse', function (e) {
//         let box = e.target.id
//         SendRequestShowAdtag(url, box, 0, "add")
//     });
//     $("#BoxAdtag").find('.collapse').on('hide.bs.collapse', function (e) {
//         let box = e.target.id
//         SendRequestHideAdtag(url, box, 0, "add")
//     });
// }

function SendRequestHideAdtag(url, box, id, type) {
    $.ajax({
        url: url,
        type: "POST",
        dataType: "JSON",
        contentType: "application/json",
        data: JSON.stringify({
            is_collapse: 1,
            box_collapse: box,
            page_type: type,
            page_id: id
        }),
        beforeSend: function (xhr) {
            xhr.overrideMimeType("text/plain; charset=x-user-defined");
        },
        error: function (jqXHR, exception) {
        }
    }).done(function (result) {
    });
}

function SendRequestShowAdtag(url, box, id, type) {
    $.ajax({
        url: url,
        type: "POST",
        dataType: "JSON",
        contentType: "application/json",
        data: JSON.stringify({
            is_collapse: 2,
            box_collapse: box,
            page_type: type,
            page_id: id
        }),
        beforeSend: function (xhr) {
            xhr.overrideMimeType("text/plain; charset=x-user-defined");
        },
        error: function (jqXHR, exception) {
            // var msg = AjaxErrorMessage(jqXHR, exception);
            // new AlertError("AJAX ERROR: " + msg);
        }
    }).done(function (result) {
        // switch (result.status) {
        //     case "error":
        //         // new AlertError(result.errors[0].message);
        //         console.log("error")
        //         break
        //     case "success":
        //         // new AlertSuccess("success",);
        //         console.log("success")
        //         break
        //     default:
        //         // new AlertError("Undefined");
        //         console.log("Undefined")
        //         break
        // }
    });
}

// khi change val - input, select, textarea sẽ bỏ phần hiển thị lỗi
function CheckFormError() {
    var Modal = $("#BoxAdtag")
    Modal.on("change", "input", function (e) {
        let inputElement = $(this)
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    Modal.on("change", "textarea", function (e) {
        let inputElement = $(this)
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    Modal.on("change", "select", function (e) {
        let inputElement = $(this)
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    Modal.on("change", "select2", function (e) {
        let inputElement = $(this)
        inputElement.parent().next("span.invalid-feedback").empty();
        inputElement.parent().find(".select2-container--bootstrap .select2-selection--single").css("border-color", "#ebedf2");
    });
    Modal.on("change", ".selectpicker", function (e) {
        let inputElement = $(this)
        inputElement.closest(".is-invalid").next("span.invalid-feedback").empty();
        inputElement.parent().find(".dropdown-toggle").css("border-color", "#ebedf2");
    });
}

// **********************  Edit Adtag  **************************

function LoadEditAdtag(e) {
    var id = e.attr("data-id")
    var inventoryID = e.attr("data-inventory")
    if (!inventoryID || !id) {
        return
    }

    $("#BoxAdtag").append(Loading())
    $("#BoxAdtag").find("h3 span.sidebar-title").text("Edit Ad Tag")
    $("#BoxAdtag").find(".load-history").removeClass("d-none").attr("data-id", id)

    $.ajax({
        type: 'GET',
        url: '/adtag-v2/edit',
        data: {
            id: id,
            inventoryId: inventoryID
        }
    })
        .done(function (result) {
            if (result.error) {
                return;
            }
            $("#BoxAdtag").find(".result-adtag").html(result)
            $("#BoxAdtag").find("._blur").remove()
            $("#BoxAdtag").find(".selectpicker").selectpicker("refresh")

            HandleCollapseEditAdtag();
            SelectContentSource();
            SelectPrimaryAdSize();
            SelectAdSizeAdditional();
            // SelectSizeStickBanner();
            EventSelectPositionStickyBanner();
            selectAdSize();
            // inputCheckBox()
            checkPassbackType();
            changeRenderer();
            selectColorPlayZone();
            checkEnableStickyDesktopAndMobile(true);
            selectContentType();
            selectTemplatePlayZone();
            selectAdRefresh();
            SubmitFormEdit("EditAdTag", EditAdTag, "/adtag/edit");
        })
}

function SubmitFormEdit(formID, functionCallback, ajxURL = "") {
    let formElement = $("#BoxAdtag").find("#" + formID);
    formElement.find(".selectpicker").on("changed.bs.select", function (e, clickedIndex, newValue, oldValue) {
        $(this).closest(".box-selectpicker").find(".bs-placeholder").attr("style", "border-color:#a0acc2 !important");
        $(this).closest(".box-selectpicker").removeClass("is-invalid");
    });
    formElement.find(".submit").on("click", function (e) {
        e.preventDefault();
        const buttonElement = $(this);
        const submitButtonText = buttonElement.text();
        const submitButtonTextLoading = "Loading...";
        var postData = formElement.serializeObject();
        let additional_ad_size = []
        if (postData.additional_ad_size && postData.additional_ad_size.constructor === Array) {
            $.each(postData.additional_ad_size, function (index, value) {
                additional_ad_size.push(parseInt(value))
            })
        } else {
            additional_ad_size.push(parseInt(postData.additional_ad_size))
        }
        postData.additional_ad_size = additional_ad_size

        let additional_ad_size_mobile = []
        if (postData.additional_ad_size_mobile && postData.additional_ad_size_mobile.constructor === Array) {
            $.each(postData.additional_ad_size_mobile, function (index, value) {
                additional_ad_size_mobile.push(parseInt(value))
            })
        } else {
            additional_ad_size_mobile.push(parseInt(postData.additional_ad_size_mobile))
        }
        postData.additional_ad_size_mobile = additional_ad_size_mobile

        let additional_ad_size_mobile_stick = []
        if (postData.size_sticky_mobile && postData.size_sticky_mobile.constructor === Array) {
            $.each(postData.size_sticky_mobile, function (index, value) {
                additional_ad_size_mobile_stick.push(parseInt(value))
            })
        } else {
            additional_ad_size_mobile_stick.push(parseInt(postData.size_sticky_mobile))
        }
        postData.additional_ad_size_mobile_stick = additional_ad_size_mobile_stick

        let additional_ad_size_desktop_stick = []
        if (postData.size_sticky && postData.size_sticky.constructor === Array) {
            $.each(postData.size_sticky, function (index, value) {
                additional_ad_size_desktop_stick.push(parseInt(value))
            })
        } else {
            additional_ad_size_desktop_stick.push(parseInt(postData.size_sticky))
        }
        postData.additional_ad_size_desktop_stick = additional_ad_size_desktop_stick

        postData.id = parseInt(postData.id)
        postData.ad_tag_type = parseInt($('#ad_tag_type').val())
        postData.content_source_articles = parseInt(postData.content_source_articles)
        postData.content_type = parseInt(postData.content_type)
        postData.content_source = parseInt(postData.content_source)
        postData.playlist = parseInt(postData.playlist)
        if (!postData.primary_ad_size) {
            postData.primary_ad_size = 0;
        } else {
            postData.primary_ad_size = parseInt(postData.primary_ad_size)
        }
        if (!postData.size_on_mobile) {
            postData.size_on_mobile = 0;
        } else {
            postData.size_on_mobile = parseInt(postData.size_on_mobile)
        }
        postData.passback_type_outstream = parseInt(postData.passback_type_outstream)
        postData.inline_tag_outstream = parseInt(postData.inline_tag_outstream)
        postData.template = parseInt(postData.template)
        postData.template_articles = parseInt(postData.template_articles)
        postData.template_outstream = parseInt(postData.template_outstream)
        postData.inventory_id = parseInt(postData.inventory_id)
        postData.position_sticky = parseInt(postData.position_sticky)
        postData.position_sticky_mobile = parseInt(postData.position_sticky_mobile)
        postData.size_sticky = parseInt(postData.size_sticky)
        postData.size_sticky_mobile = parseInt(postData.size_sticky_mobile)
        postData.renderer_instream = parseInt(postData.renderer_instream)
        postData.renderer_outstream = parseInt(postData.renderer_outstream)
        postData.template_play_zone = parseInt(postData.template_play_zone);
        postData.content_source_play_zone = parseInt(postData.content_source_play_zone);
        postData.total_ads = parseInt(postData.total_ads);
        postData.ad_size = parseInt(postData.ad_size);
        postData.responsive_type = parseInt(postData.responsive_type);
        postData.template_native = parseInt(postData.template_native);
        postData.ad_refresh_time = parseInt(postData.ad_refresh_time);
        if (postData.close_button_sticky === "on") {
            postData.close_button_sticky = 1
        } else {
            postData.close_button_sticky = 2
        }
        if (postData.close_button_sticky_mobile) {
            if (postData.close_button_sticky_mobile === "on") {
                postData.close_button_sticky_mobile = 1
            } else {
                postData.close_button_sticky_mobile = 2
            }
        }
        if (postData.bid_out_stream) {
            if (postData.bid_out_stream === "on") {
                postData.bid_out_stream = 1
            } else {
                postData.bid_out_stream = 2
            }
        }
        if (postData.status === "on" || postData.status === "ON" || postData.status === "running") {
            postData.status = 1
        } else {
            postData.status = 2
        }
        if (postData.shift_content === "on") {
            postData.shift_content = 1
        } else {
            postData.shift_content = 2
        }
        if ($("#enable_sticky_desktop").is(":checked")) {
            postData.enable_sticky_desktop = 1
        } else {
            postData.enable_sticky_desktop = 2
        }
        if ($("#enable_sticky_mobile").is(":checked")) {
            postData.enable_sticky_mobile = 1
        } else {
            postData.enable_sticky_mobile = 2
        }
        if (postData.banner_ad === "on") {
            postData.banner_ad = 1
        } else {
            postData.banner_ad = 2
        }
        if (postData.video_ad === "on") {
            postData.video_ad = 1
        } else {
            postData.video_ad = 2
        }

        if (postData.ad_tag_type === 1 || postData.ad_tag_type === 5) {
            if (postData.ad_refresh === "domain_configuration") {
                delete postData.ad_refresh_time;
            }
        } else {
            delete postData.ad_refresh;
            delete postData.ad_refresh_time;
        }
        $.ajax({
            url: ajxURL,
            type: "POST",
            dataType: "JSON",
            contentType: "application/json",
            data: JSON.stringify(postData),
            beforeSend: function (xhr) {
                buttonElement.prop('disabled', true).text(submitButtonTextLoading);
            },
            error: function (jqXHR, exception) {
                const msg = AjaxErrorMessage(jqXHR, exception);
                new AlertError("AJAX ERROR: " + msg);
                buttonElement.prop('disabled', false).text(submitButtonText);
            },
            success: function (responseJSON) {
                buttonElement.prop('disabled', false).text(submitButtonText);
            },
            complete: function (res) {
                functionCallback(res.responseJSON, formElement);
            }
        });
    });
}

function HandleCollapseEditAdtag() {
    let url = "/adtag/collapse"
    const urlSearchParams = new URLSearchParams(window.location.search);
    const params = Object.fromEntries(urlSearchParams.entries());
    $("#BoxAdtag").find('.collapse').on('show.bs.collapse', function (e) {
        let box = e.target.id
        SendRequestShowAdtag(url, box, parseInt(params.id), "edit")
    });
    $("#BoxAdtag").find('.collapse').on('hide.bs.collapse', function (e) {
        let box = e.target.id
        SendRequestHideAdtag(url, box, parseInt(params.id), "edit")
    });
}

function EditAdTag(response, formElement) {
    switch (response.status) {
        case "error":
            if (response.errors.length === 1 && response.errors[0].id === "") {
                new AlertError(response.errors[0].message);
            } else if (response.errors.length > 0) {
                $.each(response.errors, function (key, value) {
                    let inputElement = $("#BoxAdtag").find("#" + value.id);
                    inputElement.addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
                    inputElement.closest(".box-selectpicker").find(".bs-placeholder").attr("style", "border-color:#e35d6a!important")
                    inputElement.closest(".box-selectpicker").addClass("is-invalid").nextAll("span.invalid-feedback").html(value.message);
                    // if (value.id === "box_primary_ad_size") {
                    //     inputElement.find(".select2-container--bootstrap .select2-selection--single").css("border-color", "#e35d6a")
                    // }
                });
                $("#BoxAdtag").find("#" + response.errors[0].id).focus();
                // new AlertError(response.errors[0].message, function () {
                //     $("#" + response.errors[0].id).focus();
                //     $("#" + response.errors[0].id).prev('label').focus();
                // })
            } else {
                new AlertError("Error!");
            }
            break
        case "success":
            let id = response.data_object.inventory_id
            NoticeSuccess("Ad tag has been updated successfully")
            GetTable(true);
            // $("#BoxAdtag").find('[data-bs-dismiss=modal]').click()
            LoadAdTag($("#inventoryTabs").attr('data-id'))
            break
        default:
            new AlertError("Undefined");
            break
    }
}

function changeRenderer() {
    // var renderer = $("#BoxAdtag").find("#renderer_instream").val()
    // if (parseInt(renderer) !== 1) {
    //     $("#BoxAdtag").find("#template").closest("div.sidebar-content").addClass("d-none");
    //     $("#BoxAdtag").find("#content_source").closest("div.sidebar-content").addClass("d-none");
    //     $("#BoxAdtag").find("#output").closest("div.sidebar-content").removeClass("d-none");
    //     $("#BoxAdtag").find(".content_source_sub").attr("hidden", true);
    //     if (parseInt(renderer) === 6) {
    //         $(".renderer_instream_overlay_ad").removeClass("d-none");
    //     } else {
    //         $(".renderer_instream_overlay_ad").addClass("d-none");
    //     }
    // } else {
    //     $("#BoxAdtag").find("#output").closest("div.sidebar-content").addClass("d-none");
    //     $("#BoxAdtag").find("#template").closest("div.sidebar-content").removeClass("d-none");
    //     $("#BoxAdtag").find("#content_source").closest("div.sidebar-content").removeClass("d-none");
    //     $(".renderer_instream_overlay_ad").addClass("d-none");
    //     SelectContentSource()
    // }

    // var renderer_outstream = $("#BoxAdtag").find("#renderer_outstream").val()
    // if (parseInt(renderer_outstream) !== 1) {
    //     $("#BoxAdtag").find("#template_outstream").closest("div.sidebar-content").addClass("d-none");
    //     $("#BoxAdtag").find("#passback_type_outstream").closest("div.sidebar-content").addClass("d-none");
    //     $("#BoxAdtag").find(".collapsePassbackType").addClass("d-none");
    // } else {
    //     $("#BoxAdtag").find("#template_outstream").closest("div.sidebar-content").removeClass("d-none");
    //     $("#BoxAdtag").find("#passback_type_outstream").closest("div.sidebar-content").removeClass("d-none");
    //     checkPassbackType();
    // }

    var renderer = $("#BoxAdtag").find("#renderer_video").val()

    if (parseInt(renderer) !== 1) {
        $("#BoxAdtag").find("#template").closest("div.sidebar-content").addClass("d-none");
        $("#BoxAdtag").find("#content_source").closest("div.sidebar-content").addClass("d-none");
        // $("#BoxAdtag").find("#output").closest("div.sidebar-content").removeClass("d-none");
        $("#BoxAdtag").find(".instream-show-wrapper").addClass("d-none");
        $("#BoxAdtag").find(".outstream-show-wrapper").addClass("d-none");

        if (parseInt(renderer) === 6) {
            $(".renderer_instream_overlay_ad").removeClass("d-none");
        } else {
            $(".renderer_instream_overlay_ad").addClass("d-none");
        }
    } else {
        // $("#BoxAdtag").find("#output").closest("div.sidebar-content").addClass("d-none");
        $("#BoxAdtag").find("#template").closest("div.sidebar-content").removeClass("d-none");
        $(".renderer_instream_overlay_ad").addClass("d-none");

        let templateNow = $("#ConfigVideoBox").find("#template").find(":selected").data("type");
        if (templateNow == "Outstream") {
            $("#BoxAdtag").find(".outstream-show-wrapper").removeClass("d-none");
            checkPassbackType();
        }

        if (templateNow == "Instream") {
            $("#BoxAdtag").find(".instream-show-wrapper").removeClass("d-none");
            SelectContentSource()
        }
    }
}

// ************************** Ads.txt ************************** //
$(document).ready(function () {
    idInventory = $("#inventory_id").attr("data-inventory-id")
    $("#resultLoadMissingLine").on("click", ".copy-line", function () {
        let copyType = $(this).data("type");
        Copy(copyType);
    });

    SubmitForm("scanAds", function (resp) {
        Load();
        if (resp.status === "success") {
            NoticeSuccess(resp.message)
        } else {
            new NoticeError(resp.message);
        }
    }, "/ads_txt/scan");

    SubmitForm("SaveAdsTxt", SaveAdsTxt, "/ads_txt/detail?did=" + idInventory);
})

function SaveAdsTxt(response) {
    switch (response.status) {
        case "error":
            new AlertError(response.message)
            break
        case "success":
            NoticeSuccess(response.message)
            $("#ads_txt").val(response.data_object);
            break
        default:
            new AlertError("Undefined");
            break
    }
    Load();
}

function Load() {
    let postData = {
        did: $("#inventoryTabs").attr('data-id')
    }
    console.log(postData);
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
            console.log(res);
            $("#resultLoadMissingLine").html(res.responseJSON.data_object.html);
            if ($("#resultLoadMissingLine").get(0).clientHeight > 600) {
                $("#resultLoadMissingLine").css({ "max-height": "500px", "height": "600px", "overflow-y": "auto" })
            } else {
                $("#resultLoadMissingLine").css({ "max-height": "500px", "overflow-y": "auto" })
            }
            $("#lastScanAds").html(res.responseJSON.data_object.lastScanAdsTxt);
        }
    });
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

// ************************** Ads.txt ************************** //
$(document).ready(function () {
    // $("#inventoryTabs").on("click", "a.copy-tag", function (e) {
    //     const id = $(this).data("id");
    //     LoadAdTag(id)
    // })
    LoadAdTag($("#inventoryTabs").attr('data-id'))
    // change js load type => show js
    $("#nav-integration").on("change", ".js-load-type", function () {
        var type = $(this).val()
        if (!type) {
            return
        }
        if (type == "asynchronous") {
            $("#nav-integration").find(".js-type-asynchronous").removeClass("d-none").css("opacity", "0").animate({ opacity: "1" })
            $("#nav-integration").find(".js-type-normal").addClass("d-none")
        } else if (type == "normal") {
            $("#nav-integration").find(".js-type-asynchronous").addClass("d-none")
            $("#nav-integration").find(".js-type-normal").removeClass("d-none").css("opacity", "0").animate({ opacity: "1" })
        }
    })

    $("#nav-integration").on("click", ".copy-adtag", function (e) {
        CopyAdTag($(this))
    })

    $(".customize-sidebar").on("click", ".load-history", function () {
        $("body").removeClass("customize-box")
    })
})

function CopyAdTag(el) {
    /* Get the text field */
    var id = el.attr("data-id")
    var input = el.closest("tr").find("#pw_" + id);
    /* Select the text field */
    input.select();
    document.execCommand("copy");
    el.text("Copies")
    setTimeout(function () {
        el.html('<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" id="ds-icon-dashboard-copy">\n' +
            '                                <path d="M19.75 8.69h-2.37V6.31a2 2 0 0 0-2-2H4.25a2 2 0 0 0-2 2v7a2 2 0 0 0 2 2h2.37v2.38a2 2 0 0 0 2 2h11.13a2 2 0 0 0 2-2v-7a2 2 0 0 0-2-2zm-15.5 4.62v-7h11.13v2.38H8.62a2 2 0 0 0-2 2v2.62zm15.5 4.38H8.62v-7h11.13z">\n' +
            '                                </path>\n' +
            '                            </svg>')
    }, 1500);
}

function LoadAdTag(id) {
    let url = "/supply/copyAdTag"
    // $(`#copy-adtag-modal`).css("display", "block")
    $.ajax({
        url: url,
        type: "GET",
        contentType: "application/json",
        data: { id: id },
        beforeSend: function (xhr) {
        },
        error: function (jqXHR, exception) {
            const msg = AjaxErrorMessage(jqXHR, exception);
            new AlertError("AJAX ERROR: " + msg);
        }
    }).done(function (result) {
        $('#nav-integration').html(result);
        selectColor();
        changePlaceHolderTool();
        changeDesktopTag();
        changeVastType();
        HandleDesktopTag()
    });
}

function selectColor() {
    $("#nav-integration #select_text_color").on("input", function () {
        $("#nav-integration").find("#text_color").val($(this).val())
        $(this).closest("label").css("background", $(this).val()).css({
            "border": "1px solid #aab4c8",
            "border-right": 0
        })
    });
    $("#nav-integration #select_border_color").on("input", function () {
        $("#nav-integration").find("#border_color").val($(this).val())
        $(this).closest("label").css("background", $(this).val()).css({
            "border": "1px solid #aab4c8",
            "border-right": 0
        })
    });

    $("#nav-integration .form-control-color").on("input", function () {
        $(this).next("input").val($(this).val());
    });

    $("#nav-integration .form-control-color").next("input").on("input", function (e) {
        $(this).prev(".form-control-color").val(e.target.value)
    });
    // $("#nav-integration").on("change","#select_text_color", function () {
    //     alert('zzz');
    //     console.log($(this).val());
    //     $("#nav-integration").find("#text_color").val($(this).val())
    //     $(this).next("input").val($(this).val());
    // });
}

function changePlaceHolderTool() {
    $("#nav-integration #place_holder").on("change", function () {
        if (this.checked) {
            $(".box_place_holder").removeClass("d-none");
        } else {
            $(".box_place_holder").addClass("d-none");
        }
    });

}

function changeVastType() {
    $("#nav-integration .vast-type").on("change", function () {
        console.log($(this).val());

        let boxVast = $(this).closest(".box-vast");
        if ($(this).val() === "vpaid") {
            boxVast.find(".vast-type-vast").addClass("d-none");
            boxVast.find(".vast-type-vpaid").removeClass("d-none");
        } else if ($(this).val() === "vast") {
            boxVast.find(".vast-type-vpaid").addClass("d-none");
            boxVast.find(".vast-type-vast").removeClass("d-none");
        }
    })
}

function changeDesktopTag() {
    $("#desktop_tag").on("change", function () {
        HandleDesktopTag()
    });
}

function HandleDesktopTag() {
    $("#mobile_tag_sticky").removeAttr('disabled')
    $("#mobile_tag_sticky").find("option").removeAttr("hidden")
    $("#mobile_tag_display").find("option").removeAttr("hidden")
    let type = $("#desktop_tag").find(':selected').data("type")
    let valueDesktopTag = $("#desktop_tag").val();
    // let DesktopTag = $("#desktop_tag").val()
    if (type === 5) {
        let valueMobileTag = $("#mobile_tag_sticky").val();
        if (valueDesktopTag === valueMobileTag) {
            $("#mobile_tag_sticky").val('default');
            $("#mobile_tag_sticky").selectpicker("refresh");
        }
        $("#mobile_tag_sticky").attr("name", "mobile_tag")
        $("#toolCopyTag .box-mobile-sticky").removeClass("d-none")
        $("#toolCopyTag .box-mobile-display").addClass("d-none")
        $("#mobile_tag_display").removeAttr("name")
        $("#mobile_tag_sticky").find("option").attr("disabled", false)
        $("#mobile_tag_sticky").find("option[value=" + $("#desktop_tag").val() + "]").attr("disabled", true)
        $("#mobile_tag_sticky").selectpicker('refresh')
        $(".box_display").addClass("d-none");
    } else if (type === 1) {
        let valueMobileTag = $("#mobile_tag_display").val();
        if (valueDesktopTag === valueMobileTag) {
            $("#mobile_tag_display").val('default');
            $("#mobile_tag_display").selectpicker("refresh");
        }
        $("#mobile_tag_display").attr("name", "mobile_tag")
        $("#toolCopyTag .box-mobile-display").removeClass("d-none")
        $("#toolCopyTag .box-mobile-sticky").addClass("d-none")
        $("#mobile_tag_sticky").removeAttr("name")
        $("#mobile_tag_display").find("option").attr("disabled", false)
        $("#mobile_tag_display").find("option[value=" + $("#desktop_tag").val() + "]").attr("disabled", true)
        $("#mobile_tag_display").selectpicker('refresh')
        $(".box_display").removeClass("d-none");
    }
}

function HandleTabTool() {
    $('#nav-integration').on("click", ".build-tool", function () {
        loadScriptTool(this)
    })
}

function loadScriptTool(element) {
    const formElement = $("#toolCopyTag");
    const buttonElement = $(element);
    const submitButtonText = buttonElement.text();
    const submitButtonTextLoading = "Loading...";
    var postData = formElement.serializeObject();
    postData.desktop_tag = parseInt(postData.desktop_tag)
    postData.mobile_tag = parseInt(postData.mobile_tag)
    // postData.user_id = parseInt(postData.user_id)

    if (postData.place_holder === "on") {
        postData.place_holder = 1
    } else {
        postData.place_holder = 2
    }

    let url = "/supply/buildScript"
    $.ajax({
        url: url,
        type: "POST",
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
    }).done(function (result) {
        $('#nav-integration #script_build').text(result);
    });
}

function HandleTabPlayZone() {
    $('#nav-integration').on("click", ".build-playzone", function () {
        loadCodeQuiz(this)
    })
}

function loadCodeQuiz(element) {
    const formElement = $("#toolGenerateQuiz");
    const buttonElement = $(element);
    const submitButtonText = buttonElement.text();
    const submitButtonTextLoading = "Loading...";
    var postData = formElement.serializeObject();


    buttonElement.attr('disabled', true).text(submitButtonTextLoading);
    let text = '';
    if (postData.quiz) {
        text = `
<div class="adsbyvli" data-ad-slot='pw_${postData.playzone_tag}'></div> 
<script type="text/javascript"> 
    (vitag.Init = window.vitag.Init || []).push(function () { 
        viAPItag.initPlayZone('pw_${postData.playzone_tag}', { "quizIds": [${postData.quiz}] }) 
    }) 
</script>`;
    }
    $('#toolGenerateQuiz #code_snippet_quiz').text(text.trim());
    buttonElement.attr('disabled', false).text(submitButtonText);
}

function selectColorPlayZone() {
    let elementInputColor = $(".input-color");
    elementInputColor.on("input", function () {
        $(this).next("input").val($(this).val());
    });
    elementInputColor.next("input").on("input", function (e) {
        $(this).prev("input").val(e.target.value);
    });
}

function selectAdSizeCopyTag() {
    $("#nav-integration").on("change", ".select-ad-size", function () {
        var ad_size = $(this).val()
        if (ad_size === "all") {
            $(".ad-size-all").removeClass("d-none")
            $(".ad-size-responsive").addClass("d-none")
            $(".ad-size-fixed").addClass("d-none")
        } else if (ad_size === "fixed") {
            $(".ad-size-fixed").removeClass("d-none")
            $(".ad-size-all").addClass("d-none")
            $(".ad-size-responsive").addClass("d-none")
        } else if (ad_size === "responsive") {
            $(".ad-size-responsive").removeClass("d-none")
            $(".ad-size-all").addClass("d-none")
            $(".ad-size-fixed").addClass("d-none")
        }
    })
}