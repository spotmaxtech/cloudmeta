package main

import (
	"fmt"
	"testing"
)

//func Test_getALiInstanceFamily(t *testing.T) {
//	v := getALiInstanceFamily("Compute")
//	fmt.Print(v)
//}

func Test_validInstance(t *testing.T) {
	var s = "ecs.gn6i-c8g1.2xlarge"
	fmt.Print(validInstance(s, "gn6i"))

	fmt.Print(validInstance(s, "scc"))
}
