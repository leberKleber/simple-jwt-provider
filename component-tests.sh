#!/bin/bash -e

BUILD_ID=$RANDOM
export BUILD_ID

networkName="component-tests-$BUILD_ID"

docker build -f component-tests.Dockerfile -t ct-runner:$BUILD_ID .
docker network create $networkName

docker-compose -f component-tests.docker-compose.yml up -d --build

set +e
docker run --rm -e "WUC_EXPECTED=200" \
      -e "WUC_WRITE_OUT=%{http_code}" \
      -e "WUC_URL=http://simple-jwt-provider/v1/internal/alive" \
      -e "WUC_MAX_ITERATIONS=20" \
      --network "${networkName}" leberkleber/wait_until_curl
sjp_alive=$?
set -e

if [[ "$sjp_alive" -ne 0 ]]; then
  echo "provider didn't start succesfull"
  test_result=1
else
  set +e
      docker run --rm --network "${networkName}" ct-runner:${BUILD_ID}
      test_result=$?
  set -e
fi

if [[ "$test_result" -gt 0 ]]; then
  docker-compose -f component-tests.docker-compose.yml logs
fi

docker-compose -f component-tests.docker-compose.yml down
docker network rm $networkName

exit $test_result
