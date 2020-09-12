package spider

import (
	"fmt"
	"rssreader/src/model"
	"time"
)

func Spider() {
	for true {
		fmt.Println("Update Feed")
		getUpdate()
		time.Sleep(time.Minute * 15)
	}
}

func getFeed() []model.Feed {
	db := model.Initdatabase()
	_ = db.AutoMigrate(&model.Feed{})
	feeds := []model.Feed{}
	db.Find(&feeds)
	return feeds
}

func getUpdate() {
	feeds := getFeed()
	for _, i := range feeds {
		i.Update()
	}
}
