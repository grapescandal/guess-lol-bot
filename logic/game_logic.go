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
var currentScore int = 65
var maxScore int = 65
var pieceScore int = 64
var championData *model.ChampionData
var hint string = ""

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
	maxTurn = 0
	openPieces = []int{}
	currentScore = maxScore
	hint = ""
}

func StartGame() {
	if !isStart {
		isStart = true
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

		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		randomNumber := r1.Intn(maxTurn)
		turn = randomNumber
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
	if turn >= maxTurn {
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

func ReadSkinImage(fileName string) *os.File {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	return file
}

func CreatePuzzleImage() {

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

func UpdatePuzzleImage(inputX int, inputY int, inputWidth int, inputHeight int, piece image.Image) *os.File {

	file := ReadSkinImage("puzzle.jpg")
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	width := 1215
	height := 717

	imgRGBA := imageToRGBA(img)

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
	if index > col {
		indexY = (index / col)
		indexX = index % col
		if indexX == 0 {
			indexX = col - 1
		}
	} else {
		indexX = index - 1
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

	finalFile := UpdatePuzzleImage(x, y, width, height, croppedImg)

	openPieces = append(openPieces, index)
	return finalFile, nil
}

func DecreaseScore(decrease int) {
	currentScore -= decrease
}

func GetHint() (string, int) {
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

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}
