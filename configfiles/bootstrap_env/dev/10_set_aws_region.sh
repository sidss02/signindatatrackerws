#!/bin/bash

# Attempt to read the ec2 region metadata. 15s timeout, 1 attempt.
aws_region=$(wget -q -O - -T 15 -t 1 http://169.254.169.254/latest/meta-data/placement/region)

if [ -z "$aws_region" ]
then
    aws_region="us-east-1"
fi

sudo /bin/systemctl set-environment aws_region="$aws_region"