package tibiachars

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"bytes"
	"time"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/yagpdb/commands"
	"github.com/araddon/dateparse"
)

type Tibia struct {
	Characters struct {
		Error	string	`json:"error"`
		Data struct {
			Name              string   `json:"name"`
			FormerNames       []string `json:"former_names"`
			Title             string   `json:"title"`
			Sex               string   `json:"sex"`
			Vocation          string   `json:"vocation"`
			Level             int      `json:"level"`
			AchievementPoints int      `json:"achievement_points"`
			World             string   `json:"world"`
			FormerWorld       string   `json:"former_world"`
			Residence         string   `json:"residence"`
			MarriedTo         string   `json:"married_to"`
			House             struct {
				Name    string `json:"name"`
				Town    string `json:"town"`
				Paid    string `json:"paid"`
				World   string `json:"world"`
				Houseid int    `json:"houseid"`
			} `json:"house"`
			Guild             struct {
				Name string `json:"name"`
				Rank string `json:"rank"`
			} `json:"guild"`
			LastLogin []struct {
				Date         string `json:"date"`
				TimezoneType int    `json:"timezone_type"`
				Timezone     string `json:"timezone"`
			} `json:"last_login"`
			Comment	string	`json:comment`
			AccountStatus string `json:"account_status"`
			Status        string `json:"status"`
		} `json:"data"`
		Achievements []struct {
			Stars int    `json:"stars"`
			Name  string `json:"name"`
		} `json:"achievements"`
		Deaths []struct {
			Date struct {
				Date         string `json:"date"`
				TimezoneType int    `json:"timezone_type"`
				Timezone     string `json:"timezone"`
			} `json:"date"`
			Level    int    `json:"level"`
			Reason   string `json:"reason"`
			Involved []struct {
				Name string `json:"name"`
			} `json:"involved"`
		} `json:"deaths"`
		AccountInformation ActInfo `json:"account_information"`
		OtherCharacters []struct {
			Name   string `json:"name"`
			World  string `json:"world"`
			Status string `json:"status"`
		} `json:"other_characters"`
	} `json:"characters"`
	Information struct {
		APIVersion    int     `json:"api_version"`
		ExecutionTime float64 `json:"execution_time"`
		LastUpdated   string  `json:"last_updated"`
		Timestamp     string  `json:"timestamp"`
	} `json:"information"`
}

type ActInfo struct {
	LoyaltyTitle string `json:"loyalty_title"`
	Created      struct {
		Date         string `json:"date"`
		TimezoneType int    `json:"timezone_type"`
		Timezone     string `json:"timezone"`
	} `json:"created"`
}

type TibiaWorld struct {
	World struct {
		WorldInformation struct {
			Name          string `json:"name"`
			PlayersOnline int    `json:"players_online"`
			OnlineRecord  struct {
				Players int `json:"players"`
				Date    struct {
					Date         string `json:"date"`
					TimezoneType int    `json:"timezone_type"`
					Timezone     string `json:"timezone"`
				} `json:"date"`
			} `json:"online_record"`
			CreationDate     string   `json:"creation_date"`
			Location         string   `json:"location"`
			PvpType          string   `json:"pvp_type"`
			WorldQuestTitles []string `json:"world_quest_titles"`
			BattleyeStatus   string   `json:"battleye_status"`
			GameWorldType    string   `json:"Game World Type:"`
		} `json:"world_information"`
		PlayersOnline []struct {
			Name     string `json:"name"`
			Level    int    `json:"level"`
			Vocation string `json:"vocation"`
		} `json:"players_online"`
	} `json:"world"`
	Information struct {
		APIVersion    int     `json:"api_version"`
		ExecutionTime float64 `json:"execution_time"`
		LastUpdated   string  `json:"last_updated"`
		Timestamp     string  `json:"timestamp"`
	} `json:"information"`
}

