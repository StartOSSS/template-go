# Terraform

Terraform is split into a reusable module (`modules/app`) plus component folders under
`components/`. Components wire one or more modules together, while environment-specific
values live in `*.tfvars` files (for example `components/app/dev.tfvars`). This layout keeps
module code DRY while letting each environment override project IDs, regions, or database
tiers.

## Usage

```bash
cd terraform/components/app
cp dev.tfvars.example dev.tfvars # customise project/region/etc.
terraform init
terraform plan -var-file=dev.tfvars
```

`make plan` (and the `Terraform` GitHub workflow) iterate through `TF_ENVS` and run
`terraform init -backend=false` plus `plan -refresh=false` from
`terraform/components/$(TF_COMPONENT)` so syntax changes stay safe across environments.
Provide a `GCP_SA_KEY` secret in CI to authenticate the `google` provider.

## Module outputs

The module exposes:

- `service_account_email` – Cloud Run service account for Connect RPC
- `database_instance` – Cloud SQL instance name
- `database_secret` – Secret Manager secret containing the DATABASE_URL
- `database_url` – connection string (sensitive)

You can consume those outputs to wire Cloud Run, Skaffold, or GitHub Actions secrets.
