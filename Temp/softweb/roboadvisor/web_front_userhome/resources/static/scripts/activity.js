var ActivityState_NotBegin=1;
var ActivityState_Beginning=2;
var ActivityState_Finish=3;
var HasMobile=false;
$(function () {
    initActivity(ActivityState_Beginning,true);

    //初始化账户状态
    initMobileState();

    //当期活动
    $("#btn_activity_current").click(function () {
        $("#showActFrame").hide()
        initActivity(ActivityState_Beginning);
    });

    //即将开始的活动
    $("#btn_activity_future").click(function () {
        $("#showActFrame").hide()
        initActivity(ActivityState_NotBegin);
    });

    //已结束的活动
    $("#btn_activity_finish").click(function () {
        $("#showActFrame").hide()
        initActivity(ActivityState_Finish);
    });
});

//初始化活动列表 autoShowAct 是否自动展示指定的活动
function initActivity(state,autoShowAct) {
    $.post(www+"activity/getactivitylist",{
        state:state
    },function(data){
        if(data.RetCode == "0"){
            switch(state){
                case ActivityState_NotBegin:
                    var $div_activity_future = $("#div_activity_future");
                    var html = template_activity_future(data);
                    $div_activity_future.html(html);
                    break;
                case ActivityState_Beginning:
                    var $div_activity_current = $("#div_activity_current");
                    var html = template_activity_current(data);
                    $div_activity_current.html(html);
                    initActivityClick();
                    initActivityPageClick();
                    if(autoShowAct){
                        autoShowActivity();
                    }
                    break;
                case ActivityState_Finish:
                    var $div_activity_finish = $("#div_activity_finish");
                    var html = template_activity_finish(data);
                    $div_activity_finish.html(html);
                    break;
            }
        }else{

        }
    });
}

function initActivityClick() {
    $(".activity_begining").click(function () {
        var activityEle=$(this);
        var activityId=activityEle.attr("activity_id");
        var needBind=activityEle.attr("need_bind");
        var needSSO=activityEle.attr("need_sso");
        //没有绑定手机号则跳转到注册页
        if(HasMobile==false&&needBind=="1"){
            //询问框
            layer.confirm('您还没有绑定手机，请绑定手机后再参与活动', {
                title:"",
                btn: ['去绑定手机','再看看'] //按钮
            }, function(){
                // 绑定手机
                window.location.href=$("#a_home_page").attr("href")+"&openreg=1"
            }, function(){
                // 放弃绑定
            });
        }else{
            var activityUrl=activityEle.find(".viewimg").attr("url");
            if(activityUrl!=undefined&&activityUrl.length>0){
                //判断是否要附加SSO
                if(needSSO=="1"){
                    activityUrl=appendSSO(activityUrl);
                }
                //判断打开方式 0 当前iframe 1 新窗口弹出
                var openMode=activityEle.find(".viewimg").attr("open_model");
                if(openMode=="1"){
                    window.open(activityUrl)
                }else{
                    $.showActities(activityUrl);
                }
            }
        }
    });
}

function initActivityPageClick() {
    var App = "10150";
    var Module = "zy_userhome";//模块名称
    var Remark = ParamFunc.getParam("uid")+"|"+ParamFunc.getParam("pid");     //备注uid|pid
    var ClickFlag = true;//默认为true
    $.EMoneyAnalytics.Init(App, Module, Remark, false, ClickFlag);
}

// 附加SSO串
function appendSSO(url) {
    if(SSOStr==undefined||SSOStr==null||SSOStr.length==0||url==undefined||url==null||url.length==0){
        return url;
    }
    var sso=SSOStr;
    if(sso[0]=="?"||sso[0]=="&"){
        sso=sso.substring(1)
    }
    if(url.indexOf("?")>=0){
        url=url+"&"+sso;
    }else{
        url=url+"?"+sso;
    }
    return url
}
//初始化账户状态
function initMobileState() {
    $.get(www+"mine/getbindmobilebyuid",{
    },function(data){
        if(data.RetCode == "0") {
            if (data.Message != undefined && data.Message.EncryptMobile != undefined && data.Message.EncryptMobile.length>0) {
                HasMobile = true;
            }
        }
    });
}

//自动跳转到指定的活动
function autoShowActivity() {
    var activityId=ParamFunc.getParam("showActivity")
    activityId=parseInt(activityId)
    if(isNaN(activityId)){
        return
    }
    var activityUrl=$(".activity_begining[activity_id='"+activityId+"'] .viewimg").attr("url");
    if(activityUrl==undefined||activityId==null||activityId.length==0){
        return
    }
    $.showActities(activityUrl);
}