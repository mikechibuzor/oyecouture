package payments

import "context"

type Provider interface {
	Name() string
	Initialize(context.Context, InitializeRequest) (InitializeResult, []byte, error)
	Verify(context.Context, string) (Verification, error)
}
