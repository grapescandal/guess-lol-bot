package logic

import (
	"fmt"
	"guess-lol-bot/api"
	"guess-lol-bot/model"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/oliamb/cutter"
)

var answer model.Answer
var isStart bool
var isItemPhase bool
var isOpenPiece bool
var openingCount int = 0
var isAnswer bool
var answerCount int = 0
var remainingPieces map[int]bool
var openPieces []int

var currentRound int = 0
var maxRound int = 5
var startTurn int = 0
var turn int = 0
var maxTurn int = 0
var currentScore int = 64
var maxScore int = 64
var additionalScore int = 0
var pieceScore int = 64
var championData *model.ChampionData
var hint string = ""
var row int = 8
var col int = 8

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
	currentScore = maxScore
	hint = ""
}

func StartGame(channelID string) string {
	if isStart {
		return "Game not start yet"
	}

	isOpenPiece = false
	isItemPhase = true
	isAnswer = false
	championName := GetRamdomChampion()
	skin := GetRandomSkin(championName)

	skinNum := strconv.Itoa(skin.Num)
	GetSkinImage(championName + "_" + skinNum + ".jpg")

	if skin.Name == "default" {
		skin.Name = skin.Name + " " + championName
	}
	answer = model.Answer{
		Name: skin.Name,
	}

	CreatePuzzleImage()

	isStart = true

	lengthCounter := 0
	for _, a := range answer.Name {
		if isAlphabets(a) {
			lengthCounter++
		}
	}
	maxScore = lengthCounter + pieceScore
	currentScore = maxScore

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	randomNumber := r1.Intn(maxTurn)
	startTurn = randomNumber
	if randomNumber == startTurn {
		turn += 1
		if turn >= maxTurn {
			turn = 0
		}
	} else {
		turn = randomNumber
	}

	players := GetPlayers(channelID)
	player := players[turn]
	message := fmt.Sprintf("%v's turn\n", player.Name)
	return message

}

func SkipItemPhase() {
	isItemPhase = false
}

func GetTurn(channelID string) (int, *model.Player) {
	players := GetPlayers(channelID)
	player := players[turn]
	return turn, player
}

func SetMaxTurn(number int) {
	maxTurn = number
}

func NextTurn(channelID string) (int, string) {
	players := GetPlayers(channelID)
	turn += 1
	if turn >= maxTurn {
		turnIndex := turn - maxTurn
		turn = turnIndex
	}

	isOpenPiece = false
	isItemPhase = true
	isAnswer = false
	openingCount = 0
	additionalScore = 0
	player := players[turn]
	message := fmt.Sprintf("%v's turn\n", player.Name)
	return turn, message
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

func ReadSkinImage(fileName string) *os.File {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	return file
}

func CreatePuzzleImage() {
	remainingPieces = make(map[int]bool)
	openPieces = []int{}
	piecesLength := row * col
	for i := 0; i < piecesLength; i++ {
		remainingPieces[i] = false
	}

	width := 1215
	height := 717

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	black := color.Black

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, black)
		}
	}

	file, err := os.Create("puzzle.jpg")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	jpeg.Encode(file, img, nil)
	defer file.Close()
}

func UpdatePuzzleImage(inputX int, inputY int, inputWidth int, inputHeight int, piece image.Image, open bool) *os.File {

	file := ReadSkinImage("puzzle.jpg")
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	width := 1215
	height := 717

	imgRGBA := imageToRGBA(img)

	if open {
		// Set color for each pixel.
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				if x >= inputX && x < inputX+inputWidth && y >= inputY && y < inputY+inputHeight {
					pieceX := x
					pieceY := y
					imgRGBA.Set(x, y, piece.At(pieceX, pieceY))
				}
			}
		}
	} else {
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				if x >= inputX && x < inputX+inputWidth && y >= inputY && y < inputY+inputHeight {
					imgRGBA.Set(x, y, color.Black)
				}
			}
		}
	}

	finalFile, err := os.Create("puzzle.jpg")
	if err != nil {
		log.Fatalf("failed to create: %s", err)
	}
	jpeg.Encode(finalFile, imgRGBA, nil)
	defer finalFile.Close()

	finalFile, err = os.Open("puzzle.jpg")
	if err != nil {
		log.Fatalf("failed to create: %s", err)
	}
	return finalFile
}

func imageToRGBA(src image.Image) *image.RGBA {

	// No conversion needed if image is an *image.RGBA.
	if dst, ok := src.(*image.RGBA); ok {
		return dst
	}

	// Use the image/draw package to convert to *image.RGBA.
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}

