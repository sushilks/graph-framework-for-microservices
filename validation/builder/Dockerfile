FROM harbor-repo.vmware.com/nexus/golang:1.19

ENV GOPATH /go
ENV PATH $PATH:/usr/local/go/bin:/go/bin
ENV GO111MODULE=on

RUN tdnf -y --refresh install \
  make=4.3-2.ph4 \
  coreutils=8.32-3.ph4 \
  build-essential=0.1-3.ph4 \
  openssh-clients \
  jq=1.6-2.ph4 \
  && go install github.com/onsi/ginkgo/ginkgo@v1.16.0 \
  && go install github.com/onsi/gomega/...@v1.17.0 \
  && go install golang.org/x/tools/cmd/goimports@latest \
  && go install github.com/mikefarah/yq/v4@latest \
  && go install k8s.io/kube-openapi/cmd/openapi-gen@1a6458611d189dc17e98a0824dc92536365efedf \
  && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go env GOPATH)"/bin v1.49.0