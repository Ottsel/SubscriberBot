package main

import (
	"encoding/csv"
	"github.com/bwmarrin/discordgo"
	"os"
)

/*
Save the user's ID and DM channel ID to the user's file.
*/
func subscribeUser(session *discordgo.Session, user *discordgo.User, channelID string) bool {
	//Check if the users file exists.
	if !writeToCSV("users.csv", [][]string{{"UserID", "ChannelID"}, {user.ID, channelID}}) {
		{
			file, err := os.OpenFile("users.csv",os.O_CREATE,os.ModePerm)
			if logError(err, "Could not open file") {
				os.Exit(1)
			}

			defer file.Close()

			//If the user isn't already subscribed, subscribe them.
			if !userAlreadySubscribed(file, user, channelID) {
				appendCSV(file, [][]string{{user.ID, channelID}})
			} else {
				sendMessage(session, channelID, "You are already subscribed.")
				return false
			}
		}
	}
	return true
}

/*
Check if the user is already present in the user's file with the same channel ID.
*/
func userAlreadySubscribed(file *os.File, user *discordgo.User, channelID string) bool {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if logError(err, "Could not read file") {
		os.Exit(1)
	}

	for _, row := range records {
		if user.ID == row[0] && channelID == row[1] {
			return true
		}
	}
	return false
}
func getUsersAndChannels() [][]string {
	file, openErr := os.OpenFile("users.csv",os.O_CREATE,os.ModePerm)
	if logError(openErr, "Could not open file") {
		os.Exit(1)
	}
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if logError(err, "Could not read file") {
		os.Exit(1)
	}

	return records
}
