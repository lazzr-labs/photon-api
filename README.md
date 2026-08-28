# API

## Brew

**commands**

    brew update
    brew upgrade
    brew install act
    brew install git
    brew install k3d
    brew install asdf
    brew install helm
    brew install argocd
    brew install opentofu
    brew install ariga/tap/atlas
    brew install google-cloud-sdk

## ASDF

**install**

    asdf plugin add golang

    asdf install golang 1.27.0

**commands**

    asdf list

## Go

**commands**

    go run api.go --env=[env] --server=[server]

    go build
    go run api-go --env=[env] --server=[server]

    go get -u
    go mod download
    go mod tidy

## Atlas

**commands**

    atlas schema inspect --env [env]
    atlas migrate status --env [env]
    atlas migrate diff --env [env]
    atlas migrate apply --env [env]

## ENT

**commands**

    go run entgo.io/ent/cmd/ent new [name]

    go generate ./ent
