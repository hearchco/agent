# Hearchco agent repository built using Go

## Installation

### Docker
[https://github.com/hearchco/agent/pkgs/container/agent](https://github.com/hearchco/agent/pkgs/container/agent)

```bash
docker pull ghcr.io/hearchco/agent
```

### Binary
<details>
    <summary>Binary file - Linux</summary>

Download the latest release from the [releases page](https://github.com/hearchco/agent/releases) manually, or automatically like below and set the permissions for the files.

```bash
# Replace the 'match' part with your own ARCH
curl -L -o /opt/hearchco <<< echo $(curl -sL https://api.github.com/repos/hearchco/agent/releases/latest | jq -r '.assets[] | select(.name? | match("linux_amd64$")) | .browser_download_url')
```

### Create a user and modify the rights.

```bash
sudo useradd --shell /bin/bash --system --user-group hearchco
sudo chown hearchco:hearchco /opt/hearchco
```

## Start/Stop/Status

### Create a Systemd Unit

Save example systemd unit file into `/etc/systemd/system/hearchco.service` [docs](../docs/hearchco.service).

### Start the hearchco Service

Reload the service daemon, start the newly create service and check status.

```bash
sudo systemctl daemon-reload
sudo systemctl start hearchco
sudo systemctl status hearchco
```

### Debug

```bash
sudo journalctl -u hearchco -b --reverse
```

### Start hearchco on Startup

```bash
sudo systemctl enable hearchco.service
```

</details>
