ALTER TABLE podcasts
    ADD CONSTRAINT check_podcasts_year CHECK (year BETWEEN 2003 AND date_part('year', NOW()));

ALTER TABLE podcasts
    ADD CONSTRAINT check_tags_length CHECK (array_length(tags, 1) BETWEEN 1 AND 10);

ALTER TABLE podcasts
    ADD CONSTRAINT check_guest_speakers_length CHECK (array_length(guest_speakers, 1) BETWEEN 1 AND 10);