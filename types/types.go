package types

type Login struct {
	ClientID int    `json:"clientID"`
	Username string `json:"username"`
}

type WSMessage struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PlayerState struct {
	Health   int      `json:"health"`
	Position Position `json:"position"`
}
