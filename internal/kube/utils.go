package kube

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"

	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

func GenerateClusterSecret() (string, error) {
	const secretSize = 32 // bytes

	buf := make([]byte, secretSize)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}

type BootstrapKeyPair struct {
	PrivateKey string
	PeerID     string
}

func GenerateBootstrapPrivateKey() (*BootstrapKeyPair, error) {
	// Генерируем Ed25519 ключ
	privKey, pubKey, err := libp2pcrypto.GenerateEd25519Key(nil)
	if err != nil {
		return nil, err
	}

	// Marshal private key (protobuf)
	privBytes, err := libp2pcrypto.MarshalPrivateKey(privKey)
	if err != nil {
		return nil, err
	}

	// Получаем PeerID из публичного ключа
	peerID, err := peer.IDFromPublicKey(pubKey)
	if err != nil {
		return nil, err
	}

	return &BootstrapKeyPair{
		PrivateKey: base64.StdEncoding.EncodeToString(privBytes),
		PeerID:     peerID.String(),
	}, nil
}
