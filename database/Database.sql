select * from recurring_bills

SELECT SUM(amount) FROM recurring_bills;

select * from menus;
delete from menus;
select * from menus;

INSERT INTO menus (Id, Menu, Url, Is_Active) VALUES (0, 'recurring bills', '/recurring-debts', 1);
