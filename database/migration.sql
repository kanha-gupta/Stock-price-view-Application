CREATE TABLE stocks (
                        id INT AUTO_INCREMENT PRIMARY KEY,
                        code VARCHAR(10),
                        name VARCHAR(255),
                        open DECIMAL(10, 2),
                        high DECIMAL(10, 2),
                        low DECIMAL(10, 2),
                        close DECIMAL(10, 2)
);

CREATE TABLE favourites (
                            id INT AUTO_INCREMENT PRIMARY KEY,
                            code VARCHAR(255) NOT NULL
);
