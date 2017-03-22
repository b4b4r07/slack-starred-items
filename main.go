package main

import (
	"fmt"
	"html/template"
	"os"
	// "strings"

	"github.com/nlopes/slack"
	"github.com/pelletier/go-toml"
	"github.com/russross/blackfriday"
)

// Data to put into  template
type Page struct {
	Title string
	Body  string
}

// The template
var templateText string = `
<head>
  <title>{{.Title}}</title>
</head>

<body>
  {{.Body | markDown}}
</body>
`

// Real blackfriday functionality commented out, using strings.ToLower for demo
func markDowner(args ...interface{}) template.HTML {
	return template.HTML(blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...))))
	// return template.HTML(strings.ToLower(fmt.Sprintf("%s", args...)))
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
	text := ""

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
				text += fmt.Sprintf("## @%s\n\n%s\n\n", user, star.Message.Msg.Text)
			}
		}
		pages = paging.Pages
		starParams.Page++
	}

	// Create a page
	p := &Page{Title: "A Test Demo", Body: text}

	// Parse the template and add the function to the funcmap
	tmpl := template.Must(template.New("page.html").Funcs(template.FuncMap{"markDown": markDowner}).Parse(templateText))

	// Execute the template
	err = tmpl.ExecuteTemplate(os.Stdout, "page.html", p)
	if err != nil {
		fmt.Println(err)
	}
}
