package chat

/*
chat 中事件的关键字，以 Server 开头的事件为服务端发送的事件，以 client 开头的事件为客户端发送的事件
*/

// 服务端推送的事件
const (
	ServerSendMsg            = "send_msg"             // 推送消息
	ServerReadMsg            = "read_msg"             // 已读消息
	ServerAccountLogin       = "account_login"        // 账户上线
	ServerAccountLogout      = "account_logout"       // 账户离线
	ServerUpdateAccount      = "update_account"       // 更新账户信息
	ServerApplication        = "application"          // 好友申请
	ServerDeleteRelation     = "delete_relation"      // 删除关系
	ServerUpdateNickName     = "update_nickname"      // 更新昵称
	ServerUpdateSettingState = "update_setting_state" // 更新关系状态
	ServerUpdateEmail        = "update_email"         // 更新邮箱
	ServerUpdateMsgState     = "update_msg_state"     // 更新消息状态
	ServerGroupTransferred   = "group_transferred"    // 群被转让
	ServerGroupDissolved     = "group_dissolved"      // 群被解散
	ServerInviteAccount      = "group_new_account"    // 新人进群
	ServerQuitGroup          = "group_quit_account"   // 新人进群
	ServerCreateNotify       = "create_notify"        // 创建群通知
	ServerUpdateNotify       = "update_notify"        // 更新群通知
	ServerError              = "error"                // 错误
)

// 客户端推送的事件
const (
	ClientSendMsg = "send_msg" // 发送消息
	ClientReadMsg = "read_msg" // 已读消息
	ClientTest    = "test"     // 测试
	ClientAuth    = "auth"     // 认证
)
