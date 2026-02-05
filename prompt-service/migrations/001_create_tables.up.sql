-- Create prompt templates table
CREATE TABLE IF NOT EXISTS prompt_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    version VARCHAR(50) DEFAULT 'v1',
    content TEXT NOT NULL,
    variables JSONB DEFAULT '[]',
    tags JSONB DEFAULT '[]',
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_templates_name ON prompt_templates(name);
CREATE INDEX IF NOT EXISTS idx_templates_status ON prompt_templates(status);

-- Create experiments table
CREATE TABLE IF NOT EXISTS prompt_experiments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    prompt_template_id UUID NOT NULL REFERENCES prompt_templates(id),
    model_id UUID NOT NULL,
    config JSONB DEFAULT '{}',
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW()
);
