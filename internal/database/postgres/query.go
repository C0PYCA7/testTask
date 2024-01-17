package postgres

import (
	"fmt"
	_ "github.com/lib/pq"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"
)

type User struct {
	ID          uint   `gorm:"primaryKey;type:serial"`
	Name        string `gorm:"not null;type:varchar(255)"`
	Surname     string `gorm:"not null;type:varchar(255)"`
	Patronymic  string `gorm:"type:varchar(255)"`
	Age         int    `gorm:"not null;type:integer"`
	Gender      string `gorm:"not null;type:varchar(255)"`
	Nationality string `gorm:"not null;type:varchar(255)"`
}

func (d *Database) Migrate() error {
	const op = "database/postgres/postgres/Migrate"

	err := d.db.AutoMigrate(&User{})
	if err != nil {
		return fmt.Errorf("%s: failed to create database: %w", op, err)
	}
	return nil
}

//func (d *Database) CreateUser(name, surname, patronymic, nationalize string, age int, gender bool) (int64, error) {
//	query := `INSERT INTO `
//}
