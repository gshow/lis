package point


type Point struct {
	id     uint64
	lat    float64
	lng    float64
	role   uint8
	update uint32
	expire uint32
	ext    uint64
}

type Points struct {
	pt     *point
	expire uint32
}

var PointsCollector = []*points


type QueryObject struct{
	
	role uint8
	id uint64
	
}

