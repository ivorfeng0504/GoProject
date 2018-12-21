var $StrategyPage = {
    fundsmart: 'http://www.fundsmart.com.cn:8888/fundsmart/api/emoney/',
    fundCode: null,
    checkedType: 'W',
    piecolors: ['#884DE6', '#8C93F6', '#E65362', '#F1BC5B', '#91c7ae', '#749f83', '#ca8622', '#bda29a', '#6e7074', '#546570', '#c4ccd3'],
    timer: null,
    //组合走势 折线图
    groupTechChart: null,
    //策略配置 饼图
    strategyChart: null,
    //实时估值 走势图
    estimationChart: null,
    //策略配置的全部信息
    strategyAllInfoList: null,
    //策略配置 的历史信息
    strategyConfigInfoList: null,
    //事件绑定
    bindEvent: function () {
        var _this = this;
        var $container = $('#groupTechContainer');
        var $investType = $('.invest-type');
        var $poptype = $('#poptype');

        //组合走势 事件绑定
        $container.on('click', '[data-type]', function () {
            $(this).addClass("actived").siblings().removeClass('actived');

            var _type = $('.switch-header-i.actived', $container).attr("data-type");
            var $checked = $('.charttype .actived', $container);
            var _date = $checked.attr('data-type');

            $StrategyPage.checkedType = _type;
            $('#spchartType').text($checked.html());
            if (!!$StrategyPage.fundCode) {
                _this.loadGroupTech(_type, _date);
            }
        });

        $investType.on('click', 'dt', function () {
            $investType.toggleClass('active');
        });
        $poptype.on('click', 'li', function () {
            $(this).addClass('active').siblings().removeClass('active');
            $investType.toggleClass('active');
        });
        $('#selAdjustDates').change(function () {
            var _date = $(this).val();
            _this.loadStrategyInfos(_date);
        });
        $('[data-type]', $container).eq(0).trigger('click');
    },

    //实时估值
    loadEstimation: function () {
        $StrategyPage.getEstimationChart(function (data) {
            $StrategyPage.initEstimationChart(data);
        });
    },
    //实时估值 右侧数据绑定
    loadEstimationInfos: function (data) {
        if (!data) return;

        var allInfo = $StrategyPage.strategyAllInfoList;
        if (!allInfo) return;
        allInfo = JSON.parse(allInfo);
        var value = null;
        $.each(allInfo, function () {
            var strategyList = this.StrategyList;
            if (!!value) return false;
            $.each(strategyList, function () {
                var item = this;
                if (item.code == $StrategyPage.fundCode) {
                    value = item;
                    return false;
                }
            });
        });
        if (!value) return;
        var levelInfo = Common.GetMyLevel(value.style)
        if (!levelInfo) return;
        var latestInfo = value.LatestInfo;
        var rightInfo = {};
        rightInfo.name = value.name;
        rightInfo.style = levelInfo.name;
        rightInfo.lables = value.labels.join(" ");
        rightInfo.level = levelInfo.tag;
        rightInfo.levelInfo = levelInfo.level + "-" + levelInfo.name;
        rightInfo.navDate = !!latestInfo.navDate && latestInfo.navDate.length >= 10 ? latestInfo.navDate.substring(5, 10) : "--";

        if (!!latestInfo.nav) {
            rightInfo.nav_class = Number(latestInfo.nav) > 0 ? "fred" : "fgreen";
            rightInfo.nav = latestInfo.nav;
        }
        else {
            rightInfo.current = "--";

        }
        if (!!latestInfo.navChange) {
            rightInfo.navChange = Common.ConvertPercent(latestInfo.navChange);
            rightInfo.navChange_class = latestInfo.navChange > 0 ? "fred" : "fgreen";
        }
        else {
            rightInfo.navChange_class = "--";
        }

        if (!!value.LatestInfo.annualChangeAll) {
            rightInfo.annualChangeAll = Common.ConvertPercent(value.LatestInfo.annualChangeAll);
            rightInfo.annualChangeAll_class = value.LatestInfo.annualChangeAll > 0 ? "fred" : "fgreen";
        }
        else {
            rightInfo.annualChangeAll = "--";
            rightInfo.annualChangeAll_class = "";
        }
        Common.CommonBind(rightInfo, "estimationRightInfo", false);
    },
    //实时估值 获取数据
    getEstimationChart: function (callback) {
        if (!$StrategyPage.fundCode) return;
        Common.CommonAjax("home/GetFundTimeLineByCode", {code: $StrategyPage.fundCode}, function (response) {
                if (response.RetCode !== 0) {
                    layer.msg("获取基金实时估值数据失败");
                    console.log(response.RetMsg);
                    return;
                }
                var data = JSON.parse(response.Message);

                var _timedate = data.timeline;
                var retdata = {last: data.last, allData: [], xAxisData: [], yAxisData: [], yAxisData1: []};
                var timestamps = Common.CreateTimeArray();
                var _data = [];
                for (var index = 0; index < timestamps.length; index++) {
                    var element, _nav = null, rate = null;
                    if (!!_timedate[index]) {
                        _nav = _timedate[index].nav;
                        rate = ((_nav / data.last - 1) * 100);
                        rate = Number(rate.toFixed(2));
                        if (isNaN(rate) || rate === 0) {
                            rate = 0;
                        }
                    }
                    var currDate = new Date();
                    var _timestamp = new Date(currDate.getFullYear() + '/' + (currDate.getMonth() + 1) + '/' + currDate.getDay() + ' ' + timestamps[index]).getTime();
                    retdata.yAxisData.push({x: _timestamp, y: _nav});
                    retdata.yAxisData1.push({x: _timestamp, y: _nav});
                    retdata.allData.push([_timestamp, _nav, rate]);
                }

                $('#spLastEstimation').html(data.current);
                var spzfEstimation = (data.current / data.last - 1) * 100;
                spzfEstimation = Number(spzfEstimation.toFixed(2));
                var spzfEstimation_class = "";
                if (!isNaN(spzfEstimation) && spzfEstimation !== 0) {
                    if (spzfEstimation > 0) {
                        spzfEstimation_class = "fred";
                    }
                    else {
                        spzfEstimation_class = "fgreen";
                    }
                }
                else {
                    spzfEstimation = 0;
                }
                $('#spzfEstimation').html(spzfEstimation + '%').addClass(spzfEstimation_class);

                $StrategyPage.loadEstimationInfos(data);

                if (callback && (typeof callback == 'function')) {
                    callback(retdata);
                }
            }
            ,
            "GET"
        );
    },
    //实时估值 走势图绑定
    initEstimationChart: function (retdata) {
        Highcharts.setOptions({global: {useUTC: false}});
        Highcharts.stockChart("EstimationChart", {
            chart: {
                plotBorderWidth: 1,
                plotBorderColor: "#e6e6e6",
                backgroundColor: "none",
                animation: !1,
                marginLeft: 45,
                marginRight: 60,
                alignTicks: !0,
                style: {
                    fontFamily: '"Helvetica Neue", Arial, "Microsoft YaHei"',
                    fontSize: "12px"
                }
            },
            credits: {
                enabled: !1
            },
            legend: {
                enabled: !1
            },
            navigator: {
                enabled: !1,
                maskFill: "#fff",
                maskInside: !1,
                series: {
                    type: "areaspline",
                    color: "#4572A7",
                    fillOpacity: .4,
                    dataGrouping: {
                        smoothed: !0
                    },
                    lineWidth: 1,
                    marker: {
                        enabled: !1
                    }
                },
                xAxis: {
                    xAxis: {
                        type: "datetime",
                        labels: {
                            overflow: "justify"
                        }
                    },
                    labels: {
                        style: {},
                        formatter: function () {
                            return Highcharts.dateFormat("%b-%e", this.value)
                        },
                        x: -15,
                        y: 28
                    }
                },
                yAxis: {
                    gridLineWidth: 0,
                    startOnTick: !1,
                    endOnTick: !1,
                    minPadding: .1,
                    maxPadding: .1,
                    labels: {
                        enabled: !1
                    },
                    title: {
                        text: null
                    },
                    tickWidth: 0
                }
            },
            plotOptions: {
                line: {
                    color: "#67ace9",
                    lineWidth: 1.4,
                    animation: !1,
                    states: {
                        hover: {
                            lineWidth: .4
                        }
                    },
                    pointStart: 0
                }
            },
            rangeSelector: {
                enabled: !1
            },
            scrollbar: {
                enabled: !1
            },
            series: [{
                type: "line",
                name: "\u5f53\u524d\u4ef7",
                enabledCrosshairs: !0,
                data: retdata.yAxisData,
                yAxis: 0
            }, {
                type: "line",
                name: "\u5f53\u524d\u4ef7",
                enabledCrosshairs: !0,
                data: retdata.yAxisData1,
                yAxis: 1
            }],
            tooltip: {
                crosshairs: [{
                    color: "#ffcbcc",
                    width: 1
                }, {
                    color: "#ffcbcc",
                    width: 1
                }],
                followPointer: !1,
                useHTML: !0,
                borderColor: "#ccc",
                style: {},
                formatter: function () {
                    var e = this.points[0].point.index,
                        a = null;
                    a = retdata.allData[e] ? retdata.allData[e][2] : 0;
                    var i = "";
                    i = "<span>" + this.y + "</span>", a = 0 > a ? '<span style="color:#3fb539">' + a + "%</span>" :
                        '<span style="color:#ff5256">' + a + "%</span>";
                    var n = "";
                    return n = '<span style="line-height: 20px; padding: 8px 10px;">\u65f6\u95f4\uff1a' +
                        Highcharts.dateFormat("%H:%M", this.x) +
                        '</span><br><span style="line-height: 20px; padding: 8px 10px;">\u4f30\u503c\uff1a' + i +
                        '\u5143</span><br><span style="line-height: 20px; padding: 8px 10px;">\u6da8\u5e45\uff1a' +
                        a + "</span>"
                }
            },
            yAxis: [{
                offset: -370,
                labels: {
                    x: 12,
                    y: 4,
                    style: {
                        color: "#a5a5a5",
                        fontSize: "12px"
                    },
                    formatter: function () {
                        var t = this.value.toFixed(4);
                        return t = t > retdata.last ? '<span style="color:#f52f3e">' + t + "</span>" :
                            '<span style="color:#69cd8e">' + t + "</span>"
                    }
                },
                plotLines: [{
                    value: retdata.last,
                    color: "#5e5e5e",
                    dashStyle: "shortdash",
                    width: 1,
                    label: {
                        align: "right",
                        x: -5,
                        y: -4,
                        text: '<span class="add-bor">0%</span>',
                        style: {
                            backgroundColor: "#fff"
                        }
                    },
                    zIndex: 1
                }, {
                    value: retdata.last,
                    color: "#5e5e5e",
                    dashStyle: "shortdash",
                    width: 1,
                    label: {
                        align: "left",
                        x: 5,
                        text: retdata.last,
                        style: {
                            backgroundColor: "#fff"
                        }
                    },
                    zIndex: 1
                }],
                gridLineColor: "#f0f0f0",
                showLastLabel: !0
            }, {
                offset: -12,
                labels: {
                    y: 4,
                    style: {
                        color: "#a5a5a5",
                        fontSize: "12px"
                    },
                    formatter: function () {
                        var t = ((this.value - retdata.last) / retdata.last * 100).toFixed(2);
                        return t = t > 0 ? '<span style="color:#f52f3e">' + t + "%</span>" :
                            '<span style="color:#69cd8e">' + t + "%</span>"
                    }
                },
                gridLineColor: "#f6f6f6",
                showLastLabel: !0
            }],
            xAxis: {
                gridLineColor: "#f0f0f0",
                gridLineWidth: 1,
                gridLineDashStyle: "dash",
                showLastLabel: !0,
                showFirstLabel: !0,
                tickWidth: 0,
                tickPixelInterval: 50,
                labels: {
                    formatter: function () {
                        var t = Highcharts.dateFormat("%H:%M", this.value);
                        return "13:00" == t ? "11:30/13:00" : "09:30" == t || "10:30" == t || "14:00" == t ||
                        "15:00" == t ? t : void 0
                    },
                    rotation: 0,
                    staggerLines: 2,
                    style: {
                        color: "#a5a5a5",
                        whiteSpace: "normal"
                    }
                }
            }
        });
    },

    //组合走势 获取数据
    loadGroupTech: function (type, date) {
        $.ajax({
            url: $StrategyPage.fundsmart + 'fund/trends/' + $StrategyPage.fundCode + '_' + type + '_' + date + '.json',
            type: "get",
            dataType: "jsonp",
            data: {},
            success: function (data, status) {
            },
            error: function (jqXHR, textStatus, errorThrown) {
            }
        });
    },
    //组合走势 左侧数据绑定
    initGroupTechChart: function (data) {
        if (!$StrategyPage.groupTechChart) {
            $StrategyPage.groupTechChart = echarts.init(document.getElementById('groupTechChart'));
        }
        var myChart = $StrategyPage.groupTechChart;
        var xAxisData = [], yAxisData1 = [], yAxisData2 = [], yAxisData3 = [];
        var _trendData = data.trends;
        if (!_trendData || !_trendData.length) {
            return;
        }
        var min = 100 * 10000;
        for (var index = 0; index < _trendData.length; index++) {
            var element = _trendData[index];
            xAxisData.push(element.date);
            yAxisData1.push(element.mv);
            min = min > element.mv ? element.mv : min;
            yAxisData2.push(element.bmv);
            min = min > element.bmv ? element.bmv : min;
            yAxisData3.push(element.cost);
        }
        var options = {
            tooltip: {
                trigger: 'axis'
            },
            //X轴设置
            xAxis: [{
                type: 'category',
                boundaryGap: false,
                data: xAxisData
            }],
            grid: {
                right: 10
            },
            color: ['#F80000', '#6FACE7'],
            legend: {
                x: 'center',
                y: 'bottom',
                icon: 'rect',
                data: ['本策略', '比较基准']
            },
            yAxis: [
                {
                    type: 'value',
                    axisLabel: {
                        formatter: '{value}'
                    },
                    splitLine: {
                        lineStyle: {
                            type: 'dashed'
                        }
                    },
                    min: Number(((min - 100) / 100).toFixed(0)) * 100
                }
            ],
            series: [{
                name: '本策略',
                type: 'line',
                z: '3',
                showSymbol: false,
                data: yAxisData1
            }, {
                name: '比较基准',
                type: 'line',
                z: '2',
                showSymbol: false,
                data: yAxisData2
            }]
        };
        if ($StrategyPage.checkedType === 'K') {
            options.series.unshift({
                name: '投资额',
                type: 'line',
                color: '#E5E5E5',
                z: '1',
                showSymbol: false,
                areaStyle: {},
                data: yAxisData3
            });
        }
        myChart.clear();
        myChart.setOption(options);
        $('#mdStart').text(data.mdStart);
        window.onresize = myChart.resize;
    },
    //组合走势 右侧数据绑定
    loadGroupTechInfos: function (data) {
        $('#yield').attr('class', this.getClass(data.percent)).html(data.percent + '%');
        $('#maxretrace').attr('class', this.getClass(data.maxDown)).html(data.maxDown + '%');
        $('#xpRate').attr('class', this.getClass(data.sharpe)).html(data.sharpe + '%');

        var run = "";
        $('#percentDiff').attr('class', this.getClass(data.percentDiff)).html(Math.abs(data.percentDiff) + '%');
        if (Number(data.percentDiff) > 0) {
            run = "超出";
        }
        else {
            run = "跑输";
        }
        $('#percentRun').html(run);

        $('#sharpeDiff').attr('class', this.getClass(data.sharpeDiff)).html(Math.abs(data.sharpeDiff) + '%');
        if (Number(data.sharpeDiff) > 0) {
            run = "超出";
        }
        else {
            run = "跑输";
        }
        $('#sharpeRun').html(run);

        $('#mdStart-End').html((data.mdStart + "~" + data.mdEnd).replace(/-/g, "/"));
    },
    getClass: function (percent) {
        if (Number(percent) > 0) {
            return 'fred';
        } else if (Number(percent) < 0) {
            return 'fgreen';
        }
        return 'fgray';
    },

    //策略配置 绑定
    loadStrategyInfos: function (date) {
        var _this = this;
        var jsonData = $StrategyPage.strategyConfigInfoList;
        if (!jsonData) return;
        jsonData = JSON.parse(jsonData);
        var data = null;
        $.each(jsonData, function () {
            var item = this;
            if (item.date == date) {
                data = item;
                return false;
            }
        })
        if (!data) return;
        var _sampledata = {}, chartData = [], _list, _typelist = [];
        if (!data.samples || !data.samples.length) {
            return;
        }
        _list = data.samples;
        for (var index = 0; index < _list.length; index++) {
            var element = _list[index];
            if (!_sampledata[element.type]) {
                _typelist.push(element.type);
            }
            _sampledata[element.type] = (Number((_sampledata[element.type] || 0)) + Number(element.weight)).toFixed(2);
        }

        for (index = 0; index < _typelist.length; index++) {
            var element = _typelist[index];
            chartData.push({value: _sampledata[element], name: element});
        }
        _this.loadStrategyinfo(chartData, data);
        _this.initStrategyChart(chartData);

    },
    //策略配置 调仓日期 绑定
    loadAdjustDates: function (dateList) {
        var _tploption = '<option value="{{date}}">{{date}}</option>';
        var _arrHtmls = [];
        dateList = dateList.sort();
        for (var index = 0; index < dateList.length; index++) {
            var element = dateList[dateList.length - index - 1];
            _arrHtmls.push(_tploption.replaceFormat({date: element}));
        }
        $('#selAdjustDates').html(_arrHtmls.join(''));
        this.loadStrategyInfos(dateList[dateList.length - 1]);
    },
    //策略配置 右侧数据绑定
    loadStrategyinfo: function (chartData, retdata) {
        var infoData = [], data = retdata.samples;
        for (var index = 0; index < chartData.length; index++) {
            var _dttpl = '<dt><span style="color: ' + $StrategyPage.piecolors[index] + ';">{{name}}</span><span>总占比 {{value}}%</span></dt>';
            var ddArray = [];
            for (var j = 0; j < data.length; j++) {
                var element = data[j];
                if (chartData[index].name == element.type) {
                    var _ddtpl = '<dd><span>{{name}} {{ticker}}</span><span>{{weight}}%</span></dd>';
                    ddArray.push(_ddtpl.replaceFormat(element));
                }

            }
            infoData.push('<dl>' + _dttpl.replaceFormat(chartData[index]) + ddArray.join('') + '</dl>');
        }

        $('#strategyinfo .s-tips').html(retdata.memo);
        var $details = $('#strategyinfo .s-details');
        $details.html(infoData.join(''));
        $('dl', $details).each(function (i, n) {
            if (i % 2 === 0) {
                var _idx = $(this).index(), _nextidx = _idx + 1 > $('dl', $details).length ? 0 : _idx + 1;
                var _height = $(this).height();
                if (_nextidx != 0) {
                    var $nextobj = $('dl', $details).eq(_nextidx), _nextheight = $nextobj.height();
                    if (_nextheight > _height) {
                        $(this).css({height: _nextheight});
                    } else {
                        $nextobj.css({height: $(this).outerHeight()})
                    }
                }
            }
        });
    },
    //策略配置 左侧饼图绑定
    initStrategyChart: function (data) {
        if (!$StrategyPage.strategyChart) {
            $StrategyPage.strategyChart = echarts.init(document.getElementById('StrategyChart'));
        }
        var myChart = $StrategyPage.strategyChart;
        var options = {
            tooltip: {
                trigger: 'item',
                formatter: "{a} <br/>{b}: {c}%"
            },
            color: $StrategyPage.piecolors,
            series: [
                {
                    name: '策略配置',
                    type: 'pie',
                    radius: ['50%', '90%'],
                    avoidLabelOverlap: false,
                    label: {
                        normal: {
                            show: false,
                            position: 'center'
                        },
                        emphasis: {
                            show: true,
                            textStyle: {
                                fontSize: '30',
                                fontWeight: 'bold'
                            }
                        }
                    },
                    labelLine: {
                        normal: {
                            show: false
                        }
                    },
                    data: data
                }
            ]
        };
        myChart.setOption(options);
    }
};
