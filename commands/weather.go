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
		embed := discordgo.MessageEmbed{}
		embed.Title = fmt.Sprintf("%sの天気情報", args[0])
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "天気",
			Value: fmt.Sprintf("WIP"),
		})
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "気温",
			Value: fmt.Sprintf("%.1f℃", amedas.Temp[0]),
		})
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "降水量(10分)",
			Value: fmt.Sprintf("%.1fmm", amedas.Precipitation10m[0]),
		})
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "風速",
			Value: fmt.Sprintf("%.1fm/s", amedas.WindDirection[0]),
		})
		embed.Color = 0x00ff00
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		if err != nil {
			fmt.Println(err)
		}
	}
}
