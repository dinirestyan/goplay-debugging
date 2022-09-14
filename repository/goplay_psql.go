package repository

import (
	"github.com/dinirestyan/goplay-debugging/models"
	"github.com/jinzhu/gorm"
)

func NewPSQLGoplayRepo(conn *gorm.DB) models.GoplayRepo {
	return &psqlGoplayrRepo{
		Conn: conn,
	}
}

type psqlGoplayrRepo struct {
	Conn *gorm.DB
}

func (p *psqlGoplayrRepo) Login(email string) (models.User, error) {

	var user models.User
	err := p.Conn.Model(&user).
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return models.User{}, err
		}
		return models.User{}, err
	}

	return user, nil

}

func (p *psqlGoplayrRepo) UpdateToken(email string, token string) error {

	err := p.Conn.Where("email = ?", email).
		Update("token", token).
		Error
	if err != nil {
		return err
	}
	return nil
}
