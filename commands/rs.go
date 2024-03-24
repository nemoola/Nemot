package commands

import (
	"Nemot/utils"
	"github.com/bwmarrin/discordgo"
)

type Rs struct{}

func (Rs) Help() string {
	return "現在の震度を表示します"
}

func (Rs) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	eew := utils.NewEEW()
	eew.GetEEW()

	_, _ = s.ChannelFileSend(m.ChannelID, "eew.gif", &eew.Img)
}
