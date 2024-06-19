<h1 align="center">
  <img height="180px" src="https://docs.raito.io/assets/images/raito-logo-half.png" alt="Raito" />
</h1>

<h4 align="center">
  AWS Organization plugin for the Raito CLI
</h4>

<p align="center">
    <a href="/LICENSE.md" target="_blank"><img src="https://img.shields.io/badge/license-Apache%202-brightgreen.svg?label=License" alt="Software License" /></a>
    <img src="https://img.shields.io/github/v/release/raito-io/cli-plugin-aws-organization?sort=semver&label=Release&color=651FFF" />
    <a href="https://github.com/raito-io/cli-plugin-aws-organization/actions/workflows/build.yml" target="_blank"><img src="https://img.shields.io/github/actions/workflow/status/raito-io/cli-plugin-aws-organization/build.yml?branch=main" alt="Build status" /></a>
    <a href="https://codecov.io/gh/raito-io/cli-plugin-aws-organization" target="_blank"><img src="https://img.shields.io/codecov/c/github/raito-io/cli-plugin-aws-organization?label=Coverage" alt="Code Coverage" /></a>
    <a href="https://golang.org/"><img src="https://img.shields.io/github/go-mod/go-version/raito-io/cli-plugin-aws-organization?color=7fd5ea" /></a>
</p>

<hr/>

# Raito CLI Plugin - AWS

**Note: This repository is still in an early stage of development.
At this point, no contributions are accepted to the project yet.**

This Raito CLI plugin is used to synchronize the identity store information of an AWS organization (IAM Identity Center).
This identity store can then be linked to AWS Account data sources (or set as Master Identity Store) so that permission sets can be visualized correctly.


## Prerequisites
To use this plugin, you will need

1. The Raito CLI to be correctly installed. You can check out our [documentation](http://docs.raito.io/docs/cli/installation) for help on this.
2. A Raito Cloud account to synchronize your AWS organization with. If you don't have this yet, visit our webpage at (https://www.raito.io/trial) and request a trial account.
3. Access to your AWS environment. Minimal required permissions still need to be defined. Right now we assume that you're set up with one of the default
authentication options: https://docs.aws.amazon.com/sdkref/latest/guide/standardized-credentials.html#credentialProviderChain. 

A full example on how to start using Raito Cloud with Snowflake can be found as a [guide in our documentation](http://docs.raito.io/docs/guide/cloud).

## Usage
To use the plugin, add the following snippet to your Raito CLI configuration file (`raito.yml`, by default) under the `targets` section:

```yaml
- name: aws-account
  connector-name: raito-io/cli-plugin-aws-organization
  identity-store-id: "{{IDENTITYSTORE_ID}}"

  aws-account-id: "{{AWS_ACCOUNT_ID}}"
```

Next, replace the values of the indicated fields with your specific values, or use [environment variables](https://docs.raito.io/docs/cli/configuration):
- `identity-store-id`: The ID of the Identity Store you created in Raito Cloud.
- `aws-profile` (optional): The AWS SDK profile to use for connecting to your AWS account that managed the organization to synchronize. When not specified, the default profile is used (or what is defined in the AWS_PROFILE environment variable).
- `aws-region` (optional): The AWS region to use for connecting to the AWS account to synchronize. When not specified, the default region as found by the AWS SDK is used.

You will also need to configure the Raito CLI further to connect to your Raito Cloud account, if that's not set up yet.
A full guide on how to configure the Raito CLI can be found on (http://docs.raito.io/docs/cli/configuration).

## Trying it out

As a first step, you can check if the CLI finds this plugin correctly. In a command-line terminal, execute the following command:
```bash
$> raito info raito-io/cli-plugin-aws-organization
```

This will download the latest version of the plugin (if you don't have it yet) and output the name and version of the plugin, together with all the plugin-specific parameters to configure it.

When you are ready to try out the synchronization for the first time, execute:
```bash
$> raito run
```
This will take the configuration from the `raito.yml` file (in the current working directory) and start a single synchronization.

Note: if you have multiple targets configured in your configuration file, you can run only this target by adding `--only-targets aws-account` at the end of the command.

## Authentication
To authenticate the AWS plugin, the AWS default provider chain will be used:
1. Environment variables: The environment variables `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, and `AWS_SESSION_TOKEN` are used.
2. Shared credentials file. Credentials defined on `~/.aws/credentials` will be used. A profile can be defined with `aws-region`.
3. If running on an Amazon EC2 instance, IAM role for Amazon EC2.

More information can be found on the [AWS SDK documentation](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials).