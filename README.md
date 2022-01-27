# Threat Intelligence / Anonymization Service

## IP Anonymization Application

This application has two (2) core functions.
-   **Αnonymizing** an IP address, given an original IP address and return its obfuscated version.
-   **De-anonymizing** an IP address, given an obfuscated IP address and return its original version.

Application Features,
-   Usage of [CryptoPan](https://github.com/Yawning/cryptopan) algorithm for anonymization.
    -   Anonymizes IP address keeping subnet structure.
-   Implemented a REST-ful anonymization service.
    -   Anonymize IP addresses making POST requests.
    -   De-anonymize IP addresses making POST requests.
    -   Usage of [Go](https://golang.org/) programming language, due to its good concurrency, using goroutines.
-   Delay time measurement
    - Starting measure the moment that Spark sends POST request to anonymization service, until the service respond back with the obfuscated IP.

## API Routes
### Anonymize
Method: POST 
Endpoint http://IP:PORT/anonymize
#### Request Body
```
{
    "IpAddr": String The IP address to be anonymized.
}
```
Example
```
{
    "IpAddr": "206.5.95.138"
}
```

#### Response Body
```
{
    "execution_time": String Time needed (in seconds) to execute the anonymization query.
    "obfuscatedIp": String Obfuscated IP address.
    "originalIp": String Original IP address.
    "result": String Query status.
}
```
Example
```
{
    "execution_time": "0.000917",
    "obfuscatedIp": "206.5.95.138",
    "originalIp": "255.186.223.5",
    "result":"success"
}
```

### De-anonymize
Method: POST 
Endpoint http://IP:PORT/deanonymize
#### Request Body
```
{
    "IpAddr": String. The IP address to be de-anonymized
}
```
Example
```
{
    "IpAddr": "206.5.95.138"
}
```

#### Response Body
```
{
    "execution_time": String Time needed (in seconds) to execute the anonymization query.
    "obfuscatedIp": String Obfuscated IP address to be de-anonymized.
    "originalIp": String The de-anonymized IP address.
    "result": String Query status.
}
```
Example
```
{
    "execution_time": "0.000156",
    "obfuscatedIp": "206.5.95.138",
    "originalIp": "255.186.223.5",
    "result":"success"
}
```

## Data Storing

-   Usage of [Redis](https://redis.io/) database to keep mapping between original & obfuscated IP addresses.
    -   Runs in-memory (Extremely fast)
    -   Data Persistence (RDB + AOF).

-   Usage of [RedisInsight GUI](https://redis.com/redis-enterprise/redis-insight/).
    -   RedisInsight provides an intuitive Redis admin GUI and helps optimize the usage of Redis in our application.


## Deployment
It is a fully distributed project thanks to [Docker](https://www.docker.com/) and [Kubernetes](https://kubernetes.io/). Executing,
```bash
sudo bash deploy-in-k8s.sh
```
-   Three (3) custom docker images will be built (IP anonymization, RedisDB, RedisInsight).
-   Afterwards, these 3 images become K8s deployments with 3 K8s services respectively will be in a Kubernetes namespace.

If everything gone well, the ip-anonymize app logger will show this:

```
INFO: yy/mm/dd hh:mm:ss logger.go:44: /IpAnonymizationService
INFO: yy/mm/dd hh:mm:ss logger.go:44: PING -> PONG
INFO: yy/mm/dd hh:mm:ss logger.go:44: Anonymization Service Running on [an-ip-address]:8100 (Press CTRL+C to quit)
```
