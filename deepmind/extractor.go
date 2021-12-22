package deepmind

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/figment-networks/extractor-tendermint"
	ext "github.com/figment-networks/extractor-tendermint"
	"github.com/tendermint/tendermint/types"
)

var (
	enabled            bool
	currentBlockHeight int64
	inBlock            bool
)

func Enable() {
	enabled = true
}

func IsEnabled() bool {
	return enabled
}

func Disable() {
	enabled = false
}

func Initialize(config *extractor.Config) {
	extractor.SetWriterFromConfig(config)
	enabled = true
}

func Shutdown(ctx context.Context) {
	defer func() {
		enabled = false
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(time.Second):
			if !inBlock {
				return
			}
		}
	}
}

func BeginBlock(height int64) error {
	if !enabled {
		return nil
	}

	if height == currentBlockHeight {
		panic("cannot initialize the same height more than once")
	}

	currentBlockHeight = height
	inBlock = true

	err := ext.SetHeight(height)
	if err != nil {
		return err
	}
	return ext.WriteLine(ext.MsgBegin, "%d", height)
}

func FinalizeBlock(height int64) error {
	if !enabled {
		return nil
	}

	if height != currentBlockHeight {
		panic("finalize block on invalid height")
	}

	inBlock = false

	return ext.WriteLine(ext.MsgEnd, "%d", height)
}

func AddBlockData(data types.EventDataNewBlock) error {
	if !enabled {
		return nil
	}

	buff, err := encodeBlock(data)
	if err != nil {
		return err
	}

	return ext.WriteLine(ext.MsgBlock, "%s", base64.StdEncoding.EncodeToString(buff))
}

func AddBlockHeaderData(data types.EventDataNewBlockHeader) error {
	if !enabled {
		return nil
	}
	// Skipped for now
	return nil
}

func AddEvidenceData(data types.EventDataNewEvidence) error {
	if !enabled {
		return nil
	}
	// Skipped for now
	return nil
}

func AddTransactionData(tx types.EventDataTx) error {
	if !enabled {
		return nil
	}

	data, err := encodeTx(&tx.TxResult)
	if err != nil {
		return err
	}

	return ext.WriteLine(ext.MsgTx, "%s", base64.StdEncoding.EncodeToString(data))
}

func AddValidatorSetUpdatesData(updates types.EventDataValidatorSetUpdates) error {
	if !enabled {
		return nil
	}

	data, err := encodeValidatorSetUpdates(&updates)
	if err != nil {
		return err
	}

	return ext.WriteLine(ext.MsgValidatorSetUpdate, "%s", base64.StdEncoding.EncodeToString(data))
}
