# driftwatch

> Monitors infrastructure state files and alerts on unexpected drift between environments.

---

## Installation

```bash
go install github.com/yourorg/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/driftwatch.git && cd driftwatch && go build ./...
```

---

## Usage

Point driftwatch at your state files and define the environments you want to compare:

```bash
driftwatch --config driftwatch.yaml
```

**Example `driftwatch.yaml`:**

```yaml
environments:
  - name: staging
    state_file: ./states/staging.tfstate
  - name: production
    state_file: ./states/production.tfstate

alerts:
  slack_webhook: https://hooks.slack.com/services/your/webhook/url
  threshold: 3
```

Run a one-off drift check:

```bash
driftwatch check --from staging --to production
```

Watch for drift continuously on a schedule:

```bash
driftwatch watch --interval 5m
```

Output example:

```
[DRIFT DETECTED] 3 resource(s) differ between staging and production
  ~ aws_instance.web_server  (ami-id mismatch)
  + aws_s3_bucket.logs       (missing in production)
  - aws_security_group.old   (removed in staging)
```

---

## Contributing

Pull requests and issues are welcome. Please open an issue before submitting large changes.

---

## License

[MIT](LICENSE)