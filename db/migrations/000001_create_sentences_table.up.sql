CREATE TABLE IF NOT EXISTS categories(
   category_id serial PRIMARY KEY,
   category_name VARCHAR (100) NOT NULL
);

INSERT INTO "categories" ("category_id", "category_name")
VALUES ('0', '未分類');

CREATE TABLE IF NOT EXISTS sentences(
   id serial PRIMARY KEY,
   sentence VARCHAR (200) NOT NULL,
   category_id INT DEFAULT 0,
   cite VARCHAR (50),
   author VARCHAR (50),
   created_at TIMESTAMP DEFAULT NOW(),
   CONSTRAINT fk_category
      FOREIGN KEY(category_id) 
         REFERENCES categories(category_id)
);