package tibiaguilds

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common/templates"
)

var SpecificGuildCommand = &commands.YAGCommand{
	CmdCategory: commands.CategoryFun,
	Name:        "Guild",
	Description: "Retorna informações da guild especificada.",
	Arguments: []*dcmd.ArgDef{
		&dcmd.ArgDef{Name: "Nome da Guild", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		guildName := data.Args[0].Str()

		guild, err := templates.GetSpecificGuild(guildName)
		if err != nil {
			if len(guildName) <= 0 {
				return "Você tem que especificar uma guild.", err
			} else {
				return "Algo deu errado ao pesquisar esse char.", err
			}
		} else if len(guild.Guild.Error) >= 1 {
			return "Essa guild não existe.", err
		}

		desc := "Guild sem descrição."
		if len(guild.Guild.Data.Description) >= 1 && len(guild.Guild.Data.Description) < 2048 {
			desc = guild.Guild.Data.Description
		}

		guildHall := "Nenhuma."
		if len(guild.Guild.Data.Guildhall.Name) > 1 {
			guildHall = fmt.Sprintf("**%s** que fica em %s", guild.Guild.Data.Guildhall.Name, guild.Guild.Data.Guildhall.Town)
		}

		guerra := "Não."
		if guild.Guild.Data.War {
			guerra = "Sim."
		}

		embed := &discordgo.MessageEmbed{
			Title: guild.Guild.Data.Name,
			Color: int(rand.Int63n(16777215)),
			Description: desc,
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "Número de membros",
					Value:  strconv.Itoa(guild.Guild.Data.Totalmembers),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Mundo",
					Value:  guild.Guild.Data.World,
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Guild Hall",
					Value:  guildHall,
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Está em Guerra?",
					Value:  guerra,
					Inline: true,
				},
			},
		}

		return embed, nil
	},
}
