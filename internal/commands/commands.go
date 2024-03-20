package commands

import (
	"debug/buildinfo"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"

	"github.com/o7q2ab/goxm/internal/build"
	"github.com/o7q2ab/goxm/internal/xmmod"
	"github.com/o7q2ab/goxm/internal/xmpath"
)

const (
	// logo created here: https://patorjk.com/software/taag/#p=display&f=ANSI%20Shadow&t=goxm
	logo = `
 ██████╗  ██████╗ ██╗  ██╗███╗   ███╗  version: %s
██╔════╝ ██╔═══██╗╚██╗██╔╝████╗ ████║  runtime: %s
██║  ███╗██║   ██║ ╚███╔╝ ██╔████╔██║  gh repo: https://github.com/o7q2ab/goxm
██║   ██║██║   ██║ ██╔██╗ ██║╚██╔╝██║
╚██████╔╝╚██████╔╝██╔╝ ██╗██║ ╚═╝ ██║
 ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝

`
)

func NewRootCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "goxm",
		Run: func(*cobra.Command, []string) {
			fmt.Printf(
				logo,
				build.Version(),
				runtime.Version(),
			)
		},
	}
	c.AddCommand(
		newBinaryCmd(),
		newPathCmd(),
		newProcCmd(),
		newModuleCmd(),
	)
	return c
}

func newBinaryCmd() *cobra.Command {
	var showDeps, showLatest, showBuildSettings bool

	c := &cobra.Command{
		Use:     "binary [<file-path> | <dir-path>]",
		Aliases: []string{"bin", "b"},
		Short:   "Examine binary file(s) at given path",
		Run: func(cmd *cobra.Command, args []string) {
			printFiles(
				xmpath.List(args[0]),
				showDeps,
				showLatest,
				showBuildSettings,
			)
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

func newPathCmd() *cobra.Command {
	var showDeps, showLatest, showBuildSettings bool

	c := &cobra.Command{
		Use:   "path",
		Short: "Examine all Go binaries found in directories added to PATH environment variable",
		Run: func(cmd *cobra.Command, args []string) {
			printFiles(
				xmpath.ListPathEnv(),
				showDeps,
				showLatest,
				showBuildSettings,
			)
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
				latest := xmmod.GetLatest(r.Mod.Path)
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

func printFiles(names []string, showDeps, showLatest, showBuildSettings bool) {
	short := len(names) == 1

	idx := 0
	for _, name := range names {
		info, err := buildinfo.ReadFile(name)
		if err != nil {
			continue
		}

		if idx != 0 {
			fmt.Println("---------------")
		}
		idx++

		if short {
			fmt.Printf(
				"%s [%s | %d deps | mod: %s]\n",
				info.Path, info.GoVersion, len(info.Deps), info.Main.Path,
			)
		} else {
			fmt.Printf(
				"%d | %s\n%s [%s | %d deps | mod: %s]\n",
				idx, name, info.Path, info.GoVersion, len(info.Deps), info.Main.Path,
			)
		}

		if showDeps || showBuildSettings {
			latest := xmmod.GetLatest(info.Main.Path)
			fmt.Printf("\ncurrent: %s\nlatest: %s\n", info.Main.Version, latest)
		}
		if showDeps {
			fmt.Printf("\nDependencies:\n")
			for _, d := range info.Deps {
				suffix := ""
				if showLatest {
					latest := xmmod.GetLatest(d.Path)
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
	}

	if idx == 0 {
		fmt.Println("No Go binary files were found.")
	}
}
