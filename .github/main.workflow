workflow "docker" {
  on = "push"
  resolves = ["Docker Registry"]
}

action "Docker Registry" {
  uses = "bjornmagnusson/actions/balena-deployer@balena-entrypoint"
  secrets = ["BALENA_TOKEN"]
  args = "pi1led"
}
