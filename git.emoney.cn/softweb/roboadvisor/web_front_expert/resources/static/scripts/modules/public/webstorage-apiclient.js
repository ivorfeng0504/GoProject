//缓存等级枚举
CacheLevel = {
    //页面级缓存，页面刷新时，缓存失效
    Page: 1,
    //会话级别缓存，浏览器关闭或者会话变更时，缓存失效
    Session: 2,
    //本地存储缓存，在本地持久保存一段时间，有效期内缓存一直有效
    Storage: 3
}

//API请求类型
APIMethod = {
    POST: "POST",
    JSONP: "JSONP"
}

APIClient = {
    
    //是否启用客户端缓存
    EnablesAllowClientCache: false,
    //强制关闭缓存
    ForceAllowClientCacheClosed: false,
    //设置缓存等级
    CacheLevel: CacheLevel.Page,
    //默认callApi的调用方式，全局控制callApi实际执行的方法，可设置为JSONP和POST
    DefaultCallMethod: APIMethod.POST,
    //请求基地址
    BaseUrl: undefined,
    //请求公共参数
    ApiRequestParams: ["AllowClientCache", "Hash", "JSONPCallback"],
    //用户请求数据前缀
    RequestData: "RequestData",


    //使用defaultCallMethod选定的方式请求服务器
    callApi: function (apiUrl, reqData, callback, method) {
        if (typeof (method) == "undefined") {
            method = this.DefaultCallMethod;
        }
        if (method == APIMethod.JSONP) {
            this.getApi(apiUrl, reqData, callback);
        } else {
            this.postApi(apiUrl, reqData, callback);
        }
    },

    //使用JSONP方式请求服务器，不提供通用的GET请求，非JSONP请求，均强制使用POST方式请求
    //apiUrl    api请求地址
    //reqData   api请求数据
    //callback  api执行成功后的回调函数
    getApi: function (apiUrl, reqData, callback) {
        apiUrl = this.buildReqUrl(apiUrl);
        if (apiUrl.indexOf("?") > 0) {
            apiUrl += "&JSONPCallback=?";
        } else {
            apiUrl += "?JSONPCallback=?";
        }
        var result = this.initAjax(apiUrl, reqData, callback);
        var requestHash = result.RequestHash;
        reqData = result.ReqData;
        $.ajax({
            url: apiUrl,
            type: "GET",
            dataType: "JSONP",
            data: reqData,
            success: APIClient.apiCallback(callback, requestHash)
        });
    },

    //使用POST方式请求服务器
    //apiUrl    api请求地址
    //reqData   api请求数据
    //callback  api执行成功后的回调函数
    postApi: function (apiUrl, reqData, callback) {
        apiUrl = this.buildReqUrl(apiUrl);
        var result = this.initAjax(apiUrl, reqData, callback);
        var requestHash = result.RequestHash;
        reqData = result.ReqData;
        $.ajax({
            url: apiUrl,
            type: "POST",
            dataType: "JSON",
            data: reqData,
            success: APIClient.apiCallback(callback, requestHash)
        });
    },

    //api调用回调处理函数
    //callback  用户的回调函数
    //requestHash   请求的哈希值
    apiCallback: function (callback, requestHash) {
        var func = function (data, textStatus, jqXHR) {
            if (data == undefined || data == null) {
                alert("请求异常，请稍后重试!");
                return;
            }
            if (data.ClientCached == true) {
                data = APIClient.getData(data.Hash)
                console.log("从本地缓存读取数据,读取到的数据为：" + getJSON(data));
            } else {
                APIClient.setData(data.Hash, data);
                APIClient.setData(requestHash, data.Hash);
                console.log("添加了缓存:" + data.Hash);
            }
            console.log("当前缓存等级为：" + APIClient.CacheLevel);
            callback(data, textStatus, jqXHR);
        }
        return func;
    },
    //发起api调用前的一些初始化工作,并返回请求的RequestHash值和处理后的请求参数ReqData
    //apiUrl    api请求地址
    //reqData   api请求数据
    //callback  api执行成功后的回调函数
    initAjax: function (apiUrl, reqData, callback) {
        if (reqData == undefined || reqData == null) {
            reqData = {};
        }
        else {
            //辅助用户处理用户请求数据RequestData
            //用户传参数时，可以使用RequestData作为前缀，也可以直接省略
            var tmpReq = reqData;
            reqData = {};
            //假设用户传递了RequestData
            reqData[this.RequestData] = tmpReq[this.RequestData];
            //如果用户没有传递RequestData
            if (reqData[this.RequestData] == undefined || reqData[this.RequestData] == null) {
                reqData[this.RequestData] = {};
            }
            for (var item in tmpReq) {
                //组装公共请求参数
                if (this.ApiRequestParams.indexOf(item) >= 0) {
                    reqData[item] = tmpReq[item];
                }
                    //组装用户请求数据
                else if (item != this.RequestData) {
                    reqData[this.RequestData][item] = tmpReq[item];
                }
            }
        }
        //手动清理用户设置的Hash
        reqData.Hash = "";
        var json = getJSON(reqData);
        var requestHash = md5(json + apiUrl);
        var dataHash = APIClient.getData(requestHash, true);
        if (dataHash != undefined && dataHash != null && trim(dataHash) != '') {
            //发起正式请求前，查询一下本地数据，防止缓存了错误的数据
            //如果有错误数据，则清理相关的缓存数据
            var cacheValue = APIClient.getData(dataHash);
            if (cacheValue != undefined && cacheValue != null) {
                reqData.Hash = dataHash;
            } else {
                APIClient.clearData(requestHash);
                APIClient.clearData(dataHash);
            }
        }
        //如果强制关闭缓存，则每次覆盖用户或全局的设置
        if (APIClient.ForceAllowClientCacheClosed == true) {
            reqData.AllowClientCache = false;
        }
        else if (reqData.AllowClientCache == undefined) {
            reqData.AllowClientCache = APIClient.EnablesAllowClientCache;
        }
        //返回结果
        var returnValue = {
            //用户的请求Hash值
            RequestHash: requestHash,
            //处理后的请求数据
            ReqData: reqData
        };
        return returnValue;
    },

    //构建请求地址
    //apiUrl    api地址
    buildReqUrl: function (apiUrl) {
        if (this.BaseUrl != undefined && trim(this.BaseUrl).length > 0 && apiUrl.toLowerCase().search("http://") < 0 && apiUrl.toLowerCase().search("https://") < 0) {
            apiUrl = this.BaseUrl + apiUrl;
        }
        return apiUrl;
    },

    //获取缓存数据
    //k 当前查询的键
    //ingoreCheck   是否忽略缓存降级检测
    getData: function (k, ingoreCheck) {
        var value = this.storage().getData(k);
        //如果从缓存中查询数据为空，则直接返回空，并且尝试清除该key
        if (value == null || value == undefined) {
            this.storage().clearData(k);
            return null;
        }
        //是否忽略检查
        if (ingoreCheck) {
            return value;
        }
        //如果从缓存中没有获取到数据
        //1.如果当前缓存不是页面级别（CacheLevel.Page），则降级将缓存降级到页面级别
        //2.如果当前缓存级别已经是页面级别（CacheLevel.Page），则取消全局缓存设置
        //3.如果当前缓存级别已经是页面级别（CacheLevel.Page），并且全局缓存开关已经关闭，则强制关闭缓存
        if (value == null || value == undefined) {
            if (APIClient.CacheLevel != CacheLevel.Page) {
                APIClient.CacheLevel = CacheLevel.Page
            } else if (APIClient.EnablesAllowClientCache == true) {
                APIClient.EnablesAllowClientCache = false;
            } else {
                APIClient.ForceAllowClientCacheClosed = true;
            }
        }
        return value;
    },

    //设置缓存数据
    //k 要设置的键
    //v 要设置的值
    setData: function (k, v) {
        this.storage().setData(k, v);
    },

    //清除指定的键
    //k 要清除的键
    clearData: function (k) {
        this.storage().clearData(k);
    },

    //获取当前存储器
    storage: function () {
        switch (APIClient.CacheLevel) {
            case CacheLevel.Storage:
                if (LocalStorage.isSupport()) {
                    APIClient.CacheLevel = CacheLevel.Storage;
                    return LocalStorage;
                }
            case CacheLevel.Session:
                if (SessionStorage.isSupport()) {
                    APIClient.CacheLevel = CacheLevel.Session;
                    return SessionStorage;
                }
            case CacheLevel.Page:
                if (PageStorage.isSupport()) {
                    APIClient.CacheLevel = CacheLevel.Page;
                    return PageStorage;
                }
            default: APIClient.CacheLevel = CacheLevel.Page; return PageStorage;
        }
    }
}


