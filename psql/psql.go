package psql

import (
	"github.com/hrabalvojta/micro-dvdrental/users"
	"gorm.io/driver/postgres"
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

type User struct {
	gorm.Model
	FirstName     string `gorm:"type:varchar(100)"`
	LastName      string `gorm:"type:varchar(100)"`
	FavoriteColor string `gorm:"type:varchar(100)"`
}

// inMemUserRepository is an implementation of a user repository for storage in local memory
type psqlUserRepository struct {
	conn *gorm.DB
}

func NewPsqlUserRepository(host, port, dbname, user, pass, ssl, timezone string) (users.Repository, error) {
	dsn := "host=" + host + " port=" + port + " dbname=" + dbname + " user=" + user + " password=" + pass + " sslmode=" + ssl + " TimeZone=" + timezone

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&User{})

	return &psqlUserRepository{conn: db}, nil
}

func (d *psqlUserRepository) Store(user *users.User) error {
	result := d.conn.Create(&user)
	return result.Error
}

// Find retrieves a single user from the repository
func (d *psqlUserRepository) Find(id int) (*users.User, error) {
	var user users.User
	result := d.conn.First(&user, id)
	return &user, result.Error
}

// FindAll retrieves all users from memory
func (d *psqlUserRepository) FindAll() []*users.User {
	users := []*users.User{}
	d.conn.Find(&users)
	return users
}
