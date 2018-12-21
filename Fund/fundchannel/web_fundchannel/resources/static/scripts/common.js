var Common = function () {
    String.prototype.replaceFormat = function (obj) {
        return this.replace(/\{{\w+\}}/gi, function (matchs) {
            var returns = obj[matchs.replace(/\{{/g, "").replace(/\}}/g, "")];
            return (returns + "") == "undefined" ? "" : returns;
        });
    };
    var styleList = [
        {
            no: 1
            , level: "C1"
            , tag: "低"
            , name: "保守型"
            , remark: "根据您的风险测评结果-C1，您适配的投资风格是保守型，为您推荐保守型基金策略组合。该类型风险容忍程度非常低，在做投资决定时，尽量保护本金不受损失是首要目标。"
        }, {
            no: 2
            , level: "C2"
            , tag: "中低"
            , name: "稳健型"
            , remark: "根据您的风险测评结果-C2，您适配的投资风格是稳健型，为您推荐稳健型基金策略组合。该类型风险容忍程度中等偏低。在做投资决定时，希望在保障本金的基础上能有一些增值收入。"
        }, {
            no: 3
            , level: "C3"
            , tag: "中"
            , name: "平衡型"
            , remark: "根据您的风险测评结果-C3，您适配的投资风格是平衡型，为您推荐平衡型基金策略组合。该类型风险容忍程度比较中等。在做投资决定时，会考虑到在风险相对可控的的情况下获得一定的收益。"
        }, {
            no: 4
            , level: "C4"
            , tag: "中高"
            , name: "成长型"
            , remark: "根据您的风险测评结果-C4，您适配的投资风格是成长型，为您推荐成长型基金策略组合。该类型风险容忍程度中等偏高。在做投资决定时，希望有较高的投资收益，并因此愿意承受一定的投资波动。"
        }, {
            no: 5
            , level: "C5"
            , tag: "高"
            , name: "进取型"
            , remark: "根据您的风险测评结果-C5，您匹配的投资风格是进取型，为您推荐进取型基金策略组合。该类型风险容忍程度极高，在做投资决定时，对风险的考量较少，而更关注于投资收益，并愿意为此承受较大的风险。"
        }
    ];
    /**
     * Ajax公用方法
     * @param {string} url：请求地址
     * @param {object} data：请求数据
     * @param {function} successFunc:成功回调
     * @param {string} type:默认POST
     * @param {boolean} async
     */
    var commonAjax = function (url, data, successFunc, type, async) {
        async = async == undefined ? true : async;
        type = !!type ? type : "POST";
        $.ajax({
            url: www + url,
            type: type,
            async: async,
            dataType: "json",
            data: data,
            beforeSend: function () {
                layer.load(2);
            },
            success: function (response) {
                try {
                    if (typeof(successFunc) === "function") {
                        successFunc(response);
                    }
                } catch (e) {
                    layer.msg("操作异常,请稍后重试");
                    console.log(e);
                }
            },
            error: function (XMLHttpRequest, errorInfo) {
                layer.msg("请求失败,请重新登录or联系管理员");
                console.log(XMLHttpRequest);
                console.log(errorInfo);
            },
            complete: function () {
                layer.closeAll("loading");
            }
        });

    }
    var commonGet = function (url, successFunc, async) {
        async = async == undefined ? true : async;
        $.ajax({
            url: url,
            type: "GET",
            async: async,
            dataType: "json",
            beforeSend: function () {
                layer.load(2);
            },
            success: function (response) {
                try {
                    if (typeof(successFunc) === "function") {
                        successFunc(response);
                    }
                } catch (e) {
                    layer.msg("操作异常,请稍后重试");
                    console.log(e);
                }
            },
            error: function (XMLHttpRequest, errorInfo) {
                layer.msg("请求失败,请重新登录or联系管理员");
                console.log(XMLHttpRequest);
                console.log(errorInfo);
            },
            complete: function () {
                layer.closeAll("loading");
            }
        });
    }
    var commonPost = function (url, data, successFunc) {
        $.ajax({
            url: url,
            type: "Post",
            data: data,
            dataType: "json",
            beforeSend: function () {
                layer.load(2);
            },
            success: function (response) {
                try {
                    if (typeof(successFunc) === "function") {
                        successFunc(response);
                    }
                } catch (e) {
                    layer.msg("操作异常,请稍后重试");
                    console.log(e);
                }
            },
            error: function (XMLHttpRequest, errorInfo) {
                layer.msg("请求失败,请重新登录or联系管理员");
                console.log(XMLHttpRequest);
                console.log(errorInfo);
            },
            complete: function () {
                layer.closeAll("loading");
            }
        });
    }

    /**
     * Handlebars绑定公用方法
     * @param {object} data：绑定数据
     * @param {string} selecter：容器
     * @param {boolean} isArry:数据是否集合
     * @param {function} func:回调
     */
    var commonBind = function (data, selecter, isArry, func) {
        try {

            var tpl = $("#" + selecter + "_temp").html();
            var template = Handlebars.compile(tpl);
            if (isArry) {
                data = {dataList: data};
            }
            var html = template(data);
            $("#" + selecter).html(html);

            if (!!func && typeof (func) === "function") {
                func();
            }
        }
        catch (e) {
            console.log(e)
        }
    }

    /**
     * 获取格式化时间
     * @date {Date} date：日期，默认当前日期
     * @separator {string} separator：分隔符，默认“-”
     */
    var getNowFormatDate = function (date, separator) {
        if (!date) {
            date = new Date();
        }
        if (!separator) {
            separator = "-";
        }
        var year = date.getFullYear();
        var month = date.getMonth() + 1;
        var strDate = date.getDate();
        if (month >= 1 && month <= 9) {
            month = "0" + month;
        }
        if (strDate >= 0 && strDate <= 9) {
            strDate = "0" + strDate;
        }
        var currentData = year + separator + month + separator + strDate;
        return currentData;
    }

    /**
     * 获取Url的参数值
     * @name {string} name：参数名称
     */
    function getQueryString(name) {
        var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
        var r = window.location.search.substr(1).match(reg);
        if (r != null) return unescape(r[2]);
        return null;
    }

    /**
     * 获取newGuid
     */
    function newGuid() {
        var guid = "";
        for (var i = 1; i <= 32; i++) {
            var n = Math.floor(Math.random() * 16.0).toString(16);
            guid += n;
            if ((i == 8) || (i == 12) || (i == 16) || (i == 20))
                guid += "-";
        }
        return guid;
    }

    /**
     * 去除HTML标签
     * @html {string} html：html
     */
    function replaceImg(html) {
        html = html.replace(/\<img[^\>]*\>/g, "【图片需进入查看】");
        html = html.replace(/<[^>]+>/g, "");
        html = html.replace(/&nbsp;/g, "").trim();
        return html;
    }

    /**
     * 字符串 百分化
     * @obj {string} obj：需要转化的数字
     * @isEnlarge {boolean} isEnlarge：是否乘以 100
     */
    function convertPercent(obj, isEnlarge) {
        try {
            if (!obj) return "--";

            var x = Number(obj)
            if (isNaN(x)) return "--";

            if (isEnlarge) {
                x = x * 100;
                x = x.toFixed(2);
            }
            return x + "%";
        } catch (e) {
            console.log("数字转换成字符串 异常");
            console.log(e);
            return "--";
        }
    }

    /**
     * 展示我的风格和头像
     */
    function showMyStyle(func) {
        var pid = getQueryString("pid");
        $(".my-style,.style-desc-wrapper").hide();
        if (pid == "888010000" || pid == "888020000") {

            var userId = getQueryString("uid");
            if (!userId) return;
            commonAjax("home/GetEncryptMobile", {userId: userId}, function (response) {
                if (response.RetCode !== 0) {
                    layer.msg("获取用户加密手机号失败");
                    console.log(response.RetMsg)
                    return;
                }
                var data = JSON.parse(response.Message);
                if (data.RetCode != "0") {
                    layer.msg("获取用户加密手机号失败");
                    console.log(response.RetMsg)
                    return;
                }
                data = JSON.parse(data.Message);
                if (data.RetCode !== 0) {
                    layer.msg("获取用户加密手机号失败");
                    console.log(response.RetMsg)
                    return;
                }

                var mobileNumber = null;
                $.each(data.Message, function () {
                    var item = this;
                    if (item.AccountType === 1 && !!item.EncryptMobile) {
                        mobileNumber = item.EncryptMobile;
                        return false;
                    }
                })
                if (!mobileNumber) return;

                commonAjax("home/QueryUserRiskInfo", {mobileNumber: mobileNumber}, function (response) {
                    if (response.RetCode !== 0) {
                        layer.msg("获取用户风格失败");
                        console.log(response.RetMsg)
                        return;
                    }
                    var data = JSON.parse(response.Message);
                    if (data.RetCode != 0) {
                        layer.msg("获取用户风格失败");
                        console.log(response.RetMsg)
                        return;
                    }
                    data = JSON.parse(data.Message);
                    if (!data.IsSuccess) {
                        layer.msg("获取用户风格失败");
                        console.log(response.RetMsg)
                        return;
                    }
                    if (!!data.RiskName) {
                        var levelInfo = getMyLevel(data.RiskName);
                        if (!!levelInfo) {
                            $("#myStyle_name").html("【" + levelInfo.name + "】");
                            $("#myStyle_name_").html(levelInfo.name + "：");
                            $("#myStyle_remark").html(levelInfo.remark);

                            $(".my-style,.style-desc-wrapper").show();
                            if (!!func && typeof (func) === "function") {
                                func(levelInfo);
                            }
                        }
                    }
                }, "GET");

            }, "GET");
        }
        //获取用户头像
        commonAjax("home/GetUserInfoByUID", {}, function (response) {
            if (response.RetCode !== 0) {
                layer.msg("获取用户信息失败");
                console.log(response.RetMsg)
                return;
            }
            if (!!response.Message) {
                var data = JSON.parse(response.Message);
                if (!!data.Headportrait) {
                    $("img.header-pic").attr("src", data.Headportrait);
                }
            }
        }, "GET");
    }

    function formatTimeStr(timeStr, mstep) {
        var currDate = new Date();
        var timeSrc = currDate.getFullYear() + '/' + (currDate.getMonth() + 1) + '/' + currDate.getDay() + ' ' + timeStr + ':00';
        var timeDriv = new Date(timeSrc);
        var timeNext = timeDriv.getTime() + mstep * 60 * 1000;
        var resTime = new Date(timeNext);
        var H = '00' + resTime.getHours(), M = '00' + resTime.getMinutes();
        var resStr = H.substr(H.length - 2, 2) + ':' + M.substr(M.length - 2, 2);
        return resStr;
    }

    function getTimePoint(respoonArr, startTime, endTIme, minutesStep) {
        var result = formatTimeStr(startTime, minutesStep);
        respoonArr.push(result);
        if (result == endTIme) {
            return respoonArr;
        } else {
            return getTimePoint(respoonArr, result, endTIme, minutesStep);
        }
    }

    function createTimeArray() {
        var TimeArray = new Array();
        TimeArray.push('09:30');
        TimeArray.concat(getTimePoint(TimeArray, '09:30', '11:30', 1));
        //TimeArray.push('13:00');
        TimeArray.concat(getTimePoint(TimeArray, '13:00', '15:00', 1));
        return TimeArray;
    }

    /**
     * 获取我的风险等级
     * @level {string} level：等级
     */
    function getMyLevel(level) {
        var levelInfo = null;
        $.each(styleList, function () {
            var style = this;
            if (typeof (level) === "number") {
                if (style.no === level) {
                    levelInfo = style;
                    return false;
                }
            }
            else if (level.indexOf("型") > 0) {
                if (style.name === level) {
                    levelInfo = style;
                    return false;
                }
            }
            else {
                if (style.level === level) {
                    levelInfo = style;
                    return false;
                }
            }
        })
        return levelInfo;
    }

    function getShortTime(time, flag) {
        //"2018-12-09T15:10:44"
        if (!flag) {
            flag = "T";
        }
        if (time.indexOf(flag) > 0) {
            time = time.split(flag);
            if (time.length !== 2) return "00:00";
            return time[1].substring(0, 5);
        }
        return "00:00";
    }

    return {
        CommonAjax: commonAjax
        , CommonBind: commonBind
        , CommonGet: commonGet
        , CommonPost: commonPost
        , GetNowFormatDate: getNowFormatDate
        , GetQueryString: getQueryString
        , NewGuid: newGuid
        , ReplaceImg: replaceImg
        , ConvertPercent: convertPercent
        , ShowMyStyle: showMyStyle
        , CreateTimeArray: createTimeArray
        , GetMyLevel: getMyLevel
        , GetShortTime: getShortTime
    }
}();