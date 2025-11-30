CREATE TABLE IF NOT EXISTS device_data
(
    id String,
    latitude Float64,
    longitude Float64,
    altitude Float64,
    battery Float64,
    timestamp DateTime64(0),
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, id);
