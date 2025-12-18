#!/bin/sh
set -eu

SERVICE_NAME=agent
BIN_NAME=agent
BIN_SRC=./agent
BIN_DST=/usr/local/bin/agent
UNIT_FILE=/etc/systemd/system/agent.service
SERVICE_USER=agent

echo "[1/6] Checking binary"
if [ ! -x "$BIN_SRC" ]; then
    echo "ERROR: ./agent binary not found or not executable"
    exit 1
fi

echo "[2/6] Installing binary"
sudo install -m 0755 "$BIN_SRC" "$BIN_DST"

echo "[3/6] Ensuring system user exists"
if ! id "$SERVICE_USER" >/dev/null 2>&1; then
    sudo useradd \
        --system \
        --no-create-home \
        --shell /usr/sbin/nologin \
        "$SERVICE_USER"
fi

echo "[4/6] Writing systemd unit"
sudo tee "$UNIT_FILE" >/dev/null <<EOF
[Unit]
Description=Agent Service
After=network.target
Wants=network.target

[Service]
Type=simple
ExecStart=$BIN_DST
User=$SERVICE_USER
Group=$SERVICE_USER
Restart=on-failure
RestartSec=2

NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true

StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

echo "[5/6] Reloading systemd"
sudo systemctl daemon-reload

echo "[6/6] Enabling and starting service"
sudo systemctl enable "$SERVICE_NAME"
sudo systemctl restart "$SERVICE_NAME"

echo "Done."
echo "Check logs with:"
echo "  journalctl -u agent.service -f"
