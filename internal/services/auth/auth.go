package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spruceid/siwe-go"

	v1 "github.com/sambarnes/seaport-go/pkg/proto/seaport/v1"
)

const (
	AuthDomain = "opensea.io"
	AuthURI    = "https://opensea.io"
)

// Service holds state for the auth.Service.
type Service struct {
	// TODO: all in-memory state for now
	AuthState map[string]*common.Address
}

// Nonce returns an EIP-4361 nonce for session and invalidates existing session.
func (s *Service) Nonce(_ context.Context, _ *connect.Request[v1.Empty]) (
	*connect.Response[v1.NonceResponse], error,
) {
	// TODO: Double check cookie values here?
	nonce := siwe.GenerateNonce()
	res := connect.NewResponse(&v1.NonceResponse{Nonce: nonce})
	cookie := http.Cookie{
		Name:     "opensea-siwe",
		Value:    nonce,
		Domain:   AuthDomain,
		MaxAge:   60 * 60,
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
	}
	res.Header().Add("Set-Cookie", cookie.String())
	res.Header().Set("Seaport-Version", "v1")
	return res, nil
}

// Verify verifies the SignedMessage is valid and returns the verified address
func (s *Service) Verify(_ context.Context, req *connect.Request[v1.VerifyRequest]) (
	*connect.Response[v1.H160], error,
) {
	message, err := siwe.ParseMessage(req.Msg.Message)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// Actual signature verification, using present time validity
	// TODO: get current known nonce from some storage
	domain := AuthDomain
	publicKey, err := message.Verify(req.Msg.Signature, &domain, nil, nil)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	authenticatedAddress := v1.AddressToH160(crypto.PubkeyToAddress(*publicKey))
	res := connect.NewResponse(authenticatedAddress)
	res.Header().Set("Seaport-Version", "v1")
	return res, nil
}

// Authenticate checks if a given connection is authenticated, and returns
// the address which is authenticated for the nonce cookie.
func (s *Service) Authenticate(_ context.Context, _ *connect.Request[v1.Empty]) (
	*connect.Response[v1.H160], error,
) {
	// TODO: return the authenticated address (assume this is in headers?)
	authenticatedAddress := &v1.H160{
		Hi: &v1.H128{Hi: 17, Lo: 38},
		Lo: 420,
	}

	// Placeholder behavior for easy dev initially. Flip this to true to see the
	// unauthenticated response.
	flipToUnauthenticate := false
	if flipToUnauthenticate {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	res := connect.NewResponse(authenticatedAddress)
	res.Header().Set("Seaport-Version", "v1")
	return res, nil
}
