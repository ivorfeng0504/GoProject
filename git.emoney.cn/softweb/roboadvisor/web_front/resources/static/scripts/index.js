/*
* @Author: gitshilly
* @Date:   2018-03-14 09:24:03
* @Version:   1.0
* @Last Modified by:   gitshilly
* @Last Modified time: 2018-03-16 15:26:17
*/

//初始化图片点击事件
function  initImgClickEvent() {
    $('.module-item img').unbind("click");
    $('.module-item img').click(function(event) {
        /* Act on the event */
        $('#ImgSrc').attr("src",$(this).attr("src"));
        $('#picPop').show();
    });
}
$(document).ready(function() {

    initImgClickEvent();
	
	$("#picPop").click(function () {
        $(this).hide();
    });
	// switchlab
	$('.switchlabpack').switchlab();
	
	// resize
	resize();
	$(window).resize(function(event) {
		/* Act on the event */
		resize();
	});

	
});
function resize(obj) {
    /* Act on the event */
    if(obj!=undefined){
        setHeight(obj);
        return;
	}
    $('.modulelist').each(function(index, el) {
        setHeight($(el));
    });
    function setHeight($obj){
    	var $top = $('.moduleitemTop',$obj.parent().parent());
        $obj.height($(window).height()-$('.switchlab-tit').height()-($top.is(":hidden")?0:$top.height())-3);
	}
};