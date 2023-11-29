const AdsenseAdSlot = require("./inc.adsense-ad-slot")
// import * as AdsenseAdSlot from "./adsense-ad-slot"
// const {SelectSize, RemoveSize} = require("./adsense-ad-slot");
// import {RemoveSize, SelectSize} from "./adsense-ad-slot";
import PrebidClass from "../../../../../www/themes/muze/assets/js/prebid.js";

const info = {
    targetPrebid: "#tab-prebid",
    formId: "AddBidder",
    urlLoadParam: "/line-item/loadParam",
    urlSubmit: "",
    showHideSelectBidder: true,
    typeSelect: "selectpicker",
    ajaxAddParam: "/line-item/addParamBidder",
    isCheckParamRequired: "/line-item/checkParam",
}

let Prebid = new PrebidClass(info);

let optionDomainsLoadMore = {
    idSelect: "#list_inventory",
    idBox: ".box_inventory",
    idIsMoreData: "#is_more_inventory",
    isMoreData: true,
    lastPage: false,
    page: 1,
    isSearch: false,
    search: "",
    idSearch: "#search_inventory",
    urlAjax: "/target/loadInventory",
    optionAjax: "domain",
    filterAjax: [],
    list_selected: [],
    btnInclude: "add_inventory",
    btnRemove: "remove_inventory",
    idEmpty: "empty_ad_domain",
    idEmptyLoad: "load_empty_domain",
    checkLoadMore: false,
    text: "#text_for_domain",
    btnRemoveAll: "remove_all_domain",
    btnSelectAll: "select_all_domain",
    block: ".block_domain",
    container: ".container_domain",
    boxEmpty: "box_empty_domain"
}

let optionAdTagsLoadMore = {
    idSelect: "#list_ad_tag",
    idBox: ".box_ad_tag",
    idIsMoreData: "#is_more_tag",
    isMoreData: true,
    lastPage: false,
    page: 1,
    isSearch: false,
    search: "",
    idSearch: "#search_ad_tag",
    urlAjax: "/target/loadInventory",
    optionAjax: "adtag",
    filterAjax: [],
    list_selected: [],
    btnInclude: "add_tag",
    btnRemove: "remove_tag",
    idEmpty: "empty_ad_tag",
    idEmptyLoad: "load_empty_ad_tag",
    checkLoadMore: false,
    text: "#text_for_tag",
    btnRemoveAll: "remove_all_tag",
    btnSelectAll: "select_all_tag",
    block: ".block_tag",
    container: ".container_tag",
    boxEmpty: "box_empty_tag"
}

let optionAdFormatsLoadMore = {
    idSelect: "#list_ad_format",
    idBox: ".box_ad_format",
    idIsMoreData: "#is_more_adformat",
    isMoreData: true,
    lastPage: false,
    page: 1,
    isSearch: false,
    search: "",
    idSearch: "#search_ad_format",
    urlAjax: "/target/loadInventory",
    optionAjax: "adformat",
    filterAjax: [],
    list_selected: [],
    btnInclude: "add_format",
    btnRemove: "remove_adFormat",
    idEmpty: "empty_ad_format",
    text: "#text_for_format",
    btnRemoveAll: "remove_all_format",
    btnSelectAll: "select_all_format",
    block: ".block_format",
    container: ".container_format",
    boxEmpty: "box_empty_format",
}

let optionAdSizesLoadMore = {
    idSelect: "#list_ad_size",
    idBox: ".box_ad_size",
    idIsMoreData: "#is_more_adsize",
    isMoreData: true,
    lastPage: false,
    page: 1,
    isSearch: false,
    search: "",
    idSearch: "#search_ad_size",
    urlAjax: "/target/loadInventory",
    optionAjax: "adsize",
    filterAjax: [],
    list_selected: [],
    btnInclude: "add_size",
    btnRemove: "remove_adSize",
    idEmpty: "empty_ad_size",
    idEmptyLoad: "load_empty_ad_size",
    checkLoadMore: false,
    text: "#text_for_size",
    btnRemoveAll: "remove_all_size",
    btnSelectAll: "select_all_size",
    block: ".block_size",
    container: ".container_size",
    boxEmpty: "box_empty_size",
}

let optionGeographyLoadMore = {
    idSelect: "#list_geography",
    idBox: ".box_geography",
    idIsMoreData: "#is_more_geography",
    isMoreData: true,
    lastPage: false,
    page: 1,
    isSearch: false,
    search: "",
    idSearch: "#search_geography",
    urlAjax: "/target/loadInventory",
    optionAjax: "geography",
    list_selected: [],
    btnInclude: "add_country",
    btnRemove: "remove_Geo",
    idEmpty: "empty_ad_geo",
    idEmptyLoad: "load_empty_geo",
    checkLoadMore: false,
    text: "#text_for_geo",
    btnRemoveAll: "remove_all_geo",
    btnSelectAll: "select_all_geo",
    block: ".block_geo",
    container: ".container_geo",
    boxEmpty: "box_empty_geo"
}

let optionDeviceLoadMore = {
    idSelect: "#list_device",
    idBox: ".box_device",
    idIsMoreData: "#is_more_device",
    isMoreData: true,
    lastPage: false,
    page: 1,
    isSearch: false,
    search: "",
    idSearch: "#search_device",
    urlAjax: "/target/loadInventory",
    optionAjax: "device",
    list_selected: [],
    list_original: [],
    btnInclude: "add_device",
    btnRemove: "remove_device",
    idEmpty: "empty_ad_device",
    idEmptyLoad: "load_empty_device",
    checkLoadMore: false,
    text: "#text_for_device",
    btnRemoveAll: "remove_all_device",
    btnSelectAll: "select_all_device",
    block: ".block_device",
    container: ".container_device",
    boxEmpty: "box_empty_device"
}

