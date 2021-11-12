package example

import "fmt"


type Game struct {
	Game map[string]interface{} `json:"game"`
}

func GetGame(key string) (*Game, error) {
	response, err := httpClient.Get("http://localhost:5005/games/" + key, nil)

	if err != nil {
		return nil, err
	}

	fmt.Println(response.StatusCode())
	fmt.Println(response.String())
	var game Game

	if err := response.UnmarshalJson(&game); err != nil {
		return nil, err
	}
	fmt.Println(game)

	return &game, nil
}
