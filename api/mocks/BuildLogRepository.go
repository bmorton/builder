package mocks

import "github.com/stretchr/testify/mock"

import "github.com/bmorton/builder/builds"

type BuildLogRepository struct {
	mock.Mock
}

func (m *BuildLogRepository) FindByBuildID(_a0 string, _a1 string) (*builds.BuildLog, error) {
	ret := m.Called(_a0, _a1)

	var r0 *builds.BuildLog
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*builds.BuildLog)
	}
	r1 := ret.Error(1)

	return r0, r1
}
