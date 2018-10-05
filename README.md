# Snowflake API
ID generation service inspired by Twitter's [Snowflake](https://github.com/twitter-archive/snowflake/tree/b3f6a3c6ca8e1b6847baa6ff42bf72201e2c2231)



## Using with Docker
You can run the following command to run it locally from [Docker Hub](https://hub.docker.com/r/everettcaleb/snowflake/):

    docker run -d --name snowflake -e MACHINE_ID=0 -p 8080:8080 everettcaleb/snowflake

Then you can test it with:

    curl http://localhost:8080/v1/snowflake

## Using with Kubernetes
You can deploy it as a [StatefulSet](k8s/statefulset.yaml) and [Service](k8s/service.yaml) (a Deployment won't work because the `MACHINE_ID` values need to be unique):

    kubectl create -f k8s/statefulset.yaml
    kubectl create -f k8s/service.yaml

The defaults should be suitable for most users. You can use it from within the cluster via service DNS as `http://snowflake.default`.

## How It Works
Generates IDs like so (highest-to-lowest bit order):

`[1b:unused][41b: time in ms since epoch][10b: machine ID][12b: counter]`

Machine ID is a number from 0 to 1023 (inclusive) that identifies the snowflake server. It is retrieved from an environment variable or the end of the hostname (ex: `snowflake-0` or `snowflake-2`). Epoch is 2018-01-01T00:00:00Z. Counter is the number of IDs generated this millisecond between 0 and 4095 (inclusive). If the counter rotates down to 0 then the server waits until the next clock millisecond. If the clock runs backwards, the previous millisecond timestamp is used.

Note: Machine ID is automatically populated if you're using a StatefulSet in Kubernetes and the `MACHINE_ID` environment variable is set to `HOST`.

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
