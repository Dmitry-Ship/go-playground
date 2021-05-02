package chat

import (
	"io"
	"net/http"
	"os"
	"text/template"

	"golang.org/x/net/websocket"
)

type socket struct {
	io.ReadWriter
	done chan bool
}

func (s socket) Close() error {
	s.done <- true
	return nil
}

func SocketHandler(ws *websocket.Conn) {
	s := socket{ws, make(chan bool)}
	go Match(s)
	<-s.done
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	rootTemplate.Execute(w, host+":"+port)
}

var rootTemplate = template.Must(template.New("root").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<script>
	const onMessage = (m) => console.log("Received:", m.data)
	const onClose = () => console.log("tying to close")
    websocket = new WebSocket("ws://{{.}}/socket");
    websocket.onmessage = onMessage;
    websocket.onclose = onClose;
</script>
</html>
`))
