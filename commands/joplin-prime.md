---
description: 'Entry point for the Joplin CLI — load the full schema and route to the right command'
targets: ["*"]
---

# Prime — Joplin CLI Router

Load the authoritative CLI schema first, then route the user's request to the right slash command or CLI call.

## Step 1 — Load the schema

Run:

```sh
joplin-cli describe
```

The output is JSON: `{"binary","version","commands":[{name, short, long, flags:[{name, shorthand, type, default, description}], output, example}]}`. Every command documents both its **inputs** (flags) and its **output** JSON shape. Keep this in mind as the source of truth for the rest of the conversation — prefer it over `--help`.

If the binary is missing (`command not found: joplin-cli`), tell the user to run `make install` from the repo root (or `make build` then use `./bin/joplin-cli`) and stop.

## Step 2 — Check reachability & auth (when a task needs the API)

- Joplin Desktop must be running with the **Web Clipper service enabled**. A quick check: `joplin-cli ping` → `{"ok":true,"message":"JoplinClipperServer",...}`.
- Commands that touch data need a token, resolved as: `--token` flag → `JOPLIN_TOKEN` env → the credential stored by `login`. If any command errors with `no Joplin token found …`, have the user authenticate: `echo "$JOPLIN_TOKEN" | joplin-cli login` (stored AES-encrypted, machine-bound), or set `JOPLIN_TOKEN`.

## Step 3 — Route the user's request

Match the user's task against the table below. Multiple rows can match — pick the most specific. Rows with a `/slash-command` delegate to a dedicated skill; the rest are run directly against the CLI.

| Task signal | Action |
|---|---|
| "Find / search / list notes matching …", full-text or Joplin filters | `/joplin-search-note` |
| "Create / add a new note or to-do" | `/joplin-add-note` |
| "Edit / update / change / move a note", set or clear a to-do flag | `/joplin-update-note` |
| "Delete / remove / trash a note" | `/joplin-delete-note` |
| "Read / show the full body of a note by id" | `joplin-cli note get <id>` |
| "Import a Markdown file as a note" | `joplin-cli note import <file.md>` |
| "List / create / delete notebooks (folders)", "notes in a notebook" | `joplin-cli notebook {list,create,delete,notes}` |
| "List / create / delete tags", "tag or untag a note", "notes with a tag" | `joplin-cli tag {list,create,delete,add,remove,notes}` |
| "Is Joplin reachable / health check" | `joplin-cli ping` |
| "Log in / save my token / auth" | `joplin-cli login` (token via stdin or `--token`) |
| "Log out / clear stored token" | `joplin-cli logout` |
| "Show me the schema / all commands / what can this do" | `joplin-cli describe` |

## Step 4 — Rules

- For the four note actions, delegate to the slash command — do not duplicate its content here.
- For everything else, drive the CLI directly using the schema from Step 1: read the target command's `flags` and `output` fields to build the call and parse the result. Every command prints JSON to stdout on success; errors go to stderr with a non-zero exit.
- Global flags apply to every command: `--format json|text` (default `json`), `--token`, `--base-url`.
- When resolving a note/notebook/tag by name rather than id, list first (`note search`, `notebook list`, `tag list`) to get the id, and confirm the match before any write or delete.
- If no row matches, ask the user to clarify what they want to do with Joplin.
- The `commands/` directory is the static reference; `joplin-cli describe` is the runtime reference — prefer the latter when they disagree.

## Step 5 — Common pitfalls (read before piping to `jq`)

- **List outputs are wrapped objects, not arrays.** Extract the inner array by name; `jq '.[]'` on these silently returns nothing and looks like "empty":
  - `note search "<q>"` → `{"count","has_more","notes":[…]}` → use `jq '.notes[]'`
  - `notebook list` → `{"count","notebooks":[…]}` → use `jq '.notebooks[]'`
  - `tag list` → `{"count","tags":[…]}` → use `jq '.tags[]'`
  - If a list looks empty, re-check with the raw JSON (no `jq`) before concluding there's no data.
- **`note search` query is positional** — `joplin-cli note search "text"`. The `--query`/`-q` flag is also accepted. There is no other flag name for it.
