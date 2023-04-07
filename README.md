# go-test-integration
[![Go Report Card](https://goreportcard.com/badge/github.com/zev-zakaryan/go-test-integration)](https://goreportcard.com/report/github.com/zev-zakaryan/go-test-integration)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/zev-zakaryan/go-test-integration/blob/main/LICENSE)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=zev-zakaryan_go-test-integration&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=zev-zakaryan_go-test-integration)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=zev-zakaryan_go-test-integration&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=zev-zakaryan_go-test-integration)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=zev-zakaryan_go-test-integration&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=zev-zakaryan_go-test-integration)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=zev-zakaryan_go-test-integration&metric=coverage)](https://sonarcloud.io/summary/new_code?id=zev-zakaryan_go-test-integration)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=zev-zakaryan_go-test-integration&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=zev-zakaryan_go-test-integration)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=zev-zakaryan_go-test-integration&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=zev-zakaryan_go-test-integration)

1. Create new project from empty or another project as template (not clone)
2. Change project name related string in README.md, sonar-project.properties and language related files such as go.mod
3. We use "bound" sonarcloud, require import project in sonarcloud (analyze action will fail after copy repo+workflow) https://sonarcloud.io/projects/create
	Don't create SONAR_TOKEN yourself but get it from auto-created after import project. Add that SONAR_TOKEN in the project.
4. GitHub actions cannot trigger another GitHub action. Use personal access token so the PRs are opened by you thereby triggering the required workflows.
https://github.com/googleapis/release-please/issues/922
https://github.com/marketplace/actions/release-please-action#github-credentials
Add REPO_ACTION_TOKEN: ...token from... https://github.com/settings/tokens
5. Re-run failed jobs. Wait until completely success
6. Settings > Branches > Add branches rule on main branch as needed:
	Require status checks: SonarCloudChk[Github Actions],GosecChk[Github Actions],
  SonarCloud[SonarCloud],gosec[Github Code Scanning],SonarCloud Code Analysis[Github Code Scanning]
7. Settings > Actions > "Allow GitHub Actions to create and approve pull requests" to use "release-please" action.
8. Merge
