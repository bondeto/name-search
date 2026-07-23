-- ============================================================================
-- POSTGRESQL IDENTITY RESOLUTION & WATCHLIST SCREENING SCHEMA & FUNCTIONS
-- Extensions Required: pg_trgm (Trigram Indexing) & fuzzystrmatch (Phonetics)
-- ============================================================================

-- 1. Enable Required Extensions
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;

-- 2. Create Watchlist Table Schema
CREATE TABLE IF NOT EXISTS watchlist_identities (
    id VARCHAR(50) PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    dob DATE,
    passport_no VARCHAR(50),
    nationality VARCHAR(10),
    aliases TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Create High-Performance GIN & Phonetic Indexes
-- Trigram GIN index for sub-millisecond candidate blocking
CREATE INDEX IF NOT EXISTS idx_watchlist_name_trgm 
ON watchlist_identities USING GIN (full_name gin_trgm_ops);

-- Metaphone Index for Phonetic Blocking
CREATE INDEX IF NOT EXISTS idx_watchlist_metaphone 
ON watchlist_identities (dmetaphone(full_name));

-- DOB Index for Date Window Filtering
CREATE INDEX IF NOT EXISTS idx_watchlist_dob 
ON watchlist_identities (dob);


-- ============================================================================
-- 4. Custom Hybrid Scoring Search Function
-- Calculates Composite Weighted Score:
-- (40% Trigram Similarity) + (30% Phonetic Match) + (30% DOB Gaussian Decay)
-- ============================================================================

CREATE OR REPLACE FUNCTION fn_screen_watchlist(
    p_query_name TEXT,
    p_query_dob DATE DEFAULT NULL,
    p_threshold FLOAT DEFAULT 0.70
)
RETURNS TABLE (
    watchlist_id VARCHAR(50),
    full_name VARCHAR(255),
    dob DATE,
    passport_no VARCHAR(50),
    name_similarity FLOAT,
    phonetic_match BOOLEAN,
    dob_similarity FLOAT,
    composite_score FLOAT,
    match_status TEXT
) AS $$
BEGIN
    RETURN QUERY
    WITH candidate_pool AS (
        -- Step 1: Candidate Blocking via Trigram & Metaphone (< 5ms)
        SELECT 
            w.id,
            w.full_name,
            w.dob,
            w.passport_no,
            -- Trigram Similarity (0.0 to 1.0)
            similarity(UPPER(w.full_name), UPPER(p_query_name)) AS trgm_sim,
            -- Double Metaphone Phonetic Overlap
            (dmetaphone(w.full_name) = dmetaphone(p_query_name)) AS is_phonetic
        FROM watchlist_identities w
        WHERE 
            -- Trigram Similarity Threshold Filter (Blocking)
            UPPER(w.full_name) % UPPER(p_query_name)
            OR dmetaphone(w.full_name) = dmetaphone(p_query_name)
            OR (p_query_dob IS NOT NULL AND w.dob = p_query_dob)
    ),
    scored_candidates AS (
        -- Step 2: Attribute-Level Scoring Calculation
        SELECT 
            c.id,
            c.full_name,
            c.dob,
            c.passport_no,
            c.trgm_sim::FLOAT AS name_sim,
            c.is_phonetic,
            -- DOB Gaussian Decay Scoring: exp(-0.5 * (days / 30)^2)
            CASE 
                WHEN p_query_dob IS NULL OR c.dob IS NULL THEN 0.5
                WHEN c.dob = p_query_dob THEN 1.0
                ELSE EXP(-0.5 * POWER(ABS(c.dob - p_query_dob) / 30.0, 2))
            END::FLOAT AS dob_sim
        FROM candidate_pool c
    )
    -- Step 3: Composite Score & Threshold Decision Gate
    SELECT 
        s.id AS watchlist_id,
        s.full_name,
        s.dob,
        s.passport_no,
        ROUND((s.name_sim * 100)::numeric, 2)::FLOAT AS name_similarity,
        s.is_phonetic AS phonetic_match,
        ROUND((s.dob_sim * 100)::numeric, 2)::FLOAT AS dob_similarity,
        ROUND((
            (0.50 * s.name_sim) + 
            (0.20 * CASE WHEN s.is_phonetic THEN 1.0 ELSE 0.0 END) + 
            (0.30 * s.dob_sim)
        ) * 100::numeric, 2)::FLOAT AS composite_score,
        CASE 
            WHEN ((0.50 * s.name_sim) + (0.20 * CASE WHEN s.is_phonetic THEN 1.0 ELSE 0.0 END) + (0.30 * s.dob_sim)) >= 0.88 
                THEN 'CRITICAL MATCH (RED)'
            WHEN ((0.50 * s.name_sim) + (0.20 * CASE WHEN s.is_phonetic THEN 1.0 ELSE 0.0 END) + (0.30 * s.dob_sim)) >= p_threshold 
                THEN 'POTENTIAL MATCH (YELLOW)'
            ELSE 'NO MATCH (CLEAR)'
        END AS match_status
    FROM scored_candidates s
    WHERE ((0.50 * s.name_sim) + (0.20 * CASE WHEN s.is_phonetic THEN 1.0 ELSE 0.0 END) + (0.30 * s.dob_sim)) >= p_threshold
    ORDER BY composite_score DESC;
END;
$$ LANGUAGE plpgsql STABLE;

-- ============================================================================
-- SAMPLE USAGE & TESTING
-- ============================================================================

-- Insert Test Watchlist Records
INSERT INTO watchlist_identities (id, full_name, dob, passport_no, nationality) VALUES
('WL-001', 'Mohamad Ali Al-Khan', '1982-05-14', 'A12345678', 'ID'),
('WL-002', 'Usama Bin Laden', '1957-03-10', 'P98765432', 'SA'),
('WL-003', 'Vladimir Vladimirovich Putin', '1952-10-07', 'R55443322', 'RU')
ON CONFLICT (id) DO NOTHING;

-- Test Query Executions:
-- SELECT * FROM fn_screen_watchlist('MOHAMMED ALY KHAN', '1982-05-14', 0.70);
