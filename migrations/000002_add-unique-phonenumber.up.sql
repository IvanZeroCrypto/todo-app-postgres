ALTER TABLE todoapp.users
ADD CONSTRAINT users_unique_phone_number UNIQUE(phone_number)