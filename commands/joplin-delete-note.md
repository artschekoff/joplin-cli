---
description: 'Delete a Joplin note by id (to trash by default, or permanently)'
targets: ["*"]
---

request = $ARGUMENTS

Extract from `request`:
- **note id** (required) — the id of the note to delete. If the user describes the note by title instead of id, run `/joplin-search-note` first to resolve the id, and **confirm the exact note with the user before deleting** — deletion is destructive.
- **permanent** (optional) — pass `--permanent` only if the user explicitly wants to bypass the trash and delete irreversibly. Default is to move the note to the trash.

Run:

```sh
joplin-cli note delete <note-id> [--permanent]
```

The output is JSON: `{"deleted":true,"id":"<note-id>","permanent":<bool>}`. Confirm to the user which note was deleted and whether it went to the trash or was removed permanently.

If the command errors with `no Joplin token found …`, tell the user to authenticate first: `echo "$JOPLIN_TOKEN" | joplin-cli login` (or set `JOPLIN_TOKEN`). On other errors (e.g. not found), surface the stderr message verbatim.
