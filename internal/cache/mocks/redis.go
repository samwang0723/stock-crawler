// Code generated by MockGen. DO NOT EDIT.
// Source: redis.go

// Package cache is a generated GoMock package.
package cache

import (
	context "context"
	reflect "reflect"
	time "time"

	redislock "github.com/bsm/redislock"
	gomock "github.com/golang/mock/gomock"
)

// MockRedis is a mock of Redis interface.
type MockRedis struct {
	ctrl     *gomock.Controller
	recorder *MockRedisMockRecorder
}

// MockRedisMockRecorder is the mock recorder for MockRedis.
type MockRedisMockRecorder struct {
	mock *MockRedis
}

// NewMockRedis creates a new mock instance.
func NewMockRedis(ctrl *gomock.Controller) *MockRedis {
	mock := &MockRedis{ctrl: ctrl}
	mock.recorder = &MockRedisMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedis) EXPECT() *MockRedisMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockRedis) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockRedisMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRedis)(nil).Close))
}

// ObtainLock mocks base method.
func (m *MockRedis) ObtainLock(ctx context.Context, key string, expire time.Duration) *redislock.Lock {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ObtainLock", ctx, key, expire)
	ret0, _ := ret[0].(*redislock.Lock)
	return ret0
}

// ObtainLock indicates an expected call of ObtainLock.
func (mr *MockRedisMockRecorder) ObtainLock(ctx, key, expire any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ObtainLock", reflect.TypeOf((*MockRedis)(nil).ObtainLock), ctx, key, expire)
}

// SAdd mocks base method.
func (m *MockRedis) SAdd(ctx context.Context, key, value string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SAdd", ctx, key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// SAdd indicates an expected call of SAdd.
func (mr *MockRedisMockRecorder) SAdd(ctx, key, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SAdd", reflect.TypeOf((*MockRedis)(nil).SAdd), ctx, key, value)
}

// SMembers mocks base method.
func (m *MockRedis) SMembers(ctx context.Context, key string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SMembers", ctx, key)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SMembers indicates an expected call of SMembers.
func (mr *MockRedisMockRecorder) SMembers(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SMembers", reflect.TypeOf((*MockRedis)(nil).SMembers), ctx, key)
}

// SetExpire mocks base method.
func (m *MockRedis) SetExpire(ctx context.Context, key string, expired time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetExpire", ctx, key, expired)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetExpire indicates an expected call of SetExpire.
func (mr *MockRedisMockRecorder) SetExpire(ctx, key, expired any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetExpire", reflect.TypeOf((*MockRedis)(nil).SetExpire), ctx, key, expired)
}
