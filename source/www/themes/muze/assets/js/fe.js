$(document).ready(function () {
	var accountManager = $('.account-manager');
	accountManager.on('click', '.copy-mail', function (event) {
		// Ngăn chặn hành động mặc định của liên kết (để không chuyển đến URL "#")
		event.preventDefault();
		// Lấy nội dung bên trong thẻ <div>
		var contentToCopy = accountManager.find("#email-manager").text();
		// Tạo một phần tử input ẩn để sao chép nội dung
		var inputElement = $('<input>').attr('value', contentToCopy);
		$('body').append(inputElement);
		// Chọn và sao chép nội dung vào clipboard
		inputElement.select();
		document.execCommand('copy');
		// Loại bỏ phần tử input ẩn
		inputElement.remove();
		// Thay đổi nội dung của liên kết
		accountManager.find('.copy-mail i').addClass('d-none');
		accountManager.find('.copy-mail .copies').removeClass("d-none").addClass('d-inline-block');
		setTimeout(function () {
			accountManager.find('.copy-mail i').removeClass('d-none');
			accountManager.find('.copy-mail .copies').removeClass('d-inline-block').addClass("d-none");
		}, 1000);
	});

	accountManager.on('click', '.dropdown-menu', function (e) {
		if ($(this).hasClass('dropdown-menu')) {
			e.stopPropagation();
		}
	});
});

GetAccountManagerInfo()

function GetAccountManagerInfo() {
	var token = getCookie('mctehj');
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
			if (!result.includes('not found')) {
				$(".account-manager").html(result)
			}
			else {
				console.log(result)
			}
			// }
		})
}

function getCookie(name) {
	const value = `; ${document.cookie}`;
	const parts = value.split(`; ${name}=`);
	if (parts.length === 2) return parts.pop().split(';').shift();
}