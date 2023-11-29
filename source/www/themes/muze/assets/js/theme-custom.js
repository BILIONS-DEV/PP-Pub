//Muze Onhover Get Inline CSS JavaScript
var muzeNavItem = document.querySelectorAll('.navbar-vertical .nav-item').forEach((muzeNavItem) => {
    muzeNavItem.addEventListener('mouseover', (e) => {
        var muzePosition = muzeNavItem.getBoundingClientRect();
        muzeNavItem.style.top = muzePosition.top + 'px';
    });
});

//Muze Dropdown Stop Propagation JavaScript
// document.querySelectorAll('.dropdown-menu, .dropdown-toggle').forEach(dropdownMenu => {
//   dropdownMenu.addEventListener('click', (e) => {
//     e.stopPropagation();
//   });
// });

//Muze Toggle Class JavaScript
document.querySelectorAll('.back-arrow').forEach(muzeSidebarSwitcher => {
    muzeSidebarSwitcher.addEventListener('click', (e) => {
        if (document.body.className.match('sidebar-compact')) {
            document.querySelector('body').classList.toggle('sidebar-compact');
            document.querySelector('body').classList.add('previous-compact');
        } else if (document.body.className.match('previous-compact')) {
            document.querySelector('body').classList.toggle('sidebar-compact');
        } else {
            document.querySelector('body').classList.toggle('sidebar-icons');
        }
    });
});
document.querySelectorAll('.muze-hamburger').forEach(muzeHamburger => {
    muzeHamburger.addEventListener('click', (e) => {
        document.querySelector('body').classList.toggle('sidebar-menu');
    });
});
document.querySelectorAll('body').forEach(muzeHamburger => {
    muzeHamburger.addEventListener('click', (e) => {
        if (document.querySelector('.customize-sidebar') && !document.querySelector('.customize-sidebar').contains(e.target)) {
            document.querySelector('body').classList.remove('customize-box');
        }
    });
});
document.querySelectorAll('.customize-btn, .customize-close').forEach(muzeCustomizerToggle => {
    muzeCustomizerToggle.addEventListener('click', (e) => {
        e.stopPropagation();
        document.querySelector('body').classList.toggle('customize-box');
        if (localStorage.getItem('selectedSkin')) {
            document.getElementById(localStorage.getItem('selectedSkin')).checked = true;
        }
        var headerStyleStorage = localStorage.getItem('headerStyle');
        if (headerStyleStorage != '' && headerStyleStorage != null) {
            document.getElementById(headerStyleStorage).checked = true;
        }
        var sidebarStyleStorage = localStorage.getItem('sidebarStyle')
        if (sidebarStyleStorage != '' && sidebarStyleStorage != null) {
            document.getElementById(sidebarStyleStorage).checked = true;
        }
    });
});
document.querySelectorAll('.muze-search').forEach(muzeSearch => {
    muzeSearch.addEventListener('click', (e) => {
        document.querySelector('body').classList.toggle('search-open');
    });
});
document.querySelectorAll('.customize-sidebar').forEach(muzeSearch => {
    muzeSearch.addEventListener('click', (e) => {
        // e.stopPropagation();
    });
});

//Muze Set Item Value By Local Storage JavaScript
document.querySelectorAll('[name="MuzeSkins"]').forEach(skins => {
    skins.addEventListener('change', (e) => {
        localStorage.setItem('skins', skins.value);
    });
});
document.querySelectorAll('[name="HeaderStyles"]').forEach(headerStyles => {
    headerStyles.addEventListener('change', (e) => {
        localStorage.setItem('headerStyles', headerStyles.value);
        var sidebarCheckboxes = document.getElementsByName('SidebarStyles');
        Array.prototype.forEach.call(sidebarCheckboxes, function (sidebarCheckboxes) {
            sidebarCheckboxes.checked = false;
            sidebarCheckboxes.removeAttribute('checked');
        });
        localStorage.removeItem('sidebarStyles');
    });
});
document.querySelectorAll('[name="SidebarStyles"]').forEach(sidebarStyles => {
    sidebarStyles.addEventListener('change', (e) => {
        localStorage.setItem('sidebarStyles', sidebarStyles.value);
        var headerCheckboxes = document.getElementsByName('HeaderStyles');
        Array.prototype.forEach.call(headerCheckboxes, function (headerCheckboxes) {
            headerCheckboxes.checked = false;
            headerCheckboxes.removeAttribute('checked');
        });
        localStorage.removeItem('headerStyles');
    });
});
var customizerBody = document.querySelector('body').classList.add(localStorage.getItem('skins'), localStorage.getItem('headerStyles'), localStorage.getItem('sidebarStyles'));

