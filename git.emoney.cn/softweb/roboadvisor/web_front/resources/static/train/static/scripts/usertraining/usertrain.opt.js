!function(){var l={nonedataHtml:'<div class="no-info">暂无课程安排，课程回放也很精彩哦</div>',curPageOpts:{NewsID:0,ColumnID:16,StrategyID:0,currpage:1,pageSize:8,LiveID:0,ExpertLiveroom:"http://yqq.emoney.cn/Live/UserLive?lid=93",isbg:!1,newsidzan:[],zanparams:[]},init:function(){var e=this;e.getUID=$.cookie("expertnews.uid")||utils.GetQueryString("uuid")||0,$.support.cors=!0,e.getAppId=appIDTrainingInfo,utils.EMSSO(),utils.appendSSO(),$.cookie("expertnews.uid")?e.getUID=$.cookie("expertnews.uid"):(e.getUID=utils.GetQueryString("uid")||0,$.cookie("expertnews.uid",e.getUID)),e.resetLiveContboxSize(),e.scrollContent(),e.bindEvents(),e.getLatestgData(),e.latestData=[],e.getTagData=[],e.getTag3Data=[],e.getonemonthData=[],e.getStrategyData=[],e.format(),e.pulldownLoad()},pulldownLoad:function(){var i,s=this,r=!0;window.liveFrmNS1=null,$(function(){liveFrmNS1=$("#menuScroll").getNiceScroll(0),liveFrmNS1.scrollend(function(){var e=$("#uthNav .current").find("span").text(),a=$("#strategyList .current").attr("data-classid"),t=$("#uthNav .current").attr("data-date");liveFrmNS1.newscrolly>=liveFrmNS1.page.maxh-50&&0!=liveFrmNS1.page.maxh&&(s.curPageOpts.isbg?s.dataLength<=s.curPageOpts.pageSize&&s.curPageOpts.totalpage==s.curPageOpts.currpage&&!s.flag&&(s.flag=!0,s.curPageOpts.isbg=!0,layer.tips("已加载全部","#menuScroll",{offset:["400px","50px"],tips:[2,"#cb4b4b"]})):r&&(r=!r,clearTimeout(i),i=setTimeout(function(){s.curPageOpts.currpage=parseInt(s.curPageOpts.currpage)+1,"解盘晨会"==e?s.getusertrainingData(a,"jpch",t):"指标学习"==e?s.getusertrainingData(a,"zbxx",t):"当日最新"==e||"策略实战"==e&&$("#strategyList li").eq(0).hasClass("current")&&s.get1monthData(t),r=!r},500)))})})},format:function(){var a=this;Handlebars.registerHelper("parentNameStr",function(e){return e="基本面.策略"==e?"价值":"　"+e.slice(0,2),new Handlebars.SafeString(e)}),Handlebars.registerHelper("firstItem",function(e){return e=0==e&&1==a.curPageOpts.currpage?"current":"",new Handlebars.SafeString(e)}),Handlebars.registerHelper("first",function(e){return e=0==e&&1==a.curPageOpts.currpage?"first-item":"",new Handlebars.SafeString(e)}),Handlebars.registerHelper("tagNameLength",function(e){return e=4<=e.length?"four-word":"three-word",new Handlebars.SafeString(e)}),Handlebars.registerHelper("renderTime",function(e){return e=(e=e.replace("T"," "))==undefined||null==e||0==e.length?"":e.slice(0,10)==utils.getNowDate()?e.slice(11,16):e.slice(5,11)+" "+e.slice(11,16),new Handlebars.SafeString(e)}),Handlebars.registerHelper("renderTimeStrategy",function(e,a){return e=(e=(e="0001-01-01T00:00:00Z"==e||"0001-01-01T00:00:00"==e?a:e).replace("T"," "))==undefined||null==e||0==e.length?"":e.slice(0,10)==utils.getNowDate()?e.slice(11,16):e.slice(5,11)+" "+e.slice(11,16),new Handlebars.SafeString(e)}),Handlebars.registerHelper("decompose",function(e){if(e==undefined||null==e||0==e.length)e="";else{for(var a=JSON.parse(e),t=new Array,i=0;i<a.length;i++){var s=a[i];t.push(s.TagName)}e=t.join("，")}return new Handlebars.SafeString(e)}),Handlebars.registerHelper("controlLines",function(e){return e==undefined||null==e||0==e.length?time="":10<=e.length&&e.length<=17?e=$.trim(e).substr(0,10)+"...":18<=e.length&&e.length<=22?e=e:22<e.length&&(e=$.trim(e).substr(0,22)+"..."),new Handlebars.SafeString(e)}),Handlebars.registerHelper("renderDate",function(e){return e=e==undefined||null==e||0==e.length?"":"　"+e.slice(5,10)+" "+e.slice(11,16),new Handlebars.SafeString(e)}),Handlebars.registerHelper("substr30",function(e){return e=e==undefined||null==e||0==e.length?"":30<=e.length?e.substr(0,30)+"...":e,new Handlebars.SafeString(e)}),Handlebars.registerHelper("renderHTML",function(e){return e=e==undefined||null==e?"":e}),Handlebars.registerHelper("isSpreadThismenu",function(e){var a="";return a=0==e.length?"&#xe72e;":" ",new Handlebars.SafeString(a)}),Handlebars.registerHelper("isexistence",function(e){return e=e==undefined||null==e||0==e.length?"":e,new Handlebars.SafeString(e)}),Handlebars.registerHelper("decompose",function(e){if(e==undefined||null==e||0==e.length)e="";else{for(var a=JSON.parse(e),t=new Array,i=0;i<a.length;i++){var s=a[i];t.push(s.TagName)}e=t.join("，")}return new Handlebars.SafeString(e)}),Handlebars.registerHelper("showstatus",function(e,a,t,i){var s=l.showStatusFunc(e,a,t,i);return new Handlebars.SafeString(s)}),Handlebars.registerHelper("statusClass",function(e,a,t,i,s){var r=null,n="",l=(new Date).getTime();if(e&&a){var d=t?new Date((t||"").replace(/-/g,"/").replace(/[TZ]/g," ").slice(0,19)).getTime():l,o=l-d,c=new Date((e||"").replace(/-/g,"/").replace(/[TZ]/g," ").slice(0,19)).getTime();c+=o;var g=new Date((a||"").replace(/-/g,"/").replace(/[TZ]/g," ").slice(0,19)).getTime();if((g+=o)<l)i?n="replay":a!=e&&(n="recording");else if(c<l&&l<g){n="living",r=setTimeout(function(){$("#livestat"+s).text("录制中"),clearTimeout(r)},g-d)}else l<c&&(n="livenotstart");return new Handlebars.SafeString(n)}n=""}),Handlebars.registerHelper("actralTraingStatusClass",function(e,a,t,i,s,r){var n=null,l="",d=(new Date).getTime();if(e&&a){var o=t?new Date((t||"").replace(/-/g,"/").replace(/[TZ]/g," ").slice(0,19)).getTime():d,c=d-o,g=new Date((e||"").replace(/-/g,"/").replace(/[TZ]/g," ").slice(0,19)).getTime();g+=c;var u=new Date((a||"").replace(/-/g,"/").replace(/[TZ]/g," ").slice(0,19)).getTime();if((u+=c)<d)"0001-01-01T00:00:00Z"!=a&&"0001-01-01T00:00:00"!=a||!s||(l="replay"),i?l="replay":a!=e&&(l="recording");else if(g<d&&d<u){l="living",n=setTimeout(function(){$("#livestat"+r).text("录制中"),clearTimeout(n)},u-o)}else d<g&&(l="livenotstart");return new Handlebars.SafeString(l)}l=""}),Handlebars.registerHelper("canIClick",function(e,a,t,i){var s=!0,r=(new Date).getTime();if(e&&a){var n=r-(t?new Date((t||"").replace(/-/g,"/").replace(/[TZ]/g," ").slice(0,19)).getTime():r),l=new Date((e||"").replace(/-/g,"/").replace(/[TZ]/g," ").slice(0,19)).getTime();l+=n;var d=new Date((a||"").replace(/-/g,"/").replace(/[TZ]/g," ").slice(0,19)).getTime();((d+=n)<r&&!i||r<l)&&(s=!1)}else s=!1;return s})},resetLiveContboxSize:function(){var e=this;e.reCountLiveCntSize(),$(window).resize(function(){e.reCountLiveCntSize()})},reCountLiveCntSize:function(){var e=$(window).height(),a=$(window).width(),t=$("#strategyListBox").height();$("#menuScroll").height(e-76-t-6),$("#utMainBox").height(e-60),$(".right-box").height(e-76),$(".right-box").width(a-344),$("iframe").height(e-107-110),$(".iframe-box").height($(window).height()-107-110)},scrollContent:function(){this.srcollBar=function(){$.each(["#menuScroll"],function(e,a){0<$(a).length&&$(a).niceScroll({cursorcolor:"#666",cursoropacitymax:.5,touchbehavior:!1,cursorwidth:"6px",cursorborder:"0",cursorborderradius:"8px",autohidemode:!1})})},this.srcollBar()},showStatusFunc:function(e,a,t,i){var s="",r=(new Date).getTime();if(e&&a){var n=r-(t?new Date((t||"").replace(/-/g,"/").replace(/[TZ]/g," ").slice(0,19)).getTime():r),l=new Date((e||"").replace(/-/g,"/").replace(/[TZ]/g," ").slice(0,19)).getTime();l+=n;var d=new Date((a||"").replace(/-/g,"/").replace(/[TZ]/g," ").slice(0,19)).getTime();(d+=n)<r?i?s="":a!=e&&(s="录制中"):l<r&&r<d?s="直播中":r<l&&(s="直播未开始")}else s="";return s},getLikes:function(){var a=this,t=[],i=$("#videoDetailBox");$(".encourage-icon",i).each(function(){var e=$(this).parents(".atticle-show").attr("data-artid");t.push({uid:a.getUID,newsId:e,appId:a.getAppId})}),$.ajax({type:"get",timeout:ajaxTimeout,contentType:"text/plain",dataType:"jsonp",url:window.gConfig.likeListDataServerJsonp,data:{jsondata:JSON.stringify(t)},success:function(t,e){t.isSucess?$(".encourage-icon",i).each(function(){for(var e=$(this).parents(".atticle-show").attr("data-artid"),a=0;a<t.message.length;a++)t.message[a].newsId==e&&($("#encourageed_"+e).attr("data-encourageednum",t.message[a].likes),!0===t.message[a].liked?(0!=t.message[a].likes&&$("#encourageed_"+e).text(t.message[a].likes),$("#encourageicon_"+e).html("&#xe66b;"),$("#encourageed_"+e).addClass("liked")):!1===t.message[a].liked&&($("#encourageicon_"+e).html("&#xe61b;"),$("#encourageed_"+e).removeClass("liked")))}):$(".encourage",i).each(function(){$(this).find("span").text(0)})},error:function(e,a,t){$(".encourage",i).each(function(){$(this).find("span").text(0)})}})},clickLike:function(){var s=this,e=$("#videoDetailBox");$(".encourage-icon",e).on("click",function(){var e=$(this).parents(".atticle-show").attr("data-artid"),t=$("#encourageed_"+e),i=$("#encourageicon_"+e);if(!t.hasClass("liked")){t.addClass("liked");var a=t.attr("data-encourageednum");t.html(Number(a)+1),i.html("&#xe66b;"),$.ajax({url:window.gConfig.likeDataServerJsonp,type:"get",dataType:"jsonp",timeout:ajaxTimeout,cache:!1,data:{uid:s.getUID,newsId:e,appId:s.getAppId},success:function(e,a){"success"===a&&(e.isSucess?t.attr("data-channelid"):(i.html("&#xe61b;"),t.removeClass("liked")))},error:function(e,a,t){}})}})},bindEvents:function(){var l=this,r=$("#strategyMenuIcon"),n=$("#strategyList"),d=$("#riliBox"),o=($("#datePick"),$("#riliBox .jrzx")),c=$("#riliBox .other");$("#uthNav").on("click","li",function(){var e=$(this),a=e.find("span").text();l.curPageOpts.currpage=1,l.curPageOpts.isbg=!1,l.flag=!1,l.headerNavCurrentIndex=e.index();var t=e.attr("data-tagindex"),i=e.attr("data-date"),s=e.attr("data-tagid");$("#d11").val(i),e.hasClass("current")||($("#menuList").find("li").remove(),$("#menuList").find(".noContent").remove(),$("#uthNav li").removeClass("current"),$("#menuList .no-info").show(),e.addClass("current"),r.hasClass("up")&&("策略实战"!=a?n.removeClass("height-auto height-fixed"):n.addClass("height-auto height-fixed")),"解盘晨会"==a?(d.removeClass("today"),r.removeClass("open"),o.hide(),c.show(),l.getZbxxTagList("jpch"),$("#strategyList li").eq(t).addClass("current"),l.getusertrainingData(s,"jpch",i)):"指标学习"==a?(d.removeClass("today"),r.removeClass("open"),o.hide(),c.show(),l.getZbxxTagList("zbxx"),l.getusertrainingData(s,"zbxx",i)):"当日最新"==a?(d.addClass("today"),r.removeClass("open"),n.html(""),o.show(),c.hide(),$("#riqi").text("今日"),l.getLatestgData(e.attr("data-date"))):"策略实战"==a&&(d.removeClass("today"),r.addClass("open"),o.hide(),c.show(),l.getStrategyList(),0==t?l.get1monthData():l.getFirstTypeStrategyList(s,"",i)))}),$("#strategyList").on("click",".strategy-item",function(){var e=$(this);if(!e.hasClass("current")){var a=e.attr("data-classId"),t=e.index(),i=(e.attr("data-className"),$("#uthNav .current")),s=i.index(),r=e.attr("data-classid"),n=i.attr("data-date");i.attr("data-tagindex",t),i.attr("data-tagid",r),0!=t?i.attr("data-date",n):($("#d11").val(""),i.attr("data-date","")),l.curPageOpts.currpage=1,l.curPageOpts.isbg=!1,l.flag=!1,$("#menuList").find("li").remove(),$("#menuList").find(".noContent").remove(),$("#menuList .no-info").show(),$("#strategyList li").removeClass("current"),e.addClass("current"),i.attr("data-tagindex",t),1==s?0==t?l.getusertrainingData(r,"jpch",""):l.getusertrainingData(r,"jpch",n):3==s?0==t?l.getusertrainingData(r,"zbxx",""):l.getusertrainingData(r,"zbxx",n):2==s&&(0<t?l.getFirstTypeStrategyList(a,"",n):l.get1monthData(""))}}),$("#menuList").on("click",".menu-item",function(){var e=$(this);if(!e.hasClass("current")){var a=e.index()-1,t=$("#uthNav .current").find("span").text();e.hasClass("recording")||e.hasClass("livenotstart")||($("#menuList li").removeClass("current"),e.addClass("current"),"解盘晨会"==t?l.articleListShow(l.getTagData[a],"usertraining"):"指标学习"==t?l.articleListShow(l.getTagData[a],"usertraining"):"当日最新"==t?l.articleListShow(l.latestData[a],"usertraining"):"策略实战"==t&&($("#strategyList li").eq(0).hasClass("current")?l.articleListShow(l.getonemonthData[a]):l.articleListShow(l.getStrategyData[a])))}}),r.on("click",function(){var e=$(this);e.hasClass("down")?(e.html("&#xe6b1;"),e.addClass("up").removeClass("down"),n.addClass("height-auto"),7==$("#strategyList li").length&&n.addClass("height-fixed")):(e.html("&#xe6b2;"),e.addClass("down").removeClass("up"),n.removeClass("height-auto").removeClass("height-fixed")),l.reCountLiveCntSize(),$("#menuScroll").getNiceScroll(0).resize()})},getStrategyList:function(){var i=this,s=$("#strategyList"),r=$("#uthNav .current"),n=["{{#each this}}",'<li class="strategy-item" data-classId="{{ParentId}}" data-className="{{ShowName}}"  clickkey="clsz|celuelist" clickdata="{{ParentId}}|usertraining" clickremark="">{{ShowName}}</li>',"{{/each}}"].join("");$.ajax({url:window.gConfig.trainapiHost+"gettrainclientinfo",type:"GET",timeout:ajaxTimeout,dataType:"JSON",cache:!1,success:function(e){if(e){var a=e;i.strategyList=a;var t=Handlebars.compile(n)(a);s.html(" "),s.append('<li class="strategy-item" data-classId="0" clickkey="clsz|celuelist" clickdata="zonghe|usertraining" clickremark="" >综合</li>'),s.append(t),$("#strategyList li").eq(r.attr("data-tagindex")).addClass("current")}},error:function(e,a,t){layer.msg("该分类的暂无实战培训")}})},getZbxxTagList:function(i){var s=this,r=$("#strategyList"),n=$("#uthNav .current"),l=["{{#each this}}",'<li class="strategy-item" data-classId="{{TrainTagID}}" data-className="{{ShowName}}"  clickkey="clsz|celuelist" clickdata="{{TrainTagID}}|usertraining" clickremark="">{{ShowName}}</li>',"{{/each}}"].join("");$.ajax({url:window.gConfig.trainapiHost+"gettraintaginfo?type="+i,type:"GET",timeout:ajaxTimeout,dataType:"JSON",cache:!1,success:function(e){if(e){var a=e;s.strategyTagList=a;var t=Handlebars.compile(l)(a);r.html(" "),r.append('<li class="strategy-item" data-classId="0" clickkey="'+i+'|celuelist" clickdata="zonghe|usertraining" clickremark="" >综合</li>'),r.append(t),$("#strategyList li").eq(n.attr("data-tagindex")).addClass("current")}},error:function(e,a,t){layer.msg("该分类的暂无实战培训")}})},getFirstTypeStrategyList:function(e,a,t){var s=this,r=["{{#each Message}}",'<li class="menu-item {{actralTraingStatusClass Live_StartTime Live_EndTime ../SystemTime LiveVideoURL VideoPlayURL ID}} {{firstItem @index}} {{first @index}}" data-newsid={{ID}}  clickkey="'+e+'|leftlist" clickdata="{{ID}}|usertraining" clickremark="">','<div class="tag three-word">{{Froms}}</div>','<div class="menu-content">','<div class="title">{{substr30 Title}}</div>','<div class="autor-time">','<span class="autor">{{From}}</span>','<span class="time">{{renderTimeStrategy Live_StartTime CreateTime}}</span>','<span class="live-status" id="livestat{{ID}}">{{showstatus Live_StartTime Live_EndTime ../SystemTime LiveVideoURL}}</span>',"</div>","</div>","</li>","{{/each}}"].join("");$.ajax({url:window.gConfig.trainapiHost+"getstrategynewslist_multimedia",type:"GET",timeout:ajaxTimeout,dataType:"JSON",cache:!1,data:{ClientStrategyId:e,date:t},success:function(e){0==e.RetCode&&e.Message&&0<e.Message.length?($("#menuList .no-info").hide(),function(e){for(var a=0;a<s.strategyList.length;a++)for(var t=0;t<e.Message.length;t++)s.getStrategyData.push(e.Message[t]),e.Message[t].From==s.strategyList[a].ParentName&&(e.Message[t].Froms=s.strategyList[a].ShowName);1==s.curPageOpts.currpage&&(s.getonemonthData=[],$("#menuList li").remove(),$("#menuList .no-info").hide());var i=Handlebars.compile(r)(e);$("#menuList").append(i),s.firstScreemLiveStatus(e,"actralTraining")}(e),setTimeout(function(){utils.InitNiceScroll("#menuScroll")},0)):($("#menuList li").remove(),$("#menuList .no-info").show().text("暂无相关培训内容..."),$("#videoDetailBox").html('<div class="noContent">暂无相关培训内容...</div>'))},error:function(e,a,t){$("#menuList li").remove(),$("#menuList .no-info").show().text("暂无相关培训内容..."),$("#videoDetailBox").html('<div class="noContent">暂无相关培训内容...</div>')}})},articleListShow:function(e,a,t){var i=this,s=$("#videoDetailBox"),r="";if(!t){Handlebars.registerHelper("getstats",function(e){var a="live";return e&&(a="vod"),a}),Handlebars.registerHelper("actraltraininggetstat",function(e){var a="vod";return e&&(a="live"),a}),i.detailCode="usertraining"==a?['<div class=\'atticle-show section yaowen-article\' id="atticleShow{{ID}}" data-artid={{ID}} data-type="{{decompose TagInfo}}" >',"        <div class='a-iner'>","            <div class='art-title'>",'                <h3 id="artcleTitle" title="">{{Mtg_name}}</h3>',"            </div>",'            <div class="art-subtitle">',"                <div class=' pull-right'>",'                  <div class="encourage"><span id="encourageed_{{ID}}" data-encourageedNum="0">喜欢就鼓励下！</span><i class="icon-1 encourage-icon" id="encourageicon_{{ID}}" data-channelid={{ID}}>&#xe61b;</i></div>  ',"                </div>",'                <div class="pull-left ac-submain">','                    <span id="articleComefrom">{{decompose TagInfo}}</span>',"                    <span>{{Txtteachar}}</span>","                </div>",'                <div class="clearfloat"></div>',"            </div>",'            <div class="art-cnt">',"{{#if hasNoPlayLive}} ",'                   <div class="videoPack" data-VideoUrl = "{{Video_url}}">','                       <div class="vedio cover-img" data-VideoUrl = "{{ Video_url }}" data-midtype="vod" data-newsid = "{{ID}}"><img calss="cover" src="'+hasNoLive+'" /</div>','                       <div class="vedio iframe-box"><iframe src="" data-VideoUrl = "{{Video_url}}" frameborder="0"></iframe></div>',"                   </div> ","{{else}} ",'                <div class="ac-iner">',"                   {{#if Video_url}} ",'                   <div class="videoPack" data-VideoUrl = "{{Video_url}}">',"                       {{#if CoverImg}}",'                       <div class="vedio cover-img" data-VideoUrl = "{{ Video_url }}" data-midtype="vod" data-newsid = "{{ID}}"><img class="cover" src="{{CoverImg}}" /><img class="liveIcon" src="'+liveImg+'" /></div>{{else}}','                       <div class="vedio cover-img" data-VideoUrl = "{{ Video_url }}" data-midtype="vod" data-newsid = "{{ID}}"><img calss="cover" src="'+imgDefault+'" /><img class="liveIcon" src="'+liveImg+'" /></div>{{/if}}','                       <div class="vedio iframe-box"><iframe src="" data-VideoUrl = "{{Video_url}}" frameborder="0"></iframe></div>',"                   </div> ","                  {{else}}",' <div class="videoPack livepack" data-VideoUrl = "{{Gensee_URL}}">',"                       {{#if CoverImg}}",'                       <div class="vedio cover-img" data-VideoUrl = "{{ Gensee_URL }}"  data-newsid = "{{ID}}" data-vodurl="{{Gensee_URL}}" data-midtype="{{actraltraininggetstat Gensee_URL }}"><div style="height:100%;"><img class="cover" src="{{CoverImg}}" /><img class="liveIcon" src="'+liveImg+'" /></div></div>{{else}}','                       <div class="vedio cover-img" data-VideoUrl = "{{ Gensee_URL }}"  data-newsid = "{{ID}}" data-vodurl="{{Gensee_URL}}" data-midtype="{{actraltraininggetstat Gensee_URL }}"><div style="height:100%;"><img calss="cover" src="'+imgDefault+'" /><img class="liveIcon" src="'+liveImg+'" /></div></div>{{/if}}','                       <div class="vedio iframe-box"><iframe src="" frameborder="0"></iframe></div>',"                   </div> ","                  {{/if}} ","                </div>","                  {{/if}} ","            </div>","        </div>","    </div>"].join(""):['<div class=\'atticle-show section yaowen-article\' id="atticleShow{{ID}}" data-artid={{ID}} data-type="{{decompose TagInfo}}" >',"        <div class='a-iner'>","            <div class='art-title'>",'                <h3 id="artcleTitle" title="">{{Title}}</h3>',"            </div>",'            <div class="art-subtitle">',"                <div class=' pull-right'>",'                  <div class="encourage"><span id="encourageed_{{ID}}" data-encourageedNum="0">喜欢就鼓励下！</span><i class="icon-1 encourage-icon" id="encourageicon_{{ID}}" data-channelid={{ID}}>&#xe61b;</i></div>  ',"                </div>",'                <div class="pull-left ac-submain">',"                    <span>{{From}}</span>","                </div>",'                <div class="clearfloat"></div>',"            </div>",'            <div class="art-cnt">',"{{#if hasNoPlayLive}} ",'                   <div class="videoPack" data-VideoUrl = "{{VideoPlayURL}}">','                       <div class="vedio cover-img" data-VideoUrl = "{{ VideoPlayURL }}" data-midtype="vod" data-newsid = "{{ID}}"><img calss="cover" src="'+hasNoLive+'" /></div>','                       <div class="vedio iframe-box"><iframe src="" data-VideoUrl = "{{VideoPlayURL}}" frameborder="0"></iframe></div>',"                   </div> ","{{else}} ",'                <div class="ac-iner">',"                   {{#if VideoPlayURL}} ",'                   <div class="videoPack" data-VideoUrl = "{{VideoPlayURL}}">',"                       {{#if CoverImg}}",'                       <div class="vedio cover-img" data-VideoUrl = "{{ VideoPlayURL }}" data-midtype="vod" data-newsid = "{{ID}}"><img class="cover" src="{{CoverImg}}" /><img class="liveIcon" src="'+liveImg+'" /></div>{{else}}','                       <div class="vedio cover-img" data-VideoUrl = "{{ VideoPlayURL }}" data-midtype="vod" data-newsid = "{{ID}}"><img calss="cover" src="'+imgDefault+'" /><img class="liveIcon" src="'+liveImg+'" /></div>{{/if}}','                       <div class="vedio iframe-box"><iframe src="" data-VideoUrl = "{{VideoPlayURL}}" frameborder="0"></iframe></div>',"                   </div> ","                  {{else}}",'                   <div class="videoPack livepack" data-VideoUrl = "{{LiveURL}}">',"                       {{#if CoverImg}}",'                       <div class="vedio cover-img {{canIClick Live_StartTime Live_EndTime ../SystemTime LiveVideoURL}}" data-VideoUrl = "{{ LiveURL }}"  data-newsid = "{{ID}}" data-vodurl="{{LiveVideoURL}}" data-midtype="{{getstats LiveVideoURL }}"><div style="height:100%;"><img class="cover" src="{{CoverImg}}" /><img class="liveIcon" src="'+liveImg+'" /></div></div>{{else}}','                       <div class="vedio cover-img {{canIClick Live_StartTime Live_EndTime ../SystemTime LiveVideoURL}}" data-VideoUrl = "{{ LiveURL }}"  data-newsid = "{{ID}}" data-vodurl="{{LiveVideoURL}}" data-midtype="{{getstats LiveVideoURL }}"><div style="height:100%;"><img calss="cover" src="'+imgDefault+'" /><img class="liveIcon" src="'+liveImg+'" /></div></div>{{/if}}','                       <div class="vedio iframe-box"><iframe src="" frameborder="0"></iframe></div>',"                   </div> ","                  {{/if}} ","                </div>"," {{/if}} ","            </div>","        </div>","    </div>"].join("");!function(e){if(e){var a=Handlebars.compile(i.detailCode);r=a(e),s.html(r),i.hasNoPlayLive||i.changeUrl($("#videoDetailBox .cover-img")),$("iframe").height($(window).height()-107-110),$(".iframe-box").height($(window).height()-107-110),i.getLikes(),i.clickLike()}else s.html('<div class="noContent">暂无内容...</div>')}(e)}},changeUrl:function(e){if(!e.hasClass("false")){e.parents("#videoDetailBox").find(".iframe-box").removeClass("open").find("iframe").attr("src",""),e.parents("#videoDetailBox").find(".cover-img").removeClass("close");var a=e.siblings(".iframe-box").addClass("open").find("iframe"),t=e.attr("data-vodurl");t||(t=e.attr("data-VideoUrl"));var i=t.substr(t.indexOf("-")+1,t.length),s=videoMediaSDK+"?midtype="+e.attr("data-midtype")+"&ownerid="+i+"&uid="+this.getUID+"&uname=&authcode=888888";a.attr("src",s),e.addClass("close")}},getusertrainingData:function(e,a,t){var i=this,s=["{{#each Message}}",'<li class="menu-item {{statusClass Class_date EndDate ../SystemTime Video_url}} {{firstItem @index}} {{first @index}}"  data-newsid={{ID}}  clickkey="tag'+e+'|leftlist" clickdata="{{ID}}|usertraining" clickremark="">','<div class="tag {{tagNameLength TrainTagName}}">{{TrainTagName}}</div>','<div class="menu-content">','<div class="title">{{substr30 Mtg_name}}</div>','<div class="autor-time">','<span class="autor">{{Txtteachar}}</span>','<span class="time">{{renderTime Class_date}}</span>','<span class="live-status" id="livestat{{ID}}">{{showstatus Class_date EndDate ../SystemTime Video_url}}</span>',"</div>","</div>","</li>","{{/each}}"].join("");$.ajax({url:window.gConfig.trainapiHost+"gettrainlistbydateandtag",type:"GET",timeout:ajaxTimeout,dataType:"JSON",cache:!1,data:{currpage:i.curPageOpts.currpage,pageSize:i.curPageOpts.pageSize,trainTag:e,trainType:a,date:t},success:function(e){if(0==e.RetCode&&e.Message&&0<e.Message.length){1==i.curPageOpts.currpage&&(i.getTagData=[],$("#menuList li").remove(),$("#menuList .no-info").hide(),i.curPageOpts.totalpage=Math.ceil(e.TotalCount/i.curPageOpts.pageSize));for(var a=0;a<e.Message.length;a++)i.getTagData.push(e.Message[a]);var t=Handlebars.compile(s)(e);$("#menuList").append(t),setTimeout(function(){utils.InitNiceScroll("#menuScroll")},0),i.firstScreemLiveStatus(e,"usertraining"),i.dataLength=e.Message.length,e.Message.length<=i.curPageOpts.pageSize&&i.curPageOpts.totalpage==i.curPageOpts.currpage&&(i.curPageOpts.isbg=!0)}else $("#menuList li").remove(),$("#menuList .no-info").show().text("暂无相关培训内容..."),$("#videoDetailBox").html('<div class="noContent">暂无相关培训内容...</div>')},error:function(e,a,t){$("#menuList .no-info").hide(),$("#menuList").html('<div class="noContent">暂无相关培训内容...</div>'),$("#videoDetailBox").html('<div class="noContent">暂无相关培训内容...</div>')}})},getLatestgData:function(e){var i=this,s=["{{#each Message}}",'<li class="menu-item {{statusClass Class_date EndDate ../SystemTime Video_url ID}} {{firstItem @index}} {{first @index}}" data-vodurl="" data-newsid={{ID}} clickkey="drzx|leftlist" clickdata="{{ID}}|usertraining" clickremark=""><div class="menu-item-wrapper">','<div class="tag {{tagNameLength TrainTagName}}">{{TrainTagName}}</div>','<div class="menu-content">','<div class="title">{{substr30 Mtg_name}}</div>','<div class="autor-time">','<span class="autor">{{Txtteachar}}</span>','<span class="time">{{renderTime Class_date}}</span>','<span class="live-status" id="livestat{{ID}}">{{showstatus Class_date EndDate ../SystemTime Video_url}}</span>',"</div>","</div>","</div></li>","{{/each}}"].join("");$.ajax({url:window.gConfig.trainapiHost+"gettrainlistbydate",type:"GET",timeout:ajaxTimeout,dataType:"JSON",cache:!1,data:{date:e},success:function(e){if(0==e.RetCode&&e.Message&&0<e.Message.length){i.latestData=[],$("#menuList li").remove(),$("#menuList .no-info").hide();for(var a=0;a<e.Message.length;a++)i.latestData.push(e.Message[a]);var t=Handlebars.compile(s)(e);$("#menuList").append(t),setTimeout(function(){utils.InitNiceScroll("#menuScroll")},0),i.firstScreemLiveStatus(e,"usertraining")}else $("#menuList li").remove(),$("#menuList .no-info").show().text("暂无相关培训内容..."),$("#videoDetailBox").html('<div class="noContent">暂无相关培训内容...</div>')},error:function(e,a,t){$("#menuList li").remove(),$("#menuList .no-info").show().text("暂无相关培训内容..."),$("#videoDetailBox").html('<div class="noContent">暂无相关培训内容...</div>')}})},get1monthData:function(e){var s=this,r=["{{#each Message}}",'<li class="menu-item {{actralTraingStatusClass Live_StartTime Live_EndTime ../SystemTime LiveVideoURL VideoPlayURL ID}} {{firstItem @index}} {{first @index}}" data-newsid={{ID}} clickkey="zonghe|leftlist" clickdata="{{ID}}|usertraining" clickremark="">','<div class="tag three-word">{{Froms}}</div>','<div class="menu-content">','<div class="title">{{substr30 Title}}</div>','<div class="autor-time">','<span class="autor">{{From}}</span>','<span class="time">{{renderTimeStrategy Live_StartTime CreateTime}}</span>','<span class="live-status" id="livestat{{ID}}">{{showstatus Live_StartTime Live_EndTime ../SystemTime LiveVideoURL}}</span>',"</div>","</div>","</li>","{{/each}}"].join("");$.ajax({url:window.gConfig.trainapiHost+"getstrategynewslist_multimedia1month",type:"GET",timeout:ajaxTimeout,dataType:"JSON",cache:!1,data:{currpage:s.curPageOpts.currpage,pageSize:s.curPageOpts.pageSize,date:e},success:function(e){if(1==s.curPageOpts.currpage&&(s.getonemonthData=[],$("#menuList li").remove(),$("#menuList .no-info").hide(),s.curPageOpts.totalpage=Math.ceil(e.TotalCount/s.curPageOpts.pageSize)),0==e.RetCode&&e.Message&&0<e.Message.length){for(var a=0;a<s.strategyList.length;a++)for(var t=0;t<e.Message.length;t++)e.Message[t].From==s.strategyList[a].ParentName&&(e.Message[t].Froms=s.strategyList[a].ShowName);for(a=0;a<e.Message.length;a++)s.getonemonthData.push(e.Message[a]);var i=Handlebars.compile(r)(e);$("#menuList").append(i),setTimeout(function(){utils.InitNiceScroll("#menuScroll")},0),s.firstScreemLiveStatus(e,"actralTraining"),s.dataLength=e.Message.length,e.Message.length<=s.curPageOpts.pageSize&&s.curPageOpts.totalpage==s.curPageOpts.currpage&&(s.curPageOpts.isbg=!0)}else $("#menuList li").remove(),$("#menuList .no-info").show().text("暂无相关培训内容..."),$("#videoDetailBox").html('<div class="noContent">暂无相关培训内容...</div>')},error:function(e,a,t){$("#menuList li").remove(),$("#menuList .no-info").show().text("暂无相关培训内容..."),$("#videoDetailBox").html('<div class="noContent">暂无相关培训内容...</div>')}})},firstScreemLiveStatus:function(e,a){var t=this;if(1==t.curPageOpts.currpage){var i=$("#menuList .living").length,s=$("#menuList .replay").length;if(i){t.hasNoPlayLive=!1;for(var r=0,n=e.Message.length;r<n;r++)e.Message[r].ID==$("#menuList .living").eq(0).attr("data-newsid")&&("usertraining"==a?t.articleListShow(e.Message[r],"usertraining"):t.articleListShow(e.Message[r],"actralTraining"),$("#menuList .living").eq(0).addClass("current").siblings().removeClass("current"))}else if(s){t.hasNoPlayLive=!1;for(r=0,n=e.Message.length;r<n;r++)e.Message[r].ID==$("#menuList .replay").eq(0).attr("data-newsid")&&("usertraining"==a?t.articleListShow(e.Message[r],"usertraining"):t.articleListShow(e.Message[r],"actralTraining"),$("#menuList .replay").eq(0).addClass("current").siblings().removeClass("current"))}else t.hasNoPlayLive=!0,$.extend(e.Message[0],{hasNoPlayLive:!0}),"usertraining"==a?t.articleListShow(e.Message[0],"usertraining"):t.articleListShow(e.Message[0],"actralTraining")}},selectDatePick:function(){var e=$("#riqi"),a=$("#d11").val(),t=$("#uthNav .current"),i=t.index(),s=t.attr("data-tagindex"),r=$("#strategyList .current"),n=r.attr("data-classid");r.attr("data-classname");0!=r.index()&&0!=i&&t.attr("data-date",a),l.curPageOpts.currpage=1,l.curPageOpts.isbg=!1,l.flag=!1,0==i?(l.getLatestgData(a),e.text(a)):1==i?l.getusertrainingData(n,"jpch",a):3==i?l.getusertrainingData(n,"zbxx",a):2==i&&(0==s?l.get1monthData(a):l.getFirstTypeStrategyList(n,"",a))}};l.init(),window.selectDatePick=l.selectDatePick,window.pickedFunc=function(){window.selectDatePick()},$.emoneyAanalytics().Init(tjAppid.usertraining,"userTraining","")}();