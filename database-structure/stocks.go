package database_structure

type Stock struct {
	ID    int
	Code  string
	Name  string
	Open  float64
	High  float64
	Low   float64
	Close float64
}
