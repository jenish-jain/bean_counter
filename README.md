# Bean counter 🫘

An Accountant application to help with monthly GST filing for our traditional business

### latest build status
[![Deploy to Cloud Run from Source](https://github.com/jenish-jain/bean_counter/actions/workflows/google-cloudrun-source.yml/badge.svg)](https://github.com/jenish-jain/bean_counter/actions/workflows/google-cloudrun-source.yml)

Libraries Used

Google API go client : https://github.com/googleapis/google-api-go-client

### Curl to generate report

```curl
 curl --location 'https://bean-counter-t2xsqgseuq-uc.a.run.app/gstReport/monthly' \
--header "Authorization: Bearer $(gcloud auth print-identity-token)" \
--header 'Content-Type: application/json' \
--data '{
    "month" : 6,
    "year" : 2023
}
```
