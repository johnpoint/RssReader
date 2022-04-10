package infra

import "RssReader/pkg/errorhelper"

var (
	// 请求异常 401xx
	ReqParseError = &errorhelper.Err{Code: 40100, Message: "请求参数异常"}
	DataBaseError = &errorhelper.Err{Code: 40200, Message: "数据库异常"}
	LoginError    = &errorhelper.Err{Code: 40300, Message: "登录失败"}
)
