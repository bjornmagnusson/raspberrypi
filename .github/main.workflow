workflow "docker" {
  on = "push"
  resolves = [
    "dinfo",
    "env",
  ]
}

action "dinfo" {
  uses = "actions/docker/cli@master"
  args = "info"
}

action "env" {
  uses = "docker://docker:stable"
  args = "env"
}

action "Deploy app on Balena Cloud" {
  uses = "docker://bjornmagnusson/balena-deployer:canary"
  secrets = ["BALENA_TOKEN"]
  args = "pi1led"
}
