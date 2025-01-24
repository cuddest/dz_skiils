CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);


CREATE TABLE sub_cats (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category_id INTEGER REFERENCES categories(id) ON DELETE CASCADE
);


CREATE TABLE teachers (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    picture VARCHAR(255),
    skills TEXT,
    degrees VARCHAR(255),
    experience TEXT
);


CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    picture VARCHAR(255)
);


CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    pricing VARCHAR(255),
    duration VARCHAR(255),
    image VARCHAR(255),
    language VARCHAR(50),
    level VARCHAR(50),
    teacher_id INTEGER REFERENCES teachers(id),
    category_id INTEGER REFERENCES categories(id)
);

CREATE TABLE student_courses (
    student_id INTEGER REFERENCES students(id),
    course_id INTEGER REFERENCES courses(id),
    grade VARCHAR(10),
    enrollment TIMESTAMP,
    certificate VARCHAR(255),
    issued BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (student_id, course_id)
);


CREATE TABLE cratings (
    course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE,
    student_id INTEGER REFERENCES students(id) ON DELETE CASCADE,
    rating DECIMAL(3,2),
    PRIMARY KEY (course_id, student_id)
);


CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    link VARCHAR(255),
    description TEXT,
    course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE
);

CREATE TABLE videos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    link VARCHAR(255) NOT NULL,
    course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE
);


CREATE TABLE feedbacks (
    id SERIAL PRIMARY KEY,
    description TEXT,
    review INTEGER,
    student_id INTEGER REFERENCES students(id) ON DELETE CASCADE
);


CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE,
    student_id INTEGER REFERENCES students(id) ON DELETE CASCADE,
    question TEXT NOT NULL
);

CREATE TABLE answers (
    id SERIAL PRIMARY KEY,
    answer TEXT NOT NULL,
    question_id INTEGER REFERENCES questions(id) ON DELETE CASCADE
);

CREATE TABLE exams (
    id SERIAL PRIMARY KEY,
    description TEXT,
    course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE UNIQUE
);


CREATE TABLE course_quizzes (
    id SERIAL PRIMARY KEY,
    question TEXT NOT NULL,
    option1 VARCHAR(255) NOT NULL,
    option2 VARCHAR(255) NOT NULL,
    option3 VARCHAR(255) NOT NULL,
    option4 VARCHAR(255) NOT NULL,
    answer VARCHAR(255) NOT NULL,
    course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE
);


CREATE TABLE exam_quizzes (
    id SERIAL PRIMARY KEY,
    question TEXT NOT NULL,
    option1 VARCHAR(255) NOT NULL,
    option2 VARCHAR(255) NOT NULL,
    option3 VARCHAR(255) NOT NULL,
    option4 VARCHAR(255) NOT NULL,
    answer INTEGER NOT NULL,
    exam_id INTEGER REFERENCES exams(id) ON DELETE CASCADE