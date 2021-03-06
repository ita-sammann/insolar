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

/*
Package example provides smart contracts for building example of system work.

Entities:

	MemberDomain - domain that allows to add new members to system.

	Usage:

		factory := NewMemberDomainFactory(factoryParent)
		mDomain, error := factory.Create(domainParent)

		record, err := mDomain.CreateMember()
		memberProxy, err := mDomain.GetMember(record)

	Member - smart contract that represent user of the system.

	Usage:

		// because type of member is a child of memberDomain
		factory := NewMemberFactory(memberDomain)
		m, error := factory.Create(parent)

////

	WalletDomain - domain that allows to add new wallets to system

	Usage:

		factory := NewWalletDomainFactory(factoryParent)
		wDomain, err := factory.Create(domainParent)

		// since Create inject composite to member
		err := wDomain.CreateWallet(member)

	Wallet - smart contract that represent wallet

	Usage:
		factory := NewWalletFactory(walletDomain)
		w, err := factory.Create(parent)

////

	Get user balance example:
		wFactory := NewWalletFactory(walletDomain)

		w := m.GetOrCreateComposite(wFactory)
		w.GetBalance()

*/
package example
