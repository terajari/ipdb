ALTER TABLE podcasts
    DROP CONSTRAINT check_podcasts_year;

ALTER TABLE podcasts
    DROP CONSTRAINT check_tags_length;

ALTER TABLE podcasts
    DROP CONSTRAINT check_guest_speakers_length;