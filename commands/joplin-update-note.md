---
description: 'Update fields on an existing Joplin note by id (only the fields you pass change)'
targets: ["*"]
---

request = $ARGUMENTS

Extract from `request`:
- **note id** (required) — the id of the note to update. If the user describes the note by title instead of id, run `/joplin-search-note` first to resolve the id, and confirm the match before editing.
- **title** (optional) — new title.
- **body** (optional) — new Markdown body. Note: this **replaces** the whole body, it does not append. If the user wants to append, first read the current body with `joplin-cli note get <id>`, combine, then pass the full new body.
- **notebook** (optional) — new parent notebook (folder) id to move the note.
- **todo** (optional) — `--todo` to mark it a to-do, `--todo=false` to clear the to-do flag.

Only pass the flags the user actually wants to change — omitted flags are left untouched.

Run:

```sh
joplin-cli note update <note-id> [--title "<title>"] [--body "<markdown>"] [--notebook "<folder-id>"] [--todo | --todo=false]
```

The output is JSON: `{"id","parent_id","title","body","created_time","updated_time","is_todo"}` — the updated note. Confirm the change back to the user.

If the command errors with `no Joplin token found …`, tell the user to authenticate first: `echo "$JOPLIN_TOKEN" | joplin-cli login` (or set `JOPLIN_TOKEN`). On other errors (e.g. not found), surface the stderr message verbatim.
