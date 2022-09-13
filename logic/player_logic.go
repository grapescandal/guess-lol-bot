package logic

import "guess-lol-bot/model"

func CreatePlayer(name string, userID string) *model.Player {
	player := new(model.Player)
	player.Name = name
	player.Score = 0
	player.UserID = userID
	player.OpeningCount = 1
	return player
}

func GetPlayer(chanelID string, userID string) *model.Player {
	var player *model.Player
	for _, p := range lobbies[chanelID].Player {
		if p.UserID == userID {
			player = p
		}
	}

	return player
}

func GetPlayers(chanelID string) []*model.Player {
	return lobbies[chanelID].Player
}

func RemovePlayer(chanelID string, userID string) (bool, string) {
	removedPlayerName := ""
	players := GetPlayers(chanelID)
	found := false
	for i := 0; i < len(players); i++ {
		player := players[i]
		if player.UserID == userID {
			removedPlayerName = player.Name
			players = append(players[:i], players[i+1:]...)
			i-- // Important: decrease index
			found = true
			break
		}
	}

	if !found {
		return found, ""
	}

	lobby := GetLobby(chanelID)
	lobby.Player = players
	return found, removedPlayerName
}
