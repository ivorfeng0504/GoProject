// 埋点对象
var clickData = {},
doc = window.document,
EMoneyAnalytics = {

  Init: function (App, Module, Remark) {
    clickData.App = App;
    clickData.Module = Module;
    clickData.Remark = Remark;
    
    $(document).on('click', '[clickkey]', function (event) {
      event.stopPropagation();
      var $this = $(this);
        clickData._clickkey = $this.attr("clickkey");
        clickData._clickdata = $this.attr("clickdata");
        clickData._clickremark = $this.attr("clickremark");
        clickData._htmltype = '';
        clickData.senddate = $this.attr("data-senddate");
      var lastDateforKey = clickData['date' + clickData._clickkey];
      var currDate = new Date().getTime();
      // 时间差10s以上或者首次请求，发起请求
      if (new Date(lastDateforKey).getTime() - currDate > 1000 * 10 || !lastDateforKey) {
        clickData['date' + clickData._clickkey] = currDate;
        $this.attr("data-senddate",currDate)
        EMoneyAnalytics.sendRequest();
      }
    });
  },
  sendRequest: function () {
    var Host = "http://api2.tongji.emoney.cn";
    var ClickUrl = Host + "/Page/PageClick";
    var PageViewUrl = Host + "/Page/PageView";
    var pageUrl = window.top.location.href;
    // pageUrl = pageUrl.replace(window.location.search, '');
    // 还需对比下时间
    if (clickData.App != "" && clickData._clickdata != "") {
      var src = ClickUrl + "?v=" + Math.random()
        + "&app=" + clickData.App
        + "&module=" + clickData.Module
        + "&clickkey=" + clickData._clickkey
        + "&clickdata=" + clickData._clickdata
        + "&clickremark=" + clickData._clickremark
        + "&htmltype=" + clickData._htmltype
        + "&pageurl=" + encodeURIComponent(pageUrl)
        + "&remark=" + clickData.Remark;
      var elm = document.createElement("script");
      elm.src = src;
      elm.style.display = "none";
      document.body.appendChild(elm);
    }
  }
};
