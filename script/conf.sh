Target="os-webui"
Package="github.com/os-webui/os-webui"
Docker="king011/os-webui"
Dir=$(cd "$(dirname $BASH_SOURCE)/.." && pwd)
Version="v0.0.1"
Platforms=(
    darwin/amd64
    windows/amd64
    linux/arm64
    linux/amd64
)