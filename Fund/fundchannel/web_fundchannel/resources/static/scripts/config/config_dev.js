//网站根地址
var www="http://127.0.0.1:8186/fundchannel/";
//socket地址
var wsurl = "ws://wsproxy.emoney.cn/ws/onsocket";

//应用AppId
var AppId = '10094';
var GroupId = "2";

//longpoll地址
var longpollurl = "http://wsproxy.emoney.cn/poll/onpolling";

//获取直播室信息URL
var LiveRoomUrl="http://pre.dsclient.emoney.cn/yqqapi/api/Live/GetLiveRoom?lid=2";

//获取直播内容
var LiveContentUrl="http://pre.dsclient.emoney.cn/yqqapi/api/Live/GetLiveContentInfoByType?lid=2&page=0&type=1,2,8&date=";

//获取我的问答内容URL
var LiveMyQuestion="http://pre.dsclient.emoney.cn/yqqapi/api/Live/GetMyLiveQuestionByDay?lid=2&page=0";

//获取全部问答内容URL
var LiveAllQuestion="http://pre.dsclient.emoney.cn/yqqapi/api/Live/GetTodayLiveQuestionDes?lid=2&page=0";

//提问URL
var LiveSubmitQuestion="http://pre.dsclient.emoney.cn/yqqapi/api/WeChat/SetUserQuestion?lid=2&source=7&question=";

var LiveAllDataUrl="http://pre.dsclient.emoney.cn/yqqapi/api/Live/SSOAndGetFundLiveData?lid=2&page=0&type=1,2,8&date=";


