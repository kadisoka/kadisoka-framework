--

-- convention:
-- prefix md_ indicates metadata of a record.
-- prefix md_c_ indicates metadata about a record's creation.
-- prefix md_d_ indicates metadata about a record's (soft) deletion.
-- we use suffix _dt (data table) for table names for consistency and
--   so that it can be generated in the future. using names like 'users'
--   and 'classes' are not generator-friendly.

\set ON_ERROR_STOP true

BEGIN;
------

CREATE TABLE user_dt (
    id_num  bigint PRIMARY KEY,

    md_c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    md_c_tid  bigint,
    md_c_uid  bigint,

    md_d_ts     timestamp with time zone,
    md_d_tid    bigint,
    md_d_uid    bigint,
    md_d_notes  jsonb,

    CHECK (id_num > 0)
);

CREATE TABLE terminal_dt (
    id_num          bigint PRIMARY KEY,
    application_id  integer NOT NULL,
    user_id         bigint NOT NULL, -- use zero if it's for a non-user

    md_c_ts              timestamp with time zone NOT NULL DEFAULT now(),
    md_c_tid             bigint,
    md_c_uid             bigint,
    md_c_origin_address  text,
    md_c_origin_env      text,

    md_d_ts   timestamp with time zone,
    md_d_tid  bigint,
    md_d_uid  bigint,

    secret           text NOT NULL,
    display_name     text,
    accept_language  text NOT NULL DEFAULT '', --TODO: list?

    verification_type  text NOT NULL,
    verification_id    bigint NOT NULL,
    verification_ts    timestamp with time zone,

    CHECK (application_id > 0 AND id_num > 0)
);
CREATE INDEX ON terminal_dt (user_id)
    WHERE md_d_ts IS NULL
    AND verification_ts IS NOT NULL;

CREATE TABLE session_dt (
    terminal_id  bigint NOT NULL,
    id_num       bigint NOT NULL,

    md_c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    md_c_tid  bigint,
    md_c_uid  bigint,
    md_c_sid  bigint, -- if a session was assumed, this field contains the parent session's ID.

    md_d_ts   timestamp with time zone,
    md_d_tid  bigint,
    md_d_uid  bigint,

    expiry  timestamp with time zone,

    PRIMARY KEY (terminal_id, id_num),
    CHECK (terminal_id > 0 AND id_num > 0)
);

CREATE TABLE terminal_fcm_registration_token_dt (
    terminal_id  bigint NOT NULL,
    user_id      bigint NOT NULL,
    token        text NOT NULL,

    md_c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    md_c_tid  bigint NOT NULL,
    md_c_uid  bigint NOT NULL,

    md_d_ts   timestamp with time zone,
    md_d_tid  bigint,
    md_d_uid  bigint
);
CREATE UNIQUE INDEX terminal_fcm_registration_token_dt_pidx
    ON terminal_fcm_registration_token_dt (terminal_id)
    WHERE md_d_ts IS NULL;
CREATE INDEX terminal_fcm_registration_token_dt_user_idx
    ON terminal_fcm_registration_token_dt (user_id)
    WHERE user_id is NOT NULL AND md_d_ts IS NULL;

CREATE TABLE user_key_phone_number_dt (
    user_id          bigint NOT NULL,
    country_code     integer NOT NULL,
    national_number  bigint NOT NULL, -- libphonenumber says uint64
    raw_input        text NOT NULL, 

    md_c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    md_c_tid  bigint NOT NULL,
    md_c_uid  bigint NOT NULL,

    md_d_ts   timestamp with time zone,
    md_d_tid  bigint,
    md_d_uid  bigint,

    verification_id  bigint NOT NULL DEFAULT 0,
    verification_ts  timestamp with time zone
);
-- Each user can only have one reference to the same phone number
CREATE UNIQUE INDEX user_key_phone_number_dt_pidx
    ON user_key_phone_number_dt (user_id, country_code, national_number)
    WHERE md_d_ts IS NULL;
-- Each user has only one non-deleted-verified phone number
CREATE UNIQUE INDEX user_key_phone_number_dt_user_id_uidx
    ON user_key_phone_number_dt (user_id)
    WHERE md_d_ts IS NULL AND verification_ts IS NOT NULL;
