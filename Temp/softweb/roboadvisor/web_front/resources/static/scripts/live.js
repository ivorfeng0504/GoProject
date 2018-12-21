var noContentHtml = "<p style=\"margin-top:10px;\">直播暂未开始，请耐心等待</p>";
//当前视觉焦点标识
//0 直播  1所有问答 2仅看自己
var viewFocus="0";
//页面定位-直播内容
var VF_Live="0"
//页面定位-所有问答
var VF_All_Question="1"
//页面定位-仅看自己
var VF_My_Question="2"

//获取直播内容
//isIncre 是否增量 1是 0不是
function getLiveContent(isIncre) {
    date = $("#txt_newst_content_date").val();
    $.post(www + "live/livecontent", {
        IsIncre: isIncre,
        Date: date
    }, function (data) {
        if (data == null) {
            console.log("查询直播内容失败，响应内容为空")
            return;
        }
        if (data.RetCode == 0) {

            if (data.RetMsg != "SUCCESS") {
                //直播未开始时 开始倒计时 到点后刷新页面
                var clientRefreshCountDown = parseInt(data.RetMsg);
                if (isNaN(clientRefreshCountDown) == false) {
                    console.log("直播未开始，将在" + clientRefreshCountDown + "毫秒后刷新页面")
                    setTimeout(function () {
                        window.location.reload();
                    }, clientRefreshCountDown)
                }
            }

            //直播内容列表
            var $contentdiv = $("#div_live_content_list");
            if (data.Message != null && data.Message != undefined && data.Message.length > 0) {
                var newestDate = data.Message[data.Message.length - 1].CreateTime
                $("#txt_newst_content_date").val(newestDate);
                var html = template_live_content(data)
                if ($contentdiv.html().trim() == noContentHtml) {
                    $contentdiv.html("");
                }
                $contentdiv.append(html);
                //初始化置顶消息 非增量获取的时候
                if (isIncre != "1"){
                    showTopLiveContent(data.Message);
                }
                resize($contentdiv);
                //初始化图片点击事件
                initImgClickEvent();
                //增量的时候 滚动到底部
                if (isIncre == "1"||IsHistory()==false) {
                    $(".modulelist").prop("scrollTop", $(".modulelist").prop("scrollHeight"))
                }
                //初始化股票点击事件
                $(".stock_info").unbind("click").click(OpenStock);
                console.log("当前直播内容最新时间为：" + newestDate)
            } else {
                //如果全量获取的时候 没有直播内容则给予文字提示
                if (isIncre == "0") {
                    $contentdiv.append(noContentHtml);
                }
            }
            console.log("查询直播内容成功：" + JSON.stringify(data))
        } else {
            //如果全量获取的时候 没有直播内容则给予文字提示
            if (isIncre == "0") {
                $contentdiv.append(noContentHtml);
            }
            console.log("查询直播内容失败：" + JSON.stringify(data))
        }

    });
}

// getTopContent 获取置顶的直播消息
function getTopContent() {
    date = $("#txt_newst_content_date").val();
    $.post(www + "live/livecontent", {
        IsQueryTop: "1",
        Date: date
    }, function (data) {
        if (data == null) {
            console.log("查询置顶直播内容失败，响应内容为空")
            return;
        }
        if (data.RetCode == 0) {
            //初始化置顶消息
            showTopLiveContent(data.Message);
            console.log("查询置顶直播内容成功：" + JSON.stringify(data))
        } else {
            console.log("查询置顶直播内容失败：" + JSON.stringify(data))
        }
    });
}

//获取当前时间的字符串 "2018/4/3 19:5:59"
function getCurrentDate() {
    var now = new Date();
    var dateStr = now.getFullYear() + "/" + (now.getMonth() + 1) + "/" + now.getDate() + " " + now.getHours() + ":" + now.getMinutes() + ":" + now.getSeconds();
    return dateStr;
}