$.APIClient = APIClient;

//页面级存储
PageStorage = {
    //页面级别缓存数据
    cacheData: [],

    //获取缓存数据
    getData: function (k) {
        var value = this.cacheData[k];
        return value;
    },

    //设置缓存数据
    setData: function (k, v) {
        this.cacheData[k] = v;
    },

    //清除指定的键
    clearData: function (k) {
        this.cacheData[k] = null;
    },
    //当前浏览器是否支持本存储器
    isSupport: function () {
        return true;
    }
}

//会话存储
SessionStorage = {

    //获取缓存数据
    getData: function (k) {
        var json = sessionStorage.getItem(k);
        var obj = getObj(json);
        return obj;
    },

    //设置缓存数据
    setData: function (k, v) {
        var json = getJSON(v);
        if (json == null) {
            console.warn("存储数据失败，对应的键为：" + k);
            return;
        }
        var len = getJSON(sessionStorage).length;
        try {
            sessionStorage.setItem(k, json);
        } catch (e) {
            sessionStorage.clear();
            if (json.length < len) {
                sessionStorage.setItem(k, json);
            }
            console.log("sessionStorage数据库超过最大容量，数据库被重置");
        }
    },

    //清除指定的键
    clearData: function (k) {
        sessionStorage.removeItem(k);
    },
    //当前浏览器是否支持本存储器
    isSupport: function () {
        if (typeof (sessionStorage) == "undefined") {
            return false;
        }
        return true;
    }
}

