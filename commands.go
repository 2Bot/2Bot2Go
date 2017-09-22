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
		Help: "Args: [emoji | name]\n\nReturns all the emojis that match the given emoji or emoji name in all the servers you are in",
		Exec: msgFindEmoji,
	}.add()
	imageCommand = command{
		Name: "image",
		Help: "Args: [save,recall,delete,list,status] [name]\n\nSave images and recall them at anytime! Everyone gets 8MB of image storage. Any name counts so long theres no `/` in it." +
		"Only you can 'recall' your saved images. There's a review process to make sure nothing illegal is being uploaded but we're fairly relaxed for the most part\n\n" +
		"Example:\n`!owo image save 2B Happy`\n2Bot downloads the image and sends it off for reviewing\n\n" +
		"`!owo image recall 2B Happy`\nIf your image was confirmed, 2Bot will send the image named `2B Happy`\n\n" +
		"`!owo image delete 2B Happy`\nThis will delete the image you saved called `2B Happy`\n\n" +
		"`!owo image list`\nThis will list your saved images along with a preview!\n\n" +
		"`!owo image status`\nShows some details on your saved images and quota",
		Exec: msgImageRecall,
	}.add()
	helpComm = command{"help",
		"", msgHelp,
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

func (c command) msgHelp(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	if len(msglist) == 0 {
		listCommands(s, m)
		return
	}
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
