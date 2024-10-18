package subcommands

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sort"

	"os"
	"runtime"

	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/cpu"
	dsk "github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"github.com/spf13/cobra"
)

var (
	help     bool
	all      bool
	cpuInfo  bool
	cpuUsage bool
	disk     bool
	localIP  bool
	publicIP bool
	ram      bool
	top5RAM  bool
)

var CmdCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Run the analysis on a given repository",
	Long:  "The 'cmd' command allows you to analyze public or private Git repositories for various statistics.",
	// Custom help flag handling in the Run function
	Run: func(cmd *cobra.Command, args []string) {
		// Custom help logic
		if cmd.Flags().Lookup("help").Changed {
			printCmdHelpTable() // Print help table
			os.Exit(0)
		}
		// Other flag handling logic
		if all {
			fmt.Println("Showing all stats for the repository...")
			getCPUInfo()	
			getCPUUsage()
			getDiskUsage()
			getLocalIP()
			getPublicIP()
			getCurrentRamUsage()
			getTop5RamConsumption()
			

		} else {
			if cpuInfo {
				getCPUInfo()

			}
			if cpuUsage {
				getCPUUsage()

			}
			if disk {
				getDiskUsage()

			}
			if localIP {
				getLocalIP()

			}
			if publicIP {
				getPublicIP()

			}
			if ram {
				getCurrentRamUsage()

			}
			if top5RAM {
				getTop5RamConsumption()

			}
		}
	},
}

func init() {
	// Disable Cobra's default help flag
	CmdCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		printCmdHelpTable() // Custom help function in table format
	})
	CmdCmd.Flags().BoolVarP(&all, "all", "a", false, "Show all stats for the repository")
	CmdCmd.Flags().BoolVarP(&cpuInfo, "cpu-info", "i", false, "Show CPU information")
	CmdCmd.Flags().BoolVarP(&cpuUsage, "cpu-usage", "u", false, "Show CPU usage")
	CmdCmd.Flags().BoolVarP(&disk, "disk", "d", false, "Show disk usage")
	CmdCmd.Flags().BoolVarP(&localIP, "local-ip", "l", false, "Show local IP address")
	CmdCmd.Flags().BoolVarP(&publicIP, "public-ip", "p", false, "Show public IP address")
	CmdCmd.Flags().BoolVarP(&ram, "ram", "r", false, "Show RAM usage")
	CmdCmd.Flags().BoolVarP(&top5RAM, "top5-ram", "t", false, "Show top 5 processes consuming the most RAM")

	// Add examples to the command
	CmdCmd.Example = `  # Analyze all statistics
  ./vc-analyze cmd --all

  # Analyze CPU information
  ./vc-analyze cmd --cpu-info

  # Analyze disk usage
  ./vc-analyze cmd --disk

  # Analyze local IP address
  ./vc-analyze cmd --local-ip

  # Analyze RAM usage
  ./vc-analyze cmd --ram`
}

// Function to print command help in table format
func printCmdHelpTable() {
	// Create a new table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetCaption(true, "Available flags for the 'cmd' command")
	table.SetHeader([]string{"Flag", "short-hand", "Description"})
	//Set minimum column width for 2nd column
	table.SetColMinWidth(2, 35)
	//// Disable automatic wrapping of text
	table.SetAutoWrapText(false)

	// Add data to the table
	table.Append([]string{"--all", "-a", "Show all stats for the repository"})
	table.Append([]string{"--cpu-info", "-i", "Show CPU information"})
	table.Append([]string{"--cpu-usage", "-u", "Show CPU usage"})
	table.Append([]string{"--disk", "-d", "Show disk usage"})
	table.Append([]string{"--local-ip", "-l", "Show local IP address"})
	table.Append([]string{"--public-ip", "-p", "Show public IP address"})
	table.Append([]string{"--ram", "-r", "Show RAM usage"})
	table.Append([]string{"--top5-ram", "-t", "Show top 5 processes consuming the most RAM"})

	// Setting border and colors for each column
	table.SetBorder(true)
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor}, // Color for "Flag" header
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor}, // Color for "Short-hand" header
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor}, //  Color for Description header
	)
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgHiBlackColor}, // Color for "Flag" column
		tablewriter.Colors{tablewriter.FgHiBlackColor}, // Color for "Description" column
		tablewriter.Colors{tablewriter.FgHiBlackColor}, // Color for Short-hand column
	)

	// Render the table to the console
	table.Render()
}

