BEGIN;

CREATE TABLE IF NOT EXISTS artists
(
    "artist_id" SERIAL NOT NULL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tracks
(
    "track_id" SERIAL PRIMARY KEY,
    "artist_id" INTEGER NOT NULL,
    "title" VARCHAR(255) NOT NULL,
    "link" VARCHAR(2048) NOT NULL,
    "released_at" TIMESTAMP NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS lyrics
(
    "lyric_id" SERIAL NOT NULL PRIMARY KEY,
    "track_id" INTEGER NOT NULL,
    "verse_text" VARCHAR(2048) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE IF EXISTS lyrics
    ADD FOREIGN KEY ("track_id") REFERENCES tracks ("track_id")
;

ALTER TABLE IF EXISTS tracks
    ADD FOREIGN KEY ("artist_id") REFERENCES artists ("artist_id")
;

END;
