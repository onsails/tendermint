package deepmind

import (
	"fmt"
	"time"

	"github.com/figment-networks/tendermint-protobuf-def/codec"
	abci "github.com/tendermint/tendermint/abci/types"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"
	"github.com/tendermint/tendermint/types"
)

func mapBlockID(bid types.BlockID) *codec.BlockID {
	return &codec.BlockID{
		Hash: bid.Hash,
		PartSetHeader: &codec.PartSetHeader{
			Total: bid.PartSetHeader.Total,
			Hash:  bid.PartSetHeader.Hash,
		},
	}
}

func mapProposer(val *types.Validator) *codec.Validator {
	nPK := &codec.PublicKey{}

	return &codec.Validator{
		Address:          val.Address,
		PubKey:           nPK,
		ProposerPriority: 0,
	}
}

func mapEvent(ev abci.Event) *codec.Event {
	cev := &codec.Event{Eventtype: ev.Type}

	for _, at := range ev.Attributes {
		cev.Attributes = append(cev.Attributes, &codec.EventAttribute{
			Key:   string(at.Key),
			Value: string(at.Value),
			Index: at.Index,
		})
	}

	return cev
}

func mapVote(edv *types.Vote) *codec.EventVote {
	return &codec.EventVote{
		Eventvotetype: codec.SignedMsgType(edv.Type),
		Height:        uint64(edv.Height),
		Round:         edv.Round,
		BlockId:       mapBlockID(edv.BlockID),
		Timestamp:     mapTimestamp(edv.Timestamp),
		ValidatorAddress: &codec.Address{
			Address: edv.ValidatorAddress,
		},
		ValidatorIndex: edv.ValidatorIndex,
		Signature:      edv.Signature,
	}
}

func mapSignatures(commitSignatures []types.CommitSig) ([]*codec.CommitSig, error) {
	signatures := make([]*codec.CommitSig, len(commitSignatures))
	for i, commitSignature := range commitSignatures {
		signature, err := mapSignature(commitSignature)
		if err != nil {
			return nil, err
		}
		signatures[i] = signature
	}
	return signatures, nil
}

func mapSignature(s types.CommitSig) (*codec.CommitSig, error) {
	return &codec.CommitSig{
		BlockIdFlag:      codec.BlockIDFlag(s.BlockIDFlag),
		ValidatorAddress: &codec.Address{Address: s.ValidatorAddress.Bytes()},
		Timestamp:        mapTimestamp(s.Timestamp),
		Signature:        s.Signature,
	}, nil
}

func mapValidatorUpdate(v abci.ValidatorUpdate) (*codec.Validator, error) {
	nPK := &codec.PublicKey{}
	var address []byte

	switch key := v.PubKey.Sum.(type) {
	case *crypto.PublicKey_Ed25519:
		nPK.Sum = &codec.PublicKey_Ed25519{Ed25519: key.Ed25519}
		address = tmcrypto.AddressHash(nPK.GetEd25519())
	case *crypto.PublicKey_Secp256K1:
		nPK.Sum = &codec.PublicKey_Secp256K1{Secp256K1: key.Secp256K1}
		address = tmcrypto.AddressHash(nPK.GetSecp256K1())
	default:
		return nil, fmt.Errorf("given type %T of PubKey mapping doesn't exist ", key)
	}

	// NOTE on ProposerPriority field: Priority value seems to be calcuulated in the
	// context of the validator set. We're already processing this as a separate event.
	// More info in here: https://docs.tendermint.com/v0.34/spec/consensus/proposer-selection.html

	return &codec.Validator{
		Address:          address,
		PubKey:           nPK,
		VotingPower:      v.Power,
		ProposerPriority: 0,
	}, nil
}

func mapValidators(srcValidators []*types.Validator) ([]*codec.Validator, error) {
	validators := make([]*codec.Validator, len(srcValidators))
	for i, validator := range srcValidators {
		val, err := mapValidator(validator)
		if err != nil {
			return nil, err
		}
		validators[i] = val
	}
	return validators, nil
}

func mapValidator(v *types.Validator) (*codec.Validator, error) {
	nPK := &codec.PublicKey{}

	key := v.PubKey
	switch key.Type() {
	case ed25519.KeyType:
		nPK = &codec.PublicKey{
			Sum: &codec.PublicKey_Ed25519{Ed25519: key.Bytes()}}
	case secp256k1.KeyType:
		nPK = &codec.PublicKey{
			Sum: &codec.PublicKey_Secp256K1{Secp256K1: key.Bytes()}}
	default:
		return nil, fmt.Errorf("given type %T of PubKey mapping doesn't exist ", key)
	}

	// NOTE: See note in mapValidatorUpdate() about ProposerPriority

	return &codec.Validator{
		Address:          v.Address,
		PubKey:           nPK,
		VotingPower:      v.VotingPower,
		ProposerPriority: 0,
	}, nil
}

func mapTimestamp(time time.Time) *codec.Timestamp {
	return &codec.Timestamp{
		Seconds: time.Unix(),
		Nanos:   int32(time.UnixNano() - time.Unix()*1000000000),
	}
}
