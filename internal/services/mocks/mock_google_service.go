package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/felipe1496/open-wallet/internal/services"
)

type MockGoogleService struct {
	mock.Mock
}

func (m *MockGoogleService) GetUserInfo(accessToken string) (*services.GoogleUserInfo, error) {
	args := m.Called(accessToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.GoogleUserInfo), args.Error(1)
}

func (m *MockGoogleService) GetUserAccessToken(code string) (*string, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}
