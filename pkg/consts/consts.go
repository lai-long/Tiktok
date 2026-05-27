// Package consts 放置错误码以及错误信息等常量
package consts

// 通用动作/类型常量
const (
	// UserIDKey 和 UsernameKey 用于在请求上下文中存储用户信息
	UserIDKey   = "userid"
	UsernameKey = "username"

	// Action types (社交模块)
	ActionFollow   = "0"
	ActionUnfollow = "1"

	// Target types (互动模块)
	TargetVideo   = "1"
	TargetComment = "2"

	// Like/Dislike actions
	ActionLike    = "1"
	ActionDislike = "2"

	// Message types (WebSocket)
	MsgTypePrivate  = "1" // 私信
	MsgTypeOffline  = "2" // 未在线时搜到的消息
	MsgTypeHistory  = "3" // 历史消息请求
	MsgTypeGroupMsg = "4" // 群组消息
)

// 00通用
const (
	Success       int32 = 0
	FileError     int32 = 100002 // FILE错误
	IOOsError     int32 = 200003 // io os错误
	SentinelBlock int32 = 100003 // sentinel限流
)

// 01 用户模块
const (
	UserNameExists     int32 = 101001 // 用户名已存在
	UserPasswordError  int32 = 101002 // 密码错误
	UserNotExists      int32 = 101003 // 用户不存在
	MfaLack            int32 = 101004 // 缺少mfa code
	MfaCodeFalse       int32 = 101005 // mfa code错误
	ImageFalse         int32 = 101006 // 图片格式错误
	UserReqValidError  int32 = 101007 // 用户绑定错误
	UserDBSelectError  int32 = 201001 // 用户模块数据库select错误
	UserDBInsertError  int32 = 201002 // 用户模块数据库insert错误
	UserDBUpdateError  int32 = 201003 // 用户模块数据库update错误
	UserDBDeleteError  int32 = 201004 // 用户模块数据库delete错误
	UserHashError      int32 = 201005 // 加密密码失败
	GenerateTokenError int32 = 201006 // 生成token错误
	UserRedisSetError  int32 = 201007 // 用户模块redis set错误
	UserRedisGetError  int32 = 201008 // 用户模块redis get错误
	UserRedisDelError  int32 = 201009 // 用户模块redis del错误
	MfaGenerateError   int32 = 201010 // 生成mfa错误
	MfaBindError       int32 = 201011 // mfa绑定错误
)

// 02 视频
const (
	VideoReqValidError int32 = 102001 // 视频模块参数绑定错误
	VideoRedisSetError int32 = 202001 // 视频模块redis set错误
	VideoDBInsertError int32 = 202002 // 视频模块模块数据库insert错误
	VideoDBUpdateError int32 = 202003 // 视频模块模块数据库update错误
	VideoDBDeleteError int32 = 202004 // 视频模块模块数据库delete错误
	VideoDBSelectError int32 = 202005 // 视频模块模块数据库select错误
	VideoRedisGetError int32 = 202006 // 视频模块redis get错误
)

// 03 react模块
const (
	ReactReqValidError int32 = 103001 // 互动模块视频绑定错误
	ReactReqValueError int32 = 103002 // 互动模块传入参数错误
	ReactError         int32 = 103003 // nil case error
	ReactDBInsertError int32 = 203002 // 互动模块模块数据库insert错误
	ReactDBUpdateError int32 = 203003 // 互动模块模块数据库update错误
	ReactDBDeleteError int32 = 203004 // 互动模块模块数据库delete错误
	ReactDBSelectError int32 = 203005 // 互动模块模块数据库select错误
)

// 04social模块
const (
	SocialReqValidError int32 = 104001 // 社交模块视频绑定错误
	SocialReqValueError int32 = 104002 // 社交模块传入参数错误
	SocialDBInsertError int32 = 204002 // 社交模块模块数据库insert错误
	SocialDBUpdateError int32 = 204003 // 社交模块模块数据库update错误
	SocialDBDeleteError int32 = 204004 // 社交模块模块数据库delete错误
	SocialDBSelectError int32 = 204005 // 社交模块模块数据库select错误
)

