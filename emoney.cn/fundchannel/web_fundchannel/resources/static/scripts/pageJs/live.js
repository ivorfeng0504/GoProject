var $Live = function () {
    var Times2Rollback = 2;
    var WSOn = 0;
    var ClientId = Common.NewGuid();
    var roomInfo = null;
    var LookDate = null;

    //获取直播间信息
    function GetRoomInfo() {
        Common.CommonGet(LiveRoomUrl, function (response) {
            if (response.RetCode !== "0") {
                layer.msg("获取直播室信息失败");
                console.log(response.RetMsg);
                return;
            }
            var data = response.Data;
            roomInfo = data;
            $("#todayTopic").html(data.TopicName);
            Common.CommonBind(data, "liveRoomInfo", false);
        });
    }

    function bindRoomInfo(data) {
        if (!data.TopicName) {
            data.TopicName = "今日无话题";
        }
        roomInfo = data;

        $("#todayTopic").html(data.TopicName);
        Common.CommonBind(data, "liveRoomInfo", false);
    }

    //获取全部直播内容
    function GetLiveContent(date) {
        if (!date) {
            date = Common.GetNowFormatDate();
        }
        Common.CommonGet(LiveContentUrl + date, function (response) {
            if (response.RetCode !== "0") {
                layer.msg("获取直播室内容失败");
                console.log(response.RetMsg);
                return;
            }
            var data = response.Data;
            bindLiveContent(data, "");
        });
    }

    function bindLiveContent(data, date) {
        $.each(data, function () {
            var item = this;
            if (item.LiveFrom == 8) {
                item.ShortTime = '<span class="fred">置顶消息</span>';
            }
            else {
                item.ShortTime = date + Common.GetShortTime(item.LastModifyTime);
            }
            if (item.Content.indexOf("【投资顾问") > 0) {
                item.Content = item.Content.replace("【投资顾问", "<br/>【投资顾问");
            }
            item.Content = item.Content.replace(/_big./g, "_small.");
        });
        Common.CommonBind(data, "liveContentList", true, function () {
            $("img", "#liveContentList").click(function () {
                var $img = $(this);
                var src = $img.attr("src");
                if (!src) return;
                src = src.replace("_small.", "_big.");
                showPic(src);
            });
        });
    }

    //获取全部问答内容
    function GetLiveAllQuestion(date) {
        var url = LiveAllQuestion;
        if (!!date) {
            url += "&date=" + date;
        }
        Common.CommonGet(url, function (response) {
            if (response.RetCode !== "0") {
                layer.msg("获取全部问答失败");
                console.log(response.RetMsg);
                return;
            }
            var data = response.Data;
            bindLiveAllQuestion(data, "");
        });
    }

    //绑定全部问答内容
    function bindLiveAllQuestion(data, date, room) {
        $.each(data, function () {
            var item = this;
            item.ShortTime = date + Common.GetShortTime(item.AskTime);
            if (!item.MasterNumber) {
                item.MasterNumber = room.AdviserNo;
                item.MasterName = room.AdviserName;
            }
        });
        Common.CommonBind(data, "allQuestionList", true);
    }

    //绑定直播内容+全部问答
    function BindLiveAllContent(date, isSelcte) {
        var flag = true;
        if (!date) {
            flag = false;
            date = Common.GetNowFormatDate();
        }
        LookDate = date;
        Common.CommonGet(LiveAllDataUrl + date, function (response) {
            if (response.RetCode !== "0") {
                layer.msg("获取直播内容失败");
                console.log(response.RetMsg);
                return;
            }

            var contentList = response.Data.contentlist;
            var questionList = response.Data.questionlist;
            var room = response.Data.room;

            if ((date !== Common.GetNowFormatDate() || !flag) && (!room || !room.TopicName) && isSelcte !== 1) {
                date = new Date(date);
                var lastTime = date.getTime() - (24 * 60 * 60 * 1000);
                var lastDay = new Date(lastTime);
                date = Common.GetNowFormatDate(lastDay);
                BindLiveAllContent(date)
            }
            else {
                if (date !== Common.GetNowFormatDate()) {
                    $("#date").html(date).attr("data-date", date);
                    date = date.substring(5) + " ";
                }
                else {
                    date = "";
                }
                bindRoomInfo(room);
                bindLiveContent(contentList, date);
                bindLiveAllQuestion(questionList, date, room);
            }
        });
    }

    //获取My问答内容
    function GetLiveMyQuestion(date) {
        var shortDate = "";
        var query = location.search;
        query = "&" + query.substr(1);
        if (!!date) {
            query += "&date=" + date;
            if (date != Common.GetNowFormatDate()) {
                shortDate = date.substring(5) + " ";
            }
        }

        Common.CommonGet(LiveMyQuestion + query, function (response) {
            if (response.RetCode !== "0") {
                layer.msg("获取我的问答失败");
                console.log(response.RetMsg);
                return;
            }
            var data = response.Data;
            $.each(data, function () {
                var item = this;
                item.ShortTime = shortDate + Common.GetShortTime(item.AskTime);
                ;
            });
            $.each(data, function () {
                var item = this;
                if (!item.AnswerUserName) {
                    item.style = "display: none;"
                }
                if (!item.MasterNumber) {
                    item.MasterNumber = roomInfo.AdviserNo;
                    item.MasterName = roomInfo.AdviserName;
                }

            });

            Common.CommonBind(data, "myQuestionList", true);
        })
    }

    function showPic(src) {
        $(".pic-pop img").attr("src", src);
        $(".pic-pop").show();
        $(".pic-pop").click(function () {
            $(this).hide();
        })
    }

    //注册事件
    function Register() {
        $("#subQuestion").click(function () {

            var question = $("#id_expertquestion").val();
            if (!question) {
                layer.msg("请填写问题~");
                return;
            }
            question = replace_em(question);
            question = encodeURI(question);

            var query = location.search;
            query = "&" + query.substr(1);
            var url = LiveSubmitQuestion + question + query;
            Common.CommonPost(url, {}, function (response) {
                if (response.RetCode !== "0") {
                    layer.msg("提问失败");
                    console.log(response.RetMsg);
                }
                $("#id_expertquestion").val("");

                $('#tabBox .question').removeClass('current');
                $('#tabBox .question').eq(1).addClass('current');
                $('#interactionContentBox .interaction-content-list').removeClass('open');
                $('#interactionContentBox .interaction-content-list').eq(1).addClass('open');

                GetLiveMyQuestion();
            })
        });
    }

    /*注册消息处理方法*/
    function OnMessage(event) {
        try {
            //延迟拉取最新数据
            setTimeout(function () {
                console.log("OnMessage");
                console.log(event.data);
                var data = JSON.parse(event.data)
                if (Common.GetNowFormatDate() == LookDate) {
                    if (data.ContentNum > 0) {
                        GetLiveContent();
                    }
                    if (data.RepliedNum > 0) {
                        GetLiveAllQuestion();
                    }
                }
                else {
                    var message = "有新的";
                    if (data.ContentNum > 0) {
                        message += "直播";
                    }
                    if (data.RepliedNum > 0) {
                        message += " 问答";
                    }
                    layer.msg(message + "消息", {offset: 'b'});
                }
            }, 3000);

        } catch (e) {
            console.log(e)
        }
    }

    /*注册ws关闭事件处理办法*/
    var OnClose = function (event) {
        try {
            console.log('OnClose');
            MBiz._FailTimes++;
            if (MBiz._FailTimes >= Times2Rollback) {
                MBiz.stop();
            }
        } catch (e) {
            console.log(e)
        }
    }

    /*注册ws错误事件处理办法*/
    var OnError = function (event) {
        try {
            console.log('OnError');
            MBiz._FailTimes++;
            if (MBiz._FailTimes >= Times2Rollback) {
                MBiz.stop();
            }
        } catch (e) {
            console.log(e)
        }
    }
    /*注册LongPoll消息处理方法*/
    var OnLongPollMessage = function (event) {
        try {
            console.log("OnLongPollMessage");

            console.log(event.data)
        } catch (e) {
            console.log(e)
        }
    }

    /*返回轮询模式*/
    var GoBackProcess = function () {
        console.log("GoBackProcess");
    }

    function replace_em(str) {
        str = str.replace(/\</g, '&lt;');
        str = str.replace(/\>/g, '&gt;');
        str = str.replace(/\n/g, '<br/>');
        str = str.replace(/\[em_([0-9]*)\]/g, '<img src="http://static.emoney.cn/live/Content/images/arclist/$1.gif" border="0" />');
        return str;
    }

    function init() {
        Register();
        Common.ShowMyStyle();
        BindLiveAllContent();
        if (window.WebSocket) {
            WSOn = 1;
        }
        if (WSOn == 1) {
            MBiz.init(AppId, GroupId, ClientId, OnMessage, OnClose, OnError);
        } else {
            MBiz.initlongpoll(AppId, GroupId, ClientId, OnLongPollMessage, OnClose, OnError, '', "2018-12-01", GoBackProcess);
        }
    }

    return {
        Init: init,
        BindLiveAllContent: BindLiveAllContent,
        GetLiveMyQuestion: GetLiveMyQuestion
    }
}();

$(function () {
    $Live.Init();
});