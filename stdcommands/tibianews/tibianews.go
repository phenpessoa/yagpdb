package tibianews

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"math/rand"
	"regexp"
	"strings"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/yagpdb/commands"
	"github.com/araddon/dateparse"
)

type TibiaNews struct {
	Newslist struct {
		Type string `json:"type"`
		Data []struct {
			ID       int    `json:"id"`
			Type     string `json:"type"`
			News     string `json:"news"`
			Apiurl   string `json:"apiurl"`
			Tibiaurl string `json:"tibiaurl"`
			Date     struct {
				Date         string `json:"date"`
				TimezoneType int    `json:"timezone_type"`
				Timezone     string `json:"timezone"`
			} `json:"date"`
		} `json:"data"`
	} `json:"newslist"`
	Information struct {
		APIVersion    int     `json:"api_version"`
		ExecutionTime float64 `json:"execution_time"`
		LastUpdated   string  `json:"last_updated"`
		Timestamp     string  `json:"timestamp"`
	} `json:"information"`
}

type TibiaSpecificNews struct {
	News struct {
		Error 	string `json:"error"`
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
		Date    struct {
			Date         string `json:"date"`
			TimezoneType int    `json:"timezone_type"`
			Timezone     string `json:"timezone"`
		} `json:"date"`
	} `json:"news"`
	Information struct {
		APIVersion    int     `json:"api_version"`
		ExecutionTime float64 `json:"execution_time"`
		LastUpdated   string  `json:"last_updated"`
		Timestamp     string  `json:"timestamp"`
	} `json:"information"`
}

var NewsCommand = &commands.YAGCommand{
	CmdCategory: commands.CategoryFun,
	Name:			"News",
	Aliases:		[]string{"noticia"},
	Description:	"Última noticia do tibia, ou alguma específica.",
	Arguments: []*dcmd.ArgDef{
		&dcmd.ArgDef{Name: "ID da notícia", Type: dcmd.Int},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {

		tibia, err := getNews("")
		if err != nil {
			return "Algo deu errado com essa pesquisa.", err
		}

		url := tibia.Newslist.Data[0].Tibiaurl
		inside := tibia.Newslist.Data[0].ID

		if data.Args[0].Value != nil {
			inside = data.Args[0].Int()
		}

		url = fmt.Sprintf("https://www.tibia.com/news/?subtopic=newsarchive&id=%d", inside)

		tibiaInside, err := insideNews(inside)
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

		tibia, err := getNews("ticker")
			if err != nil {
				return "Algo deu errado com essa pesquisa.", err
		}

		tibiaInside, err := insideNews(tibia.Newslist.Data[0].ID)
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

func getNews(name string) (*TibiaNews, error) {
	tibia := TibiaNews{}
	queryUrl := "https://api.tibiadata.com/v2/latestnews.json"

	if name == "ticker" {
		queryUrl = "https://api.tibiadata.com/v2/newstickers.json"
	}

	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "curl/7.65.1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	queryErr := json.Unmarshal(body, &tibia)
	if queryErr != nil {
		return nil, err
	}

	return &tibia, nil
}


func insideNews(number int) (*TibiaSpecificNews, error) {
	tibiaInside := TibiaSpecificNews{}
	queryUrl := fmt.Sprintf("https://api.tibiadata.com/v2/news/%d.json", number)

	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "curl/7.65.1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	queryErr := json.Unmarshal(body, &tibiaInside)
	if queryErr != nil {
		return nil, err
	}

	return &tibiaInside, nil
}
