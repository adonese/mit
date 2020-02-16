

SELECT
    tb.*, tc.FldCityName, tl.FldLocalityName, ts.FldStateName, tn.FldNeighborhoodName
FROM TblBakery tb
    INNER JOIN TblCity tc on tc.FldCityNo = tb.FldCityNo
    INNER JOIN TblLocality tl on tl.FldLocalityNo = tb.FldLocalityNo
    INNER JOIN TblState ts on ts.FldStateNo = tb.FldStateNo
    INNER JOIN TblNeighborhood tn on tn.FldNeighborhoodNo = tb.FldNeighborhoodNo
where tc.FldLocalityNo in (1,2) 

