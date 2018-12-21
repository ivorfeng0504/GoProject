(function(){
  function GetExternal() {
    return window.external.EmObj;
  }
  function PC_JH(type, c) {
    try {
      var obj =
        GetExternal();
      return obj.EmFunc(type, c);
    } catch (e) { }
  }
  function LoadComplete() {
    try {
      PC_JH("EM_FUNC_DOWNLOAD_COMPLETE", "");
    } catch (ex) { }
  }
  function EM_FUNC_HIDE() {
    try {
      PC_JH("EM_FUNC_HIDE", "");
    } catch (ex) { }
  }
  function EM_FUNC_SHOW() {
    try {
      PC_JH("EM_FUNC_SHOW", "");
    } catch (ex) { }
  }
  function IsShow() {
    try { return PC_JH("EM_FUNC_WND_ISSHOW", ""); }
    catch (ex) { return "0"; }
  }
  function openWindow() {
    LoadComplete();
    if (IsShow() != "1") {
      PC_JH("EM_FUNC_WND_SIZE", "w=1106,h=711,mid");
      EM_FUNC_SHOW();
    }
  }
  openWindow();
})();