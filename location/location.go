package location

var GeohashPrecision int = 6

type QueryObject struct {
	lat    float64
	lng    float64
	radius uint32
	limit  uint32
	role   uint8
	order  string
}
