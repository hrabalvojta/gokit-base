package psql

import (
	"github.com/hrabalvojta/micro-dvdrental/films"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// psqlUserRepository is an implementation of GORM PSQL
type psqlUserRepository struct {
	conn *gorm.DB
}

func NewPsqlUserRepository(host, port, dbname, user, pass, ssl, timezone string) (films.Repository, error) {
	// Create connection
	dsn := "host=" + host + " port=" + port + " dbname=" + dbname + " user=" + user + " password=" + pass + " sslmode=" + ssl + " TimeZone=" + timezone
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Create enums
	result := db.Exec(Enum_mpaa_rating_select)
	switch {
	case result.RowsAffected == 0:
		if err := db.Exec(Enum_mpaa_rating_sql).Error; err != nil {
			return nil, err
		}
		return nil, err
	case result.Error != nil:
		return nil, result.Error
	}

	// Automigrate
	err = db.AutoMigrate(&User{}, &Category{}, &Film_category{}, &Film{}, &Language{})
	if err != nil {
		return nil, err
	}

	return &psqlUserRepository{conn: db}, nil
}

func (d *psqlUserRepository) Store(user *films.User) error {
	result := d.conn.Create(&user)
	return result.Error
}

// Find retrieves a single user from the repository
func (d *psqlUserRepository) Find(id int) (*films.User, error) {
	var user films.User
	result := d.conn.First(&user, id)
	return &user, result.Error
}

// FindAll retrieves all users from memory
func (d *psqlUserRepository) FindAll() []*films.User {
	users := []*films.User{}
	d.conn.Find(&users)
	return users
}
