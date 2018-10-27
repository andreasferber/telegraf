package concat

import (
	"testing"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/influxdata/telegraf/plugins/processors"
	"github.com/stretchr/testify/assert"
)

func newM1() telegraf.Metric {
	m1, _ := metric.New("concat_test",
		map[string]string{
			"tag_a": "value_a",
			"tag_b": "value_b",
		},
		map[string]interface{}{
			"field": "value",
		},
		time.Now(),
	)
	return m1
}

func newM2() telegraf.Metric {
	m1, _ := metric.New("concat_test",
		map[string]string{
			"tag_a": "value_a",
			"tag_b": "value_b",
			"tag_c": "value_c",
			"tag_d": "value_d",
		},
		map[string]interface{}{
			"field": "value",
		},
		time.Now(),
	)
	return m1
}

func TestRegistration(t *testing.T) {
	assert.Contains(t, processors.Processors, "concat")
	assert.IsType(t, &Concat{}, processors.Processors["concat"]())
}

func TestProcessorInterface(t *testing.T) {
	concat := NewConcat()
	assert.IsType(t, "", concat.SampleConfig())
	assert.IsType(t, "", concat.Description())
}

func TestTagConcatenations(t *testing.T) {
	tests := []struct {
		message      string
		converter    converter
		expectedTags map[string]string
	}{
		{
			message: "Should add new tag",
			converter: converter{
				Keys:      []string{"tag_a", "tag_b"},
				ResultKey: "new_tag",
			},
			expectedTags: map[string]string{
				"tag_a":   "value_a",
				"tag_b":   "value_b",
				"new_tag": "value_avalue_b",
			},
		},
		{
			message: "Should use separator",
			converter: converter{
				Keys:      []string{"tag_a", "tag_b"},
				Separator: "-",
				ResultKey: "new_tag",
			},
			expectedTags: map[string]string{
				"tag_a":   "value_a",
				"tag_b":   "value_b",
				"new_tag": "value_a-value_b",
			},
		},
		{
			message: "Should skip if input tag missing",
			converter: converter{
				Keys:      []string{"tag_a", "tag_b", "tag_c"},
				ResultKey: "new_tag",
			},
			expectedTags: map[string]string{
				"tag_a": "value_a",
				"tag_b": "value_b",
			},
		},
	}

	for _, test := range tests {
		concat := NewConcat()
		concat.Tags = []converter{
			test.converter,
		}

		processed := concat.Apply(newM1())

		expectedFields := map[string]interface{}{
			"field": "value",
		}

		assert.Equal(t, expectedFields, processed[0].Fields(), test.message, "Should not change fields")
		assert.Equal(t, test.expectedTags, processed[0].Tags(), test.message)
		assert.Equal(t, "concat_test", processed[0].Name(), "Should not change name")
	}
}

func TestMultipleConcatenations(t *testing.T) {
	concat := NewConcat()
	concat.Tags = []converter{
		{
			Keys:      []string{"tag_a", "tag_b"},
			Separator: "-",
			ResultKey: "result_1",
		},
		{
			Keys:      []string{"tag_a", "tag_c"},
			Separator: "+",
			ResultKey: "result_2",
		},
		{
			Keys:      []string{"tag_c", "tag_d"},
			Separator: "@",
			ResultKey: "result_3",
		},
	}

	processed := concat.Apply(newM2())

	expectedTags := map[string]string{
		"tag_a":    "value_a",
		"tag_b":    "value_b",
		"tag_c":    "value_c",
		"tag_d":    "value_d",
		"result_1": "value_a-value_b",
		"result_2": "value_a+value_c",
		"result_3": "value_c@value_d",
	}

	assert.Equal(t, expectedTags, processed[0].Tags())
}

func TestMultipleMetrics(t *testing.T) {
	concat := NewConcat()
	concat.Tags = []converter{
		{
			Keys:      []string{"tag_a", "tag_b"},
			Separator: "-",
			ResultKey: "result",
		},
	}

	processed := concat.Apply(newM1(), newM1())

	expectedTags := map[string]string{
		"tag_a":  "value_a",
		"tag_b":  "value_b",
		"result": "value_a-value_b",
	}

	assert.Len(t, processed, 2)
	assert.Equal(t, expectedTags, processed[0].Tags())
	assert.Equal(t, expectedTags, processed[1].Tags())
}

func BenchmarkConcatenations(b *testing.B) {
	concat := NewConcat()
	concat.Tags = []converter{
		{
			Keys:      []string{"tag_a", "tag_b"},
			Separator: "-",
			ResultKey: "result",
		},
	}

	for n := 0; n < b.N; n++ {
		processed := concat.Apply(newM1())
		_ = processed
	}
}
