-- Table: movies
CREATE TABLE movies
(
    id          UUID PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    launch_year INT          NOT NULL,
    genre       VARCHAR(100),
    duration    INT          NOT NULL,
    status      VARCHAR(50)  NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE images
(
    id                 UUID PRIMARY KEY,
    file_path          VARCHAR(500) NOT NULL,
    storage_location   VARCHAR(255) NOT NULL,
    content_type       VARCHAR(100) NOT NULL,
    storage_provider   VARCHAR(100) NOT NULL,
    original_file_name VARCHAR(255) NOT NULL,
    movie_id           UUID         NOT NULL REFERENCES movies (id) ON DELETE CASCADE,
    type               VARCHAR(50)  NOT NULL,
    created_at         TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE videos
(
    id                 UUID PRIMARY KEY,
    file_path          VARCHAR(500) NOT NULL,
    storage_location   VARCHAR(255) NOT NULL,
    content_type       VARCHAR(100) NOT NULL,
    storage_provider   VARCHAR(100) NOT NULL,
    original_file_name VARCHAR(255) NOT NULL,
    movie_id           UUID         NOT NULL REFERENCES movies (id) ON DELETE CASCADE,
    type               VARCHAR(50)  NOT NULL,
    created_at         TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMP    NOT NULL DEFAULT NOW(),
    status             VARCHAR(50)  NOT NULL
);
