package consts

const (
	Success   int32 = 0
	FileError int32 = 100002 //FILE错误
	IOOsError int32 = 100003 //io os错误
)
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

var mapErrorMsg = map[int32]string{
	Success:            "success",
	FileError:          "FILE错误",
	IOOsError:          "io os错误",
	UserNameExists:     "用户名已存在",
	UserPasswordError:  "密码错误",
	UserNotExists:      "用户不存在",
	MfaLack:            "缺少mfa code",
	MfaCodeFalse:       "mfa code错误",
	ImageFalse:         "图片格式错误",
	UserReqValidError:  "用户绑定参数错误",
	UserDBSelectError:  "用户模块数据库select错误",
	UserDBInsertError:  "用户模块数据库insert错误",
	UserDBUpdateError:  "用户模块数据库update错误",
	UserDBDeleteError:  "用户模块数据库delete错误",
	UserHashError:      "加密密码失败",
	GenerateTokenError: "生成token错误",
	UserRedisSetError:  "用户模块redis set错误",
	UserRedisGetError:  "用户模块redis get错误",
	UserRedisDelError:  "用户模块redis del错误",
	MfaGenerateError:   "生成mfa错误",
	MfaBindError:       "mfa绑定错误",
}

func GetErrorCodeMsg(code int32) string {
	if msg, ok := mapErrorMsg[code]; ok {
		return msg
	}
	return `未知错误`
}
