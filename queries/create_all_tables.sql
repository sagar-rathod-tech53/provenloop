-- Enable the uuid-ossp extension if it's not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS otps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL,
    otp VARCHAR(10) NOT NULL,
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL
);


CREATE TABLE IF NOT EXISTS content_post (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    post_content TEXT NOT NULL,
    media_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS post_likes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    post_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user_id, post_id), -- Prevent duplicate likes
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES content_post(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS post_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    post_id UUID REFERENCES content_post(id) ON DELETE CASCADE,
    comment TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);




CREATE TABLE IF NOT EXISTS job_post (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    job_title VARCHAR(255) NOT NULL,             -- Title of the job (e.g., Frontend Developer)
    company_name VARCHAR(255) NOT NULL,          -- Name of the company posting the job
    job_description TEXT,                       -- Description of the job responsibilities
    job_apply_url VARCHAR(255),                  -- URL where users can apply for the job
    location VARCHAR(255),                      -- Job location (e.g., Remote)
    last_date_to_apply DATE,                    -- Last date to apply for the job
    post_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Timestamp when the job listing was created
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS user_profile (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),           -- Auto-generated UUID
    user_id UUID NOT NULL,                                    -- Reference to the user
    profile_image TEXT,                                       --Profile Image
    full_name VARCHAR(255) NOT NULL,                          -- Full name of the user
    designation VARCHAR(255),                                 -- Designation (e.g., Backend Engineer)
    organization VARCHAR(255),                                -- Organization name (e.g., TechStars)
    professional_summary TEXT,                                -- Summary of professional experience
    location VARCHAR(255),                                    -- User's location (e.g., Mumbai, India)
    email VARCHAR(255) UNIQUE NOT NULL,                       -- Email of the user
    contact_number VARCHAR(20),                               -- Contact number of the user
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,           -- Timestamp when the profile was created
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,           -- Timestamp when the profile was last updated
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS video_profile (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    video_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_education (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),  -- Unique identifier
    user_id UUID NOT NULL,                           -- Foreign key to user
    degree VARCHAR(255) NOT NULL,
    institution_name VARCHAR(255) NOT NULL,
    field_of_study VARCHAR(255) NOT NULL,
    grade VARCHAR(10) NOT NULL,
    year VARCHAR(10),                                -- Store plain "2023"
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS user_experience (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    job_title VARCHAR(255) NOT NULL,        -- Title of the job (e.g., Front Desk Manager)
    company_name VARCHAR(255) NOT NULL,     -- Name of the company
    location VARCHAR(255),                  -- Job location (e.g., Mumbai, India)
    job_description TEXT,                   -- Description of job responsibilities
    achievements TEXT,                      -- Achievements during the role
    start_date DATE NOT NULL,               -- Start date of the position
    end_date DATE,                          -- End date of the position (nullable if the user is still in the role)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,   -- Timestamp when the record was created
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,   -- Timestamp when the record was last updated
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);



CREATE TABLE IF NOT EXISTS followers (
    follower_id UUID NOT NULL,   -- The user who is following
    followed_id UUID NOT NULL,   -- The user being followed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, followed_id),
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (followed_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS followings (
    followed_id UUID NOT NULL,   -- The user being followed
    follower_id UUID NOT NULL,   -- The user who is following
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, followed_id),
    FOREIGN KEY (followed_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    recipient_user_id UUID NOT NULL,
    sender_user_id UUID NOT NULL,
    type TEXT NOT NULL,
    entity_id UUID,
    entity_type TEXT,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (recipient_user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_user_id) REFERENCES users(id) ON DELETE CASCADE
);