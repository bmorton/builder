package mocks

import "github.com/stretchr/testify/mock"

import "github.com/bmorton/builder/builds"

type BuildQueue struct {
	mock.Mock
}

func (m *BuildQueue) Add(_a0 *builds.Build) string {
	ret := m.Called(_a0)

	r0 := ret.Get(0).(string)

	return r0
}