// CPU Information
func getCPUInfo() {
	fmt.Println("Fetching CPU information...")
	switch runtime.GOOS {
	case "darwin", "windows":
		// Fallback on macOS/Windows
		fmt.Printf("CPU Model: %s\n", runtime.GOARCH)
		cores, err := cpu.Counts(true)
		if err != nil {
			log.Printf("Error fetching CPU core count: %v", err)
			return
		}
		fmt.Printf("Number of Cores: %d\n", cores)
	case "linux":
		info, err := cpu.Info()
		if err != nil {
			log.Printf("Error fetching CPU info: %v", err)
			return
		}
		for _, cpu := range info {
			fmt.Printf("CPU: %s, Cores: %d\n", cpu.ModelName, cpu.Cores)
		}
	default:
		fmt.Println("Unsupported platform")
	}
	seperator()
}

// CPU Usage
func getCPUUsage() {
	fmt.Println("Fetching CPU usage...")
	usage, err := cpu.Percent(0, false)
	if err != nil {
		log.Printf("Error fetching CPU usage: %v", err)
		return
	}
	fmt.Printf("CPU Usage: %.2f%%\n", usage[0])
	seperator()
}

// Disk Usage
func getDiskUsage() {
	fmt.Println("Fetching disk usage...")
	usage, err := dsk.Usage("/")
	if err != nil {
		log.Printf("Error fetching disk usage: %v", err)
		return
	}
	fmt.Printf("Total: %v GB, Used: %v GB, Free: %v GB, Usage: %.2f%%\n",
		usage.Total/1024/1024/1024, usage.Used/1024/1024/1024, usage.Free/1024/1024/1024, usage.UsedPercent)
	seperator()
}

// Local IP
func getLocalIP() {
	fmt.Println("Fetching local IP address...")
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Printf("Error fetching local IP: %v", err)
		return
	}
	for _, iface := range interfaces {
		if addrs, err := iface.Addrs(); err == nil {
			for _, addr := range addrs {
				if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() && ip.IP.To4() != nil {
					fmt.Println("Local IP:", ip.IP.String())
				}
			}
		}
	}
	seperator()
}

// Public IP
func getPublicIP() {
	fmt.Println("Fetching public IP address...")
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		log.Printf("Error fetching public IP: %v", err)
		return
	}
	defer resp.Body.Close()
	ip, _ := io.ReadAll(resp.Body)
	fmt.Println("Public IP:", string(ip))
	seperator()
}

// RAM Usage
func getCurrentRamUsage() {
	fmt.Println("Fetching RAM usage...")
	vMem, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error fetching RAM usage: %v", err)
		return
	}
	fmt.Printf("Total: %v MB, Used: %v MB, Free: %v MB, Usage: %.2f%%\n",
		vMem.Total/1024/1024, vMem.Used/1024/1024, vMem.Free/1024/1024, vMem.UsedPercent)
	seperator()
}

// Top 5 RAM-consuming processes
func getTop5RamConsumption() {
	fmt.Println("Fetching top 5 processes by RAM usage...")
	processes, err := process.Processes()
	if err != nil {
		log.Printf("Error fetching processes: %v", err)
		return
	}

	var topProcesses []struct {
		PID    int32
		Name   string
		Memory uint64
	}

	for _, p := range processes {
		memInfo, err := p.MemoryInfo()
		if err != nil {
			continue
		}
		name, err := p.Name()
		if err != nil {
			continue
		}
		topProcesses = append(topProcesses, struct {
			PID    int32
			Name   string
			Memory uint64
		}{PID: p.Pid, Name: name, Memory: memInfo.RSS})
	}

	// Sort processes by memory usage
	sort.Slice(topProcesses, func(i, j int) bool {
		return topProcesses[i].Memory > topProcesses[j].Memory
	})

	// Display top 5 processes
	fmt.Println("Top 5 processes by RAM usage:")
	for i := 0; i < 5 && i < len(topProcesses); i++ {
		fmt.Printf("PID: %d, Name: %s, Memory: %d MB\n", topProcesses[i].PID, topProcesses[i].Name, topProcesses[i].Memory/1024/1024)
	}
	seperator()
}

// Helper function for separators
func seperator() {
	fmt.Println("⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘⋘")
}