-- One instance for a non-deleted-verified phone number
CREATE UNIQUE INDEX user_key_phone_number_dt_country_code_national_number_uidx
    ON user_key_phone_number_dt (country_code, national_number)
    WHERE md_d_ts IS NULL AND verification_ts IS NOT NULL;

CREATE TABLE phone_number_verification_dt (
    id_num              bigserial PRIMARY KEY,
    country_code        integer NOT NULL,
    national_number     bigint NOT NULL,

    code                text NOT NULL,
    code_expiry         timestamp with time zone,

    md_c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    md_c_tid  bigint,
    md_c_uid  bigint,

    confirmation_attempts_remaining  smallint NOT NULL DEFAULT 3,
    confirmation_ts   timestamp with time zone,
    confirmation_tid  bigint,
    confirmation_uid  bigint
);

-- user profile
CREATE TABLE user_display_name_dt (
    user_id       bigint NOT NULL,
    display_name  text NOT NULL,

    md_c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    md_c_tid  bigint NOT NULL,
    md_c_uid  bigint NOT NULL,

    md_d_ts   timestamp with time zone,
    md_d_tid  bigint,
    md_d_uid  bigint
);
CREATE UNIQUE INDEX user_display_name_dt_pidx
    ON user_display_name_dt (user_id)
    WHERE md_d_ts IS NULL;

CREATE TABLE user_profile_image_key_dt (
    user_id            bigint NOT NULL,
    profile_image_key  text NOT NULL,

    md_c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    md_c_tid  bigint NOT NULL,
    md_c_uid  bigint NOT NULL,

    md_d_ts   timestamp with time zone,
    md_d_tid  bigint,
    md_d_uid  bigint
);
CREATE UNIQUE INDEX user_profile_image_key_dt_pidx
    ON user_profile_image_key_dt (user_id)
    WHERE md_d_ts IS NULL;

-- user email
CREATE TABLE user_key_email_address_dt (
    user_id      bigint NOT NULL,
    domain_part  text NOT NULL,
    local_part   text NOT NULL,
    raw_input    text NOT NULL,

    md_c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    md_c_tid  bigint NOT NULL,
    md_c_uid  bigint NOT NULL,

    md_d_ts   timestamp with time zone,
    md_d_tid  bigint,
    md_d_uid  bigint,

    verification_id  bigint NOT NULL DEFAULT 0,
    verification_ts  timestamp with time zone
);
-- Each user can only have one reference to the same email address
CREATE UNIQUE INDEX user_key_email_address_dt_pidx
    ON user_key_email_address_dt (user_id, domain_part, local_part)
    WHERE md_d_ts IS NULL;
-- Each user has only one non-deleted-verified email address
CREATE UNIQUE INDEX user_key_email_address_dt_user_id_uidx
    ON user_key_email_address_dt (user_id)
    WHERE md_d_ts IS NULL AND verification_ts IS NOT NULL;
-- One instance for a non-deleted-verified email address
CREATE UNIQUE INDEX user_key_email_address_dt_domain_part_local_part_uidx
    ON user_key_email_address_dt (domain_part, local_part)
    WHERE md_d_ts IS NULL AND verification_ts IS NOT NULL;

CREATE TABLE email_address_verification_dt (
    id_num              bigserial PRIMARY KEY,
    domain_part         text NOT NULL,
    local_part          text NOT NULL,

    code                text NOT NULL,
    code_expiry         timestamp with time zone,

    md_c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    md_c_tid  bigint,
    md_c_uid  bigint,

    confirmation_attempts_remaining  smallint NOT NULL DEFAULT 3,
    confirmation_ts   timestamp with time zone,
    confirmation_tid  bigint,
    confirmation_uid  bigint
);

-- user password
--TODO: passwords for different purposes?
CREATE TABLE user_password_dt (
    user_id   bigint NOT NULL,
    password  text NOT NULL,

    md_c_ts   timestamp with time zone NOT NULL DEFAULT now(),
    md_c_tid  bigint NOT NULL,
    md_c_uid  bigint NOT NULL,

    md_d_ts   timestamp with time zone,
    md_d_tid  bigint,
    md_d_uid  bigint
);
CREATE UNIQUE INDEX ON user_password_dt (user_id)
    WHERE md_d_ts IS NULL;

----
END;
