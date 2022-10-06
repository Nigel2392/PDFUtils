package PDFUtils

import (
	"embed"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	PDF_GENERATOR string = "pdf.exe"
	PDF_PY        string = "pdf.py"
)

//go:embed pdf.exe
//go:embed pdf.css
//go:embed pdf.py
var files embed.FS

// Generate a PDF file with raw css supplied
func GeneratePDFFromRawCSS(html string, css string) ([]byte, error) {
	return pdf_generator(html, css, false)
}

// Generate PDF with CSS from a file
func GeneratePDFWithCSSFile(html string, css_file string) ([]byte, error) {
	return pdf_generator(html, css_file, true)
}

// Generate PDF with basic CSS.
func GeneratePDF(html string) ([]byte, error) {
	return pdf_generator(html, "pdf.css", true)
}

func pdf_generator(html string, css_data string, isFile bool) ([]byte, error) {
	var HOST string = "127.0.0.1"
	var PORT int = 22239
	var BUFSIZE int = 1024 * 128 // 128 KB

	_, err := start_program(PDF_GENERATOR)
	if err != nil {
		return nil, err
	}
	// Get css from a file or a string
	var css string
	if isFile {
		css, err = GetCSS(css_data)
		if err != nil {
			return nil, err
		}
	} else {
		css_data = strings.ReplaceAll(css_data, "<style>", "")
		css_data = strings.ReplaceAll(css_data, "</style>", "")
		css = css_data
	}
	css = WrapCSS(css)
	// Get HTML from a string
	html = GetHTMLbase(html, css)
	// Convert HTML to PDF via python client
	pdf := get_PDF_from_client(html, HOST, PORT, BUFSIZE)
	return pdf, nil
}

func WrapCSS(css string) string {
	return "<style>\n" + css + "\n</style>"
}

func GetCSS(filename string) (string, error) {
	css, err := files.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(css), nil
}

func GetHTMLbase(html string, extra string) string {
	prefix := "<html>\n\t<head>\n\t\t<meta charset=\"UTF-8\">\n" + extra + "\n\t</head>\n\t<body>\n"
	sufix := "\n\t</body>\n</html>"
	html = prefix + html + sufix
	return html
}

func start_program(program string) (*exec.Cmd, error) {
	var PACKAGE_READY bool = false
	// Check if python is installed
	if CheckPythonInstalled() {
		// Check if python is version 3
		if CheckPythonVersion("Python 3") {
			if CheckPythonPackageInstalled("xhtml2pdf.pisa as pisa") {
				PACKAGE_READY = true
			} else {
				err := InstallPythonPackage("xhtml2pdf")
				if err != nil {
					create_client()
				} else {
					PACKAGE_READY = true
				}
			}
		} else {
			fmt.Println("Python version is not 3")
			// Create client file for PDF generation
			create_client()
		}
	} else {
		fmt.Println("Python is not installed")
		// Create client file for PDF generation
		create_client()

	}
	var process *exec.Cmd
	if PACKAGE_READY {
		py_file_bytes, err := files.ReadFile(PDF_PY)
		if err != nil {
			return nil, err
		}
		process = exec.Command("python", "-c", string(py_file_bytes))
		fmt.Println("Starting python client")
	} else {
		process = exec.Command(".\\" + program)
		fmt.Println("Starting client")
	}
	if err := process.Start(); err != nil {
		return nil, err
	}
	return process, nil
}

func create_client() error {
	if _, err := os.Stat(PDF_GENERATOR); os.IsNotExist(err) {
		pdf_client, err := files.ReadFile(PDF_GENERATOR)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = ioutil.WriteFile(PDF_GENERATOR, pdf_client, 0777)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func get_PDF_from_client(html string, host string, port int, bufsize int) []byte {
	var str_port string = strconv.Itoa(port)
	// Listen
	listener, err := net.Listen("tcp", host+":"+str_port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer listener.Close()
	// Listen for an incoming connection.
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting: ", err.Error())
		os.Exit(1)
	}
	// Handle connections in a new goroutine.
	buf := make([]byte, bufsize)
	// Read the incoming connection into the buffer.
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	// Close the connection when you're done with it.
	// Send HTML to client
	conn.Write([]byte(html))
	// Flush old buffer, read new buffer
	buf = make([]byte, bufsize)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	// Close the connection when you're done with it.
	conn.Close()
	return buf[:reqLen]
}
