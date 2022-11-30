# search-github

Use Search-GithubIssues command to seach github issue based on resource provider or keywords in issue title.

1. Get the command:
```
go install github.com/xiaxyi/Search-GithubIssues@latest
```

2. Example Usage:

Please pass your github token via environment variable, for example `set GITHUB_TOKEN=xxx`

- Search the issue based on Resource provider:
`Search-GithubIssues ResourceProvider --repoOwner hashicorp --repoName terraform-provider-azurerm --startDate 2022-07-01 --endDate 2022-11-01 --resourceProvider "Microsoft.Web" --toCSV c:\\users\\xiaxyi\\desktop\\issues\\rp-web.csv --append "false"`, 
``Search-GithubIssues ResourceProvider --repoOwner hashicorp --repoName terraform-provider-azurerm --startDate 2022-07-01 --endDate 2022-11-01 --resourceProvider "Microsoft.Web" --toCSV c:\\users\\xiaxyi\\desktop\\issues\\rp-eventhub.csv --append "false"``

- Search the issue based on Time range:
`Search-GithubIssues searchOnTimeRange --repoOwner hashicorp --repoName terraform-provider-azurerm --startDate 2022-07-01 --endDate 2022-11-01 --toCSV c:\\users\\xiaxyi\\desktop\\issues\\timerange.csv --append "false"`

- search the issue based on keyWords:
`Search-GithubIssues keyWords --repoOwner hashicorp --repoName terraform-provider-azurerm --startDate 2022-07-01 --endDate 2022-11-01 --keyWords "app service" "Microsoft.Web" --toCSV c:\\users\\xiaxyi\\desktop\\issues\\timerange.csv --append "false"`
