workflow "docker" {
  on = "push"
  resolves = ["Deploy app on Balena Cloud"]
}

action "Deploy app on Balena Cloud" {
  uses = "bjornmagnusson/actions/balena-deployer@balena-entrypoint"
  secrets = ["BALENA_TOKEN"]
  args = "pi1led bjornmagnusson/pi-led"
}
