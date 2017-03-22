package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/nlopes/slack"
	"github.com/pelletier/go-toml"
)

func init() {
	color.NoColor = false
}

func main() {
	config, err := toml.LoadFile("config.toml")
	if err != nil {
		panic(err)
	}

	apiToken := config.Get("slack.access_token").(string)
	userId := config.Get("slack.user_id").(string)

	api := slack.New(apiToken)

	users, err := api.GetUsers()
	if err != nil {
		panic(err)
	}
	users_map := make(map[string]string, len(users))
	for _, user := range users {
		users_map[user.ID] = user.Name
	}
	fmt.Printf("%#v\n", users_map)

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
				if star.Message.Msg.User == "" {
					continue
				}
				user := users_map[star.Message.Msg.User]
				printStarredMessage(user, star.Message.Msg.Text)
			}
		}
		pages = paging.Pages
		starParams.Page++
	}
}

func printStarredMessage(username, text string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("%s\t%s\n", yellow("@"+username), text)
}
