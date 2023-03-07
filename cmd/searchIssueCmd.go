package cmd

import (
	"context"
	"github.com/google/go-github/v47/github"
	"github.com/spf13/cobra"
	"github.com/xiaxyi/Search-GithubIssues/helper/model"
	"github.com/xiaxyi/Search-GithubIssues/searches"
	"golang.org/x/oauth2"
	"net/url"
	"os"
)

var repoOwner, repoName string
var rpToSearch, startDate, endDate, csvFileLocation string
var appendOption string
var CmdGetIssueOnTimeRange = &cobra.Command{
	Use:   "searchOnTimeRange",
	Short: "Get GH issue based on specified time range",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		token := oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")}
		ts := oauth2.StaticTokenSource(&token)

		tokenClient := oauth2.NewClient(ctx, ts)
		ghClient := github.NewClient(tokenClient)
		ghClient.BaseURL, _ = url.Parse("https://api.github.com/")
		toFile := false
		if csvFileLocation != "" {
			toFile = true
		}
		searchConditionOnTimeOnly := model.Conditions{
			RepoOwner:        repoOwner,
			RepoName:         repoName,
			CreatedTimeStart: startDate,
			CreatedTimeEnd:   endDate,
		}

		//filter := fmt.Sprintf("repo:%s/%s", repoOwner, repoName)
		//repo, _, err := ghClient.Search.Repositories(ctx, filter, nil)
		//if repo.Repositories == nil {
		//	return fmt.Errorf("connection to Github repo error: %+v", err)
		//}

		if err := searches.GetIssueOnCreationTimeRange(ctx, ghClient, searchConditionOnTimeOnly, csvFileLocation, appendOption, toFile); err != nil {
			return err
		}
		return nil
	},
}

var CmdGetIssueOnTimeRangeAndResourceProvider = &cobra.Command{
	Use:   "ResourceProvider",
	Short: "Get GH issue based on resource provider",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		token := oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")}
		ts := oauth2.StaticTokenSource(&token)

		tokenClient := oauth2.NewClient(ctx, ts)
		ghClient := github.NewClient(tokenClient)
		ghClient.BaseURL, _ = url.Parse("https://api.github.com/")

		searchConditionWithRp := model.Conditions{
			RepoOwner:        repoOwner,
			RepoName:         repoName,
			ResourceProvider: rpToSearch,
			CreatedTimeStart: startDate,
			CreatedTimeEnd:   endDate,
		}
		KeyWordsList := InitiatingResourceProvider(rpToSearch).SysKeyWordsList

		if err := searches.GetIssueBasedOnKeyWordListAndTimeRange(ctx, ghClient, searchConditionWithRp, KeyWordsList, csvFileLocation, appendOption); err != nil {
			return err
		}
		return nil
	},
}

var keywordsList []string
var CmdGetIssueOnTimeRangeAndKeyWords = &cobra.Command{
	Use:   "keyWords",
	Short: "search issue based on user defined key words",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		token := oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")}
		ts := oauth2.StaticTokenSource(&token)

		tokenClient := oauth2.NewClient(ctx, ts)
		ghClient := github.NewClient(tokenClient)
		ghClient.BaseURL, _ = url.Parse("https://api.github.com/")

		searchConditionWithRp := model.Conditions{
			RepoOwner:        repoOwner,
			RepoName:         repoName,
			CreatedTimeStart: startDate,
			CreatedTimeEnd:   endDate,
		}

		if err := searches.GetIssueBasedOnKeyWordListAndTimeRange(ctx, ghClient, searchConditionWithRp, keywordsList, csvFileLocation, appendOption); err != nil {
			return err
		}
		return nil
	},
}
