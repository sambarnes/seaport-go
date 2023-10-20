package seaportv1_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/matryer/is"

	seaportv1 "github.com/sambarnes/seaport-go/pkg/proto/seaport/v1"
)

// Should be able to convert from Address to H160 and back
func TestH160_ToAddress(t *testing.T) {
	is := is.New(t)

	_, wantAddress := seaportv1.MustGenRandomAddress(t)
	gotH160 := seaportv1.AddressToH160(wantAddress)
	gotAddress := gotH160.ToAddress()
	is.Equal(gotAddress, wantAddress)
}

// Documenting expected message size comparisons with a test. Sizes
// should be strictly smaller than their naive string representation.
func TestH160_Sizes(t *testing.T) {
	is := is.New(t)

	address := common.HexToAddress("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	h160 := seaportv1.AddressToH160(address)
	data, err := seaportv1.MarshalProto(h160)
	is.NoErr(err)
	is.Equal(len(data), 30) // Largest possible address should be 30 bytes
	naiveSize := len([]byte(address.String()))
	is.Equal(naiveSize, 42)

	address = common.HexToAddress("0x00000000000000ADc04C56Bf30aC9d3c0aAF14dC")
	h160 = seaportv1.AddressToH160(address)
	data, err = seaportv1.MarshalProto(h160)
	is.NoErr(err)
	is.Equal(len(data), 21) // Seaport 1.5 address should be 21 bytes
	naiveSize = len([]byte(address.String()))
	is.Equal(naiveSize, 42)

	h160 = &seaportv1.H160{}
	data, err = seaportv1.MarshalProto(h160)
	is.NoErr(err)
	is.Equal(len(data), 0) // Null address should be zero bytes
	naiveSize = len([]byte(address.String()))
	is.Equal(naiveSize, 42)
}
