package service

import (
	"github.com/authgear/authgear-server/pkg/api/apierrors"
)

var ErrBadNFTCollection = apierrors.Forbidden.WithReason("BadNFTCollection")
