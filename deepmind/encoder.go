package deepmind

import (
	"fmt"

	"github.com/figment-networks/tendermint-protobuf-def/codec"
	"github.com/golang/protobuf/proto"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/types"
)

func encodeBlock(bh types.EventDataNewBlock) ([]byte, error) {
	mappedCommitSignatures, err := mapSignatures(bh.Block.LastCommit.Signatures)
	if err != nil {
		return nil, err
	}

	nb := &codec.EventBlock{
		Block: &codec.Block{
			Header: &codec.Header{
				Version: &codec.Consensus{
					Block: bh.Block.Header.Version.Block,
					App:   bh.Block.Header.Version.App,
				},
				ChainId:            bh.Block.Header.ChainID,
				Height:             uint64(bh.Block.Header.Height),
				Time:               mapTimestamp(bh.Block.Header.Time),
				LastBlockId:        mapBlockID(bh.Block.LastBlockID),
				LastCommitHash:     bh.Block.Header.LastCommitHash,
				DataHash:           bh.Block.Header.DataHash,
				ValidatorsHash:     bh.Block.Header.ValidatorsHash,
				NextValidatorsHash: bh.Block.Header.NextValidatorsHash,
				ConsensusHash:      bh.Block.Header.ConsensusHash,
				AppHash:            bh.Block.Header.AppHash,
				LastResultsHash:    bh.Block.Header.LastResultsHash,
				EvidenceHash:       bh.Block.Header.EvidenceHash,
				ProposerAddress:    bh.Block.Header.ProposerAddress,
			},
			LastCommit: &codec.Commit{
				Height:     bh.Block.LastCommit.Height,
				Round:      bh.Block.LastCommit.Round,
				BlockId:    mapBlockID(bh.Block.LastCommit.BlockID),
				Signatures: mappedCommitSignatures,
			},
			Evidence: &codec.EvidenceList{},
			Data: &codec.Data{
				Txs: mapTxs(bh.Block.Data.Txs),
			},
		},
	}

	nb.BlockId = &codec.BlockID{
		Hash: bh.Block.Header.Hash(),
		PartSetHeader: &codec.PartSetHeader{
			Total: bh.Block.LastBlockID.PartSetHeader.Total,
			Hash:  bh.Block.LastBlockID.PartSetHeader.Hash,
		},
	}

	if len(bh.Block.Evidence.Evidence) > 0 {
		for _, ev := range bh.Block.Evidence.Evidence {

			newEv := &codec.Evidence{}
			switch evN := ev.(type) {
			case *types.DuplicateVoteEvidence:
				newEv.Sum = &codec.Evidence_DuplicateVoteEvidence{
					DuplicateVoteEvidence: &codec.DuplicateVoteEvidence{
						VoteA:            mapVote(evN.VoteA),
						VoteB:            mapVote(evN.VoteB),
						TotalVotingPower: evN.TotalVotingPower,
						ValidatorPower:   evN.ValidatorPower,
						Timestamp:        mapTimestamp(evN.Timestamp),
					},
				}
			case *types.LightClientAttackEvidence:
				mappedSetValidators, err := mapValidators(evN.ConflictingBlock.ValidatorSet.Validators)
				if err != nil {
					return nil, err
				}

				mappedByzantineValidators, err := mapValidators(evN.ByzantineValidators)
				if err != nil {
					return nil, err
				}

				mappedCommitSignatures, err := mapSignatures(evN.ConflictingBlock.Commit.Signatures)
				if err != nil {
					return nil, err
				}

				newEv.Sum = &codec.Evidence_LightClientAttackEvidence{
					LightClientAttackEvidence: &codec.LightClientAttackEvidence{
						ConflictingBlock: &codec.LightBlock{
							SignedHeader: &codec.SignedHeader{
								Header: &codec.Header{
									Version: &codec.Consensus{
										Block: evN.ConflictingBlock.Version.Block,
										App:   evN.ConflictingBlock.Version.App,
									},
									ChainId:            evN.ConflictingBlock.Header.ChainID,
									Height:             uint64(evN.ConflictingBlock.Header.Height),
									Time:               mapTimestamp(evN.ConflictingBlock.Header.Time),
									LastBlockId:        mapBlockID(evN.ConflictingBlock.Header.LastBlockID),
									LastCommitHash:     evN.ConflictingBlock.Header.LastCommitHash,
									DataHash:           evN.ConflictingBlock.Header.DataHash,
									ValidatorsHash:     evN.ConflictingBlock.Header.ValidatorsHash,
									NextValidatorsHash: evN.ConflictingBlock.Header.NextValidatorsHash,
									ConsensusHash:      evN.ConflictingBlock.Header.ConsensusHash,
									AppHash:            evN.ConflictingBlock.Header.AppHash,
									LastResultsHash:    evN.ConflictingBlock.Header.LastResultsHash,
									EvidenceHash:       evN.ConflictingBlock.Header.EvidenceHash,
									ProposerAddress:    evN.ConflictingBlock.Header.ProposerAddress,
								},
								Commit: &codec.Commit{
									Height:     evN.ConflictingBlock.Commit.Height,
									Round:      evN.ConflictingBlock.Commit.Round,
									BlockId:    mapBlockID(evN.ConflictingBlock.Commit.BlockID),
									Signatures: mappedCommitSignatures,
								},
							},
							ValidatorSet: &codec.ValidatorSet{
								Validators:       mappedSetValidators,
								Proposer:         mapProposer(evN.ConflictingBlock.ValidatorSet.Proposer),
								TotalVotingPower: evN.ConflictingBlock.ValidatorSet.TotalVotingPower(),
							},
						},
						CommonHeight:        evN.CommonHeight,
						ByzantineValidators: mappedByzantineValidators,
						TotalVotingPower:    evN.TotalVotingPower,
						Timestamp:           mapTimestamp(evN.Timestamp),
					},
				}

			default:
				return nil, fmt.Errorf("given type %T of EvidenceList mapping doesn't exist ", ev)
			}

			nb.Block.Evidence.Evidence = append(nb.Block.Evidence.Evidence, newEv)
		}
	}

	if len(bh.ResultBeginBlock.Events) > 0 {
		nb.ResultBeginBlock = &codec.ResponseBeginBlock{}
		for _, ev := range bh.ResultBeginBlock.Events {
			nb.ResultBeginBlock.Events = append(nb.ResultBeginBlock.Events, mapEvent(ev))
		}
	}

	if len(bh.ResultEndBlock.Events) > 0 || len(bh.ResultEndBlock.ValidatorUpdates) > 0 || bh.ResultEndBlock.ConsensusParamUpdates != nil {
		nb.ResultEndBlock = &codec.ResponseEndBlock{
			ConsensusParamUpdates: &codec.ConsensusParams{},
		}

		for _, ev := range bh.ResultEndBlock.Events {
			nb.ResultEndBlock.Events = append(nb.ResultEndBlock.Events, mapEvent(ev))
		}

		for _, v := range bh.ResultEndBlock.ValidatorUpdates {
			val, err := mapValidatorUpdate(v)
			if err != nil {
				return nil, err
			}
			nb.ResultEndBlock.ValidatorUpdates = append(nb.ResultEndBlock.ValidatorUpdates, val)
		}
	}

	return proto.Marshal(nb)
}

