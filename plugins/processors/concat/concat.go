package concat

import (
	"strings"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/processors"
)

// Concat contains the main plugin configuration
type Concat struct {
	Tags []converter `toml:"tags"`
}

// converter describes a particular tag concatenation
type converter struct {
	Keys      []string `toml:"keys"`
	Separator string   `toml:"separator"`
	ResultKey string   `toml:"result_key"`
}

const sampleConfig = `
  ## Tag concatenations defined in a separate sub-table
  # [[processors.concat.tags]]
  #   ## Tags that should be concatenated
  #   keys = [ "host", "if_name" ]
  #   ## Separator used to join the values together
  #   separator = "-"
  #   ## New tag key
  #   result_key = "global_if_name"
`

// NewConcat creates a new Concat object
func NewConcat() *Concat {
	return &Concat{}
}

// SampleConfig returns the sample configuration of the plugin
func (r *Concat) SampleConfig() string {
	return sampleConfig
}

// Description returns a descriptive string for the plugin
func (r *Concat) Description() string {
	return "Concatenate tag values into a new tag"
}

// Apply applies the configured concatenations to the tags associated with metrics
func (r *Concat) Apply(in ...telegraf.Metric) []telegraf.Metric {
	for _, metric := range in {
	ConverterLoop:
		for _, converter := range r.Tags {
			values := make([]string, len(converter.Keys))
			for idx, key := range converter.Keys {
				value, ok := metric.GetTag(key)
				if !ok {
					continue ConverterLoop
				}
				values[idx] = value
			}
			newValue := strings.Join(values, converter.Separator)
			metric.AddTag(converter.ResultKey, newValue)
		}
	}

	return in
}

func init() {
	processors.Add("concat", func() telegraf.Processor {
		return NewConcat()
	})
}
