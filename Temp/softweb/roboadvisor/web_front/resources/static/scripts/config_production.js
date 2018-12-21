//网站根地址
var www="http://dsclient.emoney.cn/zt/";
//socket地址
var wsurl = "ws://ztwsproxy.emoney.cn/ws/onsocket";
//应用AppId
var AppId = '10150';
//longpoll地址
var longpollurl = "http://ztwsproxy.emoney.cn/poll/onpolling";
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