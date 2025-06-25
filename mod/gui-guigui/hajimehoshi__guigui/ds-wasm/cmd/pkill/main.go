package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

func main() {
	var (
		byName = flag.String("name", "", "Kill processes by name (partial match)")
		byPort = flag.Int("port", 0, "Kill processes using specific port")
		byPID  = flag.Int("pid", 0, "Kill specific process by PID")
		force  = flag.Bool("force", false, "Force kill (SIGKILL)")
		list   = flag.Bool("list", false, "List matching processes without killing")
		help   = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *byName == "" && *byPort == 0 && *byPID == 0 {
		fmt.Println("❌ Must specify -name, -port, or -pid")
		showHelp()
		os.Exit(1)
	}

	processes, err := process.Processes()
	if err != nil {
		log.Fatalf("❌ Failed to get processes: %v", err)
	}

	var targets []*process.Process

	// Find processes by criteria
	for _, p := range processes {
		if *byPID != 0 {
			if int(p.Pid) == *byPID {
				targets = append(targets, p)
			}
			continue
		}

		if *byName != "" {
			name, err := p.Name()
			if err != nil {
				continue
			}
			if strings.Contains(strings.ToLower(name), strings.ToLower(*byName)) {
				targets = append(targets, p)
			}
		}

		if *byPort != 0 {
			connections, err := p.Connections()
			if err != nil {
				continue
			}
			for _, conn := range connections {
				if int(conn.Laddr.Port) == *byPort {
					targets = append(targets, p)
					break
				}
			}
		}
	}

	if len(targets) == 0 {
		fmt.Println("🔍 No matching processes found")
		return
	}

	// List or kill processes
	for _, p := range targets {
		name, _ := p.Name()
		cmdline, _ := p.Cmdline()

		if *list {
			fmt.Printf("📋 PID: %d, Name: %s, Cmd: %s\n", p.Pid, name, cmdline)
			continue
		}

		fmt.Printf("🔪 Killing PID: %d, Name: %s\n", p.Pid, name)

		var err error
		if *force {
			err = p.Kill()
		} else {
			err = p.Terminate()
		}

		if err != nil {
			fmt.Printf("❌ Failed to kill PID %d: %v\n", p.Pid, err)
		} else {
			fmt.Printf("✅ Killed PID %d\n", p.Pid)
		}
	}
}

func showHelp() {
	fmt.Println("Cross-Platform Process Killer using gopsutil v4")
	fmt.Println("===============================================")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  pkill -name <process_name>    # Kill by process name (partial match)")
	fmt.Println("  pkill -port <port_number>     # Kill processes using specific port")
	fmt.Println("  pkill -pid <process_id>       # Kill specific process by PID")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -force                        # Force kill (SIGKILL instead of SIGTERM)")
	fmt.Println("  -list                         # List matching processes without killing")
	fmt.Println("  -help                         # Show this help")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  pkill -name caddy             # Kill all processes with 'caddy' in name")
	fmt.Println("  pkill -port 8081              # Kill process using port 8081")
	fmt.Println("  pkill -name ds-server -list   # List ds-server processes")
	fmt.Println("  pkill -pid 1234 -force        # Force kill PID 1234")
}
