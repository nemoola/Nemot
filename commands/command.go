package commands

import (
	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string)
	Help() string
}

func CommandList() map[string]Command {
	return map[string]Command{
		"args": &Args{},
		"j2c":  &J2c{},
		"j2co": &J2co{},
		"c2j":  &C2j{},
		"ping": &Ping{},
		"wiki": &Wiki{},
		"help": &Help{},
	}
}
