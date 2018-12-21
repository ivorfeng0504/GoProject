/*
 * @Author: gitshilly
 * @Date:   2018-03-14 16:30:27
 * @Version:   1.0
 * @Last Modified by:   zhangye
 * @Last Modified time: 2018-04-13 09:21:42
 */
var GroupID = ['newscol3'];
var MBizList = {};
var querykey = "";
var WSOn = 0;
var xxzh_columnID = 3;//策略详情-课堂
var StrategyID = ParamFunc.getParam("StrategyID");// $.fn.getUrlVar('StrategyID');
var clickflag_newsid=0;
if (window.WebSocket) {
    WSOn = 1;
}
/*注册消息处理方法*/
var OnMessage = function (event) {
    var data = $.parseJSON(event.data);
    var columnid = data.columnID;
    if (columnid == "3") {
        pageData.init();
        getTopNewsList();
    }
    try {
        console.log("接收到socket消息：" + event.data);
    } catch (e) {}
}

/*注册ws关闭事件处理办法*/
var OnClose = function (event) {
    try {
        var index = event.currentTarget.bizIndex;
        console.log('OnClose' + index);
        var biz = MBizList[index];
        biz.FailTimes++;
        if (biz.FailTimes >= biz.Times2Rollback) {
            biz.stop();
        }
    } catch (e) {}
}

/*注册ws错误事件处理办法*/
var OnError = function (event) {
    try {
        var index = event.currentTarget.bizIndex;
        console.log('OnError' + index);
        var biz = MBizList[index];
        TimerFlagCount();
        biz.FailTimes++;
        if (biz.FailTimes >= biz.Times2Rollback) {
            biz.stop();
        }
    } catch (e) {}
}


/*注册LongPoll消息处理方法*/
var OnLongPollMessage = function (event) {
    var columnid = event.data.columnID;
    if (columnid == "3") {
        pageData.init();
        getTopNewsList();
    }

    //更新消息最后时间 防止重复循环
    if (typeof (event.data.LastUpdateTime) != "undefined") {
        event.$MBiz.initlongpoll(OnLongPollMessage, OnClose, OnError, '', event.data.LastUpdateTime, GoBackProcess);
    } else {
        event.$MBiz.initlongpoll(OnLongPollMessage, OnClose, OnError, '', getCurrentDate(), GoBackProcess);
    }

    try {
        console.log(event.data);
    } catch (e) {}
}


/*返回轮询模式*/
var GoBackProcess = function () {
    console.log("切换到轮询模式")
    TimerFlagCount();
};

function TimerFlagCount() {
    setInterval(function () {
        pageData.init();
        getTopNewsList();
    }, loopUpdateDelaySeconds * 1000)
}



var source_top = $("#xxzhTopNewsList-template").html();
var template_top = Handlebars.compile(source_top);

var source = $("#xxzhNewsList-template").html();
var template = Handlebars.compile(source);

var newsInfosource = $("#NewsInfo-template").html();
var newsInfotemplate = Handlebars.compile(newsInfosource);

function lessonlist(currpage) {
    $("#ul_xxzhList").html("");
    $.ajax({
        url:www + "panwai/newslist",
        type: 'GET',
        dataType: 'json',
        data: {ColumnID:xxzh_columnID,StrategyID:StrategyID,currpage:currpage,pageSize:100},
        success: function(data){
            HiddenLoadingDocument();
            if(data.RetCode=="0") {
                var newsList = data.Message;
                //置顶资讯
                var IsTopList = newsList.filter(function(e) {return e.IsTop==true;});
                if(IsTopList.length>0)
                {
                    var topnewsHTML = template_top(IsTopList.shift());//返回第一条数据并排除列表中
                    $("#ul_xxzhList").append(topnewsHTML);
                }

                //非置顶资讯
                var NotTopList = newsList.filter(function(e) { return e.IsTop==false;});
                //根据时间排序
                var retObj = {NewsList: NotTopList.concat(IsTopList).sort(function (a,b) {
                        return a.CreateTime<b.CreateTime;
                    })};
                var html = template(retObj);
                $("#ul_xxzhList").append(html);

                //点击标题查看资讯详情
                $(".lesson-menu-i").click(function () {
                    $(".lesson-menu-i").removeClass("active");
                    $(this).addClass("active");
                    var newsid = $(this).attr("newsid");
                    $.get(www + "panwai/newsinfo",{NewsId:newsid},function(data) {
                        if(data.RetCode=="0") {
                            var retObj = data.Message;
                            var html = newsInfotemplate(retObj);
                            $("#div_article").html(html);
                        }

                        resize();
                    });
                });

                $("#ul_xxzhList li:eq(0)").trigger("click");

                resize();
            }
            if(data.RetCode=="-2"){
                NoneDataDocument($(".lesson-menu"));
            }
        },
        beforeSend: function () {
            LoadingDocument();
        },
        error: function(XMLHttpRequest, textStatus, errorThrown){

        }
    });
}

//置顶资讯
function getTopNewsList() {
    $("#top_NewsList").html("");
    $.get(www + "panwai/topnewslist", {
        ColumnID: xxzh_columnID,
        StrategyID: StrategyID,
        r: Math.random()
    }, function (data) {
        if (data.RetCode == "0") {
            var newsList = data.Message;
            var topnewsHTML = template_top(newsList);//返回第一条数据并排除列表中
            $("#top_NewsList").html(topnewsHTML);

        }
    });
}

