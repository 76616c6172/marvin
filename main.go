package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var WELCOME_MSG = "I think you ought to know I'm feeling very depressed.. what's your p(doom|agi)?\n"
var SERVER_ID = "1027322537622913044"
var DOOM_ROLES = map[string]string{
	"50%": "1027338125820833802",
	"20%": "1027338737337761842",
}

var USER_IS_USING_INTERACTIVE_CMD = map[string]bool{}

const (
	DEFAULT = "\033[0m"
	RED     = "\033[31m"
	GREEN   = "\033[32m"
)

func set_log_file(s string) {
	log_file, err := os.OpenFile(s, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		os.Exit(1)
	}
	log.SetOutput(log_file)
}

func print_and_log_error(s string, err error) {
	fmt.Println(s, err)
	log.Println(s, err)
}

func set_discord_api_token() string {
	if len(os.Args) < 2 || len(os.Args) > 2 {
		print_and_log_error("error: wrong argument count. Need authorization token.", errors.New("wrong arg count"))
		os.Exit(1)
	}
	return os.Args[1]
}

// Treat the next user message content as if it is the reply message to an interactive command
func handle_interactive_command(s *discordgo.Session, m *discordgo.MessageCreate) {
	USER_IS_USING_INTERACTIVE_CMD[m.Author.ID] = false

	user_input := m.Content
	if !strings.Contains(user_input, "%") {
		user_input += "%"
	}

	if _, exists := DOOM_ROLES[user_input]; !exists {
		_, err := s.ChannelMessageSend(m.ChannelID, "That's not a valid number from my list..")
		if err != nil {
			print_and_log_error("error sending message", err)
		}
		return
	}

	// Unassign all other doom roles
	for _, role := range DOOM_ROLES {
		err := s.GuildMemberRoleRemove(SERVER_ID, m.Author.ID, role)
		if err != nil {
			print_and_log_error("error unassigning role", err)
		}

	}
	// Assign new doom role
	new_role_id := DOOM_ROLES[user_input]
	err := s.GuildMemberRoleAdd(SERVER_ID, m.Author.ID, new_role_id)
	if err != nil {
		print_and_log_error("error assigning role", err)
	}
	_, err = s.ChannelMessageSend(m.ChannelID, "Thanks.")
	if err != nil {
		print_and_log_error("error sending message", err)
	}
}

func handle_server_event(s *discordgo.Session, m *discordgo.MessageCreate) {

	if USER_IS_USING_INTERACTIVE_CMD[m.Author.ID] {
		handle_interactive_command(s, m)
	} else if m.Message.Content == "/choose" || m.Type == discordgo.MessageTypeGuildMemberJoin {

		new_message := m.Author.Mention() + " " + WELCOME_MSG
		new_message += "\n**Choose from:**\n"
		for role_name, _ := range DOOM_ROLES {
			new_message += fmt.Sprintf("%s\n", role_name)
		}

		_, err := s.ChannelMessageSend(m.ChannelID, new_message)
		if err != nil {
			print_and_log_error("error sending message", err)
		}

		USER_IS_USING_INTERACTIVE_CMD[m.Author.ID] = true
	}
}

func main() {
	set_log_file("marvin.log")
	api_token := set_discord_api_token()

	discord_session, err := discordgo.New("Bot " + api_token)
	if err != nil {
		print_and_log_error("error creating session", err)
		os.Exit(1)
	}
	discord_session.Identify.Intents = discordgo.IntentsAll // receive all server events

	discord_session.AddHandler(handle_server_event)

	if err = discord_session.Open(); err != nil {
		print_and_log_error("error opening connection", err)
		os.Exit(1)
	}

	fmt.Printf("%s[marvin] %srunning..\n", GREEN, DEFAULT)
	keep_running := make(chan os.Signal, 1)
	signal.Notify(keep_running, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-keep_running
	discord_session.Close()
	fmt.Printf("\n%s[marvin] %sshutdown gracefully.\n", GREEN, DEFAULT)
}
