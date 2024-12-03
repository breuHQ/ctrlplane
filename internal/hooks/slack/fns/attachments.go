package fns

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"

	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

func LineExceedFields(event *events.Event[eventsv1.RepoHook, eventsv1.Diff]) []slack.AttachmentField {
	fields := []slack.AttachmentField{
		{
			Title: "*Repository*",
			Value: fmt.Sprintf("<%s|%s>", event.Context.Source, ExtractRepoName(event.Context.Source)),
			Short: true,
		}, {
			Title: "*Branch*",
			Value: fmt.Sprintf("<%s/tree/%s|%s>", event.Context.Source, "", ""),
			Short: true,
		}, {
			Title: "*Threshold*",
			Value: fmt.Sprintf("%d", 0),
			Short: true,
		}, {
			Title: "*Total Lines Count*",
			Value: fmt.Sprintf("%d", 0),
			Short: true,
		}, {
			Title: "*Lines Added*",
			Value: fmt.Sprintf("%d", event.Payload.GetLines().GetAdded()),
			Short: true,
		}, {
			Title: "*Lines Deleted*",
			Value: fmt.Sprintf("%d", event.Payload.GetLines().GetRemoved()),
			Short: true,
		}, {
			Title: "Affected Files",
			Value: fmt.Sprintf("%s", FormatFilesList(event.Payload.GetFiles().GetModified())),
			Short: false,
		}, {
			Title: "Rename Files",
			Value: fmt.Sprintf("%s", FormatFilesList(event.Payload.GetFiles().GetRenamed())),
			Short: false,
		},
	}

	return fields
}

func ExtractRepoName(repoURL string) string {
	parts := strings.Split(repoURL, "/")
	return parts[len(parts)-1]
}

func FormatFilesList(files []string) string {
	result := ""
	for _, file := range files {
		result += "- " + file + "\n"
	}

	return result
}