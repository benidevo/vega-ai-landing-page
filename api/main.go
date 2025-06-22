package main

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"vega.ai/landing-api/internal"
)

func init() {
	functions.HTTP("HandleRequest", internal.Application)
}
