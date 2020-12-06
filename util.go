package main

import (
	"encoding/csv"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

//Error logging abstraction.
func logError(err error, comment string) bool {
	if err != nil {
		if comment != "" {
			log.Println(comment)
		}
		log.Println("Error:", err)
		return true
	} else {
		return false
	}
}

//Message abstraction.
func sendMessage(session *discordgo.Session, channelID string, content string) (bool, *discordgo.Message) {
	message, err := session.ChannelMessageSend(channelID, content)
	if logError(err, "Unable to send message in channel: "+channelID) {
		return false, message
	}
	return true, message
}

/*
Return true if the file was created, false if it already exists.
*/
func writeToCSV(filepath string, initData [][]string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		file, err := os.OpenFile(filepath,os.O_CREATE,os.ModePerm)
		if logError(err, "Could not create file") {
			os.Exit(1)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		err = writer.WriteAll(initData)
		if logError(err, "Could not write to file") {
			os.Exit(1)
		}
		return true
	}
	return false
}

/*
Append data to a two-column csv file.
*/
func appendCSV(file *os.File, additionalData [][]string) {

	reader := csv.NewReader(file)
	existingData, err := reader.ReadAll()
	if logError(err, "Could not read file") {
		os.Exit(1)
	}
	for _, value := range additionalData {
		existingData = append(existingData, value)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(existingData)
	if logError(err, "Could not write to file") {
		os.Exit(1)
	}
}
