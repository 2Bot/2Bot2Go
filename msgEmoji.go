package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//Thanks to iopred
func emojiFile(s string) string {
	var found string
	var filename string
	for _, r := range s {
		if filename != "" {
			filename = fmt.Sprintf("%s-%x", filename, r)
		} else {
			filename = fmt.Sprintf("%x", r)
		}

		if _, err := os.Stat(fmt.Sprintf("emoji/%s.png", filename)); err == nil {
			found = filename
		} else if found != "" {
			return found
		}
	}
	return found
}

func msgEmoji(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	submatch := emojiRegex.FindStringSubmatch(msglist[0])

	if len(submatch) == 3 {
		emojiID := submatch[2]

		resp, err := http.Get(fmt.Sprintf("https://cdn.discordapp.com/emojis/%s.png", emojiID))
		if err != nil {
			errorLog.Println("Custom emoji err:", err.Error())
			return
		}
		defer resp.Body.Close()

		s.ChannelFileSend(m.ChannelID, "emoji.png", resp.Body)
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		return
	}

	emoji := emojiFile(msglist[0])
	if emoji != "" {
		file, err := os.Open(fmt.Sprintf("emoji/%s.png", emoji))
		if err != nil {
			errorLog.Println("Twemoji emoji err:", err.Error())
			return
		}
		defer file.Close()

		s.ChannelFileSend(m.ChannelID, "emoji.png", file)

		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func msgFindEmoji(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	var emojiName string
	submatch := emojiRegex.FindStringSubmatch(strings.Join(msglist, " "))

	if len(submatch) < 2 {
		emojiName = msglist[0]
	} else {
		emojiName = submatch[1]
	}

	var emojis []string
	var lenEmojiNames int
	var lenEmoji int
	var done bool
	for _, guild := range s.State.Guilds {
		for _, emoji := range guild.Emojis {
			if strings.Contains(strings.ToLower(emoji.Name), strings.ToLower(emojiName)) {
				if !done && lenEmojiNames+len("<:"+emoji.APIName()+">") >= 900 {
					done = true
					lenEmoji = 0
				}
				if done {
					lenEmoji++
					continue
				}
				lenEmojiNames += len("<:" + emoji.APIName() + ">")
				emojis = append(emojis, "<:"+emoji.APIName()+">")
			}
		}
	}

	emojisEmbed := []*discordgo.MessageEmbedField{
		{Name: "​", Value: strings.Join(emojis, " ") + func() string {
			if lenEmoji > 0 {
				return fmt.Sprintf(" and %d more...", lenEmoji)
			}
			return "​"
		}()},
	}

	userColor := s.State.UserColor(s.State.User.ID, m.ChannelID)

	if len(emojis) == 0 {
		edit := newEdit(s, m, userColor)
		edit.setTitle("No emojis found!")
		edit.send()
		return
	}

	_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:      m.Message.ID,
		Channel: m.ChannelID,
		Content: &content,
		Embed: &discordgo.MessageEmbed{
			Title: "Emojis with the substring `" + emojiName + "`",

			Color: userColor,

			Fields: emojisEmbed,
		},
	})
	if err != nil {
		errorLog.Println(err)
	}
}
