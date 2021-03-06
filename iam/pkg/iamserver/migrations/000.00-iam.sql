--
\set ON_ERROR_STOP true

BEGIN;
------

CREATE TABLE user_dt (
    id       bigint PRIMARY KEY,

    c_ts     timestamp with time zone NOT NULL DEFAULT now(),
    c_uid    bigint,
    c_tid    bigint,

    d_ts     timestamp with time zone,
    d_uid    bigint,
    d_tid    bigint,
    d_notes  jsonb,

    CHECK (id > 0)
);

CREATE TABLE user_key_phone_number_dt (
    user_id          bigint NOT NULL,
    country_code     integer NOT NULL,
    national_number  bigint NOT NULL, -- libphonenumber says uint64
    raw_input        text NOT NULL, 

    c_ts    timestamp with time zone NOT NULL DEFAULT now(),
    c_uid   bigint NOT NULL,
    c_tid   bigint NOT NULL,

    d_ts    timestamp with time zone,
    d_uid   bigint,
    d_tid   bigint,

    verification_id    bigint NOT NULL DEFAULT 0,
    verification_time  timestamp with time zone
);
-- Each user can only have one reference to the same phone number
CREATE UNIQUE INDEX user_key_phone_number_dt_pidx
    ON user_key_phone_number_dt (user_id, country_code, national_number)
    WHERE d_ts IS NULL;
-- Each user has only one non-deleted-verified phone number
CREATE UNIQUE INDEX user_key_phone_number_dt_user_id_uidx
    ON user_key_phone_number_dt (user_id)
    WHERE d_ts IS NULL AND verification_time IS NOT NULL;
-- One instance for an non-deleted-verified phone number
CREATE UNIQUE INDEX user_key_phone_number_dt_country_code_national_number_uidx
    ON user_key_phone_number_dt (country_code, national_number)
    WHERE d_ts IS NULL AND verification_time IS NOT NULL;

CREATE TABLE terminal_dt (
    application_id     integer NOT NULL,
    user_id            bigint, --TODO: not null. use zero if it's for non-user
    id                 bigint PRIMARY KEY,
    secret             text NOT NULL,
    display_name       text,
    accept_language    text NOT NULL DEFAULT '', --TODO: list?

    c_ts               timestamp with time zone NOT NULL DEFAULT now(),
    c_uid              bigint,
    c_tid              bigint,
    c_origin_address   text,
    c_origin_env       text,

    verification_type  text NOT NULL,
    verification_id    bigint NOT NULL,
    verification_time  timestamp with time zone,

    CHECK (application_id > 0 AND id > 0)
);
CREATE INDEX ON terminal_dt (user_id)
    WHERE verification_time IS NOT NULL;

CREATE TABLE session_dt (
    terminal_id  bigint NOT NULL,
    id           bigint NOT NULL,

    c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    c_uid  bigint,
    c_tid  bigint,

    d_ts   timestamp with time zone,
    d_uid  bigint,
    d_tid  bigint,

    PRIMARY KEY (terminal_id, id),
    CHECK (terminal_id > 0 AND id > 0)
);

CREATE TABLE terminal_fcm_registration_token_dt (
    terminal_id  bigint NOT NULL,
    user_id      bigint NOT NULL,
    token        text NOT NULL,

    c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    c_uid  bigint NOT NULL,
    c_tid  bigint NOT NULL,

    d_ts   timestamp with time zone,
    d_uid  bigint,
    d_tid  bigint
);
CREATE UNIQUE INDEX terminal_fcm_registration_token_dt_pidx
    ON terminal_fcm_registration_token_dt (terminal_id)
    WHERE d_ts IS NULL;
CREATE INDEX terminal_fcm_registration_token_dt_user_idx
    ON terminal_fcm_registration_token_dt (user_id)
    WHERE user_id is NOT NULL AND d_ts IS NULL;

