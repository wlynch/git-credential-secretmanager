FROM golang:1.13
ENV GO111MODULE=on
WORKDIR /git-credential-k8s
COPY . .
RUN go build .

FROM gcr.io/google.com/cloudsdktool/cloud-sdk:slim
COPY --from=0 /git-credential-k8s/git-credential-k8s /usr/bin/
ENTRYPOINT [ "git" ]