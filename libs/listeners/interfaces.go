package listeners

type PortListener interface {
	Info() string
	Run(port int) error
	Stop() error
}
type BackgroundWorker interface {
	Start() error
	Stop() error
	Info() string
}
