class General {
  constructor(opts = {}) {
    this.isEdit = opts.isEdit || false;
    this.playTypeId = opts.playTypeId || "#type";
    this.nativeAdStyleId = opts.nativeAdStyleId || "#ad_style";
    this.videoPlayerLayoutClass = opts.videoPlayerLayoutClass || "#player_layout";
    this.xhrMakePreview = null;
    this.baseUrl = "/player-v2/template";
    this.ajaxPreview = opts.ajaxPreview || this.baseUrl + "/preview";
    this._validFileImageExtensions = [".jpg", ".jpeg", ".bmp", ".gif", ".png"];
    this.maxSize = 15 * 1024 * 1024;
    this.domainUpload = "https://ul.pubpowerplatform.io";
    this.api_upload_image = `${this.domainUpload}/api/v1/image`;
    this.formId = "#formTemplate";
    this.ajxURL = opts.ajxURL || "/player/template/add";
    this.isFirst = this.isEdit ? true : false; //phục vụ nếu edit lần đầu ko sửa dữ liệu -> những lần sau change player layout -> đổi về default
  }

  init() {
    if (this.isEdit) {
      $(this.formId).append(Loading({ topLoading: "15%" }));
    }

    $(document).ready(() => {
      this.startEvents();
    });
  }

  startEvents() {
    this.firstLoad();
    if (this.isFirst) {
      this.isFirst = false;
    }

    if (this.isEdit) {
      $(this.formId).find("._blur").remove();
    }

    this.eventClickTab();
    this.eventChangePlayerType();
    this.videoEvents();
    this.nativeEvents();
    this.eventSelectColor();
    $(this.formId).on("change", () => {
      this.makePreview()
    });
    this.formEvents();
  }

  formEvents() {
    let formElement = $(this.formId);
    formElement.find("input").on("input", function (e) {
      let inputElement = $(this);
      if (inputElement.hasClass("is-invalid")) {
        inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
        inputElement.parent().removeClass("is-invalid").next(".invalid-feedback").empty();
      }
    });

    formElement.find("textarea").on("input", function (e) {
      let inputElement = $(this);
      if (inputElement.hasClass("is-invalid")) {
        inputElement.removeClass("is-invalid").next(".invalid-feedback").empty();
      }
    });

    let thisClass = this;

    formElement.find(".submit").on("click", function (e) {
      e.preventDefault();
      const buttonElement = $(this);
      const submitButtonText = buttonElement.text();
      const submitButtonTextLoading = "Loading...";
      var postData = formElement.serializeObject();
      postData = thisClass.makeData(postData);

      $.ajax({
        url: thisClass.ajxURL,
        type: "POST",
        dataType: "JSON",
        contentType: "application/json",
        data: JSON.stringify(postData),
        beforeSend: function (xhr) {
          buttonElement.attr('disabled', true).text(submitButtonTextLoading);
        },
        error: function (jqXHR, exception) {
          const msg = AjaxErrorMessage(jqXHR, exception);
          new AlertError("AJAX ERROR: " + msg);
          buttonElement.attr('disabled', false).text(submitButtonText);
        },
        success: function (responseJSON) {
          buttonElement.attr('disabled', false).text(submitButtonText);
        },
        complete: function (res) {
          thisClass.added(res.responseJSON);
        }
      });
    });
  }

  added(response) {
    switch (response.status) {
      case "error":
        if (response.errors.length === 1 && response.errors[0].id === "") {
          new AlertError(response.errors[0].message);
        } else if (response.errors.length > 0) {
          $.each(response.errors, function (key, value) {
            let inputElement = $("#" + value.id);
            if (key === 0) {
            }
            inputElement.addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
            inputElement.parent().addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
          });
          if (response.errors[0].id === "floating_width_mobile") {
            $(".play_mode[data-type='mobile']").click();
          } else if (response.errors[0].id === "floating_width") {
            $(".play_mode[data-type='desktop']").click();
          }
          $("#" + response.errors[0].id).focus();
          $("#" + response.errors[0].id).prev('label').focus();
          new AlertError(response.errors[0].message, function () {
            $("#" + response.errors[0].id).focus();
            $("#" + response.errors[0].id).prev('label').focus();
          });
        } else {
          new AlertError("Error!");
        }
        break;
      case "success":
        if (this.isEdit) {
          NoticeSuccess("Template has been updated successfully");
          return;
        }
        NoticeSuccess("Template has been created successfully");
        setTimeout(function () {
          window.location.replace(this.baseUrl);
        }, 1000);
        break;
      default:
        new AlertError("Undefined");
        break;
    }
  }

