package searches

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"Search-GithubIssues/helper/model"

	"github.com/google/go-github/v47/github"
)

func GetIssueOnCreationTimeRange(ctx context.Context, ghClient *github.Client, condition model.Conditions, desFile string, appendOption string) error {
	filter := fmt.Sprintf("repo:%s/%s is:issue created:%s..%s", condition.RepoOwner, condition.RepoName, condition.CreatedTimeStart, condition.CreatedTimeEnd)
	totalIssueSet, _, err := ghClient.Search.Issues(ctx, filter, &github.SearchOptions{Order: "desc", ListOptions: github.ListOptions{PerPage: 100}})
	if totalIssueSet == nil {
		return fmt.Errorf("no issue Found by consition: %s, suggested to check your github token", filter)
	}
	if err != nil {
		return fmt.Errorf("retrieving issues error: %+v", err)
	}

	issueList := map[int]interface{}{}
	var totalIssueCount int
	if totalIssueSet.Total != nil {
		totalIssueCount = *totalIssueSet.Total
		pageCount := totalIssueCount / 100
		pageLeft := (totalIssueCount - 100*pageCount) % 100
		if pageLeft > 0 {
			pageCount += 1
		}
		for i := 1; i <= pageCount; i++ {
			issueResult, _, err := ghClient.Search.Issues(ctx, filter, &github.SearchOptions{Order: "desc", ListOptions: github.ListOptions{Page: i, PerPage: 100}})
			if err != nil {
				return fmt.Errorf("retrieving issues error: %+v", err)
			}
			AddIssueToList(issueResult, issueList)
		}
	}

	if desFile != "" {
		ExportResultToLocalCsv(issueList, desFile, appendOption)
	}
	return nil
}

func GetIssueBasedOnKeyWordListAndTimeRange(ctx context.Context, ghClient *github.Client, condition model.Conditions, keyWordsList []string, desFile string, appendOption string) error {
	issueList := map[int]interface{}{}
	for _, kw := range keyWordsList {
		filter := fmt.Sprintf("repo:%s/%s is:issue in:title %s created:%s..%s", condition.RepoOwner, condition.RepoName, kw, condition.CreatedTimeStart, condition.CreatedTimeEnd)
		err := KeyWordAndTimeRangeBasedIssueQueries(ctx, ghClient, filter, issueList)
		if err != nil {
			log.Printf("querying issue error:%+v", err)
			return err
		}
	}
	if desFile != "" {
		ExportResultToLocalCsv(issueList, desFile, appendOption)
	}
	return nil
}

func KeyWordAndTimeRangeBasedIssueQueries(ctx context.Context, ghClient *github.Client, filter string, issueList map[int]interface{}) error {
	totalIssueSet, _, err := ghClient.Search.Issues(ctx, filter, &github.SearchOptions{Order: "asc", ListOptions: github.ListOptions{PerPage: 100}})
	if totalIssueSet == nil {
		return fmt.Errorf("no issue Found by consition: %s, suggested to check your github token", filter)
	}
	var totalIssueCount int
	if totalIssueSet.Total != nil {
		totalIssueCount = *totalIssueSet.Total
		pageCount := totalIssueCount / 100
		pageLeft := (totalIssueCount - 100*pageCount) % 100
		if pageLeft > 0 {
			pageCount += 1
		}
		for i := 1; i <= pageCount; i++ {
			issueResult, _, err := ghClient.Search.Issues(ctx, filter, &github.SearchOptions{Order: "desc", ListOptions: github.ListOptions{Page: i, PerPage: 100}})
			if err != nil {
				return fmt.Errorf("retrieving issues error: %+v", err)
			}
			AddIssueToList(issueResult, issueList)
		}
	}
	if err != nil {
		return fmt.Errorf("retrieving issues error: %+v", err)
	}

	return nil
}

