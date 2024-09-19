create table if not exists students(
    id SERIAL PRIMARY KEY,
    name text not null,
    age int not null,
    version int not null default 1
);

insert into students (name, age)
values ('张三', 18);