module.exports = {
    Redirect: function (url) {
        window.location.replace(url);  //=> similar behavior as an HTTP redirect
        window.location.href = url;  //=> similar behavior as clicking on a link
    },
    IsFunction: function (functionToCheck) {
        return functionToCheck && {}.toString.call(functionToCheck) === '[object Function]';
    },
    CopyTextToClipboard: function (text) {
        const textArea = document.createElement("textarea");
        // Place in the top-left corner of screen regardless of scroll position.
        textArea.style.position = 'fixed';
        textArea.style.top = 0;
        textArea.style.left = 0;
        // Ensure it has a small width and height. Setting to 1px / 1em
        // doesn't work as this gives a negative w/h on some browsers.
        textArea.style.width = '2em';
        textArea.style.height = '2em';
        // We don't need padding, reducing the size if it does flash render.
        textArea.style.padding = 0;
        // Clean up any borders.
        textArea.style.border = 'none';
        textArea.style.outline = 'none';
        textArea.style.boxShadow = 'none';
        // Avoid flash of the white box if rendered for any reason.
        textArea.style.background = 'transparent';
        textArea.value = text;
        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();
        let successful;
        try {
            successful = document.execCommand('copy');
            let msg = successful ? 'successful' : 'unsuccessful';
            // console.log('Copying text command was ' + msg);
        } catch (err) {
            console.log('Oops, unable to copy');
        }
        document.body.removeChild(textArea);
        return successful;
    }
}