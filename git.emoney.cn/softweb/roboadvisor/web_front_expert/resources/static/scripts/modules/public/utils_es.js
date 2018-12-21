/*
 * 通用方法
 *
 */
(function($) {

    $.support.cors = true;
    if (!Array.prototype.indexOf) {
        Array.prototype.indexOf = function(elt) {
            var len = this.length >>> 0;
            var from = Number(arguments[1]) || 0;
            from = from < 0 ? Math.ceil(from) : Math.floor(from);
            if (from < 0) from += len;
            for (; from < len; from++) {
                if (from in this && this[from] === elt) return from;
            }
            return -1;
        };
    }

    //移除数组元素
    Array.prototype.remove = function(val) {
        var index = this.indexOf(val);
        if (index > -1) {
            this.splice(index, 1);
        }
    };

    var _toQueryPair = function(key, value) {
        if (typeof value === "undefined") {
            return key;
        }
        return (
            key + "=" + encodeURIComponent(value === null ? "" : String(value))
        );
    };

    $.callEvent = function(name, func, proxy, relativeElement, params) {
        var event = $.Event(name);
        var result;
        if (relativeElement) {
            $(relativeElement).trigger(event, params);
            result = event.result;
        }
        if ($.isFunction(func)) {
            result = func.apply(proxy, $.isArray(params) ? params : [params]);
        }
        return result;
    };

    $.fn.callEvent = function(component, name, params) {
        return $.callEvent(
            name,
            component.options ? component.options[name] : null,
            component,
            this,
            params
        );
    };

    $.bindFn = function(fnName, _Constructor, defaultOptions) {
        var old = $.fn[fnName];
        var NAME = _Constructor.NAME || "feui." + fnName;

        $.fn[fnName] = function(option, params) {
            return this.each(function() {
                var $this = $(this);
                var data = $this.data(NAME);

                var options = typeof option == "object" && option;
                if (defaultOptions) options = $.extend(options, defaultOptions);

                if (!data)
                    $this.data(NAME, (data = new _Constructor($this, options)));
                if (typeof option == "string")
                    data[option].apply(
                        data,
                        $.isArray(params) ? params : [params]
                    );
            });
        };

        $.fn[fnName].Constructor = _Constructor;

        $.fn[fnName].noConflict = function() {
            $.fn[fnName] = old;
            return this;
        };
    };

    $.formatDate = function(date, format) {

        if (!(date instanceof Date)){ 
            date = date.replace("T"," ");
            date = new Date(date)
        };

        var dateInfo = {
            "M+": date.getMonth() + 1,
            "d+": date.getDate(),
            "h+": date.getHours(),
            // 'H+': date.getHours() % 12,
            "m+": date.getMinutes(),
            "s+": date.getSeconds(),
            // 'q+': Math.floor((date.getMonth() + 3) / 3),
            "S+": date.getMilliseconds()
        };
        if (/(y+)/i.test(format)) {
            format = format.replace(
                RegExp.$1,
                (date.getFullYear() + "").substr(4 - RegExp.$1.length)
            );
        }
        for (var k in dateInfo) {
            if (new RegExp("(" + k + ")").test(format)) {
                format = format.replace(
                    RegExp.$1,
                    RegExp.$1.length == 1 ?
                    dateInfo[k] :
                    ("00" + dateInfo[k]).substr(("" + dateInfo[k]).length)
                );
            }
        }
        return format;
    };

    $.format = function(str, args) {
        if (str instanceof Date) return $.formatDate(str, args);

        if (arguments.length > 1) {
            var reg;
            if (arguments.length == 2 && typeof args == "object") {
                for (var key in args) {
                    if (args[key] !== undefined) {
                        reg = new RegExp("({" + key + "})", "g");
                        str = str.replace(reg, args[key]);
                    }
                }
            } else {
                for (var i = 1; i < arguments.length; i++) {
                    if (arguments[i] !== undefined) {
                        reg = new RegExp("({[" + (i - 1) + "]})", "g");
                        str = str.replace(reg, arguments[i]);
                    }
                }
            }
        }
        return str;
    };

    $.calValue = function(anything, proxy, params) {
        return $.isFunction(anything) ?
            anything.apply(proxy, $.isArray(params) ? params : [params]) :
            anything;
    };

    $.is$ = function(obj) {
        return window.jQuery === $ ? obj instanceof jQuery : $.zepto.isZ(obj);
    };

    $.isStr = function(x) {
        return typeof x == "string";
    };

    $.isNum = function(x) {
        return typeof x == "number";
    };

    $.setFontSize = function(mode, obj) {
            var changeSize = 2,
                contentfontSize = 16;
            var elmContent = $("[data-fontsizecnt]")[0];
            if (mode < 0) {
                changeSize = -2;
            }
            contentfontSize = parseInt((elmContent.style.fontSize == "" ? (contentfontSize + "") : elmContent.style.fontSize).replace("px", "")) + changeSize;

            contentfontSize = contentfontSize < 12 ? 12 : contentfontSize;
            contentfontSize = contentfontSize > 32 ? 32 : contentfontSize;

            elmContent.style.lineHeight = (contentfontSize * 1.7) + "px";
            elmContent.style.fontSize = contentfontSize + "px";

        },

        $.TapName =
        "ontouchstart" in document.documentElement ? "touchstart" : "click";

    if (!$.uuid) $.uuid = (Math.random() + "").slice(-10);

    // 36bit random string
    $.getUUID =function () {
        var s = [];
        var hexDigits = "0123456789abcdef";
        for (var i = 0; i < 36; i++) {
            s[i] = hexDigits.substr(Math.floor(Math.random() * 0x10), 1);
        }
        s[14] = "4"; 
        s[19] = hexDigits.substr((s[19] & 0x3) | 0x8, 1); 
                                                            
        s[8] = s[13] = s[18] = s[23] = "-";
    
        var uuid = s.join("");
        return uuid;
    };


    //
    !(function($) {
        var pluses = /\+/g;

        function encode(s) {
            return config.raw ? s : encodeURIComponent(s);
        }

        function decode(s) {
            return config.raw ? s : decodeURIComponent(s);
        }

        function stringifyCookieValue(value) {
            return encode(config.json ? JSON.stringify(value) : String(value));
        }

        function parseCookieValue(s) {
            if (s.indexOf('"') === 0) {
                // This is a quoted cookie as according to RFC2068, unescape...
                s = s
                    .slice(1, -1)
                    .replace(/\\"/g, '"')
                    .replace(/\\\\/g, "\\");
            }

            try {
                // Replace server-side written pluses with spaces.
                // If we can't decode the cookie, ignore it, it's unusable.
                // If we can't parse the cookie, ignore it, it's unusable.
                s = decodeURIComponent(s.replace(pluses, " "));
                return config.json ? JSON.parse(s) : s;
            } catch (e) {}
        }

        function read(s, converter) {
            var value = config.raw ? s : parseCookieValue(s);
            return $.isFunction(converter) ? converter(value) : value;
        }

        var config = ($.cookie = function(key, value, options) {
            // Write

            if (value !== undefined && !$.isFunction(value)) {
                options = $.extend({}, config.defaults, options);

                if (typeof options.expires === "number") {
                    var days = options.expires,
                        t = (options.expires = new Date());
                    t.setTime(+t + days * 864e5);
                }

                return (document.cookie = [
                    encode(key),
                    "=",
                    stringifyCookieValue(value),
                    options.expires ? "; expires=" + options.expires.toUTCString() : "", // use expires attribute, max-age is not supported by IE
                    options.path ? "; path=" + options.path : "",
                    options.domain ? "; domain=" + options.domain : "",
                    options.secure ? "; secure" : ""
                ].join(""));
            }

            // Read

            var result = key ? undefined : {};

            // To prevent the for loop in the first place assign an empty array
            // in case there are no cookies at all. Also prevents odd result when
            // calling $.cookie().
            var cookies = document.cookie ? document.cookie.split("; ") : [];

            for (var i = 0, l = cookies.length; i < l; i++) {
                var parts = cookies[i].split("=");
                var name = decode(parts.shift());
                var cookie = parts.join("=");

                if (key && key === name) {
                    // If second argument (value) is a function it's a converter...
                    result = read(cookie, value);
                    break;
                }

                // Prevent storing a cookie that we couldn't decode.
                if (!key && (cookie = read(cookie)) !== undefined) {
                    result[name] = cookie;
                }
            }

            return result;
        });

        config.defaults = {};

        $.removeCookie = function(key, options) {
            if ($.cookie(key) === undefined) {
                return false;
            }

            // Must not alter options, thus extending a fresh object...
            $.cookie(key, "", $.extend({}, options, { expires: -1 }));
            return !$.cookie(key);
        };
    })($);


    return {
        /*
         * 处理文本输入框的placeholder
         */
        placeholder: function() {
            _placeholderHandle();
        },

        ieVersion: function(callback) {
            _ieVersion(callback);
        },

        /**
         * 字符串截取
         */
        subStr: function(str, length) {
            if (str.length > length) {
                return str.substr(0, parseInt(length)) + "...";
            }
            return str;
        },
        /*
         * 去掉前后空格
         */
        strTrim: function(s) {
            return s.replace(/(^\s+)|(\s+$)/g, "");
        },
        GetQueryString: function(name) {
            var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)");
            var r = window.location.search.substr(1).match(reg);
            if (r != null) return unescape(r[2]);
            return null;
        },

        /*
         * 解析RUI参数
         * str: uri字符串
         */
        parseURIParams: function(str) {
            var params = {},
                e,
                a = /\+/g,
                r = /([^&=]+)=?([^&]*)/g,
                d = function(s) {
                    return decodeURIComponent(s.replace(a, " "));
                };

            while ((e = r.exec(str))) {
                params[d(e[1])] = d(e[2]);
            }
            return params;
        },

        /*
         * 对像转成URI
         */
        objToQuery: function(obj) {
            var ret = [];
            for (var key in obj) {
                key = encodeURIComponent(key);
                var values = obj[key];
                if (values && values.constructor === Array) {
                    var queryValues = [];
                    for (var i = 0, len = values.length, value; i < len; i++) {
                        value = values[i];
                        queryValues.push(_toQueryPair(key, value));
                    }
                    ret = ret.concat(queryValues);
                } else {
                    ret.push(_toQueryPair(key, values));
                }
            }
            return ret.join("&");
        },
        /*
         * 取当前路径的参数值
         * arg: 参数名
         */
        parseLocation: function(arg) {
            var uri = location.search;
            if (uri !== "") {
                var argsObj = this.parseURIParams(uri.substr(1));
                return argsObj[arg] || "";
            }
            return "";
        },
        /*
         * 中文链接编码
         */
        b64EncodeUrl: function(string) {
            if (window.BASE64) {
                return BASE64.encoder(string.replace("风格", ""))
                    .replace("+", "-")
                    .replace("/", "_")
                    .replace("=", "");
            }
            return string;
        },
        /*
         * Timeago 相对时间美化  2011-05-06 12:30:22  ---> 三分钟之前
         */
        prettyDate: function(time) {
            var date = new Date(
                    (time || "").replace(/-/g, "/").replace(/[TZ]/g, " ")
                ),
                diff = (new Date().getTime() - date.getTime()) / 1000,
                day_diff = Math.floor(diff / 86400);

            if (isNaN(day_diff) || day_diff < 0) {
                return;
            } else if (day_diff >= 31) {
                return time;
            }

            return (
                (day_diff === 0 &&
                    ((diff < 60 && "刚刚") ||
                        (diff < 120 && "1分钟前") ||
                        (diff < 3600 && Math.floor(diff / 60) + "分钟前") ||
                        (diff < 7200 && "1个小时前") ||
                        (diff < 86400 &&
                            Math.floor(diff / 3600) + "小时前"))) ||
                (day_diff === 1 && "昨天") ||
                (day_diff < 7 && day_diff + "天前") ||
                (day_diff < 31 && Math.ceil(day_diff / 7) + "周前")
            );
        },
        /**
         * 切换城市刷新URL
         */
        changeURLArg: function(url, arg, arg_val) {
            if (url.indexOf("#")) {
                url = url.split("#")[0];
            }
            var pattern = arg + "=([^&]*)";
            var replaceText = arg + "=" + arg_val;
            if (url.match(pattern)) {
                var tmp = "/(" + arg + "=)([^&]*)/gi";
                tmp = url.replace(eval(tmp), replaceText);
                return tmp;
            } else {
                if (url.match("[?]")) {
                    return url + "&" + replaceText;
                } else {
                    return url + "?" + replaceText;
                }
            }
            return url + "\n" + arg + "\n" + arg_val;
        },
        /**
         * url跳转
         */
        locationUrl: function(url) {
            var w = window.open();
            return (w.location = url);
        },
        /**
         * 计算总页数
         * total 记录总数
         * size 每页显示的记录个数
         */
        pageCount: function(total, size) {
            var count = Math.floor(total / size),
                vod = total % size;
            if (vod > 0) {
                count += 1;
            }
            return count;
        },
        /**
         * 金额格式化
         * money 数额
         * split 是否每3位添加一个分隔，通常是','，不分不要传
         */
        formatCurrency: function(money, split) {
            split = split || "";
            var num = money.toString().replace(/\$|\,/g, ""),
                sign;
            if (isNaN(num)) {
                num = "0";
            }
            sign = num == (num = Math.abs(num));
            num = Math.floor(num * 100 + 0.50000000001);
            cents = num % 100;
            num = Math.floor(num / 100).toString();
            if (cents < 10) cents = "0" + cents;
            for (var i = 0; i < Math.floor((num.length - (1 + i)) / 3); i++)
                num =
                num.substring(0, num.length - (4 * i + 3)) +
                split +
                num.substring(num.length - (4 * i + 3));
            return (sign ? "" : "-") + num + "." + cents;
        },
        /**
         * 链接中的Next参数
         */
        uriNext: function(def) {
            uriObj = this.parseURIParams(location.search.substr(1));
            return uriObj.next || (def || "");
        },

        //优化url，去掉url中不合法的token
        optimizeUrl: function(url) {
            var re = new RegExp("<[^>]*>", "gi");
            url = url.replace(re, "");
            return url;
        },

        //判断是否邮件
        isEmail: function(str) {
            return this.reEmail.test(str);
        },

        /*
         * 检查发布内容是否包含链接
         */
        checkContentUrl: function(content) {
            var matchStr = "fuwo";
            var flag = false;
            var indexResult;
            var re_http = new RegExp(
                "(http[s]{0,1}|ftp)?(:)?(//)?[a-zA-Z0-9\\.\\-]+\\.([a-zA-Z]{2,4})(:\\d+)?(/[a-zA-Z0-9\\.\\-~!@#$%^&*+?:_/=<>]*)?",
                "gi"
            );
            var pic_re = new RegExp(".+.(png|PNG|jpg|JPG|bmp|gif|GIF)$");
            if (content.match(re_http) === null) {
                return true;
            } else {
                var result_http =
                    content.match(re_http) === null ?
                    "" :
                    content.match(re_http).toString();
                var resultArray_http = [];
                resultArray_http = result_http.split(",");
                //http验证
                if (resultArray_http !== "") {
                    for (var i = 0; i < resultArray_http.length; i++) {
                        resultArray_http[i] = this.optimizeUrl(
                            resultArray_http[i]
                        );
                        if (!pic_re.test(resultArray_http[i])) {
                            if (!this.isEmail(resultArray_http[i])) {
                                indexResult =
                                    resultArray_http[i].indexOf(matchStr) >= 0 ?
                                    true :
                                    false;
                                if (!indexResult) {
                                    flag = true;
                                    break;
                                }
                            }
                        }
                    }
                }

                if (flag) {
                    return false;
                }
                return true;
            }
        },

        //判断发布内容中是否有广告链接
        checkUrl: function(content) {
            if (this.checkContentUrl(content) === false) {
                alert("发布内容中包含非本站点链接，请检查您的发布内容！");
                return false;
            }
            return true;
        },

        //改变锚点标签颜色
        changeAnchorColor: function(content) {
            var re_color = new RegExp("<a", "gi");
            return content.replace(
                re_color,
                '<a style="color:rgb(120,120,200)"'
            );
        },

        //首字母大写
        ucFirst: function(word) {
            return word.substring(0, 1).toUpperCase() + word.substring(1);
        },

        //解析url
        parseUrl: function(url) {
            var regUrl = {
                protocol: /([^\/]+:)\/\/(.*)/i,
                host: /(^[^\:\/]+)((?:\/|:|$)?.*)/,
                port: /\:?([^\/]*)(\/?.*)/,
                pathname: /([^\?#]+)(\??[^#]*)(#?.*)/
            };
            var tmp,
                res = {};

            res["href"] = url;
            for (p in regUrl) {
                tmp = regUrl[p].exec(url);
                res[p] = tmp[1];
                url = tmp[2];
                if (url === "") {
                    url = "/";
                }
                if (p === "pathname") {
                    res["pathname"] = tmp[1];
                    res["search"] = tmp[2];
                    res["hash"] = tmp[3];
                }
            }
            return res;
        },

        //将毫秒时间转化为xx-xx-xx的格式
        getStyleTime: function(time) {
            var oDate = new Date(time * 1000),
                oYear = oDate.getFullYear(),
                oMonth = oDate.getMonth() + 1,
                oDay = oDate.getDate(),
                oTime = oYear + "-" + getzf(oMonth) + "-" + getzf(oDay); //最后拼接时间
            return oTime;

            function getzf(num) {
                if (parseInt(num) < 10) {
                    num = "0" + num;
                }
                return num;
            }
        },

        // 股票跳转
        goThisStock: function(aStkCode) {
            function GetExternal() {
                return window.external.EmObj
            };

            //调用客户端接口
            function PC_JH(type, c) {
                try {
                    var obj = GetExternal();
                    return obj.EmFunc(type, c)
                } catch (e) {}
            };
            
            try { PC_JH("EM_FUNC_GOTO_TECH_VIEW", "455," + aStkCode + ""); } catch (ex) {}
        },

        // EMSSO处理
        EMSSO: function () { 
            var _this = this;
            var ssoOBJ = {},ssoSTR="";
            if (!$.cookie("expertnews.ssourl")) {
                ssoOBJ.rand = _this.GetQueryString('rand');
                ssoOBJ.Version = _this.GetQueryString('Version');
                ssoOBJ.uid = _this.GetQueryString('uid');
                ssoOBJ.pid = _this.GetQueryString('pid');
                ssoOBJ.sid = _this.GetQueryString('sid');
                ssoOBJ.tid = _this.GetQueryString('tid');
                ssoOBJ.agentid = _this.GetQueryString('agentid');
                ssoOBJ.clienttype = _this.GetQueryString('clienttype');
                ssoOBJ.OutOfDate = _this.GetQueryString('OutOfDate');
                ssoOBJ.pluglet = _this.GetQueryString('pluglet');
                ssoOBJ.token = _this.GetQueryString('token');
                ssoOBJ.bata = _this.GetQueryString('bata');
                ssoSTR = _this.objToQuery(ssoOBJ);
                $.cookie("expertnews.ssourl", encodeURIComponent(ssoSTR));
            } else { 
                ssoSTR = $.cookie("expertnews.ssourl");
            }
            return ssoSTR;
        },
        PC_JH: function (type, c) {

            //params : string string 
            // type:  string  "EM_FUNC_ZNTG_CLZL",
            // c:  string "aimHREF"
            //获取客户端接口
            function GetExternal() {
                return window.external.EmObj
            };

            //调用客户端接口
            
            function clientJSBridge(type, c) {
                try {
                    var obj = GetExternal();
                    return obj.EmFunc(type, c)
                } catch (e) {}
            };

            clientJSBridge(type, c);
        }


    };
})(jQuery);
 