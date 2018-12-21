if (!window.console) {
    window.console = {};
}
if (!window.console.log) {
    window.console.log = function(msg) {};
}
// window.onerror = killErrors; function killErrors() {  return true;}
function preppendZero(num) {
    if (num < 10) {
        return '0' + num;
    }
    return num;
}
var _date = new Date();
var version = _date.getFullYear() + preppendZero(_date.getMonth() + 1) + preppendZero(_date.getDate()) + preppendZero(_date.getHours());

window.gConfig = {
    env: "serv",
    apiHost: "http://dsclient1.emoney.cn:8081/expert/",     // api地址
    trainapiHost:"http://dsclient1.emoney.cn:8081/train/",// 用户培训接口地址
    expertBaseUrl: "http://dsclient1.emoney.cn:8081/",  // 专家资讯地址
    myoptUrl: "http://dsclient1.emoney.cn:8082/",    // 我的自选地址
    strategyBaseUrl: "http://dsclient1.emoney.cn:8081/",   // 策略副窗地址
    resHost: "http://dsclient1.emoney.cn:8081/res/",     // 用户培训策略列表接口地址
    staticPath: "",
    likeDataServer: "http://ds.emoney.cn/newsapi/api/zan",
    likeListDataServer: "http://ds.emoney.cn/newsapi/api/zans/getlist",
    likeListDataServerJsonp: "http://ds.emoney.cn/newsapi/api/zans/getlistjsonp",   // 点赞列表查询（ie8以下）
    likeDataServerJsonp: 'http://ds.emoney.cn/newsapi/api/zans/submitjsonp',
    zyapiyqq:"http://dsclient.emoney.cn/yqqapi/",
    advertising:'http://static.emoney.cn/ds/jrtpad/jrtpadtongji_iv.js'  // 广告统计
};
var wsurl = "ws://wsproxy.emoney.cn/ws/onsocket";
var longpollurl = "http://wsproxy.emoney.cn/poll/onpolling";
var de = parseInt(Math.random() * 10000);
var nowtime = "";
var appIDnewsinfo = 1015001;    // 策略资讯
var appIDcloudinfo = 1015002;   // 云端资讯
var appIDstrategyinfo = 1015003;    // 策略学习资讯
var appIDweiguba = 1015004;  // 微股吧资讯
var appIDstocknewsinfo = 1015005; // 股票相关资讯
var appIDlive = 1015006;    // // 直播间
var appIDLiveInfo = 1015007;    // 直播内容
var appIDTrainingInfo = 1015008;	// 策略培训
var configkey = '752cdab6-671f-4628-8511-f2dab291a3e3';// 用户培训

var confColumnID = 16;
var ajaxTimeout = 10000;    // ajax超时时间
var TagID =1; 
// 统计埋点：expert 专家策略appid,strategy 策略appid
var tjAppid = {
    expert: '10163', strategy: '10167',  usertraining: '10150'
};

var liveImg = window.gConfig.staticPath + "static/images/video.png";    // 策略副窗缺省图  
var imgplaceholder = window.gConfig.staticPath + "static/images/zhuanjiazixun/placeholder.png";
var imgDefault = window.gConfig.staticPath + "static/images/default-pic.png";
var defaultAvatar = window.gConfig.staticPath + "static/images/defaultavatar.png";
var imgLoading2 = window.gConfig.staticPath + "static/libs/layer/theme/default/loading-2.gif";
var imgLoadings = window.gConfig.staticPath + "static/images/loading.gif";
var defaultAudioImg = window.gConfig.staticPath + "static/images/audioplay_cover.png";
var hasNoLive = window.gConfig.staticPath + "static/images/hasnolive.png";
var videoMediaSDK = window.gConfig.staticPath + "static/htmls/videoMediaSDK.html";   // 媒体播放
var imageURL = window.gConfig.staticPath + "static"

var pagerouter = {
    home: window.gConfig.expertBaseUrl +"home.html",
    liveVideo: window.gConfig.expertBaseUrl +"live_video.html",
    yaowenArticle: window.gConfig.expertBaseUrl +"yaowen_article.html",
    yaowenHome: window.gConfig.expertBaseUrl + "yaowen_home.html",
    yqqHome: window.gConfig.expertBaseUrl +"yqq_home.html",
    yuceArticle: window.gConfig.expertBaseUrl +"yuce_article.html",
    // zhuanjiacelueArticle: window.gConfig.expertBaseUrl +"zhuanjiacelue_article.html",
    celueArticle: window.gConfig.expertBaseUrl +"celue_article.html",
    zhuanjiacelueHome: window.gConfig.expertBaseUrl +"zhuanjiacelue_home.html",
    zhuanjiacelueList: window.gConfig.expertBaseUrl +"zhuanjiacelue_list.html",
    zhuanjiacelueVideo: window.gConfig.expertBaseUrl +"zhuanjiacelue_video.html",
    zhuanjiacelueArticle: window.gConfig.expertBaseUrl +"zhuanjiacelue_article.html",
    zhuti: window.gConfig.expertBaseUrl +"zhuti.html",
    zhutiArticle: window.gConfig.expertBaseUrl +"zhuti_article.html",
    yqqUserLive: "http://yqq.emoney.cn/Live/UserLive",
    yqqPage: "http://yqq.emoney.cn",
    share: window.gConfig.expertBaseUrl +"shareArticle.html",
    celueHome: window.gConfig.expertBaseUrl +"celue_home.html",
    celueArticle: window.gConfig.expertBaseUrl + "celue_article.html",
    articleZixun: window.gConfig.myoptUrl + "relatedArticle.html",
    liveVideoOnline: 'http://ds.emoney.cn/Video/live3/playeritemlist',
    
    //策略副窗
    celueReading: window.gConfig.strategyBaseUrl + "celueReading.html",
    celueVideo: window.gConfig.strategyBaseUrl + "celueVideo.html",
    celueTraining: window.gConfig.strategyBaseUrl + "celueTraining.html"
};