func AddIssueToList(ghResult *github.IssuesSearchResult, issueList map[int]interface{}) {
	var issueNum int

	for _, issue := range ghResult.Issues {
		var dataInBuffer []string
		item := map[string]interface{}{
			"serviceLabel": "NA",
			"labeledBug":   "NA",
		}
		issueType := "bugs"
		if issue.Number != nil {
			issueNum = *issue.Number
			item["id"] = issueNum
		}
		if issue.Title != nil {
			issueTitle := *issue.Title
			item["title"] = strings.ReplaceAll(*issue.Title, "'", "\\'")
			dataInBuffer = append(dataInBuffer, issueTitle)
			if strings.Contains(issueTitle, "support") || strings.Contains(issueTitle, "Support") {
				issueType = "enhancement"
				item["type"] = issueType
			}
		}
		if issue.URL != nil {
			urlPrefix := "https://github.com/"
			apiUrl := *issue.URL
			apiUrl = strings.TrimPrefix(apiUrl, "https://api.github.com/repos/")
			endpoint := urlPrefix + apiUrl
			item["issueLink"] = endpoint
			dataInBuffer = append(dataInBuffer, endpoint)
		}
		if issue.State != nil {
			item["state"] = *issue.State
			dataInBuffer = append(dataInBuffer, *issue.State)
		}
		if issue.CreatedAt != nil {
			created := (*issue.CreatedAt).Format("2006-01-02")
			item["createdAt"] = created
			dataInBuffer = append(dataInBuffer, created)
		}
		if issue.UpdatedAt != nil {
			updated := (*issue.UpdatedAt).Format("2006-01-02")
			item["updatedAt"] = updated
			dataInBuffer = append(dataInBuffer, updated)
		}

		votes := 0
		if reactions := issue.Reactions; reactions != nil {
			if reactions.PlusOne != nil {
				votes = *reactions.PlusOne
				item["votes"] = votes
				dataInBuffer = append(dataInBuffer, strconv.Itoa(item["votes"].(int)))
			}
		}

		labels := issue.Labels
		if labels != nil && len(labels) > 0 {
			for _, label := range labels {
				if label != nil {
					if label.Name != nil {
						if strings.Contains(*label.Name, "service") {
							item["serviceLabel"] = *label.Name
						} else if strings.Contains(*label.Name, "bug") {
							item["labeledBug"] = *label.Name
						}
					}
				}
			}
		}
		dataInBuffer = append(dataInBuffer, issueType)
		dataInBuffer = append(dataInBuffer, item["serviceLabel"].(string))
		dataInBuffer = append(dataInBuffer, item["labeledBug"].(string))
		fmt.Println(dataInBuffer)
		issueList[issueNum] = dataInBuffer
	}
}

func ExportResultToLocalCsv(issueList map[int]interface{}, fp string, appendOption string) error {
	var dataSet [][]string
	csvHeader := []string{"id", "title", "issueLink", "state", "creation", "updatedAt", "votes", "type", "serviceLabel", "labeledBug"}
	isFileExists := true
	_, err := os.Stat(fp)
	if errors.Is(err, os.ErrNotExist) {
		isFileExists = false
	}

	f, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE, 0600)

	if err != nil {
		panic(err)
	}
	defer f.Close()

	if isFileExists {
		if appendOption == "true" {
			csvLinesReader := csv.NewReader(f)
			csvLinesReader.FieldsPerRecord = -1
			csvLines, err := csvLinesReader.ReadAll()
			if err != nil {
				return fmt.Errorf("reading data from existed csv file error: %+v", err)
			}

			csvLines = csvLines[1:]
			for _, row := range csvLines {
				issueNum, err := strconv.Atoi(row[0])
				if err != nil {
					return fmt.Errorf("converting issue num error: %+v", err)
				}
				if issueList[issueNum] == nil {
					issueList[issueNum] = row[1:]
				}
			}
		}
	}

	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("writing data to csv file error: %+v", err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("writing data to csv file error: %+v", err)
	}

	for k, v := range issueList {
		var linesToAdd []string
		linesToAdd = append(linesToAdd, strconv.Itoa(k))
		for _, item := range v.([]string) {
			linesToAdd = append(linesToAdd, item)
		}
		dataSet = append(dataSet, linesToAdd)
	}

	if appendOption == "false" || !isFileExists {
		w := csv.NewWriter(f)
		if err := w.Write(csvHeader); err != nil {
			return fmt.Errorf("writing data to csv file error: %+v", err)
		}
		w.Flush()
	}

	wd := csv.NewWriter(f)
	defer func() {
		wd.Flush()
	}()
	if err := wd.WriteAll(dataSet); err != nil {
		return fmt.Errorf("writing data to csv file %s error: %+v", fp, err)
	}

	return nil
}
