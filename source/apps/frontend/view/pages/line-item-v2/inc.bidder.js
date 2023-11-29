module.exports = {
    AddParam, ChangeTypeParam, LoadBidderParam, AddBidder, SelectParamPB
};

function AddParam() {
    $("#bidder-params").on("click", ".bidder-box .btn-add-param", function (e) {

        var el = $(this)
        var demand = el.closest(".bidder-box").attr("data-name")

        $.ajax({
            type: 'GET',
            url: '/line-item-v2/addParamBidder',
            data: { demand: demand }
        })
            .done(function (result) {
                el.closest(".bidder-box").find(".list_param").append(result)
                handleSelectParam()
            })

    });
}

function handleSelectParam() {
    $("#bidder-params").find(".param-bidder-pb").each(function () {
        var select_param = $(this)
        var list_param = []
        select_param.closest(".list_param").find(".param_value").each(function () {
            // get param already exist
            if ($(this).attr("data-name")) {
                list_param.push($(this).attr("data-name"))
            }
        })
        if (list_param) {
            list_param.forEach(function (value, index) {
                select_param.find('[value="' + value + '"]').attr("disabled", true)
            })
            select_param.selectpicker('refresh');
        }
    })
}

function SelectParamPB() {
    $("#bidder-params").on("changed.bs.select", ".param-bidder-pb", function (e) {
        var param = $(this).val()
        if (!param) {
            return;
        }

        var BidderId = $(this).closest(".bidder-box").attr("data-id")
        var BidderIndex = $(this).closest(".bidder-box").attr("data-index")
        var example = $(this).find('[value="' + param + '"]').attr("data-example")
        var type = $(this).find('[value="' + param + '"]').attr("data-type")
        var data_type = $(this).find('[value="' + param + '"]').attr("data-type")
        switch (type) {
            case "int":
            case "integer":
            case "number":
            case "float":
            case "decimal":
                type = "number"
                data_type = "int";
                example = type
                break;
            case "float":
            case "decimal":
                type = "number"
                data_type = "float";
                example = type
                break;
            case "boolean":
            case "bool":
                type = "boolean"
                data_type = "boolean";
                break;
            case "string":
                type = "text"
                data_type = "string";
                example = "string"
                break;
            default:
                type = "text"
                break
        }

        var html = ''
        if (type == "boolean") {
            html = `<div class="center-selectpicker pp-9">
                            <select id="` + BidderId + `-` + param + `-` + BidderIndex + `"
                                    class="form-control selectpicker param_value add-param"
                                    data-type="boolean"
                                    data-name="` + param + `">
                                <option value="true">true</option>
                                <option value="false">false</option>
                            </select>
                        </div>`
        } else {
            html = `<input style="border-left: 0;border-radius: 0px;" type="` + type + `"
                               class="pp-10 param_value add-param"
                               id="` + BidderId + `-` + param + `-` + BidderIndex + `"
                               placeholder="` + example + `"
                               data-type="` + data_type + `"
                                data-name="` + param + `">`
        }
        if (html) {
            $(this).closest(".box-param").find('.result-param-value').html(html)
            InitSelectpicker()
        }
    })

}

function AddBidder() {
    $("#tab-prebid").on("click", ".add-bidder", function (e) {
        $("#tab-prebid").find(".box-select-bidder").removeClass("d-none")
        InitSelectpicker()
    });
}

function InitSelectpicker() {
    $("#bidder-params").find('.selectpicker').selectpicker();
}

// function addParamHtml() {
//     return `<div class="row my-4 box-param">
//                 <div class="d-flex align-items-center">
//                     <div class="w-input-group">
//                         <input class="form-control param_name" style="max-width: 240px" placeholder="Param name">
//                     </div>
//                     <div class="w-50-custom d-flex align-items-center box-value">
//                         <input type="number" class="form-control param_value add-param" placeholder="Value">
//                         <div class="ms-2 ps-0 align-items-center w-25">
//                             <select class="form-select param_type">
//                                 <option value="string">String</option>
//                                 <option selected="" value="int">Int</option>
//                                 <option value="float">Float</option>
//                                 <option value="json">Json</option>
//                                 <option value="boolean">Boolean</option>
//                             </select>
//                         </div>
//                     </div>
//
//                     <div class="ms-2 ps-0 d-flex align-items-center">
//                         <button type="button" class="btn btn-outline-danger btn-sm px-2 rounded-2 btn-remove-param d-none">
//                             <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-trash" viewBox="0 0 16 16">
//                                 <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z"></path>
//                                 <path fill-rule="evenodd" d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z"></path>
//                             </svg>
//                         </button>
//                     </div>
//                     <span class="form-text ms-2 w-auto invalid-feedback"></span>
//                 </div>
//             </div>`
// }

function ChangeTypeParam() {
    $("#bidder-params").on("changed.bs.select", ".param_type", function (e) {
        let type = $(this).val()
        if (!type) {
            return
        }
        $(this).closest(".box-param").find(".param_value").remove()
        if (type === "string" || type === "json") {
            $(this).closest(".box-param").find('.box-value').prepend(inputText())
        } else if (type === "int" || type === "float") {
            $(this).closest(".box-param").find('.box-value').prepend(inputNumber())
        } else if (type === "boolean") {
            $(this).closest(".box-param").find('.box-value').prepend(selectBool())
            InitSelectpicker()
        }
    });
}


function inputText() {
    return '<input style="border-left: 0;border-radius: 0 2px 2px 0;" class="pp-10 param_value" type="text" placeholder="Value...">'
}

function inputNumber() {
    return '<input style="border-left: 0;border-radius: 0 2px 2px 0;" class="pp-10 param_value" type="number" placeholder="Value...">'
}

function selectBool() {
    return `<div class="end-selectpicker">
                <select class="form-control selectpicker param_value">
                    <option value="true">true</option>
                    <option value="false">false</option>
                </select>
            </div>`
}


function LoadBidderParam(id, name, index, typ) {
    $.ajax({
        url: "/line-item-v2/loadParam",
        type: "GET",
        data: {
            id: id,
            name: name,
            index: index,
            type: typ,
        },
        success: function (json) {
            $("#bidder-params").prepend(json.data)
            $('#select-bidder').val("");
            $('#select-bidder').selectpicker('refresh')
            // $('#select-bidder').val(null).trigger('change')
            InitSelectpicker()
        },
        error: function (xhr) {
            console.log(xhr)
        }
    })
}