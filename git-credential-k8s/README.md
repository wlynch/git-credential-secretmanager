# git-credential-k8s

A Git credential helper for accessing GCP Secret Manager secrets.

## Usage

```sh
$ git config --global credential.https://github.com.helper k8s
$ export GIT_K8S_SECRET="<secret manager version>"
```

## Future ideas

- Determine namespace from Downwards API.
- Scope allowed hosts in Secret labels/annotations.