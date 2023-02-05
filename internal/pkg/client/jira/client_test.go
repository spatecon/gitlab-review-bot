package jira

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/andygrunwald/go-jira"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestJiraLibrary(t *testing.T) {
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USERNAME"),
		Password: os.Getenv("JIRA_API_TOKEN"),
	}

	client, err := jira.NewClient(tp.Client(), os.Getenv("JIRA_URL"))
	require.NoError(t, err, "failed to create jira client")

	//b, _, err := client.Board.GetBoard(71)
	//require.NoError(t, err, "failed to get board")
	//
	//t.Log(string(lo.Must(json.MarshalIndent(b, "", " "))))
	//
	//bc, _, err := client.Board.GetBoardConfiguration(71)
	//require.NoError(t, err, "failed to get board configuration")
	//
	//t.Log(string(lo.Must(json.MarshalIndent(bc, "", " "))))
	//
	//sprints, _, err := client.Board.GetAllSprints(strconv.Itoa(b.ID))
	//require.NoError(t, err, "failed to get sprints")
	//
	//t.Log(string(lo.Must(json.MarshalIndent(sprints, "", " "))))

	sprintID := 738

	issues, _, err := client.Sprint.GetIssuesForSprint(sprintID)
	require.NoError(t, err, "failed to get issues for sprint")

	t.Log(string(lo.Must(json.MarshalIndent(issues, "", " "))))
}