  eventSelectColor() {
    let elementInputColor = $(".input-color");
    elementInputColor.on("input", function () {
      $(this).next("input").val($(this).val());
    });
    elementInputColor.next("input").on("input", function (e) {
      $(this).prev("input").val(e.target.value);
    });
  }

  videoEvents() {
    this.eventChangeVideoLayout();
    this.eventChangePlayerSize();
    this.eventChangePlayerMode();
    this.eventChangeFloat();
    this.eventChangeFloat("floating_on_mobile");
    this.eventChangeVideoPosition();
    this.eventChangeVideoPosition("position_mobile");
    this.eventChangeMainTitle();
    this.eventAddLogo();
    this.eventChangeAutoSkip();
  }

  nativeEvents() {
    this.eventChangeNativeAdStyle();
    this.eventChangeNativeOptimizeLayout();
    this.eventChangeNativeMultiDimension();
    this.eventChangeNativeAdSize();
  }

  eventChangeAutoSkip() {
    $("#auto_skip").change(() => {
      this.checkShowAdvertisement();
    });
  }

  eventAddLogo() {
    $(".btn-add-logo").on("click", function () {
      $("#logo").click();
    });

    let thisClass = this;
    $("#logo").on("change", function (e) {
      thisClass.validateAndUploadLogo($(this), thisClass._validFileImageExtensions, e);
    });
  }

  validateAndUploadLogo(element, validFile, event) {
    let oInput = element[0];
    if (oInput.type === "file") {
      let sFileName = oInput.value;
      if (sFileName.length > 0) {
        let blnValid = false;
        for (let j = 0; j < validFile.length; j++) {
          let sCurExtension = validFile[j];
          if (sFileName.substr(sFileName.length - sCurExtension.length, sCurExtension.length).toLowerCase() === sCurExtension.toLowerCase()) {
            blnValid = true;
            this.uploadFile(event, element);
            break;
          }
        }
        if (!blnValid) {
          new AlertError("Sorry, " + sFileName + " is invalid, allowed extensions are: " + validFile.join(", "));
          return false;
        }

      }
    }
    return true;
  }

  uploadFile(event, element) {
    let thisClass = this;
    var fd = new FormData();
    var file = event.target.files[0];
    if (file.size > this.maxSize) {
      new AlertError("You uploaded file over 10mb, please choose another file!");
      return;
    }
    fd.append('file', file);
    $.ajax({
      url: thisClass.api_upload_image, type: "POST", // dataType: 'json',
      contentType: false, processData: false, data: fd, beforeSend: function (xhr) {
        // buttonElement.attr('disabled', true).text(submitButtonTextLoading);
      }, error: function (jqXHR, exception) {
        const msg = AjaxErrorMessage(jqXHR, exception);
        new AlertError("AJAX ERROR: " + msg);
        // buttonElement.attr('disabled', false).text(submitButtonText);
      }, success: function (responseJSON) {
        // buttonElement.attr('disabled', false).text(submitButtonText);
      }, complete: function (res) {
        thisClass.afterUpload(res.responseJSON, element);
      }
    });
  }

  afterUpload(response, element) {
    switch (response.status) {
      case "error":
        if (response.errors.length === 1 && response.errors[0].id === "") {
          new AlertError(response.errors[0].message);
        } else if (response.errors.length > 0) {
          $.each(response.errors, function (key, value) {
            let inputElement = $("#" + value.id);
            if (key === 0) {
            }
            inputElement.addClass("is-invalid").nextAll("span.invalid-feedback").text(value.message);
          });
          new AlertError(response.errors[0].message);
        } else {
          new AlertError("Error!");
        }
        break;
      case "success":
        $("#preview-logo").attr("src", this.domainUpload + response.data_object.thumb[0]);
        $("#logo").attr("value", response.data_object.thumb[0]);
        break;
      default:
        new AlertError("Undefined");
        break;
    }
  }

