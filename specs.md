# Specs

- currently, password is stored without encryption. Fix it
- auth is db-auth. i usually prefer jwt as it does auth without hitting databse
- username is missing from most of users\*-tables

## TODOS

- Auth part

  - authorization (jwt)
  - logout
  - refresh token

- Agent services
  - get grinders
- distributor services
- bakery services
- auditor services

- Adding chained query methods. Will greatly clean the api when we are building our search engine

## Auth part

These are important fields and should be encoded in `jwt`

- StateNo
- LocalityNo
- CityNo
- NeighborhoodNo
- UserNo
- UserType

### Flour submission

This describe the general workflow for flour submission, and the api used for so.

- First the user log in /login
  We should generate a token for them, to make step two works. TODO
- In the app, the user (in this case, it is a flour agent), would like to submit they have received a quantity of flour. They need to _get_ Grinder info (to make the drop down works)
  `/get_grinders` returns a _list of Grinders for this current logged in agent_, hence why we need a token, to get that.
- After that, you can submit a `Agent Flour Receive` request, which will add this record to `TblFlourAgentReceive`

(i have made another api to get all submitted `Agent Flour Receive` requests for management)

## Bakery services

- Record Received Flour from Flour Agent `TblFlourBakeryReceive` [use tblbakeyshare as lookup]
- Record Baked Flour `TblFlourBaking` [set flddate,fldbakeryno, fldqunatity, fldnote]

### Endpoints (will be updated)

    /login
    /get_grinders
    /get_grinder
    /submit_flour
    /_get_flour

### Agent endpoints

- You need to login first, to get agentID no
  /login

#### Example

`curl -X POST https://mit.soluspay.net/login -d'{"username": "admin", "password": "admin"}' -v`

Agent Receive APIs

    /get_grinders
    /get_grinder
    /submit_flour
    /_get_flour // admin api to get all submitted data, maybe will be useful for dashboard

Agent Distribute API

    /get_bakery
    /submit_bakery

### Agent Receive APIs

#### /get_grinders [GET]

