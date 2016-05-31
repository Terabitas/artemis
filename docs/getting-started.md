# How this works ?

`artemis` is HA service running on `DigitalOcean`. This service allows you over API to create `Auto Scaling Group` 
configuration which defines what is your `Desired` number of healthy droplets, what is `HealthyThreshold`, which DO `Snapshot`,  
`Region`, `Size`, `SSHKey` should be used when droplet has to be created. Once this configuration is created `artemis` 
will start monitoring health of your droplets and when droplet goes down it will spin up a new one.

# How do you know that droplet is unhealthy ?

Every droplet will have to run this daemon `https://github.com/nildev/watcher`. It is dead simple program which takes data 
from DO metadata service and reports it back to `artemis` service.

Currently it just pings back to `artemis`, but future plans are that this watcher will be reporting `CPU`, `Memory` metrics so that
`artemis` could scale based on it. For actual data collection i am planning to use https://github.com/firehol/netdata.git.

# What is `HealthyThreshold` param ?

It can happen that `watcher` deployed on droplet can fail register itself 1 of 5 times, but that does not mean that droplet is unhealthy.
`artemis` default policy takes metrics of last `CheckInterval` amount and calculates average. If after three consecutive checks 
result is lower than `HealthyThreshold` then `artemis` will mark instance as unhealthy and will initiate `relaunch` action. 
`Relaunch` action will first of all launch new instance, will wait until it will become healthy and only then will try to terminate old 
unhealthy droplet.

For example setting `HealthyThreshold` at 0.7 and `CheckInterval` at 5 seconds it will allow your watcher to fail register it self 1 of 5 times.

# `artemisctl` client

This is cli client that communicates over API with `artemis` service. Download it from release section.

# How to create ASG ?

With this one liner you can create `ASG` which will make sure you always have at least one droplet up and running:

```
artemisctl asg create \ 
  		--name "my-test-asg" \
        --desired 1 \
        --check-interval 5 \
        --consecutive-checks 3 \
        --healthy-threshold 0.7 \
  		--api-key "your-do-api-key-which-has-write-perms" \
  		--region "ams3" \
  		--size "1gb" \
  		--image "img-id-to-be-used" \
  		--ssh-key "finger-print-of-ssh-key-to-use"
```
