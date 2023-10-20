package auth_test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/matryer/is"
	"github.com/spruceid/siwe-go"

	"github.com/sambarnes/seaport-go/internal/services/auth"
	v1 "github.com/sambarnes/seaport-go/pkg/proto/seaport/v1"
	"github.com/sambarnes/seaport-go/pkg/proto/seaport/v1/seaportv1connect"
)

func TestAuth_Nonce(t *testing.T) {
	is := is.New(t)
	mux := http.NewServeMux()
	mux.Handle(seaportv1connect.NewAuthServiceHandler(&auth.Service{}))
	server := httptest.NewServer(mux)
	defer server.Close()

	client := seaportv1connect.NewAuthServiceClient(server.Client(), server.URL)
	resp, err := client.Nonce(context.Background(), connect.NewRequest(&v1.Empty{}))
	is.NoErr(err) // Failed Nonce call

	got := resp.Msg.Nonce
	is.True(len(got) >= 8)
	is.True(regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(got))
	is.True(resp.Header().Get("Set-Cookie") != "")
}

// Should be able to verify a correctly signed SIWE signed message
func TestAuth_Verify_Success(t *testing.T) {
	is, ctx := is.New(t), context.Background()
	mux := http.NewServeMux()
	mux.Handle(seaportv1connect.NewAuthServiceHandler(&auth.Service{}))
	server := httptest.NewServer(mux)
	defer server.Close()

	client := seaportv1connect.NewAuthServiceClient(server.Client(), server.URL)
	nonceResp, err := client.Nonce(ctx, connect.NewRequest(&v1.Empty{}))
	is.NoErr(err) // Failed Nonce call

	// Challenge setup
	nonce := nonceResp.Msg.Nonce
	privateKey, wantAddress := v1.MustGenRandomAddress(t)
	msg, err := siwe.InitMessage(
		auth.AuthDomain,
		wantAddress.String(),
		auth.AuthURI,
		nonce,
		nil,
	)
	is.NoErr(err) // Failed to initialize SIWE message
	signature, err := signMessage(msg.String(), privateKey)
	is.NoErr(err)                // Failed to sign message
	is.Equal(len(signature), 65) // Invalid signature length

	// Verify the challenge against the AuthService
	verifyResp, err := client.Verify(ctx, connect.NewRequest(&v1.VerifyRequest{
		Message:   msg.String(),
		Signature: fmt.Sprintf("0x%x", signature),
	}))
	is.NoErr(err) // Failed Verify call

	gotAddress := verifyResp.Msg.ToAddress()
	is.Equal(gotAddress, wantAddress)
}

func signHash(data []byte) common.Hash {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256Hash([]byte(msg))
}

func signMessage(message string, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	sign := signHash([]byte(message))
	signature, err := crypto.Sign(sign.Bytes(), privateKey)

	if err != nil {
		return nil, err
	}

	signature[64] += 27
	return signature, nil
}
