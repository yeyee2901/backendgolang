CREATE TABLE rideindego_bikes(
    fetch_id uuid,
    feature_id integer,
    properties_id integer,
    id SERIAL,
    dock_number integer NOT NULL,
    is_electric boolean NOT NULL,
    is_available boolean NOT NULL,
    battery integer,

    PRIMARY KEY(fetch_id, feature_id, properties_id, id)
);
