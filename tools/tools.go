// Package tools provides small helper utilities used across the project.
package tools

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

func AssertLinux() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("not supported operating system")
	}
	return nil
}

func RunRuntimeDiagnostics(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {

		fmt.Print("\033[2J")
		fmt.Print("\033[H")

		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		fmt.Println("runtime diagnostics")
		fmt.Println("-------------------")

		fmt.Printf("%-20s : %s\n", "go version", runtime.Version())
		fmt.Printf("%-20s : %s\n", "os", runtime.GOOS)
		fmt.Printf("%-20s : %s\n", "architecture", runtime.GOARCH)
		fmt.Printf("%-20s : %s\n", "compiler", runtime.Compiler)

		fmt.Printf("%-20s : %d\n", "cpu cores", runtime.NumCPU())
		fmt.Printf("%-20s : %d\n", "gomaxprocs", runtime.GOMAXPROCS(0))
		fmt.Printf("%-20s : %d\n", "goroutines", runtime.NumGoroutine())

		fmt.Println()
		fmt.Println("memory")
		fmt.Println("-----------------")

		fmt.Printf("%-20s : %d KB\n", "heap alloc", m.HeapAlloc/1024)
		fmt.Printf("%-20s : %d KB\n", "heap sys", m.HeapSys/1024)
		fmt.Printf("%-20s : %d KB\n", "heap in use", m.HeapInuse/1024)
		fmt.Printf("%-20s : %d KB\n", "heap idle", m.HeapIdle/1024)
		fmt.Printf("%-20s : %d KB\n", "stack in use", m.StackInuse/1024)
		fmt.Printf("%-20s : %d KB\n", "stack sys", m.StackSys/1024)

		fmt.Printf("%-20s : %d\n", "gc cycles", m.NumGC)
		fmt.Printf("%-20s : %d ms\n", "gc pause total", m.PauseTotalNs/1e6)

		if info, ok := debug.ReadBuildInfo(); ok {
			fmt.Println()
			fmt.Println("information")
			fmt.Println("-----------------")

			fmt.Printf("%-20s : %s\n", "module", info.Path)
			fmt.Printf("%-20s : %s\n", "go version", info.GoVersion)

			if info.Main.Version != "(devel)" {
				fmt.Printf("%-20s : %s\n", "version", info.Main.Version)
			}
		}
	}
}
