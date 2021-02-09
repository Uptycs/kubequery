# Copyright (c) 2020-present, The kubequery authors
#
# This source code is licensed as defined by the LICENSE file found in the
# root directory of this source tree.
#
# SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)

FROM ubuntu:20.04

ARG BASEQUERY_VERSION=4.6.0
ARG KUBEQUERY_VERSION

LABEL \
  name="kubequery" \
  description="kubequery powered by Osquery" \
  version="${KUBEQUERY_VERSION}" \
  url="https://github.com/Uptycs/kubequery"

ADD https://uptycs-basequery.s3.amazonaws.com/${BASEQUERY_VERSION}/basequery_${BASEQUERY_VERSION}-1.linux_amd64.deb /tmp/basequery.deb

COPY kubequery /usr/local/bin/kubequery.ext

RUN set -ex; \
    DEBIAN_FRONTEND=noninteractive apt-get update -y && \
    DEBIAN_FRONTEND=noninteractive apt-get upgrade -y && \
    DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends curl jq -y && \
    dpkg -i /tmp/basequery.deb && \
    /etc/init.d/osqueryd stop && \
    rm -rf /var/osquery/* /var/log/osquery/* /var/lib/apt/lists/* /var/cache/apt/* /tmp/* && \
    groupadd -g 1000 kubequery && \
    useradd -m -g kubequery -u 1000 -d /opt/kubequery -s /bin/bash kubequery && \
    mkdir /opt/kubequery/var && \
    echo "/usr/local/bin/kubequery.ext" > /opt/kubequery/autoload.exts && \
    chmod 700 /usr/local/bin/kubequery.ext && \
    chown kubequery:kubequery /usr/bin/osquery? /usr/local/bin/kubequery.ext /opt/kubequery/autoload.exts /opt/kubequery/var

COPY entrypoint.sh /opt/kubequery/entrypoint.sh

# NOTE: Not running as root breaks bunch of Osquery tables. But Osquery tables are meaningless
#       in the context of kubequery as the pod is ephemeral in nature
USER kubequery

WORKDIR /opt/kubequery

ENTRYPOINT ["/opt/kubequery/entrypoint.sh"]
