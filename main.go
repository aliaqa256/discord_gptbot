package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	gpt3 "github.com/PullRequestInc/go-gpt3"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {
	
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}


func main() {
	

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + "token")
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {


apiKey :="token"
client := gpt3.NewClient(apiKey)
ctx := context.Background()


	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// if not m.Content starts with gpt then return
	if m.Content[0:3] != "gpt" {
		return
	} 

	var res string
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
			Prompt: []string{
				m.Content,
			},
			MaxTokens:   gpt3.IntPtr(1000),
			Temperature: gpt3.Float32Ptr(0),
		}, func(resp *gpt3.CompletionResponse) {
			fmt.Print(resp.Choices[0].Text)
			res = res + resp.Choices[0].Text
		})
		
		
		s.ChannelMessageSend(m.ChannelID,res)

	
	
	
		if err != nil {
		fmt.Println("error creating gpt session,", err)
		return
	}



}


