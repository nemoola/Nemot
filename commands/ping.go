package commands

import "github.com/bwmarrin/discordgo"

type Ping struct{}

func (_ *Ping) Execute(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
}
