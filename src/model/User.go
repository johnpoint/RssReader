package model

import (
	"errors"
	"fmt"
)

type User struct {
	ID        int64 `gorm:"AUTO_INCREMENT"`
	Mail      string
	Password  string
	subscribe []subscribe
}

type Read struct {
	ID  int64 `gorm:"AUTO_INCREMENT"`
	PID int64
	UID int64
}

type Star struct {
	ID  int64 `gorm:"AUTO_INCREMENT"`
	PID int64
	UID int64
}

type subscribe struct {
	ID  int64 `gorm:"AUTO_INCREMENT"`
	UID int64
	FID int64
}

func (u *User) GetSub() error {
	if u.ID == 0 {
		return errors.New("Incomplete parameters")
	}
	db, err := Initdatabase()
	if err != nil {
		return err
	}
	defer db.Close()
	_ = db.AutoMigrate(&subscribe{})
	subscribes := []subscribe{}
	db.Where(subscribe{UID: u.ID}).Find(&subscribes)
	u.subscribe = subscribes
	return nil
}

func (u *User) Sub() []subscribe {
	return u.subscribe
}

func (u *User) AddSub(sub int64) error {
	if u.ID == 0 {
		return errors.New("Incomplete parameters")
	}
	err := u.GetSub()
	if err != nil {
		return err
	}
	subs := u.Sub()
	for _, i := range subs {
		if i.FID == sub {
			return errors.New("Already subscribed")
		}
	}
	db, err := Initdatabase()
	if err != nil {
		return err
	}
	defer db.Close()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	_ = tx.AutoMigrate(&subscribe{})
	subscribe := subscribe{UID: u.ID, FID: sub}
	if err := tx.Create(&subscribe).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	f := Feed{ID: sub}
	_ = f.Get()
	f.Num = f.Num + 1
	_ = f.save()
	return nil
}

func (u *User) DelSub(sub int64) error {
	if u.ID == 0 {
		return errors.New("Incomplete parameters")
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
		return errors.New("Feed does not exist")
	}
	db, err := Initdatabase()
	if err != nil {
		return err
	}
	defer db.Close()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	_ = tx.AutoMigrate(&subscribe{})
	if err := tx.Where(subscribe{FID: sub, UID: u.ID}).Delete(subscribe{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	f := Feed{ID: sub}
	_ = f.Get()
	p := f.Post()
	for _, i := range p {
		_ = u.UnRead(i.ID)
	}
	f.Num = f.Num - 1
	if f.Num <= 0 {
		_ = f.Detele()
	} else {
		_ = f.save()
	}
	return nil
}

func (u *User) Get() error {
	db, err := Initdatabase()
	if err != nil {
		return err
	}
	defer db.Close()
	_ = db.AutoMigrate(&User{})
	Users := []User{}
	db.Where(User{Mail: u.Mail, ID: u.ID}).Find(&Users)
	if len(Users) == 0 {
		return errors.New("Not Found")
	}
	u.ID = Users[0].ID
	u.Password = Users[0].Password
	u.Mail = Users[0].Mail
	u.GetSub()
	return nil
}

func (u *User) New() error {
	if u.Mail == "" || u.Password == "" {
		return errors.New("Incomplete parameters")
	}
	err := u.Get()
	if err == nil {
		return errors.New("Email has been used")
	}
	db, err := Initdatabase()
	if err != nil {
		return err
	}
	defer db.Close()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	_ = tx.AutoMigrate(&User{})
	if err := tx.Create(u).Error; err != nil {
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
	db, err := Initdatabase()
	if err != nil {
		return err
	}
	defer db.Close()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		fmt.Println(tx.Error)
		return tx.Error
	}
	_ = tx.AutoMigrate(&User{})
	where := User{ID: u.ID}
	if err := tx.Model(&where).Updates(u).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (u *User) ReadPost() ([]int64, error) {
	if u.ID == 0 {
		return nil, errors.New("Incomplete parameters")
	}
	db, err := Initdatabase()
	if err != nil {
		return []int64{}, err
	}
	defer db.Close()
	_ = db.AutoMigrate(&Read{})
	reads := []Read{}
	db.Where(Read{UID: u.ID}).Find(&reads)
	readList := []int64{}
	for _, i := range reads {
		readList = append(readList, i.PID)
	}
	return readList, nil
}

func (u *User) Read(pid int64) error {
	if pid == 0 {
		return errors.New("Incomplete parameters")
	}
	p := Post{ID: pid}
	err := p.Get()
	if err != nil {
		return err
	}
	db, err := Initdatabase()
	if err != nil {
		return err
	}
	defer db.Close()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}
	_ = tx.AutoMigrate(&Read{})
	if err := tx.Create(&Read{UID: u.ID, PID: pid}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (u *User) UnRead(pid int64) error {
	if pid == 0 {
		return errors.New("Incomplete parameters")
	}
	p := Post{ID: pid}
	err := p.Get()
	if err != nil {
		return err
	}
	db, err := Initdatabase()
	if err != nil {
		return err
	}
	defer db.Close()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	_ = tx.AutoMigrate(&Read{})
	if err := tx.Where(Read{UID: u.ID, PID: pid}).Delete(Read{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
