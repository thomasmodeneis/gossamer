// Copyright 2021 ChainSafe Systems (ON)
// SPDX-License-Identifier: LGPL-3.0-only

package keystore

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/crypto"
	"github.com/ChainSafe/gossamer/lib/crypto/ed25519"
	"github.com/ChainSafe/gossamer/lib/crypto/secp256k1"
	"github.com/ChainSafe/gossamer/lib/crypto/sr25519"
	"github.com/ChainSafe/gossamer/lib/utils"
)

// PrivateKeyToKeypair returns a public, private keypair given a private key
func PrivateKeyToKeypair(priv crypto.PrivateKey) (kp crypto.Keypair, err error) {
	if key, ok := priv.(*sr25519.PrivateKey); ok {
		kp, err = sr25519.NewKeypairFromPrivate(key)
	} else if key, ok := priv.(*ed25519.PrivateKey); ok {
		kp, err = ed25519.NewKeypairFromPrivate(key)
	} else if key, ok := priv.(*secp256k1.PrivateKey); ok {
		kp, err = secp256k1.NewKeypairFromPrivate(key)
	} else {
		return nil, errors.New("cannot decode key: invalid key type")
	}

	return kp, err
}

// DecodePrivateKey turns input bytes into a private key based on the specified key type
func DecodePrivateKey(in []byte, keytype crypto.KeyType) (priv crypto.PrivateKey, err error) {
	if keytype == crypto.Ed25519Type {
		priv, err = ed25519.NewPrivateKey(in)
	} else if keytype == crypto.Sr25519Type {
		priv, err = sr25519.NewPrivateKey(in)
	} else if keytype == crypto.Secp256k1Type {
		priv, err = secp256k1.NewPrivateKey(in)
	} else {
		return nil, errors.New("cannot decode key: invalid key type")
	}

	return priv, err
}

// DecodeKeyPairFromHex turns an hex-encoded private key into a keypair
func DecodeKeyPairFromHex(keystr []byte, keytype crypto.KeyType) (kp crypto.Keypair, err error) {
	switch keytype {
	case crypto.Sr25519Type:
		kp, err = sr25519.NewKeypairFromSeed(keystr)
	case crypto.Ed25519Type:
		kp, err = ed25519.NewKeypairFromSeed(keystr)
	default:
		return nil, errors.New("cannot decode key: invalid key type")
	}

	return kp, err
}

// GenerateKeypair create a new keypair with the corresponding type and saves
// it to basepath/keystore/[public key].key in json format encrypted using the
// specified password and returns the resulting filepath of the new key
func GenerateKeypair(keytype string, kp crypto.Keypair, basepath string, password []byte) (string, error) {
	if keytype == "" {
		keytype = crypto.Sr25519Type
	}
	var err error

	if kp == nil {
		if keytype == crypto.Sr25519Type {
			kp, err = sr25519.GenerateKeypair()
			if err != nil {
				return "", fmt.Errorf("failed to generate sr25519 keypair: %s", err)
			}
		} else if keytype == crypto.Ed25519Type {
			kp, err = ed25519.GenerateKeypair()
			if err != nil {
				return "", fmt.Errorf("failed to generate ed25519 keypair: %s", err)
			}
		} else if keytype == crypto.Secp256k1Type {
			kp, err = secp256k1.GenerateKeypair()
			if err != nil {
				return "", fmt.Errorf("failed to generate secp256k1 keypair: %s", err)
			}
		}
	}

	keyPath, err := utils.KeystoreDir(basepath)
	if err != nil {
		return "", fmt.Errorf("failed to get keystore directory: %s", err)
	}

	pub := hex.EncodeToString(kp.Public().Encode())
	fp, err := filepath.Abs(keyPath + "/" + pub + ".key")
	if err != nil {
		return "", fmt.Errorf("failed to create absolute filepath: %s", err)
	}

	err = EncryptAndWriteToFile(fp, kp.Private(), password)
	if err != nil {
		return "", fmt.Errorf("failed to write key to file: %s", err)
	}

	return fp, nil
}

// LoadKeystore loads a new keystore and inserts the test key into the keystore
func LoadKeystore(key string, ks Keystore) error {
	if key != "" {

		var kr Keyring
		var err error

		switch ks.Type() {
		case crypto.Ed25519Type:
			kr, err = NewEd25519Keyring()
			if err != nil {
				return fmt.Errorf("failed to create keyring: %s", err)
			}
		default:
			kr, err = NewSr25519Keyring()
			if err != nil {
				return fmt.Errorf("failed to create keyring: %s", err)
			}
		}

		switch strings.ToLower(key) {
		// Insert can error only if kestore type do not match with key
		// type do not match. Since we have created keyring based on ks.Type(),
		// Insert would never error here. Thus, ignoring those errors.
		case "alice":
			_ = ks.Insert(kr.Alice())
		case "bob":
			_ = ks.Insert(kr.Bob())
		case "charlie":
			_ = ks.Insert(kr.Charlie())
		case "dave":
			_ = ks.Insert(kr.Dave())
		case "eve":
			_ = ks.Insert(kr.Eve())
		case "ferdie":
			_ = ks.Insert(kr.Ferdie())
		case "george":
			_ = ks.Insert(kr.George())
		case "heather":
			_ = ks.Insert(kr.Heather())
		case "ian":
			_ = ks.Insert(kr.Ian())
		default:
			return fmt.Errorf("invalid test key provided")
		}
	}

	return nil
}

