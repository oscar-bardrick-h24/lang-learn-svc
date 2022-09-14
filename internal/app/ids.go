package app

type IDTool interface {
	New() (string, error)
}
