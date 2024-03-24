package main

import (
	"Nemot/commands"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"io"
	"math/big"
	"net/http"
	"os"
	"sort"
	"strings"
)

var (
	prefix     = ""
	commandMap = commands.CommandList()
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	prefix = os.Getenv("PREFIX")

	if _, err := os.Stat("assets"); !os.IsExist(err) {
		_ = os.Mkdir("assets", 0777)
	}

	if _, err := os.Stat("assets/base_map_w.gif"); !os.IsExist(err) {
		res, _ := http.Get("http://www.kmoni.bosai.go.jp/data/map_img/CommonImg/base_map_w.gif")
		body, _ := io.ReadAll(res.Body)
		file, _ := os.Create("assets/base_map_w.gif")
		_, _ = file.Write(body)
	}

	if _, err := os.Stat("assets/nied_jma_s_w_scale.gif"); !os.IsExist(err) {
		res, _ := http.Get("http://www.kmoni.bosai.go.jp/data/map_img/ScaleImg/nied_jma_s_w_scale.gif")
		body, _ := io.ReadAll(res.Body)
		file, _ := os.Create("assets/nied_jma_s_w_scale.gif")
		_, _ = file.Write(body)
	}
}

func main() {
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

func DisplayName(s *discordgo.Session, guildID string, userID string) string {
	m, _ := s.GuildMember(guildID, userID)
	if m.Nick != "" {
		return m.Nick
	} else if m.User.Bot {
		return m.User.Username
	}
	return m.User.GlobalName
}

func onReady(_ *discordgo.Session, _ *discordgo.Ready) {
	fmt.Println("Bot is ready")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if strings.HasPrefix(m.Content, prefix) {
		commandName := strings.Split(m.Content, " ")[0][len(prefix):]
		args := strings.Split(m.Content, " ")[1:]
		if command, exist := commandMap[commandName]; exist {
			command.Execute(s, m, args)
		} else {
			_, _ = s.ChannelMessageSend(m.ChannelID, "コマンドが見つかりませんでした。")
		}
	} else {
		if rnd, _ := rand.Int(rand.Reader, big.NewInt(100)); rnd.Int64() < 10 {
			msgs, _ := s.ChannelMessages(m.ChannelID, 3, "", "", "")
			sort.Slice(msgs, func(i, j int) bool {
				return msgs[j].Timestamp.After(msgs[i].Timestamp)
			})

			var messages []openai.ChatCompletionMessage
			messages = append(messages, openai.ChatCompletionMessage{
				Role: openai.ChatMessageRoleSystem,
				Content: strings.Join([]string{
					fmt.Sprintf("あなたの名前は%sです。", DisplayName(s, m.GuildID, s.State.User.ID)),
					"あなたは人が複数人いるグループチャットで会話を取得し、会話できます。",
					"あなたは与えられた会話に対して自然な返信を生成してください。",
					"会話の意味や意図が理解できない場合は、適当に反応してください。",
					"会話の口調は友達と話しているような感じでお願いします。",
					"返信は短く、自然であることを心がけてください。",
					"返信はテキスト形式で返してください。",
					"テンプレートは以下の通りです。",
					"名前: 会話内容",
					"以下は例です。",
					"佐藤: こんにちは",
					"山田: <@佐藤>こんにちは",
					"<@名前>はそのユーザーに対するメンションです。その人について言及する際は名前だけを使ってください。",
					`例えば"<@山田>"と書いてあった場合、その人の名前は"山田"です。`,
					`"佐藤: <@Johnny>"と書いてあった場合は"佐藤"が"Johnny"に言及しています。`,
					"以下は例です。",
					"input: 佐藤: こんにちは",
					"output: こんにちは",
				}, "\n"),
			})

			for _, msg := range msgs {
				for _, user := range msg.Mentions {
					msg.Content = strings.ReplaceAll(msg.Content, fmt.Sprintf(func() string {
						if m.Author.Bot {
							return "<@!%s>"
						} else {
							return "<@%s>"
						}
					}(), user.ID), fmt.Sprintf("<@%s>", DisplayName(s, m.GuildID, user.ID)))
				}
				if msg.Author.ID == s.State.User.ID {
					messages = append(messages, openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleAssistant,
						Name:    msg.Author.ID,
						Content: msg.Content,
					})
				} else {
					messages = append(messages, openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleUser,
						Name:    msg.Author.ID,
						Content: fmt.Sprintf("%s: %s", DisplayName(s, m.GuildID, msg.Author.ID), msg.Content),
					})
				}
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
			_, _ = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content:   resp.Choices[0].Message.Content,
				Reference: m.Reference(),
				AllowedMentions: &discordgo.MessageAllowedMentions{
					RepliedUser: false,
				},
			})
		}
	}

	//var onMention = false
	//for _, user := range m.Mentions {
	//	if user.ID == s.State.User.ID {
	//		onMention = true
	//		break
	//	}
	//}
	//if onMention {
	//}
}
