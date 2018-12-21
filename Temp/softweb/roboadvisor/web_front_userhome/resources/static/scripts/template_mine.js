//已发放未领取奖品
var html_mine_grant =$("#html_mine_grant").html();
//已领取的奖品
var html_mine_receive =$("#html_mine_receive").html();
//我的产品
var html_mine_product =$("#html_mine_product").html();

//已发放未领取奖品模板
var template_mine_grant= Handlebars.compile(html_mine_grant);
//已领取的奖品模板
var template_mine_receive= Handlebars.compile(html_mine_receive);
//我的产品模板
var template_mine_product= Handlebars.compile(html_mine_product);


Handlebars.registerHelper("addOne",function(index,options){
    return parseInt(index)+1;
});