// Code generated by MockGen. DO NOT EDIT.
// Source: producer.go

// Package kafka is a generated GoMock package.
package kafka

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockKafka is a mock of Kafka interface.
type MockKafka struct {
	ctrl     *gomock.Controller
	recorder *MockKafkaMockRecorder
}

// MockKafkaMockRecorder is the mock recorder for MockKafka.
type MockKafkaMockRecorder struct {
	mock *MockKafka
}

// NewMockKafka creates a new mock instance.
func NewMockKafka(ctrl *gomock.Controller) *MockKafka {
	mock := &MockKafka{ctrl: ctrl}
	mock.recorder = &MockKafkaMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKafka) EXPECT() *MockKafkaMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockKafka) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockKafkaMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockKafka)(nil).Close))
}

// WriteMessages mocks base method.
func (m *MockKafka) WriteMessages(ctx context.Context, topic string, message []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteMessages", ctx, topic, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteMessages indicates an expected call of WriteMessages.
func (mr *MockKafkaMockRecorder) WriteMessages(ctx, topic, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteMessages", reflect.TypeOf((*MockKafka)(nil).WriteMessages), ctx, topic, message)
}