#graphql's and compiler's builder/Dockerfile should be in sync
FROM harbor-repo.vmware.com/nexus/golang:1.18

ENV GOPATH /go
ENV PATH $PATH:/usr/local/go/bin:/go/bin
ENV GO111MODULE=on
ENV GOPRIVATE=*.vmware.com
ENV GOPROXY=https://build-artifactory.eng.vmware.com/srp-mds-go-remote,direct
ENV GONOPROXY=*.eng.vmware.com
ENV GOINSECURE=golang-appnet.eng.vmware.com

RUN tdnf -y --refresh install \
  make=4.3-2.ph4 \
  coreutils=8.32-3.ph4 \
  build-essential=0.1-3.ph4 \
  openssh-clients=8.8p1-5.ph4 \
  jq=1.6-2.ph4 \
  && go install github.com/onsi/ginkgo/ginkgo@v1.16.0 \
  && go install github.com/onsi/gomega/...@v1.18.0 \
  && go install golang.org/x/tools/cmd/goimports@latest \
  && go install github.com/mikefarah/yq/v4@latest \
  && go install k8s.io/kube-openapi/cmd/openapi-gen@1a6458611d189dc17e98a0824dc92536365efedf \
  && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go env GOPATH)"/bin v1.49.0
