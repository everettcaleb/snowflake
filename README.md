# Snowflake API
ID generation service inspired by Twitter's [Snowflake](https://github.com/twitter-archive/snowflake/tree/b3f6a3c6ca8e1b6847baa6ff42bf72201e2c2231)

[![Build Status](https://travis-ci.com/everettcaleb/snowflake.svg?branch=master)](https://travis-ci.com/everettcaleb/snowflake)
[![Maintainability](https://api.codeclimate.com/v1/badges/a3af91de1b11806ea09e/maintainability)](https://codeclimate.com/github/everettcaleb/snowflake/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/a3af91de1b11806ea09e/test_coverage)](https://codeclimate.com/github/everettcaleb/snowflake/test_coverage)
[![Docker Pulls](https://img.shields.io/docker/pulls/everettcaleb/snowflake.svg?style=flat)](https://hub.docker.com/r/everettcaleb/snowflake)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](LICENSE)

## Using with Docker
You can run the following command to run it locally from [Docker Hub](https://hub.docker.com/r/everettcaleb/snowflake/):

    docker run -d --name snowflake -e REDIS_URI="redis://host:port/" -p 8080:8080 everettcaleb/snowflake

Then you can test it with:

    curl http://localhost:8080/id

## How It Works
Generates IDs like so (highest-to-lowest bit order):

`[1b:unused][41b: time in ticks since epoch][10b: machine ID][12b: counter]`

Ticks are either seconds (default) or milliseconds (use the `SNOWFLAKE_USE_MILLISECONDS` environment variable and `"true"` or `"false"` to override). Machine ID is a number from 0 to 1023 (inclusive) that identifies the snowflake server. It is picked at random (using math/rand package) and maintained using `SETNX` as a lock in Redis. Default epoch is `2016-01-01T00:00:00Z`, it can be set using the `SNOWFLAKE_EPOCH` environment variable (unit is seconds since Unix epoch). Counter is the number of IDs generated this tick between 0 and 4095 (inclusive). If the counter rotates down to 0 then the server waits until the next clock tick. If the clock runs backwards for any reason, the previous tick timestamp is used.

## TODO
I need to improve the documentation in the code and perhaps provide example requests.

## License
MIT License

Copyright 2019 Caleb Everett

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
