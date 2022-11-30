package cli

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "GetGHIssues",
	Short: "Used to search Gh issues",
}

var cmdListKeyWords = &cobra.Command{
	Use:   "listKeyWords",
	Short: "listing key words of the specified resource provider",
}

func init() {
	RootCmd.PersistentFlags().StringVar(&repoOwner, "repoOwner", "", "specifying the repo owner of the github repo")
	RootCmd.MarkFlagRequired("repoOwner")
	RootCmd.PersistentFlags().StringVar(&repoName, "repoName", "", "specifying the repo name of the github repo")
	RootCmd.MarkFlagRequired("repoName")
	RootCmd.PersistentFlags().StringVar(&startDate, "startDate", "", "search start date")
	RootCmd.MarkFlagRequired("startDate")
	RootCmd.PersistentFlags().StringVar(&endDate, "endDate", "", "search end date")
	RootCmd.MarkFlagRequired("endDate")
	RootCmd.PersistentFlags().StringVar(&csvFileLocation, "toCSV", "", "csv file location for data export")
	RootCmd.PersistentFlags().StringVar(&appendOption, "append", "true", "whether to append result to the csv file")

	CmdGetIssueOnTimeRangeAndResourceProvider.Flags().StringVar(&rpToSearch, "resourceProvider", "", "search based on azure resource provider. i.e. Microsoft.Web")
	CmdGetIssueOnTimeRangeAndResourceProvider.MarkFlagRequired("resourceProvider")

	CmdGetIssueOnTimeRangeAndKeyWords.Flags().StringSliceVar(&keywordsList, "keyWords", nil, "search gh issues beased on user defined keywords")
	CmdGetIssueOnTimeRangeAndKeyWords.MarkFlagRequired("keyWords")

	RootCmd.AddCommand(CmdGetIssueOnTimeRange, CmdGetIssueOnTimeRangeAndKeyWords, CmdGetIssueOnTimeRangeAndResourceProvider)
}
