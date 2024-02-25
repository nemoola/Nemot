package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/api/option"
	"math/big"
	"os"
	"sort"
	"strings"
)

type Chinese struct {
	Chinese  string `json:"chinese"`
	Reading  string `json:"reading"`
	Meanings []struct {
		Word    string `json:"word"`
		Meaning string `json:"meaning"`
	} `json:"meanings"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	bot, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("DISCORD_TOKEN")))
	if err != nil {
		panic(err)
	}

	bot.AddHandler(messageCreate)
	bot.AddHandler(onReady)

	if err := bot.Open(); err != nil {
		panic(err)
	}

	<-make(chan struct{})
	if err := bot.Close(); err != nil {
		panic(err)
	}
}

func onReady(_ *discordgo.Session, _ *discordgo.Ready) {
	fmt.Println("Bot is ready")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "n!") {
		command := strings.Split(m.Content, " ")[0][len("n!"):]
		args := strings.Split(m.Content, " ")[1:]
		switch command {
		case "gj2c":
			ctx := context.Background()
			client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_KEY")))
			if err != nil {
				return
			}

			resp, err := client.GenerativeModel("gemini-pro").GenerateContent(ctx, genai.Text(strings.Join([]string{
				"あなたは私が簡体字の中国語を学習するために存在しています。",
				"中文とピンインとそれぞれ意味を回答してください。",
				"回答はすべてjsonで返してください。",
				"テンプレートは以下の通りです。",
				`{"chinese":"", "reading":"", "meanings":[{"word":"", "meaning":""}]}`,
				"以下は例です。",
				"input: あなたは猫が好きですか？",
				`output: {"chinese":"你喜欢猫吗？", "reading":"nǐ xǐ huān māo ma？", "meanings":[{"word":"你", "meaning":"あなた"}, {"word":"喜欢", "meaning":"好き"}, {"word":"猫", "meaning":"猫"}, {"word":"吗", "meaning":"〜ですか？"}]}`,
				"input: 私は中国語を勉強中です",
				`output: {"chinese":"我正在学习中文", "reading":"wǒ zhèng zài xué xí zhōng wén", "meanings":[{"word":"我", "meaning":"私"}, {"word":"正在", "meaning":"している"}, {"word":"学习", "meaning":"勉強する"}, {"word":"中文", "meaning":"中国語"}]}`,
			}, "\n")), genai.Text(strings.Join(args, " ")))

			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}

			var chinese Chinese
			if err = json.Unmarshal([]byte(fmt.Sprint(resp.Candidates[0].Content.Parts[0])), &chinese); err != nil {
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
		case "j2c":
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
								`output: {"chinese":"我正在学习中文", "reading":"wǒ zhèng zài xué xí zhōng wén", "meanings":[{"word":"我", "meaning":"私"}, {"word":"正在", "meaning":"している"}, {"word":"学习", "meaning":"勉強する"}, {"word":"中文", "meaning":"中国語"}]}`,
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
		case "ping":
			_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
		case "args":
			_, _ = s.ChannelMessageSend(m.ChannelID, strings.Join(args, " "))
		default:
			_, _ = s.ChannelMessageSend(m.ChannelID, "Unknown command")
		}
	}

	if rnd, _ := rand.Int(rand.Reader, big.NewInt(100)); rnd.Int64() <= 30 {
		msgs, _ := s.ChannelMessages(m.ChannelID, 3, "", "", "")
		sort.Slice(msgs, func(i, j int) bool {
			return msgs[j].Timestamp.After(msgs[i].Timestamp)
		})

		var messages []openai.ChatCompletionMessage
		messages = append(messages, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleSystem,
			Content: strings.Join([]string{
				"あなたはランダムの確率でチャットを取得するbotです。その与えられたチャットに対して文章の意味がわからない場合は適当に返信して短く返信をするbotです。",
				"その与えられたチャットに対して文章の意味がわからない場合は適当に返信してください。",
				"意味がわかる場合はその会話にまじるように返信してください。",
				"返信は短くお願いします。",
			}, "\n"),
		})

		for _, msg := range msgs {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: msg.Content,
			})
		}

		client := openai.NewClient(os.Getenv("CHATGPT_KEY"))
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    openai.GPT3Dot5Turbo0125,
				Messages: messages,
			},
		)

		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}

		fmt.Println(resp.Choices[0].Message.Content)
		_, _ = s.ChannelMessageSendReply(m.ChannelID, resp.Choices[0].Message.Content, m.Reference())
	}

	var onMention = false
	for _, user := range m.Mentions {
		if user.ID == s.State.User.ID {
			onMention = true
			break
		}
	}
	if onMention {
		//msg, _ := s.ChannelMessageSend(m.ChannelID, "回答を生成中...")
		//gemini.GenerateContentStream(context.Background(), m.Content)
		//for gemini.Guwaa() {
		//	_, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, "回答を生成中...")
		//}
	}
}
