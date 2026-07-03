#!/usr/bin/env bash
# dev.sh - AuxiTalk Development Environment Launcher
# Usage: ./dev.sh [command]

set -e

CORE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DASHBOARD_DIR="$CORE_DIR/../plugin-dashboard"
OPENAI_DIR="$CORE_DIR/../plugin-openai"

print_header() {
    echo -e "\033[1;34m▶ $1\033[0m"
}

print_success() {
    echo -e "\033[0;32m✓ $1\033[0m"
}

print_error() {
    echo -e "\033[0;31m✗ $1\033[0m"
}

usage() {
    echo "AuxiTalk Development Launcher"
    echo ""
    echo "Commands:"
    echo "  build         Build core + plugins"
    echo "  run           Build and run everything (core + dashboard)"
    echo "  core          Run only auxitalkd with dev config"
    echo "  dashboard     Run only dashboard"
    echo "  clean         Remove build artifacts"
    echo "  help          Show this help"
    echo ""
}

build_core() {
    print_header "Building auxitalk-core..."
    cd "$CORE_DIR"
    go build -o bin/auxitalkd ./cmd/auxitalkd
    print_success "Core built"
}

build_dashboard() {
    print_header "Building dashboard..."
    cd "$DASHBOARD_DIR"
    ./dev.sh build 2>/dev/null || go build -o dashboard ./cmd/dashboard
    print_success "Dashboard built"
}

build_openai() {
    print_header "Building OpenAI plugin..."
    cd "$OPENAI_DIR"
    go build -o bin/plugin-openai ./cmd/plugin
    print_success "OpenAI plugin built"
}

build_all() {
    build_core
    build_dashboard
    build_openai
}

run_all() {
    build_all
    print_header "Starting AuxiTalk development environment..."

    # Kill previous
    pkill -f auxitalkd || true
    pkill -f dashboard || true

    # Start core
    cd "$CORE_DIR"
    mkdir -p bin
    ./bin/auxitalkd --config configs/auxitalk.dev.json &
    CORE_PID=$!
    echo $CORE_PID > /tmp/auxitalk-core.pid
    print_success "auxitalkd started (pid $CORE_PID)"

    sleep 2

    # Start dashboard
    cd "$DASHBOARD_DIR"
    PORT=8080 ./dashboard > dashboard.log 2>&1 &
    DASH_PID=$!
    echo $DASH_PID > /tmp/auxitalk-dashboard.pid
    print_success "Dashboard started on http://localhost:8080 (pid $DASH_PID)"

    echo ""
    print_success "Everything is running!"
    echo ""
    echo "  Dashboard:  http://localhost:8080"
    echo "  Core logs:  tail -f $CORE_DIR/bin/auxitalkd.log (if configured)"
    echo "  Stop all:   ./dev.sh stop"
    echo ""
}

run_core() {
    build_core
    cd "$CORE_DIR"
    ./bin/auxitalkd --config configs/auxitalk.dev.json
}

run_dashboard() {
    cd "$DASHBOARD_DIR"
    ./dev.sh run
}

stop_all() {
    print_header "Stopping all services..."
    pkill -f auxitalkd || true
    pkill -f dashboard || true
    rm -f /tmp/auxitalk-*.pid
    print_success "Stopped"
}

clean() {
    print_header "Cleaning build artifacts..."
    rm -rf "$CORE_DIR/bin"
    rm -rf "$DASHBOARD_DIR/dashboard" "$DASHBOARD_DIR/dashboard.log" "$DASHBOARD_DIR/.dashboard.pid"
    print_success "Cleaned"
}

case "${1:-help}" in
    build)
        build_all
        ;;
    run)
        run_all
        ;;
    core)
        run_core
        ;;
    dashboard)
        run_dashboard
        ;;
    stop)
        stop_all
        ;;
    clean)
        clean
        ;;
    help|*)
        usage
        ;;
esac
