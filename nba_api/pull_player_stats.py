from nba_api.stats.endpoints import commonallplayers, playercareerstats
import pandas as pd
import time
from sqlalchemy import create_engine, text

# Connect to PostgreSQL
engine = create_engine('postgresql://nba:nba@172.21.192.1:5432/nbadb')

# Get all active players
players = commonallplayers.CommonAllPlayers(is_only_current_season=1)
players_df = players.get_data_frames()[0]

# Check which players are already in the database
try:
    with engine.connect() as conn:
        existing = pd.read_sql('SELECT DISTINCT "PLAYER_ID" FROM career_stats', conn)
        existing_ids = set(existing['PLAYER_ID'].tolist())
        print(f"Found {len(existing_ids)} players already in database, skipping them...")
except:
    existing_ids = set()
    print("No existing data found, fetching all players...")

all_career_stats = []

for index, player in players_df.iterrows():
    player_id = player['PERSON_ID']
    player_name = player['DISPLAY_FIRST_LAST']

    # Skip if already in database
    if player_id in existing_ids:
        print(f"Skipping {player_name} (already in database)")
        continue

    try:
        career = playercareerstats.PlayerCareerStats(player_id=player_id)
        df = career.season_totals_regular_season.get_data_frame()
        df['PLAYER_NAME'] = player_name
        all_career_stats.append(df)
        print(f"Fetched: {player_name}")
        time.sleep(0.6)
    except Exception as e:
        print(f"Failed for {player_name}: {e}")
        continue

if all_career_stats:
    final_df = pd.concat(all_career_stats, ignore_index=True)
    final_df.to_sql('career_stats', engine, if_exists='append', index=False)
    print("Done! Data saved to database.")
else:
    print("No new players to fetch!")