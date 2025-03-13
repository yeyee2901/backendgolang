CREATE TABLE openweather_master(
    fetch_id uuid NOT NULL,
    coord geometry NOT NULL,
    base varchar,
    main_temp float,
    main_feels_like float,
    main_temp_min float,
    main_temp_max float,
    main_pressure integer,
    main_humidity integer,
    main_sea_level integer,
    main_grnd_level integer,
    visibility integer,
    wind_speed float,
    wind_deg integer,
    wind_gust float,
    clouds_all integer,
    dt integer,
    sys_type integer,
    sys_id integer,
    sys_country varchar,
    sys_sunrise integer,
    sys_sunset integer,

    PRIMARY KEY(fetch_id)
);

CREATE TABLE openweather_details(
    fetch_id uuid NOT NULL,
    idx serial,
    id integer,
    main varchar,
    description varchar,
    icon varchar,

    PRIMARY KEY(fetch_id, idx)
);
