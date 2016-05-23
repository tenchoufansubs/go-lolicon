package lolicon

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type Message struct {
	Date      time.Time `json:"date"`
	Content   string    `json:"content"`
	ChannelId string    `json:"channel_id"`

	Command  string `json:"command"`
	Trailing string `json:"trailing"`

	Raw *RawMessageData `json:"raw"`
}

type RawMessageData struct {
	Session *discordgo.Session       `json:"-"`
	Event   *discordgo.MessageCreate `json:"event"`
}
