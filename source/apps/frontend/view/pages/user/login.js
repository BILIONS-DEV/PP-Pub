$(document).ready(function () {
    new SubmitForm("login", Callback);
});

function Callback(response, formElement) {
    switch (response.status) {
        case "error":
            if (response.errors.length === 1 && response.errors[0].id === "") {
                AlertError(response.errors[0].message);
            } else if (response.errors.length > 0) {
                $.each(response.errors, function (key, value) {
                    let inputElement = $("input#" + value.id);
                    if (key === 0) {
                        inputElement.select().focus();
                    }
                    inputElement.addClass("is-invalid").nextAll("div.invalid-feedback").text(value.message);
                });
            }
            break;
        case "success":
            let buttonElement = formElement.find(".submit");
            buttonElement.addClass("btn-success is-valid")
                .attr('disabled', true).text("Congrats!")
                .next(".valid-feedback").text("is redirecting to the home page...");
            setTimeout(function () {
                let backURL = formElement.find("#BackURL").val();
                new Redirect(backURL);
            }, 800);
            break;
        default:
            AlertError("undefined");
            break;
    }
}