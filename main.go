package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	EnvSlackWebhook   = "SLACK_WEBHOOK"
	EnvSlackIcon      = "SLACK_ICON"
	EnvSlackIconEmoji = "SLACK_ICON_EMOJI"
	EnvSlackChannel   = "SLACK_CHANNEL"
	EnvSlackTitle     = "SLACK_TITLE"
	EnvSlackMessage   = "SLACK_MESSAGE"
	EnvSlackColor     = "SLACK_COLOR"
	EnvSlackUserName  = "SLACK_USERNAME"
	EnvGithubActor    = "GITHUB_ACTOR"
	EnvSiteName       = "SITE_NAME"
	EnvSiteURL       = "SITE_URL"
	EnvShowActionsUrl = "SHOW_ACTIONS_URL"
	EnvShowRef = "SHOW_REF"
)

type Webhook struct {
	Text        string       `json:"text,omitempty"`
	UserName    string       `json:"username,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	UnfurlLinks bool         `json:"unfurl_links"`
	Attachments []Attachment `json:"attachments,omitmepty"`
}

type Attachment struct {
	Fallback   string  `json:"fallback"`
	Pretext    string  `json:"pretext,omitempty"`
	Color      string  `json:"color,omitempty"`
	AuthorName string  `json:"author_name,omitempty"`
	AuthorLink string  `json:"author_link,omitempty"`
	AuthorIcon string  `json:"author_icon,omitempty"`
	Footer     string  `json:"footer,omitempty"`
	Fields     []Field `json:"fields,omitempty"`
	
}

type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

func main() {
	endpoint := os.Getenv(EnvSlackWebhook)
	if endpoint == "" {
		fmt.Fprintln(os.Stderr, "URL is required")
		os.Exit(1)
	}
	text := os.Getenv(EnvSlackMessage)
	if text == "" {
		fmt.Fprintln(os.Stderr, "Message is required")
		os.Exit(1)
	}

	fields:= []Field{
		{
			Title: "Event",
			Value: os.Getenv("GITHUB_EVENT_NAME"),
			Short: true,
		},
		{
			Title: os.Getenv(EnvSlackTitle),
			Value: envOr(EnvSlackMessage, "EOM"),
			Short: false,
		},
	}

	showRef := os.Getenv(EnvShowRef)
	if showRef == "true" {
		refField := Field{
			Title: "Ref",
			Value: os.Getenv("GITHUB_REF"),
			Short: true,
		}
		fields = append([]Field { refField }, fields...)
	}

	githubEventName := os.Getenv("GITHUB_EVENT_NAME")
	
	if githubEventName == "pull_request" {
		branchNameField := Field{
			Title: "Branch",
			Value: os.Getenv("GITHUB_BRANCH"),
			Short: true,
		}
		fields = append(fields, branchNameField)
	} else if githubEventName == "push" {
		commitMessageField := Field{
			Title: "Commit Message",
			Value: os.Getenv("COMMIT_MESSAGE"),
			Short: true,
		}
		fields = append(fields, commitMessageField)
	}

	showActionsUrl := os.Getenv(EnvShowActionsUrl)
	if showActionsUrl == "true" {
		actionsUrlField := Field{
			Title: "Actions URL",
			Value: "https://github.com/" + os.Getenv("GITHUB_REPOSITORY") + "/actions/runs/" + os.Getenv("GITHUB_RUN_ID"),
			Short: false,
		}
		fields = append(fields, actionsUrlField)
	}

	siteName := os.Getenv(EnvSiteName)
	if siteName != "" {
		siteNameField := Field{
			Title: "Site",
			Value: os.Getenv(EnvSiteName),
			Short: true,
		}
		fields = append(fields, siteNameField)
	}

	siteUrl := os.Getenv(EnvSiteURL)
	if siteUrl != "" {
		siteUrlField:= Field{
			Title: "Site URL",
			Value: os.Getenv(EnvSiteURL),
			Short: true,
		}
		fields = append(fields, siteUrlField)
	}

	msg := Webhook{
		UserName:  os.Getenv(EnvSlackUserName),
		IconURL:   os.Getenv(EnvSlackIcon),
		IconEmoji: os.Getenv(EnvSlackIconEmoji),
		Channel:   os.Getenv(EnvSlackChannel),
		Attachments: []Attachment{
			{
				Fallback: envOr(EnvSlackMessage, "GITHUB_ACTION=" + os.Getenv("GITHUB_ACTION") + " \n GITHUB_ACTOR=" + os.Getenv("GITHUB_ACTOR") + " \n GITHUB_EVENT_NAME=" + os.Getenv("GITHUB_EVENT_NAME") + " \n GITHUB_REF=" + os.Getenv("GITHUB_REF") + " \n GITHUB_REPOSITORY=" + os.Getenv("GITHUB_REPOSITORY") + " \n GITHUB_WORKFLOW=" + os.Getenv("GITHUB_WORKFLOW")),
				Color:      envOr(EnvSlackColor, "good"),
				AuthorName: envOr(EnvGithubActor, ""),
				AuthorLink: "http://github.com/" + os.Getenv(EnvGithubActor),
				AuthorIcon: "http://github.com/" + os.Getenv(EnvGithubActor) + ".png?size=32",
				Footer: "<https://github.com/partnerhero/action-webhook-slack-notify|Powered By PartnerHero's GitHub Actions Library>",
				Fields: fields,
			},
		},
	}

	if err := send(endpoint, msg); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending message: %s\n", err)
		os.Exit(2)
	}
}

func envOr(name, def string) string {
	if d, ok := os.LookupEnv(name); ok {
		return d
	}
	return def
}

func send(endpoint string, msg Webhook) error {
	enc, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(enc)
	res, err := http.Post(endpoint, "application/json", b)
	if err != nil {
		return err
	}

	if res.StatusCode >= 299 {
		return fmt.Errorf("Error on message: %s\n", res.Status)
	}
	fmt.Println(res.Status)
	return nil
}
