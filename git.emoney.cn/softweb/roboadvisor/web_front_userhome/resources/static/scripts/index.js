/*
 * @Author: gitshilly
 * @Date:   2018-03-14 16:30:27
 * @Version:   1.0
 * @Last Modified by:   qiaoli
 * @Last Modified time: 2018-05-16 09:21:42
 */
var StrategyID = 0;//策略id 0不区分策略
var HasMobile=false;

//用户信息信息模板
var source_user = $("#userInfo-template").html();
var template_user = Handlebars.compile(source_user);

//我的资料模板
var source_myProfile= $("#myProfile-template").html();
var template_myProfile = Handlebars.compile(source_myProfile);

//置顶资讯模板
var source_top = $("#xxzhNewsTop-template").html();
var template_top = Handlebars.compile(source_top);

//非置顶资讯模板
var source_xxzh = $("#xxzhNewsList-template").html();
var template_xxzh = Handlebars.compile(source_xxzh);

//资讯详情
var source_newsinfo = $("#newsinfo-template").html();
var template_newsinfo = Handlebars.compile(source_newsinfo);

var pid = ParamFunc.getParam("pid");

function TimerNewsListFlagCount() {
    pageData.init();
    getTopNewsList();
    setTimeout(function () {
        TimerNewsListFlagCount();
    }, 600 * 1000)
}

$.ajaxSetup({ cache: false });
//打开消息阅读窗口
openWindow();

//置顶通知
function getTopNewsList() {
    $("#top_NewsList").html("");
    $.get(www + "homepage/topnewslist", {ColumnID: userHome_columnID, PID: pid, r: Math.random()}, function (data) {
        if (data.RetCode == "0") {
            var newsList = data.Message;
            var retObj = {NewsList: newsList};
            var topnewsHTML = template_top(retObj);
            $("#top_NewsList").html(topnewsHTML);
        }
    });
}

//通知信息详情
function getNewsInfoByID(NewsID) {
    $.ajax({
        url:www + "homepage/newsinfo" ,
        type: 'GET',
        dataType: 'json',
        data: {NewsId:NewsID,r:Math.random()},
        success: function(data){
            if(data.RetCode=="0") {
                var retObj = data.Message;
                var html = template_newsinfo(retObj);
                $("#div_article").html(html);
            }
        },
        beforeSend: function () {

        },
        error: function(XMLHttpRequest, textStatus, errorThrown){

        }
    });
}

//分页获取通知列表
var pageData = {
    pagecount: 0,
    pageSize: 6,
    isCompleted: false,
    init: function() {
        $("#PageIndex").val("1");
        pageData.isCompleted = false;
        pageData.load(1);
        $("#li_NewsList").html("");
    },
    load: function(pagenum) {
        $.ajax({
            url: www + "homepage/newslist",
            type: 'GET',
            dataType: 'json',
            data: {ColumnID:userHome_columnID, PID: pid,currpage:pagenum,pageSize:pageData.pageSize,r:Math.random()},
            success: function (data) {
                HiddenLoadingDocument();
                if (data.RetCode == "0") {
                    var newsList = data.Message;
                    var retObj = {NewsList: newsList};
                    var html = template_xxzh(retObj);
                    $("#li_NewsList").append(html);
                    $("#li_NewsList").height(150);

                    pageData.pagecount = data.TotalCount;
                    pageData.setPageIndex();
                    if (Math.ceil(pageData.pagecount/pageData.pageSize) == pagenum) {
                        pageData.isCompleted = true;
                    }else{
                        pageData.isCompleted = false;
                    }
                }
                if (data.RetCode == "-2") {
                    NoneDataDocument($("#li_NewsList"));
                }
            },
            beforeSend: function () {
                LoadingDocument();
                // 禁用按钮防止重复提交
                pageData.isCompleted = true;
            },
            error: function (XMLHttpRequest, textStatus, errorThrown) {

            }
        });
    },
    scroll: function(page) { //滚动到底部加载数据
        if (pageData.isCompleted) {
            return false;
        }
        var top = $('#li_NewsList').scrollTop();
        var win = $('#li_NewsList').height();
        var doc = document.getElementById("li_NewsList").scrollHeight;

        if ((top + win) >= doc) {
            pageData.load(page);
        }
    },
    setPageIndex: function() { //数据载入成功，设置下一页索引
        var $page = $("#PageIndex");
        var index = parseInt($page.val()) + 1;
        $page.val(index);
    },
    appendHtml: function(data) {
    }
};

