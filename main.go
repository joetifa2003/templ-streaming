package main

import (
	"bytes"
	"context"
	"fmt"

	"templ-streaming/templates"
)

func main() {
	ctx := context.Background()
	susCtx := templates.NewSuspenseCtx()
	ctx = templates.WithSuspenseCtx(ctx, susCtx)

	out := bytes.Buffer{}
	main := templates.Main()
	main.Render(ctx, &out)
	susCtx.Stream(ctx, &out)

	fmt.Println(string(out.Bytes()))
}
