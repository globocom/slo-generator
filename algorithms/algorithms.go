package algoritms

import (
	"fmt"

	"github.com/prometheus/prometheus/pkg/labels"
)

func MultiBurnRateForPage(metric string, lbs labels.Labels, operator string, value float64) string {
	result := ""
	result += fmt.Sprintf(`(%s:ratio_rate_1h%s %s (14.4 * %.3f) and `, metric, lbs.String(), operator, value)
	result += fmt.Sprintf(`%s:ratio_rate_5m%s %s (14.4 * %.3f))`, metric, lbs.String(), operator, value)

	result += " or "

	result += fmt.Sprintf(`(%s:ratio_rate_6h%s %s (6 * %.3f) and `, metric, lbs.String(), operator, value)
	result += fmt.Sprintf(`%s:ratio_rate_30m%s %s (6 * %.3f))`, metric, lbs.String(), operator, value)

	return result
}

func MultiBurnRateForTicket(metric string, lbs labels.Labels, operator string, value float64) string {
	result := ""
	result += fmt.Sprintf(`(%s:ratio_rate_1d%s %s (3 * %.3f) and `, metric, lbs.String(), operator, value)
	result += fmt.Sprintf(`%s:ratio_rate_2h%s %s (3 * %.3f))`, metric, lbs.String(), operator, value)

	result += " or "

	result += fmt.Sprintf(`(%s:ratio_rate_3d%s %s %.3f and `, metric, lbs.String(), operator, value)
	result += fmt.Sprintf(`%s:ratio_rate_6h%s %s %.3f)`, metric, lbs.String(), operator, value)

	return result
}
