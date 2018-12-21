# Web_API - 智投API接口模块
包含图文直播、盘外功夫等相关接口

#### version 1.2.5
* GetStrategyLiveRoomConfig 添加缓存
* 2018-08-27 17:15

#### version 1.2.4
* 更新连登接口
* 2018-08-27 14:15

#### version 1.2.3
* update  market_redis & captchastore_redis for production and uat
* 2018-08-09 13:15

#### version 1.2.2
* UserHome新增用户领取股票接口 /api/activity/insertuserstock
* UserHome新增获取用户当日领取的股票接口 /api/activity/getuserstocktoday
* UserHome新增获取用户领取的股票历史接口 /api/activity/getuserstockhistory
* UserHome新增添加股票池接口 /api/activity/insertstockpool
* UserHome新增获取最新的股票池接口 /api/activity/getnewststockpool
* 2018-06-05 12:00

#### version 1.2.1
* 更新日志配置
* 2018-06-01 11:40

#### version 1.2.0
* UserHome新增获取活动列表接口    /api/activity/getactivitylist
* UserHome新增获取用户参与活动Map接口   /api/activity/getuserinactivitymap
* UserHome新增获取用户连登奖励接口    /api/activity/getuserserialaward
* UserHome新增获取用户奖品列表接口    /api/activity/getuserawardlist
* UserHome新增获取历史信息-包括我的抽奖记录&往期中奖号码&中奖名单接口    /api/activity/getguessmarkethistoryinfo
* UserHome新增获取指定活动的奖品列表接口    /api/activity/getactivityawardlist
* UserHome新增用户获取抽奖号码，如果用户已经抽过奖则返回抽过的号码接口    /api/activity/getguessmarketnumber
* UserHome新增获取当前抽奖活动信息及用户参与信息接口    /api/activity/getcurrentguessinfo
* 2018-05-29 9:17

#### version 1.1.23
* 增加智投账号相关API接口（注册、找密、绑定）
* 2018-05-25 15:45

#### version 1.1.22
* 修改Redis连接的创建方式，支持自定义最大连接数
* 2018-05-09 20:00

#### version 1.1.21
* 数据库同步使用API网关进行调用 
* 2018-05-08 18:00

#### version 1.1.13
* 默认直播间列表修改为210
* 软件股池行情接口变更，定制逻辑策略对应不同lua脚本
* 2018-05-04 10:07

#### version 1.1.12
* 用户提问接口支持UID   /api/live/addquestion
* 查询直播内容接口使用IsTradeDay来判断交易日    /api/live/livecontent
* 2018-05-03 18:10

#### version 1.1.11
* 根据环境变量参数设置开发模式
* 增加回踩英雄池接口
* 2018-05-02 10:45

#### version 1.1.10
* 如果是非交易日显示上一个交易日的直播内容 /api/live/livecontent
* 交易日判断添加缓存
* 新增用户记录的时候使用redis进行加锁
* 2018-04-28 11:45

#### version 1.1.9
* 添加用户权限注销接口 /api/live/removeroompermit
* 2018-04-27 13:30

#### version 1.1.8
* 调整提问接口，添加来源以及昵称字段，如果用户未注册则自动添加用户信息 /api/live/addquestion
* 使用tokenapi.QueryToken来进行toke校验 /api/live/checktoken
* 接入数据库同步接口
* 手机加密添加缓存，减少接口调用
* 2018-04-26 18:03

#### version 1.1.7
* 直播间权限注册接口调整 /api/live/addroompermit
* 手机号加密处理
* 2018-04-25 18:31

#### version 1.1.6
* 用户问答添加屏蔽字处理 /api/live/addquestion
* 修复dotweb、dotlog日志相关的BUG
* 2018-04-24 13:30

#### version 1.1.5
* checktoken使用TokenServer来校验token  /api/live/checktoken
* 对用户提问内容进行html编码 初步防御xss攻击 /api/live/addquestion
* 添加用户查询、注册相关接口
* 2018-04-23 20:15

#### version 1.1.4
* 添加websocket授权回调接口  /api/live/checktoken
* 添加websocket授权消息回调查询接口 /api/live/hasnewmessage
* 2018-04-20 11:24

#### version 1.1.3
* 获取用户问答接口支持传递多个直播间查询 /api/live/livequestion
* 获取用户问答接口支持手机号查询 /api/live/livequestion
* 获取直播内容接口支持多直播间查询 /api/live/livecontent
* 如果指定日期中有一个直播间有开播记录 则忽略DisPlayLastTopic相关逻辑 /api/live/livecontent
* 用户直播间权限注册接口添加PID支持，有限使用RoomId注册，如果没有RoomId则尝试用PID查找RoomId的映射关系 /api/live/addroompermit
* 用户提问添加手机号支持 /api/live/addquestion
* 2018-04-18 19:45

#### version 1.1.2
* 添加用户直播间权限注册接口 /api/live/addroompermit
* 2018-04-16 20:26
* 添加获取当前用户拥有的直播间权限接口 /api/live/getuserroomlist
* 2018-04-17 14:35

#### version 1.1.1
* 修复Redis缓存模式下无法查询历史问答列表
* 2018-04-13 20:06

#### version 1.1.0
* 完善程序业务日志
* 支持客户端数据缓存
* 2018-04-13 15:15

#### version 1.0
* Init version
* 2018-04-12 17:51