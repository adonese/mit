# Agent app

-

# Bakery app

# auditor app

# Known bugs

- login logic
- bakery structure (get bakeries)
- get bakeries (for agent)
- user type
- fix agent get bakeries to include the geo naming as well

USER TYPE

- auth (user id, or system no)
- fldbakeryno, fldagentno
- userid = {fldbakeryno, fldagentno} => systemno

        /get_bakeries?agent=userid
        /get_agents?agent=userid

/submit_flour

- fldflouragentreceiveno autonumber
- fldupdate is today's date anyway
- auditor:
  - city
  - neighborhood
  - locality
  - addendum / adminno
    ARE ALL inferenced from respective user's table
- subsequent queries for localities
- fldworkingstatusno:

  - bet fldworkingstatusno
    - 0: working
    - 1: status_1
    - 2: status_2

- any fld\*check is for the remaining stock
- fldflourbakingno: autoincrememt
- auditor:

Check for existence of record _today_ with fldbakeryno

- all for actor to add a note
- record does not exist:
  - set quantity to ZERO
- if exists: ignore (for fldquantity)

Auditor report

- fldreportby (name)
- auditstatusno (violationno)
- unit price optional (zero value)
- fldreferneceno (string)

- add user name
- add mobile number
- check public variables
