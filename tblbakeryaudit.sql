-- 


-- set ansi_warnings OFF


insert into tblbakeryaudit
    (FldDate,FldBakeyNo,FldAuditBy,FldAuditType,FldAuditStatusNo,FldNote,FldUserNo, FldAlternativeName, FldPhone)
values
    ('2020-02-25', 3, 4, 2, 1, 'my note', 3, 'my text', '0912141679')

-- select * from TblBakeryAudit

-- alter TABLE tblbakeryaudit ADD  FldPhone text, FldAlternativeName text
-- alter TABLE tblbakeryaudit drop COLUMN  FldPhone, FldAlternativeName




-- select * from tblauditor

-- select * from tbl

-- alter table tblauditor add FldUserNo INT

-- update tblauditor set FldUserNo = 3, FldPhone = '0912141679'



-- Select tu.*, ta.FldPhone from tblusers tu
-- INNER JOIN TblAuditor ta on ta.FldUserNo = tu.FldSystemNo
-- where tu.FldsystemNo = 3


-- select * from tblauditor

select *
from TblBakeryAudit




