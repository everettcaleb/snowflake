# Snowflake API
ID generation service inspired by Twitter's [Snowflake](https://github.com/twitter-archive/snowflake/tree/b3f6a3c6ca8e1b6847baa6ff42bf72201e2c2231)

[![Build Status](https://travis-ci.com/everettcaleb/snowflake.svg?branch=master)](https://travis-ci.com/everettcaleb/snowflake)
[![Coverage Status](https://coveralls.io/repos/github/everettcaleb/snowflake/badge.svg?branch=master)](https://coveralls.io/github/everettcaleb/snowflake?branch=master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/7c8e46edb5444b29b21cb1b9b2cbe25e)](https://www.codacy.com/app/everettcaleb/snowflake?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=everettcaleb/snowflake&amp;utm_campaign=Badge_Grade)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](LICENSE)

## Using with Docker
You can run the following command to run it locally from [Docker Hub](https://hub.docker.com/r/everettcaleb/snowflake/):

    docker run -d --name snowflake -e REDIS_URI="redis://host:port/" -p 8080:8080 everettcaleb/snowflake

Then you can test it with:

    curl http://localhost:8080/id

## Using with Kubernetes
You can deploy it as a [Deployment](k8s/deployment.yaml) and [Service](k8s/service.yaml):

    kubectl create -f k8s/deployment.yaml
    kubectl create -f k8s/service.yaml

The defaults should be suitable for most users. You will need to give it a `REDIS_URI` environment variable. You can use it from within the cluster via service DNS as `http://snowflake.default`.

## How It Works
Generates IDs like so (highest-to-lowest bit order):

`[1b:unused][41b: time in ticks since epoch][10b: machine ID][12b: counter]`

Ticks are either seconds (default) or milliseconds (use the `SNOWFLAKE_USE_MILLISECONDS` environment variable and `"true"` or `"false"` to override). Machine ID is a number from 0 to 1023 (inclusive) that identifies the snowflake server. It is picked at random (using math/rand package) and maintained using `SETNX` as a lock in Redis. Default epoch is `2016-01-01T00:00:00Z`, it can be set using the `SNOWFLAKE_EPOCH` environment variable (unit is seconds since Unix epoch). Counter is the number of IDs generated this tick between 0 and 4095 (inclusive). If the counter rotates down to 0 then the server waits until the next clock tick. If the clock runs backwards for any reason, the previous tick timestamp is used.

## TODO
I need to improve the documentation in the code and perhaps provide example requests and a link to the spec in a OpenAPI editor or something.

## License
MIT License

Copyright 2018 Caleb Everett

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
