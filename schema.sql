PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT OR IGNORE INTO users (id, email, username, password)
VALUES (1, 'test@example.com', 'testuser', 'placeholder_hash');

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

INSERT OR IGNORE INTO categories (name)
VALUES ('test1'), ('test2'), ('test3');

CREATE TABLE IF NOT EXISTS post_categories (
    post_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER,
    comment_id INTEGER,
    reaction_type TEXT NOT NULL CHECK (reaction_type IN ('like', 'dislike')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, post_id, comment_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE
);

INSERT OR IGNORE INTO posts (id, user_id, title, content, created_at)
VALUES
    (1, 1, 'test post 1', 'test1', CURRENT_TIMESTAMP),
    (2, 1, 'test post 2', 'test2', CURRENT_TIMESTAMP),
    (3, 1, 'test post 3', 'test3', CURRENT_TIMESTAMP);

INSERT OR IGNORE INTO post_categories (post_id, category_id)
SELECT 1, id FROM categories WHERE name = 'test1'
UNION ALL
SELECT 2, id FROM categories WHERE name = 'test2'
UNION ALL
SELECT 3, id FROM categories WHERE name = 'test3';

INSERT OR IGNORE INTO comments (id, post_id, user_id, content, created_at)
VALUES
    (1, 1, 1, 'test1', CURRENT_TIMESTAMP),
    (2, 2, 1, 'test2', CURRENT_TIMESTAMP),
    (3, 3, 1, 'test3', CURRENT_TIMESTAMP);

INSERT OR IGNORE INTO reactions (id, user_id, post_id, reaction_type, created_at)
VALUES
    (1, 1, 1, 'like', CURRENT_TIMESTAMP),
    (2, 1, 2, 'dislike', CURRENT_TIMESTAMP),
    (3, 1, 3, 'like', CURRENT_TIMESTAMP);

    CREATE TABLE IF NOT EXISTS sessions (
        token TEXT PRIMARY KEY,
        user_id INTEGER NOT NULL,
        expiry DATETIME NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    )