package handler

import (
	"fmt"

	"github.com/authgear/authgear-server/pkg/lib/ratelimit"
	"github.com/authgear/authgear-server/pkg/util/duration"
)

func AntiSpamContractMetadataRequestBucket(appID string) ratelimit.Bucket {
	return ratelimit.Bucket{
		Key:         fmt.Sprintf("contract-metadata-request-%s", appID),
		Size:        10,
		ResetPeriod: duration.PerHour,
	}
}

func AntiSpamProbeCollectionRequestBucket(appID string) ratelimit.Bucket {
	return ratelimit.Bucket{
		Key:         fmt.Sprintf("probe-collection-request-%s", appID),
		Size:        60,
		ResetPeriod: duration.PerHour,
	}
}
