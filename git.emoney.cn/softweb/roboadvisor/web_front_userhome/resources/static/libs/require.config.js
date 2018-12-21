requirejs.config({
    baseUrl: window.g_config.staticPath,
    waitSeconds: 30,
    map: {
        '*': {
            'css': 'require.plugin.css'
        }
    },
    paths: {
        //lib
        'jquery': 'libs/modules/public/jquery',  
        'json2':'libs/modules/public/json2.min',
        //template engine
         'handlebars':'libs/modules/public/handlebars',
        'shuangxiangpao':'scripts/shuangxiangpao.min'
    },
    shim: {
        
    }
});
require([
    'shuangxiangpao',
],function(
    shuangxiangpao
){});
