let idx = 0;
$(document).ready(function () {
    $('[data-bs-toggle="popover"]').popover({
        html: true,
        sanitize: false,
    });
});

$(document).ready(function () {
    $('#summernote').summernote({
        toolbar: [
            ['style', ['style']],
            ['font', ['bold', 'underline', 'clear']],
            ['fontname', ['fontname']],
            ['color', ['color']],
            ['para', ['ul', 'ol', 'paragraph']],
            ['table', ['table']],
            ['insert', ['link', 'picture']],
            ['view', ['fullscreen', 'codeview', 'help']],
        ],
    });

    // change message when change  Dialog Template
    $('.adb-modal-fixed').change(function () {
        var value = $(this).val();
        if (value == 'modal') {
            var message = '<h1 style="margin-bottom: 15px; padding: 0 30px; color: #252b37; font-size: 28px; line-height: 1.25; text-align: center;">Adblocker detected! Please consider reading this notice.</h1>'
                + '<p>We' + "'" + 've detected that you are using AdBlock Plus or some other adblocking software which is preventing the page from fully loading.</p>'
                + '<p>We don' + "'" + 't have any banner, Flash, animation, obnoxious sound, or popup ad. We do not implement these annoying types of ads!</p>'
                + '<p>We need money to operate the site, and almost all of it comes from our online advertising.</p>'
                + '<p> <strong> Please add <a href="http://yourdomain.com" target="_blank" rel="noopener"> yourdomain.com </a> to your ad blocking whitelist or disable your adblocking software.</strong></p>';
            $('#summernote').summernote('code', message);
        }
        if (value == 'fixed') {
            var message = '<p>Our website is made possible by displaying online advertisements to our visitors.</p>'
                + '<p>Please consider supporting us by disabling your ad blocker.</p>';
            $('#summernote').summernote('code', message);
        }
    })

    $('.generator').click(function () {
        var adb_domain = $('#adb-domain').val();
        var adb_content = $('#summernote').summernote('code');
        var adb_modal_fixed = $('#adb-modal-fixed').val();
        var adb_display_time = $('#adb-display-time').val();
        var adb_hide_close_button = $('#adb-hide-close-button').val();
        var adb_close_background = $('#adb-close-background').val();

        if (!adb_domain) {
            new notifiWarning("Domain is required!");
            return;
        }
        if (!adb_content) {
            new notifiWarning("Message is required!");
            return;
        }

        var template = '';
        if (adb_modal_fixed == 'modal') {
            template = '<style> #__vliadb83 {' +
                '        display: none;' +
                '        position: fixed;' +
                '        background: rgb(221, 221, 221);' +
                '        z-index: 9999999;' +
                '        opacity: 1;' +
                '        visibility: visible;' +
                '        top: 100px;' +
                '        right: 0px;' +
                '        left: 0px;' +
                '        max-width: 640px;' +
                '        margin-right: auto;' +
                '        margin-left: auto;' +
                '        box-shadow: rgba(0, 0, 0, 0.25) 0px 3px 5px 2px;' +
                '        font-family: Arial, Helvetica, sans-serif;' +
                '    }' +
                '    #__vliadb83 .__vliadb83-content {' +
                '        padding: 30px 30px 15px;' +
                '    }' +
                '    #__vliadb83 #__vliadb83-cls {' +
                '        display: inline-block;' +
                '        position: absolute;' +
                '        top: 15px;' +
                '        right: 15px;' +
                '        width: 30px;' +
                '        height: 30px;' +
                '        color: #bbb;' +
                '        font-size: 32px;' +
                '        font-weight: 700;' +
                '        line-height: 30px;' +
                '        text-align: center;' +
                '        cursor: pointer;' +
                '        -webkit-transition: 0.3s;' +
                '        transition: 0.3s;' +
                '    }' +
                '    #__vliadb83 #__vliadb83-cls:hover {' +
                '        color: #5f5e5e;' +
                '    }' +
                '    #__vliadb83-bg {' +
                '        display: none;' +
                '        position: fixed;' +
                '        z-index: 999999;' +
                '        background: rgba(0, 0, 0, 0.8);' +
                '        top: 0px;' +
                '        left: 0px;' +
                '        width: 100%;' +
                '        height: 100%;' +
                '    } </style>' +
                '<div id="__vliadb83">' +
                '<div class="__vliadb83-content">' +
                adb_content +
                '</div>' +
                '<a id="__vliadb83-cls">×</a>' +
                '</div>' +
                '<div id="__vliadb83-bg"></div>' +
                '<script src="//cdn.jsdelivr.net/gh/vli-platform/adb-analytics@29f6e17/v1.0.min.js"></script>' +
                '<script>' +
                '(function () {' +
                '(window.adblockDetector = window.adblockDetector || []).push(function () {' +
                'window.adbDetector.init({' +
                'id: "' + adb_domain + '",' +
                'debug: true,' +
                'cookieExpire: ' + adb_display_time + ',' +
                'found: function () {' +
                'window.adbDetector.alert({' +
                'hiddenCloseButton: ' + adb_hide_close_button + ',' +
                'clickBackgroundToClose: ' + adb_close_background +
                '});' +
                '}' +
                '});' +
                '})' +
                '}());' +
                '</script>';
        }

        if (adb_modal_fixed == 'fixed') {
            template = '<style>' +
                '#__vliadb83 {' +
                'position: fixed;' +
                'top: 0;' +
                'left: 0;' +
                'width: 100%;' +
                ' padding: 10px 0;' +
                'text-align: center;' +
                'background: #eca6a6;' +
                'color: #b31717;' +
                'font-size: 15px;' +
                'font-weight: 700;' +
                'display: none;' +
                'z-index: 99999;' +
                '}' +
                '#__vliadb83-cls {' +
                'border: none;' +
                'display: inline-block;' +
                'padding: 4px 8px;' +
                'border-radius: 7px;' +
                '-moz-border-radius: 7px;' +
                '-webkit-border-radius: 7px;' +
                'margin-left: 5px;' +
                'background: none;' +
                'cursor: pointer;' +
                'color: #fff;' +
                'font-size: 12px;' +
                'font-weight: bold;' +
                'position: absolute;' +
                'top: -5px;' +
                'right: 3px;' +
                '}' +
                '</style>' +
                '<div id="__vliadb83">' +
                '<div style="position:relative">' +
                '<div>' +
                adb_content +
                '</div>' +
                '<button id="__vliadb83-cls">' +
                '<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" id="Capa_1" x="0px" y="0px" width="12px" height="12px" viewBox="0 0 348.333 348.334" style="enable-background:new 0 0 348.333 348.334;" xml:space="preserve" class="">' +
                '<g>' +
                '<g>' +
                '<path d="M336.559,68.611L231.016,174.165l105.543,105.549c15.699,15.705,15.699,41.145,0,56.85   c-7.844,7.844-18.128,11.769-28.407,11.769c-10.296,0-20.581-3.919-28.419-11.769L174.167,231.003L68.609,336.563   c-7.843,7.844-18.128,11.769-28.416,11.769c-10.285,0-20.563-3.919-28.413-11.769c-15.699-15.698-15.699-41.139,0-56.85   l105.54-105.549L11.774,68.611c-15.699-15.699-15.699-41.145,0-56.844c15.696-15.687,41.127-15.687,56.829,0l105.563,105.554   L279.721,11.767c15.705-15.687,41.139-15.687,56.832,0C352.258,27.466,352.258,52.912,336.559,68.611z"' +
                'data-original="#000000" class="active-path" data-old_color="#000000" fill="#3e3d3d" />' +
                '</g>' +
                '</g>' +
                '</svg>' +
                '</button>' +
                '</div>' +
                '</div>' +
                '<script src="//cdn.jsdelivr.net/gh/vli-platform/adb-analytics@29f6e17/v1.0.min.js"></script>' +
                '<script>' +
                '(function () {' +
                '(window.adblockDetector = window.adblockDetector || []).push(function () {' +
                'window.adbDetector.init({' +
                'id: ' + adb_domain + ',' +
                'debug: true,' +
                'cookieExpire: ' + adb_display_time + ',' +
                'found: function () {' +
                'window.adbDetector.alert({' +
                'hiddenCloseButton: ' + adb_hide_close_button + ',' +
                'clickBackgroundToClose: ' + adb_close_background +
                '});' +
                '}' +
                '});' +
                '})' +
                '}());' +
                '</script>';
        }
        if (template != "") {
             template = template.replace(/\s+/g, " ");
            $(".result-adblock-analytics").val(template)
            // $('#summernote').summernote('code', template);

        }


//         $.ajax({
//             type: 'POST',
//             url: '/Adblock?action=ajax_generator',
//             data: {
//                 adb_domain: adb_domain,
//                 adb_content: adb_content,
//                 adb_modal_fixed: adb_modal_fixed,
//                 adb_display_time: adb_display_time,
//                 adb_hide_close_button: adb_hide_close_button,
//                 adb_close_background: adb_close_background,
//             }
//         }).done(function (result) {
//             $("#generator-manager").find(".BG-blur").hide()
//             console.log(result);
//             if (result.error) {
// //                    alert(result.error);
//             } else {
//                 $('.result-adblock-analytics').val(result);
//             }
//         })
    })


});