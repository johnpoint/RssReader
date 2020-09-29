package spider

import (
	"log"
	"rssreader/src/model"
	"time"
)

func Spider() {
	for true {
		log.Println("Spider Running!")
		getUpdate()
		time.Sleep(time.Minute * 15)
	}
}

func getFeed() []model.Feed {
	db, err := model.Initdatabase()
	if err != nil {
		return []model.Feed{}
	}
	_ = db.AutoMigrate(&model.Feed{})
	feeds := []model.Feed{}
	db.Find(&feeds)
	return feeds
}

func getUpdate() {
	feeds := getFeed()
	log.Println("=== Update Start ===")
	for _, i := range feeds {
		err := i.Update()
		log.Println(i.Title)
		if err != nil {
			log.Print(err.Error())
		}
	}
	log.Println("=== Update Finish ===")
}
