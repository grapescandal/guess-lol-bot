package logic

import (
	"encoding/json"
	"fmt"
	"guess-lol-bot/model"
	"io/ioutil"
	"os"
)

var itemList *model.ItemList

func ReadJsonItem() *model.ItemList {
	jsonFile, err := os.Open("itemList.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened itemList.json")
	// defer the closing of our jsonFile so that we can parse it later on

	byteValue, _ := ioutil.ReadAll(jsonFile)
	defer jsonFile.Close()

	itemList = new(model.ItemList)
	json.Unmarshal(byteValue, &itemList)
	return itemList
}

func GetItem(itemID int) model.Item {
	item := new(model.Item)

	for _, i := range itemList.ItemList {
		if i.Id == itemID {
			item = &i
		}
	}

	return *item
}

func UseItem(itemID int, user *model.Player, players []*model.Player) string {
	item := GetItem(itemID)

	switch item.Id {
	case 0:
		user.Score += 50
	case 1:
		user.OpeningCount = 2
	default:
	}

	return "Item used"
}
