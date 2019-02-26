// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/seibert-media/go-kafka/schema"
)

type SchemaRegistry struct {
	SchemaIdStub        func(string, string) (uint32, error)
	schemaIdMutex       sync.RWMutex
	schemaIdArgsForCall []struct {
		arg1 string
		arg2 string
	}
	schemaIdReturns struct {
		result1 uint32
		result2 error
	}
	schemaIdReturnsOnCall map[int]struct {
		result1 uint32
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *SchemaRegistry) SchemaId(arg1 string, arg2 string) (uint32, error) {
	fake.schemaIdMutex.Lock()
	ret, specificReturn := fake.schemaIdReturnsOnCall[len(fake.schemaIdArgsForCall)]
	fake.schemaIdArgsForCall = append(fake.schemaIdArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("SchemaId", []interface{}{arg1, arg2})
	fake.schemaIdMutex.Unlock()
	if fake.SchemaIdStub != nil {
		return fake.SchemaIdStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.schemaIdReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *SchemaRegistry) SchemaIdCallCount() int {
	fake.schemaIdMutex.RLock()
	defer fake.schemaIdMutex.RUnlock()
	return len(fake.schemaIdArgsForCall)
}

func (fake *SchemaRegistry) SchemaIdCalls(stub func(string, string) (uint32, error)) {
	fake.schemaIdMutex.Lock()
	defer fake.schemaIdMutex.Unlock()
	fake.SchemaIdStub = stub
}

func (fake *SchemaRegistry) SchemaIdArgsForCall(i int) (string, string) {
	fake.schemaIdMutex.RLock()
	defer fake.schemaIdMutex.RUnlock()
	argsForCall := fake.schemaIdArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *SchemaRegistry) SchemaIdReturns(result1 uint32, result2 error) {
	fake.schemaIdMutex.Lock()
	defer fake.schemaIdMutex.Unlock()
	fake.SchemaIdStub = nil
	fake.schemaIdReturns = struct {
		result1 uint32
		result2 error
	}{result1, result2}
}

func (fake *SchemaRegistry) SchemaIdReturnsOnCall(i int, result1 uint32, result2 error) {
	fake.schemaIdMutex.Lock()
	defer fake.schemaIdMutex.Unlock()
	fake.SchemaIdStub = nil
	if fake.schemaIdReturnsOnCall == nil {
		fake.schemaIdReturnsOnCall = make(map[int]struct {
			result1 uint32
			result2 error
		})
	}
	fake.schemaIdReturnsOnCall[i] = struct {
		result1 uint32
		result2 error
	}{result1, result2}
}

func (fake *SchemaRegistry) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.schemaIdMutex.RLock()
	defer fake.schemaIdMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *SchemaRegistry) recordInvocation(key string, args []interface{}) {
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

var _ schema.Registry = new(SchemaRegistry)
