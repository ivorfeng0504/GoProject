

Handlebars.registerHelper('HasCoverImgFormat', function (items, options) {
    if(items!=null && items!="")
    {
        //满足添加继续执行
        return options.fn(this);
    } else {
        //不满足条件执行{{else}}部分
        return options.inverse(this);
    }
});


Handlebars.registerHelper('DateFormat', function (items,dateType, options) {
    //dateType:yyyy-MM-dd hh:mm:ss
    return dateFormat(items.replace(/T/g, ' ').replace(/-/g,'/'),dateType);
});

Handlebars.registerHelper('IsTodayFormat', function (items, options) {
    var lastmodifytime = new Date(items.replace(/T/g, ' ').replace(/-/g,'/'));
    var now = new Date();

    if(now.toDateString()==lastmodifytime.toDateString())
    {
        //满足添加继续执行
        return options.fn(this);
    } else {
        //不满足条件执行{{else}}部分
        return options.inverse(this);
    }
});

Handlebars.registerHelper('CoverImgFormat', function (items, options) {
    if(items == "")
    {
        var defaultimgurl = www + "/static/images/thumbs.png";
        return defaultimgurl;
    }else{
        return items;
    }
});

//是否是手机账号用户
Handlebars.registerHelper('IsMobileUserFormat', function (items, options) {
    if(items=="1")
    {
        //满足添加继续执行
        return options.fn(this);
    } else {
        //不满足条件执行{{else}}部分
        return options.inverse(this);
    }
});

Handlebars.registerHelper('HeadImgFormat', function (items, options) {
    if(items != "")
    {
        var defaultimgurl = www + "/static/images/Arena_" + items + ".png";
        return defaultimgurl;
    }else{
        return www + "/static/images/Arena_13.png";
    }
});

Handlebars.registerHelper('IsFreeUserFormat', function (items, options) {
    if(items=="0")
    {
        //满足添加继续执行
        return options.fn(this);
    } else {
        //不满足条件执行{{else}}部分
        return options.inverse(this);
    }
});

Handlebars.registerHelper('HasValueFormat', function (items, options) {
    if(items!=null && items!="")
    {
        //满足添加继续执行
        return options.fn(this);
    } else {
        //不满足条件执行{{else}}部分
        return options.inverse(this);
    }
});

Handlebars.registerHelper('VIPImgFormat', function (items, options) {
    var cls="";
    if (items=="1"){
        cls="viplv0";
    }
    if (items=="2"){
        cls="viplv1";
    }
    if (items=="3"){
        cls="viplv2";
    }
    return cls;
});