let filterAjax = {
    inventory: [],
    format: [],
    size: [],
}

$(document).ready(function () {
    loadTarget();

    AdsenseAdSlot.SelectSize()
    AdsenseAdSlot.RemoveSize()
    AdsenseAdSlot.AddAdSlot()
    SelectAccount();

    PrebidLoad();

    if (!info.isEdit) {
        SelectServerType();
    }

    EventToggleSetupTargeting();
    EventCollapseTarget();
    HandleCollapse();

    $('[data-bs-toggle="popover"]').popover({
        html: true,
        sanitize: false,
    });

    FormEvent();
})

function FormEvent() {
    let formElement = $("#" + info.formId);

    //Xóa bỏ text error khi nhập data vào ô input
    formElement.on("input", ".list_param .param_value", function (e) {
        let inputElement = $(this)
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
        let parent = $(this).parent()
        if (parent.hasClass("is-invalid")) {
            parent.removeClass("is-invalid").nextAll("span.invalid-feedback").empty();
        }
    });
    formElement.on('keypress', "input", function (e) {
        var keyCode = e.keyCode || e.which;
        if (keyCode === 13) {
            e.preventDefault();
            submit()
        }
    });
    formElement.find(".submit").on("click", function (e) {
        e.preventDefault();
        $.each(formElement.find("input.is-invalid"), function () {
            $(this).removeClass("is-invalid").next(".invalid-feedback").empty();
        });
        submit()
    });

    function submit() {
        let postData = formElement.serializeObject();
        postData.listInventory = optionDomainsLoadMore.list_selected
        postData.listGeo = optionGeographyLoadMore.list_selected
        postData.listDevice = optionDeviceLoadMore.list_selected
        postData.listAdtag = optionAdTagsLoadMore.list_selected
        postData.listAdFormat = optionAdFormatsLoadMore.list_selected
        postData.listAdSize = optionAdSizesLoadMore.list_selected
        postData.priority = parseInt(postData.priority)
        postData.type = parseInt(postData.type)
        postData.connection_type = parseInt(postData.connection_type)
        postData.line_item_type = parseInt(postData.line_item_type)
        postData.server_type = info.isEdit ? parseInt($("#server_type").val()) : parseInt(postData.server_type)
        postData.rate = parseInt(postData.rate)
        postData.linked_gam = parseInt(postData.linked_gam)
        postData.bidder_params = addParamBidder()
        postData.adsense_ad_slots = addAdsenseAdSlot()
        if (!Array.isArray(postData.select_account)) {
            const selectAccount = [];
            if (postData.select_account !== "") {
                selectAccount.push(parseInt(postData.select_account))
            }
            postData.select_account = selectAccount
        }

        if (info.isEdit) {
            postData.id = parseInt(postData.id)
        }

        const data = JSON.stringify(postData);
        const buttonElement = $(".submit");
        $.ajax({
            url: info.urlSubmit,
            type: "POST",
            dataType: "JSON",
            contentType: "application/json",
            data: data,
            beforeSend: function (xhr) {
                buttonElement.attr('disabled', true).text("Loading...");
            },
            error: function (jqXHR, exception) {
                const msg = AjaxErrorMessage(jqXHR, exception);
                new AlertError("AJAX ERROR: " + msg);
                buttonElement.attr('disabled', false).text("Submit");
            },
            success: function (responseJSON) {
                buttonElement.attr('disabled', false).text("Submit");
            },
            complete: function (res) {
                Added(res.responseJSON, formElement);
            }
        });
    }
}

function PrebidLoad() {
    Prebid.buttonAddBidder(); //Add New Bidder -> OK
    Prebid.EventLoadParamsAfterSelectBidder(); //Load Params After Select Bidder -> OK
    Prebid.EventClickPrebidTabConfig(); //Click Prebid Tab -> OK
    Prebid.EventAddParam(); //Add New Param -> OK
    Prebid.SelectParamPB(); //Select Param -> OK
    Prebid.EventSelectBidderType(); //Select Bidder Type -> OK
    Prebid.EventRemoveParam(); //Remove Param -> OK
    Prebid.EventRemoveBidder(); //Remove Bidder -> OK
}

function loadTarget() {
    //Load more domains
    setUpLoadMoreData(optionDomainsLoadMore)

    //Load more ad tag
    setUpLoadMoreData(optionAdTagsLoadMore)

    //Load more ad format
    setUpLoadMoreData(optionAdFormatsLoadMore)

    //Load more ad size
    setUpLoadMoreData(optionAdSizesLoadMore)

    //Load more ad size
    setUpLoadMoreData(optionGeographyLoadMore)

    //Load more ad size
    setUpLoadMoreData(optionDeviceLoadMore)
}

function EventToggleSetupTargeting() {
    $(".line-item-tab").on("click", ".nav-link", function () {
        $(".line-item-tab").find(".nav-link").removeClass("pp-4")
        $(this).addClass("pp-4")
        var tab = $(this).attr("data-tab")
        if (tab != "1") {
            $(this).addClass("at-1")
        } else {
            $(".line-item-tab").find('.at-1').removeClass('at-1')
        }
    })
}

