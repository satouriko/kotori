package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"time"
)

const (
	CommentBonus = 50
)

type Index struct {
	ID    uint `gorm:"AUTO_INCREMENT"`
	Class string `gorm:"not null"`
	Title string
	Attr  string
}

type User struct {
	ID      uint `gorm:"AUTO_INCREMENT"`
	Name    string
	Email   string `gorm:"not null;unique"`
	Website string
	Rank    int64
	Honor   string
}

type Comment struct {
	ID            uint `gorm:"AUTO_INCREMENT"`
	CommentZoneID uint
	FatherID      uint
	UserID        uint
	User          User
	Content       string
	Type          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Post struct {
	ID        uint `gorm:"AUTO_INCREMENT"`
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ListComments(db *gorm.DB, commentZoneID uint, fatherID uint, offsetID uint) (comments []Comment, err error) {
	var order string
	var offset string
	if fatherID != 0 {
		order = "id asc"
		offset = "id > ?"
	} else {
		order = "id desc"
		offset = "id < ?"
	}
	if offsetID == 0 {
		err = db.Where("comment_zone_id = ?", commentZoneID).Where("father_id = ?", fatherID).
			Preload("User").Order(order).Limit(10).Find(&comments).Error
	} else {
		err = db.Where("comment_zone_id = ?", commentZoneID).Where("father_id = ?", fatherID).
			Where(offset, offsetID).
			Preload("User").Order(order).Limit(10).Find(&comments).Error
	}
	if err != nil {
		err = errors.Wrap(err, "ListComments")
		return
	}
	return
}

func SaveComment(db *gorm.DB, comment Comment) (err error) {
	var users []User
	var user_cnt uint
	err = db.Model(&User{}).Where("email = ?", comment.User.Email).Find(&users).Count(&user_cnt).Error
	if err != nil {
		err = errors.Wrap(err, "SaveComment")
		return
	}
	if user_cnt != 0 {
		comment.UserID = users[0].ID
		users[0].Name = comment.User.Name
		users[0].Website = comment.User.Website
		users[0].Rank += CommentBonus
		db.Model(&User{}).Updates(&users[0])
	} else {
		db.Create(&comment.User)
		comment.UserID = comment.User.ID
	}
	db.Set("gorm:save_associations", false).Create(&comment)
	return
}

func RemoveComment(db *gorm.DB, id uint) (err error) {
	var comment Comment
	err = db.Model(&Comment{}).Where("id = ?", id).Preload("User").First(&comment).Error
	if err != nil {
		err = errors.Wrap(err, "RemoveComment")
		return
	}
	comment.User.Rank -= CommentBonus
	db.Model(&User{}).Updates(&comment.User)
	db.Delete(&comment)
	return
}
