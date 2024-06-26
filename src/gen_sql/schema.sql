-- BEGIN models
CREATE SCHEMA cew;
CREATE TABLE cew.University (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::VARCHAR,
    title VARCHAR NOT NULL,
    about VARCHAR NOT NULL
);
CREATE TABLE cew.BaseUser(
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::VARCHAR,
    name_first VARCHAR NOT NULL,
    name_last VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    password_hash VARCHAR NOT NULL
);
CREATE TABLE cew.Student (
    id VARCHAR PRIMARY KEY REFERENCES cew.BaseUser(id) ON UPDATE CASCADE ON DELETE CASCADE,
    university_id VARCHAR NOT NULL REFERENCES cew.University(id) ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE TABLE cew.SuperAdmin (
    id VARCHAR PRIMARY KEY REFERENCES cew.BaseUser(id) ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE TABLE cew.Rso (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::VARCHAR,
    title VARCHAR NOT NULL,
    about VARCHAR NOT NULL,
    university_id VARCHAR NOT NULL REFERENCES cew.University(id) ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE TABLE cew.RsoMember (
    id VARCHAR NOT NULL REFERENCES cew.Student(id) ON UPDATE CASCADE ON DELETE CASCADE,
    rso_id VARCHAR NOT NULL REFERENCES cew.Rso(id) ON UPDATE CASCADE ON DELETE CASCADE,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (rso_id, id)
);
CREATE TABLE cew.BaseEvent (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::VARCHAR,
    title VARCHAR NOT NULL,
    about VARCHAR NOT NULL,
    university_id VARCHAR NOT NULL REFERENCES cew.University(id) ON UPDATE CASCADE ON DELETE CASCADE,
    post_time TIMESTAMP NOT NULL,
    start_time TIMESTAMP NOT NULL,
    contact_phone VARCHAR NOT NULL,
    contact_email VARCHAR NOT NULL,
    event_location VARCHAR NOT NULL
);
CREATE TABLE cew.PublicEvent (
    id VARCHAR PRIMARY KEY REFERENCES cew.BaseEvent(id),
    approved BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE TABLE cew.PrivateEvent (
    id VARCHAR PRIMARY KEY REFERENCES cew.BaseEvent(id)
);
CREATE TABLE cew.RsoEvent (
    id VARCHAR PRIMARY KEY REFERENCES cew.BaseEvent(id),
    rso_id VARCHAR NOT NULL REFERENCES cew.Rso(id)
);
CREATE TABLE cew.Tag (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::VARCHAR,
    title VARCHAR NOT NULL
);
CREATE TABLE cew.TaggedEvent (
    tag_id VARCHAR NOT NULL REFERENCES cew.Tag(id) ON UPDATE CASCADE ON DELETE CASCADE,
    base_event_id VARCHAR NOT NULL REFERENCES cew.BaseEvent(id) ON UPDATE CASCADE ON DELETE CASCADE,
    PRIMARY KEY (tag_id, base_event_id)
);
CREATE TABLE cew.TaggedRso (
    tag_id VARCHAR NOT NULL REFERENCES cew.Tag(id) ON UPDATE CASCADE ON DELETE CASCADE,
    rso_id VARCHAR NOT NULL REFERENCES cew.Rso(id) ON UPDATE CASCADE ON DELETE CASCADE,
    PRIMARY KEY (tag_id, rso_id)
);
CREATE TABLE cew.Comment (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::VARCHAR,
    body VARCHAR NOT NULL,
    student_id VARCHAR REFERENCES cew.Student(id) ON UPDATE CASCADE ON DELETE
    SET NULL,
        base_event_id VARCHAR NOT NULL REFERENCES cew.BaseEvent(id) ON UPDATE CASCADE ON DELETE CASCADE,
        post_timestamp TIMESTAMP NOT NULL
);
CREATE TABLE cew.Rating (
    id VARCHAR PRIMARY KEY DEFAULT gen_random_uuid()::VARCHAR,
    stars INTEGER NOT NULL CHECK (
        stars > 0
        AND stars <= 5
    ),
    student_id VARCHAR REFERENCES cew.Student(id) ON UPDATE CASCADE ON DELETE
    SET NULL,
        base_event_id VARCHAR NOT NULL REFERENCES cew.BaseEvent(id) ON UPDATE CASCADE ON DELETE CASCADE,
        post_timestamp TIMESTAMP NOT NULL
);
-- END models