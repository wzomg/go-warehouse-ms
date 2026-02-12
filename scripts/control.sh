#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_DIR="$ROOT_DIR/bin"
RUN_DIR="$ROOT_DIR/run"
BIN="$BIN_DIR/server"
PID_FILE="$RUN_DIR/server.pid"
LOG_FILE="$RUN_DIR/server.log"
COMPOSE_FILE="$ROOT_DIR/docker-compose.yml"

OS_NAME="$(uname -s)"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    return 1
  fi
  return 0
}

install_with_brew() {
  if ! require_cmd brew; then
    echo "缺少 brew，无法自动安装 $1"
    exit 1
  fi
  brew install "$1"
}

install_with_brew_cask() {
  if ! require_cmd brew; then
    echo "缺少 brew，无法自动安装 $1"
    exit 1
  fi
  brew install --cask "$1"
}

install_with_apt() {
  if ! require_cmd apt-get; then
    echo "缺少 apt-get，无法自动安装 $1"
    exit 1
  fi
  if require_cmd sudo; then
    sudo apt-get update
    sudo apt-get install -y "$@"
  else
    apt-get update
    apt-get install -y "$@"
  fi
}

ensure_go() {
  if require_cmd go; then
    return 0
  fi
  if [[ "$OS_NAME" == "Darwin" ]]; then
    install_with_brew go
  else
    install_with_apt golang
  fi
}

ensure_docker() {
  if require_cmd docker; then
    return 0
  fi
  if [[ "$OS_NAME" == "Darwin" ]]; then
    install_with_brew_cask docker
    echo "已安装 Docker Desktop，请手动启动后重试"
    exit 1
  else
    install_with_apt docker.io
  fi
}

ensure_compose() {
  if docker compose version >/dev/null 2>&1; then
    return 0
  fi
  if require_cmd docker-compose; then
    return 0
  fi
  if [[ "$OS_NAME" == "Darwin" ]]; then
    install_with_brew docker-compose
  else
    install_with_apt docker-compose-plugin
  fi
}

ensure_env() {
  ensure_go
  ensure_docker
  ensure_compose
}

start() {
  ensure_env
  mkdir -p "$BIN_DIR" "$RUN_DIR"
  if [[ -f "$PID_FILE" ]] && kill -0 "$(cat "$PID_FILE")" 2>/dev/null; then
    echo "服务已在运行"
    exit 0
  fi
  (cd "$ROOT_DIR" && docker compose -f "$COMPOSE_FILE" up -d db)
  for i in {1..30}; do
    if docker exec go_warehouse_db mysqladmin ping -uroot -proot --silent >/dev/null 2>&1; then
      break
    fi
    sleep 1
  done
  (cd "$ROOT_DIR" && go build -o "$BIN" ./cmd/server)
  nohup "$BIN" >"$LOG_FILE" 2>&1 &
  echo $! >"$PID_FILE"
  echo "服务已启动"
}

stop() {
  ensure_env
  if [[ ! -f "$PID_FILE" ]]; then
    echo "服务未运行"
    (cd "$ROOT_DIR" && docker compose -f "$COMPOSE_FILE" stop db)
    exit 0
  fi
  PID="$(cat "$PID_FILE")"
  if kill -0 "$PID" 2>/dev/null; then
    kill "$PID"
    sleep 1
  fi
  rm -f "$PID_FILE"
  (cd "$ROOT_DIR" && docker compose -f "$COMPOSE_FILE" stop db)
  echo "服务已停止"
}

status() {
  ensure_env
  if [[ -f "$PID_FILE" ]] && kill -0 "$(cat "$PID_FILE")" 2>/dev/null; then
    echo "服务运行中"
  else
    echo "服务未运行"
  fi
  if docker ps --format '{{.Names}}' | grep -q '^go_warehouse_db$'; then
    echo "数据库运行中"
  else
    echo "数据库未运行"
  fi
}

restart() {
  ensure_env
  stop
  start
}

case "${1:-}" in
  start) start ;;
  stop) stop ;;
  restart) restart ;;
  status) status ;;
  *)
    echo "用法: $0 {start|stop|restart|status}"
    exit 1
    ;;
esac
