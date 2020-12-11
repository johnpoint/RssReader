package model

import (
	"encoding/json"
	"errors"
	"github.com/gilliek/go-opml/opml"
	"log"
)

type User struct {
	ID         int64 `gorm:"autoIncrement"`
	Mail       string
	Password   string
	opml       opml.OPML
	Reads      []Read      `gorm:"foreignKey:UID;constraint:OnDelete:CASCADE;"`
	subscribes []subscribe `gorm:"foreignKey:UID;constraint:OnDelete:CASCADE;"`
	ReadNum    int64
}

type Read struct {
	ID  int64 `gorm:"AUTO_INCREMENT"`
	PID int64
	UID int64
}

type subscribe struct {
	ID  int64 `gorm:"AUTO_INCREMENT"`
	UID int64
	FID int64
}

func (u *User) Export() error {
	return nil
}

func (u *User) Import(opmlStr string) error {
	type errItem struct {
		Url  string
		Info string
	}
	doc, err := opml.NewOPML([]byte(opmlStr))
	if err != nil {
		return errors.New("opml parsing error:" + err.Error())
	}
	e := 0
	var errorItem []errItem
	for _, i := range doc.Body.Outlines {
		f := Feed{Url: i.XMLURL}
		err := f.Get([]string{"id", "url"})
		log.Println("import:" + i.XMLURL)
		if err != nil && err.Error() == "not Found" {
			err = nil
			err := f.New()
			if err != nil {
				log.Println("new:" + err.Error() + "/" + i.XMLURL)
				errorItem = append(errorItem, errItem{Url: i.XMLURL, Info: err.Error()})
				e = 1
				continue
			}
		}
		_ = f.Get([]string{"id"})
		err = nil
		err = u.AddSub(f.ID)
		if err != nil {
			errorItem = append(errorItem, errItem{Url: i.XMLURL, Info: err.Error()})
			e = 1
		}
		log.Println("imported:" + i.XMLURL)
	}
	if e == 0 {
		return nil
	}
	jsonMsg, _ := json.Marshal(errorItem)
	return errors.New("Imported successfully, but something went wrong:\n" + string(jsonMsg))
}

func (u *User) GetSub() error {
	if u.ID == 0 {
		return errors.New("incomplete parameters")
	}
	if db == nil {
		return errors.New("database connection failed")
	}
	var subscribes []subscribe
	db.Where(subscribe{UID: u.ID}).Find(&subscribes)
	u.subscribes = subscribes
	return nil
}

func (u *User) Sub() []subscribe {
	return u.subscribes
}

func (u *User) AddSub(sub int64) error {
	if u.ID == 0 {
		return errors.New("incomplete parameters")
	}
	err := u.GetSub()
	if err != nil {
		return err
	}
	subs := u.Sub()
	for _, i := range subs {
		if i.FID == sub {
			return errors.New("already subscribed")
		}
	}
	if db == nil {
		return errors.New("database connection failed")
	}
	// defer db.Close()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		l := Log{Type: "DB", Level: 1, Message: tx.Error.Error()}
		_ = l.New()
		return tx.Error
	}

	subscribe := subscribe{UID: u.ID, FID: sub}
	if err := tx.Create(&subscribe).Error; err != nil {
		tx.Rollback()
		l := Log{Type: "DB", Level: 1, Message: err.Error()}
		_ = l.New()
		return err
	}
	tx.Commit()
	f := Feed{ID: sub}
	_ = f.Get([]string{"num"})
	f.Num = f.Num + 1
	_ = f.save()
	return nil
}

