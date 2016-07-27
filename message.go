package lolicon

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// Message is a message from Discord.
type Message struct {
	Date      time.Time `json:"date"`
	Content   string    `json:"content"`
	ChannelId string    `json:"channel_id"`
	UserId    string    `json:"user_id"`

	Command  string `json:"command"`
	Trailing string `json:"trailing"`

	IsAdmin bool `json:"is_admin"`

	Raw *RawMessageData `json:"raw"`
}

// RawMessageData contains the raw message received
// from Discord.
type RawMessageData struct {
	Session *discordgo.Session       `json:"-"`
	Event   *discordgo.MessageCreate `json:"event"`
}
