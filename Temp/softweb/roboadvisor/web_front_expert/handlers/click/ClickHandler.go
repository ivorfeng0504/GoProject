package click

import (
	"git.emoney.cn/softweb/roboadvisor/agent"
	"git.emoney.cn/softweb/roboadvisor/contract"
	"github.com/devfeel/dotweb"
	"strings"
	"strconv"
)

// 记录点击量
func AddClick(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	identity := ctx.QueryString("identity")
	clickType := ctx.QueryString("clickType")
	extParam:=ctx.QueryString("extParam")
	if len(identity) == 0 {
		response.RetCode = -1
		response.RetMsg = "identity不能为空"
		return ctx.WriteJson(response)
	}
	if len(clickType) == 0 {
		response.RetCode = -2
		response.RetMsg = "clickType不能为空"
		return ctx.WriteJson(response)
	}

	if len(extParam) > 0 {
		extStrs := strings.Split(extParam, "|")
		if len(extStrs) > 0 {
			for _, v := range extStrs {
				if v != "" {
					extStrs2 := strings.Split(v, ",")
					if len(extStrs2) == 2 {
						agent.AddClick(extStrs2[1], extStrs2[0])
					}
				}
			}
		}
	}

	count, err := agent.AddClick(clickType, identity)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = count
	return ctx.WriteJson(response)
}

// QueryClick 查询点击量
func QueryClick(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	identity := ctx.QueryString("identity")
	clickType := ctx.QueryString("clickType")
	if len(identity) == 0 {
		response.RetCode = -1
		response.RetMsg = "identity不能为空"
		return ctx.WriteJson(response)
	}
	if len(clickType) == 0 {
		response.RetCode = -2
		response.RetMsg = "clickType不能为空"
		return ctx.WriteJson(response)
	}
	identitys := strings.Split(identity, ",")
	count, err := agent.QueryClick(clickType, identitys...)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = count
	return ctx.WriteJson(response)
}


// HandleClick 手动处理点击量
func HandleClick(ctx dotweb.Context) error {
	response := contract.NewResonseInfo()
	identity := ctx.QueryString("identity")
	clickType := ctx.QueryString("clickType")
	strclicknum := ctx.QueryString("clicknum")
	if len(identity) == 0 {
		response.RetCode = -1
		response.RetMsg = "identity不能为空"
		return ctx.WriteJson(response)
	}
	if len(clickType) == 0 {
		response.RetCode = -2
		response.RetMsg = "clickType不能为空"
		return ctx.WriteJson(response)
	}
	if len(strclicknum) == 0 {
		response.RetCode = -3
		response.RetMsg = "clickType不能为空"
		return ctx.WriteJson(response)
	}
	clicknum, _ := strconv.Atoi(strclicknum)
	count, err := agent.HandleClick(clickType, identity, clicknum)
	if err != nil {
		response.RetCode = -2
		response.RetMsg = err.Error()
		return ctx.WriteJson(response)
	}
	response.RetCode = 0
	response.RetMsg = "SUCCESS"
	response.Message = count
	return ctx.WriteJson(response)
}
