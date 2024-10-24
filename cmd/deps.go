package main

import "os"

type app interface {
	Run(stopChan chan os.Signal)
}
