// const func = require("./select_content")
var listContentSelected = []

$(document).ready(function () {
    // func.firstLoad();
    checkContentSelected()
    HandleCollapse()
    SubmitFormCreate("editPlaylist", Added, "/playlist/edit");
    // $.fn.select2.defaults.set("theme", "bootstrap");
    // $('.select2').select2({
        // theme: "bootstrap-5",
        // selectionCssClass: "select2", // For Select2 v4.1
        // dropdownCssClass: "select2",
        // closeOnSelect: false
        // style: "text-transform: capitalize"
    // });

    $("div > #content").on("select2:select", function (e) {
        selectContent(e, $(this))
    });

    $(".playlist-tab").on("click", ".nav-link", function () {
        $(".playlist-tab").find(".nav-link").removeClass("pp-4")
        $(this).addClass("pp-4")
        var tab = $(this).attr("data-tab")
        if (tab != "1") {
            $(this).addClass("at-1")
        } else {
            $(".playlist-tab").find('.at-1').removeClass('at-1')
        }
    })

    // Target
    handleTarget()

    // $("#show-content-select").on("click", ".content .remove-content", function () {
    //     removeContent($(this))
    // })

    checkedIncludeExcludeConfig();
    // LoadHistory();
})

function handleTarget() {
    $('[data-toggle="collapse"]').on("click", function () {
        var element = $(this)
        console.log($(this));

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


$(document).ready(function () {
    //Load more language
    setUpLoadMoreData(optionlanguageLoadMore)

    //Load more channels
    setUpLoadMoreData(optionChannelsLoadMore)

    //Load more category
    setUpLoadMoreData(optionCategoryLoadMore)

    //Load more keywords
    setUpLoadMoreData(optionKeywordsLoadMore)

    //Load more videos
    setUpLoadMoreData(optionVideosLoadMore)
});

function setUpLoadMoreData(optionTarget) {
    // firstLoad(optionTarget)
    getListChecked(optionTarget)
    firstLoad(optionTarget)
    checkBoxSelectedEmpty(optionTarget)
    DisplayTextSelected(optionTarget)
    setUpClickAndSearch(optionTarget)
    scrollToBottom(optionTarget)
}

function checkedIncludeExcludeConfig() {
    $('input[type="radio"]').on("change", function () {
        getFilterAjax()
        if ($(this).hasClass(optionlanguageLoadMore.idType.replace(".", ""))) {
            reloadData(optionChannelsLoadMore)
            reloadData(optionCategoryLoadMore)
            reloadData(optionKeywordsLoadMore)
            reloadData(optionVideosLoadMore)
        }
        if ($(this).hasClass(optionChannelsLoadMore.idType.replace(".", ""))) {
            reloadData(optionCategoryLoadMore)
            reloadData(optionKeywordsLoadMore)
            reloadData(optionVideosLoadMore)
        }
        if ($(this).hasClass(optionKeywordsLoadMore.idType.replace(".", ""))) {
            reloadData(optionVideosLoadMore)
        }
    })

    var option = [
        optionlanguageLoadMore,
        optionChannelsLoadMore,
        optionCategoryLoadMore,
        optionKeywordsLoadMore,
        optionVideosLoadMore
    ]
    option.forEach(function(optionTarget){
        if ($(optionTarget.idBox).find("." + optionTarget.btnRemove).length == 0) {
            $(optionTarget.idType + '[value="1"]').prop("checked", true)
            $(optionTarget.idType).closest("div").addClass("cursor-not-allowed").find("label").addClass("disabled")
        }
    })
}

function getFilterAjax() {
    // optionlanguageLoadMore
    let divslanguage = $(optionlanguageLoadMore.idBox).find("." + optionlanguageLoadMore.btnRemove)
    var type = parseInt($(optionlanguageLoadMore.idType + ":checked").val())
    divslanguage.each(function (index, elm) {
        var language_id = parseInt(elm.id)
        if (type == 1) {
            if (!filterAjax.language.includes(language_id)) {
                filterAjax.language.push(language_id)
            }
            if (filterAjax.exclude_language.indexOf(language_id) > -1) {
                filterAjax.exclude_language.splice(filterAjax.exclude_language.indexOf(language_id))
            }
        } else {
            if (!filterAjax.exclude_language.includes(language_id)) {
                filterAjax.exclude_language.push(language_id)
            }
            if (filterAjax.language.indexOf(language_id) > -1) {
                filterAjax.language.splice(filterAjax.language.indexOf(language_id))
            }
        }
    })

    // optionChannelsLoadMore
    let divsChannels = $(optionChannelsLoadMore.idBox).find("." + optionChannelsLoadMore.btnRemove)
    var type = parseInt($(optionChannelsLoadMore.idType + ":checked").val())
    divsChannels.each(function (index, elm) {
        var channels_id = parseInt(elm.id)
        if (type == 1) {
            if (!filterAjax.channels.includes(channels_id)) {
                filterAjax.channels.push(channels_id)
            }
            if (filterAjax.exclude_channels.indexOf(channels_id) > -1) {
                filterAjax.exclude_channels.splice(filterAjax.exclude_channels.indexOf(channels_id))
            }
        } else {
            if (!filterAjax.exclude_channels.includes(channels_id)) {
                filterAjax.exclude_channels.push(channels_id)
            }
            if (filterAjax.channels.indexOf(channels_id) > -1) {
                filterAjax.channels.splice(filterAjax.channels.indexOf(channels_id))
            }
        }
    })

    // optionlanguageLoadMore
    let divsKeywords = $(optionKeywordsLoadMore.idBox).find("." + optionKeywordsLoadMore.btnRemove)
    var type = parseInt($(optionKeywordsLoadMore.idType + ":checked").val())
    divsKeywords.each(function (index, elm) {
        var keyword_id = parseInt(elm.id)
        if (type == 1) {
            if (!filterAjax.keywords.includes(keyword_id)) {
                filterAjax.keywords.push(keyword_id)
            }
            if (filterAjax.exclude_keywords.indexOf(keyword_id) > -1) {
                filterAjax.exclude_keywords.splice(filterAjax.exclude_keywords.indexOf(keyword_id))
            }
        } else {
            if (!filterAjax.exclude_keywords.includes(keyword_id)) {
                filterAjax.exclude_keywords.push(keyword_id)
            }
            if (filterAjax.keywords.indexOf(keyword_id) > -1) {
                filterAjax.keywords.splice(filterAjax.keywords.indexOf(keyword_id))
            }
        }
    })
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
        // $("." + optionTarget.btnRemoveAll).attr("hidden", false)
        let text = "#" + optionTarget.idEmpty
        $(text).remove()
    }
}

function setUpClickAndSearch(optionTarget) {
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

    $(optionTarget.idBox).on("click", "button." + optionTarget.btnRemove, function (e) {
        RemoveInclude($(this), e, optionTarget)
        // RemoveNotifyNoData(optionTarget)
    })
}

function getListChecked(optionTarget) {
    let divs = $(optionTarget.idBox).find("." + optionTarget.btnRemove)
    var type = parseInt($(optionTarget.idType + ":checked").val())

    divs.each(function (index, elm) {
        let id = parseInt(elm.id)
        let name = elm.name
        optionTarget.list_selected.push({id: id, name: name})
        switch (optionTarget.optionAjax) {
            case optionlanguageLoadMore.optionAjax:
                if (type == 1) {
                    if (!filterAjax.language.includes(id)) {
                        filterAjax.language.push(id)
                    }
                } else {
                    if (!filterAjax.exclude_language.includes(id)) {
                        filterAjax.exclude_language.push(id)
                    }
                }
                break
            case optionChannelsLoadMore.optionAjax:
                if (type == 1) {
                    if (!filterAjax.channels.includes(id)) {
                        filterAjax.channels.push(id)
                    }
                } else {
                    if (!filterAjax.exclude_channels.includes(id)) {
                        filterAjax.exclude_channels.push(id)
                    }
                }
                break
            case optionKeywordsLoadMore.optionAjax:
                if (type == 1) {
                    if (!filterAjax.keywords.includes(id)) {
                        filterAjax.keywords.push(id)
                    }
                } else {
                    if (!filterAjax.exclude_keywords.includes(id)) {
                        filterAjax.exclude_keywords.push(id)
                    }
                }
                break
        }
    })
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
    if (currentPage == 1) {
        $(optionTarget.idSelect).html("")
    }
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
            if (currentPage == 1) {
                $(optionTarget.idSelect).html("")
            }
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
            DisplayTextSelected(optionTarget)
            optionTarget.checkLoadMore = false
            optionTarget.lastPage = json.last_page
            DisplayNoDataAvailable(optionTarget, json.total, json.current_page)
            if (optionTarget.idSelect == optionVideosLoadMore.idSelect) {
                if (json.total_all && json.total_all > 0) {
                    $("#total_video").removeClass("d-none")
                    $(".number-v").text($("#list_videos").find(".target-item").length);
                    $(".total-v").text(json.total_all);
                    $(".box_empty_video").addClass("d-none")
                } else {
                    $("#total_video").addClass("d-none")
                    $(".box_empty_video").removeClass("d-none")
                }
            }
        },
        error: function (xhr) {
            console.log(xhr)
        }
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

function RemoveNotifyNoData(optionTarget) {
    let a_elm = $(optionTarget.idSelect).find("div.target-item")
    let a_empty = $("." + optionTarget.boxEmpty)
    if (a_elm.length > 0) {
        if (a_empty.length > 0) {
            a_empty.attr("hidden", true)
        } else {
            $(optionTarget.idSelect).append(`<div class="list-group list-group-flush my-n3 pt-3 ${optionTarget.boxEmpty}" hidden>
        <div class="list-group-item">
            <div class="d-flex flex-row align-items-center px-md-2">
                <div class="col p-0 text-center"><h6 class="m-0 font-weight-semibold fs-12">No data available</h6></div>
            </div>
        </div>
    </div>`)
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

function removeAll(optionTarget) {
    optionTarget.list_selected.map(item => {
        AppendList(item.id, item.name, optionTarget)
    })

    optionTarget.list_selected = []
    reloadData(optionTarget)

    switch (optionTarget.idSelect) {
        case optionlanguageLoadMore.idSelect:
            filterAjax.language = []
            filterAjax.exclude_language = []
            reloadData(optionChannelsLoadMore)
            reloadData(optionCategoryLoadMore)
            reloadData(optionKeywordsLoadMore)
            reloadData(optionVideosLoadMore)
            break
        case optionChannelsLoadMore.idSelect:
            filterAjax.channels = []
            filterAjax.exclude_channels = []
            reloadData(optionCategoryLoadMore)
            reloadData(optionKeywordsLoadMore)
            reloadData(optionVideosLoadMore)
            break
        case optionKeywordsLoadMore.idSelect:
            filterAjax.keywords = []
            filterAjax.exclude_keywords = []
            reloadData(optionVideosLoadMore)
            break
    }
    // $(optionTarget.idSelect).html("")
    // optionTarget.page = 1
    // optionTarget.isSearch = false
    // optionTarget.isMoreData = true
    // LoadMoreData(optionTarget)
    // checkBoxSelectedEmpty(optionTarget)
    // DisplayTextSelected(optionTarget)
}

function reloadData(optionTarget) {
    $(optionTarget.idSelect).html("")
    $(optionTarget.idBox).html("")
    $(optionTarget.idType + '[value="1"]').prop("checked", true)
    optionTarget.page = 1
    optionTarget.isSearch = false
    optionTarget.isMoreData = true
    optionTarget.list_selected = []
    LoadMoreData(optionTarget)
    checkBoxSelectedEmpty(optionTarget)
    DisplayTextSelected(optionTarget)
    if ($(optionTarget.idBox).find("." + optionTarget.btnRemove).length == 0) {
        $(optionTarget.idType + '[value="1"]').prop("checked", true)
        $(optionTarget.idType).closest("div").addClass("cursor-not-allowed").find("label").addClass("disabled")
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
            case optionlanguageLoadMore:
                $(optionTarget.text).html("all language")
                break
            case optionChannelsLoadMore:
                $(optionTarget.text).html("all channels")
                break
            case optionCategoryLoadMore:
                $(optionTarget.text).html("all categories")
                break
            case optionVideosLoadMore:
                $(optionTarget.text).html("all videos")
                break
            case optionKeywordsLoadMore:
                $(optionTarget.text).html("all keywords")
                break
        }

    }
}

function ClickInclude(element, event, optionTarget) {
    event.preventDefault();
    let id = parseInt(element.data("id"))
    let name = element.attr("name")
    let div = $(optionTarget.idSelect).find(`div[id = '${id}']`)
    // div.attr("hidden", true)
    div.remove()
    UpdateListSelected(id, name, optionTarget)
    UpdateFilterAjax(id, -1, optionTarget)

    if ($(optionTarget.idBox).find("." + optionTarget.btnRemove).length > 0) {
        $(optionTarget.idType).closest("div").removeClass("cursor-not-allowed").find("label").removeClass("disabled")
    }
    if (optionTarget.optionAjax === optionlanguageLoadMore.optionAjax) {
        reloadData(optionChannelsLoadMore)
        reloadData(optionCategoryLoadMore)
        reloadData(optionKeywordsLoadMore)
        reloadData(optionVideosLoadMore)
    }
    if (optionTarget.optionAjax === optionChannelsLoadMore.optionAjax) {
        reloadData(optionCategoryLoadMore)
        reloadData(optionKeywordsLoadMore)
        reloadData(optionVideosLoadMore)
    }
    if (optionTarget.optionAjax === optionKeywordsLoadMore.optionAjax) {
        reloadData(optionVideosLoadMore)
    }
}

function RemoveInclude(element, event, optionTarget) {
    event.preventDefault();
    let id = parseInt(element.attr("id"))
    let name = element.attr("name")
    let div = $(optionTarget.idBox).find(`div[id = '${id}']`)
    div.remove()
    // AppendList(id, name, optionTarget)
    UpdateListSelected(id, name, optionTarget)
    if ($(optionTarget.idBox).find("." + optionTarget.btnRemove).length == 0) {
        $(optionTarget.idType + '[value="1"]').prop("checked", true)
        $(optionTarget.idType).closest("div").addClass("cursor-not-allowed").find("label").addClass("disabled")
    }

    var type = parseInt($(optionTarget.idType + ":checked").val())
    switch (optionTarget.optionAjax) {
        case optionlanguageLoadMore.optionAjax:
            if (type == 1) {
                let currentIdLanguage = filterAjax.language.indexOf(id)
                if (currentIdLanguage > -1) {
                    filterAjax.language.splice(currentIdLanguage, 1)
                }
            } else {
                let currentIdLanguage = filterAjax.exclude_language.indexOf(id)
                if (currentIdLanguage > -1) {
                    filterAjax.exclude_language.splice(currentIdLanguage, 1)
                }
            }
            optionlanguageLoadMore.page = 1
            optionlanguageLoadMore.isSearch = true
            LoadMoreData(optionlanguageLoadMore)
            reloadData(optionChannelsLoadMore)
            reloadData(optionCategoryLoadMore)
            reloadData(optionKeywordsLoadMore)
            reloadData(optionVideosLoadMore)
            break

        case optionChannelsLoadMore.optionAjax:
            if (type == 1) {
                let currentIdChannels = filterAjax.channels.indexOf(id)
                if (currentIdChannels > -1) {
                    filterAjax.channels.splice(currentIdChannels, 1)
                }
            } else {
                let currentIdChannels = filterAjax.exclude_channels.indexOf(id)
                if (currentIdChannels > -1) {
                    filterAjax.exclude_channels.splice(currentIdChannels, 1)
                }
            }
            optionChannelsLoadMore.page = 1
            optionChannelsLoadMore.isSearch = true
            LoadMoreData(optionChannelsLoadMore)
            reloadData(optionCategoryLoadMore)
            reloadData(optionKeywordsLoadMore)
            reloadData(optionVideosLoadMore)
            break

        case optionKeywordsLoadMore.optionAjax:
            if (type == 1) {
                let currentIdKeywords = filterAjax.keywords.indexOf(id)
                if (currentIdKeywords > -1) {
                    filterAjax.keywords.splice(currentIdKeywords, 1)
                }
            } else {
                let currentIdKeywords = filterAjax.exclude_keywords.indexOf(id)
                if (currentIdKeywords > -1) {
                    filterAjax.exclude_keywords.splice(currentIdKeywords, 1)
                }
            }
            optionKeywordsLoadMore.page = 1
            optionKeywordsLoadMore.isSearch = true
            LoadMoreData(optionKeywordsLoadMore)
            reloadData(optionVideosLoadMore)
            break
        default:
            LoadMoreData(optionTarget)
            break

    }
}

function CheckIfScrollBottom(scrollTop, clientHeight, scrollHeight = 0) {
    if ((scrollTop + clientHeight + 1) >= scrollHeight) {
        return true
    } else {
        return false
    }
}

function HideSelected(optionTarget) {
    optionTarget.list_selected.map(item => {
        let divCountry = $(optionTarget.idSelect).find(`div[id = '${item.id}']`)
        divCountry.remove()
    })
}

function DisplayNoDataAvailable(optionTarget, total, page) {
    if (total < 30 && page === 1) {

    } else if (optionTarget.lastPage || (total > 30 || total === 0)) {
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

function DisplaySelected(id, name, optionTarget) {
    $(optionTarget.idBox).append(`<div class="dm20 target-item item_selected" id="${id}">
            <span>${name}</span>
            <button class="${optionTarget.btnRemove}" id="${id}" name="${name}" type="button">
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor" class="bi bi-dash-lg" viewBox="0 0 16 16">
                    <path d="M0 8a1 1 0 0 1 1-1h14a1 1 0 1 1 0 2H1a1 1 0 0 1-1-1z"></path>
                </svg>
            </button>
        </div>`)
}

function AppendList(id, name, option) {
    $(option.idSelect).append(`<div class="dm20 target-item" id="${id}">
        <span>${name}</span>
        <button class="${option.btnInclude}" data-id="${id}" name="${name}" type="button">
            <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor" class="bi bi-plus-lg" viewBox="0 0 16 16">
                <path d="M8 0a1 1 0 0 1 1 1v6h6a1 1 0 1 1 0 2H9v6a1 1 0 1 1-2 0V9H1a1 1 0 0 1 0-2h6V1a1 1 0 0 1 1-1z"></path>
            </svg>
        </button>
    </div>`)
}

function UpdateFilterAjax(id, currentIdx, optionTarget) {
    var type = parseInt($(optionTarget.idType + ":checked").val())
    switch (optionTarget.optionAjax) {
        case optionlanguageLoadMore.optionAjax:
            if (type == 1) {
                if (currentIdx === -1) {
                    filterAjax.language.push(id)
                }
            } else {
                if (currentIdx === -1) {
                    filterAjax.exclude_language.push(id)
                }
            }
            break
        case optionChannelsLoadMore.optionAjax:
            if (type == 1) {
                if (currentIdx === -1) {
                    filterAjax.channels.push(id)
                }
            } else {
                if (currentIdx === -1) {
                    filterAjax.exclude_channels.push(id)
                }
            }
            break
        case optionKeywordsLoadMore.optionAjax:
            if (type == 1) {
                if (currentIdx === -1) {
                    filterAjax.keywords.push(id)
                }
            } else {
                if (currentIdx === -1) {
                    filterAjax.exclude_keywords.push(id)
                }
            }
            break
    }
}

function HandleCollapse() {
    let url = "/playlist/collapse"
    $('#createPlaylist').on('show.bs.collapse', '.collapse', function (e) {
        let box = e.target.id
        SendRequestShow(url, box, 0, "add")
    });
    $('#createPlaylist').on('hide.bs.collapse', '.collapse', function (e) {
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
        }
    }).done(function (result) {
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
        }
    }).done(function (result) {
    });
}

function checkContentSelected() {
    $("#show-content-select .content").each(function () {
        let id = $(this).attr('data-content-id')
        listContentSelected.push(id)
    })
    hiddenSelect2()
}

function hiddenSelect2(){
    $("#content").find("option").removeAttr("disabled")
    $("#content").find("option").each(function(){
        var value = $(this).attr("value")
        if ( listContentSelected.includes(value) ){
            $(this).attr("disabled","disabled")
        }
    })
    // $('#content').select2();
}

function selectContent(e, select) {
    var data = e.params.data;
    let isCheck = (listContentSelected.indexOf(data.id) > -1)
    if (isCheck) {
        return
    }
    let description = $('#content option:selected').attr("data-desc")
    let thumb = $('#content option:selected').attr("data-thumb")
    listContentSelected.push(data.id)
    select.parent().append('<input type="hidden" name="content" value="${data.id}">')
    func.appendContentSelect(data.id, data.text, thumb)
    hiddenSelect2()
}

// function removeContent(el) {
//     let idRemove = el.attr("data-content-id")
//     let indexRemove = listContentSelected.indexOf(idRemove)
//     if (indexRemove > -1) {
//         listContentSelected.splice(indexRemove, 1);
//     }
//     el.closest('.content').remove()
//     func.refreshListContent();
//     hiddenSelect2()
// }

function SubmitFormCreate(formID, functionCallback, ajxURL = "") {
    let formElement = $("#" + formID);
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
    formElement.find(".select2").on("select2:open", function (e) {
        let selectElement = $(this)
        if (selectElement.next().find(".select2-selection").hasClass("select2-is-invalid")) {
            selectElement.next().find(".select2-selection").removeClass("select2-is-invalid")
            selectElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        }
    });
    formElement.find(".submit").on("click", function (e) {
        e.preventDefault();
        const buttonElement = $(this);
        const submitButtonText = buttonElement.text();
        const submitButtonTextLoading = "Loading...";
        var postData = formElement.serializeObject();
        postData.videos_limit = parseInt(postData.videos_limit)
        postData.listLanguage = optionlanguageLoadMore.list_selected
        postData.listChannels = optionChannelsLoadMore.list_selected
        postData.listCategory = optionCategoryLoadMore.list_selected
        postData.listKeywords = optionKeywordsLoadMore.list_selected
        postData.listVideos = optionVideosLoadMore.list_selected
        postData.type_category = parseInt(postData.type_category)
        postData.type_channels = parseInt(postData.type_channels)
        postData.type_keywords = parseInt(postData.type_keywords)
        postData.type_language = parseInt(postData.type_language)
        postData.type_videos = parseInt(postData.type_videos)
        postData.id = parseInt(postData.id)
        // postData.user_id = parseInt(postData.user_id)

        $.ajax({
            url: ajxURL,
            type: "POST",
            dataType: "JSON",
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
            complete: function (res) {
                functionCallback(res.responseJSON, formElement);
            }
        });
    });
}

function Added(response) {
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
                    inputElement.next().find(".select2-selection").addClass("select2-is-invalid")
                });
                $("#" + response.errors[0].id).focus();
                $("#" + response.errors[0].id).prev('label').focus();
                // new AlertError(response.errors[0].message, function () {
                //     $("#" + response.errors[0].id).focus();
                //     $("#" + response.errors[0].id).prev('label').focus();
                // })
            } else {
                new AlertError("Error!");
            }
            break
        case "success":
            NoticeSuccess("Playlist has been updated successfully")
            break
        default:
            new AlertError("Undefined");
            break
    }
}

