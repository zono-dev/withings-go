# withings-go

withings-go is UNOFFICIAL go client to acess Withings API easily.
More Withings API document can be found in [Withings developer documentation](https://developer.withings.com/oauth2/#).

## Supported Resources

- Offline Authorization
- [Measure - Getmeas](https://developer.withings.com/oauth2/#operation/measure-getmeas)
- [Measure v2 - Getactivity](https://developer.withings.com/oauth2/#operation/measurev2-getactivity)
- [Sleep v2 - Get](https://developer.withings.com/oauth2/#operation/sleepv2-get)
- [Sleep v2 - Getsummary](https://developer.withings.com/oauth2/#operation/sleepv2-getsummary)


## Installation

`go get github.com/zono-dev/withings-go`

## Getting started

### Create your Withings account

If you have any Withings gear, you should already have your Withings account.

### Register your app

[Register your app](https://account.withings.com/partner/add_oauth2) and get your client ID and consumer secret.

### Setup your settings file

Create `.test_settings.yaml` in `cmd/auth`.

Tip: In the `cmd` directory, there is examples of how to use withings-go.

`.test_settings.yaml` sample is here.

```
CID: "YOUR CLIENT ID HERE"
Secret: "YOUR CONSUMER SECRET HERE"
RedirectURL: "https://example.com/"
```

`cd cmd/auth`, and `go build` then `auth` will be generated.

### Run

When you run `auth`, some messages will be displayed in the console as follows.

```
map[CID:YOUR-CLIENT-ID RedirectURL:https://example.com/ Secret:YOUR-SECRET]
[user.activity,user.metrics,user.info]
URL to authorize:http://account.withings.com/oauth2_user/authorize2?access_type=offline&client_id=yourclientid&redirect_uri=https%3A%2F%2Fexample.com&response_type=code&scope=user.activity%2Cuser.metrics%2Cuser.info&state=state
Open url your browser and Enter your grant code here.
 Grant Code:
```

### Get token

Go to the `URL to authorize` URL in your browser and allow your app connect to your withings account.

Then your browser will be redirect to your `RedirectURL` automatically and you will find your `Grant code` in redirected URL.

Type this into the console followed by the `Grant Code:`.

Lastly, your `access_token.json` will be generated in `cmd/auth`.

## Get Measurements 

Create `.test_settings.yaml` in `cmd/getMeasurements`.

`.test_settings.yaml` sample is here.

```
CID: "YOUR CLIENT ID HERE"
Secret: "YOUR CONSUMER SECRET HERE"
RedirectURL: "https://example.com/"
```

`cd cmd/getMeasurements` and `go build` then `getMeasurements` will be generated.

When you run `getMeasurements`, your measurements will be displayed in the console.

###
