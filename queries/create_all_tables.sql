CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_profile (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID NOT NULL UNIQUE,

    profile_image TEXT,
    full_name VARCHAR(255) NOT NULL,
    designation VARCHAR(255),
    organization VARCHAR(255),
    professional_summary TEXT,
    location VARCHAR(255),
    contact_number VARCHAR(20),

    college_name VARCHAR(255),
    university_name VARCHAR(255),
    degree VARCHAR(255),
    field_of_study VARCHAR(255),
    graduation_year INT,

    profile_video_url TEXT,
    video_created_at TIMESTAMP,
    video_updated_at TIMESTAMP,

    last_active_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    is_verified BOOLEAN DEFAULT FALSE,
    is_public BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS user_education (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,

    institution_name VARCHAR(255) NOT NULL,
    degree VARCHAR(255),                  -- 10th / 12th / Diploma / BTech / MTech etc
    field_of_study VARCHAR(255),
    start_year INT,
    end_year INT,

    grade VARCHAR(50),
    description TEXT,

    is_verified BOOLEAN DEFAULT FALSE,    -- admin/college verified

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_experience (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,

    company_name VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL,
    employment_type VARCHAR(100),   -- Full-time, Internship, Contract
    location VARCHAR(255),

    start_date DATE,
    end_date DATE,

    is_current BOOLEAN DEFAULT FALSE,

    description TEXT,

    is_verified BOOLEAN DEFAULT FALSE, -- HR verification

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS user_projects(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID NOT NULL,

    title VARCHAR(255) NOT NULL,
    description TEXT,

    tech_stack TEXT[],
    
    github_url TEXT,
    live_url TEXT,

    project_video_url TEXT,

    start_date DATE,
    end_date DATE,

    is_verified BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY(user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS user_skills(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID NOT NULL,

    skill_name VARCHAR(255) NOT NULL,

    proficiency_level VARCHAR(50),
    -- Beginner / Intermediate / Advanced / Expert

    endorsements_count INT DEFAULT 0,

    is_verified BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY(user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS user_certifications(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID NOT NULL,

    title VARCHAR(255) NOT NULL,
    issuing_organization VARCHAR(255) NOT NULL,

    credential_id VARCHAR(255),

    credential_url TEXT,

    issue_date DATE,
    expiry_date DATE,

    is_verified BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY(user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);