define([
    "jquery",
    "cors", 
    "utils",
    "layer",
    "handlebars",
    "nicescroll",
    "localstorage",
    "webstorageapi"
], function ($, cors, utils, layer, Handlebars, nicescroll, localstorage, webstorageapi) {
    
    // layer 引用补丁：
    layer.config({ path: window.gConfig.staticPath + "static/libs/layer/" });

    var page = {
        layerindex: 0,
        curPageOpts: {
            NewsID: 0,
            ColumnID: 16,
            StrategyID: 0,
            LiveID:0,
            currpage: 1,
            pageSize: 10,
            LiveID: 0,
            isg: true
        },
        init: function() {
            var _this = this;
            $.support.cors = true;
            $.each(
                [
                    "#medialistScroll",
                    "#medialistScroll2",
                    "#iveNlPack",
                    "#iveNlPack2",
                    "#acScrollbox",
                    "#yaowenNlPack",
                    "#iveYwNlPack"
                ],
                function(idx, item) {
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



            _this.StrategyList(16);
            _this.hotArticlesList();
            _this.hotBigCastor();

                    window.pulldownLoad = function () {
                        var ns = null;
                        var nstimer;
                        var waite = true;
                        $(function () {
                            console.log("load$$");
                            ns = $("#medialistScroll").getNiceScroll(0);
                            ns.jqbind("#medialistScroll", "scroll", function () {
                                if (ns.newscrolly == ns.page.maxh) {

                                    if (waite) {
                                        waite = !waite;
                                        clearTimeout(nstimer);
                                        _this.curPageOpts.currpage += 1;
                                        nstimer = setTimeout(function () {
                                            _this.StrategyList();
                                            waite = !waite;
                                        },500);
                                    }
                                }
                            });
                        });
            };
            window.pulldownLoad();

        },
        // 策略首页列表
        strategyOpts: {
            ColumnID: 16,
            StrategyID: 0,
            LiveID: 0,
            currpage: 1,
            pageSize: 20
        },

        StrategyList: function(ColumnID, StrategyID, currpage, pageSize) {
            var _this = this;
            //ColumnID=16&StrategyID=1&currpage=1&pageSize=10&
            var jumpUrl = '';

            var htmlCodes = [
                "<div data-page='{{ currpage }}'>",
                "{{#each Message}}",
                "<div class='item section'>",
                '<div class="item-media"><a class="alink" href="'+ pagerouter.zhuanjiacelueList +'?NewsID={{NewsInfo.ID}}&LiveID={{StrategyInfo.LiveID}}&ColumnID=' +
                _this.curPageOpts.ColumnID +
                "&StrategyID={{StrategyInfo.ID}}&currpage=" +
                _this.curPageOpts.currpage +
                "&pageSize=" +
                _this.curPageOpts.pageSize +
                '"><img src="{{rendernewsimg NewsInfo.CoverImg }}" alt=""></a></div>',
                '<div class="item-content">',
                '	<div class="ic-iner">',
                '		<h3 class="ic-title"><a class="alink" href="'+ pagerouter.zhuanjiacelueArticle+'?NewsID={{NewsInfo.ID}}&LiveID={{StrategyInfo.LiveID}}&ColumnID=' +
                            _this.curPageOpts.ColumnID +
                            "&StrategyID=0&currpage=" +
                            _this.curPageOpts.currpage +
                            "&pageSize=" +
                            _this.curPageOpts.pageSize +
                            '">{{ NewsInfo.Title }}</a></h3>',
                        '<div class="ic-txt"><a class="alink" href="'+ pagerouter.zhuanjiacelueArticle +'?NewsID={{NewsInfo.ID}}&LiveID={{StrategyInfo.LiveID}}&ColumnID=' +
                        _this.curPageOpts.ColumnID +
                        "&StrategyID=0&currpage=" +
                        _this.curPageOpts.currpage +
                        "&pageSize=" +
                        _this.curPageOpts.pageSize +
                        '">{{rendernewsSummary NewsInfo.Summary NewsInfo.NewsContent}}</a></div>',

                        '<div class="ic-footer">',
                        '	<div class="pull-right">',
                        '		<b class="clickbtn-like" data-chanelid="{{NewsInfo.ID}}"><i class="icon-1">&#xe61b;</i><span class="clickLikes"> </span></b>',
                        '		<b>　<i class="icon-1">&#xe602;</i> {{NewsInfo.ClickNum}}</b>',
                        "	</div>",
                        "	<div><a class='alink' href='"+ pagerouter.zhuanjiacelueList +"?NewsID={{NewsInfo.ID}}&LiveID={{StrategyInfo.LiveID}}&ColumnID=" +
                                    _this.curPageOpts.ColumnID +
                                    "&StrategyID={{StrategyInfo.ID}}&currpage=" +
                                    _this.curPageOpts.currpage +
                                    "&pageSize=" +
                                    _this.curPageOpts.pageSize +
                                    "'><span>{{ StrategyInfo.StrategyName }} | <span class='vnum'>投顾编号：{{ StrategyInfo.StrategyTGNo }}</span></span></a>",
                        "   </div>",
                        "</div>",

                "</div>",
                "</div><div class='clearfloat'></div></div>",
                "{{/each}}",
                "</div>"
            ].join("");

            Handlebars.registerHelper("renderTime", function(time) {
                if (time == undefined || time == null || time.length == 0) {
                    time = "";
                } else {
                    time = time.substr(0, 10);
                }
                return new Handlebars.SafeString(time);
            });
            Handlebars.registerHelper("rendernewsimg", function(imgsrc) {
                if (imgsrc == "" || imgsrc == null || imgsrc == undefined) {
                    imgsrc = window.gConfig.staticPath + "static/images/default1.png";
                }
                return new Handlebars.SafeString(imgsrc);
            });
            Handlebars.registerHelper("rendernewsSummary", function(
                summary,
                content
            ) {
                var tempTxt ="",moreSign="",srcLen;
                if (summary == undefined || summary == null || summary == "") {
                    //summary = $("<div>" + content + "</div>").text();
                    // summary = $("<div>" + content + "</div>")
                    //     .text()
                    //     .slice(0, 65);
                    summary = " - ";
                } else {
                    tempTxt = $("<div>" + summary + "</div>")
                        .text();
                    srcLen = tempTxt.length;
                    moreSign = srcLen>65?"…":"";
                    summary = tempTxt.slice(0, 65) + moreSign;
  
                }

                return new Handlebars.SafeString(summary);
            });
            var compileNaveTpl = Handlebars.compile(htmlCodes);

            $.ajax({
                url:
                    window.gConfig.apiHost + "strategy/GetStrategyNewsList",
                    contentType: "application/json",
                    type: "GET",
                    timeout: 5000,
                    dataType: "JSON",
                    data: {
                        ColumnID: 16,
                        StrategyID: 0,
                        currpage: page.curPageOpts.currpage,
                        pageSize: page.curPageOpts.pageSize
                    },
                beforeSend: function() {},
                success: function(data) {
                    if (data.RetCode == 0 && data.Message && data.Message.length > 0) {
                        $.extend(data, { currpage: page.curPageOpts.currpage });
                        var html = compileNaveTpl(data);

                        if (page.curPageOpts.currpage == 1) {
                            $("#tacticsList").html(html);
                        } else {
                            $("#tacticsList").append(html);
                        }
                        
                        layer.close(window.layerIndex);

                        if (_this.curPageOpts.currpage > 1){
                            layer.tips('底部有新的内容加载', '#medialistScroll', {
                                offset: ['400px', '450px'],
                                tips: [2, '#cb4b4b']
                            });
                        }

                        var uid = $.cookie("expertnews.uid") || utils.GetQueryString("uid") || 0;
                        var newsId;
                        var appId = appIDnewsinfo;
                        var idDatalist = [];

                        //增加点赞处理
                        $.each($("#tacticsList [data-page=" + page.curPageOpts.currpage + "] .section"), function (idx, elm) {
                            
                            newsId = $(".clickbtn-like", elm).data("chanelid");

                            // _this.getLikesnum($(".clickLikes", elm), uid, newsId,appId)
                            idDatalist.push({
                                newsId: newsId,
                                uid: uid,
                                appId: appId
                            })
                        });

                        _this.getLikesCollect(idDatalist);

                    } else { 

                        if (page.curPageOpts.isg) {
                            layer.tips('已加载全部', '#medialistScroll', {

                                offset: ['400px', '450px'],
                                tips: [2, '#cb4b4b']
                            });
                        }

                        page.curPageOpts.isg = !1;

                    }
                },
                error: function(jqXHR, textStatus, errorThrown) {
                    //_this.errorPop('服务器发生异常，请稍后再试~')
                }
            });
        },

        getLikesCollect: function (idDatalist) { 
            var showzanData = {};
        
            $.post({
                type: "post",
                contentType: 'text/plain',
                dataType: 'json',
                url: window.gConfig.likeListDataServer,
                data: JSON.stringify(idDatalist)
            }).done(function (data) {
                var showlikesbox,curNewsid;
                if (data.isSucess) {
                     $.each(data.message, function (idx, elm) {
                        showzanData[data.message[idx]['newsId']] = {
                            liked: data.message[idx]['liked'],
                            likes: data.message[idx]['likes']
                        }
                    });

                      $.each($("#tacticsList [data-page=" + page.curPageOpts.currpage + "] .section"), function (idx, elm) {
                          
                          showlikesbox = $(".clickbtn-like", elm)
                          curNewsid = showlikesbox.attr("data-chanelid");
                          $(".clickLikes", showlikesbox).html(showzanData[curNewsid]['likes']);
                          
                      });
                }
                


            })
        },
        // 获取点赞数
        getLikesnum:function(domobj, uid, newsId, appId){
            var _this = this;
            var getUID = uid;
            var getNewsId = newsId;
            var getAppId = appId;
            var likeBtn = $(".clickbtn-like");
            var likesicon = $('.icon-1',likeBtn);
            var likesspan = domobj;
            var reqData = { uid: getUID, newsId: getNewsId, appId: getAppId };
            var reqDataStr = utils.objToQuery(reqData);
            function showLikes(data){
                 likesspan.html(data.message.likes);  
                console.log(data.message.likes);
            }
 
                $.ajax({
                    type: "get",
                    url: window.gConfig.likeDataServer, //获取点赞数和是否点赞
                    data: reqData,
                    success: function (data, status) {
                        if (status === "success") {
                            if (data.isSucess) {
                                showLikes(data);                            
                            }
                        }
                    }
                });            
        },
        // 热门资讯
        hotArticlesList: function() {
            var _this = this;

            var hotnewsOpts =  {
                ColumnID: 16,
                StrategyID: 0,
                currpage: 1,
                pageSize: 10
            };

            var htmlCodes = [
                '<ul class="ive-nl-cnt">',
                "{{#each Message}}",
                '    <li>  <a class="alink" href="'+ pagerouter.zhuanjiacelueArticle +'?NewsID={{ ID }}&ColumnID=' +
                _this.strategyOpts.ColumnID +
                "&StrategyID={{ExpertStrategyID}}&currpage=" +
                _this.strategyOpts.currpage +
                "&pageSize=" +
                _this.strategyOpts.pageSize +
                '"><i class="icon-1">&#xe637;</i>{{ Title }}</a>',
                "    </li>",
                "{{/each}}",
                "</ul>"
            ].join("");

            Handlebars.registerHelper("renderTime", function(time) {
                if (time == undefined || time == null || time.length == 0) {
                    time = "";
                } else {
                    time = time.substr(0, 10);
                }
                return new Handlebars.SafeString(time);
            });
            var compileNaveTpl = Handlebars.compile(htmlCodes);

            $.get(
                //http://test.roboadvisor.emoney.cn/expert/strategy/GetHotStrategyNewsList_Top10?ColumnID=16
                window.gConfig.apiHost +
                    "strategy/GetHotStrategyNewsList_Top10",
                hotnewsOpts,
                function(data) {
                    if (data.RetCode == 0 && data.Message.length > 0) {
                        var html = compileNaveTpl(data);
                        $("#iveNlPack").html(html);
                    }
                }
            );
        },

        // 人气大咖
        hotBigCastor: function() {
            var htmlCodes = [
                '<div class="medialist">',
                "{{#each Data}}",
                '	<div class="item">',
                '		<a class="alink" href="'+pagerouter.yqqUserLive+'?lid={{ Id }}&random=634832180" target="_blank"><div class="item-heading">',
                '			<div class="media pull-left">',
                '				<img src="{{ LiveImg }}" alt="">',
                "			</div>",
                "		</div></a>",
                '		<div class="item-content">',
                '			<div class="pull-right">',
                '				<i class="icon-1">&#xe635;</i>',
                "				<span> {{ FansNum }} </span>",
                "			</div>",
                "			<h4>",
                '			<a class="alink" href="'+pagerouter.yqqUserLive+'?lid={{ Id }}&random=634832180" target="_blank"> {{ LiveName }} </a>',
                "			</h4>",
                "       {{#unless @index}}",

                '			<div class="text"><a class="alink" href="'+pagerouter.yqqUserLive+'?lid={{ Id }}&random=634832180" target="_blank">',
                "				 {{ LiveIntro }}",
                "			</a></div>",
                "       {{/unless}}",
                "		</div>",
                '	<div class="clearfloat"></div></div>',
                "{{/each}}",
                "</div>"
            ].join("");

            Handlebars.registerHelper("isfirst", function(idx) {
                if (v1 == v2) {
                    //满足添加继续执行
                    return options.fn(this);
                } else {
                    //不满足条件执行{{else}}部分
                    return options.inverse(this);
                }
                var rst = idx == 0;
                // new Handlebars.SafeString(time)
                return rst;
            });
            var compileNaveTpl = Handlebars.compile(htmlCodes);

            var compare = function(property) {
                return function(a, b) {
                    var value1 = a[property];
                    var value2 = b[property];
                    return value2 - value1;
                };
            };

            $.get(
                window.gConfig.apiHost + "yqq/getliveroomlist",
                {
                    ColumnID: 16,
                    StrategyID: 0,
                    currpage: 1,
                    pageSize: 10
                },
                function(data) {
                    var data = JSON.parse(data);
                    var curData,html,spNum;
                    if (data.RetCode == "0" && data.Data.length > 0) {
                        data.Data = data.Data.sort(compare("FansNum"));
                        
                        spNum = ($(window).width() > 960) ? 5 : 3;

                         curData = {
                             Data: data.Data.slice(0, spNum)
                         };
                         html = compileNaveTpl(curData);
                        $("#iveNlPack1").html(html);
                    }
                }
            );
        }
    };
    page.init();
});
