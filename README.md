Dispatch
========

[![Build Status](https://travis-ci.org/jakedahn/dispatch.png?branch=master)](https://travis-ci.org/jakedahn/dispatch)


Dispatch is an experimental system for running small programs in a distributed environment.

## The Idea

The idea behind dispatch is experimental and likely flawed, but potentially interesting.

Dispatch is a service that manages building, packaging, and running small composable services written in golang -- we'll call them "Dispatch Jobs"

The Dispatch Job binary must be be able to run without the ability to pass in command line arguments -- all inputs must be defined/passed-in as environment variables.

The first example of what a Dispatch Job in the wild looks like can be found here: https://github.com/jakedahn/echo


#### The Dispatchfile (aka .dispatch.yml)

The Dispatchfile is meant to do a few things.

1. Provide build/packaging instructions

    * how do I build the go binary?
    * how do I get this thing in a docker image?

2. Specify the interface that the Dispatch Job has

    * this can be used to generate documentation about a Dispatch Job
    * this allows us to verify and validate input parameters to the Dispatch Job (if we make a malformed request, we can immediately return an error so we don't have to wait for the process to crash on the cluster...)

Here is an example of what a Dispatchfile should look like (this is lifted from the [Echo Dispatch Job](https://github.com/jakedahn/echo/blob/master/.dispatch.yml):

```yaml
build:
  - CGO_ENABLED=0 go build -o ./bin/echo -a main.go
  - docker build .

arguments:
  - {key: GOECHO, type: string, presence: required}
```

#### Queuing Jobs

To start with we're going to have an HTTP API that allows us to queue 1 job per request. Shortly after things are working it probably makes sense to allow batch job queuing. When you make a request to this API all it is going to do reach into redis and run the RPUSH command to push onto the right side of a queued jobs list.

The way we can reference a Dispatch Job is by its git repo url, and the sha that we are trying to run. Eventually it would be cool to be able to reference branch name (like master), but for now a git sha is good enough.


```

POST /v0/jobs/queue
BODY: {
    "repo_url": "https://github.com/jakedahn/echo.git",
    "ref": "b2f8590e29173333482016681b7f06f32474a3cd",
    "arguments": {
        "GOECHO": "wheeeeeeee"
    }
}

```

#### Running Jobs

To run a job we will dequeue it form redis by running LPOP to pop (dequeue) the job off the left side of queued jobs list.

Once we have the information about the job we want to run, and the environment variables that we need to build the environment to run the go binary we will use [fleetctl](http://github.com/coreos/fleet) to schedule the job to be run inside of a container somewhere in the CoreOS cluster.


## Project Principles

* Community: If a newbie has a bad time, it's a bug.
* Software: Make it work, then make it right, then make it fast.
* Technology: If it doesn't do a thing today, we can make it do it tomorrow.

(I saw these in the [logstash readme](https://github.com/elasticsearch/logstash/blob/master/README.md#project-principles), and thought it was a nice gesture to have in this readme)


## License

https://tldrlegal.com/license/apache-license-2.0-(apache-2.0)
