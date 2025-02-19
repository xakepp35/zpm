package zpm

import dto "github.com/prometheus/client_model/go"

type (
	Metric = dto.Metric
	LabelPair = dto.LabelPair
	LabelPairs = []*LabelPair
)

// NewLabelPairs param keyValues is an interleaved key-value-key-value... slice.
func NewLabelPairs(keyValues ...string) LabelPairs {
	res := make(LabelPairs, len(keyValues)/2)
	for i := range res {
		j := i*2
		res[i] = &LabelPair{
			Name: &keyValues[j],
			Value: &keyValues[j+1],
		}
	}
	return res
}

type state struct {
	Dto  *Metric
	Data any
}

type StateInitFunc = func(metricState *state)

func newState(timestampMs int64, labels ...*LabelPair) *state {
	return &state{
		Dto: &Metric{
			Label:       labels,
			TimestampMs: &timestampMs,
		},
	}
}
