package validate

import (
	"regexp"
)

const (
	regular = "^1[0-9]{10}$"
)


//func IsNilString(val string, errCode int, errMsg string) (*contract.ResponseInfo, error) {
//	if val != "" {
//		return contract.NewResonseInfo(), nil
//	} else {
//		return contract.CreateResponse(errCode, errMsg, nil), errors.New("val is nil")
//	}
//}

func IsMobile(mobileNum string) bool {
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}