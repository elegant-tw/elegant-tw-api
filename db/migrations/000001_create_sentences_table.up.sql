CREATE TABLE IF NOT EXISTS categories(
   category_id serial PRIMARY KEY,
   category_name VARCHAR (100) NOT NULL
);
CREATE TABLE IF NOT EXISTS sentences(
   id serial PRIMARY KEY,
   sentence VARCHAR (200) NOT NULL,
   category_id INT,
   cite VARCHAR (50) NOT NULL,
   author VARCHAR (50) NOT NULL,
   created_at TIMESTAMP DEFAULT NOW(),
   CONSTRAINT fk_category
      FOREIGN KEY(category_id) 
         REFERENCES categories(category_id)
);