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

type C2j struct{}

type Japanese struct {
	Japanese string `json:"japanese"`
	Reading  string `json:"reading"`
	Meanings []struct {
		Word    string `json:"word"`
		Meaning string `json:"meaning"`
	} `json:"meanings"`
}

func (C2j) Help() string {
	return "中国語を日本語に変換します。"
}

func (C2j) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	client := openai.NewClient(os.Getenv("CHATGPT_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo0125,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: strings.Join([]string{
						"あなたは私が中国語を学習するために存在しています。",
						"日本語とピンインとそれぞれ意味を回答してください。",
						"意味は日本語で答えてください。",
						"回答はすべてjsonで返してください。",
						"テンプレートは以下の通りです。",
						`{"japanese":"", "reading":"", "meanings":[{"word":"", "meaning":""}]}`,
						"以下は例です。",
						"input: 你喜欢猫吗？",
						`output: {"japanese":"あなたは猫が好きですか？", "reading":"nǐ xǐ huān māo má？", "meanings":[{"word":"你", "meaning":"あなた"}, {"word":"喜欢", "meaning":"好き"}, {"word":"猫", "meaning":"猫"}, {"word":"吗", "meaning":"〜ですか？"}]}`,
						"input: 我在学习中文",
						`output: {"japanese":"私は中国語を勉強中です", "reading":"wǒ zài xué xí zhōng wé", "meanings":[{"word":"我", "meaning":"私"}, {"word":"在", "meaning":"している"}, {"word":"学习", "meaning":"学ぶ"}, {"word":"中文", "meaning":"中国語"}]}`,
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

	var japanese Japanese
	if err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &japanese); err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	embed := discordgo.MessageEmbed{}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "日本語",
		Value: japanese.Japanese,
	})
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "読み方",
		Value: japanese.Reading,
	})
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name: "意味",
		Value: func() string {
			var str string
			for _, v := range japanese.Meanings {
				str += fmt.Sprintf("%s - %s\n", v.Word, v.Meaning)
			}
			return str
		}(),
	})
	embed.Color = 0x1a5fb4
	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}
