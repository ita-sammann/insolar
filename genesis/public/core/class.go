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
	"fmt"

	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/insolar/insolar/genesis/model/factory"
	"github.com/insolar/insolar/genesis/model/object"
)

// ClassDomainName is a name for class domain.
const ClassDomainName = "ClassDomain"

// ClassDomain is a contract that allows to publish new classes (e.g. new contract types).
type ClassDomain interface {
	// Base domain implementation.
	domain.Domain
	// RegisterClass is used to publish new .
	RegisterClass(factory.Factory) (string, error)
	// GetClass provides factory instance from record.
	GetClass(string) (factory.Factory, error)
}

type classDomain struct {
	domain.BaseDomain
}

// newClassDomain creates new instance of ClassDomain
func newClassDomain(parent object.Parent) (*classDomain, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent must not be nil")
	}

	classDomain := &classDomain{
		BaseDomain: *domain.NewBaseDomain(parent, ClassDomainName),
	}
	return classDomain, nil
}

// GetClassID return string representation of ClassDomain's class.
func (cd *classDomain) GetClassID() string {
	return class.ClsDomainID
}

// RegisterClass method used to create new public Class.
func (cd *classDomain) RegisterClass(fc factory.Factory) (string, error) {
	recordID, err := cd.AddChild(fc)
	if err != nil {
		return "", fmt.Errorf("class registration error")
	}

	return recordID, nil
}

// GetClass method used for retrieve class information from record.
func (cd *classDomain) GetClass(recordID string) (factory.Factory, error) {
	cls, err := cd.GetChild(recordID)
	if err != nil {
		return nil, err
	}

	result, ok := cls.(factory.Factory)
	if !ok {
		return nil, fmt.Errorf("object with record `%s` is not a Class", recordID)
	}

	return result, nil
}

type classDomainProxy struct {
	instance *classDomain
}

// newClassDomainProxy creates new proxy and associate it with new instance of ClassDomain.
func newClassDomainProxy(parent object.Parent) (*classDomainProxy, error) {
	instance, err := newClassDomain(parent)
	if err != nil {
		return nil, err
	}

	return &classDomainProxy{
		instance: instance,
	}, nil
}

// RegisterClass proxy call for instance method.
func (cdp *classDomainProxy) RegisterClass(fc factory.Factory) (string, error) {
	return cdp.instance.RegisterClass(fc)
}

// GetClass proxy call for instance method.
func (cdp *classDomainProxy) GetClass(record string) (factory.Factory, error) {
	return cdp.instance.GetClass(record)
}

// GetReference proxy call for instance method.
func (cdp *classDomainProxy) GetReference() *object.Reference {
	return cdp.instance.GetReference()
}

// GetParent proxy call for instance method.
func (cdp *classDomainProxy) GetParent() object.Parent {
	return cdp.instance.GetParent()
}

// GetClassID proxy call for instance method.
func (cdp *classDomainProxy) GetClassID() string {
	return class.ClsDomainID
}

type classDomainFactory struct {
	parent object.Parent
}

// NewClassDomainFactory creates new factory for ClassDomain.
func NewClassDomainFactory(parent object.Parent) factory.Factory {
	return &classDomainFactory{
		parent: parent,
	}
}

// GetParent is a method that return parent link
func (cdf *classDomainFactory) GetParent() object.Parent {
	// TODO: return real parent, fix tests
	return nil
}

// GetClassID return string representation of ClassDomain's class.
func (cdf *classDomainFactory) GetClassID() string {
	return class.ClsDomainID
}

// GetReference returns nil for not published factory
func (cdf *classDomainFactory) GetReference() *object.Reference {
	return nil
}

// Create factory method for new ClassDomain instances.
func (cdf *classDomainFactory) Create(parent object.Parent) (object.Proxy, error) {
	proxy, err := newClassDomainProxy(parent)
	if err != nil {
		return nil, err
	}

	_, err = parent.AddChild(proxy)
	if err != nil {
		return nil, err
	}
	return proxy, nil
}
