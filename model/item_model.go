package model

type Rank int

const (
	Bronze Rank = iota
	Silver
	Gold
	Platinum
	Diamond
	Challenger
)

type ItemList struct {
	ItemList []Item `json:"itemList"`
}

type Item struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ImagePath   string `json:"imagePath"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Rank        int    `json:"rank"`
}
