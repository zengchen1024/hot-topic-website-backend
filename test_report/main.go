package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/opensourceways/hot-topic-website-backend/common/infrastructure/mongodb"
	"github.com/opensourceways/hot-topic-website-backend/hottopicmanagement/domain"
)

func handleHotTopic() {
	fmt.Println("-------- handle hot topic --------")
	// add hot topic
	file, err := os.Open("topics.json")
	if err != nil {
		log.Fatalf("无法打开文件: %v", err)
	}
	defer file.Close()

	fmt.Println("Load json file successful.")

	//
	var hotpic domain.HotTopic

	// decode json
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&hotpic); err != nil {
		log.Fatalf("Can not decode JSON: %v", err)
	}

	fmt.Printf("Decode json successful.")
	fmt.Printf("Topic: v%", hotpic.Title)
}

func main() {
	err := mongodb.Init(&mongodb.Config{
		Conn:    "mongodb://localhost:27017",
		DBName:  "community-hot-topic",
		CAFile:  "",
		Timeout: 10,
	})
	if err != nil {
		fmt.Printf("Connect mongo db failed. err %s", err)
		return
	}
	defer mongodb.Close()

}