CREATE TABLE phone_number_verification_dt (
    id                  bigserial PRIMARY KEY,
    country_code        integer NOT NULL,
    national_number     bigint NOT NULL,
    code                text NOT NULL,
    code_expiry         timestamp with time zone,
    attempts_remaining  smallint NOT NULL DEFAULT 3,

    c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    c_uid  bigint,
    c_tid  bigint,

    confirmation_time         timestamp with time zone,
    confirmation_user_id      bigint,
    confirmation_terminal_id  bigint
);

CREATE TABLE user_contact_phone_number_dt (
    user_id                  bigint NOT NULL,
    contact_country_code     integer NOT NULL,
    contact_national_number  bigint NOT NULL,

    c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    c_uid  bigint NOT NULL,
    c_tid  bigint NOT NULL,

    PRIMARY KEY (user_id, contact_country_code, contact_national_number)
);

-- user profile
CREATE TABLE user_display_name_dt (
    user_id       bigint NOT NULL,
    display_name  text NOT NULL,

    c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    c_uid  bigint NOT NULL,
    c_tid  bigint NOT NULL,

    d_ts   timestamp with time zone,
    d_uid  bigint,
    d_tid  bigint
);
CREATE UNIQUE INDEX user_display_name_dt_pidx
    ON user_display_name_dt (user_id)
    WHERE d_ts IS NULL;

CREATE TABLE user_profile_image_key_dt (
    user_id            bigint NOT NULL,
    profile_image_key  text NOT NULL,

    c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    c_uid  bigint NOT NULL,
    c_tid  bigint NOT NULL,

    d_ts   timestamp with time zone,
    d_uid  bigint,
    d_tid  bigint
);
CREATE UNIQUE INDEX user_profile_image_key_dt_pidx
    ON user_profile_image_key_dt (user_id)
    WHERE d_ts IS NULL;

-- user email
CREATE TABLE user_key_email_address_dt (
    user_id      bigint NOT NULL,
    local_part   text NOT NULL,
    domain_part  text NOT NULL,
    raw_input    text NOT NULL,

    c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    c_uid  bigint NOT NULL,
    c_tid  bigint NOT NULL,

    d_ts   timestamp with time zone,
    d_uid  bigint,
    d_tid  bigint,

    verification_id    bigint NOT NULL DEFAULT 0,
    verification_time  timestamp with time zone
);
-- Each user can only have one reference to the same email address
CREATE UNIQUE INDEX user_key_email_address_dt_pidx
    ON user_key_email_address_dt (user_id, local_part, domain_part)
    WHERE d_ts IS NULL;
-- Each user has only one non-deleted-verified email address
CREATE UNIQUE INDEX user_key_email_address_dt_user_id_uidx
    ON user_key_email_address_dt (user_id)
    WHERE d_ts IS NULL AND verification_time IS NOT NULL;
-- One instance for an non-deleted-verified email address
CREATE UNIQUE INDEX user_key_email_address_dt_local_part_domain_part_uidx
    ON user_key_email_address_dt (local_part, domain_part)
    WHERE d_ts IS NULL AND verification_time IS NOT NULL;

CREATE TABLE email_address_verification_dt (
    id                  bigserial PRIMARY KEY,
    local_part          text NOT NULL,
    domain_part         text NOT NULL,
    code                text NOT NULL,
    code_expiry         timestamp with time zone,
    attempts_remaining  smallint NOT NULL DEFAULT 3,

    c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    c_uid  bigint,
    c_tid  bigint,

    confirmation_time         timestamp with time zone,
    confirmation_user_id      bigint,
    confirmation_terminal_id  bigint
);

-- user password
--TODO: passwords for different purposes?
CREATE TABLE user_password_dt (
    user_id   bigint NOT NULL,
    password  text NOT NULL,

    c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    c_uid  bigint NOT NULL,
    c_tid  bigint NOT NULL,

    d_ts   timestamp with time zone,
    d_uid  bigint,
    d_tid  bigint
);
CREATE UNIQUE INDEX ON user_password_dt (user_id)
    WHERE d_ts IS NULL;

----
END;
