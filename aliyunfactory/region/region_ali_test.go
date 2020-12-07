package main

import (
	connections "github.com/spotmaxtech/cloudconnections"
	"testing"
)

func TestRegion_getAvailableRegions(t *testing.T) {
	conn := *connections.NewAli("cn-hangzhou", "", "")
	r := Region{Conn: &conn}
	r.getAvailableRegions()
}
