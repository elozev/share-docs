package main

import "fmt"

type ApiKeyCmd struct {
	Create ApiKeyCreateCmd `cmd:"create"`
}

type ApiKeyCreateCmd struct {
	Name string `arg:"" name:"name" help:"create an API for a client"`
}

func (ak *ApiKeyCreateCmd) Run(ctx *Context) error {
	clientName := ak.Name

	fmt.Printf("creating an API key for client: %s \n", clientName)
	return nil
}