Send `agent` ID in url query, e.g. `GET https://mit.soluspay.net/get_grinders?agent=2

NOTE: Agent is to be get from the login response, with name `FldUserNo`.

NOTE: AGENT with ID `2` has associated grinder

##### RESPONSE

`[]Grinder` object. Grinder is:

| field name     | type    |
| -------------- | ------- |
| FldGrinderNo   | int     |
| FldGrinderName | string  |
| FldIsActive    | bool    |
| FldStateNo     | int     |
| FldContactName | string  |
| FldPhone       | string  |
| FldEmail       | string  |
| FldAddress     | string  |
| FldVolume      | float32 |
| FldUserNo      | int     |
| FldLogNo       | int     |
| FldUpdateDate  | string  |

#### /submit_flour [POST]

An agent submits the flour the received from the grinder

#### Request

| field name                | type    |
| ------------------------- | ------- |
| FldFlourAgentReceiveNo    | int     |
| FldDate                   | string  |
| FldFlourAgentNo           | int     |
| FldGrinderNo              | int     |
| FldQuantity               | float32 |
| FldUnitPrice              | float32 |
| FldTotalAmount            | float32 |
| FldRefNo                  | int     |
| FldNFCFlourAgentReceiveNo | int     |
| FldNFCStatusNo            | int     |
| FldNFCNote                | string  |
| FldUserNo                 | int     |
| FldUpdateDate             | string  |

#### Response

2xx (successful response)

`{"result": "ok"}`

400 (Bad request)

{"message": "A user friendly message you can show", "code": "error_message"}

NOTE: we can use `code` to localize errors messages, they can be used as a hashmap keys and they point to arabic and english messages

### Agent Distribute API

#### /get_bakery [GET]

Get associated Bakeries to agent.

##### Request

TODO ADD UNIT TESTING

URL query: agentID
example: /get_bakery?agent=2

FIXME: queries are not yet supported. However, you can use the apis, and add queries as you want, while I will later implement them.

##### Example

GET /get_bakery?agent=2

```json
[
  {
    "FldBakeryNo": 2,
    "FldBakeryName": "مخبز الاول",
    "FldIsActive": true,
    "FldStateNo": 1,
    "FldLocalityNo": 1,
    "FldCityNo": 2,
    "FldNeighborhoodNo": 1,
    "FldContactName": "N/A",
    "FldPhone": "N/A",
    "FldEmail": "N/A",
    "FldAddress": "N/A",
    "FldVolume": 1550,
    "FldLong": " ",
    "FldLat": " ",
    "FldUserNo": 1,
    "FldLogNo": 10065,
    "FldUpdateDate": "2020-01-30T15:39:00Z",
    "FldImage": "",
    "FldNFCBakeryNo": 0
  }
]
```

#### /bakery/submit [POST]

TODO ADD UNIT TESTING

##### Request

You should send a [`FlourAgentReceive`](https://github.com/adonese/mit/blob/master/db.go#L319-L339) Object

| field name                 | type    |
| -------------------------- | ------- |
| FldFlourBakeryReceiveNo    | int     |
| FldFlourAgentDistributeNo  | int     |
| FldDate                    | string  |
| FldFlourAgentNo            | int     |
| FldBakeryNo                | int     |
| FldQuantity                | float32 |
| FldUnitPrice               | float32 |
| FldTotalAmount             | float32 |
| FldRefNo                   | int     |
| FldNFCFlourBakeryReceiveNo | int     |
| FldNFCStatusNo             | int     |
| FldNFCNote                 | string  |
| FldUserNo                  | int     |
| FldDriverName              | string  |
| FldCarPlateNo              | string  |
| FldUpdateDate              | string  |

##### Response

2xx (successful response)

```json
{ "result": "ok" }
```

400 (Bad request)

```json
{ "message": "A user friendly message you can show", "code": "error_message" }
```

#### /bakery/get_agents [GET]

Get asscoiated agents to _this_ bakery, using bakery ID

##### Request

URL query `agent`

##### Response

Use this example
/bakery/get_agents?agent=2&state=2

```json
[
  {
    "FldFlourAgentNo": 1,
    "FldFlourAgentName": "ابراهيم محمد علي",
    "FldIsActive": true,
    "FldStateNo": 2,
    "FldContactName": "N/A",
    "FldPhone": "N/A",
    "FldEmail": "N/A",
    "FldAddress": "N/A",
    "FldVolume": 3125,
    "FldLong": " ",
    "FldLat": " ",
    "FldUserNo": 1,
    "FldLogNo": 10048,
    "FldUpdateDate": "2020-01-30T14:26:00Z"
  },
  {
    "FldFlourAgentNo": 2,
    "FldFlourAgentName": "توكيل عمر ابراهيم",
    "FldIsActive": true,
    "FldStateNo": 1,
    "FldContactName": "N/A",
    "FldPhone": "N/A",
    "FldEmail": "N/A",
    "FldAddress": "N/A",
    "FldVolume": 1797,
    "FldLong": " ",
    "FldLat": " ",
    "FldUserNo": 1,
    "FldLogNo": 10054,
    "FldUpdateDate": "2020-01-30T15:06:00Z"
  }
]
```

#### /bakery/baked [POST]

This api is used by baker to submit their baked bread.

##### Request

The request object is of type [FlourBaking](https://github.com/adonese/blob/master/db.go/L387:L417)

| field name           | type    |
| -------------------- | ------- |
| FldFlourBakingNo     | int     |
| FldDate              | string  |
| FldBakeryNo          | int     |
| FldWorkingStatusNo   | int     |
| FldQuantity          | float32 |
| FldNote              | string  |
| FldLocalityCheck     | float32 |
| FldLocalityUserNo    | int     |
| FldLocalityNote      | string  |
| FldSecurityCheck     | float32 |
| FldSecurityUserNo    | int     |
| FldSecurityNote      | string  |
| FldGovernmentalCheck | float32 |
| FldGovermentalUserNo | int     |
| FldGovernmentalNote  | int     |
| FldCommunityCheck    | float32 |
| FldComuunityUserNo   | int     |
| FldCommunityNote     | int     |
| FldNFCFlourBakingNo  | int     |
| FldNFCStatusNo       | int     |
| FldNFCNote           | string  |
| FldUserNo            | int     |
| FldUpdateDate        | string  |

2xx (successful response)

`{"result": "ok"}`

400 (Bad request)

{"message": "A user friendly message you can show", "code": "error_message"}

## ِAuditor APIs

- Auditor App
- Flour Auditing
  - Update [TblFlourBaking]
  - Record Baked Flour according to UserType as follows:
    - If UserType =3 [ Set FldLocalityCheck, FldLoclityUserno, FldLocalitynote]
    - If UserType=4 [Set FlSecurityCheck, FldSecurityUserNo, FldSecurityNote]
    - If USerType=5 [Set FldGovernmentalCheck, FldGovernmentalUser, FldGovernmentalNote]
    - If UserType=6 [ Set FldCommunityCheck, FldCommunityUserNo, FldCommunityNote]
    - Check = Flour Quanityt
    - FldxxxxUserNo= Current Logged UserNo
- Violation and Reporting:
  - Use TblBakeryAudit

### Auditor specs

- update only TblFlourBaking
- Get User Type from User Profile
- Request Params:

      FldLocalityCheck
      FldLoclityUserno
      FldLocalitynote
      FlSecurityCheck
      FldSecurityUserNo
      FldSecurityNote
      FldGovernmentalCheck
      FldGovernmentalUser
      FldGovernmentalNote
      FldCommunityCheck
      FldCommunityUserNo
      FldCommunityNote

      (Should send Flour Quantity)

#### /auditor/report [POST]

##### Request

| field name          | type   |
| ------------------- | ------ |
| FldBakeryAuditNo    | int    |
| FldDate             | string |
| FldBakeyNo          | int    |
| FldAuditBy          | int    |
| FldAuditType        | int    |
| FldAuditStatusNo    | int    |
| FldNote             | string |
| FldAuditResponseNo  | int    |
| FldNFCBakeryAuditNo | int    |
| FldNFCStatusNo      | int    |
| FldNFCNote          | string |
| FldUserNo           | int    |
| FldUpdateDate       | string |

##### Response

2xx (successful response)

```json
{ "result": "ok" }
```

400 (Bad request)

```json
{ "message": "A user friendly message you can show", "code": "error_message" }
```

#### /auditor/check [POST]

This api _should_ send the baked amount of bread in bakeries, by perspective auditors.

##### Request

| field name           | type    |
| -------------------- | ------- |
| FldFlourBakingNo     | int     |
| FldDate              | string  |
| FldBakeryNo          | int     |
| FldWorkingStatusNo   | int     |
| FldQuantity          | float32 |
| FldNote              | string  |
| FldLocalityCheck     | float32 |
| FldLocalityUserNo    | int     |
| FldLocalityNote      | string  |
| FldSecurityCheck     | float32 |
| FldSecurityUserNo    | int     |
| FldSecurityNote      | string  |
| FldGovernmentalCheck | float32 |
| FldGovermentalUserNo | int     |
| FldGovernmentalNote  | int     |
| FldCommunityCheck    | float32 |
| FldComuunityUserNo   | int     |
| FldCommunityNote     | int     |
| FldNFCFlourBakingNo  | int     |
| FldNFCStatusNo       | int     |
| FldNFCNote           | string  |
| FldUserNo            | int     |
| FldUpdateDate        | string  |

##### Response

2xx (successful response)

```json
{ "result": "ok" }
```

400 (Bad request)

```json
{ "message": "A user friendly message you can show", "code": "error_message" }
```

#### /auditor/complains [GET]

#### Request

Empty body

#### Response

```json
{ "FldAuditStatusNo": 0, "FldAuditStatusName": "No Bakery" }
```

## Admin Services

These are services for use by mit system adminstrators only.

- create users
- change / assign roles to users
