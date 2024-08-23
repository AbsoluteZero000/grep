set -e

(
  cd "$(dirname "$0")" # Ensure compile steps are run within the repository directory
  go build -o /tmp/codecrafters-build-grep-go cmd/mygrep/main.go
)

exec /tmp/codecrafters-build-grep-go "$@"
