from sqlalchemy import create_engine, text

engine = create_engine('postgresql://nba:nba@172.21.192.1:5432/nbadb')

with engine.connect() as conn:
    result = conn.execute(text('SELECT 1'))
    print("Connection successful!", result.fetchone())