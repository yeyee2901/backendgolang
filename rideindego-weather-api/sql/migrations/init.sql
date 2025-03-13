CREATE EXTENSION postgis;

CREATE TABLE IF NOT EXISTS rideindego_master (
    fetch_id uuid PRIMARY KEY,
    feature_type varchar(25) NOT NULL,
    last_updated timestamp WITH TIME ZONE NOT NULL
);