//获取用户信息
function getUserInfo()
{
    $.get(www + "homepage/userinfo", {}, function (data) {
        if (data.RetCode == "0") {
            var userInfo = data.Message;
            //首页左上角用户信息
            var userInfoHTML = template_user(userInfo);
            $("#div_userInfo").html(userInfoHTML);

            //我的资料-弹窗
            var myProfileHTML = template_myProfile(userInfo);
            $("#div_myProfile").html(myProfileHTML);

            //手机掩码
            $("[flag=div_mobilemask]").html(userInfo.MobileMask);
            $("#hid_mobilex").val(userInfo.MobileX);

            //初始化用户等级图标
            initUserLevelICO(userInfo.UserLevel);

            //初始化连登奖励
            initLoginAward();

            //初始化勋章
            initMyMedal();

            //判断是否有手机号
            initMobileState();
        }
    });
}

//初始化账户状态
function initMobileState() {
    $.get(www+"mine/getbindmobilebyuid",{
    },function(data){
        if(data.RetCode == "0") {
            if (data.Message != undefined && data.Message.EncryptMobile != undefined && data.Message.EncryptMobile.length>0) {
                HasMobile = true;
            }else{
                $(".usr-mobile").removeClass("active").css("cursor","pointer");
            }
        }
    });
}

//打开交互弹窗
function openWindow()
{

    //通知消息窗口
    $('#top_NewsList,#li_NewsList').on("click","[flag=openWindow]",function() {
        var id = $(this).data("id");
        getNewsInfoByID(id);
        openPop("newsArticleshow","消息阅读",600,420);
    });
    //我的资料窗口
    $('#div_userInfo').on("click","[flag=btn_MyProfile]",function() {
        location.href = www+"mine/myprofile";
    }).on("click","[flag=btn_ModifyPwd]",function() {
        //修改密码窗口
        var usertype = $(this).data("usertype");
        if(usertype=="0"){
            //游客、微信、QQ-弹出注册页
            //LoadValidCode();
            //openPop("userRegist","用户注册",430,395);
        }
        else{
            $.get(www+"mine/getbindmobilebyuid",{},function (data) {
                if(data.RetCode=="0" && data.Message!="" && data.Message.EncryptMobile!="") {
                    $("#hid_mobilex").val(data.Message.EncryptMobile);
                    $("[flag=div_mobilemask]").html(data.Message.AccountName);
                    openPop("usrModyPwd", "修改密码", 430, 350);
                }else{
                    LoadValidCode();
                    openPop("bindPhone","手机绑定",430,310);
                }
            });
        }
    }).on("click","[flag=btn_BindAccount]",function() {
        //绑定账号窗口
        var usertype = $(this).data("usertype");
        if(usertype=="0" ){
            //游客-弹出注册页
            //LoadValidCode();
            //openPop("userRegist","用户注册",430,395);
        }
        else{
            $.get(www+"mine/getbindmobilebyuid",{},function (data) {
                if(data.RetCode=="0" && data.Message!="" && data.Message.EncryptMobile!="") {
                    //手机-掩码tips
                    var mobilemask = data.Message.AccountName;
                    layer.tips(mobilemask, '.usr-mobile',{
                        tips:[2, '#0FA6D8']
                    });
                }else{
                    LoadValidCode();
                    openPop("bindPhone","手机绑定",430,310);
                }
            });
        }
    }).on("click","#usrAvatarShow",function () {
        //修改头像
        ModifyHeadportrait();
    }).on("click","#modyUsrname",function () {
        var $editBtn = $("a", this),
            $editbox = $("#editUsrname"),
            $editInp = $("input", $editbox),
            $editSpan = $("b", $editbox);

        if ($editBtn.text() === "[编辑]") {
            $("#editUsrname").addClass("editing");
            $editInp.val($editSpan.text());
            $editBtn.text("[保存]");
        } else {
            if($editInp.val()==""){
                alert("请输入昵称");
                return false;
            }
            if($editInp.val().length>20){
                alert("昵称格式不正确，建议在20个字符内");
                return false;
            }

            $("#editUsrname").removeClass("editing");
            $editBtn.text("[编辑]");
            //$editSpan.text($editInp.val()).attr("title", $editInp.val());

            SaveNickName($editInp.val());
        }
    });
}