//初始化指定的直播消息展示
function showTopLiveContent(liveContentList) {
    if (liveContentList == null || liveContentList == undefined || liveContentList.length == 0) {
        $("#div_live_top").html("");
        return
    }
    //支持多个置顶消息
    var contentList = [];
    $(liveContentList).each(function (index, item) {
        if (item.IsTop && item.Content != null && item.Content.length > 0) {
            contentList.push(item);
        }
    });

    var contentHtml = $(template_live_content_top({
        Message: contentList
    }));
    $("img", contentHtml).remove();
    $("#div_live_top").html(contentHtml);
    $("#div_live_top").show();
    resize($("#div_live_content_list"));
}

//查看历史直播
function queryHistoryLiveContent() {
    var date = $dp.cal.getDateStr('yyyy-MM-dd');
    var url = www + "live/index?Date=" + date;
    if (viewFocus!=undefined){
        url+="&VF="+viewFocus;
    }
    var keys = ParamFunc.getKeys();
    $(keys).each(function (index, key) {
        if (key != "Date"&&key!="VF") {
            url += "&" + key + "=" + ParamFunc.getParam(key)
        }
    })
    window.location.href = url;
}

//获取问答列表
//isSelf 是否只看自己的1 or 0
function getQuestionList(isSelf) {
    date = $("#txt_newst_content_date").val();
    $.post(www + "live/livequestion", {
        Date: date,
        IsSelf: isSelf
    }, function (data) {
        if (data == null) {
            console.log("查询问答列表失败，响应内容为空")
            return;
        }
        if (data.RetCode == 0) {
            if (data.Message != null && data.Message != undefined && data.Message.length > 0) {
                var html = template_html_live_answer(data)
                //问答列表
                $("#div_live_question_list").html(html);
                resize($("#div_live_question_list"));
            } else {
                //没有问答内容
                $("#div_live_question_list").html("");
            }
        } else {
            console.log("查询问答列表失败：" + JSON.stringify(data))
        }

    });
}

//轮训更新页面数据
function TimerFlagCount() {
    setInterval(function () {
        //更新直播内容
        getLiveContent("1");
        //更新全部问答
        $("#btn_all_question").click();
    }, loopUpdateDelaySeconds * 1000)
}


//更新直播内容或者用户问答模块
function updateLiveContentOrQuestion(data) {
    if (data == undefined || data == null) {
        return
    }
    data.ContentNum = parseInt(data.ContentNum)
    if (data.ContentNum > 0) {
        //延迟拉取最新数据
        setTimeout(function () {
            getLiveContent("1");
        }, liveContentUpdateDelaySeconds * 1000);
    }
    data.RepliedNum = parseInt(data.RepliedNum)
    if (data.RepliedNum > 0) {
        //延迟拉取最新数据
        setTimeout(function () {
            $("#btn_all_question").click();
        }, liveQuestionUpdateDelaySeconds * 1000);
    }
    //更新置顶消息
    if (data.SetTop != undefined) {
        //置顶或取消置顶
        //延迟拉取最新数据
        setTimeout(function () {
            getTopContent();
        }, liveTopContentUpdateDelaySeconds * 1000);
    }
}