//初始化-建立ws连接
function initNewsList() {
    //WSOn = 0;
    var groupInfo = {};
    groupInfo.IDs = [];
    for (var i = 0; i < GroupID.length; i++) {
        groupInfo.IDs.push(GroupID[i] + "")
    }
    var groupIds = encodeURI(JSON.stringify(groupInfo));
    //console.log(groupIds);
    MBizList = new MyMBiz(AppId, '', groupIds, newGuid(), 0, 0, 0);
    if (WSOn == 1) {
        var connectionid = MBizList.init(OnMessage, OnClose, OnError, undefined, undefined);
    } else {
        querykey = getCurrentDate();
        MBizList.initlongpoll(OnLongPollMessage, OnClose, OnError, '', querykey, GoBackProcess, undefined);
    }
}


$(document).ready(function() {
    EM_FUNC_DOWNLOAD_COMPLETE();

    //初始化-建立ws连接
    initNewsList();
    //置顶资讯
    getTopNewsList();
    //分页加载资讯
    pageData.init();

    var $page = $("#PageIndex"); //页索引
    // 滚动加载-学习综合
    $('#div_xxzhList').scroll(function () {
        pageData.scroll($page.val());
    });

    $(window).resize(function(event) {
        /* Act on the event */
        resize();
    });
    resize();
});

function resize(){
    var _height = $(window).height()-$('.lesson-nav').outerHeight()-4;
    $('.lesson-menu,.lesson-cnt').height(_height);
    $('.article-cnt').height(_height-$('.article-tit').outerHeight()-4-10);

    $('#div_xxzhList').height(_height-$("#top_NewsList").outerHeight());
}

//获取当前时间的字符串 "2018/4/3 19:5:59"
function getCurrentDate() {
    var now = new Date();
    var dateStr = now.getFullYear() + "/" + (now.getMonth() + 1) + "/" + now.getDate() + " " + now.getHours() + ":" + now.getMinutes() + ":" + now.getSeconds();
    return dateStr;
}

// 分页获取学习综合
var pageData = {
    pagecount: 0,
    pageSize: 25,
    isCompleted: false,
    init: function () {
        $("#PageIndex").val("1");
        pageData.isCompleted = false;
        pageData.load(1);

        $("#ul_xxzhList").html("");
    },
    load: function (pagenum) {
        $.ajax({
            url: www + "panwai/newslist",
            type: 'GET',
            dataType: 'json',
            data: {
                ColumnID: xxzh_columnID,
                StrategyID: StrategyID,
                currpage: pagenum,
                pageSize: pageData.pageSize,
                r: Math.random()
            },
            success: function (data) {
                HiddenLoadingDocument();
                if (data.RetCode == "0") {
                    var newsList = data.Message;

                    //不区分置顶非置顶
                    var retObj = {NewsList: newsList};
                    var html = template(retObj);
                    $("#ul_xxzhList").append(html);


                    //点击标题查看资讯详情
                    $(".lesson-menu-i").click(function () {
                        $(".lesson-menu-i").removeClass("active");
                        $(this).addClass("active");
                        var newsid = $(this).attr("newsid");
                        clickflag_newsid = newsid;
                        $.get(www + "panwai/newsinfo", {NewsId: newsid}, function (data) {
                            if (data.RetCode == "0") {
                                var retObj = data.Message;
                                var html = newsInfotemplate(retObj);
                                $("#div_article").html(html);
                            }

                            resize();
                        });
                    });

                    //触发click
                    if (clickflag_newsid == "0") {
                        $("#top_NewsList ul li").length ? $("#top_NewsList ul li:eq(0)").trigger("click") : $("#ul_xxzhList li:eq(0)").trigger("click");
                    } else {
                        $("[newsid=" + clickflag_newsid + "]").trigger("click");
                    }

                    resize();

                    pageData.pagecount = data.TotalCount;
                    pageData.setPageIndex();
                    if (Math.ceil(pageData.pagecount / pageData.pageSize) == pagenum) {
                        pageData.isCompleted = true;
                    } else {
                        pageData.isCompleted = false;
                    }
                }
                if (data.RetCode == "-2") {
                    NoneDataDocument($(".switchlab-cnt-i"));
                }
            },
            beforeSend: function () {
                LoadingDocument();
                // 禁用按钮防止重复提交
                pageData.isCompleted = true;
            },
            error: function (XMLHttpRequest, textStatus, errorThrown) {

            }
        });
    },
    scroll: function (page) { //滚动到底部加载数据
        if (pageData.isCompleted) {
            return false;
        }
        var top = $('#div_xxzhList').scrollTop();
        var win = $('#div_xxzhList').height();
        var doc = document.getElementById("div_xxzhList").scrollHeight;
        if ((top + win) >= doc) {
            pageData.load(page);
        }
    },
    setPageIndex: function () { //数据载入成功，设置下一页索引
        var $page = $("#PageIndex");
        var index = parseInt($page.val()) + 1;
        $page.val(index);
    },
    appendHtml: function (data) {
        var $container = $('.container');
        // handlerbar compile
        // var tpl =  $("#tpl").html();
        // var template = Handlebars.compile(tpl);
        // $container.append(template(data));
    }
}