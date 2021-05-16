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
	ReadAfter  []ReadAfter `gorm:"foreignKey:UID;constraint:OnDelete:CASCADE;"`
	subscribes []subscribe `gorm:"foreignKey:UID;constraint:OnDelete:CASCADE;"`
	Status     Status      `gorm:"foreignKey:UID;constraint:OnDelete:CASCADE;"`
	ReadNum    int64
}

type Read struct {
	PID int64 `gorm:"primaryKey"`
	UID int64 `gorm:"primaryKey"`
}

type Status struct {
	UID    int64 `gorm:"primaryKey"`
	Status int64
}

type ReadAfter struct {
	UID int64 `gorm:"primaryKey"`
	PID int64 `gorm:"primaryKey"`
}

type subscribe struct {
	UID int64 `gorm:"primaryKey"`
	FID int64 `gorm:"primaryKey"`
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
		err = u.Subscribe(f.ID)
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

func (u *User) Subscribes() (error, []subscribe) {
	if u.ID == 0 {
		return errors.New("incomplete parameters"), []subscribe{}
	}
	if Db == nil {
		return errors.New("database connection failed"), []subscribe{}
	}
	var subscribes []subscribe
	Db.Where(subscribe{UID: u.ID}).Find(&subscribes)
	u.subscribes = subscribes
	return nil, u.subscribes
}

func (u *User) Subscribe(sub int64) error {
	if u.ID == 0 {
		return errors.New("incomplete parameters")
	}
	err, subs := u.Subscribes()
	if err != nil {
		return err
	}
	for _, i := range subs {
		if i.FID == sub {
			return errors.New("already subscribed")
		}
	}
	if Db == nil {
		return errors.New("database connection failed")
	}
	// defer db.Close()
	tx := Db.Begin()
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

func (u *User) Unsubscribe(sub int64) error {
	if u.ID == 0 {
		return errors.New("incomplete parameters")
	}
	err, subs := u.Subscribes()
	if err != nil {
		return err
	}
	flag := 0
	for _, i := range subs {
		if i.FID == sub {
			flag = 1
		}
	}
	if flag == 0 {
		return errors.New("feed does not exist")
	}
	if Db == nil {
		return errors.New("database connection failed")
	}
	// defer db.Close()
	tx := Db.Begin()
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
	var postIDs []int64
	for _, i := range p {
		postIDs = append(postIDs, i.ID)
	}
	_ = u.UnRead(postIDs)
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
	if Db == nil {
		return errors.New("database connection failed")
	}
	// defer db.Close()
	var Users []User
	Db.Where(User{Mail: u.Mail, ID: u.ID}).Find(&Users)
	if len(Users) == 0 {
		return errors.New("not Found")
	}
	u.ID = Users[0].ID
	u.Password = Users[0].Password
	u.Mail = Users[0].Mail
	return nil
}

func (u *User) GetReadAfter() (error, []ReadAfter) {
	if Db == nil {
		return errors.New("database connection failed"), nil
	}
	var ReadAfters []ReadAfter
	Db.Where(ReadAfter{UID: u.ID}).Find(&ReadAfters)
	return nil, ReadAfters
}

func (u *User) AddReadAfter(pid int64) error {
	if Db == nil {
		return errors.New("database connection failed")
	}
	// defer db.Close()
	tx := Db.Begin()
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

	readAfter := ReadAfter{UID: u.ID, PID: pid}
	if err := tx.Create(&readAfter).Error; err != nil {
		tx.Rollback()
		l := Log{Type: "DB", Level: 1, Message: err.Error()}
		_ = l.New()
		return err
	}
	tx.Commit()
	return nil
}

func (u *User) RemoveReadAfter(pid int64) error {
	tx := Db.Begin()
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
	if err := tx.Where(ReadAfter{PID: pid, UID: u.ID}).Delete(ReadAfter{}).Error; err != nil {
		l := Log{Type: "DB", Level: 1, Message: err.Error()}
		_ = l.New()
		tx.Rollback()
		return err
	}
	tx.Commit()
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
	if Db == nil {
		return errors.New("database connection failed")
	}
	// defer db.Close()
	tx := Db.Begin()
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
	if Db == nil {
		return errors.New("database connection failed")
	}
	// defer db.Close()
	tx := Db.Begin()
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
	if Db == nil {
		return []int64{}, errors.New("database connection failed")
	}
	// defer db.Close()
	var reads []Read
	Db.Where(Read{UID: u.ID}).Find(&reads)
	var readList []int64
	for _, i := range reads {
		readList = append(readList, i.PID)
	}
	return readList, nil
}

func (u *User) Read(pid []int64) error {
	if len(pid) == 0 {
		return errors.New("incomplete parameters")
	}
	// defer db.Close()
	tx := Db.Begin()
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
	var reads []Read
	for i := range pid {
		_p := Read{UID: u.ID, PID: pid[i]}
		reads = append(reads, _p)
	}
	if err := tx.Create(&reads).Error; err != nil {
		l := Log{Type: "DB", Level: 1, Message: err.Error()}
		_ = l.New()
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (u *User) UnRead(pid []int64) error {
	if len(pid) == 0 {
		return errors.New("incomplete parameters")
	}
	if Db == nil {
		return errors.New("database connection failed")
	}
	// defer db.Close()
	tx := Db.Begin()
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

	if err := tx.Where("id in ?", pid).Delete(Read{}).Error; err != nil {
		l := Log{Type: "DB", Level: 1, Message: err.Error()}
		_ = l.New()
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
