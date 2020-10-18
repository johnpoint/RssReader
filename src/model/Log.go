package model

type Log struct {
	ID      int64 `gorm:"autoIncrement"`
	Code    int64
	Message string
	Type    string
	Created int64 `gorm:"autoCreateTime"`
}
