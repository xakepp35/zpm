package zpm

import dto "github.com/prometheus/client_model/go"

type state struct {
	Dto  *dto.Metric
	Data any
}

type StateInitFunc = func(metricState *state)

func newState(timestampMs int64, labels ...*dto.LabelPair) *state {
	return &state{
		Dto: &dto.Metric{
			Label:       labels,
			TimestampMs: &timestampMs,
		},
	}
}