  checkDisableControls() {
    // Controls Color
    let controlsColor = $(".controls");
    if (this.isInstream()) {
      controlsColor.removeClass("disabled");
      controlsColor.prop("disabled", false);
      controlsColor.next().removeClass("disabled");
    } else {
      controlsColor.addClass("disabled");
      controlsColor.prop("disabled", true);
      controlsColor.next().addClass("disabled");
    }
  }

  checkDisableLogoConfig() {
    // Logo Config
    let logoConfig = $(".logo-config");
    if (this.isInstream()) {
      logoConfig.removeClass("disabled");
      logoConfig.prop("disabled", false);
      logoConfig.next().removeClass("disabled");
    } else {
      logoConfig.addClass("disabled");
      logoConfig.prop("disabled", true);
      logoConfig.next().addClass("disabled");
    }
  }

  // checkCustomizeControls() {
  //   if (this.isOutstream()) {
  //     $("#customize-controls").addClass("d-none");
  //     $("#fullscreen_button").prop("disabled", true);
  //     $("#next_prev_arrows").prop("disabled", true);
  //     $("#next_prev_10_sec").prop("disabled", true);
  //     $("#video_config").prop("disabled", true);
  //     $("#show_view_likes").prop("disabled", true);
  //     $("#share_button").prop("disabled", true);
  //   } else {
  //     $("#customize-controls").removeClass("d-none");
  //     $("#fullscreen_button").prop("disabled", false);
  //     $("#next_prev_arrows").prop("disabled", false);
  //     $("#next_prev_10_sec").prop("disabled", false);
  //     $("#video_config").prop("disabled", false);
  //     $("#show_view_likes").prop("disabled", false);
  //     $("#share_button").prop("disabled", false);
  //   }
  // }

  checkDisableColor() {
    let layout = $("#player_layout").val();

    // Controls Color
    let controlsColor = $(".box-controls-color");
    if (this.isInstream()) {
      controlsColor.removeClass("d-none");
      controlsColor = $("#control_color");
      controlsColor.removeClass("disabled");
      controlsColor.prop("disabled", false);
      controlsColor.prev().removeClass("disabled");
    } else {
      controlsColor.addClass("d-none");
    }

    // Description Color
    let descriptionColor = $(".box-description-color");
    if (this.isInstream()) {
      if (layout === "3" || layout === "4" || layout === "5" || layout === "7") {
        descriptionColor.removeClass("d-none");
        descriptionColor = $("#description_color");
        descriptionColor.removeClass("disabled");
        descriptionColor.prop("disabled", false);
        descriptionColor.prev().removeClass("disabled");
      } else {
        descriptionColor.addClass("d-none");
      }
    } else {
      descriptionColor.addClass("d-none");
    }

  }

  eventChangeMainTitle() {
    $("#show_main_title").change(() => {
      this.changeMainTitleAction();
    });
  }

  changeMainTitleAction() {
    // Show main title
    let showMainTitleElement = $("#show_main_title");
    if (this.isInstream()) {
      showMainTitleElement.removeClass("disabled");
      showMainTitleElement.prop("disabled", false);
      showMainTitleElement.next().removeClass("disabled");
    } else {
      showMainTitleElement.addClass("disabled");
      showMainTitleElement.prop("disabled", true);
      showMainTitleElement.next().addClass("disabled");
    }

    // Main title text
    let mainTitleTextElement = $("#main_title_text");
    if (this.isInstream() && showMainTitleElement.is(":checked")) {
      mainTitleTextElement.removeClass("disabled");
      mainTitleTextElement.prop("disabled", false);
      mainTitleTextElement.next().removeClass("disabled");
    } else {
      mainTitleTextElement.addClass("disabled");
      mainTitleTextElement.prop("disabled", true);
      mainTitleTextElement.next().addClass("disabled");
    }
  }

