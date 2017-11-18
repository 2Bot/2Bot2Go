package main

import (
	"github.com/bwmarrin/discordgo"
)

func guildDetails(channelID string, s *discordgo.Session) (*discordgo.Guild, error) {
	channelInGuild, err := s.State.Channel(channelID)
	if err != nil {
		return nil, err
	}
	guildDetails, err := s.State.Guild(channelInGuild.GuildID)
	if err != nil {
		return nil, err
	}
	return guildDetails, nil
}
