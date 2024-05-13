package mock

import (
	"github.com/stretchr/testify/mock"
)

type Pool struct {
	mock.Mock
}

func (m *Pool) Submit(task func()) error {
	return m.Called(task).Error(0)
}

func (m *Pool) StopAndWait() {
	m.Called()
}

func (m *Pool) GetWorkerCount() int {
	args := m.Called()
	return args.Get(0).(int)
}

func (m *Pool) GetQueueSize() int {
	args := m.Called()
	return args.Get(0).(int)
}
