var captchaid = "";//验证码使用
$(document).ready(function() {
    initProfileInfo();
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
//个人资料初始化
function initProfileInfo() {
    $.getJSON(www + "mine/getprofile", {}, function (data) {
        if (data.RetCode == "0") {
            var Profile = data.Message;
            var account = Profile.Account;
            var headimg = Profile.Headportrait;
            var nickname = Profile.NickName;
            var bindAccountList = Profile.BindAccountList;
            var province=Profile.Province;
            var city=Profile.City;
            var serviceAgentName=Profile.ServiceAgentName;

            if (headimg != "") {
                headimg = www + "/static/images/Arena_" + headimg + ".png";
            } else {
                headimg = www + "/static/images/Arena_13.png";
            }
            //$("#pro_account").html(account);
            //$("#pid_xy").attr("href","http://fuwu.emoney.cn/Register/RiskControl?pid="+Profile.PID);
            $("#pro_nickname").html(nickname);
            $("#pro_headimg").attr("src",headimg);
            if(province==""){
                $("#div_province").hide();
            }else{
                $("#pro_province").html(province+"、"+city);
            }

            if(serviceAgentName==""){
                $("#div_serviceAgentName").hide();
            }else{
                $("#pro_serviceAgentName").html(serviceAgentName);
            }

            if(bindAccountList!=null) {
                var i = 0;
                bindAccountList.forEach(function (item, index, arr) {
                    var accountType = arr[index].AccountType;//0: em帐号; 1: 手机号; 2: 微信帐号; 3: QQ帐号
                    var customerID = arr[index].CustomerID;//CID
                    var accountName = arr[index].AccountName;//账号显示名
                    var mobilex = arr[index].EncryptMobile;//加密手机号

                    var div_bindPhone = $("#pro_PhoneBind");
                    var div_bindEMCard = $("#pro_EMBind");
                    var div_bindWeChat = $("#pro_WeChatBind");
                    var div_bindQQ = $("#pro_QQBind");
                    switch (accountType) {
                        case 0:
                            if(i==0) {
                                if (accountName==Profile.UserName)
                                {
                                    div_bindEMCard.html(accountName);
                                }else {
                                    div_bindEMCard.html(accountName + "     <a href='###' name='btn_EMBind' data-type='remove' data-cid='" + customerID + "'>解除绑定</a>");
                                }
                            }
                            else{
                                if (accountName==Profile.UserName) {
                                    div_bindEMCard.html(accountName);
                                }else {
                                    div_bindEMCard.append(accountName + "     <a href='###' name='btn_EMBind' data-type='remove' data-cid='" + customerID + "'>解除绑定</a>");
                                }
                            }
                            i++;
                            break;
                        case 1:
                            if (mobilex==Profile.UserName) {
                                div_bindPhone.html(accountName);
                            }else {
                                div_bindPhone.html(accountName + "     <a href='###' name='btn_PhoneBind' data-type='remove' data-cid='" + customerID + "'>解除绑定</a>");
                            }
                            break;
                        // case 2:
                        //     if (accountName==Profile.UserName) {
                        //         div_bindWeChat.html(accountName);
                        //     }else {
                        //         div_bindWeChat.html("已绑定微信账号(" + accountName + ")     <a href='###' name='btn_WeChatBind' data-type='remove' data-cid='" + customerID + "'>解除绑定</a>");
                        //     }
                        //     break;
                        // case 3:
                        //     if (accountName==Profile.UserName) {
                        //         div_bindQQ.html(accountName);
                        //     }else {
                        //         div_bindQQ.html("已绑定QQ账号(" + accountName + ")     <a href='###' name='btn_QQBind' data-type='remove' data-cid='" + customerID + "'>解除绑定</a>");
                        //     }
                        //     break;
                    }
                });
            }
        }
    });
}

//解除绑定
var removeflag = true;
function RemoveBind(cid,type,typename) {
    var popcnt = layer.open({
        type: 1
        , title: ['温馨提示', '']
        , skin: 'ucpop-layer'
        , area: ['430px', '280px']
        , scrollbar: false
        , btn: ['确定解绑']
        , btnAlign: 'c'
        , yes: function (index) {
            if (removeflag) {
                removeflag = false;
                $.getJSON(www + "mine/removebind", {cid: cid}, function (data) {
                    removeflag = true;
                    if (data.RetCode == "0") {
                        layer.msg("成功解除绑定");
                        setTimeout(function () {
                            window.location.reload();
                            layer.closeAll();
                        },3000);
                        //$("#pro_" + type).html("<a href='javascritp:;' name='btn_" + type + "' data-type='add' data-cid='0'>绑定" + typename + "</a>");
                    } else {
                        layer.msg(data.retMsg);
                        return;
                    }
                });
            }
        }
        , content: $('#div_removeTips')
    });
}

//点击绑定或解除绑定
$(".usr-infobox").on("click","[name=btn_PhoneBind]",function () {
    var cid=$(this).data("cid");
    var type=$(this).data("type");
    if (type=="add"){
        //弹出手机绑定窗口
        LoadValidCode();
        openPop("bindPhone","手机绑定",430,320);
    }
    if(type=="remove"){
        //解除绑定
        RemoveBind(cid,"PhoneBind","手机号")
    }
}).on("click","[name=btn_EMBind]",function () {
    var cid=$(this).data("cid");
    var type=$(this).data("type");
    if (type=="add"){
        //弹出em绑定窗口
        LoadValidCode();
        openPop("bindEM","EM账号绑定",430,320);
    }
    if(type=="remove"){
        //解除绑定
        RemoveBind(cid,"EMBind","EM账号")
    }
}).on("click","[name=btn_WeChatBind]",function () {
    var cid=$(this).data("cid");
    var type=$(this).data("type");
    if (type=="add"){
        window.location.href = www + "mine/wechatlogin";
    }
    if(type=="remove"){
        //解除绑定
        RemoveBind(cid,"WeChatBind","微信账号")
    }
}).on("click","[name=btn_QQBind]",function () {
    var cid=$(this).data("cid");
    var type=$(this).data("type");
    if (type=="add"){
        window.location.href = www + "mine/qqlogin";
    }
    if(type=="remove"){
        //解除绑定
        RemoveBind(cid,"QQBind","QQ账号")
    }
});

//获取密码
var tnum = 0;
var flagtime;
var flag = true;
$("#btn_getPwd").click(function () {
    if (flag) {
        var mobile = $("#usrbindMobile").val();
        if (mobile == "") {
            layer.msg("请输入手机号码！");
            $("#usrbindMobile").focus();
            return;
        }
        var telReg = /^[1][3,4,5,7,8][0-9]{9}$/;
        if (mobile.length != 11 || (telReg.test(mobile)) == false) {
            layer.msg("请输入正确的手机号！");
            $("#usrbindMobile").focus();
            return;
        }

        flag = false;
        $.getJSON(www + "mine/getpwdbymobile", {Mobile: do_encrypt(mobile)}, function (data) {
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
    $('#btn_getPwd').html('重发(' + num + 'S)');
    flagtime = setTimeout("timejs()", 1000);
    if (tnum > 60) {
        stopjs();
        $('#btn_getPwd').html("获取验证码");
        flag = true;
        tnum = 0;
    }
}
function stopjs() {
    clearTimeout(flagtime);
}

//提交-绑定（绑定手机号）
var phonebind_flag=true;
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
                },3000);
            } else {
                layer.msg(data.RetMsg);
                RefreshValidCode();
                return false;
            }
        });
    }
});

