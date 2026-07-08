# joplin-cli

A console utility for [Joplin](https://joplinapp.org/), rewritten in Go from the
`dweigend/joplin-mcp-server` MCP. It reads and writes notes, notebooks, and tags
through the Joplin **Web Clipper REST API**, and emits JSON for machine (LLM)
consumption.

## Requirements

- Joplin Desktop with the **Web Clipper Service enabled** (Tools → Options → Web Clipper).
- Go 1.24+ to build.

## Install

```bash
make build            # -> ./bin/joplin-cli
make install          # -> /usr/local/bin/joplin-cli (uses sudo)
```

## Authentication

The token comes from (in priority order): `--token`, `JOPLIN_TOKEN`, or the
encrypted store written by `login`.

```bash
# copy the token from Joplin: Tools -> Options -> Web Clipper -> Advanced
echo "your_joplin_web_clipper_token_here" | joplin-cli login
joplin-cli ping        # -> {"ok":true,"message":"JoplinClipperServer",...}
```

## The schema (for agents)

Every command documents its inputs **and** output JSON shape. Ask the binary:

```bash
joplin-cli describe | jq '.commands[] | {name, output}'
```

Or read a single command's rich help:

```bash
joplin-cli note create --help
```

## Commands

| Command | Purpose |
|---------|---------|
| `note search <query> [--limit]`        | Full-text search |
| `note get <id>`                        | Read one note |
| `note create --title [--body --notebook --todo]` | Create a note |
| `note update <id> [--title --body --notebook --todo]` | Update a note |
| `note delete <id> [--permanent]`       | Trash / delete a note |
| `note import <file.md>`                | Import a Markdown file as a note |
| `notebook list`                        | List notebooks |
| `notebook create --title [--parent]`   | Create a notebook |
| `notebook delete <id>`                 | Delete a notebook |
| `notebook notes <id> [--limit]`        | Notes in a notebook |
| `tag list`                             | List tags |
| `tag create --title`                   | Create a tag |
| `tag delete <id>`                      | Delete a tag |
| `tag add <tag-id> <note-id>`           | Attach a tag to a note |
| `tag remove <tag-id> <note-id>`        | Detach a tag from a note |
| `tag notes <id> [--limit]`             | Notes carrying a tag |
| `ping`                                 | Health check |
| `login` / `logout`                     | Store / remove the token |
| `describe`                             | Machine-readable schema (JSON) |

## Global flags

- `--format json|text` — output format (default `json`).
- `--token STRING` — override the token for one call.
- `--base-url STRING` — override `http://localhost:41184`.

## Environment variables

| Variable | Default | Meaning |
|----------|---------|---------|
| `JOPLIN_TOKEN` | — | Web Clipper token |
| `JOPLIN_BASE_URL` | `http://localhost:41184` | API base URL |
| `JOPLIN_TIMEOUT_SECONDS` | `30` | HTTP timeout |
| `JOPLIN_HTTP_RETRIES` | `3` | Retry attempts for transient failures |
| `JOPLIN_HTTP_RETRY_BACKOFF_SECONDS` | `1.0` | Linear backoff base |

## Output contract

Success payloads are pretty JSON on **stdout**. Errors print to **stderr** and
the process exits non-zero — there is no `status` field to inspect. Timestamps
are RFC3339 UTC (or `null`); `is_todo` is a boolean.
