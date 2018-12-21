package service

import (
	"git.emoney.cn/softweb/roboadvisor/protected/service"
	"git.emoney.cn/softweb/roboadvisor/config"
	"strings"
	"net/http"
	"io/ioutil"
	"strconv"
)

type YqqLiveService struct {
	service.BaseService
}

// ModifyHeadportrait 更新用户头像
func (service *YqqLiveService) GetLiveRoomList()(datastr string,err error)  {
	apiUrl := config.CurrentConfig.MarketPriceApi
	if len(apiUrl) == 0 {
		shareMarketPriceServiceLogger.ErrorFormat(err, "行情接口地址配置不正确 configkey=MarketPriceApi")
		//return price, tradeDate, errors.New("行情接口地址配置不正确")
	}
	resp, err := http.Get(apiUrl)
	defer func() { resp.Body.Close() }()
	if err != nil {
		return "error",nil
	}
	if resp.StatusCode != http.StatusOK {
		return "0",nil
	}
	data, err := ioutil.ReadAll(resp.Body)


	return string(data),err
}