package inputport

type InputPortServer interface {
	Run()
	Shutdown()
}
