create table if not exists doctors (
  id int generated always as identity,
  name text,
  on_call bool
);

insert into doctors (name, on_call)
values 
('Ada', true),
('Elvin', true),
('Todd', false);