func (u *User) DelSub(sub int64) error {
	if u.ID == 0 {
		return errors.New("incomplete parameters")
	}
	err := u.GetSub()
	if err != nil {
		return err
	}
	subs := u.Sub()
	flag := 0
	for _, i := range subs {
		if i.FID == sub {
			flag = 1
		}
	}
	if flag == 0 {
		return errors.New("feed does not exist")
	}
	if db == nil {
		return errors.New("database connection failed")
	}
	// defer db.Close()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		l := Log{Type: "DB", Level: 1, Message: tx.Error.Error()}
		_ = l.New()
		return tx.Error
	}
	if err := tx.Where(subscribe{FID: sub, UID: u.ID}).Delete(subscribe{}).Error; err != nil {
		l := Log{Type: "DB", Level: 1, Message: err.Error()}
		_ = l.New()
		tx.Rollback()
		return err
	}
	tx.Commit()
	f := Feed{ID: sub}
	_ = f.Get([]string{"id"})
	p := f.Post(-1)
	for _, i := range p {
		_ = u.UnRead(i.ID)
	}
	f.Num = f.Num - 1
	if f.Num <= 0 {
		if err := f.Delete(); err != nil {
			log.Println(err)
		}
	} else {
		_ = f.save()
	}
	return nil
}

func (u *User) Get() error {
	if db == nil {
		return errors.New("database connection failed")
	}
	// defer db.Close()
	var Users []User
	db.Where(User{Mail: u.Mail, ID: u.ID}).Find(&Users)
	if len(Users) == 0 {
		return errors.New("not Found")
	}
	u.ID = Users[0].ID
	u.Password = Users[0].Password
	u.Mail = Users[0].Mail
	if err := u.GetSub(); err != nil {
		return err
	}
	return nil
}

func (u *User) New() error {
	if u.Mail == "" || u.Password == "" {
		return errors.New("incomplete parameters")
	}
	err := u.Get()
	if err == nil {
		return errors.New("email has been used")
	}
	if db == nil {
		return errors.New("database connection failed")
	}
	// defer db.Close()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		l := Log{Type: "DB", Level: 1, Message: tx.Error.Error()}
		_ = l.New()
		return tx.Error
	}

	if err := tx.Create(&u).Error; err != nil {
		l := Log{Type: "DB", Level: 1, Message: err.Error()}
		_ = l.New()
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (u *User) VerPassword(getPassword string) bool {
	if u.Password == getPassword {
		return true
	}
	return false
}

func (u *User) Save() error {
	if db == nil {
		return errors.New("database connection failed")
	}
	// defer db.Close()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		l := Log{Type: "DB", Level: 1, Message: tx.Error.Error()}
		_ = l.New()
		return tx.Error
	}
	where := User{ID: u.ID}
	if err := tx.Model(&where).Where(where).Updates(u).Error; err != nil {
		l := Log{Type: "DB", Level: 1, Message: err.Error()}
		_ = l.New()
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (u *User) ReadPost() ([]int64, error) {
	if u.ID == 0 {
		return nil, errors.New("incomplete parameters")
	}
	if db == nil {
		return []int64{}, errors.New("database connection failed")
	}
	// defer db.Close()
	var reads []Read
	db.Where(Read{UID: u.ID}).Find(&reads)
	var readList []int64
	for _, i := range reads {
		readList = append(readList, i.PID)
	}
	return readList, nil
}

func (u *User) Read(pid int64) error {
	if pid == 0 {
		return errors.New("incomplete parameters")
	}
	p := Post{ID: pid}
	err := p.Get([]string{"id"})
	if err != nil {
		return err
	}
	if db == nil {
		return errors.New("database connection failed")
	}
	// defer db.Close()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		l := Log{Type: "DB", Level: 1, Message: tx.Error.Error()}
		_ = l.New()
		return tx.Error
	}
	if err := tx.Create(&Read{UID: u.ID, PID: pid}).Error; err != nil {
		l := Log{Type: "DB", Level: 1, Message: err.Error()}
		_ = l.New()
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (u *User) UnRead(pid int64) error {
	if pid == 0 {
		return errors.New("incomplete parameters")
	}
	p := Post{ID: pid}
	err := p.Get([]string{"id"})
	if err != nil {
		return err
	}
	if db == nil {
		return errors.New("database connection failed")
	}
	// defer db.Close()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		l := Log{Type: "DB", Level: 1, Message: tx.Error.Error()}
		_ = l.New()
		return tx.Error
	}

	if err := tx.Where(Read{UID: u.ID, PID: pid}).Delete(Read{}).Error; err != nil {
		l := Log{Type: "DB", Level: 1, Message: err.Error()}
		_ = l.New()
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
