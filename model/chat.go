package model

type HandleSendMsg struct {
	AccessToken string     // AccessToken
	RelationID  int64      // 关系 ID
	AccountID   int64      // 账户 ID
	MsgContent  string     // 消息内容
	MsgExtend   *MsgExtend // 消息扩展信息
	RlyMsgID    int64      // 回复消息 ID（如果是回复消息，则此字段大于 0）
}

type HandleReadMsg struct {
	AccessToken string  // AccessToken
	MsgIDs      []int64 // 消息 IDs
	RelationID  int64   // 这些消息所处的关系 ID
	ReaderID    int64   // 读者账号 ID
}
