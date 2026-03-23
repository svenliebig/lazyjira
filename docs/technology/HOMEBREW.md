# Homebrew

For the application to be available as a homebrew/tap, we use the [svenliebig/homebrew-tap](https://github.com/svenliebig/homebrew-tap) repository. The projects CI/CD workflow is configured to update the homebrew-tap repository with the new release.

For that to work we needed to:

1. Create a Token that's able to write to the homebrew-tap repository [example](https://github.com/settings/tokens/new?scopes=repo&description=HOMEBREW_TAP_GITHUB_TOKEN)
2. Add the token to the GitHub Actions secrets as `HOMEBREW_TAP_GITHUB_TOKEN`: `gh secret set HOMEBREW_TAP_GITHUB_TOKEN --repo yourname/repository`
3. Adjust the CI/CD workflow to use the token in the release.yml workflow.
