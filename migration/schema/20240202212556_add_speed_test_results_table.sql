-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS speed_test_results (
        id UUID NOT NULL,

        download_speed INT NOT NULL,
        max_download_speed INT DEFAULT NULL,
        min_download_speed INT DEFAULT NULL,
        total_download INT DEFAULT NULL,

        upload_speed INT NOT NULL,
        max_upload_speed INT DEFAULT NULL,
        min_upload_speed INT DEFAULT NULL,
        total_upload INT DEFAULT NULL,

        latency INT NOT NULL,
        loaded_latency INT DEFAULT NULL,
        unloaded_latency INT DEFAULT NULL,
        download_latency INT DEFAULT NULL,
        upload_latency INT DEFAULT NULL,

        connection_type VARCHAR(50) DEFAULT NULL,
        connection_device VARCHAR(50) DEFAULT NULL,
        isp VARCHAR(50) DEFAULT NULL,
        client_ip VARCHAR(50) DEFAULT NULL,
        client_id VARCHAR(50) DEFAULT NULL,
        city VARCHAR(50) DEFAULT NULL,
        server_name VARCHAR(50) DEFAULT NULL,

        longitude DECIMAL DEFAULT NULL,
        latitude DECIMAL DEFAULT NULL,

        location_access BOOLEAN DEFAULT false,

        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        test_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS speed_test_results
-- +goose StatementEnd
