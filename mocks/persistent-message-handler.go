// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"context"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/boltdb/bolt"
	"github.com/seibert-media/go-kafka/persistent"
)

type PersistentMessageHandler struct {
	ConsumeMessageStub        func(context.Context, *bolt.Tx, *sarama.ConsumerMessage) error
	consumeMessageMutex       sync.RWMutex
	consumeMessageArgsForCall []struct {
		arg1 context.Context
		arg2 *bolt.Tx
		arg3 *sarama.ConsumerMessage
	}
	consumeMessageReturns struct {
		result1 error
	}
	consumeMessageReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *PersistentMessageHandler) ConsumeMessage(arg1 context.Context, arg2 *bolt.Tx, arg3 *sarama.ConsumerMessage) error {
	fake.consumeMessageMutex.Lock()
	ret, specificReturn := fake.consumeMessageReturnsOnCall[len(fake.consumeMessageArgsForCall)]
	fake.consumeMessageArgsForCall = append(fake.consumeMessageArgsForCall, struct {
		arg1 context.Context
		arg2 *bolt.Tx
		arg3 *sarama.ConsumerMessage
	}{arg1, arg2, arg3})
	fake.recordInvocation("ConsumeMessage", []interface{}{arg1, arg2, arg3})
	fake.consumeMessageMutex.Unlock()
	if fake.ConsumeMessageStub != nil {
		return fake.ConsumeMessageStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.consumeMessageReturns
	return fakeReturns.result1
}

func (fake *PersistentMessageHandler) ConsumeMessageCallCount() int {
	fake.consumeMessageMutex.RLock()
	defer fake.consumeMessageMutex.RUnlock()
	return len(fake.consumeMessageArgsForCall)
}

func (fake *PersistentMessageHandler) ConsumeMessageCalls(stub func(context.Context, *bolt.Tx, *sarama.ConsumerMessage) error) {
	fake.consumeMessageMutex.Lock()
	defer fake.consumeMessageMutex.Unlock()
	fake.ConsumeMessageStub = stub
}

func (fake *PersistentMessageHandler) ConsumeMessageArgsForCall(i int) (context.Context, *bolt.Tx, *sarama.ConsumerMessage) {
	fake.consumeMessageMutex.RLock()
	defer fake.consumeMessageMutex.RUnlock()
	argsForCall := fake.consumeMessageArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *PersistentMessageHandler) ConsumeMessageReturns(result1 error) {
	fake.consumeMessageMutex.Lock()
	defer fake.consumeMessageMutex.Unlock()
	fake.ConsumeMessageStub = nil
	fake.consumeMessageReturns = struct {
		result1 error
	}{result1}
}

func (fake *PersistentMessageHandler) ConsumeMessageReturnsOnCall(i int, result1 error) {
	fake.consumeMessageMutex.Lock()
	defer fake.consumeMessageMutex.Unlock()
	fake.ConsumeMessageStub = nil
	if fake.consumeMessageReturnsOnCall == nil {
		fake.consumeMessageReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.consumeMessageReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *PersistentMessageHandler) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.consumeMessageMutex.RLock()
	defer fake.consumeMessageMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *PersistentMessageHandler) recordInvocation(key string, args []interface{}) {
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

var _ persistent.MessageHandler = new(PersistentMessageHandler)
