// Copyright SecureKey Technologies Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0

module github.com/trustbloc/orb/test/bdd

require (
	github.com/cenkalti/backoff/v4 v4.1.1
	github.com/containerd/containerd v1.5.5 // indirect
	github.com/cucumber/godog v0.9.0
	github.com/cucumber/messages-go/v10 v10.0.3
	github.com/fsouza/go-dockerclient v1.6.5
	github.com/google/uuid v1.3.0
	github.com/hyperledger/aries-framework-go v0.1.7-0.20210816171017-5da380dba24e
	github.com/hyperledger/aries-framework-go-ext/component/storage/couchdb v0.0.0-20210826164831-40568174ea45
	github.com/hyperledger/aries-framework-go-ext/component/storage/mongodb v0.0.0-20210903215754-11447fcf4d91
	github.com/hyperledger/aries-framework-go-ext/component/vdr/orb v0.0.0-20210908194819-812cc66d2834
	github.com/hyperledger/aries-framework-go/component/storageutil v0.0.0-20210807121559-b41545a4f1e8
	github.com/hyperledger/aries-framework-go/spi v0.0.0-20210902194940-97c6f2cded6c
	github.com/igor-pavlenko/httpsignatures-go v0.0.21
	github.com/ipfs/go-ipfs-api v0.2.0
	github.com/jamiealquiza/tachymeter v2.0.0+incompatible
	github.com/moby/sys/mount v0.2.0 // indirect
	github.com/mr-tron/base58 v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/tidwall/gjson v1.7.4
	github.com/trustbloc/orb v0.1.3-0.20210826224204-8f7cf7841ff2
	github.com/trustbloc/sidetree-core-go v0.6.1-0.20210910132742-a2e8795453c1
)

replace github.com/trustbloc/orb => ../../

go 1.16
