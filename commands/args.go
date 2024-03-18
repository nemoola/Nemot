package commands

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

type Args struct{}

func (Args) Help() string {
	return "引数をそのまま返します。"
}

func (Args) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	_, _ = s.ChannelMessageSend(m.ChannelID, strings.Join(args, " "))
}
