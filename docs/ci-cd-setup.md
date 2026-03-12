# CI/CD Setup: GitHub Actions → ECR → EKS

How to set up automated build and deploy for a repo that pushes Docker images
to ECR and deploys to EKS via Kustomize overlays.

## Architecture

```
merge to main/master
  → GitHub Actions builds clean images (dev-<sha>)
  → pushes to ECR
  → generates kustomization.yaml with new tags
  → commits overlay back to repo
  → kubectl apply + rollout restart on EKS
```

Authentication uses GitHub OIDC federation — no long-lived AWS keys.
The IAM role is scoped to a single repo/branch and a single EKS namespace.

## Quick Setup (automated)

```bash
./bin/setup-ci-cd \
  --repo tokonoma-ai/microservices-demo \
  --branch master \
  --cluster tokonoma \
  --namespace sock-shop \
  --region us-west-2
```

First time in an AWS account, add `--setup-oidc` to register the GitHub OIDC provider.

Use `--dry-run` to see what commands would run without executing them.

## What the Script Does

### 1. GitHub OIDC Provider (one-time per AWS account)

Registers GitHub as an identity provider so IAM roles can trust GitHub Actions tokens.

```bash
aws iam create-open-id-connect-provider \
  --url https://token.actions.githubusercontent.com \
  --client-id-list sts.amazonaws.com \
  --thumbprint-list 6938fd4d98bab03faadb97b34396831e3780aea1
```

### 2. IAM Role with OIDC Trust (per repo)

Creates a role that only GitHub Actions from a specific repo/branch can assume.

```bash
aws iam create-role \
  --role-name github-actions-tokonoma-ai-microservices-demo \
  --assume-role-policy-document '{
    "Version": "2012-10-17",
    "Statement": [{
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::<ACCOUNT>:oidc-provider/token.actions.githubusercontent.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
        },
        "StringLike": {
          "token.actions.githubusercontent.com:sub": "repo:<OWNER>/<REPO>:ref:refs/heads/<BRANCH>"
        }
      }
    }]
  }'
```

### 3. ECR Push Permissions (per role)

```bash
aws iam attach-role-policy \
  --role-name github-actions-tokonoma-ai-microservices-demo \
  --policy-arn arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser
```

### 4. EKS Describe Permission (per role + cluster)

Needed for `aws eks update-kubeconfig` in CI.

```bash
aws iam put-role-policy \
  --role-name github-actions-tokonoma-ai-microservices-demo \
  --policy-name eks-deploy \
  --policy-document '{
    "Version": "2012-10-17",
    "Statement": [{
      "Effect": "Allow",
      "Action": ["eks:DescribeCluster"],
      "Resource": "arn:aws:eks:<REGION>:<ACCOUNT>:cluster/<CLUSTER>"
    }]
  }'
```

### 5. EKS Access Entry (per role + cluster + namespace)

Maps the IAM role to Kubernetes RBAC with edit access scoped to a namespace.

```bash
aws eks create-access-entry \
  --cluster-name <CLUSTER> \
  --principal-arn arn:aws:iam::<ACCOUNT>:role/<ROLE_NAME> \
  --type STANDARD

aws eks associate-access-policy \
  --cluster-name <CLUSTER> \
  --principal-arn arn:aws:iam::<ACCOUNT>:role/<ROLE_NAME> \
  --policy-arn arn:aws:eks::aws:cluster-access-policy/AmazonEKSEditPolicy \
  --access-scope type=namespace,namespaces=<NAMESPACE>
```

### 6. GitHub Secret (per repo)

```bash
gh secret set AWS_ROLE_ARN \
  --repo <OWNER>/<REPO> \
  --body "arn:aws:iam::<ACCOUNT>:role/<ROLE_NAME>"
```

## GitHub Actions Workflow Template

Add this as `.github/workflows/build-and-update.yaml`:

