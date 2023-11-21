# Load Test


## Locust

Locust is an open source tool written in python for load testing apis


## Test

This test handles mixed load, with
```
80% read
20% write
```

Each user will generate a request between `0.1` and `0.5` seconds. 


## Run

start master:
```
locust -f locustfile.py --master
```

start workers:
```
locust -f locustfile.py --worker
```


## Dashboard

To access the dashboard, open your web browser to `0.0.0.0:8089`. Here, the test can be configured for both concurrency (num of users) and time to ramp up users to max concurrency.