module github.com/tendermint/tendermint

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/ChainSafe/go-schnorrkel v0.0.0-20200405005733-88cbf1b4c40d
	github.com/Workiva/go-datastructures v1.0.52
	github.com/adlio/schema v1.1.13
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/confio/ics23/go v0.6.6
	github.com/cosmos/cosmos-proto v1.0.0-alpha7 // indirect
	github.com/cosmos/iavl v0.17.1
	github.com/figment-networks/extractor-cosmos v0.1.0
	github.com/figment-networks/proto-cosmos v0.1.0
	github.com/fortytw2/leaktest v1.3.0
	github.com/go-kit/kit v0.10.0
	github.com/go-logfmt/logfmt v0.5.0
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/google/orderedcode v0.0.1
	github.com/gorilla/websocket v1.4.2
	github.com/gtank/merlin v0.1.1
	github.com/lib/pq v1.2.0
	github.com/libp2p/go-buffer-pool v0.0.2
	github.com/minio/highwayhash v1.0.1
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.10.0
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0
	github.com/rs/cors v1.7.0
	github.com/sasha-s/go-deadlock v0.2.1-0.20190427202633-1595213edefa
	github.com/snikch/goodman v0.0.0-20171125024755-10e37e294daa
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tm-db v0.6.4
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/net v0.0.0-20210903162142-ad29c8ab022f
	google.golang.org/grpc v1.40.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
