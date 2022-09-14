CREATE TABLE users (
    id VARCHAR(36),
    email VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255),
    password VARCHAR(500) NOT NULL,
    profile_pic VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT uc_email UNIQUE (email)
);

CREATE TABLE languages (
    code VARCHAR(2),
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (code)
);

CREATE TABLE lessons (
    id VARCHAR(36),
    title VARCHAR(255) NOT NULL,
    text VARCHAR,
    language VARCHAR(2) NOT NULL,
    created_by VARCHAR(36) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_lessons_language FOREIGN KEY (language) REFERENCES languages(code),
    CONSTRAINT fk_lessons_created_by FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE courses (
    id VARCHAR(36),
    title VARCHAR(255) NOT NULL,
    lesson_ids JSON,
    created_by VARCHAR(36) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_courses_created_by FOREIGN KEY (created_by) REFERENCES users(id)
);

-- course_lessons records tell us which lessons are associated with which courses
-- a lesson can be associated with multiple courses allowing users to mix and match and compose custom courses 
-- as they so wish
CREATE TABLE course_lessons (
    course_id VARCHAR(36),
    lesson_id VARCHAR(36),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (course_id, lesson_id),
    CONSTRAINT fk_course_lessons_course_id FOREIGN KEY (course_id) REFERENCES courses(id),
    CONSTRAINT fk_course_lessons_lesson_id FOREIGN KEY (lesson_id) REFERENCES lessons(id)
);

-- user_courses records tell us which courses a user has enrolled in and which lesson the user has reached 
CREATE TABLE user_courses (
    user_id VARCHAR(36) NOT NULL,
    course_id VARCHAR(36) NOT NULL,
    active_lesson_number VARCHAR(36) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, course_id),
    CONSTRAINT fk_user_courses_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_user_courses_course FOREIGN KEY (course_id) REFERENCES courses(id)
);

