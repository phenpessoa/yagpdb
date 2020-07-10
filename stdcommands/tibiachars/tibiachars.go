package tibiachars

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common/templates"
	"github.com/araddon/dateparse"
)

var MainCharCommand = &commands.YAGCommand{
	CmdCategory: commands.CategoryFun,
	Name:        "Char",
	Description: "Retorna informações do personagem especificado.",
	Arguments: []*dcmd.ArgDef{
		&dcmd.ArgDef{Name: "Nome do Char", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		char := data.Args[0].Str()

		tibia, err := templates.GetChar(char)
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

		world, err := templates.GetWorld(tibia.Characters.Data.World)
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

		tibia, err := templates.GetChar(char)
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
							checkOutros, _ := regexp.MatchString(`e outros.\z`, motivo)
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
					checkOutras, _ := regexp.MatchString(`\.\.\. entre outras \.\.\.\z`, deaths)
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

		world, err := templates.GetWorld(mundo)
		if err != nil {
			if len(mundo) <= 0 {
				return "Você tem que especificar um mundo.", err
			} else {
				return "Algo deu errado ao pesquisar esse mundo.", err
			}
		} else {
			if len(world.World.WorldInformation.CreationDate) == 0 {
				return "Esse mundo não existe.", err
			}
		}

		if data.Source == dcmd.DMSource {
			m := make([]map[string]interface{}, len(world.World.PlayersOnline))
			for k, v := range world.World.PlayersOnline {
				m[k] = map[string]interface{}{
					"Name": v.Name,
					"Level": v.Level,
					"Vocation": v.Vocation,
				}
			}
			return m, nil
		}

		desc := ""

		if len(world.World.PlayersOnline) > 0 {
			for _, v := range world.World.PlayersOnline {
				checkEnd, _ := regexp.MatchString(`e outros.\z`, desc)
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
