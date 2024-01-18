package postgres

import (
	"fmt"
	_ "github.com/lib/pq"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"
	"testTask/internal/handler/user/update"
	"testTask/internal/models"
)

// todo: 255 -> ?
type User struct {
	ID          int64  `gorm:"primaryKey;type:serial"`
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

func (d *Database) CreateUser(name, surname, patronymic, nationalize, gender string, age int) (int64, error) {

	const op = "database/postgres/postgres/CreateUser"

	user := &User{
		Name:        name,
		Surname:     surname,
		Patronymic:  patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationalize,
	}

	err := d.db.Create(user).Error
	if err != nil {
		return 0, fmt.Errorf("%s: failed to create user: %w", op, err)
	}
	return user.ID, nil
}

func (d *Database) DeleteUser(id int64) error {
	const op = "database/postgres/postgres/DeleteUser"

	user := &User{ID: id}

	err := d.db.Delete(user).Error
	if err != nil {
		return fmt.Errorf("%s: failed to delete user: %w", op, err)
	}

	return nil
}

func (d *Database) UpdateUser(id int64, request update.Request) (int64, error) {
	const op = "database/postgres/postgres/UpdateUser"
	user := &User{
		ID:          id,
		Name:        request.Name,
		Surname:     request.Surname,
		Patronymic:  request.Patronymic,
		Age:         request.Age,
		Gender:      request.Gender,
		Nationality: request.Nationality,
	}

	err := d.db.Save(user).Error
	if err != nil {
		return 0, fmt.Errorf("%s: failed to update user: %w", op, err)
	}

	return 0, nil
}

func (d *Database) GetUsers(filter models.Filter, pageSize, page int) ([]User, error) {
	var users []User
	err := d.db.Model(&User{}).Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get users")
	}
	return users, nil
}
