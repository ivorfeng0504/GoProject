var UserAwardState_Grant=1;
var UserAwardState_Receive=2;
var HasMobile=false;

$(function () {
    initAward(UserAwardState_Grant);
    initProduct();

    //我的产品
    $("#btn_product").click(function () {
        initProduct();
    });

    //我的福利
    $("#btn_award").click(function () {
        initAward(UserAwardState_Receive);
    });

});

//加载福利数据
function initAward(state) {
    $.post(www+"mine/getuserawardlist",{
        state:state
    },function(data){
        if(data.RetCode == "0"){
            switch(state){
                case UserAwardState_Grant:
                    var $div_mine_grant = $("#div_mine_grant");
                    var html = template_mine_grant(data);
                    $div_mine_grant.html(html);
                    if(data!=null&&data.Message!=undefined&&data.Message.length>0){
                        $("#txt_award_count").html("<i class=\"icon\">&#xe87a;</i>您还有个"+data.Message.length+"【特权】未领取")
                    }else{
                        // $("#txt_award_count").html("");
                        //没有可领取的特权则展示广告位
                        initMineCenterAd()
                    }
                    initAwiper();
                    break;
                case UserAwardState_Receive:
                    var $div_mine_receive = $("#div_mine_receive");
                    var html = template_mine_receive(data);
                    $div_mine_receive.html(html);
                    initAdvisoryClass();
                    break;
            }
        }else{

        }
    });
}


//加载产品数据
function initProduct() {
    $.post(www+"mine/getuserproductlist",{
    },function(data){
        if(data.RetCode == "0"){
            var $div_mine_product = $("#div_mine_product");
            var html = template_mine_product(data);
            $div_mine_product.html(html);
            initAdvisoryClassProduct();
        }else{

        }
    });
}

//初始化
function initAwiper() {
    var swiper = new Swiper('.swiper-container', {
        slidesPerView: 1,
        paginationClickable: true,
        loop: false

    });

    $('.swiper-button-prev').click(function () {
        swiper.swipePrev();
    })
    $('.swiper-button-next').click(function () {
        swiper.swipeNext();
    })
}

//初始化广告中的活动链接
function initActivityAd() {
    setTimeout(function () {
        var adHref=$("#ad_mine_center a").attr("href")
        if(adHref!=undefined&&adHref.indexOf("showActivity")>0){
            var url=$("#ad_mine_center a").attr("href");
            url=getActivityUrlWithSSO(url);
            $("#ad_mine_center a").attr("href",url);
            if(HasMobile==false){
                $("#ad_mine_center a").click(function () {
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