-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       full_name VARCHAR(255) NOT NULL,
                       free_generations_left INT DEFAULT 2,
                       is_premium BOOLEAN DEFAULT FALSE,
                       created_at TIMESTAMP DEFAULT NOW(),
                       updated_at TIMESTAMP DEFAULT NOW()
);

-- Profiles table
CREATE TABLE profiles (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          user_id UUID REFERENCES users(id) ON DELETE CASCADE UNIQUE,
                          phone VARCHAR(50),
                          location VARCHAR(255),
                          linkedin_url VARCHAR(500),
                          github_url VARCHAR(500),
                          website_url VARCHAR(500),
                          summary TEXT,
                          created_at TIMESTAMP DEFAULT NOW(),
                          updated_at TIMESTAMP DEFAULT NOW()
);

-- Experiences table
CREATE TABLE experiences (
                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                             profile_id UUID REFERENCES profiles(id) ON DELETE CASCADE,
                             company VARCHAR(255) NOT NULL,
                             position VARCHAR(255) NOT NULL,
                             start_date DATE NOT NULL,
                             end_date DATE,
                             is_current BOOLEAN DEFAULT FALSE,
                             description TEXT,
                             achievements TEXT[],
                             created_at TIMESTAMP DEFAULT NOW()
);

-- Education table
CREATE TABLE education (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           profile_id UUID REFERENCES profiles(id) ON DELETE CASCADE,
                           institution VARCHAR(255) NOT NULL,
                           degree VARCHAR(255) NOT NULL,
                           field_of_study VARCHAR(255),
                           start_date DATE NOT NULL,
                           end_date DATE,
                           gpa DECIMAL(3,2),
                           created_at TIMESTAMP DEFAULT NOW()
);

-- Skills table
CREATE TABLE skills (
                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                        profile_id UUID REFERENCES profiles(id) ON DELETE CASCADE,
                        name VARCHAR(100) NOT NULL,
                        category VARCHAR(50),
                        proficiency_level VARCHAR(50),
                        created_at TIMESTAMP DEFAULT NOW()
);

-- Documents table
CREATE TABLE documents (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                           type VARCHAR(50) NOT NULL,
                           title VARCHAR(255) NOT NULL,
                           content JSONB NOT NULL,
                           template_id VARCHAR(50),
                           job_title VARCHAR(255),
                           company_name VARCHAR(255),
                           job_description TEXT,
                           status VARCHAR(50) DEFAULT 'draft',
                           created_at TIMESTAMP DEFAULT NOW(),
                           updated_at TIMESTAMP DEFAULT NOW()
);

-- Generation history table
CREATE TABLE generation_history (
                                    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                                    document_id UUID REFERENCES documents(id) ON DELETE CASCADE,
                                    prompt_tokens INT,
                                    completion_tokens INT,
                                    total_cost DECIMAL(10,4),
                                    generation_time_ms INT,
                                    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_profiles_user_id ON profiles(user_id);
CREATE INDEX idx_experiences_profile_id ON experiences(profile_id);
CREATE INDEX idx_education_profile_id ON education(profile_id);
CREATE INDEX idx_skills_profile_id ON skills(profile_id);
CREATE INDEX idx_documents_user_id ON documents(user_id);
CREATE INDEX idx_documents_created_at ON documents(created_at DESC);
CREATE INDEX idx_generation_history_user_id ON generation_history(user_id);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_profiles_updated_at BEFORE UPDATE ON profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_documents_updated_at BEFORE UPDATE ON documents
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();