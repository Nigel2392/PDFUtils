################################################################
# Easily create a PDF from HTML.
# 
# Usage:
#
# Create a server listening on the loopback address, on port 22239.
# Send the HTML from the server to this client.
# We will return the PDF data as bytes to the socket server.
# 
# Compile with:
# pyinstaller --onefile --noconsole --collect-all reportlab.graphics.barcode pdf.py
################################################################

from socket import socket, AF_INET, SOCK_STREAM
import xhtml2pdf.pisa as pisa
from io import BytesIO

HOST = "127.0.0.1"  # Standard loopback interface address (localhost)
PORT = 22239  # Port to listen on (non-privileged ports are > 1023)
BUFSIZE = 1024 * 128 # Big enough buffer to hold HTML data.

def pad(s):
    if len(s) > 28:
        raise Exception("String too long")
    return "$$" + s + "_" * (28 - len(s)) + "$$"

with socket(AF_INET, SOCK_STREAM) as s:
    s.connect((HOST, PORT))
    s.send(bytes(pad("CONNECTED").encode("utf-8")))
    HTML = s.recv(BUFSIZE).decode("utf-8")
    print(HTML)
    buf = BytesIO()
    pisa.CreatePDF(HTML, buf)
    s.send(buf.getvalue())
    s.close()
