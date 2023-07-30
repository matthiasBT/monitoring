package web

import (
	"errors"
	"github.com/matthiasBT/monitoring/internal/storage"
	"strconv"
	"strings"
)

var (
	ErrInvalidMetricType = errors.New("invalid metric type")
	ErrMissingMetricName = errors.New("missing metric name")
	ErrInvalidMetricVal  = errors.New("invalid metric value")
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
		return ErrInvalidMetricType // treating the first token as the type; no type - return error
	} else if tokensCnt == 1 || tokensCnt == 2 && (*tokens)[1] == "" {
		return ErrMissingMetricName
	} else if tokensCnt == 2 {
		return ErrInvalidMetricVal
	} else if tokensCnt > 3 {
		return ErrInvalidMetricVal // treating the rest of the URL as the value
	}
	return nil
}

func validateMetricTypeVal(metricType string, val string) error {
	switch metricType {
	case storage.TypeGauge:
		if _, err := strconv.ParseFloat(val, 64); err != nil {
			return ErrInvalidMetricVal
		}
	case storage.TypeCounter:
		if _, err := strconv.ParseInt(val, 10, 64); err != nil {
			return ErrInvalidMetricVal
		}
	default:
		return ErrInvalidMetricType
	}
	return nil
}
