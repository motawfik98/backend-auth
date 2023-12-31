package database

import (
	"backend-auth/pkg/models"
	"fmt"
	"strings"
)

func (db *DB) CreateUser(user *models.User) error {
	err := db.connection.Create(user).Error
	if err != nil && strings.HasPrefix(err.Error(), "Error 1062") {
		return &DuplicateEmailError{}
	}
	return err
}

func (db *DB) GetUserByEmail(email string) (*models.User, error) {
	user := new(models.User)
	err := db.connection.Where("email = ?", email).Find(user).Error
	return user, err
}

func (db *DB) SaveAccessRefreshTokens(userTokens *models.UserTokens) error {
	return db.connection.Save(userTokens).Error
}

func (db *DB) MarkRefreshTokenAsUsed(token *models.UsedRefreshToken) error {
	return db.connection.Create(token).Error
}

func (db *DB) GetUsersCount() int64 {
	var count int64 = 0
	db.connection.Model(&models.User{}).Count(&count)
	return count
}

func (db *DB) IsUsedRefreshToken(refreshToken string) (bool, error) {
	var count int64 = 0
	err := db.connection.Model(&models.UsedRefreshToken{}).Where("refresh_token = ?", refreshToken).Count(&count).Error
	return count > 0, err
}

func (db *DB) GetGeneratedRefreshToken(token string) *models.GeneratedRefreshToken {
	generatedToken := new(models.GeneratedRefreshToken)
	db.connection.Where("refresh_token = ?", token).First(generatedToken)
	return generatedToken
}

func (db *DB) GetGeneratedRefreshTokenChain(refreshToken string) []models.GeneratedRefreshToken {
	var generatedTokens []models.GeneratedRefreshToken
	query := fmt.Sprintf(
		`
			WITH RECURSIVE TokenHierarchy(refresh_token, id, parent_refresh_token_id) AS
			(
            	SELECT refresh_token, id, parent_refresh_token_id FROM generated_refresh_tokens WHERE refresh_token = "%s"
                UNION ALL
                SELECT grt.refresh_token, grt.id, grt.parent_refresh_token_id FROM generated_refresh_tokens as grt
					JOIN TokenHierarchy AS th ON grt.parent_refresh_token_id = th.id
		    )
			SELECT refresh_token FROM TokenHierarchy`, refreshToken)
	db.connection.Raw(query).Scan(&generatedTokens)
	return generatedTokens
}
