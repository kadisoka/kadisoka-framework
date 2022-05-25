--

-- convention:
-- prefix _m indicates metadata of a record.
-- prefix _mc_ indicates metadata about a record's creation.
-- prefix _md_ indicates metadata about a record's (soft) deletion.

\set ON_ERROR_STOP true

BEGIN;
------

CREATE TABLE user_dt (
    id_num  bigint PRIMARY KEY,

    _mc_ts   timestamp with time zone NOT NULL DEFAULT now(),
    _mc_tid  bigint,
    _mc_uid  bigint,

    _md_ts     timestamp with time zone,
    _md_tid    bigint,
    _md_uid    bigint,
    _md_notes  jsonb,

    CHECK (id_num > 0)
);

CREATE TABLE terminal_dt (
    id_num          bigint PRIMARY KEY,
    application_id  integer NOT NULL,
    user_id         bigint NOT NULL, -- use zero if it's for a non-user

    _mc_ts              timestamp with time zone NOT NULL DEFAULT now(),
    _mc_tid             bigint,
    _mc_uid             bigint,
    _mc_origin_address  text,
    _mc_origin_env      text,

    _md_ts   timestamp with time zone,
    _md_tid  bigint,
    _md_uid  bigint,

    secret           text NOT NULL,
    display_name     text,
    accept_language  text NOT NULL DEFAULT '', --TODO: list?

    verification_type  text NOT NULL,
    verification_id    bigint NOT NULL,
    verification_ts    timestamp with time zone,

    CHECK (application_id > 0 AND id_num > 0)
);
CREATE INDEX ON terminal_dt (user_id)
    WHERE _md_ts IS NULL
    AND verification_ts IS NOT NULL;

CREATE TABLE session_dt (
    terminal_id  bigint NOT NULL,
    id_num       bigint NOT NULL,

    _mc_ts   timestamp with time zone NOT NULL DEFAULT now(),
    _mc_tid  bigint,
    _mc_uid  bigint,

    _md_ts   timestamp with time zone,
    _md_tid  bigint,
    _md_uid  bigint,

    expiry  timestamp with time zone,

    PRIMARY KEY (terminal_id, id_num),
    CHECK (terminal_id > 0 AND id_num > 0)
);

CREATE TABLE terminal_fcm_registration_token_dt (
    terminal_id  bigint NOT NULL,
    user_id      bigint NOT NULL,
    token        text NOT NULL,

    _mc_ts   timestamp with time zone NOT NULL DEFAULT now(),
    _mc_tid  bigint NOT NULL,
    _mc_uid  bigint NOT NULL,

    _md_ts   timestamp with time zone,
    _md_tid  bigint,
    _md_uid  bigint
);
CREATE UNIQUE INDEX terminal_fcm_registration_token_dt_pidx
    ON terminal_fcm_registration_token_dt (terminal_id)
    WHERE _md_ts IS NULL;
CREATE INDEX terminal_fcm_registration_token_dt_user_idx
    ON terminal_fcm_registration_token_dt (user_id)
    WHERE user_id is NOT NULL AND _md_ts IS NULL;

CREATE TABLE user_key_phone_number_dt (
    user_id          bigint NOT NULL,
    country_code     integer NOT NULL,
    national_number  bigint NOT NULL, -- libphonenumber says uint64
    raw_input        text NOT NULL, 

    _mc_ts   timestamp with time zone NOT NULL DEFAULT now(),
    _mc_tid  bigint NOT NULL,
    _mc_uid  bigint NOT NULL,

    _md_ts   timestamp with time zone,
    _md_tid  bigint,
    _md_uid  bigint,

    verification_id  bigint NOT NULL DEFAULT 0,
    verification_ts  timestamp with time zone
);
-- Each user can only have one reference to the same phone number
CREATE UNIQUE INDEX user_key_phone_number_dt_pidx
    ON user_key_phone_number_dt (user_id, country_code, national_number)
    WHERE _md_ts IS NULL;
-- Each user has only one non-deleted-verified phone number
CREATE UNIQUE INDEX user_key_phone_number_dt_user_id_uidx
    ON user_key_phone_number_dt (user_id)
    WHERE _md_ts IS NULL AND verification_ts IS NOT NULL;
-- One instance for a non-deleted-verified phone number
CREATE UNIQUE INDEX user_key_phone_number_dt_country_code_national_number_uidx
    ON user_key_phone_number_dt (country_code, national_number)
    WHERE _md_ts IS NULL AND verification_ts IS NOT NULL;

