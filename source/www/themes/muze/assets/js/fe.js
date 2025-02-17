$(document).ready(function () {
	
});

GetAccountManagerInfo()

function GetAccountManagerInfo() {
	var token = getCookie('mcflgi');
	$.ajax({
		type: 'GET',
		// url: 'http://127.0.0.1:8542/api/get-account-manager?q=pp',
		url: '/api/get-account-manager?q=pp',
		headers: {
			'Token': token,
		},
		success: function (result) {
		},
		error: function (jqXHR, exception) {
			const msg = AjaxErrorMessage(jqXHR, exception);
			// NoticeError("AJAX ERROR: " + msg);
		}
	})
		.done(function (result) {
			// console.log(result);
			// if (result.error) {
			// 	alert(result.error);
			// } else {
			$(".account-manager").html(result)
			// }
		})
}

function getCookie(name) {
	const value = `; ${document.cookie}`;
	const parts = value.split(`; ${name}=`);
	if (parts.length === 2) return parts.pop().split(';').shift();
}