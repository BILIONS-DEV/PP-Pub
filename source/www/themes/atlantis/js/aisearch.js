jQuery(document).ready(function ($) {
    var ai_search = $('#ai-search');
    var ai_input = ai_search.find('input');
    document.onkeyup = function (e) {
        var key = e.which || e.keyCode;
        if (e.ctrlKey && key == 32) {
            var status = ai_search.attr('status');
            status == 'on' ? (ai_search.hide(), ai_search.attr('status', 'off')) : (ai_search.show(), ai_input.focus(), ai_search.attr('status', 'on'));
        }
    };
})