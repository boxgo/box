#!/usr/bin/env python3
"""
Extract @box_module sections from docs/prometheus_alerts_template.yaml
for generate_alerts.sh. Writes YAML fragment: indented rule lines (no top-level groups:).
"""
from __future__ import annotations

import re
import sys
from typing import List, Set, Tuple

MOD_LINE = re.compile(r"^(\s*)#\s*@box_module:\s*(\S+)\s*$")

# Canonical id -> set of allowed aliases
ALIASES = {
    "http_server": {"http_server", "httpserver"},
    "http_client": {"http_client", "httpclient", "wukong"},
    "grpc_server": {"grpc_server", "grpc"},
    "db_client": {"db_client", "db", "database", "gorm"},
    "redis": {"redis"},
    "mongodb": {"mongo", "mongodb", "mongo_client"},
    "schedule": {"schedule", "cron"},
    "go_runtime": {"go_runtime", "go", "golang", "runtime"},
}

CANONICAL: dict[str, str] = {}
for c, syns in ALIASES.items():
    for s in syns:
        CANONICAL[s.lower()] = c


def canonicalize(token: str) -> str:
    t = token.strip().lower()
    if t in CANONICAL:
        return CANONICAL[t]
    raise ValueError(
        f"Unknown module: {token!r}. Use one of: {sorted(set(ALIASES))} (or aliases: go, db, grpc, mongo, ...)"
    )


def parse_modules(template_path: str) -> Tuple[List[str], List[Tuple[str, int, int]]]:
    with open(template_path, encoding="utf-8") as f:
        lines = f.readlines()
    blocks: List[Tuple[str, int, int]] = []
    i = 0
    n = len(lines)
    while i < n:
        m = MOD_LINE.match(lines[i])
        if m:
            name = m.group(2).strip()
            start = i
            i += 1
            while i < n and not MOD_LINE.match(lines[i]):
                i += 1
            blocks.append((name, start, i))
        else:
            i += 1
    return lines, blocks


def extract(
    template_path: str, wanted: Set[str], all_modules: bool
) -> str:
    lines, blocks = parse_modules(template_path)
    if not blocks:
        sys.stderr.write("No @box_module sections found in template\n")
        sys.exit(1)
    have = {b[0] for b in blocks}
    if all_modules:
        select = {b[0] for b in blocks}
    else:
        select = set(wanted)
    missing = select - have
    if missing:
        sys.stderr.write(f"Unknown or missing in template: {sorted(missing)}. Have: {sorted(have)}\n")
        sys.exit(1)
    out: List[str] = []
    for name, s, e in blocks:
        if name in select:
            out.extend(lines[s:e])
    if not out:
        sys.stderr.write("No content selected\n")
        sys.exit(1)
    return "".join(out)


def main() -> None:
    if len(sys.argv) < 3:
        print(
            "Usage: extract_alert_modules.py <template.yaml> ALL|module1,module2,...",
            file=sys.stderr,
        )
        sys.exit(1)
    path = sys.argv[1]
    modstr = sys.argv[2]
    if modstr.strip().upper() == "ALL":
        text = extract(path, set(), all_modules=True)
    else:
        wanted = {canonicalize(p.strip()) for p in modstr.split(",") if p.strip()}
        text = extract(path, wanted, all_modules=False)
    sys.stdout.write(text)


if __name__ == "__main__":
    try:
        main()
    except ValueError as e:
        sys.stderr.write(f"{e}\n")
        sys.exit(1)