//本地存储
LocalStorage = {

    //获取缓存数据
    getData: function (k) {
        var json = localStorage.getItem(k);
        var obj = getObj(json);
        return obj;
    },

    //设置缓存数据
    setData: function (k, v) {
        var json = getJSON(v);
        if (json == null) {
            console.warn("存储数据失败，对应的键为：" + k);
            return;
        }
        var len = getJSON(localStorage).length;
        try {
            localStorage.setItem(k, json);
        } catch (e) {
            localStorage.clear();
            if (json.length < len) {
                localStorage.setItem(k, json);
            }
            console.log("localStorage数据库超过最大容量，数据库被重置");
        }
    },

    //清除指定的键
    clearData: function (k) {
        localStorage.removeItem(k);
    },
    //当前浏览器是否支持本存储器
    isSupport: function () {
        if (typeof (localStorage) == "undefined") {
            return false;
        }
        return true;
    }
}

//去除字符串两端空格
function trim(str) {
    if (str == null || str == undefined) {
        return str;
    }
    if (typeof (str.trim) == "undefined") {
        return str.replace(/^\s+|\s+$/gm, '');
    } else {
        return str.trim();
    }
}

//序列化对象为JSON字符串
function getJSON(obj) {
    if (typeof (JSON) == "undefined") {
        console.error("当前浏览器不支持JSON操作");
        return;
    }
    try {
        var json = JSON.stringify(obj);
        return json;
    } catch (e) {
        console.error("getJSON 对象执行序列化失败");
        return null;
    }
}

//JSON字符串反序列化为对象
function getObj(json) {
    if (typeof (JSON) == "undefined") {
        console.error("当前浏览器不支持JSON操作");
        return;
    }
    try {
        var obj = JSON.parse(json);
        return obj;
    } catch (e) {
        console.warn("反序列化出错，原始字符串为：" + json);
        return null;
    }
}