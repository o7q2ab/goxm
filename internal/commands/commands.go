package commands

import (
	"debug/buildinfo"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
)

func NewRootCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "goxm",
		Run: func(*cobra.Command, []string) {
			fmt.Println("goxm v0.1.0")
			fmt.Println(runtime.Version())
			fmt.Println("https://github.com/o7q2ab/goxm")
		},
	}
	c.AddCommand(
		newBinaryCmd(),
		newProcCmd(),
		newModuleCmd(),
	)
	return c
}

func newBinaryCmd() *cobra.Command {
	var showDeps, showLatest, showBuildSettings bool

	c := &cobra.Command{
		Use:     "binary <file-path>",
		Aliases: []string{"bin", "b"},
		Short:   "Examine binary file",
		Run: func(cmd *cobra.Command, args []string) {
			p := args[0]
			stat, err := os.Stat(p)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			if stat.IsDir() {
				fmt.Println("error: directory")
				return
			}
			info, err := buildinfo.ReadFile(p)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Printf(
				"%s [%s | %d deps | mod: %s]\n",
				info.Path, info.GoVersion, len(info.Deps), info.Main.Path,
			)
			if showDeps {
				fmt.Printf("\nDependencies:\n")
				for _, d := range info.Deps {
					suffix := ""
					if showLatest {
						latest := getLatest(d.Path)
						if latest == "" {
							suffix = " (latest: unknown)"
						} else {
							suffix = fmt.Sprintf(" (latest: %s)", latest)
						}
					}
					fmt.Printf("    %s %s%s\n", d.Path, d.Version, suffix)
				}
			}
			if showBuildSettings {
				fmt.Printf("\nBuild settings:\n")
				for _, s := range info.Settings {
					fmt.Printf("    %s=%s\n", s.Key, s.Value)
				}
			}
		},
	}

	c.Flags().BoolVarP(
		&showDeps, "deps", "d", false, "show all the dependency modules",
	)
	c.Flags().BoolVar(
		&showLatest, "latest", false, "show latest versions for all the dependency modules",
	)
	c.Flags().BoolVarP(
		&showBuildSettings, "build", "b", false, "show the build settings used to build the binary",
	)

	return c
}

