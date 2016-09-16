package location

type QueryObject struct {
	lat    float64
	lng    float64
	radius uint32
	limit  uint32
	role   uint8
	order  string
}
