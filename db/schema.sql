create table tickets
(
    id         int auto_increment primary key,
    title      varchar(255)                       not null,
    status     tinyint(1)                         not null,
    created_at datetime default CURRENT_TIMESTAMP not null,
    constraint tickets_id_uindex unique (id)
);

INSERT INTO todo.tickets (title, status) VALUES ('Ticket 1', 1);
INSERT INTO todo.tickets (title, status) VALUES ('Ticket 2', 1);

