// Package platform provides small helper utilities used across the project.
package platform

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func AssertLinux() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("not supported operating system")
	}
	return nil
}

type tickMsg time.Time

type model struct {
	ip       string
	port     int
	interval time.Duration
}

func (m model) Init() tea.Cmd {
	return tick(m.interval)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tickMsg:
		return m, tick(m.interval)
	}
	return m, nil
}

func (m model) View() string {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		MarginBottom(1)

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginTop(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("white"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1)

	row := func(label, value string) string {
		return fmt.Sprintf("%s %s",
			labelStyle.Width(20).Render(label+":"),
			valueStyle.Render(value))
	}

	var s string

	s += titleStyle.Render("Runtime Diagnostics") + "\n"

	s += sectionStyle.Render("System") + "\n"
	s += row("Go Version", runtime.Version()) + "\n"
	s += row("OS", runtime.GOOS) + "\n"
	s += row("Architecture", runtime.GOARCH) + "\n"
	s += row("Compiler", runtime.Compiler) + "\n"
	s += row("CPU Cores", fmt.Sprintf("%d", runtime.NumCPU())) + "\n"
	s += row("GOMAXPROCS", fmt.Sprintf("%d", runtime.GOMAXPROCS(0))) + "\n"
	s += row("Goroutines", fmt.Sprintf("%d", runtime.NumGoroutine())) + "\n"

	s += sectionStyle.Render("Memory") + "\n"
	s += row("Heap Alloc", fmt.Sprintf("%d KB", mem.HeapAlloc/1024)) + "\n"
	s += row("Heap Sys", fmt.Sprintf("%d KB", mem.HeapSys/1024)) + "\n"
	s += row("Heap In Use", fmt.Sprintf("%d KB", mem.HeapInuse/1024)) + "\n"
	s += row("Heap Idle", fmt.Sprintf("%d KB", mem.HeapIdle/1024)) + "\n"
	s += row("Stack In Use", fmt.Sprintf("%d KB", mem.StackInuse/1024)) + "\n"
	s += row("Stack Sys", fmt.Sprintf("%d KB", mem.StackSys/1024)) + "\n"
	s += row("GC Cycles", fmt.Sprintf("%d", mem.NumGC)) + "\n"
	s += row("GC Pause Total", fmt.Sprintf("%d ms", mem.PauseTotalNs/1e6)) + "\n"
	s += row("GC Next", fmt.Sprintf("%d KB", mem.NextGC/1024)) + "\n"

	s += sectionStyle.Render("Server Info") + "\n"
	s += row("IP Address", m.ip) + "\n"
	s += row("Port", fmt.Sprintf("%d", m.port)) + "\n"

	if info, ok := debug.ReadBuildInfo(); ok {
		s += sectionStyle.Render("Build Info") + "\n"
		s += row("Module", info.Path) + "\n"
		s += row("Go Version", info.GoVersion) + "\n"
		if info.Main.Version != "(devel)" {
			s += row("Version", info.Main.Version) + "\n"
		}
	}

	s += helpStyle.Render("\nPress 'q' or 'ctrl+c' to quit")

	return s
}

func tick(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func RunRuntimeDiagnostics(interval time.Duration, ip string, port int) {
	p := tea.NewProgram(model{
		ip:       ip,
		port:     port,
		interval: interval,
	})

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running diagnostics: %v\n", err)
	}
}
