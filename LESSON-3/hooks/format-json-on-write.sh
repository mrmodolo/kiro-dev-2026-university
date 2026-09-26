#!/usr/bin/env bash
#
# postToolUse hook: validate + format JSON after a write
#
# Fires after the `write` tool runs (matcher: "write"). If the written file is
# a .json, it validates the JSON and pretty-formats it in place.
#
# Event shape (postToolUse for the write tool), via STDIN:
#   {
#     "hook_event_name": "postToolUse",
#     "cwd": "...",
#     "tool_name": "fs_write",
#     "tool_input": { "command": "create|strReplace|insert", "path": "...", ... },
#     "tool_response": { "success": true, ... }
#   }
#
# TIMING (documented): postToolUse runs AFTER the file is already written and
# CANNOT block the write. So this hook is reactive: it warns on invalid JSON
# (STDERR + non-zero exit -> shown to user) and reformats valid JSON on disk.
#
# Exit codes: 0 = ok (silent). Non-zero = warning shown to user.

set -uo pipefail

EVENT="$(cat)"

# jq is required to parse the event and validate/format JSON.
if ! command -v jq >/dev/null 2>&1; then
  echo "json-format hook: 'jq' not found; skipping validation/format." >&2
  exit 1
fi

FILE_PATH="$(printf '%s' "$EVENT" | jq -r '.tool_input.path // empty')"

# Nothing to do if we can't determine the path.
[ -z "$FILE_PATH" ] && exit 0

# Only act on .json files.
case "$FILE_PATH" in
  *.json) ;;
  *) exit 0 ;;
esac

# File must exist (write already happened).
if [ ! -f "$FILE_PATH" ]; then
  exit 0
fi

# Validate JSON.
if ! jq empty "$FILE_PATH" >/dev/null 2>&1; then
  ERR="$(jq empty "$FILE_PATH" 2>&1 || true)"
  {
    echo "Invalid JSON written to: $FILE_PATH"
    echo "  $ERR"
    echo "  (postToolUse runs after the write and cannot block it; please fix the file.)"
  } >&2
  exit 1
fi

# Valid JSON: pretty-format in place (2-space indent) via safe temp + move.
TMP="$(mktemp)"
if jq --indent 2 '.' "$FILE_PATH" > "$TMP" 2>/dev/null; then
  # Only overwrite if content actually changed (avoids needless writes).
  if ! cmp -s "$TMP" "$FILE_PATH"; then
    mv "$TMP" "$FILE_PATH"
    echo "Formatted JSON: $FILE_PATH"   # STDOUT (exit 0) -> model context only
  else
    rm -f "$TMP"
  fi
else
  rm -f "$TMP"
fi

exit 0
