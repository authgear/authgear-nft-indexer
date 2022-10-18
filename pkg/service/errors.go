package service

import (
	"github.com/authgear/authgear-server/pkg/api/apierrors"
)

var ErrAlchemyError = apierrors.ServiceUnavailable.WithReason("ServiceUnavailable")

var ErrBadNFTCollection = apierrors.Forbidden.WithReason("BadNFTCollection")
