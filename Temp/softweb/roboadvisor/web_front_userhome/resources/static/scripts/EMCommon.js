/*
 * @Author: zhangye
 * @Date:   2018-02-07 18:07:58
 * @Version:   1.0
 * @Last Modified by:   zhangye
 * @Last Modified time: 2018-02-13 15:07:57
 */
;
(function() {
    $.fn.extend({
        /**
         * 插件标志性的外框样式名：'switchlabpack'
         * param : '.switchlabpack' className
         * HTML结构 ：外框：.witchlabpack
         * lab切换标签：'.lab_tit'
         * 应用场景：页面加载完整后
         * eg: $(".class").switchlab();
         * eg: $('.switchlabpack').switchlab();
         */
        switchlab: function() {
            // this  ->  当前选择元素 - 外框
            var that = this;
            
            function switchAction(title, tits, conts) {
                var indx = tits.index(title);
                $(title).addClass('active').siblings().removeClass('active');
                conts.eq(indx).show().siblings().hide();
            }
            this.each(function(index, el) {
                var _switchPack = $(el);
                var _switchLabTitle = _switchPack.find(">.switchlab-tit");
                var _switchLabContent = _switchPack.find(">.switchlab-cnt");

                var _labTitIterms = _switchLabTitle.find(".switchlab-tit-i");
                var _labCntIterms = _switchLabContent.find('.switchlab-cnt-i');

                function callBackThis() {
                    var that = this;
                    switchAction(that, _labTitIterms, _labCntIterms);
                }

                _labTitIterms.off('click', switchAction);
                _labTitIterms.on('click', callBackThis);

            });
        },
        DateStringFormat: function (aStr, fmStr) {
            if (typeof aStr === "string" && typeof fmStr === "string") {
                return (new Date(parseInt(aStr.slice(6, -2))).format(fmStr));
            }
        },
        getUrlVars: function() {
            var vars = [],
                hash;
            var hashes = window.location.href.slice(window.location.href.indexOf('?') + 1).split('&');
            for (var i = 0; i < hashes.length; i++) {
                hash = hashes[i].split('=');
                vars.push(hash[0]);
                vars[hash[0]] = hash[1];
            }
            return vars;
        },
        getUrlVar: function(name) {
            // var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)","i");
            // var r = window.location.search.substr(1).match(reg);
            // if (r!=null) return unescape(r[2]); return null;
            return $.fn.getUrlVars()[name];
        }
    });

})();


/* ========================================================================
 * ZUI: date.js
 * Date polyfills
 * http://zui.sexy
 * ========================================================================
 * Copyright (c) 2014 cnezsoft.com; Licensed MIT
 * ======================================================================== */


(function() {
    'use strict';

    /**
     * Ticks of a whole day
     * @type {number}
     */
    Date.ONEDAY_TICKS = 24 * 3600 * 1000;

    /**
     * Format date to a string
     *
     * @param  string   format
     * @return string
     */
    if(!Date.prototype.format) {
        Date.prototype.format = function(format) {
            var date = {
                'M+': this.getMonth() + 1,
                'd+': this.getDate(),
                'h+': this.getHours(),
                'm+': this.getMinutes(),
                's+': this.getSeconds(),
                'q+': Math.floor((this.getMonth() + 3) / 3),
                'S+': this.getMilliseconds()
            };
            if(/(y+)/i.test(format)) {
                format = format.replace(RegExp.$1, (this.getFullYear() + '').substr(4 - RegExp.$1.length));
            }
            for(var k in date) {
                if(new RegExp('(' + k + ')').test(format)) {
                    format = format.replace(RegExp.$1, RegExp.$1.length == 1 ? date[k] : ('00' + date[k]).substr(('' + date[k]).length));
                }
            }
            return format;
        };
    }

    /**
     * Add milliseconds to the date
     * @param {number} value
     */
    if(!Date.prototype.addMilliseconds) {
        Date.prototype.addMilliseconds = function(value) {
            this.setTime(this.getTime() + value);
            return this;
        };
    }


    /**
     * Add days to the date
     * @param {number} days
     */
    if(!Date.prototype.addDays) {
        Date.prototype.addDays = function(days) {
            this.addMilliseconds(days * Date.ONEDAY_TICKS);
            return this;
        };
    }


    /**
     * Clone a new date instane from the date
     * @return {Date}
     */
    if(!Date.prototype.clone) {
        Date.prototype.clone = function() {
            var date = new Date();
            date.setTime(this.getTime());
            return date;
        };
    }


    /**
     * Judge the year is in a leap year
     * @param  {integer}  year
     * @return {Boolean}
     */
    if(!Date.isLeapYear) {
        Date.isLeapYear = function(year) {
            return(((year % 4 === 0) && (year % 100 !== 0)) || (year % 400 === 0));
        };
    }

    if(!Date.getDaysInMonth) {
        /**
         * Get days number of the date
         * @param  {integer} year
         * @param  {integer} month
         * @return {integer}
         */
        Date.getDaysInMonth = function(year, month) {
            return [31, (Date.isLeapYear(year) ? 29 : 28), 31, 30, 31, 30, 31, 31, 30, 31, 30, 31][month];
        };
    }


    /**
     * Judge the date is in a leap year
     * @return {Boolean}
     */
    if(!Date.prototype.isLeapYear) {
        Date.prototype.isLeapYear = function() {
            return Date.isLeapYear(this.getFullYear());
        };
    }


    /**
     * Clear time part of the date
     * @return {date}
     */
    if(!Date.prototype.clearTime) {
        Date.prototype.clearTime = function() {
            this.setHours(0);
            this.setMinutes(0);
            this.setSeconds(0);
            this.setMilliseconds(0);
            return this;
        };
    }


    /**
     * Get days of this month of the date
     * @return {integer}
     */
    if(!Date.prototype.getDaysInMonth) {
        Date.prototype.getDaysInMonth = function() {
            return Date.getDaysInMonth(this.getFullYear(), this.getMonth());
        };
    }


    /**
     * Add months to the date
     * @param {date} value
     */
    if(!Date.prototype.addMonths) {
        Date.prototype.addMonths = function(value) {
            var n = this.getDate();
            this.setDate(1);
            this.setMonth(this.getMonth() + value);
            this.setDate(Math.min(n, this.getDaysInMonth()));
            return this;
        };
    }


    /**
     * Get last week day of the date
     * @param  {integer} day
     * @return {date}
     */
    if(!Date.prototype.getLastWeekday) {
        Date.prototype.getLastWeekday = function(day) {
            day = day || 1;

            var d = this.clone();
            while(d.getDay() != day) {
                d.addDays(-1);
            }
            d.clearTime();
            return d;
        };
    }


    /**
     * Judge the date is same day as another date
     * @param  {date}  date
     * @return {Boolean}
     */
    if(!Date.prototype.isSameDay) {
        Date.prototype.isSameDay = function(date) {
            return date.toDateString() === this.toDateString();
        };
    }


    /**
     * Judge the date is in same week as another date
     * @param  {date}  date
     * @return {Boolean}
     */
    if(!Date.prototype.isSameWeek) {
        Date.prototype.isSameWeek = function(date) {
            var weekStart = this.getLastWeekday();
            var weekEnd = weekStart.clone().addDays(7);
            return date >= weekStart && date < weekEnd;
        };
    }


    /**
     * Judge the date is in same year as another date
     * @param  {date}  date
     * @return {Boolean}
     */
    if(!Date.prototype.isSameYear) {
        Date.prototype.isSameYear = function(date) {
            return this.getFullYear() === date.getFullYear();
        };
    }
}());


