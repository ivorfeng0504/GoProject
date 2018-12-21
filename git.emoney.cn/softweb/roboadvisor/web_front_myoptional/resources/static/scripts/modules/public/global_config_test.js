if (!window.console) {
    window.console = {};
}
if (!window.console.log) {
    window.console.log = function(msg) {};
}
// window.onerror = killErrors; function killErrors() {  return true;}

window.gConfig = {
    apiHost: "http://test.roboadvisor.emoney.cn/myoptional/",
    myoptUrl: "http://test.roboadvisor.emoney.cn/myoptional/",
    staticPath: "http://test.static.emoney.cn:8081/myoptional/",
    expertBaseUrl: "http://test.roboadvisor.emoney.cn/expert/",
    strategyBaseUrl: "http://test.roboadvisor.emoney.cn/expert/",
    readApiHost: 'http://test.roboadvisor.emoney.cn/expert/',
    likeDataServer: "http://ds.emoney.cn/newsapi/api/zan",
    likeListDataServer:"http://ds.emoney.cn/newsapi/api/zans/getlist",
    expertBaseUrl: "http://test.roboadvisor.emoney.cn/expert/",
    likeDataServerJsonp: 'http://ds.emoney.cn/newsapi/api/zans/submitjsonp',
};

var wsurl = "ws://wsproxy.emoney.cn/ws/onsocket";
var longpollurl = "http://wsproxy.emoney.cn/poll/onpolling";
var de = parseInt(Math.random() * 10000);
var nowtime = "";
var appIDnewsinfo = 1015001;
var appIDstrategyinfo = 1015003;
var appIDcloudinfo = 1015002;   // 资讯
var appIDstocknewsinfo = 1015005; // 股票相关资讯
var appIDlive = 1015006;     // 直播
var confColumnID = 16;
var ajaxTimeout = 10000;    // ajax超时时间
var myoptTjAppid = 10170;     // 我的自选埋点 appid

var defaultAvatar = window.gConfig.staticPath + "static/images/defaultavatar.png";

var pagerouter = {
    yqqUserLive: "http://yqq.emoney.cn/Live/UserLive",
    yqqPage: "http://yqq.emoney.cn",
    zhuanjiacelueList: window.gConfig.expertBaseUrl + "page/zhuanjiaceluelist",
    celueArticle: window.gConfig.expertBaseUrl + "page/celuearticle",
    articleZixun: window.gConfig.myoptUrl + "myoptional/relatedarticle",
    expertHome: window.gConfig.expertBaseUrl + "page/home"
};