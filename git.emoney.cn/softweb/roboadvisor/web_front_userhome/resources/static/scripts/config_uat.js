//网站根地址
var www="http://pre.dsclient.emoney.cn/userhome/";
//静态服务器地址
var StaticServerHost="http://test.static.emoney.cn:8081/pre/userhome"
//socket地址
var wsurl = "ws://wsproxy.emoney.cn/ws/onsocket";
//应用AppId
var AppId = '10150';
//longpoll地址
var longpollurl = "http://wsproxy.emoney.cn/poll/onpolling";
//提问限制时间（秒）
var questionSubmitLimitSeconds=60;
//延迟拉取最新数据-直播内容（秒）
var liveContentUpdateDelaySeconds=3;
//延迟拉取最新数据-问答内容（秒）
var liveQuestionUpdateDelaySeconds=3;
//延迟拉取最新数据-置顶消息（秒）
var liveTopContentUpdateDelaySeconds=3;
//数据轮询更新时间，如果socket和longpoll都失效才会启用（秒）
var loopUpdateDelaySeconds=10;
//我的产品广告位代码
var Ad_AC_Product="AC2018060105";
//我的福利广告位代码
var Ad_AC_Award="AC2018060104";
//头部-在线咨询广告位代码
var Ad_AC_Online="AC2018060103";
//我的-中心部分广告位代码
var Ad_AC_Mine_Center="AC2018060106";
//PID
var Ad_PID=""
//首页右上广告位
var Ad_AC_Index_Up="AC2018060101";
//首页右下广告位
var Ad_AC_Index_DwFree="AC2018101501";//免费
var Ad_AC_Index_Dw="AC2018060102";//付费

//用户中心活动-领数字
var GNACT_getCodeUrl = www+"activity/getguessmarketnumberwithdetail";
var GNACT_historyListUrl = www+"activity/getguessmarkethistoryinfo";
var GNACT_getAwardUrl = www+"activity/getguessmarketaward";
var GNACT_getActStateUrl = www+"activity/getcurrentguessinfo";

//ajax请求方式
var requestType="POST"
var userHome_columnID = 15