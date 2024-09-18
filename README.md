# Infrapolicy Concourse resource

The goal of this resource is to perform a check of a generated Terraform Plan against Cycloid infrapolicies (link doc) in order to add more control on what's deployed with Terraform. It has been tested
with the Terraform concourse resource `ljfranklin/terraform-resource` but it can be used with any commands / resources providing a JSON terraform plan.

## Usage

First, declare your new resource type:

```yaml
resource_types:
  - name: cycloid-resource
    type: docker-image
    source:
      repository: cycloid/cycloid-resource
      tag: latest
```

Then configure your resource:

```yaml
resources:

# Infrapolicy resource
  - name: infrapolicy
    type: cycloid-resource
    source:
      feature: infrapolicy
      api_key: <api-key>
      env: ((env))
      org: ((org))
      project: ((project))

# Terracost resource
  - name: terracost
    type: cycloid-resource
    source:
      feature: terracost
      api_key: <api-key>
      env: ((env))
      org: ((org))
      project: ((project))

# Event resource
  - name: event
    type: cycloid-resource
    source:
      feature: event
      api_key: <api-key>
      env: ((env))
      org: ((org))
      project: ((project))
```


## Parameters 

### Source configuration

`feature`: _required_. The name of Cycloid feature to use, `terracost`, `infrapolicy` or `event`

`api_key`: _required_. The Cycloid API key used to authenticate the resource against Cycloid APIs

`project`: _required_. The name of the Cycloid project

`env`: _required_. The environment name of the Cycloid project

`org`: _required_. The organization name of the Cycloid project

`api_url`: _optional_. Override the default API URL for infrapolicy validation 

### Put parameters for `event`

`title`: _required_. The title of the event.

`message` or `message_file`: _required_. One have to be specified, `message` message in the event body or `message_file` file path which contain the message for event body.

`type`: _optional_. The type of the event. Currently, only Cycloid, Custom, AWS or Monitoring are allowed.

`severity`: _optional_. The severity of the event. Currently, only info, warn, err or crit are allowed.

`icon`: _optional_. Icon to display. The icons are the ones from Font Awesome. Example: fa-cubes https://fontawesome.com/search?o=r&m=free&f=classic

`vars_file`: _optional_. Load vars from a file that you can use in event message or title. format MYKEY: value usage my title containing vars $MYKEY.

`tags`: _optional_. The tags allow filtering. Example:


### Put parameters for `terracost`, `infrapolicy`

`tfplan_path`: _required_. The path to the JSON terraform plan result (this should be updated since we know the name of the JSON terraform plan)

## Output files

Used with `get`, the resource will populate one output file:

  * `version.json`: Which contain the same json output provided to Concourse for the version


## Usage

Finally, add the `put` step right after the terraform plan and don't forget to the `output_planfile: true` in order to generate a terraform plan JSON file:

```yaml
# Terracost and infrapolicy
- put: tfstate
  get_params:
    output_planfile: true
  ...
- put: infrapolicy
  params:
    tfplan_path: tfstate/plan.json
- put: terracost
  params:
    tfplan_path: tfstate/plan.json

# Event
- put: event
  params:
    title: "my event"
    message: "This is my message"
```

## Tips - run the resource as a task (Advanced/troubleshooting)

If you need to obtain detailed json file. You can run it as a task to populate the following json files:

  * `output.json`: JSON formatted output used also as stdout
  * `cy-output.json`: Raw json output from Cycloid CLI

```YAML
      - task: cost
        config:
          platform: linux
          image_resource:
            type: registry-image
            source:
              repository: cycloid/cycloid-resource
              tag: dev
          run:
            path: /bin/bash
            args:
              - '-ec'
              - |
                ls
                cp ${src_tfplan_path} /tmp; echo "${resource_config}" > source.json
                /opt/resource/out $PWD/terracost-json/ < source.json
          inputs:
            - name: tfstate
          outputs:
            - name: terracost-json
        params:
          src_tfplan_path: tfstate/plan.json
          resource_config:
            source:
              api_key: ((custom_api-key-admin.key))
              api_url: 'https://http-api.cycloid.io'
              env: dev
              feature: terracost
              org: cycloid-demo
              project: accenture-mi
            params:
              tfplan_path: /tmp/plan.json
```

## Contributing

If you want to contribute or to have more information on the workflow: [CONTRIBUTIING.md](./CONTRIBUTING.md)

