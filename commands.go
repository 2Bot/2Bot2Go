package main

import "github.com/bwmarrin/discordgo"
import "strings"

type command struct {
	Name string
	Help string

	Exec func(*discordgo.Session, *discordgo.MessageCreate, []string)
}

var (
	commMap    = make(map[string]command)
	gitCommand = command{
		Name: "git",
		Help: "Args: none\n\nLinks 2Bots github page.\n\nExample:\n`!owo git`",
		Exec: msgGit,
	}.add()
	emojiCommand = command{
		Name: "",
		Help: "Args: [emoji]\n\nSends a large image of the given emoji.\n\nExample:\n`!owo :smile:`",
		Exec: msgEmoji,
	}.add()
	gameCommand = command{
		Name: "setGame",
		Help: "Args: [game]\n\nSets your current game to 'game'",
		Exec: msgGame,
	}.add()
	findEmojiCommand = command{
		Name: "findEmoji",
		Help: "Args: [emoji]\n\nReturns all the emojis that match the given emoji in all the servers you are in",
		Exec: msgFindEmoji,
	}.add()
)

//Small wrapper function to reduce clutter
func l(s string) (r string) {
	return strings.ToLower(s)
}

func parseCommand(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	msglist := strings.Fields(message)
	command := func() string {
		if strings.HasPrefix(message, " ") {
			return " " + msglist[0]
		}
		return msglist[0]
	}()

	if command == "help" {
		if len(msglist) == 2 {
			if val, ok := commMap[l(msglist[1])]; ok {
				val.helpCommand(s, m)
				return
			}
		}

		listCommands(s, m)

		return
	}

	if command == l(commMap[command].Name) {
		commMap[command].Exec(s, m, msglist[1:])
		return
	}

	//if data passed as command isnt a valid command,
	//check if its an emoji
	emojiCommand.Exec(s, m, msglist)
}

func listCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	var commands []string
	for _, val := range commMap {
		if val.Name == "" {
			continue
		}
		commands = append(commands, "`"+val.Name+"`")
	}

	userColor := s.State.UserColor(s.State.User.ID, m.ChannelID)

	s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: "Help",

		Color: userColor,

		Fields: []*discordgo.MessageEmbedField{
			{Name: "List", Value: strings.Join(commands, ", ")},
			{Name: "Info", Value: "\n\nUse `" + conf.Prefix + "help [command]` for detailed info about a command."},
		},
	})
}

func (c command) helpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	userColor := s.State.UserColor(s.State.User.ID, m.ChannelID)

	s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: "Help: " + c.Name,

		Color: userColor,

		Fields: []*discordgo.MessageEmbedField{
			{Name: "Details", Value: c.Help},
		},
	})
	return
}

func (c command) add() command {
	commMap[l(c.Name)] = c
	return c
}
