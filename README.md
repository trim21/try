# ♻️ retry

Tool to execute terminal commands with retries.

## 💡 Idea

```bash
$ try [options] -- curl example.com
```

## 🤼‍♂️ How to

```
Usage: try [flags] -- command

flags:
      --delay duration       retry delay (default 100ms)
      --delay-type string    delay type, can 'fixed' / 'backoff' / 'off' (default "fixed")
      --limit uint           max retry, set limit to 0 to disable limit (default 5)
      --max-delay duration   max retry delay when using non-fixed delay type (default 1s)
      --quiet                hide command stdout/stderr
```