func OpenPieceImage(index int) (*os.File, error) {
	index = index - 1

	isAlreadyOpen := false
	for key, value := range remainingPieces {
		if index == key {
			if value {
				isAlreadyOpen = true
				break
			}
		}
	}

	displayIndex := index + 1
	if isAlreadyOpen {
		err := fmt.Errorf("%v is already open", displayIndex)
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

	indexY := 0
	indexX := 0
	if index >= col {
		indexY = index / col
		indexX = index % col
	} else {
		indexX = index
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

	finalFile := UpdatePuzzleImage(x, y, width, height, croppedImg, true)

	remainingPieces[index] = true
	fmt.Printf("Add: %v\n", index)
	openPieces = append(openPieces, index)
	return finalFile, nil
}

func ClosePieceImage(index int) (*os.File, error) {
	indexY := 0
	indexX := 0
	if index >= col {
		indexY = index / col
		indexX = index % col
	} else {
		indexX = index
	}
	width := 152
	height := 90
	x := width * (indexX)
	y := height * (indexY)

	finalFile := UpdatePuzzleImage(x, y, width, height, nil, false)

	remainingPieces[index] = false
	for i := 0; i < len(openPieces); i++ {
		openPiece := openPieces[i]
		if openPiece == index {
			openPieces = append(openPieces[:i], openPieces[i+1:]...)
			i--
			break
		}
	}

	return finalFile, nil
}

func IncreaseOpeningCount(player *model.Player) {
	openingCount++
	if openingCount == player.OpeningCount {
		isOpenPiece = true
		player.OpeningCount = 1
	}
}

func IncreaseAnswerCount(player *model.Player) {
	answerCount++
	if answerCount == player.AnswerCount {
		isAnswer = true
		player.AnswerCount = 1
	}
}

func IncreaseScore(increase int) {
	currentScore += increase
}

func DecreaseScore(decrease int) {
	currentScore -= decrease
}

func GetHint() (string, string) {
	for _, a := range answer.Name {
		isAlphabets := isAlphabets(a)
		if isAlphabets {
			hint += "-"
		} else {
			hint += string(a)
		}
	}

	strList := strings.Split(hint, " ")
	answerSplitLength := "Length is"

	for _, s := range strList {
		answerSplitLength += " " + strconv.Itoa(len(s))
	}

	return hint, answerSplitLength
}

func isAlphabets(c rune) bool {
	return unicode.IsLetter(c)
}

func Answer(message string) (bool, bool, string, *model.Answer) {
	answerLower := strings.ToLower(answer.Name)
	if len(message) != len(answer.Name) {
		return false, false, fmt.Sprintf("Please be sure your answer length is %v", len(answer.Name)), nil
	}

	decreaseCounter := 0
	hintLower := strings.ToLower(hint)

	for i, a := range answerLower {

		if string(message[i]) != string(hintLower[i]) {
			if string(message[i]) == string(a) {
				hint = replaceAtIndex(hint, rune(answer.Name[i]), i)
				hintLower = replaceAtIndex(hintLower, rune(answer.Name[i]), i)
				hintLower = strings.ToLower(hintLower)
				decreaseCounter++
			}
		}
	}

	if hintLower == answerLower {
		return true, true, "", &answer
	} else {
		DecreaseScore(decreaseCounter)
		return false, true, hint, nil
	}
}

func GetFirstAndLastAlphabet() string {
	return fmt.Sprintf("%v, %v\n", string(answer.Name[0]), string(answer.Name[len(answer.Name)-1]))
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func GetRemainingPieces() *map[int]bool {
	return &remainingPieces
}

func NextRound() bool {
	currentRound++
	canNextRound := true
	if currentRound == maxRound {
		canNextRound = false
		currentRound = 0
	}
	return canNextRound
}

func EndRound(players []*model.Player) string {
	scoreboard := "----------Scoreboard----------\n"
	sort.SliceStable(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})
	for _, player := range players {
		scoreboard += fmt.Sprintf("%v : %v\n", player.Name, player.Score)
	}

	return scoreboard
}

func EndGame(players []*model.Player) string {
	isStart = false
	scoreboard := "----------Scoreboard----------\n"
	sort.SliceStable(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})
	for i, player := range players {
		scoreboard += fmt.Sprintf("%v : %v\n", player.Name, player.Score)
		players[i].Score = 0
	}

	scoreboard += fmt.Sprintf("Game ended\n%v win!\n", players[0].Name)

	return scoreboard
}

func RandomItem() int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	randomNumber := r1.Intn(len(itemList.ItemList))
	return randomNumber
}
