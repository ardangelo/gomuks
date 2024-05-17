// gomuks - A terminal Matrix client written in Go.
// Copyright (C) 2020 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	flag "maunium.net/go/mauflag"

	"maunium.net/go/gomuks/debug"
	"maunium.net/go/gomuks/initialize"
	ifc "maunium.net/go/gomuks/interface"
	"maunium.net/go/gomuks/matrix"
	"maunium.net/go/gomuks/ui"
)

// Information to find out exactly which commit gomuks was built from.
// These are filled at build time with the -X linker flag.
var (
	Tag       = "unknown"
	Commit    = "unknown"
	BuildTime = "unknown"
)

var (
	// Version is the version number of gomuks. Changed manually when making a release.
	Version = "0.3.0"
	// VersionString is the gomuks version, plus commit information. Filled in init() using the build-time values.
	VersionString = ""
)

func init() {
	if len(Tag) > 0 && Tag[0] == 'v' {
		Tag = Tag[1:]
	}
	if Tag != Version {
		suffix := ""
		if !strings.HasSuffix(Version, "+dev") {
			suffix = "+dev"
		}
		if len(Commit) > 8 {
			Version = fmt.Sprintf("%s%s.%s", Version, suffix, Commit[:8])
		} else {
			Version = fmt.Sprintf("%s%s.unknown", Version, suffix)
		}
	}
	VersionString = fmt.Sprintf("gomuks %s (%s with %s)", Version, BuildTime, runtime.Version())
}

var MainUIProvider ifc.UIProvider = ui.NewGomuksUI

var wantVersion = flag.MakeFull("v", "version", "Show the version of gomuks", "false").Bool()
var clearCache = flag.MakeFull("c", "clear-cache", "Clear the cache directory instead of starting", "false").Bool()
var skipVersionCheck = flag.MakeFull("s", "skip-version-check", "Skip the homeserver version checks at startup and login", "false").Bool()
var printLogPath = flag.MakeFull("l", "print-log-path", "Print the log path instead of starting", "false").Bool()
var clearData = flag.Make().LongKey("clear-all-data").Usage("Clear all data instead of starting").Default("false").Bool()
var headless = flag.Make().LongKey("headless").Usage("Update new messages and exit").Default("false").Bool()
var benchmarkMode = flag.Make().LongKey("benchmark").Usage("Run benchmark in given cache directory").Default("false").Bool()
var wantHelp, _ = flag.MakeHelpFlag()

func main() {
	flag.SetHelpTitles(
		"gomuks - A terminal Matrix client written in Go.",
		"gomuks [-vcsh] [--clear-all-data|--print-log-path]",
	)
	err := flag.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else if *wantHelp {
		flag.PrintHelp()
		return
	} else if *wantVersion {
		fmt.Println(VersionString)
		return
	}

	debugDir := os.Getenv("DEBUG_DIR")
	if len(debugDir) > 0 {
		debug.LogDirectory = debugDir
	}
	debugLevel := strings.ToLower(os.Getenv("DEBUG"))
	if debugLevel == "1" || debugLevel == "t" || debugLevel == "true" {
		debug.RecoverPrettyPanic = false
		debug.DeadlockDetection = true
		debug.WriteLogs = true
	}
	if *printLogPath {
		fmt.Println(debug.LogFile())
		return
	}

	debug.Initialize()
	defer debug.Recover()

	if *benchmarkMode {
		tempDir, err := os.MkdirTemp("", "gomuks-benchmark")

		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Failed to create temporary directory:", err)
			os.Exit(3)
		}

		defer func() {
			err := os.RemoveAll(tempDir)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Failed to remove temporary directory:", err)
			}
		}()

		os.Setenv("GOMUKS_ROOT", tempDir)
	}

	var configDir, dataDir, cacheDir, downloadDir string

	configDir, err = initialize.UserConfigDir()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to get config directory:", err)
		os.Exit(3)
	}
	dataDir, err = initialize.UserDataDir()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to get data directory:", err)
		os.Exit(3)
	}
	cacheDir, err = initialize.UserCacheDir()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to get cache directory:", err)
		os.Exit(3)
	}
	downloadDir, err = initialize.UserDownloadDir()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to get download directory:", err)
		os.Exit(3)
	}

	debug.Print("Config directory:", configDir)
	debug.Print("Data directory:", dataDir)
	debug.Print("Cache directory:", cacheDir)
	debug.Print("Download directory:", downloadDir)

	gmx := initialize.NewGomuks(MainUIProvider, configDir, dataDir, cacheDir, downloadDir)
	if *skipVersionCheck {
		gmx.Matrix().(*matrix.Container).SetSkipVersionCheck()
	}
	if *headless {
		gmx.Matrix().(*matrix.Container).SetHeadless()
	}
	if *benchmarkMode {
		gmx.Matrix().(*matrix.Container).SetBenchmarkMode()
	}

	if *clearCache {
		debug.Print("Clearing cache as requested by CLI flag")
		gmx.Config().Clear()
		fmt.Printf("Cleared cache at %s\n", gmx.Config().CacheDir)
		return
	} else if *clearData {
		debug.Print("Clearing all data as requested by CLI flag")
		gmx.Config().Clear()
		gmx.Config().ClearData()
		_ = os.RemoveAll(gmx.Config().Dir)
		fmt.Printf("Cleared cache at %s, data at %s and config at %s\n", gmx.Config().CacheDir, gmx.Config().DataDir, gmx.Config().Dir)
		return
	}

	gmx.Start()

	// We use os.Exit() everywhere, so exiting by returning from Start() shouldn't happen.
	time.Sleep(5 * time.Second)
	fmt.Println("Unexpected exit by return from gmx.Start().")
	os.Exit(2)
}
