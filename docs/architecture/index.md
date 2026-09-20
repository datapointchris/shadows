# Architecture

How a shadow file moves between the work repo, shadow storage and another machine, and which of
that is built. Read the overview before changing a package, and the roadmap before starting a
feature.

## Documentation

- [Architecture Overview](OVERVIEW.md) — Storage layout under `~/.shadows/`, the data flow through add, sync and promote, the package-by-package breakdown, and the decisions behind SQLite and `.git/info/exclude`
- [Development Roadmap](ROADMAP.md) — The eight phases from MVP to TUI, the feature checklist inside each, the backlog, and where the work currently stands