  checkDisableDisplayOption() {
    let layout = $(this.videoPlayerLayoutClass).val();

    this.changeMainTitleAction();

    // Show content title
    let showContentTitleElement = $("#show_content_title");
    if (this.isInstream()) {
      showContentTitleElement.removeClass("disabled");
      showContentTitleElement.prop("disabled", false);
      showContentTitleElement.next().removeClass("disabled");

      $("#show-content-title").closest('div').removeClass("d-none");
      showContentTitleElement.removeClass("disabled");
      showContentTitleElement.prop("disabled", false);
      showContentTitleElement.next().removeClass("disabled");
    } else {
      showContentTitleElement.addClass("disabled");
      showContentTitleElement.prop("disabled", true);
      showContentTitleElement.next().addClass("disabled");
    }

    // Show content description
    let showContentDescElement = $("#show_content_description");
    let boxShowContentDesc = $(".box-show-content-desc");
    if (this.isInstream()) {
      if (layout === "3" || layout === "4" || layout === "5" || layout === "7") {
        boxShowContentDesc.removeClass("d-none");
        showContentDescElement.removeClass("disabled");
        showContentDescElement.prop("disabled", false);
        showContentDescElement.next().removeClass("disabled");
      } else {
        boxShowContentDesc.addClass("d-none");
        showContentDescElement.addClass("disabled");
        showContentDescElement.prop("disabled", true);
        showContentDescElement.next().addClass("disabled");
      }
    } else {
      boxShowContentDesc.addClass("d-none");
      showContentDescElement.addClass("disabled");
      showContentDescElement.prop("disabled", true);
      showContentDescElement.next().addClass("disabled");
    }
    // Show controls
    let showControlsElement = $("#show_controls");
    if (this.isInstream()) {
      showControlsElement.removeClass("disabled");
      showControlsElement.prop("disabled", false);
      showControlsElement.next().removeClass("disabled");
      showControlsElement.closest("div").removeClass("d-none");
      showControlsElement.removeClass("disabled");
      showControlsElement.prop("disabled", false);
      showControlsElement.next().removeClass("disabled");
    } else {
      showControlsElement.addClass("disabled");
      showControlsElement.prop("disabled", true);
      showControlsElement.next().addClass("disabled");
    }
  }

  changeVideoPositionAction(input = "position_desktop") {
    let position = $(`#${input}`).val();
    let elementPosition = $(`.${input}`);
    elementPosition.addClass("d-none");

    let showElems;
    switch (position) {
      case "1":
        showElems = $(`.${input}_bottom_right`);
        break;
      case "2":
        showElems = $(`.${input}_bottom_left`);
        break;
      case "3":
        showElems = $(`.${input}_top_right`);
        break;
      case "4":
        showElems = $(`.${input}_top_left`);
        break;
    }
    showElems.removeClass("d-none");
  }

  eventChangeVideoPosition(input = "position_desktop") {
    $(`#${input}`).on("change", () => {
      this.changeVideoPositionAction(input);
    });
  }

  resizePreviewNative() {
    try {
      const ele = document.getElementById("resizeWrapper");
      if (!ele) {
        return;
      }
      let x = 0;
      let w = 0;
      const mouseDownHandler = (e) => {
        x = e.clientX;
        const styles = window.getComputedStyle(ele);
        w = parseInt(styles.width, 10);
        document.addEventListener("mousemove", mouseMoveHandler);
        document.addEventListener("mouseup", mouseUpHandler);
      };

      let thisClass = this;

      const mouseMoveHandler = (e) => {
        const dx = e.clientX - x;
        ele.style.width = (w + dx) + "px";
        thisClass.setPreviewSize();
      };

      const mouseUpHandler = function () {
        document.removeEventListener("mousemove", mouseMoveHandler);
        document.removeEventListener("mouseup", mouseUpHandler);
      };

      const resizer = document.querySelector(".resizer");
      resizer.addEventListener("mousedown", mouseDownHandler);
    } catch (error) {
      // Helper.log(error);
    }
    this.setPreviewSize();
  }

