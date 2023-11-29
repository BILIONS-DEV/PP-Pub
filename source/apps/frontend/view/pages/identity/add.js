let filterAjax = {
    inventory: [],
    format: [],
    size: [],
};

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
};

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
};

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
};

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
};

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
};

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
};

let list_select_module = [];
let moduleParams = [];
let moduleParamsDefault = [];

$(document).ready(function () {
    $('[data-bs-toggle="popover"]').popover({
        html: true,
        sanitize: false,
    });
    $(".identity-tab").on("click", ".nav-link", function () {
        $(".identity-tab").find(".nav-link").removeClass("pp-4");
        $(this).addClass("pp-4");
        var tab = $(this).attr("data-tab");
        if (tab != "1") {
            $(this).addClass("at-1");
        } else {
            $(".identity-tab").find('.at-1').removeClass('at-1');
        }
    })

    // Target
    $('[data-toggle="collapse"]').on("click", function () {
        var element = $(this);

        element.toggleClass("dm14");
        if (element.hasClass("dm14")) {
            // element.collapse('show')
            element.find(".dm23").find("button").html(
                '<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor" class="bi bi-dash-lg" viewBox="0 0 16 16">\n' +
                '<path d="M0 8a1 1 0 0 1 1-1h14a1 1 0 1 1 0 2H1a1 1 0 0 1-1-1z">\n' +
                '</path>\n' +
                '</svg>');
        } else {
            element.find(".dm23").find("button").html(
                '<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor" class="bi bi-plus-lg" viewBox="0 0 16 16">\n' +
                '<path d="M8 0a1 1 0 0 1 1 1v6h6a1 1 0 1 1 0 2H9v6a1 1 0 1 1-2 0V9H1a1 1 0 0 1 0-2h6V1a1 1 0 0 1 1-1z"></path>\n' +
                '</svg>');
        }

        $("#nav-target").find('[data-toggle="collapse"]').each(function () {
            if (element[0] != $(this)[0]) {
                $(this).removeClass("dm14");
                $(this).find(".dm23").find("button").html(
                    '<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor" class="bi bi-plus-lg" viewBox="0 0 16 16">\n' +
                    '<path d="M8 0a1 1 0 0 1 1 1v6h6a1 1 0 1 1 0 2H9v6a1 1 0 1 1-2 0V9H1a1 1 0 0 1 0-2h6V1a1 1 0 0 1 1-1z"></path>\n' +
                    '</svg>');
            }
        });
    });

    //Submit
    new SubmitFormIdentity("submitIdentity", Added);
    EventClickCollapseTarget();
    HandleModule();
    AddModule();
});

function AddModule() {
    $("#tab-module").on("click", ".add-module", function (e) {
        $("#tab-module").find(".box-select-module").removeClass("d-none");
        InitAfterLoadModule();
    });
}

function HandleModule() {
    $('#select-module').on('changed.bs.select', function (e) {
        var param = $(this).val();
        if (!param) {
            return;
        }

        let value = {};
        value.id = $(this).val();
        value.text = $('option:selected', this).text();
        // let value = e.params.data;
        let currentIdx = list_select_module.map(item => {
            return item.id;
        }).lastIndexOf(value.id);
        if (currentIdx > -1) {
        } else {
            list_select_module.push({id: value.id, index: 1});
            LoadModuleParam(value.id, value.text, 1);
        }
        $("#tab-module").find(".box-select-module").addClass("d-none");
    });

    $("#module-params").on("click", "div.bidder-box .d a.rm_c", function () {
        const module_id = $(this).attr("data-id");
        RemoveModule(module_id);
    });
}

function RemoveModule(moduleID) {
    $("#ModuleParamBox-" + moduleID).remove();
    const modules = [];
    $("#module-params").find('.module_param_box').each(function () {
        modules.push($(this).attr('data-id'));
    });
    $('#module').val(modules);
    $('#module').trigger('change');
}

function RemoveModuleParams(id) {
    list_select_module = list_select_module.filter((item) => {
        return !(item.id === id.toString());
    });
    let list = document.querySelectorAll(`#ModuleParamBox-${id}`);
    list.forEach(item => {
        item.remove();
    });
}