function EventCollapseTarget() {
    // Target
    $('[data-toggle="collapse"]').on("click", function () {
        var element = $(this)

        element.toggleClass("dm14")
        if (element.hasClass("dm14")) {
            // element.collapse('show')
            element.find(".dm23").find("button").html(
                '<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor" class="bi bi-dash-lg" viewBox="0 0 16 16">\n' +
                '<path d="M0 8a1 1 0 0 1 1-1h14a1 1 0 1 1 0 2H1a1 1 0 0 1-1-1z">\n' +
                '</path>\n' +
                '</svg>')
        } else {
            element.find(".dm23").find("button").html(
                '<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor" class="bi bi-plus-lg" viewBox="0 0 16 16">\n' +
                '<path d="M8 0a1 1 0 0 1 1 1v6h6a1 1 0 1 1 0 2H9v6a1 1 0 1 1-2 0V9H1a1 1 0 0 1 0-2h6V1a1 1 0 0 1 1-1z"></path>\n' +
                '</svg>')
        }

        $("#nav-target").find('[data-toggle="collapse"]').each(function () {
            if (element[0] != $(this)[0]) {
                $(this).removeClass("dm14")
                $(this).find(".dm23").find("button").html(
                    '<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor" class="bi bi-plus-lg" viewBox="0 0 16 16">\n' +
                    '<path d="M8 0a1 1 0 0 1 1 1v6h6a1 1 0 1 1 0 2H9v6a1 1 0 1 1-2 0V9H1a1 1 0 0 1 0-2h6V1a1 1 0 0 1 1-1z"></path>\n' +
                    '</svg>')
            }
        })
    })
}

//Chọn Type Prebid / Google -> OK
function SelectServerType() {
    $("#server_type").on("changed.bs.select", function () {
        let selectValue = $(this).val();
        if (selectValue === "2") {
            $(info.targetPrebid).addClass("d-none")
            $(".w3-bar--config").addClass("d-none")
            $("#tab-google").removeClass("d-none")
            $(".box_type_google_hide").addClass("d-none")
            checkSelectConnectionType($("#select_connection_type").val())
        } else {
            $(info.targetPrebid).removeClass("d-none")
            $(".w3-bar--config").removeClass("d-none")
            $("#tab-google").addClass("d-none")
            $(".box_type_google_hide").removeClass("d-none")
        }
    })
}

function SelectAccount() {
    $('#select_account').on("changed.bs.select", function () {
        let typeAccount = $('#select_account').find('option:selected', this).attr("data-type").trim()
        let isSystem = $('#select_account').find('option:selected', this).attr("data-system")
        if (typeAccount === "Adx" && isSystem === "true") {
            $("#select_connection_type option").remove()
            var optionMCM = new Option("MCM", '3', false, true);
            // Append it to the select
            $('#select_connection_type').append(optionMCM).selectpicker('refresh');
            getListLinkedGam(true, $(this).val())
        } else {
            if (!$('#select_connection_type').find("option[value='1']").length) {
                $("#select_connection_type option").remove()
                var optionAdUnit = new Option("Ad units", '1', false, true);
                var optionLineItem = new Option("Line items", '2', false, false);
                // Append it to the select
                $('#select_connection_type').append(optionAdUnit).append(optionLineItem).selectpicker('refresh');
                getListLinkedGam(false, $(this).val())
            }
        }
        checkSelectConnectionType($("#select_connection_type").val())
        checkShowAdsenseSlot()
    })
    $('#select_connection_type').on("changed.bs.select", function () {
        checkSelectConnectionType($(this).val())
        checkShowAdsenseSlot()
    })
    $('#select_line_item_type').on("changed.bs.select", function () {
        checkShowAdsenseSlot()
    })
}

function getListLinkedGam(ppAdx = false, bidderAdxId) {
    $.ajax({
        url: "/line-item/listLinkedGam",
        type: "GET",
        data: {
            ppAdx: ppAdx,
            bidderAdxId: bidderAdxId,
        },
        success: function (json) {
            $("#linked_gam option").remove();
            $("#linked_gam").selectpicker('refresh');
            $.each(json.list_gam, function (i, gam) {
                var option = new Option(gam.network_name, gam.id, false, false);
                // Append it to the select
                $('#linked_gam').append(option).selectpicker('refresh');
            })
        },
        error: function (xhr) {
            console.log(xhr)
        }
    })
}

function checkSelectConnectionType(value) {
    if (value === "2") {
        $('.box_line_item_type').removeClass('d-none')
    } else {
        $('.box_line_item_type').addClass('d-none')
    }
}

function checkShowAdsenseSlot() {
    if ($("#server_type").val() === "2" && $('#select_account').find('option:selected', this).attr("data-type") === "Adsense" && $("#select_connection_type").val() === "2" && $("#select_line_item_type").val() === "1") {
        $('.box-adsense-ad-slot').removeClass('d-none')
    } else {
        $('.box-adsense-ad-slot').addClass('d-none')
    }
}

function SelectAdsenseSize(size, itemLabel, itemPlaceholder) {
    return `<div class="box-c">
                <hr class="mt-5 hr_custom bg-gray-400">
                <div class="bidder-box">
                    <div class="row my-4">
                        <div class="d-flex align-items-center bidder-name">
                            <label class="col-form-label form-label form-label-lg w-input-group text-uppercase">
                                ${size}
                            </label>
                            <div class="d">
                                <a class="btn p-1 ms-2 rm_c remove-ad-slot" data-size="${size}">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="15"
                                         height="15" fill="currentColor" class="bi bi-x-lg"
                                         viewBox="0 0 16 16">
                                        <path d="M1.293 1.293a1 1 0 0 1 1.414 0L8 6.586l5.293-5.293a1 1 0 1 1 1.414 1.414L9.414 8l5.293 5.293a1 1 0 0 1-1.414 1.414L8 9.414l-5.293 5.293a1 1 0 0 1-1.414-1.414L6.586 8 1.293 2.707a1 1 0 0 1 0-1.414z"></path>
                                    </svg>
                                </a>
                            </div>
                        </div>
                    </div>
    
                    <div class="row my-4">
                        <div class="d-flex align-items-center">
                            <label class="col-form-label form-label form-label-lg w-input-group">
                                ${itemLabel}
                            </label>
                            <input id="adsense-ad-slot-${size}" type="text"
                                   class="form-control w-50-custom param_value"
                                   placeholder="${itemPlaceholder}"
                                   data-size="${size}" value="">
                            <span class="form-text ms-2 w-auto invalid-feedback"></span>
                        </div>
                    </div>
                </div>
            </div>`
}

