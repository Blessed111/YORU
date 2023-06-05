-- Drop table if EXISTS participants_team;
-- Update tournaments set current_round_number = 0, winner_team_id =0 where tournament_id = 2;
-- update teams set role = 'admin' where role = 'user'
CREATE TABLE IF NOT EXISTS teams (
    team_id SERIAL PRIMARY KEY,
    team_name VARCHAR(255) UNIQUE,
    player1 VARCHAR(255),
    player2 VARCHAR(255),
    player3 VARCHAR(255),
    player4 VARCHAR(255),
    player5 VARCHAR(255),
    role VARCHAR(5) DEFAULT 'user',
    password VARCHAR(255)
);


CREATE TABLE IF NOT EXISTS tournaments (
    tournament_id SERIAL PRIMARY KEY,
    tournament_name VARCHAR(255),
    description TEXT,
    start_date date,
    end_date date,
    teams_count INT,
    total_round_number INT,
    active BOOLEAN DEFAULT true,
    winner_team_id INT,
    current_teams_count INT DEFAULT 0,
    current_round_number INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS participants_team (
	team_id INT,
	tournament_id INT,
    CONSTRAINT unique_participant_team UNIQUE (team_id, tournament_id)
);


CREATE TABLE IF NOT EXISTS matches (
    match_id SERIAL PRIMARY KEY,
    tournament_id INT REFERENCES tournaments(tournament_id),
    round_number INT,
    participant1_id INT REFERENCES teams(team_id),
    participant2_id INT REFERENCES teams(team_id),
    winner_id INT,
    loser_id INT,
    CONSTRAINT unique_match_entry UNIQUE (tournament_id, round_number, participant1_id, participant2_id)
);


-- INSERT INTO teams (team_name, player1, player2, player3, player4, player5, role, password)
-- VALUES ('Team A', 'Player 1', 'Player 2', 'Player 3', 'Player 4', 'Player 5', 'user', 'qwe');
-- INSERT INTO teams (team_name, player1, player2, player3, player4, player5, role, password)
-- VALUES ('qwe', 'Player 1', 'Player 2', 'Player 3', 'Player 4', 'Player 5', 'user', 'qwe');
-- INSERT INTO teams (team_name, player1, player2, player3, player4, player5, role, password)
-- VALUES ('NAVI', 'Boom', 'ele', 'Simple', 'bit', 'perfect', 'user', 'password123');