// ImportKeypair imports a key specified by its filename into a subdirectory
// by the name "keystore" and saves it under the filename "[publickey].key",
// returns the absolute path of the imported key file
func ImportKeypair(fp, dir string) (string, error) {
	keyDir, err := utils.KeystoreDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to create keystore directory: %s", err)
	}

	keyData, err := os.ReadFile(filepath.Clean(fp))
	if err != nil {
		return "", fmt.Errorf("failed to read keystore file: %s", err)
	}

	keystore := new(EncryptedKeystore)
	err = json.Unmarshal(keyData, keystore)
	if err != nil {
		return "", fmt.Errorf("failed to read import keystore data: %s", err)
	}

	keyFilePath, err := filepath.Abs(keyDir + "/" + keystore.PublicKey[2:] + ".key")
	if err != nil {
		return "", fmt.Errorf("failed to create keystore filepath: %s", err)
	}

	err = os.WriteFile(keyFilePath, keyData, 0600)
	if err != nil {
		return "", fmt.Errorf("failed to write to keystore file: %s", err)
	}

	return keyFilePath, nil
}

// ImportRawPrivateKey imports a raw private key and saves it to the keystore directory
func ImportRawPrivateKey(key, keytype, basepath string, password []byte) (string, error) {
	var kp crypto.Keypair
	var err error

	if keytype == "" {
		keytype = crypto.Sr25519Type
	}

	if keytype == crypto.Sr25519Type {
		kp, err = sr25519.NewKeypairFromPrivateKeyString(key)
		if err != nil {
			return "", fmt.Errorf("failed to import sr25519 keypair: %s", err)
		}
	} else if keytype == crypto.Ed25519Type {
		kp, err = ed25519.NewKeypairFromPrivateKeyString(key)
		if err != nil {
			return "", fmt.Errorf("failed to generate ed25519 keypair: %s", err)
		}
	} else if keytype == crypto.Secp256k1Type {
		kp, err = secp256k1.NewKeypairFromPrivateKeyString(key)
		if err != nil {
			return "", fmt.Errorf("failed to generate secp256k1 keypair: %s", err)
		}
	}

	return GenerateKeypair(keytype, kp, basepath, password)
}

// UnlockKeys unlocks keys specified by the --unlock flag with the passwords given by --password
// and places them into the keystore
func UnlockKeys(ks Keystore, dir, unlock, password string) error {
	var indices []int
	var passwords []string
	var err error

	keyDir, err := utils.KeystoreDir(dir)
	if err != nil {
		return err
	}

	if unlock != "" {
		// indices of keys to unlock
		indices, err = common.StringToInts(unlock)
		if err != nil {
			return err
		}
	}

	if password != "" {
		// passwords corresponding to the keys
		passwords = strings.Split(password, ",")
	}

	if len(passwords) != len(indices) {
		return fmt.Errorf("number of passwords given does not match number of keys to unlock")
	}

	// get paths to key files
	keyFiles, err := utils.KeystoreFiles(dir)
	if err != nil {
		return err
	}

	if len(keyFiles) < len(indices) {
		return fmt.Errorf("number of accounts to unlock is greater than number of accounts in keystore")
	}

	// for each key to unlock, read its file and decrypt contents and add to keystore
	for i, idx := range indices {
		if idx >= len(keyFiles) {
			return fmt.Errorf("invalid account index: %d", idx)
		}

		keyFile := keyFiles[idx]
		priv, err := ReadFromFileAndDecrypt(keyDir+"/"+keyFile, []byte(passwords[i]))
		if err != nil {
			return fmt.Errorf("failed to decrypt key file %s: %s", keyFile, err)
		}

		kp, err := PrivateKeyToKeypair(priv)
		if err != nil {
			return fmt.Errorf("failed to create keypair from private key %d: %s", idx, err)
		}

		if err = ks.Insert(kp); err != nil {
			return fmt.Errorf("failed to insert key in keystore: %v", err)
		}
	}

	return nil
}

// DetermineKeyType takes string as defined in https://github.com/w3f/PSPs/blob/psp-rpc-api/psp-002.md#Key-types
//  and returns the crypto.KeyType
func DetermineKeyType(t string) crypto.KeyType {
	switch t {
	case "babe":
		return crypto.Sr25519Type
	case "gran":
		return crypto.Ed25519Type
	case "acco":
		return crypto.Sr25519Type
	case "aura":
		return crypto.Sr25519Type
	case "imon":
		return crypto.Sr25519Type
	case "audi":
		return crypto.Sr25519Type
	case "dumy":
		return crypto.Sr25519Type
	}
	return crypto.UnknownType
}

// HasKey returns true if given hex encoded public key string is found in keystore, false otherwise, error if there
// are issues decoding string
func HasKey(pubKeyStr, keyType string, keystore Keystore) (bool, error) {
	keyBytes, err := common.HexToBytes(pubKeyStr)
	if err != nil {
		return false, err
	}
	cKeyType := DetermineKeyType(keyType)

	var pubKey crypto.PublicKey
	switch cKeyType {
	case crypto.Sr25519Type:
		pubKey, err = sr25519.NewPublicKey(keyBytes)
	case crypto.Ed25519Type:
		pubKey, err = ed25519.NewPublicKey(keyBytes)
	default:
		err = fmt.Errorf("unknown key type: %s", keyType)
	}

	if err != nil {
		return false, err
	}
	key := keystore.GetKeypairFromAddress(pubKey.Address())
	return key != nil, nil
}
