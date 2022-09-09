package logic

import (
	"fmt"
	"guess-lol-bot/api"
	"guess-lol-bot/model"
	"image"
	"image/jpeg"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/oliamb/cutter"
)

var answer model.Answer
var isStart bool
var openPieces []int
var turn int = 0
var maxTurn int = 0
var currentScore int = 10
var maxScore int = 10
var pieceScore int = 9
var championData *model.ChampionData

func PrepareChampionData() {
	championResponse, err := api.GetChampionResponse()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	championData = api.GetChampionData(championResponse)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}

func InitGame() {
	isStart = false
	openPieces = []int{}
	currentScore = maxScore
}

func StartGame() {
	if !isStart {
		isStart = true
		championName := GetRamdomChampion()
		skin := GetRandomSkin(championName)
		skinNum := strconv.Itoa(skin.Num)
		GetSkinImage(championName + "_" + skinNum + ".jpg")
		answer = model.Answer{
			Name: skin.Name,
		}
		turn = 0
		lengthCounter := 0
		for _, a := range answer.Name {
			if isAlphabets(a) {
				lengthCounter++
			}
		}
		maxScore = lengthCounter + pieceScore
		currentScore = maxScore
	}
}

func GetTurn() int {
	return turn
}

func SetMaxTurn(number int) {
	maxTurn = number
}

func NextTurn() {
	turn += 1
	if turn > maxTurn {
		turn = 0
	}
}

func GetRamdomChampion() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	randomNumber := r1.Intn(len(championData.Name))
	champion := championData.Name[randomNumber]
	fmt.Printf("Pick champion Name: %v\n", champion)
	return champion
}

func GetRandomSkin(championName string) model.Skins {
	championDataResponse, err := api.GetChampionDataReponse(championName)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	randomNumber := r1.Intn(len(championDataResponse.Data[championName].Skins))
	skin := championDataResponse.Data[championName].Skins[randomNumber]
	fmt.Printf("Pick skin Name: %v\n", skin.Name)
	return skin
}

func GetSkinImage(url string) {
	_, err := api.GetSkinImage(url)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	_, err = os.Open("image.jpg")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}

func ReadCardImage() *os.File {
	file, err := os.Open("image.jpg")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	return file
}

func GetPieceCardImage(index int) (*os.File, error) {

	isAlreadyOpen := false
	for _, i := range openPieces {
		if index == i {
			isAlreadyOpen = true
			break
		}
	}

	if isAlreadyOpen {
		err := fmt.Errorf("%v is already open", index)
		return nil, err
	}

	file, err := os.Open("image.jpg")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Failed to decode: %v", err)
	}

	col := 8
	indexY := 0
	indexX := 0
	actualIndex := index - 1
	if index > col {
		indexY = actualIndex / col
		indexX = actualIndex % col
	} else {
		indexX = actualIndex
	}
	width := 152
	height := 90
	x := width * (indexX)
	y := height * (indexY)
	croppedImg, err := cutter.Crop(img, cutter.Config{
		Width:  width,
		Height: height,
		Anchor: image.Point{x, y},
		Mode:   cutter.TopLeft,
	})
	fmt.Printf("x: %v, y: %v\n", x, y)
	if err != nil {
		fmt.Printf("Error croppedImg: %v", err)
	}

	f, err := os.Create("piece.jpg")
	if err != nil {
		fmt.Printf("Failed to create: %v", err)
	}
	defer f.Close()

	err = jpeg.Encode(f, croppedImg, nil)
	if err != nil {
		fmt.Printf("Failed to encode: %v", err)
	}

	finalFile, err := os.Open("piece.jpg")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	openPieces = append(openPieces, index)
	return finalFile, nil
}

func DecreaseScore() {
	currentScore -= 1
}

func GetHint() (string, int) {
	hint := ""
	for _, a := range answer.Name {
		isAlphabets := isAlphabets(a)
		if isAlphabets {
			hint += "-"
		} else {
			hint += string(a)
		}
	}

	return hint, len(answer.Name)
}

func isAlphabets(c rune) bool {
	return unicode.IsLetter(c)
}

func Answer(message string) (bool, bool, string, *model.Answer) {
	answerLower := strings.ToLower(answer.Name)
	if len(message) != len(answer.Name) {
		return false, false, fmt.Sprintf("Please be sure your answer length is %v", len(answer.Name)), nil
	}
	if message == answerLower {
		return true, true, "", &answer
	} else {
		hint := ""
		for i, a := range answerLower {

			if string(message[i]) == string(a) {
				hint += string(answer.Name[i])
			} else {
				hint += "-"
			}
		}
		return false, true, hint, nil
	}
}