function SelectModuleUserId(e) {
    const data = e.params.data;
    let currentIdx = list_select_module.map(item => {
        return item.id;
    }).lastIndexOf(data.id);
    if (currentIdx > -1) {
        let index = list_select_module[currentIdx].index + 1;
        list_select_module.push({id: data.id, index: index});
        LoadModuleParam(data.id, data.text, index);
    } else {
        list_select_module.push({id: data.id, index: 1});
        LoadModuleParam(data.id, data.text, 1);
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
            $("#module-params").prepend(json.data);
            InitAfterLoadModule();
        },
        error: function (xhr) {
            console.log(xhr);
        }
    });
}

function InitAfterLoadModule() {
    $('#select-module option').prop("selected", false).trigger('change');
    $("#select-module").selectpicker('refresh');
}

function EventClickCollapseTarget() {
    $('.btn-collapse').click(function (e) {
        if ($(this).hasClass('collapsed')) {
            $(this).find(".text-notify-target").removeClass("d-none");
        } else {
            $(".text-notify-target").removeClass("d-none");
            $(this).find(".text-notify-target").addClass("d-none");
        }
    });
}

function addParamModule() {
    moduleParams = []
    $('#module-params > div.module_param_box').each(function () {
        if ($(this).find(".bidder-box").attr('data-name') === "criteo") {
            let currentId = $(this).find(".bidder-box").attr('data-id')
            let bidderParam = {}
            bidderParam.id = parseInt(currentId)
            bidderParam.name = $(this).find(".bidder-box").attr('data-name')
            bidderParam.storage = []
            bidderParam.params = []
            moduleParams.push(bidderParam)
        } else {
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
        }
    });
    return Object.values(moduleParams);
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

function Added(response, formElement) {
    switch (response.status) {
        case "error":
            if (response.errors.length === 1 && response.errors[0].id === "") {
                new AlertError(response.errors[0].message);
            } else if (response.errors.length > 0) {
                $.each(response.errors, function (key, value) {
                    let inputElement = $("#" + value.id);
                    inputElement.addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
                    inputElement.next().find(".select2-selection").addClass("select2-is-invalid")
                    inputElement.closest('.box-select2').addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message)
                    if (value.id === "text_for_domain") {
                        let msg = value.message
                        inputElement.html(`<small class="text-danger">${msg}</small>`)
                        $('.domain-card').addClass("domain-card-invalid")
                    }
                    if (value.id === "text_for_domain") {
                        $("#tab-target").click()
                        let msg = value.message
                        inputElement.html(`<small class="text-danger">${msg}</small>`)
                        $('.domain-card').addClass("domain-card-invalid")
                        $("#collapseDomains").collapse('show')
                        $('div[href="#collapseDomains"]').addClass("dm14")
                    }
                    NoticeError(value.message)
                });
                $("#" + response.errors[0].id).focus();
                $("#" + response.errors[0].id).prev('label').focus();
            } else {
                new AlertError("Error!");
            }
            break
        case "success":
            NoticeSuccess("Identity has been created successfully")
            setTimeout(function () {
                window.location.replace("/identity");
            }, 1000);
            break
        default:
            new AlertError("Undefined");
            break
    }
}

