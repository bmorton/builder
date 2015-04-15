package mocks

import "github.com/stretchr/testify/mock"

import "github.com/bmorton/builder/builds"

type BuildRepository struct {
	mock.Mock
}

func (m *BuildRepository) All() []*builds.Build {
	ret := m.Called()

	var r0 []*builds.Build
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]*builds.Build)
	}

	return r0
}
func (m *BuildRepository) Create(_a0 *builds.Build) {
	m.Called(_a0)
}
func (m *BuildRepository) Save(_a0 *builds.Build) {
	m.Called(_a0)
}
func (m *BuildRepository) Find(_a0 string) (*builds.Build, error) {
	ret := m.Called(_a0)

	var r0 *builds.Build
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*builds.Build)
	}
	r1 := ret.Error(1)

	return r0, r1
}
