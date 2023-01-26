# koyo-cli

CLI for the Koyo screenshot service api.

## Install

```sh
go install github.com/tbdsux/koyo/cli/koyo@latest
```

## Usage

```sh
koyo --api https://your-koyo-app-instance-url.deta.app https://alpha.deta.space
```

### Config management

The CLI creates a default config file which can be configured accordingly at `$HOME/.koyo.yaml`

```yaml
# Default Config
api: ''
driver: playwright
fullPage: false
height: 800
imageType: png
whiteHole: ''
width: 1280
```

#### Update config

The keys from the default config can be updated as such...

```
koyo set config.[key] [value]
```

```
koyo set config.width 1400
```

##

**tbdsux | &copy; 2023**
