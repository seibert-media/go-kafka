// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/Shopify/sarama"
	"github.com/seibert-media/go-kafka/schema"
)

type SchemaAvroEncoderBuilder struct {
	BuildEncoderStub        func(string, schema.Avro) (sarama.Encoder, error)
	buildEncoderMutex       sync.RWMutex
	buildEncoderArgsForCall []struct {
		arg1 string
		arg2 schema.Avro
	}
	buildEncoderReturns struct {
		result1 sarama.Encoder
		result2 error
	}
	buildEncoderReturnsOnCall map[int]struct {
		result1 sarama.Encoder
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *SchemaAvroEncoderBuilder) BuildEncoder(arg1 string, arg2 schema.Avro) (sarama.Encoder, error) {
	fake.buildEncoderMutex.Lock()
	ret, specificReturn := fake.buildEncoderReturnsOnCall[len(fake.buildEncoderArgsForCall)]
	fake.buildEncoderArgsForCall = append(fake.buildEncoderArgsForCall, struct {
		arg1 string
		arg2 schema.Avro
	}{arg1, arg2})
	fake.recordInvocation("BuildEncoder", []interface{}{arg1, arg2})
	fake.buildEncoderMutex.Unlock()
	if fake.BuildEncoderStub != nil {
		return fake.BuildEncoderStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.buildEncoderReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *SchemaAvroEncoderBuilder) BuildEncoderCallCount() int {
	fake.buildEncoderMutex.RLock()
	defer fake.buildEncoderMutex.RUnlock()
	return len(fake.buildEncoderArgsForCall)
}

func (fake *SchemaAvroEncoderBuilder) BuildEncoderCalls(stub func(string, schema.Avro) (sarama.Encoder, error)) {
	fake.buildEncoderMutex.Lock()
	defer fake.buildEncoderMutex.Unlock()
	fake.BuildEncoderStub = stub
}

func (fake *SchemaAvroEncoderBuilder) BuildEncoderArgsForCall(i int) (string, schema.Avro) {
	fake.buildEncoderMutex.RLock()
	defer fake.buildEncoderMutex.RUnlock()
	argsForCall := fake.buildEncoderArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *SchemaAvroEncoderBuilder) BuildEncoderReturns(result1 sarama.Encoder, result2 error) {
	fake.buildEncoderMutex.Lock()
	defer fake.buildEncoderMutex.Unlock()
	fake.BuildEncoderStub = nil
	fake.buildEncoderReturns = struct {
		result1 sarama.Encoder
		result2 error
	}{result1, result2}
}

func (fake *SchemaAvroEncoderBuilder) BuildEncoderReturnsOnCall(i int, result1 sarama.Encoder, result2 error) {
	fake.buildEncoderMutex.Lock()
	defer fake.buildEncoderMutex.Unlock()
	fake.BuildEncoderStub = nil
	if fake.buildEncoderReturnsOnCall == nil {
		fake.buildEncoderReturnsOnCall = make(map[int]struct {
			result1 sarama.Encoder
			result2 error
		})
	}
	fake.buildEncoderReturnsOnCall[i] = struct {
		result1 sarama.Encoder
		result2 error
	}{result1, result2}
}

func (fake *SchemaAvroEncoderBuilder) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.buildEncoderMutex.RLock()
	defer fake.buildEncoderMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *SchemaAvroEncoderBuilder) recordInvocation(key string, args []interface{}) {
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

var _ schema.AvroEncoderBuilder = new(SchemaAvroEncoderBuilder)
