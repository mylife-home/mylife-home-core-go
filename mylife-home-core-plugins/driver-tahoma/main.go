//go:generate go run ../../mylife-home-core-generator/cmd/main.go .

// @Module(version="1.0.3")
package plugin_entry

import (
	_ "mylife-home-core-plugins-driver-tahoma/plugin"
)