  setPreviewSize() {
    if ($("#preview_size").hasClass("d-none")) {
      return;
    }
    if (document.querySelector("#adsPlacement")) {
      setTimeout(() => {
        let w = document.querySelector("#adsPlacement").offsetWidth;
        let h = document.querySelector("#adsPlacement").offsetHeight;
        document.getElementById("p_height").innerHTML = h;
        document.getElementById("p_width").innerHTML = w;
      }, 50);
    }
  }

  eventChangeWidth() {
    $("#width").change(() => {
      this.makePreview();
    });
  }

  eventChangeFloat(input = "floating_on_desktop") {
    $("#" + input).change(() => {
      this.changeFloatAction(input);
    });
  }

  changeFloatAction(input = "floating_on_desktop") {
    let isFloating = $("#" + input).is(":checked");
    let elementFloating = $("." + input);
    if (isFloating) {
      elementFloating.removeClass("disabled");
      elementFloating.prop("disabled", false);
      elementFloating.next().removeClass("disabled");
    } else {
      elementFloating.addClass("disabled");
      elementFloating.prop("disabled", true);
      elementFloating.next().addClass("disabled");
    }
  }

  firstLoad() {
    this.changePlayerSizeAction();
    this.changeVideoPositionAction();
    this.changeVideoPositionAction("position_mobile");
    this.checkDisableDisplayOption();
    this.checkDisableColor();
    this.checkDisableControls();
    this.checkDisableLogoConfig();
    this.checkShowColumnsNumber();
    this.checkShowColumnsPosition();
    this.checkShowAdvertisement();
    this.changeFloatAction();
    this.changeFloatAction("floating_on_mobile");
    // this.checkCustomizeControls();
    this.checkValueimpressionLogo();
    this.changeNativeAdStyleAction();

    this.changePlayerTypeAction();
  }

  checkValueimpressionLogo() {
    if (this.isInstream()) {
      $("#pubpower_logo").closest("div").removeClass("d-none");
    } else {
      $("#pubpower_logo").closest("div").addClass("d-none");
    }
  }

  checkShowAdvertisement() {
    let isAutoSkip = $("#auto_skip").is(":checked");
    let boxAutoSkip = $(".box-auto-skip");
    let boxTimeToSkip = $(".box-time-to-skip");
    let boxShowAutoSkipButton = $(".box-show-auto-skips-button");
    let boxNumberOfPreRoll = $(".box-number-of-pre-roll");
    let boxDelay = $(".box-delay");
    if (this.isInstream()) {
      boxAutoSkip.removeClass("d-none");
      if (isAutoSkip) {
        boxTimeToSkip.removeClass("d-none");
        boxShowAutoSkipButton.removeClass("d-none");
        boxNumberOfPreRoll.removeClass("d-none");
        boxDelay.addClass("d-none");
      } else {
        boxTimeToSkip.addClass("d-none");
        boxShowAutoSkipButton.addClass("d-none");
        boxNumberOfPreRoll.addClass("d-none");
        boxDelay.removeClass("d-none");
      }
    } else {
      boxAutoSkip.addClass("d-none");
      boxTimeToSkip.addClass("d-none");
      boxShowAutoSkipButton.addClass("d-none");
      boxNumberOfPreRoll.addClass("d-none");
      boxDelay.removeClass("d-none");
    }
  }

  checkShowColumnsPosition() {
    let layout = $("#player_layout").val();
    let box = $(".box-columns-position");
    let positionRight = $("#columns_position_right");
    if (this.isInstream() && (layout === "6" || layout === "7")) {
      box.removeClass("d-none");

      if (this.isEdit && this.isFirst) {
        return;
      }
      positionRight.prop("checked", true);
    } else {
      box.addClass("d-none");
    }
  }