func newProcCmd() *cobra.Command {
	var showDeps, showBuildSettings, showConn bool
	var filter string

	addrFamilies := []string{
		"AF_UNSPEC",
		"AF_UNIX",
		"AF_INET",
		"AF_AX25",
		"AF_IPX",
		"AF_APPLETALK",
		"AF_NETROM",
		"AF_BRIDGE",
		"AF_ATMPVC",
		"AF_X25",
		"AF_INET6",
		"AF_ROSE",
		"AF_DECnet",
		"AF_NETBEUI",
		"AF_SECURITY",
		"AF_KEY",
		"AF_NETLINK",
		"AF_PACKET",
		"AF_ASH",
		"AF_ECONET",
		"AF_ATMSVC",
		"AF_RDS",
		"AF_SNA",
		"AF_IRDA",
		"AF_PPPOX",
		"AF_WANPIPE",
		"AF_LLC",
		"AF_IB",
		"AF_MPLS",
		"AF_CAN",
		"AF_TIPC",
		"AF_BLUETOOTH",
		"AF_IUCV",
		"AF_RXRPC",
		"AF_ISDN",
		"AF_PHONET",
		"AF_IEEE802154",
		"AF_CAIF",
		"AF_ALG",
		"AF_NFC",
		"AF_VSOCK",
		"AF_KCM",
		"AF_QIPCRTR",
		"AF_SMC",
		"AF_XDP",
		"AF_MCTP",
	}

	c := &cobra.Command{
		Use:     "process [<pid>]",
		Aliases: []string{"proc", "ps", "p"},
		Short:   "Examine currently running Go processes",
		Run: func(cmd *cobra.Command, args []string) {
			var all []*process.Process
			var err error
			if len(args) == 0 {
				all, err = process.Processes()
				if err != nil {
					fmt.Println("error:", err)
					return
				}
			} else {
				pid, err := strconv.Atoi(args[0])
				if err != nil {
					fmt.Println("error:", err)
					return
				}
				pr, err := process.NewProcess(int32(pid))
				if err != nil {
					fmt.Println("error:", err)
					return
				}
				all = []*process.Process{pr}
				filter = ""
			}

			idx := 0
			for _, p := range all {
				path, err := p.Exe()
				if err != nil {
					continue
				}
				info, err := buildinfo.ReadFile(path)
				if err != nil {
					continue
				}

				if filter != "" && !strings.Contains(info.Main.Path, filter) {
					continue
				}

				if idx != 0 {
					fmt.Println("---------------")
				}
				idx++

				parent, err := p.Ppid()
				if err != nil {
					fmt.Printf("%d | PID: %d [parent: %v]\n", idx, p.Pid, err)
				} else {
					fmt.Printf("%d | PID: %d [parent: %d]\n", idx, p.Pid, parent)
				}
				fmt.Printf(
					"%s [%s | %d deps | mod: %s]\n",
					info.Path, info.GoVersion, len(info.Deps), info.Main.Path,
				)

				if showDeps {
					fmt.Printf("\nDependencies:\n")
					for _, d := range info.Deps {
						fmt.Printf("    %s %s\n", d.Path, d.Version)
					}
				}

				if showBuildSettings {
					fmt.Printf("\nBuild settings:\n")
					for _, s := range info.Settings {
						fmt.Printf("    %s=%s\n", s.Key, s.Value)
					}
				}

				if showConn {
					fmt.Printf("\nConnections:\n")
					conns, err := p.Connections()
					if err != nil {
						fmt.Println("    error:", err)
						continue
					}
					if len(conns) == 0 {
						fmt.Println("    no connections")
					}
					for _, c := range conns {
						fmt.Printf(
							"    %s %s:%d - %s:%d\n",
							addrFamilies[c.Family],
							c.Laddr.IP, c.Laddr.Port,
							c.Raddr.IP, c.Raddr.Port,
						)
					}
				}
			}
		},
	}

	c.Flags().BoolVarP(
		&showDeps, "deps", "d", false, "show all the dependency modules",
	)
	c.Flags().BoolVarP(
		&showBuildSettings, "build", "b", false, "show the build settings used to build the binary",
	)
	c.Flags().BoolVar(
		&showConn, "conn", false, "show all the connections (TCP, UDP, Unix) used by the process",
	)
	c.Flags().StringVar(
		&filter, "filter", "", "filter by the package name",
	)

	return c
}

func newModuleCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "module [<file-path>]",
		Aliases: []string{"mod", "m"},
		Short:   "Examine Go module",
		Run: func(cmd *cobra.Command, args []string) {
			var p string
			if len(args) != 0 {
				p = args[0]
			} else {
				p, _ = os.Getwd()
			}

			stat, err := os.Stat(p)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			if stat.IsDir() {
				p = filepath.Join(p, "go.mod")
				_, err = os.Stat(p)
				if err != nil {
					fmt.Println("error:", err)
					return
				}
			}
			f, err := os.ReadFile(p)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			modf, err := modfile.Parse(p, f, nil)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println(modf.Module.Mod.Path)
			for _, r := range modf.Require {
				latest := getLatest(r.Mod.Path)
				suffix := ""
				if latest == "" {
					suffix = " (latest: unknown)"
				} else if latest != r.Mod.Version {
					suffix = fmt.Sprintf(" (latest: %s)", latest)
				}
				if r.Indirect {
					fmt.Printf("    [indirect] %s%s\n", r.Mod, suffix)
				} else {
					fmt.Printf("    %s%s\n", r.Mod, suffix)
				}
			}
		},
	}
	return c
}

func getLatest(modpath string) string {
	// GOPROXY protocol: https://go.dev/ref/mod#goproxy-protocol

	resp, err := http.Get("https://proxy.golang.org/" + modpath + "/@latest")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	latest := map[string]any{}
	err = json.NewDecoder(resp.Body).Decode(&latest)
	if err != nil {
		return ""
	}

	v, ok := latest["Version"].(string)
	if !ok {
		return ""
	}
	return v
}
