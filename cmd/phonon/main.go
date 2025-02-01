package phonon

import (
	"phonon/pkg/config"
	"phonon/pkg/instrumentation"
)

func main() {
	config.Initialize()
	instrumentation.InitializeLogging()
}
