

-- select tb.*, tb.fldbakeryno, tb.fldbakeryname from tblbakeryaudit ta
-- inner join tblbakery tb on tb.FldBakeryNo = ta.FldBakeyNo

-- select * from tblbakery

-- select * from TblBakeryAudit

-- exec sp_columns tblbakeryaudit

-- describe tblbakery

-- SELECT 
--     is_identity
-- FROM sys.columns
-- WHERE 
--     object_id = object_id('tblbakeryaudit')
--     AND name = 'fldbakeryauditno'

-- select * from tblbakeryshare


-- select tb.FldBakeryName, tb.FldBakeryNo from tblbakeryaudit ta inner join tblbakery tb on tb.FldBakeryNo = ta.FldBakeyNo

-- select * from tblbakeryaudit

-- alter table tblbakeryaudit add FldAlternativeName NVARCHAR

-- select * from tblusers

-- tblflourbaking get quantity

-- SELECT * from TblFlourBakeryReceive

-- select * from TblFlourBaking
-- select * from TblFlourBakeryReceive
-- update TblFlourBakeryReceive set FldBakeryNo = 1

-- select sum(fb.FldQuantity) as quantity, sum(br.FldQuantity) as received_quantity, tb.FldBakeryNo, tb.FldBakeryName
-- from TblFlourBaking fb
-- inner join TblBakery tb on tb.FldBakeryNo = fb.FldBakeryNo
-- inner join TblFlourBakeryReceive br on br.FldBakeryNo = tb.FldBakeryNo
-- where  tb.FldStateNo = 1 AND tb.FldLocalityNo = 1 AND tb.FldAdminNo = 1 AND fb.FldDate BETWEEN '2020-02-10' AND '2020-02-17' AND br.FldDate BETWEEN '2020-02-10' AND '2020-02-17'
-- group by tb.FldBakeryNo, tb.FldBakeryName

-- select fb.FldQuantity, br.FldQuantity
-- from TblBakery tb
--     inner join TblFlourBaking fb on fb.FldBakeryNo = tb.FldBakeryNo
--     inner join TblFlourBakeryReceive br on br.FldBakeryNo = tb.FldBakeryNo

-- where  tb.FldStateNo = 1 AND tb.FldLocalityNo = 1 AND tb.FldAdminNo = 1 AND fb.FldDate BETWEEN '2020-02-10' AND '2020-02-17'
-- group by tb.FldBakeryNo, tb.FldBakeryName

-- select * from tblbakery
-- select * from TblFlourBaking

-- update tblflourbaking set FldBakeryNo = 1


-- select * from TblUsers
-- update  tblusers set fldusertype = 4 where flduserno = 2
-- update tblusers set fldstateno = 1, FldLocaliyNo = 1, fldadminno = 1 where FldUserNo = 2

-- alter table tblusers add FldAdminNo INT

-- select * from TblBakeryAudit


-- select * from tblauditor

select *
from TblBakery

-- update tblusers set FldSystemNo = 3 where FldUserNo = 2