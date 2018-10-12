#!/bin/bash

n=$[1 - $(cat n || echo 1)]
echo ${n} > n

export OPENSHIFT_INSTALL_BASE_DOMAIN="devcluster.openshift.com"
export OPENSHIFT_INSTALL_CLUSTER_NAME="mfojtik-dev-${n}"
export OPENSHIFT_INSTALL_EMAIL_ADDRESS="mfojtik@redhat.com"
export OPENSHIFT_INSTALL_PASSWORD="openshiftdev"
export OPENSHIFT_INSTALL_PLATFORM="aws"
export OPENSHIFT_INSTALL_AWS_REGION="us-east-1"
export OPENSHIFT_INSTALL_SSH_PUB_KEY="$(<~/.ssh/id_rsa.pub)"
export OPENSHIFT_INSTALL_PULL_SECRET="$(<~/.openshift_pull_secret.json)"

# aws ec2 describe-instances --query 'Reservations[*].Instances[*].[Tags[?Key==`Name`].Value|[0],PublicIpAddress]' --output table