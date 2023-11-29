$(document).ready(function () {
    new SubmitForm("register", Callback);
});

function Callback(response, formElement) {
    switch (response.status) {
        case "error":
            if (response.errors.length === 1 && response.errors[0].id === "") {
                new AlertError(response.errors[0].message);
            } else if (response.errors.length > 0) {
                // new AlertError(response.errors[0].message);
                $.each(response.errors, function (key, value) {
                    let inputElement = $("input#" + value.id);
                    if (key === 0) {
                        inputElement.select().focus();
                    }
                    inputElement.addClass("is-invalid").nextAll("div.invalid-feedback").text(value.message);
                });
            }
            break
        case "success":
            var queryString = window.location.search;
            if ( queryString.search("reddit") != -1 && queryString.search("utm_campain") != -1 ){
                rdt('track', 'SignUp');
            }
            let buttonElement = formElement.find(".submit");
            buttonElement.addClass("btn-success is-valid")
                .attr('disabled', true).text("Congrats!")
                .next(".valid-feedback").text("is redirecting to the home page...");
            setTimeout(function () {
                let backURL = formElement.find("#BackURL").val();
                new Redirect(backURL);
            }, 800);
            break
        default:
            AlertError("undefined");
            break
    }
}