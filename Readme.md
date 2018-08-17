# Github Webhook Tester

## Introduction

This repo should help you test github webhooks

## Contents

- [Install](#install)

## Install

```bash
brew cask install ngrok
```

## Run

In one terminal:

```bash
go run main.go
```

In another terminal:

```bash
ngrok http 4444
```

Update the settings in your repo so the webhook points to the ngrok URL