//Muze Customizer Radio And Checkbox Selected JavaScript
if (localStorage.getItem('rtlModeCheck') == 'true') {
    document.getElementById('RTLMode').checked = true;
    var customizerHtml = document.querySelector('html').setAttribute('dir', 'rtl');
    var customizerBody = document.querySelector('body').classList.add(localStorage.getItem('rtlModeValue'));
    var rtlCss = document.createElement('link');
    rtlCss.href = '../assets/vendor/bootstrap/dist/css/bootstrap.rtl.min.css';
    rtlCss.rel = 'stylesheet';
    rtlCss.type = 'text/css';
    rtlCss.media = 'all';
    document.getElementsByTagName('head')[0].prepend(rtlCss);
}
if (localStorage.getItem('fluidLayoutCheck') == 'true') {
    document.getElementById('FluidLayout').checked = true;
    var customizerBody = document.querySelector('body').classList.remove(localStorage.getItem('fluidLayoutValue'));
} else {
    var customizerBody = document.querySelector('body').classList.add(localStorage.getItem('fluidLayoutValue'));
}

//Muze Get Selected Radio Button ID By Name Javascript
var getSelectedRadio = function getSelectedRadio(radioName) {
    var ele = document.getElementsByName(radioName);
    var value = '';
    for (i = 0; i < ele.length; i++) {
        if (ele[i].checked) {
            value = ele[i].id;
            break;
        }
    }
    return value;
}

//Muze Reload Page And Set item By LocalStorage JavaScript
// document.querySelector('#CustomizerPreview').addEventListener('click', event => {
//   var skin = getSelectedRadio('MuzeSkins');
//   localStorage.setItem('selectedSkin', skin);
//   var rtlMode = document.getElementById('RTLMode');
//   localStorage.setItem('rtlModeCheck', rtlMode.checked);
//   localStorage.setItem('rtlModeValue', rtlMode.value);
//   var headerStyle = getSelectedRadio('HeaderStyles');
//   localStorage.setItem('headerStyle', headerStyle);
//   var sidebarStyle = getSelectedRadio('SidebarStyles');
//   localStorage.setItem('sidebarStyle', sidebarStyle);
//   var fluidLayout = document.getElementById('FluidLayout');
//   localStorage.setItem('fluidLayoutCheck', fluidLayout.checked);
//   localStorage.setItem('fluidLayoutValue', fluidLayout.value);
//   window.location.reload();
// });

//Muze Remove Local Storage And Reload Page JavaScript
// document.querySelector('#ResetCustomizer').addEventListener('click', event => {
//   localStorage.removeItem('skins');
//   localStorage.removeItem('rtlModeCheck');
//   localStorage.removeItem('rtlModeValue');
//   localStorage.removeItem('headerStyle');
//   localStorage.removeItem('selectedSkin');
//   localStorage.removeItem('headerStyles');
//   localStorage.removeItem('sidebarStyle');
//   localStorage.removeItem('sidebarStyles');
//   localStorage.removeItem('fluidLayoutCheck');
//   localStorage.removeItem('fluidLayoutValue');
//   window.location.reload();
// });

//Muze Custom JavaScript
function _objectWithoutProperties(source, excluded) {
    if (source == null) return {};
    var target = _objectWithoutPropertiesLoose(source, excluded);
    var key, i;
    if (Object.getOwnPropertySymbols) {
        var sourceSymbolKeys = Object.getOwnPropertySymbols(source);
        for (i = 0; i < sourceSymbolKeys.length; i++) {
            key = sourceSymbolKeys[i];
            if (excluded.indexOf(key) >= 0) continue;
            if (!Object.prototype.propertyIsEnumerable.call(source, key)) continue;
            target[key] = source[key];
        }
    }
    return target;
}

