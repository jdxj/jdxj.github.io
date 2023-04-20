---
title: "terminal prompts disabled"
date: 2023-01-04T21:55:38+08:00
tags:
  - git
  - git config
---

[解决could not read Username for 'https://github.com': terminal prompts disabled](https://stackoverflow.com/questions/32232655/go-get-results-in-terminal-prompts-disabled-error-for-github-private-repo)

```bash
$ git config --global --add url."git@github.com:".insteadOf "https://github.com/"
```
