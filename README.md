# TenK Backend

## Setup

### Environment

1. Create a `.env.development` at the root level
2. Go to our discord server, look for `secret` channel in category `Development`, find and copy the `.env.development` content
3. Paste that content into the `.env.development`

### Air for Live Reload

Check if your machine has installed Air already or not

```sh
air -v
```

- [Install Air on your local machine](https://github.com/cosmtrek/air)

### Make

Check if your machine has installed Make already or not

```sh
make -v
```

- [Install Make for Windows](https://stackoverflow.com/questions/32127524/how-to-install-and-use-make-in-windows)
- [Install Make for MacOS](https://stackoverflow.com/questions/10265742/how-to-install-make-and-gcc-on-a-mac)

### VSCode tools

1. Install `Go` extension
2. Install `gofumpt` in your local machine: `go install mvdan.cc/gofumpt@latest`
3. Press `Ctrl/Cmd + Shift + P`, search `Preferences: Open User Settings (JSON)` and paste these line at the end of the json

```json
"go.useLanguageServer": true,
"gopls": {
 "formatting.gofumpt": true,
},
```

## Development

To run a specific service with Air

```sh
# Run core service with air
make core

# Run analytics service with air
make analytics
```
