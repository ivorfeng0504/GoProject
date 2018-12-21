//当前活动模板页面
var html_activity_current =$("#html_activity_current").html();
//未开始的活动模板页面
var html_activity_future =$("#html_activity_future").html();
//已结束的活动模板页面
var html_activity_finish =$("#html_activity_finish").html();

//当前活动模板
var template_activity_current= Handlebars.compile(html_activity_current);
//未开始的活动模板
var template_activity_future= Handlebars.compile(html_activity_future);
//已结束的活动模板
var template_activity_finish= Handlebars.compile(html_activity_finish);
