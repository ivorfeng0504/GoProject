
(function(){
    var page = {
        nonedataHtml: '<div class="no-info">暂无课程安排，课程回放也很精彩哦</div>',
        curPageOpts: {
            NewsID: 0,
            ColumnID: 16,
            StrategyID: 0,
            currpage: 1,
            pageSize: 8,
            LiveID: 0,
            ExpertLiveroom: "http://yqq.emoney.cn/Live/UserLive?lid=93",
            isbg: false,
            newsidzan: [],
            zanparams: []
        },
        init: function() {
            page.getUID = $.cookie('expertnews.uid') || utils.GetQueryString('uuid') || 0;
            $.support.cors = true;
            utils.EMSSO();
            utils.appendSSO();
            // page.scrollBar();
            page.bindLeftNavEnvent();
            page.liveScroll();
        },
        scrollBar:function(){
            $('html').niceScroll({
                cursorcolor: "#9c9797",
                cursoropacitymax: 0.5,
                touchbehavior: false,
                cursorwidth: "6px",
                cursorborder: "0",
                cursorborderradius: "8px",
                autohidemode: false
            });
        },
        bindLeftNavEnvent:function(){
            var $styleItem = $('#styleList li');
            var $styleBgItem = $('#styleBgBox .style-bg');
            $styleItem.on('click',function(){
                var $this = $(this);
                if($this.hasClass('first')) {
                    return;
                }
                var index = $this.index();
                $styleItem.removeClass('current');
                $this.addClass('current');
                $styleBgItem.hide().eq(index).show();
            })
        },
         //直播滚动
         liveScroll: function () {
            $.fn.extend({
                newScroll: function () {
                    var obj = $(this);
                    var lifirstHeight = $(obj).find("li:first").height();
                    setInterval(function () {
                        $(obj).find("li:first").animate({ "marginTop": lifirstHeight }, 500, function () {
                            $(this).css({ "marginTop": "0" });
                            $(obj).prepend($(obj).find("li:last"));
                        });
                    }, 8000);
                }
            });
            $('#liveContentList').newScroll();
            // $('.interaction-content-list').newScroll();
        }
       
    };
    page.init();
    // $.emoneyAanalytics().Init(tjAppid.usertraining, 'ifund','');
})();

