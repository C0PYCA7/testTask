package postgres

import (
	"fmt"
	_ "github.com/lib/pq"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"
	"testTask/internal/database"
	"testTask/internal/handler/user/update"
	"testTask/internal/models"
)

type User struct {
	ID          int64  `gorm:"primaryKey;type:serial"`
	Name        string `gorm:"not null;type:varchar(50)"`
	Surname     string `gorm:"not null;type:varchar(50)"`
	Patronymic  string `gorm:"type:varchar(50)"`
	Age         int    `gorm:"not null;type:integer"`
	Gender      string `gorm:"not null;type:varchar(50)"`
	Nationality string `gorm:"not null;type:varchar(50)"`
}

func (d *Database) Migrate() error {
	const op = "database/postgres/postgres/Migrate"

	err := d.db.AutoMigrate(&User{})
	if err != nil {
		return fmt.Errorf("%s: failed to create database: %w", op, err)
	}
	return nil
}

func (d *Database) CreateUser(user *User) (int64, error) {

	const op = "database/postgres/postgres/CreateUser"

	err := d.db.Create(&user).Error
	if err != nil {
		return 0, fmt.Errorf("%s: failed to create user: %w", op, err)
	}
	return user.ID, nil
}

func (d *Database) DeleteUser(id int64) error {
	const op = "database/postgres/postgres/DeleteUser"

	user := &User{ID: id}

	results := d.db.Delete(user)
	if results.Error != nil {
		return fmt.Errorf("%s: %w", op, database.ErrInternal)
	}
	if results.RowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, database.ErrUserNotFound)
	}
	return nil
}

func (d *Database) UpdateUser(id int64, request update.Request) error {
	const op = "database/postgres/postgres/UpdateUser"
	user := &User{ID: id}

	results := d.db.Model(user).Updates(request)
	if results.RowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, database.ErrUserNotFound)
	}
	if results.Error != nil {
		return fmt.Errorf("%s: %w", op, database.ErrInternal)
	}

	return nil
}

func (d *Database) GetUsers(filter models.Filter, pageSize, page int) ([]User, error) {
	var users []User

	if page == 0 && pageSize == 0 {
		err := d.db.Where(filter).Model(&User{}).Find(&users).Error
		if err != nil {
			return nil, fmt.Errorf("failed to get users")
		}
	} else {
		err := d.db.Where(filter).Model(&User{}).Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
		if err != nil {
			return nil, fmt.Errorf("failed to get users")
		}
	}

	return users, nil
}
