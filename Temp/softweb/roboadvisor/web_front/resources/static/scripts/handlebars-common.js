Handlebars.registerHelper('IsTopFormat', function (items, options) {
    if (items) {
        //满足添加继续执行
        return options.fn(this);
    } else {
        //不满足条件执行{{else}}部分
        return options.inverse(this);
    }
});
Handlebars.registerHelper('IsNotArticleFormat', function (items, options) {
    if(items!="0")
    {
        //满足添加继续执行
        return options.fn(this);
    } else {
        //不满足条件执行{{else}}部分
        return options.inverse(this);
    }
});
Handlebars.registerHelper('IsOneLessonFormat', function (items, options) {
    if(items!="1")
    {
        //满足添加继续执行
        return options.fn(this);
    } else {
        //不满足条件执行{{else}}部分
        return options.inverse(this);
    }
});
Handlebars.registerHelper('NewsTypeFormat', function (items, options) {
    if(items=="0")
    {
        return "img-icon";
    }else{
        return "img-icon video";
    }
});

Handlebars.registerHelper('bannerOlindexformat', function (items, options) {
    var ret = "";
    //items++;
    if(items == 0)
    {
        ret = "active t" + items;
    }else{
        ret = "t" + items;
    }
    return ret;
});

Handlebars.registerHelper('bannerPointformat', function (items, options) {
    if(items == "1")
    {
        //满足添加继续执行
        return options.fn(this);
    } else {
        //不满足条件执行{{else}}部分
        return options.inverse(this);
    }
});

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

Handlebars.registerHelper('CoverImgFormat', function (items, options) {
    if(items == "")
    {
        var defaultimgurl = www + "/static/images/thumbs.png";
        return defaultimgurl;
    }else{
        return items;
    }
});
