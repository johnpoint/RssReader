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
	var feeds []model.Feed
	model.Db.Find(&feeds)
	return feeds
}

func getUpdate() {
	var t int64 = 0
	for true {
		if t == 1025 {
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
					if i.Status == 1024 {
						i.Status = -1
					} else {
						if i.Status == 0 {
							i.Status = 1
						}
						if i.Status%2 != 0 {
							i.Status--
						} else {
							i.Status *= 2
						}
					}
				} else {
					if i.Status != 1 {
						if i.Status%2 != 0 {
							i.Status++
						} else {
							i.Status /= 2
						}
					}
				}
				_ = i.Save()
			}
		}
		log.Println("=== Update Finish ===")
		t2 := time.Now()
		fmt.Println("usage time: ", t2.Sub(t1))
		time.Sleep(time.Minute * 60)
	}
}
