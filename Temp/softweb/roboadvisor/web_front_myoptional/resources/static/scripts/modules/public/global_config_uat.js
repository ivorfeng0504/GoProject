if (!window.console) {
    window.console = {};
}
if (!window.console.log) {
    window.console.log = function(msg) {};
}
// window.onerror = killErrors; function killErrors() {  return true;}

window.gConfig = {
    apiHost: "http://pre.dsclient.emoney.cn:8080/myoptional/",
    myoptUrl: "http://pre.dsclient.emoney.cn/myoptional/",      // 我的自选地址
    expertBaseUrl: "http://pre.dsclient.emoney.cn/expert/",    // 专家资讯地址
    strategyBaseUrl: "http://pre.dsclient.emoney.cn/expert/",   // 策略地址
    // apilocal: "http://dsclient.emoney.cn:8080/",
    staticPath: "http://static.dsclient.emoney.cn:8080/myoptional/",
    likeDataServer: "http://ds.emoney.cn/newsapi/api/zan",
    likeListDataServer:"http://ds.emoney.cn/newsapi/api/zans/getlist",
    likeDataServerJsonp: 'http://ds.emoney.cn/newsapi/api/zans/submitjsonp',
    expertBaseUrl:"http://pre.dsclient.emoney.cn:8080/expert/",
    readApiHost:'http://pre.dsclient.emoney.cn:8080/expert/'
};

var wsurl = "ws://wsproxy.emoney.cn/ws/onsocket";
var longpollurl = "http://wsproxy.emoney.cn/poll/onpolling";
var de = parseInt(Math.random() * 10000);
var nowtime = "";
var appIDnewsinfo = 1015001;    // 策略资讯
var appIDstrategyinfo = 1015003;    // 策略
var appIDcloudinfo = 1015002;   // 资讯
var appIDstocknewsinfo = 1015005; // 股票相关资讯
var appIDlive = 1015006;     // 直播
var confColumnID = 39;
var ajaxTimeout = 10000;    // ajax超时时间
var myoptTjAppid = 10170;     // 我的自选埋点 appid

var defaultAvatar = window.gConfig.staticPath + "static/images/defaultavatar.png";

var pagerouter = {
    yqqUserLive: "http://yqq.emoney.cn/Live/UserLive",
    yqqPage: "http://yqq.emoney.cn",
    zhuanjiacelueList: window.gConfig.expertBaseUrl +"page/zhuanjiaceluelist",
    celueArticle: window.gConfig.expertBaseUrl +"page/celuearticle",
    articleZixun: window.gConfig.myoptUrl +"myoptional/relatedarticle",
    expertHome: window.gConfig.expertBaseUrl +"page/home"
};