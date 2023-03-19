# ♻️ retry

Tool to execute terminal commands with retries.

## 💡 Idea

```bash
$ try [--limit=3 --delay=100ms] -- curl example.com
```

## 🤼‍♂️ How to

```
Usage: try [flags] -- command

Flags:
      --delay duration   retry delay (default 100ms)
      --limit uint       max retry (default 5)
```
