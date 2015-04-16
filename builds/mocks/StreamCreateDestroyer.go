package mocks

import "github.com/stretchr/testify/mock"

import "github.com/bmorton/builder/streams"

type StreamCreateDestroyer struct {
	mock.Mock
}

func (m *StreamCreateDestroyer) Create(_a0 *streams.BuildStream) {
	m.Called(_a0)
}
func (m *StreamCreateDestroyer) Destroy(_a0 string) {
	m.Called(_a0)
}
