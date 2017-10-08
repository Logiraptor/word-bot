package main

import "github.com/gopherjs/gopherjs/js"
import "github.com/Logiraptor/word-bot/web"
import "encoding/json"

func main() {
	js.Global.Set("core", map[string]interface{}{
		"RenderBoard": Render,
	})
}

func Render(body string) string {
	var req web.MoveRequest
	err := json.Unmarshal([]byte(body), &req)
	if err != nil {
		js.Debugger()
	}
	output := web.Render(req)
	buf, err := json.Marshal(output)
	if err != nil {
		js.Debugger()
	}
	return string(buf)
}
