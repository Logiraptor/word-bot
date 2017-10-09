package main

import (
	"encoding/json"
	"reflect"

	"github.com/Logiraptor/word-bot/web"
	"github.com/gopherjs/gopherjs/js"
)

func main() {
	js.Global.Set("core", map[string]interface{}{
		"RenderBoard": Bridge(web.Render, &web.MoveRequest{}),
	})
}

func Bridge(f interface{}, arg interface{}) func(string) string {
	val := reflect.ValueOf(f)
	return func(body string) string {
		err := json.Unmarshal([]byte(body), arg)
		if err != nil {
			js.Debugger()
		}
		results := val.Call([]reflect.Value{reflect.ValueOf(arg).Elem()})
		buf, err := json.Marshal(results[0].Interface())
		if err != nil {
			js.Debugger()
		}
		return string(buf)
	}
}
