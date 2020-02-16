SELECT
    tu.*, tb.FldPhone
FROM TblUsers tu
    INNER JOIN TblBakery tb on tb.FldSystemNo = tu.FldSystemNo

