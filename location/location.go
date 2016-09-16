package location



type QueryObject struct{
	lat float64
	lng float64
	radius uint32
	limit uint32
	role uint8 nil
	order string const("update", "distance")
	
	
}