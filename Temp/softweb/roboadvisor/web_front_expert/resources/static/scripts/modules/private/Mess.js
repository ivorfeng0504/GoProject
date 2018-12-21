 
var MBiz = {
    _client: {
        AppId: "",
        GroupId: "",
        UserId: "",
        ConnectTime: "",
        Status: ""
    },

    _FailTimes: 0,
    _run: 1,

    stop: function () {
        MBiz._run = 0;
    },

    _sock: [],

    _sockupdate: function (cid, c) {
        MBiz._sock.forEach(function (e, index) {
            if (e[0] == cid) {
                MBiz._sock.splice(index, 1, c);
            }
        })
    },

    _sockdelete: function (cid) {
        MBiz._sock.forEach(function (e, index) {
            if (e[0] == cid) {
                MBiz._sock.splice(index, 1);
            }
        })
    },


    init: function (appid, groupid, userid, OnMessage, OnClose, OnError, OnOpen, Token) {

        if (MBiz._run != 1) { return; }

        var wsc; 
        try {
            if (typeof (Token) != "undefined") {
                wsc = new WebSocket(wsurl + '?appid=' + appid + '&groupid=' + groupid + '&userid=' + userid + '&token=' + Token);
            } else {
                wsc = new WebSocket(wsurl + '?appid=' + appid + '&groupid=' + groupid + '&userid=' + userid);
            }
        } catch (e) {
            console.log("WebSocket Initialize Error");
            return;
        }
        
        wsc.onopen = OnOpen;
        wsc.onmessage = OnMessage;
        wsc.onerror = OnError;

        var cid = parseInt(Math.random() * 100000);

        MBiz._sock.push([cid, wsc]);

        /*如果没有自定义关闭事件，则自动重连*/
        if (OnClose == null) {
            wsc.onclose = function () {
                /*被动关闭后2s自动重连*/
                setTimeout(function () {
                    if (MBiz._run != 1) { return; }
                    MBiz.init(appid, groupid, userid, OnMessage, OnClose, OnError, OnOpen);
                    MBiz._sockdelete(cid);
                }, 2000);

            };
        } else {
            wsc.onclose = OnClose;
        }

        return cid;
    },

    //LongPoll
    initlongpoll: function (appid, groupid, userid, OnMessage, OnClose, OnError, OnOpen, querykey, GoBackProcess,Token) {
        
        var pollingurl = longpollurl + '?appid=' + appid + '&groupid=' + groupid + '&userid=' + userid + '&querykey=' + escape(querykey) + '&_r='+Math.random();

        if (typeof (Token) != "undefined") {
            pollingurl = longpollurl + '?appid=' + appid + '&groupid=' + groupid + '&userid=' + userid + '&querykey=' + escape(querykey) + '&token=' + Token + '&_r=' + Math.random();
        }

        $.ajax({
            type: "GET",
            cache: false,
            url: pollingurl,
            dataType: "jsonp",
            jsonp: "jsonpcallback",
            success: function (data, textStatus, jqXHR) {
                if (data != null && data != "" && data != "null") {

                    if (data.RetCode == "0") {
                        var ev = {};
                        eval('var obj=' + data.Message + ';');
                        if (typeof(obj)!="undefined") {
                            ev.data = obj;
                            OnMessage(ev);
                        }
                        
                    } else {
                        if (data.RetCode == "-100009" ) {
                            //100毫秒后重连
                            setTimeout(function () {
                                MBiz.initlongpoll(appid, groupid, userid, OnMessage, OnClose, OnError, OnOpen, querykey , GoBackProcess);
                            }, 100);
                        } else if (data.RetCode == "-200009") {
                            //应用消息接口异常等待一分钟后重连
                            setTimeout(function () {
                                MBiz.initlongpoll(appid, groupid, userid, OnMessage, OnClose, OnError, OnOpen, querykey, GoBackProcess);
                            }, 60*1000);
                        } else {
                            //返回轮询模式
                            GoBackProcess();
                        }
                    } 
                        
                }
                    
            },
            complete: function (XMLHttpRequest, textStatus) {
            },
            error: function (XMLHttpRequest, textStatus, errorThrown) {
                if (textStatus == "timeout") {
                    MBiz.initlongpoll(appid, groupid, userid, OnMessage, OnClose, OnError, OnOpen, querykey, GoBackProcess);
                }
                
            }

        });

    }
    
};


var browser = {
    versions: function () {
        var u = navigator.userAgent, app = navigator.appVersion;
        return {//移动终端浏览器版本信息
            trident: u.indexOf('Trident') > -1, //IE内核
            presto: u.indexOf('Presto') > -1, //opera内核
            webKit: u.indexOf('AppleWebKit') > -1, //苹果、谷歌内核
            gecko: u.indexOf('Gecko') > -1 && u.indexOf('KHTML') == -1, //火狐内核
            mobile: !!u.match(/AppleWebKit.*Mobile.*/) || !!u.match(/AppleWebKit/), //是否为移动终端
            ios: !!u.match(/\(i[^;]+;( U;)? CPU.+Mac OS X/), //ios终端
            android: u.indexOf('Android') > -1 || u.indexOf('Linux') > -1, //android终端或者uc浏览器
            iPhone: u.indexOf('iPhone') > -1 || u.indexOf('Mac') > -1, //是否为iPhone或者QQHD浏览器
            iPad: u.indexOf('iPad') > -1, //是否iPad
            webApp: u.indexOf('Safari') == -1 //是否web应该程序，没有头部与底部
        };
    }(),
    language: (navigator.browserLanguage || navigator.language).toLowerCase()
}