function _objectWithoutPropertiesLoose(source, excluded) {
    if (source == null) return {};
    var target = {};
    var sourceKeys = Object.keys(source);
    var key, i;
    for (i = 0; i < sourceKeys.length; i++) {
        key = sourceKeys[i];
        if (excluded.indexOf(key) >= 0) continue;
        target[key] = source[key];
    }
    return target;
}

function ownKeys(object, enumerableOnly) {
    var keys = Object.keys(object);
    if (Object.getOwnPropertySymbols) {
        var symbols = Object.getOwnPropertySymbols(object);
        if (enumerableOnly) symbols = symbols.filter(function (sym) {
            return Object.getOwnPropertyDescriptor(object, sym).enumerable;
        });
        keys.push.apply(keys, symbols);
    }
    return keys;
}

function _objectSpread(target) {
    for (var i = 1; i < arguments.length; i++) {
        var source = arguments[i] != null ? arguments[i] : {};
        if (i % 2) {
            ownKeys(Object(source), true).forEach(function (key) {
                _defineProperty(target, key, source[key]);
            });
        } else if (Object.getOwnPropertyDescriptors) {
            Object.defineProperties(target, Object.getOwnPropertyDescriptors(source));
        } else {
            ownKeys(Object(source)).forEach(function (key) {
                Object.defineProperty(target, key, Object.getOwnPropertyDescriptor(source, key));
            });
        }
    }
    return target;
}

function _defineProperty(obj, key, value) {
    if (key in obj) {
        Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true });
    } else {
        obj[key] = value;
    }
    return obj;
}

var docReady = function docReady(fn) {
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', fn);
    } else {
        setTimeout(fn, 1);
    }
};

//Muze Resize JavaScript
var resize = function resize(fn) {
    return window.addEventListener('resize', fn);
};

//Muze Getdata JavaScript
var getData = function getData(el, data) {
    try {
        return JSON.parse(el.dataset[data]);
    } catch (e) {
        return el.dataset[data];
    }
};

var isScrolledIntoView = function isScrolledIntoView(el) {
    var top = el.offsetTop;
    var left = el.offsetLeft;
    var width = el.offsetWidth;
    var height = el.offsetHeight;

    while (el.offsetParent) {
        el = el.offsetParent;
        top += el.offsetTop;
        left += el.offsetLeft;
    }

    return {
        all: top >= window.pageYOffset && left >= window.pageXOffset && top + height <= window.pageYOffset + window.innerHeight && left + width <= window.pageXOffset + window.innerWidth,
        partial: top < window.pageYOffset + window.innerHeight && left < window.pageXOffset + window.innerWidth && top + height > window.pageYOffset && left + width > window.pageXOffset
    };
};

//Muze Color JavaScript
var hexToRgb = function hexToRgb(hexValue) {
    var hex;
    hexValue.indexOf('#') === 0 ? hex = hexValue.substring(1) : hex = hexValue; // Expand shorthand form (e.g. "03F") to full form (e.g. "0033FF")

    var shorthandRegex = /^#?([a-f\d])([a-f\d])([a-f\d])$/i;
    var result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex.replace(shorthandRegex, function (m, r, g, b) {
        return r + r + g + g + b + b;
    }));
    return result ? [parseInt(result[1], 16), parseInt(result[2], 16), parseInt(result[3], 16)] : null;
};
var rgbaColor = function rgbaColor() {
    var color = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : '#fff';
    var alpha = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : 0.5;
    return "rgba(".concat(hexToRgb(color), ", ").concat(alpha, ")");
};
var colors = {
    primary: '#2c7be5',
    secondary: '#748194',
    success: '#00d27a',
    info: '#27bcfd',
    warning: '#f5803e',
    danger: '#e63757',
    light: '#f9fafd',
    dark: '#0b1727'
};
var grays = {
    white: '#fff',
    100: '#f9fafd',
    200: '#edf2f9',
    300: '#d8e2ef',
    400: '#b6c1d2',
    500: '#9da9bb',
    600: '#748194',
    700: '#5e6e82',
    800: '#4d5969',
    900: '#344050',
    1000: '#232e3c',
    1100: '#0b1727',
    black: '#000'
};
var utils = {
    docReady: docReady,
    resize: resize,
    getData: getData,
    hexToRgb: hexToRgb,
    rgbaColor: rgbaColor,
    isScrolledIntoView: isScrolledIntoView,
    colors: colors,
    grays: grays,
};

