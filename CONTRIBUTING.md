# Contributing

## Concourse workflow

Before contributing, we need to understand the Concourse resource workflow.

A concourse resource is composed by three components:
  * `in`
  * `out`
  * `check`

This three programs must be executable, you can use bash, python or binary the important thing is to be able to understand the input params and to output the correct information. This [doc](https://concourse-ci.org/implementing-resource-types.html) will help you to understand the params expected for each component.

`infrapolicy-resource` will mainly run the `out` program, since it's a `put` in the pipeline definition, the worfklow is the following:

1. put:
2. put will call the `./out` program
3. validation is done using the CLI and source params
4. once the `./out` exits, we fallback on an implicit `get` so we call the `./in` program
5. once the `./in` exits, we call the `check` program

## Testing

If you want to test, you need:
  * [cy](https://github.com/cycloidio/cycloid-cli) in your `$PATH`
  * access to cycloid console through the API
  * a terraform plan JSON output (`terraform plan -out=./plan; terraform show -json ./plan  > plan.json`)

You first need to define a `source` JSON:

```json
{
        "source": {
                "email": "your-email",
                "password": "your-password",
                "env": "your-env",
                "project": "your-project",
                "org": "your-org",
                "api_url": "https://api.staging.cycloid.io"
        },
        "params": {
                "tfplan_path": "/path/to/terraform/plan.json"
        }
}
```

Then you can perform an infrapolicy validation:

```shell
$ // build_id is provided by concourse in build environment
$ export BUILD_ID=1234
$ make
$ ./resource/out . < source.json
{"version":{"build_id":"1234","criticals":"0","warnings":"0","advisories":"0"},"metadata":null}
```
