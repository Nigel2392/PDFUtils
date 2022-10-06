To use:

This package works only for Windows, unless Python3 is installed. With Python3 it works on any system.
When python is installed, the package will check if python3 has the module "xhtml2pdf" installed.
If not, it will install it.

If Python3 is not installed, the package will use the embedded executable. 
It will place this in your current working directory to be able to execute it.
This EXE is a compiled Python script, which is around 6 seconds slower than when using Python3.

Creates a local socket server on port 22239. 
Golang then sends a packet with the HTML encoded in bytes, which is then converted to a PDF in the Python script.

```
go get github.com/Nigel2392/GoPDF
```
Then import the package in go:
```
import (
	"github.com/Nigel2392/GoPDF/PDFUtils"
)
```
To get a PDF from the HTML:
```
pdf, err := PDFUtils.GeneratePDF("<h1>test</h1><code>TEST CODE</code><h1>test</h1><h1>test</h1><h1>test</h1>")
pdf, err := PDFUtils.GeneratePDFFromRawCSS("<h1>test</h1><code>TEST CODE</code><h1>test</h1><h1>test</h1><h1>test</h1>", "h1{color:red;}")

// To use the file option, first you have to embed it in your go file.
// Then you can use it like this:
pdf, err := PDFUtils.GeneratePDFWithCSSFile("<h1>test</h1><code>TEST CODE</code><h1>test</h1><h1>test</h1><h1>test</h1>", "css_file_name.css")
```
where pdf is a []byte object. 
You can then write this to a file or send it to a client.
