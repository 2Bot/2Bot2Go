package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
)

type config struct {
	Token  string `toml:"token"`
	Prefix string `toml:"prefix"`
}

var (
	dg         *discordgo.Session
	conf       = &config{}
	emojiRegex = regexp.MustCompile("<:(.*?):(.*?)>")
	loginTime  time.Time
	errorLog   *log.Logger
	infoLog    *log.Logger
	logF       *os.File
	// Zero width whitespace to replace message content
	content = "​"
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
	fmt.Println("\nFirst, I'll need your user token. To do that, follow these instructions:\n1. Type Ctrl-Shift-i\n2. Click on the tab labelled 'Application'\n3. Click 'Local Storage', and then https://discordapp.com")
	fmt.Print("4. Then copy paste the long string of random characters here, but WITHOUT THE QUOTATION MARKS! Thats very important\nPaste your token here: ")

	fmt.Scanln(&conf.Token)
	err := testLogin()
	if err != nil {
		fmt.Println("There was an issue logging you in. Check your token and try again")
		return err
	}

	return nil
}

func inputPrefix() {
	fmt.Println("\nNext up, I'll need a prefix of your preference! This will be used to call your commands.\nYou can choose to have a space between your prefix and command after you input your prefix")
	fmt.Println("Example: prefix => ||\n||help => shows the help menu")
	fmt.Print("\nType your chosen prefix here: ")

	fmt.Scanln(&conf.Prefix)
	var ws string
	for {
		fmt.Print("Do you want a space at the end of your prefix? (y/n) ")
		fmt.Scanln(&ws)
		switch strings.ToLower(ws) {
		case "y":
			conf.Prefix += " "
			return
		case "n":
			return
		default:
			fmt.Println("\nPlease type y or n")
		}
	}
}

func testLogin() error {
	fmt.Println("Trying to login...")
	infoLog.Println("Trying to login...")

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
	infoLog.Println("\rToken is valid!")

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

func openLog() *os.File {
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	return f
}

func main() {
	logF = openLog()
	defer logF.Close()

	log.SetOutput(logF)

	infoLog = log.New(logF, "INFO: ", log.Ldate|log.Ltime)
	errorLog = log.New(logF, "ERROR: ", log.Ldate|log.Ltime)

	infoLog.Println("log opened")

	_, err := toml.DecodeFile("config.toml", &conf)
	if os.IsNotExist(err) {
		if err = createConfig(); err != nil {
			fmt.Println(err)
			errorLog.Fatalln(err)
		}
	} else {
		if err = loadConfig(); err != nil || conf.Prefix == "" || conf.Token == "" {
			fmt.Println(err)
			errorLog.Fatalln(err)
		}
	} 

	fmt.Println("Prefix is "+conf.Prefix)
	infoLog.Println("Prefix is "+conf.Prefix)

	if dg == nil {
		if err = testLogin(); err != nil {
			fmt.Println(err)
			errorLog.Fatalln(err)
		}
	}

	loginTime = time.Now()
	if err = dg.Open(); err != nil {
		fmt.Println(err)
		errorLog.Fatalln(err)
	}
	defer dg.Close()

	prepareCommands()

	dg.AddHandlerOnce(ready)
	dg.AddHandler(message)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func ready(s *discordgo.Session, m *discordgo.Ready) {
	fmt.Printf("Log-in successful! Log-in time: %.2f\n", time.Since(loginTime).Seconds())
	fmt.Printf("Joined %d guilds\n", len(m.Guilds))
	fmt.Println("Type Ctrl+C to quit 2Bot2Go")

	infoLog.Printf("Log-in successful! Log-in time: %.2f\n", time.Since(loginTime).Seconds())
	infoLog.Printf("Joined %d guilds\n", len(m.Guilds))
}

func message(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != s.State.User.ID || !strings.HasPrefix(m.Content, conf.Prefix) || len(strings.TrimPrefix(m.Content, conf.Prefix)) == 0 {
		return
	}

	parseCommand(s, m, strings.TrimPrefix(m.Content, conf.Prefix))
}
