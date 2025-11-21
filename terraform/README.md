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

Checked-in `dev.tfvars`, `preprod.tfvars`, and `prod.tfvars` point at `sandbox-project-tc` so
you can run `make plan` and review all three environments in one command. The
`Terraform Plan` GitHub workflow mirrors the same tfvars, aggregates the plans, and posts
them on every pull request touching `terraform/**`.

`make plan` (and the workflow above) iterate through `TF_ENVS` and run
`terraform init -backend=false` plus `plan -refresh=false` from
`terraform/components/$(TF_COMPONENT)` so syntax changes stay safe across environments.
Provide a `GCP_SA_KEY` secret in CI to authenticate the `google` provider.

## Secrets

The module generates a unique database password per environment via the
`random_password` resource. Terraform writes the password to a dedicated Secret Manager
secret (`<env>-todo-database-password`) and combines it with the Cloud SQL host/user to
produce a standard `DATABASE_URL` secret (`<env>-todo-database-url`). Cloud Run revisions
mount the `DATABASE_URL` secret as an environment variable during deployment and the
service account receives `roles/secretmanager.secretAccessor` automatically.

## Module outputs

The module exposes:

- `service_account_email` – Cloud Run service account for Connect RPC
- `database_instance` – Cloud SQL instance name
- `database_secret` – Secret Manager secret containing the DATABASE_URL
- `database_password_secret` – Secret Manager secret containing the raw DB password
- `database_url` – connection string (sensitive)

You can consume those outputs to wire Cloud Run, Skaffold, or GitHub Actions secrets.
