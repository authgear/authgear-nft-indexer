package web3

import (
	"github.com/authgear/authgear-server/pkg/api/apierrors"
)

var ErrAlchemyProtocol = apierrors.InternalError.WithReason("AlchemyProtocol")
