# Infrapolicy Concourse resource

The goal of this resource is to perform a check of a generated Terraform Plan against Cycloid infrapolicies (link doc) in order to add more control on what's deployed with Terraform. It has been tested
with the Terraform concourse resource `ljfranklin/terraform-resource` but it can be used with any commands / resources providing a JSON terraform plan.

## Usage

First, declare your new resource type:

```yaml
resource_types:
  - name: infrapolicy
    type: docker-image
    source:
      repository: cycloid/infrapolicy-resource
      tag: latest
```

Then configure your resource:

```yaml
resources:
  - name: check
    type: infrapolicy
    source:
      email: <user-email>
      password: <user-password>
      env: ((env))
      org: ((org))
      project: ((project))
```

Finally, add the `put` step right after the terraform plan and don't forget to the `output_planfile: true` in order to generate a terraform plan JSON file:

```yaml
- put: tfstate
  get_params:
    output_planfile: true
  ...
- put: check
  params:
    tfplan_path: tfplan/plan.json
```

## Parameters 

### Source configuration

`email`: _required_. The user email address used to authenticate the resource against Cycloid APIs (this should be soon replaced in favor of API keys)

`password`: _required_. The user password used to authenticate the resource against Cycloid APIs (this should be soon replaced in favor of API keys)

`project`: _required_. The name of the Cycloid project

`env`: _required_. The environment name of the Cycloid project

`org`: _required_. The organization name of the Cycloid project

`api_url`: _optional_. Override the default API URL for infrapolicy validation 

### Put parameters

`tfplan_path`: _required_. The path to the JSON terraform plan result (this should be updated since we know the name of the JSON terraform plan)

## Contributing

If you want to contribute or to have more information on the workflow: [CONTRIBUTIING.md](./CONTRIBUTING.md)