function SubmitFormIdentity(formID, functionCallback, ajxURL = "") {
    let formElement = $("#" + formID);
    let button = formElement.find(".submit");
    let submitButtonText = button.text();
    let submitButtonTextLoading = "Loading...";
    let acceptChangeStatus = false;
    formElement.find(".select2").on("select2:open", function (e) {
        let selectElement = $(this)
        if (selectElement.next().find(".select2-selection").hasClass("select2-is-invalid")) {
            selectElement.next().find(".select2-selection").removeClass("select2-is-invalid")
            selectElement.removeClass("is-invalid").next(".invalid-feedback").empty();
            selectElement.closest('.box-select2').next(".invalid-feedback").empty();
        }
    });
    formElement.find("input").on("input", function (e) {
        let inputElement = $(this)
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    $("#module-params").on("input", ".module_param_box .param_value", function (e) {
        let inputElement = $(this)
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    formElement.find("textarea").on("input", function (e) {
        let inputElement = $(this)
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    $(formElement).on('keypress', "input", function (e) {
        const keyCode = e.keyCode || e.which;
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
        postData()
        // let profileTargetAllInventory = formElement.data("profile_target_all_inventory")
        // let isHaveAProfile = formElement.data("have-profile")
        //
        // // trường hợp target all domain
        // if (optionDomainsLoadMore.list_selected.length === 0) {
        //     if (isHaveAProfile && !acceptChangeStatus) {
        //         swal({
        //             title: "Are you sure?",
        //             text: "Each domain can only target 1 corresponding profile. When you choose to target the entire domain, the active profiles will be disabled, you will have to re-enable them manually when needed. Would you like to do this?",
        //             icon: "warning",
        //             buttons: true,
        //             dangerMode: true,
        //         }).then((willDelete) => {
        //             if (willDelete) {
        //                 postData()
        //                 acceptChangeStatus = true
        //             }
        //         });
        //     } else {
        //         postData()
        //     }
        // } else {
        //     if (profileTargetAllInventory && !acceptChangeStatus) {
        //         swal({
        //             title: "Are you sure?",
        //             text: "This Identity Profile is targeting all domains. Because each domain can only use 1 Identity Profile, so you need to turn off all other Identity Profiles. Would you like to do this?",
        //             icon: "warning",
        //             buttons: true,
        //             dangerMode: true,
        //         }).then((willDelete) => {
        //             if (willDelete) {
        //                 postData()
        //                 acceptChangeStatus = true
        //             }
        //         });
        //     } else {
        //         postData()
        //     }
        // }
    }

    function postData() {
        let postData = formElement.serializeObject();
        postData.listInventory = optionDomainsLoadMore.list_selected
        postData.priority = parseInt(postData.priority)
        postData.sync_delay = parseInt(postData.sync_delay)
        postData.auction_delay = parseInt(postData.auction_delay)
        if (postData.status === "on") {
            postData.status = 1
        } else {
            postData.status = 2
        }
        postData.module_params = addParamModule()
        const data = JSON.stringify(postData);
        const buttonElement = $(".submit");
        $.ajax({
            url: ajxURL,
            type: "POST",
            dataType: "JSON",
            contentType: "application/json",
            data: data,
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
            complete: function (res) {
                functionCallback(res.responseJSON, formElement);
            }
        });
    }
}

//Target
$(document).ready(function () {
    //Load more domains
    setUpLoadMoreData(optionDomainsLoadMore)
});

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
    // $(optionTarget.idSelect).on("click", "div.list-group div.list-group-item div.col-auto button." + optionTarget.btnInclude, function (e) {
    $(optionTarget.idSelect).on("click", "button." + optionTarget.btnInclude, function (e) {
        ClickInclude($(this), e, optionTarget)
        CheckAddAll1Page(optionTarget)
        RemoveNotifyNoData(optionTarget)

        //Xóa báo lỗi card domain
        if ($(this).hasClass("add_inventory")) {
            $(".domain-card").removeClass("domain-card-invalid")
        }
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
        removeAll(optionTarget)
    })

    // $(optionTarget.idBox).on("click", ".list-group .list-group-item .flex-row .col-auto button." + optionTarget.btnRemove, function (e) {
    $(optionTarget.idBox).on("click", "button." + optionTarget.btnRemove, function (e) {
        RemoveInclude($(this), e, optionTarget)
        RemoveNotifyNoData(optionTarget)
    })
}

function getListChecked(optionTarget) {
    let divs = $(optionTarget.idBox).find("." + optionTarget.btnRemove)
    divs.each(function (index, elm) {
        let id = parseInt(elm.id)
        let name = elm.name
        optionTarget.list_selected.push({id: id, name: name})
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
            optionTarget.list_selected.push({id: item.id, name: item.name})
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
            optionTarget.list_selected.push({id: item.id, name: item.name})
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
            optionTarget.list_selected.push({id: item.id, name: item.name})
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
        optionTarget.list_selected.push({id: id, name: name})
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
                $(optionTarget.text).html("all domains")
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