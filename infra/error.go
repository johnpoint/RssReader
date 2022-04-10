package infra

import "RssReader/pkg/errorhelper"

var (
	// 请求异常 401xx
	ReqParseError = &errorhelper.Err{Code: 40101, Message: "请求参数异常"}
	DataBaseError = &errorhelper.Err{Code: 40102, Message: "数据库异常"}
	//LoginError    = &errorhelper.Err{Code: 40103, Message: "登录失败"}
	ImportError = &errorhelper.Err{Code: 40104, Message: "导入异常"}
	ExportError = &errorhelper.Err{Code: 40105, Message: "导出异常"}
)
