package cli

import (
	"github.com/spf13/cobra"
)

var KeyWordsList = &cobra.Command{
	Short: "List key words based on resource provider",
}

type ResourceProviderObj struct {
	RpName           string
	SysKeyWordsList  []string
	UserKeyWordsList []string
}

var (
	AppServiceKeyWordList = []string{"app service", "web", "function_app", "web_app", "linux_web_app", "app_service", "FUNCTIONS_WORKER_RUNTIME", "WEBSITE_NODE_DEFAULT_VERSION", "stack", "app_settings"}
	EventHubKeyWordList   = []string{"event hub", "eventhub namespace"}
	ServiceBusKeyWordList = []string{"service bus", "servicebus "}
)

var KeyWordsListSet = map[string]interface{}{
	"Microsoft.Web":        AppServiceKeyWordList,
	"Microsoft.EventHub":   EventHubKeyWordList,
	"Microsoft.ServiceBus": ServiceBusKeyWordList,
}

//todo: check returned value is {"Microsoft.Web": {"app service", "web"...}}
func InitiatingResourceProvider(rpName string) *ResourceProviderObj {
	keyWordsToUse := KeyWordsListSet[rpName]
	return &ResourceProviderObj{
		RpName:          rpName,
		SysKeyWordsList: keyWordsToUse.([]string),
	}
}

func AddKeyWordsList(rpName string, keywordsToAdd []string) []string {
	systemAssignedKeyWords := KeyWordsListSet[rpName].([]string)

	for _, kw := range keywordsToAdd {
		systemAssignedKeyWords = append(systemAssignedKeyWords, kw)
	}

	return systemAssignedKeyWords
}
