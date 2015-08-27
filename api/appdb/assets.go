package appdb

import (
	"golang.org/x/net/context"

	"chain/database/pg"
	"chain/fedchain/wire"
)

// Asset represents an asset type in the blockchain.
// It is made up of extended keys, and paths (indexes) within those keys.
// Assets belong to wallets.
type Asset struct {
	Keys           []*Key
	WIndex, AIndex []uint32
	Hash           wire.Hash20 // the raw Asset ID
	RedeemScript   []byte
	WalletID       string
}

// AssetByID loads an asset from the database using its ID.
func AssetByID(ctx context.Context, id string) (*Asset, error) {
	const q = `
		SELECT keys, redeem_script, wallet_id,
			key_index(wallets.key_index), key_index(assets.key_index),
		FROM assets
		INNER JOIN wallets ON wallets.id=assets.wallet_id
		WHERE assets.id=$1
	`
	var (
		keyIDs []string
		a      = new(Asset)
	)
	var err error
	a.Hash, err = wire.NewHash20FromStr(id)
	if err != nil {
		return nil, err
	}
	err = pg.FromContext(ctx).QueryRow(q, id).Scan(
		(*pg.Strings)(&keyIDs),
		&a.RedeemScript,
		&a.WalletID,
		(*pg.Uint32s)(&a.WIndex),
		(*pg.Uint32s)(&a.AIndex),
	)
	if err != nil {
		return nil, err
	}

	a.Keys, err = getKeys(ctx, keyIDs)
	if err != nil {
		return nil, err
	}

	return a, nil
}
