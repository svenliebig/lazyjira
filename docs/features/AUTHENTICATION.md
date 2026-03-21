# Authentication

The user need to authenticate to Jira Cloud to use the tool. There are multiple ways to authenticate and we will elaborate them in this document.

## Necessary Information

We need two information from the user to authenticate:

- Jira Cloud URL
- API Token

## Authentication Methods

The application checks multiple locations for the authentication information in the following order:

1. Command Line Arguments
2. Environment Variables
3. Configuration File

If not authenticated, the application will prompt the user for the authentication information and save it to the configuration file.

### Command Line Arguments

We can use command line arguments to authenticate. We will use the following command line arguments:

- `--jira-cloud-url`
- `--jira-api-token`

### Environment Variables

We can use environment variables to authenticate. We will use the following environment variables:

- `JIRA_CLOUD_URL`
- `JIRA_API_TOKEN`

### Configuration File

We can use a configuration file to authenticate. We will use the following configuration file:

- `$XDG_CONFIG_HOME/jira-cli/config.json`

```json
{
  "jiraCloudUrl": "https://your-jira-cloud-url.com",
  "jiraApiToken": "your-api-token"
}
```
