/*
 * @Author: gitshilly
 * @Date:   2018-03-14 16:30:27
 * @Version:   1.0
 * @Last Modified by:   qiaoli
 * @Last Modified time: 2018-06-05 09:21:42
 */
var GroupIDs = ['newscol1','newscol2'];
var MBizList = {};
var querykey = "";
var WSOn = 0;
if (window.WebSocket) {
    WSOn = 1;
}
/*注册消息处理方法*/
var OnMessage = function (event) {
    var data = $.parseJSON(event.data);
    var columnid = data.columnID;
    if (columnid == "1") {
        pageData.init();
        getTopNewsList();
    }
    if (columnid == "2") {
        pageData_jxk.init();
    }
    try {
        console.log("接收到socket消息：" + event.data);
    } catch (e) {
    }
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
    if (columnid == "1") {
        pageData.init();
        getTopNewsList();
    }
    if (columnid == "2") {
        pageData_jxk.init();
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
        pageData_jxk.init();
    }, loopUpdateDelaySeconds * 1000)
}

var xxzh_columnID = 1;//学习综合
var jxk_columnID = 2;//精选课
var StrategyID = 0;//策略id 0不区分策略

// banner模板
var source_banner = $("#bannerList-template").html();
var template_banner = Handlebars.compile(source_banner);

//置顶资讯模板
var source_top = $("#xxzhNewsTop-template").html();
var template_top = Handlebars.compile(source_top);

//非置顶资讯模板
var source_xxzh = $("#xxzhNewsList-template").html();
var template_xxzh = Handlebars.compile(source_xxzh);

//精选课模板
var source_jxk = $("#jxkNewsList-template").html();
var template_jxk = Handlebars.compile(source_jxk);

//列表点击打开客户端窗口
openWindow();

//获取banner列表
$.get(www + "panwai/bannerlist",{r:Math.random()},function (data) {
    if(data.RetCode=="0") {
        var bannerList = data.Message;
        var retObj = {BannerList: bannerList};
        var html = template_banner(retObj);
        $("#bannerlist").html(html);

        slider();
    }
});

//置顶资讯
function getTopNewsList() {
    $("#top_NewsList").html("");
    $.get(www + "panwai/topnewslist", {ColumnID: xxzh_columnID, StrategyID: 0, r: Math.random()}, function (data) {
        if (data.RetCode == "0") {
            var newsList = data.Message;
            var topnewsHTML = template_top(newsList);//返回第一条数据并排除列表中
            $("#top_NewsList").html(topnewsHTML);
        }
    });
}

function openWindow()
{
    $('#bannerlist').on("click","[flag=openWindow_banner]",function(){
        var id = $(this).data("id");
        var newsType = $(this).data("newstype");

        openwin(id,newsType);
    });
    $('.switchlab-cnt').on("dblclick","[flag=openSoftWindow]",function(){
        var id = $(this).data("id");
        var newsType = $(this).data("newstype");

        openwin(id,newsType);
    }).on("click","[flag=openSoftWindow]",function () {
        $("[flag=openSoftWindow]").removeClass("on");
        $(this).addClass("on");
    });
    function openwin(id,newsType) {
        var newsUrl = www + "panwai/article?newsid="+id;
        if(newsType == "2"){
            newsUrl = www + "panwai/serieslesson?newsid="+id;
        }
        PC_JH("EM_FUNC_DATAINFO_DATAIL","{|*详情页*|}{|*" + newsUrl + "*|}");
    }
}
// resize
resize();
$(window).resize(function(event) {
/* Act on the event */
    resize();
});

function resize(event) {
    /* Act on the event */
    $('.courselist').each(function (index, el) {
        $(el).height($(window).height() - $('.focuswrapper').height() - $('.switchlab-tit').height() - $('.courseitemTop', $(el).parent().parent()).height() - 3);
    });
}

// banner 图切换
function slider(){
var todo = 0;
var maintimer = null;
var focusItem = $(".focus ul li");
var btnItem = $(".focus ol li");
focusItem.hide();
focusItem.eq(0).show();
btnItem.eq(0).addClass("active");
var length = focusItem.length;
var curr = 0;

btnItem.mouseover(function(evt) {
    var kidx = btnItem.index(this);
    curr = kidx;
    focusItem.eq(kidx).show().siblings("li").hide();
    btnItem.eq(kidx).addClass("active").siblings().removeClass("active");
});

//自动翻
function goloop() {
    maintimer = setInterval(function() {
        //length = $(".focus ul li").length;
        todo = (++curr) % length;
        focusItem.eq(todo).show().siblings("li").hide();
        btnItem.eq(todo).addClass("active").siblings().removeClass("active");
        //btnItem.eq(todo).mouseover();
    }, 5000);
}
goloop();
//鼠标悬停在触发器上时停止自动翻
$(".focus").mouseenter(function() {
    clearInterval(maintimer);
}).mouseout(function() {
    clearInterval(maintimer);
    goloop();
});
}