  checkShowColumnsNumber() {
    let layout = $("#player_layout").val();
    let box = $(".box-columns-number");
    let columnsTwo = $("#columns_number_two");
    if (this.isInstream() && layout === "6") {
      box.removeClass("d-none");
      if (this.isEdit && this.isFirst) {
        return;
      }
      columnsTwo.prop("checked", true);
    } else {
      box.addClass("d-none");
    }
  }

  eventChangePlayerMode() {
    $(".play_mode").on("click", function () {
      $(".play_mode").removeClass("pt52");
      $(".box-play-mode").addClass("d-none");

      $(this).addClass("pt52");
      $(".box-" + $(this).data("type")).removeClass("d-none");
    });
  }

  eventClickTab() {
    $(".template-tab").on("click", ".nav-link", function () {
      $(".template-tab").find(".nav-link").removeClass("pp-4 at-1");
      $(this).addClass("pp-4");
      if ($(this).attr("data-tab") !== "1") {
        $(this).addClass("at-1");
      }
    });
  }

  eventChangeVideoLayout() {
    let thisClass = this;
    $(".slide").on("click", function () {
      $(".slide").removeClass("selected");
      $(this).addClass("selected");

      $(thisClass.videoPlayerLayoutClass).val($(this).data("value"));

      switch (true) {
        case thisClass.isOutstream():
          $("#tab-customize").addClass("d-none");
          $(".no-preview-available").removeClass("d-none");
          $("#videocontainer").empty();
          thisClass.resizePreviewNative();
          break;
        case thisClass.isInstream():
          $(".no-preview-available").addClass("d-none");
          $("#tab-customize").removeClass("d-none");
          break;
      }

      thisClass.firstLoad();
    });
  }

  isOutstream() {
    return $(this.playTypeId).val() === "9" && $(this.videoPlayerLayoutClass).val() === "9";
  }

  isInstream() {
    return $(this.playTypeId).val() === "9" && ["1", "2", "3", "4", "5", "6", "7", "8"].includes($(this.videoPlayerLayoutClass).val());
  }

  eventChangePlayerType() {
    $(this.playTypeId).on("change", () => {
      this.changePlayerTypeAction();
    });
  }

  eventChangePlayerSize() {
    $("input[type=radio][name=size]").change(() => {
      this.changePlayerSizeAction();
    });
  }

  changePlayerSizeAction() {
    let type = $(`input[type=radio][name=size]:checked`).val();
    let boxWidth = $(".box-width");
    let boxRatio = $(".box-ratio");

    switch (type) {
      //responsive
      case "1":
        boxWidth.addClass("d-none");
        boxRatio.removeClass("d-none");
        break;
      //fixed
      case "2":
        boxRatio.addClass("d-none");
        boxWidth.removeClass("d-none");
        break;
    }
  }

  eventChangeNativeAdStyle() {
    $(this.nativeAdStyleId).change(() => {
      this.changeNativeAdStyleAction();
      this.updateNativePreview()
    });
  }

  eventChangeNativeOptimizeLayout() {
    $("input[name=mode]").change(() => {
      this.changeNativeOptimizeLayoutAction();
      this.updateNativePreview()
    });
  }

  eventChangeNativeMultiDimension() {
    $(".multi-row-input" + "," + ".multi-column-input").change(() => {
      this.updateNativePreview();
    });
  }

  eventChangeNativeAdSize() {
    $("#ad-size").change(() => {
      this.updateNativePreview()
    });
  }

  changeNativeOptimizeLayoutAction() {
    switch ($(`input[name=mode]:checked`).val()) {
      case "auto":
        $(".box-optimize-layout").addClass("d-none");
        break;
      case "fixed":
        $(".box-optimize-layout").removeClass("d-none");
        break;
    }
  }

  changeNativeAdStyleAction() {
    switch (true) {
      case this.isMultiplexAdNative():
        this.showMultiplexAdNative();
        break;
      case this.isSingleAdNative():
        this.showSingleAdNative();
        break;
    }
  }

  isMultiplexAdNative() {
    return $(this.nativeAdStyleId).val() === "multiple";
  }

  isSingleAdNative() {
    return $(this.nativeAdStyleId).val() === "single";
  }

