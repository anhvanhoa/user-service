package authUC

import (
	"cms-server/domain/entity"
	"cms-server/domain/repository"
	"cms-server/domain/service/argon"
	serviceJwt "cms-server/domain/service/jwt"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct{}

func (m *mockUserRepo) CreateUser(entity.User) (entity.UserInfor, error) {
	return entity.UserInfor{ID: "1"}, nil
}
func (m *mockUserRepo) GetUserByEmailOrPhone(val string) (entity.User, error) {
	if val == "error@example.com" {
		return entity.User{}, assert.AnError
	}
	return entity.User{ID: "1", Email: "test@example.com", Password: "hashed"}, nil
}
func (m *mockUserRepo) GetUserByID(id string) (entity.User, error) { return entity.User{ID: id}, nil }
func (m *mockUserRepo) CheckUserExist(val string) (bool, error)    { return true, nil }
func (m *mockUserRepo) GetUserByEmail(email string) (entity.User, error) {
	return entity.User{ID: "1", Email: email}, nil
}
func (m *mockUserRepo) UpdateUser(Id string, data entity.User) (entity.UserInfor, error) {
	return entity.UserInfor{ID: Id}, nil
}
func (m *mockUserRepo) UpdateUserByEmail(email string, data entity.User) (bool, error) {
	return true, nil
}
func (m *mockUserRepo) Tx(ctx context.Context) repository.UserRepository { return m }

type mockArgon struct{}

func (m *mockArgon) HashPassword(password string) (string, error) { return password, nil }
func (m *mockArgon) VerifyPassword(hashedPassword, password string) (bool, error) {
	if password == "err" {
		return false, assert.AnError
	}
	return hashedPassword == password, nil
}
func (m *mockArgon) SetParams(memory uint32, iterations uint32, parallelism uint8, saltLength uint32, keyLength uint32) argon.Argon {
	return m
}
func (m *mockArgon) SetMemory(memory uint32) argon.Argon          { return m }
func (m *mockArgon) SetIterations(iterations uint32) argon.Argon  { return m }
func (m *mockArgon) SetParallelism(parallelism uint8) argon.Argon { return m }
func (m *mockArgon) SetSaltLength(saltLength uint32) argon.Argon  { return m }
func (m *mockArgon) SetKeyLength(keyLength uint32) argon.Argon    { return m }

type mockJwtService struct{}

func (m *mockJwtService) GenRegisterToken(id, code string, exp time.Time) (string, error) {
	return "token", nil
}
func (m *mockJwtService) VerifyRegisterToken(token string) (*serviceJwt.VerifyClaims, error) {
	return &serviceJwt.VerifyClaims{}, nil
}
func (m *mockJwtService) GenAuthToken(id, fullName string, exp time.Time) (string, error) {
	if id == "err" {
		return "", assert.AnError
	}
	return "token", nil
}
func (m *mockJwtService) VerifyAuthToken(token string) (*serviceJwt.AuthClaims, error) {
	return &serviceJwt.AuthClaims{}, nil
}
func (m *mockJwtService) GenForgotPasswordToken(id, code string, exp time.Time) (string, error) {
	return "token", nil
}
func (m *mockJwtService) VerifyForgotPasswordToken(token string) (*serviceJwt.ForgotPasswordClaims, error) {
	return &serviceJwt.ForgotPasswordClaims{}, nil
}

type mockSessionRepo struct {
	createSessionCalled *bool
}

func (m *mockSessionRepo) CreateSession(data entity.Session) error {
	if m.createSessionCalled != nil {
		*m.createSessionCalled = true
	}
	if data.UserID == "errSession" {
		return assert.AnError
	}
	return nil
}
func (m *mockSessionRepo) GetSessionAliveByToken(typeSession entity.SessionType, token string) (entity.Session, error) {
	return entity.Session{}, nil
}
func (m *mockSessionRepo) GetSessionAliveByTokenAndIdUser(typeSession entity.SessionType, token, idUser string) (entity.Session, error) {
	return entity.Session{}, nil
}
func (m *mockSessionRepo) GetSessionForgotAliveByTokenAndIdUser(token, idUser string) (entity.Session, error) {
	return entity.Session{}, nil
}
func (m *mockSessionRepo) TokenExists(token string) bool {
	return false
}
func (m *mockSessionRepo) DeleteSessionByTypeAndUserID(sessionType entity.SessionType, userID string) error {
	return nil
}
func (m *mockSessionRepo) DeleteSessionByTypeAndToken(sessionType entity.SessionType, token string) error {
	return nil
}
func (m *mockSessionRepo) DeleteSessionVerifyByUserID(userID string) error {
	return nil
}
func (m *mockSessionRepo) DeleteSessionAuthByToken(token string) error {
	return nil
}
func (m *mockSessionRepo) DeleteSessionVerifyByToken(token string) error {
	return nil
}
func (m *mockSessionRepo) DeleteSessionForgotByToken(token string) error {
	return nil
}
func (m *mockSessionRepo) DeleteAllSessionsExpired() error {
	return nil
}
func (m *mockSessionRepo) DeleteSessionForgotByTokenAndIdUser(token, idUser string) error {
	return nil
}
func (m *mockSessionRepo) DeleteAllSessionsForgot() error {
	return nil
}

func (m *mockSessionRepo) Tx(ctx context.Context) repository.SessionRepository {
	return m
}

type mockCache struct {
	setErr bool
}

func (m *mockCache) Get(key string) ([]byte, error) {
	return nil, nil
}
func (m *mockCache) Set(key string, val []byte, exp time.Duration) error {
	if m.setErr {
		return assert.AnError
	}
	return nil
}
func (m *mockCache) Delete(key string) error {
	return nil
}
func (m *mockCache) Reset() error {
	return nil
}
func (m *mockCache) Close() error {
	return nil
}

func TestLoginUsecase_GetUserByEmailOrPhone(t *testing.T) {
	uc := NewLoginUsecase(&mockUserRepo{}, nil, nil, nil, &mockArgon{}, nil)
	user, err := uc.GetUserByEmailOrPhone("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestLoginUsecase_CheckHashPassword(t *testing.T) {
	uc := NewLoginUsecase(&mockUserRepo{}, nil, nil, nil, &mockArgon{}, nil)
	ok := uc.CheckHashPassword("hashed", "hashed")
	assert.True(t, ok)
	ok = uc.CheckHashPassword("wrong", "hashed")
	assert.False(t, ok)
}

func TestLoginUsecase_GengerateAccessToken(t *testing.T) {
	mockJwt := &mockJwtService{}
	uc := NewLoginUsecase(&mockUserRepo{}, nil, mockJwt, nil, &mockArgon{}, nil)
	token, err := uc.GengerateAccessToken("1", "Test", time.Now().Add(time.Hour))
	assert.NoError(t, err)
	assert.Equal(t, "token", token)
}

func TestLoginUsecase_GetUserByEmailOrPhone_Error(t *testing.T) {
	uc := NewLoginUsecase(&mockUserRepo{}, nil, nil, nil, &mockArgon{}, nil)
	_, err := uc.GetUserByEmailOrPhone("error@example.com")
	assert.Error(t, err)
}

func TestLoginUsecase_CheckHashPassword_Error(t *testing.T) {
	uc := NewLoginUsecase(&mockUserRepo{}, nil, nil, nil, &mockArgon{}, nil)
	ok := uc.CheckHashPassword("hashed", "err")
	assert.False(t, ok)
}

func TestLoginUsecase_GengerateAccessToken_Error(t *testing.T) {
	mockJwt := &mockJwtService{}
	uc := NewLoginUsecase(&mockUserRepo{}, nil, mockJwt, nil, &mockArgon{}, nil)
	_, err := uc.GengerateAccessToken("err", "Test", time.Now().Add(time.Hour))
	assert.Error(t, err)
}

func TestLoginUsecase_GengerateRefreshToken_CacheOK(t *testing.T) {
	called := false
	uc := NewLoginUsecase(&mockUserRepo{}, &mockSessionRepo{createSessionCalled: &called}, nil, &mockJwtService{}, &mockArgon{}, &mockCache{})
	token, err := uc.GengerateRefreshToken("1", "Test", time.Now().Add(time.Hour), "win")
	assert.NoError(t, err)
	assert.Equal(t, "token", token)
	time.Sleep(10 * time.Millisecond) // Đợi goroutine chạy xong
	assert.True(t, called)            // go routine vẫn gọi CreateSession
}

func TestLoginUsecase_GengerateRefreshToken_CacheSetErr_SessionOK(t *testing.T) {
	uc := NewLoginUsecase(&mockUserRepo{}, &mockSessionRepo{}, nil, &mockJwtService{}, &mockArgon{}, &mockCache{setErr: true})
	token, err := uc.GengerateRefreshToken("1", "Test", time.Now().Add(time.Hour), "win")
	assert.NoError(t, err)
	assert.Equal(t, "token", token)
}

func TestLoginUsecase_GengerateRefreshToken_CacheSetErr_SessionErr(t *testing.T) {
	uc := NewLoginUsecase(&mockUserRepo{}, &mockSessionRepo{}, nil, &mockJwtService{}, &mockArgon{}, &mockCache{setErr: true})
	_, err := uc.GengerateRefreshToken("err", "Test", time.Now().Add(time.Hour), "win")
	assert.Error(t, err)
}

func TestLoginUsecase_GengerateRefreshToken_CacheSetErr_SessionCreateSessionFail(t *testing.T) {
	repo := &mockSessionRepo{}
	// ép CreateSession trả về lỗi khi UserID == "errSession"
	uc := NewLoginUsecase(&mockUserRepo{}, repo, nil, &mockJwtService{}, &mockArgon{}, &mockCache{setErr: true})
	_, err := uc.GengerateRefreshToken("errSession", "Test", time.Now().Add(time.Hour), "win")
	assert.Error(t, err)
}

func TestLoginUsecase_CheckHashPassword_Fail(t *testing.T) {
	uc := NewLoginUsecase(&mockUserRepo{}, nil, nil, nil, &mockArgon{}, nil)
	ok := uc.CheckHashPassword("hashed", "notmatch")
	assert.False(t, ok)
}

func TestLoginUsecase_CheckHashPassword_VerifyPasswordError(t *testing.T) {
	uc := NewLoginUsecase(&mockUserRepo{}, nil, nil, nil, &mockArgon{}, nil)
	ok := uc.CheckHashPassword("err", "hashed")
	assert.False(t, ok)
}

func TestLoginUsecase_CheckHashPassword_TableDriven(t *testing.T) {
	uc := NewLoginUsecase(&mockUserRepo{}, nil, nil, nil, &mockArgon{}, nil)

	tests := []struct {
		name     string
		password string
		hash     string
		expect   bool
	}{
		{"match", "hashed", "hashed", true},
		{"not match", "wrong", "hashed", false},
		{"verify error", "err", "hashed", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := uc.CheckHashPassword(tt.password, tt.hash)
			assert.Equal(t, tt.expect, ok)
		})
	}
}
