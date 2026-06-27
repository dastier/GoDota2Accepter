[![Go CI](https://github.com/dastier/GoDota2Accepter/actions/workflows/go22.yml/badge.svg)](https://github.com/dastier/GoDota2Accepter/actions/workflows/go22.yml)

# GoDota2Accepter

GoDota2Accepter is a tiny tray app for Dota 2 players who are tired of missing the match-ready prompt while alt-tabbed, watching a stream, tweaking builds, or waiting through a long queue.

Turn it on, keep playing around on your desktop, and let it watch for the Dota 2 notification. When the game-ready or unpause notification appears, it brings Dota 2 forward and presses Enter for you.

<p>
  <img src="assets/d_disabled.svg" alt="Disabled tray icon" width="48">
  <img src="assets/d_enabled.svg" alt="Enabled tray icon" width="48">
</p>

## Why Gamers Use It

- Stop missing ready checks when Dota 2 is not the focused window.
- Keep browsing, chatting, or watching videos during queue time.
- Avoid losing a queue because you stepped away for a few seconds.
- Toggle the listener from a simple system tray menu.
- Keep control: the app only acts while listening is enabled.

## How It Works

The app listens to Linux desktop notifications over DBus and looks for Dota 2 messages such as:

- `Your game is ready`
- `The game is unpausing...`

When one of those messages is detected, GoDota2Accepter searches for the Dota 2 window, activates it, maximizes it, waits briefly, and sends an Enter keypress.

## Tray Icons

The included SVG assets are used as tray-state icons:

- `assets/d_disabled.svg`: listener is disabled.
- `assets/d_enabled.svg`: listener is enabled and watching DBus notifications.

The icons are embedded into the binary, so the app does not need asset files beside it at runtime.

## Install From a Release

The easiest way to use GoDota2Accepter is to download a prebuilt Linux binary from the [GitHub Releases](https://github.com/dastier/GoDota2Accepter/releases) page.

Download the latest `GoDota2Accepter-v0.1.0-linux-amd64.tar.gz` and `.sha256` files, replacing `v0.1.0` with the release version you downloaded:

```bash
sha256sum -c GoDota2Accepter-v0.1.0-linux-amd64.tar.gz.sha256
tar -xzf GoDota2Accepter-v0.1.0-linux-amd64.tar.gz
./GoDota2Accepter
```

The release binary is built for Linux amd64. If your desktop is missing tray, X11, or AppIndicator runtime libraries, install the matching packages from your distribution or build from source with the dependencies below.

To install it for your user account:

```bash
mkdir -p ~/.local/bin
install -m 0755 GoDota2Accepter ~/.local/bin/GoDota2Accepter
```

Make sure `~/.local/bin` is in your `PATH`, then run:

```bash
GoDota2Accepter
```

## Runtime Requirements

- Linux desktop session with DBus notifications.
- X11 or XWayland-compatible desktop session.
- Dota 2 installed and running.
- System tray/AppIndicator support.

## Build Requirements

- Go 1.26 or newer.
- Native build dependencies used by `systray` and `robotgo`.

On Ubuntu, the project CI installs:

```bash
sudo apt-get update
sudo apt-get install -y \
  gcc \
  gir1.2-appindicator3-0.1 \
  libc6-dev \
  libappindicator3-dev \
  libpng++-dev \
  libx11-dev \
  libx11-xcb-dev \
  libxcb-xkb-dev \
  libxkbcommon-dev \
  libxkbcommon-x11-dev \
  libxtst-dev \
  x11-xkb-utils \
  xcb \
  xorg-dev
```

## Build From Source

```bash
go mod tidy
go build -v ./...
```

## Run From Source

```bash
go run .
```

Use the tray icon to enable or disable listening. Quit from the same tray menu when you are done.

## Publishing a Release

Maintainers can publish binaries by pushing a version tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release workflow builds the Linux amd64 binary, packages it as a `.tar.gz`, generates a `.sha256` checksum, and publishes both files to the matching GitHub Release. The archive includes the binary, README, MIT license, and tray icon assets.

## Security and Behavior

- The app listens to session DBus notification messages while the listener is enabled.
- The app only reacts to notification text matching Dota 2 game-ready or unpause messages.
- When a matching notification is detected, the app activates the Dota 2 window and sends one Enter keypress.
- The app does not modify Dota 2 files, game memory, Steam files, credentials, or account data.
- The app does not store notification data.
- The app does not make network requests during normal use. The tray menu can open this GitHub project page in your browser if you click it.

## Notes

GoDota2Accepter is not affiliated with Valve or Dota 2. Use it responsibly and make sure it fits the rules of the games and platforms you play on.

## License

GoDota2Accepter is released under the MIT License. See [LICENSE](LICENSE).
