// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/trustbloc/orb/pkg/activitypub/service/spi"
	"github.com/trustbloc/orb/pkg/activitypub/vocab"
)

type ActivityHandler struct {
	StartStub        func()
	startMutex       sync.RWMutex
	startArgsForCall []struct{}
	StopStub         func()
	stopMutex        sync.RWMutex
	stopArgsForCall  []struct{}
	StateStub        func() spi.State
	stateMutex       sync.RWMutex
	stateArgsForCall []struct{}
	stateReturns     struct {
		result1 spi.State
	}
	stateReturnsOnCall map[int]struct {
		result1 spi.State
	}
	HandleActivityStub        func(activity *vocab.ActivityType) error
	handleActivityMutex       sync.RWMutex
	handleActivityArgsForCall []struct {
		activity *vocab.ActivityType
	}
	handleActivityReturns struct {
		result1 error
	}
	handleActivityReturnsOnCall map[int]struct {
		result1 error
	}
	SubscribeStub        func() <-chan *vocab.ActivityType
	subscribeMutex       sync.RWMutex
	subscribeArgsForCall []struct{}
	subscribeReturns     struct {
		result1 <-chan *vocab.ActivityType
	}
	subscribeReturnsOnCall map[int]struct {
		result1 <-chan *vocab.ActivityType
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ActivityHandler) Start() {
	fake.startMutex.Lock()
	fake.startArgsForCall = append(fake.startArgsForCall, struct{}{})
	fake.recordInvocation("Start", []interface{}{})
	fake.startMutex.Unlock()
	if fake.StartStub != nil {
		fake.StartStub()
	}
}

func (fake *ActivityHandler) StartCallCount() int {
	fake.startMutex.RLock()
	defer fake.startMutex.RUnlock()
	return len(fake.startArgsForCall)
}

func (fake *ActivityHandler) Stop() {
	fake.stopMutex.Lock()
	fake.stopArgsForCall = append(fake.stopArgsForCall, struct{}{})
	fake.recordInvocation("Stop", []interface{}{})
	fake.stopMutex.Unlock()
	if fake.StopStub != nil {
		fake.StopStub()
	}
}

func (fake *ActivityHandler) StopCallCount() int {
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	return len(fake.stopArgsForCall)
}

func (fake *ActivityHandler) State() spi.State {
	fake.stateMutex.Lock()
	ret, specificReturn := fake.stateReturnsOnCall[len(fake.stateArgsForCall)]
	fake.stateArgsForCall = append(fake.stateArgsForCall, struct{}{})
	fake.recordInvocation("State", []interface{}{})
	fake.stateMutex.Unlock()
	if fake.StateStub != nil {
		return fake.StateStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.stateReturns.result1
}

func (fake *ActivityHandler) StateCallCount() int {
	fake.stateMutex.RLock()
	defer fake.stateMutex.RUnlock()
	return len(fake.stateArgsForCall)
}

func (fake *ActivityHandler) StateReturns(result1 spi.State) {
	fake.StateStub = nil
	fake.stateReturns = struct {
		result1 spi.State
	}{result1}
}

func (fake *ActivityHandler) StateReturnsOnCall(i int, result1 spi.State) {
	fake.StateStub = nil
	if fake.stateReturnsOnCall == nil {
		fake.stateReturnsOnCall = make(map[int]struct {
			result1 spi.State
		})
	}
	fake.stateReturnsOnCall[i] = struct {
		result1 spi.State
	}{result1}
}

func (fake *ActivityHandler) HandleActivity(activity *vocab.ActivityType) error {
	fake.handleActivityMutex.Lock()
	ret, specificReturn := fake.handleActivityReturnsOnCall[len(fake.handleActivityArgsForCall)]
	fake.handleActivityArgsForCall = append(fake.handleActivityArgsForCall, struct {
		activity *vocab.ActivityType
	}{activity})
	fake.recordInvocation("HandleActivity", []interface{}{activity})
	fake.handleActivityMutex.Unlock()
	if fake.HandleActivityStub != nil {
		return fake.HandleActivityStub(activity)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.handleActivityReturns.result1
}

func (fake *ActivityHandler) HandleActivityCallCount() int {
	fake.handleActivityMutex.RLock()
	defer fake.handleActivityMutex.RUnlock()
	return len(fake.handleActivityArgsForCall)
}

func (fake *ActivityHandler) HandleActivityArgsForCall(i int) *vocab.ActivityType {
	fake.handleActivityMutex.RLock()
	defer fake.handleActivityMutex.RUnlock()
	return fake.handleActivityArgsForCall[i].activity
}

func (fake *ActivityHandler) HandleActivityReturns(result1 error) {
	fake.HandleActivityStub = nil
	fake.handleActivityReturns = struct {
		result1 error
	}{result1}
}

func (fake *ActivityHandler) HandleActivityReturnsOnCall(i int, result1 error) {
	fake.HandleActivityStub = nil
	if fake.handleActivityReturnsOnCall == nil {
		fake.handleActivityReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.handleActivityReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ActivityHandler) Subscribe() <-chan *vocab.ActivityType {
	fake.subscribeMutex.Lock()
	ret, specificReturn := fake.subscribeReturnsOnCall[len(fake.subscribeArgsForCall)]
	fake.subscribeArgsForCall = append(fake.subscribeArgsForCall, struct{}{})
	fake.recordInvocation("Subscribe", []interface{}{})
	fake.subscribeMutex.Unlock()
	if fake.SubscribeStub != nil {
		return fake.SubscribeStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.subscribeReturns.result1
}

func (fake *ActivityHandler) SubscribeCallCount() int {
	fake.subscribeMutex.RLock()
	defer fake.subscribeMutex.RUnlock()
	return len(fake.subscribeArgsForCall)
}

func (fake *ActivityHandler) SubscribeReturns(result1 <-chan *vocab.ActivityType) {
	fake.SubscribeStub = nil
	fake.subscribeReturns = struct {
		result1 <-chan *vocab.ActivityType
	}{result1}
}

func (fake *ActivityHandler) SubscribeReturnsOnCall(i int, result1 <-chan *vocab.ActivityType) {
	fake.SubscribeStub = nil
	if fake.subscribeReturnsOnCall == nil {
		fake.subscribeReturnsOnCall = make(map[int]struct {
			result1 <-chan *vocab.ActivityType
		})
	}
	fake.subscribeReturnsOnCall[i] = struct {
		result1 <-chan *vocab.ActivityType
	}{result1}
}

func (fake *ActivityHandler) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.startMutex.RLock()
	defer fake.startMutex.RUnlock()
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	fake.stateMutex.RLock()
	defer fake.stateMutex.RUnlock()
	fake.handleActivityMutex.RLock()
	defer fake.handleActivityMutex.RUnlock()
	fake.subscribeMutex.RLock()
	defer fake.subscribeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ActivityHandler) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ spi.ActivityHandler = new(ActivityHandler)
