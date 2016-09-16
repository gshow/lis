package main

import (
	"github.com/gshow/obis/location"
	"github.com/gshow/obis/point"
	"github.com/gshow/obis/command"
)


query := location.QueryObject{}

ret := command.CommandLocation.Query(query)



query2 := point.QueryObject{}
command.CommandPoint.Query(query2)

