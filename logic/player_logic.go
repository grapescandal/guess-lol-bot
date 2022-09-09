package logic

import "guess-lol-bot/model"

func CreatePlayer(name string, userID string) *model.Player {
	player := new(model.Player)
	player.Name = name
	player.Score = 0
	player.UserID = userID
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
