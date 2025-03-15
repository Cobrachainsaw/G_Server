package types

type Login struct {
	ClientID int    `json:"clientID"`
	Username string `json:"username"`
}

type WSMessage struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}
