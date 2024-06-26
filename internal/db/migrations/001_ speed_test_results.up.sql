CREATE TABLE 
     IF NOT EXISTS users (
        id UUID NOT NULL PRIMARY KEY,
        
        username VARCHAR(50) NOT NULL UNIQUE,
        email VARCHAR(100) NOT NULL UNIQUE,

        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP DEFAULT NULL
    );

CREATE TABLE 
    IF NOT EXISTS test_servers (
        id UUID NOT NULL PRIMARY KEY,
        
        name VARCHAR(100),
        identifier VARCHAR(100) UNIQUE NOT NULL,
        city VARCHAR(100),
        country VARCHAR(100) NOT NULL,
        url VARCHAR(255),

        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP DEFAULT NULL
    );

CREATE TABLE 
    IF NOT EXISTS devices (
       id UUID NOT NULL PRIMARY KEY,

        user_id UUID DEFAULT NULL,
        identifier VARCHAR(100) UNIQUE NOT NULL,
        os VARCHAR(50),
        device_type VARCHAR(50),
        manufacturer VARCHAR(50),
        model VARCHAR(50),
        screen_resolution VARCHAR(50),

        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP DEFAULT NULL,

        FOREIGN KEY (user_id) REFERENCES users(id)
    );

CREATE TABLE IF NOT EXISTS speed_test_results (
    id UUID NOT NULL PRIMARY KEY,

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

    device_id UUID,
    isp VARCHAR(50) DEFAULT NULL,
    isp_code VARCHAR(15) DEFAULT NULL,
    connection_type VARCHAR(50) DEFAULT NULL,
    connection_device VARCHAR(50) DEFAULT NULL,
    test_platform VARCHAR(50) DEFAULT NULL,
    server_id UUID,

    city VARCHAR(50) DEFAULT NULL,
    state VARCHAR(50) DEFAULT NULL,
    country_code VARCHAR(5) DEFAULT NULL,
    country_name VARCHAR(50) DEFAULT NULL,
    longitude DECIMAL DEFAULT NULL,
    latitude DECIMAL DEFAULT NULL,
    location_access BOOLEAN DEFAULT false,

    test_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY (device_id) REFERENCES devices(id),
    FOREIGN KEY (server_id) REFERENCES test_servers(id)
);
