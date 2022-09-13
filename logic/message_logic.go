package logic

import (
	"fmt"
	"guess-lol-bot/helper"
	"guess-lol-bot/model"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const PREFIX = ".lol"

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	content := m.Content
	if len(content) <= len(PREFIX) {
		return
	}
	if content[:len(PREFIX)] != PREFIX {
		return
	}
	content = content[len(PREFIX):]
	if len(content) < 1 {
		return
	}
	args := strings.Fields(content)
	command := strings.ToLower(args[0])

	if command == "create" {
		CreateCommand(s, m)
	} else if command == "lobby" {
		LobbyCommand(s, m)
	} else if command == "join" {
		JoinCommand(s, m)
	} else if command == "start" {
		StartCommand(s, m)
	} else if command == "answer" {
		AnswerCommand(s, m)
	} else if command == "open" {
		OpenCommand(s, m)
	} else if command == "pass" {
		PassCommand(s, m)
	} else if command == "leave" {
		LeaveCommand(s, m)
	} else if command == "item" {
		ItemCommand(s, m)
	}
}

func CreateCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	lobby := CreateLobby(m.ChannelID)
	message := fmt.Sprintf("Lobby: %s was created", lobby.Id)
	_, err := s.ChannelMessageSend(m.ChannelID, message)
	if err != nil {
		fmt.Println(err)
	}
}

func LobbyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	lobby := GetLobby(m.ChannelID)
	message := ""
	if lobby == nil {
		message += "Lobby not found"
	} else {
		message += fmt.Sprintf("Lobby:%s \n", lobby.Id)
	}

	_, err := s.ChannelMessageSend(m.ChannelID, message)
	if err != nil {
		fmt.Println(err)
	}

}

func JoinCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	playerName := helper.FilterInput(m.Content, PREFIX+" "+"join")
	player := CreatePlayer(playerName, m.Author.ID)
	lobbyID := JoinLobby(m.ChannelID, player)

	message := fmt.Sprintf("Player: %v has joined to lobby: %v", playerName, lobbyID)
	_, err := s.ChannelMessageSend(m.ChannelID, message)
	if err != nil {
		fmt.Println(err)
	}

}

func StartCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	message := ""
	if !isStart {
		lobby := GetLobby(m.ChannelID)
		if lobby == nil {
			message += "Please create lobby first"
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				fmt.Println(err)
			}
			return
		} else if lobby != nil && len(lobby.Player) == 0 {
			message += "Please join lobby first"
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		message = "Game Started!\n"
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	players := GetPlayers(m.ChannelID)
	SetMaxTurn(len(players))

	InitGame()
	_, turnMessage := NextTurn(m.ChannelID, 1)

	StartGame(m.ChannelID)

	hint, length := GetHint()

	message = fmt.Sprintf("Answer is %s, \n%v\n", hint, length)
	message += turnMessage
	_, err := s.ChannelMessageSend(m.ChannelID, message)
	if err != nil {
		fmt.Println(err)
	}
}

func AnswerCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	message := ""
	if !isStart {
		message += "Game is not start yet"
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	players := GetPlayers(m.ChannelID)
	_, player := GetTurn(m.ChannelID)

	if m.Author.ID != player.UserID {
		message += "Not your turn"
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if isAnswer {
		message += "Can't use answer command anymore"
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	answerFromUser := strings.ToLower(helper.FilterInput(m.Content, PREFIX+" "+"answer"))
	result, success, status, answer := Answer(answerFromUser)
	if success {
		SkipItemPhase()
		IncreaseAnswerCount(player)

		if result {
			player.Score += currentScore + additionalScore
			message += fmt.Sprintf("Player: %v win, Answer is %v \n", player.Name, answer.Name)
			scoreboard := EndRound(players)
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				fmt.Println(err)
			}

			_, err = s.ChannelMessageSend(m.ChannelID, scoreboard)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			message += fmt.Sprintf("Try again, Answer is %v", status)
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				fmt.Println(err)
			}
			_, turnMessage := NextTurn(m.ChannelID, 1)
			_, err = s.ChannelMessageSend(m.ChannelID, turnMessage)
			if err != nil {
				fmt.Println(err)
			}
		}

		if result {
			cardImage := ReadSkinImage("image.jpg")
			_, err := s.ChannelFileSend(m.ChannelID, "image.jpg", cardImage)
			if err != nil {
				fmt.Println(err)
			}

			canNextRound := NextRound()
			if canNextRound {
				StartCommand(s, m)
			} else {
				scoreboard := EndGame(players)
				_, err := s.ChannelMessageSend(m.ChannelID, scoreboard)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	} else {
		message += fmt.Sprintf("%v", status)
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func OpenCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	message := ""

	if !isStart {
		message += "Game is not start yet"
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	_, player := GetTurn(m.ChannelID)

	if m.Author.ID != player.UserID {
		message += "Not your turn"
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if isOpenPiece {
		message += "Can't use open command anymore"
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	openPiece := strings.ToLower(helper.FilterInput(m.Content, PREFIX+" "+"open"))
	index, err := strconv.Atoi(openPiece)
	if err != nil {
		message += "Please input only 1-64"
		fmt.Printf("Failed to convert openPiece: %v", err)
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	if index <= 0 || index > 64 {
		message += "Please input only 1-64"
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	hintImage, err := OpenPieceImage(index)
	if err != nil {
		message += err.Error()
		fmt.Printf("%s\n", err)
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	_, err = s.ChannelFileSend(m.ChannelID, "card.jpg", hintImage)
	if err != nil {
		fmt.Println(err)
	}
	defer hintImage.Close()

	SkipItemPhase()
	DecreaseScore(1)
	IncreaseOpeningCount(player)
}

func PassCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	message := ""

	if !isStart {
		message += "Game is not start yet"
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	_, player := GetTurn(m.ChannelID)

	if m.Author.ID != player.UserID {
		message += "Not your turn"
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	SkipItemPhase()
	_, turnMessage := NextTurn(m.ChannelID, 1)
	_, err := s.ChannelMessageSend(m.ChannelID, turnMessage)
	if err != nil {
		fmt.Println(err)
	}
}

func LeaveCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	message := ""
	found, playerName := RemovePlayer(m.ChannelID, m.Author.ID)
	if found {
		message = fmt.Sprintf("%v leave", playerName)
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		message = "Player not found"
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func ItemCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	message := ""

	if !isStart {
		message += "Game is not start yet"
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	players := GetPlayers(m.ChannelID)
	_, player := GetTurn(m.ChannelID)

	if m.Author.ID != player.UserID {
		message += "Not your turn"
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if !isItemPhase {
		message += "Not item phase"
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	randomNumber := RandomItem()
	item := GetItem(randomNumber)
	itemRank := GetItemRank(model.Rank(item.Rank))
	description := fmt.Sprintf("Rank: %s\nDescription: %s", itemRank, item.Description)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       item.Name,
		Description: description,
		Image: &discordgo.MessageEmbedImage{
			URL:    item.ImagePath,
			Width:  64,
			Height: 64,
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	UseItem(s, m, &item, player, players)
	SkipItemPhase()
}