func encodeTx(result *abci.TxResult) ([]byte, error) {
	tx := &codec.EventTx{
		TxResult: &codec.TxResult{
			Height: uint64(result.Height),
			Index:  result.Index,
			Tx:     mapTx(result.Tx),
			Result: &codec.ResponseDeliverTx{
				Code:      result.Result.Code,
				Data:      result.Result.Data,
				Log:       result.Result.Log,
				Info:      result.Result.Info,
				GasWanted: result.Result.GasWanted,
				GasUsed:   result.Result.GasUsed,
				Codespace: result.Result.Codespace,
			},
		},
	}

	for _, ev := range result.Result.Events {
		tx.TxResult.Result.Events = append(tx.TxResult.Result.Events, mapEvent(ev))
	}

	return proto.Marshal(tx)
}

func encodeValidatorSetUpdates(updates *types.EventDataValidatorSetUpdates) ([]byte, error) {
	result := &codec.EventValidatorSetUpdates{}

	for _, update := range updates.ValidatorUpdates {
		nPK := &codec.PublicKey{}

		switch update.PubKey.Type() {
		case "ed25519":
			nPK.Sum = &codec.PublicKey_Ed25519{Ed25519: update.PubKey.Bytes()}
		case "secp256k1":
			nPK.Sum = &codec.PublicKey_Secp256K1{Secp256K1: update.PubKey.Bytes()}
		default:
			return nil, fmt.Errorf("unsupported pubkey type: %T", update.PubKey)
		}

		result.ValidatorUpdates = append(result.ValidatorUpdates, &codec.Validator{
			Address:          update.Address.Bytes(),
			VotingPower:      update.VotingPower,
			ProposerPriority: update.ProposerPriority,
			PubKey:           nPK,
		})
	}

	return proto.Marshal(result)
}
