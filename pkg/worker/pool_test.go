package worker

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type WorkerPoolTestSuite struct {
	suite.Suite
}

func (suite *WorkerPoolTestSuite) TestWorkerPool_NewWorkerPool_Success() {
	require := suite.Require()
	expectedMinWorkers := 1
	expectedMinQueueSize := 0
	testCases := []struct {
		workerCount int
		queueSize   int
	}{
		{
			workerCount: 1,
			queueSize:   10,
		},
		{
			workerCount: -1,
			queueSize:   1,
		},
		{
			workerCount: 0,
			queueSize:   -3,
		},
	}

	for _, tc := range testCases {
		wp := NewWorkerPool(logrus.New(), tc.workerCount, tc.queueSize)

		require.NotNil(wp)
		require.LessOrEqual(expectedMinWorkers, wp.GetWorkerCount())
		require.LessOrEqual(expectedMinQueueSize, wp.GetQueueSize())
	}
}

func (suite *WorkerPoolTestSuite) TestWorkerPool_Submit_Success() {
	require := suite.Require()
	wp := NewWorkerPool(logrus.New(), 5, 10)

	done := make(chan bool)
	err := wp.Submit(func() {
		done <- true
	})

	select {
	case <-time.After(200 * time.Millisecond):
		require.Fail("Submit timed out")
	case <-done:
		wp.StopAndWait()
	}
	require.NoError(err)
}

func (suite *WorkerPoolTestSuite) TestWorkerPool_Submit_StopAndWait_Success() {
	require := suite.Require()
	wp := NewWorkerPool(logrus.New(), 5, 10)

	taskExecuted := false
	err := wp.Submit(func() {
		time.Sleep(200 * time.Millisecond)
		taskExecuted = true
	})
	time.Sleep(100 * time.Millisecond)
	wp.StopAndWait()

	require.NoError(err)
	require.True(taskExecuted)
}

func TestWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(WorkerPoolTestSuite))
}
