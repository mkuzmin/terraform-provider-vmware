FROM alpine:3.4

ARG TERRAFORM_VERSION=0.8.1
ARG VSPHERE_VERSION=1.1

RUN apk update && \
    apk add ca-certificates wget git && \
    update-ca-certificates && \
    wget -nv https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip -O terraform.zip && \
    unzip terraform.zip -d /usr/local/bin/ && \
    rm terraform.zip && \
    wget -nv https://github.com/mkuzmin/terraform-vsphere/releases/download/v${VSPHERE_VERSION}/terraform-provider-vmware.linux -O /usr/local/bin/terraform-provider-vmware && \
    chmod +x /usr/local/bin/terraform-provider-vmware

ENTRYPOINT ["terraform"]
