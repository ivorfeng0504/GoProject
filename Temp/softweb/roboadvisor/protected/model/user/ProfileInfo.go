package model

type ProfileInfo struct {
	/// <summary>
	/// 用户姓名
	/// </summary>
	Name string
	/// <summary>
	/// 年龄
	/// </summary>
	Age string
	/// <summary>
	/// 地址
	/// </summary>
	Address       string
	UserAttribute string
	/// <summary>
	/// 邮编
	/// </summary>
	Description string

	/// <summary>
	/// 固话号码 
	/// </summary>
	Phone string
	/// <summary>
	/// 手机号码
	/// </summary>
	Mobile string
	/// <summary>
	/// 用户ID
	/// </summary>
	UserID string
	/// <summary>
	/// 用户账号名
	/// </summary>
	UserName string
	/// <summary>
	/// 身份证号码
	/// </summary>
	IDCard string
	/// <summary>
	/// 地址
	/// </summary>
	Email string
	/// <summary>
	/// 性别
	/// </summary>
	Sex string
	/// <summary>
	/// 省
	/// </summary>
	CoName string
	/// <summary>
	/// 市
	/// </summary>
	ProvinceName string
	/// <summary>
	/// 区
	/// </summary>
	ProvinceName1 string
	/// <summary>
	/// 操作风格分数
	/// </summary>
	TZvalue string
	/// <summary>
	/// QQ号码
	/// </summary>
	QQ string
	/// <summary>
	/// 服务代理商ID，默认:100000000
	/// </summary>
	ServiceAgentId string

	//服务商
	ServiceAgentName string
	/// <summary>
	/// ? 默认 1
	/// </summary>
	AgreementA string
	AgreementB string

	/// <summary>
	/// 国际区号
	/// </summary>
	InternationalNumber string
	/// <summary>
	/// 区号
	/// </summary>
	ProvinceNumber string
	/// <summary>
	/// 分机号
	/// </summary>
	SmallNumber string
	/// <summary>
	/// 头像
	/// </summary>
	Picture string

	/// <summary>
	/// 昵称
	/// </summary>
	Nickname string

	/// <summary>
	/// 微信号
	/// </summary>
	WeixinID string

	/// <summary>
	/// 生日年份
	/// </summary>
	Birth_year string

	/// <summary>
	/// 生日月份
	/// </summary>
	Birth_month string

	/// <summary>
	/// 生日日期
	/// </summary>
	Birth_day string

	/// <summary>
	/// 常用联系方式
	/// </summary>
	Contact string

	/// <summary>
	/// 投资偏好
	/// </summary>
	Tzph string

	/// <summary>
	/// 资金规模
	/// </summary>
	Zjgm string

	/// <summary>
	/// 职业
	/// </summary>
	Profession string

	/// <summary>
	/// 学历
	/// </summary>
	Degree string

	/// <summary>
	/// 是否有盯盘时间
	/// </summary>
	Dp_time string
}
