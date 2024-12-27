package key

import (
	"encoding/hex"
	"fmt"

	"encoding/base64"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	legacybech32 "github.com/cosmos/cosmos-sdk/types/bech32/legacybech32"
)

// pubkey, err := getPubKeyFromRawString("c6OnDc85QAwfND3URSqIE/7eZUPbZ+F4Fje8bqG8x44=", "ed25519")
// if err != nil {
// 	panic(err)
// }

func BytesToPubkey(bz []byte, keytype string) (cryptotypes.PubKey, bool) {
	if keytype == "ed25519" {
		if len(bz) == ed25519.PubKeySize {
			return &ed25519.PubKey{Key: bz}, true
		}
	}

	if len(bz) == secp256k1.PubKeySize {
		return &secp256k1.PubKey{Key: bz}, true
	}
	return nil, false
}

func GetPubKeyFromRawString(pkstr, keytype string) (cryptotypes.PubKey, error) {
	// Try hex decoding
	bz, err := hex.DecodeString(pkstr)
	if err == nil {
		pk, ok := BytesToPubkey(bz, keytype)
		if ok {
			return pk, nil
		}
	}

	bz, err = base64.StdEncoding.DecodeString(pkstr)
	if err == nil {
		pk, ok := BytesToPubkey(bz, keytype)
		if ok {
			return pk, nil
		}
	}

	pk, err := legacybech32.UnmarshalPubKey(legacybech32.AccPK, pkstr) //nolint:staticcheck // we do old keys, they're keys after all.
	if err == nil {
		return pk, nil
	}

	pk, err = legacybech32.UnmarshalPubKey(legacybech32.ValPK, pkstr) //nolint:staticcheck // we do old keys, they're keys after all.
	if err == nil {
		return pk, nil
	}

	pk, err = legacybech32.UnmarshalPubKey(legacybech32.ConsPK, pkstr) //nolint:staticcheck // we do old keys, they're keys after all.
	if err == nil {
		return pk, nil
	}

	return nil, fmt.Errorf("pubkey '%s' invalid; expected hex, base64, or bech32 of correct size", pkstr)
}