var captchaid = "";//验证码使用
$(document).ready(function() {
    EM_FUNC_DOWNLOAD_COMPLETE();

    if (pid=="888010000" || pid=="888020000")
    {
        $("#div_adver_dw").show();
    }else{
        $("#div_free").show();
    }
    getUserInfo();
    TimerNewsListFlagCount();

    var $page = $("#PageIndex"); //页索引

    // 滚动加载-学习综合
    $('#li_NewsList').scroll(function () {
        pageData.scroll($page.val());
    });
});

function LoadValidCode(){
    $.getJSON(www + "captcha/fetchid", {}, function (data) {
        captchaid = data.CaptchaId;

        setTimeout(function () {
            var imgcodeUrl = www + "captcha/image?captchid=" + captchaid;
            $("[name=img_code]").attr("src", imgcodeUrl);
        },200);

    });
}
$("[name=img_code]").click(function () {
    RefreshValidCode();
});
function RefreshValidCode() {
    var imgcodeUrl = www + "captcha/reloadimage?captchid=" + captchaid;
    $("[name=img_code]").attr("src", imgcodeUrl);
}
/*********************************数据提交begin*******************************/
var tnum = 0;
var flagtime;
var flag = true;
var flagtext = "获取验证码";
$('.btn-getcheckcode').click(function () {
    var posturl="";
    /* Act on the event */
    //$(this).addClass('disabled');
    if (flag) {
        var type = $(this).data("type");
        var mobile = "";
        if (type == "encrypt") {
            //密文手机号，发送短信
            mobile=$("#hid_mobilex").val();
        }
        else {
            //明文手机号，发送短信
            var mobileid = $(this).data("mobileid");
            mobile = $("#" + mobileid).val();
            var telReg = /^[1][3,4,5,7,8][0-9]{9}$/;
            if (mobile.length != 11 || (telReg.test(mobile)) == false) {
                layer.msg("请输入正确的手机号！");
                $("#" + mobileid).focus();
                return;
            }
        }
        flag = false;

        posturl = www + "mine/getValidateCode";
        if(type == "getpwd") {
            flagtext = "获取密码";
            posturl = www + "mine/getpwdbymobile";
        }
        $.getJSON(posturl, { Mobile: do_encrypt(mobile) }, function (data) {
            if (data.RetCode == "0") {
                timejs();
            } else {
                layer.msg(data.RetMsg);
                flag = true;
            }
        });
    }
});

