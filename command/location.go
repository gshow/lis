package command

import (
	"github.com/gshow/obis/point"
)

type CommandLocation interface{}

func (*CommandLocation) Query() []*point.Point {

}
