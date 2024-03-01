package main

import (
	"Nemot/commands"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"math/big"
	"os"
	"sort"
	"strings"
)

var (
	prefix     = ""
	commandMap = map[string]commands.Command{
		"args": &commands.Args{},
		"j2c":  &commands.J2c{},
		"ping": &commands.Ping{},
	}
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	prefix = os.Getenv("PREFIX")

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

	DisplayName := func() string {
		if m.Member.Nick != "" {
			return m.Member.Nick
		} else {
			return m.Author.GlobalName
		}
	}

	if strings.HasPrefix(m.Content, prefix) {
		commandName := strings.Split(m.Content, " ")[0][len(prefix):]
		args := strings.Split(m.Content, " ")[1:]
		if command, exist := commandMap[commandName]; exist {
			command.Execute(s, m, args)
		} else {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Unknown command")
		}
	} else {
		if rnd, _ := rand.Int(rand.Reader, big.NewInt(100)); rnd.Int64() < 15 {
			msgs, _ := s.ChannelMessages(m.ChannelID, 3, "", "", "")
			sort.Slice(msgs, func(i, j int) bool {
				return msgs[j].Timestamp.After(msgs[i].Timestamp)
			})

			var messages []openai.ChatCompletionMessage
			messages = append(messages, openai.ChatCompletionMessage{
				Role: openai.ChatMessageRoleSystem,
				Content: strings.Join([]string{
					fmt.Sprintf("あなたの名前は%sです。", DisplayName()),
					"あなたは会話の一部を取得できます。",
					"その与えられた会話に対して文章の意味がわからない場合は適当に返信してください。",
					"意味がわかる場合はその会話に交じるように返信してください。",
					"返信は短くお願いします。",
					"会話は以下のように与えられます。",
					"名前: 会話内容",
					"名前: <@名前>",
					"<@名前>はそのユーザーに対するメンションです。その人について言及する際は名前だけを使ってください。",
					"例えば<@山田>と書いてあった場合、その人の名前は山田です。",
				}, "\n"),
			})

			for _, msg := range msgs {
				for _, user := range msg.Mentions {
					member, _ := s.GuildMember(m.GuildID, user.ID)
					msg.Content = strings.ReplaceAll(msg.Content, fmt.Sprintf(func() string {
						if m.Author.Bot {
							return "<@!%s>"
						} else {
							return "<@%s>"
						}
					}(), user.ID), fmt.Sprintf("<@%s>", member.DisplayName()))
				}
				if msg.Author.ID == s.State.User.ID {
					messages = append(messages, openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleAssistant,
						Content: msg.Content,
					})
				} else {
					messages = append(messages, openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleUser,
						Content: fmt.Sprintf("%s: %s", DisplayName(), msg.Content),
					})
				}
				fmt.Printf("%s: %s\n", DisplayName(), msg.Content)
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