function timejs() {
    var num = 60 - tnum;
    tnum = tnum + 1;
    $('.btn-getcheckcode').html('重发(' + num + 'S)');
    flagtime = setTimeout("timejs()", 1000);
    if (tnum > 60) {
        stopjs();
        $('.btn-getcheckcode').html("获取验证码");
        flag = true;
        tnum = 0;
    }
}
function stopjs() {
    clearTimeout(flagtime);
}
function resetTimejs(){
    stopjs();
    $('.btn-getcheckcode').html("获取验证码");
    flag = true;
    tnum = 0;
}
//提交-注册
var regflag = true;
$("#submit_reg").click(function () {
    if(regflag) {
        /* Act on the event */
        errMsg = "";
        if (!Validate(1)) {
            layer.msg(errMsg);
            //layer.msg(errMsg);
            return false;
        }
        // Next
        var mobile = $("#usrMobile").val();
        var code = $("#usrCheckcode").val();
        var usrPassword = $("#usrPassword").val();
        var usrVeryfypwd = $("#usrVeryfypwd").val();

        if (usrPassword != usrVeryfypwd) {
            layer.msg("两次密码输入不一致，请确认后输入");
            $("#usrVeryfypwd").focus();
            return false;
        }
        regflag = false;
        $.getJSON(www + "mine/regTyAccount", {
            Mobile: do_encrypt(mobile),
            code: code,
            password: do_encrypt(usrPassword),
            confirmpassword: do_encrypt(usrVeryfypwd)
        }, function (data) {
            regflag = true;
            if (data.RetCode == "0") {
                layer.msg("注册成功");
                resetTimejs();
                //关闭弹窗
                closePop();
                return false;
            } else {
                layer.msg(data.RetMsg);
                return false;
            }
        });
    }
});
//提交-绑定（绑定手机号）
var phonebind_flag = true;
$("#submit_phonebind").click(function () {
    if(phonebind_flag) {
        /* Act on the event */
        errMsg = "";
        if (!Validate(1)) {
            layer.msg(errMsg);
            return false;
        }
        // Next
        var mobile = $("#usrbindMobile").val();
        var pwd = $("#usrbindPwd").val();
        var code = $("#bindphone_code").val();

        phonebind_flag = false;
        $.getJSON(www + "mine/bindaccountphone", {
            Mobile: do_encrypt(mobile),
            pwd: do_encrypt(pwd),
            code: code,
            captchaId: captchaid
        }, function (data) {
            phonebind_flag = true;
            if (data.RetCode == "0") {
                layer.msg("手机号绑定成功");
                setTimeout(function () {
                    window.location.reload();
                }, 1500);
                return false;
            } else {
                layer.msg(data.RetMsg);
                RefreshValidCode();
                return false;
            }
        });
    }
});
//提交-找密
var modifyflag = true;
$("#submit_modifypwd").click(function () {
    if(modifyflag) {

        /* Act on the event */
        errMsg = "";
        if (!Validate(0)) {
            layer.msg(errMsg);
            return false;
        }
        // Next
        var mobile = $("#hid_mobilex").val();//已绑定密文手机
        var code = $("#usrmodyphoneCheckcode").val();
        var usrPassword = $("#usrnewPassword").val();
        var usrVeryfypwd = $("#usrVeryfyndwpwd").val();

        if(isNaN(usrPassword))
        {
            layer.msg("密码只能是6位数字");
            $("#usrnewPassword").focus();
            return false;
        }

        if(isNaN(usrVeryfypwd))
        {
            layer.msg("密码只能是6位数字");
            $("#usrVeryfyndwpwd").focus();
            return false;
        }

        if(usrPassword.length>6||usrPassword.length<3){
            layer.msg("密码只能是6位数字");
            $("#usrnewPassword").focus();
            return false;
        }

        if(usrVeryfypwd.length>6||usrVeryfypwd.length<3){
            layer.msg("密码只能是6位数字");
            $("#usrVeryfyndwpwd").focus();
            return false;
        }

        if (usrPassword != usrVeryfypwd) {
            layer.msg("两次密码输入不一致，请确认后输入");
            $("#usrVeryfyndwpwd").focus();
            return false;
        }
        modifyflag = false;
        $.getJSON(www + "mine/modifypassword", {
            Mobile: mobile,
            code: code,
            password: do_encrypt(usrPassword),
            confirmpassword: do_encrypt(usrVeryfypwd)
        }, function (data) {
            modifyflag = true;
            if (data.RetCode == "0") {
                layer.msg("密码修改成功");
                setTimeout(function () {
                    closePop();
                }, 3000);
                resetTimejs();
                return false;
            } else {
                layer.msg(data.RetMsg);
                return false;
            }
        });
    }
});

