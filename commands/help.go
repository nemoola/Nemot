package commands

import (
	"github.com/bwmarrin/discordgo"
	"sort"
)

type Help struct{}

func (Help) Help() string {
	return "コマンド一覧とヘルプを表示します。"
}

func (Help) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	commandList := CommandList()

	var cmds []string
	for cmd := range commandList {
		cmds = append(cmds, cmd)
	}
	sort.Strings(cmds)

	embed := &discordgo.MessageEmbed{}
	embed.Title = "コマンド一覧"
	for _, command := range cmds {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  command,
			Value: commandList[command].Help(),
		})
	}
	embed.Color = 0x00ff00
	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
