# Task - 智投Task模块
智投相关Task

#### version 1.1.23
* Add Task Task_ProcessStockTalkMsgQueue 0 */10 * * * * 处理直播精选消息队列，将队列中的数据插入数据库
* 2018-09-07 15:40

#### version 1.1.18
* Task_ProcessUpdateQueue 每隔7分钟（35 */7 * * * *）执行一次 与Task_UpdateNewsInformationByDate Task时间尽量错开
* 2018-09-07 15:40

#### version 1.1.7
* 猜涨跌活动初始化Task 9:30进行重试
* 2018-08-10 10:55

#### version 1.1.6
* 修复err变量定义产生歧义的问题
* 2018-08-10 10:00

#### version 1.1.5
* TaskLogger Debug日志取消发邮件
* 2018-08-09 14:15

#### version 1.1.4
* update  market_redis & captchastore_redis for production and uat
* 2018-08-09 13:15

#### version 1.1.3
* 上线生产Task->Task_UpdateHotNewsInfoByClickNum
* 上线生产Task->Task_UpdateTopicInfo
* 上线生产Task->Task_SyncClientStrategyInfo
* 上线生产Task->Task_UpdateHotArticleListByClickNum
* 上线生产Task->Task_UpdateHotVideoListByClickNum
* 上线生产Task->Task_RefreshAllIndexStrategyNewsCache
* 2018-08-09 11:25

#### version 1.1.2
* 生产上线资讯要闻相关task
* 2018-08-08 11:40

#### version 1.1.1
* 更新日志配置
* 2018-06-01 11:40

#### version 1.1.0
* 新增Task_InitNextGuessMarketActivity任务    初始化下一期的活动
* 新增Task_InitCurrentGuessMarketActivity任务   初始化本期的活动（用于第一次初始化或者初始化下期失败后重试）
* 新增Task_GrantAward任务 发放奖励
* 2018-05-29 9:17

#### version 1.0.0
* Init version
* 2018-05-22 17:51