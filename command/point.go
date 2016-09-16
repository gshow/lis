package command

import (
	"github.com/gshow/obis/point"
)

type CommandPoint interface{}

func (*CommandPoint) Set(pt *point) bool {

}

func (*CommandPoint) Delete(pt *point) bool {

}
func (*CommandPoint) Query() bool {

}
