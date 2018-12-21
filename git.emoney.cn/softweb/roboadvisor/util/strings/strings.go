package _strings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Trim 去除两边的空格换行等空白符
func Trim(str string) string {
	str = strings.TrimSpace(str)
	return str
}

// Replace 替换字符串
func Replace(str, oldStr, newStr string) string {
	str = strings.Replace(str, oldStr, newStr, -1)
	return str
}

// SubString 截取字符串
func SubString(str string, start int, end int) string {
	if start <= 0 && end <= 0 {
		return str
	}
	if start <= 0 {
		str = string([]rune(str)[:end])
	} else if end <= 0 {
		str = string([]rune(str)[start:])
	} else {
		str = string([]rune(str)[start:end])
	}
	return str
}

// StringLen 计算字符串长度（一个中文算一个字符）
func StringLen(str string) int {
	length := utf8.RuneCountInString(str)
	return length
}

// LastString 获取最后一个字符
func LastString(str string) string {
	if len(str) <= 1 {
		return str
	}
	length := StringLen(str)
	str = SubString(str, length-1, -1)
	return str
}

// RoundStr 简单的四舍五入计算
func RoundStr(num string, digits int) (string, error) {
	return roundStrV1(num, digits)
}

// roundStrV1 简单的四舍五入计算
func roundStrV1(num string, digits int) (string, error) {
	if digits <= 0 {
		digits = 0
	}
	result := ""
	source := num
	if len(num) == 0 {
		return source, errors.New("无效的数字")
	}
	symbol := ""
	if num[:1] == "-" {
		symbol = "-"
		num = num[1:]
	}
	//整数
	if strings.Contains(num, ".") == false {
		integer, err := strconv.Atoi(num)
		if err != nil {
			return source, err
		} else {
			result = symbol + strconv.Itoa(integer)
			if digits > 0 {
				result = result + "." + strings.Repeat("0", digits)
			}
			return result, nil
		}
	}

	//带小数
	parts := strings.Split(num, ".")
	interger, err := strconv.Atoi(parts[0])
	if err != nil {
		return source, err
	}
	decimal := parts[1]
	//不要小数
	if digits == 0 {
		if decimal[0:1] > "4" {
			interger++
		}
		result = symbol + strconv.Itoa(interger)
		return result, nil
	}

	//保留digits位小数
	if len(decimal) <= digits {
		result = symbol + strconv.Itoa(interger) + "." + decimal + strings.Repeat("0", digits-len(decimal))
		return result, nil
	}
	digitsAfter, err := strconv.Atoi(decimal[digits : digits+1])
	if err != nil {
		return source, err
	}
	if digitsAfter <= 4 {
		//四舍
		result = symbol + strconv.Itoa(interger) + "." + decimal[0:digits]
	} else {
		//五入
		numSrc := decimal[0:digits]
		num, err := strconv.Atoi(numSrc)
		if err != nil {
			return source, err
		}
		//整数位是否要进1
		if len(strconv.Itoa(num+1)) > len(numSrc) {
			interger++
			num = 0
		} else {
			num++
		}
		numStr := strconv.Itoa(num)
		//处理前置的0
		preRepeatZero := ""
		if len(numSrc) > len(numStr) {
			preRepeatZero = strings.Repeat("0", len(numSrc)-len(numStr))
		}
		//处理后置的0
		sufRepeatZero := ""
		if digits-len(numStr)-len(preRepeatZero) > 0 {
			sufRepeatZero = strings.Repeat("0", digits-len(numStr)-len(preRepeatZero))
		}
		result = symbol + strconv.Itoa(interger) + "." + preRepeatZero + numStr + sufRepeatZero
	}
	return result, nil
}

// roundStrV2 四舍六入五成双
func roundStrV2(num string, digits int) (string, error) {
	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%."+strconv.Itoa(digits)+"f", f), nil
}

// RoundStr 简单的四舍五入计算 如果出现错误，返回原始字符串，不抛出错误
func RoundStrWithNoError(num string, digits int) string {
	result, err := RoundStr(num, digits)
	if err != nil {
		return num
	}
	return result
}

// StartWith 是否以指定字符开头
func StartWith(str string, start string) bool {
	if strings.Index(str, start) == 0 {
		return true
	} else {
		return false
	}
}
