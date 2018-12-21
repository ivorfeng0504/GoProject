var $Index = function () {
    var onGetYqqContent = true;

    /*绑定背景图片列表*/
    function bindBgImgList(data) {
        var dataList = [{
            index: 1,
            bgImg: "http://test.static.emoney.cn:8081/localtest/ifund/static/images/ifundmiddlebg1.png"
        }];
        $.each(data, function (i) {
            dataList.push({index: (i + 2), bgImg: this.bgImg});
        });

        Common.CommonBind(dataList, "styleBgBox", true);
    }

    /*绑定策略信息 card*/
    function bindStrategyCard(data) {
        var dataList = [];
        $.each(data, function () {
            var cardList = this.StrategyList;
            $.each(cardList, function () {
                var itme = this;
                var percent = Number(itme.LatestInfo.percents[0].percentM3);
                var _class = percent > 0 ? "fred" : "fgreen";
                var s = Common.GetMyLevel(itme.style)
                var label = s.tag + "风险";
                dataList.push(
                    {
                        style: itme.style
                        , percentM: percent
                        , label: label
                        , title: "近3个月收益率"
                        , name: itme.name
                        , summary: itme.summary
                        , code: itme.code
                        , _class: _class
                    })
            });

        });

        Common.CommonBind(dataList, "latestInfoList", true, function () {
            var selecter = "[data-style='1']";
            $(selecter, "#latestInfoList").show();
            $("li a", "#latestInfoList").click(function () {
                var fundStyle = $(this).data("style");
                var fundCode = $(this).data("code");
                var url = location.href;
                var indexOf = url.indexOf("&fundStyle");
                if (indexOf > 0) {
                    url = url.substr(0, indexOf);
                }
                url += "&fundStyle=" + fundStyle + "&fundCode=" + fundCode;
                location.href = url.replace("/index?", "/strategy?");
            });
        });
    }

    /* 绑定基金配置信息*/
    function GetStrategyInfo() {
        Common.CommonAjax("home/GetStrategyList", {}, function (response) {
            if (response.RetCode !== 0) {
                layer.msg("获取基金策略信息失败");
                return;
            }
            var strategyData = JSON.parse(response.Message);
            // 绑定背景图片
            bindBgImgList(strategyData);
            bindStrategyCard(strategyData);
        }, "GET")
    }

    //绑定益圈圈直播内容
    function GetYqqContent(date) {
        if (!date) {
            date = Common.GetNowFormatDate();
        }
        Common.CommonGet(LiveContentUrl + date, function (response) {
            if (response.RetCode !== "0") {
                layer.msg("获取直播室信息失败");
                return;
            }
            if (response.Data.length === 0) {
                date = new Date(date);
                var lastTime = date.getTime() - (24 * 60 * 60 * 1000);
                var lastDay = new Date(lastTime);
                date = Common.GetNowFormatDate(lastDay);
                GetYqqContent(date)
            }
            else {
                var liveList = [];
                var questionList = [];
                $.each(response.Data, function () {
                    var item = this;
                    item.Content = Common.ReplaceImg(item.Content);
                    if (item.LiveFrom == 1) {
                        var ShortTime = "--";
                        if (date === Common.GetNowFormatDate()) {
                            ShortTime = Common.GetShortTime(item.LastModifyTime);
                        }
                        else {
                            date = date.substring(5);
                            ShortTime = date + " " + Common.GetShortTime(item.LastModifyTime);
                        }
                        liveList.push({
                            ShortTime: ShortTime,
                            Content: item.Content
                        })
                    }
                    else if (item.LiveFrom == 2) {
                        var qAndA = getQandA(item.Content);
                        if (!!qAndA) {
                            questionList.push(qAndA);
                        }
                    }
                });

                Common.CommonBind(liveList, "liveContentList", true);
                Common.CommonBind(questionList, "questionContentList", true);
            }
        })
    }


    function getQandA(content) {
        var ask, answer, answerInfo;
        if (content.indexOf("【答】") > 0 && content.indexOf("【投资顾问") > 0) {
            var arr = content.split("【答】");
            ask = arr[0].substring(3);

            arr = arr[1].split("【投资顾问");
            answer = arr[0];
            answerInfo = "【投资顾问" + arr[1];
            return {
                ask: ask
                , answer: answer
                , answerInfo: answerInfo
            };
        }
        return null;
    }

    //获取直播间信息
    function GetRoomInfo() {
        Common.CommonGet(LiveRoomUrl, function (response) {
            if (response.RetCode !== "0") {
                layer.msg("获取直播室信息失败");
                console.log(response.RetMsg);
                return;
            }
            var data = response.Data;
            $("#liveName_Title").html(data.TopicName);
        });
    }

    //注册事件
    function Register() {
        $("#styleList li").click(function () {
            var $this = $(this);
            var style = $this.data("style");
            if (!style) return;
            $("li", "#latestInfoList").hide();
            $(".style-bg", "#styleBgBox").hide();
            var selecter = "[data-style='" + style + "']";
            $(selecter, "#latestInfoList").show();
            selecter = ".style-bg" + style;
            $(selecter, "#styleBgBox").show();
        });
    }

    function init() {
        GetStrategyInfo();
        GetYqqContent();
        GetRoomInfo();
        Register();
        Common.ShowMyStyle(function (levelInfo) {
            if (!!levelInfo) {
                $("#styleList li[data-style='" + levelInfo.no + "']").trigger("click");
            }
        });
    }

    return {
        Init: init
    }
}();

$(function () {
    $Index.Init();
});