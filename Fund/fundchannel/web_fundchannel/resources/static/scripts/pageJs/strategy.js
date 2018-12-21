var $Strategy = function () {

    var setIntervalIndex = null;

    //获取 策略类别 并绑定
    function GetStrategyList() {
        Common.CommonAjax("home/GetStrategyList", {}, function (response) {
            if (response.RetCode !== 0) {
                layer.msg("获取基金策略信息失败");
                console.log(response.RetMsg);
                return;
            }
            var strategyData = JSON.parse(response.Message);
            $("#strategyAllInfoList").html(response.Message);
            $StrategyPage.strategyAllInfoList = response.Message;
            localStorageSet("strategyInfo", strategyData);

            var dataList = [];
            $.each(strategyData, function () {
                var strategyList = this.StrategyList;
                $.each(strategyList, function () {
                    var item = this;
                    dataList.push(
                        {
                            style: item.style
                            , name: item.name
                            , code: item.code
                        })
                });

            });
            Common.CommonBind(dataList, "strategyList", true, function () {
                pageInit();
                //策略 选中事件
                $("dd", "#strategyList").click(function () {
                    var oldCode = $(".actived", "#strategyList").data("code");
                    var code = $(this).data("code");
                    if (oldCode == code) return;
                    $("dd", "#strategyList").removeClass("actived");
                    $(this).addClass("actived");
                    localStorageSet("fundCode",code);
                    bindPageData(code);
                });
            });
        }, "GET")
    }

    //注册事件
    function Register() {
        //策略类别 选中事件
        $("li", "#poptype").click(function () {
            var oldStyle = $(".active", "#poptype").data("style");
            var style = $(this).data("style");
            if (oldStyle == style) return;

            $("dd", "#strategyList").removeClass("actived").hide();

            var selecter = "li[data-style='" + style + "']";
            var typeName = $(selecter, "#poptype").html();
            $("#selectedStrategyName").html(typeName);

            selecter = "dd[data-style='" + style + "']";
            $(selecter, "#strategyList").show();
            var code = $(selecter, "#strategyList").eq(0).data("code");

            selecter += "[data-code='" + code + "']";
            $(selecter, "#strategyList").addClass("actived");
            localStorageSet("fundStyle",style);
            localStorageSet("fundCode",code);

            bindPageData(code);
        });

        $('#selAdjustDates').change(function () {
            var date = $(this).val();
            $StrategyPage.loadStrategyInfos(date);
        });

        $StrategyPage.bindEvent();
    }

    //页面初始化
    function pageInit() {
        var style = Common.GetQueryString("fundStyle");
        var code = Common.GetQueryString("fundCode");
        if (!style) {
            style = localStorageGet("fundStyle");
            code = localStorageGet("fundCode");
            if (!style) {
                style = 1;
                code = null;
            }
        }

        var selecter = "li[data-style='" + style + "']";
        $(selecter, "#poptype").addClass("active");
        var typeName = $(selecter, "#poptype").html();
        $("#selectedStrategyName").html(typeName);

        selecter = "dd[data-style='" + style + "']";
        $(selecter, "#strategyList").show();
        if (!code) {
            code = $(selecter, "#strategyList").eq(0).data("code");
        }

        selecter += "[data-code='" + code + "']";
        $(selecter, "#strategyList").addClass("actived");

        bindPageData(code);
    }

    //绑定页面数据
    function bindPageData(code) {
        $StrategyPage.fundCode = code;
        var type = $('.switch-header-i.actived', "#groupTechContainer").attr("data-type");
        var date = $('.charttype .actived', "#groupTechContainer").attr('data-type');

        BindLatestInfo(code);
        GetStrategyConfigs(code)

        $StrategyPage.loadGroupTech(type, date);
        $StrategyPage.loadEstimation();

        if (!!setIntervalIndex) {
            clearInterval(setIntervalIndex);
        }
        setIntervalIndex = setInterval(function () {
            $StrategyPage.loadEstimation()
        }, 1000 * 60);
    }

    //绑定 组合涨幅
    function bindLatestInfo(data) {
        if (!data) return;
        var dataList = [];
        $.each(data.LatestInfo.percents, function () {
            var item = this;
            var class_percentW1 = Number(item.percentW1) >= 0 ? "fred" : "fgreen";
            var class_percentM1 = Number(item.percentM1) >= 0 ? "fred" : "fgreen";
            var class_percentM3 = Number(item.percentM3) >= 0 ? "fred" : "fgreen";
            var class_percentM6 = Number(item.percentM6) >= 0 ? "fred" : "fgreen";
            var class_percentYtd = Number(item.percentYtd) >= 0 ? "fred" : "fgreen";
            var class_percentY1 = Number(item.percentY1) >= 0 ? "fred" : "fgreen";
            var class_percentY3 = Number(item.percentY3) >= 0 ? "fred" : "fgreen";
            dataList.push({
                name: item.name
                , class_percentW1: class_percentW1
                , class_percentM1: class_percentM1
                , class_percentM3: class_percentM3
                , class_percentM6: class_percentM6
                , class_percentYtd: class_percentYtd
                , class_percentY1: class_percentY1
                , class_percentY3: class_percentY3
                , percentW1: Common.ConvertPercent(item.percentW1)
                , percentM1: Common.ConvertPercent(item.percentM1)
                , percentM3: Common.ConvertPercent(item.percentM3)
                , percentM6: Common.ConvertPercent(item.percentM6)
                , percentYtd: Common.ConvertPercent(item.percentYtd)
                , percentY1: Common.ConvertPercent(item.percentY1)
                , percentY3: Common.ConvertPercent(item.percentY3)
            });
        });
        Common.CommonBind(dataList, "groupRiseTbody", true);
    }

    function BindLatestInfo(code) {
        if (!code) return null;
        var data = localStorageGet("strategyInfo", code)
        if (!!data) {
            bindLatestInfo(data);
        }
        else {
            Common.CommonAjax("home/GetStrategyList", {}, function (response) {
                if (response.RetCode !== 0) {
                    layer.msg("获取基金策略信息失败");
                    return;
                }
                var strategyData = JSON.parse(response.Message);

                localStorageSet("strategyInfo", strategyData);
                var value = null;
                $.each(strategyData, function () {
                    if (!!value) return false;
                    var strategyList = this.StrategyList;
                    $.each(strategyList, function () {
                        var item = this;
                        if (item.code == code) {
                            value = item;
                            return false;
                        }
                    });
                });
                bindLatestInfo(value);
            }, "GET")
        }

    }

    function GetStrategyConfigs(code) {
        Common.CommonAjax("home/GetStrategyInfoByCode", {code: code}, function (response) {
            if (response.RetCode !== 0) {
                layer.msg("获取基金 策略配置 信息失败");
                return;
            }
            var strategyInfoList = JSON.parse(response.Message);
            $StrategyPage.strategyConfigInfoList = response.Message;
            var dateList = [];
            $.each(strategyInfoList, function () {
                var item = this;
                dateList.push(item.date);
            });
            $StrategyPage.loadAdjustDates(dateList);
        }, "GET");
    }

    function localStorageGet(storageName, key) {
        try {
            var data = null;
            switch (storageName) {
                case "strategyInfo": {
                    data = localStorage.strategyInfo;
                    if (!data) return null;
                    var value = null;
                    data = JSON.parse(data);
                    $.each(data, function () {
                        var strategyList = this.StrategyList;
                        if (!!value) return false;
                        $.each(strategyList, function () {
                            var item = this;
                            if (item.code == key) {
                                value = item;
                                return false;
                            }
                        });
                    });
                    return value;
                }
                case "fundStyle": {
                    try {
                        return localStorage.fundStyle;
                    } catch (e) {
                        return $.cookie("fundStyle");
                    }
                }
                case "fundCode": {
                    try {
                        return localStorage.fundCode;
                    } catch (e) {
                        return $.cookie("fundCode");
                    }
                }
                default:
                    console.log("不存在" + storageName + "仓库！");

            }
        }
        catch (e) {
            console.log("localStorageGet 异常");
            console.log(e);
            return null;
        }
    }

    function localStorageSet(storageName, value) {
        try {
            if (typeof (value) == "object") {
                value = JSON.stringify(value);
            }
            switch (storageName) {
                case "strategyInfo": {
                    localStorage.strategyInfo = value;
                    return;
                }
                case "fundStyle": {
                    try {
                        localStorage.fundStyle = value;
                    }
                    catch (e) {
                        $.cookie(storageName, value);
                    }
                }
                case "fundCode": {
                    try {
                        localStorage.fundCode = value;
                    }
                    catch (e) {
                        $.cookie(storageName, value);
                    }
                }
                default:
                    console.log("不存在" + storageName + "仓库！");
            }
        } catch (e) {
            console.log("localStorageSet 异常");
            console.log(e);
        }
    }

    function init() {
        GetStrategyList();
        Register();
        Common.ShowMyStyle();
        window.OutputData = function (data) {
            $StrategyPage.initGroupTechChart(data);
            $StrategyPage.loadGroupTechInfos(data);
        };
    }

    return {
        Init: init
    }
}();

$(function () {
    $Strategy.Init();
});