//修改头像
function ModifyHeadportrait() {
    var usrAvartar = $("#usrAvatarShow img");
    usrAvartar.data('oldavatarurl', usrAvartar.attr("src"));

    var popcnt = layer.open({
        type: 1
        , title: ['用户头像选择', '']
        , skin: 'ucpop-layer'
        , area: ['600px', '440px']
        , scrollbar: false
        , btn: ['保存选择']
        , yes: function (index) {
            var headimg_num = $("#hid_headimg").val();
            $.getJSON(www + "mine/modifyheadportrait", {headimg: headimg_num}, function (data) {
                if (data.RetCode == "0") {
                    $("#usrAvatarShow img").attr('src', www + '/static/images/Arena_' + headimg_num + '.png');
                    //关闭弹窗
                    layer.close(index);
                    return false;
                } else {
                    layer.msg(data.RetMsg);
                    return false;
                }
            });
        }
        , cancel: function (index, layero) {
            usrAvartar.attr("src", usrAvartar.data("oldavatarurl"));
            layer.close(index);
        }
        , content: $('#userAvatarSelect')
    });
}
$("#userAvatarSelect").on("click",'.user-avatar li', function(){
    var $this = $(this);
    var num = $this.find("div>img").data("num");

    $this.addClass('selected').siblings().removeClass('selected');
    $("#hid_headimg").val(num);
});
//点击编辑
function EditNickName(){}
//点击保存
function SaveNickName(nickname) {
    $.getJSON(www + "mine/modifynickname", {nickname: encodeURIComponent(nickname)}, function (data) {
        if (data.RetCode == "0") {
            $("#editUsrname b").text(data.Message).attr("title", data.Message);
            //关闭弹窗
            layer.close();
            return false;
        } else {
            layer.msg(data.RetMsg);
            return false;
        }
    });
}

$("#btn_sure").click(function () {
   closePop();
});
// 验证
function Validate(hastel) {
    var telReg = /^[1][3,4,5,7,8][0-9]{9}$/; // 格式设定
    var otel = $('#usrMobile'),
        strtel = $.trim(otel.val());

    var form_ele="";
    if(!$("#userRegist").is(":hidden"))
    {
        // 非空
        form_ele="form_reg input";
    }
    if(!$("#bindPhone").is(":hidden"))
    {
        // 非空
        otel = $('#usrbindMobile');
        strtel = $.trim(otel.val());
        form_ele="form_bindPhone input";
    }
    if(!$("#usrModyPwd").is(":hidden"))
    {
        // 非空
        form_ele="form_modifypwd input";
    }

    //非空
    $("#" + form_ele).each(function (index, el) {
        if ($.trim($(this).val()) == "") {
            errMsg = $(this).data("type") + "不能为空";
            $(this).focus();
            return false;
        }
    });

    if (errMsg != "") {
        return false;
    }
    // 格式
    if(hastel==1) {
        if (!telReg.test(strtel)) {
            otel.focus();
            errMsg = "请输入正确的手机号！";
            return false;
        }
    }
    return true;
}


