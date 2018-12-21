function initAdvisoryId() {
    var uid = ParamFunc.getParam("uid");
    var pid = ParamFunc.getParam("pid");
    var sid = ParamFunc.getParam("sid");
    var tid = ParamFunc.getParam("tid");
    var online = AddjrptAd_Click("span_online", sid, tid, pid, Ad_AC_Award, AppId, uid, ""); //在线咨询
}

function initAdvisoryIndex_up() {
    var uid = ParamFunc.getParam("uid");
    var pid = ParamFunc.getParam("pid");
    var sid = ParamFunc.getParam("sid");
    var tid = ParamFunc.getParam("tid");
    var online = AddjrptAd_Click("div_adver_up", sid, tid, pid, Ad_AC_Index_Up, AppId, uid, ""); //首页-右上角
}

function initAdvisoryIndex_down() {
    var uid = ParamFunc.getParam("uid");
    var pid = ParamFunc.getParam("pid");
    var sid = ParamFunc.getParam("sid");
    var tid = ParamFunc.getParam("tid");
    var online = AddjrptAd_Click("div_adver_dw", sid, tid, pid, Ad_AC_Index_Dw, AppId, uid, ""); //首页-右下角
}

function initAdvisoryIndex_downFree() {
    var uid = ParamFunc.getParam("uid");
    var pid = ParamFunc.getParam("pid");
    var sid = ParamFunc.getParam("sid");
    var tid = ParamFunc.getParam("tid");
    var online = AddjrptAd_Click("div_adver_dwnew", sid, tid, pid, Ad_AC_Index_DwFree, AppId, uid, ""); //首页-右下角（免费）
}


function initAdvisoryClass() {
    var handle=setInterval(function () {
        var adHtml=$("#span_online").html();
        if(adHtml!="<a></a>"){
            $(".span_online").html(adHtml);
            clearInterval(handle)
        }
    },100)
}

function initAdvisoryIdProduct() {
    var uid = ParamFunc.getParam("uid");
    var pid = ParamFunc.getParam("pid");
    var sid = ParamFunc.getParam("sid");
    var tid = ParamFunc.getParam("tid");
    var online = AddjrptAd_Click("span_online_product", sid, tid, pid, Ad_AC_Product, AppId, uid, ""); //在线咨询
}

function initAdvisoryClassProduct() {
    var handle=setInterval(function () {
        var adHtml=$("#span_online_product").html();
        if(adHtml!="<a></a>"){
            $(".span_online_product").html(adHtml);
            clearInterval(handle)
        }
    },100)
}

//初始化我的-中间广告位
function initMineCenterAd() {
    var uid = ParamFunc.getParam("uid");
    var pid = ParamFunc.getParam("pid");
    var sid = ParamFunc.getParam("sid");
    var tid = ParamFunc.getParam("tid");
    var online = AddjrptAd_Click("ad_mine_center", sid, tid, pid, Ad_AC_Mine_Center, AppId, uid, "");
}