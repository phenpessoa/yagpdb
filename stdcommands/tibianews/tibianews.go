package tibianews

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common/templates"
	"github.com/araddon/dateparse"
)

var NewsCommand = &commands.YAGCommand{
	CmdCategory: commands.CategoryFun,
	Name:			"News",
	Aliases:		[]string{"noticia"},
	Description:	"Última noticia do tibia, ou alguma específica.",
	Arguments: []*dcmd.ArgDef{
		&dcmd.ArgDef{Name: "ID da notícia", Type: dcmd.Int},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		tibia, err := templates.GetNews("")
		if err != nil {
			return "Algo deu errado com essa pesquisa.", err
		}

		url := tibia.Newslist.Data[0].Tibiaurl
		inside := tibia.Newslist.Data[0].ID

		if data.Args[0].Value != nil {
			inside = data.Args[0].Int()
		}

		url = fmt.Sprintf("https://www.tibia.com/news/?subtopic=newsarchive&id=%d", inside)

		tibiaInside, err := templates.InsideNews(inside)
		if err != nil {
			return "Algo deu errado com essa pesquisa.", err
		}

		if len(tibiaInside.News.Error) >= 1 {
			return "Essa notícia não existe.", err
		}

		t, err := dateparse.ParseLocal(tibiaInside.News.Date.Date)
		if err != nil {
			return "Algo deu errado ao pesquisar essa notícia, por causa da data de criação.", err
		}

		re := regexp.MustCompile(`<(.*?)>`)
		output := re.ReplaceAllString(tibiaInside.News.Content, "")
		finalOut := ""

		if len(output) > 1400 {
			split := strings.Split(output, " ")
			for i := range split {
				checkEtc, _ := regexp.MatchString(`\.\.\.\z`, finalOut)
				if len(finalOut) < 1400 {
					finalOut += fmt.Sprintf("%s ", split[i])
				} else {
					if !checkEtc {
						finalOut += "..."
					}
				}
			}
		}

		if finalOut == "" {
			finalOut = output
		}

		embed := &discordgo.MessageEmbed{
			Title:	fmt.Sprintf("%s", tibiaInside.News.Title),
			Color:       int(rand.Int63n(16777215)),
			Description:	fmt.Sprintf("%s\n[Clique para ver mais](%s)", finalOut, url),
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Notícia mais recente do Tibia. | ID: %d\nData: %s", tibiaInside.News.ID, t.Format("02/01/2006")),
			},
		}

		return embed, nil

	},
}

var NewsTickerCommand = &commands.YAGCommand{
	CmdCategory: commands.CategoryFun,
	Name:        "NewsTicker",
	Description: "Último newsticker do tibia.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		tibia, err := templates.GetNews("ticker")
			if err != nil {
				return "Algo deu errado com essa pesquisa.", err
		}

		tibiaInside, err := templates.InsideNews(tibia.Newslist.Data[0].ID)
		if err != nil {
			return "Algo deu errado com essa pesquisa.", err
		}

		if len(tibiaInside.News.Error) >= 1 {
			return "Essa notícia não existe.", err
		}

		t, err := dateparse.ParseLocal(tibiaInside.News.Date.Date)
		if err != nil {
			return "Algo deu errado ao pesquisar essa notícia, por causa da data de criação.", err
		}

		re := regexp.MustCompile(`<(.*?)>`)
		output := re.ReplaceAllString(tibiaInside.News.Content, "")
		finalOut := ""

		if len(output) > 1400 {
			split := strings.Split(output, " ")
			for i := range split {
				checkEtc, _ := regexp.MatchString(`\.\.\.\z`, finalOut)
				if len(finalOut) < 1400 {
					finalOut += fmt.Sprintf("%s ", split[i])
				} else {
					if !checkEtc {
						finalOut += "..."
					}
				}
			}
		}

		if finalOut == "" {
			finalOut = output
		}

		embed := &discordgo.MessageEmbed{
			Title:	fmt.Sprintf("%s", tibiaInside.News.Title),
			Color:       int(rand.Int63n(16777215)),
			Description:	fmt.Sprintf("%s\n[Clique para ver mais](%s)", finalOut, tibia.Newslist.Data[0].Tibiaurl),
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Notícia mais recente do Tibia. | ID: %d\nData: %s", tibiaInside.News.ID, t.Format("02/01/2006")),
			},
		}

		return embed, nil

	},
}
