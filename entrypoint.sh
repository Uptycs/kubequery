#!/bin/bash

ADDITIONAL_FLAGS=""

# See if kube-system namespaces can be retrieved
KUBE_SYSTEM=$(curl -sf \
        --cacert /var/run/secrets/kubernetes.io/serviceaccount/ca.crt \
        -H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" \
        -H "Accept: application/json" \
        https://${KUBERNETES_SERVICE_HOST}:${KUBERNETES_SERVICE_PORT}/api/v1/namespaces/kube-system) && {

    # Get UUID of the kube-system
    UUID=$(echo ${KUBE_SYSTEM} | jq -r ".metadata.uid")
    if [ -n "${UUID}" ]; then
        # Use kube-system UUID as the host identifier
        ADDITIONAL_FLAGS="--host_identifier=specified --specified_identifier=${UUID}"
    fi
}

exec /usr/bin/osqueryd \
    ${ADDITIONAL_FLAGS} \
    --flagfile=/opt/kubequery/etc/osquery.flags \
    --config_path=/opt/kubequery/etc/osquery.conf \
    --ephemeral \
    --disable_logging \
    --database_path=/opt/kubequery/osquery.db \
    --enroll_tables=osquery_info,kubernetes_info \
    --extensions_socket=/opt/kubequery/var/osquery.em \
    --extensions_autoload=/opt/kubequery/autoload.exts \
    --extensions_require=kubequery \
    --extension_event_tables=kubernetes_events
