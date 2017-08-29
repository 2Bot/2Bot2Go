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
	found := ""
	filename := ""
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
	fmt.Println(submatch)
	if len(submatch) == 3 {
		emojiID := submatch[2]

		resp, err := http.Get(fmt.Sprintf("https://cdn.discordapp.com/emojis/%s.png", emojiID))
		if err != nil {
			fmt.Println("Custom emoji err:", err.Error())
			return
		}
		defer resp.Body.Close()

		s.ChannelFileSend(m.ChannelID, "emoji.png", resp.Body)

		if m != nil {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}
	} else {
		emoji := emojiFile(msglist[0])
		if emoji != "" {
			file, err := os.Open(fmt.Sprintf("emoji/%s.png", emoji))
			if err != nil {
				fmt.Println("Twemoji emoji err:", err.Error())
				return
			}
			defer file.Close()

			s.ChannelFileSend(m.ChannelID, "emoji.png", file)

			if m != nil {
				s.ChannelMessageDelete(m.ChannelID, m.ID)
			}
		}
	}
	return
}

func msgFindEmoji(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	var emojiName string
	submatch := emojiRegex.FindStringSubmatch(strings.Join(msglist, " "))

	if len(submatch) < 2 {
		emojiName = msglist[0]
	} else {
		emojiName = submatch[1]
	}

	emojis := []string{}
	lenEmojiNames := 0
	lenEmoji := 0
	done := false
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

	if len(emojis) == 0 {
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: "No emojis found!",
		})
		return
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: "Emojis with the name " + emojiName,

		Fields: emojisEmbed,
	})
	if err != nil {
		fmt.Println(err)
	}
}
