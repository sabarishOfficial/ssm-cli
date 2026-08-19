# ssm-cli

A terminal tool to interactively select an AWS EC2 instance and connect to it via AWS Systems Manager (SSM) Session Manager — no SSH keys or open ports required.

## How it works

1. Loads AWS credentials from your environment or `~/.aws/credentials`
2. Lists all EC2 instances in the given region
3. Opens a fuzzy finder (`fzf`) so you can search and select an instance
4. Starts an SSM session to the selected instance using `session-manager-plugin`

```
┌─────────────┐     DescribeInstances     ┌─────────┐
│   ssm-cli   │ ─────────────────────────▶│   EC2   │
│             │                            └─────────┘
│             │     StartSession           ┌─────────┐
│             │ ─────────────────────────▶│   SSM   │
│             │                            └─────────┘
│             │     WebSocket tunnel
│             │ ◀────────────────────────▶ session-manager-plugin
└─────────────┘
```

---

## Prerequisites

| Requirement | Install |
|---|---|
| Go 1.21+ | https://go.dev/dl |
| AWS CLI v2 | https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html |
| session-manager-plugin | https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html |
| AWS credentials configured | `aws configure` or env vars |

### macOS install via Homebrew

```bash
brew install awscli
brew install --cask session-manager-plugin
```

### Linux install (Debian/Ubuntu)

```bash
# AWS CLI v2
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip && sudo ./aws/install

# session-manager-plugin
curl "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/ubuntu_64bit/session-manager-plugin.deb" -o "session-manager-plugin.deb"
sudo dpkg -i session-manager-plugin.deb
```

### Required IAM permissions

Your AWS identity needs at minimum:

```json
{
  "Effect": "Allow",
  "Action": [
    "ec2:DescribeInstances",
    "ssm:StartSession",
    "ssm:TerminateSession",
    "ssm:DescribeSessions"
  ],
  "Resource": "*"
}
```

The EC2 instance must have:
- The **AmazonSSMManagedInstanceCore** IAM policy attached to its instance role
- The SSM Agent running (pre-installed on Amazon Linux 2, Ubuntu 20.04+, and Windows AMIs)

---

## Installation

### Run directly

```bash
git clone https://github.com/your-username/ssm-cli.git
cd ssm-cli
go run main.go [region]
```

### Build a binary

```bash
# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o ssm-cli-darwin-arm64 .

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o ssm-cli-darwin-amd64 .

# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o ssm-cli-linux-amd64 .

# Linux (arm64)
GOOS=linux GOARCH=arm64 go build -o ssm-cli-linux-arm64 .
```

Move the binary to your PATH:

```bash
mv ssm-cli-darwin-arm64 /usr/local/bin/ssm-cli
chmod +x /usr/local/bin/ssm-cli
```

---

## Usage

```
ssm-cli [region]
```

`region` is optional. Defaults to `us-east-1` if not provided.

### Examples

```bash
# Connect to an instance in us-east-2
ssm-cli us-east-2

# Connect to an instance in the default region (us-east-1)
ssm-cli

# Using go run
go run main.go eu-west-1
```

---

## Example output

### Successful connection

```
$ ssm-cli us-east-2

AWS Region: us-east-2
AWS configuration loaded successfully.

> i-0a1b2c3d4e5f  web-server-prod
  i-0f9e8d7c6b5a  bastion
  i-0123456789ab  worker-node-1

  3/3
─────────────────────────────

Selected EC2 Instance: i-0a1b2c3d4e5f

Starting session with SessionId: user-abc123xyz456 ...

sh-4.2$
```

Once connected you land in a shell on the remote instance. Type `exit` to end the session.

```
sh-4.2$ whoami
ssm-user
sh-4.2$ exit

Exiting session with sessionId: user-abc123xyz456.
Connected to EC2 Instance: i-0a1b2c3d4e5f
```

---

## Error reference

### No region provided → defaults to us-east-1

```
$ ssm-cli
AWS Region: us-east-1
AWS configuration loaded successfully.
```

### Invalid / missing AWS credentials

```
Failed to set up AWS config: failed to refresh cached credentials, ...
exit status 1
```

Fix: run `aws configure` or export credentials:

```bash
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
export AWS_SESSION_TOKEN=...   # if using temporary credentials / SSO
```

### No EC2 instances found

The fuzzy finder opens empty. Press `Esc` or `Ctrl+C` to exit:

```
Failed to select EC2 instance: user cancelled
exit status 1
```

Fix: verify the region is correct and your IAM role has `ec2:DescribeInstances`.

### Instance not managed by SSM

```
Failed to connect to EC2 instance: operation error SSM: StartSession,
https response error StatusCode: 400, api error TargetNotConnected:
i-0a1b2c3d4e5f is not connected to Systems Manager
exit status 1
```

Fix checklist:
- Instance has an IAM role with **AmazonSSMManagedInstanceCore**
- SSM Agent is running: `sudo systemctl status amazon-ssm-agent`
- Instance has outbound HTTPS (port 443) to `ssm.<region>.amazonaws.com`

### `session-manager-plugin` not installed

```
Failed to connect to EC2 instance: exec: "session-manager-plugin":
executable file not found in $PATH
exit status 1
```

Fix: install `session-manager-plugin` (see Prerequisites above).

### Permission denied (IAM)

```
Failed to connect to EC2 instance: operation error SSM: StartSession,
https response error StatusCode: 403, api error AccessDeniedException:
User is not authorized to perform: ssm:StartSession
exit status 1
```

Fix: attach the required IAM permissions to your user/role (see Prerequisites above).

---

## Project structure

```
ssm-cli/
├── main.go               # Entry point: wires region → config → list → select → connect
├── go.mod
├── go.sum
├── awsConfig/
│   └── main.go           # Loads AWS SDK config for the given region
└── cmd/
    ├── ec2/
    │   └── main.go       # DescribeInstances → returns IDs and Names
    ├── fzf/
    │   └── main.go       # Interactive fuzzy instance picker (go-fzf)
    └── ssmConnect/
        └── main.go       # StartSession + invokes session-manager-plugin
```