function EventClickCollapseTarget() {
    $('.btn-collapse').click(function (e) {
        if ($(this).hasClass('collapsed')) {
            $(this).find(".text-notify-target").removeClass("d-none")
        } else {
            $(".text-notify-target").removeClass("d-none")
            $(this).find(".text-notify-target").addClass("d-none")
        }
    })
}

function addAdsenseAdSlot() {
    let adsenseAdSlots = []
    $("#adsense_ad_slot_item > div.box-c > .param_value").each(function () {
        let adSlot = {}
        let size = $(this).data("size")
        let adSlotId = $(this).val().trim()
        adSlot.size = size
        adSlot.ad_slot_id = adSlotId
        adsenseAdSlots.push(adSlot)
    })
    return adsenseAdSlots
}

//Get data để submit line item -> OK
function addParamBidder() {
    let bidderParams = []
    $(`${info.targetPrebid} div.box-c > div.bidder-box`).each(function () {
        let configType = $(this).closest(".tab-prebid-wrapper").attr("id");

        var bidderBox = $(this)
        if ($(this).attr('data-name') !== "") {
            let bidderParam = {}
            bidderParam["id"] = parseInt($(this).attr('data-id'))
            bidderParam["name"] = $(this).attr('data-name').toLowerCase()
            bidderParam["bidder_type"] = parseInt($(this).find("select.type-bidder").val())
            bidderParam["bidder_index"] = parseInt($(this).attr('data-index'))
            bidderParam["config_type"] = configType.replace("tab-prebid-", "");
            let params = {}
            bidderBox.find('.param_value').each(function () {
                if (!$(this).attr('data-name')) {
                    return true;
                }
                if (!$(this).hasClass("add-param")) {
                    let type = $(this).attr('data-type').toLowerCase()
                    let name = $(this).attr('data-name')
                    params[name] = {
                        type: type,
                        value: this.value.trim(),
                        addParam: 2,
                    }
                } else {
                    let type = $(this).attr('data-type').toLowerCase()
                    let name = $(this).attr('data-name').trim()
                    $(this).attr("id", bidderParam["id"] + "-" + name + "-" + bidderParam["bidder_index"])
                    if (params[name] === undefined) {
                        params[name] = {
                            type: type,
                            value: this.value.trim(),
                            addParam: 1,
                        }
                    }
                }
            })
            bidderParam["params"] = params
            bidderParams.push(bidderParam)
        }
    });
    return bidderParams
}

//Handle sau khi submit xong line item -> OK
function Added(response) {
    switch (response.status) {
        case "error":
            if (response.errors.length === 1 && response.errors[0].id === "") {
                new AlertError(response.errors[0].message);
            } else if (response.errors.length > 0) {
                $.each(response.errors, function (key, value) {
                    if (value.id === "select-adsense-ad-slot-size") {
                        let inputElement = $("body").find("#adsense_ad_slot_item");
                        inputElement.addClass("is-invalid").find("span.invalid-feedback").text(value.message);
                        inputElement.find(".param_value").addClass("is-invalid");
                    }

                    let inputElement = $("body").find("#" + value.id);
                    inputElement.addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
                    inputElement.parent().addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
                    // inputElement.next().find(".select2-selection").addClass("select2-is-invalid")
                    // inputElement.closest('.box-select2').addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message)
                    if (value.id === "text_for_domain") {
                        let msg = value.message
                        inputElement.html(`<small class="text-danger">${msg}</small>`)
                        $('.domain-card').addClass("domain-card-invalid")
                        $("#collapseDomains").collapse('show')
                        $('div[href="#collapseDomains"]').addClass("dm14")
                        if (response.errors.length == 1) {
                            $("#tab-target").click()
                        }
                    }
                });


                //Xử lí focus sang tab có error config của prebid
                let isPrebid = !$(info.targetPrebid).hasClass("d-none");

                if (isPrebid) {
                    try {
                        let wrapper = $("#" + response.errors[0].id).closest(".tab-prebid-wrapper");

                        if (wrapper.length && wrapper.hasClass("d-none")) {
                            $(`.w3-bar--config-item[data-id="${wrapper.attr("id")}"]`).click();
                        }
                    } catch (error) { }
                }

                $("#" + response.errors[0].id).focus();
                $("#" + response.errors[0].id).prev('label').focus();

            } else {
                new AlertError("Error!");
            }
            break
        case "success":
            if ($("#select_line_item_type").val() !== "") {
                $("#select_line_item_type").prop("disabled", true);
            }
            NoticeSuccess("Line item has been updated successfully")
            setTimeout(function () {
                window.location.href = "/line-item";
            }, 1000)
            break
        default:
            new AlertError("Undefined");
            break
    }
}

function setUpLoadMoreData(optionTarget) {
    firstLoad(optionTarget)

    checkBoxSelectedEmpty(optionTarget)

    setUpClickAndSearch(optionTarget)

    getListChecked(optionTarget)

    scrollToBottom(optionTarget)
}

