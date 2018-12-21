/* ========================================================================
 * app.js
 * ======================================================================== */
if (!window.console) {
    window.console = {};
}
if (!window.console.log) {
    window.console.log = function (msg) {};
}
// window.onerror = killErrors; function killErrors() {  return true;}

window.gConfig = {
    env: "serv",
    // apiHost: "http://dsclient.emoney.cn:8080/userhome",
    apiHost: www,
    staticPath: StaticServerHost
};

(function ($, window, layer) {

    // placeholder
    $(function () {
        $("input, textarea").placeholder({
            customClass: "pad-placeholder"
        });
    });

    var page = {

        refundListData:null,
        proListTable : null,
        refundOrderId :null,
        refundMode:null,
        refundProductList:null,
        init: function () {
            var _self = this;
            _self.prodouctList();

            $("#submitRefundBtn").click(function(){
                _self.submitRefundInfo(_self.refundOrderId, _self.refundMode);
            });

            $(".btn-gobacktoorderlist").click(function(){
                $("#userProListshowbox").show();
                $("#userRefundshowbox").hide();
                $(".btn-goback").show();
                $(".btn-gobacktoorderlist").hide();
                _self.refundOrderId = null;
                _self.refundMode = null;
            });
        },
        prodouctList: function () {
            var _self = this;
            _self.refundOrderId = null;
            _self.refundMode = null;

            var htmlCodes = [
                '{{#each Message}}',
                '<div class="upro-tablepack normal-tablepack"><div id="OrderId">',
                '<table class="table form-table" >',
                '    <caption class="">',
                '        <div class="row">',
                '            <div class="col-xs-6">',
                '                <span class="OrderId">订单号：{{ OrderId }}</span>',
                '            </div>',
                '            <div class="col-xs-6 text-right">',
                '                <span id="allCountPrice">优惠总价：',
                '                    <b>{{priceFixed2 Price }}元</b>',
                '                </span>',
                '            </div>',
                '        </div>',
                '    </caption>',
                '    <colgroup>',
                '        <col width="10%"></col>',
                '        <col width="35%"></col>',
                '        <col width="10%"></col>',
                '        <col width="8%"></col>',
                '        <col width="20%"></col>',
                '        <col width=""></col>',
                '    </colgroup>',
                '    <thead>',
                '        <tr>',
                '            <th>序号</th>',
                '            <th>商品</th>',
                '            <th>单价</th>',
                '            <th>数量</th>',
                '            <th>购买日期</th>',
                '            <th>最新状态</th>',
                '        </tr>',
                '    </thead>',
                '    <tbody>',
                '{{#each ProductList}}',
                '        <tr>',
                '            <td>{{xuhao @index }}</td>',
                '            <td> {{ ProductName }} </td>',
                '            <td> {{priceFixed2 Price }}元</td>',
                '            <td> {{ Count }} </td>',
                '            <td> {{timeformat1 CreateTime }}</td>',
                '            <td> {{ StateDesc }}</td>',
                '        </tr>',
                '{{/each}}',
                '    </tbody>',
                '</table>',
                '</div><div class="usr-opt-line text-right">',
                '{{#if CanRefund}}',
                '<span class="btn btn-primary refundBtn" data-pindex="{{@index}}" data-paramOrderID="{{OrderId}}"> 申请退款 </span> ',
                '{{/if}}',
                '{{#if IsRefund}}',
                '<a class="btn btn-primary refundDetailBtn" data-pindex="{{@index}}" href="' + window.gConfig.apiHost + 'mine/myorderdetail?orderid={{OrderId}}"> 退款详情 </a> ',
                '{{/if}}',
                '</div>',
                '</div>',
                '{{/each}}'
            ].join("");          

            Handlebars.registerHelper("priceFixed2", function (price) {
                var timestr3 = "";
                if (price == undefined || price == null || price.length == 0) {
                    timestr3 = "";
                } else {
                    timestr3 = parseFloat(""+price).toFixed(2);
                }
                return timestr3;
            });

            Handlebars.registerHelper("timeformat1", function (atime) {

                var timestr3 = "";

                if (atime == undefined || atime == null || atime.length == 0) {
                    timestr3 = "-";
                } else {
                    timestr3 = atime.slice(0,16);
                }
                return timestr3;
            });

            Handlebars.registerHelper("xuhao", function (idx) {

                var timestr3 = "";

                if (idx == undefined || idx == null || idx.length == 0) {
                    timestr3 = "-";
                } else {
                    timestr3 = idx+1;
                }
                return timestr3;
            });

            var template = Handlebars.compile(htmlCodes);
            $.ajax({
                url: window.gConfig.apiHost + "mine/queryorderlist",
                //    url: "../static/libs/mock/usrOrderList.json",
                type: "POST",
                timeout: 5000,
                dataType: "JSON",
                data: {
                    ColumnID: 16,
                    StrategyID: 0,
                    currpage: 1,
                    pageSize: 1
                },
                beforeSend: function () {},
                success: function (data) {
                    if (data.RetCode == 0 && data.Message != null && data.Message.length >0) {
                        $.extend(data, {});

                        var html = template(data);
                        $("#usersProductsList").html(html);
                        _self.proListTable = JSON.parse(JSON.stringify(data.Message));

                        $("#usersProductsList").on("click", '.refundBtn',function(){
                            var $this = $(this);
                            var dom = this;

                            _self.refundThisOrder.call(_self,  dom, $this.attr("data-pindex"), $this.attr("data-paramorderid"), _self.proListTable);
                        })
                    }
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    //_this.errorPop('服务器发生异常，请稍后再试~')
                }
            });

        },
        refundThisOrder:function( dom, indx,orderid,proddata){

            var $this = $(dom);

            var htmlCodes = "";

            page.refundOrderId = orderid;


            var opts = {
                totalPrice: parseFloat(proddata[indx].Price).toFixed(2),
                ProductList: proddata[indx].ProductList

            }

            var prolistHtmlCodes1 = [
                '<ul>',
                '{{#each ProductList}}',
                '   <li> {{ProductName}} </li>',
                '{{/each}}',
                '</ul>'
            ].join("");

            var ProlisthtmlCodes = [
                '<div class="upro-tablepack normal-tablepack">',
                '<table class="table form-table" >',
                '    <caption class="">',
                '    </caption>',
                '    <colgroup>',
                '        <col width="15%"></col>',
                '        <col width="45%"></col>',
                '        <col width="15%"></col>',
                '        <col width=""></col>',
                '    </colgroup>',
                '    <thead>',
                '        <tr>',
                '            <th>序号</th>',
                '            <th>商品</th>',
                '            <th>单价</th>',
                '            <th>数量</th>',
                '        </tr>',
                '    </thead>',
                '    <tbody>',
                '{{#each ProductList}}',
                '        <tr>',
                '            <td>{{xuhao @index }}</td>',
                '            <td> {{ ProductName }} </td>',
                '            <td> {{priceFixed2 Price }}元</td>',
                '            <td> {{ Count }} </td>',
                '        </tr>',
                '{{/each}}',
                '    </tbody>',
                '</table>',
                '</div>'
            ].join("");

            var prolistTemplate = Handlebars.compile(ProlisthtmlCodes);

            var banlistHtmlCodes = '{{#each BankInfoList}}<option value="{{BankValue}}">{{ BankName }}</option>{{/each}}';
            var banklistTemplate = Handlebars.compile(banlistHtmlCodes);


            

            $.ajax({
                url: window.gConfig.apiHost + "mine/validateorder",
                // url: "../static/libs/mock/isRefundOk.json?"+Math.random(),
                type: "POST",

                timeout: 5000,
                dataType: "JSON",
                data: {
                    OrderId: orderid
                },
                beforeSend: function () {},
                success: function (data) {

                    if (data.RetCode == 0) {
                        $.extend(data, opts);

                        page.refundMode = data.Message.ReFundMode;

                        if (data.Message.ReFundMode == 0) {

                            layer.msg("当前订单不能退款！");

                            $("#usrBankinfobox").hide();

                        } else if (data.Message.ReFundMode == 2) {

                            $("#userProListshowbox").hide();
                            $("#userRefundshowbox").show();
                            $(".btn-goback").hide();
                            $(".btn-gobacktoorderlist").show();

                            $("#usrBankinfobox").hide();
                            $("#nickName").html(page.refundOrderId);
                            $("#totalPrice").html("<span>" + opts.totalPrice + "元</span>");
                            $("#refundBank").html(data.Message.ReFundModeDesc);
                            if (data.Message.ReFundMode == 2) {
                                $("#refundBank + .inp-tips").html("退款将根据您购买时的支付方式，原路退回。");
                            } else {
                                $("#refundBank + .inp-tips").html("");
                            }
                            var prohtml = prolistTemplate(opts);
                            $("#refProductsList").html(prohtml);


                        }  else if (data.Message.ReFundMode == 1){

                            $("#userProListshowbox").hide();
                            $("#userRefundshowbox").show();
                            $(".btn-goback").hide();
                            $(".btn-gobacktoorderlist").show();


                            $("#usrBankinfobox").show();


                            $("#nickName").html(page.refundOrderId);
                            $("#totalPrice").html("<span>" +  opts.totalPrice + "元</span>" );

                            $("#refundBank").html(data.Message.ReFundModeDesc);
                            if (data.Message.ReFundMode == 1){
                                $("#refundBank + .inp-tips").html("由于您的支付方式不支持原路退回，请您准确地输入收款账户相关信息，以利于您及时收到退款。");
                            }else{
                                $("#refundBank + .inp-tips").html("");
                            }
                            var prohtml = prolistTemplate(opts);
                            $("#refProductsList").html(prohtml);

                            var bankhtml = banklistTemplate(data.Message);
                            $("#bankListshow").html(bankhtml);

                        }

                        //增加点赞处理

                    }else{
                        layer.msg(data.RetMsg);
                    }
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    //_this.errorPop('服务器发生异常，请稍后再试~')
                }
            });
        },
        dealProlist:function(){

        },
        myRefundList: function () { 
            var _self = this;
            _self.refundOrderId = null;
            _self.refundMode = null;
            var htmlCodes = [
                '{{#each Message}}',
                '<div class="upro-tablepack normal-tablepack"><div id="OrderId">',
                '<table class="table form-table" >',
                '    <caption class="">',
                '        <div class="row">',
                '            <div class="col-xs-6">',
                '                <span class="OrderId">111订单号1：{{ OrderId }}</span>',
                '            </div>',
                '            <div class="col-xs-6 text-right">',
                '                <span id="allCountPrice">优惠总价：',
                '                    <b>{{priceFixed2 Price }}元</b>',
                '                </span>',
                '            </div>',
                '        </div>',
                '    </caption>',
                '    <colgroup>',
                '        <col width="10%"></col>',
                '        <col width="35%"></col>',
                '        <col width="10%"></col>',
                '        <col width="8%"></col>',
                '        <col width="20%"></col>',
                '        <col width=""></col>',
                '    </colgroup>',
                '    <thead>',
                '        <tr>',
                '            <th>序号</th>',
                '            <th>商品</th>',
                '            <th>单价</th>',
                '            <th>数量</th>',
                '            <th>购买日期</th>',
                '            <th>最新状态</th>',
                '        </tr>',
                '    </thead>',
                '    <tbody>',
                '{{#each ProductList}}',
                '        <tr>',
                '            <td>{{xuhao @index }}</td>',
                '            <td> {{ ProductName }} </td>',
                '            <td> {{priceFixed2 Price }} </td>',
                '            <td> {{ Count }} </td>',
                '            <td> {{timeformat1 CreateTime }}</td>',
                '            <td> {{ StateDesc }}</td>',
                '        </tr>',
                '{{/each}}',
                '    </tbody>',
                '</table>',
                '</div><div class="usr-opt-line text-right">',
                '{{#if CanRefund}}',
                '<span class="btn btn-primary refundBtn" data-pindex="{{@index}}" data-paramOrderID="{{OrderId}}"> 退款 </span> ',
                '{{/if}}',
               
                '</div>',
                '</div>',
                '{{/each}}'
            ].join("");
            
            Handlebars.registerHelper("priceFixed2", function (price) {
                var timestr3 = "";
                if (price == undefined || price == null || price.length == 0) {
                    timestr3 = "";
                } else {
                    timestr3 = parseFloat(""+price).toFixed(2);
                }
                return timestr3;
            });

            Handlebars.registerHelper("timeformat1", function (atime) {

                var timestr3 = "";

                if (atime == undefined || atime == null || atime.length == 0) {
                    timestr3 = "-";
                } else {
                    timestr3 = atime.slice(0,10);
                }
                return timestr3;
            });

            Handlebars.registerHelper("xuhao", function (idx) {

                var timestr3 = "";

                if (idx == undefined || idx == null || idx.length == 0) {
                    timestr3 = "-";
                } else {
                    timestr3 = idx+1;
                }
                return timestr3;
            });

            var template = Handlebars.compile(htmlCodes);
            $.ajax({
                url: window.gConfig.apiHost + "mine/queryorderlist",
                //    url: "../static/libs/mock/usrOrderList.json",
                type: "POST",
                timeout: 5000,
                dataType: "JSON",
                data: {
                    ColumnID: 16,
                    StrategyID: 0,
                    currpage: 1,
                    pageSize: 1
                },
                beforeSend: function () {},
                success: function (data) {
                    if (data.RetCode == 0 && data.Message != null && data.Message.length >0) {
                        $.extend(data, {});

                        var html = template(data);
                        $("#usersProductsList").html(html);
                        _self.proListTable = JSON.parse(JSON.stringify(data.Message));

                        $("#usersProductsList").on("click", '.refundBtn',function(){
                            var $this = $(this);
                            var dom = this;


                            _self.refundThisOrder.call(_self,  dom, $this.attr("data-pindex"), $this.attr("data-paramorderid"), _self.proListTable);
                        })


                    }
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    //_this.errorPop('服务器发生异常，请稍后再试~')
                }
            });
        },
        submitRefundInfo: function (OrderId, RefundMode) {

            // 退还理由： rf_reasontype
            // 退货描述： rf_resoninfo
            // 退款姓名： rf_userName
            // 退款账号： rf_bkcount
            // 退款银行： bankListshow
            // 开户行：   rf_bkaddr

            if (!OrderId || !RefundMode){
                layer.msg("请选择退款订单！");
                return;
            }

            var rf_reasontype = $("#rf_reasontype").val();
            var rf_resoninfo = $("#rf_resoninfo").val();
            var rf_userName = $("#rf_userName").val();
            var rf_bkcount = $("#rf_bkcount").val();
            var bankListshow = $("#bankListshow").val();
            var rf_bkaddr = $("#rf_bkaddr").val();

            if (RefundMode == 1){
                
                if ($.trim(rf_resoninfo) == "") {
                    layer.msg("请填写退款理由!");
                    $("#rf_resoninfo").focus();
                    return;
                }
                if ($.trim(rf_resoninfo).length>200) {
                    layer.msg("退款理由不能超过200字!");
                    $("#rf_resoninfo").focus();
                    return;
                }

                if ($.trim(rf_userName) == "") {
                    layer.msg("请填写退款人姓名!");
                    $("#rf_userName").focus();
                    return;
                }

                if ($.trim(rf_bkcount) == "") {
                    layer.msg("请填写退款账号!");
                    $("#rf_bkcount").focus();
                    return;
                }

                if ($.trim(rf_bkaddr) == "") {
                    layer.msg("请填写开户行!");
                    $("#rf_bkaddr").focus();
                    return;
                }
            }
            var submitParam = {
                "OrderId": OrderId,
                "Reason": rf_reasontype + "|" + rf_resoninfo,
                "Name": rf_userName,
                "BankAccount": do_encrypt(rf_bkcount),
                "BankValue": bankListshow,
                "BankDetail": rf_bkaddr,
                "RefundMode": RefundMode
            }

            $.post(window.gConfig.apiHost + "mine/refundsubmit", submitParam, function (data) {

                if (data.RetCode == 0){
                    var popcnt = layer.open({
                        type: 1,
                        title: ['退款申请成功', 'text-align:center;'],
                        skin: 'ucpop-layer', //加上边框
                        area: ['540px', '285px'], //宽高
                        content: $('#usrOptInfor'), //'弹窗内容'
                        success: function (layero, index) {
                            $("#closeLayer1").one("click", function () {
                                layer.close(index);
                                window.location.reload();
                            })
                        }
                    });

                }else{
                    layer.msg(data.RetMsg);
                }

            })

        }

    };
    page.init();

})(jQuery, window, layer);