package auth

import "github.com/stretchr/testify/mock"

type mockIDTool struct {
	mock.Mock
}

func (m *mockIDTool) New() (string, error) {
	ret := m.Called()

	r0 := ret.Get(0).(string)
	r1 := ret.Error(1)

	return r0, r1
}
