BEGIN;

ALTER TABLE IF EXISTS artists
    DROP CONSTRAINT "artists_name_unique"
;

ALTER TABLE IF EXISTS lyrics
    DROP CONSTRAINT "lyrics_track_id_fkey"
;

ALTER TABLE IF EXISTS tracks
    DROP CONSTRAINT "tracks_title_artist_id_unique"
;

ALTER TABLE IF EXISTS tracks
    DROP CONSTRAINT "tracks_artist_id_fkey"
;


DROP TABLE IF EXISTS
    artists,
    lyrics,
    track
;

END;