function removeAll(optionTarget) {
    optionTarget.list_selected.map(item => {
        AppendList(item.id, item.name, optionTarget)
    })
    optionTarget.list_selected.map(item => {
        // let divCountry = $(optionTarget.idSelect).find(`div[id = '${item.id}']`)
        // divCountry.attr("hidden", false)
        // UpdateListSelected(item.id, item.name, optionTarget)
        switch (optionTarget.optionAjax) {
            case optionDomainsLoadMore.optionAjax:
                let currentIdxInventory = filterAjax.inventory.indexOf(item.id)
                if (currentIdxInventory > -1) {
                    filterAjax.inventory.splice(currentIdxInventory, 1)
                }
                optionAdTagsLoadMore.page = 1
                optionAdTagsLoadMore.isSearch = true
                LoadMoreData(optionAdTagsLoadMore)
                break

            case optionAdFormatsLoadMore.optionAjax:
                let currentIdxFormat = filterAjax.format.indexOf(item.id)
                if (currentIdxFormat > -1) {
                    filterAjax.format.splice(currentIdxFormat, 1)
                }
                optionAdTagsLoadMore.page = 1
                optionAdTagsLoadMore.isSearch = true
                LoadMoreData(optionAdTagsLoadMore)
                break
            case optionAdSizesLoadMore.optionAjax:
                let currentIdxSize = filterAjax.size.indexOf(item.id)
                if (currentIdxSize > -1) {
                    filterAjax.size.splice(currentIdxSize, 1)
                }
                optionAdTagsLoadMore.page = 1
                optionAdTagsLoadMore.isSearch = true
                LoadMoreData(optionAdTagsLoadMore)
                break
            case optionAdTagsLoadMore.optionAjax:
                let currentIdxAdTag = optionAdTagsLoadMore.filterAjax.indexOf(item.id)
                if (currentIdxAdTag > -1) {
                    optionAdTagsLoadMore.filterAjax.splice(currentIdxAdTag, 1)
                }
                // FilterAdTagChange(optionAdTagsLoadMore,"remove")
                break
        }
    })
    optionTarget.list_selected = []
    $(optionTarget.idSelect).html("")
    optionTarget.page = 1
    optionTarget.isSearch = false
    optionTarget.isMoreData = true
    LoadMoreData(optionTarget)
    checkBoxSelectedEmpty(optionTarget)
    DisplayTextSelected(optionTarget)
}

function firstLoad(optionTarget) {
    optionTarget.isMoreData = $(optionTarget.idSearch).val()
    optionTarget.page = 1
    LoadMoreData(optionTarget)
}

function checkBoxSelectedEmpty(optionTarget) {
    if (optionTarget.list_selected.length === 0) {
        // $("." + optionTarget.btnRemoveAll).attr("hidden", true)
        $(optionTarget.idBox).html("")
        $(optionTarget.idBox).append(`<div id="${optionTarget.idEmpty}" class="d-flex flex-row align-items-center px-md-2" style="height: 25px">
            <div class="col p-0"><h6
                        class="m-0 font-weight-semibold fs-12 text-center">
                    No data selected</h6></div>
        </div>`)
    } else {
        $("." + optionTarget.btnRemoveAll).attr("hidden", false)
        let text = "#" + optionTarget.idEmpty
        $(text).remove()
    }
}

function setUpClickAndSearch(optionTarget) {
    let fullInfo = [optionDomainsLoadMore, optionAdTagsLoadMore, optionAdFormatsLoadMore, optionAdSizesLoadMore, optionGeographyLoadMore];
    // $(optionTarget.idSelect).on("click", "div.list-group div.list-group-item div.col-auto button." + optionTarget.btnInclude, function (e) {
    $(optionTarget.idSelect).on("click", "button." + optionTarget.btnInclude, function (e) {
        ClickInclude($(this), e, optionTarget)
        CheckAddAll1Page(optionTarget)
        RemoveNotifyNoData(optionTarget)

        //Xóa báo lỗi card domain
        if ($(this).hasClass("add_inventory")) {
            $(".domain-card").removeClass("domain-card-invalid")
        }

        Prebid.handleDisabledTabs(optionTarget, fullInfo);
    })

    $(optionTarget.idSearch).on("input", function () {
        optionTarget.page = 1
        optionTarget.isSearch = true
        optionTarget.lastPage = false
        optionTarget.search = $(this).val()
        LoadMoreData(optionTarget)
        scrollToBottom(optionTarget)
    })

    $(optionTarget.idSearch).on('keyup keypress', function (e) {
        var keyCode = e.keyCode || e.which;
        if (keyCode === 13) {
            e.preventDefault();
            return false;
        }
    });

    $(optionTarget.block).on("click", "a." + optionTarget.btnRemoveAll, function (e) {
        e.preventDefault()
        removeAll(optionTarget);

        Prebid.handleDisabledTabs(optionTarget, fullInfo);
    })

    // $(optionTarget.idBox).on("click", ".list-group .list-group-item .flex-row .col-auto button." + optionTarget.btnRemove, function (e) {
    $(optionTarget.idBox).on("click", "button." + optionTarget.btnRemove, function (e) {
        RemoveInclude($(this), e, optionTarget)
        RemoveNotifyNoData(optionTarget)

        Prebid.handleDisabledTabs(optionTarget, fullInfo);
    })
}

function getListChecked(optionTarget) {
    let divs = $(optionTarget.idBox).find("." + optionTarget.btnRemove)
    divs.each(function (index, elm) {
        let id = parseInt(elm.id)
        let name = elm.name
        optionTarget.list_selected.push({ id: id, name: name })
    })
}

function CheckAddAll1Page(optionTarget) {
    let div = $(optionTarget.idSelect).get(0)
    if (div.scrollHeight === 350 && optionTarget.isMoreData) {
        optionTarget.search = $(optionTarget.idSearch).val()
        optionTarget.isSearch = false
        LoadMoreData(optionTarget)
    }
}

