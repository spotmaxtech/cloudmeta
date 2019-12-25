package main

import (
	. "github.com/smartystreets/goconvey/convey"
	connections "github.com/spotmaxtech/cloudconnections"
	"testing"
)

func Test_validInstance(t *testing.T) {
	Convey("test", t, func() {
		var instances []string
		var inst1, inst2, inst3, inst4, inst5 string
		inst1 = "g5se.xlarge"
		inst2 = "sn2ne.large"
		inst3 = "c6.large"
		inst4 = "re4e.40xlarge"
		inst5 = "d1ne.2xlarge"
		instances = append(instances, inst1, inst2, inst3, inst4, inst5)
		for i, v := range instances {
			if i !=4 {
				So(validInstance(v), ShouldBeTrue)
			} else {
				So(validInstance(v), ShouldBeFalse)
			}
		}
	})
}

func TestInstUtil_FetchInstance(t *testing.T) {
	conn := *connections.NewAli("cn-hangzhou","","")
	util := InstUtil{Conn:&conn}
	util.FetchInstance("cn-hongkong","cn-hongkong-b")
}