package web

import (
	"errors"
	"github.com/matthiasBT/monitoring/internal/storage"
	"strconv"
	"strings"
)

var (
	InvalidMetricType = errors.New("invalid metric type")
	MissingMetricName = errors.New("missing metric name")
	InvalidMetricVal  = errors.New("invalid metric value")
)

func ParseMetricUpdate(url string, prefix string) (*storage.MetricUpdate, error) {
	url = strings.TrimPrefix(url, prefix)
	tokens := strings.Split(url, "/")
	if err := validateMetricUpdateURL(&tokens); err != nil {
		return nil, err
	}
	var metricType, name, val = tokens[0], tokens[1], tokens[2]
	if err := validateMetricTypeVal(metricType, val); err != nil {
		return nil, err
	}
	res := storage.MetricUpdate{Type: metricType, Name: name, Value: val}
	return &res, nil
}

func validateMetricUpdateURL(tokens *[]string) error {
	tokensCnt := len(*tokens)
	if tokensCnt == 0 || tokensCnt == 1 && (*tokens)[0] == "" {
		return InvalidMetricType // treating the first token as the type; no type - return error
	} else if tokensCnt == 1 || tokensCnt == 2 && (*tokens)[1] == "" {
		return MissingMetricName
	} else if tokensCnt == 2 {
		return InvalidMetricVal
	} else if tokensCnt > 3 {
		return InvalidMetricVal // treating the rest of the URL as the value
	}
	return nil
}

func validateMetricTypeVal(metricType string, val string) error {
	switch metricType {
	case storage.TypeGauge:
		if _, err := strconv.ParseFloat(val, 64); err != nil {
			return InvalidMetricVal
		}
	case storage.TypeCounter:
		if _, err := strconv.ParseInt(val, 10, 64); err != nil {
			return InvalidMetricVal
		}
	default:
		return InvalidMetricType
	}
	return nil
}
