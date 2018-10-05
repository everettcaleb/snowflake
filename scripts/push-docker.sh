#!/bin/bash
if [[ $TRAVIS_TAG ]]; then
  echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
  docker tag everettcaleb/snowflake everettcaleb/snowflake:$TRAVIS_TAG
  docker push everettcaleb/snowflake:$TRAVIS_TAG
  docker push everettcaleb/snowflake
fi
