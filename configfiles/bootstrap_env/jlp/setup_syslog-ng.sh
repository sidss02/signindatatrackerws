#!/bin/bash


ud_file="/tmp/userdata"

## Source the file to get the exported variables
. $ud_file

### Test uf CLOUD_STACK is set.
test -z $CLOUD_STACK && { $logger "Warning: no CLOUD_STACK specified, assuming dev"; CLOUD_STACK="dev"; }

### Boil the stack down to dev or prod for receiver endpoint
if [[ ${ENVIRONMENT_NAME} =~ "prod" ]]; then
  RECEIVER_STACK=prod
else
  RECEIVER_STACK=dev
fi
# Is it in VPC?
# Get MAC address
MAC=$(/usr/bin/curl --silent http://169.254.169.254/latest/meta-data/network/interfaces/macs/ | cut -d / -f 1)
# Curl metadate
if /usr/bin/curl --output /dev/null --silent --head --fail  http://169.254.169.254/latest/meta-data/network/interfaces/macs/"${MAC}"/vpc-id; then
  VPC=1
fi
# Get the region
REGION=`/usr/bin/ec2metadata --availability-zone  | sed 's/[a-z]$//'`
## INstance ID
INSTANCE_ID=$(/usr/bin/ec2metadata --instance-id)


# Dest syslog reciever per address
case "$REGION" in
  us-east*)
    if [ "${VPC}" == "1" ] ; then
      SYSLOG_RECEIVER=syslog-relay-useast.wsroute.mathworks.com

    else
      SYSLOG_RECEIVER=syslog-$RECEIVER_STACK-useast.wsroute.mathworks.com
    fi
    ;;
    eu-west*)
    if [ "${VPC}" == "1" ] ; then
      SYSLOG_RECEIVER=syslog-relay-euwest.wsroute.mathworks.com
    else
      SYSLOG_RECEIVER=syslog-$RECEIVER_STACK-useast.wsroute.mathworks.com
    fi
    ;;
  *)
    if [ "${VPC}" == "1" ] ; then
      SYSLOG_RECEIVER=syslog-relay-useast.wsroute.mathworks.com
    else
    SYSLOG_RECEIVER=syslog-$RECEIVER_STACK-useast.wsroute.mathworks.com
    fi
    logger "BOOTSTRAP ERROR: Unknown region assuming us-east"
    ;;
esac


### Replace syslog variables with real ones
echo "@define cloud_en $CLOUD_ENVIRONMENT" > /etc/syslog-ng/conf.d/00_defines.conf
echo "@define cloud_app $CLOUD_APP" >> /etc/syslog-ng/conf.d/00_defines.conf
echo "@define instance_id $INSTANCE_ID" >> /etc/syslog-ng/conf.d/00_defines.conf
echo "@define cloud_stack $CLOUD_STACK"  >> /etc/syslog-ng/conf.d/00_defines.conf
echo "@define SYSLOG_RECEIVER $SYSLOG_RECEIVER" >> /etc/syslog-ng/conf.d/00_defines.conf


### Reload syslog config to ensure it has latest data.
service syslog-ng restart
