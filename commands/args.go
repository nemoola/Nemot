package commands

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

type Args struct{}

func (_ *Args) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	_, _ = s.ChannelMessageSend(m.ChannelID, strings.Join(args, " "))
}
