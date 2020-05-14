# Slack Notify - GitHub Action
[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#wip)


A [GitHub Action](https://github.com/features/actions) to send a message to a Slack channel.

**Screenshot**

<img width="485" alt="action-slack-notify-rtcamp" src="https://i.ibb.co/tP48n6n/Screenshot-from-2020-05-13-14-35-10.png">

## Usage

You can use this action after any other action. Here is an example setup of this action:

1. Create a `.github/workflows/slack-notify.yml` file in your GitHub repo.
2. Add the following code to the `slack-notify.yml` file.

```yml
on: push
name: Slack Notification Demo
jobs:
  slackNotification:
    name: Slack Notification
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Slack Notification
      uses: partnerhero/action-slack-notify@v2.1.0
      env:
        SLACK_WEBHOOK: ${{ secrets.SLACK_WEBHOOK }}
```

3. Create `SLACK_WEBHOOK` secret using [GitHub Action's Secret](https://developer.github.com/actions/creating-workflows/storing-secrets). You can [generate a Slack incoming webhook token from here](https://slack.com/apps/A0F7XDUAZ-incoming-webhooks).


## Environment Variables

By default, action is designed to run with minimal configuration but you can alter Slack notification using following environment variables:

Variable          | Default                                               | Purpose
------------------|-------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------
SLACK_CHANNEL     | Set during Slack webhook creation                     | Specify Slack channel in which message needs to be sent
SLACK_USERNAME    | `DevOpsBot`                                           | The name of the sender of the message. Does not need to be a "real" username
SLACK_ICON_EMOJI  | -                                                     | User/Bot icon shown with Slack message, in case you do not wish to add a URL for slack icon as above, you can set slack emoji in this env variable. Example value: `:bell:` or any other valid slack emoji.
SLACK_COLOR       | `good` (green)                                        | You can pass an RGB value like `#efefef` which would change color on left side vertical line of Slack message.
SLACK_MESSAGE     | Generated from git commit message.                    | The main Slack message in attachment.
SLACK_TITLE       | Message                                               | Title to use before main Slack message.
SITE_NAME         | -                                                     | Environment site name
SITE_URL          | -                                                     | Site URL
SHOW_ACTIONS_URL  | true                                                  | Show the actions url field
SHOW_REF          | true                                                  | Show the ref field
SHOW_EVENT        | true                                                  | Show the event that triggered the workflow

You can see the action block with all variables as below:

```yml
    - name: Slack Notification
      uses: partnerhero/action-slack-notify@v2.1.0
      env:
        SLACK_CHANNEL: general
        SLACK_COLOR: '#3278BD'
        SLACK_MESSAGE: 'Post Content :rocket:'
        SLACK_TITLE: Post Title
        SLACK_USERNAME: Devops Bot
        SLACK_WEBHOOK: ${{ secrets.SLACK_WEBHOOK }}
        SITE_NAME: staging
        SITE_URL: http://myapp.staging.com
        SHOW_ACTIONS_URL: true
        SHOW_REF: true
        SHOW_EVENT: true
```

## Hashicorp Vault (Optional)

This GitHub action supports [Hashicorp Vault](https://www.vaultproject.io/). 

To enable Hashicorp Vault support, please define following GitHub secrets:

Variable      | Purpose                                                                       | Example Vaule
--------------|-------------------------------------------------------------------------------|-------------
`VAULT_ADDR`  | [Vault server address](https://www.vaultproject.io/docs/commands/#vault_addr) | `https://example.com:8200`
`VAULT_TOKEN` | [Vault token](https://www.vaultproject.io/docs/concepts/tokens.html)          | `s.gIX5MKov9TUp7iiIqhrP1HgN`

You will need to change `secrets` line in `slack-notify.yml` file to look like below.

```yml
on: push
name: Slack Notification Demo
jobs:
  slackNotification:
    name: Slack Notification
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Slack Notification
      uses: rtCamp/action-slack-notify@v2.1.0
      env:
        VAULT_ADDR: ${{ secrets.VAULT_ADDR }}
        VAULT_TOKEN: ${{ secrets.VAULT_TOKEN }}
```

GitHub action uses `VAULT_TOKEN` to connect to `VAULT_ADDR` to retrieve slack webhook from Vault.

In the Vault, the Slack webhook should be setup as field `webhook` on path `secret/slack`.

## License

[MIT](LICENSE) © 2019 rtCamp

This action was forked from https://github.com/rtCamp/action-slack-notify created by [rtCamp](https://github.com/rtCamp/).

### Made with love and &#127861; @ [PartnerHero](https://partnerhero.com/)

<a
  href="https://partnerhero.com/">
  <img
    src="https://partnerhero.com/static/ph-logo-red-74c934089759f0f99f33e26cd77baf06.svg"
    alt="Outsourcing for Startups. From bootstrap through IPO and beyond. We’ve got you.">
</a>
