Thank you for considering making contributions to the Cosmos modules repositories!

### Adding a Module:
To get your module added to this repo make sure you meet the following requirements:

1. Your module has a unique name from any other module in the `cosmos/modules` repository
2. Your module contains a `LICENSE` AND `README`
3. Your module has thorough documentation in a `docs/` subfolder
4. Your module has its own `go.mod` and `go.sum` files
5. Your module has unit tests along with a `simapp.go` and `sim_test.go` for fuzz testing
6. Your module contains a `module.yaml` file containing the following fields:

```yaml
name: my_module_name
version: 1.0.0
sdk_versions: 
- 0.32.0
- 0.33.0
- 0.34.0
tm_versions: 
- 0.28.0
- 0.29.0
- 0.30.0
deprecated: {true|false}
description: short description of what the module does
owners:
- @AdityaSripal
- @asripal
website: (optional url)
audits:
- optional link to audit report 1
- optional link to audit report 2
keywords:
- optional keyword 1
- optional keyword 2...
```

### Module Specific Issues/PRs

To make an issue or pull request to a specific module (e.g. `modules/poa`), prefix the issue/pr name with the module in question. For example, we can create a PR to update SDK version in poa module like so, `poa: update SDK version`.
