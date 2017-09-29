package main

import (
	"github.com/bwmarrin/discordgo"
)

type editComplex struct {
	Content *string
	Embed   *discordgo.MessageEmbed

	MessageCreate *discordgo.MessageCreate
	Session       *discordgo.Session
}

func newEdit(s *discordgo.Session, m *discordgo.MessageCreate, color int) *editComplex {
	return &editComplex{
		Content: &content,

		MessageCreate: m,
		Session:       s,

		Embed: &discordgo.MessageEmbed{
			Color: color,
		},
	}
}

func (e *editComplex) setTitle(title string) {
	e.Embed.Title = title
}

func (e *editComplex) setImage(url string) {
	e.Embed.Image = &discordgo.MessageEmbedImage{
		URL: url,
	}
}

func (e *editComplex) setFields(fields []*discordgo.MessageEmbedField) {
	e.Embed.Fields = fields
}

func (e *editComplex) setDescription(text string) {
	e.Embed.Description = text
}

func (e *editComplex) send() *discordgo.Message {
	msg, err := e.Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel: e.MessageCreate.ChannelID,
		ID:      e.MessageCreate.Message.ID,

		Content: &content,

		Embed: &discordgo.MessageEmbed{
			Color:       e.Embed.Color,
			Title:       e.Embed.Title,
			Description: e.Embed.Description,

			Fields: e.Embed.Fields,

			Image: e.Embed.Image,
		},
	})

	if err != nil {
		errorLog.Println(err)
		return nil
	}

	return msg
}