function scrollToBottom(optionTarget) {
    $(optionTarget.idSelect).on("scroll", function () {
        let div = $(this).get(0);
        let flag = CheckIfScrollBottom(Math.round(div.scrollTop), div.clientHeight, div.scrollHeight)
        if (flag && optionTarget.isMoreData) {
            optionTarget.search = $(optionTarget.idSearch).val()
            optionTarget.isSearch = false
            if (!optionTarget.checkLoadMore) {
                LoadMoreData(optionTarget)
            }
        }
    });
}

function CheckIfScrollBottom(scrollTop, clientHeight, scrollHeight = 0) {
    if ((scrollTop + clientHeight + 1) >= scrollHeight) {
        return true
    } else {
        return false
    }
}

function RemoveNotifyNoData(optionTarget) {
    let a_elm = $(optionTarget.idSelect).find("div.target-item")
    let a_empty = $("." + optionTarget.boxEmpty)
    if (a_elm.length > 0) {
        if (a_empty.length > 0) {
            a_empty.attr("hidden", true)
        } else {
            // $(optionTarget.idSelect).append(`<div class="list-group list-group-flush my-n3 pt-3 ${optionTarget.boxEmpty}" hidden>
            //     <div class="list-group-item">
            //         <div class="d-flex flex-row align-items-center px-md-2">
            //             <div class="col p-0 text-center"><h6 class="m-0 font-weight-semibold fs-12">No data available</h6></div>
            //         </div>
            //     </div>
            // </div>`)
        }
    } else {
        if (a_empty.length > 0) {
            a_empty.attr("hidden", false)
        } else {
            $(optionTarget.idSelect).append(`<div class="list-group list-group-flush my-n3 pt-3 ${optionTarget.boxEmpty}">
                <div class="list-group-item">
                    <div class="d-flex flex-row align-items-center px-md-2">
                        <div class="col p-0 text-center"><h6 class="m-0 font-weight-semibold fs-12">No data available</h6></div>
                    </div>
                </div>
            </div>`)
        }
    }
}

function LoadMoreData(optionTarget) {
    // if (optionTarget.isMoreData === false) {
    //     return
    // }
    optionTarget.checkLoadMore = true
    let currentPage
    if (optionTarget.page === 0) {
        currentPage = optionTarget.page + 1
    } else {
        currentPage = optionTarget.page
    }
    let listId = []
    optionTarget.list_selected.map(item => {
        listId.push(item.id)
    })
    $.ajax({
        url: optionTarget.urlAjax,
        type: "GET",
        data: {
            key: currentPage,
            search: optionTarget.search,
            option: optionTarget.optionAjax,
            filter: JSON.stringify(filterAjax),
            selected: JSON.stringify(listId)
        },
        success: function (json) {
            if (optionTarget.isSearch) {
                $(optionTarget.idSelect).html("")
            }
            $(optionTarget.idSelect).append(json.data)
            optionTarget.isMoreData = json.is_more_data
            // $(optionTarget.idIsMoreData).val(json.is_more_data)
            if (json.is_more_data === true) {
                optionTarget.page = json.current_page + 1
            }
            HideSelected(optionTarget)
            RemoveNotInFilter(optionTarget, json.list_filter)
            DisplayTextSelected(optionTarget)
            optionTarget.checkLoadMore = false
            optionTarget.lastPage = json.last_page
            DisplayNoDataAvailable(optionTarget, json.total, json.current_page)
        },
        error: function (xhr) {
            console.log(xhr)
        }
    })
}

function FilterAdTagChange(optionTarget, type) {
    $.ajax({
        url: "/target/filterAdTag",
        type: "POST",
        contentType: "application/json",
        data: JSON.stringify(
            {
                type: type,
                tagId: optionTarget.filterAjax,
                serverType: parseInt($("#server_type").val())
            }
        ),
        success: function (json) {
            CheckIncludeAfterChangeAdTag(json)
        },
        error: function (xhr) {
            console.log(xhr)
        }
    })
}

function CheckIncludeAfterChangeAdTag(json) {
    let optionTarget
    json.inventory_id.map(item => {
        optionTarget = optionDomainsLoadMore
        let currentIdx = optionTarget.list_selected.map(item => {
            return item.id
        }).indexOf(item.id)
        if (currentIdx === -1) {
            optionTarget.list_selected.push({ id: item.id, name: item.name })
            DisplaySelected(item.id, item.name, optionTarget)
        }
        HideSelected(optionTarget)
        DisplayTextSelected(optionTarget)
        checkBoxSelectedEmpty(optionTarget)
        ChangeAdTagWithFilter(item.id, currentIdx, optionTarget)
    })
    json.ad_format_id.map(item => {
        optionTarget = optionAdFormatsLoadMore
        let currentIdx = optionTarget.list_selected.map(item => {
            return item.id
        }).indexOf(item.id)
        if (currentIdx === -1) {
            optionTarget.list_selected.push({ id: item.id, name: item.name })
            DisplaySelected(item.id, item.name, optionTarget)
        }
        HideSelected(optionTarget)
        DisplayTextSelected(optionTarget)
        checkBoxSelectedEmpty(optionTarget)
        ChangeAdTagWithFilter(item.id, currentIdx, optionTarget)
    })
    json.ad_size_id.map(item => {
        optionTarget = optionAdSizesLoadMore
        let currentIdx = optionTarget.list_selected.map(item => {
            return item.id
        }).indexOf(item.id)
        if (currentIdx === -1) {
            optionTarget.list_selected.push({ id: item.id, name: item.name })
            DisplaySelected(item.id, item.name, optionTarget)
        }
        HideSelected(optionTarget)
        DisplayTextSelected(optionTarget)
        checkBoxSelectedEmpty(optionTarget)
        ChangeAdTagWithFilter(item.id, currentIdx, optionTarget)
    })
    ReloadAdTag()
}