var MainCharCommand = &commands.YAGCommand{
	CmdCategory: commands.CategoryFun,
	Name:        "Char",
	Description: "Retorna informações do personagem especificado.",
	Arguments: []*dcmd.ArgDef{
		&dcmd.ArgDef{Name: "Nome do Char", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {

		char := data.Args[0].Str()

		tibia, err := getChar(char)
		if err != nil {
			if len(char) <= 0 {
				return "Você tem que especificar um char.", err
			} else {
				return "Algo deu errado ao pesquisar esse char.", err
			}
		} else {
			matched, err := regexp.MatchString(`Character does not exist.`, tibia.Characters.Error)
			if matched {
				return "Esse char não existe.", err
			}
		}

		comentario := "Char sem comentário."

		if len(tibia.Characters.Data.Comment) >= 1 {
			comentario = tibia.Characters.Data.Comment
		}

		lealdade := "Sem lealdade"

		if len(tibia.Characters.AccountInformation.LoyaltyTitle) > 0 {
			lealdade = tibia.Characters.AccountInformation.LoyaltyTitle
		}

		world, err := getWorld(tibia.Characters.Data.World)
		if err != nil {
			return "Algo deu errado com o mundo desse char.", err
		}

		level := tibia.Characters.Data.Level
		for _, v := range world.World.PlayersOnline {
			if v.Name == tibia.Characters.Data.Name {
				if v.Level > tibia.Characters.Data.Level {
					level = v.Level
				}
			}
		}

		embed := &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%s", tibia.Characters.Data.Name),
			Color: int(rand.Int63n(16777215)),
			Description: comentario,
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "Level",
					Value:  strconv.Itoa(level),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Mundo",
					Value:  tibia.Characters.Data.World,
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Vocação",
					Value:  tibia.Characters.Data.Vocation,
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Templo",
					Value:  tibia.Characters.Data.Residence,
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Status",
					Value:  tibia.Characters.Data.AccountStatus,
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "On/Off",
					Value:  strings.Title(tibia.Characters.Data.Status),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Lealdade",
					Value:  lealdade,
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Pontos de Achievement",
					Value:  strconv.Itoa(tibia.Characters.Data.AchievementPoints),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Gênero",
					Value:  strings.Title(tibia.Characters.Data.Sex),
					Inline: true,
				},
			},
		}

		if len(tibia.Characters.Data.Guild.Name) >= 1 {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:  "Guild",
				Value: tibia.Characters.Data.Guild.Name,
				Inline: true,
			})
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:  "Cargo na Guild",
				Value: tibia.Characters.Data.Guild.Rank,
				Inline: true,
			})
		} else {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:  "Guild",
				Value: "Nenhuma",
				Inline: true,
			})
		}

		if len(tibia.Characters.Data.House.Name) >= 1 {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:  "Casa",
				Value: tibia.Characters.Data.House.Name,
				Inline: true,
			})
		}

		if len(tibia.Characters.Data.MarriedTo) >= 1 {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:  "Casado",
				Value: tibia.Characters.Data.MarriedTo,
				Inline: true,
			})
		}

		if len(tibia.Characters.AccountInformation.Created.Date) > 0 {
			t, err := dateparse.ParseLocal(tibia.Characters.AccountInformation.Created.Date)
			if err != nil {
				return "Algo deu errado ao pesquisar esse char, por causa da data de criação.", err
			}
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:  "Criado",
				Value: (t.Add(time.Hour * -5)).Format("02/01/2006 15:04:05 BRT"),
				Inline: true,
			})
		}

		return embed, nil
	},
}

