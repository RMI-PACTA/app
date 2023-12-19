package pactasrv

import (
	"fmt"
	"math"

	"github.com/RMI/pacta/oapierr"
	"go.uber.org/zap"
)

const errorID = "UGC_INPUT_TOO_LARGE"

func checkStringLimitMediumPtr(site string, value *string) error {
	if value == nil {
		return nil
	}
	return checkStringLimitMedium(site, *value)
}

func checkStringLimitMedium(site string, value string) error {
	return checkStringLimit(site, value, 10000)
}

func checkStringLimitSmallPtr(site string, value *string) error {
	if value == nil {
		return nil
	}
	return checkStringLimitSmall(site, *value)
}

func checkStringLimitSmall(site string, value string) error {
	return checkStringLimit(site, value, 1000)
}

func checkStringLimit(site string, value string, byteLimit int) error {
	byteLength := len(value)
	return checkLimit(
		site,
		byteLength, formatByteSize(byteLength),
		byteLimit, formatByteSize(byteLimit))
}

func checkIntLimit(site string, value int, limit int) error {
	return checkLimit(site, value, fmt.Sprintf("%d", value), limit, fmt.Sprintf("%d", limit))
}

func checkLimit(site string, value int, valueStr string, limit int, limitStr string) error {
	if value > limit {
		return oapierr.BadRequest(
			"Input too large",
			zap.String("site", site),
			zap.Int("value", value),
			zap.Int("limit", limit),
		).WithMessage(
			fmt.Sprintf("the input for %s is too large (%s exceeds the limit, %s)",
				site,
				valueStr,
				limitStr),
		).WithErrorID(errorID)
	}
	return nil
}

func formatByteSize(bytes int) string {
	if bytes <= 0 {
		return ""
	}
	k := 1000.0
	i := int(math.Floor(math.Log(float64(bytes)) / math.Log(k)))
	return fmt.Sprintf("%.2f %s", float64(bytes)/math.Pow(k, float64(i)), dataSizes[i])
}

var dataSizes []string = []string{
	"Bytes", "kB", "MB", "GB", "TB",
}

func anyError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