//Muze Popover JavaScript
var popoverInit = function popoverInit() {
    var popoverTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="popover"]'));
    popoverTriggerList.map(function (popoverTriggerEl) {
        return new window.bootstrap.Popover(popoverTriggerEl, {
            // trigger: 'focus'
            html: true
        });
    });
    $('body').on('click', function (e) {
        $('[data-bs-toggle=popover]').each(function () {
            // hide any open popovers when the anywhere else in the body is clicked
            if (!$(this).is(e.target) && $(this).has(e.target).length === 0 && $('.popover').has(e.target).length === 0) {
                $(this).popover('hide');
            }
        });
    });
};


//Muze Toast JavaScript
var toastInit = function toastInit() {
    var toastElList = [].slice.call(document.querySelectorAll('.toast'));
    toastElList.map(function (toastEl) {
        return new window.bootstrap.Toast(toastEl);
    });
};

//Muze Tooltip JavaScript
var tooltipInit = function tooltipInit() {
    var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-toggle="tooltip"]'));
    tooltipTriggerList.map(function (tooltipTriggerEl) {
        return new window.bootstrap.Tooltip(tooltipTriggerEl);
    });
};


//Muze Theme Initialization
docReady(tooltipInit);
docReady(popoverInit);
docReady(toastInit);

$(document).click(function (event) {
    // if ($(event.target).closest('.box-notifi .dropdown-menu').length) {
    //     readAllNotification()
    // }
    $('.close-notification').click(function () {
        $("#notification").dropdown("hide")
    });
});

if (window.location.origin.search("127.0.0.1:8123") >= 0 || window.location.origin.search("apps.valueimpression.com") >= 0 || (window.location.origin.search("127.0.0.1:8542") >= 0 || window.location.origin.search("be.pubpower.io") >= 0)) {
    var myDropdown = document.getElementById('notification')
    if (myDropdown) {
        myDropdown.addEventListener('show.bs.dropdown', function () {
            // do something...
            readAllNotification()
        })
    }
}

$(function () {
    if (window.location.origin.search("127.0.0.1:8123") >= 0 || window.location.origin.search("apps.valueimpression.com") >= 0 || (window.location.origin.search("127.0.0.1:8542") >= 0 || window.location.origin.search("be.pubpower.io") >= 0)) {
        getNotification()
    }
    $(".toggle-collapse").click(function () {
        $(this).find('a[data-bs-toggle="collapse"]')[0].click()
    })
})

function getNotification() {
    $.ajax({
        type: 'GET',
        url: '/notification',
    })
        .done(function (result) {
            if (result.NumberNew > 0) {
                $("#number-notification").text(result.NumberNew).removeClass("d-none")
                // $("#not-new-notification").addClass("d-none")
            } else {
                $("#number-notification").addClass("d-none")
                // $("#not-new-notification").removeClass("d-none")
            }
            if (result.Template) {
                $("#result-notification").html(result.Template)
            }
        })
}

function readAllNotification() {
    $.ajax({
        type: 'GET',
        url: '/notification/ReadAll',
    })
        .done(function (result) {
            setTimeout(function () {
                $("#number-notification").remove()
            }, 300)
        })
}

function readNotification(e) {
    var id = e.attr("data-id")
    $.ajax({
        type: 'GET',
        url: '/notification/read',
        data: { id: id }
    })
        .done(function (result) {
            // console.log(result);
            getNotification()
        })
}

// align top left bottom right
function Loading() {
    return '<div class="_blur"><div class="spinner-border text-secondary" style="position: absolute; top: 45%" role="status"><span class="sr-only">Loading...</span></div></div>';
}