var embind_flag=true;
$("#submit_embind").click(function(){
    if(embind_flag) {
        $("#bindEM").each(function (index, el) {
            if ($.trim($(this).val()) == "") {
                errMsg = $(this).data("type") + "不能为空";
                $(this).focus();
                return false;
            }
        });
        // Next
        var emcard = $("#usrbindEM").val();
        var pwd = $("#bindem_pwd").val();
        var code = $("#bindem_code").val();
        embind_flag=false;
        $.getJSON(www + "mine/bindaccountem", {
            emcard: do_encrypt(emcard),
            pwd: do_encrypt(pwd),
            code: code,
            captchaId: captchaid
        }, function (data) {
            embind_flag=true;
            if (data.RetCode == "0") {
                layer.msg("EM账号绑定成功");
                setTimeout(function () {
                    window.location.reload();
                },3000);
            } else {
                layer.msg(data.RetMsg);
                RefreshValidCode();
                return false;
            }
        });
    }
});

// 验证
function Validate(hastel) {
    var telReg = /^[1][3,4,5,7,8][0-9]{9}$/; // 格式设定

    var otel = $('#usrMobile'),
        strtel = $.trim(otel.val());

    var form_ele="";

    if(!$("#form_bindPhone").is(":hidden")) {
        // 非空
        otel = $('#usrbindMobile');
        strtel = $.trim(otel.val());
        form_ele = "form_bindPhone input";
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

function openPop(obj,title,width,height) {
    $("input[ type='text']").val("");
    $("input[ type='password']").val("");
    layer.open({
        type: 1,
        title: [title, ''],
        skin: 'ucpop-layer', //加上边框
        area: [width + 'px', height + 'px'], //宽高
        content: $('#' + obj) //'弹窗内容'
    });
}
function closePop() {
    layer.closeAll();
}