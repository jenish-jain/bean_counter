#!/bin/bash

# follow https://github.com/google-github-actions/auth#setting-up-workload-identity-federation
# else run the below script after changing the top 3 variables

# replace with your project ID
export PROJECT_ID="realm-asgard"

# replace with your repo name you are creation deployment action on
export REPO="jenish-jain/bean_counter"

export SERVICE_ACCOUNT_ID="cloud-function-deployer"

gcloud iam service-accounts create "${SERVICE_ACCOUNT_ID}" \
  --project "${PROJECT_ID}"

gcloud iam workload-identity-pools create "github-actions-pool" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --display-name="Github actions pool"


WPID=$(gcloud iam workload-identity-pools providers describe "github-actions" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --workload-identity-pool="github-actions-pool" \
  --format="value(name)")

export WORKLOAD_IDENTITY_POOL_ID=$WPID


gcloud iam workload-identity-pools providers create-oidc "github-actions" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --workload-identity-pool="github-actions-pool" \
  --display-name="Github actions provider" \
  --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" \
  --issuer-uri="https://token.actions.githubusercontent.com"

gcloud iam service-accounts add-iam-policy-binding "${SERVICE_ACCOUNT_ID}@${PROJECT_ID}.iam.gserviceaccount.com" \
  --project="${PROJECT_ID}" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/${WORKLOAD_IDENTITY_POOL_ID}/attribute.repository/${REPO}"

