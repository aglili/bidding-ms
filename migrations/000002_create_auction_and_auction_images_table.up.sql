CREATE TABLE auctions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    starting_price NUMERIC(12,2) NOT NULL CHECK (starting_price > 0),
    current_price NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (current_price >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'closed', 'cancelled')),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_auction_times CHECK (end_time > start_time)
);

CREATE TABLE auction_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auction_id UUID NOT NULL REFERENCES auctions(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL CHECK (image_url != ''),
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_auctions_seller_id ON auctions(seller_id);
CREATE INDEX idx_auctions_status ON auctions(status);
CREATE INDEX idx_auctions_end_time ON auctions(end_time);
CREATE INDEX idx_auctions_status_end_time ON auctions(status, end_time);
CREATE INDEX idx_auction_images_auction_id ON auction_images(auction_id);
