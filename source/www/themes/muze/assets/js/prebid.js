class PrebidClass {
  constructor(opts) {
    this.targetPrebid = opts.targetPrebid;
    this.showHideSelectBidder = opts.showHideSelectBidder || false;
    this.urlLoadParam = opts.urlLoadParam;
    this.typeSelect = opts.typeSelect || "select2";
    this.ajaxAddParam = opts.ajaxAddParam;
    this.isCheckParamRequired = opts.isCheckParamRequired || false;
    this.wrapperBox = opts.wrapperBox || "";

    this.initCss();
    this.initTabConfig();
  }

  initCss() {
    let style = document.createElement('style');
    style.type = 'text/css';
    style.id = 'css-for-prebid-tabs';

    style.innerHTML = `
    .col-md-4 .w-50-custom{
      width:100%!important;
    }
    .w-input-group.w-input-group-small-width{
      min-width: 150px;
    }

    .nav-item-tab-setup{
        font-size: 0.75rem;
    }

    .tab-prebid-wrapper{
        position: relative;
    }

    .btn-light.fs-12 {
        font-size: 0.75rem!important;
    }

    .muze-collapes .card .card-header .btn-light.fs-12:before {
        height:  0.75rem; 
    }

    .muze-collapes .card .card-header .btn-light.fs-12:after {
        width: 0.75rem;
        right: 4px;
    }

    .dm2-tab-item{
        width: 24px;
        height: 24px;
        margin: 0 auto;
    }

    .submit-btn-add-bidder{
        margin-left: auto;
        align-items: center;
    }

    .submit-btn-add-bidder .submit{
        font-size: 0.75rem;
        height: 2.5rem;
        padding: 0.5rem 1rem;
    }

    .box-c .bidder-box div.bidder-name label.w-input-group.w-input-group-small-width {
        margin-left: 5px;
    }

    .select2-custom-width-css{
        min-width: 250px;
    }

    
    .w3-button--config:not(.disabled-config):hover {
        background-color: #ffc107!important;
    }

    .w3-button--config.w3-bar--config-item.disabled-config{
        cursor: not-allowed;
        opacity: 0.5;
    }

    .w3-bar--config .w3-button--config {
        white-space: normal;
    }

    .w3-bar--config .w3-bar--config-item {
        padding: 8px 16px;
        float: left;
        width: auto;
        border: none;
        display: block;
        outline: 0;
        text-transform: uppercase;
        font-weight: bold;
        font-size: 0.65rem;

        cursor: pointer;
        background-color: rgb(240, 240, 240);
        border-left: 1px solid #e2dada;
        border-bottom: 1px solid #e2dada;
        border-top: 1px solid #e2dada;
    }

    .w3-bar--config .w3-bar--config-item:last-child{
        border-right: 1px solid #e2dada;
    }

    .w3-bar--config {
        overflow: hidden;        
    }

    .active--config:not(.disabled-config){
        border-color: #ffd657;
        background: #ffd657;
    }

    .tooltip-inner{
        text-align: left;
    }

    .tooltip{
        font-size: 12px;
    }`;
    document.getElementsByTagName('head')[0].appendChild(style);
  }

  initTabConfig() {
    let temp = $("#tab-info-js").text();

    try {
      window.Tab_Config = JSON.parse(temp) || [];

      window.Tab_Check_Active = [];

      window.Tab_Config.forEach(element => {
        for (const key in element.Available) {
          if (Object.hasOwnProperty.call(element.Available, key)) {
            if (!window.Tab_Check_Active.includes(key)) {
              window.Tab_Check_Active.push(key);
            }
          }
        }
      });
    } catch (error) {
      console.error(error);
    }
  }

  idWrapperBox(target) {
    if (!this.wrapperBox) {
      return "";
    }

    let id = target.closest(this.wrapperBox).attr("id");

    return id ? ("#" + id) : "";
  }

  isSelect2() {
    return this.typeSelect == "select2";
  }

  isSelectPicker() {
    return this.typeSelect == "selectpicker";
  }

  selectBidderElement() {
    switch (true) {
      case this.isSelect2():
        return "select2:select";
      case this.isSelectPicker():
        return "changed.bs.select";
      default:
        return "";
    }
  }

  //Click chuyển tab config ở type prebid
  EventClickPrebidTabConfig() {
    let thisClass = this;
    $(".w3-bar--config").on("click", ".w3-bar--config-item", function (e) {
      let id = $(this).data("id");

      if ($(this).hasClass("disabled-config")) {
        return
      }

      let tabs = $(`${thisClass.idWrapperBox($(this))} .w3-bar--config .w3-bar--config-item`);

      let wrappers = $(` ${thisClass.targetPrebid} ${thisClass.idWrapperBox($(this))} .tab-prebid-wrapper`);

      try {
        wrappers.each(function () {
          if (!$(this).hasClass("d-none")) {
            $(this).append(Loading());
          }
        });

        setTimeout(() => {
          tabs.each(function () {
            if ($(this).data("id") == id) {
              $(this).addClass("active--config");
            } else {
              $(this).removeClass("active--config");

            }
          })

          wrappers.each(function () {
            if ($(this).attr("id") == id) {
              $(this).removeClass("d-none");
            } else {
              $(this).addClass("d-none");

            }
          })

          wrappers.find("._blur").remove();
        }, 200);

      } catch (error) { }

    });
  }

  //Select a bidder -> load từ backend -> show list params của bidder
  EventLoadParamsAfterSelectBidder(opts = {}) {
    let thisClass = this;

    $(this.targetPrebid).on(this.selectBidderElement(), function (e) {
      let selectBidder;
      if ($(e.target).hasClass('select-bidder')) {
        selectBidder = $(e.target);
      } else {
        selectBidder = $(e.target).closest(".select-bidder");
      }

      if (!selectBidder.length) {
        return;
      }

      let wrapper = $(e.target).closest(".tab-prebid-wrapper");

      let value = {
        id: selectBidder.find(":selected").val(),
        text: selectBidder.find(":selected").text()
      };

      let index = 1 + $(thisClass.targetPrebid + thisClass.idWrapperBox(selectBidder)).find(`.bidder-box[data-id="${value.id}"]`).length;
      let bidders = wrapper.find(`.bidder-box[data-id="${value.id}"]`);
      let hasClient = "0";

      bidders.each(function () {
        if ($(this).find(":selected").attr("value") == "1") {
          hasClient = "1";
          return false;
        }
      });

      thisClass.LoadBidderParam({
        id: value.id, name: value.text, index: index, typ: hasClient,
        wrapper: wrapper
      });

      if (thisClass.showHideSelectBidder) {
        wrapper.find(".box-select-bidder").addClass("d-none");
      }
    });
  }

  LoadBidderParam(opts = {}) {
    if (!opts.id || !opts.name || !opts.index || !opts.typ || !opts.wrapper) {
      // console.error('LoadBidderParam: Missing params', opts);
      return;
    }

    let thisClass = this;
    $.ajax({
      url: thisClass.urlLoadParam,
      type: "GET",
      data: {
        id: opts.id,
        name: opts.name,
        index: opts.index,
        type: opts.typ,
      },
      success: function (json) {
        opts.wrapper.find("#bidder-params").prepend(json.data)

        thisClass.handleSelect(opts);
      },
      error: function (xhr) {
        console.error(xhr)
      }
    })
  }

  handleSelect(opts = {}) {
    if (!opts.wrapper) {
      console.error('handleSelect: Wrapper is required');
      return;
    }

    if (this.isSelect2()) {
      opts.wrapper.find('.select-bidder').val(null).trigger('change')

      this.InitSelect2({ wrapper: opts.wrapper })
    }

    if (this.isSelectPicker()) {


      //Reset select-bidder (Add New Bidder)
      opts.wrapper.find('.select-bidder').val("");
      opts.wrapper.find('.select-bidder').selectpicker('refresh');

      if (opts.changeBidderType) {
        //reset select-bidder (Change Bidder Type) -> trường hợp chọn client / s2s
        opts.wrapper.find('.selectpicker.type-bidder').selectpicker('refresh');
      }

      this.InitSelectpicker({ wrapper: opts.wrapper })
    }
  }

  InitSelectpicker(opts = {}) {
    if (!opts.wrapper) {
      console.error('InitSelectpicker: Wrapper is required');
      return;
    }

    opts.wrapper.find('#bidder-params').find('.selectpicker').selectpicker();
  }

  InitSelect2(opts = {}) {
    if (!opts.wrapper) {
      console.error('InitSelect2: Wrapper is required');
      return;
    }
    opts.wrapper.find('.select-bidder').select2({
      // style: "text-transform: capitalize"
    });
    $.fn.select2.defaults.set("theme", "bootstrap");
  }

  EventAddParam() {
    let thisClass = this;
    $(this.targetPrebid).on("click", ".bidder-box .btn-add-param", function (e) {

      //dùng ở app
      if (thisClass.ajaxAddParam) {
        var el = $(this)
        var demand = el.closest(".bidder-box").attr("data-name")

        $.ajax({
          type: 'GET',
          url: thisClass.ajaxAddParam,
          data: { demand: demand }
        })
          .done(function (result) {
            el.closest(".bidder-box").find(".list_param").append(result)
            let wrapper = el.closest(".bidder-params-wrapper");

            thisClass.handleSelectParam({ wrapper: wrapper });
          })
        return;
      }

      //Default -> dùng ở be
      $(this).closest(".bidder-box").find(".list_param").append(thisClass.addParamHtml())
    });
  }
  //Xử lí param đã có thì disabled đi -> OK
  handleSelectParam(opts = {}) {
    if (!opts.wrapper) {
      return;
    }

    opts.wrapper.find(".param-bidder-pb").each(function () {

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

  addParamHtml() {
    return `<div class="row my-4 box-param">
                <div class="d-flex align-items-center">
                    <div class="w-input-group w-input-group-small-width">
                        <input class="form-control param_name"
                               style="width: 150px" placeholder="Param name">
                    </div>
                    <div class="w-50-custom d-flex align-items-center box-value">
                        <input type="number"
                               class="form-control param_value add-param" placeholder="Value">
                        <div class="ms-2 ps-0 align-items-center w-25">
                            <select class="form-select param_type">
                                <option value="string">String</option>
                                <option selected="" value="int">Int</option>
                                <option value="float">Float</option>
                                <option value="json">Json</option>
                                <option value="boolean">Boolean</option>
                            </select>
                        </div>
                    </div>
  
                    <div class="ms-2 ps-0 d-flex align-items-center">
                        <button type="button"
                                class="btn btn-outline-danger btn-sm px-2 rounded-2 btn-remove-param d-none">
                            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
                                 class="bi bi-trash" viewBox="0 0 16 16">
                                <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z"></path>
                                <path fill-rule="evenodd"
                                      d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z"></path>
                            </svg>
                        </button>
                    </div>
                    <span class="form-text ms-2 w-auto invalid-feedback"></span>
                </div>
            </div>`
  }

  inputParam(opts = {}) {
    if (!opts.type) {
      console.error("inputParam: Type is required");
      return;
    }

    switch (opts.type) {
      case "string":
      case "json":
        return `<input type="text" class="form-control param_value add-param" placeholder="Value">`;
      case "int":
      case "float":
        return `<input type="number" class="form-control param_value add-param" placeholder="Value">`;
      case "boolean":
        return `
        <select class="form-select param_value add-param">
            <option value="true">true</option>
            <option value="false">false</option>
        </select>
        `;
    }

  }

  //dùng ở be
  EventChangeTypeParam() {
    let thisClass = this;
    $(this.targetPrebid).on("change", ".bidder-box .param_type", function (e) {
      let type = $(this).val()
      $(this).closest(".box-value").find(".param_value").remove()
      $(this).closest(".box-value").prepend(thisClass.inputParam({ type: type }));
    });
  }

  EventSelectBidderType() {
    let thisClass = this;
    $(this.targetPrebid).on("change", "#bidder-params div.bidder-box select.type-bidder", function (e) {

      let wrapper = $(this).closest(".tab-prebid-wrapper");
      let bidderId = $(this).data("id")
      let index = $(this).data("index")
      let thisBidder = wrapper.find(`#type_${bidderId}_${index} option[value='1']`)
      let prebidClient = wrapper.find(`select.b-t-${bidderId} option[value="1"]`)

      if ($(this).val() === "1") {
        //Nếu là client -> disable tất cả select client khác của bidder đó

        //ở setup cần timeout 1s để có tác dụng
        prebidClient.attr("disabled", true)

        thisBidder.attr("disabled", false)
      } else {
        //Mở khóa lựa chọn client
        prebidClient.attr("disabled", false)
      }

      thisClass.handleSelect({ wrapper: wrapper, changeBidderType: true });
    })
  }

  EventRemoveParam() {
    let thisClass = this;
    $(".bidder-params-wrapper").on("click", ".list_param .btn-remove-param", function (e) {
      thisClass.RemoveParam({ element: this });
      // $(this).closest('.box-param').remove();
    });
  }

  RemoveParam(opts = {}) {
    if (!opts.element) {
      console.error("RemoveParam: Missing element");
      return;
    }

    let isRequired;
    if (this.isCheckParamRequired) {
      isRequired = this.checkParamRequired(opts);
    }

    if (!isRequired) {
      let thisClass = this;
      if ($(opts.element).closest('.list_param').find('div.box-param').length === 1) {
        swal({
          title: "Are you sure?",
          text: "This is the last Params if you delete this demand you will delete!",
          icon: "warning",
          buttons: true,
          dangerMode: true,
        })
          .then((willDelete) => {
            if (willDelete) {
              thisClass.RemoveBidder({ element: $(opts.element).closest('.bidder-box').find('.d a.rm_c') })
            }
          });
      } else {
        $(opts.element).closest('.box-param').remove();
      }
    }
  }

  checkParamRequired(opts = {}) {
    if (!opts.element) {
      console.error("checkParamRequired: Missing opts");
      return;
    }

    let thisClass = this;
    let isRequired = false;
    let name = $(opts.element).closest('.bidder-box').data("name");
    let param = $(opts.element).closest('.box-param').find('.param_value').data("name");
    $.ajax({
      url: thisClass.isCheckParamRequired,
      type: "GET",
      async: false,
      data: {
        name: name,
        param: param,
      },
      success: function (json) {
        if (json.required) {
          isRequired = true
          $(opts.element).closest('.box-param').find('.param_value').addClass("is-invalid").nextAll("span.invalid-feedback").text("(*) param required don't remove");
        }
      },
      error: function (xhr) {
        console.log(xhr);
      }
    })
    return isRequired;
  }

  EventRemoveBidder() {
    let thisClass = this;
    $(".bidder-params-wrapper").on("click", "div.bidder-box .d a.rm_c", function () {
      thisClass.RemoveBidder({ element: $(this) });
    })
  }

  RemoveBidder(opts = {}) {
    if (!opts.element) {
      console.error("RemoveBidder: Missing element");
      return;
    }
    let wrapper = opts.element.closest(".tab-prebid-wrapper");

    let typeElement = opts.element.closest(".bidder-box").find(".type-bidder");

    //Nếu element xóa bỏ là client thì mở khóa các bidder khác đc chọn client
    if (typeElement.find(":selected").attr("value") == "1") {
      let bidderId = typeElement.data("id");
      let prebidClient = wrapper.find(`select.b-t-${bidderId} option[value="1"]`)
      prebidClient.attr("disabled", false);
    }

    opts.element.closest(".box-c").remove();

    //Nếu hết bidder -> show ô chọn bidder lên
    if (this.showHideSelectBidder && wrapper.find(".box-c").length === 0) {
      wrapper.find(".box-select-bidder").removeClass("d-none");
    }
  }

  EventToggleBttRemoveParam() {
    $(".bidder-params-wrapper").on({
      mouseenter: function () {
        $(this).find(".btn-remove-param").removeClass("d-none")
      },
      mouseleave: function () {
        $(this).find(".btn-remove-param").addClass("d-none")
      }
    }, ".list_param .box-param");
  }

  newTabActive() {
    let Tab_Active = {};

    window.Tab_Config.forEach(element => {
      Tab_Active[element.Id] = { list: {} };
      for (const key in element.Available) {
        if (Object.hasOwnProperty.call(element.Available, key)) {
          Tab_Active[element.Id].list[key] = true;
        }
      }
    });

    return Tab_Active;
  }

  handleDisabledTabs(optionTarget, fullInfo, opts = {}) {
    let tabActive = this.newTabActive();

    if (!window.Tab_Check_Active.includes(optionTarget.optionAjax) && !optionTarget.force) {
      return;
    }

    fullInfo.forEach(info => {
      if (!window.Tab_Check_Active.includes(info.optionAjax)) {
        return;
      }

      let listSelected = info.list_selected.map(a => a.name);

      if (listSelected.length > 0) {
        window.Tab_Config.forEach(tabConfig => {
          let availableNow = tabConfig.Available;
          let tabConfigNow = availableNow[info.optionAjax].List;

          let avail = false;

          if (tabConfigNow) {
            if (tabConfigNow.length == 1 && tabConfigNow[0] == "all") {
              avail = true;
            } else {

              for (let index = 0; index < listSelected.length; index++) {
                if (tabConfigNow.includes(listSelected[index])) {
                  avail = true;
                  break;
                }
              }
            }
          }

          if (!avail) {
            tabActive[tabConfig.Id].list[info.optionAjax] = false;
            tabActive[tabConfig.Id].tabConfig = tabConfig;
          }
        });
      }
    });

    //ĐOẠN CHECK MULTICONDITIONAL => Tab active hay ko
    window.Tab_Config.forEach(tabConfig => {
      let onByCondi;

      let multiConditional = tabConfig.MultiConditonal;

      if (multiConditional) {
        let size = Object.keys(multiConditional).length;
        let count = 0;

        loopcondi: for (const key in multiConditional) {
          if (Object.hasOwnProperty.call(multiConditional, key)) {
            let has = false;

            for (let index = 0; index < fullInfo.length; index++) {
              if (fullInfo[index].optionAjax == key) {
                has = true;
                const condi = multiConditional[key].List;

                if (fullInfo[index].list_selected.length == 0) {
                  count++;
                  continue loopcondi;
                } else {
                  for (let i = 0; i < fullInfo[index].list_selected.length; i++) {
                    if (condi.includes(fullInfo[index].list_selected[i].name)) {
                      count++;
                      continue loopcondi;
                    }
                  }
                  has = false;
                  break;
                }
              }
            }

            if (!has) {
              onByCondi = false;
              break;
            }
          }
        }

        if (count == size) {
          onByCondi = true;
        } else {
          onByCondi = false;
        }

      } else {
        onByCondi = false;
      }

      if (onByCondi && tabActive[tabConfig.Id]) {
        tabActive[tabConfig.Id].active = true;
      }
    });

    for (const key in tabActive) {
      if (Object.hasOwnProperty.call(tabActive, key)) {
        const element = tabActive[key];

        let isActive = true;

        //nếu tab ko active ở chế độ multiconditional thì mới check đến active bình thường
        if (!element.active) {
          for (const k in element.list) {
            if (Object.hasOwnProperty.call(element.list, k)) {
              if (element.list[k] == false) {
                isActive = false;
                break;
              }
            }
          }
        }

        if (isActive) {
          this.enabledTab({ Id: key, target: opts.target });
        } else {
          this.disabledTab({ Id: key, tabConfig: element.tabConfig, target: opts.target });
        }
      }
    }

    //Change tab active do tab đang active bị disabled
    let needChangeActive = $(`${this.idWrapperBox(opts.target)} .active--config.disabled-config`);

    if (needChangeActive.length != 0) {

      let hasTabToChange = false;

      $(`${this.idWrapperBox(opts.target)} .w3-bar--config-item`).each(function () {
        if (!$(this).hasClass("disabled-config")) {
          $(this).click();
          hasTabToChange = true;
          return false;
        }
      });

      if (hasTabToChange) {
        needChangeActive.removeClass("active--config");
      }
    }
  }

  disabledTab(opts = {}) {
    let element = $(`${this.idWrapperBox(opts.target)} .w3-bar--config-item[data-id="tab-prebid-${opts.Id}"]`);
    element.addClass("disabled-config");

    let title = `This config is only available with:`;

    if (opts.tabConfig) {

      if (opts.tabConfig.Available) {
        title += "<br>";

        let plus = "";

        for (const key in opts.tabConfig.Available) {
          if (Object.hasOwnProperty.call(opts.tabConfig.Available, key)) {
            const element = opts.tabConfig.Available[key];

            title += plus + element.Name + ": " + element.List.join(", ") + ".";

            if (plus == "") {
              plus = "<br>";
            }
          }
        }
      }

      if (opts.tabConfig.MultiConditonal) {
        title += "<br>";

        let plus = "";

        for (const key in opts.tabConfig.MultiConditonal) {
          if (Object.hasOwnProperty.call(opts.tabConfig.MultiConditonal, key)) {
            const element = opts.tabConfig.MultiConditonal[key];

            title += plus + element.Name + ": " + element.List.join(", ") + "";

            if (plus == "") {
              plus = " & ";
            }
          }
        }

        title += "."
      }

    }

    element.tooltip("dispose");
    element.tooltip({
      html: true,
      title: title,
    });
  }
  enabledTab(opts = {}) {
    let element = $(`${this.idWrapperBox(opts.target)} .w3-bar--config-item[data-id="tab-prebid-${opts.Id}"]`);
    element.tooltip("dispose");

    element.removeClass("disabled-config");
  }

  buttonAddBidder() {
    let thisClass = this;
    $(this.targetPrebid).on("click", ".add-bidder", function (e) {
      let wrapper = $(e.target).closest(".tab-prebid-wrapper");

      if (!wrapper) {
        return;
      }

      wrapper.find(".box-select-bidder").removeClass("d-none");
      thisClass.InitSelectpicker({ wrapper: wrapper });
    });
  }

  SelectParamPB() {
    let thisClass = this;
    $(this.targetPrebid).on("changed.bs.select", ".param-bidder-pb", function (e) {
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
        thisClass.InitSelectpicker({ wrapper: $(this).closest(".tab-prebid-wrapper") })
      }
    })

  }

  firstLoadBidderDisabled() {
    let wrappers = $(this.targetPrebid).find(".tab-prebid-wrapper");

    wrappers.each(function () {

      let wrapper = $(this);

      let typeBidder = wrapper.find(".type-bidder");

      typeBidder.each(function () {
        if ($(this).val() == 1) {

          let id = $(this).data("id");

          let sameBidder = wrapper.find(`.type-bidder[data-id='${id}'`);

          sameBidder.each(function () {
            if ($(this).val() == 2) {
              $(this).find("[value='1']").attr("disabled", true);
            }
          });
        }
      });
    });

  }

  //Click vào tab ko bị disabled đầu tiên
  activeInFirstLoad() {
    $(".w3-bar--config").each(function () {
      let active = $(this).find(".active--config")

      if (active.length == 0) {
        try {
          $(this).find(".w3-bar--config-item.w3-button--config:not(.disabled-config)").first().click();
        } catch (error) { }
      }
    })
  }
}

export default PrebidClass;