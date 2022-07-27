package xerr

var message map[uint32]string

func init()  {
	message = make(map[uint32]string)

	message[OK] = "SUCCESS"
	message[ServerCommonError] = "服务器开小差啦，请稍后重试"
	message[RequestParamError] = "参数错误"
	message[TokenGenerateError] = "生成Token失败"
	message[TokenExpireError] = "token失效，请重新登陆"
	message[DBError] = "服务器失败，请稍后再试"

	message[ProjectNotExists] = "日志项目不存在"
	message[ReadQuotaExceed] = "超过读取日志限额"
	message[InternalServerError] = "服务器内部错误"
	message[ServerBusy] = "服务器正忙，请稍后再试"

}

func MapErrMsg(errCode uint32) string  {
	if msg, ok := message[errCode]; ok {
		return msg
	} else {
		return "服务器开小差啦，请稍后重试"
	}
}

func isCodeErr(errCode uint32) bool  {
	if _, ok := message[errCode]; ok {
		return true
	} else {
		return false
	}
}
