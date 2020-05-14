#!/usr/bin/env bash

export GITHUB_BRANCH=${GITHUB_HEAD_REF##*heads/}
export SLACK_USERNAME=${SLACK_USERNAME:-"DevOps Bot"}
export CI_SCRIPT_OPTIONS="ci_script_options"
export SLACK_TITLE=${SLACK_TITLE:-"Message"}
export SITE_NAME=${SITE_NAME:-""}
export SITE_URL=${SITE_URL:-""}
export SHOW_ACTIONS_URL=${SHOW_ACTIONS_URL:-"true"}
export SHOW_REF=${SHOW_REF:-"true"}
export COMMIT_MESSAGE=$(cat "/github/workflow/event.json" | jq .commits | jq '.[0].message' -r)

hosts_file="$GITHUB_WORKSPACE/.github/hosts.yml"

if [[ -z "$SLACK_CHANNEL" ]]; then
	if [[ -f "$hosts_file" ]]; then
		user_slack_channel=$(cat "$hosts_file" | shyaml get-value "$CI_SCRIPT_OPTIONS.slack-channel" | tr '[:upper:]' '[:lower:]')
	fi
fi

if [[ -n "$user_slack_channel" ]]; then
	export SLACK_CHANNEL="$user_slack_channel"
fi

# Login to vault using GH Token
if [[ -n "$VAULT_GITHUB_TOKEN" ]]; then
	unset VAULT_TOKEN
	vault login -method=github token="$VAULT_GITHUB_TOKEN" > /dev/null
fi

if [[ -n "$VAULT_GITHUB_TOKEN" ]] || [[ -n "$VAULT_TOKEN" ]]; then
	export SLACK_WEBHOOK=$(vault read -field=webhook secret/slack)
fi

if [[ -f "$hosts_file" ]]; then
	hostname=$(cat "$hosts_file" | shyaml get-value "$GITHUB_BRANCH.hostname")
	user=$(cat "$hosts_file" | shyaml get-value "$GITHUB_BRANCH.user")
	export HOST_NAME="\`$user@$hostname\`"
	export DEPLOY_PATH=$(cat "$hosts_file" | shyaml get-value "$GITHUB_BRANCH.deploy_path")

	temp_url=${DEPLOY_PATH%%/app*}
	export SITE_NAME="${temp_url##*sites/}"
    export HOST_TITLE="SSH Host"
fi

k8s_site_hostname="$GITHUB_WORKSPACE/.github/kubernetes/hostname.txt"

if [[ -f "$k8s_site_hostname" ]]; then
    export SITE_NAME="$(cat $k8s_site_hostname)"
    export HOST_NAME="\`$CLUSTER_NAME\`"
    export HOST_TITLE="Cluster"
fi


if [[ -z "$SLACK_MESSAGE" ]]; then
	export SLACK_MESSAGE="$COMMIT_MESSAGE"
fi

webhook-slack-notify "$@"
