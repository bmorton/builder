package mocks

import "github.com/stretchr/testify/mock"

import "github.com/bmorton/builder/streams"

type StreamRepository struct {
	mock.Mock
}

func (m *StreamRepository) Find(_a0 string) (*streams.BuildStream, error) {
	ret := m.Called(_a0)

	var r0 *streams.BuildStream
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*streams.BuildStream)
	}
	r1 := ret.Error(1)

	return r0, r1
}
