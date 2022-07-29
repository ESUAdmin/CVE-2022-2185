# CVE-2022-2185
wo ee cve-2022-2185 gitlab authenticated rce

read: https://starlabs.sg/blog/2022/07-gitlab-project-import-rce-analysis-cve-2022-2185/

## how to use

First spawn a gitlab instance. Log in, create a group and project with a unique name. Create an access token.

Edit these lines in main.go and compile it:

```go
const importProjectName = "projectwtf"
const runCmd = "/bin/sleep inf"
const proxyTo = "http://localhost:8000/"
```

This mitm runs on `*:8100`. Expose it to the Internet.

Log in to target server. Navigate to create a group - import and enter the local server details. When you are going to import, intercept the request and change it:

- `source_type` to `project_entity`
- `source_full_path` to `your_group/your_project`
- if `destination_namespace` is empty, change it to any non-empty name
- `destination_name` is not empty by design

Pass the modified request to server. Wait 255s to get rce.

Note: the command may be run multiple times.
