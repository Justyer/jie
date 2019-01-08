package jie

type IProtocol interface {
	IsWhole([]byte) bool
	Parse([]byte)
	Get([]byte) (int, error)
	New() IProtocol
	Data() []byte
}
