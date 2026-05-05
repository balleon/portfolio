```bash
export AWS_ACCESS_KEY_ID="<REDACTED>"
export AWS_SECRET_ACCESS_KEY="<REDACTED>"
export PULUMI_CONFIG_PASSPHRASE=""

pulumi login --local

pulumi stack init dev
pulumi config set aws:region eu-west-3
pulumi install
pulumi preview
pulumi up
```