  updateNativePreview() {
    console.log("updateNativePreview");
    $(".adsbypocpoc").remove();
    $("#adsPlacement").append('<div class="adsbypocpoc" data-ad-slot="preview"></div>');

    var background = $("#native_background").val();
    var colorAdvertiser = $("#native_advertiser_name").val();
    var colorButtonCTA = $("#native_cta_button").val();
    var colorHeadline = $("#native_title_color").val();
    var columns = parseInt($(".multi-column-input").val());
    var rows = parseInt($(".multi-row-input").val());
    var size = $("#ad-size").find("option:selected").text();
    var mode = $('input[name="mode"]:checked').val();
    var type = $("#ad_style").val();
    var template = "grid";

    switch (true) {
      case this.isMultiplexAdNative():
        // mode = "auto"
        size = "1x1";
        template = "grid";
        break;
      case this.isSingleAdNative():
        columns = 1;
        rows = 1;
        template = "standard";
        break;
    }

    try {
      (ppocTag.Init = window.ppocTag.Init || []).push(function () {
        ppocAPITag.previewTemplate({
          "configs": {
            "appearance": {
              "background": background,
              "colorAdvertiser": colorAdvertiser,
              "colorButtonCTA": colorButtonCTA,
              "colorHeadline": colorHeadline,
              "sponsoredBrand": false
            },
            "layout": {
              "template": template,
              "columns": columns,
              "rows": rows,
              "size": size
            },
            "mode": mode,
            "type": type
          }
        });
      })
    } catch (error) {
      console.error(error);
    }

    console.log({
      appearance: {
        background: background,
        colorAdvertiser: colorAdvertiser,
        colorButtonCTA: colorButtonCTA,
        colorHeadline: colorHeadline,
        sponsoredBrand: false
      }, "layout": {
        "template": template, "columns": columns, "rows": rows, "size": size
      }, "mode": mode, "type": type,
    });

    this.setPreviewSize();
  }

  showMultiplexAdNative() {
    $(".box-style-multiple").removeClass("d-none");
    $(".box-style-single").addClass("d-none");
    $("#preview_size").removeClass("d-none");
  }

  showSingleAdNative() {
    $(".box-style-multiple").addClass("d-none");
    $(".box-style-single").removeClass("d-none");
    $("#preview_size").addClass("d-none");
  }

  changePlayerTypeAction() {
    switch (true) {
      case this.isNative():
        this.showNativeOptions();
        break;
      case this.isVideo():
        this.showVideoOptions();
        break;
    }
    this.makePreview();
  }

  isNative() {
    return $(this.playTypeId).val() === "8";
  }

  isVideo() {
    return $(this.playTypeId).val() === "9";
  }

  showNativeOptions() {
    $("#tab-customize").removeClass("d-none");
    $(".no-preview-available").addClass("d-none");

    $(".video-template-element").addClass("d-none");
    $(".native-template-element").removeClass("d-none");
  }

  showVideoOptions() {
    if (this.isOutstream()) {
      $("#tab-customize").addClass("d-none");
      $(".no-preview-available").removeClass("d-none");
    }

    $(".video-template-element").removeClass("d-none");
    $(".native-template-element").addClass("d-none");

  }

  abortXhrMakePreview() {
    // Hủy lệnh gọi AJAX trước đó
    if (this.xhrMakePreview && this.xhrMakePreview.readyState !== 4) {
      this.xhrMakePreview.abort();
    }
  }

  makePreview() {
    this.resizePreviewNative();

    let adType = $(this.playTypeId).val();
    switch (true) {
      case adType === "8":
        this.updateNativePreview()
        break;
      case this.isInstream():
        this.abortXhrMakePreview();

        let postData = $(this.formId).serializeObject();
        postData = this.makeData(postData);

        this.xhrMakePreview = $.ajax({
          url: this.ajaxPreview,
          type: "POST",
          dataType: "JSON",
          contentType: "application/json",
          data: JSON.stringify(postData),
          success: (json) => {
            $("#videocontainer").empty();
            this.callVideo(json);
          },
        });
        break;
    }
  }

