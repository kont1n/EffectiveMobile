create table if not exists songs (
    id uuid primary key default gen_random_uuid(),
    song_name varchar(255) not null,
    artist_name varchar(255) not null,
    release_date date,
    song_text text,
    link varchar(255),
    created_at timestamp not null  default (now() at time zone 'utc'),
    unique (song_name, artist_name)
);

create unique index if not exists songs__id_created__idx on songs (id,created_at);
create unique index if not exists songs__artist_song__idx on songs (song_name,artist_name);