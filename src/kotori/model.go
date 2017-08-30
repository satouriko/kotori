package main

type Index struct {
	ID        uint `gorm:"AUTO_INCREMENT"`
	Class	  string `gorm:"not null"`
	Title	  string
	Attr      string
}

type User struct {
	ID        uint `gorm:"AUTO_INCREMENT"`
	Name	  string
	Email	  string `gorm:"not null"`
	Website   string
	Rank	  int64
	Honor	  string
}

type Comment struct {
	ID        uint `gorm:"AUTO_INCREMENT"`
	FatherID  uint
	UserID    uint
	Content   string
	Type	  string
}

type Post struct {
	ID        uint `gorm:"AUTO_INCREMENT"`
	Title	  string
	Content   string
}