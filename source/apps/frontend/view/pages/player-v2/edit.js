import general from "./general";
const act = new general({ ajxURL: "/player/template/edit", isEdit: true });

act.init();

$(document).ready(function () {
    $('[data-bs-toggle="popover"]').popover({
        html: true, sanitize: false,
    });
});