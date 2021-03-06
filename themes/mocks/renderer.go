package mocks_test

import (
	"io"

	"github.com/srgrn/postmanerator/postman"
	. "github.com/srgrn/postmanerator/themes"
	"github.com/stretchr/testify/mock"
)

type MockThemeRenderer struct {
	mock.Mock
}

func (m *MockThemeRenderer) Render(w io.Writer, theme *Theme, collection postman.Collection) error {
	return m.Called(w, theme, collection).Error(0)
}
