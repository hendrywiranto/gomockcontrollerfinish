package examples

import (
	"testing"

	gomick "github.com/golang/mock/gomock"
)

func TestRenamedFinishCall(t *testing.T) {
	mock := gomick.NewController(t)
	mock.Finish() // want "since go1.14, if you are passing a testing.T to NewController then calling Finish on gomock.Controller is no longer needed"
}

func TestRenamedFinishCallDefer(t *testing.T) {
	mock := gomick.NewController(t)
	defer mock.Finish() // want "since go1.14, if you are passing a testing.T to NewController then calling Finish on gomock.Controller is no longer needed"
}

func TestRenamedNoFinishCall(t *testing.T) {
	gomick.NewController(t)
}
