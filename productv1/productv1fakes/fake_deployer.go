// This file was generated by counterfeiter
package productv1fakes

import (
	"sync"

	"github.com/enaml-ops/pluginlib/cred"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/productv1"
)

type FakeDeployer struct {
	GetProductStub        func(args []string, cloudConfig []byte, cs cred.Store) ([]byte, error)
	getProductMutex       sync.RWMutex
	getProductArgsForCall []struct {
		args        []string
		cloudConfig []byte
		cs          cred.Store
	}
	getProductReturns struct {
		result1 []byte
		result2 error
	}
	GetMetaStub        func() product.Meta
	getMetaMutex       sync.RWMutex
	getMetaArgsForCall []struct{}
	getMetaReturns     struct {
		result1 product.Meta
	}
	GetFlagsStub        func() []pcli.Flag
	getFlagsMutex       sync.RWMutex
	getFlagsArgsForCall []struct{}
	getFlagsReturns     struct {
		result1 []pcli.Flag
	}
}

func (fake *FakeDeployer) GetProduct(args []string, cloudConfig []byte, cs cred.Store) ([]byte, error) {
	fake.getProductMutex.Lock()
	fake.getProductArgsForCall = append(fake.getProductArgsForCall, struct {
		args        []string
		cloudConfig []byte
		cs          cred.Store
	}{args, cloudConfig, cs})
	fake.getProductMutex.Unlock()
	if fake.GetProductStub != nil {
		return fake.GetProductStub(args, cloudConfig, cs)
	} else {
		return fake.getProductReturns.result1, fake.getProductReturns.result2
	}
}

func (fake *FakeDeployer) GetProductCallCount() int {
	fake.getProductMutex.RLock()
	defer fake.getProductMutex.RUnlock()
	return len(fake.getProductArgsForCall)
}

func (fake *FakeDeployer) GetProductArgsForCall(i int) ([]string, []byte, cred.Store) {
	fake.getProductMutex.RLock()
	defer fake.getProductMutex.RUnlock()
	return fake.getProductArgsForCall[i].args, fake.getProductArgsForCall[i].cloudConfig, fake.getProductArgsForCall[i].cs
}

func (fake *FakeDeployer) GetProductReturns(result1 []byte, result2 error) {
	fake.GetProductStub = nil
	fake.getProductReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeDeployer) GetMeta() product.Meta {
	fake.getMetaMutex.Lock()
	fake.getMetaArgsForCall = append(fake.getMetaArgsForCall, struct{}{})
	fake.getMetaMutex.Unlock()
	if fake.GetMetaStub != nil {
		return fake.GetMetaStub()
	} else {
		return fake.getMetaReturns.result1
	}
}

func (fake *FakeDeployer) GetMetaCallCount() int {
	fake.getMetaMutex.RLock()
	defer fake.getMetaMutex.RUnlock()
	return len(fake.getMetaArgsForCall)
}

func (fake *FakeDeployer) GetMetaReturns(result1 product.Meta) {
	fake.GetMetaStub = nil
	fake.getMetaReturns = struct {
		result1 product.Meta
	}{result1}
}

func (fake *FakeDeployer) GetFlags() []pcli.Flag {
	fake.getFlagsMutex.Lock()
	fake.getFlagsArgsForCall = append(fake.getFlagsArgsForCall, struct{}{})
	fake.getFlagsMutex.Unlock()
	if fake.GetFlagsStub != nil {
		return fake.GetFlagsStub()
	} else {
		return fake.getFlagsReturns.result1
	}
}

func (fake *FakeDeployer) GetFlagsCallCount() int {
	fake.getFlagsMutex.RLock()
	defer fake.getFlagsMutex.RUnlock()
	return len(fake.getFlagsArgsForCall)
}

func (fake *FakeDeployer) GetFlagsReturns(result1 []pcli.Flag) {
	fake.GetFlagsStub = nil
	fake.getFlagsReturns = struct {
		result1 []pcli.Flag
	}{result1}
}

var _ product.Deployer = new(FakeDeployer)
