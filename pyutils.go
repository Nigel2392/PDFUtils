package PDFUtils

import (
	"errors"
	"os/exec"
	"strings"
)

func CheckPythonInstalled() bool {
	_, err := exec.LookPath("python")
	// fmt.Println("Python installed: " + strconv.FormatBool(err == nil))
	return err == nil
}

func CheckPythonPackageInstalled(pkg string) bool {
	cmd := exec.Command("python", "-c", `import `+pkg)
	err := cmd.Run()
	// fmt.Println("Python package " + pkg + " installed: " + strconv.FormatBool(err == nil))
	return err == nil
}

func InstallPythonPackage(pkg string) error {
	cmd := exec.Command("python", "-m", "pip", "install", pkg)
	err := cmd.Run()
	if err != nil {
		err = errors.New("failed to install package " + pkg)
		// fmt.Println(err)
		return err
	}
	// fmt.Println("Package " + pkg + " installed")
	return nil
}

func CheckPythonVersion(version string) bool {
	version_to_check := strings.ToLower(version)
	cmd := exec.Command("python", "--version")
	out, err := cmd.Output()
	if err != nil {
		return false
	} else {
		// fmt.Println("Version: " + string(out))
		version := strings.ToLower(string(out))
		if strings.Contains(version, version_to_check) {
			return true
		}
	}
	return false
}
