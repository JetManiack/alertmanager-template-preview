package template

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/prometheus/common/helpers/templates"
)

// Shared template functions used in both Alertmanager and Prometheus modes.

func round(v any) (float64, error) {
	f, err := templates.ConvertToFloat(v)
	if err != nil {
		return 0, err
	}
	return math.Round(f), nil
}

func toTime(v any) (time.Time, error) {
	f, err := templates.ConvertToFloat(v)
	if err != nil {
		return time.Time{}, err
	}
	// Round to the nearest millisecond to avoid float64 precision artifacts,
	// as Prometheus timestamps are commonly millisecond-precision.
	return time.Unix(0, int64(math.Round(f*1000))*1e6).UTC(), nil
}

func toDuration(v any) (time.Duration, error) {
	f, err := templates.ConvertToFloat(v)
	if err != nil {
		return 0, err
	}
	return time.Duration(f * float64(time.Second)), nil
}

func toJson(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func humanizePercentage(v any) (string, error) {
	f, err := templates.ConvertToFloat(v)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.4g%%", f*100), nil
}

func humanize(v any) (string, error) {
	f, err := templates.ConvertToFloat(v)
	if err != nil {
		return "", err
	}
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return fmt.Sprintf("%.4g", f), nil
	}
	if math.Abs(f) < 1 {
		return fmt.Sprintf("%.4g", f), nil
	}

	units := []string{"", "k", "M", "G", "T", "P", "E"}
	i := 0
	for math.Abs(f) >= 1000 && i < len(units)-1 {
		f /= 1000
		i++
	}
	return strconv.FormatFloat(f, 'g', 4, 64) + units[i], nil
}

func humanize1024(v any) (string, error) {
	f, err := templates.ConvertToFloat(v)
	if err != nil {
		return "", err
	}
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return fmt.Sprintf("%.4g", f), nil
	}
	if math.Abs(f) < 1 {
		return fmt.Sprintf("%.4g", f), nil
	}

	units := []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei"}
	i := 0
	for math.Abs(f) >= 1024 && i < len(units)-1 {
		f /= 1024
		i++
	}
	return strconv.FormatFloat(f, 'g', 4, 64) + units[i], nil
}

func humanizeTimestamp(v any) (string, error) {
	return templates.HumanizeTimestamp(v)
}

func date(fmt string, t time.Time) string {
	return t.Format(fmt)
}

func tz(name string, t time.Time) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}

func list(args ...any) ([]any, error) {
	if args == nil {
		return []any{}, nil
	}
	return args, nil
}

func appendFunc(slice []any, args ...any) []any {
	return append(slice, args...)
}

func dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("dict requires an even number of arguments")
	}

	res := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		res[key] = values[i+1]
	}

	return res, nil
}
