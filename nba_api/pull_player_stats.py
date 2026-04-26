from nba_api.stats.endpoints import commonallplayers, playercareerstats
import pandas as pd
import time
from sqlalchemy import create_engine, text
from dotenv import load_dotenv
import os

# Connect to PostgreSQL
load_dotenv()
engine = create_engine(os.getenv("DB_URL_PYTHON"))

# Get all active players
players = commonallplayers.CommonAllPlayers(is_only_current_season=1)
players_df = players.get_data_frames()[0]

for index, player in players_df.iterrows():
    player_id = player['PERSON_ID']
    player_name = player['DISPLAY_FIRST_LAST']

    try:
        career = playercareerstats.PlayerCareerStats(player_id=player_id)
        df = career.season_totals_regular_season.get_data_frame()
        df['PLAYER_NAME'] = player_name

        # Delete existing rows for this player then re-insert
        with engine.connect() as conn:
            conn.execute(text('DELETE FROM career_stats WHERE "PLAYER_ID" = :pid'), {"pid": player_id})
            conn.commit()

        df.to_sql('career_stats', engine, if_exists='append', index=False)
        print(f"Fetched: {player_name}")
        time.sleep(0.6)

    except Exception as e:
        print(f"Failed for {player_name}: {e}")
        continue

print("Done! All player stats updated.")