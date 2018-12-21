/* ========================================================================
 * app.js
 * ======================================================================== */

(function($, window, layer) {
    // placeholder
    $(function() {
        $("input, textarea").placeholder({ customClass: "pad-placeholder" });
    });

    // $('#topNav .navmenu').on('click','li',function(){
    //   var $this = $(this);
    //   var target = $this.data('item');
    //   $this.addClass('active').siblings().removeClass('active');
    //   $('#MainCont .pageItem').hide();
    //   $('#MainCont .'+target).show();
    // });

    // 活动切换
    (function() {
        var tabsPack = $(".tabs-pack");
        var tabsNav = $(">.tabs-nav", tabsPack);
        var tabsCnt = $(">.tabs-nav + .tabs-cnt", tabsPack);

        tabsPack.on("click", ">.tabs-nav li", function() {
            var tabsNav = $(this)
                .parent()
                .parent();
            var tabsCnt = $("+.tabs-cnt > .tabsitem", tabsNav);

            var idx = $("li", tabsNav)
                .removeClass("active")
                .index(this);
            $(this).addClass("active");

            // $(this).addClass('active').siblings().removeClass('active');

            // console.log(pparent);

            tabsCnt
                .removeClass("active")
                .eq(idx)
                .addClass("active");
        });
    })();

    //处理滚动高度
    $.extend({
        setHomeListScroll: function() {
            var $actList = $("#actList");
            if ($(".topstar", $actList)) {
                $(".list-scrollbox", $actList).height(
                    182 - $(".topstar", $actList).height()
                );
            } else {
                $(".list-scrollbox", $actList).height(182);
            }
        },

        // 选择
        setUserStarsShow: function(t, g, s) {
            var opt = { tlv: 0, glv: 0, slv: 0 };
            var param = { tlv: t, glv: g, slv: s };
            opt = $.extend(opt, param);
            var classStr = "starlv";
            function lvnumber(num) {
                var cnum = parseInt(num);
                if ($.isNumeric(cnum)) {
                    return cnum < 0 ? 0 : cnum > 10 ? 10 : cnum;
                } else {
                    return 0;
                }
            }

            param = { tlv: lvnumber(t), glv: lvnumber(g), slv: lvnumber(s) };

            $.each(param, function(key, val) {
                classStr += " " + key + "" + val;
            });

            $("#starlv").attr("class", classStr);
        },

        // 编辑用户名
        editUserName: function(cbA, cbB) {
            $("#modyUsrname").click(function() {
                var $editBtn = $("a", this),
                    $editbox = $("#editUsrname"),
                    $editInp = $("input", $editbox),
                    $editSpan = $("b", $editbox);

                if ($editBtn.text() === "[编辑]") {
                    $("#editUsrname").addClass("editing");
                    $editInp.val($editSpan.text());
                    $editBtn.text("[保存]");

                    if ($.isFunction(cbB)) {
                        cbA();
                    }
                } else {
                    $("#editUsrname").removeClass("editing");
                    $editBtn.text("[编辑]");
                    $editSpan.text($editInp.val()).attr("title", $editInp.val());
                    if ($.isFunction(cbB)) {
                        cbB();
                    }
                }
            });
        },

        userDecoration: function(udoptions, uid) {
            var opts = {}; //
            var decorationName = {
                wybb: "文娱标兵",
                hwbs: "合格投资者",
                hddr: "互动达人",
                gpmz: "高朋满座",
                seed: "种子用户"
            };
            var returnCls = "";
            var usrDcShow = $("#userDecorationShow");
            var allDecors = "";
            var oldCookies = null;
            var cookieKey = "UYhIKq9uZDecoration" + uid;
            var oldDecoration = null;

            $.extend(opts, udoptions);
            var curclas = usrDcShow.attr("class");
            returnCls = curclas;

            function clearThisDecorations(keyw) {
                var regss = new RegExp("\\s" + keyw + "-ud\\d", "g");
                if(returnCls!=undefined){
                    returnCls = returnCls.replace(regss, "");
                }
                return returnCls;
            }

            function chkNum(str) {
                var numm = parseInt(str);
                numm = isNaN(numm) ? 0 : numm;
                return numm;
            }

            // 是否存在cookies
            oldCookies = $.cookie(cookieKey);
            $.cookie(cookieKey, JSON.stringify(opts), { expires: 365 });
            if (!!oldCookies) {
                //存在cookies

                oldDecoration = JSON.parse(oldCookies);

                $.each(opts, function(key, itm) {
                    if (
                        !!oldDecoration.hasOwnProperty(key) &&
                        chkNum(itm) > chkNum(oldDecoration[key])
                    ) {
                        allDecors += decorationName[key] + "、";
                    }
                });

                // 显示升级
                if (!!allDecors) {
                    layer.msg("你获得了新的" + allDecors.slice(0, -1) + "勋章！", {
                        icon: 1,
                        anim: 1,
                        tipsMore: true
                    });
                }

                // 过滤提示内容
            } else {
            }

            // 修改样式
            $.each(opts, function(key, itm) {
                clearThisDecorations(key);
                returnCls += " " + key + "-ud" + itm;
            });

            usrDcShow.attr("class", returnCls);
        },

        // 展示活动
        showActities: function(actUrl) {
            var showActFrame = $("#showActFrame");
            var actFrame = $("iframe", showActFrame);
            actFrame.attr("src", actUrl);
            showActFrame.show();
        }
    });
})(jQuery, window, layer);

$.setHomeListScroll();