//初始化广告中的活动链接
function initActivityAd() {
 setTimeout(function () {
     var adHref=$(".adver-dw a").attr("href")
     if(adHref!=undefined&&adHref.indexOf("showActivity")>0){
         var url=$(".adver-dw a").attr("href");
         url=getActivityUrlWithSSO(url);
         $(".adver-dw a").attr("href",url);
         if(HasMobile==false){
             var querenLayer;
             $(".adver-dw a").click(function () {
                 //询问框
                querenLayer = layer.confirm('您还没有绑定手机，请绑定手机后再参与活动', {
                     title:"",
                     btn: ['去绑定手机','再看看'] //按钮
                 }, function(){
                     // 绑定手机
                     LoadValidCode();
                     openPop("bindPhone","手机绑定",430,310);
                     layer.close(querenLayer);
                 }, function(){
                     // 放弃绑定
                 });
                 return false;
             });
         }
     }
 },600)

    //上方的广告
    setTimeout(function () {
        var adHref=$(".adver-up a").attr("href")
        if(adHref!=undefined&&adHref.indexOf("showActivity")>0){
            var url=$(".adver-up a").attr("href");
            url=getActivityUrlWithSSO(url);
            $(".adver-up a").attr("href",url);
            if(HasMobile==false){
                var querenLayer;
                $(".adver-up a").click(function () {
                    //询问框
                    querenLayer= layer.confirm('您还没有绑定手机，请绑定手机后再参与活动', {
                        title:"",
                        btn: ['去绑定手机','再看看'] //按钮
                    }, function(){
                        // 绑定手机
                        LoadValidCode();
                        openPop("bindPhone","手机绑定",430,310);
                        layer.close(querenLayer);
                    }, function(){
                        // 放弃绑定
                    });
                    return false;
                });
            }
        }
    },600)
}

function getActivityUrlWithSSO(url) {
    var ssoStr=window.location.search;
    if(ssoStr.length>0){
        ssoStr=ssoStr.substr(1)
        if(ssoStr.length>0){
            ssoStr=ssoStr.replace("&openreg=1","")
            url= url+"&"+ssoStr;
        }
    }
    return url;
}

//初始化用户等级图标
function initUserLevelICO(level) {
    if(level==undefined){
        level=1;
    }
    //银币
    var silverCoinCount=0;
    //金币
    var goldCoinCount=0;
    //元宝
    var ingotCoinCount=0;
    if(level<=3){
        silverCoinCount=level;
    }else if(level<=7){
        silverCoinCount=level-4;
        goldCoinCount=1;
    }else if(level<=9){
        goldCoinCount=level-6;
    }else{
        ingotCoinCount=level-9;
    }
    $.setUserStarsShow(ingotCoinCount,goldCoinCount,silverCoinCount)
}

function initLoginAward() {
    $.post(www+"activity/getserialloginrule",{},function (data) {
        if(data!=undefined&&data.Message!=undefined&&data.RetCode==0){
            $(".getState span").html(data.Message.UserPrompt);
            if(data.Message.LoginDay==3){
                $(".getState button").text("立即领取");
            }else{
                $(".getState button").text("了解详情");
            }
            $(".getState button").click(function () {
                if(HasMobile==false){
                    var querenLayer = layer.confirm('您还没有绑定手机，请绑定手机后再参与活动', {
                            title:"",
                            btn: ['去绑定手机','再看看'] //按钮
                        }, function(){
                            // 绑定手机
                            LoadValidCode();
                            openPop("bindPhone","手机绑定",430,310);
                            layer.close(querenLayer);
                        }, function(){
                            // 放弃绑定
                        });
                }else{
                    var url=getActivityUrlWithSSO(data.Message.ActivityUrl);
                    window.location.href=url;
                }
            });
            $(".getState").show();
        }else{
            $(".getState").hide();
        }
    });
}

function initMyMedal() {
    $.post(www+"activity/getusermedal",{},function (data) {
        if(data!=undefined&&data.Message!=undefined&&data.RetCode==0){
            //文娱标兵
            var wybb=data.Message.WenYuBiaoBingLevel;
            //合格投资者
            var hwbs=data.Message.QualifiedInvestorLevel;
            //互动达人
            var hddr=data.Message.HuDongDaRenLevel;
            //高朋满座
            var gpmz=data.Message.GaoPengManZuoLevel;
            //种子用户
            var seed=data.Message.SeedUserLevel;
            var uid=ParamFunc.getParam("uid")
            $.userDecoration({wybb:wybb,hwbs:hwbs,hddr:hddr,gpmz:gpmz,seed:seed},uid)
        }
    });
}
/*********************************end****************************************/