// 05 websocket模块
const (
	WsReqValidError   int32 = 105001 // websocket参数错误
	WsClientNotOnline int32 = 105002 // 对方不在线
	WsClientOnline    int32 = 105003 // 对方在线
	WsDisconnect      int32 = 105004 // 连接中断
	WsConnectSuccess  int32 = 105005 // 连接成功
	WsAIReplyEmpty    int32 = 105006 // AI无响应
	WsGetOfflineError int32 = 205001 // 获取离线消息错误
	WsGetHistoryError int32 = 205002 // 获取历史消息错误
)

// 06 mfa模块
const (
	MfaReqValidError int32 = 106001 // mfa请求参数错误
	MfaDBSelectError int32 = 206001 // mfa模块数据库select错误
	MfaDBInsertError int32 = 206002 // mfa模块数据库insert错误
	MfaDBUpdateError int32 = 206003 // mfa模块数据库update错误
	MfaDBDeleteError int32 = 206004 // mfa模块数据库delete错误
)

var mapErrorMsg = map[int32]string{
	Success:             "success",
	FileError:           "文件处理失败",
	IOOsError:           "文件写入失败",
	SentinelBlock:       "请求过于频繁，请稍后再试",
	UserNameExists:      "用户名已存在",
	UserPasswordError:   "密码错误",
	UserNotExists:       "用户不存在",
	MfaLack:             "请输入验证码",
	MfaCodeFalse:        "验证码错误",
	ImageFalse:          "图片格式不支持",
	UserReqValidError:   "用户请求参数错误",
	UserDBSelectError:   "查询用户信息失败",
	UserDBInsertError:   "用户注册失败",
	UserDBUpdateError:   "更新用户信息失败",
	UserDBDeleteError:   "删除用户失败",
	UserHashError:       "密码处理失败",
	GenerateTokenError:  "登录凭证生成失败",
	UserRedisSetError:   "用户信息缓存失败，请重试",
	UserRedisGetError:   "用户信息读取失败，请重试",
	UserRedisDelError:   "用户信息更新失败，请重试",
	MfaGenerateError:    "生成验证信息失败",
	MfaBindError:        "验证器绑定失败",
	VideoRedisSetError:  "视频信息缓存失败，请重试",
	VideoDBInsertError:  "视频发布失败",
	VideoDBUpdateError:  "视频更新失败",
	VideoDBDeleteError:  "视频删除失败",
	VideoDBSelectError:  "视频查询失败",
	VideoRedisGetError:  "视频信息读取失败，请重试",
	VideoReqValidError:  "视频请求参数错误",
	ReactReqValidError:  "互动请求参数错误",
	ReactReqValueError:  "互动参数值无效",
	ReactError:          "互动异常，请重试",
	ReactDBInsertError:  "互动失败",
	ReactDBUpdateError:  "互动更新失败",
	ReactDBDeleteError:  "取消互动失败",
	ReactDBSelectError:  "互动查询失败",
	SocialReqValidError: "社交请求参数错误",
	SocialReqValueError: "社交参数值无效",
	SocialDBInsertError: "关注失败",
	SocialDBUpdateError: "关注状态更新失败",
	SocialDBDeleteError: "取消关注失败",
	SocialDBSelectError: "社交关系查询失败",
	MfaReqValidError:    "验证请求参数错误",
	MfaDBSelectError:    "验证信息查询失败",
	MfaDBInsertError:    "验证信息保存失败",
	MfaDBUpdateError:    "验证信息更新失败",
	MfaDBDeleteError:    "验证信息删除失败",
	WsReqValidError:     "消息请求参数错误",
	WsClientNotOnline:   "对方不在线",
	WsClientOnline:      "对方在线",
	WsDisconnect:        "连接中断",
	WsConnectSuccess:    "连接成功",
	WsAIReplyEmpty:      "AI暂时无法回复",
	WsGetOfflineError:   "获取离线消息失败",
	WsGetHistoryError:   "获取历史消息失败",
}

// GetErrorCodeMsg 根据code获取对应错误信息
func GetErrorCodeMsg(code int32) string {
	if msg, ok := mapErrorMsg[code]; ok {
		return msg
	}
	return `未知错误`
}
