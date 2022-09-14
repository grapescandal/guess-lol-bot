package logic

import (
	"encoding/json"
	"flag"
	"fmt"
	"guess-lol-bot/model"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"
)

var itemList *model.ItemList
var fileName = flag.String("itemlist", "itemlist.json", "Location of itemlist file")

func ReadJsonItem() {
	file, err := os.Open(*fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	fmt.Println("Successfully Opened itemlist.json")

	jsonFile, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	itemList = new(model.ItemList)
	json.Unmarshal(jsonFile, &itemList)
}

func GetItem(itemID int) model.Item {
	item := new(model.Item)

	for i, j := range itemList.ItemList {
		if j.Id == itemID {
			item = &itemList.ItemList[i]
			break
		}
	}

	return *item
}

func GetItemRank(r model.Rank) string {
	switch r {
	case model.Bronze:
		return "Bronze"
	case model.Silver:
		return "Silver"
	case model.Gold:
		return "Gold"
	case model.Platinum:
		return "Platinum"
	case model.Diamond:
		return "Diamond"
	case model.Challenger:
		return "Challenger"
	}
	return "Iron"
}

func UseItem(s *discordgo.Session, m *discordgo.MessageCreate, item *model.Item, user *model.Player, players []*model.Player) {
	result := fmt.Sprintf("Item used!\n%s\n", item.Description)
	_, err := s.ChannelMessageSend(m.ChannelID, result)
	if err != nil {
		fmt.Println(err)
	}

	message := ""

	switch item.Name {
	case "Renewal Tunic":
		user.Score += 20
		message = fmt.Sprintf("%v: %v", user.Name, user.Score)
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}

		scoreboard := "----------Scoreboard----------\n"
		sort.SliceStable(players, func(i, j int) bool {
			return players[i].Score > players[j].Score
		})
		for _, player := range players {
			scoreboard += fmt.Sprintf("%v : %v\n", player.Name, player.Score)
		}

		_, err = s.ChannelMessageSend(m.ChannelID, scoreboard)
		if err != nil {
			fmt.Println(err)
		}
	case "Sword of the Divine":
		user.OpeningCount = 3
	case "Ionic Spark":
		for i, p := range players {
			if p.UserID != user.UserID {
				players[i].Score -= 30
				if players[i].Score <= 0 {
					players[i].Score = 0
				}
			}
		}

		scoreboard := "----------Scoreboard----------\n"
		sort.SliceStable(players, func(i, j int) bool {
			return players[i].Score > players[j].Score
		})
		for _, player := range players {
			scoreboard += fmt.Sprintf("%v : %v\n", player.Name, player.Score)
		}

		_, err := s.ChannelMessageSend(m.ChannelID, scoreboard)
		if err != nil {
			fmt.Println(err)
		}

	case "Stealth Ward":
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		index := r1.Intn(col * row)

		hintImage, err := OpenPieceImage(index)
		if err != nil {
			result += err.Error()
			fmt.Printf("%s\n", err)
			_, err := s.ChannelMessageSend(m.ChannelID, result)
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
	case "Vision Ward":
		if len(openPieces) == 0 {
			message = "Nothing happen"
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		index := openPieces[r1.Intn(len(openPieces))]
		fmt.Println(index)
		hintImage, err := ClosePieceImage(index)
		if err != nil {
			message = err.Error()
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
		IncreaseScore(1)
	case "Rod of Ages":
		hint := GetFirstAndLastAlphabet()
		message = fmt.Sprintf("First and last alphabet is %s", hint)
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	case "Will of the Ancients":
		remainingAlphabets := 0
		for _, a := range hint {
			if string(a) == "-" {
				remainingAlphabets++
			}
		}
		additionalScore += (remainingAlphabets * 10)
	case "Morellonomicon":
		currentScore = currentScore / 2
	case "Guardian Angel":
		user.AnswerCount = 2
	case "Deathfire Grasp":
		SkipItemPhase()
		_, turnMessage := NextTurn(m.ChannelID)
		_, err := s.ChannelMessageSend(m.ChannelID, turnMessage)
		if err != nil {
			fmt.Println(err)
		}
	case "Zhonya's Hourglass":
		SkipItemPhase()
		_, turnMessage := NextTurn(m.ChannelID)
		_, err := s.ChannelMessageSend(m.ChannelID, turnMessage)
		if err != nil {
			fmt.Println(err)
		}
	default:
		return
	}
}
