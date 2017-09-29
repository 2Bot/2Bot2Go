package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/crypto/blake2b"
	"net/http"
	"strings"
)

func msgImageRecall(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	if len(msglist) < 1 {
		s.ChannelMessageEdit(m.ChannelID, m.Message.ID, "Available sub-commands for `image`:\n`save`, `delete`, `recall`, `list`, `status`\n"+
			"Type `"+conf.Prefix+"help image` to see more info about this command")
		return
	}

	switch msglist[0] {
	case "recall":
		fimageRecall(s, m, msglist[1:])
		/* 	case "save":
		   		fimageSave(s, m, msglist[1:])
		   	case "list":
		   		fimageList(s, m, msglist[1:])
		   	case "delete":
		   		fimageDelete(s, m, msglist[1:])
		   	case "info":
		   		fimageInfo(s, m, msglist[1:]) */
	}
}

func fimageRecall(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	prefixedImgName := m.Author.ID + "_" + strings.Join(msglist, " ")
	hash := blake2b.Sum256([]byte(prefixedImgName))
	imgFileName := hex.EncodeToString(hash[:])

	URL := fmt.Sprintf("https://api.2bot.ml/image/%s/recall/%s", m.Author.ID, imgFileName)
	resp, err := http.Get(URL)
	if err != nil {
		errorLog.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorLog.Println(fmt.Sprintf("%d %s", resp.StatusCode, URL))
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	imgURL := buf.String()

	edit := newEdit(s, m, s.State.UserColor(s.State.User.ID, m.ChannelID))
	edit.setDescription(strings.Join(msglist, " "))
	edit.setImage(imgURL)
	edit.send()
}

func fimageSave(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {
	resp, err := http.Get("https://api.2bot.ml/inServer?id=" + m.Author.ID)
	if err != nil {
		errorLog.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		errorLog.Println("Need to be in 2Bot server to use this command https://discord.gg/9T34Y6u")
		return
	}

	if len(m.Attachments) == 0 || m.Attachments[0].Height == 0 {
		errorLog.Println("Need to send an image to be saved")
		return
	}
}

func fimageList(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {

}

func fimageDelete(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {

}

func fimageInfo(s *discordgo.Session, m *discordgo.MessageCreate, msglist []string) {

}
