//获取客户端接口
function GetExternal() {
    return window.external.EmObj
};

//调用客户端接口
function PC_JH(type, c) {
    try {
        var obj = GetExternal();
        return obj.EmFunc(type, c)
    } catch (e) { }
};

//页面加载完成
function EM_FUNC_DOWNLOAD_COMPLETE() {
    try {
        PC_JH("EM_FUNC_DOWNLOAD_COMPLETE", "");
    }
    catch (ex) { }
};
function IsShow() {
    try { return PC_JH("EM_FUNC_WND_ISSHOW", ""); }
    catch (ex) { return "0"; }
}
//打开窗口
function openPCWindow() {
    EM_FUNC_DOWNLOAD_COMPLETE();
    if (IsShow() != "1") {
        PC_JH("EM_FUNC_WND_SIZE", "w=760,h=515,mid");
        EM_FUNC_SHOW();
    }
}
//关闭窗口
function closePCWindow() {
    try {
        PC_JH("EM_FUNC_DELETECLOSE", "");
    }
    catch (ex) { }
}
//页面显示
function EM_FUNC_SHOW() {
    try {
        PC_JH("EM_FUNC_SHOW", "");
    } catch (ex) {}
};

//打开商城页(客户端打开)
function EM_FUNC_BUYACCOUNT() {
    try {
        PC_JH("EM_FUNC_DISCOUNT", "");
    } catch (ex) {}
}
//个股
//EM_FUNC_GOTO_TECH_VIEW("600600,444") // 跳转到个股分时页面
//EM_FUNC_GOTO_TECH_VIEW("600600,455") // 跳转到个股日线页面
function GoKLine(stock) {
    if (stock.length == 6) {
        if (parseInt(stock) < 600000) {
            stock = "1" + stock;
        }
    }
    var par = "455," + stock;
    try {
        PC_JH("EM_FUNC_GOTO_TECH_VIEW", par)
    } catch (err) { }
    return false;
};

//板块

function GoBKLine(code) {
    if (code.length == 6) {
        code = code.substr(2, 4);
    }
    if (code.length == 4) {
        code = "BK" + code;
    }
    var par = "455," + code;
    try {
        PC_JH("EM_FUNC_GOTO_TECH_VIEW", par)
    } catch (err) { }
    return false
};

//打开个股页面
function OpenStock() {
    var stockCode = $(this).attr('stockCode');
    var stockType = $(this).attr('stockType');
    if (stockType == "A股") {
        GoKLine(stockCode);
    } else if (stockType == "行业板块") {
        GoBKLine(stockCode);
    }
    return false;
}

//打开个股页面
function NewsOpenStock(obj) {
    var stockCode = $(obj).attr('stockcode');
    var stockType = $(obj).attr('stocktype');
    if (stockType == "A股") {
        GoKLine(stockCode);
    } else if (stockType == "行业板块") {
        GoBKLine(stockCode);
    }
    return false;
}

//正在加载中
function LoadingDocument() {
    var obj_loading = $("#div_loading");
    if (obj_loading.length) {
        obj_loading.show();
    } else {
        $("#actList").prepend("<div id=\"div_loading\" class=\"loading\">正在加载...</div>");
    }
}
//隐藏加载
function HiddenLoadingDocument() {
    var obj_loading = $("#div_loading");
    obj_loading.hide();
}

function NoneDataDocument(obj) {
    obj.html("<div class=\"nonedata\">暂无数据</div>");
}

