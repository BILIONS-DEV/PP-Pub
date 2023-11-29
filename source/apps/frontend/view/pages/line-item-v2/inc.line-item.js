const ajax = require("../../jspkg/ajax");

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
    boxEmpty: "box_empty_format"
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
    boxEmpty: "box_empty_size"
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

module.exports = {
    optionDomainsLoadMore,
    optionAdFormatsLoadMore,
    optionAdSizesLoadMore,
    optionGeographyLoadMore,
    optionDeviceLoadMore,
    SubmitForm,
}

function SubmitForm(formID, functionCallback) {
    let formElement = $("#" + formID);
    // Xử lý Select2
    formElement.find(".select2").on("select2:open", function (e) {
        let selectElement = $(this)
        if (selectElement.next().find(".select2-selection").hasClass("select2-is-invalid")) {
            selectElement.next().find(".select2-selection").removeClass("select2-is-invalid")
            selectElement.removeClass("is-invalid").next(".invalid-feedback").empty();
            selectElement.closest('.box-select2').next(".invalid-feedback").empty();
        }
    });
    //
    formElement.find("input").on("input", function (e) {
        let inputElement = $(this)
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    $("#bidder-params").on("input", ".list_param .param_value", function (e) {
        let inputElement = $(this)
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    // Nếu bấm vào input thì sẽ xóa báo lỗi
    formElement.find("textarea").on("input", function (e) {
        let inputElement = $(this)
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    // xử lý khi bấm enter
    $(formElement).on('keypress', "input", function (e) {
        const keyCode = e.keyCode || e.which;
        if (keyCode === 13) {
            e.preventDefault();
            submit()
        }
    });
    // Xử lý khi bấm Save
    formElement.find(".submit").on("click", function (e) {
        e.preventDefault();
        $.each(formElement.find("input.is-invalid"), function () {
            $(this).removeClass("is-invalid").next(".invalid-feedback").empty();
        });
        submit()
    });

    function submit() {
        let postData = formElement.serializeObject();
        const data = JSON.stringify(makePostData(postData))
        ajax.Post(formElement, data, functionCallback)
    }

    function makePostData(postData) {
        postData.listInventory = optionDomainsLoadMore.list_selected
        postData.listGeo = optionGeographyLoadMore.list_selected
        postData.listDevice = optionDeviceLoadMore.list_selected
        postData.listAdtag = optionAdTagsLoadMore.list_selected
        postData.listAdFormat = optionAdFormatsLoadMore.list_selected
        postData.listAdSize = optionAdSizesLoadMore.list_selected
        postData.priority = parseInt(postData.priority)
        postData.type = parseInt(postData.type)
        postData.linked_gam = parseInt(postData.linked_gam)
        postData.server_type = parseInt($("#server_type").val())
        postData.rate = parseInt(postData.rate)
        postData.id = parseInt(postData.id)
        if (!Array.isArray(postData.select_account)) {
            const selectAccount = [];
            if (postData.select_account !== "") {
                selectAccount.push(postData.select_account)
            }
            postData.select_account = selectAccount
        }
        postData.bidder_params = addParamBidder()
        postData.adsense_ad_slots = addAdsenseAdSlot()
        return postData
    }

    function addAdsenseAdSlot() {
        let adsenseAdSlots = []
        $("#adsense_ad_slot_item > div.box-c > div.bidder-box").each(function () {
            let adSlot = {}
            $("input", this).each(function () {
                let size = $(this).data("size")
                let adSlotId = $(this).val().trim()
                adSlot.size = size
                adSlot.ad_slot_id = adSlotId
            })
            adsenseAdSlots.push(adSlot)
        })
        return adsenseAdSlots
    }

    function addParamBidder() {
        let bidderParams = []
        $('#bidder-params > div.box-c > div.bidder-box').each(function () {
            let params = {}
            $('input', this).each(function () {
                let type = $(this).attr('data-type').toLowerCase()
                let name = $(this).attr('data-name')
                params[name] = this.value.trim()
            })
            let bidderParam = {}
            bidderParam["id"] = parseInt($(this).attr('data-id'))
            bidderParam["name"] = $(this).attr('data-name').toLowerCase()
            bidderParam["bidder_type"] = parseInt($(this).find(".type-bidder").val())
            bidderParam["bidder_index"] = parseInt($(this).attr('data-index'))
            bidderParam["params"] = params
            bidderParams.push(bidderParam)
        });
        return bidderParams
    }
}