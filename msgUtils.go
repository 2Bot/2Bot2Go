package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func msgGit(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	s.ChannelMessageEdit(m.ChannelID, m.Message.ID, "Check out 2Bot2Go, the Selfbot by Strum355 https://github.com/Strum355/2Bot2Go\nGive it star to make my creators day! ⭐")
}

func msgGame(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	err := s.UpdateStatus(0, strings.Join(msglist, " "))
	if err != nil {
		s.ChannelMessageEdit(m.ChannelID, m.Message.ID, "Error setting game")
		return
	}
	s.ChannelMessageEdit(m.ChannelID, m.Message.ID, "Game set to `"+strings.Join(msglist, " ")+"`")
}
