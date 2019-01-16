workflow "docker" {
  on = "push"
  resolves = ["Deploy app on Balena Cloud"]
}

action "dockerizing" {
  uses = "actions/docker/cli@master"
  args = "docker info"
}

action "Deploy app on Balena Cloud" {
  uses = "bjornmagnusson/actions/balena-deployer@balena-entrypoint"
  needs "dockerizing"
  secrets = ["BALENA_TOKEN"]
  args = "pi1led bjornmagnusson/pi-led"
}
