-- Migration để thêm trường image_urls vào bảng review
-- Chạy file này nếu database hiện tại chưa có trường image_urls

-- Kiểm tra xem trường image_urls đã tồn tại chưa
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'review' 
        AND column_name = 'image_urls'
    ) THEN
        -- Thêm trường image_urls nếu chưa có
        ALTER TABLE review ADD COLUMN image_urls TEXT;
        
        -- Cập nhật các review cũ để có image_urls = NULL
        UPDATE review SET image_urls = NULL WHERE image_urls IS NULL;
        
        RAISE NOTICE 'Đã thêm trường image_urls vào bảng review';
    ELSE
        RAISE NOTICE 'Trường image_urls đã tồn tại trong bảng review';
    END IF;
END $$;

-- Kiểm tra kết quả
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'review' 
ORDER BY ordinal_position;
