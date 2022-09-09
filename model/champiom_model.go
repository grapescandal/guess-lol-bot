package model

type ChampionResponse struct {
	Data map[string]interface{} `json:"data"`
}

type ChampionData struct {
	Name []string
}

type ChampionDataReponse struct {
	Data map[string]ChampionSkin `json:"data"`
}

type ChampionSkin struct {
	Skins []Skins `json:"skins"`
}

type Skins struct {
	Num  int    `json:"num"`
	Name string `json:"name"`
}
