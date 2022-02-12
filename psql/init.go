package psql

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

const (
	DefaultDatabase = "dvdrental"
	DefaultDBUser   = "postgres"
	DefaultHost     = "localhost"
	DefaultPort     = "5432"
	DefaultPassword = "secret"
	DefaultSSLMode  = "disable"
	DefaultTimeZone = "Europe/Prague"
)

type mpaa_rating string

const (
	Enum_mpaa_rating_select = "SELECT 1 FROM pg_type WHERE typname = 'mpaa_rating';"
	Enum_mpaa_rating_sql    = "CREATE TYPE mpaa_rating AS ENUM ('G','PG','PG-13','R','NC-17');"
)

const (
	G     mpaa_rating = "G"
	PG    mpaa_rating = "PG"
	PG_13 mpaa_rating = "PG-13"
	R     mpaa_rating = "R"
	NC_17 mpaa_rating = "NC-17"
)

func (mpaa *mpaa_rating) Scan(value interface{}) error {
	*mpaa = mpaa_rating(value.([]byte))
	return nil
}
func (mpaa mpaa_rating) Value() (driver.Value, error) {
	return string(mpaa), nil
}

type User struct {
	gorm.Model
	FirstName     string `gorm:"type:varchar(100)"`
	LastName      string `gorm:"type:varchar(100)"`
	FavoriteColor string `gorm:"type:varchar(100)"`
}

type Category struct {
	gorm.Model
	Name string `gorm:"type:varchar(100)"`
}

func (Category) TableName() string {
	return "category"
}

type Film_category struct {
	gorm.Model
	Film_id     uint32 `gorm:"type:integer"`
	Category_id uint32 `gorm:"type:integer"`
	Name        string `gorm:"type:varchar(100)"`
}

func (Film_category) TableName() string {
	return "film_category"
}

type Film struct {
	gorm.Model
	Title            string      `gorm:"type:varchar(255)"`
	Description      string      `gorm:"type:text"`
	Release_year     uint8       `gorm:"type:smallint"`
	Language_id      uint8       `gorm:"type:smallint"`
	Rental_duration  uint8       `gorm:"type:smallint;default:3"`
	Rental_rate      float32     `gorm:"type:numeric(4,2);default:4.99"`
	Length           uint8       `gorm:"type:smallint"`
	Replacement_cost float32     `gorm:"type:numeric(5,2);default:19.99"`
	Rating           mpaa_rating `gorm:"type:mpaa_rating;default:G"`
	Special_features []string    `gorm:"type:text[]"`
	Fulltext         string      `gorm:"type:tsvector"`
}

func (Film) TableName() string {
	return "film"
}

type Language struct {
	gorm.Model
	Name string `gorm:"type:varchar(100)"`
}

func (Language) TableName() string {
	return "language"
}
