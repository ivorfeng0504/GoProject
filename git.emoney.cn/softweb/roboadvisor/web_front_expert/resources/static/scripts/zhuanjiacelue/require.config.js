require.config({
    baseUrl: window.gConfig.staticPath +"static/",
    waitSeconds: 30,
    paths: {
        jquery: "libs/jquery/jquery.min",
        layer: "libs/layer/layer",
        cors: "libs/jquery/jquery.xdomainrequest.min",
        moment: "libs/moment/moment.min",
        handlebars: "libs/handlebars/handlebars.min",
        swiper: "libs/swiper/swiper.min",
        nicescroll: "libs/nicescroll/jquery.nicescroll.min",
        localstorage: "modules/public/localstorage",
        webstorageapi: "modules/public/webstorage-apiclient",

        //global
        // config: "../../../config",
        utils: "scripts/modules/public/utils",
        // custom
        testdefault: "scripts/zhuanjiacelue/testdefault"
    },
    shim: {
        // modal: {
        //   deps: ['jquery', 'transition'],
        //   exports: 'modal'
        // }
        cors: ["jquery"],
        nicescroll: ["jquery"],
        utils: ["jquery"]
    }
});
require(["testdefault"], function (testdefault) {});
