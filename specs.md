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
