//直播内容模板
var html_live_content =$("#html_live_content").html();
//置顶直播内容模板
var html_live_content_top =$("#html_live_content_top").html();
//问答列表
var html_live_answer =$("#html_live_answer").html();

var template_live_content= Handlebars.compile(html_live_content);
var template_live_content_top= Handlebars.compile(html_live_content_top);
var template_html_live_answer= Handlebars.compile(html_live_answer);

//渲染消息类型
Handlebars.registerHelper('renderMessageType', function(messageType) {
    var html="";
    //1利好 2利空
    if(messageType==1){
        html+="<span class='boon' style='color:#ffffff'>利好</span>"
    }else if(messageType==2){
        html+="<span class='bad' style='color:#ffffff'>利空</span>"
    }
    return new Handlebars.SafeString(html);
});


//渲染消息内容
Handlebars.registerHelper('renderContent', function(content) {
    return new Handlebars.SafeString(content);
});

//渲染问答内容
Handlebars.registerHelper('renderAnswer', function(answer) {
    return new Handlebars.SafeString(answer);
});

//渲染股票列表
Handlebars.registerHelper('renderStockList', function(stockList) {
    var html="";
    if (stockList==null||stockList==undefined){ return html;}
    var stockCount=stockList.length;
    $(stockList).each(function (index,stock) {
        var changeClass="f-red";
        if(stock.Change.indexOf("-")==0){
            changeClass="f-green";
        }else if(stock.Change.trim()==("0%")){
            changeClass="";
        }

        var spliteline="<i class=\"spliteline\"></i>"
        if(index==stockCount-1){
            spliteline="";
        }
        html+="<li><a class='stock_info' href='javascript:void(0)' data-code='"+stock.StockCode+"' stockCode='"+stock.StockCode+"' stockType='"+stock.StockType+"'>"+stock.StockName+"</a><span class='"+changeClass+"'>"+stock.Change+"</span>"+spliteline+"</li>";
    });
    return new Handlebars.SafeString(html);
});