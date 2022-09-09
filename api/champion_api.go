package api

import (
	"encoding/json"
	"fmt"
	"guess-lol-bot/model"
	"io"
	"log"
	"net/http"
	"os"
)

const championUrl = "http://ddragon.leagueoflegends.com/cdn/12.17.1/data/en_US/champion.json"
const championDataUrl = "http://ddragon.leagueoflegends.com/cdn/12.17.1/data/en_US/champion/"
const imageUrl = "http://ddragon.leagueoflegends.com/cdn/img/champion/splash/"

func GetChampionResponse() (*model.ChampionResponse, error) {
	resp := new(model.ChampionResponse)
	body, err := http.Get(championUrl)
	if err != nil {
		fmt.Println(err)
	}
	defer body.Body.Close()

	if body.StatusCode == 200 {
		bodyBytes, err := io.ReadAll(body.Body)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(bodyBytes, &resp)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	} else {
		fmt.Printf("Error: %v\n", err)
	}

	return resp, err
}

func GetChampionData(resp *model.ChampionResponse) *model.ChampionData {
	championData := new(model.ChampionData)
	for key := range resp.Data {
		championData.Name = append(championData.Name, key)
	}

	return championData
}

func GetChampionDataReponse(championName string) (*model.ChampionDataReponse, error) {
	resp := new(model.ChampionDataReponse)
	body, err := http.Get(championDataUrl + championName + ".json")
	if err != nil {
		fmt.Println(err)
	}
	defer body.Body.Close()

	if body.StatusCode == 200 {
		bodyBytes, err := io.ReadAll(body.Body)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(bodyBytes, &resp)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	} else {
		fmt.Printf("Error: %v\n", err)
	}

	return resp, err
}

func GetSkinImage(url string) (*os.File, error) {
	// don't worry about errors
	response, e := http.Get(imageUrl + url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create("image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Download skin image Success!")
	return file, err
}
