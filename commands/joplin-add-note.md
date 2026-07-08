---
description: 'Create a new Joplin note (optionally a to-do, in a specific notebook)'
targets: ["*"]
---

request = $ARGUMENTS

Extract from `request`:
- **title** (required) — the note title.
- **body** (optional) — the note content, in Markdown.
- **notebook** (optional) — the parent notebook (folder) **id** to file the note under. If the user names a notebook instead of giving an id, run `joplin-cli notebook list` first to resolve the id.
- **todo** (optional) — pass `--todo` when the user wants a to-do item rather than a plain note.

Run:

```sh
joplin-cli note create --title "<title>" [--body "<markdown>"] [--notebook "<folder-id>"] [--todo]
```

The output is JSON: `{"id","parent_id","title","body","created_time","updated_time","is_todo"}`. Report the new note's `id` and `title` back to the user (they'll need the `id` for later edits).

If the command errors with `no Joplin token found …`, tell the user to authenticate first: `echo "$JOPLIN_TOKEN" | joplin-cli login` (or set `JOPLIN_TOKEN`). On other errors, surface the stderr message verbatim.
