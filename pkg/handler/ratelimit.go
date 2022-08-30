package handler

import (
	"fmt"

	"github.com/authgear/authgear-server/pkg/lib/ratelimit"
	"github.com/authgear/authgear-server/pkg/util/duration"
)

func AntiSpamContractMetadataRequestBucket() ratelimit.Bucket {
	return ratelimit.Bucket{
		Key:         fmt.Sprintf("contract-metadata-request"),
		Size:        10,
		ResetPeriod: duration.PerHour,
	}
}
