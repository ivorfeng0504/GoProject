define(["jquery","cors","utils","layer","handlebars","nicescroll"],function($,$cors,utils,layer,Handlebars,nicescroll){layer.config({path:window.gConfig.staticPath+"static/libs/layer/"}),$.support.cors=!0;var page={nonedataHtml:'<div class="no-info">暂无数据显示，请稍后重试</div>',init:function(){utils.appendSSO();var n=this;n.clickhotNews=!1,n.flag=!0,n.index=0,n.firstArticle=!0,n.getUID=utils.getUID(),n.getAppId=appIDcloudinfo,n.NewsInfoId=utils.GetQueryString("NewsID"),n.source=utils.GetQueryString("source"),n.from=utils.GetQueryString("from"),n.scrollContent(),n.hotNews(0),n.getContentList(0),utils.thumbUP(".like",n.getUID,n.getAppId,"#yaowenScroll"),n.share(),window.pulldownLoad=function(){var e,t=null,a=!0;$(function(){(t=$("#yaowenScroll").getNiceScroll(0)).jqbind("#yaowenScroll","scroll",function(){t.newscrolly>=t.page.maxh-500&&a&&(a=!a,clearTimeout(e),e=setTimeout(function(){n.index>n.newDataArrayLength-1&&n.flag?(layer.tips("已加载全部","#yaowenScroll",{offset:["400px","50px"],tips:[2,"#cb4b4b"]}),n.flag=!1):n.flag&&n.updateInfo(n.newDataArray[n.index].DataUrl),a=!a},500))}),t.scrollend(function(){var s,e=$("#yaowenScroll .atticle-show");$.each(e,function(e,t){var a=$(t),n=a.offset().top,i=a.height();(n<=100&&i*(2/3)<n+i-60||100<n&&i*(2/3)<553-n+86)&&(s=a.attr("data-artid"))}),s!=undefined&&n.getNewsIDreadingCount(s)})})},window.pulldownLoad(),$("#goBack").on("click",function(){1<history.length?window.history.back():(window.location.href=pagerouter.yaowenHome,window.event.returnValue=!1)})},scrollContent:function(){$.each(["#yaowenScroll"],function(e,t){0<$(t).length&&$(t).niceScroll({cursorcolor:"#666",cursoropacitymax:.5,touchbehavior:!1,cursorwidth:"6px",cursorborder:"0",cursorborderradius:"8px",autohidemode:!1})})},getNewDataArray:function(e,t){var a=this;if(a.detailCode=["<div class='atticle-show section yaowen-article' id=\"atticleShow{{newsId}}\" data-artid={{newsId}}>","        <div class='a-iner'>","            <div class='art-title'>",'                <h3 id="artcleTitle" title="">{{article_title}}</h3>',"            </div>",'            <div class="art-subtitle">',"                <div class=' pull-right'>","                    <span class='view-tips'>",'                        <span class="vt-cell like" id="like{{newsId}}" data-channelid={{newsId}}>','                            <i class=\'icon-1\' id="likeicon{{newsId}}">&#xe61b;</i><span class="liked" id="liked{{newsId}}"></span></span>','                    <span class="vt-cell">',"                            <i class='icon-1'>&#xe602;</i><span id=\"reading{{newsId}}\"></span></span>",'                    <span class="vt-cell">','                        <div class="fl ml10 bdsharebuttonbox" data-tag="share_{{newsId}}" style=" display: inline-block; vertical-align: middle;  margin-top: -5px;">','                            <a href="javascript:;" class="bds_weixin" data-cmd="weixin" title="分享到微信" data-id="{{newsId}}" data-title="{{article_title}}"></a>','                            <a href="javascript:" class="bds_tsina" data-cmd="tsina" title="分享到新浪微博" data-id="{{newsId}}" data-title="{{article_title}}"></a>',"                        </div>","                    </span>","                    </span>","                </div>",'                <div class="ac-submain">',"                    <span>来源：{{article_source}}</span>","                    <span>{{renderTime publish_time}}</span>","                </div>","            </div>",'            <div class="art-cnt">','                <div class="ac-iner">','                    <div class="ac-content  ywac-content" id="acScrollbox" data-fontsizecnt>',"                        {{htmled content}}","                    </div>","                </div>","            </div>","        </div>","    </div>"].join(""),a.newDataArrayLength=e.Message.length,a.newDataArray=[],a.preDataArray=[],a.nextDataArray=[],0!=a.from||a.clickhotNews)for(i=0;i<e.Message.length;i++)e.Message[i].NewsInformationId==t&&(a.preDataArray=e.Message.slice(0,i),a.nextDataArray=e.Message.slice(i+1,e.Message.length),a.newDataArray.push(e.Message[i]),a.newDataArray=a.newDataArray.concat(a.nextDataArray).concat(a.preDataArray));else{for(var n=0;n<e.Message.length;n++)for(var i=0;i<a.hotNewsData.Message.length;i++)e.Message[n].NewsInformationId==a.hotNewsData.Message[i].NewsInformationId&&a.hotNewsData.Message.splice(i,1);for(var i=0;i<e.Message.length;i++)e.Message[i].NewsInformationId==t&&(a.preDataArray=e.Message.slice(0,i),a.nextDataArray=e.Message.slice(i+1,e.Message.length),a.newDataArray.push(e.Message[i]),a.newDataArray=a.newDataArray.concat(a.nextDataArray).concat(a.preDataArray).concat(a.hotNewsData.Message),a.newDataArrayLength=e.Message.length+a.hotNewsData.Message.length)}},updateInfo:function(e){var a=this,n=$("#yaowenScroll");if(a.readTheFirstFourArticles=function(){a.index<4&&a.updateInfo(a.newDataArray[a.index].DataUrl)},$.support.cors=!0,a.currentArticleId=a.newDataArray[a.index].NewsInformationId,Handlebars.registerHelper("htmled",function(e){return e!=undefined&&null!=e&&0!=e.length||(e=""),new Handlebars.SafeString(e)}),Handlebars.registerHelper("renderTime",function(e){return e=e==undefined||null==e||0==e.length?"":"　"+e.slice(5,10)+" "+e.slice(11,16),new Handlebars.SafeString(e)}),window.XDomainRequest){var i=new XDomainRequest;i.open("get",e),i.onprogress=function(){},i.ontimeout=function(){},i.onerror=function(){n.html(page.nonedataHtml)},i.onload=function(){var e=JSON.parse(i.responseText);if(e.id!=undefined&&null!=e.id){$.extend(e,{newsId:a.currentArticleId});var t=Handlebars.compile(a.detailCode)(e);n.find(".no-info").hide(),n.removeClass("no-info-cnt").append(t),a.getReadCount(a.currentArticleId),a.clickLike(a.getUID,a.currentArticleId,a.getAppId),a.index=a.index+1,a.readTheFirstFourArticles(),a.share(),window._bd_share_main&&window._bd_share_main.init&&window._bd_share_main.init()}else n.html(page.nonedataHtml)},setTimeout(function(){i.send()},0)}else $.ajax({url:e,type:"GET",timeout:ajaxTimeout,dataType:"json",success:function(e){if(e.id!=undefined&&null!=e.id){a.tip=null,$.extend(e,{newsId:a.currentArticleId});var t=Handlebars.compile(a.detailCode)(e);n.find(".no-info").hide(),n.removeClass("no-info-cnt").append(t),a.getReadCount(a.currentArticleId),a.clickLike(a.getUID,a.currentArticleId,a.getAppId),a.index=a.index+1,a.readTheFirstFourArticles(),a.share(),window._bd_share_main&&window._bd_share_main.init&&window._bd_share_main.init()}else n.html(page.nonedataHtml)},error:function(e,t,a){n.html(page.nonedataHtml)}})},getContentList:function(a){var n=this,i=$("#yaowenScroll");$.ajax({url:window.gConfig.apiHost+"expertnews/getnewstop30",type:"GET",timeout:ajaxTimeout,dataType:"json",success:function(e,t){0==e.RetCode&&5==n.source&&(n.getNewDataArray(e,n.NewsInfoId),n.updateInfo(n.newDataArray[a].DataUrl))},error:function(e,t,a){i.html(page.nonedataHtml)}})},hotNews:function(a){var n=this,i=$("#hotNewsList"),s=["{{#each Message}}",'<li cid="{{NewsInformationId}}" clickkey="articleAside|list" clickdata="{{NewsInformationId}}|yaowen_article" clickremark="" >','<a href="javascript:;" >','<i class="icon-1">&#xe637;</i>',"<span>{{ArticleTitle}}</span>","</a>","</li>","{{/each}}"].join("");$.ajax({url:window.gConfig.apiHost+"expertnews/gethotnews",contentType:"application/json",type:"GET",timeout:ajaxTimeout,dataType:"JSON",success:function(t){if(0==t.RetCode){n.hotNewsData=JSON.parse(JSON.stringify(t));var e=Handlebars.compile(s)(t);$("#hotNewsList").empty().append(e),3==n.source&&0!=n.from&&(n.getNewDataArray(t,n.NewsInfoId),n.updateInfo(n.newDataArray[a].DataUrl)),3==n.source&&0==n.from&&n.todayTopNews(0),$("#hotNewsList li").on("click",function(){n.clickhotNews=!0,n.index=0;var e=$(this);$("#yaowenScroll").addClass("no-info-cnt").html('<div class="no-info"><span class="inerloading"></span>正在加载...</div>'),n.hotNewsId=e.attr("cid"),n.firstArticle=!0,n.getNewDataArray(t,n.hotNewsId),n.updateInfo(n.newDataArray[a].DataUrl)})}else i.html(page.nonedataHtml)},error:function(e,t,a){i.html(page.nonedataHtml)}})},todayTopNews:function(a){var n=$("#yaowenScroll"),i=this;$.ajax({url:window.gConfig.apiHost+"expertnews/gettodaynews",type:"GET",timeout:ajaxTimeout,dataType:"json",success:function(e,t){0==e.RetCode?3==i.source&&0==i.from&&(i.getNewDataArray(e,i.NewsInfoId),i.updateInfo(i.newDataArray[a].DataUrl)):n.html(page.nonedataHtml)},error:function(e,t,a){n.html(page.nonedataHtml)}})},getReadCount:function(t){var a=this;$.ajax({url:window.gConfig.apiHost+"click/queryclick",type:"GET",timeout:ajaxTimeout,dataType:"JSON",cache:!1,data:{identity:t,clickType:"news.information"},success:function(e){$("#reading"+t).text(e.Message.Result[t]),1==a.firstArticle&&(a.getNewsIDreadingCount(a.newDataArray[0].NewsInformationId),a.firstArticle=!1)},error:function(e,t,a){}})},getNewsIDreadingCount:function(e){var a;e!=undefined&&""!=e&&null!=e&&0!=e&&((a=$("#reading"+e)).attr("viewed")||(a.text(parseInt(a.text())+1),a.attr("viewed","1"),$.ajax({type:"get",url:window.gConfig.apiHost+"click/addclick",contentType:"text/plain",cache:!1,timeout:ajaxTimeout,data:{identity:e,clickType:"news.information"},dataType:"json",success:function(e,t){"success"===t&&0==e.RetCode&&"SUCCESS"==e.RetMsg&&a.text(e.Message)},error:function(e,t,a){}})))},clickLike:function(e,n,t){$.ajax({type:"get",timeout:ajaxTimeout,url:window.gConfig.likeDataServer,contentType:"text/plain",dataType:"json",data:{uid:e,newsId:n,appId:t},success:function(e,t){e.isSucess?$("#like"+n).attr("data-channelid")==n&&($("#liked"+n).text(e.message.likes),!0===e.message.liked?($("#likeicon"+n).html("&#xe66b;"),$("#like"+n).addClass("liked")):($("#likeicon"+n).html("&#xe61b;"),$("#like"+n).removeClass("liked"))):($("#likeicon"+n).html("&#xe61b;"),$("#like"+n).removeClass("liked"))},error:function(e,t,a){$("#likeicon"+n).html("&#xe61b;"),$("#like"+n).removeClass("liked")}})},share:function(){var _this=this,shareId="",title="";function SetConf(e,t){return shareId&&(t.bdUrl=url,t.bdText=title),t}with($(function(){$(".bdsharebuttonbox a").mouseover(function(){shareId=$(this).attr("data-id"),title=$(this).attr("data-title"),url=pagerouter.share+"?type=1&NewsID="+shareId+"&source="+_this.source})}),window._bd_share_config={common:{onBeforeClick:SetConf},share:[{bdSize:12}]},document)(getElementsByTagName("head")[0]||body).appendChild(createElement("script")).src="http://bdimg.share.baidu.com/static/api/js/share.js?v=89860593.js?cdnversion="+~(-new Date/36e5)}};page.init(),$.emoneyAanalytics().Init(tjAppid.expert,"yaowen_article","")});