```yaml
name: build-and-update-overlay

on:
  push:
    branches: [master]  # or main
    paths-ignore:
      - 'deploy/kubernetes/overlays/*/kustomization.yaml'

permissions:
  id-token: write
  contents: write

jobs:
  build-and-update:
    runs-on: ubuntu-latest
    if: github.actor != 'github-actions[bot]'

    steps:
    - uses: actions/checkout@v4

    - uses: docker/setup-buildx-action@v3

    - name: Configure AWS credentials
      uses: aws-actions/configure-aws-credentials@v4
      with:
        role-to-assume: ${{ secrets.AWS_ROLE_ARN }}
        aws-region: us-west-2

    - name: Log in to ECR
      uses: aws-actions/amazon-ecr-login@v2

    - name: Build and push images
      run: ./bin/build --eks

    - name: Generate overlay
      run: ./bin/deploy --eks --no-apply

    - name: Commit updated overlay
      run: |
        git config user.name "github-actions[bot]"
        git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
        git add deploy/kubernetes/overlays/eks/kustomization.yaml
        if git diff --cached --quiet; then
          echo "No overlay changes to commit."
        else
          TAG=$(python3 -c "import json; print(json.load(open('.build-manifest.json'))['tag'])")
          git commit -m "deploy: update eks overlay to ${TAG} [skip ci]"
          git push
        fi

    - name: Configure kubectl
      run: aws eks update-kubeconfig --name <CLUSTER> --region us-west-2

    - name: Deploy to EKS
      run: |
        kubectl apply -k deploy/kubernetes/overlays/eks
        for dep in <DEPLOYMENT_NAMES>; do
          kubectl rollout restart deployment/"${dep}" -n <NAMESPACE>
        done
        for dep in <DEPLOYMENT_NAMES>; do
          echo "Waiting for ${dep}..."
          kubectl rollout status deployment/"${dep}" -n <NAMESPACE> --timeout=300s || \
            echo "WARNING: ${dep} did not become ready within 300s"
        done
```

Replace `<CLUSTER>`, `<NAMESPACE>`, and `<DEPLOYMENT_NAMES>` for your repo.

## Requirements

| Requirement | Scope |
|---|---|
| `bin/build` writes `.build-manifest.json` with `dev-<sha>` tags | per repo |
| `bin/deploy --no-apply` generates overlay without kubectl | per repo |
| Kustomize overlay with `newTag` field | per repo |
| GitHub OIDC provider in AWS account | one-time per account |
| IAM role with OIDC trust | per repo |
| ECR push policy | per role |
| EKS describe policy + access entry | per role + cluster |
| `AWS_ROLE_ARN` secret | per repo |

## Verifying the Setup

Check the IAM role exists and has the right trust:
```bash
aws iam get-role --role-name <ROLE_NAME>
```

Check attached policies:
```bash
aws iam list-attached-role-policies --role-name <ROLE_NAME>
aws iam get-role-policy --role-name <ROLE_NAME> --policy-name eks-deploy
```

Check EKS access:
```bash
aws eks list-access-entries --cluster-name <CLUSTER>
aws eks list-associated-access-policies --cluster-name <CLUSTER> \
  --principal-arn arn:aws:iam::<ACCOUNT>:role/<ROLE_NAME>
```

Check workflow runs:
```bash
gh run list --repo <OWNER>/<REPO> --workflow build-and-update.yaml --limit 5
gh run view <RUN_ID> --log-failed   # if a run fails
```

Check deployed images:
```bash
kubectl get deployments -n <NAMESPACE> \
  -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.template.spec.containers[0].image}{"\n"}{end}' \
  | column -t
```

## Infinite Loop Prevention

The workflow commits back to the repo, which could re-trigger itself. Three safeguards:

1. `paths-ignore` — overlay-only changes don't trigger the workflow
2. `if: github.actor != 'github-actions[bot]'` — skips bot-authored pushes
3. `[skip ci]` in commit message — belt-and-suspenders fallback
