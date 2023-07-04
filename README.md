# ♻️ retry

Tool to execute terminal commands with retries.

## 💡 Idea

```bash
$ try [--limit=3 --delay=100ms] -- curl example.com
```

## 🤼‍♂️ How to

```
Usage: try [flags] -- command

flags:
      --delay duration   retry delay (default 100ms)
      --limit uint       max retry, set limit to 0 to disable limit (default 5)
      --quiet            hide command stdout/stderr
```
