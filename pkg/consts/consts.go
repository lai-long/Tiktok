package consts

// 00通用
const (
	Success   int32 = 0
	FileError int32 = 100002 //FILE错误
	IOOsError int32 = 200003 //io os错误
)

// 01用户模块
const (
	UserNameExists     int32 = 101001 //用户名已存在
	UserPasswordError  int32 = 101002 //密码错误
	UserNotExists      int32 = 101003 //用户不存在
	MfaLack            int32 = 101004 //缺少mfa code
	MfaCodeFalse       int32 = 101005 //mfa code错误
	ImageFalse         int32 = 101006 //图片格式错误
	UserReqValidError  int32 = 101007 //用户绑定错误
	UserDBSelectError  int32 = 201001 //用户模块数据库select错误
	UserDBInsertError  int32 = 201002 //用户模块数据库insert错误
	UserDBUpdateError  int32 = 201003 //用户模块数据库update错误
	UserDBDeleteError  int32 = 201004 //用户模块数据库delete错误
	UserHashError      int32 = 201005 //加密密码失败
	GenerateTokenError int32 = 201006 //生成token错误
	UserRedisSetError  int32 = 201007 //用户模块redis set错误
	UserRedisGetError  int32 = 201008 //用户模块redis get错误
	UserRedisDelError  int32 = 201009 //用户模块redis del错误
	MfaGenerateError   int32 = 201010 //生成mfa错误
	MfaBindError       int32 = 201011 //mfa绑定错误
)

// 02视频
const (
	VideoReqValidError int32 = 102001 //视频模块参数绑定错误
	VideoRedisSetError int32 = 202001 //视频模块redis set错误
	VideoDBInsertError int32 = 202002 //视频模块模块数据库insert错误
	VideoDBUpdateError int32 = 202003 //视频模块模块数据库update错误
	VideoDBDeleteError int32 = 202004 //视频模块模块数据库delete错误
	VideoDBSelectError int32 = 202005 //视频模块模块数据库select错误
	VideoRedisGetError int32 = 202006 //视频模块redis get错误
)

// 03react模块
const (
	ReactReqValidError int32 = 103001 //互动模块视频绑定错误
	ReactReqValueError int32 = 103002 //互动模块传入参数错误
	ReactDBInsertError int32 = 203002 //互动模块模块数据库insert错误
	ReactDBUpdateError int32 = 203003 //互动模块模块数据库update错误
	ReactDBDeleteError int32 = 203004 //互动模块模块数据库delete错误
	ReactDBSelectError int32 = 203005 //互动模块模块数据库select错误
)

// 04social模块
const (
	SocialReqValidError int32 = 104001 //社交模块视频绑定错误
	SocialReqValueError int32 = 104002 //社交模块传入参数错误
	SocialDBInsertError int32 = 204002 //社交模块模块数据库insert错误
	SocialDBUpdateError int32 = 204003 //社交模块模块数据库update错误
	SocialDBDeleteError int32 = 204004 //社交模块模块数据库delete错误
	SocialDBSelectError int32 = 204005 //社交模块模块数据库select错误
)

var mapErrorMsg = map[int32]string{
	Success:             "success",
	FileError:           "FILE错误",
	IOOsError:           "io os错误",
	UserNameExists:      "用户名已存在",
	UserPasswordError:   "密码错误",
	UserNotExists:       "用户不存在",
	MfaLack:             "缺少mfa code",
	MfaCodeFalse:        "mfa code错误",
	ImageFalse:          "图片格式错误",
	UserReqValidError:   "用户绑定参数错误",
	UserDBSelectError:   "用户模块数据库select错误",
	UserDBInsertError:   "用户模块数据库insert错误",
	UserDBUpdateError:   "用户模块数据库update错误",
	UserDBDeleteError:   "用户模块数据库delete错误",
	UserHashError:       "加密密码失败",
	GenerateTokenError:  "生成token错误",
	UserRedisSetError:   "用户模块redis set错误",
	UserRedisGetError:   "用户模块redis get错误",
	UserRedisDelError:   "用户模块redis del错误",
	MfaGenerateError:    "生成mfa错误",
	MfaBindError:        "mfa绑定错误",
	VideoRedisSetError:  "视频模块redis set错误",
	VideoDBInsertError:  "视频模块模块数据库insert错误",
	VideoDBUpdateError:  "视频模块模块数据库update错误",
	VideoDBDeleteError:  "视频模块模块数据库delete错误",
	VideoDBSelectError:  "视频模块模块数据库select错误",
	VideoRedisGetError:  "视频模块redis get错误",
	VideoReqValidError:  "视频模块参数绑定错误",
	ReactReqValidError:  " //互动模块视频绑定错误",
	ReactReqValueError:  " //互动模块传入参数错误",
	ReactDBInsertError:  " //互动模块模块数据库insert错误",
	ReactDBUpdateError:  " //互动模块模块数据库update错误",
	ReactDBDeleteError:  " //互动模块模块数据库delete错误",
	ReactDBSelectError:  " //互动模块模块数据库select错误",
	SocialReqValidError: " //社交模块视频绑定错误",
	SocialReqValueError: " //社交模块传入参数错误",
	SocialDBInsertError: " //社交模块模块数据库insert错误",
	SocialDBUpdateError: " //社交模块模块数据库update错误",
	SocialDBDeleteError: " //社交模块模块数据库delete错误",
	SocialDBSelectError: " //社交模块模块数据库select错误",
}

func GetErrorCodeMsg(code int32) string {
	if msg, ok := mapErrorMsg[code]; ok {
		return msg
	}
	return `未知错误`
}
