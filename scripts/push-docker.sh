#!/bin/bash
if [[ $TRAVIS_TAG ]]; then
  TAG_VERSION=`echo -n $TRAVIS_TAG | sed 's/v//'`
  TAG_W_MINOR=`echo -n $TAG_VERSION | sed 's/\.[0-9]*$//g'`
  TAG_MAJOR=`echo -n $TAG_VERSION | sed 's/\.[0-9]*\.[0-9]*$//g'`
  echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
  docker tag everettcaleb/snowflake everettcaleb/snowflake:$TAG_VERSION
  docker tag everettcaleb/snowflake everettcaleb/snowflake:$TAG_W_MINOR
  docker tag everettcaleb/snowflake everettcaleb/snowflake:$TAG_MAJOR
  docker push everettcaleb/snowflake:$TAG_VERSION
  docker push everettcaleb/snowflake:$TAG_W_MINOR
  docker push everettcaleb/snowflake:$TAG_MAJOR
  docker push everettcaleb/snowflake
fi
