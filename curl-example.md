# Curl request example

## User signup

```shell
curl --location 'http://localhost:8080/users/v1/signup' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Lucas",
    "birthday": "2020-10-23T00:00:00Z",
    "gender": "male",
    "photo_url": "photo_url"
}'
```
