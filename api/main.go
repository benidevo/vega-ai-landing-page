package main

import (
	"api/internal"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("HandleRequest", internal.Application)
}
