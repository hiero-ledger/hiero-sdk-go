package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// ProtoConfig holds the configuration for Proto files.
type ProtoConfig struct {
	ProtoFiles []string `json:"proto_files"`
}

func main() {
	// Define source and destination directories.
	sourceDir := "../services/hapi/hedera-protobufs/services/state" // Replace with your source directory
	destDir := "../services/hapi/hedera-protobufs/services"         // Replace with your destination directory

	// Move Proto files from source to destination.
	if err := moveProtoFiles(sourceDir, destDir); err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Define source and destination directories.
	sourceDir1 := "../services/hapi/hedera-protobufs/services/auxiliary" // Replace with your source directory

	// Move Proto files from source to destination.
	if err := moveProtoFiles(sourceDir1, destDir); err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Load the Proto configuration from the provided JSON file.
	config, err := loadProtoConfig("../proto/proto.json")
	if err != nil {
		fmt.Println("Error loading proto config:", err)
		return
	}

	// Create a filter for Proto files based on the loaded config.
	protoFilter := createProtoFileFilter(config.ProtoFiles)

	// Remove all existing ".pb.go" files in the "proto/services" directory.
	removeFilesWithExtension(getProjectRootPath(), "proto/services", ".pb.go")

	// Invoke the build process for the Proto files.
	buildProtos(protoFilter)
}

// ensureDirectoryExists checks if the destination directory exists.
func ensureDirectoryExists(destDir string) error {
	_, err := os.Stat(destDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("destination directory does not exist: %v", destDir)
	}
	return nil
}

// moveFile moves a single file from the source path to the destination directory.
func moveFile(sourcePath, destDir string) error {
	destPath := filepath.Join(destDir, filepath.Base(sourcePath))
	if err := os.Rename(sourcePath, destPath); err != nil {
		return fmt.Errorf("error moving file %s: %v", sourcePath, err)
	}
	fmt.Printf("Moved: %s -> %s\n", sourcePath, destPath)
	return nil
}

// moveProtoFiles moves all files from the source directory to the destination.
func moveProtoFiles(sourceDir, destDir string) error {
	if err := ensureDirectoryExists(destDir); err != nil {
		return err
	}
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		return moveFile(path, destDir)
	})
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

// createProtoFileFilter creates a set of Proto files from the provided list.
func createProtoFileFilter(protoFiles []string) map[string]struct{} {
	filter := make(map[string]struct{})
	for _, file := range protoFiles {
		filter[file] = struct{}{}
	}
	return filter
}

// getProjectRootPath returns the root path of the project.
func getProjectRootPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(filename, "../..")
}

// buildProtos builds the Proto files using the provided filter.
func buildProtos(protoFilter map[string]struct{}) {
	var protoFilePaths []string
	var moduleDecls []string

	// Collect Proto file paths and module declarations.
	addProtoFilePaths(&protoFilePaths, &moduleDecls, path.Join(getProjectRootPath(), "services/hapi/hedera-protobufs/services"), protoFilter)

	// Prepare the command arguments for building Proto files.
	cmdArgs := []string{
		"--go_out=proto/services/",
		"--go_opt=paths=source_relative",
		"--go-grpc_out=proto/services/",
		"--go-grpc_opt=paths=source_relative",
		"-Iservices/hapi/hedera-protobufs/services",
	}

	cmdArgs = append(cmdArgs, moduleDecls...)
	cmdArgs = append(cmdArgs, protoFilePaths...)

	// Execute the Proto build command.
	cmd := exec.Command("protoc", cmdArgs...)
	cmd.Dir = getProjectRootPath()

	if err := executeCommand(cmd); err != nil {
		fmt.Println("Error running the build command:", err)
		return
	}

	// Rename package declarations in the generated gRPC files.
	renameGrpcPackageDeclarations(getProjectRootPath(), "proto", "services")
}

// addProtoFilePaths collects Proto file paths and their corresponding module declarations.
func addProtoFilePaths(protoFilePaths *[]string, moduleDecls *[]string, protoPath string, protoFilter map[string]struct{}) {
	err := filepath.Walk(protoPath, func(fullpath string, info fs.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(fullpath, ".proto") {
			return err
		}

		relativePath := strings.TrimPrefix(fullpath, getProjectRootPath()+"/")
		filename := path.Base(fullpath)

		// Skip files that are not in the filter.
		if _, exists := protoFilter[filename]; exists {
			*protoFilePaths = append(*protoFilePaths, relativePath)
			*moduleDecls = append(*moduleDecls,
				fmt.Sprintf("--go_opt=M%s=github.com/hashgraph/hedera-go-sdk/services", filename),
				fmt.Sprintf("--go-grpc_opt=M%s=github.com/hashgraph/hedera-go-sdk/services", filename),
			)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
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
			if err := os.Remove(filename); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

// renameGrpcPackageDeclarations updates the package declaration in gRPC files.
func renameGrpcPackageDeclarations(rootDir, oldPkg, newPkg string) {
	err := filepath.Walk(path.Join(rootDir, newPkg), func(filename string, info fs.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(filename, "_grpc.pb.go") {
			return err
		}

		data, err := os.ReadFile(filename)
		if err != nil {
			return err
		}

		contents := string(data)
		contents = strings.Replace(contents, fmt.Sprintf("package %s", oldPkg), fmt.Sprintf("package %s", newPkg), 1)

		return os.WriteFile(filename, []byte(contents), info.Mode())
	})
	if err != nil {
		panic(err)
	}
}
