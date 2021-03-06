// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"context"
	"sync"

	"github.com/seibert-media/go-kafka/consumer"
)

type Consumer struct {
	ConsumeStub        func(context.Context) error
	consumeMutex       sync.RWMutex
	consumeArgsForCall []struct {
		arg1 context.Context
	}
	consumeReturns struct {
		result1 error
	}
	consumeReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Consumer) Consume(arg1 context.Context) error {
	fake.consumeMutex.Lock()
	ret, specificReturn := fake.consumeReturnsOnCall[len(fake.consumeArgsForCall)]
	fake.consumeArgsForCall = append(fake.consumeArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	fake.recordInvocation("Consume", []interface{}{arg1})
	fake.consumeMutex.Unlock()
	if fake.ConsumeStub != nil {
		return fake.ConsumeStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.consumeReturns
	return fakeReturns.result1
}

func (fake *Consumer) ConsumeCallCount() int {
	fake.consumeMutex.RLock()
	defer fake.consumeMutex.RUnlock()
	return len(fake.consumeArgsForCall)
}

func (fake *Consumer) ConsumeCalls(stub func(context.Context) error) {
	fake.consumeMutex.Lock()
	defer fake.consumeMutex.Unlock()
	fake.ConsumeStub = stub
}

func (fake *Consumer) ConsumeArgsForCall(i int) context.Context {
	fake.consumeMutex.RLock()
	defer fake.consumeMutex.RUnlock()
	argsForCall := fake.consumeArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Consumer) ConsumeReturns(result1 error) {
	fake.consumeMutex.Lock()
	defer fake.consumeMutex.Unlock()
	fake.ConsumeStub = nil
	fake.consumeReturns = struct {
		result1 error
	}{result1}
}

func (fake *Consumer) ConsumeReturnsOnCall(i int, result1 error) {
	fake.consumeMutex.Lock()
	defer fake.consumeMutex.Unlock()
	fake.ConsumeStub = nil
	if fake.consumeReturnsOnCall == nil {
		fake.consumeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.consumeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Consumer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.consumeMutex.RLock()
	defer fake.consumeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Consumer) recordInvocation(key string, args []interface{}) {
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

var _ consumer.Consumer = new(Consumer)
