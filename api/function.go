package function

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/benidevo/vega-ai-landing-page/api/internal"
)

func init() {
	functions.HTTP("HandleRequest", internal.Application)
}