  callVideo(json) {
    let tryTimes = 0;
    let itv = setInterval(() => {
      if (tryTimes > 10) {
        clearInterval(itv);
        console.error("Can not init power video container");
      }

      if (typeof viAPItag == "undefined") {
        tryTimes++;
      } else {
        clearInterval(itv);
        viAPItag.initPowerVideoContainer(json);
        this.setPreviewSize();
      }
    }, 100);
  }

  makeData(postData) {
    if (this.isEdit && postData.id) {
      postData.id = parseInt(postData.id);
    }
    postData.type = parseInt(postData.type);
    postData.vast_retry = parseInt(postData.vast_retry);
    postData.delay = parseInt(postData.delay);
    postData.time_to_skip = parseInt(postData.time_to_skip);
    postData.max_width = parseInt(postData.max_width);
    postData.width = parseInt(postData.width);
    postData.floating_width = parseInt(postData.floating_width);
    postData.floating_width_mobile = parseInt(postData.floating_width_mobile);
    postData.margin_bottom_desktop = parseInt(postData.margin_bottom_desktop);
    postData.margin_bottom_mobile = parseInt(postData.margin_bottom_mobile);
    postData.margin_left_desktop = parseInt(postData.margin_left_desktop);
    postData.margin_left_mobile = parseInt(postData.margin_left_mobile);
    postData.margin_right_desktop = parseInt(postData.margin_right_desktop);
    postData.margin_right_mobile = parseInt(postData.margin_right_mobile);
    postData.margin_top_desktop = parseInt(postData.margin_top_desktop);
    postData.columns_position = parseInt(postData.columns_position);
    postData.columns_number = parseInt(postData.columns_number);
    postData.player_layout = parseInt(postData.player_layout);
    postData.show_auto_skip_button = parseInt(postData.show_auto_skip_button);
    postData.number_of_pre_roll_ads = parseInt(postData.number_of_pre_roll_ads);
    postData.floating_position_desktop = parseInt(postData.floating_position_desktop);
    postData.floating_position_mobile = parseInt(postData.floating_position_mobile);
    postData.play_mode = parseInt(postData.play_mode);
    postData.advertisement_scenario = parseInt(postData.advertisement_scenario);
    postData.size = parseInt(postData.size);
    postData.auto_start = parseInt(postData.auto_start);
    postData.link = $("#preview-logo").attr("src");

    const onOffArr = ["close_floating_button_mobile", "close_floating_button_desktop",
      "main_title", "default_sound_mode", "description_enable", "show_controls",
      "floating_on_desktop", "floating_on_mobile", "sub_title", "action_button",
      "title_enable", "pubpower_logo", "powered_by", "powered_by_top_article",
      "share_button", "video_config", "show_stats", "fullscreen_button",
      "next_prev_arrows_button", "next_prev_time", "enable_logo", "auto_skip",
      "float_on_bottom", "floating_on_view", "float_on_bottom_mobile",
      "floating_on_view_mobile", "floating_on_impression", "floating_on_ad_fetched",
      "floating_on_ad_fetched_mobile", "wait_for_ad", "pre_roll", "mid_roll", "post_roll"
    ];

    onOffArr.forEach(element => {
      if (postData[element] === "on") {
        postData[element] = 1;
      } else {
        postData[element] = 2;
      }
    });

    postData.custom_logo_top_article = 1;
    postData.custom_logo = 1;

    postData.enable_logo_top_article = postData.enable_logo;

    // Native
    if (postData.type == "8" || postData.type == 8 ||
      (this.isEdit && $(this.playTypeId).attr("data-type") === "8")) {
      postData.title_color = $("#native_title_color").val()
      postData.background_color = $("#native_background").val()
      postData.action_button_color = $("#native_cta_button").val()
      postData.advertiser_color = $("#native_advertiser_name").val()
    }
    postData.ad_size = parseInt(postData.ad_size)
    postData.rows = parseInt(postData.rows)
    postData.columns = parseInt(postData.columns)
    return postData;
  }
}

export default General;