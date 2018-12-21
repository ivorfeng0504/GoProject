define([
  "jquery",
  'utils',
  "handlebars",
  "nicescroll",
], function (jquery, utils, Handlebars, nicescroll) {

  // 客户端调用
  window.SetMyZXGChanged = function (context) {
    //获取客户端接口
    function GetExternal() {
      return window.external.EmObj
    };
    //调用客户端接口
    function PC_JH(type, c) {
      try {
        var obj = GetExternal();
        return obj.EmFunc(type, c)
      } catch (e) { }
    };
    // 从客户端获取股票列表
    var stocklist = "";
    try {
      stocklist = PC_JH("EM_FUNC_GET_MY_ZXG", context);
    } catch (ex) { }
    // 为空，不请求
    if (stocklist === "" || stocklist === undefined) {
      return;
    }
    setTimeout(function () {
      stocklist = stocklist.replace(/\[/g, '').replace(/]/g, '');
      page.StockList = stocklist;
      page.relationsInfoList(context);
      // 重置数据
      page.resetPageData();

    }, 300);
  };

  var page = {
    curPageOpts: {
      currpage: 0,
      pageSize: 10,
      isFinished: false
    },
    StockList: '',
    EMSSOSTR: '',
    init: function () {
      var _this = this;
      if (!!$.cookie('expertnews.uid')) {
        _this.getUID = $.cookie('expertnews.uid');
      } else {
        _this.getUID = utils.GetQueryString('uid') || 1;
        $.cookie('expertnews.uid', _this.getUID);
      }
      _this.BlockCode = utils.GetQueryString("BlockCode");
      _this.EMSSOSTR = utils.EMSSO();
      _this.getUID = $.cookie('expertnews.uid') || utils.GetQueryString('uid') || 1;
      _this.getAppId = 1015004;
      _this.handerbarsFormat();
      _this.scrollContent();

      _this.godeyesInfoList('');    // 测试

      _this.EM_HREF_JUMP(); // 跳转客户端绑定

      _this.jumpToStockHQ(); // 跳转股票行情

    },
    jumpToStockHQ: function () {
      var _this = this;
      // 绑定跳转
      $(document).on('click', '[data-stockCode]', function () {
        var stockcode = _this.preStockcode($(this).attr('data-stockCode'));
        utils.goThisStock(stockcode);
      });
    },

    handerbarsFormat: function () {
      Handlebars.registerHelper("renderTime", function (time) {
        var myDate = new Date();
        var getMonth = myDate.getMonth() + 1;
        var getMonthstr = (getMonth <= 9) ? "0" + getMonth : getMonth;
        var nowDate = myDate.getFullYear() + "-" + getMonthstr + "-" + myDate.getDate();
        if (time == undefined || time == null || time.length == 0) {
          time = "";
        } else {
          if (time.slice(0, 10) == nowDate) {
            time = time.slice(11, 16);
          } else {
            time = time.slice(5, 10) + ' ' + time.slice(11, 16);
          }

        }
        return new Handlebars.SafeString(time);
      });

      Handlebars.registerHelper("renderTimed", function (time) {
        var myDate = new Date();
        var getMonth = myDate.getMonth() + 1;
        var getMonthstr = (getMonth <= 9) ? "0" + getMonth : getMonth;
        var nowDate = myDate.getFullYear() + "-" + getMonthstr + "-" + myDate.getDate();
        if (time == undefined || time == null || time.length == 0) {
          time = "";
        } else {
          if (time.slice(0, 10) == nowDate) {
            time = time.slice(11, 16);
          } else {
            time = time.slice(5, 7) + '/' + time.slice(8, 10);
          }

        }
        return new Handlebars.SafeString(time);
      });

    },

    //自定义滚动条
    scrollContent: function () {
      var _this = this;
      _this.srcollBar = function () {
        $('.down-part').css('height', $(window).height()- 1);
        $.each(
          [".down-part"],
          function (idx, item) {
            if ($(item).length > 0) {
              $(item).niceScroll({
                cursorcolor: "#666",
                cursoropacitymax: 0.5,
                touchbehavior: false,
                cursorwidth: "6px",
                cursorborder: "0",
                cursorborderradius: "8px",
                autohidemode: true
              });
            }
          }
        );
      },
        _this.srcollBar();
      $(window).on('resize', function () {
        _this.srcollBar();
      });
    },

    //相关资讯列表
    godeyesInfoList: function (context) {
      var _this = this;
      $('#newsList').addClass('no-info-list').html('<div class="no-info">板块资讯加载中……</div>');
      var htmlCodes = [
        '{{#each Message}}',
        '<li class="news-item"><a class="news-item-link" href="#" data-ahref="' + pagerouter.articleZixun + '?newsid={{matchid DataUrl}}&source=cloud" data-title="{{ArticleTitle}}"><i class="icon-1">&#xe637;</i>{{ArticleTitle}}</a><span class="time">{{renderTimed PublishTime}}</span></li>',
        '{{/each}}'
      ].join("");
      Handlebars.registerHelper("matchid", function (dataurl) {
        var idReg = /\/(\d*?)\.data$/gm;
        var sdataurl = "" + dataurl;
        var newsid = (idReg.exec(sdataurl))[1];

        return new Handlebars.SafeString(newsid);
      });
      $.ajax({
        url: window.gConfig.myoptUrl + 'blocknews/getblocknews',
        type: 'GET',
        timeout: ajaxTimeout,
        dataType: 'JSON',
        data: { BlockCode: _this.BlockCode, PageSize: 30 },
        success: function (data) {
          if (data.RetCode == 0 && data.Message && data.Message.length > 0) {
            var $newsList = $('#newsList');
            var template = Handlebars.compile(htmlCodes);
            var html = template(data);
            $newsList.removeClass('no-info-list').html(html);
          } else {
            $('#newsList').addClass('no-info-list').html('<div class="no-info">请添加自选股，追踪您的个股动态</div>');
          }
        },
        error: function (jqXHR, textStatus, errorThrown) {
          //_this.errorPop('服务器发生异常，请稍后再试~')
        }
      });
    },

    EM_PC_JSBRIDGE: function (type, c) {
      //获取客户端接口
      function GetExternal() {
        return window.external.EmObj
      };

      //调用客户端接口
      function PC_JH(type, c) {
        try {
          var obj = GetExternal();
          return obj.EmFunc(type, c)
        } catch (e) { }
      };

      PC_JH(type, c);
    },

    EM_HREF_JUMP: function () {
      var _this = this;

      //页面加载完成
      (function () {

        $('#newsList').on("dblclick", ".news-item", function () {
          var $this = $(this);
          $this.addClass("actived");
          if ($this.hasClass('disabled')) {
            return;
          }
          $this.addClass('disabled');
          var $child = $this.find('.news-item-link');
          var _href = $child.attr("data-ahref");
          var _title = $child.attr("data-title");
          var _aimHREF = _href + "&";
          try {
            // EM_FUNC_GOODS_NEWSINFO   {|*Title*|}{|*url*|}
            _this.EM_PC_JSBRIDGE("EM_FUNC_GOODS_NEWSINFO", "{|*" + _title + "*|}{|*" + _aimHREF + "*|}");
          } catch (ex) { }
          $this.removeClass('disabled');
        })

      })();
    }
  };
  page.init();
});