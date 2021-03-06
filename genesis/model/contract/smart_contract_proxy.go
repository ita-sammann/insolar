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

package contract

import (
	"fmt"

	"github.com/insolar/insolar/genesis/model/factory"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/insolar/insolar/genesis/model/resolver"
)

// BaseSmartContractProxy is a implementation of Proxy with Instance.
type BaseSmartContractProxy struct {
	resolver.BaseProxy
	Instance object.Child
}

// GetClassID is a proxy call for instance method.
func (bp *BaseSmartContractProxy) GetClassID() string {
	return bp.Instance.GetClassID()
}

// GetParent is a proxy call for instance method.
func (bp *BaseSmartContractProxy) GetParent() object.Parent {
	return bp.Instance.GetParent()
}

// GetOrCreateComposite is a proxy call for instance method if instance inherit ComposingContainer.
func (bp *BaseSmartContractProxy) GetOrCreateComposite(interfaceKey string, compositeFactory factory.CompositeFactory) (factory.Composite, error) {
	sc, ok := bp.Instance.(factory.ComposingContainer)
	if !ok {
		return nil, fmt.Errorf("Instance is not Composing Container")
	}
	return sc.GetOrCreateComposite(compositeFactory)
}
