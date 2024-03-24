package commands

import (
	"Nemot/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type Wiki struct{}

func (Wiki) Help() string {
	return "Wikipediaの検索結果を返します"
}

func (Wiki) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	wc := utils.NewWikiClient()
	result := wc.Search(args[0])
	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		URL:         fmt.Sprintf("https://ja.wikipedia.org/wiki/%s", result.Title),
		Title:       result.Title,
		Description: result.Content,
		Color:       0x00ff00,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Wikipedia",
			IconURL: "https://upload.wikimedia.org/wikipedia/commons/thumb/8/80/Wikipedia-logo-v2.svg/103px-Wikipedia-logo-v2.svg.png",
		},
	})
}
