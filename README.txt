1. __________________________________________________
my connection string is:
psql "postgres://postgres:@localhost:5433/chirpy"

2. __________________________________________________
to connect to the chirpy database:
sudo -u postgres psql chirpy

3. __________________________________________________
to migrate to the new schema:
cd into the sql/schema directory and run:
goose postgres "postgres://postgres:<pw>@localhost:5433/chirpy" up

4. __________________________________________________
to generate the new go methods like create user run (from project root):
sqlc generate 