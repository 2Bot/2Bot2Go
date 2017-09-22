package main

import (
	"io/ioutil"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"fmt"
)

func msgImageRecall(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	if len(msglist)  < 1 {
		s.ChannelMessageEdit(m.ChannelID, m.Message.ID, "Available sub-commands for `image`:\n`save`, `delete`, `recall`, `list`, `status`\n"+
			"Type `"+conf.Prefix+"help image` to see more info about this command")
		return
	}

	switch msglist[0] {
	case "recall":
		fimageRecall(s, m, msglist)
	case "save":
		fimageSave(s, m, msglist)
	case "list":
		fimageList(s, m, msglist)
	case "delete":
		fimageDelete(s, m, msglist)
	case "info":
		fimageInfo(s, m, msglist)
	}
}

func fimageRecall(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	resp, err := http.Get("http://strum355.netsoc.co/inServer?id="+m.Author.ID)
	if err != nil {
		errorLog.Println(err)
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorLog.Println(err)
		return
	}

	if string(bytes) == "1" {
		fmt.Println("in server")
		return
	}
	fmt.Println("not in server")
}

func fimageSave(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	
}

func fimageList(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	
}

func fimageDelete(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	
}

func fimageInfo(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	
}