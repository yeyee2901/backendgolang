CREATE EXTENSION postgis;

CREATE TABLE rideindego_master (
    fetch_id uuid,
    feature_type varchar(25) NOT NULL,
    last_updated timestamp WITH TIME ZONE NOT NULL,

    PRIMARY KEY(fetch_id)
);

CREATE TABLE rideindego_features(
    fetch_id uuid,
    feature_id integer NOT NULL,
    feature_type varchar(25) NOT NULL,
    geo_type varchar(25) NOT NULL,
    geo_coord geometry NOT NULL,
    PRIMARY KEY(fetch_id, feature_id)
);

CREATE TABLE rideindego_properties(
    fetch_id uuid,
    feature_id integer NOT NULL,
    properties_id integer NOT NULL,
    coordinates geometry NOT NULL,
    name varchar NOT NULL,
    total_docks integer ,
    docks_available integer ,
    bikes_available integer ,
    classic_bikes_available integer ,
    smart_bikes_available integer ,
    eletric_bikes_available integer ,
    reward_bikes_available integer ,
    reward_docks_available integer ,
    kiosk_status varchar ,
    kiosk_public_status varchar ,
    kiosk_connection_status varchar ,
    kiosk_type integer ,
    address_street varchar,
    address_city varchar,
    address_state varchar,
    address_zip_code varchar,
    close_time varchar DEFAULT NULL,
    event_end varchar,
    event_start varchar,
    is_event_based boolean NOT NULL,
    is_virtual boolean NOT NULL,
    kiosk_id integer NOT NULL,
    notes varchar,
    open_time varchar DEFAULT NULL,
    public_text varchar,
    timezone varchar,
    trikes_available integer,

    PRIMARY KEY(fetch_id, feature_id, properties_id)
);

CREATE INDEX idx_rideindego_kioskid ON rideindego_properties(kiosk_id);
