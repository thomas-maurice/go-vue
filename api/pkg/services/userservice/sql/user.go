package sqluserservice

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/thomas-maurice/api/go-vue/pkg/services/userservice"
	"github.com/thomas-maurice/api/go-vue/pkg/services/userservice/sql/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CustomClaims struct {
	SessionId string `json:"session_id"`
	Name      string `json:"name"`
	Admin     bool   `json:"admin"`

	jwt.RegisteredClaims
}

type UserService struct {
	DB         *gorm.DB
	SigningKey *ecdsa.PrivateKey
}

func userFromModel(input *models.User) *userservice.User {
	return &userservice.User{
		Id:          input.Id,
		Username:    input.Username,
		Email:       input.Email,
		DisplayName: input.DisplayName,
		Admin:       input.Admin,
		Active:      input.Active,
		Kind:        userservice.UserKind(input.Kind),
		Created:     input.Created,
		LastLogin:   input.LastLogin,
	}
}

func sessionFromModel(input *models.Session) *userservice.Session {
	return &userservice.Session{
		Id:      input.Id,
		Expires: input.Expires,
	}
}

func userToModel(input *models.User) *models.User {
	return &models.User{
		Id:          input.Id,
		Username:    input.Username,
		Email:       input.Email,
		DisplayName: input.DisplayName,
		Admin:       input.Admin,
		Active:      input.Active,
		Kind:        userservice.UserKind(input.Kind),
		Created:     input.Created,
		LastLogin:   input.LastLogin,
	}
}

func sessionToModel(input *userservice.Session) *models.Session {
	return &models.Session{
		Id:      input.Id,
		Expires: input.Expires,
	}
}

func NewUserService(db *gorm.DB, sk *ecdsa.PrivateKey) (*UserService, error) {
	if err := db.AutoMigrate(models.User{}); err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(models.Session{}); err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(models.APIKey{}); err != nil {
		return nil, err
	}

	return &UserService{
		DB:         db,
		SigningKey: sk,
	}, nil
}

func (s *UserService) GetUserByUsername(username string) (*userservice.User, error) {
	var user models.User
	if err := s.DB.Where(&models.User{Username: username}).First(&user).Error; err == gorm.ErrRecordNotFound {
		return nil, userservice.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return userFromModel(&user), nil
}

func (s *UserService) GetUserById(id string) (*userservice.User, error) {
	var user models.User
	if err := s.DB.Where(&models.User{Id: id}).First(&user).Error; err != nil {
		return nil, err
	}

	return userFromModel(&user), nil
}

func (s *UserService) UpdateUser(id string, email string, admin bool, displayName string) error {
	updates := make(map[string]any)

	updates["admin"] = admin

	if email != "" {
		updates["email"] = email
	}

	if displayName != "" {
		updates["display_name"] = displayName
	}

	return s.DB.Model(&models.User{}).Where(&models.User{Id: id}).Updates(updates).Error
}

func (s *UserService) CreateUser(username string, email string, password string, kind string, admin bool, displayName string) (*userservice.User, error) {
	hashed := ""
	if password != "" {
		b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hashed = string(b)
	}

	user := models.User{
		Username:    username,
		Password:    hashed,
		Email:       email,
		Admin:       admin,
		Kind:        userservice.UserKind(kind),
		DisplayName: displayName,
		Created:     time.Now(),
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return userFromModel(&user), nil
}

func (s *UserService) ListUsers() ([]userservice.User, error) {
	var users []models.User
	if err := s.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	ulist := make([]userservice.User, 0)
	for _, u := range users {
		ulist = append(ulist, *(userFromModel(&u)))
	}

	return ulist, nil
}

func (s *UserService) Authenticate(username, password string) (*userservice.User, error) {
	var user models.User
	if err := s.DB.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		return nil, err
	}

	if user.Kind != "local" {
		return nil, fmt.Errorf("cannot authenticate with a password on a non-local user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return userFromModel(&user), nil
}

func (s *UserService) LogoutFromToken(token string) error {
	var claims CustomClaims
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return &s.SigningKey.PublicKey, nil
	})
	if err != nil {
		return fmt.Errorf("failed to verify auth token: %w", err)
	}

	if claims.SessionId == "" {
		return fmt.Errorf("no session id provided")
	}
	if err := s.DB.Delete(&models.Session{Id: claims.SessionId}).Error; err != nil {
		return err
	}

	return nil
}

func (s *UserService) GenerateSessionToken(user *userservice.User) (string, error) {
	sessionId := uuid.NewString()

	c := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "webapp",
			Audience:  jwt.ClaimStrings{"webapp"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Subject:   user.Username,
			ID:        sessionId,
		},
		SessionId: sessionId,
		Admin:     user.Admin,
	}

	if err := s.DB.Where("expires < ?", time.Now()).Delete(&models.Session{}).Error; err != nil {
		fmt.Println("failed to cleanup old sessions")
	}

	t := jwt.New(jwt.SigningMethodES512)
	t.Claims = &c

	sig, err := t.SignedString(s.SigningKey)
	if err != nil {
		return "", err
	}

	if err := s.DB.Create(&models.Session{
		Id:      sessionId,
		UserId:  user.Id,
		Expires: time.Now().Add(time.Hour),
	}).Error; err != nil {
		return "", err
	}

	if err := s.DB.Model(&user).Update("last_login", time.Now()).Error; err != nil {
		return "", err
	}

	return sig, nil
}

func (s *UserService) VerifySessionToken(token string) (*userservice.Session, *userservice.User, error) {
	var claims CustomClaims
	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return &s.SigningKey.PublicKey, nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to verify auth token: %w", err)
	}

	if !parsedToken.Valid {
		return nil, nil, fmt.Errorf("failed to verify token claims: %w", err)
	}

	var session models.Session
	if err := s.DB.Where(&models.Session{Id: claims.SessionId}).First(&session).Error; err != nil {
		return nil, nil, err
	}

	var user models.User
	if err := s.DB.Where(&models.User{Username: claims.RegisteredClaims.Subject}).First(&user).Error; err != nil {
		return nil, nil, err
	}

	return sessionFromModel(&session), userFromModel(&user), nil
}
