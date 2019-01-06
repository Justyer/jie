package jie

type Protocol interface {
	IsWhole([]byte) bool
	Parse([]byte)
	Get([]byte) (int, error)
}
