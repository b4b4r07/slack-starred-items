package main

import (
	"fmt"

	"github.com/nlopes/slack"
	"github.com/pelletier/go-toml"
)

func main() {
	config, err := toml.LoadFile("config.toml")
	if err != nil {
		panic(err)
	}

	apiToken := config.Get("slack.access_token").(string)
	userId := config.Get("slack.user_id").(string)

	api := slack.New(apiToken)

	// Get all starred items.
	starParams := slack.NewStarsParameters()
	pages := 1
	starParams.User = userId
	starParams.Page = 1
	for starParams.Page <= pages {
		stars, paging, err := api.ListStars(starParams)
		if err != nil {
			panic(err)
		}
		for _, star := range stars {
			if star.Type == "message" {
				fmt.Println(star.Message.Msg)
			}
		}
		pages = paging.Pages
		starParams.Page++
	}
}