//初始化-建立ws连接
function initNewsList() {
    //WSOn = 0;
    var groupInfo = {};
    groupInfo.IDs = [];
    for (var i = 0; i < GroupIDs.length; i++) {
        groupInfo.IDs.push(GroupIDs[i] + "")
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

//获取当前时间的字符串 "2018/4/3 19:5:59"
function getCurrentDate() {
    var now = new Date();
    var dateStr = now.getFullYear() + "/" + (now.getMonth() + 1) + "/" + now.getDate() + " " + now.getHours() + ":" + now.getMinutes() + ":" + now.getSeconds();
    return dateStr;
}

$(document).ready(function() {
    EM_FUNC_DOWNLOAD_COMPLETE();

    initNewsList();

    getTopNewsList();

    var $page = $("#PageIndex"); //页索引
    var $page_1 = $("#PageIndex_1"); //页索引

    // 滚动加载-学习综合
    $('#dl_NewsList').scroll(function () {
        pageData.scroll($page.val());
    });
    // 滚动加载-精选课
    $('#div_jxkNewsList').scroll(function () {
        pageData_jxk.scroll($page_1.val());
    });

    //切换（学习综合、精选课切换）
    var flag_xxzh = true;
    var flag_jxk = true;
    $(".switchlab-tit ul li").click(function () {
        var ind = $(".switchlab-tit ul li").index(this);
        if(ind == 0){
            if(flag_xxzh)
            {pageData.init();}
        }else{
            if(flag_jxk)
            {pageData_jxk.init();}
        }
    });

    $('.switchlabpack').switchlab();
    $('.switchlabpack .switchlab-tit .switchlab-tit-i:eq(0)').trigger('click');
});

// 分页获取学习综合
var pageData = {
    pagecount: 0,
    pageSize: 10,
    isCompleted: false,
    init: function() {
        $("#PageIndex").val("1");
        pageData.isCompleted = false;
        pageData.load(1);
        $("#dl_NewsList").html("");
    },
    load: function(pagenum) {
        $.ajax({
            url: www + "panwai/newslist",
            type: 'GET',
            dataType: 'json',
            data: {ColumnID: xxzh_columnID, StrategyID: StrategyID,currpage:pagenum,pageSize:pageData.pageSize,r:Math.random()},
            success: function (data) {
                HiddenLoadingDocument();
                if (data.RetCode == "0") {
                    xxzh_flag=true;
                    var newsList = data.Message;
                    if (newsList.length > 0) {
                        xxzh_flag = false;
                    }

                    var retObj = {NewsList: newsList};
                    var html = template_xxzh(retObj);
                    $("#dl_NewsList").append(html);

                    resize();

                    pageData.pagecount = data.TotalCount;
                    //pageData.appendHtml(data.data);
                    pageData.setPageIndex();
                    if (Math.ceil(pageData.pagecount/pageData.pageSize) == pagenum) {
                        pageData.isCompleted = true;
                    }else{
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
    scroll: function(page) { //滚动到底部加载数据
        if (pageData.isCompleted) {
            return false;
        }
        var top = $('#dl_NewsList').scrollTop();
        var win = $('#dl_NewsList').height();
        var doc = document.getElementById("dl_NewsList").scrollHeight;
        if ((top + win) >= doc) {
            pageData.load(page);
        }
    },
    setPageIndex: function() { //数据载入成功，设置下一页索引
        var $page = $("#PageIndex");
        var index = parseInt($page.val()) + 1;
        $page.val(index);
    },
    appendHtml: function(data) {
        var $container = $('.container');
        // handlerbar compile
        // var tpl =  $("#tpl").html();
        // var template = Handlebars.compile(tpl);
        // $container.append(template(data));
    }
}

// 分页获取精选课
var pageData_jxk = {
    pagecount: 0,
    pageSize: 12,
    isCompleted: false,
    init: function() {
        $("#PageIndex_1").val("1");
        pageData_jxk.isCompleted = false;
        pageData_jxk.load(1);
        $("#jxkNewsList").html("");
    },
    load: function(pagenum) {
        $.ajax({
            url: www + "panwai/newslist",
            type: 'GET',
            dataType: 'json',
            data: {ColumnID: jxk_columnID, StrategyID: StrategyID,currpage:pagenum,pageSize:pageData_jxk.pageSize,r:Math.random()},
            success: function (data) {
                HiddenLoadingDocument();
                if (data.RetCode == "0") {
                    var newsList = data.Message;
                    jxk_flag = true;
                    if (newsList.length > 0) {
                        jxk_flag = false;
                    }

                    var retObj = {NewsList: newsList};
                    var html = template_jxk(retObj);
                    $("#jxkNewsList").append(html);

                    resize();

                    pageData_jxk.pagecount = data.TotalCount;
                    pageData_jxk.setPageIndex();
                    if (Math.ceil(pageData_jxk.pagecount/pageData_jxk.pageSize) == pagenum) {
                        pageData_jxk.isCompleted = true;
                    }else{
                        pageData_jxk.isCompleted = false;
                    }
                }
                if (data.RetCode == "-2") {
                    NoneDataDocument($(".switchlab-cnt-i"));
                }
            },
            beforeSend: function () {
                LoadingDocument();
                // 禁用按钮防止重复提交
                pageData_jxk.isCompleted = true;
            },
            error: function (XMLHttpRequest, textStatus, errorThrown) {

            }
        });
    },
    scroll: function(page) { //滚动到底部加载数据
        if (pageData_jxk.isCompleted) {
            return false;
        }
        var top = $('#div_jxkNewsList').scrollTop();
        var win = $('#div_jxkNewsList').height();
        var doc = document.getElementById("div_jxkNewsList").scrollHeight;
        if ((top + win) >= doc) {
            pageData_jxk.load(page);
        }
    },
    setPageIndex: function() { //数据载入成功，设置下一页索引
        var $page_1 = $("#PageIndex_1");
        var index = parseInt($page_1.val()) + 1;
        $page_1.val(index);
    },
    appendHtml: function(data) {
        var $container = $('.container');
        // handlerbar compile
        // var tpl =  $("#tpl").html();
        // var template = Handlebars.compile(tpl);
        // $container.append(template(data));
    }
}