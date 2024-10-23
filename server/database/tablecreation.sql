DROP TABLE IF EXISTS Users;
DROP TABLE IF EXISTS Posts;
DROP TABLE IF EXISTS Comments;
DROP TABLE IF EXISTS Messages;
DROP TABLE IF EXISTS Sessions;
DROP TABLE IF EXISTS Notifications;

CREATE TABLE IF NOT EXISTS Users (
    Id BLOB PRIMARY KEY NOT NULL,
    Nickname STRING NOT NULL,
    Age INTEGER NOT NULL,
    Gender STRING NOT NULL,
    FirstName STRING NOT NULL,
    LastName STRING NOT NULL,
    Email STRING NOT NULL,
    Password STRING NOT NULL  
);

CREATE TABLE IF NOT EXISTS Posts (
    Id BLOB PRIMARY KEY NOT NULL,
    UserId BLOB NOT NULL,
    Title STRING NOT NULL,
    Category STRING,
    Content STRING NOT NULL  
);

CREATE TABLE IF NOT EXISTS Comments (
    Id BLOB PRIMARY KEY NOT NULL,
    PostId BLOB NOT NULL,
    UserId BLOB NOT NULL,
    Content STRING NOT NULL  
);

CREATE TABLE IF NOT EXISTS Messages (
    Id BLOB PRIMARY KEY NOT NULL,
    SenderId BLOB NOT NULL,
    ReceiverId BLOB NOT NULL,
    Message STRING NOT NULL,
    CreatedAt DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS Sessions (
    ID TEXT PRIMARY KEY,
    UserID TEXT,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Notifications (
    Id BLOB PRIMARY KEY NOT NULL,
    CurrentUserId BLOB NOT NULL,
    SenderId BLOB NOT NULL,
    NumberOfUnread INTEGER NOT NULL
);