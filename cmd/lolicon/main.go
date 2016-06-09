package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/tenchoufansubs/go-lolicon"
)

var (
	interrupt = make(chan struct{})

	loliconUser *discordgo.User
)

func init() {
	initCache()
	initToken()

	initPlugins()
}

func main() {
	defer cache.Close()

	var (
		err error

		s *discordgo.Session
	)

	//
	// Capture SIGINT.
	//
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		// First ^C triggers a graceful shutdown.
		<-c
		close(interrupt)

		// Second ^C forces program termination.
		<-c
		os.Exit(1)
	}()

	s, err = discordgo.New(Token)
	if err != nil {
		panic(err)
	}

	s.AddHandler(onMessageCreate)
	s.AddHandler(onReady)

	log.Print("connect")

	err = s.Open()
	if err != nil {
		panic(err)
	}
	defer func() {
		err := s.Close()
		if err != nil {
			log.Print(err)
			return
		}
	}()

	<-interrupt
	log.Print("interrupt")

	for _, p := range plugins {
		log.Printf("save: %s", p.Id())

		err = p.Close()
		if err != nil {
			log.Print(err.Error())
		}
	}
}

func onReady(s *discordgo.Session, ready *discordgo.Ready) {
	log.Print("ready")

	var (
		err error
	)

	loliconUser = ready.User

	err = cache.Set("userId", loliconUser.ID)
	if err != nil {
		panic(err)
	}

	for _, p := range plugins {
		log.Printf("start: %s", p.Id())

		err = p.Open(s)
		if err != nil {
			log.Print(err.Error())
			return
		}
	}
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == loliconUser.ID {
		return
	}
	if m.Content == "" {
		return
	}

	var (
		err error
		ts  time.Time
	)
	defer func() {
		if err != nil {
			log.Print(err.Error())
		}
	}()

	ts, err = time.Parse(time.RFC3339Nano, m.Timestamp)
	if err != nil {
		return
	}

	msg := &lolicon.Message{
		Date:      ts,
		Content:   m.Content,
		ChannelId: m.ChannelID,
		UserId:    m.Author.ID,

		Raw: &lolicon.RawMessageData{
			Session: s,
			Event:   m,
		},
	}

	for _, p := range plugins {
		handler, ok := p.(lolicon.MessageHandler)
		if !ok {
			continue
		}

		var (
			done bool
		)

		done, err = handler.HandleMessage(msg)
		if err != nil {
			return
		}
		if done {
			break
		}

		if msg.Command == "help" {
			outputStr := ""
			for _, p2 := range plugins {
				helpProvider, ok := p2.(lolicon.HelpProvider)
				if !ok {
					continue
				}

				commandHelp := helpProvider.Help()
				if commandHelp == nil {
					continue
				}

				for command, helpStr := range commandHelp {
					outputStr += fmt.Sprintf("* **%s** -- %s\n", command, helpStr)
				}
			}

			_, err = s.ChannelMessageSend(msg.ChannelId, outputStr)

			break
		}
	}
}
