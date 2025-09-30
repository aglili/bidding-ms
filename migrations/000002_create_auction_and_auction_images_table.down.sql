

-- Drop indexes
DROP INDEX IF EXISTS idx_auction_images_auction_id;
DROP INDEX IF EXISTS idx_auctions_status_end_time;
DROP INDEX IF EXISTS idx_auctions_end_time;
DROP INDEX IF EXISTS idx_auctions_status;
DROP INDEX IF EXISTS idx_auctions_seller_id;

-- Drop tables (child tables first due to foreign key constraints)
DROP TABLE IF EXISTS auction_images;
DROP TABLE IF EXISTS auctions;
