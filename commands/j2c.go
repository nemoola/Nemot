package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sashabaranov/go-openai"
	"os"
	"strings"
)

type J2c struct{}

type Chinese struct {
	Chinese  string `json:"chinese"`
	Reading  string `json:"reading"`
	Meanings []struct {
		Word    string `json:"word"`
		Meaning string `json:"meaning"`
	} `json:"meanings"`
}

func (_ *J2c) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	client := openai.NewClient(os.Getenv("CHATGPT_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo0125,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: strings.Join([]string{
						"あなたは私が簡体字の中国語を学習するために存在しています。",
						"中文とピンインとそれぞれ意味を回答してください。",
						"意味は日本語で答えてください。",
						"回答はすべてjsonで返してください。",
						"テンプレートは以下の通りです。",
						`{"chinese":"", "reading":"", "meanings":[{"word":"", "meaning":""}]}`,
						"以下は例です。",
						"input: あなたは猫が好きですか？",
						`output: {"chinese":"你喜欢猫吗？", "reading":"nǐ xǐ huān māo ma？", "meanings":[{"word":"你", "meaning":"あなた"}, {"word":"喜欢", "meaning":"好き"}, {"word":"猫", "meaning":"猫"}, {"word":"吗", "meaning":"〜ですか？"}]}`,
						"input: 私は中国語を勉強中です",
						`output: {"chinese":"我正在学习中文", "reading":"wǒ zài xué xí zhōng wén", "meanings":[{"word":"我", "meaning":"私"}, {"word":"正在", "meaning":"している"}, {"word":"学习", "meaning":"勉強する"}, {"word":"中文", "meaning":"中国語"}]}`,
					}, "\n"),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: strings.Join(args, " "),
				},
			},
		},
	)

	fmt.Println(resp.Choices[0].Message.Content)

	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	var chinese Chinese
	if err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &chinese); err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	embed := discordgo.MessageEmbed{}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "中文",
		Value: chinese.Chinese,
	})
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "読み方",
		Value: chinese.Reading,
	})
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name: "意味",
		Value: func() string {
			var str string
			for _, v := range chinese.Meanings {
				str += fmt.Sprintf("%s - %s\n", v.Word, v.Meaning)
			}
			return str
		}(),
	})
	embed.Color = 0x1a5fb4
	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}
