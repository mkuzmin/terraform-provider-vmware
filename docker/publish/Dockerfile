FROM alpine:3.4

ARG TERRAFORM_VERSION=0.8.1
ARG VSPHERE_VERSION=1.0.pre2

RUN apk update
RUN apk add ca-certificates wget
RUN update-ca-certificates
RUN wget -nv https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip -O terraform.zip
RUN unzip terraform.zip -d /usr/local/bin/
RUN rm terraform.zip
RUN wget -nv https://github.com/mkuzmin/terraform-vsphere/releases/download/v${VSPHERE_VERSION}/terraform-provider-vmware.linux -O /usr/local/bin/terraform-provider-vsphere
RUN chmod +x /usr/local/bin/terraform-provider-vsphere
