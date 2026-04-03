package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUsernameTaken      = errors.New("username already exists")
)

const defaultTokenTTL = 24 * time.Hour

type CarryOverService interface {
	CarryOverPendingTasks(userID uint64, today time.Time) (int64, error)
}

type Service struct {
	repo      *Repository
	secretKey []byte
	tokenTTL  time.Duration
	now       func() time.Time
	taskSvc   CarryOverService
}

func NewService(repo *Repository, secret string) *Service {
	if strings.TrimSpace(secret) == "" {
		secret = os.Getenv("AUTH_TOKEN_SECRET")
	}
	if strings.TrimSpace(secret) == "" {
		secret = "dev-only-change-me"
	}

	return &Service{
		repo:      repo,
		secretKey: []byte(secret),
		tokenTTL:  defaultTokenTTL,
		now:       time.Now,
	}
}

func (s *Service) SetTaskService(taskSvc CarryOverService) {
	s.taskSvc = taskSvc
}

func (s *Service) Register(username, password string) (string, uint64, string, error) {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return "", 0, "", ErrInvalidCredentials
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", 0, "", err
	}

	user, err := s.repo.CreateUser(context.Background(), strings.TrimSpace(username), string(hash))
	if err != nil {
		if isUniqueConstraintError(err) {
			return "", 0, "", ErrUsernameTaken
		}
		return "", 0, "", err
	}

	token, err := s.IssueToken(user.ID)
	if err != nil {
		return "", 0, "", err
	}

	return token, user.ID, user.Username, nil
}

func (s *Service) Login(username, password string) (string, uint64, string, error) {
	user, err := s.repo.GetUserByUsername(context.Background(), strings.TrimSpace(username))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", 0, "", ErrInvalidCredentials
		}
		return "", 0, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", 0, "", ErrInvalidCredentials
	}

	token, err := s.IssueToken(user.ID)
	if err != nil {
		return "", 0, "", err
	}

	if s.taskSvc != nil {
		_, _ = s.taskSvc.CarryOverPendingTasks(user.ID, s.now())
	}

	return token, user.ID, user.Username, nil
}

func (s *Service) ChangePassword(userID uint64, oldPassword, newPassword string) error {
	user, err := s.repo.GetUserByID(context.Background(), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidCredentials
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePasswordHash(context.Background(), userID, string(hash))
}

func (s *Service) GetUser(userID uint64) (uint64, string, error) {
	user, err := s.repo.GetUserByID(context.Background(), userID)
	if err != nil {
		return 0, "", err
	}
	return user.ID, user.Username, nil
}

func (s *Service) IssueToken(userID uint64) (string, error) {
	expiresAt := s.now().Add(s.tokenTTL).Unix()
	payload := fmt.Sprintf("%d:%d", userID, expiresAt)
	sig := signHMAC(s.secretKey, payload)

	payloadPart := base64.RawURLEncoding.EncodeToString([]byte(payload))
	sigPart := base64.RawURLEncoding.EncodeToString([]byte(sig))
	return payloadPart + "." + sigPart, nil
}

func (s *Service) VerifyToken(token string) (uint64, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return 0, ErrInvalidToken
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, ErrInvalidToken
	}
	sigBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, ErrInvalidToken
	}

	payload := string(payloadBytes)
	expectedSig := signHMAC(s.secretKey, payload)
	if !hmac.Equal(sigBytes, []byte(expectedSig)) {
		return 0, ErrInvalidToken
	}

	payloadParts := strings.Split(payload, ":")
	if len(payloadParts) != 2 {
		return 0, ErrInvalidToken
	}

	userID, err := strconv.ParseUint(payloadParts[0], 10, 64)
	if err != nil {
		return 0, ErrInvalidToken
	}
	expiresAt, err := strconv.ParseInt(payloadParts[1], 10, 64)
	if err != nil {
		return 0, ErrInvalidToken
	}
	if s.now().Unix() > expiresAt {
		return 0, ErrInvalidToken
	}

	return userID, nil
}

func signHMAC(secret []byte, payload string) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

func isUniqueConstraintError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique") || strings.Contains(msg, "duplicate")
}
