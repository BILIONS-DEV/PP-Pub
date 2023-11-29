module.exports = {
    firstLoad, appendContentSelect, refreshListContent
    // getDataQuestion
}

var listContentSelected = []

function firstLoad() {
    refreshListContent();
}

function appendContentSelect(id, title, thumb) {
    $("#show-content-select").append(`
        <div class="box-c content" data-content-id="` + id + `">
            <hr class="mt-5 hr_custom bg-gray-400">
            <div class="bidder-box">
                <div class="row my-4">
                    <div class="d-flex align-items-center bidder-name">
                        <label class="col-form-label form-label form-label-lg w-input-group text-uppercase lb-content" style="min-width: 50px">100</label>
                        <div class="d">
                            <a class="btn p-1 ms-2 rm_c remove-content" data-content-id="` + id + `">
                                <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" fill="currentColor" class="bi bi-x-lg" viewBox="0 0 16 16">
                                    <path d="M1.293 1.293a1 1 0 0 1 1.414 0L8 6.586l5.293-5.293a1 1 0 1 1 1.414 1.414L9.414 8l5.293 5.293a1 1 0 0 1-1.414 1.414L8 9.414l-5.293 5.293a1 1 0 0 1-1.414-1.414L6.586 8 1.293 2.707a1 1 0 0 1 0-1.414z"></path>
                                </svg>
                            </a>
                        </div>
                    </div>
                </div>
                <div class="row my-4">
                    <div class="d-flex align-items-center">
                        <span class="avatar avatar-lg rounded-circle d-flex align-items-center justify-content-center me-3">
                            <img src="` + thumb + `">
                        </span> 
                        <span>` + title + `</span>
                    </div>
                </div>
            </div>
        </div>`
    )
    refreshListContent();
}

function refreshListContent() {
    indexContent = 1
    $('#show-content-select > div.box-c > div.bidder-box').each(function () {
        $(this).find(".lb-content").html(indexContent)
        indexContent++
    })
}
