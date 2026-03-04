## Installation

```bash
brew tap iam-bkpl/tap
brew install time-pugcha
```

### macOS Security Warning
If you see a Gatekeeper warning, run:
```bash
xattr -d com.apple.quarantine $(which time-pugcha)
```

---

## Usuage

With Flag:
```bash
time-pugcha -t=5:30 || time-pugcha -t 5:30
```

Without Flag:
```bash
time-pugcha 5:30
```

Help:
```bash
time-pugcha -h
```
