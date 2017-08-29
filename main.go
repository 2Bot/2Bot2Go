package main

import (
	"log"
	"io/ioutil"
	"strings"
	"time"
	"os"
	"fmt"
	"syscall"
	"os/signal"
	"regexp"
	"github.com/bwmarrin/discordgo"
	"github.com/BurntSushi/toml"	
)

type config struct {
	Token string `toml:"token"`
	Prefix string `toml:"prefix"`
}

var (
	dg *discordgo.Session
	conf = &config{}
	emojiRegex   = regexp.MustCompile("<:.*?:(.*?)>")	
	loginTime time.Time
)

func createConfig() error {
	fmt.Println("Welcome and thanks for downloading 2Bot2Go!")
	fmt.Println("As I couldn't find a file called config.toml, I'll assume this is your first time starting the bot, so lets get it setup!")

	err := inputToken()
	if err != nil {
		return err
	}
	
	inputPrefix()

	file, err := os.OpenFile("config.toml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		fmt.Println("Error creating config\n", err)
		return err
	}
	defer file.Close()

	err = toml.NewEncoder(file).Encode(conf)
	if err != nil {
		fmt.Println("Error creating config\n", err)
		return err		
	}

	return nil
}

func inputToken() error {
	fmt.Println("\nFirst, I'll need your user token. To do that, follow these instructions:")
	fmt.Println("1. Type Ctrl-Shift-i\n2. Click on the tab labelled 'Application'\n3. Click 'Local Storage', and then https://discordapp.com")
	fmt.Println("4. Then copy paste the long string of random characters here, but WITHOUT THE QUOTATION MARKS! Thats very important")
	
	fmt.Print("\nPaste your token here: ")
	fmt.Scanln(&conf.Token)
	err := testLogin()
	if err != nil {
		fmt.Println("There was an issue logging you in. Check your token and try again")
		return err
	}

	return nil
}

func inputPrefix() {
	fmt.Println("\nNext up, I'll need a prefix of your preference! This will be used to call your commands.")
	fmt.Println("Example: prefix => ||\n||help => shows the help menu")
	fmt.Print("\nType your chosen prefix here: ")
	fmt.Scanln(&conf.Prefix)
	fmt.Print("Do you want a space at the end of your prefix? (y/n) ")
	var ws string
	for {
		fmt.Scanln(&ws)
		switch strings.ToLower(ws) {
			case "y":
				conf.Prefix += " "
				return
			case "n":
				return
			default:
				fmt.Println("Please type y or n")
		}
	}

}

func testLogin() error {
	fmt.Println("Trying to login... ")

	var err error
	dg, err = discordgo.New(conf.Token)
	if err != nil {
		return err
	}

	_, err = dg.User("@me")
	if err != nil {
		return err
	}
	
	fmt.Println("\rToken is valid!")
	return nil
}

func loadConfig() error {
	bytes, err := ioutil.ReadFile("config.toml")
	if err != nil {
		return err
	}

	_, err = toml.Decode(string(bytes), conf)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	
	_, err := toml.DecodeFile("config.toml", &conf)
	if os.IsNotExist(err) {
		err = createConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}else {
		err = loadConfig()
		if err != nil || conf.Prefix == "" || conf.Token == "" {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if dg == nil {
		if err = testLogin(); err != nil {
			fmt.Println(err)
		}
	}

	loginTime = time.Now()
	err = dg.Open()
	if err != nil {
		log.Fatalln(err)
	}

	dg.AddHandlerOnce(ready)
	dg.AddHandler(message)

	loadCommands()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func ready(s *discordgo.Session, m *discordgo.Ready) {
	fmt.Printf("Log-in successful! Log-in time: %.2f\n", time.Since(loginTime).Seconds())
	fmt.Printf("Joined %d guilds\n", len(m.Guilds))
}

func message(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != s.State.User.ID || !strings.HasPrefix(m.Content, conf.Prefix) || len(strings.TrimPrefix(m.Content, conf.Prefix)) == 0 {
		return
	}

	parseCommand(s, m, strings.TrimPrefix(m.Content, conf.Prefix))
}