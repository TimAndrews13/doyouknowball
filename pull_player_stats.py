from nba_api.stats.endpoints import commonallplayers, playercareerstats
import pandas as pd
import time
from sqlalchemy import create_engine

# Connect to PostgreSQL
engine = create_engine('postgresql://nba:nba@172.21.192.1:5432/nbadb')

# Get all active players
players = commonallplayers.CommonAllPlayers(is_only_current_season=1)
players_df = players.get_data_frames()[0]

all_career_stats = []

for index, player in players_df.iterrows():
    player_id = player['PERSON_ID']
    player_name = player['DISPLAY_FIRST_LAST']
    
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

# Combine into one dataframe
final_df = pd.concat(all_career_stats, ignore_index=True)

# Save to PostgreSQL
final_df.to_sql('career_stats', engine, if_exists='replace', index=False)
print("Done! Data saved to database.")