var DeathsCommand = &commands.YAGCommand{
	CmdCategory: commands.CategoryFun,
	Name:        "Mortes",
	Description: "Retorna as mortes recentes do personagem especificado.",
	Arguments: []*dcmd.ArgDef{
		&dcmd.ArgDef{Name: "Nome do Char", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {

		char := data.Args[0].Str()

		tibia, err := getChar(char)
		if err != nil {
			if len(char) <= 0 {
				return "Você tem que especificar um char.", err
			} else {
				return "Algo deu errado ao pesquisar esse char.", err
			}
		} else {
			matched, err := regexp.MatchString(`Character does not exist.`, tibia.Characters.Error)
			if matched {
				return "Esse char não existe.", err
			}
		}

		mortes := tibia.Characters.Deaths
		deaths := "\n"
		motivo := ""

		if data.Source == dcmd.DMSource {
			if len(mortes) >= 1 {
				t, err := dateparse.ParseLocal(mortes[0].Date.Date)
				if err != nil {
					return "Algo deu errado ao pesquisar esse char, por causa da data de criação.", err
				}
				embedCC := &discordgo.MessageEmbed{
					Title: fmt.Sprintf("Mortes recentes de %s", tibia.Characters.Data.Name),
					Description: fmt.Sprintf("**Data**: %s\n**Level**: %d\n**Motivo**: %s\n\n", (t.Add(time.Hour * -5)).Format("02/01/2006 15:04:05 BRT"), mortes[0].Level, mortes[0].Reason),
					Color: int(rand.Int63n(16777215)),
				}
				return embedCC, nil
			} else {
				return "Esse char não tem mortes recentes.", nil
			}
		}

		if len(mortes) >= 1 {
			for _, v := range mortes {
				t2, err := dateparse.ParseLocal(v.Date.Date)
				if err != nil {
					return "Algo deu errado ao pesquisar esse char, por causa da data de criação.", err
				}
				if len(deaths) < 1800 {
					checkKillByMonster, _ := regexp.MatchString(`Died by a`, v.Reason)
					if checkKillByMonster {
						deaths += fmt.Sprintf("**Data**: %s\n**Level**: %d\n**Motivo**: %s\n\n", (t2.Add(time.Hour * -5)).Format("02/01/2006 15:04:05 BRT"), v.Level, v.Reason)
					} else {
						split := strings.Split(v.Reason, ",")
						for i := range split {
							checkOutros, _ := regexp.MatchString(`e outros.`, motivo)
							if len(motivo) < 150 {
								motivo += fmt.Sprintf("%s, ", split[i])
							} else {
								if !checkOutros {
									motivo += "e outros."
								}
							}
						}
						re := regexp.MustCompile(`, \z`)
						motivo = re.ReplaceAllString(motivo, ".")
						deaths += fmt.Sprintf("**Data**: %s\n**Level**: %d\n**Motivo**: %s\n\n", (t2.Add(time.Hour * -5)).Format("02/01/2006 15:04:05 BRT"), v.Level, motivo)
						motivo = ""
					}
				} else {
					checkOutras, _ := regexp.MatchString(`... entre outras ...`, deaths)
					if !checkOutras {
						deaths += "... entre outras ..."
					}
				}
			}
		} else {
			deaths = "Sem mortes recentes."
		}

		embed := &discordgo.MessageEmbed{
			Title: fmt.Sprintf("Mortes recentes de %s", tibia.Characters.Data.Name),
			Description: deaths,
			Color: int(rand.Int63n(16777215)),
		}

		return embed, nil

	},
}

var CheckOnlineCommand = &commands.YAGCommand{
	CmdCategory: commands.CategoryFun,
	Name:        "CheckOnline",
	Description: "Mostra quem está online no mundo especificado.",
	Aliases:		[]string{"co"},
	Arguments: []*dcmd.ArgDef{
		&dcmd.ArgDef{Name: "Nome do Mundo", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {

		mundo := data.Args[0].Str()

		world, err := getWorld(mundo)
		if err != nil {
			if len(mundo) <= 0 {
				return "Você tem que especificar um mundo.", err
			} else {
				return "Algo deu errado ao pesquisar esse char.", err
			}
		} else {
			if len(world.World.WorldInformation.CreationDate) == 0 {
				return "Esse mundo não existe.", err
			}
		}

		if data.Source == dcmd.DMSource {
			m := make([]map[string]interface{}, len(world.World.PlayersOnline))
			for k, v := range world.World.PlayersOnline {
				m[k] = make(map[string]interface{})
				m[k]["Name"] = v.Name
				m[k]["Level"] = v.Level
				m[k]["Vocation"] = v.Vocation
			}
			return m, nil
		}

		desc := ""

		if len(world.World.PlayersOnline) > 0 {
			for _, v := range world.World.PlayersOnline {
				checkEnd, _ := regexp.MatchString(`e outros.`, desc)
				if len(desc) < 1948 {
					desc += fmt.Sprintf("%s, ", v.Name)
				} else {
					if !checkEnd {
						desc += "e outros."
					}
				}
			}
			re := regexp.MustCompile(`, \z`)
			desc = re.ReplaceAllString(desc, ".")
		} else {
			desc = "Nenhum player online."
		}

		embed := &discordgo.MessageEmbed{
			Title: fmt.Sprintf("Players online em %s", world.World.WorldInformation.Name),
			Description: desc,
			Color: int(rand.Int63n(16777215)),
		}

		return embed, nil

	},
}

func (ai *ActInfo) UnmarshalJSON(data []byte) error {
	if bytes.HasPrefix(data, []byte("{")) {
		type actInfoNoMethods ActInfo
		return json.Unmarshal(data, (*actInfoNoMethods)(ai))
	}
	return nil
}

func getChar(name ...string) (*Tibia, error) {
	tibia := Tibia{}
	queryUrl := ""

	if len(name) >= 1 {
		queryUrl = fmt.Sprintf("https://api.tibiadata.com/v2/characters/%s.json", name[0])
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
		return nil, queryErr
	}

	return &tibia, nil
}

func getWorld(name ...string) (*TibiaWorld, error) {
	world := TibiaWorld{}
	queryUrl := ""

	if len(name) >= 1 {
		queryUrl = fmt.Sprintf("https://api.tibiadata.com/v2/world/%s.json", name[0])
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

	queryErr := json.Unmarshal(body, &world)
	if queryErr != nil {
		return nil, queryErr
	}

	return &world, nil
}
