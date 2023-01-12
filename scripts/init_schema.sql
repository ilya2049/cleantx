create table if not exists doctors (
  id int generated always as identity,
  name text,
  is_on_shift bool
);

insert into doctors (name, is_on_shift)
values 
('Ada', true),
('Elvin', true),
('Todd', false);