function ClickInclude(element, event, optionTarget) {
    event.preventDefault();
    let id = parseInt(element.data("id"))
    let name = element.attr("name")
    let div = $(optionTarget.idSelect).find(`div[id = '${id}']`)
    // div.attr("hidden", true)
    div.remove()
    UpdateListSelected(id, name, optionTarget)
    ChangeAdTagWithFilter(id, -1, optionTarget)
    if (optionTarget.optionAjax !== optionAdTagsLoadMore.optionAjax) {
        ReloadAdTag()
    }
}

function ChangeAdTagWithFilter(id, currentIdx, optionTarget) {
    switch (optionTarget.optionAjax) {
        case optionDomainsLoadMore.optionAjax:
            if (currentIdx === -1) {
                filterAjax.inventory.push(id)
            }
            break
        case optionAdFormatsLoadMore.optionAjax:
            if (currentIdx === -1) {
                filterAjax.format.push(id)
            }
            break
        case optionAdSizesLoadMore.optionAjax:
            if (currentIdx === -1) {
                filterAjax.size.push(id)
            }
            break
        case optionAdTagsLoadMore.optionAjax:
            optionAdTagsLoadMore.filterAjax.push(id)
            FilterAdTagChange(optionAdTagsLoadMore, "include")
            break
    }
}

function ReloadAdTag() {
    optionAdTagsLoadMore.page = 1
    optionAdTagsLoadMore.isSearch = true
    LoadMoreData(optionAdTagsLoadMore)
}

function UpdateListSelected(id, name, optionTarget) {
    let currentIdx = optionTarget.list_selected.map(item => {
        return item.id
    }).indexOf(id)
    if (currentIdx === -1) {
        optionTarget.list_selected.push({ id: id, name: name })
        DisplaySelected(id, name, optionTarget)
    } else {
        optionTarget.list_selected.splice(currentIdx, 1)
    }
    DisplayTextSelected(optionTarget)
    checkBoxSelectedEmpty(optionTarget)
}

function RemoveInclude(element, event, optionTarget) {
    event.preventDefault();
    let id = parseInt(element.attr("id"))
    let name = element.attr("name")
    element.closest(".item_selected").remove()
    // let div = $(optionTarget.idBox).find(`div[id ='${id}']`)
    // div.remove()

    AppendList(id, name, optionTarget)
    UpdateListSelected(id, name, optionTarget)
    switch (optionTarget.optionAjax) {
        case optionDomainsLoadMore.optionAjax:
            let currentIdxInventory = filterAjax.inventory.indexOf(id)
            if (currentIdxInventory > -1) {
                filterAjax.inventory.splice(currentIdxInventory, 1)
            }
            optionAdTagsLoadMore.page = 1
            optionAdTagsLoadMore.isSearch = true
            LoadMoreData(optionAdTagsLoadMore)
            break

        case optionAdFormatsLoadMore.optionAjax:
            let currentIdxFormat = filterAjax.format.indexOf(id)
            if (currentIdxFormat > -1) {
                filterAjax.format.splice(currentIdxFormat, 1)
            }
            optionAdTagsLoadMore.page = 1
            optionAdTagsLoadMore.isSearch = true
            LoadMoreData(optionAdTagsLoadMore)
            break
        case optionAdSizesLoadMore.optionAjax:
            let currentIdxSize = filterAjax.size.indexOf(id)
            if (currentIdxSize > -1) {
                filterAjax.size.splice(currentIdxSize, 1)
            }
            optionAdTagsLoadMore.page = 1
            optionAdTagsLoadMore.isSearch = true
            LoadMoreData(optionAdTagsLoadMore)
            break
        case optionAdTagsLoadMore.optionAjax:
            let currentIdxAdTag = optionAdTagsLoadMore.filterAjax.indexOf(id)
            if (currentIdxAdTag > -1) {
                optionAdTagsLoadMore.filterAjax.splice(currentIdxAdTag, 1)
            }
            // FilterAdTagChange(optionAdTagsLoadMore,"remove")
            break
    }
}

function DisplayNoDataAvailable(optionTarget, total, page) {
    if (total < 7 && page === 1) {

    } else if (optionTarget.lastPage || (total > 7 || total === 0)) {
        $("." + optionTarget.boxEmpty).remove()
        $(optionTarget.idSelect).append(`<div class="list-group list-group-flush my-n3 pt-3 ${optionTarget.boxEmpty}">
        <div class="list-group-item">
            <div id="${optionTarget.idEmptyLoad}" class="d-flex flex-row align-items-center px-md-2">
                <div class="col p-0 text-center"><h6 class="m-0 font-weight-semibold fs-12">No data available</h6></div>
            </div>
        </div>
    </div>`)
    }
}

