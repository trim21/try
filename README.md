# ♻️ retry

Tool to execute terminal commands with retries.

## 💡 Idea

```bash
$ retry --limit=3 -- curl example.com
```

## 🤼‍♂️ How to

[![asciicast](https://asciinema.org/a/150367.png)](https://asciinema.org/a/150367)

```
Usage: retry --limit=N -- command

The strategy flags
    -limit=N
        Limit creates a Strategy that limits the number of attempts that Retry will make.
        If N<=0, default value 3 will be used.

Examples:
    retry --limit=3 -- curl http://example.com
```
