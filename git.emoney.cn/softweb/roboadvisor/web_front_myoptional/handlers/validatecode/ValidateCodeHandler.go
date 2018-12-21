package validatecode

import (
	"time"
	"errors"
	"net/http"
	"text/template"
	"io"
	"github.com/devfeel/dotweb"
	"bytes"
	"git.emoney.cn/softweb/roboadvisor/util/captcha"
	"git.emoney.cn/softweb/roboadvisor/protected/service/expertnews"

	"git.emoney.cn/softweb/roboadvisor/protected/service/strategyservice"
)

const (
	StdWidth  = 80
	StdHeight = 30
)

var (
	ErrNotFound = errors.New("captcha: id not found")

)

var formTemplate = template.Must(template.New("example").Parse(formTemplateSrc))


func GetCaptchaId(ctx dotweb.Context)error{
	d := struct {
		CaptchaId string
	}{
		captcha.New(),
	}
	return ctx.WriteJson(d)
}


//验证码页面
func ShowCaptchaPage(ctx dotweb.Context) error {

	strategySrv := strategyservice.NewRedisStrategyInfoService()
	//strategySrv.GetStrategyGroup()
	//strategySrv.GetExpertStrategyList()

	sdata :=strategySrv.GetStrategyData()
	return ctx.WriteJson(sdata)


	d := struct {
		CaptchaId string
	}{
		captcha.New(),
	}

	err := formTemplate.Execute(ctx.Response().Writer(), &d)
	if err != nil {
		http.Error(ctx.Response().Writer(), err.Error(), http.StatusInternalServerError)
	}


	//获取主题涉及股票关注次数的统计列表
	topicSrv := expertnews.NewExpertNews_TopicService()
	ranklist, _ :=  topicSrv.StatTopicFocusStock()
	ctx.WriteString("<br><br>")
	ctx.WriteString(ranklist)
	return err
}

//输出图片
func BuffImage(ctx dotweb.Context) error {
	w := ctx.Response().Writer()
	r := ctx.Request().Request
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	captchaid := ctx.QueryString("captchid")
	var content bytes.Buffer

	w.Header().Set("Content-Type", "image/png")
	if err := captcha.WriteImage(&content, captchaid,StdWidth ,StdHeight);err != nil{
		http.Error(ctx.Response().Writer(), err.Error(), http.StatusInternalServerError)
	}else{
		http.ServeContent(w,r, captchaid, time.Time{}, bytes.NewReader(content.Bytes()))
	}

	return nil
}


func BuffNewImage(ctx dotweb.Context) error {
	w := ctx.Response().Writer()
	r := ctx.Request().Request
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	captchaid := ctx.QueryString("captchid")
	var content bytes.Buffer

	w.Header().Set("Content-Type", "image/png")
	captcha.Reload(captchaid)

	if err := captcha.WriteImage(&content, captchaid,StdWidth ,StdHeight);err != nil{
		http.Error(ctx.Response().Writer(), err.Error(), http.StatusInternalServerError)
	}else{
		http.ServeContent(w,r, captchaid, time.Time{}, bytes.NewReader(content.Bytes()))
	}

	return nil
}

//验证图形验证码
func VerifyCaptcha(ctx dotweb.Context) error {
	w := ctx.Response().Writer()
	r := ctx.Request().Request

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
		io.WriteString(w, "Wrong captcha solution! No robots allowed!\n")
	} else {
		io.WriteString(w, "Great job, human! You solved the captcha.\n")
	}
	_,err := io.WriteString(w, "<br><a href='/captcha/page'>Try another one</a>")

	return err
}


const formTemplateSrc = `<!doctype html>
<head>
    <title>Captcha Example</title>
</head>
<body>
<script>

</script>
<form action="/captcha/verify" method=post>
    <p>输入验证码</p>
    <p><img id=image src="/captcha/image?captchid={{.CaptchaId}}" ></p>
    <input type=hidden name=captchaId value="{{.CaptchaId}}"><br>
    <input name=captchaSolution>
    <input type=submit value=Submit>
</form>`




















