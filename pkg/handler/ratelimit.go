package handler

import (
	"fmt"

	"github.com/authgear/authgear-server/pkg/lib/ratelimit"
	"github.com/authgear/authgear-server/pkg/util/duration"
)

func AntiSpamContractMetadataRequestBucket(appID string, ip string) ratelimit.Bucket {
	return ratelimit.Bucket{
		Key:         fmt.Sprintf("request:%s-%s", appID, ip),
		Size:        10,
		ResetPeriod: duration.PerHour,
	}
}
