package commands

import "github.com/bwmarrin/discordgo"

type Ping struct{}

func (p Ping) Help() string {
	return "Pong!"
}

func (Ping) Execute(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
}
