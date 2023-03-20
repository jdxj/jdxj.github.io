---
title: "Config"
date: 2023-01-04T21:55:38+08:00
---

# could not read Username for 'https://github.com': terminal prompts disabled

[解决](could not read Username for 'https://github.com': terminal prompts disabled)

```shell
$ git config --global --add url."git@github.com:".insteadOf "https://github.com/"
```