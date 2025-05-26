---
title: "Deploying Applications on Kubernetes"
description: ""
date: 2025-05-26
---

!TODO: This article doesn't explain Kubernetes basics, we are only going to explain
configuration specifications.

When it comes to deploying Applications on Kubernetes, we usually deal with two
most important resources. `Deployment` and `Service`.

## Deployment

Please check the deployment file (yaml) in the current directory to understand
what an example deployment file.

One important thing is to configure `readiness` and `liveness` probes. For which
we added separate `/health` route.
