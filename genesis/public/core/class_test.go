/*
 *    Copyright 2018 INS Ecosystem
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package core

import (
	"testing"

	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/contract"
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestClassDomain_GetClassID(t *testing.T) {
	parent := &mockParent{}
	clsDom, err := newClassDomain(parent)
	assert.NoError(t, err)

	domainID := clsDom.GetClassID()

	assert.Equal(t, class.ClsDomainID, domainID)
}

func TestClassDomain_GetClass_NoSuchRecord(t *testing.T) {
	parent := &mockParent{}
	classDom, err := newClassDomain(parent)
	assert.NoError(t, err)
	cl, err := classDom.GetClass("test")

	assert.EqualError(t, err, "object with record test does not exist")
	assert.Nil(t, cl)
}

func TestClassDomain_GetClass(t *testing.T) {
	parent := &mockParent{}
	classDom, err := newClassDomain(parent)
	assert.NoError(t, err)

	classFactory := NewClassDomainFactory(parent)
	recordId, regErr := classDom.RegisterClass(classFactory)
	assert.NoError(t, regErr)

	resolved, err := classDom.GetClass(recordId)

	assert.NoError(t, err)
	assert.Equal(t, classFactory, resolved)
}

func TestNewClassDomain(t *testing.T) {
	parent := &mockParent{}
	classDom, err := newClassDomain(parent)

	assert.NoError(t, err)
	assert.Equal(t, &classDomain{
		BaseDomain: *domain.NewBaseDomain(parent, ClassDomainName),
	}, classDom)
}

func TestNewClassDomain_WithNoParent(t *testing.T) {
	_, err := newClassDomain(nil)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestClassDomain_RegisterClass(t *testing.T) {
	parent := &mockParent{}
	classDom, err := newClassDomain(parent)

	assert.NoError(t, err)

	recordId, regErr := classDom.RegisterClass(NewClassDomainFactory(parent))

	assert.NoError(t, regErr)

	_, err = uuid.FromString(recordId)
	assert.NoError(t, err)
}

func TestNewClassDomainFactory(t *testing.T) {
	parent := &mockParent{}
	expected := &classDomainFactory{parent: parent}
	factory := NewClassDomainFactory(parent)

	assert.Equal(t, expected, factory)
}

func TestClassDomainFactory_GetParent(t *testing.T) {
	parent := &mockParent{}
	factory := NewClassDomainFactory(parent)

	assert.Nil(t, factory.GetParent())
}

func TestClassDomainFactory_GetClassID(t *testing.T) {
	parent := &mockParent{}
	factory := NewClassDomainFactory(parent)
	classId := factory.GetClassID()
	assert.Equal(t, class.ClsDomainID, classId)
}

func TestClassDomainFactory_Create(t *testing.T) {
	parent := &mockParent{}
	factory := NewClassDomainFactory(parent)
	proxy, err := factory.Create(parent)
	assert.NoError(t, err)

	classDmn, err := newClassDomain(parent)
	assert.NoError(t, err)
	assert.Equal(t, &classDomainProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: classDmn,
		},
	}, proxy)
}

func TestClassDomainFactory_Create_NoParent(t *testing.T) {
	factory := NewClassDomainFactory(nil)
	_, err := factory.Create(nil)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestNewClassDomainProxy(t *testing.T) {
	parent := &mockParent{}
	clDomainProxy, err := newClassDomainProxy(parent)

	assert.NoError(t, err)

	newClDomain, clErr := newClassDomain(parent)
	assert.NoError(t, clErr)

	assert.Equal(t, &classDomainProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: newClDomain,
		},
	}, clDomainProxy)
}

func TestNewClassDomainProxy_Error(t *testing.T) {
	_, err := newClassDomainProxy(nil)

	assert.EqualError(t, err, "parent must not be nil")
}

func TestClassDomainProxy_GetClass(t *testing.T) {

	parent := &mockParent{}
	clDomainProxy, err := newClassDomainProxy(parent)
	assert.NoError(t, err)

	classFactory := NewClassDomainFactory(parent)
	recordId, regErr := clDomainProxy.RegisterClass(classFactory)

	assert.NoError(t, regErr)

	resolved, err := clDomainProxy.GetClass(recordId)

	assert.NoError(t, err)
	assert.Equal(t, classFactory, resolved)
}

func TestClassDomainProxy_GetParent(t *testing.T) {
	parent := &mockParent{}
	clDomainProxy, err := newClassDomainProxy(parent)
	assert.NoError(t, err)

	actualParent := clDomainProxy.GetParent()
	assert.Equal(t, parent, actualParent)
}

func TestClassDomainProxy_RegisterClass(t *testing.T) {
	parent := &mockParent{}
	clDomainProxy, err := newClassDomainProxy(parent)
	assert.NoError(t, err)

	regist, err := clDomainProxy.RegisterClass(NewClassDomainFactory(parent))
	assert.NoError(t, err)

	_, err = uuid.FromString(regist)
	assert.NoError(t, err)
}

func TestClassDomainProxy_GetClassID(t *testing.T) {
	parent := &mockParent{}
	clDomainProxy, err := newClassDomainProxy(parent)
	assert.NoError(t, err)
	assert.Equal(t, class.ClsDomainID, clDomainProxy.GetClassID())
}
