function MyMBiz(appId, groupId, groupIds, userId, index, connectTime, status) {
    this.AppId = appId;
    this.GroupId = groupId;
    this.GroupIds = groupIds;
    this.UserId = userId;
    this.ConnectTime = connectTime;
    this.Status = status;
    this.FailTimes = 0;
    this.Times2Rollback = 2;
    //记录索引编号
    this.Index = index;
    this._run = 1;
    this._sock = [];
    this.stop = function () {
        this._run = 0;
    };

    this._sockupdate = function (cid, c) {
        var $this = this;
        $this._sock.forEach(function (e, index) {
            if (e[0] == cid) {
                $this._sock.splice(index, 1, c);
            }
        })
    };

    this._sockdelete = function (cid) {
        var $this = this;
        $this._sock.forEach(function (e, index) {
            if (e[0] == cid) {
                $this._sock.splice(index, 1);
            }
        })
    };

    this.init = function (OnMessage, OnClose, OnError, OnOpen, Token) {
        var $this = this;
        var appid = $this.AppId;
        var groupid = $this.GroupId;
        var userid = $this.UserId;
        var groupIds = $this.GroupIds;

        if ($this._run != 1) {
            return;
        }

        var wsc;
        try {
            if (typeof (Token) != "undefined") {
                wsc = new WebSocket(wsurl + '?appid=' + appid + '&groupid=' + groupid + '&groupids=' + groupIds + '&userid=' + userid + '&token=' + Token);
            } else {
                wsc = new WebSocket(wsurl + '?appid=' + appid + '&groupid=' + groupid + '&groupids=' + groupIds + '&userid=' + userid);
            }
        } catch (e) {
            console.log("WebSocket Initialize Error");
        }

        //使用bizIndex在WebSocket实例中记录当前MyMBiz实例的索引
        wsc.bizIndex = $this.Index;
        wsc.onopen = OnOpen;
        wsc.onmessage = OnMessage;
        wsc.onerror = OnError;

        var cid = parseInt(Math.random() * 100000);

        $this._sock.push([cid, wsc]);

        /*如果没有自定义关闭事件，则自动重连*/
        if (OnClose == null) {
            wsc.onclose = function () {
                /*被动关闭后2s自动重连*/
                setTimeout(function () {
                    if ($this._run != 1) {
                        return;
                    }
                    $this.init(OnMessage, OnClose, OnError, OnOpen);
                    $this._sockdelete(cid);
                }, 2000);

            };
        } else {
            wsc.onclose = OnClose;
        }

        return cid;
    };

    //LongPoll
    this.initlongpoll = function (OnMessage, OnClose, OnError, OnOpen, querykey, GoBackProcess, Token) {
        var $this = this;
        var appid = $this.AppId;
        var groupid = $this.GroupId;
        var userid = $this.UserId;
        var groupIds = $this.GroupIds;

        var pollingurl = longpollurl + '?appid=' + appid + '&groupid=' + groupid + '&groupids=' + groupIds + '&userid=' + userid + '&querykey=' + escape(querykey) + '&_r=' + Math.random();

        if (typeof (Token) != "undefined") {
            pollingurl = longpollurl + '?appid=' + appid + '&groupid=' +
                groupid + '&groupids=' +
                groupIds + '&userid=' +
                userid + '&querykey=' + escape(querykey) +
                '&token=' + Token + '&_r=' + Math.random();
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
                        if (typeof (obj) != "undefined") {
                            ev.data = obj;
                            ev.$MBiz = $this;
                            OnMessage(ev);
                        }

                    } else {
                        if (data.RetCode == "-100009") {
                            //100毫秒后重连
                            setTimeout(function () {
                                $this.initlongpoll(OnMessage, OnClose, OnError, OnOpen, querykey, GoBackProcess);
                            }, 100);
                        } else if (data.RetCode == "-200009") {
                            //应用消息接口异常等待一分钟后重连
                            setTimeout(function () {
                                $this.initlongpoll(OnMessage, OnClose, OnError, OnOpen, querykey, GoBackProcess);
                            }, 60 * 1000);
                        } else {
                            //返回轮询模式
                            GoBackProcess();
                        }
                    }

                }

            },
            complete: function (XMLHttpRequest, textStatus) {},
            error: function (XMLHttpRequest, textStatus, errorThrown) {
                if (textStatus == "timeout") {
                    $this.initlongpoll(OnMessage, OnClose, OnError, OnOpen, querykey, GoBackProcess);
                }

            }

        });


    }

}

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

var browser = {
    versions: function () {
        var u = navigator.userAgent,
            app = navigator.appVersion;
        return { //移动终端浏览器版本信息
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