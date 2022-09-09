package logic

import (
	"fmt"
	"guess-lol-bot/helper"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const PREFIX = ".lol"

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
		lobby := CreateLobby(m.ChannelID)
		message := fmt.Sprintf("Lobby: %s was created", lobby.Id)
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	} else if command == "lobby" {
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
	} else if command == "join" {
		playerName := helper.FilterInput(m.Content, PREFIX+" "+"join")
		player := CreatePlayer(playerName, m.Author.ID)
		lobbyID := JoinLobby(m.ChannelID, player)

		message := fmt.Sprintf("Player: %v has joined to lobby: %v", playerName, lobbyID)
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	} else if command == "start" {
		message := ""
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

		players := GetPlayers(m.ChannelID)
		SetMaxTurn(len(players))

		StartGame()
		turn := GetTurn()

		playerName := players[turn].Name
		hint, length := GetHint()
		message += "Game Started!\n"

		message += fmt.Sprintf("Answer is %s, Length is %v", hint, length)
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}

		message1 := fmt.Sprintf("%v's turn", playerName)
		_, err = s.ChannelMessageSend(m.ChannelID, message1)
		if err != nil {
			fmt.Println(err)
		}
	} else if command == "answer" {

		message := ""
		if !isStart {
			message += "Game is not start yet"
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		turn := GetTurn()
		players := GetPlayers(m.ChannelID)

		if m.Author.ID != players[turn].UserID {
			message += "Not your turn"
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		answerFromUser := strings.ToLower(helper.FilterInput(m.Content, PREFIX+" "+"answer"))
		player := GetPlayer(m.ChannelID, m.Author.ID)
		result, success, status, answer := Answer(answerFromUser)
		if success {
			if result {
				player.Score += currentScore
				message += fmt.Sprintf("Player: %v win, Answer is %v \n", player.Name, answer.Name)
				message += "----------Scoreboard----------\n"
				sort.SliceStable(players, func(i, j int) bool {
					return players[i].Score > players[j].Score
				})
				for _, player := range players {
					message += fmt.Sprintf("%v : %v\n", player.Name, player.Score)
				}
				_, err := s.ChannelMessageSend(m.ChannelID, message)
				if err != nil {
					fmt.Println(err)
				}

				InitGame()
			} else {
				message += fmt.Sprintf("Try again, Answer is %v", status)
				_, err := s.ChannelMessageSend(m.ChannelID, message)
				if err != nil {
					fmt.Println(err)
				}
				NextTurn()
				turn = GetTurn()
				player := players[turn]
				message1 := fmt.Sprintf("%v's turn", player.Name)
				_, err = s.ChannelMessageSend(m.ChannelID, message1)
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
			}

		} else {
			message += fmt.Sprintf("%v", status)
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				fmt.Println(err)
			}
		}

	} else if command == "open" {
		message := ""

		if !isStart {
			message += "Game is not start yet"
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		turn := GetTurn()
		players := GetPlayers(m.ChannelID)

		if m.Author.ID != players[turn].UserID {
			message += "Not your turn"
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
		hintImage, err := GetPieceCardImage(index)
		if err != nil {
			message += err.Error()
			fmt.Printf("%s", err)
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

		DecreaseScore(1)
	} else if command == "pass" {
		message := ""

		if !isStart {
			message += "Game is not start yet"
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		turn := GetTurn()
		players := GetPlayers(m.ChannelID)

		if m.Author.ID != players[turn].UserID {
			message += "Not your turn"
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		NextTurn()
		turn = GetTurn()
		player := players[turn]
		message = fmt.Sprintf("%v's turn", player.Name)
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}
}
