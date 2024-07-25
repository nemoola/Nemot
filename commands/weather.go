package commands

import (
	"Nemot/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type Weather struct{}

func (Weather) Help() string {
	return "天気情報を表示します。"
}

func (Weather) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "地名を指定してください。")
		return
	} else {
		weather := utils.NewWeather()
		amedas, err := weather.GetWeather(args[0])
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%.1f℃", amedas.Temp[0]))
	}
}
