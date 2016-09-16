package command

import (
	"github.com/obis/point"
)

type CommandLocation interface{}

func (*CommandLocation) Query() []*point.Point {

}
