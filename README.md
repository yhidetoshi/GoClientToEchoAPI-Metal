# GoClientToEchoAPI-MetalToMackerel

■ デプロイ
```
export MKRKEY=XXX
export APIURL=XXX

curl -X POST https://api.mackerelio.com/api/v0/services \
    -H "X-Api-Key: ${MKRKEY}" \
    -H "Content-Type: application/json" \
    -d '{"name": "Metal", "memo": "metal"}'

make build
sls deploy --aws-profile <PROFILE> --mkrkey ${MKRKEY} --apiurl ${APIURL}
```