let filterAjax = {
    language: [],
    exclude_language: [],
    channels: [],
    exclude_channels: [],
    keywords: [],
    exclude_keywords: [],
}

let optionlanguageLoadMore = {
    idSelect: "#list_language",
    idBox: ".box_language",
    idType: ".type_language",
    idIsMoreData: "#is_more_language",
    isMoreData: false,
    lastPage: false,
    page: 2,
    isSearch: false,
    search: "",
    idSearch: "#search_language",
    urlAjax: "/target/loadPlaylist",
    optionAjax: "language",
    filterAjax: [],
    list_selected: [],
    btnInclude: "add_language",
    btnRemove: "remove_language",
    idEmpty: "empty_language",
    idEmptyLoad: "load_empty_language",
    checkLoadMore: false,
    text: "#text_for_language",
    btnRemoveAll: "remove_all_language",
    btnSelectAll: "select_all_language",
    block: ".block_language",
    container: ".container_language",
    boxEmpty: "box_empty_language"
}
let optionChannelsLoadMore = {
    idSelect: "#list_channels",
    idBox: ".box_channels",
    idType: ".type_channels",
    idIsMoreData: "#is_more_channels",
    isMoreData: false,
    lastPage: false,
    page: 2,
    isSearch: false,
    search: "",
    idSearch: "#search_channels",
    urlAjax: "/target/loadPlaylist",
    optionAjax: "channels",
    filterAjax: [],
    list_selected: [],
    btnInclude: "add_channels",
    btnRemove: "remove_channels",
    idEmpty: "empty_channels",
    idEmptyLoad: "load_empty_channels",
    checkLoadMore: false,
    text: "#text_for_channels",
    btnRemoveAll: "remove_all_channels",
    btnSelectAll: "select_all_channels",
    block: ".block_channels",
    container: ".container_channels",
    boxEmpty: "box_empty_channels"
}
let optionCategoryLoadMore = {
    idSelect: "#list_category",
    idBox: ".box_category",
    idType: ".type_category",
    idIsMoreData: "#is_more_category",
    isMoreData: false,
    lastPage: false,
    page: 2,
    isSearch: false,
    search: "",
    idSearch: "#search_category",
    urlAjax: "/target/loadPlaylist",
    optionAjax: "category",
    filterAjax: [],
    list_selected: [],
    btnInclude: "add_category",
    btnRemove: "remove_category",
    idEmpty: "empty_category",
    text: "#text_for_category",
    btnRemoveAll: "remove_all_category",
    btnSelectAll: "select_all_category",
    block: ".block_category",
    container: ".container_category",
    boxEmpty: "box_empty_category"
}
let optionKeywordsLoadMore = {
    idSelect: "#list_keywords",
    idBox: ".box_keywords",
    idType: ".type_keywords",
    idIsMoreData: "#is_more_keywords",
    isMoreData: false,
    lastPage: false,
    page: 2,
    isSearch: false,
    search: "",
    idSearch: "#search_keywords",
    urlAjax: "/target/loadPlaylist",
    optionAjax: "keywords",
    filterAjax: [],
    list_selected: [],
    btnInclude: "add_keywords",
    btnRemove: "remove_keywords",
    idEmpty: "empty_keywords",
    idEmptyLoad: "load_empty_keywords",
    checkLoadMore: false,
    text: "#text_for_keywords",
    btnRemoveAll: "remove_all_keywords",
    btnSelectAll: "select_all_keywords",
    block: ".block_keywords",
    container: ".container_keywords",
    boxEmpty: "box_empty_keywords"
}
let optionVideosLoadMore = {
    idSelect: "#list_videos",
    idBox: ".box_videos",
    idType: ".type_videos",
    idIsMoreData: "#is_more_videos",
    isMoreData: false,
    lastPage: false,
    page: 2,
    isSearch: false,
    search: "",
    idSearch: "#search_videos",
    urlAjax: "/target/loadPlaylist",
    optionAjax: "videos",
    list_selected: [],
    btnInclude: "add_videos",
    btnRemove: "remove_videos",
    idEmpty: "empty_videos",
    idEmptyLoad: "load_empty_videos",
    checkLoadMore: false,
    text: "#text_for_videos",
    btnRemoveAll: "remove_all_videos",
    btnSelectAll: "select_all_videos",
    block: ".block_videos",
    container: ".container_videos",
    boxEmpty: "box_empty_videos"
}