package main

import "time"

// Borrowed from default pod termination grace period of Kubernetes:
// https://kubernetes.io/docs/reference/kubernetes-api/core/pod-v1/#PodSpec
const terminationGracePeriod = 30 * time.Second
