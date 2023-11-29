import general from "./general";
const act = new general();

act.init();

$(document).ready(function () {
    $('[data-bs-toggle="popover"]').popover({
        html: true, sanitize: false,
    });
});