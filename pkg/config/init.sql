CREATE DATABASE IF NOT EXISTS tiktok;
USE tiktok;

CREATE TABLE users
(
    id          VARCHAR(64) PRIMARY KEY NOT NULL,
    username    VARCHAR(50)  NOT NULL,
    password    VARCHAR(255) NOT NULL,
    avatar_url  VARCHAR(255) DEFAULT '' NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NULL,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL,
    mfa_secret  VARCHAR(100) DEFAULT '' NULL,
    mfa_enabled TINYINT(1) DEFAULT 0 NULL,
    UNIQUE INDEX uk_users_username (username)
);
CREATE TABLE videos
(
    id            VARCHAR(64) PRIMARY KEY NOT NULL,
    user_id       VARCHAR(64) NULL,
    video_url     VARCHAR(255) NULL,
    cover_url     VARCHAR(255) DEFAULT '' NOT NULL,
    title         VARCHAR(255) NULL,
    description   TEXT NULL,
    visit_count   BIGINT DEFAULT 0 NOT NULL,
    like_count    BIGINT DEFAULT 0 NOT NULL,
    comment_count BIGINT DEFAULT 0 NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP NULL,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP NULL,
    INDEX idx_videos_userid (user_id)
);
CREATE TABLE relations
(
    user_id      VARCHAR(64) NOT NULL,
    follower_id  VARCHAR(64) NULL,
    following_id VARCHAR(64) NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP NULL,
    UNIQUE INDEX uk_user_follower (user_id, follower_id),
    UNIQUE INDEX uk_user_following (user_id, following_id)
);
CREATE TABLE likes
(
    user_id       VARCHAR(64) NOT NULL,
    to_video_id   VARCHAR(64) NULL,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP NULL,
    to_comment_id VARCHAR(64) NULL,
    UNIQUE INDEX uk_comment_user (user_id, to_comment_id),
    UNIQUE INDEX uk_video_user (to_video_id, user_id)
);
CREATE TABLE comments
(
    comment_id VARCHAR(64) PRIMARY KEY NOT NULL,
    user_id    VARCHAR(64) NOT NULL,
    video_id   VARCHAR(64) NOT NULL,
    content    TEXT NOT NULL,
    like_count INT DEFAULT 0 NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_video_id (video_id)
);
CREATE TABLE message
(
    session_id VARCHAR(128),
    content TEXT,
    sender_id VARCHAR(64),
    receiver_id VARCHAR(64),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL
);
CREATE TABLE friends
(
    user_id varchar(64),
    friend_id varchar(64),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL
)
