# linux-log-analyzer — Kiro Power / Agent Plugin

A portable [Agent Plugin](https://agent-plugins.org) (v1.0.0) that turns Kiro
into a specialist in **Linux log analysis**. It bundles a reusable Agent Skill
that investigates, correlates, and diagnoses incidents from system, application,
and container logs — read-only and non-destructive.

Derived from the `log-analyzer` custom agent
(`.kiro/agents/log-analyzer.json`) built in LESSON-7.

## What's inside

```text
linux-log-analyzer/
├── plugin.json                         # Agent Plugins manifest (identity + metadata)
├── skills/
│   └── log-analysis/
│       ├── SKILL.md                    # Skill instructions (methodology, sources, output)
│       └── references/
│           └── cheatsheet.md           # Ready-to-use read-only commands per source
├── README.md
└── LICENSE
```

- **plugin.json** — declares the plugin identity and targets Agent Plugins
  `1.0.0`. Closed schema: only `$schema`, `name`, `version`, `description`,
  `author`, `homepage`, `repository`, `license`, `keywords`, `extensions`.
- **skills/log-analysis** — the capability. A compatible client discovers the
  skill by scanning `skills/` and loads `SKILL.md`. The `description` frontmatter
  is what the agent matches against a user request to activate the skill.
- No `mcp.json` is included: this power is pure guidance + skills and relies on
  the client's shell/read tools. Add an `mcp.json` at the plugin root to bundle
  MCP servers (e.g., a Loki or Kubernetes MCP server) if desired.

## Capabilities

- Sources: journald, `/var/log` (syslog, kern, auth/secure, dmesg), nginx/apache,
  Docker/Podman, Kubernetes.
- Tools: `lnav` (multi-file correlation + SQL), `logcli`/Loki (LogQL).
- Methodology: scope → collect → correlate → root cause → remediate.
- Strictly read-only: never mutates logs or touches service/container state.

## Install (Kiro)

Kiro Powers support the Agent Plugins spec natively. Point Kiro at this plugin
directory to install it, then ask Kiro to analyze logs — the `log-analysis`
skill activates on matching requests.

## Extending

- **Add MCP servers**: create `mcp.json` at the plugin root
  (`$schema: https://agent-plugins.org/schemas/1.0.0/mcp.schema.json`).
- **Add skills**: create another `skills/<name>/SKILL.md` (e.g., a
  `postgres-log-analysis` skill).
