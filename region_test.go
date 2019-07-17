package cloudmeta

import (
	"github.com/spotmaxtech/gokit"
	"testing"
)

func TestNewAWSRegion(t *testing.T) {
	aws := NewAWSRegion()
	t.Log(gokit.Prettify(aws.data))
}
