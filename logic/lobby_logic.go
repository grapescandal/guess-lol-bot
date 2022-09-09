package logic

import (
	"guess-lol-bot/model"

	"github.com/google/uuid"
)

var lobbies map[string]*model.Lobby

func CreateLobby(channelID string) *model.Lobby {
	if lobbies == nil {
		lobbies = make(map[string]*model.Lobby)
	}
	lobby := new(model.Lobby)
	lobby.Id = uuid.New().String()
	lobby.ChannelID = channelID
	lobbies[channelID] = lobby
	return lobby
}

func GetLobby(channelID string) *model.Lobby {
	return lobbies[channelID]
}

func JoinLobby(channelID string, player *model.Player) string {
	lobby, found := lobbies[channelID]
	if !found {
		return "Lobby not found"
	} else {
		lobby.Player = append(lobby.Player, player)
		return lobby.Id
	}
}