// center left right
function loading() {
    return '<div class="_blur"><div class="spinner-border text-secondary" role="status"><span class="sr-only">Loading...</span></div></div>';
}

btnResetCache();

function btnResetCache() {
    $("#btn_reset_cache").on("click", function () {
        $.ajax({
            url: "/rebuild-cache",
            type: "GET",
            dataType: "JSON",
            contentType: "application/json",
            data: "",
            beforeSend: function (xhr) {
            },
            error: function (jqXHR, exception) {
                const msg = AjaxErrorMessage(jqXHR, exception);
                new AlertError("AJAX ERROR: " + msg);
            },
            success: function (responseJSON) {
            },
            complete: function (res) {
                if (res.responseJSON.status === "success") {
                    NoticeSuccess("Reset Cache Success!")
                } else {
                    new AlertError("Reset Cache Fail!")
                }
            }
        });
    })
}

// ********** AI Search ***********
jQuery(document).ready(function ($) {
    var ai_search = $('#ai-search');
    var ai_input = ai_search.find('input');
    document.onkeyup = function (e) {
        var key = e.which || e.keyCode;
        if (e.ctrlKey && key == 32) {
            var status = ai_search.attr('status');
            status == 'on' ? (ai_search.hide(), ai_search.attr('status', 'off')) : (ai_search.show(), ai_input.focus(), ai_search.attr('status', 'on'));
        }
        if (key == 27) {
            var status = ai_search.attr('status');
            if (status == 'on') {
                ai_search.hide()
                ai_search.attr('status', 'off')
            }
        }
        if (key == 13) {
            e.preventDefault();
            var search = ai_input.val()
            if (!search) {
                return;
            }
            $.ajax({
                type: 'GET',
                url: '/AiSearch/query',
                data: { s: search }
            })
                .done(function (result) {
                    if (result.url) {
                        if (result.blank) {
                            window.open(result.url, '_blank');
                        } else {
                            window.location.href = result.url;
                        }
                    }
                })
        }
    };

    removeError();
    checkHistoryType()
})

function removeError() {
    $("body").on("click", "input", function (e) {
        let inputElement = $(this);
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
            // inputElement.removeClass("is-invalid").parent().removeClass("is-invalid").next(".invalid-feedback").empty();
            inputElement.closest(".is-invalid").removeClass("is-invalid").next(".invalid-feedback").empty().find(".invalid-feedback").empty();
        }
    });
    $("body").on("click", "textarea", function (e) {
        let inputElement = $(this);
        if (inputElement.hasClass("is-invalid")) {
            inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
            // inputElement.removeClass("is-invalid").parent().removeClass("is-invalid").next(".invalid-feedback").empty();
            inputElement.closest(".is-invalid").removeClass("is-invalid").next(".invalid-feedback").empty().find(".invalid-feedback").empty();
        }
    });
    $("body").find('.btn-group-radio').on('change', function (e) {
        let box = $(this).closest(".box-radio-btn");
        if (box.hasClass("is-invalid")) {
            box.removeClass("is-invalid").next(".invalid-feedback").empty();
            box.closest(".is-invalid").removeClass("is-invalid").next(".invalid-feedback").empty().find(".invalid-feedback").empty();
        }
    });
}

function checkHistoryType() {
    if ($(".load-history").length) {
        var objectType = $(".load-history").attr("data-object")
        var objectID = $(".load-history").attr("data-id")
        console.log(objectType);
        if (!objectType || !objectID) {
            $(".load-history").addClass("d-none")
        } else {
            $(".load-history").removeClass("d-none")
        }
    }
}

$(function () {
    $(".thu-menu").on("mouseover", ".dropdown", function () {
        $(".thu-menu").find(".dropdown").removeClass("show").find(".dropdown-menu").removeClass("show")
        $(this).addClass("show").find(".dropdown-menu").addClass("show")
    })

    document.onclick = function (event) {
        $(".thu-menu").find(".dropdown").removeClass("show").find(".dropdown-menu").removeClass("show")
    };
})