CREATE TABLE IF NOT EXISTS rooms (
    id SERIAL PRIMARY KEY,
    office_id INT NOT NULL REFERENCES offices(id),
    type VARCHAR(100) NOT NULL,
    price_hour NUMERIC(10,2) NOT NULL,
    qty_desks INT NOT NULL,
    qty_person INT NOT NULL,
    access_code INT, 
    has_air_conditioner BOOLEAN NOT NULL DEFAULT FALSE,
    has_prayer_room BOOLEAN NOT NULL DEFAULT FALSE,
    

    CONSTRAINT chk_access_code_required_for_private 
    CHECK (
        (type = 'private' AND access_code IS NOT NULL) OR 
        (type != 'private')
    )
);