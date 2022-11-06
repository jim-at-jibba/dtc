<h1 align="center">Welcome to DevTools CLI 🛠️</h1>
<p>
  <img alt="Version" src="https://img.shields.io/badge/version-0.1-blue.svg?cacheSeconds=2592000" />
  <a href="https://twitter.com/jimgbest" target="_blank">
    <img alt="Twitter: jimgbest" src="https://img.shields.io/twitter/follow/jimgbest.svg?style=social" />
  </a>
</p>

> A collection of tools that you would normally reach for a browser for.

Written in Go and leaning heavily on Cobra and the suite of tools from [Charm](https://charm.sh/):

- [Bubbletea](https://github.com/charmbracelet/bubbletea)
- [Lipgloss](https://github.com/charmbracelet/bubbletea)
- [Bubbles](https://github.com/charmbracelet/bubbles)

<h1>👷 This project is WIP and will have more added to it in the future.</h1>

## Install

```sh
go install github.com/jim-at-jibba/dev-tools-cli@latest
```

## Commands

### UUID

> Generate UUID v4

```bash
dev-tools-cli uuid generate
```

### Base64

> Encode and Decode base64 strings in both standard and URL compatible formats

#### Encode

**Standard**

```bash
dev-tools-cli base64 encode
```

**URL Compatible**

```bash
dev-tools-cli base64 encode -u
```

### File Share

> Ephemeral file sharing, the link provided will expire, after a given time or when the file is downloaded. Makes use of [https://file.io](file.io)

Note that there's a limit of 100mb on files

**defaults to 14 days expiry**

```bash
dev-tools-cli file-share
```

**Pass in an expiry time frame**

```bash
dev-tools-cli file-share

# current dir, expires in 3 days
dev-tools-cli file-share --expiry=3d

# expires in 4 weeks
dev-tools-cli file-share --expiry=4w
```

## Author

👤 **James Best**

- Website: jamesbest.uk
- Twitter: [@jimgbest](https://twitter.com/jimgbest)
- Github: [@jim-at-jibba](https://github.com/jim-at-jibba)

## Show your support

Give a ⭐️ if this project helped you!
