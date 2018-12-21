/* ========================================================================
 * app.js
 * ======================================================================== */
if (!window.console) {
    window.console = {};
}
if (!window.console.log) {
    window.console.log = function(msg) {};
}
// window.onerror = killErrors; function killErrors() {  return true;}

window.gConfig = {
    env: "serv",
    // apiHost: "http://dsclient.emoney.cn:8080/userhome",
    apiHost: www,
    staticPath: StaticServerHost
};

(function($, window, layer) {
    // placeholder
    $(function() {
        $("input, textarea").placeholder({
            customClass: "pad-placeholder"
        });
    });
    $.extend({
        GetQueryString: function(name) {
            var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)");
            var r = window.location.search.substr(1).match(reg);
            if (r != null) return unescape(r[2]);
            return null;
        }
    });

    var page = {
        proListTable: null,
        refundOrderId: null,
        refundMode: null,
        init: function() {
            var _self = this;
            _self.myRefundList();

            $(".btn-gobacktoorderlist").click(function() {
                $("#userProListshowbox").show();
                $("#userRefundshowbox").hide();
                $(".btn-goback").show();
                $(".btn-gobacktoorderlist").hide();
                _self.refundOrderId = null;
                _self.refundMode = null;
            });
        },

        myRefundList: function() {
            var _self = this;
            _self.refundOrderId = null;
            _self.refundMode = null;
            var OrderID = $.GetQueryString("orderid");
            if (OrderID == undefined || OrderID == "" || OrderID == null) {
                layer.msg("订单号有误！");
                window.history.back();
            }
            var htmlCodes = [
                '<div class="upro-tablepack normal-tablepack"><div id="OrderId">',
                '<table class="table form-table" >',
                '    <caption class="">',
                '        <div class="row">',
                '            <div class="col-xs-6">',
                '                <span class="OrderId">订单号：{{ Message.OrderId }}</span>',
                "            </div>",
                '            <div class="col-xs-6 text-right">',
                "            </div>",
                "        </div>",
                "    </caption>",
                "    <colgroup>",
                '        <col width="10%"></col>',
                '        <col width="35%"></col>',
                '        <col width="10%"></col>',
                '        <col width="8%"></col>',
                '        <col width="20%"></col>',
                '        <col width=""></col>',
                "    </colgroup>",
                "    <thead>",
                "        <tr>",
                "            <th>序号</th>",
                "            <th>商品</th>",
                "            <th>单价</th>",
                "            <th>数量</th>",
                "            <th>退货日期</th>",
                "        </tr>",
                "    </thead>",
                "    <tbody>",
                "{{#each Message.ProductList }}",
                "        <tr>",
                "            <td>{{xuhao @index }}</td>",
                "            <td> {{ ProductName }} </td>",
                "            <td> {{ priceFixed2 Price }} </td>",
                "            <td> {{ Count }} </td>",
                "            <td> {{timeformat1 RefundTime }}</td>",
                "        </tr>",
                "{{/each}}",
                "    </tbody>",
                "</table>",
                "</div>",
                "</div>"
            ].join("");
            //timeformat1  ../RefundTime

            var htmlCodes2 = [
                '<div class="upro-tablepack normal-tablepack"><div id="OrderId">',
                '<table class="table form-table" >',
                '    <caption class="">',
                '        <div class="row">',
                '            <div class="col-xs-6">',
                '                <span class="OrderId text-primary">您的退款状态：</span>',
                "            </div>",
                '            <div class="col-xs-6 text-right">',
                "            </div>",
                "        </div>",
                "    </caption>",
                "    <colgroup>",
                '        <col width="20%"></col>',
                '        <col width="35%"></col>',
                '        <col width="20%"></col>',
                '        <col width=""></col>',
                "    </colgroup>",
                "    <thead>",
                "        <tr>",
                "            <th>名称</th>",
                "            <th>退款方式</th>",
                "            <th>退款金额</th>",
                '            <th class="color-red">退款状态</th>',
                "        </tr>",
                "    </thead>",
                "    <tbody>",
                "{{#each Message.RefundDetailList}}",
                "        <tr>",
                "            <td> {{ ReturnModeName }} </td>",
                "            <td> {{ RefundTypeDesc }} </td>",
                "            <td> {{ RealbackPrice }} </td>",
                "            <td> {{ RefundStatusDesc }}</td>",
                "        </tr>",
                "{{/each}}",
                "    </tbody>",
                "</table>",
                '</div><div class="usr-opt-line text-right">',
                "</div>",
                "</div>"
            ].join("");
            Handlebars.registerHelper("priceFixed2", function(price) {
                var timestr3 = "";
                if (price == undefined || price == null || price.length == 0) {
                    timestr3 = "";
                } else {
                    timestr3 = parseFloat("" + price).toFixed(2);
                }
                return timestr3;
            });
            Handlebars.registerHelper("timeformat1", function(atime) {
                var timestr3 = "";
                if (atime == undefined || atime == null || atime.length == 0) {
                    timestr3 = "-";
                } else {
                    timestr3 = atime.slice(0, 16);
                }
                return timestr3;
            });

            Handlebars.registerHelper("xuhao", function(idx) {
                var timestr3 = "";

                if (idx == undefined || idx == null || idx.length == 0) {
                    timestr3 = "-";
                } else {
                    timestr3 = idx + 1;
                }
                return timestr3;
            });

            var template = Handlebars.compile(htmlCodes);
            var refundtemplate = Handlebars.compile(htmlCodes2);
            $.ajax({
                url: window.gConfig.apiHost + "mine/getorderinfo",
                type: "POST",
                timeout: 5000,
                dataType: "JSON",
                data: {
                    OrderId: OrderID
                },
                beforeSend: function() {},
                success: function(data) {

                    if (data.RetCode == 0 && data.Message != null) {

                            if(data.Message.RefundDetailList==null||data.Message.RefundDetailList==undefined||data.Message.RefundDetailList.length==0){

                                $("#usersProductsList").html("<div class='panel'><div class='panel-content' style='padding:1em;font-size:13px;line-height:2;'>　　尊敬的客户，由于您选择的退货途径不支持在此页面查看退货详情，请您谅解，并知悉，当您的退货申请被成功受理之后，会开始为您安排退款事宜。关于退款，将根据不同的支付方式和退款期限为您安排，请耐心等待您的退款。如需了解退货退款详情，请联系您的服务专员或者拨打服务热线10108688进行咨询。</div></div>");
                                $("#MainCont .usr-anounce .text-primary").html("　");
                                
                            }else{
                                $.extend(data, {});

                                var html = template(data);
                                var refundhtml = refundtemplate(data);
    
                                $("#usersProductsList").html(html);
    
                                $("#usersRefundsList").html(refundhtml);
    
                                _self.proListTable = JSON.parse(JSON.stringify(data.Message));
    
                                $("#usersProductsList").on(
                                    "click",
                                    ".refundBtn",
                                    function() {   

                                        var $this = $(this);
                                        var dom = this;
                                        _self.refundThisOrder.call(_self,dom,
                                            $this.attr("data-pindex"),
                                            $this.attr("data-paramorderid"),
                                            _self.proListTable
                                        );
                                    }
                                );

                            }
                           

                        }
                        
                   
                },
                error: function(jqXHR, textStatus, errorThrown) {
                    //_this.errorPop('服务器发生异常，请稍后再试~')
                }
            });
        }
    };
    page.init();
})(jQuery, window, layer);
