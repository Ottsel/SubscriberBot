package main

import (
	"flag"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	var (
		Token = flag.String("t", "", "Discord Authentication Token")
	)
	flag.Parse()

	/*
		Create session and connect to Discord API
	*/
	session, err := discordgo.New("Bot " + *Token)
	if logError(err, "Unable to authenticate with discord API") {
		return
	}
	session.AddHandler(messageCreate)
	err = session.Open()
	if logError(err, "Unable to open session") {
		return
	}
	log.Println("Successfully connected to Discord API")

	//Start media scan.
	go scanForNewMedia(session)

	//Persist
	<-make(chan struct{})
	return
}

/*
MessageCreate event fired when someone sends a message.
*/
func messageCreate(session *discordgo.Session, messageCreate *discordgo.MessageCreate) {
	if strings.ToLower(messageCreate.Content) == "subscribe" {
		if subscribeUser(session, messageCreate.Author, messageCreate.ChannelID) {
			log.Println("Subscribed user:", messageCreate.Author.Username, "on channel", messageCreate.ChannelID)
			succeeded, _ := sendMessage(session, messageCreate.ChannelID, "Success! You are now subscribed.")
			if !succeeded {
				return
			}
		}
	}
}

/*
Goroutine scans for new files in each user's temp media directory.
*/
func scanForNewMedia(session *discordgo.Session) {
	for {
		time.Sleep(5)
		for i, row := range getUsersAndChannels() {

			//Skip header row of CSV file.
			if i == 0 {
				continue
			}

			userID := row[0]
			channelID := row[1]

			_, err := os.Stat("media/" + userID)
			if os.IsNotExist(err) {
				err := os.MkdirAll("media/"+userID, os.ModePerm)
				if logError(err, "Could not create directory") {
					os.Exit(1)
				}
				log.Println("Media directory created for user:", userID)
			}

			files, err := ioutil.ReadDir("./media/" + userID + "/")
			if logError(err, "Could not read directory") {
				os.Exit(1)
			}
			for _, f := range files {
				if f.Name() == "urls.csv" {
					continue
				}
				postMedia(session, userID, channelID, "./media/"+userID+"/", f.Name())
			}
		}
	}
}

/*
Post media in the target user's DM channel.
*/
func postMedia(session *discordgo.Session, userID string, channelID string, mediaPath string, fileName string) {

	r, err := os.Open(mediaPath + fileName)
	if logError(err, "Could not read media") {
		os.Exit(1)
	}

	message, err := session.ChannelFileSend(channelID, fileName, r)
	if logError(err, "Could not post media") {
		os.Exit(1)
	}
	log.Println("Posted new media for user:", userID)

	if storeMediaUrl(userID, fileName, message.Attachments[0].URL) {
		err = os.Remove(mediaPath + fileName)
		if logError(err, "Could not delete file") {
			os.Exit(1)
		}
		log.Println("Removed cached media")
	}
}

/*
Do something with the media URL
*/
func storeMediaUrl(userID string, name string, url string) bool {
	if !writeToCSV("./media/"+userID+"/urls.csv", [][]string{{"Name", "URL"}, {name, url}}) {
		file, err := os.Open("./media/" + userID + "/urls.csv")
		if logError(err, "Could not open file") {
			os.Exit(1)
		}

		defer file.Close()

		appendCSV(file, [][]string{{name, url}})
	}
	log.Println("Saved media URL to file")
	return true
}