// function AppendList(id, name, option) {
//     $(option.idSelect).append(`<div class="target-item border-bottom border-gray-50 ms-2 me-3" id="${id}">
//         <div class="list-group list-group-flush my-n3">
//             <div class="list-group-item {{if eq $row.Status.Int 3}}bg-gray-200{{end}}">
//                 <div class="d-flex flex-row align-items-center">
//                     <div class="col p-0">
//                         <span>${name}</span>
//                     </div>
//                     <div class="col-auto">
//                         <button type="button" data-id="${id}" name="${name}"
//                                 class="btn d-flex align-items-center btn-outline-secondary btn-icon rounded-circle p-0 ${option.btnInclude}">
//                             <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor"
//                                  class="bi bi-plus-lg" viewBox="0 0 16 16">
//                                 <path d="M8 0a1 1 0 0 1 1 1v6h6a1 1 0 1 1 0 2H9v6a1 1 0 1 1-2 0V9H1a1 1 0 0 1 0-2h6V1a1 1 0 0 1 1-1z"/>
//                             </svg>
//                         </button>
//                     </div>
//                 </div>
//             </div>
//         </div>
//     </div>`)
// }
function AppendList(id, name, option) {
    $(option.idSelect).append(`
    <div class="dm20 target-item" id="${id}">
        <span>${name}</span>
        <button class="${option.btnInclude}" data-id="${id}" name="${name}" type="button">
            <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor" class="bi bi-plus-lg" viewBox="0 0 16 16">
                <path d="M8 0a1 1 0 0 1 1 1v6h6a1 1 0 1 1 0 2H9v6a1 1 0 1 1-2 0V9H1a1 1 0 0 1 0-2h6V1a1 1 0 0 1 1-1z"></path>
            </svg>
        </button>
    </div>`)
}

// function DisplaySelected(id, name, optionTarget) {
//     $(optionTarget.idBox).append(`
//     <div class="target-item border-bottom border-gray-50 ms-2 me-3 item_selected" id="${id}">
//         <div class="list-group list-group-flush my-n3">
//             <div class="list-group-item">
//                 <div class="d-flex flex-row align-items-center">
//                     <div class="col p-0">
//                         <span>${name}</span>
//                     </div>
//                     <div class="col-auto">
//                     <button type="button" id="${id}" name="${name}" class="btn d-flex align-items-center btn-outline-danger btn-icon rounded-circle p-0 ${optionTarget.btnRemove}">
//                         <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor" class="bi bi-dash-lg" viewBox="0 0 16 16">
//                         <path d="M0 8a1 1 0 0 1 1-1h14a1 1 0 1 1 0 2H1a1 1 0 0 1-1-1z"></path>
//                         </svg>
//                      </button>
//                     </div>
//                 </div>
//             </div>
//         </div>
//     </div>`)
// }

function DisplaySelected(id, name, optionTarget) {
    $(optionTarget.idBox).append(`
        <div class="dm20 target-item item_selected" id="${id}">
            <span>${name}</span>
            <button class="${optionTarget.btnRemove}" id="${id}" name="${name}" type="button">
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor" class="bi bi-dash-lg" viewBox="0 0 16 16">
                    <path d="M0 8a1 1 0 0 1 1-1h14a1 1 0 1 1 0 2H1a1 1 0 0 1-1-1z"></path>
                </svg>
            </button>
        </div>`)
}

function HideSelected(optionTarget) {
    optionTarget.list_selected.map(item => {
        let divCountry = $(optionTarget.idSelect).find(`div[id = '${item.id}']`)
        divCountry.remove()
    })
}

function RemoveNotInFilter(optionTarget, listFilter) {
    switch (optionTarget.optionAjax) {
        case optionAdTagsLoadMore.optionAjax:
            let listTagRemove = []
            optionAdTagsLoadMore.list_selected.map((item, index) => {
                let indexSelected = listFilter.indexOf(item.id)
                if (indexSelected === -1) {
                    let div = $(optionAdTagsLoadMore.idBox).find(`div[id = '${item.id}']`)
                    div.remove()
                    listTagRemove.push(index)
                }
            })
            for (let i = listTagRemove.length - 1; i >= 0; i--) {
                optionAdTagsLoadMore.list_selected.splice(listTagRemove[i], 1)
            }

            let listTagFilterRemove = []
            $.each(optionAdTagsLoadMore.filterAjax, function (i, v) {
                let currentIdxAdTag = listFilter.indexOf(v)
                if (currentIdxAdTag === -1) {
                    listTagFilterRemove.push(i)
                }
            })
            for (let i = listTagFilterRemove.length - 1; i >= 0; i--) {
                optionAdTagsLoadMore.filterAjax.splice(listTagFilterRemove[i], 1)
            }

            checkBoxSelectedEmpty(optionAdTagsLoadMore)
            break
    }
}

function DisplayTextSelected(optionTarget) {
    const str = optionTarget.list_selected.map(item => {
        return ` ${item.name}`
    }).join();
    const lenString = str.length;
    if (lenString > 40) {
        $(optionTarget.text).prev().attr("hidden", true)
        const newStr = str.substring(0, 39);
        $(optionTarget.text).html(newStr + "....")
    } else if (lenString > 0) {
        $(optionTarget.text).prev().attr("hidden", true)
        $(optionTarget.text).html(str)
    } else {
        $(optionTarget.text).prev().removeAttr('hidden')
        switch (optionTarget) {
            case optionDomainsLoadMore:
                // $(optionTarget.text).html("all domains")
                $(optionTarget.text).html(`Choose at least one domain`)
                break
            case optionAdFormatsLoadMore:
                $(optionTarget.text).html("all formats")
                break
            case optionAdSizesLoadMore:
                $(optionTarget.text).html("all sizes")
                break
            case optionAdTagsLoadMore:
                $(optionTarget.text).html("all ad tags")
                break
            case optionGeographyLoadMore:
                $(optionTarget.text).html("all geographies")
                break
            case optionDeviceLoadMore:
                $(optionTarget.text).html("all devices")
                break
        }

    }
}

function HandleCollapse() {
    let url = "/line-item/collapse"
    $('#AddBidder').on('show.bs.collapse', '.collapse', function (e) {
        let box = e.target.id
        SendRequestShow(url, box, 0, "add")
    });
    $('#AddBidder').on('hide.bs.collapse', '.collapse', function (e) {
        let box = e.target.id
        SendRequestHide(url, box, 0, "add")
    });
}

function SendRequestHide(url, box, id, type) {
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

function SendRequestShow(url, box, id, type) {
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