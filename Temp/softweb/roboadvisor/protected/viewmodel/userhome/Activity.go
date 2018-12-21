package viewmodel

type Activity struct {
	ActivityId  int64
	Title       string
	Summary     string
	BeginTime   string
	EndTime     string
	ImageUrl    string
	ActivityUrl string
	//是否已经参加
	IsJoin bool

	//打开方式（0 当前页打开 1 浏览器新窗口外链）
	OpenMode int

	//是否需要绑定手机号 0 否 1 是
	NeedBind int

	//是否需要附加SSO 0 否 1 是
	NeedSSO int
}