/* ========================================================================
 * ZUI: string.js
 * String Polyfill.
 * http://zui.sexy
 * ========================================================================
 * Copyright (c) 2014-2016 cnezsoft.com; Licensed MIT
 * ======================================================================== */


(function() {
    'use strict';

    /**
     * Format string with argument list or object
     * @param  {object | arguments} args
     * @return {String}
     */
    // if(!String.prototype.format) {
    //     String.prototype.format = function(fmt) {
    //         var o = {
    //             "M+" : this.getMonth()+1,                 //月份
    //             "d+" : this.getDate(),                    //日
    //             "h+" : this.getHours(),                   //小时
    //             "m+" : this.getMinutes(),                 //分
    //             "s+" : this.getSeconds(),                 //秒
    //             "q+" : Math.floor((this.getMonth()+3)/3), //季度
    //             "S"  : this.getMilliseconds()             //毫秒
    //         };
    //         if(/(y+)/.test(fmt)) {
    //             fmt=fmt.replace(RegExp.$1, (this.getFullYear()+"").substr(4 - RegExp.$1.length));
    //         }
    //         for(var k in o) {
    //             if(new RegExp("("+ k +")").test(fmt)){
    //                 fmt = fmt.replace(RegExp.$1, (RegExp.$1.length==1) ? (o[k]) : (("00"+ o[k]).substr((""+ o[k]).length)));
    //             }
    //         }
    //         return fmt;
    //     }
    // }

    /**
     * Judge the string is a integer number
     *
     * @access public
     * @return bool
     */
    // if(!String.prototype.isNum) {
    //     String.prototype.isNum = function(s) {
    //         if(s !== null) {
    //             var r, re;
    //             re = /\d*/i;
    //             r = s.match(re);
    //             return(r == s) ? true : false;
    //         }
    //
    //     };
    // }

    /**
     * filter
    */
    if (!Array.prototype.filter)
    {
        Array.prototype.filter = function(fn) {
            var newArray = [];
            var length = this.length;
            var i = 0;
            for(;i<length;i++){
                if(fn(this[i])){
                    newArray.push(this[i])
                }
            }
            return newArray;
        };
    }

    if (!Date.prototype.toISOString) {
        Date.prototype.toISOString = function() {
            function pad(n) { return n < 10 ? '0' + n : n }
            return this.getUTCFullYear() + '-'
                + pad(this.getUTCMonth() + 1) + '-'
                + pad(this.getUTCDate()) + 'T'
                + pad(this.getUTCHours()) + ':'
                + pad(this.getUTCMinutes()) + ':'
                + pad(this.getUTCSeconds()) + '.'
                + pad(this.getUTCMilliseconds()) + 'Z';
        }
    }


    if (!Array.prototype.forEach) {
        Array.prototype.forEach = function forEach(callback, thisArg) {
            var T, k;
            if (this == null) {
                throw new TypeError("this is null or not defined");
            }
            var O = Object(this);
            var len = O.length >>> 0;
            if (typeof callback !== "function") {
                throw new TypeError(callback + " is not a function");
            }
            if (arguments.length > 1) {
                T = thisArg;
            }
            k = 0;
            while (k < len) {
                var kValue;
                if (k in O) {
                    kValue = O[k];
                    callback.call(T, kValue, k, O);
                }
                k++;
            }
        };
    }

    })();