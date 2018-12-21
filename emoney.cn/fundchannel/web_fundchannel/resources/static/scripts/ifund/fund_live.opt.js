(function () {
    var page = {
        nonedataHtml: '<div class="no-info">暂无课程安排，课程回放也很精彩哦</div>',
        clickFalg: false,
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
        init: function () {
            page.getUID = $.cookie('expertnews.uid') || utils.GetQueryString('uuid') || 0;
            $.support.cors = true;
            $('textarea').placeholder();
            utils.EMSSO();
            utils.appendSSO();
            page.bindaskAnswerEnvent();
            // page.scrollBar();
            $("#inpFacemark").qqFace({
                assign: "id_expertquestion", //给输入框赋值
                path: "http://static.emoney.cn/live/Content/images/arclist/" //表情图片存放的路径
            });
        },
        bindaskAnswerEnvent: function () {
            var $tabBoxItem = $('#tabBox .question');
            var $tabBoxContentItem = $('#interactionContentBox .interaction-content-list');
            $tabBoxItem.on('click', function () {
                var $this = $(this);
                var index = $this.index();
                $tabBoxItem.removeClass('current');
                $this.addClass('current');
                $tabBoxContentItem.removeClass('open');
                $tabBoxContentItem.eq(index).addClass('open');
                if (index == 1 && !page.clickFalg) {
                    page.clickFalg = true;
                    $Live.GetLiveMyQuestion();
                }
            });
        },

        scrollBar: function () {
            $.each(
                [
                    ".live-content-box",
                    ".interaction-content-inner",
                    "html"
                ],
                function (idx, item) {
                    if ($(item).length > 0) {
                        $(item).niceScroll({
                            cursorcolor: "#afadad",
                            cursoropacitymax: 0.5,
                            touchbehavior: false,
                            cursorwidth: "6px",
                            cursorborder: "0",
                            cursorborderradius: "8px",
                            autohidemode: false
                        });
                    }
                }
            );
        },
        pickedFuncRight: function () {

        },
        pickedFuncLeft: function (dp) {
            var $d11Left = $('#d11Left');
            var $date = $('#date');
            var today = utils.getNowDate();
            var selectedDate = $d11Left.val();
            if (selectedDate == today) {
                $date.text("今日直播");
            } else {
                $date.text(selectedDate);
            }
            var oldDate = $date.attr("data-date");
            if (selectedDate === oldDate) {
                return;
            }
            $date.attr("data-date", selectedDate);
            $Live.BindLiveAllContent(selectedDate, 1);
            page.clickFalg = true;
            $Live.GetLiveMyQuestion(selectedDate);
        }

    };
    page.init();
    window.pickedFuncRight = page.pickedFuncRight;
    window.pickedFuncLeft = page.pickedFuncLeft;
})();

