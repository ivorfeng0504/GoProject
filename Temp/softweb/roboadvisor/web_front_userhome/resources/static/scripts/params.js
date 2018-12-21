//请求参数解析器
var ParamFunc = {
    //参数对象
    ParamObj: {},
    //所有的键
    Keys:[],
    initParams: function () {
        var searchStr = window.location.search;
        if (searchStr.length == 0) {
            return;
        }
        //去除问号
        var hasQuestion = true;
        while (hasQuestion) {
            if (searchStr.indexOf("?") == 0) {
                searchStr = searchStr.substr(1)
            } else {
                hasQuestion = false;
            }
        }
        if (searchStr.length == 0) {
            return;
        }
        var kvParams = searchStr.split("&")
        if (kvParams.length == 0) {
            return;
        }

        for (var i = 0; i < kvParams.length; i++) {
            var kv = kvParams[i].split("=")
            if (kv.length == 2) {
                this.ParamObj[kv[0]] = kv[1];
                this.Keys.push(kv[0])
            }
        }
    },

    //获取所有参数对象
    getParams: function () {
        return this.ParamObj;
    },

    //根据key获取对应的参数值
    getParam: function (key) {
        var value = this.ParamObj[key];
        return value;
    },

    //获取所有的键
    getKeys:function () {
        return this.Keys;
    }
}

ParamFunc.initParams();