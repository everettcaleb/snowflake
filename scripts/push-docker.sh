#!/bin/bash
if [[ $TRAVIS_TAG ]]; then
  TAG_W_MINOR=`echo -n $TRAVIS_TAG | sed 's/\.[0-9]*$//g'`
  TAG_MAJOR=`echo -n $TRAVIS_TAG | sed 's/\.[0-9]*\.[0-9]*$//g'`
  echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
  docker tag everettcaleb/snowflake everettcaleb/snowflake:$TRAVIS_TAG
  docker tag everettcaleb/snowflake everettcaleb/snowflake:$TAG_W_MINOR
  docker tag everettcaleb/snowflake everettcaleb/snowflake:$TAG_MAJOR
  docker push everettcaleb/snowflake:$TRAVIS_TAG
  docker push everettcaleb/snowflake:$TAG_W_MINOR
  docker push everettcaleb/snowflake:$TAG_MAJOR
  docker push everettcaleb/snowflake
fi