CREATE TABLE phone_number_verification_dt (
    id_num              bigserial PRIMARY KEY,
    country_code        integer NOT NULL,
    national_number     bigint NOT NULL,
    code                text NOT NULL,
    code_expiry         timestamp with time zone,
    attempts_remaining  smallint NOT NULL DEFAULT 3,

    _mc_ts   timestamp with time zone NOT NULL DEFAULT now(),
    _mc_tid  bigint,
    _mc_uid  bigint,

    confirmation_ts   timestamp with time zone,
    confirmation_tid  bigint,
    confirmation_uid  bigint
);

CREATE TABLE user_contact_phone_number_dt (
    user_id                  bigint NOT NULL,
    contact_country_code     integer NOT NULL,
    contact_national_number  bigint NOT NULL,

    _mc_ts   timestamp with time zone NOT NULL DEFAULT now(),
    _mc_tid  bigint NOT NULL,
    _mc_uid  bigint NOT NULL,

    PRIMARY KEY (user_id, contact_country_code, contact_national_number)
);

-- user profile
CREATE TABLE user_display_name_dt (
    user_id       bigint NOT NULL,
    display_name  text NOT NULL,

    _mc_ts   timestamp with time zone NOT NULL DEFAULT now(),
    _mc_tid  bigint NOT NULL,
    _mc_uid  bigint NOT NULL,

    _md_ts   timestamp with time zone,
    _md_tid  bigint,
    _md_uid  bigint
);
CREATE UNIQUE INDEX user_display_name_dt_pidx
    ON user_display_name_dt (user_id)
    WHERE _md_ts IS NULL;

CREATE TABLE user_profile_image_key_dt (
    user_id            bigint NOT NULL,
    profile_image_key  text NOT NULL,

    _mc_ts   timestamp with time zone NOT NULL DEFAULT now(),
    _mc_tid  bigint NOT NULL,
    _mc_uid  bigint NOT NULL,

    _md_ts   timestamp with time zone,
    _md_tid  bigint,
    _md_uid  bigint
);
CREATE UNIQUE INDEX user_profile_image_key_dt_pidx
    ON user_profile_image_key_dt (user_id)
    WHERE _md_ts IS NULL;

-- user email
CREATE TABLE user_key_email_address_dt (
    user_id      bigint NOT NULL,
    domain_part  text NOT NULL,
    local_part   text NOT NULL,
    raw_input    text NOT NULL,

    _mc_ts   timestamp with time zone NOT NULL DEFAULT now(),
    _mc_tid  bigint NOT NULL,
    _mc_uid  bigint NOT NULL,

    _md_ts   timestamp with time zone,
    _md_tid  bigint,
    _md_uid  bigint,

    verification_id  bigint NOT NULL DEFAULT 0,
    verification_ts  timestamp with time zone
);
-- Each user can only have one reference to the same email address
CREATE UNIQUE INDEX user_key_email_address_dt_pidx
    ON user_key_email_address_dt (user_id, domain_part, local_part)
    WHERE _md_ts IS NULL;
-- Each user has only one non-deleted-verified email address
CREATE UNIQUE INDEX user_key_email_address_dt_user_id_uidx
    ON user_key_email_address_dt (user_id)
    WHERE _md_ts IS NULL AND verification_ts IS NOT NULL;
-- One instance for a non-deleted-verified email address
CREATE UNIQUE INDEX user_key_email_address_dt_domain_part_local_part_uidx
    ON user_key_email_address_dt (domain_part, local_part)
    WHERE _md_ts IS NULL AND verification_ts IS NOT NULL;

CREATE TABLE email_address_verification_dt (
    id_num              bigserial PRIMARY KEY,
    domain_part         text NOT NULL,
    local_part          text NOT NULL,
    code                text NOT NULL,
    code_expiry         timestamp with time zone,
    attempts_remaining  smallint NOT NULL DEFAULT 3,

    _mc_ts   timestamp with time zone NOT NULL DEFAULT now(),
    _mc_tid  bigint,
    _mc_uid  bigint,

    confirmation_ts   timestamp with time zone,
    confirmation_tid  bigint,
    confirmation_uid  bigint
);

-- user password
--TODO: passwords for different purposes?
CREATE TABLE user_password_dt (
    user_id   bigint NOT NULL,
    password  text NOT NULL,

    _mc_ts   timestamp with time zone NOT NULL DEFAULT now(),
    _mc_tid  bigint NOT NULL,
    _mc_uid  bigint NOT NULL,

    _md_ts   timestamp with time zone,
    _md_tid  bigint,
    _md_uid  bigint
);
CREATE UNIQUE INDEX ON user_password_dt (user_id)
    WHERE _md_ts IS NULL;

----
END;
