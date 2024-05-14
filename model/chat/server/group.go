package server

type TransferGroup struct {
	EnToken   string `json:"en_token,omitempty"`
	AccountID int64  `json:"account_id,omitempty"`
}

type DissolveGroup struct {
	EnToken    string `json:"en_token,omitempty"`
	RelationID int64  `json:"relation_id,omitempty"`
}

type InviteGroup struct {
	EnToken   string `json:"en_token,omitempty"`
	AccountID int64  `json:"account_id,omitempty"`
}

type QuitGroup struct {
	EnToken   string `json:"en_token,omitempty"`
	AccountID int64  `json:"account_id,omitempty"`
}
