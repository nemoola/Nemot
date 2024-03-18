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

type J2co struct{}

type Cantonese struct {
	Cantonese string `json:"cantonese"`
	Reading   string `json:"reading"`
	Meanings  []struct {
		Word    string `json:"word"`
		Meaning string `json:"meaning"`
	} `json:"meanings"`
}

func (c J2co) Help() string {
	return "日本語を広東語に変換します。"
}

func (J2co) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	client := openai.NewClient(os.Getenv("CHATGPT_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo0125,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: strings.Join([]string{
						"あなたは私が広東語を学習するために存在しています。",
						"広東語の文章とピンインとそれぞれ意味を回答してください。",
						"意味は日本語で答えてください。",
						"回答はすべてjsonで返してください。",
						"テンプレートは以下の通りです。",
						`{"cantonese":"", "reading":"", "meanings":[{"word":"", "meaning":""}]}`,
						"以下は例です。",
						"input: あなたは猫が好きですか？",
						`output: {"cantonese":"你鍾意貓嗎？", "reading":"nei5 zung1 ji3 maau1 maa1？", "meanings":[{"word":"你", "meaning":"あなた"}, {"word":"鍾意", "meaning":"好き"}, {"word":"猫", "meaning":"猫"}, {"word":"吗", "meaning":"〜ですか？"}]}`,
						"input: 私は中国語を勉強中です",
						`output: {"cantonese":"我喺學習中文", "reading":"ngo5 hai2 hok6 zaap6 zung1 man4", "meanings":[{"word":"我", "meaning":"私"}, {"word":"喺", "meaning":"している"}, {"word":"學習", "meaning":"学ぶ"}, {"word":"中文", "meaning":"中国語"}]}`,
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

	var cantonese Cantonese
	if err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &cantonese); err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	embed := discordgo.MessageEmbed{}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "広東語",
		Value: cantonese.Cantonese,
	})
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "読み方",
		Value: cantonese.Reading,
	})
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name: "意味",
		Value: func() string {
			var str string
			for _, v := range cantonese.Meanings {
				str += fmt.Sprintf("%s - %s\n", v.Word, v.Meaning)
			}
			return str
		}(),
	})
	embed.Color = 0x1a5fb4
	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}
