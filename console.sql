CREATE TABLE snippets(
    id SERIAL NOT NULL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMP NOT NULL,
    expires TIMESTAMP NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'Не имей сто рублей',
    'Не имей сто рублей,\nа имей сто друзей.',
    now(),
    now() + INTERVAL'365 day'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'Не откладывай на завтра',
    'Не откладывай на завтра,\nчто можешь сделать сегодня.',
    now(),
    now() + INTERVAL'7 day'
);

CREATE USER web WITH PASSWORD 'q';
GRANT SELECT, INSERT, UPDATE ON snippets TO web;

SELECT id, title, expires FROM  snippets



