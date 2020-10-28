package spider

import (
	"fmt"
	"log"
	"rssreader/src/model"
	"time"
)

func Spider() {
	log.Println("Spider Running!")
	getUpdate()
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
	var t int64 = 0
	for true {
		if t == 101 {
			t = 0
		}
		t++
		feeds := getFeed()
		t1 := time.Now()
		log.Println("=== Update Start", t, " ===")
		for _, i := range feeds {
			if i.Status%t == 0 && i.Status >= 0 {
				err := i.Update()
				log.Println(i.Status, " ", i.Title)
				if err != nil {
					log.Print(err.Error())
					if i.Status == 64 {
						i.Status = -1
					} else {
						if i.Status == 0 {
							i.Status = 1
						}
						i.Status *= 2
					}
					i.Save()
				}
			}
		}
		log.Println("=== Update Finish ===")
		t2 := time.Now()
		fmt.Println("usage time: ", t2.Sub(t1))
		time.Sleep(time.Minute * 15)
	}
}
