package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

const goPackage = "github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

// importPattern matches a proto import statement and captures the imported path.
var importPattern = regexp.MustCompile(`(?m)^(\s*import\s+(?:public\s+|weak\s+)?")([^"]+)(")`)

// ProtoConfig holds the configuration for Proto files.
type ProtoConfig struct {
	ProtoFiles []string `json:"proto_files"`
}

func main() {
	var protoSourceDirs []string
	protoServicesDir := flag.String("dest", "", "Path to the services proto directory")

	flag.Func("source", "Comma-separated list of additional source directories", func(s string) error {
		protoSourceDirs = strings.Split(s, ",")
		return nil
	})

	flag.Parse()

	if *protoServicesDir == "" {
		fmt.Println("Usage: go run scripts/generators/proto/generator.go -dest <services_dir> [-source <dir1,dir2,...>]")
		os.Exit(2)
	}

	if err := run(*protoServicesDir, protoSourceDirs); err != nil {
		log.Fatal(err)
	}
}

// run stages the proto sources, then regenerates proto/services from the files listed in proto/proto.json.
func run(protoServicesDir string, protoSourceDirs []string) error {
	config, err := loadProtoConfig(path.Join(getProjectRootPath(), "proto/proto.json"))
	if err != nil {
		return fmt.Errorf("loading proto config: %w", err)
	}

	// Stage a flattened copy of the proto tree so the submodule checkout is never modified.
	stagingDir, err := os.MkdirTemp("", "hiero-protos-")
	if err != nil {
		return fmt.Errorf("creating staging directory: %w", err)
	}
	defer os.RemoveAll(stagingDir)

	if err := stageProtoFiles(append([]string{protoServicesDir}, protoSourceDirs...), stagingDir); err != nil {
		return fmt.Errorf("staging proto files: %w", err)
	}

	protoFiles, err := selectProtoFiles(config.ProtoFiles, stagingDir)
	if err != nil {
		return err
	}

	removeFilesWithExtension(getProjectRootPath(), "proto/services", ".pb.go")

	if err := buildProtos(protoFiles, stagingDir); err != nil {
		return fmt.Errorf("running protoc: %w", err)
	}

	return nil
}

// stageProtoFiles copies every .proto file under the given directories into stagingDir,
// flattened to its basename, and rewrites nested imports to those basenames.
// Two different files with the same basename are reported instead of overwritten.
func stageProtoFiles(sourceDirs []string, stagingDir string) error {
	staged := make(map[string]string)

	for _, dir := range sourceDirs {
		dir = filepath.Clean(dir)
		if _, err := os.Stat(dir); err != nil {
			return fmt.Errorf("source directory %s: %w", dir, err)
		}

		err := filepath.WalkDir(dir, func(sourcePath string, entry fs.DirEntry, err error) error {
			if err != nil || entry.IsDir() || !strings.HasSuffix(sourcePath, ".proto") {
				return err
			}

			name := filepath.Base(sourcePath)
			if previous, exists := staged[name]; exists {
				if previous == sourcePath {
					return nil
				}
				return fmt.Errorf("%s and %s both flatten to %s", previous, sourcePath, name)
			}
			staged[name] = sourcePath

			// G122/G703: build-time generator reading the local submodule checkout into its own temp directory.
			content, err := os.ReadFile(sourcePath) //nolint:gosec
			if err != nil {
				return err
			}

			return os.WriteFile(filepath.Join(stagingDir, name), flattenImports(content), 0o600) //nolint:gosec
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// flattenImports rewrites every non-Google nested import to its basename.
func flattenImports(content []byte) []byte {
	return importPattern.ReplaceAllFunc(content, func(match []byte) []byte {
		parts := importPattern.FindSubmatch(match)
		imported := string(parts[2])
		if strings.HasPrefix(imported, "google/protobuf/") {
			return match
		}
		return []byte(string(parts[1]) + path.Base(imported) + string(parts[3]))
	})
}

// selectProtoFiles returns the sorted, de-duplicated allowlisted files, failing on any entry that was not staged.
func selectProtoFiles(allowlist []string, stagingDir string) ([]string, error) {
	unique := make(map[string]struct{})
	var missing []string

	for _, name := range allowlist {
		if _, seen := unique[name]; seen {
			continue
		}
		unique[name] = struct{}{}
		if _, err := os.Stat(filepath.Join(stagingDir, name)); err != nil {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("proto/proto.json lists files that do not exist in the sources: %s", strings.Join(missing, ", "))
	}

	files := make([]string, 0, len(unique))
	for name := range unique {
		files = append(files, name)
	}
	sort.Strings(files)

	return files, nil
}

// loadProtoConfig loads the Proto configuration from a JSON file.
func loadProtoConfig(filename string) (ProtoConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return ProtoConfig{}, err
	}
	defer file.Close()

	var config ProtoConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return ProtoConfig{}, err
	}

	return config, nil
}

// getProjectRootPath returns the root path of the project.
//
//nolint:dogsled
func getProjectRootPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(filename, "../../../..")
}

// buildProtos compiles the given staged proto files into proto/services.
func buildProtos(protoFiles []string, stagingDir string) error {
	cmdArgs := []string{
		"--go_out=proto/services/",
		"--go_opt=paths=source_relative",
		"--go-grpc_out=proto/services/",
		"--go-grpc_opt=paths=source_relative",
		"-I" + stagingDir,
	}

	for _, name := range protoFiles {
		cmdArgs = append(cmdArgs,
			fmt.Sprintf("--go_opt=M%s=%s", name, goPackage),
			fmt.Sprintf("--go-grpc_opt=M%s=%s", name, goPackage),
		)
	}

	for _, name := range protoFiles {
		cmdArgs = append(cmdArgs, filepath.Join(stagingDir, name))
	}

	cmd := exec.Command("protoc", cmdArgs...)
	cmd.Dir = getProjectRootPath()

	return executeCommand(cmd)
}

// executeCommand runs a command and checks for errors.
func executeCommand(cmd *exec.Cmd) error {
	_, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			fmt.Print(string(exitErr.Stderr))
			return fmt.Errorf("command failed with exit code %d", exitErr.ExitCode())
		}
		return err
	}
	return nil
}

// removeFilesWithExtension removes all files with the given extension in the specified directory.
func removeFilesWithExtension(rootDir, module, ext string) {
	err := filepath.Walk(path.Join(rootDir, module), func(filename string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(filename, ext) {
			// G122: This is a build-time code generator script, not production code.
			// The TOCTOU risk from symlink traversal is acceptable in this context.
			if err := os.Remove(filename); err != nil { //nolint:gosec
				return err
			}
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}
