CREATE TABLE contacts (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),       -- Unique ID for each contact
  user_id UUID NOT NULL,                                -- Foreign key to users table
  phone VARCHAR(20),                                    -- Contact's phone number
  street VARCHAR(100),                                  -- Street address
  city VARCHAR(50),                                     -- City
  state VARCHAR(50),                                    -- State
  zip_code VARCHAR(20),                                 -- Zip code
  country VARCHAR(50),                                  -- Country
  created_at TIMESTAMP DEFAULT NOW(),                   -- Created timestamp
  updated_at TIMESTAMP DEFAULT NOW(),                   -- Updated timestamp
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE -- Foreign key constraint
);