var RoomIdList = [];
var MBizList = [];
var querykey = "";
var WSOn = 0;
if (window.WebSocket) {
    WSOn = 1;
}
//消息认证token
var MessageToken = "";
/*注册消息处理方法*/
var OnMessage = function (event) {
    try {
        updateLiveContentOrQuestion(JSON.parse(event.data))
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
    try {
        console.log(event.data);
        updateLiveContentOrQuestion(event.data);

        //更新消息最后时间 防止重复循环
        if (typeof (event.data.LastUpdateTime) != "undefined") {
            event.$MBiz.initlongpoll(OnLongPollMessage, OnClose, OnError, '', event.data.LastUpdateTime, GoBackProcess, MessageToken);
        } else {
            event.$MBiz.initlongpoll(OnLongPollMessage, OnClose, OnError, '', getCurrentDate(), GoBackProcess, MessageToken);
        }

    } catch (e) {}
}


/*返回轮询模式*/
var GoBackProcess = function () {
    console.log("切换到轮询模式")
    TimerFlagCount();
}


//初始化提问按钮
function initQuestionSubmitButton() {
    var questionLimit = parseInt($.cookie("questionLimit"));
    var submitBtn = $("#btn_submit_question");
    if (questionLimit > 0) {
        submitBtn.addClass('disabled');
        var handler = setInterval(function () {
            var questionLimit = parseInt($.cookie("questionLimit"));
            questionLimit--;
            $.cookie("questionLimit", questionLimit);
            submitBtn.html(questionLimit + "秒后提问");
            if (questionLimit <= 0) {
                submitBtn.html("提问");
                submitBtn.removeClass("disabled");
                clearInterval(handler)
            }
        }, 1000)
    }
}


//自定义的弹窗
function liveAlert(str) {
    $(".poptips-c").text(str);
    $('.popmask,.poptips').show();
}

//初始化用户直播列表
function initRoomIdList() {
    $.post(www + "live/getuserroomlist", {}, function (data) {
        if (data == null) {
            RoomIdList = [];
            console.log("查询用户直播间列表异常，响应内容为空");
        }
        if (data.RetCode == 0) {
            if (data.Message != null && data.Message != undefined&&data.Message.RoomList!=undefined && data.Message.RoomList.length > 0) {
                RoomIdList = data.Message.RoomList;
                MessageToken=data.Message.Token;
                //如果正确解析出日期并且为历史直播 则不进行监听
                if (IsHistory() == false) {
                    //初始化直播内容监听
                    initLiveListening();
                }
            }
        } else {
            liveAlert(data.RetMsg)
        }
    });
}

function IsHistory() {
    //如果正确解析出日期并且为历史直播
    var dateStr = ParamFunc.getParam("Date")
    var now = new Date();
    var date = Date.parse(dateStr)
    if (isNaN(date) == false) {
        date = new Date(date)
        if (date.getFullYear() == now.getFullYear() && date.getMonth() == now.getMonth() && date.getDate() == now.getDate()) {
            return false;
        }
        return true;
    } else {
        return false;
    }
}

//初始化直播内容监听
function initLiveListening() {
    //WSOn = 0;
    if (RoomIdList != null && RoomIdList.length > 0) {
        var groupInfo = {};
        groupInfo.IDs = [];
        for (var i = 0; i < RoomIdList.length; i++) {
            groupInfo.IDs.push(RoomIdList[i] + "")
        }
        var groupIds = encodeURI(JSON.stringify(groupInfo));
        console.log(groupIds)
        MBizList[i] = new MyMBiz(AppId, "", groupIds, newGuid(), i, 0, 0);
        if (WSOn == 1) {
            var connectionid = MBizList[i].init(OnMessage, OnClose, OnError, undefined, MessageToken);
        } else {
            querykey = getCurrentDate();
            MBizList[i].initlongpoll(OnLongPollMessage, OnClose, OnError, '', querykey, GoBackProcess, MessageToken);
        }
    }
}

//打开日历
function openWdatePicker() {
    WdatePicker({el:'d12',position:{right:10},onpicked:queryHistoryLiveContent})
	$("body div:last").hover(function(){
		$(this).attr("data-flag",1);
	},function(){
		var $obj = $(this);
		$obj.attr("data-flag",0);
		setTimeout(function(){
			if($obj.attr("data-flag")=="0"){
				$obj.hide();
			}
		},600);
	});
}

//初始化默认打开的问答
function initAskFocus() {
    //1打开所有问答  2打开仅看自己  其他则不处理
     viewFocus=ParamFunc.getParam("VF")
    if (viewFocus==VF_All_Question){
        $(".switchlab-tit-i").eq(1).click();
        return
    }
    if(viewFocus==VF_My_Question){
        $(".switchlab-tit-i").eq(1).click();
        $("#btn_myself_question").click();
    }
}

//自动定位到直播页面的某个栏目
function liveViewFocus(vf) {
    if (vf==VF_All_Question){
        $(".switchlab-tit-i").eq(1).click();
    }else if(vf==VF_My_Question){
        $(".switchlab-tit-i").eq(1).click();
        $("#btn_myself_question").click();
    }else{
        $(".switchlab-tit-i").eq(0).click();
    }
}

//初始化鼠标放到置顶消息上显示全部内容
function initTopContentLimitDisplay() {
    $('.moduleitemTop').on("hover",".module-item",function(){
        $(this).toggleClass('clamp');
    });
}
/*****************************************函数执行****************************************/
$(function () {
    //页面加载完成
    EM_FUNC_DOWNLOAD_COMPLETE();

    //从地址中获取Date参数
    var dateStr = ParamFunc.getParam("Date")
    var dateStrComp=dateStr;
    var now = new Date();
    if(dateStrComp&&dateStrComp.length>0){
        dateStrComp=dateStrComp.replace(/-/ig,"/")
    }
    var date = Date.parse(dateStrComp)
    if (isNaN(date) == false) {
        date = new Date(date)
        $("#txt_newst_content_date").val(dateStr);
        $(".his_date span").text(dateStr);

    } else {
        date = now;
        dateStr=date.getFullYear() + "-" + (date.getMonth() + 1) + "-" + date.getDate();
        $("#txt_newst_content_date").val(dateStr);
        $(".his_date span").text(dateStr);
    }

    //如果是历史直播 不展示回答输入框
    if (date.getFullYear() == now.getFullYear() && date.getMonth() == now.getMonth() && date.getDate() == now.getDate()) {
        $("#area_question").css({
            visibility: 'visible'
        });
    } else {
        $("#area_question").css({
            visibility: 'hidden'
        });
    }

    //初始化用户直播列表
    initRoomIdList();

    //初始化提问按钮
    initQuestionSubmitButton();

    //初始化鼠标放到置顶消息上显示全部内容
    initTopContentLimitDisplay();

    //页面打开时请求最新的直播内容
    getLiveContent("0");

    //提交问答
    $("#btn_submit_question").click(function () {
        var $this = $(this);
        
        if($this.hasClass("disabled")){
          return;
        }
        
        var question = $("#txt_question").val();
        if (question.length == 0) {
            liveAlert("请填写问题内容！");
            return;
        }
        var maxLen=200;
        if(question.length>maxLen){
            liveAlert("提问内容不能超过200个字！")
            return;
        }
        
        $this.addClass('disabled');

        //限制提问时间
        $.cookie("questionLimit", questionSubmitLimitSeconds);
        $this.html(questionSubmitLimitSeconds + "秒后提问");
        var handler = setInterval(function () {
            var questionLimit = parseInt($.cookie("questionLimit"));
            questionLimit--;
            $.cookie("questionLimit", questionLimit);
            $this.html(questionLimit + "秒后提问");
            if (questionLimit <= 0) {
                $this.html("提问");
                $this.removeClass('disabled');
                clearInterval(handler)
            }
        }, 1000)

        $.post(www + "live/addquestion", {
            AskContent: question
        }, function (data) {
            if (data == null) {
                liveAlert("提交失败，请稍后再试！")
                return;
            }
            if (data.RetCode == 0) {
                $("#btn_myself_question").click();
                $("#txt_question").val("");
            } else {
                liveAlert(data.RetMsg);
            }
        });

    });

    //是否已经点击过问答
    var isQuestionTabClick = false;

    //点击问答Tab 第一次点击的时候才手动刷新问答
    $(".switchlab-tit-i").click(function () {
        var $this = $(this);
        var tabTitle = $this.text();
        if (tabTitle == "问答" && isQuestionTabClick == false) {
            //请求最新的问答内容
            getQuestionList("0");
            //已经点击过问答了
            isQuestionTabClick = true;
        }
        if(tabTitle == "问答"){
            viewFocus=VF_All_Question;
        }else{
            viewFocus=VF_Live;
        }
    });


    //点击所有问答
    $("#btn_all_question").click(function () {
        $("#btn_myself_question").removeClass("active");
        $(this).addClass("active");
        $("#txt_is_self").val("0");
        getQuestionList("0")
        viewFocus=VF_All_Question;
    });

    //点击仅看自己
    $("#btn_myself_question").click(function () {
        $("#btn_all_question").removeClass("active");
        $(this).addClass("active");
        $("#txt_is_self").val("1");
        getQuestionList("1");
        viewFocus=VF_My_Question;
    });

    //初始化问答展示
    initAskFocus();

    //弹窗关闭
    $('.btnok,.pop-close').click(function (event) {
        /* Act on the event */
        $('.popmask,.poptips').hide();
    });
});
/*****************************************函数执行****************************************/