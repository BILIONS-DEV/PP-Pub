module.exports = {
    SelectSize, RemoveSize, AddAdSlot
};

function AddAdSlot() {
    $(".add-adslot").click(function () {
        $(".adslot-sample").clone().prependTo("#adsense_ad_slot_item")
        $("#adsense_ad_slot_item").find(".adslot-sample").removeClass("adslot-sample")
        $("#adsense_ad_slot_item").find('.select-adsense-ad-slot-size').each(function () {
            $(this).selectpicker('refresh');
        })
    })
}

function SelectSize() {
    // Select Size of AdsenseAdSlots
    $("#adsense_ad_slot_item").on("changed.bs.select", ".select-adsense-ad-slot-size", function (e) {
        const selectedElement = $(e.currentTarget);
        const size = selectedElement.val();
        if (!size) {
            return
        }
        if ($("#adsense_ad_slot_item").find(`#adsense-ad-slot-` + size).length) {
            return new AlertError("Size already exists")
        }
        // cancel disabled size
        var size_ipunt = selectedElement.closest(".box-c").find("input").attr("data-size")
        if (size_ipunt) {
            $("#tab-google").find('.size-' + size_ipunt).prop('disabled', false);
        }
        // disabled size
        selectedElement.closest(".box-c").find("input").attr("id", "adsense-ad-slot-" + size).attr("data-size", size)
        $("#tab-google").find('.size-' + size).prop('disabled', true);
        $("#adsense_ad_slot_item").find('.select-adsense-ad-slot-size').each(function () {
            $(this).selectpicker('refresh');
        })
        // const itemLabel = $(this).data("item-label")
        // const itemPlaceholder = $(this).data("item-placeholder")
        // console.log(itemLabel);
        // console.log(itemPlaceholder);
        $.each(this.options, function (i, item) {
            if (item.selected) {
                $(item).prop("disabled", true);
            }
        });
        // $("#adsense_ad_slot_item").prepend(SelectAdsenseSize(size, itemLabel, itemPlaceholder))
        // $(this).select2("val", "All")
    });
}

function RemoveSize() {
    $("#adsense_ad_slot_item").on("click", ".remove-ad-slot", function () {
        let size = $(this).data("size")
        $(`#size-` + size).prop('disabled', false);
        $(this).parents("div.box-c").remove()
        // console.log(`size-` + size)
        // $("#select-adsense-ad-slot-size").select2("val", "All")
    })
}

// function SelectAdsenseSize(size, itemLabel, itemPlaceholder) {
//     return `<div class="box-c">
//                 <hr class="mt-5 hr_custom bg-gray-400">
//                 <div class="bidder-box">
//                     <div class="row my-4">
//                         <div class="d-flex align-items-center bidder-name">
//                             <label class="col-form-label form-label form-label-lg w-input-group text-uppercase">
//                                 ${size}
//                             </label>
//                             <div class="d">
//                                 <a class="btn p-1 ms-2 rm_c remove-ad-slot" data-size="${size}">
//                                     <svg xmlns="http://www.w3.org/2000/svg" width="15"
//                                          height="15" fill="currentColor" class="bi bi-x-lg"
//                                          viewBox="0 0 16 16">
//                                         <path d="M1.293 1.293a1 1 0 0 1 1.414 0L8 6.586l5.293-5.293a1 1 0 1 1 1.414 1.414L9.414 8l5.293 5.293a1 1 0 0 1-1.414 1.414L8 9.414l-5.293 5.293a1 1 0 0 1-1.414-1.414L6.586 8 1.293 2.707a1 1 0 0 1 0-1.414z"></path>
//                                     </svg>
//                                 </a>
//                             </div>
//                         </div>
//                     </div>
//
//                     <div class="row my-4">
//                         <div class="d-flex align-items-center">
//                             <label class="col-form-label form-label form-label-lg w-input-group">
//                                 ${itemLabel}
//                             </label>
//                             <input id="adsense-ad-slot-${size}" type="number"
//                                    class="form-control w-50-custom param_value"
//                                    placeholder="${itemPlaceholder}"
//                                    data-size="${size}" value="">
//                             <span class="form-text ms-2 w-auto invalid-feedback"></span>
//                         </div>
//                     </div>
//                 </div>
//             </div>`
// }

