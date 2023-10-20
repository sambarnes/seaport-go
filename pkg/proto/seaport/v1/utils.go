package seaportv1

import (
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/protobuf/proto"
)

func protoMessage(message any) (proto.Message, error) {
	protoMessage, ok := message.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("%T doesn't implement proto.Message", message)
	}
	return protoMessage, nil
}

// MarshalProto will serialize a given protobuf based struct to bytes
func MarshalProto(message any) ([]byte, error) {
	_message, err := protoMessage(message)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(_message)
}

// AddressToH160 converts an ecdsa.PublicKey address to a v1.H160
func AddressToH160(address common.Address) *H160 {
	data := address.Bytes()
	return &H160{
		Hi: &H128{
			Hi: binary.BigEndian.Uint64(data[:8]),
			Lo: binary.BigEndian.Uint64(data[8:16]),
		},
		Lo: binary.BigEndian.Uint32(data[16:]),
	}
}

// ToAddress converts a H160 to an ethereum address
func (h160 *H160) ToAddress() common.Address {
	data := make([]byte, 20)
	binary.BigEndian.PutUint64(data[:8], h160.Hi.Hi)
	binary.BigEndian.PutUint64(data[8:16], h160.Hi.Lo)
	binary.BigEndian.PutUint32(data[16:], h160.Lo)
	return common.BytesToAddress(data)
}

//
// Test helpers -- TODO: should this be its own package?
//

func MustGenRandomAddress(t *testing.T) (*ecdsa.PrivateKey, common.Address) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}
	return privateKey, crypto.PubkeyToAddress(privateKey.PublicKey)
}

func MustGenRandomH160(t *testing.T) (*ecdsa.PrivateKey, *H160) {
	privateKey, address := MustGenRandomAddress(t)
	return privateKey, AddressToH160(address)
}
