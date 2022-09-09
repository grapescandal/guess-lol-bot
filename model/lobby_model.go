package model

type Lobby struct {
	Id        string
	ChannelID string
	Player    []*Player
}
