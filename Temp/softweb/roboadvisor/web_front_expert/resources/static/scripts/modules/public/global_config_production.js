if (!window.console) {
    window.console = {};
}
if (!window.console.log) {
    window.console.log = function(msg) {};
}
// window.onerror = killErrors; function killErrors() {  return true;}

window.gConfig = {
    apiHost: "http://dsclient.emoney.cn/expert/",   // API 接口
    expertBaseUrl: "http://dsclient.emoney.cn/expert/",    // 专家资讯地址
    strategyBaseUrl: "http://dsclient.emoney.cn/expert/",   // 策略副窗地址
    myoptUrl: "http://dsclient.emoney.cn/myoptional/",     // 我的自选地址
    staticPath: "http://static.dsclient.emoney.cn/expert/",
    likeDataServer: "http://ds.emoney.cn/newsapi/api/zan",
    likeListDataServer:"http://ds.emoney.cn/newsapi/api/zans/getlist",
    likeDataServerJsonp: 'http://ds.emoney.cn/newsapi/api/zans/submitjsonp',
    zyapiyqq: "http://zyapi.yqq.emoney.cn/"
};

// window.gConfig.virtualServerPath = window.gConfig.apiHost;
// window.gConfig.expertBaseUrl = window.gConfig.apiHost;
// window.gConfig.strategyBaseUrl = window.gConfig.apiHost;

var wsurl = "ws://wsproxy.emoney.cn/ws/onsocket";
var longpollurl = "http://wsproxy.emoney.cn/poll/onpolling";
var de = parseInt(Math.random() * 10000);
var nowtime = "";
var appIDnewsinfo = 1015001;    // 策略资讯
var appIDstrategyinfo = 1015003;    // 策略
var appIDcloudinfo = 1015002;   // 云端资讯
var appIDstocknewsinfo = 1015005; // 股票相关资讯
var appIDlive = 1015006;    // // 直播
var confColumnID = 39;
var ajaxTimeout = 10000;    // ajax超时时间

var liveImg = window.gConfig.staticPath + "static/images/video.png";    // 策略副窗缺省图
var imgplaceholder = window.gConfig.staticPath + "static/images/zhuanjiazixun/placeholder.png";
var defaultAvatar = window.gConfig.staticPath + "static/images/defaultavatar.png";
var imgLoading2 = window.gConfig.staticPath + "static/libs/layer/theme/default/loading-2.gif";
var imgLoadings = window.gConfig.staticPath + "static/images/loading.gif";

var pagerouter = {
    home: window.gConfig.expertBaseUrl + 'page/home',
    liveVideo: window.gConfig.expertBaseUrl + 'page/live_video',
    yaowenArticle: window.gConfig.expertBaseUrl + 'page/yaowen_article',
    yaowenHome: window.gConfig.expertBaseUrl + 'page/yaowen_home',
    yqqHome: window.gConfig.expertBaseUrl + 'page/yqq_home',
    yuceArticle: window.gConfig.expertBaseUrl + 'page/yuce_article',
    zhuanjiacelueArticle: window.gConfig.expertBaseUrl + 'page/zhuanjiacelue_article',
    zhuanjiacelueHome: window.gConfig.expertBaseUrl + 'page/zhuanjiacelue_home',
    zhuanjiacelueList: window.gConfig.expertBaseUrl + 'page/zhuanjiaceluelist',
    zhuti: window.gConfig.expertBaseUrl + 'page/zhuti',
    zhutiArticle: window.gConfig.expertBaseUrl + 'page/zhuti_article',
    yqqUserLive: 'http://yqq.emoney.cn/Live/UserLive',
    yqqPage: 'http://yqq.emoney.cn',
    share: window.gConfig.expertBaseUrl + "page/shareArticle",
    celueHome: window.gConfig.expertBaseUrl + "page/celue_home",
    celueArticle: window.gConfig.expertBaseUrl + "page/celuearticle",
    articleZixun: window.gConfig.myoptUrl + "myoptional/relatedarticle",
    liveVideoOnline: 'http://ds.emoney.cn/Video/live3/playeritemlist',
    //策略副窗
    celueReading: window.gConfig.strategyBaseUrl + "strategyservice/celuereading",
    celueVideo: window.gConfig.strategyBaseUrl + "strategyservice/celuevideo"
}