---
description: 'Search Joplin notes by full-text query and return matching notes as JSON'
targets: ["*"]
---

request = $ARGUMENTS

Extract from `request`:
- **query** (required) — the search text. Supports Joplin search filters, e.g. `tag:work`, `notebook:Journal`, `updated:day-7`, `title:meeting`.
- **limit** (optional) — max results (Joplin caps at 100). Default 100.

Run:

```sh
joplin-cli note search "<query>" [--limit <N>]
```

The output is JSON: `{"notes":[{"id","parent_id","title","body","created_time","updated_time","is_todo"}],"count","has_more"}`. Timestamps are RFC3339 UTC (or `null`); `is_todo` is a boolean.

Summarize the results for the user: report `count`, and list each note as `id  title` (mention `has_more: true` if the result was truncated so they can raise `--limit`). Keep note ids visible — the update/delete commands need them.

If the command errors with `no Joplin token found …`, tell the user to authenticate first: `echo "$JOPLIN_TOKEN" | joplin-cli login` (or set the `JOPLIN_TOKEN` env var). On other errors, surface the stderr message verbatim.
